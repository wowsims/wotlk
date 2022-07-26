package core

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Applies buffs that affect individual players.
func applyBuffEffects(agent Agent, raidBuffs proto.RaidBuffs, partyBuffs proto.PartyBuffs, individualBuffs proto.IndividualBuffs) {
	character := agent.GetCharacter()

	if raidBuffs.ArcaneBrilliance || raidBuffs.FelIntelligence > 0 {
		val := GetTristateValueFloat(raidBuffs.FelIntelligence, 48.0, 48.0*1.1)
		if raidBuffs.ArcaneBrilliance {
			val = 60.0
		}
		character.AddStat(stats.Intellect, val)
	} else if raidBuffs.ScrollOfIntellect {
		character.AddStats(stats.Stats{
			stats.Intellect: 48,
		})
	}

	gotwAmount := GetTristateValueFloat(raidBuffs.GiftOfTheWild, 54, 75)
	if gotwAmount > 0 {
		character.AddStats(stats.Stats{
			stats.Armor:     GetTristateValueFloat(raidBuffs.GiftOfTheWild, 750, 1050),
			stats.Stamina:   gotwAmount,
			stats.Agility:   gotwAmount,
			stats.Strength:  gotwAmount,
			stats.Intellect: gotwAmount,
			stats.Spirit:    gotwAmount,
		})
	} else if raidBuffs.DrumsOfTheWild {
		character.AddStats(stats.Stats{
			stats.Armor:     750,
			stats.Stamina:   37,
			stats.Agility:   37,
			stats.Strength:  37,
			stats.Intellect: 37,
			stats.Spirit:    37,
		})
	}

	if raidBuffs.Thorns == proto.TristateEffect_TristateEffectImproved {
		ThornsAura(character, 3)
	} else if raidBuffs.Thorns == proto.TristateEffect_TristateEffectRegular {
		ThornsAura(character, 0)
	}

	if raidBuffs.MoonkinAura > 0 || raidBuffs.ElementalOath {
		character.AddStat(stats.SpellCrit, 5*CritRatingPerCritChance)
	}
	if raidBuffs.MoonkinAura == proto.TristateEffect_TristateEffectImproved || raidBuffs.SwiftRetribution {
		character.PseudoStats.CastSpeedMultiplier *= 1.03
		character.PseudoStats.MeleeSpeedMultiplier *= 1.03
		character.PseudoStats.RangedSpeedMultiplier *= 1.03
	}

	if raidBuffs.LeaderOfThePack > 0 || raidBuffs.Rampage {
		character.AddStats(stats.Stats{
			stats.MeleeCrit: 5 * CritRatingPerCritChance,
		})
		if raidBuffs.LeaderOfThePack == proto.TristateEffect_TristateEffectImproved {
			// TODO: healing aura from imp LotP
		}
	}

	if raidBuffs.TrueshotAura || raidBuffs.AbominationsMight || raidBuffs.UnleashedRage {
		// Increases AP by 10%
		character.AddStatDependency(stats.AttackPower, stats.AttackPower, 1.0+0.1)
	}

	if raidBuffs.ArcaneEmpowerment || raidBuffs.FerociousInspiration || raidBuffs.SanctifiedRetribution {
		character.PseudoStats.DamageDealtMultiplier *= 1.03
	}

	if partyBuffs.HeroicPresence {
		character.AddStats(stats.Stats{
			stats.MeleeHit: 1 * MeleeHitRatingPerHitChance,
			stats.SpellHit: 1 * SpellHitRatingPerHitChance,
		})
	}

	if raidBuffs.BloodPact > 0 || raidBuffs.CommandingShout > 0 {
		health := GetTristateValueFloat(raidBuffs.BloodPact, 1330, 1330*1.3)
		health2 := GetTristateValueFloat(raidBuffs.CommandingShout, 2255, 2255*1.25)
		character.AddStat(stats.Health, MaxFloat(health, health2))
	}

	if raidBuffs.PowerWordFortitude != proto.TristateEffect_TristateEffectMissing {
		character.AddStats(stats.Stats{
			stats.Stamina: GetTristateValueFloat(raidBuffs.PowerWordFortitude, 165, 165*1.3),
		})
	} else if raidBuffs.ScrollOfStamina {
		character.AddStats(stats.Stats{
			stats.Stamina: 132,
		})
	}
	if raidBuffs.ShadowProtection {
		character.AddStats(stats.Stats{
			stats.ShadowResistance: 130,
		})
	}
	if raidBuffs.DivineSpirit || raidBuffs.FelIntelligence > 0 {
		v := GetTristateValueFloat(raidBuffs.FelIntelligence, 64.0, 64.0*1.1)
		if raidBuffs.DivineSpirit {
			v = 80.0
		}
		character.AddStats(stats.Stats{
			stats.Spirit: v,
		})
	} else if raidBuffs.ScrollOfSpirit {
		character.AddStats(stats.Stats{
			stats.Spirit: 64,
		})
	}

	// TODO: convert this to a real mana replenishment aura we can use in raid sim.
	if individualBuffs.VampiricTouch ||
		individualBuffs.HuntingParty ||
		individualBuffs.JudgementsOfTheWise ||
		individualBuffs.ImprovedSoulLeech ||
		individualBuffs.EnduringWinter {
		character.AddStatDependency(stats.Mana, stats.MP5, 1.0+0.01)
	}

	kingsAgiIntSpiAmount := 1.0
	kingsStrStamAmount := 1.0
	if individualBuffs.BlessingOfSanctuary {
		kingsStrStamAmount = 1.1
	}
	if individualBuffs.BlessingOfKings {
		kingsAgiIntSpiAmount = 1.1
		kingsStrStamAmount = 1.1
	} else if raidBuffs.DrumsOfForgottenKings {
		kingsAgiIntSpiAmount = 1.08
		kingsStrStamAmount = MaxFloat(kingsStrStamAmount, 1.08)
	}
	if kingsStrStamAmount > 0 {
		character.AddStatDependency(stats.Strength, stats.Strength, kingsStrStamAmount)
		character.AddStatDependency(stats.Stamina, stats.Stamina, kingsStrStamAmount)
	}
	if kingsAgiIntSpiAmount > 0 {
		character.AddStatDependency(stats.Agility, stats.Agility, kingsAgiIntSpiAmount)
		character.AddStatDependency(stats.Intellect, stats.Intellect, kingsAgiIntSpiAmount)
		character.AddStatDependency(stats.Spirit, stats.Spirit, kingsAgiIntSpiAmount)
	}

	if individualBuffs.BlessingOfSanctuary {
		character.PseudoStats.DamageTakenMultiplier *= 0.97
		BlessingOfSanctuaryAura(character)
	} else if individualBuffs.Vigilance || individualBuffs.RenewedHope {
		character.PseudoStats.DamageTakenMultiplier *= 0.97
	}

	// TODO: Is scroll exclusive to totem?
	if raidBuffs.StoneskinTotem != proto.TristateEffect_TristateEffectMissing {
		character.AddStats(stats.Stats{
			stats.Armor: GetTristateValueFloat(raidBuffs.StoneskinTotem, 1150, 1380),
		})
	}

	if raidBuffs.DevotionAura != proto.TristateEffect_TristateEffectMissing {
		character.AddStats(stats.Stats{
			stats.Armor: GetTristateValueFloat(raidBuffs.DevotionAura, 861, 1205),
		})
	} else if raidBuffs.ScrollOfProtection {
		character.AddStats(stats.Stats{
			stats.Armor: 750,
		})
	}

	if raidBuffs.RetributionAura {
		RetributionAura(character, raidBuffs.SanctifiedRetribution)
	}

	if raidBuffs.BattleShout > 0 || individualBuffs.BlessingOfMight > 0 {
		bonusAP := 550 * GetTristateValueFloat(raidBuffs.BattleShout, 1, 1.25)
		bomAP := 550 * GetTristateValueFloat(individualBuffs.BlessingOfMight, 1, 1.25)
		if bomAP > bonusAP {
			bonusAP = bomAP
		}
		character.AddStats(stats.Stats{
			stats.AttackPower: math.Floor(bomAP),
		})
	}
	character.AddStats(stats.Stats{
		stats.Health: GetTristateValueFloat(raidBuffs.CommandingShout, 1080, 1080*1.25),
	})

	spBonus := float64(raidBuffs.DemonicPact)
	if raidBuffs.TotemOfWrath {
		spBonus = MaxFloat(spBonus, 280)
	} else if raidBuffs.FlametongueTotem {
		spBonus = MaxFloat(spBonus, 144)
	}
	if spBonus > 0 {
		character.AddStats(stats.Stats{
			stats.SpellPower:   spBonus,
			stats.HealingPower: spBonus,
		})
	}
	if raidBuffs.WrathOfAirTotem {
		character.PseudoStats.CastSpeedMultiplier *= 1.05
	}
	if raidBuffs.StrengthOfEarthTotem > 0 || raidBuffs.HornOfWinter {
		val := MaxTristate(proto.TristateEffect_TristateEffectRegular, raidBuffs.StrengthOfEarthTotem)
		bonus := GetTristateValueFloat(val, 155, 186)
		character.AddStats(stats.Stats{
			stats.Strength: bonus,
			stats.Agility:  bonus,
		})
	} else {
		if raidBuffs.ScrollOfStrength {
			character.AddStats(stats.Stats{
				stats.Strength: 30,
			})
		}
		if raidBuffs.ScrollOfAgility {
			character.AddStats(stats.Stats{
				stats.Agility: 30,
			})
		}
	}

	if individualBuffs.BlessingOfWisdom > 0 || raidBuffs.ManaSpringTotem > 0 {
		character.AddStats(stats.Stats{
			stats.MP5: GetTristateValueFloat(MaxTristate(individualBuffs.BlessingOfWisdom, raidBuffs.ManaSpringTotem), 91, 109),
		})
	}

	if raidBuffs.IcyTalons {
		character.PseudoStats.MeleeSpeedMultiplier *= 1.2
	} else if raidBuffs.WindfuryTotem > 0 {
		character.PseudoStats.MeleeSpeedMultiplier *= GetTristateValueFloat(raidBuffs.WindfuryTotem, 1.16, 1.20)
	}

	if raidBuffs.Bloodlust {
		registerBloodlustCD(agent)
	}

	registerPowerInfusionCD(agent, individualBuffs.PowerInfusions)
	registerManaTideTotemCD(agent, partyBuffs.ManaTideTotems)
	registerInnervateCD(agent, individualBuffs.Innervates)

	character.AddStats(stats.Stats{
		stats.SpellCrit: 28 * float64(partyBuffs.AtieshMage),
	})
	character.AddStats(stats.Stats{
		stats.SpellPower:   33 * float64(partyBuffs.AtieshWarlock),
		stats.HealingPower: 33 * float64(partyBuffs.AtieshWarlock),
	})

	if partyBuffs.BraidedEterniumChain {
		character.AddStats(stats.Stats{stats.MeleeCrit: 28})
	}
	if partyBuffs.EyeOfTheNight {
		character.AddStats(stats.Stats{stats.SpellPower: 34})
	}
	if partyBuffs.ChainOfTheTwilightOwl {
		character.AddStats(stats.Stats{stats.SpellCrit: 2 * CritRatingPerCritChance})
	}
}

