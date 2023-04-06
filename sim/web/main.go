package main

import (
	"bufio"
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

	"github.com/pkg/browser"
	uuid "github.com/satori/go.uuid"
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
	var host = flag.String("host", ":3333", "URL to host the interface on.")
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

	setupAsyncServer()
	runServer(*useFS, *host, *launch, *simName, *wasm, bufio.NewReader(os.Stdin))
}

type simProgReportCreator func() (string, progReport)
type progReport func(progMetric *proto.ProgressMetrics)
type asyncAPIHandler struct {
	msg    func() googleProto.Message
	handle func(googleProto.Message, chan *proto.ProgressMetrics)
}

var asyncAPIHandlers = map[string]asyncAPIHandler{
	"/raidSimAsync": {msg: func() googleProto.Message { return &proto.RaidSimRequest{} }, handle: func(msg googleProto.Message, reporter chan *proto.ProgressMetrics) {
		core.RunRaidSimAsync(msg.(*proto.RaidSimRequest), reporter)
	}},
	"/statWeightsAsync": {msg: func() googleProto.Message { return &proto.StatWeightsRequest{} }, handle: func(msg googleProto.Message, reporter chan *proto.ProgressMetrics) {
		core.StatWeightsAsync(msg.(*proto.StatWeightsRequest), reporter)
	}},
	"/bulkSimAsync": {msg: func() googleProto.Message { return &proto.BulkSimRequest{} }, handle: func(msg googleProto.Message, reporter chan *proto.ProgressMetrics) {
		core.RunBulkSimAsync(msg.(*proto.BulkSimRequest), reporter)
	}},
}

func handleAsyncAPI(w http.ResponseWriter, r *http.Request, addNewSim simProgReportCreator) {
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

	// Generate a new async simulation, and get back the ID and reporting function.
	id, cacheProgressFunc := addNewSim()

	// Now launch a background process that pulls progress reports off the reporter channel
	// and pushes it into the async progress cache.
	go func() {
		for {
			select {
			case <-time.After(time.Hour):
				return // if we get no progress after an hour, exit
			case progMetric, ok := <-reporter:
				if !ok {
					return
				}
				cacheProgressFunc(progMetric)
				if progMetric.FinalRaidResult != nil || progMetric.FinalWeightResult != nil {
					return
				}
			}
		}
	}()

	protoResult := &proto.AsyncAPIResult{
		ProgressId: id,
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

func setupAsyncServer() {

	// Hold all state for in-flight async processes here.
	type asyncProgress struct {
		latestProgress atomic.Value
	}
	progMut := &sync.RWMutex{} // mutex for progresses map
	progresses := map[string]*asyncProgress{}

	// addNewSim just stores progress data for a new running simulation into the state above.
	addNewSim := func() (string, progReport) {
		newID := uuid.NewV4().String()
		simProgress := &asyncProgress{}
		simProgress.latestProgress.Store(&proto.ProgressMetrics{})
		progMut.Lock()
		progresses[newID] = simProgress
		progMut.Unlock()

		return newID, func(newProg *proto.ProgressMetrics) {
			// caches progress into the progress map indexed by the ID.
			// This can later be fetched by the async progress endpoint.
			simProgress.latestProgress.Store(newProg)
		}
	}

	// All async handlers here will call the addNewSim, generating a new UUID and cached progress state.
	http.HandleFunc("/statWeightsAsync", func(w http.ResponseWriter, r *http.Request) {
		handleAsyncAPI(w, r, addNewSim)
	})
	http.HandleFunc("/raidSimAsync", func(w http.ResponseWriter, r *http.Request) {
		handleAsyncAPI(w, r, addNewSim)
	})
	http.HandleFunc("/bulkSimAsync", func(w http.ResponseWriter, r *http.Request) {
		handleAsyncAPI(w, r, addNewSim)
	})

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
		progMut.RLock()
		progress, ok := progresses[msg.ProgressId]
		progMut.RUnlock()
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
		if latest.FinalRaidResult != nil || latest.FinalWeightResult != nil {
			progMut.Lock()
			delete(progresses, msg.ProgressId)
			progMut.Unlock()
		}
		w.Header().Add("Content-Type", "application/x-protobuf")
		w.Write(outbytes)
	})
}

func runServer(useFS bool, host string, launchBrowser bool, simName string, wasm bool, inputReader *bufio.Reader) {
	var fs http.Handler
	if useFS {
		log.Printf("Using local file system for development.")
		fs = http.FileServer(http.Dir("./dist"))
	} else {
		log.Printf("Embedded file server running.")
		fs = http.FileServer(http.FS(dist.FS))
	}

	http.HandleFunc("/version", func(resp http.ResponseWriter, req *http.Request) {
		msg := fmt.Sprintf(`{"version": "%s", "outdated": %d}`, Version, outdated)
		resp.Write([]byte(msg))
	})
	http.HandleFunc("/statWeights", handleAPI)
	http.HandleFunc("/computeStats", handleAPI)
	http.HandleFunc("/individualSim", handleAPI)
	http.HandleFunc("/raidSim", handleAPI)
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
		url := fmt.Sprintf("http://localhost%s/wotlk/%s", host, simName)
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
		case "quit":
			os.Exit(1)
		case "?":
			fmt.Printf("Commands:\n\tprofile - start a CPU profile for debugging performance\n\tquit - exits\n\n")
		case "":
			// nothing.
		default:
			fmt.Printf("Unknown command: '%s'", command)
		}
	}
}

type apiHandler struct {
	msg    func() googleProto.Message
	handle func(googleProto.Message) googleProto.Message
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
