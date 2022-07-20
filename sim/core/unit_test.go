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
	}

	unit.MultiplyStat(stats.Intellect, 2)
	unit.AddStatDependency(stats.Stamina, stats.Intellect, 1)
	unit.MultiplyStat(stats.Stamina, 2)
	unit.AddStatDependency(stats.Agility, stats.Stamina, 1)
	unit.finalizeStatDeps()

	expectedResult := stats.Stats{
		stats.Stamina:   6,
		stats.Intellect: 14,
		stats.Agility:   2,
	}

	unit.stats = unit.ApplyStatDependencies(baseStat)

	if !unit.stats.Equals(expectedResult) {
		t.Fatalf("Stats do not match:\nActual: %s\nExpected: %s", unit.stats, expectedResult)
	}

	unit.Env, _ = NewEnvironment(proto.Raid{}, proto.Encounter{})
	unit.MultiplyStatDynamic(nil, stats.Agility, 2)

	result2 := stats.Stats{
		stats.Stamina:   10,
		stats.Intellect: 22,
		stats.Agility:   4,
	}
	if !unit.stats.Equals(result2) {
		t.Fatalf("Updated stats do not match:\nActual: %s\nExpected: %s", unit.stats, result2)
	}
}

func TestCircularStatDependencies(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Did not catch circular dependency in stats")
		}
	}()
	unit := Unit{}
	unit.AddStatDependency(stats.Stamina, stats.Intellect, 1)
	unit.AddStatDependency(stats.Agility, stats.Stamina, 1)
	unit.AddStatDependency(stats.Intellect, stats.Agility, 1)
	unit.finalizeStatDeps()
}
