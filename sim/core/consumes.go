package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Registers all consume-related effects to the Agent.
// TODO: Classic Consumes
func applyConsumeEffects(agent Agent) {
	character := agent.GetCharacter()
	consumes := character.Consumes
	if consumes == nil {
		return
	}

	// if consumes.Flask != proto.Flask_FlaskUnknown {
	// 	switch consumes.Flask {
	// 	case proto.Flask_FlaskOfTheFrostWyrm:
	// 		character.AddStats(stats.Stats{
	// 			stats.SpellPower: 125,
	// 		})
	// 		if character.HasProfession(proto.Profession_Alchemy) {
	// 			character.AddStats(stats.Stats{
	// 				stats.SpellPower: 47,
	// 			})
	// 		}
	// 	case proto.Flask_FlaskOfEndlessRage:
	// 		character.AddStats(stats.Stats{
	// 			stats.AttackPower:       180,
	// 			stats.RangedAttackPower: 180,
	// 		})
	// 		if character.HasProfession(proto.Profession_Alchemy) {
	// 			character.AddStats(stats.Stats{
	// 				stats.AttackPower:       80,
	// 				stats.RangedAttackPower: 80,
	// 			})
	// 		}
	// 	case proto.Flask_FlaskOfPureMojo:
	// 		character.AddStats(stats.Stats{
	// 			stats.MP5: 45,
	// 		})
	// 		if character.HasProfession(proto.Profession_Alchemy) {
	// 			character.AddStats(stats.Stats{
	// 				stats.MP5: 20,
	// 			})
	// 		}
	// 	case proto.Flask_FlaskOfStoneblood:
	// 		character.AddStats(stats.Stats{
	// 			stats.Health: 1300,
	// 		})
	// 		if character.HasProfession(proto.Profession_Alchemy) {
	// 			character.AddStats(stats.Stats{
	// 				stats.Health: 650,
	// 			})
	// 		}
	// 	case proto.Flask_LesserFlaskOfToughness:
	// 		character.AddStats(stats.Stats{
	// 			stats.Resilience: 50,
	// 		})
	// 		if character.HasProfession(proto.Profession_Alchemy) {
	// 			character.AddStats(stats.Stats{
	// 				stats.Resilience: 82,
	// 			})
	// 		}
	// 	case proto.Flask_LesserFlaskOfResistance:
	// 		character.AddStats(stats.Stats{
	// 			stats.ArcaneResistance: 50,
	// 			stats.FireResistance:   50,
	// 			stats.FrostResistance:  50,
	// 			stats.NatureResistance: 50,
	// 			stats.ShadowResistance: 50,
	// 		})
	// 		if character.HasProfession(proto.Profession_Alchemy) {
	// 			character.AddStats(stats.Stats{
	// 				stats.ArcaneResistance: 40,
	// 				stats.FireResistance:   40,
	// 				stats.FrostResistance:  40,
	// 				stats.NatureResistance: 40,
	// 				stats.ShadowResistance: 40,
	// 			})
	// 		}
	// 	case proto.Flask_FlaskOfBlindingLight:
	// 		character.OnSpellRegistered(func(spell *Spell) {
	// 			if spell.SpellSchool.Matches(SpellSchoolArcane | SpellSchoolHoly | SpellSchoolNature) {
	// 				spell.BonusSpellPower += 80
	// 			}
	// 		})
	// 	case proto.Flask_FlaskOfMightyRestoration:
	// 		character.AddStats(stats.Stats{
	// 			stats.MP5: 25,
	// 		})
	// 	case proto.Flask_FlaskOfPureDeath:
	// 		character.OnSpellRegistered(func(spell *Spell) {
	// 			if spell.SpellSchool.Matches(SpellSchoolFire | SpellSchoolFrost | SpellSchoolShadow) {
	// 				spell.BonusSpellPower += 80
	// 			}
	// 		})
	// 	case proto.Flask_FlaskOfRelentlessAssault:
	// 		character.AddStats(stats.Stats{
	// 			stats.AttackPower:       120,
	// 			stats.RangedAttackPower: 120,
	// 		})
	// 	case proto.Flask_FlaskOfSupremePower:
	// 		character.AddStats(stats.Stats{
	// 			stats.SpellPower: 70,
	// 		})
	// 	case proto.Flask_FlaskOfFortification:
	// 		character.AddStats(stats.Stats{
	// 			stats.Health:  500,
	// 			stats.Defense: 10,
	// 		})
	// 	case proto.Flask_FlaskOfChromaticWonder:
	// 		character.AddStats(stats.Stats{
	// 			stats.Stamina:          18,
	// 			stats.Strength:         18,
	// 			stats.Agility:          18,
	// 			stats.Intellect:        18,
	// 			stats.Spirit:           18,
	// 			stats.ArcaneResistance: 35,
	// 			stats.FireResistance:   35,
	// 			stats.FrostResistance:  35,
	// 			stats.NatureResistance: 35,
	// 			stats.ShadowResistance: 35,
	// 		})
	// 	}
	// } else {
	// 	switch consumes.BattleElixir {
	// 	case proto.BattleElixir_ElixirOfAccuracy:
	// 		character.AddStats(stats.Stats{
	// 			stats.MeleeHit: 45,
	// 			stats.SpellHit: 45,
	// 		})
	// 	case proto.BattleElixir_SpellpowerElixir:
	// 		character.AddStats(stats.Stats{
	// 			stats.SpellPower: 58,
	// 		})
	// 	case proto.BattleElixir_WrathElixir:
	// 		character.AddStats(stats.Stats{
	// 			stats.AttackPower:       90,
	// 			stats.RangedAttackPower: 90,
	// 		})
	// 	case proto.BattleElixir_AdeptsElixir:
	// 		character.AddStats(stats.Stats{
	// 			stats.SpellCrit:  24,
	// 			stats.SpellPower: 24,
	// 		})
	// 	case proto.BattleElixir_ElixirOfDemonslaying:
	// 		if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
	// 			character.PseudoStats.MobTypeAttackPower += 265
	// 		}
	// 	case proto.BattleElixir_ElixirOfMajorAgility:
	// 		character.AddStats(stats.Stats{
	// 			stats.Agility:   35,
	// 			stats.MeleeCrit: 20,
	// 		})
	// 	case proto.BattleElixir_ElixirOfMajorStrength:
	// 		character.AddStats(stats.Stats{
	// 			stats.Strength: 35,
	// 		})
	// 	case proto.BattleElixir_ElixirOfMastery:
	// 		character.AddStats(stats.Stats{
	// 			stats.Stamina:   15,
	// 			stats.Strength:  15,
	// 			stats.Agility:   15,
	// 			stats.Intellect: 15,
	// 			stats.Spirit:    15,
	// 		})
	// 	case proto.BattleElixir_ElixirOfTheMongoose:
	// 		character.AddStats(stats.Stats{
	// 			stats.Agility:   25,
	// 			stats.MeleeCrit: 28,
	// 		})
	// 	case proto.BattleElixir_FelStrengthElixir:
	// 		character.AddStats(stats.Stats{
	// 			stats.AttackPower:       90,
	// 			stats.RangedAttackPower: 90,
	// 			stats.Stamina:           -10,
	// 		})
	// 	case proto.BattleElixir_GreaterArcaneElixir:
	// 		character.AddStats(stats.Stats{
	// 			stats.SpellPower: 35,
	// 		})
	// 	}

	// 	switch consumes.GuardianElixir {
	// 	case proto.GuardianElixir_ElixirOfMightyDefense:
	// 		character.AddStats(stats.Stats{
	// 			stats.Defense: 45,
	// 		})
	// 	case proto.GuardianElixir_ElixirOfMightyFortitude:
	// 		character.AddStats(stats.Stats{
	// 			stats.Health: 350,
	// 		})
	// 	case proto.GuardianElixir_ElixirOfMightyMageblood:
	// 		character.AddStats(stats.Stats{
	// 			stats.MP5: 30,
	// 		})
	// 	case proto.GuardianElixir_ElixirOfMightyThoughts:
	// 		character.AddStats(stats.Stats{
	// 			stats.Intellect: 45,
	// 		})
	// 	case proto.GuardianElixir_ElixirOfProtection:
	// 		character.AddStats(stats.Stats{
	// 			stats.Armor: 800,
	// 		})
	// 		if character.HasProfession(proto.Profession_Alchemy) {
	// 			character.AddStats(stats.Stats{
	// 				stats.Armor: 280,
	// 			})
	// 		}
	// 	case proto.GuardianElixir_ElixirOfSpirit:
	// 		character.AddStats(stats.Stats{
	// 			stats.Spirit: 50,
	// 		})
	// 	case proto.GuardianElixir_ElixirOfDraenicWisdom:
	// 		character.AddStats(stats.Stats{
	// 			stats.Intellect: 30,
	// 			stats.Spirit:    30,
	// 		})
	// 	case proto.GuardianElixir_ElixirOfIronskin:
	// 		character.AddStats(stats.Stats{
	// 			stats.Resilience: 30,
	// 		})
	// 	case proto.GuardianElixir_ElixirOfMajorDefense:
	// 		character.AddStats(stats.Stats{
	// 			stats.Armor: 550,
	// 		})
	// 	case proto.GuardianElixir_ElixirOfMajorFortitude:
	// 		character.AddStats(stats.Stats{
	// 			stats.Health: 250,
	// 		})
	// 	case proto.GuardianElixir_ElixirOfMajorMageblood:
	// 		character.AddStats(stats.Stats{
	// 			stats.MP5: 16,
	// 		})
	// 	case proto.GuardianElixir_GiftOfArthas:
	// 		character.AddStats(stats.Stats{
	// 			stats.ShadowResistance: 10,
	// 		})

	// 		debuffAuras := (&character.Unit).NewEnemyAuraArray(GiftOfArthasAura)

	// 		actionID := ActionID{SpellID: 11374}
	// 		goaProc := character.RegisterSpell(SpellConfig{
	// 			ActionID:    actionID,
	// 			SpellSchool: SpellSchoolNature,
	// 			ProcMask:    ProcMaskEmpty,

	// 			ThreatMultiplier: 1,
	// 			FlatThreatBonus:  90,

	// 			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
	// 				debuffAuras.Get(target).Activate(sim)
	// 				spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
	// 			},
	// 		})

	// 		character.RegisterAura(Aura{
	// 			Label:    "Gift of Arthas",
	// 			Duration: NeverExpires,
	// 			OnReset: func(aura *Aura, sim *Simulation) {
	// 				aura.Activate(sim)
	// 			},
	// 			OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
	// 				if result.Landed() &&
	// 					spell.SpellSchool == SpellSchoolPhysical &&
	// 					sim.RandomFloat("Gift of Arthas") < 0.3 {
	// 					goaProc.Cast(sim, spell.Unit)
	// 				}
	// 			},
	// 		})
	// 	}
	// }

	// switch consumes.Food {
	// case proto.Food_FoodFishFeast:
	// 	character.AddStats(stats.Stats{
	// 		stats.AttackPower:       80,
	// 		stats.RangedAttackPower: 80,
	// 		stats.SpellPower:        46,
	// 		stats.Stamina:           40,
	// 	})
	// case proto.Food_FoodGreatFeast:
	// 	character.AddStats(stats.Stats{
	// 		stats.AttackPower:       60,
	// 		stats.RangedAttackPower: 60,
	// 		stats.SpellPower:        35,
	// 		stats.Stamina:           30,
	// 	})
	// case proto.Food_FoodBlackenedDragonfin:
	// 	character.AddStats(stats.Stats{
	// 		stats.Agility: 40,
	// 		stats.Stamina: 40,
	// 	})
	// case proto.Food_FoodHeartyRhino:
	// 	character.AddStats(stats.Stats{
	// 		stats.ArmorPenetration: 40,
	// 		stats.Stamina:          40,
	// 	})
	// case proto.Food_FoodMegaMammothMeal:
	// 	character.AddStats(stats.Stats{
	// 		stats.AttackPower:       80,
	// 		stats.RangedAttackPower: 80,
	// 		stats.Stamina:           40,
	// 	})
	// case proto.Food_FoodSpicedWormBurger:
	// 	character.AddStats(stats.Stats{
	// 		stats.MeleeCrit: 40,
	// 		stats.SpellCrit: 40,
	// 		stats.Stamina:   40,
	// 	})
	// case proto.Food_FoodRhinoliciousWormsteak:
	// 	character.AddStats(stats.Stats{
	// 		stats.Expertise: 40,
	// 		stats.Stamina:   40,
	// 	})
	// case proto.Food_FoodImperialMantaSteak:
	// 	character.AddStats(stats.Stats{
	// 		stats.MeleeHaste: 40,
	// 		stats.SpellHaste: 40,
	// 		stats.Stamina:    40,
	// 	})
	// case proto.Food_FoodSnapperExtreme:
	// 	character.AddStats(stats.Stats{
	// 		stats.MeleeHit: 40,
	// 		stats.SpellHit: 40,
	// 		stats.Stamina:  40,
	// 	})
	// case proto.Food_FoodMightyRhinoDogs:
	// 	character.AddStats(stats.Stats{
	// 		stats.MP5:     16,
	// 		stats.Stamina: 40,
	// 	})
	// case proto.Food_FoodFirecrackerSalmon:
	// 	character.AddStats(stats.Stats{
	// 		stats.SpellPower: 46,
	// 		stats.Stamina:    40,
	// 	})
	// case proto.Food_FoodCuttlesteak:
	// 	character.AddStats(stats.Stats{
	// 		stats.Spirit:  40,
	// 		stats.Stamina: 40,
	// 	})
	// case proto.Food_FoodDragonfinFilet:
	// 	character.AddStats(stats.Stats{
	// 		stats.Strength: 40,
	// 		stats.Stamina:  40,
	// 	})
	// case proto.Food_FoodBlackenedBasilisk:
	// 	character.AddStats(stats.Stats{
	// 		stats.SpellPower: 23,
	// 		stats.Spirit:     20,
	// 	})
	// case proto.Food_FoodGrilledMudfish:
	// 	character.AddStats(stats.Stats{
	// 		stats.Agility: 20,
	// 		stats.Spirit:  20,
	// 	})
	// case proto.Food_FoodRavagerDog:
	// 	character.AddStats(stats.Stats{
	// 		stats.AttackPower:       40,
	// 		stats.RangedAttackPower: 40,
	// 		stats.Spirit:            20,
	// 	})
	// case proto.Food_FoodRoastedClefthoof:
	// 	character.AddStats(stats.Stats{
	// 		stats.Strength: 20,
	// 		stats.Spirit:   20,
	// 	})
	// case proto.Food_FoodSkullfishSoup:
	// 	character.AddStats(stats.Stats{
	// 		stats.SpellCrit: 20,
	// 		stats.Spirit:    20,
	// 	})
	// case proto.Food_FoodSpicyHotTalbuk:
	// 	character.AddStats(stats.Stats{
	// 		stats.MeleeHit: 20,
	// 		stats.Spirit:   20,
	// 	})
	// case proto.Food_FoodFishermansFeast:
	// 	character.AddStats(stats.Stats{
	// 		stats.Stamina: 30,
	// 		stats.Spirit:  20,
	// 	})
	// }

	// registerPotionCD(agent, consumes)
	// registerConjuredCD(agent, consumes)
	// registerExplosivesCD(agent, consumes)
}

