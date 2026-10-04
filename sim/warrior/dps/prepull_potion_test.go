package dps

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	indestructiblePotionItemID = 40093
	potionOfSpeedItemID        = 40211
)

// The generic Potion placeholder resolves to one potion per phase, so it cannot
// express two different potions before the pull. An APL can name the potion
// spells themselves instead, and the shared potion cooldown gates how many are
// used: the standard opening is an Indestructible Potion at -61 s, whose
// one-minute cooldown expires at -1 s, a Potion of Speed at -1 s, whose
// cooldown expires at +59 s, and a third potion at +59 s.
func TestMultiplePrepullPotions(t *testing.T) {
	scenario := runPotionScenario(t,
		&proto.Consumes{
			DefaultPotion: proto.Potions_PotionOfSpeed,
			PrepopPotion:  proto.Potions_IndestructiblePotion,
		},
		`{
			"prepullActions": [
				{"action":{"castSpell":{"spellId":{"itemId":40093}}},"doAtValue":{"const":{"val":"-61s"}}},
				{"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}}
			],
			"priorityList": [
				{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}}}
			]
		}`)

	// The APL UI builds its prepull action picker from this spell metadata.
	// Both potions must be offered as castable and must not be marked
	// encounter-only, or the prepull section cannot select them.
	for _, itemID := range []int32{indestructiblePotionItemID, potionOfSpeedItemID} {
		found := false
		for _, spell := range scenario.metadata.Spells {
			if spell.Id.GetItemId() != itemID {
				continue
			}
			found = true
			if !spell.IsCastable {
				t.Errorf("potion %d is not offered as a castable APL spell", itemID)
			}
			if spell.EncounterOnly {
				t.Errorf("potion %d is marked encounter-only, so the prepull section cannot select it", itemID)
			}
		}
		if !found {
			t.Errorf("no spell metadata found for potion %d", itemID)
		}
	}

	if warnings := countWarnings(scenario.prepullWarnings); warnings != 0 {
		t.Errorf("prepull actions produced %d warnings: %v", warnings, scenario.prepullWarnings)
	}

	want := []potionCast{
		{indestructiblePotionItemID, -61 * time.Second},
		{potionOfSpeedItemID, -1 * time.Second},
	}
	if len(scenario.casts) != 3 {
		t.Fatalf("expected 3 potion casts (2 during the prepull, 1 once the cooldown expires), got %v", scenario.casts)
	}
	for i, want := range want {
		if scenario.casts[i] != want {
			t.Errorf("potion cast %d: got %d at %s, want %d at %s", i, scenario.casts[i].itemID, scenario.casts[i].at, want.itemID, want.at)
		}
	}

	// The second potion's one-minute cooldown expires at +59 s; the third potion
	// is cast at the next APL evaluation after that.
	if scenario.casts[2].at < 59*time.Second || scenario.casts[2].at >= 60*time.Second {
		t.Errorf("third potion cast at %s, want it once the shared cooldown expires at +59s", scenario.casts[2].at)
	}
}

// The common cases, which must not change: a single configured potion each
// place, and no potions at all. Potions are used through the generic
// placeholder, which resolves by phase and reports a warning when its potion
// is not configured.
func TestPrepullPotionControls(t *testing.T) {
	placeholderAPL := `{
		"prepullActions": [
			{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1s"}}}
		],
		"priorityList": [
			{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}}}
		]
	}`

	for _, test := range []struct {
		name             string
		consumes         *proto.Consumes
		wantCasts        []potionCast
		wantPrepullWarns int
		wantCombatWarns  int
	}{
		{
			name:      "prepull potion only",
			consumes:  &proto.Consumes{PrepopPotion: proto.Potions_IndestructiblePotion},
			wantCasts: []potionCast{{indestructiblePotionItemID, -1 * time.Second}},
			// The combat placeholder has no combat potion to resolve to.
			wantCombatWarns: 1,
		},
		{
			name:      "combat potion only",
			consumes:  &proto.Consumes{DefaultPotion: proto.Potions_PotionOfSpeed},
			wantCasts: []potionCast{{potionOfSpeedItemID, 0}},
			// The prepull placeholder has no prepull potion to resolve to.
			wantPrepullWarns: 1,
		},
		{
			name:             "no potions",
			consumes:         &proto.Consumes{},
			wantPrepullWarns: 1,
			wantCombatWarns:  1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			scenario := runPotionScenario(t, test.consumes, placeholderAPL)

			if got := countWarnings(scenario.prepullWarnings); got != test.wantPrepullWarns {
				t.Errorf("got %d prepull warnings %v, want %d", got, scenario.prepullWarnings, test.wantPrepullWarns)
			}
			if got := len(scenario.combatWarnings); got != test.wantCombatWarns {
				t.Errorf("got %d combat warnings %v, want %d", got, scenario.combatWarnings, test.wantCombatWarns)
			}
			if len(scenario.casts) != len(test.wantCasts) {
				t.Fatalf("got potion casts %v, want %v", scenario.casts, test.wantCasts)
			}
			for i, want := range test.wantCasts {
				if scenario.casts[i] != want {
					t.Errorf("potion cast %d: got %d at %s, want %d at %s", i, scenario.casts[i].itemID, scenario.casts[i].at, want.itemID, want.at)
				}
			}
		})
	}
}

