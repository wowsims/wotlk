package warlock

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type WarlockPet struct {
	core.Pet

	config PetConfig

	owner *Warlock

	primaryAbility   *core.Spell
	secondaryAbility *core.Spell

	DemonicEmpowermentAura *core.Aura
}

func (warlock *Warlock) NewWarlockPet() *WarlockPet {

	summonChoice := warlock.Options.Summon
	preset := warlock.Rotation.Preset

	if preset == proto.Warlock_Rotation_Automatic {
		if warlock.Talents.Haunt {
			summonChoice = proto.Warlock_Options_Felhunter
		} else if warlock.Talents.Metamorphosis {
			summonChoice = proto.Warlock_Options_Felguard
		} else if warlock.Talents.ChaosBolt {
			summonChoice = proto.Warlock_Options_Imp
		}
	}

	petConfig := PetConfigs[summonChoice]

	wp := &WarlockPet{
		Pet: core.NewPet(
			petConfig.Name,
			&warlock.Character,
			petConfig.Stats,
			warlock.makeStatInheritance(),
			true,
			false,
		),
		config: petConfig,
		owner:  warlock,
	}

	// TODO: EnableManaBar should really be refactored to not assume 1 int = 15 mana
	wp.EnableManaBar()
	// the ratio multiplier affects the first 20 points as well
	wp.AddStat(stats.Mana, (15*20-20)+(20-15*20)*petConfig.ManaIntRatio/15)
	wp.AddStatDependency(stats.Intellect, stats.Mana, 1/(1+15))
	wp.AddStatDependency(stats.Intellect, stats.Mana, 1.0+petConfig.ManaIntRatio)
	wp.AddStatDependency(stats.Strength, stats.AttackPower, 1.0+2)
	wp.AddStatDependency(stats.Agility, stats.MeleeCrit, 1.0+core.CritRatingPerCritChance*0.04)
	wp.AddStats(stats.Stats{
		stats.MeleeCrit: float64(warlock.Talents.DemonicTactics) * 2 * core.CritRatingPerCritChance,
		stats.SpellCrit: float64(warlock.Talents.DemonicTactics) * 2 * core.CritRatingPerCritChance,
		stats.MeleeHit:  -float64(warlock.Talents.Suppression) * core.MeleeHitRatingPerHitChance, //Remove warlock's Suppression hit bonus from pet which he gets through stat inheritance
		stats.SpellHit:  -float64(warlock.Talents.Suppression) * core.SpellHitRatingPerHitChance, //Remove warlock's Suppression hit bonus from pet which he gets through stat inheritance
	})

	wp.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.04*float64(warlock.Talents.UnholyPower)

	if petConfig.Melee {
		switch summonChoice {
		case proto.Warlock_Options_Felguard:
			wp.EnableAutoAttacks(wp, core.AutoAttackOptions{
				MainHand: core.Weapon{
					BaseDamageMin:  412.5,
					BaseDamageMax:  412.5,
					SwingSpeed:     2,
					SwingDuration:  time.Second * 2,
					CritMultiplier: 2,
				},
				AutoSwingMelee: true,
			})
		case proto.Warlock_Options_Succubus:
			wp.EnableAutoAttacks(wp, core.AutoAttackOptions{
				MainHand: core.Weapon{
					BaseDamageMin:  412.5,
					BaseDamageMax:  412.5,
					SwingSpeed:     2,
					SwingDuration:  time.Second * 2,
					CritMultiplier: 2,
				},
				AutoSwingMelee: true,
			})
		case proto.Warlock_Options_Felhunter:
			wp.EnableAutoAttacks(wp, core.AutoAttackOptions{
				MainHand: core.Weapon{
					BaseDamageMin:  309.6,
					BaseDamageMax:  309.6,
					SwingSpeed:     2,
					SwingDuration:  time.Second * 2,
					CritMultiplier: 2,
				},
				AutoSwingMelee: true,
			})
		}
	}
	// wp.AutoAttacks.MHEffect.DamageMultiplier *= petConfig.DamageMultiplier
	switch summonChoice {
	case proto.Warlock_Options_Imp:
		// TODO: Does imp have different int->crit ratio than other casters? If so, we need to undo and then redo int->crit.
		// wp.AddStatDependency(stats.Intellect, stats.SpellCrit, 1.0+(0.0125*core.CritRatingPerCritChance/100))
		wp.PseudoStats.FireDamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
		wp.PseudoStats.BonusFireCritRating *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
	case proto.Warlock_Options_Succubus:
		wp.PseudoStats.ShadowDamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
		wp.PseudoStats.BonusShadowCritRating *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
	case proto.Warlock_Options_Felguard:
		wp.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
		// Simulates a pre-stacked demonic frenzy
		multiplier := 1.5 // demonic frenzy
		if wp.owner.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfFelguard) {
			multiplier *= 1.2
		}
		wp.AddStatDependency(stats.AttackPower, stats.AttackPower, multiplier)
	}

	if warlock.Talents.FelVitality > 0 {
		bonus := 1.0 + 0.05*float64(warlock.Talents.FelVitality)
		wp.AddStatDependency(stats.Intellect, stats.Intellect, bonus)
		wp.AddStatDependency(stats.Stamina, stats.Stamina, bonus)
	}

	if warlock.HasSetBonus(ItemSetOblivionRaiment, 2) {
		wp.AddStat(stats.MP5, 45)
	}

	core.ApplyPetConsumeEffects(&wp.Character, warlock.Consumes)

	warlock.AddPet(wp)

	return wp
}

