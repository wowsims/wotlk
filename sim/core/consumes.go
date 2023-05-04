package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Registers all consume-related effects to the Agent.
func applyConsumeEffects(agent Agent) {
	character := agent.GetCharacter()
	consumes := character.Consumes
	if consumes == nil {
		return
	}

	if consumes.Flask != proto.Flask_FlaskUnknown {
		switch consumes.Flask {
		case proto.Flask_FlaskOfTheFrostWyrm:
			character.AddStats(stats.Stats{
				stats.SpellPower: 125,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.SpellPower: 47,
				})
			}
		case proto.Flask_FlaskOfEndlessRage:
			character.AddStats(stats.Stats{
				stats.AttackPower:       180,
				stats.RangedAttackPower: 180,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.AttackPower:       80,
					stats.RangedAttackPower: 80,
				})
			}
		case proto.Flask_FlaskOfPureMojo:
			character.AddStats(stats.Stats{
				stats.MP5: 45,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.MP5: 20,
				})
			}
		case proto.Flask_FlaskOfStoneblood:
			character.AddStats(stats.Stats{
				stats.Health: 1300,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.Health: 650,
				})
			}
		case proto.Flask_LesserFlaskOfToughness:
			character.AddStats(stats.Stats{
				stats.Resilience: 50,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.Resilience: 82,
				})
			}
		case proto.Flask_LesserFlaskOfResistance:
			character.AddStats(stats.Stats{
				stats.ArcaneResistance: 50,
				stats.FireResistance:   50,
				stats.FrostResistance:  50,
				stats.NatureResistance: 50,
				stats.ShadowResistance: 50,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.ArcaneResistance: 40,
					stats.FireResistance:   40,
					stats.FrostResistance:  40,
					stats.NatureResistance: 40,
					stats.ShadowResistance: 40,
				})
			}
		case proto.Flask_FlaskOfBlindingLight:
			character.OnSpellRegistered(func(spell *Spell) {
				if spell.SpellSchool.Matches(SpellSchoolArcane | SpellSchoolHoly | SpellSchoolNature) {
					spell.BonusSpellPower += 80
				}
			})
		case proto.Flask_FlaskOfMightyRestoration:
			character.AddStats(stats.Stats{
				stats.MP5: 25,
			})
		case proto.Flask_FlaskOfPureDeath:
			character.OnSpellRegistered(func(spell *Spell) {
				if spell.SpellSchool.Matches(SpellSchoolFire | SpellSchoolFrost | SpellSchoolShadow) {
					spell.BonusSpellPower += 80
				}
			})
		case proto.Flask_FlaskOfRelentlessAssault:
			character.AddStats(stats.Stats{
				stats.AttackPower:       120,
				stats.RangedAttackPower: 120,
			})
		case proto.Flask_FlaskOfSupremePower:
			character.AddStats(stats.Stats{
				stats.SpellPower: 70,
			})
		case proto.Flask_FlaskOfFortification:
			character.AddStats(stats.Stats{
				stats.Health:  500,
				stats.Defense: 10,
			})
		case proto.Flask_FlaskOfChromaticWonder:
			character.AddStats(stats.Stats{
				stats.Stamina:          18,
				stats.Strength:         18,
				stats.Agility:          18,
				stats.Intellect:        18,
				stats.Spirit:           18,
				stats.ArcaneResistance: 35,
				stats.FireResistance:   35,
				stats.FrostResistance:  35,
				stats.NatureResistance: 35,
				stats.ShadowResistance: 35,
			})
		}
	} else {
		switch consumes.BattleElixir {
		case proto.BattleElixir_ElixirOfAccuracy:
			character.AddStats(stats.Stats{
				stats.MeleeHit: 45,
				stats.SpellHit: 45,
			})
		case proto.BattleElixir_ElixirOfArmorPiercing:
			character.AddStats(stats.Stats{
				stats.ArmorPenetration: 45,
			})
		case proto.BattleElixir_ElixirOfDeadlyStrikes:
			character.AddStats(stats.Stats{
				stats.MeleeCrit: 45,
				stats.SpellCrit: 45,
			})
		case proto.BattleElixir_ElixirOfExpertise:
			character.AddStats(stats.Stats{
				stats.Expertise: 45,
			})
		case proto.BattleElixir_ElixirOfLightningSpeed:
			character.AddStats(stats.Stats{
				stats.MeleeHaste: 45,
				stats.SpellHaste: 45,
			})
		case proto.BattleElixir_ElixirOfMightyAgility:
			character.AddStats(stats.Stats{
				stats.Agility: 45,
			})
		case proto.BattleElixir_ElixirOfMightyStrength:
			character.AddStats(stats.Stats{
				stats.Strength: 45,
			})
		case proto.BattleElixir_GurusElixir:
			character.AddStats(stats.Stats{
				stats.Agility:   20,
				stats.Strength:  20,
				stats.Stamina:   20,
				stats.Intellect: 20,
				stats.Spirit:    20,
			})
		case proto.BattleElixir_SpellpowerElixir:
			character.AddStats(stats.Stats{
				stats.SpellPower: 58,
			})
		case proto.BattleElixir_WrathElixir:
			character.AddStats(stats.Stats{
				stats.AttackPower:       90,
				stats.RangedAttackPower: 90,
			})
		case proto.BattleElixir_AdeptsElixir:
			character.AddStats(stats.Stats{
				stats.SpellCrit:  24,
				stats.SpellPower: 24,
			})
		case proto.BattleElixir_ElixirOfDemonslaying:
			if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
				character.PseudoStats.MobTypeAttackPower += 265
			}
		case proto.BattleElixir_ElixirOfMajorAgility:
			character.AddStats(stats.Stats{
				stats.Agility:   35,
				stats.MeleeCrit: 20,
			})
		case proto.BattleElixir_ElixirOfMajorStrength:
			character.AddStats(stats.Stats{
				stats.Strength: 35,
			})
		case proto.BattleElixir_ElixirOfMastery:
			character.AddStats(stats.Stats{
				stats.Stamina:   15,
				stats.Strength:  15,
				stats.Agility:   15,
				stats.Intellect: 15,
				stats.Spirit:    15,
			})
		case proto.BattleElixir_ElixirOfTheMongoose:
			character.AddStats(stats.Stats{
				stats.Agility:   25,
				stats.MeleeCrit: 28,
			})
		case proto.BattleElixir_FelStrengthElixir:
			character.AddStats(stats.Stats{
				stats.AttackPower:       90,
				stats.RangedAttackPower: 90,
				stats.Stamina:           -10,
			})
		case proto.BattleElixir_GreaterArcaneElixir:
			character.AddStats(stats.Stats{
				stats.SpellPower: 35,
			})
		}

		switch consumes.GuardianElixir {
		case proto.GuardianElixir_ElixirOfMightyDefense:
			character.AddStats(stats.Stats{
				stats.Defense: 45,
			})
		case proto.GuardianElixir_ElixirOfMightyFortitude:
			character.AddStats(stats.Stats{
				stats.Health: 350,
			})
		case proto.GuardianElixir_ElixirOfMightyMageblood:
			character.AddStats(stats.Stats{
				stats.MP5: 30,
			})
		case proto.GuardianElixir_ElixirOfMightyThoughts:
			character.AddStats(stats.Stats{
				stats.Intellect: 45,
			})
		case proto.GuardianElixir_ElixirOfProtection:
			character.AddStats(stats.Stats{
				stats.Armor: 800,
			})
		case proto.GuardianElixir_ElixirOfSpirit:
			character.AddStats(stats.Stats{
				stats.Spirit: 50,
			})
		case proto.GuardianElixir_ElixirOfDraenicWisdom:
			character.AddStats(stats.Stats{
				stats.Intellect: 30,
				stats.Spirit:    30,
			})
		case proto.GuardianElixir_ElixirOfIronskin:
			character.AddStats(stats.Stats{
				stats.Resilience: 30,
			})
		case proto.GuardianElixir_ElixirOfMajorDefense:
			character.AddStats(stats.Stats{
				stats.Armor: 550,
			})
		case proto.GuardianElixir_ElixirOfMajorFortitude:
			character.AddStats(stats.Stats{
				stats.Health: 250,
			})
		case proto.GuardianElixir_ElixirOfMajorMageblood:
			character.AddStats(stats.Stats{
				stats.MP5: 16,
			})
		case proto.GuardianElixir_GiftOfArthas:
			character.AddStats(stats.Stats{
				stats.ShadowResistance: 10,
			})

			debuffAuras := make([]*Aura, len(character.Env.Encounter.TargetUnits))
			for i, target := range character.Env.Encounter.TargetUnits {
				debuffAuras[i] = GiftOfArthasAura(target)
			}

			actionID := ActionID{SpellID: 11374}
			goaProc := character.RegisterSpell(SpellConfig{
				ActionID:    actionID,
				SpellSchool: SpellSchoolNature,
				ProcMask:    ProcMaskEmpty,

				ThreatMultiplier: 1,
				FlatThreatBonus:  90,

				ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
					debuffAuras[target.Index].Activate(sim)
					spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
				},
			})

			character.RegisterAura(Aura{
				Label:    "Gift of Arthas",
				Duration: NeverExpires,
				OnReset: func(aura *Aura, sim *Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
					if result.Landed() &&
						spell.SpellSchool == SpellSchoolPhysical &&
						sim.RandomFloat("Gift of Arthas") < 0.3 {
						goaProc.Cast(sim, spell.Unit)
					}
				},
			})
		}
	}

	switch consumes.Food {
	case proto.Food_FoodFishFeast:
		character.AddStats(stats.Stats{
			stats.AttackPower:       80,
			stats.RangedAttackPower: 80,
			stats.SpellPower:        46,
			stats.Stamina:           40,
		})
	case proto.Food_FoodGreatFeast:
		character.AddStats(stats.Stats{
			stats.AttackPower:       60,
			stats.RangedAttackPower: 60,
			stats.SpellPower:        35,
			stats.Stamina:           30,
		})
	case proto.Food_FoodBlackenedDragonfin:
		character.AddStats(stats.Stats{
			stats.Agility: 40,
			stats.Stamina: 40,
		})
	case proto.Food_FoodHeartyRhino:
		character.AddStats(stats.Stats{
			stats.ArmorPenetration: 40,
			stats.Stamina:          40,
		})
	case proto.Food_FoodMegaMammothMeal:
		character.AddStats(stats.Stats{
			stats.AttackPower:       80,
			stats.RangedAttackPower: 80,
			stats.Stamina:           40,
		})
	case proto.Food_FoodSpicedWormBurger:
		character.AddStats(stats.Stats{
			stats.MeleeCrit: 40,
			stats.SpellCrit: 40,
			stats.Stamina:   40,
		})
	case proto.Food_FoodRhinoliciousWormsteak:
		character.AddStats(stats.Stats{
			stats.Expertise: 40,
			stats.Stamina:   40,
		})
	case proto.Food_FoodImperialMantaSteak:
		character.AddStats(stats.Stats{
			stats.MeleeHaste: 40,
			stats.SpellHaste: 40,
			stats.Stamina:    40,
		})
	case proto.Food_FoodSnapperExtreme:
		character.AddStats(stats.Stats{
			stats.MeleeHit: 40,
			stats.SpellHit: 40,
			stats.Stamina:  40,
		})
	case proto.Food_FoodMightyRhinoDogs:
		character.AddStats(stats.Stats{
			stats.MP5:     16,
			stats.Stamina: 40,
		})
	case proto.Food_FoodFirecrackerSalmon:
		character.AddStats(stats.Stats{
			stats.SpellPower: 46,
			stats.Stamina:    40,
		})
	case proto.Food_FoodCuttlesteak:
		character.AddStats(stats.Stats{
			stats.Spirit:  40,
			stats.Stamina: 40,
		})
	case proto.Food_FoodDragonfinFilet:
		character.AddStats(stats.Stats{
			stats.Strength: 40,
			stats.Stamina:  40,
		})
	case proto.Food_FoodBlackenedBasilisk:
		character.AddStats(stats.Stats{
			stats.SpellPower: 23,
			stats.Spirit:     20,
		})
	case proto.Food_FoodGrilledMudfish:
		character.AddStats(stats.Stats{
			stats.Agility: 20,
			stats.Spirit:  20,
		})
	case proto.Food_FoodRavagerDog:
		character.AddStats(stats.Stats{
			stats.AttackPower:       40,
			stats.RangedAttackPower: 40,
			stats.Spirit:            20,
		})
	case proto.Food_FoodRoastedClefthoof:
		character.AddStats(stats.Stats{
			stats.Strength: 20,
			stats.Spirit:   20,
		})
	case proto.Food_FoodSkullfishSoup:
		character.AddStats(stats.Stats{
			stats.SpellCrit: 20,
			stats.Spirit:    20,
		})
	case proto.Food_FoodSpicyHotTalbuk:
		character.AddStats(stats.Stats{
			stats.MeleeHit: 20,
			stats.Spirit:   20,
		})
	case proto.Food_FoodFishermansFeast:
		character.AddStats(stats.Stats{
			stats.Stamina: 30,
			stats.Spirit:  20,
		})
	}

	registerPotionCD(agent, consumes)
	registerConjuredCD(agent, consumes)
	registerExplosivesCD(agent, consumes)
}

