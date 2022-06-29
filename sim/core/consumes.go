package core

import (
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Registers all consume-related effects to the Agent.
func applyConsumeEffects(agent Agent, raidBuffs proto.RaidBuffs, partyBuffs proto.PartyBuffs) {
	character := agent.GetCharacter()
	consumes := character.Consumes

	if consumes.Flask != proto.Flask_FlaskUnknown {
		switch consumes.Flask {
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

	switch consumes.Alchohol {
	case proto.Alchohol_AlchoholKreegsStoutBeatdown:
		character.AddStats(stats.Stats{
			stats.Intellect: -5,
			stats.Spirit:    25,
		})
	}

	// Scrolls
	character.AddStat(stats.Agility, []float64{0, 5, 9, 13, 17, 20}[consumes.ScrollOfAgility])
	character.AddStat(stats.Strength, []float64{0, 5, 9, 13, 17, 20}[consumes.ScrollOfStrength])
	if !character.HasRingEquipped(29297) {
		// Proc from Band of Eternal Defender removes scroll.
		character.AddStat(stats.Armor, []float64{0, 60, 120, 180, 240, 300}[consumes.ScrollOfProtection])
	}
	if raidBuffs.DivineSpirit == proto.TristateEffect_TristateEffectMissing {
		// Doesn't stack with DS
		character.AddStat(stats.Spirit, []float64{0, 3, 7, 11, 15, 30}[consumes.ScrollOfSpirit])
	}

	// Weapon Imbues
	allowMHImbue := character.HasMHWeapon() && (character.HasMHWeaponImbue || partyBuffs.WindfuryTotemRank == 0)
	if allowMHImbue {
		addImbueStats(character, consumes.MainHandImbue)
	}
	if character.HasOHWeapon() {
		addImbueStats(character, consumes.OffHandImbue)
	}

	registerPotionCD(agent, consumes)
	registerConjuredCD(agent, consumes)
	registerDrumsCD(agent, partyBuffs, consumes)
	registerExplosivesCD(agent, consumes)
}

func ApplyPetConsumeEffects(pet *Character, ownerConsumes proto.Consumes) {
	switch ownerConsumes.PetFood {
	case proto.PetFood_PetFoodKiblersBits:
		pet.AddStats(stats.Stats{
			stats.Strength: 20,
			stats.Spirit:   20,
		})
	}

	pet.AddStat(stats.Agility, []float64{0, 5, 9, 13, 17, 20}[ownerConsumes.PetScrollOfAgility])
	pet.AddStat(stats.Strength, []float64{0, 5, 9, 13, 17, 20}[ownerConsumes.PetScrollOfStrength])
}

func addImbueStats(character *Character, imbue proto.WeaponImbue) {
	if imbue == proto.WeaponImbue_WeaponImbueUnknown {
		return
	}

	if imbue == proto.WeaponImbue_WeaponImbueAdamantiteSharpeningStone {
		character.PseudoStats.BonusDamage += 12
		// Crit rating for sharpening stone applies to melee only.
		if character.Class != proto.Class_ClassHunter {
			// For non-hunters just give direct crit so it shows on the stats panel.
			character.AddStats(stats.Stats{stats.MeleeCrit: 14})
		} else {
			character.PseudoStats.BonusMeleeCritRating += 14
		}
	} else if imbue == proto.WeaponImbue_WeaponImbueAdamantiteWeightstone {
		character.PseudoStats.BonusDamage += 12
		character.AddStats(stats.Stats{stats.MeleeCrit: 14})
	} else if imbue == proto.WeaponImbue_WeaponImbueElementalSharpeningStone {
		character.AddStats(stats.Stats{
			stats.MeleeCrit: 28,
		})
	} else if imbue == proto.WeaponImbue_WeaponImbueBrilliantWizardOil {
		character.AddStats(stats.Stats{
			stats.SpellCrit:    14,
			stats.SpellPower:   36,
			stats.HealingPower: 36,
		})
	} else if imbue == proto.WeaponImbue_WeaponImbueSuperiorWizardOil {
		character.AddStats(stats.Stats{
			stats.SpellPower:   42,
			stats.HealingPower: 42,
		})
	} else if imbue == proto.WeaponImbue_WeaponImbueRighteousWeaponCoating {
		procAura := character.NewTemporaryStatsAura("RighteousWeaponCoatingProc", ActionID{ItemID: 34539}, stats.Stats{stats.AttackPower: 300, stats.RangedAttackPower: 300}, time.Second*10)
		ppmm := character.AutoAttacks.NewPPMManager(10.0, ProcMaskMeleeOrRanged)

		icd := Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}

		character.GetOrRegisterAura(Aura{
			Label:    "Righteous Weapon Coating",
			Duration: NeverExpires,
			OnReset: func(aura *Aura, sim *Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
				// mask 340
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(ProcMaskDirect) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if !ppmm.Proc(sim, spellEffect.ProcMask, "Righteous Weapon Coating") {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	}
}

var DrumsAuraTag = "Drums"

const DrumsCD = time.Minute * 2 // Tinnitus
var DrumsOfBattleActionID = ActionID{SpellID: 351355}
var DrumsOfRestorationActionID = ActionID{SpellID: 351358}
var DrumsOfWarActionID = ActionID{SpellID: 351360}

// Adds drums as a major cooldown to the character, if it's being used.
func registerDrumsCD(agent Agent, partyBuffs proto.PartyBuffs, consumes proto.Consumes) {
	character := agent.GetCharacter()
	drumsType := proto.Drums_DrumsUnknown

	// Whether this agent is the one casting the drums.
	drumsSelfCast := false

	if consumes.Drums != proto.Drums_DrumsUnknown {
		drumsType = consumes.Drums
		drumsSelfCast = true

		// Disable self-drums on other party members, so there is only 1 drummer.
		for _, partyMember := range agent.GetCharacter().Party.Players {
			if partyMember != agent {
				partyMember.GetCharacter().Consumes.Drums = proto.Drums_DrumsUnknown
			}
		}
	} else if partyBuffs.Drums != proto.Drums_DrumsUnknown {
		drumsType = partyBuffs.Drums
	}

	// If we aren't casting drums, and there is another real party member doing so, then we're done.
	if !drumsSelfCast {
		for _, partyMember := range agent.GetCharacter().Party.Players {
			if partyMember != agent && partyMember.GetCharacter().Consumes.Drums != proto.Drums_DrumsUnknown {
				return
			}
		}
	}

	var actionID ActionID
	var cooldownType CooldownType
	if drumsType == proto.Drums_DrumsOfBattle {
		actionID = DrumsOfBattleActionID
		cooldownType = CooldownTypeDPS
	} else if drumsType == proto.Drums_DrumsOfRestoration {
		actionID = DrumsOfRestorationActionID
		cooldownType = CooldownTypeMana
	} else if drumsType == proto.Drums_DrumsOfWar {
		actionID = DrumsOfWarActionID
		cooldownType = CooldownTypeDPS
	} else {
		return
	}

	mcd := MajorCooldown{
		Priority: CooldownPriorityDrums,
		Type:     cooldownType | CooldownTypeUsableShapeShifted,
	}

	if drumsSelfCast {
		auras := []*Aura{}
		for _, partyMember := range character.Party.Players {
			drumsAura := makeDrumsAura(partyMember.GetCharacter(), drumsType)
			auras = append(auras, drumsAura)
		}

		mcd.Spell = character.RegisterSpell(SpellConfig{
			ActionID: actionID,

			Cast: CastConfig{
				DefaultCast: Cast{
					GCD: GCDDefault,
				},
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: DrumsCD,
				},
			},

			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				// When a real player is using drums, their cast applies to the whole party.
				for _, aura := range auras {
					aura.Activate(sim)
				}
			},
		})
	} else {
		// When there is no real player using drums, each player gets a fake CD that
		// gives just themself the buff, with no cast time.
		drumsAura := makeDrumsAura(agent.GetCharacter(), drumsType)
		mcd.Spell = character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoMetrics,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: DrumsCD,
				},
			},

			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				drumsAura.Activate(sim)
			},
		})
	}

	agent.GetCharacter().AddMajorCooldown(mcd)
}

