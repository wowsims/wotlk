package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerSteadyShotSpell() {
	baseCost := 0.05 * hunter.BaseMana

	impSSProcChance := 0.05 * float64(hunter.Talents.ImprovedSteadyShot)
	if hunter.Talents.ImprovedSteadyShot > 0 {
		hunter.ImprovedSteadyShotAura = hunter.RegisterAura(core.Aura{
			Label:    "Improved Steady Shot",
			ActionID: core.ActionID{SpellID: 53220},
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				hunter.AimedShot.DamageMultiplier += .15
				hunter.AimedShot.CostMultiplier -= 0.2
				hunter.ArcaneShot.DamageMultiplier += .15
				hunter.ArcaneShot.CostMultiplier -= 0.2
				if hunter.ChimeraShot != nil {
					hunter.ChimeraShot.DamageMultiplier += .15
					hunter.ChimeraShot.CostMultiplier -= 0.2
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				hunter.AimedShot.DamageMultiplier -= .15
				hunter.AimedShot.CostMultiplier += 0.2
				hunter.ArcaneShot.DamageMultiplier -= .15
				hunter.ArcaneShot.CostMultiplier += 0.2
				if hunter.ChimeraShot != nil {
					hunter.ChimeraShot.DamageMultiplier -= .15
					hunter.ChimeraShot.CostMultiplier += 0.2
				}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell == hunter.AimedShot || spell == hunter.ArcaneShot || spell == hunter.ChimeraShot {
					aura.Deactivate(sim)
				}
			},
		})
	}

	hunter.SteadyShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49052},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.03*float64(hunter.Talents.Efficiency)) *
					(1 - 0.05*float64(hunter.Talents.MasterMarksman)),
				GCD:      core.GCDDefault,
				CastTime: 1, // Dummy value so core doesn't optimize the cast away
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.CastTime = hunter.SteadyShotCastTime()
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskRangedSpecial,

			BonusCritRating: 0 +
				2*core.CritRatingPerCritChance*float64(hunter.Talents.SurvivalInstincts) +
				core.TernaryFloat64(hunter.HasSetBonus(ItemSetRiftStalker, 4), 5*core.CritRatingPerCritChance, 0),

			DamageMultiplier: 1 *
				hunter.markedForDeathMultiplier(),

			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (hitEffect.RangedAttackPower(spell.Unit)+hitEffect.RangedAttackPowerOnTarget())*0.1 +
						hunter.AutoAttacks.Ranged.BaseDamage(sim)*2.8/hunter.AutoAttacks.Ranged.SwingSpeed +
						hunter.NormalizedAmmoDamageBonus +
						252
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, true, hunter.CurrentTarget)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() && impSSProcChance > 0 && sim.RandomFloat("Imp Steady Shot") < impSSProcChance {
					hunter.ImprovedSteadyShotAura.Activate(sim)
				}
			},
		}),

		InitialDamageMultiplier: 1 +
			.03*float64(hunter.Talents.FerociousInspiration) +
			core.TernaryFloat64(hunter.HasSetBonus(ItemSetGronnstalker, 4), .1, 0),
	})
}

func (hunter *Hunter) SteadyShotCastTime() time.Duration {
	return time.Duration(float64(time.Millisecond*2000) / hunter.RangedSwingSpeed())
}
