package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	//"time"

	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) ApplyBloodTalents() {
	// Butchery
	// Pointless to Implement - RaiN: Gives you passive 1 * rank runic power per 5 seconds so it needs to be implemented
	dk.applyButchery()

	// Subversion
	// TODO: Implement

	// Blade barrier
	dk.applyBladeBarrier()

	// Bladed Armor
	if dk.Talents.BladedArmor > 0 {
		coeff := float64(dk.Talents.BladedArmor)
		dk.AddStatDependency2(stats.Armor, stats.AttackPower, coeff/180.0)
	}

	// Two Handed Specialization
	if dk.Equip[proto.ItemSlot_ItemSlotMainHand].HandType == proto.HandType_HandTypeTwoHand {
		dk.PseudoStats.PhysicalDamageDealtMultiplier *= 1 + 0.01*float64(dk.Talents.TwoHandedWeaponSpecialization)
	}

	// Rune Tap
	// TODO: Implement

	// Dark Conviction
	dk.PseudoStats.BonusMeleeCritRating += core.CritRatingPerCritChance * float64(dk.Talents.DarkConviction)
	dk.PseudoStats.BonusSpellCritRating += core.CritRatingPerCritChance * float64(dk.Talents.DarkConviction)

	// Death Rune Mastery
	// TODO: Implement

	// Improved Rune Tap
	// TODO: Implement

	// Spell Deflection
	//dk.applySpellDeflection()

	// Vendetta
	// TODO: Implement

	// Bloody Strikes
	// TODO: Implement

	// Veteran of the Third War
	if dk.Talents.VeteranOfTheThirdWar > 0 {
		strengthCoeff := 0.02 * float64(dk.Talents.VeteranOfTheThirdWar)
		staminaCoeff := 0.01 * float64(dk.Talents.VeteranOfTheThirdWar)
		expertiseBonus := 2.0 * float64(dk.Talents.VeteranOfTheThirdWar)
		dk.MultiplyStat(stats.Strength, 1.0+strengthCoeff)
		dk.MultiplyStat(stats.Stamina, 1.0+staminaCoeff)
		dk.AddStat(stats.Expertise, expertiseBonus*core.ExpertisePerQuarterPercentReduction)
	}

	// Mark of Blood
	// TODO: Implement

	// Bloody Vengeance
	// TODO: Implement

	// Abomination's Might
	if dk.Talents.AbominationsMight > 0 {
		strengthCoeff := 0.01 * float64(dk.Talents.AbominationsMight)
		dk.MultiplyStat(stats.Strength, 1.0+strengthCoeff)
	}
}

func (dk *Deathknight) subversionThreatBonus() float64 {
	return []float64{0.0, 0.08, 0.16, 0.25}[dk.Talents.Subversion]
}

func (dk *Deathknight) subversionCritBonus() float64 {
	return 3.0 * float64(dk.Talents.Subversion)
}

func (dk *Deathknight) improvedDeathStrikeCritBonus() float64 {
	return 3.0 * float64(dk.Talents.ImprovedDeathStrike)
}

func (dk *Deathknight) applyBladeBarrier() {
	if dk.Talents.BladeBarrier == 0 {
		return
	}

	damageTakenMult := 1.0 - 0.01*float64(dk.Talents.BladeBarrier)

	actionID := core.ActionID{SpellID: 55226}

	dk.BladeBarrierAura = dk.RegisterAura(core.Aura{
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

	dk.onRuneSpendBladeBarrier = func(sim *core.Simulation) {
		if dk.AllBloodRunesSpent() {
			dk.BladeBarrierAura.Activate(sim)
		}
	}
}

/*
func (dk *Deathknight) applySpellDeflection() {
	if dk.Talents.SpellDeflection == 0 {
		return
	}

	damageTakenMult := 0.15 * float64(dk.Talents.SpellDeflection)

	actionID := core.ActionID{SpellID: 55226}

	dk.BladeBarrierAura = dk.RegisterAura(core.Aura{
		Label:    "Spell Deflection",
		ActionID: actionID,
		Duration: core.NeverExpires,
	})

	dk.onRuneSpendBladeBarrier = func(sim *core.Simulation) {
		if dk.AllBloodRunesSpent() {
			dk.BladeBarrierAura.Activate(sim)
		}
	}
}
*/

func (dk *Deathknight) applyButchery() {
	if dk.Talents.Butchery == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 49483}
	amountOfRunicPower := 1.0 * float64(dk.Talents.Butchery)
	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	dk.ButcheryAura = core.MakePermanent(dk.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Butchery",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.ButcheryPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 5,
				NumTicks: 0,
				OnAction: func(sim *core.Simulation) {
					dk.AddRunicPower(sim, amountOfRunicPower, rpMetrics)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
		},
	}))
}