// Applies buffs to pets.
func applyPetBuffEffects(petAgent PetAgent, raidBuffs proto.RaidBuffs, partyBuffs proto.PartyBuffs, individualBuffs proto.IndividualBuffs) {
	// Summoned pets, like Mage Water Elemental, aren't around to receive raid buffs.
	if !petAgent.GetPet().PermanentPet {
		return
	}

	// We need to modify the buffs a bit because some things are applied to pets by
	// the owner during combat (Bloodlust) or don't make sense for a pet.
	raidBuffs.Bloodlust = false
	individualBuffs.Innervates = 0
	individualBuffs.PowerInfusions = 0

	// For some reason pets don't benefit from buffs that are ratings, e.g. crit rating or haste rating.
	partyBuffs.BraidedEterniumChain = false

	applyBuffEffects(petAgent, raidBuffs, partyBuffs, individualBuffs)
}

func applyInspiration(character *Character, uptime float64) {
	if uptime <= 0 {
		return
	}
	uptime = MinFloat(1, uptime)

	var curBonus stats.Stats
	inspirationAura := character.RegisterAura(Aura{
		Label:    "Inspiration",
		ActionID: ActionID{SpellID: 15363},
		Duration: time.Second * 15,
		OnGain: func(aura *Aura, sim *Simulation) {
			curBonus = stats.Stats{stats.Armor: character.GetStat(stats.Armor) * 0.25}
			aura.Unit.AddStatsDynamic(sim, curBonus)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, curBonus.Multiply(-1))
		},
	})

	auraDuration := time.Second * 15
	tickLength := time.Millisecond * 2500
	ticksPerAura := float64(auraDuration) / float64(tickLength)
	chancePerTick := TernaryFloat64(uptime == 1, 1, 1.0-math.Pow(1-uptime, 1/ticksPerAura))

	character.RegisterResetEffect(func(sim *Simulation) {
		StartPeriodicAction(sim, PeriodicActionOptions{
			Period: tickLength,
			OnAction: func(sim *Simulation) {
				if sim.RandomFloat("Inspiration") < chancePerTick {
					inspirationAura.Activate(sim)
				}
			},
		})

		// Also try once at the start.
		StartPeriodicAction(sim, PeriodicActionOptions{
			Period:   1,
			NumTicks: 1,
			OnAction: func(sim *Simulation) {
				if sim.RandomFloat("Inspiration") < uptime {
					// Use random duration to compensate for increased chance collapsed into single tick.
					randomDur := tickLength + time.Duration(float64(auraDuration-tickLength)*sim.RandomFloat("InspirationDur"))
					inspirationAura.Duration = randomDur
					inspirationAura.Activate(sim)
					inspirationAura.Duration = time.Second * 15
				}
			},
		})
	})
}

