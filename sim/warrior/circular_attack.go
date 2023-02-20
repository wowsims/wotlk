package warrior

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (warrior *Warrior) registerCircularAttackSpell() {
	warrior.CircularAttack = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 319857},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RageCost: core.RageCostOptions{
			Cost:   0,
			Refund: 0,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
		},

		BonusCritRating:  core.TernaryFloat64(warrior.HasSetBonus(ItemSetSiegebreakerBattlegear, 4), 10, 0) * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 650 + 0.5*spell.MeleeAttackPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

func (warrior *Warrior) CanCircularAttack(sim *core.Simulation) bool {
	if warrior.HasActiveAura("Pouring out anger") {
		return true
	} else {
		return false
	}
}
