package core

import (
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

// A presim is a full simulation run with multiple iterations, as a preparation
// step for testing out settings before starting the recorded iterations.
//
// To use this, just implement this interface on your Agent.
//
// If you don't know what this is, you probably don't need it.
type Presimmer interface {
	GetPresimOptions(proto.Player) *PresimOptions
}

// Controls the presim behavior for 1 Agent.
type PresimOptions struct {
	// Called once before each presim round.
	//
	// Modify the player parameter to use whatever player options are desired
	// for the presim.
	SetPresimPlayerOptions func(player *proto.Player)

	// Called once after each presim round to provide the results.
	//
	// Should return true if this Agent is done running presims, and false otherwise.
	OnPresimResult func(presimResult proto.UnitMetrics, iterations int32, duration time.Duration) bool
}

func (sim *Simulation) runPresims(request proto.RaidSimRequest) *proto.RaidSimResult {
	const numPresimIterations = 100

	// Run presims if requested.
	raidPresimOptions := make([]*PresimOptions, 25)
	remainingAgents := 0
	for _, party := range sim.Raid.Parties {
		for _, player := range party.Players {
			presimmer, ok := player.(Presimmer)
			if !ok {
				continue
			}

			partyConfig := request.Raid.Parties[player.GetCharacter().Party.Index]
			playerConfig := partyConfig.Players[player.GetCharacter().PartyIndex]

			presimOptions := presimmer.GetPresimOptions(*playerConfig)
			if presimOptions == nil {
				continue
			}

			raidPresimOptions[player.GetCharacter().Index] = presimOptions
			remainingAgents++
		}
	}

	// Base presim request.
	// Define this outside the loop so that, as Agents iteratively update their
	// settings, we keep the most recent settings even after that Agent is
	// done with presims.
	presimRequest := googleProto.Clone(&request).(*proto.RaidSimRequest)
	presimRequest.SimOptions.RandomSeed = 1
	presimRequest.SimOptions.Debug = false
	presimRequest.SimOptions.DebugFirstIteration = false
	presimRequest.SimOptions.Iterations = numPresimIterations
	duration := DurationFromSeconds(presimRequest.Encounter.Duration)

	var lastResult *proto.RaidSimResult

	doOne := sim.Encounter.EndFightAtHealth > 0
	for doOne || remainingAgents > 0 {
		// ** Run a presim round. **

		// Let each Agent modify their own settings.
		for partyIdx, party := range sim.Raid.Parties {
			partyConfig := presimRequest.Raid.Parties[partyIdx]
			for _, player := range party.Players {
				playerConfig := partyConfig.Players[player.GetCharacter().PartyIndex]

				presimOptions := raidPresimOptions[player.GetCharacter().Index]
				if presimOptions == nil {
					continue
				}

				presimOptions.SetPresimPlayerOptions(playerConfig)
			}
		}

		// Run the presim.
		presimResult := runSim(*presimRequest, nil, true)
		lastResult = presimResult

		if presimResult.ErrorResult != "" {
			break
		}

		// Provide each Agent with their own results.
		for partyIdx, party := range sim.Raid.Parties {
			partyMetrics := presimResult.RaidMetrics.Parties[partyIdx]
			for _, player := range party.Players {
				playerMetrics := partyMetrics.Players[player.GetCharacter().PartyIndex]
				presimOptions := raidPresimOptions[player.GetCharacter().Index]
				if presimOptions != nil {
					done := presimOptions.OnPresimResult(*playerMetrics, numPresimIterations, duration)
					if done {
						raidPresimOptions[player.GetCharacter().Index] = nil
						remainingAgents--
					}
				}
			}
		}
		doOne = false
	}
	return lastResult
}
