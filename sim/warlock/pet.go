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

const PetExpertiseScale = 1.53

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

	wp.EnableManaBarWithModifier(petConfig.PowerModifier)
	wp.EnableResumeAfterManaWait(wp.OnGCDReady)

	wp.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	wp.AddStat(stats.AttackPower, -20)

	if summonChoice == proto.Warlock_Options_Imp {
		// imp has a slightly different agi crit scaling coef for some reason
		wp.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance*1/51.0204)
	} else {
		wp.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance*1/52.0833)
	}

	wp.AddStats(stats.Stats{
		stats.MeleeCrit: float64(warlock.Talents.DemonicTactics) * 2 * core.CritRatingPerCritChance,
		stats.SpellCrit: float64(warlock.Talents.DemonicTactics) * 2 * core.CritRatingPerCritChance,

		// Remove stats the pet incorrectly has because of the suppression talent through stat inheritance
		stats.MeleeHit:  -float64(warlock.Talents.Suppression) * core.MeleeHitRatingPerHitChance,
		stats.SpellHit:  -float64(warlock.Talents.Suppression) * core.SpellHitRatingPerHitChance,
		stats.Expertise: -float64(warlock.Talents.Suppression) * PetExpertiseScale * core.ExpertisePerQuarterPercentReduction,
	})

	wp.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.04*float64(warlock.Talents.UnholyPower)

	if petConfig.Melee {
		switch summonChoice {
		// TODO: revisit base damage once blizzard fixes JamminL/wotlk-classic-bugs#328
		case proto.Warlock_Options_Felguard:
			wp.EnableAutoAttacks(wp, core.AutoAttackOptions{
				MainHand: core.Weapon{
					BaseDamageMin:  88.8,
					BaseDamageMax:  133.3,
					SwingSpeed:     2,
					SwingDuration:  time.Second * 2,
					CritMultiplier: 2,
				},
				AutoSwingMelee: true,
			})
		case proto.Warlock_Options_Succubus:
			wp.EnableAutoAttacks(wp, core.AutoAttackOptions{
				MainHand: core.Weapon{
					BaseDamageMin:  98,
					BaseDamageMax:  147,
					SwingSpeed:     2,
					SwingDuration:  time.Second * 2,
					CritMultiplier: 2,
				},
				AutoSwingMelee: true,
			})
		case proto.Warlock_Options_Felhunter:
			wp.EnableAutoAttacks(wp, core.AutoAttackOptions{
				MainHand: core.Weapon{
					BaseDamageMin:  88.8,
					BaseDamageMax:  133.3,
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
		wp.PseudoStats.FireDamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
		wp.PseudoStats.BonusFireCritRating *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
	case proto.Warlock_Options_Succubus:
		wp.PseudoStats.ShadowDamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
		wp.PseudoStats.BonusShadowCritRating *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
	case proto.Warlock_Options_Felguard:
		wp.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)

		talentMultiplier := 1. + 0.1*float64(warlock.Talents.DemonicBrutality)
		glyphMultiplier := 1.
		if wp.owner.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfFelguard) {
			glyphMultiplier *= 1.2
		}
		wp.MultiplyStat(stats.AttackPower, talentMultiplier*glyphMultiplier)

		statDeps := []*stats.StatDependency{nil}
		for i := 0; i <= 10; i++ {
			statDeps = append(statDeps, wp.NewDynamicMultiplyStat(stats.AttackPower, (1+0.1*float64(warlock.Talents.DemonicBrutality)+0.05*float64(i))/talentMultiplier))
		}

		DemonicFrenzyAura := wp.RegisterAura(core.Aura{
			Label:     "Demonic Frenzy",
			ActionID:  core.ActionID{SpellID: 32851},
			Duration:  time.Second * 10,
			MaxStacks: 10,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				if oldStacks != 0 {
					aura.Unit.DisableDynamicStatDep(sim, statDeps[oldStacks])
				}
				if newStacks != 0 {
					aura.Unit.EnableDynamicStatDep(sim, statDeps[newStacks])
				}
			},
		})
		wp.RegisterAura(core.Aura{
			Label:    "Demonic Frenzy Hidden Aura",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			// OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// 	aura.Unit.EnableDynamicStatDep(sim, statDeps[0])
			// },
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}
				DemonicFrenzyAura.Activate(sim)
				DemonicFrenzyAura.AddStack(sim)
			},
		})
	}

	if warlock.Talents.FelVitality > 0 {
		bonus := 1.0 + 0.05*float64(warlock.Talents.FelVitality)
		wp.MultiplyStat(stats.Intellect, bonus)
		wp.MultiplyStat(stats.Stamina, bonus)
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
		highestSP := core.MaxFloat(ownerStats[stats.ArcaneSpellPower],
			core.MaxFloat(ownerStats[stats.FireSpellPower], core.MaxFloat(ownerStats[stats.FrostSpellPower],
				core.MaxFloat(ownerStats[stats.HolySpellPower], core.MaxFloat(ownerStats[stats.NatureSpellPower],
					ownerStats[stats.ShadowSpellPower])))))

		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
			stats.Armor:            ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:      (ownerStats[stats.SpellPower] + highestSP) * 0.57,
			stats.SpellPower:       (ownerStats[stats.SpellPower] + highestSP) * 0.15,
			stats.SpellPenetration: ownerStats[stats.SpellPenetration],
			stats.SpellCrit:        improvedDemonicTactics * 0.1 * ownerStats[stats.SpellCrit],
			stats.MeleeCrit:        improvedDemonicTactics * 0.1 * ownerStats[stats.SpellCrit],
			stats.MeleeHit:         ownerHitChance * core.MeleeHitRatingPerHitChance,
			stats.SpellHit:         ownerHitChance * core.SpellHitRatingPerHitChance,
			// TODO: revisit
			stats.Expertise: (ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance) *
				PetExpertiseScale * core.ExpertisePerQuarterPercentReduction,

			// Resists, 40%

			// TODO: does the pet scale with the 1% hit from draenei?
		}
	}
}

