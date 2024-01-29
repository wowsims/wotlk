package warlock

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type WarlockPet struct {
	core.Pet

	owner *Warlock

	primaryAbility   *core.Spell
	secondaryAbility *core.Spell

	DemonicEmpowermentAura *core.Aura

	manaPooling bool
}

// TODO: Classic warlock pet stats
func (warlock *Warlock) NewWarlockPet() *WarlockPet {
	var cfg struct {
		Name          string
		PowerModifier float64
		Stats         stats.Stats
		AutoAttacks   core.AutoAttackOptions
	}

	switch warlock.Options.Summon {
	// TODO: revisit base damage once blizzard fixes JamminL/wotlk-classic-bugs#328
	case proto.WarlockOptions_Imp:
		cfg.Name = "Imp"
		cfg.PowerModifier = 0.33 // GetUnitPowerModifier("pet")
		cfg.Stats = stats.Stats{
			stats.Strength:  47,
			stats.Agility:   25,
			stats.Stamina:   70,
			stats.Intellect: 94,
			stats.Spirit:    95,
			stats.Mana:      149,
			stats.MP5:       0,
			stats.MeleeCrit: 3.454 * core.CritRatingPerCritChance,
			stats.SpellCrit: 0.9075 * core.CritRatingPerCritChance,
		}
	case proto.WarlockOptions_Succubus:
		cfg.Name = "Succubus"
		cfg.PowerModifier = 0.77 // GetUnitPowerModifier("pet")
		cfg.Stats = stats.Stats{
			stats.Strength:  50,
			stats.Agility:   40,
			stats.Stamina:   100,
			stats.Intellect: 50,
			stats.Spirit:    51,
			stats.Mana:      60,
			stats.MP5:       0,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		}
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin:  23,
				BaseDamageMax:  38,
				SwingSpeed:     2,
				CritMultiplier: 2,
			},
			AutoSwingMelee: true,
		}
	case proto.WarlockOptions_Felhunter:
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
				CritMultiplier: 2,
			},
			AutoSwingMelee: true,
		}
	}

	wp := &WarlockPet{
		Pet:   core.NewPet(cfg.Name, &warlock.Character, cfg.Stats, warlock.makeStatInheritance(), true, false),
		owner: warlock,
	}

	wp.EnableManaBarWithModifier(cfg.PowerModifier)

	wp.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	wp.AddStat(stats.AttackPower, -20)

	if warlock.Options.Summon == proto.WarlockOptions_Imp {
		// imp has a slightly different agi crit scaling coef for some reason
		wp.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance*1/51.0204)

		// Imp gets 1mp/5 non casting regen per spirit
		wp.PseudoStats.SpiritRegenMultiplier = 1
		wp.PseudoStats.SpiritRegenRateCasting = 0
		wp.SpiritManaRegenPerSecond = func() float64 {
			// 1mp5 per spirit
			return wp.GetStat(stats.Spirit) / 5
		}
	} else {
		wp.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance*1/52.0833)
	}

	wp.AutoAttacks.MHConfig().DamageMultiplier *= 1.0 + 0.04*float64(warlock.Talents.UnholyPower)

	if warlock.Options.Summon != proto.WarlockOptions_Imp { // imps generally don't melee
		wp.EnableAutoAttacks(wp, cfg.AutoAttacks)
	}

	if warlock.Talents.MasterDemonologist > 0 {
		md := core.Aura{
			Label:    "Master Demonologist",
			ActionID: core.ActionID{SpellID: 23825},
			Duration: core.NeverExpires,
			OnGain: func(aura *core.Aura, _ *core.Simulation) {
				switch warlock.Options.Summon {
				case proto.WarlockOptions_Imp:
					aura.Unit.PseudoStats.ThreatMultiplier /= 1 + 0.04*float64(warlock.Talents.MasterDemonologist)
				case proto.WarlockOptions_Succubus:
					aura.Unit.PseudoStats.DamageDealtMultiplier *= 1 + 0.02*float64(warlock.Talents.MasterDemonologist)
				}
			},
			OnExpire: func(aura *core.Aura, _ *core.Simulation) {
				switch warlock.Options.Summon {
				case proto.WarlockOptions_Imp:
					aura.Unit.PseudoStats.ThreatMultiplier *= 1 + 0.04*float64(warlock.Talents.MasterDemonologist)
				case proto.WarlockOptions_Succubus:
					aura.Unit.PseudoStats.DamageDealtMultiplier /= 1 + 0.02*float64(warlock.Talents.MasterDemonologist)
				}
			},
		}

		mdLockAura := warlock.RegisterAura(md)
		mdPetAura := wp.RegisterAura(md)

		wp.OnPetEnable = func(sim *core.Simulation) {
			mdLockAura.Activate(sim)
			mdPetAura.Activate(sim)
		}

		wp.OnPetDisable = func(sim *core.Simulation) {
			mdLockAura.Deactivate(sim)
			mdPetAura.Deactivate(sim)
		}
	}

	// core.ApplyPetConsumeEffects(&wp.Character, warlock.Consumes)

	warlock.AddPet(wp)

	return wp
}

