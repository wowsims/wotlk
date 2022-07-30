package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type ProcDamageEffect struct {
	ID      int32
	Trigger ProcTrigger

	School     core.SpellSchool
	BaseDamage core.BaseDamageConfig
}

func newProcDamageEffect(config ProcDamageEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: config.ID},
			SpellSchool: config.School,
			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage:       config.BaseDamage,
				OutcomeApplier:   character.OutcomeFuncMagicHitAndCrit(character.DefaultSpellCritMultiplier()),
			}),
		})

		triggerConfig := config.Trigger
		triggerConfig.Handler = func(sim *core.Simulation, _ *core.Spell, _ *core.SpellEffect) {
			damageSpell.Cast(sim, character.CurrentTarget)
		}
		MakeProcTriggerAura(&character.Unit, triggerConfig)
	})
}

func init() {
	core.AddEffectsToTest = false

	newProcDamageEffect(ProcDamageEffect{
		ID: 37064,
		Trigger: ProcTrigger{
			Name:       "Vestige of Haldor",
			Callback:   OnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 45,
		},
		School:     core.SpellSchoolFire,
		BaseDamage: core.BaseDamageConfigRoll(1024, 1536),
	})

	newProcDamageEffect(ProcDamageEffect{
		ID: 37264,
		Trigger: ProcTrigger{
			Name:       "Pendulum of Telluric Currents",
			Callback:   OnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 45,
		},
		School:     core.SpellSchoolShadow,
		BaseDamage: core.BaseDamageConfigRoll(1168, 1752),
	})

	newProcDamageEffect(ProcDamageEffect{
		ID: 39889,
		Trigger: ProcTrigger{
			Name:       "Horn of Agent Fury",
			Callback:   OnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 45,
		},
		School:     core.SpellSchoolHoly,
		BaseDamage: core.BaseDamageConfigRoll(1024, 1536),
	})

	core.AddEffectsToTest = true

	newProcDamageEffect(ProcDamageEffect{
		ID: 40371,
		Trigger: ProcTrigger{
			Name:       "Bandit's Insignia",
			Callback:   OnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 45,
		},
		School:     core.SpellSchoolArcane,
		BaseDamage: core.BaseDamageConfigRoll(1504, 2256),
	})

	newProcDamageEffect(ProcDamageEffect{
		ID: 40373,
		Trigger: ProcTrigger{
			Name:       "Extract of Necromantic Power",
			Callback:   OnPeriodicDamageDealt,
			Harmful:    true,
			ProcChance: 0.10,
			ICD:        time.Second * 15,
		},
		School:     core.SpellSchoolShadow,
		BaseDamage: core.BaseDamageConfigRoll(788, 1312),
	})

	newProcDamageEffect(ProcDamageEffect{
		ID: 42990,
		Trigger: ProcTrigger{
			Name:       "DMC Death",
			Callback:   OnSpellHitDealt | OnPeriodicDamageDealt,
			Harmful:    true,
			ProcChance: 0.15,
			ICD:        time.Second * 45,
		},
		School:     core.SpellSchoolShadow,
		BaseDamage: core.BaseDamageConfigRoll(1750, 2250),
	})
}
