package core

type rune struct {
	exists        bool
	is_death_rune bool
}

type runeSystem struct {
	unit *Unit

	bloodRunes  [2]rune
	frostRunes  [2]rune
	unholyRunes [2]rune

	maxRunicPower     float64
	currentRunicPower float64
}