func (wp *WarlockPet) GetPet() *core.Pet {
	return &wp.Pet
}

func (wp *WarlockPet) Initialize() {
	wp.primaryAbility = wp.NewPetAbility(wp.config.PrimaryAbility, true)
	wp.secondaryAbility = wp.NewPetAbility(wp.config.SecondaryAbility, false)
}

func (wp *WarlockPet) Reset(sim *core.Simulation) {

}

func (wp *WarlockPet) OnGCDReady(sim *core.Simulation) {
	target := wp.CurrentTarget
	if wp.config.RandomSelection {
		if sim.RandomFloat("Warlock Pet Ability") < 0.5 {
			if !wp.TryCast(sim, target, wp.primaryAbility) {
				wp.TryCast(sim, target, wp.secondaryAbility)
			}
		} else {
			if !wp.TryCast(sim, target, wp.secondaryAbility) {
				wp.TryCast(sim, target, wp.primaryAbility)
			}
		}
		return
	}

	if !wp.TryCast(sim, target, wp.primaryAbility) {
		if wp.secondaryAbility != nil {
			wp.TryCast(sim, target, wp.secondaryAbility)
		} else if !wp.primaryAbility.IsReady(sim) {
			wp.WaitUntil(sim, wp.primaryAbility.CD.ReadyAt())
		} else {
			wp.WaitForMana(sim, wp.primaryAbility.CurCast.Cost)
		}
	}
}

func (warlock *Warlock) makeStatInheritance() core.PetStatInheritance {
	improvedDemonicTactics := float64(warlock.Talents.ImprovedDemonicTactics)

	return func(ownerStats stats.Stats) stats.Stats {
		ownerHitChance := math.Floor(ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance)
		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
			stats.Armor:            ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:      (ownerStats[stats.SpellPower] + ownerStats[stats.ShadowSpellPower]) * 0.57,
			stats.SpellPower:       (ownerStats[stats.SpellPower] + ownerStats[stats.ShadowSpellPower]) * 0.15,
			stats.SpellPenetration: ownerStats[stats.SpellPenetration],
			stats.SpellCrit:        improvedDemonicTactics * 0.3 * ownerStats[stats.SpellCrit],
			stats.MeleeCrit:        improvedDemonicTactics * 0.3 * ownerStats[stats.SpellCrit],
			stats.MeleeHit:         ownerHitChance * core.MeleeHitRatingPerHitChance,
			stats.SpellHit:         ownerHitChance * core.SpellHitRatingPerHitChance,
			// Resists, 40%
		}
	}
}

type PetConfig struct {
	Name string
	// DamageMultiplier float64
	Melee        bool
	Stats        stats.Stats
	ManaIntRatio float64

	// Randomly select between abilities instead of using a prio.
	RandomSelection bool

	PrimaryAbility   PetAbilityType
	SecondaryAbility PetAbilityType
}

var PetConfigs = map[proto.Warlock_Options_Summon]PetConfig{
	proto.Warlock_Options_Felguard: {
		Name:             "Felguard",
		Melee:            true,
		PrimaryAbility:   Cleave,
		SecondaryAbility: Intercept,
		ManaIntRatio:     11.5,
		Stats: stats.Stats{
			stats.Stamina:   280,
			stats.Strength:  153,
			stats.Agility:   108,
			stats.Intellect: 133,
			stats.Mana:      893,
			stats.Spirit:    122,
			stats.MP5:       48,
		},
	},
	proto.Warlock_Options_Imp: {
		Name:           "Imp",
		ManaIntRatio:   4.9,
		Melee:          false,
		PrimaryAbility: Firebolt,
		Stats: stats.Stats{
			stats.MP5:       123,
			stats.Stamina:   101,
			stats.Strength:  145,
			stats.Agility:   38,
			stats.Intellect: 327,
			stats.Mana:      756,
			stats.Spirit:    263,
		},
	},
	proto.Warlock_Options_Succubus: {
		Name:           "Succubus",
		ManaIntRatio:   11.5,
		Melee:          true,
		PrimaryAbility: LashOfPain,
		Stats: stats.Stats{
			stats.Stamina:   328,
			stats.Strength:  314,
			stats.Agility:   90,
			stats.Intellect: 150,
			stats.Mana:      1109,
			stats.Spirit:    209,
			stats.MP5:       11,
		},
	},
	proto.Warlock_Options_Felhunter: {
		Name:           "Felhunter",
		ManaIntRatio:   15 * 0.77, // GetUnitPowerModifier("pet")
		Melee:          true,
		PrimaryAbility: ShadowBite,
		Stats: stats.Stats{
			stats.Strength:  314,
			stats.Agility:   90,
			stats.Stamina:   328,
			stats.Intellect: 150,
			stats.Spirit:    209,
			stats.Mana:      1559,
			stats.MP5:       11,
			stats.SpellCrit: 0.01,
			stats.MeleeCrit: 0.03,
		},
	},
}
