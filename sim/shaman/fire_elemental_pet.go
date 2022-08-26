package shaman

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type FireElemental struct {
	core.Pet

	FireBlast     *core.Spell
	FireNova      *core.Spell
	FireShieldDot *core.Dot

	shamanOwner *Shaman
}

func (shaman *Shaman) NewFireElemental() *FireElemental {
	fireElemental := &FireElemental{
		Pet: core.NewPet(
			"Greater Fire Elemental",
			&shaman.Character,
			fireElementalPetBaseStats,
			shaman.fireElementalStatInheritance(),
			false,
			true,
		),
		shamanOwner: shaman,
	}
	fireElemental.EnableManaBar()
	fireElemental.EnableAutoAttacks(fireElemental, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  1,  // TODO find out values
			BaseDamageMax:  24, // TODO find out values
			SwingSpeed:     2,
			SwingDuration:  time.Second * 2,
			CritMultiplier: 2, // Pretty sure this is right.
			SpellSchool:    core.SpellSchoolFire,
		},
		AutoSwingMelee: true,
	})

	fireElemental.OnPetEnable = fireElemental.enable
	shaman.AddPet(fireElemental)

	return fireElemental
}

func (fireElemental *FireElemental) GetPet() *core.Pet {
	return &fireElemental.Pet
}

func (fireElemental *FireElemental) Initialize() {
	fireElemental.registerFireBlast()
	fireElemental.registerFireNova()
	fireElemental.registerFireShieldDot()
}

func (fireElemental *FireElemental) Reset(sim *core.Simulation) {
}

func (fireElemental *FireElemental) enable(sim *core.Simulation) {
	fireElemental.FireShieldDot.Apply(sim)
}

func (fireElemental *FireElemental) OnGCDReady(sim *core.Simulation) {
	target := fireElemental.CurrentTarget

	/*
		TODO Need to handle the rotation, a 50/50 split might be close enough for now, but it does not use every gcd,
		will have to account for that.
	*/

	if fireElemental.FireNova.IsReady(sim) {
		if fireElemental.FireNova.Cast(sim, target) {
			return
		}
		fireElemental.WaitForMana(sim, fireElemental.FireNova.CurCast.Cost)
		return
	}

	if fireElemental.FireBlast.IsReady(sim) {
		if fireElemental.FireBlast.Cast(sim, target) {
			return
		}
	}

	waitingOnCD := core.MinDuration(fireElemental.FireBlast.TimeToReady(sim), fireElemental.FireNova.TimeToReady(sim))

	if waitingOnCD == 0 {
		waitingOnCD = core.MaxDuration(fireElemental.FireBlast.TimeToReady(sim), fireElemental.FireNova.TimeToReady(sim))
	}

	fireElemental.WaitUntil(sim, sim.CurrentTime+waitingOnCD)
}

var fireElementalPetBaseStats = stats.Stats{
	stats.Mana:        1789,
	stats.Health:      994,
	stats.Intellect:   147,
	stats.Stamina:     327,
	stats.SpellPower:  995,
	stats.AttackPower: 1369,
}

func (shaman *Shaman) fireElementalStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ownerHitChance := ownerStats[stats.MeleeHit] / core.MeleeHitRatingPerHitChance
		hitRatingFromOwner := math.Floor(ownerHitChance) * core.MeleeHitRatingPerHitChance

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:   ownerStats[stats.Intellect] * 0.30,
			stats.SpellPower:  ownerStats[stats.SpellPower] * 0.5218,
			stats.AttackPower: ownerStats[stats.SpellPower] * 4.45,

			// TODO these need to be confirmed borrowed from Hunter/Warlock pets
			stats.MeleeHit:  hitRatingFromOwner / 2,
			stats.SpellHit:  hitRatingFromOwner,
			stats.Expertise: math.Floor((math.Floor(ownerHitChance/2) * PetExpertiseScale)) * core.ExpertisePerQuarterPercentReduction,
		}
	}
}