func makeDrumsAura(character *Character, drumsType proto.Drums) *Aura {
	if drumsType == proto.Drums_DrumsOfBattle {
		const hasteBonus = 80
		return character.NewTemporaryStatsAuraWrapped("Drums of Battle", DrumsOfBattleActionID, stats.Stats{stats.MeleeHaste: hasteBonus, stats.SpellHaste: hasteBonus}, time.Second*30, func(aura *Aura) {
			oldOnGain := aura.OnGain
			aura.OnGain = func(aura *Aura, sim *Simulation) {
				oldOnGain(aura, sim)
				// Drums of battle doesn't affect pets, ask Blizzard.
			}
		})
	} else if drumsType == proto.Drums_DrumsOfRestoration {
		// 600 mana over 15 seconds == 200 mp5
		const mp5Bonus = 200

		petAuras := []*Aura{}
		for _, petAgent := range character.Party.Pets {
			pet := petAgent.GetPet()
			petAuras = append(petAuras, pet.NewTemporaryStatsAura("Drums of Restoration", DrumsOfRestorationActionID, stats.Stats{stats.MP5: mp5Bonus}, time.Second*15))
		}

		return character.NewTemporaryStatsAuraWrapped("Drums of Restoration", DrumsOfRestorationActionID, stats.Stats{stats.MP5: mp5Bonus}, time.Second*15, func(aura *Aura) {
			oldOnGain := aura.OnGain
			aura.OnGain = func(aura *Aura, sim *Simulation) {
				oldOnGain(aura, sim)

				for i, petAgent := range character.Party.Pets {
					pet := petAgent.GetPet()
					if pet.IsEnabled() {
						petAuras[i].Activate(sim)
					}
				}
			}
		})
	} else if drumsType == proto.Drums_DrumsOfWar {
		petAuras := []*Aura{}
		for _, petAgent := range character.Party.Pets {
			pet := petAgent.GetPet()
			petAuras = append(petAuras, pet.NewTemporaryStatsAura("Drums of War", DrumsOfWarActionID, stats.Stats{stats.AttackPower: 60, stats.RangedAttackPower: 60, stats.SpellPower: 30}, time.Second*30))
		}

		return character.NewTemporaryStatsAuraWrapped("Drums of War", DrumsOfWarActionID, stats.Stats{stats.AttackPower: 60, stats.RangedAttackPower: 60, stats.SpellPower: 30}, time.Second*30, func(aura *Aura) {
			oldOnGain := aura.OnGain
			aura.OnGain = func(aura *Aura, sim *Simulation) {
				oldOnGain(aura, sim)

				for i, petAgent := range character.Party.Pets {
					pet := petAgent.GetPet()
					if pet.IsEnabled() {
						petAuras[i].Activate(sim)
					}
				}
			}
		})
	} else {
		return nil
	}
}

