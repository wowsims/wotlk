package common

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Keep these (and their functions) in alphabetical order.
func init() {
	// Proc effects. Keep these in order by item ID.

	core.NewItemEffect(28602, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Power of Arcanagos", core.ActionID{ItemID: 28602}, stats.Stats{stats.SpellPower: 130}, time.Second*10)

		// Gives a chance when your harmful spells land to increase the damage of your spells and effects by up to 130 for 10 sec. (Proc chance: 20%, 50s cooldown)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 50,
		}
		const proc = 0.2

		character.RegisterAura(core.Aura{
			Label:    "Robe of the Elder Scribes",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !spellEffect.Landed() {
					return
				}
				if !icd.IsReady(sim) || sim.RandomFloat("Robe of the Elder Scribe") > proc { // can't activate if on CD or didn't proc
					return
				}
				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(29305, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Band of the Eternal Sage Proc", core.ActionID{ItemID: 29305}, stats.Stats{stats.SpellPower: 95}, time.Second*10)

		// Your offensive spells have a chance on hit to increase your spell damage by 95 for 10 secs.
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 60,
		}
		const proc = 0.1

		character.RegisterAura(core.Aura{
			Label:    "Band of the Eternal Sage",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !spellEffect.Landed() {
					return
				}
				if !icd.IsReady(sim) || sim.RandomFloat("Band of the Eternal Sage") > proc { // can't activate if on CD or didn't proc
					return
				}
				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(34470, func(agent core.Agent) {
		character := agent.GetCharacter()

		timbalsSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 45055},
			SpellSchool: core.SpellSchoolShadow,
			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage:       core.BaseDamageConfigRoll(285, 475),
				OutcomeApplier:   character.OutcomeFuncMagicHitAndCrit(character.DefaultSpellCritMultiplier()),
			}),
		})

		// Each time one of your spells deals periodic damage,
		// there is a chance 285 to 475 additional damage will be dealt. (Proc chance: 10%, 15s cooldown)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 15,
		}
		const proc = 0.1

		character.RegisterAura(core.Aura{
			Label:    "Timbals",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !icd.IsReady(sim) || sim.RandomFloat("timbals") > proc { // can't activate if on CD or didn't proc
					return
				}
				icd.Use(sim)

				timbalsSpell.Cast(sim, spellEffect.Target)
			},
		})
	})

	core.NewItemEffect(34678, func(agent core.Agent) {
		character := agent.GetCharacter()
		const proc = 0.15

		var aldorAura *core.Aura
		var scryerSpell *core.Spell

		if character.ShattFaction == proto.ShattrathFaction_ShattrathFactionAldor {
			aldorAura = character.NewTemporaryStatsAura("Light's Wrath", core.ActionID{SpellID: 45479}, stats.Stats{stats.SpellPower: 120}, time.Second*10)
		} else if character.ShattFaction == proto.ShattrathFaction_ShattrathFactionScryer {
			scryerSpell = character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 45429},
				SpellSchool: core.SpellSchoolArcane,
				ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
					ProcMask:         core.ProcMaskEmpty,
					DamageMultiplier: 1,
					ThreatMultiplier: 1,
					BaseDamage:       core.BaseDamageConfigRoll(333, 367),
					OutcomeApplier:   character.OutcomeFuncMagicHitAndCrit(character.DefaultSpellCritMultiplier()),
				}),
			})
		}

		// Gives a chance when your harmful spells land to increase the damage of your spells and effects by up to 130 for 10 sec. (Proc chance: 20%, 50s cooldown)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}

		character.RegisterAura(core.Aura{
			Label:    "Shattered Sun Pendant of Acumen",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.ProcMask.Matches(core.ProcMaskSpellDamage) {
					return
				}
				if !spellEffect.Landed() {
					return
				}
				if !icd.IsReady(sim) || sim.RandomFloat("pendant of acumen") > proc { // can't activate if on CD or didn't proc
					return
				}
				icd.Use(sim)

				if character.ShattFaction == proto.ShattrathFaction_ShattrathFactionAldor {
					aldorAura.Activate(sim)
				} else if character.ShattFaction == proto.ShattrathFaction_ShattrathFactionScryer {
					scryerSpell.Cast(sim, spellEffect.Target)
				}
			},
		})
	})

}