type PetConfig struct {
	Name string
	// DamageMultiplier float64
	Melee         bool
	Stats         stats.Stats
	PowerModifier float64

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
		PowerModifier:    0.77, // GetUnitPowerModifier("pet")
		Stats: stats.Stats{
			stats.Strength:  314,
			stats.Agility:   90,
			stats.Stamina:   328,
			stats.Intellect: 150,
			stats.Spirit:    209,
			stats.Mana:      1559,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		},
	},
	proto.Warlock_Options_Imp: {
		Name:           "Imp",
		PowerModifier:  0.33, // GetUnitPowerModifier("pet")
		Melee:          false,
		PrimaryAbility: Firebolt,
		Stats: stats.Stats{
			stats.Strength:  297,
			stats.Agility:   79,
			stats.Stamina:   118,
			stats.Intellect: 369,
			stats.Spirit:    367,
			stats.Mana:      1174,
			stats.MeleeCrit: 3.454 * core.CritRatingPerCritChance,
			stats.SpellCrit: 0.9075 * core.CritRatingPerCritChance,
		},
	},
	proto.Warlock_Options_Succubus: {
		Name:           "Succubus",
		PowerModifier:  0.77, // GetUnitPowerModifier("pet")
		Melee:          true,
		PrimaryAbility: LashOfPain,
		Stats: stats.Stats{
			stats.Strength:  314,
			stats.Agility:   90,
			stats.Stamina:   328,
			stats.Intellect: 150,
			stats.Spirit:    209,
			stats.Mana:      1559,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		},
	},
	proto.Warlock_Options_Felhunter: {
		Name:           "Felhunter",
		PowerModifier:  0.77, // GetUnitPowerModifier("pet")
		Melee:          true,
		PrimaryAbility: ShadowBite,
		Stats: stats.Stats{
			stats.Strength:  314,
			stats.Agility:   90,
			stats.Stamina:   328,
			stats.Intellect: 150,
			stats.Spirit:    209,
			stats.Mana:      1559,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		},
	},
}