var PotionAuraTag = "Potion"

func registerPotionCD(agent Agent, consumes proto.Consumes) {
	if consumes.DefaultPotion == consumes.StartingPotion {
		// Starting potion is redundant in this case.
		consumes.StartingPotion = proto.Potions_UnknownPotion
	}
	if consumes.StartingPotion == proto.Potions_UnknownPotion {
		consumes.NumStartingPotions = 0
	}
	if consumes.NumStartingPotions == 0 {
		consumes.StartingPotion = proto.Potions_UnknownPotion
	}

	character := agent.GetCharacter()
	defaultPotion := consumes.DefaultPotion
	startingPotion := consumes.StartingPotion
	numStartingPotions := int(consumes.NumStartingPotions)

	if defaultPotion == proto.Potions_UnknownPotion && startingPotion == proto.Potions_UnknownPotion {
		return
	}

	potionCD := character.NewTimer()
	defaultMCD, defaultPotionSpell := makePotionActivation(defaultPotion, character, potionCD)
	startingMCD, startingPotionSpell := makePotionActivation(startingPotion, character, potionCD)

	numStartingPotionsUsed := 0
	var mageManaGemMCD *MajorCooldown

	// Setup for expected bonus mana management.
	expectedManaPerUsage := float64((3000 + 1800) / 2)
	remainingUsages := 0
	remainingManaPotionUsages := 0
	eitherPotionIsManaPotion := startingPotion == proto.Potions_SuperManaPotion || defaultPotion == proto.Potions_SuperManaPotion
	initExpectedBonusMana := func(sim *Simulation, potionType proto.Potions, isStartingPotion bool) {
		if potionType != proto.Potions_SuperManaPotion {
			return
		}
		mageManaGemMCD = character.GetMajorCooldown(MageManaGemMCDActionID)

		remainingUsages = int(1 + (MaxDuration(0, sim.Duration))/(time.Minute*2))
		if isStartingPotion {
			remainingManaPotionUsages = MinInt(numStartingPotions, remainingUsages)
		} else {
			remainingManaPotionUsages = MaxInt(0, remainingUsages-numStartingPotions)
		}

		character.ExpectedBonusMana += expectedManaPerUsage * float64(remainingManaPotionUsages)
	}

	updateExpectedBonusMana := func(sim *Simulation, character *Character, potionType proto.Potions, isStartingPotion bool) {
		if !eitherPotionIsManaPotion {
			return
		}
		thisPotionIsManaPotion := potionType == proto.Potions_SuperManaPotion

		newRemainingUsages := int(sim.GetRemainingDuration() / (time.Minute * 2))
		numRemainingStartingPotionUsages := MinInt(numStartingPotions-numStartingPotionsUsed, remainingUsages)
		numRemainingDefaultPotionUsages := MaxInt(0, newRemainingUsages-MaxInt(0, numStartingPotions-numStartingPotionsUsed))

		newRemainingManaPotionUsages := 0
		if thisPotionIsManaPotion == isStartingPotion {
			newRemainingManaPotionUsages = numRemainingStartingPotionUsages
		} else {
			newRemainingManaPotionUsages = numRemainingDefaultPotionUsages
		}

		character.ExpectedBonusMana -= expectedManaPerUsage * float64(remainingManaPotionUsages-newRemainingManaPotionUsages)
		remainingManaPotionUsages = newRemainingManaPotionUsages
	}

	if startingPotionSpell != nil {
		character.AddMajorCooldown(MajorCooldown{
			Spell: startingPotionSpell,
			CanActivate: func(sim *Simulation, character *Character) bool {
				if numStartingPotionsUsed >= numStartingPotions {
					return false
				}
				return startingMCD.CanActivate(sim, character)
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				if startingPotion == proto.Potions_SuperManaPotion &&
					mageManaGemMCD != nil &&
					mageManaGemMCD.IsEnabled() &&
					mageManaGemMCD.IsReady(sim) {
					return false
				}
				return startingMCD.ShouldActivate(sim, character)
			},

			Type: startingMCD.Type,

			ActivationFactory: func(sim *Simulation) CooldownActivation {
				numStartingPotionsUsed = 0
				initExpectedBonusMana(sim, startingPotion, true)

				return func(sim *Simulation, character *Character) {
					startingPotionSpell.Cast(sim, nil)
					numStartingPotionsUsed++
					updateExpectedBonusMana(sim, character, startingPotion, true)
				}
			},
		})
	}

	if defaultPotionSpell != nil {
		character.AddMajorCooldown(MajorCooldown{
			Spell: defaultPotionSpell,
			CanActivate: func(sim *Simulation, character *Character) bool {
				if numStartingPotionsUsed < numStartingPotions {
					return false
				}
				return defaultMCD.CanActivate(sim, character)
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				if defaultPotion == proto.Potions_SuperManaPotion &&
					mageManaGemMCD != nil &&
					mageManaGemMCD.IsEnabled() &&
					mageManaGemMCD.IsReady(sim) {
					return false
				}
				return defaultMCD.ShouldActivate(sim, character)
			},
			Type: defaultMCD.Type,

			ActivationFactory: func(sim *Simulation) CooldownActivation {
				initExpectedBonusMana(sim, defaultPotion, true)

				return func(sim *Simulation, character *Character) {
					defaultPotionSpell.Cast(sim, nil)
					updateExpectedBonusMana(sim, character, defaultPotion, true)
				}
			},
		})
	}
}

