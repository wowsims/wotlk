package stats

import (
	"testing"
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

func TestStatDependencyManagerApplyStatDependencies_Success(t *testing.T) {
	stats := Stats{
		Stamina:   1,
		Intellect: 1,
	}
	sdm := StatDependencyManager{}
	// Add these in the opposite order we expect them to be applied, to test the sorting.
	sdm.AddStatDependency(StatDependency{
		SourceStat:   Intellect,
		ModifiedStat: Intellect,
		Modifier: func(intellect float64, _ float64) float64 {
			return intellect * 2
		},
	})
	sdm.AddStatDependency(StatDependency{
		SourceStat:   Stamina,
		ModifiedStat: Intellect,
		Modifier: func(stamina float64, intellect float64) float64 {
			return intellect + stamina
		},
	})
	sdm.AddStatDependency(StatDependency{
		SourceStat:   Stamina,
		ModifiedStat: Stamina,
		Modifier: func(stamina float64, _ float64) float64 {
			return stamina + 1
		},
	})
	sdm.Finalize()
	expectedResult := Stats{
		Stamina:   2,
		Intellect: 6,
	}

	result := sdm.ApplyStatDependencies(stats)

	if !result.Equals(expectedResult) {
		t.Fatalf("Expected equal stats but were not equal: %s, %s", result, expectedResult)
	}
}
