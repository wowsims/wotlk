package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerChaosBoltSpell() {
	fireAndBrimstoneBonus := 0.02 * float64(warlock.Talents.FireAndBrimstone)
	actionID := core.ActionID{SpellID: 59172}
	spellSchool := core.SpellSchoolFire
	baseAdditiveMultiplier := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false)
	spellCoefficient := 0.7142 * (1 + 0.04*float64(warlock.Talents.ShadowAndFlame))

	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellCritRating: core.TernaryFloat64(warlock.Talents.Devastation, 1, 0) * 5 * core.CritRatingPerCritChance,
		DamageMultiplier:     1,
		ThreatMultiplier:     1 - 0.1*float64(warlock.Talents.DestructiveReach),
		BaseDamage:           core.BaseDamageConfigMagic(1429.0, 1813.0, spellCoefficient),
		OutcomeApplier:       warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
		OnInit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if warlock.ImmolateDot.IsActive() {
				spellEffect.DamageMultiplier = baseAdditiveMultiplier + fireAndBrimstoneBonus
			}
		},
	}

	baseCost := 0.07 * warlock.BaseMana
	costReductionFactor := 1.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReductionFactor -= 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}

	warlock.ChaosBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		Flags:        core.SpellFlagIgnoreResists,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * costReductionFactor,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*2500 - (time.Millisecond * 100 * time.Duration(warlock.Talents.Bane)),
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.GCD = time.Duration(float64(cast.GCD) * warlock.backdraftModifier())
				cast.CastTime = time.Duration(float64(cast.CastTime) * warlock.backdraftModifier())
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * (12 - 2*core.TernaryDuration(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfChaosBolt), 1, 0)),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

}
