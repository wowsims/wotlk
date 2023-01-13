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

	sigilBonus := dk.sigilOfTheFrozenConscienceBonus()
	amountOfRunicPower := 10.0 + 2.5*float64(dk.Talents.ChillOfTheGrave)
	baseCost := float64(core.NewRuneCost(uint8(amountOfRunicPower), 0, 1, 0, 0))

	sigilOfTheUnfalteringKnight := dk.sigilOfTheUnfalteringKnight()

	rs := &RuneSpell{
		Refundable: true,
	}
	dk.IcyTouch = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID:     IcyTouchActionID,
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.GetModifiedGCD()
			},
		},

		BonusCritRating:  dk.rimeCritBonus() * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + 0.05*float64(dk.Talents.ImprovedIcyTouch),
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (sim.Roll(227, 245) + sigilBonus + 0.1*dk.getImpurityBonus(spell)) *
				dk.glacielRotBonus(target) *
				dk.RoRTSBonus(target) *
				dk.mercilessCombatBonus(sim)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			rs.OnResult(sim, result)

			dk.LastOutcome = result.Outcome
			if result.Landed() {
				dk.FrostFeverExtended[target.Index] = 0
				dk.FrostFeverSpell.Cast(sim, target)

				if sigilOfTheUnfalteringKnight != nil {
					sigilOfTheUnfalteringKnight.Activate(sim)
				}
			}

			spell.DealDamage(sim, result)
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0.0, 0, 1, 0) && dk.IcyTouch.IsReady(sim)
	})
}
func (dk *Deathknight) registerDrwIcyTouchSpell() {
	sigilBonus := dk.sigilOfTheFrozenConscienceBonus()

	dk.RuneWeapon.IcyTouch = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    IcyTouchActionID,
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagIgnoreAttackerModifiers,

		BonusCritRating:  dk.rimeCritBonus() * core.CritRatingPerCritChance,
		DamageMultiplier: 0.5 * (1 + 0.05*float64(dk.Talents.ImprovedIcyTouch)),
		CritMultiplier:   dk.RuneWeapon.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 7,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(227, 245) + sigilBonus + 0.1*dk.RuneWeapon.getImpurityBonus(spell)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				dk.RuneWeapon.FrostFeverSpell.Cast(sim, target)
			}
			spell.DealDamage(sim, result)
		},
	})
}
