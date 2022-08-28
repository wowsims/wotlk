package shaman

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type RandomSpell int32

const (
	FireBlast = iota
	FireNova
	NoRandom
)

type FireElemental struct {
	core.Pet

	FireBlast     *core.Spell
	FireNova      *core.Spell
	FireShieldDot *core.Dot

	thinkChance     float64
	fireBlastChance float64

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
		shamanOwner:     shaman,
		thinkChance:     0.5,
		fireBlastChance: 0.5,
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

	if fireElemental.CurrentMana() < fireElemental.FireNova.CurCast.Cost {
		fireElemental.WaitForMana(sim, fireElemental.FireNova.CurCast.Cost)
		return
	}

	if fireElemental.tryThink(sim) {
		fireElemental.WaitUntil(sim, sim.CurrentTime+(time.Second*1))
		return
	}

	randomSpell := fireElemental.tryRandomSpellPicker(sim)

	if fireElemental.FireNova.IsReady(sim) && randomSpell != FireBlast {
		if fireElemental.FireNova.Cast(sim, target) {
			fireElemental.thinkChance = .95
			fireElemental.fireBlastChance = .95
			return
		}
	}

	if fireElemental.FireBlast.IsReady(sim) && randomSpell != FireNova {
		if fireElemental.FireBlast.Cast(sim, target) {
			fireElemental.thinkChance = .95
			fireElemental.fireBlastChance = 0.05
			return
		}
	}

	waitingOnCD := core.MinDuration(fireElemental.FireBlast.TimeToReady(sim), fireElemental.FireNova.TimeToReady(sim))
	if waitingOnCD == 0 {
		waitingOnCD = core.MaxDuration(fireElemental.FireBlast.TimeToReady(sim), fireElemental.FireNova.TimeToReady(sim))
	}

	fireElemental.WaitUntil(sim, sim.CurrentTime+waitingOnCD)
}

func (fireElemental *FireElemental) tryThink(sim *core.Simulation) bool {
	if sim.RandomFloat("Fire Ele Thinking") < fireElemental.thinkChance {
		fireElemental.thinkChance -= .15
		return true
	}

	return false
}

func (fireElemental *FireElemental) tryRandomSpellPicker(sim *core.Simulation) RandomSpell {

	if !fireElemental.FireBlast.IsReady(sim) || !fireElemental.FireNova.IsReady(sim) {
		return NoRandom
	}

	if sim.RandomFloat("Fire Ele RNG") < fireElemental.fireBlastChance {
		return FireBlast
	}

	return FireNova
}

var fireElementalPetBaseStats = stats.Stats{
	stats.Mana:        1789,
	stats.Health:      994,
	stats.Intellect:   147,
	stats.Stamina:     327,
	stats.SpellPower:  995,  //Estimated
	stats.AttackPower: 1369, //Estimated
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
