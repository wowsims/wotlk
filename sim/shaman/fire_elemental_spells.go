package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (fireElemental *FireElemental) registerFireBlast() {
	fireElemental.FireBlast = fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 13339},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			FlatCost: 276,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fireElemental.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO these are approximation, from base SP
			baseDamage := sim.Roll(323, 459) + 0.429*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (fireElemental *FireElemental) registerFireNova() {
	fireElemental.FireNova = fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 12470},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			FlatCost: 207,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second, // TODO estimated from log digging,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fireElemental.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO is this the right affect should it be Capped?
			// TODO these are approximation, from base SP
			dmgFromSP := 1.0071 * spell.SpellPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(1, 150) + dmgFromSP
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

func (fireElemental *FireElemental) registerFireShieldAura() {
	actionID := core.ActionID{SpellID: 11350}

	//dummy spell
	spell := fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplier: 1,
		CritMultiplier:   fireElemental.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "FireShield",
			},
			NumberOfTicks: 40,
			TickLength:    time.Second * 3,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// TODO is this the right affect should it be Capped?
				// TODO these are approximation, from base SP
				dmgFromSP := 0.032 * dot.Spell.SpellPower()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					baseDamage := sim.Roll(68, 70) + dmgFromSP
					//baseDamage *= sim.Encounter.AOECapMultiplier()
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeMagicCrit)
				}
			},
		},
	})

	fireElemental.FireShieldAura = fireElemental.RegisterAura(core.Aura{
		Label:    "Fire Shield",
		ActionID: actionID,
		Duration: time.Minute * 2,
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			spell.AOEDot().Apply(sim)
		},
	})
}
