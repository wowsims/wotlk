package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type WarlockPet struct {
	core.Pet

	owner *Warlock

	primaryAbility   *core.Spell
	secondaryAbility *core.Spell

	DemonicEmpowermentAura *core.Aura
}

const PetExpertiseScale = 1.53

func (warlock *Warlock) NewWarlockPet() *WarlockPet {
	var cfg struct {
		Name          string
		PowerModifier float64
		Stats         stats.Stats
		AutoAttacks   core.AutoAttackOptions
	}

	switch warlock.Options.Summon {
	// TODO: revisit base damage once blizzard fixes JamminL/wotlk-classic-bugs#328
	case proto.Warlock_Options_Felguard:
		cfg.Name = "Felguard"
		cfg.PowerModifier = 0.77 // GetUnitPowerModifier("pet")
		cfg.Stats = stats.Stats{
			stats.Strength:  314,
			stats.Agility:   90,
			stats.Stamina:   328,
			stats.Intellect: 150,
			stats.Spirit:    209,
			stats.Mana:      1559,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		}
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin:  88.8,
				BaseDamageMax:  133.3,
				SwingSpeed:     2,
				SwingDuration:  time.Second * 2,
				CritMultiplier: 2,
			},
			AutoSwingMelee: true,
		}
	case proto.Warlock_Options_Imp:
		cfg.Name = "Imp"
		cfg.PowerModifier = 0.33 // GetUnitPowerModifier("pet")
		cfg.Stats = stats.Stats{
			stats.Strength:  297,
			stats.Agility:   79,
			stats.Stamina:   118,
			stats.Intellect: 369,
			stats.Spirit:    367,
			stats.Mana:      1174,
			stats.MeleeCrit: 3.454 * core.CritRatingPerCritChance,
			stats.SpellCrit: 0.9075 * core.CritRatingPerCritChance,
		}
	case proto.Warlock_Options_Succubus:
		cfg.Name = "Succubus"
		cfg.PowerModifier = 0.77 // GetUnitPowerModifier("pet")
		cfg.Stats = stats.Stats{
			stats.Strength:  314,
			stats.Agility:   90,
			stats.Stamina:   328,
			stats.Intellect: 150,
			stats.Spirit:    209,
			stats.Mana:      1559,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		}
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin:  98,
				BaseDamageMax:  147,
				SwingSpeed:     2,
				SwingDuration:  time.Second * 2,
				CritMultiplier: 2,
			},
			AutoSwingMelee: true,
		}
	case proto.Warlock_Options_Felhunter:
		cfg.Name = "Felhunter"
		cfg.PowerModifier = 0.77 // GetUnitPowerModifier("pet")
		cfg.Stats = stats.Stats{
			stats.Strength:  314,
			stats.Agility:   90,
			stats.Stamina:   328,
			stats.Intellect: 150,
			stats.Spirit:    209,
			stats.Mana:      1559,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		}
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin:  88.8,
				BaseDamageMax:  133.3,
				SwingSpeed:     2,
				SwingDuration:  time.Second * 2,
				CritMultiplier: 2,
			},
			AutoSwingMelee: true,
		}
	}

	wp := &WarlockPet{
		Pet:   core.NewPet(cfg.Name, &warlock.Character, cfg.Stats, warlock.makeStatInheritance(), nil, true, false),
		owner: warlock,
	}

	wp.EnableManaBarWithModifier(cfg.PowerModifier)
	wp.EnableResumeAfterManaWait(wp.OnGCDReady)

	wp.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	wp.AddStat(stats.AttackPower, -20)

	if warlock.Options.Summon == proto.Warlock_Options_Imp {
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

	if warlock.Options.Summon != proto.Warlock_Options_Imp { // imps generally don't meele
		wp.EnableAutoAttacks(wp, cfg.AutoAttacks)
	}

	if warlock.Options.Summon == proto.Warlock_Options_Felguard {
		if wp.owner.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfFelguard) {
			wp.MultiplyStat(stats.AttackPower, 1.2)
		}

		statDeps := []*stats.StatDependency{nil}
		for i := 1; i <= 10; i++ {
			statDeps = append(statDeps, wp.NewDynamicMultiplyStat(stats.AttackPower,
				1+float64(i)*(0.05+0.01*float64(warlock.Talents.DemonicBrutality))))
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
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

	if warlock.Talents.MasterDemonologist > 0 {
		val := 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
		md := core.Aura{
			Label:    "Master Demonologist",
			ActionID: core.ActionID{SpellID: 35706}, // many different spells associated with this talent
			Duration: core.NeverExpires,
			OnGain: func(aura *core.Aura, _ *core.Simulation) {
				switch warlock.Options.Summon {
				case proto.Warlock_Options_Imp:
					aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= val
				case proto.Warlock_Options_Succubus:
					aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= val
				case proto.Warlock_Options_Felguard:
					aura.Unit.PseudoStats.DamageDealtMultiplier *= val
				}
			},
			OnExpire: func(aura *core.Aura, _ *core.Simulation) {
				switch warlock.Options.Summon {
				case proto.Warlock_Options_Imp:
					aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= val
				case proto.Warlock_Options_Succubus:
					aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= val
				case proto.Warlock_Options_Felguard:
					aura.Unit.PseudoStats.DamageDealtMultiplier /= val
				}
			},
		}

		mdLockAura := warlock.RegisterAura(md)
		mdPetAura := wp.RegisterAura(md)

		masterDemonologist := float64(warlock.Talents.MasterDemonologist) * core.CritRatingPerCritChance
		masterDemonologistFireCrit := core.TernaryFloat64(warlock.Options.Summon == proto.Warlock_Options_Imp, masterDemonologist, 0)
		masterDemonologistShadowCrit := core.TernaryFloat64(warlock.Options.Summon == proto.Warlock_Options_Succubus, masterDemonologist, 0)

		wp.OnPetEnable = func(sim *core.Simulation) {
			mdLockAura.Activate(sim)
			mdPetAura.Activate(sim)

			spellbook := make([]*core.Spell, 0)
			spellbook = append(spellbook, warlock.Spellbook...)
			spellbook = append(spellbook, wp.Spellbook...)

			for _, spell := range spellbook {
				if spell.SpellSchool.Matches(core.SpellSchoolFire) {
					spell.BonusCritRating += masterDemonologistFireCrit
				}

				if spell.SpellSchool.Matches(core.SpellSchoolShadow) {
					spell.BonusCritRating += masterDemonologistShadowCrit
				}
			}
		}

		wp.OnPetDisable = func(sim *core.Simulation) {
			mdLockAura.Deactivate(sim)
			mdPetAura.Deactivate(sim)

			spellbook := make([]*core.Spell, 0)
			spellbook = append(spellbook, warlock.Spellbook...)
			spellbook = append(spellbook, wp.Spellbook...)

			for _, spell := range spellbook {
				if spell.SpellSchool.Matches(core.SpellSchoolFire) {
					spell.BonusCritRating -= masterDemonologistFireCrit
				}

				if spell.SpellSchool.Matches(core.SpellSchoolShadow) {
					spell.BonusCritRating -= masterDemonologistShadowCrit
				}
			}
		}
	}

	core.ApplyPetConsumeEffects(&wp.Character, warlock.Consumes)

	warlock.AddPet(wp)

	return wp
}

func (wp *WarlockPet) GetPet() *core.Pet {
	return &wp.Pet
}

func (wp *WarlockPet) Initialize() {
	switch wp.owner.Options.Summon {
	case proto.Warlock_Options_Felguard:
		wp.primaryAbility = wp.NewPetAbility(Cleave, true)
		wp.secondaryAbility = wp.NewPetAbility(Intercept, false)
	case proto.Warlock_Options_Succubus:
		wp.primaryAbility = wp.NewPetAbility(LashOfPain, true)
	case proto.Warlock_Options_Felhunter:
		wp.primaryAbility = wp.NewPetAbility(ShadowBite, true)
	case proto.Warlock_Options_Imp:
		wp.primaryAbility = wp.NewPetAbility(Firebolt, true)
	}
}

func (wp *WarlockPet) Reset(sim *core.Simulation) {
}

func (wp *WarlockPet) OnGCDReady(sim *core.Simulation) {
	target := wp.CurrentTarget

	if !wp.TryCast(sim, target, wp.primaryAbility) {
		if !wp.primaryAbility.IsReady(sim) {
			wp.WaitUntil(sim, wp.primaryAbility.CD.ReadyAt())
		} else {
			wp.WaitForMana(sim, wp.primaryAbility.CurCast.Cost)
		}
	}
}

func (warlock *Warlock) makeStatInheritance() core.PetStatInheritance {
	improvedDemonicTactics := float64(warlock.Talents.ImprovedDemonicTactics)

	return func(ownerStats stats.Stats) stats.Stats {
		// EJ posts claim this value is passed through math.Floor, but in-game testing
		// shows pets benefit from each point of owner hit rating in WotLK Classic.
		// https://web.archive.org/web/20120112003252/http://elitistjerks.com/f80/t100099-demonology_releasing_demon_you
		ownerHitChance := ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance

		// TODO: Account for sunfire/soulfrost
		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
			stats.Armor:            ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:      ownerStats[stats.SpellPower] * 0.57,
			stats.SpellPower:       ownerStats[stats.SpellPower] * 0.15,
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
