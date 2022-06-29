package rogue

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

const RuptureEnergyCost = 25.0

func (rogue *Rogue) makeRupture(comboPoints int32) *core.Spell {
	refundAmount := 0.4 * float64(rogue.Talents.QuickRecovery)
	numTicks := int(comboPoints) + 3
	baseCost := RuptureEnergyCost

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26867, Tag: comboPoints},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists | rogue.finisherFlags(),

		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			ModifyCast:  rogue.applyDeathmantle,
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			OutcomeApplier:   rogue.OutcomeFuncMeleeSpecialHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					rogue.RuptureDot.Spell = spell
					rogue.RuptureDot.NumberOfTicks = numTicks
					rogue.RuptureDot.RecomputeAuraDuration()
					rogue.RuptureDot.Apply(sim)
					rogue.ApplyFinisher(sim, spell)
				} else {
					if refundAmount > 0 {
						rogue.AddEnergy(sim, spell.CurCast.Cost*refundAmount, rogue.QuickRecoveryMetrics)
					}
				}
			},
		}),
	})
}

func (rogue *Rogue) RuptureDuration(comboPoints int32) time.Duration {
	return time.Second*6 + time.Second*2*time.Duration(comboPoints)
}

func (rogue *Rogue) registerRupture() {
	rogue.Rupture = [6]*core.Spell{
		rogue.makeRupture(0), // Just for metrics
		rogue.makeRupture(1),
		rogue.makeRupture(2),
		rogue.makeRupture(3),
		rogue.makeRupture(4),
		rogue.makeRupture(5),
	}

	target := rogue.CurrentTarget
	rogue.RuptureDot = core.NewDot(core.Dot{
		Spell: rogue.Rupture[0],
		Aura: target.RegisterAura(core.Aura{
			Label:    "Rupture-" + strconv.Itoa(int(rogue.Index)),
			ActionID: rogue.Rupture[0].ActionID,
		}),
		NumberOfTicks: 0, // Set dynamically
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 + 0.1*float64(rogue.Talents.SerratedBlades),
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.BuildBaseDamageConfig(func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				comboPoints := rogue.ComboPoints()
				attackPower := hitEffect.MeleeAttackPower(spell.Unit) + hitEffect.MeleeAttackPowerOnTarget()

				return 70 + float64(comboPoints)*11 + attackPower*[]float64{0.01, 0.02, 0.03, 0.03, 0.03}[comboPoints-1]
			}, 0),
			OutcomeApplier: rogue.OutcomeFuncTick(),
		}),
	})
}
