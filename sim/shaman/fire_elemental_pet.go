package shaman

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Variables that control the Fire Elemental.
const (
	maxFireBlastCasts = 15
	maxFireNovaCasts  = 15
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
	fireElemental.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritRatingPerCritChance/212)
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
	fireBlastCasts := fireElemental.FireBlast.SpellMetrics[0].Casts
	fireNovaCasts := fireElemental.FireNova.SpellMetrics[0].Casts

	if fireBlastCasts == maxFireBlastCasts && fireNovaCasts == maxFireNovaCasts {
		fireElemental.CancelGCDTimer(sim)
		return
	}

	//Check for mana issues first
	if fireElemental.CurrentMana() < fireElemental.FireNova.CurCast.Cost {
		fireElemental.WaitForMana(sim, fireElemental.FireNova.CurCast.Cost)
		return
	}

	if fireBlastCasts < maxFireBlastCasts && fireElemental.FireBlast.IsReady(sim) {
		if fireElemental.FireBlast.Cast(sim, target) {
			return
		}
	}

	if fireNovaCasts < maxFireNovaCasts && fireElemental.FireNova.IsReady(sim) {
		if fireElemental.FireNova.Cast(sim, target) {
			return
		}
	}

	// Handle GCD down time.
	if !fireElemental.FireBlast.IsReady(sim) {
		fireElemental.WaitUntil(sim, fireElemental.FireBlast.CD.ReadyAt())
	} else if !fireElemental.FireNova.IsReady(sim) {
		fireElemental.WaitUntil(sim, fireElemental.FireNova.CD.ReadyAt())
	} else {
		fireElemental.WaitUntil(sim, fireElemental.AutoAttacks.NextAttackAt())
	}
}

var fireElementalPetBaseStats = stats.Stats{
	stats.Mana:        1789,
	stats.Health:      994,
	stats.Intellect:   147,
	stats.Stamina:     327,
	stats.SpellPower:  995,  //Estimated
	stats.AttackPower: 1369, //Estimated

	// TODO : Log digging and my own samples this seems to be around the 5% mark.
	stats.MeleeCrit: (5 + 1.8) * core.CritRatingPerCritChance,
	stats.SpellCrit: 2.61 * core.CritRatingPerCritChance,
}

func (shaman *Shaman) fireElementalStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ownerSpellHitChance := math.Floor(ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance)
		spellHitRatingFromOwner := ownerSpellHitChance * core.SpellHitRatingPerHitChance

		ownerHitChance := ownerStats[stats.MeleeHit] / core.MeleeHitRatingPerHitChance
		hitRatingFromOwner := math.Floor(ownerHitChance) * core.MeleeHitRatingPerHitChance

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:   ownerStats[stats.Intellect] * 0.30,
			stats.SpellPower:  ownerStats[stats.SpellPower] * 0.5218,
			stats.AttackPower: ownerStats[stats.SpellPower] * 4.45,

			// TODO tested useing pre-patch lvl 70 stats need to confirm in WOTLK at 80.
			stats.MeleeHit: hitRatingFromOwner,
			stats.SpellHit: spellHitRatingFromOwner,

			/*
				TODO these need to be confirmed borrowed from Hunter pets
			*/
			stats.Expertise: math.Floor((math.Floor(ownerHitChance/2) * PetExpertiseScale)) * core.ExpertisePerQuarterPercentReduction,
		}
	}
}
