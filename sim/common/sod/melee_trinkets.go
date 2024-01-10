package sod

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func init() {
	core.AddEffectsToTest = false

	// Proc effects. Keep these in order by item ID.

	//Hand of Justice
	core.NewItemEffect(11815, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.AutoAttacks.AutoSwingMelee {
			return
		}

		var handOfJusticeSpell *core.Spell
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 2,
		}
		procChance := 0.013333

		character.RegisterAura(core.Aura{
			Label:    "Hand of Justice",
			Duration: core.NeverExpires,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				config := *character.AutoAttacks.MHConfig()
				config.ActionID = core.ActionID{ItemID: 11815}
				handOfJusticeSpell = character.GetOrRegisterSpell(config)
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				// https://wotlk.wowhead.com/spell=15600/hand-of-justice, proc mask = 20.
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if sim.RandomFloat("HandOfJustice") > procChance {
					return
				}
				icd.Use(sim)

				aura.Unit.AutoAttacks.MaybeReplaceMHSwing(sim, handOfJusticeSpell).Cast(sim, result.Target)
			},
		})
	})

	core.AddEffectsToTest = true
}
