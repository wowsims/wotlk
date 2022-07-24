package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerShadowBoltSpell() {

	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellCritRating: core.CritRatingPerCritChance * 5 * (core.TernaryFloat64(warlock.Talents.Devastation, 1, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDeathbringerGarb, 4), 1, 0) + core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 1, 0)),
		DamageMultiplier:     1,
		ThreatMultiplier:     1 - 0.1*float64(warlock.Talents.DestructiveReach),
		BaseDamage:           core.BaseDamageConfigMagic(694.0, 775.0, 0.857*(1+0.04*float64(warlock.Talents.ShadowAndFlame))),
		OutcomeApplier:       warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
		OnInit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			spellEffect.DamageMultiplier = warlock.spellDamageMultiplierHelper(sim, spell, spellEffect)
		},
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			// ISB debuff
			if warlock.Talents.ImprovedShadowBolt > 0 {
				if sim.RandomFloat("ISB") < 0.2*float64(warlock.Talents.ImprovedShadowBolt) {
					if !core.ShadowMasteryAura(warlock.CurrentTarget).IsActive() {
						core.ShadowMasteryAura(warlock.CurrentTarget).Activate(sim)
					} else {
						core.ShadowMasteryAura(warlock.CurrentTarget).Refresh(sim)
					}
				}
			}
			// Everlasting Affliction Refresh
			if warlock.CorruptionDot.IsActive() {
				if sim.RandomFloat("EverlastingAffliction") < 0.2*float64(warlock.Talents.EverlastingAffliction) {
					warlock.CorruptionDot.Refresh(sim)
				}
			}
		},
	}

	baseCost := 0.17 * warlock.BaseMana
	costReductionFactor := 1.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReductionFactor -= 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}
	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfShadowBolt) {
		costReductionFactor -= 0.1
	}

	warlock.ShadowBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47809},
		SpellSchool:  core.SpellSchoolShadow,
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