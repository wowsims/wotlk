package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type BloodwormPet struct {
	core.Pet

	dkOwner *Deathknight
}

func (dk *Deathknight) NewBloodwormPet(index int) *BloodwormPet {
	bloodworm := &BloodwormPet{
		Pet: core.NewPet(
			"Bloodworm", //+strconv.Itoa(index),
			&dk.Character,
			bloodwormPetBaseStats,
			dk.bloodwormStatInheritance(),
			false,
			true,
		),
		dkOwner: dk,
	}

	bloodworm.EnableAutoAttacks(bloodworm, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  25,
			BaseDamageMax:  27,
			SwingSpeed:     2,
			SwingDuration:  time.Second * 2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})

	// Hit and Crit only
	bloodworm.AutoAttacks.MHEffect.OutcomeApplier = bloodworm.OutcomeFuncMeleeSpecialCritOnly(bloodworm.MeleeCritMultiplier(1.0, 0.0))

	bloodworm.AddStatDependency(stats.Strength, stats.AttackPower, 1.0+1)
	bloodworm.AddStatDependency(stats.Agility, stats.MeleeCrit, 1.0+(core.CritRatingPerCritChance/83.3))

	bloodworm.OnPetEnable = bloodworm.enable
	bloodworm.OnPetDisable = bloodworm.disable

	dk.AddPet(bloodworm)

	return bloodworm
}

func (bloodworm *BloodwormPet) GetPet() *core.Pet {
	return &bloodworm.Pet
}

func (bloodworm *BloodwormPet) Initialize() {

}

func (bloodworm *BloodwormPet) Reset(sim *core.Simulation) {
}

func (bloodworm *BloodwormPet) OnGCDReady(sim *core.Simulation) {
}

func (bloodworm *BloodwormPet) enable(sim *core.Simulation) {
	// Snapshot extra % speed modifiers from dk owner
	bloodworm.PseudoStats.MeleeSpeedMultiplier = 1
	bloodworm.MultiplyMeleeSpeed(sim, bloodworm.dkOwner.PseudoStats.MeleeSpeedMultiplier)
}

func (bloodworm *BloodwormPet) disable(sim *core.Simulation) {
	// Clear snapshot speed
	bloodworm.PseudoStats.MeleeSpeedMultiplier = 1
	bloodworm.MultiplyMeleeSpeed(sim, 1)
}

var bloodwormPetBaseStats = stats.Stats{
	stats.MeleeCrit: 8 * core.CritRatingPerCritChance,
}

func (dk *Deathknight) bloodwormStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.AttackPower: ownerStats[stats.AttackPower] * 0.05,
			stats.MeleeHaste:  ownerStats[stats.MeleeHaste],
		}
	}
}
