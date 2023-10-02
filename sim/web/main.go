package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	uuid "github.com/google/uuid"
	"github.com/pkg/browser"
	dist "github.com/wowsims/wotlk/binary_dist"
	"github.com/wowsims/wotlk/sim"
	"github.com/wowsims/wotlk/sim/core"
	proto "github.com/wowsims/wotlk/sim/core/proto"

	googleProto "google.golang.org/protobuf/proto"
)

func init() {
	sim.RegisterAll()
}

var (
	Version  string
	outdated int
)

func main() {
	if Version == "" {
		Version = "development"
	}
	var useFS = flag.Bool("usefs", false, "Use local file system for client files. Set to true during development.")
	var wasm = flag.Bool("wasm", false, "Use wasm for sim instead of web server apis. Can only be used with usefs=true")
	var simName = flag.String("sim", "", "Name of simulator to launch (ex: balance_druid, elemental_shaman, etc)")
	var host = flag.String("host", "localhost:3333", "URL to host the interface on.")
	var launch = flag.Bool("launch", true, "auto launch browser")
	var skipVersionCheck = flag.Bool("nvc", false, "set true to skip version check")

	flag.Parse()

	fmt.Printf("Version: %s\n", Version)
	if !*skipVersionCheck && Version != "development" {
		go func() {
			resp, err := http.Get("https://api.github.com/repos/wowsims/wotlk/releases/latest")
			if err != nil {
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}

			result := struct {
				Tag  string `json:"tag_name"`
				URL  string `json:"html_url"`
				Name string `json:"name"`
			}{}
			if err := json.Unmarshal(body, &result); err != nil {
				return
			}

			if result.Tag != Version {
				outdated = 2
				fmt.Printf("New version of simulator available: %s\n\tDownload at: %s\n", result.Name, result.URL)
			} else {
				outdated = 1
			}
		}()
	}

	s := &server{
		progMut:         sync.RWMutex{},
		asyncProgresses: map[string]*asyncProgress{},
	}
	s.runServer(*useFS, *host, *launch, *simName, *wasm, bufio.NewReader(os.Stdin))
}

// Handlers to decode and handle each proto function
var handlers = map[string]apiHandler{
	"/raidSim": {msg: func() googleProto.Message { return &proto.RaidSimRequest{} }, handle: func(msg googleProto.Message) googleProto.Message {
		return core.RunRaidSim(msg.(*proto.RaidSimRequest))
	}},
	"/statWeights": {msg: func() googleProto.Message { return &proto.StatWeightsRequest{} }, handle: func(msg googleProto.Message) googleProto.Message {
		return core.StatWeights(msg.(*proto.StatWeightsRequest))
	}},
	"/computeStats": {msg: func() googleProto.Message { return &proto.ComputeStatsRequest{} }, handle: func(msg googleProto.Message) googleProto.Message {
		return core.ComputeStats(msg.(*proto.ComputeStatsRequest))
	}},
}

var asyncAPIHandlers = map[string]asyncAPIHandler{
	"/raidSimAsync": {msg: func() googleProto.Message { return &proto.RaidSimRequest{} }, handle: func(msg googleProto.Message, reporter chan *proto.ProgressMetrics) {
		core.RunRaidSimAsync(msg.(*proto.RaidSimRequest), reporter)
	}},
	"/statWeightsAsync": {msg: func() googleProto.Message { return &proto.StatWeightsRequest{} }, handle: func(msg googleProto.Message, reporter chan *proto.ProgressMetrics) {
		core.StatWeightsAsync(msg.(*proto.StatWeightsRequest), reporter)
	}},
	"/bulkSimAsync": {msg: func() googleProto.Message { return &proto.BulkSimRequest{} }, handle: func(msg googleProto.Message, reporter chan *proto.ProgressMetrics) {
		// TODO: we can use context's to cancel stuff.
		// We should have all the async APIs take in context and let it be cancelled via its async ID.
		core.RunBulkSimAsync(context.Background(), msg.(*proto.BulkSimRequest), reporter)
	}},
}