var PotionAuraTag = "Potion"

func makePotionActivation(potionType proto.Potions, character *Character, potionCD *Timer) MajorCooldown {
	mcd := makePotionActivationInternal(potionType, character, potionCD)
	if mcd.Spell != nil {
		// Mark as 'Encounter Only' so that users are forced to select the generic Potion
		// placeholder action instead of specific potion spells, in APL prepull. This
		// prevents a mismatch between Consumes and Rotation settings.
		mcd.Spell.Flags |= SpellFlagEncounterOnly | SpellFlagPotion
		oldApplyEffects := mcd.Spell.ApplyEffects
		mcd.Spell.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			oldApplyEffects(sim, target, spell)
			if sim.CurrentTime < 0 {
				if potionType == proto.Potions_IndestructiblePotion {
					potionCD.Set(sim.CurrentTime + 2*time.Minute)
				} else {
					potionCD.Set(sim.CurrentTime + time.Minute)
				}
				character.UpdateMajorCooldowns()
			}
		}
	}
	return mcd
}

func makePotionActivationInternal(potionType proto.Potions, character *Character, potionCD *Timer) MajorCooldown {
	potionCast := CastConfig{
		CD: Cooldown{
			Timer:    potionCD,
			Duration: time.Minute * 60, // Infinite CD
		},
	}

	if potionType == proto.Potions_RunicHealingPotion || potionType == proto.Potions_RunicHealingInjector {
		itemId := map[proto.Potions]int32{
			proto.Potions_RunicHealingPotion:   33447,
			proto.Potions_RunicHealingInjector: 41166,
		}[potionType]
		actionID := ActionID{ItemID: itemId}
		healthMetrics := character.NewHealthMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeSurvival,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					healthGain := sim.RollWithLabel(2700, 4500, "RunicHealingPotion")

					character.GainHealth(sim, healthGain*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_RunicManaPotion || potionType == proto.Potions_RunicManaInjector {
		itemId := map[proto.Potions]int32{
			proto.Potions_RunicManaPotion:   33448,
			proto.Potions_RunicManaInjector: 42545,
		}[potionType]
		actionID := ActionID{ItemID: itemId}
		manaMetrics := character.NewManaMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				manaGain := 4400.0

				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= manaGain
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					manaGain := sim.RollWithLabel(4200, 4400, "RunicManaPotion")
					character.AddMana(sim, manaGain, manaMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_IndestructiblePotion {
		actionID := ActionID{ItemID: 40093}
		aura := character.NewTemporaryStatsAura("Indestructible Potion", actionID, stats.Stats{stats.Armor: 3500}, time.Minute*2)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_PotionOfSpeed {
		actionID := ActionID{ItemID: 40211}
		aura := character.NewTemporaryStatsAura("Potion of Speed", actionID, stats.Stats{stats.MeleeHaste: 500, stats.SpellHaste: 500}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_PotionOfWildMagic {
		actionID := ActionID{ItemID: 40212}
		aura := character.NewTemporaryStatsAura("Potion of Wild Magic", actionID, stats.Stats{stats.SpellPower: 200, stats.SpellCrit: 200, stats.MeleeCrit: 200}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_DestructionPotion {
		actionID := ActionID{ItemID: 22839}
		aura := character.NewTemporaryStatsAura("Destruction Potion", actionID, stats.Stats{stats.SpellPower: 120, stats.SpellCrit: 2 * CritRatingPerCritChance}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_SuperManaPotion {
		actionID := ActionID{ItemID: 22832}
		manaMetrics := character.NewManaMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= 3000
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					// Restores 1800 to 3000 mana. (2 Min Cooldown)
					manaGain := sim.RollWithLabel(1800, 3000, "super mana")
					character.AddMana(sim, manaGain, manaMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_HastePotion {
		actionID := ActionID{ItemID: 22838}
		aura := character.NewTemporaryStatsAura("Haste Potion", actionID, stats.Stats{stats.MeleeHaste: 400}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_MightyRagePotion {
		actionID := ActionID{ItemID: 13442}
		aura := character.NewTemporaryStatsAura("Mighty Rage Potion", actionID, stats.Stats{stats.Strength: 60}, time.Second*15)
		rageMetrics := character.NewRageMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				if character.Class == proto.Class_ClassWarrior {
					return character.CurrentRage() < 25
				}
				return true
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
					if character.Class == proto.Class_ClassWarrior {
						bonusRage := sim.RollWithLabel(45, 75, "Mighty Rage Potion")
						character.AddRage(sim, bonusRage, rageMetrics)
					}
				},
			}),
		}
	} else if potionType == proto.Potions_FelManaPotion {
		actionID := ActionID{ItemID: 31677}

		// Restores 3200 mana over 24 seconds.
		manaGain := 3200.0
		mp5 := manaGain / 24 * 5

		buffAura := character.NewTemporaryStatsAura("Fel Mana Potion", actionID, stats.Stats{stats.MP5: mp5}, time.Second*24)
		debuffAura := character.NewTemporaryStatsAura("Fel Mana Potion Debuff", ActionID{SpellID: 38927}, stats.Stats{stats.SpellPower: -25}, time.Minute*15)

		return MajorCooldown{
			Type: CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have low enough mana. The potion takes effect over 24
				// seconds so we can pop it a little earlier than the full value.
				return character.MaxMana()-character.CurrentMana() >= 2000
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					buffAura.Activate(sim)
					debuffAura.Activate(sim)
					debuffAura.Refresh(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_InsaneStrengthPotion {
		actionID := ActionID{ItemID: 22828}
		aura := character.NewTemporaryStatsAura("Insane Strength Potion", actionID, stats.Stats{stats.Strength: 120, stats.Defense: -75}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_IronshieldPotion {
		actionID := ActionID{ItemID: 22849}
		aura := character.NewTemporaryStatsAura("Ironshield Potion", actionID, stats.Stats{stats.Armor: 2500}, time.Minute*2)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_HeroicPotion {
		actionID := ActionID{ItemID: 22837}
		aura := character.NewTemporaryStatsAura("Heroic Potion", actionID, stats.Stats{stats.Strength: 70, stats.Health: 700}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else {
		return MajorCooldown{}
	}
}

var ConjuredAuraTag = "Conjured"
