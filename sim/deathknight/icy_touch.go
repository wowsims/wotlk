package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerIcyTouchSpell() {
	dk.FrostFeverDebuffAura = make([]*core.Aura, dk.Env.GetNumTargets())
	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit
		ffAura := core.FrostFeverAura(target, dk.Talents.ImprovedIcyTouch)
		ffAura.Duration = time.Second*15 + (time.Second * 3 * time.Duration(dk.Talents.Epidemic))
		dk.FrostFeverDebuffAura[target.Index] = ffAura
	}

	impIcyTouchCoeff := 1.0 + 0.05*float64(dk.Talents.ImprovedIcyTouch)
	sigilBonus := +dk.sigilOfTheFrozenConscienceBonus()

	dk.IcyTouch = dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 59131},
		SpellSchool: core.SpellSchoolFrost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: dk.rimeCritBonus() * core.CritRatingPerCritChance,
			DamageMultiplier:     impIcyTouchCoeff,
			ThreatMultiplier:     7.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					roll := (245.0-227.0)*sim.RandomFloat("Icy Touch") + 227.0 + sigilBonus
					return (roll + dk.getImpurityBonus(hitEffect, spell.Unit)*0.1) *
						dk.glacielRotBonus(hitEffect.Target) *
						dk.rageOfRivendareBonus(hitEffect.Target) *
						dk.tundraStalkerBonus(hitEffect.Target) *
						dk.mercilessCombatBonus(sim)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.killingMachineOutcomeMod(dk.OutcomeFuncMagicHitAndCrit(dk.spellCritMultiplier())),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				dk.LastCastOutcome = spellEffect.Outcome
				if spellEffect.Landed() {
					if dk.KillingMachineAura.IsActive() {
						dk.KillingMachineAura.Deactivate(sim)
					}

					dk.FrostFeverSpell.Cast(sim, spellEffect.Target)
					if dk.Talents.CryptFever > 0 {
						dk.CryptFeverAura[spellEffect.Target.Index].Activate(sim)
					}
					if dk.Talents.EbonPlaguebringer > 0 {
						dk.EbonPlagueAura[spellEffect.Target.Index].Activate(sim)
					}

					dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_F)
					dk.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 10.0 + 2.5*float64(dk.Talents.ChillOfTheGrave)
					dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (dk *Deathknight) CanIcyTouch(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 0, 1, 0) && dk.IcyTouch.IsReady(sim)
}

func (dk *Deathknight) CastIcyTouch(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanIcyTouch(sim) {
		dk.IcyTouch.Cast(sim, target)
		return true
	}
	return false
}
