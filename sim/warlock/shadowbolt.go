package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) registerShadowBoltSpell() {
	spellCoeff := 0.857 * (1 + 0.04*float64(warlock.Talents.ShadowAndFlame))
	ISBProcChance := 0.2 * float64(warlock.Talents.ImprovedShadowBolt)

	var shadowMasteryAuras core.AuraArray
	if ISBProcChance > 0 {
		shadowMasteryAuras = warlock.NewEnemyAuraArray(core.ShadowMasteryAura)
		warlock.CritDebuffCategory = shadowMasteryAuras.Get(warlock.CurrentTarget).ExclusiveEffects[0].Category
	}

	warlock.ShadowBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47809},
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 20,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.17,
			Multiplier: 1 -
				[]float64{0, .04, .07, .10}[warlock.Talents.Cataclysm] -
				core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfShadowBolt), 0.1, 0),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (3000 - 100*time.Duration(warlock.Talents.Bane)),
			},
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDeathbringerGarb, 4), 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery) +
			0.02*float64(warlock.Talents.ImprovedShadowBolt) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetMaleficRaiment, 4), 0.06, 0),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(694, 775) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					// ISB debuff
					if sim.Proc(ISBProcChance, "ISB") {
						shadowMasteryAuras.Get(target).Activate(sim)
					}
				}
			})
		},
	})
}
