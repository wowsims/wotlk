package paladin

import (
	"github.com/wowsims/classic/sim/core"
)

func (paladin *Paladin) ActivateRighteousFury() {

	var holySpells []*core.Spell
	paladin.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellSchool == core.SpellSchoolHoly {
			holySpells = append(holySpells, spell)
		}
	})

	dtmMul := 1 - 0.02*float64(paladin.Talents.ImprovedRighteousFury)

	paladin.RighteousFuryAura = paladin.RegisterAura(core.Aura{
		Label:    "Righteous Fury",
		ActionID: core.ActionID{SpellID: 25780},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.ThreatMultiplier *= 1.43
			paladin.PseudoStats.DamageTakenMultiplier *= dtmMul
			for _, spell := range holySpells {
				spell.ThreatMultiplier *= 1.8
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.ThreatMultiplier /= 1.43
			paladin.PseudoStats.DamageTakenMultiplier /= dtmMul
			for _, spell := range holySpells {
				spell.ThreatMultiplier /= 1.8
			}
		},
	})
}
