package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerIncinerateSpell() {
	baseCost := 0.14 * warlock.BaseMana
	spellCoeff := 0.713 * (1 + 0.04*float64(warlock.Talents.ShadowAndFlame))
	mcCastMod := (1.0 - 0.1*float64(warlock.Talents.MoltenCore))

	warlock.Incinerate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47838},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 24,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm]),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * time.Duration(2500-50*warlock.Talents.Emberstorm),
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				totalMod := warlock.backdraftModifier()
				if warlock.MoltenCoreAura.IsActive() {
					totalMod *= mcCastMod
				}
				cast.GCD = time.Duration(float64(cast.GCD) * totalMod)
				cast.CastTime = time.Duration(float64(cast.CastTime) * totalMod)
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistFireCrit +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDeathbringerGarb, 4), 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.Emberstorm) +
			core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfIncinerate), 0.05, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetMaleficRaiment, 4), 0.06, 0),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if warlock.ImmolateDot.IsActive() {
				baseDamage = sim.Roll(582+145, 676+169) + spellCoeff*spell.SpellPower()
			} else {
				baseDamage = sim.Roll(582, 676) + spellCoeff*spell.SpellPower()
			}

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})

			if warlock.DemonicSoulAura.IsActive() {
				warlock.DemonicSoulAura.Deactivate(sim)
			}
			if warlock.MoltenCoreAura.IsActive() {
				warlock.MoltenCoreAura.RemoveStack(sim)
			}
		},
	})
}