func RetributionAura(character *Character, sanctifiedRetribution bool) *Aura {
	actionID := ActionID{SpellID: 54043}

	damage := 112.0
	if sanctifiedRetribution {
		damage *= 1.5
	}

	procSpell := character.RegisterSpell(SpellConfig{
		ActionID:    actionID,
		SpellSchool: SpellSchoolHoly,
		Flags:       SpellFlagBinary,

		ApplyEffects: ApplyEffectFuncDirectDamage(SpellEffect{
			ProcMask:         ProcMaskEmpty,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage:     BaseDamageConfigFlat(damage),
			OutcomeApplier: character.OutcomeFuncMagicHitBinary(),
		}),
	})

	return character.RegisterAura(Aura{
		Label:    "Retribution Aura",
		ActionID: actionID,
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			if spellEffect.Landed() && spell.SpellSchool == SpellSchoolPhysical {
				procSpell.Cast(sim, spell.Unit)
			}
		},
	})
}

func ThornsAura(character *Character, points int32) *Aura {
	actionID := ActionID{SpellID: 53307}

	procSpell := character.RegisterSpell(SpellConfig{
		ActionID:    actionID,
		SpellSchool: SpellSchoolNature,
		Flags:       SpellFlagBinary,

		ApplyEffects: ApplyEffectFuncDirectDamage(SpellEffect{
			ProcMask:         ProcMaskEmpty,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage:     BaseDamageConfigFlat(73 * (1 + 0.25*float64(points))),
			OutcomeApplier: character.OutcomeFuncMagicHitBinary(),
		}),
	})

	return character.RegisterAura(Aura{
		Label:    "Thorns",
		ActionID: actionID,
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			if spellEffect.Landed() && spell.SpellSchool == SpellSchoolPhysical {
				procSpell.Cast(sim, spell.Unit)
			}
		},
	})
}

