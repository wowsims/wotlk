package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerImmolateSpell() {
	baseCost := 0.17 * warlock.BaseMana
	costReductionFactor := 1.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReductionFactor -= 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}
	spellCoefficient := 0.2
	actionID := core.ActionID{SpellID: 47811}
	spellSchool := core.SpellSchoolFire
	baseAdditiveMultiplier := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false)

	effect := core.SpellEffect{
		ProcMask: core.ProcMaskSpellDamage,

		DamageMultiplier: baseAdditiveMultiplier,

		BaseDamage:      core.BaseDamageConfigMagic(460.0, 460.0, spellCoefficient),
		OutcomeApplier:  warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
		OnSpellHitDealt: applyDotOnLanded(&warlock.ImmolateDot),
	}

	warlock.Immolate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * costReductionFactor,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (2000 - 100*time.Duration(warlock.Talents.Bane)),
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.GCD = time.Duration(float64(cast.GCD) * warlock.backdraftModifier())
				cast.CastTime = time.Duration(float64(cast.CastTime) * warlock.backdraftModifier())
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistFireCrit() +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

	target := warlock.CurrentTarget
	baseAdditiveMultiplierDot := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true)

	warlock.ImmolateDot = core.NewDot(core.Dot{
		Spell: warlock.Immolate,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Immolate-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 5 + int(warlock.Talents.MoltenCore),
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:   core.ProcMaskPeriodicDamage,
			IsPeriodic: true,

			DamageMultiplier: baseAdditiveMultiplierDot,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(785/5, spellCoefficient),

			OutcomeApplier: warlock.OutcomeFuncMagicCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
		}),
	})
}
