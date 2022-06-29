package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

type WarlockPet struct {
	core.Pet

	config PetConfig

	owner *Warlock

	primaryAbility   *core.Spell
	secondaryAbility *core.Spell
}

func (warlock *Warlock) NewWarlockPet() *WarlockPet {
	// if warlock.Options.PetUptime <= 0 {
	// 	return nil
	// }
	petConfig := PetConfigs[warlock.Options.Summon]

	wp := &WarlockPet{
		Pet: core.NewPet(
			petConfig.Name,
			&warlock.Character,
			petConfig.Stats,
			petStatInheritance,
			true,
		),
		config: petConfig,
		owner:  warlock,
	}
	wp.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Intellect,
		ModifiedStat: stats.Mana,
		Modifier: func(intellect float64, mana float64) float64 {
			return mana + intellect*petConfig.ManaIntRatio
		},
	})
	wp.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.AttackPower,
		Modifier: func(strength float64, attackPower float64) float64 {
			return attackPower + (strength-10)*2
		},
	})
	wp.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.MeleeCrit,
		Modifier: func(agility float64, meleeCrit float64) float64 {
			return meleeCrit + (agility*0.04)*core.MeleeCritRatingPerCritChance
		},
	})
	wp.AddStats(stats.Stats{
		stats.MeleeCrit: float64(warlock.Talents.DemonicTactics) * 1 * core.MeleeCritRatingPerCritChance,
		stats.SpellCrit: float64(warlock.Talents.DemonicTactics) * 1 * core.SpellCritRatingPerCritChance,
	})

	if warlock.Talents.SoulLink {
		wp.PseudoStats.DamageDealtMultiplier *= 1.05
	}
	wp.PseudoStats.DamageDealtMultiplier *= 1.0 + (0.04 * float64(warlock.Talents.UnholyPower))

	wp.EnableManaBar()

	if petConfig.Melee {
		wp.EnableAutoAttacks(wp, core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin:  83.4,
				BaseDamageMax:  123.4,
				SwingSpeed:     2,
				SwingDuration:  time.Second * 2,
				CritMultiplier: 2,
			},
			AutoSwingMelee: true,
		})
	}
	// wp.AutoAttacks.MHEffect.DamageMultiplier *= petConfig.DamageMultiplier
	switch warlock.Options.Summon {
	case proto.Warlock_Options_Imp:
		wp.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.SpellCrit,
			Modifier: func(intellect float64, spellCrit float64) float64 {
				return spellCrit + (0.0125*intellect/100)*core.SpellCritRatingPerCritChance
			},
		})
	case proto.Warlock_Options_Felgaurd:
		wp.PseudoStats.DamageDealtMultiplier *= 1.0 + (0.01 * float64(warlock.Talents.MasterDemonologist))
		// Simulates a pre-stacked demonic frenzy
		wp.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.AttackPower,
			ModifiedStat: stats.AttackPower,
			Modifier: func(ap float64, _ float64) float64 {
				return ap * 1.5 * 1.1 // demonic frenzy + hidden 10% boost
			},
		})
	case proto.Warlock_Options_Succubus:
		wp.PseudoStats.DamageDealtMultiplier *= 1.0 + (0.02 * float64(warlock.Talents.MasterDemonologist))
		wp.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.AttackPower,
			ModifiedStat: stats.AttackPower,
			Modifier: func(ap float64, _ float64) float64 {
				return ap * 1.05 // hidden 5% boost
			},
		})
	}

	if warlock.Talents.FelIntellect > 0 {
		intBonus := 1 + (0.05)*float64(warlock.Talents.FelIntellect)
		wp.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.Intellect,
			Modifier: func(in float64, _ float64) float64 {
				return in * intBonus
			},
		})
	}

	if ItemSetOblivionRaiment.CharacterHasSetBonus(&warlock.Character, 2) {
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
		} else if wp.primaryAbility.CD.Timer != nil {
			wp.WaitUntil(sim, wp.primaryAbility.CD.ReadyAt())
		} else {
			wp.WaitForMana(sim, wp.primaryAbility.CurCast.Cost)
		}
	}
}

var petStatInheritance = func(ownerStats stats.Stats) stats.Stats {
	return stats.Stats{
		stats.Stamina:          ownerStats[stats.Stamina] * 0.3,
		stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
		stats.Armor:            ownerStats[stats.Armor] * 0.35,
		stats.AttackPower:      (ownerStats[stats.SpellPower] + ownerStats[stats.ShadowSpellPower]) * 0.57,
		stats.SpellPower:       (ownerStats[stats.SpellPower] + ownerStats[stats.ShadowSpellPower]) * 0.15,
		stats.SpellPenetration: ownerStats[stats.SpellPenetration],
		// Resists, 40%
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
	proto.Warlock_Options_Felgaurd: {
		Name:             "Felguard",
		Melee:            true,
		PrimaryAbility:   Cleave,
		SecondaryAbility: Intercept,
		ManaIntRatio:     11.5,
		Stats: stats.Stats{
			stats.AttackPower: 20,
			stats.Stamina:     280,
			stats.Strength:    153,
			stats.Agility:     108,
			stats.Intellect:   133,
			stats.Mana:        893,
			stats.Spirit:      122,
			stats.MP5:         48,
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
			stats.AttackPower: 20,
			stats.Stamina:     280,
			stats.Strength:    153,
			stats.Agility:     108,
			stats.Intellect:   133,
			stats.Mana:        893,
			stats.Spirit:      122,
			stats.MP5:         48,
		},
	},
}

// Minion 		Health per bonus stamina 	Mana per bonus intellect
// Imp 			~8.4 						~4.9
// Voidwalker 	~11.0 						~11.5
// Sayaad			~9.1 						~11.5
// Felhunter 	~9.5 						~11.5
// Felguard 	~11.0 						~11.5

// Spell hit 	Spell hit, physical hit, expertise, being capped will cap your minion for all three stats, see below for details
