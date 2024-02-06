package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var MutilateSpellID int32 = 399956

func (rogue *Rogue) newMutilateHitSpell(isMH bool) *core.Spell {
	actionID := core.ActionID{SpellID: 399960}
	procMask := core.ProcMaskMeleeMHSpecial
	if !isMH {
		actionID = core.ActionID{SpellID: 399961}
		procMask = core.ProcMaskMeleeOHSpecial
	}

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder | SpellFlagColdBlooded,

		BonusCritRating: 10 * core.CritRatingPerCritChance * float64(rogue.Talents.ImprovedBackstab),

		DamageMultiplierAdditive: 1 +
			0.04*float64(rogue.Talents.Opportunity),
		DamageMultiplier: 1 *
			core.TernaryFloat64(isMH, 1, rogue.dwsMultiplier()),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		// TODO: Update base damage with level scaling amount
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = 20 + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				baseDamage = 20 + spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			}
			// TODO: Add support for all poison effects
			if rogue.DeadlyPoison.Dot(target).IsActive() || rogue.woundPoisonDebuffAuras.Get(target).IsActive() {
				baseDamage *= 1.2
			}

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})
}

func (rogue *Rogue) registerMutilateSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneMutilate) {
		return
	}

	rogue.MutilateMH = rogue.newMutilateHitSpell(true)
	rogue.MutilateOH = rogue.newMutilateHitSpell(false)

	rogue.Mutilate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: MutilateSpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.costModifier(40),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit) // Miss/Dodge/Parry/Hit
			if result.Landed() {
				rogue.AddComboPoints(sim, 2, spell.ComboPointMetrics())
				rogue.MutilateMH.Cast(sim, target)
				rogue.MutilateOH.Cast(sim, target)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
