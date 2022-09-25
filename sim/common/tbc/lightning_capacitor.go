package tbc

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func init() {
	core.AddEffectsToTest = false
	core.NewItemEffect(core.ItemIDTheLightningCapacitor, func(agent core.Agent) {
		character := agent.GetCharacter()

		tlcSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: core.ItemIDTheLightningCapacitor},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamageMagicHitAndCrit(sim, target, sim.Roll(694, 807))
			},
		})

		var charges int
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Millisecond * 2500,
		}

		character.RegisterAura(core.Aura{
			Label:    "Lightning Capacitor",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				charges = 0
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !icd.IsReady(sim) {
					return
				}

				if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}

				if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
					return
				}

				icd.Use(sim)
				charges++
				if charges >= 3 {
					tlcSpell.Cast(sim, spellEffect.Target)
					charges = 0
				}
			},
		})
	})
	core.AddEffectsToTest = true
}