type potionCast struct {
	itemID int32
	at     time.Duration
}

type potionScenario struct {
	casts           []potionCast
	metadata        *proto.UnitMetadata
	prepullWarnings [][]string
	combatWarnings  []string
}

// runPotionScenario runs a fury warrior fight with the given consumes and APL,
// and reports the potion casts and the APL warnings from the sim's own log and
// spell metadata.
func runPotionScenario(t *testing.T, consumes *proto.Consumes, aplJSON string) potionScenario {
	t.Helper()

	rotation := core.GetAplRotation("../../../ui/warrior/apls", "fury").Rotation
	var custom proto.APLRotation
	if err := protojson.Unmarshal([]byte(aplJSON), &custom); err != nil {
		t.Fatalf("failed to parse rotation: %s", err)
	}
	rotation.PrepullActions = custom.PrepullActions
	rotation.PriorityList = append(custom.PriorityList, rotation.PriorityList...)

	raid := core.SinglePlayerRaidProto(
		&proto.Player{
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarrior,
			Equipment:     core.GetGearSet("../../../ui/warrior/gear_sets", "p1_fury").GearSet,
			Consumes:      consumes,
			Spec:          PlayerOptionsFury,
			TalentsString: FuryTalents,
			Glyphs:        FuryGlyphs,
			Buffs:         core.FullIndividualBuffs,
			Rotation:      rotation,
		},
		core.FullPartyBuffs, core.FullRaidBuffs, core.FullDebuffs)
	encounter := &proto.Encounter{
		Duration: 120,
		Targets:  []*proto.Target{core.NewDefaultTarget()},
	}

	// ComputeStats reports the APL warnings, including the ones produced by the
	// fake prepull it runs, and the spell metadata the APL UI pickers use.
	stats := core.ComputeStats(&proto.ComputeStatsRequest{Raid: raid, Encounter: encounter})
	if stats.ErrorResult != "" {
		t.Fatalf("compute stats failed: %s", stats.ErrorResult)
	}
	playerStats := stats.RaidStats.Parties[0].Players[0]

	result := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       raid,
		Encounter:  encounter,
		SimOptions: &proto.SimOptions{Iterations: 1, DebugFirstIteration: true, IsTest: true, RandomSeed: 101},
	})
	if result.ErrorResult != "" {
		t.Fatalf("sim failed: %s", result.ErrorResult)
	}

	scenario := potionScenario{
		casts:           potionCastTimes(t, result.Logs),
		metadata:        playerStats.Metadata,
		prepullWarnings: make([][]string, 0, len(playerStats.RotationStats.PrepullActions)),
	}
	for _, action := range playerStats.RotationStats.PrepullActions {
		scenario.prepullWarnings = append(scenario.prepullWarnings, action.Warnings)
	}
	if len(playerStats.RotationStats.PriorityList) > 0 {
		scenario.combatWarnings = playerStats.RotationStats.PriorityList[0].Warnings
	}
	return scenario
}

func countWarnings(warningsByAction [][]string) int {
	count := 0
	for _, warnings := range warningsByAction {
		count += len(warnings)
	}
	return count
}

var potionCastRegex = regexp.MustCompile(`^\[(-?[0-9.]+)\] .*Casting \{ItemID: ([0-9]+)\}`)

// potionCastTimes returns every potion cast in the sim log, which stamps each
// cast with the time it happened, in log order.
func potionCastTimes(t *testing.T, log string) []potionCast {
	t.Helper()

	var casts []potionCast
	for _, line := range strings.Split(log, "\n") {
		match := potionCastRegex.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		itemID, err := strconv.Atoi(match[2])
		if err != nil {
			t.Fatalf("failed to parse item id %q: %s", match[2], err)
		}
		if int32(itemID) != indestructiblePotionItemID && int32(itemID) != potionOfSpeedItemID {
			continue
		}

		seconds, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			t.Fatalf("failed to parse cast time %q: %s", match[1], err)
		}

		casts = append(casts, potionCast{
			itemID: int32(itemID),
			at:     core.DurationFromSeconds(seconds),
		})
	}
	return casts
}