func BlessingOfSanctuaryAura(character *Character) {
	if !character.HasManaBar() {
		return
	}
	actionID := ActionID{SpellID: 25899}
	manaMetrics := character.NewManaMetrics(actionID)

	character.RegisterAura(Aura{
		Label:    "Blessing of Sanctuary",
		ActionID: actionID,
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			if spellEffect.Outcome.Matches(OutcomeBlock | OutcomeDodge | OutcomeParry) {
				character.AddMana(sim, 0.02*character.MaxMana(), manaMetrics, false)
			}
		},
	})
}

var (
	WindfuryTotemSpellRanks = []int32{
		8512,
		10613,
		10614,
		25585,
		25587,
	}

	windfuryBuffSpellRanks = []int32{
		8516,
		10608,
		10610,
		25583,
		25584,
	}

	windfuryAPBonuses = []float64{
		122,
		229,
		315,
		375,
		445,
	}
)

func IsEligibleForWindfuryTotem(character *Character) bool {
	return character.AutoAttacks.IsEnabled() &&
		character.HasMHWeapon() &&
		!character.HasMHWeaponImbue
}

var WindfuryTotemAuraLabel = "Windfury Totem"

func WindfuryTotemAura(character *Character, rank int32, iwtTalentPoints int32) *Aura {
	buffActionID := ActionID{SpellID: windfuryBuffSpellRanks[rank-1]}
	apBonus := windfuryAPBonuses[rank-1]
	apBonus *= 1 + 0.15*float64(iwtTalentPoints)

	var charges int32

	wfBuffAura := character.NewTemporaryStatsAuraWrapped("Windfury Buff", buffActionID, stats.Stats{stats.AttackPower: apBonus}, time.Millisecond*1500, func(config *Aura) {
		config.OnSpellHitDealt = func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			// *Special Case* Windfury should not proc on Seal of Command
			if spell.ActionID.SpellID == 20424 {
				return
			}
			if !spellEffect.ProcMask.Matches(ProcMaskMeleeWhiteHit) || spellEffect.ProcMask.Matches(ProcMaskMeleeSpecial) {
				return
			}
			charges--
			if charges == 0 {
				aura.Deactivate(sim)
			}
		}
	})

	var wfSpell *Spell
	icd := Cooldown{
		Timer:    character.NewTimer(),
		Duration: 1,
	}
	const procChance = 0.2

	return character.RegisterAura(Aura{
		Label:    WindfuryTotemAuraLabel,
		Duration: NeverExpires,
		OnInit: func(aura *Aura, sim *Simulation) {
			wfSpell = character.GetOrRegisterSpell(SpellConfig{
				ActionID:    buffActionID, // temporary buff ("Windfury Attack") spell id
				SpellSchool: SpellSchoolPhysical,
				Flags:       SpellFlagMeleeMetrics | SpellFlagNoOnCastComplete,

				ApplyEffects: ApplyEffectFuncDirectDamage(character.AutoAttacks.MHEffect),
			})
		},
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			// *Special Case* Windfury should not proc on Seal of Command
			if spell.ActionID.SpellID == 20424 {
				return
			}
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(ProcMaskMeleeMHAuto) {
				return
			}

			if wfBuffAura.IsActive() {
				return
			}
			if !icd.IsReady(sim) {
				// Checking for WF buff aura isn't quite enough now that we refactored auras.
				// TODO: Clean this up to remove the need for an instant ICD.
				return
			}

			if sim.RandomFloat("Windfury Totem") > procChance {
				return
			}

			// TODO: the current proc system adds auras after cast and damage, in game they're added after cast
			startCharges := int32(2)
			if !spellEffect.ProcMask.Matches(ProcMaskMeleeMHSpecial) {
				startCharges--
			}
			charges = startCharges
			wfBuffAura.Activate(sim)
			icd.Use(sim)

			aura.Unit.AutoAttacks.MaybeReplaceMHSwing(sim, wfSpell).Cast(sim, spellEffect.Target)
		},
	})
}

