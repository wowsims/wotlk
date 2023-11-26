package mage

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/stats"
)

func (mage *Mage) registerFingersOfFrost() {
	if !mage.HasRuneById(MageRuneChestFingersOfFrost) {
		return
	}

	bonusCrit := 0.1 * float64(mage.Talents.Shatter) * core.SpellCritRatingPerCritChance

	var proccedAt time.Duration

	mage.FingersOfFrostAura = mage.RegisterAura(core.Aura{
		Label:     "Fingers of Frost Proc",
		ActionID:  core.ActionID{SpellID: 400647},
		Duration:  time.Second * 15,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, -bonusCrit)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if proccedAt != sim.CurrentTime {
				aura.RemoveStack(sim)
			}
		},
	})

	procChance := 0.15
	mage.RegisterAura(core.Aura{
		Label:    "Fingers of Frost Rune",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Flags.Matches(SpellFlagChillSpell) && sim.RandomFloat("Fingers of Frost") < procChance {
				mage.FingersOfFrostAura.Activate(sim)
				mage.FingersOfFrostAura.SetStacks(sim, 2)
				proccedAt = sim.CurrentTime
			}
		},
	})
}
