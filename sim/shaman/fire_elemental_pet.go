package shaman

import (
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

	fireElemental.EnableAutoAttacks(fireElemental, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  25, // TODO find out values
			BaseDamageMax:  27, // TODO find out values
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
	fireElemental.registerFireShieldAura()
}

func (fireElemental *FireElemental) Reset(sim *core.Simulation) {
}

func (fireElemental *FireElemental) enable(sim *core.Simulation) {
	// TODO Snapshot stats

	fireElemental.FireShieldDot.Apply(sim)
}

func (fireElemental *FireElemental) OnGCDReady(sim *core.Simulation) {
	target := fireElemental.CurrentTarget

	waitingOnCD := core.MinDuration(fireElemental.FireBlast.TimeToReady(sim), fireElemental.FireNova.TimeToReady(sim))

	/*
		TODO Need to handle the rotation, a 50/50 split might be close enough for now, but it does not use every gcd,
		will have to account for that.
	*/

	if fireElemental.FireBlast.IsReady(sim) {
		if !fireElemental.FireBlast.Cast(sim, target) {
			fireElemental.WaitForMana(sim, fireElemental.FireBlast.CurCast.Cost)
		}
		return
	}

	if fireElemental.FireNova.IsReady(sim) {
		if !fireElemental.FireNova.Cast(sim, target) {
			fireElemental.WaitForMana(sim, fireElemental.FireNova.CurCast.Cost)
		}
		return
	}

	fireElemental.WaitUntil(sim, sim.CurrentTime+waitingOnCD)
}

var fireElementalPetBaseStats = stats.Stats{
	//TODO
	stats.Mana: 3000,
}

func (shaman *Shaman) fireElementalStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			// TODO

		}
	}
}
