package hunter

import (
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

	primaryAbility   PetAbility
	secondaryAbility PetAbility

	uptimePercent float64
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
		),
		config:      petConfig,
		hunterOwner: hunter,
	}

	// Happiness
	hp.PseudoStats.DamageDealtMultiplier *= 1.25

	hp.EnableFocusBar(1.0+0.5*float64(hunter.Talents.BestialDiscipline), func(sim *core.Simulation) {
		if hp.GCD.IsReady(sim) {
			hp.OnGCDReady(sim)
		}
	})

	hp.EnableAutoAttacks(hp, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  42,
			BaseDamageMax:  68,
			SwingSpeed:     2,
			SwingDuration:  time.Second * 2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})

	// Pet family bonus is now the same for all pets.
	hp.AutoAttacks.MHEffect.DamageMultiplier *= 1.05

	hp.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.AttackPower,
		Modifier: func(strength float64, attackPower float64) float64 {
			return attackPower + strength*2
		},
	})
	hp.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.MeleeCrit,
		Modifier: func(agility float64, meleeCrit float64) float64 {
			return meleeCrit + (agility/33)*core.CritRatingPerCritChance
		},
	})

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
	if hp.hunterOwner.Options.PetSingleAbility {
		hp.primaryAbility = hp.NewPetAbility(hp.config.SecondaryAbility, true)
		hp.config.RandomSelection = false
	} else {
		hp.primaryAbility = hp.NewPetAbility(hp.config.PrimaryAbility, true)
		hp.secondaryAbility = hp.NewPetAbility(hp.config.SecondaryAbility, false)
	}
}

func (hp *HunterPet) Reset(sim *core.Simulation) {
	hp.focusBar.reset(sim)
	if sim.Log != nil {
		hp.Log(sim, "Total Pet stats: %s", hp.GetStats())
		inheritedStats := hp.hunterOwner.makeStatInheritance()(hp.hunterOwner.GetStats())
		hp.Log(sim, "Inherited Pet stats: %s", inheritedStats)
	}

	hp.uptimePercent = core.MinFloat(1, core.MaxFloat(0, hp.hunterOwner.Options.PetUptime))
}

func (hp *HunterPet) OnGCDReady(sim *core.Simulation) {
	percentRemaining := sim.GetRemainingDurationPercent()
	if percentRemaining < 1.0-hp.uptimePercent { // once fight is % completed, disable pet.
		hp.Disable(sim)
		hp.focusBar.Cancel(sim)
		return
	}

	target := hp.CurrentTarget
	if hp.config.RandomSelection {
		if sim.RandomFloat("Hunter Pet Ability") < 0.5 {
			if !hp.primaryAbility.TryCast(sim, target, hp) {
				if !hp.secondaryAbility.TryCast(sim, target, hp) {
					hp.DoNothing()
				}
			}
		} else {
			if !hp.secondaryAbility.TryCast(sim, target, hp) {
				if !hp.primaryAbility.TryCast(sim, target, hp) {
					hp.DoNothing()
				}
			}
		}
		return
	}

	if !hp.primaryAbility.TryCast(sim, target, hp) {
		if hp.secondaryAbility.Type != Unknown {
			if !hp.secondaryAbility.TryCast(sim, target, hp) {
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
			if hp.KillCommandAura.IsActive() {
				return normalDamage * (1 + 0.2*float64(hp.KillCommandAura.GetStacks()))
			} else {
				return normalDamage
			}
		}
	})
}

