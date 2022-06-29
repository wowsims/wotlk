package common

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func NewSimpleStatItemActiveEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration, sharedCDFunc func(*core.Character) core.Cooldown) {
	core.NewItemEffect(itemID, core.MakeTemporaryStatsOnUseCDRegistration(
		"ItemActive-"+strconv.Itoa(int(itemID)),
		bonus,
		duration,
		core.SpellConfig{
			ActionID: core.ActionID{ItemID: itemID},
		},
		func(character *core.Character) core.Cooldown {
			return core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: cooldown,
			}
		},
		sharedCDFunc,
	))
}

// No shared CD
func NewSimpleStatItemEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration) {
	NewSimpleStatItemActiveEffect(itemID, bonus, duration, cooldown, func(character *core.Character) core.Cooldown {
		return core.Cooldown{}
	})
}

func NewSimpleStatOffensiveTrinketEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration) {
	NewSimpleStatItemActiveEffect(itemID, bonus, duration, cooldown, func(character *core.Character) core.Cooldown {
		return core.Cooldown{
			Timer:    character.GetOffensiveTrinketCD(),
			Duration: duration,
		}
	})
}

func NewSimpleStatDefensiveTrinketEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration) {
	NewSimpleStatItemActiveEffect(itemID, bonus, duration, cooldown, func(character *core.Character) core.Cooldown {
		return core.Cooldown{
			Timer:    character.GetDefensiveTrinketCD(),
			Duration: duration,
		}
	})
}