func ApplyPetConsumeEffects(pet *Character, ownerConsumes *proto.Consumes) {
	switch ownerConsumes.PetFood {
	case proto.PetFood_PetFoodSpicedMammothTreats:
		pet.AddStats(stats.Stats{
			stats.Strength: 30,
			stats.Stamina:  30,
		})
	case proto.PetFood_PetFoodKiblersBits:
		pet.AddStats(stats.Stats{
			stats.Strength: 20,
			stats.Stamina:  20,
		})
	}

	pet.AddStat(stats.Agility, []float64{0, 5, 9, 13, 17, 20}[ownerConsumes.PetScrollOfAgility])
	pet.AddStat(stats.Strength, []float64{0, 5, 9, 13, 17, 20}[ownerConsumes.PetScrollOfStrength])
}

var PotionAuraTag = "Potion"

func registerPotionCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	defaultPotion := consumes.DefaultPotion
	startingPotion := consumes.PrepopPotion

	if defaultPotion == proto.Potions_UnknownPotion && startingPotion == proto.Potions_UnknownPotion {
		return
	}

	potionCD := character.NewTimer()

	startingMCD := makePotionActivation(startingPotion, character, potionCD)
	if startingMCD.Spell != nil {
		character.RegisterPrepullAction(-1*time.Second, func(sim *Simulation) {
			startingMCD.Spell.Cast(sim, nil)
			if startingPotion == proto.Potions_IndestructiblePotion {
				potionCD.Set(sim.CurrentTime + 2*time.Minute)
			} else {
				potionCD.Set(sim.CurrentTime + time.Minute)
			}
			character.UpdateMajorCooldowns()
		})
	}

	defaultMCD := makePotionActivation(defaultPotion, character, potionCD)
	if defaultMCD.Spell != nil {
		character.AddMajorCooldown(defaultMCD)
	}
}

