package sim

import (
	"testing"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	RegisterAll()
}

var SimOptions = &proto.SimOptions{
	Iterations: 1,
	IsTest:     true,
}

var StandardTarget = &proto.Target{
	Stats:   stats.Stats{stats.Armor: 7684}.ToFloatArray(),
	MobType: proto.MobType_MobTypeDemon,
}

var STEncounter = &proto.Encounter{
	Duration: 300,
	Targets: []*proto.Target{
		StandardTarget,
	},
}

var BasicRaid = &proto.Raid{
	Parties: []*proto.Party{
		{
			Players: []*proto.Player{},
		},
		{
			Players: []*proto.Player{},
		},
	},
}

// Tests that we don't crash with various combinations of empty parties / blank players.
func TestSparseRaid(t *testing.T) {
	sparseRaid := &proto.Raid{
		Parties: []*proto.Party{
			{},
			{
				Players: []*proto.Player{
					{},
					{},
				},
			},
			{
				Players: []*proto.Player{
					{},
					{},
				},
			},
		},
	}

	rsr := &proto.RaidSimRequest{
		Raid:       sparseRaid,
		Encounter:  STEncounter,
		SimOptions: SimOptions,
	}

	core.RunRaidSim(rsr)
	// Don't need to check results, as long as it doesn't crash we're fine.
}

func TestBasicRaid(t *testing.T) {
	t.Skip()
	rsr := &proto.RaidSimRequest{
		Raid:       BasicRaid,
		Encounter:  STEncounter,
		SimOptions: SimOptions,
	}

	core.RaidSimTest("P1 ST", t, rsr, 6323.79)
}

// To quickly debug raid sim issues, uncomment this test and copy in a request string.
/*
func testRaidString(t *testing.T, raidString string) {
	rsr := &proto.RaidSimRequest{}

	data := []byte(raidString)
	if err := protojson.Unmarshal(data, rsr); err != nil {
		panic(err)
	}

	core.RunRaidSim(rsr)
	//core.RaidSimTest("Fixed Raid", t, rsr, 10000.00)
}

func TestFixedRaid(t *testing.T) {
 	testRaidString(t, `
 	`)
}
*/
