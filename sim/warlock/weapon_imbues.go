package warlock

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) applyWeaponImbue() {
	if warlock.Options.WeaponImbue == proto.Warlock_Options_GrandFirestone {
		warlock.AddStat(stats.SpellCrit, 49*(1+1.5*float64(warlock.Talents.MasterConjuror)))
		warlock.PseudoStats.DirectMagicDamageDealtMultiplier *= 1.01
	}
	if warlock.Options.WeaponImbue == proto.Warlock_Options_GrandSpellstone {
		warlock.AddStat(stats.SpellHaste, 60*(1+1.5*float64(warlock.Talents.MasterConjuror)))
		warlock.PseudoStats.PeriodicMagicDamageDealtMultiplier *= 1.01
	}
}
