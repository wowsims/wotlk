package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func applyRaceEffects(agent Agent) {
	character := agent.GetCharacter()

	switch character.Race {
	case proto.Race_RaceBloodElf:
		character.PseudoStats.ReducedArcaneHitTakenChance += 0.02
		character.PseudoStats.ReducedFireHitTakenChance += 0.02
		character.PseudoStats.ReducedFrostHitTakenChance += 0.02
		character.PseudoStats.ReducedNatureHitTakenChance += 0.02
		character.PseudoStats.ReducedShadowHitTakenChance += 0.02
		// TODO: Add major cooldown: arcane torrent
	case proto.Race_RaceDraenei:
		character.PseudoStats.ReducedShadowHitTakenChance += 0.02
		// TODO: Gift of the naaru for healers
	case proto.Race_RaceDwarf:
		character.PseudoStats.ReducedFrostHitTakenChance += 0.02

		// Gun specialization (+1% ranged crit when using a gun).
		if weapon := character.Equip[proto.ItemSlot_ItemSlotRanged]; weapon.ID != 0 {
			if weapon.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeGun {
				character.PseudoStats.BonusRangedCritRating += 1 * CritRatingPerCritChance
			}
		}

		applyWeaponSpecialization(
			character,
			5*ExpertisePerQuarterPercentReduction,
			[]proto.WeaponType{proto.WeaponType_WeaponTypeMace})

		// TODO: Stoneform
	case proto.Race_RaceGnome:
		character.PseudoStats.ReducedArcaneHitTakenChance += 0.02
		character.AddStatDependency(stats.Intellect, stats.Intellect, 1.0+0.05)
	case proto.Race_RaceHuman:
		character.AddStatDependency(stats.Spirit, stats.Spirit, 1.0+0.03)
		applyWeaponSpecialization(
			character,
			3*ExpertisePerQuarterPercentReduction,
			[]proto.WeaponType{proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeSword})
	case proto.Race_RaceNightElf:
		character.PseudoStats.ReducedNatureHitTakenChance += 0.02
		character.PseudoStats.ReducedPhysicalHitTakenChance += 0.02
		// TODO: Shadowmeld?
	case proto.Race_RaceOrc:
		// Command (Pet damage +5%)
		if len(character.Pets) > 0 {
			for _, petAgent := range character.Pets {
				pet := petAgent.GetPet()
				pet.PseudoStats.DamageDealtMultiplier *= 1.05
			}
		}

		// Blood Fury
		actionID := ActionID{SpellID: 33697}
		apBonus := float64(character.Level)*4 + 2
		spBonus := float64(character.Level)*2 + 3
		bloodFuryAura := character.NewTemporaryStatsAura("Blood Fury", actionID, stats.Stats{stats.AttackPower: apBonus, stats.RangedAttackPower: apBonus, stats.SpellPower: spBonus}, time.Second*15)

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				bloodFuryAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})

		// Axe specialization
		applyWeaponSpecialization(
			character,
			5*ExpertisePerQuarterPercentReduction,
			[]proto.WeaponType{proto.WeaponType_WeaponTypeAxe})
	case proto.Race_RaceTauren:
		character.PseudoStats.ReducedNatureHitTakenChance += 0.02
		character.AddStat(stats.Health, character.GetBaseStats()[stats.Health]*0.05)
	case proto.Race_RaceTroll:
		// Bow specialization (+1% ranged crit when using a bow).
		if weapon := character.Equip[proto.ItemSlot_ItemSlotRanged]; weapon.ID != 0 {
			if weapon.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeBow {
				character.PseudoStats.BonusRangedCritRating += 1 * CritRatingPerCritChance
			}
		}

		// Beast Slaying (+5% damage to beasts)
		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.DamageDealtMultiplier *= 1.05
		}

		// Berserking
		actionID := ActionID{SpellID: 26297}

		berserkingAura := character.RegisterAura(Aura{
			Label:    "Berserking",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(1.2)
				character.MultiplyAttackSpeed(sim, 1.2)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(1 / 1.2)
				character.MultiplyAttackSpeed(sim, 1/1.2)
			},
		})

		berserkingSpell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				berserkingAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: berserkingSpell,
			Type:  CooldownTypeDPS,
		})
	case proto.Race_RaceUndead:
		character.PseudoStats.ReducedShadowHitTakenChance += 0.02
	}
}

func applyWeaponSpecialization(character *Character, expertiseBonus float64, weaponTypes []proto.WeaponType) {
	mh := false
	oh := false
	isDW := false
	if weapon := character.Equip[proto.ItemSlot_ItemSlotMainHand]; weapon.ID != 0 {
		for _, wt := range weaponTypes {
			if weapon.WeaponType == wt {
				mh = true
			}
		}
	}
	if weapon := character.Equip[proto.ItemSlot_ItemSlotOffHand]; weapon.ID != 0 && weapon.WeaponType != proto.WeaponType_WeaponTypeShield {
		isDW = true
		for _, wt := range weaponTypes {
			if weapon.WeaponType == wt {
				oh = true
			}
		}
	}

	if mh && (oh || !isDW) {
		character.AddStat(stats.Expertise, expertiseBonus)
	} else {
		if mh {
			character.PseudoStats.BonusMHExpertiseRating += expertiseBonus
		}
		if oh {
			character.PseudoStats.BonusOHExpertiseRating += expertiseBonus
		}
	}
}
