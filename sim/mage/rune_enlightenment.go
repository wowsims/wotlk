package mage

import (
	"github.com/wowsims/classic/sim/core"
)

func (mage *Mage) registerEnlightenment() {
	if !mage.HasRuneById(MageRuneChestEnlightenment) {
		return
	}

	// https://www.wowhead.com/classic/spell=412326/enlightenment
	enlightmentCritAura := mage.RegisterAura(core.Aura{
		Label:    "EnlightenmentCrit",
		ActionID: core.ActionID{SpellID: 412326},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
		},
	})

	// https://www.wowhead.com/classic/spell=412325/enlightenment
	enlightmentManaAura := mage.RegisterAura(core.Aura{
		Label:    "EnlightenmentMana",
		ActionID: core.ActionID{SpellID: 412325},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpiritRegenMultiplier *= 0.1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpiritRegenMultiplier /= 0.1
		},
	})

	mage.EnlightenmentAura = mage.RegisterAura(core.Aura{
		Label:    "Enlightenment",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
			enlightmentCritAura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			percentMana := aura.Unit.CurrentManaPercent()

			if percentMana > 0.7 && !enlightmentCritAura.IsActive() {
				enlightmentCritAura.Activate(sim)
			} else if percentMana <= 0.7 {
				enlightmentCritAura.Deactivate(sim)
			}

			if percentMana < 0.3 && !enlightmentManaAura.IsActive() {
				enlightmentManaAura.Activate(sim)
			} else if percentMana >= 0.3 {
				enlightmentManaAura.Deactivate(sim)
			}
		},
	})
}
