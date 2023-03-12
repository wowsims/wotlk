package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var BloodBoilActionID = core.ActionID{SpellID: 49941}

func (dk *Deathknight) registerBloodBoilSpell() {
	// TODO: Handle blood boil correctly -
	//  There is no refund and you only get RP on at least one of the effects hitting.
	dk.BloodBoil = dk.RegisterSpell(core.SpellConfig{
		ActionID:    BloodBoilActionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: dk.bloodyStrikesBonus(dk.BloodBoil),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := (sim.Roll(180, 220) + 0.06*dk.getImpurityBonus(spell)) * dk.RoRTSBonus(aoeTarget) * core.TernaryFloat64(dk.DiseasesAreActive(aoeTarget), 1.5, 1.0)
				baseDamage *= sim.Encounter.AOECapMultiplier()

				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

				if aoeTarget == target {
					spell.SpendRefundableCost(sim, result)
					dk.LastOutcome = result.Outcome
				}
			}
		},
	})
}

func (dk *Deathknight) registerDrwBloodBoilSpell() {
	dk.RuneWeapon.BloodBoil = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    BloodBoilActionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagIgnoreAttackerModifiers,

		DamageMultiplier: 0.5 * dk.bloodyStrikesBonus(dk.BloodBoil),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := (sim.Roll(180, 220) + 0.06*dk.RuneWeapon.getImpurityBonus(spell)) * core.TernaryFloat64(dk.DrwDiseasesAreActive(aoeTarget), 1.5, 1.0)
				baseDamage *= sim.Encounter.AOECapMultiplier()

				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
