package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) SinisterStrikeEnergyCost() float64 {
	return []float64{45, 42, 40}[rogue.Talents.ImprovedSinisterStrike]
}

func (rogue *Rogue) registerSinisterStrikeSpell() {
	energyCost := rogue.SinisterStrikeEnergyCost()
	refundAmount := energyCost * 0.8

	rogue.SinisterStrike = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48638},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder,

		ResourceType: stats.Energy,
		BaseCost:     energyCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: energyCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
			ModifyCast:  rogue.CastModifier,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1 +
				0.03*float64(rogue.Talents.Aggression) +
				0.05*float64(rogue.Talents.BladeTwisting) +
				core.TernaryFloat64(rogue.Talents.SurpriseAttacks, 0.1, 0) +
				core.TernaryFloat64(rogue.HasSetBonus(ItemSetSlayers, 4), 0.06, 0),
			ThreatMultiplier: 1,
			BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(ItemSetVanCleefs, 4), 5*core.CritRatingPerCritChance, 0) +
				[]float64{0, 2, 4, 6}[rogue.Talents.TurnTheTables]*core.CritRatingPerCritChance,
			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 180, 1, 1, true),
			OutcomeApplier: rogue.OutcomeFuncMeleeSpecialHitAndCrit(rogue.MeleeCritMultiplier(true, true)),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					points := 1
					if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfSinisterStrike) {
						if sim.RandomFloat("Glyph of Sinister Strike") < 0.5 {
							points += 1
						}
					}
					rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				} else {
					rogue.AddEnergy(sim, refundAmount, rogue.EnergyRefundMetrics)
				}
			},
		}),
	})
}
