package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerConflagrateSpell() {
	baseCost := 0.16 * warlock.BaseMana
	spellCoefficient := 0.2

	actionID := core.ActionID{SpellID: 17962}
	spellSchool := core.SpellSchoolFire
	target := warlock.CurrentTarget
	hasGlyphOfConflag := warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)

	effect := core.SpellEffect{
		BaseDamage:      core.BaseDamageConfigMagicNoRoll(0.6*785/5.*float64(warlock.ImmolateDot.NumberOfTicks), 0.6*spellCoefficient*float64(warlock.ImmolateDot.NumberOfTicks)),
		OutcomeApplier:  warlock.OutcomeFuncMagicHitAndCrit(),
		OnSpellHitDealt: applyDotOnLanded(&warlock.ConflagrateDot),
	}

	warlock.Conflagrate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm]),
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 10,
			},
			OnCastComplete: func(sim *core.Simulation, spell *core.Spell) {
				if !warlock.ImmolateDot.IsActive() {
					panic("Conflagrate spell is cast while Immolate is not active.")
				}
				if !hasGlyphOfConflag {
					warlock.ImmolateDot.Deactivate(sim)
					//warlock.ShadowflameDot.Deactivate(sim)
				}
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistFireCrit() +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0) +
			5*float64(warlock.Talents.FireAndBrimstone)*core.CritRatingPerCritChance,
		DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false),
		CritMultiplier:           warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier:         1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

	warlock.ConflagrateDot = core.NewDot(core.Dot{
		Spell: warlock.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: spellSchool,
			ProcMask:    core.ProcMaskSpellDamage,

			BonusCritRating:          warlock.Conflagrate.BonusCritRating,
			DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true),
			CritMultiplier:           warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
			ThreatMultiplier:         warlock.Conflagrate.ThreatMultiplier,
		}),
		Aura: target.RegisterAura(core.Aura{
			Label:    "conflagrate-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 3,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			IsPeriodic:     true,
			BaseDamage:     core.BaseDamageConfigMagicNoRoll(0.4/3*785/5*float64(warlock.ImmolateDot.NumberOfTicks), 0.4/3*spellCoefficient*float64(warlock.ImmolateDot.NumberOfTicks)),
			OutcomeApplier: warlock.OutcomeFuncMagicCrit(),
		}),
	})
}
