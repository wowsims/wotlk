package sim

import (
	"testing"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
	"google.golang.org/protobuf/encoding/protojson"

	balanceDruid "github.com/wowsims/tbc/sim/druid/balance"
	hunter "github.com/wowsims/tbc/sim/hunter"
	shadowPriest "github.com/wowsims/tbc/sim/priest/shadow"
	elementalShaman "github.com/wowsims/tbc/sim/shaman/elemental"
	enhancementShaman "github.com/wowsims/tbc/sim/shaman/enhancement"
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

var P1BalanceDruid = &proto.Player{
	Name:      "P1 Boomkin",
	Race:      proto.Race_RaceTauren,
	Class:     proto.Class_ClassDruid,
	Equipment: balanceDruid.P1Gear,
	Consumes:  balanceDruid.FullConsumes,
	Spec:      balanceDruid.PlayerOptionsAdaptive,
	Buffs:     balanceDruid.FullIndividualBuffs,
}

var P1ElementalShaman = &proto.Player{
	Name:      "P1 Ele Shaman",
	Race:      proto.Race_RaceOrc,
	Class:     proto.Class_ClassShaman,
	Equipment: elementalShaman.P1Gear,
	Consumes:  elementalShaman.FullConsumes,
	Spec:      elementalShaman.PlayerOptionsAdaptive,
	Buffs:     elementalShaman.FullIndividualBuffs,
}

var P1ShadowPriest = &proto.Player{
	Name:      "P1 Shadow Priest",
	Race:      proto.Race_RaceUndead,
	Class:     proto.Class_ClassPriest,
	Equipment: shadowPriest.P1Gear,
	Consumes:  shadowPriest.FullConsumes,
	Spec:      shadowPriest.PlayerOptionsIdeal,
	Buffs:     shadowPriest.FullIndividualBuffs,
}

var P1EnhancementShaman = &proto.Player{
	Name:      "P1 Enh Shaman",
	Race:      proto.Race_RaceOrc,
	Class:     proto.Class_ClassShaman,
	Equipment: enhancementShaman.Phase2Gear,
	Consumes:  enhancementShaman.FullConsumes,
	Spec:      enhancementShaman.PlayerOptionsBasic,
	Buffs:     enhancementShaman.FullIndividualBuffs,
}

var P1BMHunter = &proto.Player{
	Name:      "P1 BM Hunter",
	Race:      proto.Race_RaceOrc,
	Class:     proto.Class_ClassHunter,
	Equipment: hunter.P1Gear,
	Consumes:  hunter.FullConsumes,
	Spec:      hunter.PlayerOptionsBasic,
	Buffs:     hunter.FullIndividualBuffs,
}

var BasicRaid = &proto.Raid{
	Parties: []*proto.Party{
		&proto.Party{
			Players: []*proto.Player{
				P1BalanceDruid,
				P1ElementalShaman,
				P1EnhancementShaman,
				P1ShadowPriest,
			},
		},
		&proto.Party{
			Players: []*proto.Player{
				P1BMHunter,
			},
		},
	},
	StaggerStormstrikes: true,
}

// Tests that we don't crash with various combinations of empty parties / blank players.
func TestSparseRaid(t *testing.T) {
	sparseRaid := &proto.Raid{
		Parties: []*proto.Party{
			&proto.Party{},
			&proto.Party{
				Players: []*proto.Player{
					&proto.Player{},
					P1ElementalShaman,
					&proto.Player{},
				},
			},
			&proto.Party{
				Players: []*proto.Player{
					&proto.Player{},
					&proto.Player{},
				},
			},
		},
		StaggerStormstrikes: true,
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
	rsr := &proto.RaidSimRequest{
		Raid:       BasicRaid,
		Encounter:  STEncounter,
		SimOptions: SimOptions,
	}

	core.RaidSimTest("P1 ST", t, rsr, 6323.79)
}

func testRaidString(t *testing.T, raidString string) {
	rsr := &proto.RaidSimRequest{}

	data := []byte(raidString)
	if err := protojson.Unmarshal(data, rsr); err != nil {
		panic(err)
	}

	core.RunRaidSim(rsr)
	//core.RaidSimTest("Fixed Raid", t, rsr, 10000.00)
}

// To quickly debug raid sim issues, uncomment this test and copy in a request string.
// func TestFixedRaid(t *testing.T) {
// 	testRaidString(t, `
// 	`)
// }
