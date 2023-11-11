package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (warrior *Warrior) registerSlamSpell() {
	warrior.Slam = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47475},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   15 - float64(warrior.Talents.FocusedRage),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*1500 - time.Millisecond*500*time.Duration(warrior.Talents.ImprovedSlam),
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if cast.CastTime > 0 {
					warrior.AutoAttacks.DelayMeleeBy(sim, cast.CastTime)
				}
			},
		},

		BonusCritRating:  core.TernaryFloat64(warrior.HasSetBonus(ItemSetWrynnsBattlegear, 4), 5, 0) * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + 0.02*float64(warrior.Talents.UnendingFury) + core.TernaryFloat64(warrior.HasSetBonus(ItemSetDreadnaughtBattlegear, 2), 0.1, 0),
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,
		FlatThreatBonus:  140,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 250 +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (warrior *Warrior) ShouldInstantSlam(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.Slam.DefaultCast.Cost && warrior.Slam.IsReady(sim) && warrior.isBloodsurgeActive() &&
		sim.CurrentTime > (warrior.lastBloodsurgeProc+warrior.reactionTime) && warrior.GCD.IsReady(sim)
}

func (warrior *Warrior) ShouldSlam(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.Slam.DefaultCast.Cost && warrior.Slam.IsReady(sim) && warrior.Talents.ImprovedSlam > 0
}
