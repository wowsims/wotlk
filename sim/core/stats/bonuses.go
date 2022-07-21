package stats

// Bonuses for a single stat
type Bonuses struct {
	Multiplier float64          // multiplier added to all stat gains from this stat.
	Deps       map[Stat]float64 // multiplier added to Stat when this stat is changed.
}
