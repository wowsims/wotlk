package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerImmolateSpell() {
	baseCost := 0.17 * warlock.BaseMana
	actionID := core.ActionID{SpellID: 47811}
	spellSchool := core.SpellSchoolFire

	warlock.Immolate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskSpellDamage,
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
		CritMultiplier:           warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier:         1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 460 + 0.2*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				warlock.ImmolateDot.Apply(sim)
			}
			spell.DealDamage(sim, result)
		},
	})

	target := warlock.CurrentTarget
	fireAndBrimstoneBonus := 0.02 * float64(warlock.Talents.FireAndBrimstone)

	warlock.ImmolateDot = core.NewDot(core.Dot{
		Spell: warlock.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: spellSchool,
			ProcMask:    core.ProcMaskSpellDamage,

			BonusCritRating:          warlock.Immolate.BonusCritRating,
			DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true),
			CritMultiplier:           warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
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
			IsPeriodic:     true,
			BaseDamage:     core.BaseDamageConfigMagicNoRoll(785/5, 0.2),
			OutcomeApplier: warlock.OutcomeFuncMagicCrit(),
		}),
	})
}
