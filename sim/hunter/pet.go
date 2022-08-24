package hunter

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type HunterPet struct {
	core.Pet
	focusBar

	config PetConfig

	hunterOwner *Hunter

	CobraStrikesAura *core.Aura
	KillCommandAura  *core.Aura

	specialAbility PetAbility
	focusDump      PetAbility

	uptimePercent    float64
	hasOwnerCooldown bool
}

func (hunter *Hunter) NewHunterPet() *HunterPet {
	if hunter.Options.PetType == proto.Hunter_Options_PetNone {
		return nil
	}
	if hunter.Options.PetUptime <= 0 {
		return nil
	}
	petConfig := PetConfigs[hunter.Options.PetType]

	hp := &HunterPet{
		Pet: core.NewPet(
			petConfig.Name,
			&hunter.Character,
			hunterPetBaseStats,
			hunter.makeStatInheritance(),
			true,
			false,
		),
		config:      petConfig,
		hunterOwner: hunter,

		hasOwnerCooldown: petConfig.SpecialAbility == FuriousHowl || petConfig.SpecialAbility == SavageRend,
	}

	hp.EnableFocusBar(1.0+0.5*float64(hunter.Talents.BestialDiscipline), func(sim *core.Simulation) {
		if hp.GCD.IsReady(sim) {
			hp.OnGCDReady(sim)
		}
	})

	hp.EnableAutoAttacks(hp, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  59,
			BaseDamageMax:  87,
			SwingSpeed:     2,
			SwingDuration:  time.Second * 2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})

	// Happiness
	hp.PseudoStats.DamageDealtMultiplier *= 1.25

	// Pet family bonus is now the same for all pets.
	hp.PseudoStats.PhysicalDamageDealtMultiplier *= 1.05

	hp.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	hp.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/62.77)
	core.ApplyPetConsumeEffects(&hp.Character, hunter.Consumes)

	hunter.AddPet(hp)

	return hp
}

func (hp *HunterPet) GetPet() *core.Pet {
	return &hp.Pet
}

func (hp *HunterPet) Talents() proto.HunterPetTalents {
	talents := hp.hunterOwner.Options.PetTalents
	if talents == nil {
		return proto.HunterPetTalents{}
	} else {
		return *talents
	}
}

func (hp *HunterPet) Initialize() {
	//if hp.hunterOwner.Options.PetSingleAbility {
	//	hp.specialAbility = hp.NewPetAbility(hp.config.FocusDump, true)
	//	hp.config.RandomSelection = false
	//} else {
	hp.specialAbility = hp.NewPetAbility(hp.config.SpecialAbility, true)
	hp.focusDump = hp.NewPetAbility(hp.config.FocusDump, false)
	//}
}

func (hp *HunterPet) Reset(sim *core.Simulation) {
	hp.focusBar.reset(sim)
	hp.uptimePercent = core.MinFloat(1, core.MaxFloat(0, hp.hunterOwner.Options.PetUptime))
}

func (hp *HunterPet) OnGCDReady(sim *core.Simulation) {
	percentRemaining := sim.GetRemainingDurationPercent()
	if percentRemaining < 1.0-hp.uptimePercent { // once fight is % completed, disable pet.
		hp.Disable(sim)
		hp.focusBar.Cancel(sim)
		return
	}

	if hp.hasOwnerCooldown && hp.CurrentFocus() < 50 {
		// When a major ability (Furious Howl or Savage Rend) is ready, pool enough
		// energy to use on-demand.
		hp.DoNothing()
		return
	}

	target := hp.CurrentTarget
	if hp.config.RandomSelection {
		if sim.RandomFloat("Hunter Pet Ability") < 0.5 {
			if !hp.specialAbility.TryCast(sim, target, hp) {
				if !hp.focusDump.TryCast(sim, target, hp) {
					hp.DoNothing()
				}
			}
		} else {
			if !hp.focusDump.TryCast(sim, target, hp) {
				if !hp.specialAbility.TryCast(sim, target, hp) {
					hp.DoNothing()
				}
			}
		}
		return
	}

	if hp.specialAbility.TryCast(sim, target, hp) {
		// For abilities that don't use the GCD.
		if hp.GCD.IsReady(sim) {
			if hp.focusDump.Type != Unknown {
				if !hp.focusDump.TryCast(sim, target, hp) {
					hp.DoNothing()
				}
			} else {
				hp.DoNothing()
			}
		}
	} else {
		if hp.focusDump.Type != Unknown {
			if !hp.focusDump.TryCast(sim, target, hp) {
				hp.DoNothing()
			}
		} else {
			hp.DoNothing()
		}
	}
}

