package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerImmolateSpell() {
	baseCost := 0.17 * warlock.BaseMana
	spellCoefficient := 0.2
	actionID := core.ActionID{SpellID: 47811}
	spellSchool := core.SpellSchoolFire

	effect := core.SpellEffect{
		ProcMask: core.ProcMaskSpellDamage,

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
				Cost:     baseCost * (1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm]),
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
		DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false),
		ThreatMultiplier:         1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

	target := warlock.CurrentTarget
	fireAndBrimstoneBonus := 0.02 * float64(warlock.Talents.FireAndBrimstone)

	warlock.ImmolateDot = core.NewDot(core.Dot{
		Spell: warlock.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: spellSchool,

			BonusCritRating:          warlock.Immolate.BonusCritRating,
			DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true),
			ThreatMultiplier:         warlock.Immolate.ThreatMultiplier,
		}),
		Aura: target.RegisterAura(core.Aura{
			Label:    "Immolate-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warlock.ChaosBolt.DamageMultiplierAdditive += fireAndBrimstoneBonus
				warlock.Incinerate.DamageMultiplierAdditive += fireAndBrimstoneBonus
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warlock.ChaosBolt.DamageMultiplierAdditive -= fireAndBrimstoneBonus
				warlock.Incinerate.DamageMultiplierAdditive -= fireAndBrimstoneBonus
			},
		}),
		NumberOfTicks: 5 + int(warlock.Talents.MoltenCore),
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:       core.ProcMaskPeriodicDamage,
			IsPeriodic:     true,
			BaseDamage:     core.BaseDamageConfigMagicNoRoll(785/5, spellCoefficient),
			OutcomeApplier: warlock.OutcomeFuncMagicCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
		}),
	})
}
