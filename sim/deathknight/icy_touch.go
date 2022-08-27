package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var IcyTouchActionID = core.ActionID{SpellID: 59131}

func (dk *Deathknight) registerIcyTouchSpell() {
	dk.FrostFeverDebuffAura = make([]*core.Aura, dk.Env.GetNumTargets())
	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit
		ffAura := core.FrostFeverAura(target, dk.Talents.ImprovedIcyTouch)
		ffAura.Duration = time.Second*15 + (time.Second * 3 * time.Duration(dk.Talents.Epidemic))
		dk.FrostFeverDebuffAura[target.Index] = ffAura
	}

	impIcyTouchCoeff := 1.0 + 0.05*float64(dk.Talents.ImprovedIcyTouch)
	sigilBonus := dk.sigilOfTheFrozenConscienceBonus()
	amountOfRunicPower := 10.0 + 2.5*float64(dk.Talents.ChillOfTheGrave)
	baseCost := float64(core.NewRuneCost(uint8(amountOfRunicPower), 0, 1, 0, 0))

	rs := &RuneSpell{}
	dk.IcyTouch = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID:     IcyTouchActionID,
		SpellSchool:  core.SpellSchoolFrost,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		ApplyEffects: dk.withRuneRefund(rs, core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: dk.rimeCritBonus() * core.CritRatingPerCritChance,
			DamageMultiplier:     impIcyTouchCoeff,
			ThreatMultiplier:     7.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					roll := (245.0-227.0)*sim.RandomFloat("Icy Touch") + 227.0 + sigilBonus
					return (roll + dk.getImpurityBonus(hitEffect, spell.Unit)*0.1) *
						dk.glacielRotBonus(hitEffect.Target) *
						dk.RoRTSBonus(hitEffect.Target) *
						dk.mercilessCombatBonus(sim)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.killingMachineOutcomeMod(dk.OutcomeFuncMagicHitAndCrit(dk.spellCritMultiplier())),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				dk.LastOutcome = spellEffect.Outcome
				if spellEffect.Landed() {
					if dk.KillingMachineAura.IsActive() {
						dk.KillingMachineAura.Deactivate(sim)
					}

					dk.FrostFeverExtended[spellEffect.Target.Index] = 0
					dk.FrostFeverSpell.Cast(sim, spellEffect.Target)
					if dk.Talents.CryptFever > 0 {
						dk.CryptFeverAura[spellEffect.Target.Index].Activate(sim)
					}
					if dk.Talents.EbonPlaguebringer > 0 {
						dk.EbonPlagueAura[spellEffect.Target.Index].Activate(sim)
					}
				}
			},
		}, false),
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0.0, 0, 1, 0) && dk.IcyTouch.IsReady(sim)
	}, nil)
}
func (dk *Deathknight) registerDrwIcyTouchSpell() {
	impIcyTouchCoeff := 1.0 + 0.05*float64(dk.Talents.ImprovedIcyTouch)
	sigilBonus := dk.sigilOfTheFrozenConscienceBonus()

	dk.RuneWeapon.IcyTouch = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    IcyTouchActionID,
		SpellSchool: core.SpellSchoolFrost,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: dk.rimeCritBonus() * core.CritRatingPerCritChance,
			DamageMultiplier:     impIcyTouchCoeff,
			ThreatMultiplier:     7.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					roll := (245.0-227.0)*sim.RandomFloat("Icy Touch") + 227.0 + sigilBonus
					return (roll + dk.RuneWeapon.getImpurityBonus(hitEffect, spell.Unit)*0.1)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.RuneWeapon.OutcomeFuncMagicHitAndCrit(dk.RuneWeapon.MeleeCritMultiplier(1.0, 0.0)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					dk.RuneWeapon.FrostFeverSpell.Cast(sim, spellEffect.Target)
				}
			},
		}),
	})
}