func (wp *WarlockPet) GetPet() *core.Pet {
	return &wp.Pet
}

// TODO: Classic warlock pet abilities
func (wp *WarlockPet) Initialize() {
	switch wp.owner.Options.Summon {
	case proto.WarlockOptions_Succubus:
		wp.registerLashOfPainSpell()
	case proto.WarlockOptions_Felhunter:
		// wp.registerShadowBiteSpell()
	case proto.WarlockOptions_Imp:
		wp.registerFireboltSpell()
	}
}

func (wp *WarlockPet) Reset(_ *core.Simulation) {
}

func (wp *WarlockPet) OnGCDReady(sim *core.Simulation) {
	if wp.manaPooling {
		maxPossibleCasts := sim.GetRemainingDuration().Seconds() / wp.primaryAbility.CurCast.CastTime.Seconds()

		if wp.CurrentMana() > (maxPossibleCasts*wp.primaryAbility.CurCast.Cost)*0.75 {
			wp.manaPooling = false
			wp.WaitUntil(sim, sim.CurrentTime+10*time.Millisecond)
			return
		}

		if wp.CurrentMana() >= wp.MaxMana()*0.94 {
			wp.manaPooling = false
			wp.WaitUntil(sim, sim.CurrentTime+10*time.Millisecond)
			return
		}

		if wp.manaPooling {
			return
		}
	}

	if !wp.primaryAbility.IsReady(sim) {
		wp.WaitUntil(sim, wp.primaryAbility.CD.ReadyAt())
		return
	}

	wp.primaryAbility.Cast(sim, wp.CurrentTarget)
}

func (warlock *Warlock) makeStatInheritance() core.PetStatInheritance {

	return func(ownerStats stats.Stats) stats.Stats {
		// based on testing for WotLK Classic the following is true:
		// - pets are meele hit capped if and only if the warlock has 210 (8%) spell hit rating or more
		//   - this is unaffected by suppression and by magic hit debuffs like FF
		// - pets gain expertise from 0% to 6.5% relative to the owners hit, reaching cap at 17% spell hit
		//   - this is also unaffected by suppression and by magic hit debuffs like FF
		//   - this is continious, i.e. not restricted to 0.25 intervals
		// - pets gain spell hit from 0% to 17% relative to the owners hit, reaching cap at 12% spell hit
		// spell hit rating is floor'd
		//   - affected by suppression and ff, but in weird ways:
		// 3/3 suppression => 262 hit  (9.99%) results in misses, 263 (10.03%) no misses
		// 2/3 suppression => 278 hit (10.60%) results in misses, 279 (10.64%) no misses
		// 1/3 suppression => 288 hit (10.98%) results in misses, 289 (11.02%) no misses
		// 0/3 suppression => 314 hit (11.97%) results in misses, 315 (12.01%) no misses
		// 3/3 suppression + FF => 209 hit (7.97%) results in misses, 210 (8.01%) no misses
		// 2/3 suppression + FF => 222 hit (8.46%) results in misses, 223 (8.50%) no misses
		//
		// the best approximation of this behaviour is that we scale the warlock's spell hit by `1/12*17` floor
		// the result and then add the hit percent from suppression/ff

		// does correctly not include ff/misery
		ownerHitChance := ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance

		// TODO: Classic warlock pet stat deps
		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * core.Ternary(warlock.Options.Summon == proto.WarlockOptions_Imp, 0.66, 0.75),
			stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
			stats.Armor:            ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:      ownerStats[stats.SpellPower] * 0.57,
			stats.MP5:              ownerStats[stats.MP5] * 0.3,
			stats.SpellPower:       ownerStats[stats.SpellPower] * 0.15,
			stats.FirePower:        ownerStats[stats.FirePower] * 0.15,
			stats.ShadowPower:      ownerStats[stats.ShadowPower] * 0.15,
			stats.SpellPenetration: ownerStats[stats.SpellPenetration],
			stats.MeleeHit:         ownerHitChance * core.MeleeHitRatingPerHitChance,
			stats.SpellHit:         math.Floor(ownerStats[stats.SpellHit] / 12.0 * 17.0),
		}
	}
}
