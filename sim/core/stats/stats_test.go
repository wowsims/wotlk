package stats

import (
	"strings"
	"testing"

	"github.com/wowsims/classic/sim/core/proto"
)

func TestStatsAdd(t *testing.T) {
	a := Stats{
		Intellect: 1,
	}
	b := Stats{
		Intellect: 1,
	}
	expectedResult := Stats{
		Intellect: 2,
	}

	result := a.Add(b)

	if !result.Equals(expectedResult) {
		t.Fatalf("Expected equal stats but were not equal: %s, %s", result, expectedResult)
	}
}

func TestStatsEquals_Success(t *testing.T) {
	a := Stats{
		Intellect: 1,
	}
	b := Stats{
		Intellect: 1,
	}

	if !a.Equals(b) {
		t.Fatalf("Expected equal stats but were not equal: %s, %s", a, b)
	}
}

func TestStatsEquals_Failure(t *testing.T) {
	a := Stats{
		Intellect: 1,
	}
	b := Stats{
		Intellect: 0,
	}

	if a.Equals(b) {
		t.Fatalf("Expected not equal stats but were equal: %s, %s", a, b)
	}
}

func TestStatsEqualsWithTolerance_Success(t *testing.T) {
	a := Stats{
		Intellect: 1,
	}
	b := Stats{
		Intellect: 0.5,
	}

	if !a.EqualsWithTolerance(b, 0.5) {
		t.Fatalf("Expected equal stats but were not equal: %s, %s", a, b)
	}
}

func TestStatsEqualsWithTolerance_Failure(t *testing.T) {
	a := Stats{
		Intellect: 1,
	}
	b := Stats{
		Intellect: 0.4,
	}

	if a.EqualsWithTolerance(b, 0.5) {
		t.Fatalf("Expected not equal stats but were equal: %s, %s", a, b)
	}
}

func TestStatsProtoInSync(t *testing.T) {
	d := proto.Stat_StatStrength.Descriptor().Values()
	if d.Len() != int(Len) {
		t.Fatalf("Unequal number of stats defined in proto.Stats (%d) and Go (%d)", d.Len(), Len)
	}

	for i := 0; i < d.Len(); i++ {
		enum := d.Get(i)
		protoName := enum.Name()
		goName := Stat(enum.Number()).StatName()
		sanitizedGoName := strings.ReplaceAll(goName, " ", "")
		if string(protoName) != "Stat"+sanitizedGoName {
			t.Fatalf("Encountered stat enum %d in proto.Stats with name %s differs from Go enum name %s (ignoring Stat prefix)", enum.Number(), protoName, goName)
		}
	}
}
