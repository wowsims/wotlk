package rogue

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const RuptureEnergyCost = 25.0
const RuptureSpellID = 48672

func (rogue *Rogue) makeRupture(comboPoints int32) *core.Spell {
	refundAmount := 0.4 * float64(rogue.Talents.QuickRecovery)
	numTicks := int(comboPoints) + 3 + core.TernaryInt(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfRupture), 2, 0)
	baseCost := RuptureEnergyCost

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: RuptureSpellID, Tag: comboPoints},
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists | rogue.finisherFlags(),
		ResourceType: stats.Energy,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			ModifyCast:  rogue.CastModifier,
			IgnoreHaste: true,
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			OutcomeApplier:   rogue.OutcomeFuncMeleeSpecialHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					rogue.ruptureDot.Spell = spell
					rogue.ruptureDot.NumberOfTicks = numTicks
					rogue.ruptureDot.RecomputeAuraDuration()
					rogue.ruptureDot.Apply(sim)
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
	return time.Second*6 +
		time.Second*2*time.Duration(comboPoints) +
		core.TernaryDuration(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfRupture), time.Second*4, 0)
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
	rogue.ruptureDot = core.NewDot(core.Dot{
		Spell: rogue.Rupture[0],
		Aura: target.RegisterAura(core.Aura{
			Label:    "Rupture-" + strconv.Itoa(int(rogue.Index)),
			ActionID: rogue.Rupture[0].ActionID,
		}),
		NumberOfTicks: 0, // Set dynamically
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 +
				0.15*float64(rogue.Talents.BloodSpatter) +
				0.02*float64(rogue.Talents.FindWeakness) +
				core.TernaryFloat64(rogue.HasSetBonus(ItemSetBonescythe, 2), 0.1, 0) +
				core.TernaryFloat64(rogue.HasSetBonus(ItemSetTerrorblade, 4), 0.2, 0) +
				0.1*float64(rogue.Talents.SerratedBlades),
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.BuildBaseDamageConfig(func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				comboPoints := rogue.ComboPoints()
				attackPower := hitEffect.MeleeAttackPower(spell.Unit) + hitEffect.MeleeAttackPowerOnTarget()
				return 127 + float64(comboPoints)*18 + attackPower*[]float64{0.0, 0.015, 0.024, 0.03, 0.034286, 0.0375}[comboPoints]
			}, 0),
			OutcomeApplier: rogue.OutcomeFuncTickHitAndCrit(rogue.MeleeCritMultiplier(true, false)),
		}),
	})
}
