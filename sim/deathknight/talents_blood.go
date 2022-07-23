package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	//"time"

	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) ApplyBloodTalents() {
	// Butchery
	// Pointless to Implement - RaiN: Gives you passive 1 * rank runic power per 5 seconds so it needs to be implemented
	deathKnight.applyButchery()

	// Subversion
	// TODO: Implement

	// Blade barrier
	deathKnight.applyBladeBarrier()

	// Bladed Armor
	if deathKnight.Talents.BladedArmor > 0 {
		coeff := float64(deathKnight.Talents.BladedArmor)
		deathKnight.AddStatDependency(stats.Armor, stats.AttackPower, 1.0+coeff/180.0)
	}

	// Two Handed Specialization
	if deathKnight.Equip[proto.ItemSlot_ItemSlotMainHand].HandType == proto.HandType_HandTypeTwoHand {
		deathKnight.PseudoStats.PhysicalDamageDealtMultiplier *= 1 + 0.01*float64(deathKnight.Talents.TwoHandedWeaponSpecialization)
	}

	// Rune Tap
	// TODO: Implement

	// Dark Conviction
	deathKnight.PseudoStats.BonusMeleeCritRating += core.CritRatingPerCritChance * float64(deathKnight.Talents.DarkConviction)
	deathKnight.PseudoStats.BonusSpellCritRating += core.CritRatingPerCritChance * float64(deathKnight.Talents.DarkConviction)

	// Death Rune Mastery
	// TODO: Implement

	// Improved Rune Tap
	// TODO: Implement

	// Spell Deflection
	// TODO: Implement

	// Vendetta
	// TODO: Implement

	// Bloody Strikes
	// TODO: Implement

	// Veteran of the Third War
	if deathKnight.Talents.VeteranOfTheThirdWar > 0 {
		strengthCoeff := 0.02 * float64(deathKnight.Talents.VeteranOfTheThirdWar)
		staminaCoeff := 0.01 * float64(deathKnight.Talents.VeteranOfTheThirdWar)
		expertiseBonus := 2.0 * float64(deathKnight.Talents.VeteranOfTheThirdWar)
		deathKnight.AddStatDependency(stats.Strength, stats.Strength, 1.0+strengthCoeff)
		deathKnight.AddStatDependency(stats.Stamina, stats.Stamina, 1.0+staminaCoeff)
		deathKnight.AddStat(stats.Expertise, expertiseBonus*core.ExpertisePerQuarterPercentReduction)
	}

	// Mark of Blood
	// TODO: Implement

	// Bloody Vengeance
	// TODO: Implement

	// Abomination's Might
	if deathKnight.Talents.AbominationsMight > 0 {
		strengthCoeff := 0.01 * float64(deathKnight.Talents.AbominationsMight)
		deathKnight.AddStatDependency(stats.Strength, stats.Strength, 1.0+strengthCoeff)
	}
}

func (deathKnight *DeathKnight) subversionCritBonus() float64 {
	return 3.0 * float64(deathKnight.Talents.Subversion)
}

func (deathKnight *DeathKnight) improvedDeathStrikeCritBonus() float64 {
	return 3.0 * float64(deathKnight.Talents.ImprovedDeathStrike)
}

func (deathKnight *DeathKnight) applyBladeBarrier() {
	if deathKnight.Talents.BladeBarrier == 0 {
		return
	}

	damageTakenMult := 1.0 - 0.01*float64(deathKnight.Talents.BladeBarrier)

	actionID := core.ActionID{SpellID: 55226}

	deathKnight.BladeBarrierAura = deathKnight.RegisterAura(core.Aura{
		Label:    "Blade Barrier",
		ActionID: actionID,
		Duration: time.Second * 10.0,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTakenMult
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTakenMult
		},
	})
}

func (deathKnight *DeathKnight) applyButchery() {
	if deathKnight.Talents.Butchery == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 49483}

	rpMetrics := deathKnight.NewRunicPowerMetrics(actionID)

	deathKnight.ButcheryAura = core.MakePermanent(deathKnight.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Butchery",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			amountOfRunicPower := 1.0 * float64(deathKnight.Talents.Butchery)
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 5,
				NumTicks: 0,
				OnAction: func(sim *core.Simulation) {
					deathKnight.AddRunicPower(sim, amountOfRunicPower, rpMetrics)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
		},
	}))
}
