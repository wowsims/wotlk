package core

import (
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func applyRaceEffects(agent Agent) {
	character := agent.GetCharacter()

	switch character.Race {
	case proto.Race_RaceBloodElf:
		character.AddStat(stats.ArcaneResistance, 5)
		character.AddStat(stats.FireResistance, 5)
		character.AddStat(stats.FrostResistance, 5)
		character.AddStat(stats.NatureResistance, 5)
		character.AddStat(stats.ShadowResistance, 5)
		// TODO: Add major cooldown: arcane torrent
	case proto.Race_RaceDraenei:
		character.AddStat(stats.ShadowResistance, 10)
	case proto.Race_RaceDwarf:
		character.AddStat(stats.FrostResistance, 10)

		// Gun specialization (+1% ranged crit when using a gun).
		if weapon := character.Equip[proto.ItemSlot_ItemSlotRanged]; weapon.ID != 0 {
			if weapon.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeGun {
				character.PseudoStats.BonusRangedCritRating += 1 * MeleeCritRatingPerCritChance
			}
		}
	case proto.Race_RaceGnome:
		character.AddStat(stats.ArcaneResistance, 10)

		character.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.Intellect,
			Modifier: func(intellect float64, _ float64) float64 {
				return intellect * 1.05
			},
		})
	case proto.Race_RaceHuman:
		character.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Spirit,
			ModifiedStat: stats.Spirit,
			Modifier: func(spirit float64, _ float64) float64 {
				return spirit * 1.1
			},
		})

		const expertiseBonus = 5 * ExpertisePerQuarterPercentReduction
		mh := false
		oh := false
		isDW := false
		if weapon := character.Equip[proto.ItemSlot_ItemSlotMainHand]; weapon.ID != 0 {
			if weapon.WeaponType == proto.WeaponType_WeaponTypeSword || weapon.WeaponType == proto.WeaponType_WeaponTypeMace {
				mh = true
			}
		}
		if weapon := character.Equip[proto.ItemSlot_ItemSlotOffHand]; weapon.ID != 0 && weapon.WeaponType != proto.WeaponType_WeaponTypeShield {
			isDW = true
			if weapon.WeaponType == proto.WeaponType_WeaponTypeSword || weapon.WeaponType == proto.WeaponType_WeaponTypeMace {
				oh = true
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
	case proto.Race_RaceNightElf:
		character.AddStat(stats.NatureResistance, 10)
		character.AddStat(stats.Dodge, DodgeRatingPerDodgeChance*1)
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
		mh := false
		oh := false
		isDW := false
		const expertiseBonus = 5 * ExpertisePerQuarterPercentReduction
		if weapon := character.Equip[proto.ItemSlot_ItemSlotMainHand]; weapon.ID != 0 {
			if weapon.WeaponType == proto.WeaponType_WeaponTypeAxe {
				mh = true
			}
		}
		if weapon := character.Equip[proto.ItemSlot_ItemSlotOffHand]; weapon.ID != 0 {
			isDW = true
			if weapon.WeaponType == proto.WeaponType_WeaponTypeAxe {
				oh = true
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
	case proto.Race_RaceTauren:
		character.AddStat(stats.NatureResistance, 10)
		character.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Health,
			ModifiedStat: stats.Health,
			Modifier: func(health float64, _ float64) float64 {
				return health * 1.05
			},
		})
	case proto.Race_RaceTroll10, proto.Race_RaceTroll30:
		// Bow specialization (+1% ranged crit when using a bow).
		if weapon := character.Equip[proto.ItemSlot_ItemSlotRanged]; weapon.ID != 0 {
			if weapon.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeBow {
				character.PseudoStats.BonusRangedCritRating += 1 * MeleeCritRatingPerCritChance
			}
		}

		// Beast Slaying (+5% damage to beasts)
		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.DamageDealtMultiplier *= 1.05
		}

		// Berserking
		hasteBonus := 1.1
		if character.Race == proto.Race_RaceTroll30 {
			hasteBonus = 1.3
		}
		inverseBonus := 1 / hasteBonus

		var resourceType stats.Stat
		var cost float64
		var actionID ActionID
		if character.Class == proto.Class_ClassRogue {
			resourceType = stats.Energy
			cost = 10
			actionID = ActionID{SpellID: 26297}
		} else if character.Class == proto.Class_ClassWarrior {
			resourceType = stats.Rage
			cost = 5
			actionID = ActionID{SpellID: 26296}
		} else {
			resourceType = stats.Mana
			cost = character.BaseMana() * 0.06
			actionID = ActionID{SpellID: 20554}
		}

		berserkingAura := character.RegisterAura(Aura{
			Label:    "Berserking",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(hasteBonus)
				character.MultiplyAttackSpeed(sim, hasteBonus)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(inverseBonus)
				character.MultiplyAttackSpeed(sim, inverseBonus)
			},
		})

		berserkingSpell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,

			ResourceType: resourceType,
			BaseCost:     cost,

			Cast: CastConfig{
				DefaultCast: Cast{
					Cost: cost,
				},
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
			CanActivate: func(sim *Simulation, character *Character) bool {
				if character.Class == proto.Class_ClassRogue {
					return character.CurrentEnergy() >= cost
				} else if character.Class == proto.Class_ClassWarrior {
					return character.CurrentRage() >= cost
				} else {
					return character.CurrentMana() >= cost
				}
			},
		})
	case proto.Race_RaceUndead:
		character.AddStat(stats.ShadowResistance, 10)
	}
}