var AlchStoneItemIDs = []int32{13503, 35748, 35749, 35750, 35751}

func makePotionActivation(potionType proto.Potions, character *Character, potionCD *Timer) (MajorCooldown, *Spell) {
	if potionType == proto.Potions_DestructionPotion {
		actionID := ActionID{ItemID: 22839}
		aura := character.NewTemporaryStatsAura("Destruction Potion", actionID, stats.Stats{stats.SpellPower: 120, stats.SpellCrit: 2 * SpellCritRatingPerCritChance}, time.Second*15)
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
						Timer:    potionCD,
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			})
	} else if potionType == proto.Potions_SuperManaPotion {
		alchStoneEquipped := false
		for _, itemID := range AlchStoneItemIDs {
			alchStoneEquipped = alchStoneEquipped || character.HasTrinketEquipped(itemID)
		}
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
			},
			character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 2,
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
			})
	} else if potionType == proto.Potions_HastePotion {
		actionID := ActionID{ItemID: 22838}
		aura := character.NewTemporaryStatsAura("Haste Potion", actionID, stats.Stats{stats.MeleeHaste: 400}, time.Second*15)
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
						Timer:    potionCD,
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			})
	} else if potionType == proto.Potions_MightyRagePotion {
		actionID := ActionID{ItemID: 13442}
		aura := character.NewTemporaryStatsAura("Mighty Rage Potion", actionID, stats.Stats{stats.Strength: 60}, time.Second*15)
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
			},
			character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
					if character.Class == proto.Class_ClassWarrior {
						bonusRage := 45.0 + (75.0-45.0)*sim.RandomFloat("Mighty Rage Potion")
						character.AddRage(sim, bonusRage, rageMetrics)
					}
				},
			})
	} else if potionType == proto.Potions_FelManaPotion {
		actionID := ActionID{ItemID: 31677}

		// Restores 3200 mana over 24 seconds.
		manaGain := 3200.0
		alchStoneEquipped := false
		for _, itemID := range AlchStoneItemIDs {
			alchStoneEquipped = alchStoneEquipped || character.HasTrinketEquipped(itemID)
		}
		if alchStoneEquipped {
			manaGain *= 1.4
		}
		mp5 := manaGain / 24 * 5

		buffAura := character.NewTemporaryStatsAura("Fel Mana Potion", actionID, stats.Stats{stats.MP5: mp5}, time.Second*24)
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
			},
			character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    potionCD,
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					buffAura.Activate(sim)
					debuffAura.Activate(sim)
					debuffAura.Refresh(sim)
				},
			})
	} else if potionType == proto.Potions_InsaneStrengthPotion {
		actionID := ActionID{ItemID: 22828}
		aura := character.NewTemporaryStatsAura("Insane Strength Potion", actionID, stats.Stats{stats.Strength: 120, stats.Defense: -75}, time.Second*15)
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
						Timer:    potionCD,
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			})
	} else if potionType == proto.Potions_IronshieldPotion {
		actionID := ActionID{ItemID: 22849}
		aura := character.NewTemporaryStatsAura("Ironshield Potion", actionID, stats.Stats{stats.Armor: 2500}, time.Minute*2)
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
						Timer:    potionCD,
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			})
	} else if potionType == proto.Potions_HeroicPotion {
		actionID := ActionID{ItemID: 22837}
		aura := character.NewTemporaryStatsAura("Heroic Potion", actionID, stats.Stats{stats.Strength: 70, stats.Health: 700}, time.Second*15)
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
						Timer:    potionCD,
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			})
	} else {
		return MajorCooldown{}, nil
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

				remainingUsages := int(1 + (MaxDuration(0, sim.Duration))/(time.Minute*2))

				if consumes.DefaultConjured == proto.Conjured_ConjuredDarkRune {
					character.ExpectedBonusMana += expectedManaPerUsage * float64(remainingUsages)
				}

				return func(sim *Simulation, character *Character) {
					defaultSpell.Cast(sim, nil)

					if consumes.DefaultConjured == proto.Conjured_ConjuredDarkRune {
						// Update expected bonus mana
						newRemainingUsages := int(sim.GetRemainingDuration() / (time.Minute * 2))
						character.ExpectedBonusMana -= expectedManaPerUsage * float64(remainingUsages-newRemainingUsages)
						remainingUsages = newRemainingUsages
					}
				}
			},
		})
	}
}

