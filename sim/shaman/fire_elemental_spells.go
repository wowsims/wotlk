package shaman

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (fireElemental *FireElemental) registerFireBlast() {
	var manaCost float64 = 276

	fireElemental.FireBlast = fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 13339},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: manaCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second,
			},
			OnCastComplete: func(sim *core.Simulation, _ *core.Spell) {
				fireElemental.AutoAttacks.DelayMeleeUntil(sim, sim.CurrentTime+fireElemental.AutoAttacks.MainhandSwingSpeed())
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
	var manaCost float64 = 207

	fireElemental.FireNova = fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 12470},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     manaCost,
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second, // TODO estimated from from log diggig,
			},
			ModifyCast: func(sim *core.Simulation, _ *core.Spell, _ *core.Cast) {
				fireElemental.AutoAttacks.DelayMeleeUntil(sim, sim.CurrentTime+fireElemental.AutoAttacks.MainhandSwingSpeed()*2)
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fireElemental.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO is this the right affect should it be Capped?
			// TODO these are approximation, from base SP
			dmgFromSP := 1.0071 * spell.SpellPower()
			for _, aoeTarget := range sim.Encounter.Targets {
				baseDamage := sim.Roll(1, 150) + dmgFromSP
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, spell.OutcomeMagicHitAndCrit)
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

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fireElemental.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
	})

	target := fireElemental.CurrentTarget

	fireShieldDot := core.NewDot(core.Dot{
		Spell: spell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "FireShield-" + strconv.Itoa(int(fireElemental.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 40,
		TickLength:    time.Second * 3,

		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			// TODO is this the right affect should it be Capped?
			// TODO these are approximation, from base SP
			dmgFromSP := 0.032 * dot.Spell.SpellPower()
			for _, aoeTarget := range sim.Encounter.Targets {
				baseDamage := sim.Roll(68, 70) + dmgFromSP
				//baseDamage *= sim.Encounter.AOECapMultiplier()
				dot.Spell.CalcAndDealDamageMagicCrit(sim, &aoeTarget.Unit, baseDamage)
			}
		},
	})

	fireElemental.FireShieldAura = fireElemental.RegisterAura(core.Aura{
		Label:    "Fire Shield",
		ActionID: actionID,
		Duration: time.Minute * 2,
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			fireShieldDot.Apply(sim)
		},
	})
}