type server struct {
	progMut         sync.RWMutex
	asyncProgresses map[string]*asyncProgress
}

type apiHandler struct {
	msg    func() googleProto.Message
	handle func(googleProto.Message) googleProto.Message
}
type asyncAPIHandler struct {
	msg    func() googleProto.Message
	handle func(googleProto.Message, chan *proto.ProgressMetrics)
}

type asyncProgress struct {
	id             string
	latestProgress atomic.Value
}

func (s *server) addNewSim() *asyncProgress {
	newID := uuid.NewString()
	simProgress := &asyncProgress{
		id: newID,
	}
	simProgress.latestProgress.Store(&proto.ProgressMetrics{})

	s.progMut.Lock()
	s.asyncProgresses[newID] = simProgress
	s.progMut.Unlock()

	return simProgress
}

func (s *server) handleAsyncAPI(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	endpoint := r.URL.Path
	handler, ok := asyncAPIHandlers[endpoint]
	if !ok {
		log.Printf("Invalid Endpoint: %s", endpoint)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	msg := handler.msg()
	if err := googleProto.Unmarshal(body, msg); err != nil {
		log.Printf("Failed to parse request: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// reporter channel is handed into the core simulation.
	//  as the simulation advances it will push changes to the channel
	//  these changes will be consumed by the goroutine below so the asyncProgress endpoint can fetch the results.
	reporter := make(chan *proto.ProgressMetrics, 100)
	handler.handle(msg, reporter)

	// Generate a new async simulation
	simProgress := s.addNewSim()

	// Now launch a background process that pulls progress reports off the reporter channel
	// and pushes it into the async progress cache.
	go func() {
		for {
			select {
			case <-time.After(time.Minute * 10):
				// if we get no progress after 10 minutes, delete the pending sim and exit.
				s.progMut.Lock()
				delete(s.asyncProgresses, simProgress.id)
				s.progMut.Unlock()
				return
			case progMetric := <-reporter:
				if progMetric == nil {
					return
				}
				simProgress.latestProgress.Store(progMetric)
				if progMetric.FinalRaidResult != nil || progMetric.FinalWeightResult != nil || progMetric.FinalBulkResult != nil {
					return
				}
			}
		}
	}()

	protoResult := &proto.AsyncAPIResult{
		ProgressId: simProgress.id,
	}

	outbytes, err := googleProto.Marshal(protoResult)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/x-protobuf")
	w.Write(outbytes)
}

func (s *server) setupAsyncServer() {
	// All async handlers here will call the addNewSim, generating a new UUID and cached progress state.
	for route := range asyncAPIHandlers {
		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			s.handleAsyncAPI(w, r)
		})
	}

	// asyncProgress will fetch the current progress of a simulation by its UUID.
	http.HandleFunc("/asyncProgress", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		msg := &proto.AsyncAPIResult{}
		if err := googleProto.Unmarshal(body, msg); err != nil {
			log.Printf("Failed to parse request: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Read lock the map of all progress statuses, fetching current one.
		s.progMut.RLock()
		progress, ok := s.asyncProgresses[msg.ProgressId]
		s.progMut.RUnlock()
		if !ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		latest := progress.latestProgress.Load().(*proto.ProgressMetrics)
		outbytes, err := googleProto.Marshal(latest)
		if err != nil {
			log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// If this was the last result, delete the cache for this simulation.
		if latest.FinalRaidResult != nil || latest.FinalWeightResult != nil || latest.FinalBulkResult != nil {
			s.progMut.Lock()
			delete(s.asyncProgresses, msg.ProgressId)
			s.progMut.Unlock()
		}
		w.Header().Add("Content-Type", "application/x-protobuf")
		w.Write(outbytes)
	})
}

func (s *server) runServer(useFS bool, host string, launchBrowser bool, simName string, wasm bool, inputReader *bufio.Reader) {
	s.setupAsyncServer()

	var fs http.Handler
	if useFS {
		log.Printf("Using local file system for development.")
		fs = http.FileServer(http.Dir("./dist"))
	} else {
		log.Printf("Embedded file server running.")
		fs = http.FileServer(http.FS(dist.FS))
	}

	for route := range handlers {
		http.HandleFunc(route, handleAPI)
	}

	http.HandleFunc("/version", func(resp http.ResponseWriter, req *http.Request) {
		msg := fmt.Sprintf(`{"version": "%s", "outdated": %d}`, Version, outdated)
		resp.Write([]byte(msg))
	})
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			http.Redirect(resp, req, "/wotlk/", http.StatusPermanentRedirect)
			return
		}
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("Content-Type", "application/wasm")
		}
		if strings.HasSuffix(req.URL.Path, ".js") {
			resp.Header().Set("Content-Type", "application/javascript")
		}
		if !useFS || (useFS && !wasm) {
			if strings.HasSuffix(req.URL.Path, "sim_worker.js") {
				req.URL.Path = strings.Replace(req.URL.Path, "sim_worker.js", "net_worker.js", 1)
			}
		}
		fs.ServeHTTP(resp, req)
	})

	if launchBrowser {
		if strings.HasPrefix(host, ":") {
			host = "localhost" + host
		}
		url := fmt.Sprintf("http://%s/wotlk/%s", host, simName)
		log.Printf("Launching interface on %s", url)
		go func() {
			err := browser.OpenURL(url)
			if err != nil {
				fmt.Printf("Error launching browser: %#v\n", err.Error())
				fmt.Printf("You will need to manually open your web browser to %s\n", url)
			}
		}()
	}

	go func() {
		// Launch server!
		if err := http.ListenAndServe(host, nil); err != nil {
			log.Printf("Failed to shutdown server: %s", err)
			os.Exit(1)
		}
		log.Printf("Server shutdown successfully.")
		os.Exit(0)
	}()

	// used to read a CTRL+C
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT)

	go func() {
		<-c
		log.Printf("Shutting down")
		os.Exit(0)
	}()
	fmt.Printf("Enter Command... '?' for list\n")
	for {
		fmt.Printf("> ")
		text, err := inputReader.ReadString('\n')
		if err != nil {
			// block forever
			<-c
			os.Exit(-1)
		}
		if len(text) == 0 {
			continue
		}
		command := strings.TrimSpace(text)
		switch command {
		case "profile":
			filename := fmt.Sprintf("profile_%d.cpu", time.Now().Unix())
			fmt.Printf("Running profiling for 15 seconds, output to %s\n", filename)
			f, err := os.Create(filename)
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
			go func() {
				time.Sleep(time.Second * 15)
				pprof.StopCPUProfile()
				f.Close()
				fmt.Printf("Profiling complete.\n> ")
			}()
		case "sims":
			s.progMut.RLock()
			fmt.Printf("Total Sims Running: %d\n", len(s.asyncProgresses))
			for _, v := range s.asyncProgresses {
				latest := (v.latestProgress.Load()).(*proto.ProgressMetrics)
				fmt.Printf("Process: %s (%d sims)\n\t  Progress: %d/%d\n", v.id, latest.TotalSims, latest.CompletedIterations, latest.TotalIterations)
			}
			s.progMut.RUnlock()
		case "quit":
			os.Exit(1)
		case "?":
			fmt.Printf("Commands:\n\tsims - Lists all active async sims running currently.\n\tprofile - start a CPU profile for debugging performance\n\tquit - exits\n\n")
		case "":
			// nothing.
		default:
			fmt.Printf("Unknown command: '%s'", command)
		}
	}
}

// handleAPI is generic handler for any api function using protos.
func handleAPI(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path

	body, err := io.ReadAll(r.Body)
	if err != nil {

		return
	}
	handler, ok := handlers[endpoint]
	if !ok {
		log.Printf("Invalid Endpoint: %s", endpoint)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	msg := handler.msg()
	if err := googleProto.Unmarshal(body, msg); err != nil {
		log.Printf("Failed to parse request: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := handler.handle(msg)

	outbytes, err := googleProto.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/x-protobuf")
	w.Write(outbytes)
}