func (hp *HunterPet) specialOutcomeMod(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
	return func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, attackTable *core.AttackTable) {
		if hp.CobraStrikesAura.IsActive() {
			hp.AddStatDynamic(sim, stats.MeleeCrit, 100*core.CritRatingPerCritChance)
			hp.AddStatDynamic(sim, stats.SpellCrit, 100*core.CritRatingPerCritChance)
			outcomeApplier(sim, spell, spellEffect, attackTable)
			hp.AddStatDynamic(sim, stats.MeleeCrit, -100*core.CritRatingPerCritChance)
			hp.AddStatDynamic(sim, stats.SpellCrit, -100*core.CritRatingPerCritChance)
		} else if hp.KillCommandAura.IsActive() && hp.hunterOwner.Talents.FocusedFire > 0 {
			bonusCrit := 10 * core.CritRatingPerCritChance * float64(hp.hunterOwner.Talents.FocusedFire)
			hp.AddStatDynamic(sim, stats.MeleeCrit, bonusCrit)
			hp.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
			outcomeApplier(sim, spell, spellEffect, attackTable)
			hp.AddStatDynamic(sim, stats.MeleeCrit, -bonusCrit)
			hp.AddStatDynamic(sim, stats.SpellCrit, -bonusCrit)
		} else {
			outcomeApplier(sim, spell, spellEffect, attackTable)
		}
	}
}

var hunterPetBaseStats = stats.Stats{
	stats.Agility:     127,
	stats.Strength:    162,
	stats.AttackPower: -20, // Apparently pets and warriors have a AP penalty.

	// Add 1.8% because pets aren't affected by that component of crit suppression.
	stats.MeleeCrit: (1.1515 + 1.8) * core.CritRatingPerCritChance,
}

func (hunter *Hunter) makeStatInheritance() core.PetStatInheritance {
	hvw := 0.1 * float64(hunter.Talents.HunterVsWild)

	petTalents := hunter.Options.PetTalents
	var wildHunt int32
	if petTalents != nil {
		wildHunt = petTalents.WildHunt
	}

	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina]*0.3 + 0.2*float64(wildHunt),
			stats.Armor:       ownerStats[stats.Armor] * 0.35,
			stats.AttackPower: ownerStats[stats.RangedAttackPower]*0.22 + ownerStats[stats.Stamina]*hvw + 0.15*float64(wildHunt),
			stats.SpellPower:  ownerStats[stats.RangedAttackPower]*0.128 + 0.15*float64(wildHunt),
		}
	}
}

type PetConfig struct {
	Name string

	PrimaryAbility   PetAbilityType
	SecondaryAbility PetAbilityType

	// Randomly select between abilities instead of using a prio.
	RandomSelection bool
}

// Abilities reference: https://wotlk.wowhead.com/hunter-pets
// https://wotlk.wowhead.com/guides/hunter-dps-best-pets-taming-loyalty-burning-crusade-classic
var PetConfigs = map[proto.Hunter_Options_PetType]PetConfig{
	proto.Hunter_Options_Bat: PetConfig{
		Name:             "Bat",
		PrimaryAbility:   Bite,
		SecondaryAbility: Screech,
	},
	proto.Hunter_Options_Bear: PetConfig{
		Name:             "Bear",
		PrimaryAbility:   Bite,
		SecondaryAbility: Claw,
	},
	proto.Hunter_Options_Cat: PetConfig{
		Name:             "Cat",
		PrimaryAbility:   Bite,
		SecondaryAbility: Claw,
	},
	proto.Hunter_Options_Crab: PetConfig{
		Name:           "Crab",
		PrimaryAbility: Claw,
	},
	proto.Hunter_Options_Owl: PetConfig{
		Name:             "Owl",
		PrimaryAbility:   Claw,
		SecondaryAbility: Screech,
		RandomSelection:  true,
	},
	proto.Hunter_Options_Raptor: PetConfig{
		Name:             "Raptor",
		PrimaryAbility:   Bite,
		SecondaryAbility: Claw,
	},
	proto.Hunter_Options_Ravager: PetConfig{
		Name:             "Ravager",
		PrimaryAbility:   Bite,
		SecondaryAbility: Gore,
	},
	proto.Hunter_Options_WindSerpent: PetConfig{
		Name:             "Wind Serpent",
		PrimaryAbility:   Bite,
		SecondaryAbility: LightningBreath,
	},
}
