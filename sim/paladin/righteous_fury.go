package paladin

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (paladin *Paladin) ActivateRighteousFury() {
	paladin.PseudoStats.DamageTakenMultiplier *= 1 - 0.02*float64(paladin.Talents.ImprovedRighteousFury)

	paladin.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellSchool == core.SpellSchoolHoly {
			spell.ThreatMultiplier *= 1.8
		}
	})

	// Extra threat provided to all tanks on certain buff activation, for Paladins that is RF.
	paladin.PseudoStats.ThreatMultiplier *= 1.43
}