// Used for approximating cooldowns applied by other players to you, such as
// bloodlust, innervate, power infusion, etc. This is specifically for buffs
// which can be consecutively applied multiple times to a single player.
type externalConsecutiveCDApproximation struct {
	ActionID         ActionID
	AuraTag          string
	CooldownPriority float64
	Type             CooldownType
	AuraDuration     time.Duration
	AuraCD           time.Duration

	// Callback for extra activation conditions.
	ShouldActivate CooldownActivationCondition

	// Applies the buff.
	AddAura CooldownActivation
}

// numSources is the number of other players assigned to apply the buff to this player.
// E.g. the number of other shaman in the group using bloodlust.
func registerExternalConsecutiveCDApproximation(agent Agent, config externalConsecutiveCDApproximation, numSources int32) {
	if numSources == 0 {
		panic("Need at least 1 source!")
	}
	character := agent.GetCharacter()

	var nextExternalIndex int

	externalTimers := make([]*Timer, numSources)
	for i := 0; i < int(numSources); i++ {
		externalTimers[i] = character.NewTimer()
	}
	sharedTimer := character.NewTimer()

	spell := character.RegisterSpell(SpellConfig{
		ActionID: config.ActionID,
		Flags:    SpellFlagNoOnCastComplete | SpellFlagNoMetrics | SpellFlagNoLogs,

		Cast: CastConfig{
			CD: Cooldown{
				Timer:    sharedTimer,
				Duration: config.AuraDuration, // Assumes that multiple buffs are different sources.
			},
		},

		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			config.AddAura(sim, character)
			externalTimers[nextExternalIndex].Set(sim.CurrentTime + config.AuraCD)

			nextExternalIndex = (nextExternalIndex + 1) % len(externalTimers)

			if externalTimers[nextExternalIndex].IsReady(sim) {
				sharedTimer.Set(sim.CurrentTime + config.AuraDuration)
			} else {
				sharedTimer.Set(sim.CurrentTime + externalTimers[nextExternalIndex].TimeToReady(sim))
			}
		},
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell:    spell,
		Priority: config.CooldownPriority,
		Type:     config.Type,

		CanActivate: func(sim *Simulation, character *Character) bool {
			if !externalTimers[nextExternalIndex].IsReady(sim) {
				return false
			}

			if character.HasActiveAuraWithTag(config.AuraTag) {
				return false
			}

			return true
		},
		ShouldActivate: config.ShouldActivate,
	})
}

