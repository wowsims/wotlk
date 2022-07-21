package stats

// Bonuses for a single stat
type Bonuses struct {
	Ratio float64          // ratio added to all stat gains from this stat.
	Deps  map[Stat]float64 // ratio added to Stat when this stat is changed.
}
