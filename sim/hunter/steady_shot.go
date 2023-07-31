package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (hunter *Hunter) registerSteadyShotSpell() {
	impSSProcChance := 0.05 * float64(hunter.Talents.ImprovedSteadyShot)
	if hunter.Talents.ImprovedSteadyShot > 0 {
		hunter.ImprovedSteadyShotAura = hunter.RegisterAura(core.Aura{
			Label:    "Improved Steady Shot",
			ActionID: core.ActionID{SpellID: 53220},
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				hunter.ArcaneShot.DamageMultiplierAdditive += .15
				hunter.ArcaneShot.CostMultiplier -= 0.2
				if hunter.AimedShot != nil {
					hunter.AimedShot.DamageMultiplierAdditive += .15
					hunter.AimedShot.CostMultiplier -= 0.2
				}
				if hunter.ChimeraShot != nil {
					hunter.ChimeraShot.DamageMultiplierAdditive += .15
					hunter.ChimeraShot.CostMultiplier -= 0.2
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				hunter.ArcaneShot.DamageMultiplierAdditive -= .15
				hunter.ArcaneShot.CostMultiplier += 0.2
				if hunter.AimedShot != nil {
					hunter.AimedShot.DamageMultiplierAdditive -= .15
					hunter.AimedShot.CostMultiplier += 0.2
				}
				if hunter.ChimeraShot != nil {
					hunter.ChimeraShot.DamageMultiplierAdditive -= .15
					hunter.ChimeraShot.CostMultiplier += 0.2
				}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == hunter.AimedShot || spell == hunter.ArcaneShot || spell == hunter.ChimeraShot {
					aura.Deactivate(sim)
				}
			},
		})
	}

	hunter.SteadyShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49052},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.05,
			Multiplier: 1 -
				0.03*float64(hunter.Talents.Efficiency) -
				0.05*float64(hunter.Talents.MasterMarksman),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 1, // Dummy value so core doesn't optimize the cast away
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.CastTime = hunter.SteadyShotCastTime()
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		BonusCritRating: 0 +
			2*core.CritRatingPerCritChance*float64(hunter.Talents.SurvivalInstincts),
		DamageMultiplierAdditive: 1 +
			.03*float64(hunter.Talents.FerociousInspiration) +
			core.TernaryFloat64(hunter.HasSetBonus(ItemSetGronnstalker, 4), .1, 0),
		DamageMultiplier: 1 *
			hunter.markedForDeathMultiplier(),
		CritMultiplier:   hunter.critMultiplier(true, true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.1*spell.RangedAttackPower(target) +
				hunter.AutoAttacks.Ranged.BaseDamage(sim)*2.8/hunter.AutoAttacks.Ranged.SwingSpeed +
				hunter.NormalizedAmmoDamageBonus +
				252

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			if result.Landed() && impSSProcChance > 0 && sim.RandomFloat("Imp Steady Shot") < impSSProcChance {
				hunter.ImprovedSteadyShotAura.Activate(sim)
			}
			spell.DealDamage(sim, result)
		},
	})
}

func (hunter *Hunter) SteadyShotCastTime() time.Duration {
	return time.Duration(float64(time.Millisecond*2000) / hunter.RangedSwingSpeed())
}