const BloodlustAuraTag = "Bloodlust"

const BloodlustDuration = time.Second * 40
const BloodlustCD = time.Minute * 10

func registerBloodlustCD(agent Agent) {
	bloodlustAura := BloodlustAura(agent.GetCharacter(), -1)

	// TODO: do we need consecutive CDs
	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 2825, Tag: -1},
			AuraTag:          BloodlustAuraTag,
			CooldownPriority: CooldownPriorityBloodlust,
			AuraDuration:     BloodlustDuration,
			AuraCD:           BloodlustCD,
			Type:             CooldownTypeDPS | CooldownTypeUsableShapeShifted,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Haste portion doesn't stack with Power Infusion, so prefer to wait.
				return !character.HasActiveAuraWithTag(PowerInfusionAuraTag)
			},
			AddAura: func(sim *Simulation, character *Character) { bloodlustAura.Activate(sim) },
		}, 1)
}

func BloodlustAura(character *Character, actionTag int32) *Aura {
	const bonus = 1.3
	const inverseBonus = 1 / bonus
	actionID := ActionID{SpellID: 2825, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "Bloodlust-" + actionID.String(),
		Tag:      BloodlustAuraTag,
		ActionID: actionID,
		Duration: BloodlustDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			if character.HasActiveAuraWithTag(PowerInfusionAuraTag) {
				character.MultiplyCastSpeed(1 / 1.2)
			}
			character.MultiplyCastSpeed(bonus)
			character.MultiplyAttackSpeed(sim, bonus)

			if len(character.Pets) > 0 {
				for _, petAgent := range character.Pets {
					pet := petAgent.GetPet()
					if pet.IsEnabled() && pet.IsPermanent() {
						BloodlustAura(&pet.Character, actionTag).Activate(sim)
					}
				}
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if character.HasActiveAuraWithTag(PowerInfusionAuraTag) {
				character.MultiplyCastSpeed(1.2)
			}
			character.MultiplyCastSpeed(inverseBonus)
			character.MultiplyAttackSpeed(sim, inverseBonus)
		},
	})
}

