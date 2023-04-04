package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) registerShadowBurnSpell() {
	if !warlock.Talents.Shadowburn {
		return
	}

	spellCoeff := 0.429 * (1 + 0.04*float64(warlock.Talents.ShadowAndFlame))

	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfShadowburn) {
		warlock.RegisterResetEffect(func(sim *core.Simulation) {
			sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int) {
				if isExecute == 35 {
					warlock.Shadowburn.BonusCritRating += 20 * core.CritRatingPerCritChance
				}
			})
		})
	}

	warlock.Shadowburn = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47827},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.2,
			Multiplier: 1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm],
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault, // backdraft procs don't change the GCD of shadowburn
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * time.Duration(15),
			},
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(775, 865) + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
