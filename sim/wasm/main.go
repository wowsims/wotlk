//go:build wasm
// +build wasm

package main

import (
	"log"
	"runtime/debug"
	"syscall/js"

	"github.com/wowsims/wotlk/sim"
	"github.com/wowsims/wotlk/sim/core"
	proto "github.com/wowsims/wotlk/sim/core/proto"
	protojson "google.golang.org/protobuf/encoding/protojson"
	googleProto "google.golang.org/protobuf/proto"
)

func init() {
	sim.RegisterAll()
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("computeStats", js.FuncOf(computeStats))
	js.Global().Set("computeStatsJson", js.FuncOf(computeStatsJson))
	js.Global().Set("gearList", js.FuncOf(gearList))
	js.Global().Set("gearListJson", js.FuncOf(gearListJson))
	js.Global().Set("raidSim", js.FuncOf(raidSim))
	js.Global().Set("raidSimJson", js.FuncOf(raidSimJson))
	js.Global().Set("raidSimAsync", js.FuncOf(raidSimAsync))
	js.Global().Set("statWeights", js.FuncOf(statWeights))
	js.Global().Set("statWeightsAsync", js.FuncOf(statWeightsAsync))
	js.Global().Call("wasmready")
	<-c
}

func computeStats(this js.Value, args []js.Value) (response interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errStr := ""
			switch errt := err.(type) {
			case string:
				errStr = errt
			case error:
				errStr = errt.Error()
			}

			errStr += "\nStack Trace:\n" + string(debug.Stack())
			result := &proto.ComputeStatsResult{
				ErrorResult: errStr,
			}
			outbytes, err := googleProto.Marshal(result)
			if err != nil {
				log.Printf("[ERROR] Failed to marshal error (%s) result: %s", errStr, err.Error())
				return
			}
			outArray := js.Global().Get("Uint8Array").New(len(outbytes))
			js.CopyBytesToJS(outArray, outbytes)
			response = outArray
		}
	}()
	csr := &proto.ComputeStatsRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), csr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.ComputeStats(csr)

	outbytes, err := googleProto.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}

	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	response = outArray
	return response
}

func computeStatsJson(this js.Value, args []js.Value) (response interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errStr := ""
			switch errt := err.(type) {
			case string:
				errStr = errt
			case error:
				errStr = errt.Error()
			}

			errStr += "\nStack Trace:\n" + string(debug.Stack())
			result := &proto.ComputeStatsResult{
				ErrorResult: errStr,
			}

			output, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(result)
			if err != nil {
				log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
			}
			response = js.ValueOf(string(output))
		}
	}()
	csr := &proto.ComputeStatsRequest{}
	log.Printf("Compute stats request: %s", getArgsJson(args[0]))
	if err := protojson.Unmarshal(getArgsJson(args[0]), csr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.ComputeStats(csr)

	output, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}
	response = js.ValueOf(string(output))
	return response
}

func gearList(this js.Value, args []js.Value) interface{} {
	glr := &proto.GearListRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), glr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.GetGearList(glr)

	outbytes, err := googleProto.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}

	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	return outArray
}

func gearListJson(this js.Value, args []js.Value) interface{} {
	glr := &proto.GearListRequest{}
	if err := protojson.Unmarshal(getArgsJson(args[0]), glr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.GetGearList(glr)

	output, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}
	response := js.ValueOf(string(output))
	return response
}

func raidSim(this js.Value, args []js.Value) interface{} {
	rsr := &proto.RaidSimRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), rsr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.RunRaidSim(rsr)

	outbytes, err := googleProto.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}

	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	return outArray
}

func raidSimJson(this js.Value, args []js.Value) interface{} {
	rsr := &proto.RaidSimRequest{}
	if err := protojson.Unmarshal(getArgsJson(args[0]), rsr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.RunRaidSim(rsr)

	output, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}
	response := js.ValueOf(string(output))
	return response
}

func raidSimAsync(this js.Value, args []js.Value) interface{} {
	rsr := &proto.RaidSimRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), rsr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	reporter := make(chan *proto.ProgressMetrics, 100)

	go core.RunRaidSimAsync(rsr, reporter)
	return processAsyncProgress(args[1], reporter)
}

func statWeights(this js.Value, args []js.Value) interface{} {
	swr := &proto.StatWeightsRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), swr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	result := core.StatWeights(swr)

	outbytes, err := googleProto.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
		return nil
	}

	outArray := js.Global().Get("Uint8Array").New(len(outbytes))
	js.CopyBytesToJS(outArray, outbytes)

	return outArray
}

func statWeightsAsync(this js.Value, args []js.Value) interface{} {
	rsr := &proto.StatWeightsRequest{}
	if err := googleProto.Unmarshal(getArgsBinary(args[0]), rsr); err != nil {
		log.Printf("Failed to parse request: %s", err)
		return nil
	}
	reporter := make(chan *proto.ProgressMetrics, 100)
	core.StatWeightsAsync(rsr, reporter)

	result := processAsyncProgress(args[1], reporter)
	return result
}

// Assumes args[0] is a Uint8Array
func getArgsBinary(value js.Value) []byte {
	data := make([]byte, value.Get("length").Int())
	js.CopyBytesToGo(data, value)
	return data
}

func getArgsJson(value js.Value) []byte {
	str := value.String()
	return []byte(str)
}

func processAsyncProgress(progFunc js.Value, reporter chan *proto.ProgressMetrics) js.Value {
reader:
	for {
		// TODO: cleanup so we dont collect these
		select {
		case progMetric, ok := <-reporter:
			if !ok {
				break reader
			}
			outbytes, err := googleProto.Marshal(progMetric)
			if err != nil {
				log.Printf("[ERROR] Failed to marshal result: %s", err.Error())
				return js.Undefined()
			}

			outArray := js.Global().Get("Uint8Array").New(len(outbytes))
			js.CopyBytesToJS(outArray, outbytes)
			progFunc.Invoke(outArray)

			if progMetric.FinalWeightResult != nil || progMetric.FinalRaidResult != nil {
				return outArray
			}
		}
	}

	return js.Undefined()
}