var PowerInfusionAuraTag = "PowerInfusion"

const PowerInfusionDuration = time.Second * 15
const PowerInfusionCD = time.Minute * 3

func registerPowerInfusionCD(agent Agent, numPowerInfusions int32) {
	if numPowerInfusions == 0 {
		return
	}

	piAura := PowerInfusionAura(agent.GetCharacter(), -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 10060, Tag: -1},
			AuraTag:          PowerInfusionAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     PowerInfusionDuration,
			AuraCD:           PowerInfusionCD,
			Type:             CooldownTypeDPS | CooldownTypeUsableShapeShifted,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Haste portion doesn't stack with Bloodlust, so prefer to wait.
				return !character.HasActiveAuraWithTag(BloodlustAuraTag)
			},
			AddAura: func(sim *Simulation, character *Character) { piAura.Activate(sim) },
		},
		numPowerInfusions)
}

func PowerInfusionAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 10060, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "PowerInfusion-" + actionID.String(),
		Tag:      PowerInfusionAuraTag,
		ActionID: actionID,
		Duration: PowerInfusionDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			if character.HasManaBar() {
				// TODO: Double-check this is how the calculation works.
				character.PseudoStats.CostMultiplier *= 0.8
			}
			if !character.HasActiveAuraWithTag(BloodlustAuraTag) {
				character.MultiplyCastSpeed(1.2)
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if character.HasManaBar() {
				character.PseudoStats.CostMultiplier /= 0.8
			}
			if !character.HasActiveAuraWithTag(BloodlustAuraTag) {
				character.MultiplyCastSpeed(1 / 1.2)
			}
		},
	})
}

var InnervateAuraTag = "Innervate"

const InnervateDuration = time.Second * 20
const InnervateCD = time.Minute * 6

func InnervateManaThreshold(character *Character) float64 {
	if character.Class == proto.Class_ClassMage {
		// Mages burn mana really fast so they need a higher threshold.
		return character.MaxMana() * 0.7
	} else {
		return 1000
	}
}

func registerInnervateCD(agent Agent, numInnervates int32) {
	if numInnervates == 0 {
		return
	}

	innervateThreshold := 0.0
	expectedManaPerInnervate := 0.0
	remainingInnervateUsages := 0
	var innervateAura *Aura

	character := agent.GetCharacter()
	character.Env.RegisterPostFinalizeEffect(func() {
		innervateThreshold = InnervateManaThreshold(character)
		expectedManaPerInnervate = character.SpiritManaRegenPerSecond() * 5 * 20
		remainingInnervateUsages = int(1 + (MaxDuration(0, character.Env.BaseDuration))/InnervateCD)
		character.ExpectedBonusMana += expectedManaPerInnervate * float64(remainingInnervateUsages)
		innervateAura = InnervateAura(character, expectedManaPerInnervate, -1)
	})

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 29166, Tag: -1},
			AuraTag:          InnervateAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     InnervateDuration,
			AuraCD:           InnervateCD,
			Type:             CooldownTypeMana | CooldownTypeUsableShapeShifted,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only cast innervate when very low on mana, to make sure all other mana CDs are prioritized.
				if character.CurrentMana() > innervateThreshold {
					return false
				}
				return true
			},
			AddAura: func(sim *Simulation, character *Character) {
				innervateAura.Activate(sim)

				newRemainingUsages := int(sim.GetRemainingDuration() / InnervateCD)
				// AddInnervateAura already accounts for 1 usage, which is why we subtract 1 less.
				character.ExpectedBonusMana -= expectedManaPerInnervate * MaxFloat(0, float64(remainingInnervateUsages-newRemainingUsages-1))
				remainingInnervateUsages = newRemainingUsages
			},
		},
		numInnervates)
}

