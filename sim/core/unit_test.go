package core

import (
	"testing"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func TestStatDependencies(t *testing.T) {
	unit := Unit{}

	baseStat := stats.Stats{
		stats.Stamina:   1,
		stats.Intellect: 1,
		stats.Agility:   2,
		stats.Spirit:    1,
	}

	unit.AddStatDependency(stats.Intellect, stats.Intellect, 1.0+1)
	unit.AddStatDependency(stats.Stamina, stats.Intellect, 1.0+1)
	unit.AddStatDependency(stats.Stamina, stats.Stamina, 1.0+1)
	unit.AddStatDependency(stats.Agility, stats.Stamina, 1.0+1)
	unit.finalizeStatDeps()

	expectedResult := stats.Stats{
		stats.Stamina:   6,
		stats.Intellect: 14,
		stats.Agility:   2,
		stats.Spirit:    1,
	}

	unit.stats = unit.applyStatDependencies(baseStat)

	if !unit.stats.Equals(expectedResult) {
		t.Fatalf("Stats do not match:\nActual: %s\nExpected: %s", unit.stats, expectedResult)
	}

	unit.Env, _ = NewEnvironment(proto.Raid{}, proto.Encounter{})
	unit.AddStatDependencyDynamic(nil, stats.Agility, stats.Agility, 2)

	result2 := stats.Stats{
		stats.Stamina:   10,
		stats.Intellect: 22,
		stats.Agility:   4,
		stats.Spirit:    1,
	}
	if !unit.stats.Equals(result2) {
		t.Fatalf("Updated stats do not match:\nActual: %s\nExpected: %s", unit.stats, result2)
	}

	unit.AddStatDependencyDynamic(nil, stats.Spirit, stats.Spirit, 3)

	result3 := stats.Stats{
		stats.Stamina:   10,
		stats.Intellect: 22,
		stats.Agility:   4,
		stats.Spirit:    3,
	}
	if !unit.stats.Equals(result3) {
		t.Fatalf("Updated stats do not match:\nActual: %s\nExpected: %s", unit.stats, result3)
	}

	unit.AddStatDependencyDynamic(nil, stats.Spirit, stats.Spirit, 1.0/3.0)
	result4 := stats.Stats{
		stats.Stamina:   10,
		stats.Intellect: 22,
		stats.Agility:   4,
		stats.Spirit:    1,
	}
	if !unit.stats.Equals(result4) {
		t.Fatalf("Updated stats do not match:\nActual: %s\nExpected: %s", unit.stats, result4)
	}

	unit.AddStatDependencyDynamic(nil, stats.Agility, stats.Spirit, 1.75)
	result5 := stats.Stats{
		stats.Stamina:   10,
		stats.Intellect: 22,
		stats.Agility:   4,
		stats.Spirit:    4,
	}
	if !unit.stats.Equals(result5) {
		t.Fatalf("Updated stats do not match:\nActual: %s\nExpected: %s", unit.stats, result5)
	}

	unit.AddStatDependencyDynamic(nil, stats.Agility, stats.Spirit, 1/1.75)
	result6 := stats.Stats{
		stats.Stamina:   10,
		stats.Intellect: 22,
		stats.Agility:   4,
		stats.Spirit:    1,
	}
	if !unit.stats.Equals(result6) {
		t.Fatalf("Updated stats do not match:\nActual: %s\nExpected: %s", unit.stats, result6)
	}

}

func TestCircularStatDependencies(t *testing.T) {
	unit := Unit{}
	unit.AddStatDependency(stats.Stamina, stats.Intellect, 1.0+1)
	unit.AddStatDependency(stats.Agility, stats.Stamina, 1.0+1)
	unit.AddStatDependency(stats.Intellect, stats.Agility, 1.0+1)

	validateFail := func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Did not catch circular dependency in stats")
			}
		}()
		unit.finalizeStatDeps()
	}

	validateFail()

	// clear circle
	unit.statBonuses[stats.Intellect].Deps = map[stats.Stat]float64{}
	unit.finalizeStatDeps()

	// Now add diamond dependency
	unit.AddStatDependency(stats.Agility, stats.Spirit, 1.0+1)
	unit.AddStatDependency(stats.Spirit, stats.Intellect, 1.0+1)
	unit.finalizeStatDeps()

}

func TestMultipleStatDep(t *testing.T) {
	unit := Unit{}

	baseStat := stats.Stats{
		stats.Intellect:  100,
		stats.SpellPower: 100,
	}

	unit.AddStatDependency(stats.Intellect, stats.SpellPower, 1.2)
	unit.AddStatDependency(stats.Intellect, stats.SpellPower, 1.2)
	unit.AddStatDependency(stats.Intellect, stats.Intellect, 1.2)
	unit.finalizeStatDeps()

	expectedResult := stats.Stats{
		stats.Intellect:  100 * 1.2,
		stats.SpellPower: 100 + ((100 * 1.2) * ((1.2 * 1.2) - 1)),
	}

	unit.stats = unit.applyStatDependencies(baseStat)

	if !unit.stats.Equals(expectedResult) {
		t.Fatalf("Stats do not match:\nActual: %s\nExpected: %s", unit.stats, expectedResult)
	}
}