func makeConjuredActivation(conjuredType proto.Conjured, character *Character) (MajorCooldown, *Spell) {
	if conjuredType == proto.Conjured_ConjuredDarkRune {
		actionID := ActionID{ItemID: 20520}
		manaMetrics := character.NewManaMetrics(actionID)
		damageTakenManaMetrics := character.NewManaMetrics(ActionID{SpellID: 33776})
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
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					// Restores 900 to 1500 mana. (2 Min Cooldown)
					manaGain := 900 + (sim.RandomFloat("dark rune") * 600)
					character.AddMana(sim, manaGain, manaMetrics, true)

					if character.Class == proto.Class_ClassPaladin {
						// Paladins gain extra mana from self-inflicted damage
						// TO-DO: It is possible for damage to be resisted or to crit
						// This would affect mana returns for Paladins
						manaFromDamage := manaGain * 2.0 / 3.0 * 0.1
						character.AddMana(sim, manaFromDamage, damageTakenManaMetrics, false)
					}
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

// Mana-related consumes need to know this so they can wait for the Mage to use
// gems first.
var MageManaGemMCDActionID = ActionID{ItemID: 22044}

var SuperSapperActionID = ActionID{ItemID: 23827}
var GoblinSapperActionID = ActionID{ItemID: 10646}
var FelIronBombActionID = ActionID{ItemID: 23736}
var AdamantiteGrenadeActionID = ActionID{ItemID: 23737}
var HolyWaterActionID = ActionID{ItemID: 13180}

func registerExplosivesCD(agent Agent, consumes proto.Consumes) {
	if !consumes.SuperSapper && !consumes.GoblinSapper && consumes.FillerExplosive == proto.Explosive_ExplosiveUnknown {
		return
	}
	character := agent.GetCharacter()
	explosivesTimer := character.NewTimer()

	var superSapper *Spell
	if consumes.SuperSapper {
		superSapper = character.newSuperSapperSpell()
	}
	var goblinSapper *Spell
	if consumes.GoblinSapper {
		goblinSapper = character.newGoblinSapperSpell()
	}

	var fillerExplosive *Spell
	switch consumes.FillerExplosive {
	case proto.Explosive_ExplosiveFelIronBomb:
		fillerExplosive = character.newFelIronBombSpell()
	case proto.Explosive_ExplosiveAdamantiteGrenade:
		fillerExplosive = character.newAdamantiteGrenadeSpell()
	case proto.Explosive_ExplosiveGnomishFlameTurret:
		fillerExplosive = character.newGnomishFlameTurretSpell()
	case proto.Explosive_ExplosiveHolyWater:
		fillerExplosive = character.newHolyWaterSpell()
	}

	cdAfterGoblinSapper := time.Minute
	if fillerExplosive == nil {
		if consumes.SuperSapper {
			cdAfterGoblinSapper = time.Minute * 4
		} else {
			cdAfterGoblinSapper = time.Minute * 5
		}
	}

	cdAfterSuperSapper := time.Minute
	if fillerExplosive == nil && !consumes.GoblinSapper {
		cdAfterSuperSapper = time.Minute * 5
	}

	spell := character.RegisterSpell(SpellConfig{
		ActionID: SuperSapperActionID,
		Flags:    SpellFlagNoOnCastComplete | SpellFlagNoMetrics | SpellFlagNoLogs,

		Cast: CastConfig{
			CD: Cooldown{
				Timer:    explosivesTimer,
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *Simulation, target *Unit, _ *Spell) {
			if superSapper != nil && superSapper.IsReady(sim) {
				superSapper.Cast(sim, target)
				explosivesTimer.Set(sim.CurrentTime + cdAfterSuperSapper)
			} else if goblinSapper != nil && goblinSapper.IsReady(sim) {
				goblinSapper.Cast(sim, target)
				explosivesTimer.Set(sim.CurrentTime + cdAfterGoblinSapper)
			} else {
				fillerExplosive.Cast(sim, target)
				explosivesTimer.Set(sim.CurrentTime + time.Minute)
			}
		},
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell: spell,
		Type:  CooldownTypeDPS,
	})
}

// Creates a spell object for the common explosive case.
func (character *Character) newBasicExplosiveSpellConfig(actionID ActionID, minDamage float64, maxDamage float64, cooldown Cooldown, isHolyWater bool, minSelfDamage float64, maxSelfDamage float64) SpellConfig {
	school := SpellSchoolFire
	damageMultiplier := 1.0
	if isHolyWater {
		school = SpellSchoolHoly
		if character.CurrentTarget.MobType != proto.MobType_MobTypeUndead {
			damageMultiplier = 0
		}
	}

	damageTakenManaMetrics := character.NewManaMetrics(ActionID{SpellID: 33776})
	return SpellConfig{
		ActionID:    actionID,
		SpellSchool: school,

		Cast: CastConfig{
			CD: cooldown,
		},

		ApplyEffects: ApplyEffectFuncAOEDamage(character.Env, SpellEffect{
			ProcMask: ProcMaskEmpty,
			// Explosives always have 1% resist chance, so just give them hit cap.
			BonusSpellHitRating: 100 * SpellHitRatingPerHitChance,

			DamageMultiplier: damageMultiplier,
			ThreatMultiplier: 1,

			BaseDamage:     BaseDamageConfigRoll(minDamage, maxDamage),
			OutcomeApplier: character.OutcomeFuncMagicHitAndCrit(2),
			OnSpellHitDealt: func(sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
				// Paladins gain extra mana from self-inflicted damage
				// TO-DO: Check if self-inflicted damage can be resisted or crit
				// This affects mana returns for Paladins
				if character.Class == proto.Class_ClassPaladin && maxSelfDamage > 0 {
					manaGain := (minSelfDamage + (sim.RandomFloat("sapper paladin") * (maxSelfDamage - minSelfDamage))) * 0.1
					character.AddMana(sim, manaGain, damageTakenManaMetrics, false)
				}
			},
		}),
	}
}
func (character *Character) newSuperSapperSpell() *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(SuperSapperActionID, 900, 1500, Cooldown{Timer: character.NewTimer(), Duration: time.Minute * 5}, false, 675, 1125))
}
func (character *Character) newGoblinSapperSpell() *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(GoblinSapperActionID, 450, 750, Cooldown{Timer: character.NewTimer(), Duration: time.Minute * 5}, false, 375, 625))
}
func (character *Character) newFelIronBombSpell() *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(FelIronBombActionID, 330, 770, Cooldown{}, false, 0, 0))
}
func (character *Character) newAdamantiteGrenadeSpell() *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(AdamantiteGrenadeActionID, 450, 750, Cooldown{}, false, 0, 0))
}
func (character *Character) newHolyWaterSpell() *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(HolyWaterActionID, 438, 562, Cooldown{}, true, 0, 0))
}