func InnervateAura(character *Character, expectedBonusManaReduction float64, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 29166, Tag: actionTag}
	return character.GetOrRegisterAura(Aura{
		Label:    "Innervate-" + actionID.String(),
		Tag:      InnervateAuraTag,
		ActionID: actionID,
		Duration: InnervateDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.ForceFullSpiritRegen = true
			character.PseudoStats.SpiritRegenMultiplier *= 5.0
			character.UpdateManaRegenRates()

			expectedBonusManaPerTick := expectedBonusManaReduction / 10
			StartPeriodicAction(sim, PeriodicActionOptions{
				Period:   InnervateDuration / 10,
				NumTicks: 10,
				OnAction: func(sim *Simulation) {
					character.ExpectedBonusMana -= expectedBonusManaPerTick
					character.Metrics.BonusManaGained += expectedBonusManaPerTick
				},
			})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.ForceFullSpiritRegen = false
			character.PseudoStats.SpiritRegenMultiplier /= 5.0
			character.UpdateManaRegenRates()
		},
	})
}

var ManaTideTotemAuraTag = "ManaTideTotem"

const ManaTideTotemDuration = time.Second * 12
const ManaTideTotemCD = time.Minute * 5

func ManaTideTotemAmount(character *Character) float64 {
	// Subtract 120 mana to simulate the loss of mana spring while MTT is active.
	// This isn't correct for multi-resto shaman groups, but that isnt a common case.
	return character.MaxMana()*0.24 - 120
}

func registerManaTideTotemCD(agent Agent, numManaTideTotems int32) {
	if numManaTideTotems == 0 {
		return
	}

	expectedManaPerManaTideTotem := 0.0
	remainingManaTideTotemUsages := 0
	initialDelay := time.Duration(0)
	var mttAura *Aura

	character := agent.GetCharacter()
	character.Env.RegisterPostFinalizeEffect(func() {
		// Use first MTT at 60s, or halfway through the fight, whichever comes first.
		initialDelay = MinDuration(character.Env.BaseDuration/2, time.Second*60)

		expectedManaPerManaTideTotem = ManaTideTotemAmount(character)
		remainingManaTideTotemUsages = int(1 + MaxDuration(0, character.Env.BaseDuration-initialDelay)/ManaTideTotemCD)
		character.ExpectedBonusMana += expectedManaPerManaTideTotem * float64(remainingManaTideTotemUsages)
		mttAura = ManaTideTotemAura(character, -1)
	})

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 16190, Tag: -1},
			AuraTag:          ManaTideTotemAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     ManaTideTotemDuration,
			AuraCD:           ManaTideTotemCD,
			Type:             CooldownTypeMana | CooldownTypeUsableShapeShifted,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// A normal resto shaman would wait to use MTT.
				if sim.CurrentTime < initialDelay {
					return false
				}
				return true
			},
			AddAura: func(sim *Simulation, character *Character) {
				mttAura.Activate(sim)

				newRemainingUsages := int(sim.GetRemainingDuration() / ManaTideTotemCD)
				// AddManaTideTotemAura already accounts for 1 usage, which is why we subtract 1 less.
				character.ExpectedBonusMana -= expectedManaPerManaTideTotem * MaxFloat(0, float64(remainingManaTideTotemUsages-newRemainingUsages-1))
				remainingManaTideTotemUsages = newRemainingUsages
			},
		},
		numManaTideTotems)
}

func ManaTideTotemAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 16190, Tag: actionTag}

	var metrics *ResourceMetrics
	if character.HasManaBar() {
		metrics = character.NewManaMetrics(actionID)
	}

	return character.GetOrRegisterAura(Aura{
		Label:    "ManaTideTotem-" + actionID.String(),
		Tag:      ManaTideTotemAuraTag,
		ActionID: actionID,
		Duration: ManaTideTotemDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			if character.HasManaBar() {
				manaPerTick := ManaTideTotemAmount(character) / 4
				StartPeriodicAction(sim, PeriodicActionOptions{
					Period:   ManaTideTotemDuration / 4,
					NumTicks: 4,
					OnAction: func(sim *Simulation) {
						if metrics != nil {
							character.AddMana(sim, manaPerTick, metrics, true)
							character.ExpectedBonusMana -= manaPerTick
						}
					},
				})
			}
		},
	})
}
