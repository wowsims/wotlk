package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Registers all consume-related effects to the Agent.
func applyConsumeEffects(agent Agent, raidBuffs proto.RaidBuffs, partyBuffs proto.PartyBuffs) {
	character := agent.GetCharacter()
	consumes := character.Consumes

	if consumes.Flask != proto.Flask_FlaskUnknown {
		switch consumes.Flask {
		case proto.Flask_FlaskOfTheFrostWyrm:
			character.AddStats(stats.Stats{
				stats.SpellPower:   125,
				stats.HealingPower: 125,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.SpellPower:   47,
					stats.HealingPower: 47,
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
			character.AddStats(stats.Stats{
				stats.NatureSpellPower: 80,
				stats.ArcaneSpellPower: 80,
				stats.HolySpellPower:   80,
			})
		case proto.Flask_FlaskOfMightyRestoration:
			character.AddStats(stats.Stats{
				stats.MP5: 25,
			})
		case proto.Flask_FlaskOfPureDeath:
			character.AddStats(stats.Stats{
				stats.FireSpellPower:   80,
				stats.FrostSpellPower:  80,
				stats.ShadowSpellPower: 80,
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
				stats.SpellPower:   58,
				stats.HealingPower: 58,
			})
		case proto.BattleElixir_WrathElixir:
			character.AddStats(stats.Stats{
				stats.AttackPower:       90,
				stats.RangedAttackPower: 90,
			})
		case proto.BattleElixir_AdeptsElixir:
			character.AddStats(stats.Stats{
				stats.SpellCrit:    24,
				stats.SpellPower:   24,
				stats.HealingPower: 24,
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
		case proto.BattleElixir_ElixirOfMajorFirePower:
			character.AddStats(stats.Stats{
				stats.FireSpellPower: 55,
			})
		case proto.BattleElixir_ElixirOfMajorFrostPower:
			character.AddStats(stats.Stats{
				stats.FrostSpellPower: 55,
			})
		case proto.BattleElixir_ElixirOfMajorShadowPower:
			character.AddStats(stats.Stats{
				stats.ShadowSpellPower: 55,
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

			var debuffAuras []*Aura
			for _, target := range character.Env.Encounter.Targets {
				debuffAuras = append(debuffAuras, GiftOfArthasAura(&target.Unit))
			}

			actionID := ActionID{SpellID: 11374}
			goaProc := character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				ApplyEffects: ApplyEffectFuncDirectDamage(SpellEffect{
					ProcMask:         ProcMaskEmpty,
					ThreatMultiplier: 1,
					FlatThreatBonus:  90,

					OutcomeApplier: character.OutcomeFuncAlwaysHit(),
					OnSpellHitDealt: func(sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
						debuffAuras[spellEffect.Target.Index].Activate(sim)
					},
				}),
			})

			character.RegisterAura(Aura{
				Label:    "Gift of Arthas",
				Duration: NeverExpires,
				OnReset: func(aura *Aura, sim *Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
					if spellEffect.Landed() &&
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
			stats.HealingPower:      46,
			stats.Stamina:           40,
		})
	case proto.Food_FoodGreatFeast:
		character.AddStats(stats.Stats{
			stats.AttackPower:       60,
			stats.RangedAttackPower: 60,
			stats.SpellPower:        35,
			stats.HealingPower:      35,
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
			stats.SpellPower:   46,
			stats.HealingPower: 46,
			stats.Stamina:      40,
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
			stats.SpellPower:   23,
			stats.HealingPower: 23,
			stats.Spirit:       20,
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

func ApplyPetConsumeEffects(pet *Character, ownerConsumes proto.Consumes) {
	switch ownerConsumes.PetFood {
	case proto.PetFood_PetFoodSpicedMammothTreats:
		pet.AddStats(stats.Stats{
			stats.Strength: 30,
			stats.Spirit:   30,
		})
	case proto.PetFood_PetFoodKiblersBits:
		pet.AddStats(stats.Stats{
			stats.Strength: 20,
			stats.Spirit:   20,
		})
	}

	pet.AddStat(stats.Agility, []float64{0, 5, 9, 13, 17, 20}[ownerConsumes.PetScrollOfAgility])
	pet.AddStat(stats.Strength, []float64{0, 5, 9, 13, 17, 20}[ownerConsumes.PetScrollOfStrength])
}

var PotionAuraTag = "Potion"

func registerPotionCD(agent Agent, consumes proto.Consumes) {
	character := agent.GetCharacter()
	defaultPotion := consumes.DefaultPotion
	startingPotion := consumes.PrepopPotion

	if defaultPotion == proto.Potions_UnknownPotion && startingPotion == proto.Potions_UnknownPotion {
		return
	}

	potionCD := character.NewTimer()

	prepopTime := time.Second
	startingMCD := makePotionActivation(startingPotion, character, potionCD, prepopTime)
	hasPrepopPotion := startingMCD.Spell != nil
	if hasPrepopPotion {
		startingPotionSpell := startingMCD.Spell
		character.RegisterResetEffect(func(sim *Simulation) {
			StartDelayedAction(sim, DelayedActionOptions{
				DoAt: 0,
				OnAction: func(sim *Simulation) {
					startingPotionSpell.Cast(sim, nil)
					potionCD.Set(time.Minute - prepopTime)
					character.UpdateMajorCooldowns()
				},
			})
		})
	}

	defaultMCD := makePotionActivation(defaultPotion, character, potionCD, 0)
	if defaultMCD.Spell != nil {
		usedDefaultPotion := false

		canActivate := defaultMCD.CanActivate
		defaultMCD.CanActivate = func(sim *Simulation, character *Character) bool {
			if usedDefaultPotion {
				return false
			}

			if hasPrepopPotion && sim.CurrentTime < time.Second*1 {
				// Because of prepop's StartDelayedAction call, regular potion actually gets
				// checked first so we need to make sure it doesn't activate.
				return false
			}

			if canActivate != nil {
				return canActivate(sim, character)
			} else {
				return true
			}
		}

		defaultPotionSpell := defaultMCD.Spell
		defaultMCD.ActivationFactory = func(sim *Simulation) CooldownActivation {
			usedDefaultPotion = false
			if defaultPotion == proto.Potions_SuperManaPotion {
				character.ExpectedBonusMana += float64((3000 + 1800) / 2)
			}
			if defaultPotion == proto.Potions_RunicManaPotion {
				character.ExpectedBonusMana += float64((4200 + 4400) / 2)
			}

			return func(sim *Simulation, character *Character) {
				usedDefaultPotion = true
				defaultPotionSpell.Cast(sim, nil)

				if defaultPotion == proto.Potions_SuperManaPotion {
					character.ExpectedBonusMana -= float64((3000 + 1800) / 2)
				}
				if defaultPotion == proto.Potions_RunicManaPotion {
					character.ExpectedBonusMana -= float64((4200 + 4400) / 2)
				}
			}
		}
		character.AddMajorCooldown(defaultMCD)
	}
}

var AlchStoneItemIDs = []int32{13503, 35748, 35749, 35750, 35751, 44322, 44323, 44324}

func (character *Character) HasAlchStone() bool {
	alchStoneEquipped := false
	for _, itemID := range AlchStoneItemIDs {
		alchStoneEquipped = alchStoneEquipped || character.HasTrinketEquipped(itemID)
	}
	return character.HasProfession(proto.Profession_Alchemy) && alchStoneEquipped
}

func makePotionActivation(potionType proto.Potions, character *Character, potionCD *Timer, prepopTime time.Duration) MajorCooldown {
	if potionType == proto.Potions_RunicHealingPotion {
		actionID := ActionID{ItemID: 33447}
		healthMetrics := character.NewHealthMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeSurvival,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					healthGain := 2700.0 + (4500.0-2700.0)*sim.RandomFloat("RunicHealingPotion")
					character.GainHealth(sim, healthGain, healthMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_RunicManaPotion {
		alchStoneEquipped := character.HasAlchStone()
		actionID := ActionID{ItemID: 33448}
		manaMetrics := character.NewManaMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeMana,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= 4400
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					manaGain := 4200 + (4400.0-4200.0)*sim.RandomFloat("RunicManaPotion")
					if alchStoneEquipped {
						manaGain *= 1.4
					}
					character.AddMana(sim, manaGain, manaMetrics, true)
				},
			}),
		}
	} else if potionType == proto.Potions_IndestructiblePotion {
		actionID := ActionID{ItemID: 40093}
		aura := character.NewTemporaryStatsAura("Indestructible Potion", actionID, stats.Stats{stats.Armor: 3500}, time.Minute*2-prepopTime)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_PotionOfSpeed {
		actionID := ActionID{ItemID: 40211}
		aura := character.NewTemporaryStatsAura("Potion of Speed", actionID, stats.Stats{stats.MeleeHaste: 500, stats.SpellHaste: 500}, time.Second*15-prepopTime)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_PotionOfWildMagic {
		actionID := ActionID{ItemID: 40212}
		aura := character.NewTemporaryStatsAura("Potion of Wild Magic", actionID, stats.Stats{stats.SpellPower: 200, stats.HealingPower: 200, stats.SpellCrit: 200, stats.MeleeCrit: 200}, time.Second*15-prepopTime)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_DestructionPotion {
		actionID := ActionID{ItemID: 22839}
		aura := character.NewTemporaryStatsAura("Destruction Potion", actionID, stats.Stats{stats.SpellPower: 120, stats.SpellCrit: 2 * CritRatingPerCritChance}, time.Second*15-prepopTime)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
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
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				if character.MaxMana()-(character.CurrentMana()+totalRegen) < 3000 {
					return false
				}

				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					// Restores 1800 to 3000 mana. (2 Min Cooldown)
					manaGain := 1800 + (sim.RandomFloat("super mana") * 1200)
					if alchStoneEquipped {
						manaGain *= 1.4
					}
					character.AddMana(sim, manaGain, manaMetrics, true)
				},
			}),
		}
	} else if potionType == proto.Potions_HastePotion {
		actionID := ActionID{ItemID: 22838}
		aura := character.NewTemporaryStatsAura("Haste Potion", actionID, stats.Stats{stats.MeleeHaste: 400}, time.Second*15-prepopTime)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_MightyRagePotion {
		actionID := ActionID{ItemID: 13442}
		aura := character.NewTemporaryStatsAura("Mighty Rage Potion", actionID, stats.Stats{stats.Strength: 60}, time.Second*15-prepopTime)
		rageMetrics := character.NewRageMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				if character.Class == proto.Class_ClassWarrior {
					return character.CurrentRage() < 25
				}
				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
					if character.Class == proto.Class_ClassWarrior {
						bonusRage := 45.0 + (75.0-45.0)*sim.RandomFloat("Mighty Rage Potion")
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

		buffAura := character.NewTemporaryStatsAura("Fel Mana Potion", actionID, stats.Stats{stats.MP5: mp5}, time.Second*24-prepopTime)
		debuffAura := character.NewTemporaryStatsAura("Fel Mana Potion Debuff", ActionID{SpellID: 38927}, stats.Stats{stats.SpellPower: -25}, time.Minute*15)

		return MajorCooldown{
			Type: CooldownTypeMana,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have low enough mana. The potion takes effect over 24
				// seconds so we can pop it a little earlier than the full value.
				if character.MaxMana()-character.CurrentMana() < 2000 {
					return false
				}

				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					buffAura.Activate(sim)
					debuffAura.Activate(sim)
					debuffAura.Refresh(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_InsaneStrengthPotion {
		actionID := ActionID{ItemID: 22828}
		aura := character.NewTemporaryStatsAura("Insane Strength Potion", actionID, stats.Stats{stats.Strength: 120, stats.Defense: -75}, time.Second*15-prepopTime)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_IronshieldPotion {
		actionID := ActionID{ItemID: 22849}
		aura := character.NewTemporaryStatsAura("Ironshield Potion", actionID, stats.Stats{stats.Armor: 2500}, time.Minute*2-prepopTime)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_HeroicPotion {
		actionID := ActionID{ItemID: 22837}
		aura := character.NewTemporaryStatsAura("Heroic Potion", actionID, stats.Stats{stats.Strength: 70, stats.Health: 700}, time.Second*15-prepopTime)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 1,
					},
				},
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

func registerConjuredCD(agent Agent, consumes proto.Consumes) {
	if consumes.DefaultConjured == consumes.StartingConjured {
		// Starting conjured is redundant in this case.
		consumes.StartingConjured = proto.Conjured_ConjuredUnknown
	}
	if consumes.StartingConjured == proto.Conjured_ConjuredUnknown {
		consumes.NumStartingConjured = 0
	}
	if consumes.NumStartingConjured == 0 {
		consumes.StartingConjured = proto.Conjured_ConjuredUnknown
	}
	character := agent.GetCharacter()

	defaultMCD, defaultSpell := makeConjuredActivation(consumes.DefaultConjured, character)
	startingMCD, startingSpell := makeConjuredActivation(consumes.StartingConjured, character)
	numStartingConjured := int(consumes.NumStartingConjured)
	if defaultSpell == nil && startingSpell == nil {
		return
	}

	numStartingConjuredUsed := 0

	if startingSpell != nil {
		character.AddMajorCooldown(MajorCooldown{
			Spell: startingSpell,
			Type:  startingMCD.Type,
			CanActivate: func(sim *Simulation, character *Character) bool {
				if numStartingConjuredUsed >= numStartingConjured {
					return false
				}
				return startingMCD.CanActivate(sim, character)
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return startingMCD.ShouldActivate(sim, character)
			},
			ActivationFactory: func(sim *Simulation) CooldownActivation {
				numStartingConjuredUsed = 0
				return func(sim *Simulation, character *Character) {
					startingSpell.Cast(sim, nil)
					numStartingConjuredUsed++
				}
			},
		})
	}

	if defaultSpell != nil {
		character.AddMajorCooldown(MajorCooldown{
			Spell: defaultSpell,
			Type:  defaultMCD.Type,
			CanActivate: func(sim *Simulation, character *Character) bool {
				if numStartingConjuredUsed < numStartingConjured {
					return false
				}
				return defaultMCD.CanActivate(sim, character)
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return defaultMCD.ShouldActivate(sim, character)
			},
			ActivationFactory: func(sim *Simulation) CooldownActivation {
				expectedManaPerUsage := float64((900 + 600) / 2)

				remainingUsages := int(1 + (MaxDuration(0, sim.Duration))/(darkRuneCD))

				if consumes.DefaultConjured == proto.Conjured_ConjuredDarkRune {
					character.ExpectedBonusMana += expectedManaPerUsage * float64(remainingUsages)
				}

				return func(sim *Simulation, character *Character) {
					defaultSpell.Cast(sim, nil)

					if consumes.DefaultConjured == proto.Conjured_ConjuredDarkRune {
						// Update expected bonus mana
						newRemainingUsages := int(sim.GetRemainingDuration() / (darkRuneCD))
						character.ExpectedBonusMana -= expectedManaPerUsage * float64(remainingUsages-newRemainingUsages)
						remainingUsages = newRemainingUsages
					}
				}
			},
		})
	}
}

const darkRuneCD = time.Minute * 15

func makeConjuredActivation(conjuredType proto.Conjured, character *Character) (MajorCooldown, *Spell) {
	if conjuredType == proto.Conjured_ConjuredDarkRune {
		actionID := ActionID{ItemID: 20520}
		manaMetrics := character.NewManaMetrics(actionID)
		// damageTakenManaMetrics := character.NewManaMetrics(ActionID{SpellID: 33776})
		return MajorCooldown{
				Type: CooldownTypeMana,
				CanActivate: func(sim *Simulation, character *Character) bool {
					return true
				},
				ShouldActivate: func(sim *Simulation, character *Character) bool {
					// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
					totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
					if character.MaxMana()-(character.CurrentMana()+totalRegen) < 1500 {
						return false
					}
					return true
				},
			},
			character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    character.GetConjuredCD(),
						Duration: darkRuneCD,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					// Restores 900 to 1500 mana. (2 Min Cooldown)
					manaGain := 900 + (sim.RandomFloat("dark rune") * 600)
					character.AddMana(sim, manaGain, manaMetrics, true)

					// if character.Class == proto.Class_ClassPaladin {
					// 	// Paladins gain extra mana from self-inflicted damage
					// 	// TO-DO: It is possible for damage to be resisted or to crit
					// 	// This would affect mana returns for Paladins
					// 	manaFromDamage := manaGain * 2.0 / 3.0 * 0.1
					// 	character.AddMana(sim, manaFromDamage, damageTakenManaMetrics, false)
					// }
				},
			})
	} else if conjuredType == proto.Conjured_ConjuredFlameCap {
		actionID := ActionID{ItemID: 22788}

		flameCapProc := character.RegisterSpell(SpellConfig{
			ActionID:    actionID,
			SpellSchool: SpellSchoolFire,
			ApplyEffects: ApplyEffectFuncDirectDamage(SpellEffect{
				ProcMask:         ProcMaskEmpty,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				BaseDamage:     BaseDamageConfigFlat(40),
				OutcomeApplier: character.OutcomeFuncMagicHitAndCrit(character.DefaultSpellCritMultiplier()),
			}),
		})

		const procChance = 0.185
		flameCapAura := character.NewTemporaryStatsAura("Flame Cap", actionID, stats.Stats{stats.FireSpellPower: 80}, time.Minute)
		flameCapAura.OnSpellHitDealt = func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(ProcMaskMeleeOrRanged) {
				return
			}
			if sim.RandomFloat("Flame Cap Melee") > procChance {
				return
			}

			flameCapProc.Cast(sim, spellEffect.Target)
		}

		return MajorCooldown{
				Type: CooldownTypeDPS,
				CanActivate: func(sim *Simulation, character *Character) bool {
					return true
				},
				ShouldActivate: func(sim *Simulation, character *Character) bool {
					return true
				},
			},
			character.RegisterSpell(SpellConfig{
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
	} else if conjuredType == proto.Conjured_ConjuredHealthstone {
		actionID := ActionID{ItemID: 22105}
		healthMetrics := character.NewHealthMetrics(actionID)
		return MajorCooldown{
				Type: CooldownTypeSurvival,
				CanActivate: func(sim *Simulation, character *Character) bool {
					return true
				},
				ShouldActivate: func(sim *Simulation, character *Character) bool {
					return true
				},
			},
			character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    character.GetConjuredCD(),
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					character.GainHealth(sim, 2496, healthMetrics)
				},
			})
	} else {
		return MajorCooldown{}, nil
	}
}

var ThermalSapperActionID = ActionID{ItemID: 42641}
var ExplosiveDecoyActionID = ActionID{ItemID: 40536}
var SaroniteBombActionID = ActionID{ItemID: 41119}
var CobaltFragBombActionID = ActionID{ItemID: 40771}

func registerExplosivesCD(agent Agent, consumes proto.Consumes) {
	character := agent.GetCharacter()
	if !character.HasProfession(proto.Profession_Engineering) {
		return
	}
	if !consumes.ThermalSapper && !consumes.ExplosiveDecoy && consumes.FillerExplosive == proto.Explosive_ExplosiveUnknown {
		return
	}
	explosivesTimer := character.NewTimer()
	sharedTimer := character.NewTimer()

	var explosives []*Spell

	if consumes.ThermalSapper {
		explosives = append(explosives, character.newThermalSapperSpell(sharedTimer))
	}
	if consumes.ExplosiveDecoy {
		explosives = append(explosives, character.newExplosiveDecoySpell(sharedTimer))
	}

	switch consumes.FillerExplosive {
	case proto.Explosive_ExplosiveSaroniteBomb:
		explosives = append(explosives, character.newSaroniteBombSpell(sharedTimer))
	case proto.Explosive_ExplosiveCobaltFragBomb:
		explosives = append(explosives, character.newCobaltFragBombSpell(sharedTimer))
	}

	spell := character.RegisterSpell(SpellConfig{
		ActionID: ThermalSapperActionID,
		Flags:    SpellFlagNoOnCastComplete | SpellFlagNoMetrics | SpellFlagNoLogs,

		Cast: CastConfig{
			CD: Cooldown{
				Timer:    explosivesTimer,
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *Simulation, target *Unit, _ *Spell) {
			for _, explosive := range explosives {
				if explosive.IsReady(sim) {
					explosive.Cast(sim, target)
					break
				}
			}

			nextExplosiveAt := sim.CurrentTime + time.Minute*5
			for _, explosive := range explosives {
				nextExplosiveAt = MinDuration(explosive.ReadyAt(), nextExplosiveAt)
			}
			explosivesTimer.Set(nextExplosiveAt)
		},
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell: spell,
		Type:  CooldownTypeDPS,
	})
}

// Creates a spell object for the common explosive case.
func (character *Character) newBasicExplosiveSpellConfig(sharedTimer *Timer, actionID ActionID, school SpellSchool, minDamage float64, maxDamage float64, cooldown Cooldown, minSelfDamage float64, maxSelfDamage float64) SpellConfig {
	return SpellConfig{
		ActionID:    actionID,
		SpellSchool: school,

		Cast: CastConfig{
			CD: cooldown,
			SharedCD: Cooldown{
				Timer:    sharedTimer,
				Duration: time.Minute,
			},
		},

		ApplyEffects: ApplyEffectFuncAOEDamage(character.Env, SpellEffect{
			ProcMask: ProcMaskEmpty,
			// Explosives always have 1% resist chance, so just give them hit cap.
			BonusSpellHitRating: 100 * SpellHitRatingPerHitChance,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage:     BaseDamageConfigRoll(minDamage, maxDamage),
			OutcomeApplier: character.OutcomeFuncMagicHitAndCrit(2),
			// TODO: Deal self-damage
			//OnSpellHitDealt: func(sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			//},
		}),
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