func (hp *HunterPet) specialDamageMod(baseDamageConfig core.BaseDamageConfig) core.BaseDamageConfig {
	return core.WrapBaseDamageConfig(baseDamageConfig, func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
		return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
			normalDamage := oldCalculator(sim, hitEffect, spell)
			return normalDamage * (1 + 0.2*float64(hp.KillCommandAura.GetStacks()))
		}
	})
}

var hunterPetBaseStats = stats.Stats{
	stats.Agility:     113,
	stats.Strength:    331,
	stats.AttackPower: -20, // Apparently pets and warriors have a AP penalty.

	// Add 1.8% because pets aren't affected by that component of crit suppression.
	stats.MeleeCrit: (3.2 + 1.8) * core.CritRatingPerCritChance,
}

const PetExpertiseScale = 3.25

func (hunter *Hunter) makeStatInheritance() core.PetStatInheritance {
	hvw := hunter.Talents.HunterVsWild

	petTalents := hunter.Options.PetTalents
	var wildHunt int32
	if petTalents != nil {
		wildHunt = petTalents.WildHunt
	}

	return func(ownerStats stats.Stats) stats.Stats {
		ownerHitChance := ownerStats[stats.MeleeHit] / core.MeleeHitRatingPerHitChance
		hitRatingFromOwner := math.Floor(ownerHitChance) * core.MeleeHitRatingPerHitChance

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.3 * (1 + 0.2*float64(wildHunt)),
			stats.Armor:       ownerStats[stats.Armor] * 0.35,
			stats.AttackPower: ownerStats[stats.RangedAttackPower]*0.22*(1+0.15*float64(wildHunt)) + ownerStats[stats.Stamina]*0.1*float64(hvw),

			stats.MeleeHit:  hitRatingFromOwner,
			stats.SpellHit:  hitRatingFromOwner * 2,
			stats.Expertise: math.Floor((math.Floor(ownerHitChance) * PetExpertiseScale)) * core.ExpertisePerQuarterPercentReduction,
		}
	}
}

type PetConfig struct {
	Name string

	SpecialAbility PetAbilityType
	FocusDump      PetAbilityType

	// Randomly select between abilities instead of using a prio.
	RandomSelection bool
}

