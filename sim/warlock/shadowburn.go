package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerShadowBurnSpell() {
	actionID := core.ActionID{SpellID: 47827}
	spellSchool := core.SpellSchoolShadow
	baseAdditiveMultiplier := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false)
	baseCost := 0.2 * warlock.BaseMana
	costReductionFactor := 1.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReductionFactor -= 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}

	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfShadowburn) {
		warlock.RegisterResetEffect(func(sim *core.Simulation) {
			sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int) {
				if isExecute == 35 {
					warlock.Shadowburn.BonusCritRating += 20 * core.CritRatingPerCritChance
				}
			})
		})
	}

	effect := core.SpellEffect{
		ProcMask: core.ProcMaskSpellDamage,

		BonusCritRating: 0 +
			warlock.masterDemonologistShadowCrit() +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0),
		DamageMultiplier: baseAdditiveMultiplier,
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		BaseDamage:     core.BaseDamageConfigMagic(775.0, 865.0, 0.429*(1+0.04*float64(warlock.Talents.ShadowAndFlame))),
		OutcomeApplier: warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
	}

	warlock.Shadowburn = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * costReductionFactor,
				GCD:  core.GCDDefault, // backdraft procs don't change the GCD of shadowburn
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * time.Duration(15),
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}