var AlchStoneItemIDs = []int32{44322, 44323, 44324}

func (character *Character) HasAlchStone() bool {
	alchStoneEquipped := false
	for _, itemID := range AlchStoneItemIDs {
		alchStoneEquipped = alchStoneEquipped || character.HasTrinketEquipped(itemID)
	}
	return character.HasProfession(proto.Profession_Alchemy) && alchStoneEquipped
}

func makePotionActivation(potionType proto.Potions, character *Character, potionCD *Timer) MajorCooldown {
	alchStoneEquipped := character.HasAlchStone()
	hasEngi := character.HasProfession(proto.Profession_Engineering)

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
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					healthGain := sim.RollWithLabel(2700, 4500, "RunicHealingPotion")

					if alchStoneEquipped && potionType == proto.Potions_RunicHealingPotion {
						healthGain *= 1.40
					} else if hasEngi && potionType == proto.Potions_RunicHealingInjector {
						healthGain *= 1.25
					}
					character.GainHealth(sim, healthGain, healthMetrics)
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
				if alchStoneEquipped && potionType == proto.Potions_RunicManaPotion {
					manaGain *= 1.4
				} else if hasEngi && potionType == proto.Potions_RunicManaInjector {
					manaGain *= 1.25
				}
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= manaGain
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					manaGain := sim.RollWithLabel(4200, 4400, "RunicManaPotion")
					if alchStoneEquipped && potionType == proto.Potions_RunicManaPotion {
						manaGain *= 1.4
					} else if hasEngi && potionType == proto.Potions_RunicManaInjector {
						manaGain *= 1.25
					}
					character.AddMana(sim, manaGain, manaMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_IndestructiblePotion {
		actionID := ActionID{ItemID: 40093}
		aura := character.NewTemporaryStatsAura("Indestructible Potion", actionID, stats.Stats{stats.Armor: 3500}, time.Minute*2)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.RegisterSpell(SpellConfig{
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
			Spell: character.RegisterSpell(SpellConfig{
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
			Spell: character.RegisterSpell(SpellConfig{
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
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_SuperManaPotion {
		alchStoneEquipped := character.HasAlchStone()
		actionID := ActionID{ItemID: 22832}
		manaMetrics := character.NewManaMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= 3000
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					// Restores 1800 to 3000 mana. (2 Min Cooldown)
					manaGain := sim.RollWithLabel(1800, 3000, "super mana")
					if alchStoneEquipped {
						manaGain *= 1.4
					}
					character.AddMana(sim, manaGain, manaMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_HastePotion {
		actionID := ActionID{ItemID: 22838}
		aura := character.NewTemporaryStatsAura("Haste Potion", actionID, stats.Stats{stats.MeleeHaste: 400}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.RegisterSpell(SpellConfig{
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
			Spell: character.RegisterSpell(SpellConfig{
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
		alchStoneEquipped := character.HasAlchStone()
		if alchStoneEquipped {
			manaGain *= 1.4
		}
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
			Spell: character.RegisterSpell(SpellConfig{
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
			Spell: character.RegisterSpell(SpellConfig{
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
			Spell: character.RegisterSpell(SpellConfig{
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
			Spell: character.RegisterSpell(SpellConfig{
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

func registerConjuredCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	conjuredType := consumes.DefaultConjured

	if conjuredType == proto.Conjured_ConjuredDarkRune {
		actionID := ActionID{ItemID: 20520}
		manaMetrics := character.NewManaMetrics(actionID)
		// damageTakenManaMetrics := character.NewManaMetrics(ActionID{SpellID: 33776})
		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.GetConjuredCD(),
					Duration: time.Minute * 15,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				// Restores 900 to 1500 mana. (2 Min Cooldown)
				manaGain := sim.RollWithLabel(900, 1500, "dark rune")
				character.AddMana(sim, manaGain, manaMetrics)

				// if character.Class == proto.Class_ClassPaladin {
				// 	// Paladins gain extra mana from self-inflicted damage
				// 	// TO-DO: It is possible for damage to be resisted or to crit
				// 	// This would affect mana returns for Paladins
				// 	manaFromDamage := manaGain * 2.0 / 3.0 * 0.1
				// 	character.AddMana(sim, manaFromDamage, damageTakenManaMetrics, false)
				// }
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= 1500
			},
		})
	} else if conjuredType == proto.Conjured_ConjuredFlameCap {
		actionID := ActionID{ItemID: 22788}

		flameCapProc := character.RegisterSpell(SpellConfig{
			ActionID:    actionID,
			ProcMask:    ProcMaskEmpty,
			SpellSchool: SpellSchoolFire,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				spell.CalcAndDealDamage(sim, target, 40, spell.OutcomeMagicHitAndCrit)
			},
		})

		const procChance = 0.185
		var fireSpells []*Spell
		character.OnSpellRegistered(func(spell *Spell) {
			if spell.SpellSchool.Matches(SpellSchoolFire) {
				fireSpells = append(fireSpells, spell)
			}
		})

		flameCapAura := character.RegisterAura(Aura{
			Label:    "Flame Cap",
			ActionID: actionID,
			Duration: time.Minute,
			OnGain: func(aura *Aura, sim *Simulation) {
				for _, spell := range fireSpells {
					spell.BonusSpellPower += 80
				}
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				for _, spell := range fireSpells {
					spell.BonusSpellPower -= 80
				}
			},
			OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(ProcMaskMeleeOrRanged) {
					return
				}
				if sim.RandomFloat("Flame Cap Melee") > procChance {
					return
				}

				flameCapProc.Cast(sim, result.Target)
			},
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.GetConjuredCD(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				flameCapAura.Activate(sim)
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})
	} else if conjuredType == proto.Conjured_ConjuredHealthstone {
		actionID := ActionID{ItemID: 36892}
		healthMetrics := character.NewHealthMetrics(actionID)

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.GetConjuredCD(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				character.GainHealth(sim, 4280*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeSurvival,
		})
	}
}

var ThermalSapperActionID = ActionID{ItemID: 42641}
var ExplosiveDecoyActionID = ActionID{ItemID: 40536}
var SaroniteBombActionID = ActionID{ItemID: 41119}
var CobaltFragBombActionID = ActionID{ItemID: 40771}

func registerExplosivesCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	hasFiller := consumes.FillerExplosive != proto.Explosive_ExplosiveUnknown
	if !character.HasProfession(proto.Profession_Engineering) {
		return
	}
	if !consumes.ThermalSapper && !consumes.ExplosiveDecoy && !hasFiller {
		return
	}
	sharedTimer := character.NewTimer()

	if consumes.ThermalSapper {
		character.AddMajorCooldown(MajorCooldown{
			Spell:    character.newThermalSapperSpell(sharedTimer),
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 0.03,
		})
	}

	if consumes.ExplosiveDecoy {
		character.AddMajorCooldown(MajorCooldown{
			Spell:    character.newExplosiveDecoySpell(sharedTimer),
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 0.02,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Decoy puts other explosives on 2m CD, so only use if there won't be enough
				// time to use another explosive OR there is no filler explosive.
				return sim.GetRemainingDuration() < time.Minute || !hasFiller
			},
		})
	}

	if hasFiller {
		var filler *Spell
		switch consumes.FillerExplosive {
		case proto.Explosive_ExplosiveSaroniteBomb:
			filler = character.newSaroniteBombSpell(sharedTimer)
		case proto.Explosive_ExplosiveCobaltFragBomb:
			filler = character.newCobaltFragBombSpell(sharedTimer)
		}

		character.AddMajorCooldown(MajorCooldown{
			Spell:    filler,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 0.01,
		})
	}
}

// Creates a spell object for the common explosive case.
func (character *Character) newBasicExplosiveSpellConfig(sharedTimer *Timer, actionID ActionID, school SpellSchool, minDamage float64, maxDamage float64, cooldown Cooldown, minSelfDamage float64, maxSelfDamage float64) SpellConfig {
	dealSelfDamage := actionID.SameAction(ThermalSapperActionID)

	return SpellConfig{
		ActionID:    actionID,
		SpellSchool: school,
		ProcMask:    ProcMaskEmpty,

		Cast: CastConfig{
			CD: cooldown,
			SharedCD: Cooldown{
				Timer:    sharedTimer,
				Duration: TernaryDuration(actionID.SameAction(ExplosiveDecoyActionID), time.Minute*2, time.Minute),
			},
		},

		// Explosives always have 1% resist chance, so just give them hit cap.
		BonusHitRating:   100 * SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(minDamage, maxDamage) * sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			if dealSelfDamage {
				baseDamage := sim.Roll(minDamage, maxDamage)
				spell.CalcAndDealDamage(sim, &character.Unit, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	}
}
func (character *Character) newThermalSapperSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, ThermalSapperActionID, SpellSchoolFire, 2188, 2812, Cooldown{Timer: character.NewTimer(), Duration: time.Minute * 5}, 2188, 2812))
}
func (character *Character) newExplosiveDecoySpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, ExplosiveDecoyActionID, SpellSchoolPhysical, 1440, 2160, Cooldown{Timer: character.NewTimer(), Duration: time.Minute * 2}, 0, 0))
}
func (character *Character) newSaroniteBombSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, SaroniteBombActionID, SpellSchoolFire, 1150, 1500, Cooldown{}, 0, 0))
}
func (character *Character) newCobaltFragBombSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, CobaltFragBombActionID, SpellSchoolFire, 750, 1000, Cooldown{}, 0, 0))
}