// Abilities reference: https://wotlk.wowhead.com/hunter-pets
// https://wotlk.wowhead.com/guides/hunter-dps-best-pets-taming-loyalty-burning-crusade-classic
var PetConfigs = map[proto.Hunter_Options_PetType]PetConfig{
	proto.Hunter_Options_Bat: PetConfig{
		Name:           "Bat",
		SpecialAbility: SonicBlast,
		FocusDump:      Claw,
	},
	proto.Hunter_Options_Bear: PetConfig{
		Name:           "Bear",
		SpecialAbility: Swipe,
		FocusDump:      Claw,
	},
	proto.Hunter_Options_BirdOfPrey: PetConfig{
		Name:           "Bird of Prey",
		SpecialAbility: Snatch,
		FocusDump:      Claw,
	},
	proto.Hunter_Options_Boar: PetConfig{
		Name:           "Boar",
		SpecialAbility: Gore,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_CarrionBird: PetConfig{
		Name:           "Carrion Bird",
		SpecialAbility: DemoralizingScreech,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Cat: PetConfig{
		Name:           "Cat",
		SpecialAbility: Rake,
		FocusDump:      Claw,
	},
	proto.Hunter_Options_Chimaera: PetConfig{
		Name:           "Chimaera",
		SpecialAbility: FroststormBreath,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_CoreHound: PetConfig{
		Name:           "Core Hound",
		SpecialAbility: LavaBreath,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Crab: PetConfig{
		Name:           "Crab",
		SpecialAbility: Pin,
		FocusDump:      Claw,
	},
	proto.Hunter_Options_Crocolisk: PetConfig{
		Name: "Crocolisk",
		//SpecialAbility: BadAttitude,
		FocusDump: Bite,
	},
	proto.Hunter_Options_Devilsaur: PetConfig{
		Name:           "Devilsaur",
		SpecialAbility: MonstrousBite,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Dragonhawk: PetConfig{
		Name:           "Dragonhawk",
		SpecialAbility: FireBreath,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Gorilla: PetConfig{
		Name: "Gorilla",
		//SpecialAbility: Pummel,
		FocusDump: Smack,
	},
	proto.Hunter_Options_Hyena: PetConfig{
		Name:           "Hyena",
		SpecialAbility: TendonRip,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Moth: PetConfig{
		Name: "Moth",
		//SpecialAbility:   SerentiyDust,
		FocusDump: Smack,
	},
	proto.Hunter_Options_NetherRay: PetConfig{
		Name:           "Nether Ray",
		SpecialAbility: NetherShock,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Raptor: PetConfig{
		Name:           "Raptor",
		SpecialAbility: SavageRend,
		FocusDump:      Claw,
	},
	proto.Hunter_Options_Ravager: PetConfig{
		Name:           "Ravager",
		SpecialAbility: Ravage,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Rhino: PetConfig{
		Name:           "Rhino",
		SpecialAbility: Stampede,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Scorpid: PetConfig{
		Name:           "Scorpid",
		SpecialAbility: ScorpidPoison,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Serpent: PetConfig{
		Name:           "Serpent",
		SpecialAbility: PoisonSpit,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Silithid: PetConfig{
		Name:           "Silithid",
		SpecialAbility: VenomWebSpray,
		FocusDump:      Claw,
	},
	proto.Hunter_Options_Spider: PetConfig{
		Name: "Spider",
		//SpecialAbility:   Web,
		FocusDump: Bite,
	},
	proto.Hunter_Options_SpiritBeast: PetConfig{
		Name:           "Spirit Beast",
		SpecialAbility: SpiritStrike,
		FocusDump:      Claw,
	},
	proto.Hunter_Options_SporeBat: PetConfig{
		Name:           "Spore Bat",
		SpecialAbility: SporeCloud,
		FocusDump:      Smack,
	},
	proto.Hunter_Options_Tallstrider: PetConfig{
		Name: "Tallstrider",
		//SpecialAbility:   DustCloud,
		FocusDump: Claw,
	},
	proto.Hunter_Options_Turtle: PetConfig{
		Name: "Turtle",
		//SpecialAbility: ShellShield,
		FocusDump: Bite,
	},
	proto.Hunter_Options_WarpStalker: PetConfig{
		Name: "Warp Stalker",
		//SpecialAbility:   Warp,
		FocusDump: Bite,
	},
	proto.Hunter_Options_Wasp: PetConfig{
		Name:           "Wasp",
		SpecialAbility: Sting,
		FocusDump:      Smack,
	},
	proto.Hunter_Options_WindSerpent: PetConfig{
		Name:           "Wind Serpent",
		SpecialAbility: LightningBreath,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Wolf: PetConfig{
		Name:           "Wolf",
		SpecialAbility: FuriousHowl,
		FocusDump:      Bite,
	},
	proto.Hunter_Options_Worm: PetConfig{
		Name:           "Worm",
		SpecialAbility: AcidSpit,
		FocusDump:      Bite,
	},
}
