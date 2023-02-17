package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func applyRaceEffects(agent Agent) {
	character := agent.GetCharacter()

	switch character.Race {
	case proto.Race_RaceEredar:
		/*
			+2% to crit modifier to Party and Character (+ handled in character.go)
			+2% to crit chance and 2% haste for spells (+)
			souls fragments (this is insane to make)
			+2% to armor, and some resistances (+)
		*/
		character.stats[stats.SpellHaste] += 2.0 * HasteRatingPerHastePercent
		character.stats[stats.MeleeCrit] += 2.0 * CritRatingPerCritChance
		character.stats[stats.SpellCrit] += 2.0 * CritRatingPerCritChance

		character.MultiplyStat(stats.Armor, 1.02)
	case proto.Race_RaceGoblin:
		/*
			1% to Haste (+), 1% to Crit Chance (+)
			6% to Crit modifier (this is buf, but uptime in combat is 90-100%) (+, handled in character.go)
			+3% to Dodge and +1% to miss by character (+)
			active: rocket with CD 60 second and damage == 2*AP/SPD (-)
		*/
		character.stats[stats.MeleeHaste] += 1.0 * HasteRatingPerHastePercent
		character.stats[stats.SpellHaste] += 1.0 * HasteRatingPerHastePercent
		character.stats[stats.MeleeCrit] += 1.0 * CritRatingPerCritChance
		character.stats[stats.SpellCrit] += 1.0 * CritRatingPerCritChance

		character.stats[stats.Dodge] = 3.0 * DodgeRatingPerDodgeChance

		character.PseudoStats.ReducedPhysicalHitTakenChance += 0.01
		character.PseudoStats.ReducedArcaneHitTakenChance += 0.01
		character.PseudoStats.ReducedFireHitTakenChance += 0.01
		character.PseudoStats.ReducedFrostHitTakenChance += 0.01
		character.PseudoStats.ReducedNatureHitTakenChance += 0.01
		character.PseudoStats.ReducedShadowHitTakenChance += 0.01

	case proto.Race_RaceZandalar:
		/*
			1% to AP and SPD for group and character (+)
			Various Loa bufs (?)
			4% Crit Modifier (+, handled in character.go)
			+5% Incoming heal and 2% Dodge (+)
			active: +40% HP in 10 seconds,
				increase Haste (both) to 10%  and crit modifier to 12% in 10 seconds (++)

		*/
		character.MultiplyStat(stats.AttackPower, 1.01)
		character.MultiplyStat(stats.SpellPower, 1.01)
		character.PseudoStats.HealingTakenMultiplier *= 1.05
		character.stats[stats.Dodge] = 2.0 * DodgeRatingPerDodgeChance

		actionID := ActionID{SpellID: 319326}
		battleRegenerationAura := character.RegisterAura(Aura{
			Label:    "Battle regeneration",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(1.1)
				character.MultiplyAttackSpeed(sim, 1.1)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(1 / 1.1)
				character.MultiplyAttackSpeed(sim, 1/1.1)
			},
		})

		battleRegenerationSpell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
			},

			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				battleRegenerationAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: battleRegenerationSpell,
			Type:  CooldownTypeDPS,
		})
	case proto.Race_RaceNaga:

	case proto.Race_RaceDwarfOfBlackIron:
		// this is a active bonus aura, but uptime in combat is 95-99%
		character.PseudoStats.DamageDealtMultiplier *= 1.03

		character.stats[stats.MeleeCrit] += 1.0 * CritRatingPerCritChance
		character.stats[stats.SpellCrit] += 1.0 * CritRatingPerCritChance

		actionID := ActionID{SpellID: 316162} // ability_rhyolith_lavapool

		slagEruptionAura := character.RegisterAura(Aura{
			Label:    "Slag eruption",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.DamageDealtMultiplier *= 1.06
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.DamageDealtMultiplier /= 1.06
			},
		})

		slagEruptionSpell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 30,
				},
			},

			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				slagEruptionAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: slagEruptionSpell,
			Type:  CooldownTypeDPS,
		})
	case proto.Race_RaceSindorei:
		character.PseudoStats.ReducedArcaneHitTakenChance += 0.02
		character.PseudoStats.ReducedFireHitTakenChance += 0.02
		character.PseudoStats.ReducedFrostHitTakenChance += 0.02
		character.PseudoStats.ReducedNatureHitTakenChance += 0.02
		character.PseudoStats.ReducedShadowHitTakenChance += 0.02

		var actionID ActionID

		var resourceMetrics *ResourceMetrics = nil
		if resourceMetrics == nil {
			if character.HasRunicPowerBar() {
				actionID = ActionID{SpellID: 50613}
				resourceMetrics = character.NewRunicPowerMetrics(actionID)
			} else if character.HasEnergyBar() {
				actionID = ActionID{SpellID: 25046}
				resourceMetrics = character.NewEnergyMetrics(actionID)
			} else if character.HasManaBar() {
				actionID = ActionID{SpellID: 28730}
				resourceMetrics = character.NewManaMetrics(actionID)
			}
		}

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, spell *Spell) {
				if spell.Unit.HasRunicPowerBar() {
					spell.Unit.AddRunicPower(sim, 15.0, resourceMetrics)
				} else if spell.Unit.HasEnergyBar() {
					spell.Unit.AddEnergy(sim, 15.0, resourceMetrics)
				} else if spell.Unit.HasManaBar() {
					spell.Unit.AddMana(sim, spell.Unit.MaxMana()*0.06, resourceMetrics, false)
				}
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell:    spell,
			Type:     CooldownTypeDPS,
			Priority: CooldownPriorityLow,
		})
	case proto.Race_RaceDraenei:
		character.PseudoStats.ReducedShadowHitTakenChance += 0.02
		// TODO: Gift of the naaru for healers
	case proto.Race_RaceDwarf:
		character.PseudoStats.ReducedFrostHitTakenChance += 0.02

		// Gun specialization (+1% ranged crit when using a gun).
		if character.Equip[proto.ItemSlot_ItemSlotRanged].RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeGun {
			character.AddBonusRangedCritRating(1 * CritRatingPerCritChance)
		}

		applyWeaponSpecialization(
			character,
			5*ExpertisePerQuarterPercentReduction,
			[]proto.WeaponType{proto.WeaponType_WeaponTypeMace})

		actionID := ActionID{SpellID: 20594}

		statDep := character.NewDynamicMultiplyStat(stats.Armor, 1.1)
		stoneFormAura := character.NewTemporaryStatsAuraWrapped("Stoneform", actionID, stats.Stats{}, time.Second*8, func(aura *Aura) {
			aura.ApplyOnGain(func(aura *Aura, sim *Simulation) {
				aura.Unit.EnableDynamicStatDep(sim, statDep)
			})
			aura.ApplyOnExpire(func(aura *Aura, sim *Simulation) {
				aura.Unit.DisableDynamicStatDep(sim, statDep)
			})
		})

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
				stoneFormAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})
	case proto.Race_RaceGnome:
		character.PseudoStats.ReducedArcaneHitTakenChance += 0.02
		character.MultiplyStat(stats.Intellect, 1.05)
	case proto.Race_RaceHuman:
		character.MultiplyStat(stats.Spirit, 1.03)
		applyWeaponSpecialization(
			character,
			3*ExpertisePerQuarterPercentReduction,
			[]proto.WeaponType{proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeSword})
	case proto.Race_RaceNightElf:
		character.PseudoStats.ReducedNatureHitTakenChance += 0.02
		character.PseudoStats.ReducedPhysicalHitTakenChance += 0.02
		// TODO: Shadowmeld?
	case proto.Race_RaceOrc:
		character.stats[stats.Expertise] += 2.0 * ExpertisePerQuarterPercentReduction
		character.MultiplyStat(stats.Strength, 1.02)

		// Blood Fury
		actionID := ActionID{SpellID: 316373}

		bloodFuryAura := character.RegisterAura(Aura{
			Label:    "Blood Fury",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
			},
		})

		bloodFureSpell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
			},

			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				bloodFuryAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: bloodFureSpell,
			Type:  CooldownTypeDPS,
		})

	case proto.Race_RaceTauren:
		character.PseudoStats.ReducedNatureHitTakenChance += 0.02
		character.AddStat(stats.Health, character.GetBaseStats()[stats.Health]*0.05)
	case proto.Race_RaceTroll:
		// Bow specialization (+1% ranged crit when using a bow).
		if character.Equip[proto.ItemSlot_ItemSlotRanged].RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeBow {
			character.AddBonusRangedCritRating(1 * CritRatingPerCritChance)
		}

		// Beast Slaying (+5% damage to beasts)
		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.DamageDealtMultiplier *= 1.05
		}

		// Berserking
		actionID := ActionID{SpellID: 26297}

		berserkingAura := character.RegisterAura(Aura{
			Label:    "Berserking (Troll)",
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
			character.OnSpellRegistered(func(spell *Spell) {
				if spell.ProcMask.Matches(ProcMaskMeleeMH) {
					spell.BonusExpertiseRating += expertiseBonus
				}
			})
		}
		if oh {
			character.OnSpellRegistered(func(spell *Spell) {
				if spell.ProcMask.Matches(ProcMaskMeleeOH) {
					spell.BonusExpertiseRating += expertiseBonus
				}
			})
		}
	}
}
