package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) registerAmbushSpell() {
	baseCost := rogue.costModifier(60 - 4*float64(rogue.Talents.SlaughterFromTheShadows))

	rogue.Ambush = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48691},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder,
		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(ItemSetVanCleefs, 4), 5*core.CritRatingPerCritChance, 0) +
			[]float64{0, 2, 4, 6}[rogue.Talents.TurnTheTables]*core.CritRatingPerCritChance +
			25*core.CritRatingPerCritChance*float64(rogue.Talents.ImprovedAmbush),
		// All of these use "Apply Aura: Modifies Damage/Healing Done", and stack additively.
		DamageMultiplier: 2.75 * (1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			0.1*float64(rogue.Talents.Opportunity)),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 330 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.AOEDot().Spell.AOEDot().Spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 2, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
