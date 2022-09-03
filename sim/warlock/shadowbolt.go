package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerShadowBoltSpell() {
	ISBProcChance := 0.2 * float64(warlock.Talents.ImprovedShadowBolt)
	actionID := core.ActionID{SpellID: 47809}
	spellSchool := core.SpellSchoolShadow
	baseAdditiveMultiplier := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false)
	baseCost := 0.17 * warlock.BaseMana
	costReductionFactor := 1.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReductionFactor -= 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}
	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfShadowBolt) {
		costReductionFactor -= 0.1
	}

	effect := core.SpellEffect{
		ProcMask: core.ProcMaskSpellDamage,
		BonusCritRating: 0 +
			warlock.masterDemonologistShadowCrit() +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDeathbringerGarb, 4), 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplier: baseAdditiveMultiplier,
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),
		MissileSpeed:     20,
		BaseDamage:       core.BaseDamageConfigMagic(694.0, 775.0, 0.857*(1+0.04*float64(warlock.Talents.ShadowAndFlame))),
		OutcomeApplier:   warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			// ISB debuff
			if warlock.Talents.ImprovedShadowBolt > 0 {
				if warlock.Talents.ImprovedShadowBolt < 5 { // This will return early if we 'miss' the refresh, 5 pts can't 'miss'.
					if sim.RandomFloat("ISB") > ISBProcChance {
						return
					}
				}
				core.ShadowMasteryAura(warlock.CurrentTarget).Activate(sim) // calls refresh if already active.
			}
		},
	}

	warlock.ShadowBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * costReductionFactor,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (3000 - 100*time.Duration(warlock.Talents.Bane)),
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.GCD = time.Duration(float64(cast.GCD) * warlock.backdraftModifier())
				cast.CastTime = time.Duration(float64(cast.CastTime) * warlock.backdraftModifier())
				if warlock.Talents.Nightfall > 0 {
					warlock.applyNightfall(cast)
				}
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}
