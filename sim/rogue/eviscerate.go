package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) makeEviscerate(comboPoints int32) *core.Spell {
	baseDamage := 127.0 + 370*float64(comboPoints)
	// tooltip implies 3..7% AP scaling, but testing show it's fixed at 7% (3.4.0.46158)
	apRatio := 0.07 * float64(comboPoints)

	cost := 35.0
	refundAmount := 0.4 * float64(rogue.Talents.QuickRecovery)

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48668, Tag: comboPoints},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | rogue.finisherFlags(),
		ResourceType: stats.Energy,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  time.Second,
			},
			ModifyCast:  rogue.CastModifier,
			IgnoreHaste: true,
		},

		BonusCritRating: core.TernaryFloat64(
			rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfEviscerate), 10*core.CritRatingPerCritChance, 0.0),
		DamageMultiplier: 1 +
			[]float64{0.0, 0.07, 0.14, 0.2}[rogue.Talents.ImprovedEviscerate] +
			0.02*float64(rogue.Talents.FindWeakness) +
			0.03*float64(rogue.Talents.Aggression),
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return baseDamage +
						254.0*sim.RandomFloat("Eviscerate") +
						apRatio*spell.MeleeAttackPower() +
						spell.BonusWeaponDamage()
				},
			},
			OutcomeApplier: rogue.OutcomeFuncMeleeSpecialHitAndCrit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					rogue.ApplyFinisher(sim, spell)
					rogue.ApplyCutToTheChase(sim)
				} else {
					if refundAmount > 0 {
						rogue.AddEnergy(sim, spell.CurCast.Cost*refundAmount, rogue.QuickRecoveryMetrics)
					}
				}
			},
		}),
	})
}

func (rogue *Rogue) registerEviscerate() {
	rogue.Eviscerate = [6]*core.Spell{
		nil,
		rogue.makeEviscerate(1),
		rogue.makeEviscerate(2),
		rogue.makeEviscerate(3),
		rogue.makeEviscerate(4),
		rogue.makeEviscerate(5),
	}
}
