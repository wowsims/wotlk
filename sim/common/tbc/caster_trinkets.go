package tbc

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false
	// Offensive trinkets. Keep these in order by item ID.
	core.NewSimpleStatOffensiveTrinketEffect(32483, stats.Stats{stats.SpellHaste: 175}, time.Second*20, time.Minute*2)  // Skull of Gul'dan
	core.NewSimpleStatOffensiveTrinketEffect(33829, stats.Stats{stats.SpellPower: 211}, time.Second*20, time.Minute*2)  // Hex Shrunken Head
	core.NewSimpleStatOffensiveTrinketEffect(34429, stats.Stats{stats.SpellPower: 320}, time.Second*15, time.Second*90) // Shifting Naaru Sliver

	// Even though these item effects are handled elsewhere, add them so they are
	// detected for automatic testing.
	for _, itemID := range core.AlchStoneItemIDs {
		core.NewItemEffect(itemID, func(core.Agent) {})
	}

	core.AddEffectsToTest = true
}
