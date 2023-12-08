package warrior

import (
	"github.com/wowsims/sod/sim/core"
)

func (warrior *Warrior) registerDevastateSpell() {
	// if !warrior.Talents.Devastate {
	// 	return
	// }

	flatThreatBonus := 315.0
	dynaThreatBonus := 0.05

	weaponMulti := 1.2
	overallMulti := core.TernaryFloat64(warrior.HasSetBonus(ItemSetWrynnsPlate, 2), 1.05, 1.00)

	warrior.Devastate = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47498},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.CanApplySunderAura(target)
		},

		DamageMultiplier: overallMulti,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,
		FlatThreatBonus:  flatThreatBonus,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Bonus 242 damage / stack of sunder. Counts stacks AFTER cast but only if stacks > 0.
			sunderBonus := 0.0
			saStacks := warrior.SunderArmorAuras.Get(target).GetStacks()
			if saStacks != 0 {
				sunderBonus = 242 * float64(min(saStacks+1, 5))
			}

			baseDamage := (weaponMulti * spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())) + sunderBonus

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			result.Threat = spell.ThreatFromDamage(result.Outcome, result.Damage+dynaThreatBonus*spell.MeleeAttackPower())
			spell.DealDamage(sim, result)

			if result.Landed() {
				if warrior.CanApplySunderAura(target) {
					warrior.SunderArmorDevastate.Cast(sim, target)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},

		RelatedAuras: []core.AuraArray{warrior.SunderArmorAuras},
	})
}
