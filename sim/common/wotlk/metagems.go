package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	// Keep these in order by item ID.

	// TODO: Destructive Skyflare (1% spell reflect)

	core.NewItemEffect(41333, func(agent core.Agent) {
		agent.GetCharacter().MultiplyStat(stats.Intellect, 1.02)
	})

	core.NewItemEffect(41377, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 0.98
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 0.98
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 0.98
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= 0.98
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 0.98
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 0.98
	})

	core.NewItemEffect(41380, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.ApplyEquipScaling(stats.Armor, 1.02)
	})

	core.NewItemEffect(41385, func(agent core.Agent) {
		character := agent.GetCharacter()

		healSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 55341},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskSpellHealing,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultHealingCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseHealing := 0.02 * target.MaxHealth()
				spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Invigorating Earthsiege Diamond",
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
			Outcome:    core.OutcomeCrit,
			ProcChance: 0.5,
			ICD:        time.Second * 45,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				healSpell.Cast(sim, &character.Unit)
			},
		})
	})

	core.NewItemEffect(41389, func(agent core.Agent) {
		agent.GetCharacter().MultiplyStat(stats.Mana, 1.02)
	})

	core.NewItemEffect(41395, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 0.98
	})

	// Eternal Earthsiege
	core.NewItemEffect(41396, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.BlockValueMultiplier += 0.05
	})

	core.NewItemEffect(41400, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Thundering Skyflare Diamond Proc", core.ActionID{SpellID: 55379}, stats.Stats{stats.MeleeHaste: 480}, time.Second*6)

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Thundering Skyflare Diamond",
			Callback: core.CallbackOnSpellHitDealt,
			// Mask 68, melee or ranged auto attacks.
			ProcMask: core.ProcMaskWhiteHit,
			Outcome:  core.OutcomeLanded,
			PPM:      1,
			ICD:      time.Second * 40,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
		procAura.Icd = triggerAura.Icd
	})

	core.NewItemEffect(41401, func(agent core.Agent) {
		character := agent.GetCharacter()
		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 55382})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Insightful Earthsiege Diamond",
			Callback:   core.CallbackOnCastComplete,
			ProcChance: 0.05,
			ICD:        time.Second * 15,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				character.AddMana(sim, 600, manaMetrics)
			},
			ProcMask: ^core.ProcMaskProc & ^core.ProcMaskWeaponProc,
		})
	})

	// These are handled in character.go, but create empty effects so they are included in tests.
	core.NewItemEffect(41285, func(_ core.Agent) {}) // Chaotic Skyflare Diamond
	core.NewItemEffect(41376, func(_ core.Agent) {}) // Revitalizing Skyflare Diamond
	core.NewItemEffect(41398, func(_ core.Agent) {}) // Relentless Earthsiege Diamond
}
