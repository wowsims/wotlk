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
