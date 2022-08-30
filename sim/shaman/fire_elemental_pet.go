package shaman

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Variables that control the Fire Elemental.
const (
	nFireBlastCasts = 15
	nFireNovaCasts  = 15
)

type FireElemental struct {
	core.Pet

	FireBlast *core.Spell
	FireNova  *core.Spell

	FireShieldAura *core.Aura

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
			BaseDamageMin:  1,  // Estimated from base AP
			BaseDamageMax:  24, // Estimated from base AP
			SwingSpeed:     2,
			SwingDuration:  time.Second * 2,
			CritMultiplier: 2, // Pretty sure this is right.
			SpellSchool:    core.SpellSchoolFire,
		},
		AutoSwingMelee: true,
	})

	fireElemental.OnPetEnable = fireElemental.enable
	fireElemental.OnPetDisable = fireElemental.disable

	shaman.AddPet(fireElemental)

	return fireElemental
}

func (fireElemental *FireElemental) enable(sim *core.Simulation) {
	fireElemental.FireShieldAura.Activate(sim)
}

func (fireElemental *FireElemental) disable(sim *core.Simulation) {
	fireElemental.FireShieldAura.Deactivate(sim)
}

func (fireElemental *FireElemental) GetPet() *core.Pet {
	return &fireElemental.Pet
}

func (fireElemental *FireElemental) Initialize() {
	fireElemental.registerFireBlast()
	fireElemental.registerFireNova()
	fireElemental.registerFireShieldAura()
}

func (fireElemental *FireElemental) Reset(sim *core.Simulation) {
}

func (fireElemental *FireElemental) OnGCDReady(sim *core.Simulation) {
	target := fireElemental.CurrentTarget

	//Check for mana issues first
	if fireElemental.CurrentMana() < fireElemental.FireNova.CurCast.Cost {
		fireElemental.WaitForMana(sim, fireElemental.FireNova.CurCast.Cost)
		return
	}

	//If no CD's are available on this GCD lets wait for the next spell off CD
	if !fireElemental.FireBlast.IsReady(sim) && !fireElemental.FireNova.IsReady(sim) {
		waitingOnCD := core.MinDuration(fireElemental.FireBlast.TimeToReady(sim), fireElemental.FireNova.TimeToReady(sim))
		fireElemental.WaitUntil(sim, sim.CurrentTime+waitingOnCD)
		return
	}

	numberCasts := fireElemental.FireBlast.SpellMetrics[0].Casts
	if numberCasts < nFireBlastCasts && fireElemental.FireBlast.IsReady(sim) {
		fireElemental.FireBlast.Cast(sim, target)
		return
	}

	numberCasts = fireElemental.FireNova.SpellMetrics[0].Casts
	if numberCasts < nFireNovaCasts && fireElemental.FireNova.IsReady(sim) {
		fireElemental.FireNova.Cast(sim, target)
		return
	}

}

var fireElementalPetBaseStats = stats.Stats{
	stats.Mana:        1789,
	stats.Health:      994,
	stats.Intellect:   147,
	stats.Stamina:     327,
	stats.SpellPower:  995,  //Estimated
	stats.AttackPower: 1369, //Estimated

	// TODO : No idea what his crit is at, he does not seem to gain any crit from owner.
	// Stole from spirit wolves.
	stats.MeleeCrit: (1.1515 + 1.8) * core.CritRatingPerCritChance,
}

func (shaman *Shaman) fireElementalStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ownerSpellHitChance := math.Floor(ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance)
		spellHitRatingFromOwner := math.Floor(ownerSpellHitChance) * core.SpellHitRatingPerHitChance

		ownerHitChance := ownerStats[stats.MeleeHit] / core.MeleeHitRatingPerHitChance
		hitRatingFromOwner := math.Floor(ownerHitChance) * core.MeleeHitRatingPerHitChance

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:   ownerStats[stats.Intellect] * 0.30,
			stats.SpellPower:  ownerStats[stats.SpellPower] * 0.5218,
			stats.AttackPower: ownerStats[stats.SpellPower] * 4.45,

			/*
				TODO these need to be confirmed borrowed from Hunter/Warlock pets.
			*/
			stats.MeleeHit:  hitRatingFromOwner / 2,
			stats.SpellHit:  spellHitRatingFromOwner,
			stats.Expertise: math.Floor((math.Floor(ownerHitChance/2) * PetExpertiseScale)) * core.ExpertisePerQuarterPercentReduction,
		}
	}
}
