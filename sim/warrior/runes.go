package warrior

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

func (warrior *Warrior) ApplyRunes() {
	warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= core.TernaryFloat64(warrior.HasRune(proto.WarriorRune_RuneSingleMindedFury), 1.1, 1)
	if warrior.GetMHWeapon().HandType == proto.HandType_HandTypeTwoHand {
		warrior.PseudoStats.MeleeSpeedMultiplier *= core.TernaryFloat64(warrior.HasRune(proto.WarriorRune_RuneFrenziedAssault), 1.1, 1)
	}

	warrior.applyBloodFrenzy()
	warrior.applyFlagellation()
	warrior.applyFlagellation()
	warrior.applyRagingBlow()
	warrior.applyConsumedByRage()
	warrior.applyQuickStrike()

	// Endless Rage implemented on dps_warrior.go and protection_warrior.go

}

func (warrior *Warrior) applyBloodFrenzy() {
	if !warrior.HasRune(proto.WarriorRune_RuneBloodFrenzy) {
		return
	}

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 412507})

	warrior.RegisterAura(core.Aura{
		Label:    "Blood Frenzy",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagBleed) {
				return
			}

			warrior.AddRage(sim, 3, rageMetrics)
		},
	})
}

func (warrior *Warrior) applyFlagellation() {
	if !warrior.HasRune(proto.WarriorRune_RuneFlagellation) {
		return
	}

	flagellationAura := warrior.RegisterAura(core.Aura{
		Label:    "Flagellation Trigger",
		ActionID: core.ActionID{SpellID: 402877},
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.25
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.25
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Flagellation Trigger",
		Duration: core.NeverExpires,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !(spell == warrior.Bloodrage || spell == warrior.BerserkerRage) {
				return
			}

			flagellationAura.Activate(sim)
		},
	})
}

func (warrior *Warrior) applyRagingBlow() {
	if !warrior.HasRune(proto.WarriorRune_RuneRagingBlow) {
		return
	}

	warrior.RagingBlow = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 402911},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 8,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.ConsumedByRageAura.IsActive() || warrior.BloodrageAura.IsActive() || warrior.BerserkerRageAura.IsActive()
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

		},
	})
}

func (warrior *Warrior) applyConsumedByRage() {
	if !warrior.HasRune(proto.WarriorRune_RuneConsumedByRage) {
		return
	}

	warrior.ConsumedByRageAura = warrior.RegisterAura(core.Aura{
		Label:     "Consumed By Rage",
		ActionID:  core.ActionID{SpellID: 425418},
		MaxStacks: 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.2
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Consumed By Rage Trigger",
		Duration: core.NeverExpires,
		OnRageChange: func(aura *core.Aura, sim *core.Simulation) {
			if !warrior.Above80RageCBRActive && warrior.CurrentRage() >= 80 {
				warrior.ConsumedByRageAura.Activate(sim)
				warrior.ConsumedByRageAura.SetStacks(sim, 12)
				warrior.Above80RageCBRActive = true
			} else if warrior.Above80RageCBRActive && warrior.CurrentRage() < 80 {
				warrior.Above80RageCBRActive = false
			}

		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !warrior.ConsumedByRageAura.IsActive() {
				return
			}

			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				warrior.ConsumedByRageAura.RemoveStack(sim)
			}
		},
	})
}

func (warrior *Warrior) applyQuickStrike() {
	if !warrior.HasRune(proto.WarriorRune_RuneQuickStrike) {
		return
	}

	warrior.QuickStrike = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 429765},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   15 - float64(warrior.Talents.ImprovedHeroicStrike),
			Refund: 0.8,
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,
		FlatThreatBonus:  259,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(0.15*spell.MeleeAttackPower(), 0.25*spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}
