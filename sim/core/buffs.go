package core

import (
	"math"
	"strconv"
	"time"

	"golang.org/x/exp/slices"
	googleProto "google.golang.org/protobuf/proto"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Applies buffs that affect individual players.
func applyBuffEffects(agent Agent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individualBuffs *proto.IndividualBuffs) {
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

	if raidBuffs.DrumsOfTheWild {
		raidBuffs.GiftOfTheWild = MaxTristate(raidBuffs.GiftOfTheWild, proto.TristateEffect_TristateEffectRegular)
	}
	gotwAmount := GetTristateValueFloat(raidBuffs.GiftOfTheWild, 37, 51)
	gotwArmorAmount := GetTristateValueFloat(raidBuffs.GiftOfTheWild, 750, 1050)
	gotwResistAmount := GetTristateValueFloat(raidBuffs.GiftOfTheWild, 54, 75)
	if gotwAmount > 0 {
		character.AddStats(stats.Stats{
			stats.Armor:            gotwArmorAmount,
			stats.Stamina:          gotwAmount,
			stats.Agility:          gotwAmount,
			stats.Strength:         gotwAmount,
			stats.Intellect:        gotwAmount,
			stats.Spirit:           gotwAmount,
			stats.ArcaneResistance: gotwResistAmount,
			stats.ShadowResistance: gotwResistAmount,
			stats.NatureResistance: gotwResistAmount,
			stats.FireResistance:   gotwResistAmount,
			stats.FrostResistance:  gotwResistAmount,
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
		// For now, we assume Improved Moonkin Form is maxed-out
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
		character.MultiplyStat(stats.AttackPower, 1.1)
		character.MultiplyStat(stats.RangedAttackPower, 1.1)
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
			stats.ShadowResistance: 130 - gotwResistAmount,
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

	var replenishmentActionID ActionID
	if individualBuffs.VampiricTouch {
		replenishmentActionID.SpellID = 48160
	} else if individualBuffs.HuntingParty {
		replenishmentActionID.SpellID = 53292
	} else if individualBuffs.JudgementsOfTheWise {
		replenishmentActionID.SpellID = 31878
	} else if individualBuffs.ImprovedSoulLeech {
		replenishmentActionID.SpellID = 54118
	} else if individualBuffs.EnduringWinter {
		replenishmentActionID.SpellID = 44561
	}
	if !(replenishmentActionID.IsEmptyAction()) {
		MakePermanent(replenishmentAura(&character.Unit, replenishmentActionID))
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
		character.MultiplyStat(stats.Strength, kingsStrStamAmount)
		character.MultiplyStat(stats.Stamina, kingsStrStamAmount)
	}
	if kingsAgiIntSpiAmount > 0 {
		character.MultiplyStat(stats.Agility, kingsAgiIntSpiAmount)
		character.MultiplyStat(stats.Intellect, kingsAgiIntSpiAmount)
		character.MultiplyStat(stats.Spirit, kingsAgiIntSpiAmount)
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
			stats.Armor: GetTristateValueFloat(raidBuffs.DevotionAura, 1205, 1807.5),
		})
	}

	if raidBuffs.ScrollOfProtection && raidBuffs.DevotionAura == proto.TristateEffect_TristateEffectMissing {
		character.AddStats(stats.Stats{
			stats.Armor: 750,
		})
	}

	if raidBuffs.RetributionAura {
		RetributionAura(character, raidBuffs.SanctifiedRetribution)
	}

	if raidBuffs.BattleShout > 0 || individualBuffs.BlessingOfMight > 0 {
		bonusAP := 550 * GetTristateValueFloat(MaxTristate(raidBuffs.BattleShout, individualBuffs.BlessingOfMight), 1, 1.25)
		character.AddStats(stats.Stats{
			stats.AttackPower:       math.Floor(bonusAP),
			stats.RangedAttackPower: math.Floor(bonusAP),
		})
	}

	if raidBuffs.FlametongueTotem {
		MakePermanent(FlametongueTotemAura(character))
	}
	if raidBuffs.TotemOfWrath {
		MakePermanent(TotemOfWrathAura(character))
	}
	if raidBuffs.DemonicPact > 0 {
		dpAura := DemonicPactAura(character)
		dpAura.ExclusiveEffects[0].Priority = float64(raidBuffs.DemonicPact) / 10.0
		MakePermanent(dpAura)
	}

	if raidBuffs.WrathOfAirTotem {
		character.PseudoStats.CastSpeedMultiplier *= 1.05
	}
	if raidBuffs.StrengthOfEarthTotem > 0 || raidBuffs.HornOfWinter {
		val := MaxTristate(proto.TristateEffect_TristateEffectRegular, raidBuffs.StrengthOfEarthTotem)
		bonus := GetTristateValueFloat(val, 155, 178)
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

	registerRevitalizeHotCD(agent, "Rejuvination", ActionID{SpellID: 26982}, 5, 3*time.Second, individualBuffs.RevitalizeRejuvination)
	registerRevitalizeHotCD(agent, "Wild Growth", ActionID{SpellID: 53251}, 7, time.Second, individualBuffs.RevitalizeWildGrowth)

	registerUnholyFrenzyCD(agent, individualBuffs.UnholyFrenzy)
	registerTricksOfTheTradeCD(agent, individualBuffs.TricksOfTheTrades)
	registerShatteringThrowCD(agent, individualBuffs.ShatteringThrows)
	registerPowerInfusionCD(agent, individualBuffs.PowerInfusions)
	registerManaTideTotemCD(agent, partyBuffs.ManaTideTotems)
	registerInnervateCD(agent, individualBuffs.Innervates)
	registerDivineGuardianCD(agent, individualBuffs.DivineGuardians)
	registerPainSuppressionCD(agent, individualBuffs.PainSuppressions)

	character.AddStats(stats.Stats{
		stats.SpellCrit: 28 * float64(partyBuffs.AtieshMage),
	})
	character.AddStats(stats.Stats{
		stats.SpellPower: 33 * float64(partyBuffs.AtieshWarlock),
	})

	if partyBuffs.BraidedEterniumChain {
		character.AddStats(stats.Stats{stats.MeleeCrit: 28, stats.SpellCrit: 28})
	}
	if partyBuffs.EyeOfTheNight {
		character.AddStats(stats.Stats{stats.SpellPower: 34})
	}
	if partyBuffs.ChainOfTheTwilightOwl {
		character.AddStats(stats.Stats{stats.MeleeCrit: 45, stats.SpellCrit: 45})
	}
	if individualBuffs.FocusMagic {
		FocusMagicAura(nil, character)
	}
}

// Applies buffs to pets.
func applyPetBuffEffects(petAgent PetAgent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individualBuffs *proto.IndividualBuffs) {
	// Summoned pets, like Mage Water Elemental, aren't around to receive raid buffs.
	if petAgent.GetPet().IsGuardian() {
		return
	}

	raidBuffs = googleProto.Clone(raidBuffs).(*proto.RaidBuffs)
	partyBuffs = googleProto.Clone(partyBuffs).(*proto.PartyBuffs)
	individualBuffs = googleProto.Clone(individualBuffs).(*proto.IndividualBuffs)

	// We need to modify the buffs a bit because some things are applied to pets by
	// the owner during combat (Bloodlust) or don't make sense for a pet.
	raidBuffs.Bloodlust = false
	raidBuffs.WrathOfAirTotem = false
	individualBuffs.HymnOfHope = 0
	individualBuffs.HandOfSalvation = 0
	individualBuffs.Innervates = 0
	individualBuffs.PowerInfusions = 0
	individualBuffs.UnholyFrenzy = 0
	individualBuffs.RevitalizeRejuvination = 0
	individualBuffs.RevitalizeWildGrowth = 0
	individualBuffs.TricksOfTheTrades = 0
	individualBuffs.ShatteringThrows = 0

	if !petAgent.GetPet().enabledOnStart {
		raidBuffs.ArcaneBrilliance = false
		raidBuffs.DivineSpirit = false
		raidBuffs.GiftOfTheWild = 0
		raidBuffs.PowerWordFortitude = 0
		raidBuffs.Thorns = 0
		raidBuffs.ShadowProtection = false
		raidBuffs.DrumsOfForgottenKings = false
		raidBuffs.DrumsOfTheWild = false
		raidBuffs.ScrollOfProtection = false
		raidBuffs.ScrollOfStamina = false
		raidBuffs.ScrollOfStrength = false
		raidBuffs.ScrollOfAgility = false
		raidBuffs.ScrollOfIntellect = false
		raidBuffs.ScrollOfSpirit = false
		individualBuffs.BlessingOfKings = false
		individualBuffs.BlessingOfSanctuary = false
		individualBuffs.BlessingOfMight = 0
		individualBuffs.BlessingOfWisdom = 0
	}

	// For some reason pets don't benefit from buffs that are ratings, e.g. crit rating or haste rating.
	partyBuffs.BraidedEterniumChain = false

	applyBuffEffects(petAgent, raidBuffs, partyBuffs, individualBuffs)
}

func InspirationAura(unit *Unit, points int32) *Aura {
	multiplier := 1 - []float64{0, .03, .07, .10}[points]

	return unit.GetOrRegisterAura(Aura{
		Label:    "Inspiration",
		ActionID: ActionID{SpellID: 15363},
		Duration: time.Second * 15,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= multiplier
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= multiplier
		},
	})
}

func ApplyInspiration(character *Character, uptime float64) {
	if uptime <= 0 {
		return
	}
	uptime = MinFloat(1, uptime)

	inspirationAura := InspirationAura(&character.Unit, 3)

	ApplyFixedUptimeAura(inspirationAura, uptime, time.Millisecond*2500, 1)
}

func RetributionAura(character *Character, sanctifiedRetribution bool) *Aura {
	actionID := ActionID{SpellID: 54043}

	baseDamage := 112.0
	if sanctifiedRetribution {
		baseDamage *= 1.5
	}

	procSpell := character.RegisterSpell(SpellConfig{
		ActionID:    actionID,
		SpellSchool: SpellSchoolHoly,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagBinary,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
		},
	})

	return character.RegisterAura(Aura{
		Label:    "Retribution Aura",
		ActionID: actionID,
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if result.Landed() && spell.SpellSchool == SpellSchoolPhysical {
				procSpell.Cast(sim, spell.Unit)
			}
		},
	})
}

func ThornsAura(character *Character, points int32) *Aura {
	actionID := ActionID{SpellID: 53307}
	baseDamage := 73 * (1 + 0.25*float64(points))

	procSpell := character.RegisterSpell(SpellConfig{
		ActionID:    actionID,
		SpellSchool: SpellSchoolNature,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagBinary,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
		},
	})

	return character.RegisterAura(Aura{
		Label:    "Thorns",
		ActionID: actionID,
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if result.Landed() && spell.SpellSchool == SpellSchoolPhysical {
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
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if result.Outcome.Matches(OutcomeBlock | OutcomeDodge | OutcomeParry) {
				character.AddMana(sim, 0.02*character.MaxMana(), manaMetrics)
			}
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
		ExtraCastCondition: func(sim *Simulation, target *Unit) bool {
			if !externalTimers[nextExternalIndex].IsReady(sim) {
				return false
			}

			if character.HasActiveAuraWithTag(config.AuraTag) {
				return false
			}

			return true
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

		ShouldActivate: config.ShouldActivate,
	})
}

var BloodlustActionID = ActionID{SpellID: 2825}

const BloodlustAuraTag = "Bloodlust"
const BloodlustDuration = time.Second * 40
const BloodlustCD = time.Minute * 10

func registerBloodlustCD(agent Agent) {
	character := agent.GetCharacter()
	bloodlustAura := BloodlustAura(character, -1)

	spell := character.RegisterSpell(SpellConfig{
		ActionID: bloodlustAura.ActionID,
		Flags:    SpellFlagNoOnCastComplete | SpellFlagNoMetrics | SpellFlagNoLogs,

		Cast: CastConfig{
			CD: Cooldown{
				Timer:    character.NewTimer(),
				Duration: BloodlustCD,
			},
		},

		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			bloodlustAura.Activate(sim)
		},
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell:    spell,
		Priority: CooldownPriorityBloodlust,
		Type:     CooldownTypeDPS,
		ShouldActivate: func(sim *Simulation, character *Character) bool {
			// Haste portion doesn't stack with Power Infusion, so prefer to wait.
			return !character.HasActiveAuraWithTag(PowerInfusionAuraTag)
		},
	})
}

func BloodlustAura(character *Character, actionTag int32) *Aura {
	actionID := BloodlustActionID.WithTag(actionTag)
	aura := character.GetOrRegisterAura(Aura{
		Label:    "Bloodlust-" + actionID.String(),
		Tag:      BloodlustAuraTag,
		ActionID: actionID,
		Duration: BloodlustDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.MultiplyAttackSpeed(sim, 1.3)

			if len(character.Pets) > 0 {
				for _, petAgent := range character.Pets {
					pet := petAgent.GetPet()
					if pet.IsEnabled() && !pet.IsGuardian() {
						BloodlustAura(&pet.Character, actionTag).Activate(sim)
					}
				}
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.MultiplyAttackSpeed(sim, 1.0/1.3)
		},
	})
	multiplyCastSpeedEffect(aura, 1.3)
	return aura
}

var PowerInfusionAuraTag = "PowerInfusion"

const PowerInfusionDuration = time.Second * 15
const PowerInfusionCD = time.Minute * 2

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
			AuraCD:           time.Duration(float64(PowerInfusionCD) * 0.8), // All disc priests take Ascension talent.
			Type:             CooldownTypeDPS,

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
	aura := character.GetOrRegisterAura(Aura{
		Label:    "PowerInfusion-" + actionID.String(),
		Tag:      PowerInfusionAuraTag,
		ActionID: actionID,
		Duration: PowerInfusionDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			if character.HasManaBar() {
				character.PseudoStats.CostMultiplier -= 0.2
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if character.HasManaBar() {
				character.PseudoStats.CostMultiplier += 0.2
			}
		},
	})
	multiplyCastSpeedEffect(aura, 1.2)
	return aura
}

func multiplyCastSpeedEffect(aura *Aura, multiplier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("MultiplyCastSpeed", false, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(multiplier)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(1 / multiplier)
		},
	})
}

var TricksOfTheTradeAuraTag = "TricksOfTheTrade"

const TricksOfTheTradeCD = time.Second * 3600 // CD is 30s from the time buff ends (so 40s with glyph) but that's in order to be able to set the number of TotT you'll have during the fight

func registerTricksOfTheTradeCD(agent Agent, numTricksOfTheTrades int32) {
	if numTricksOfTheTrades == 0 {
		return
	}

	// Assuming rogues have Glyph of TotT by default (which might not be the case).
	TotTAura := TricksOfTheTradeAura(agent.GetCharacter(), -1, true)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 57933, Tag: -1},
			AuraTag:          TricksOfTheTradeAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     TotTAura.Duration,
			AuraCD:           TricksOfTheTradeCD,
			Type:             CooldownTypeDPS,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) { TotTAura.Activate(sim) },
		},
		numTricksOfTheTrades)
}

func TricksOfTheTradeAura(character *Character, actionTag int32, glyphed bool) *Aura {
	actionID := ActionID{SpellID: 57933, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "TricksOfTheTrade-" + actionID.String(),
		Tag:      TricksOfTheTradeAuraTag,
		ActionID: actionID,
		Duration: TernaryDuration(glyphed, time.Second*10, time.Second*6),
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageDealtMultiplier *= 1.15
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageDealtMultiplier /= 1.15
		},
	})
}

var UnholyFrenzyAuraTag = "UnholyFrenzy"

const UnholyFrenzyDuration = time.Second * 30
const UnholyFrenzyCD = time.Minute * 3

func registerUnholyFrenzyCD(agent Agent, numUnholyFrenzy int32) {
	if numUnholyFrenzy == 0 {
		return
	}

	ufAura := UnholyFrenzyAura(agent.GetCharacter(), -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 49016, Tag: -1},
			AuraTag:          UnholyFrenzyAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     UnholyFrenzyDuration,
			AuraCD:           UnholyFrenzyCD,
			Type:             CooldownTypeDPS,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) { ufAura.Activate(sim) },
		},
		numUnholyFrenzy)
}

func UnholyFrenzyAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 49016, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "UnholyFrenzy-" + actionID.String(),
		Tag:      UnholyFrenzyAuraTag,
		ActionID: actionID,
		Duration: UnholyFrenzyDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.2
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.2
		},
	})
}

var DivineGuardianAuraTag = "DivineGuardian"

const DivineGuardianDuration = time.Second * 6
const DivineGuardianCD = time.Minute * 2

func registerDivineGuardianCD(agent Agent, numDivineGuardians int32) {
	if numDivineGuardians == 0 {
		return
	}

	dgAura := DivineGuardianAura(agent.GetCharacter(), -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 53530, Tag: -1},
			AuraTag:          DivineGuardianAuraTag,
			CooldownPriority: CooldownPriorityLow,
			AuraDuration:     DivineGuardianDuration,
			AuraCD:           DivineGuardianCD,
			Type:             CooldownTypeSurvival,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) { dgAura.Activate(sim) },
		},
		numDivineGuardians)
}

func DivineGuardianAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 53530, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "DivineGuardian-" + actionID.String(),
		Tag:      DivineGuardianAuraTag,
		ActionID: actionID,
		Duration: DivineGuardianDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier *= 0.8
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier /= 0.8
		},
	})
}

var PainSuppressionAuraTag = "PainSuppression"

const PainSuppressionDuration = time.Second * 8
const PainSuppressionCD = time.Minute * 3

func registerPainSuppressionCD(agent Agent, numPainSuppressions int32) {
	if numPainSuppressions == 0 {
		return
	}

	psAura := PainSuppressionAura(agent.GetCharacter(), -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 33206, Tag: -1},
			AuraTag:          PainSuppressionAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     PainSuppressionDuration,
			AuraCD:           PainSuppressionCD,
			Type:             CooldownTypeSurvival,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) { psAura.Activate(sim) },
		},
		numPainSuppressions)
}

func PainSuppressionAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 33206, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "PainSuppression-" + actionID.String(),
		Tag:      PainSuppressionAuraTag,
		ActionID: actionID,
		Duration: PainSuppressionDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier *= 0.6
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier /= 0.6
		},
	})
}

func registerRevitalizeHotCD(agent Agent, label string, hotID ActionID, ticks int, tickPeriod time.Duration, uptimeCount int32) {
	if uptimeCount == 0 {
		return
	}

	character := agent.GetCharacter()
	revActionID := ActionID{SpellID: 48545}

	manaMetrics := character.NewManaMetrics(revActionID)
	energyMetrics := character.NewEnergyMetrics(revActionID)
	rageMetrics := character.NewRageMetrics(revActionID)
	rpMetrics := character.NewRunicPowerMetrics(revActionID)

	// Calculate desired downtime based on selected uptimeCount (1 count = 10% uptime, 0%-100%)
	totalDuration := time.Duration(ticks) * tickPeriod
	uptimePercent := float64(uptimeCount) / 100.0

	aura := character.GetOrRegisterAura(Aura{
		Label:    "Revitalize-" + label,
		ActionID: hotID,
		Duration: totalDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			pa := NewPeriodicAction(sim, PeriodicActionOptions{
				Period:   tickPeriod,
				NumTicks: ticks,
				OnAction: func(s *Simulation) {
					if s.RandomFloat("Revitalize Proc") < 0.15 {
						cpb := aura.Unit.GetCurrentPowerBar()
						if cpb == ManaBar {
							aura.Unit.AddMana(s, 0.01*aura.Unit.MaxMana(), manaMetrics)
						} else if cpb == EnergyBar {
							aura.Unit.AddEnergy(s, 8, energyMetrics)
						} else if cpb == RageBar {
							aura.Unit.AddRage(s, 4, rageMetrics)
						} else if cpb == RunicPower {
							aura.Unit.AddRunicPower(s, 16, rpMetrics)
						}
					}
				},
			})
			sim.AddPendingAction(pa)
		},
	})

	ApplyFixedUptimeAura(aura, uptimePercent, totalDuration, 1)
}

const ShatteringThrowCD = time.Minute * 5

func registerShatteringThrowCD(agent Agent, numShatteringThrows int32) {
	if numShatteringThrows == 0 {
		return
	}

	stAura := ShatteringThrowAura(agent.GetCharacter().Env.Encounter.TargetUnits[0])

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 64382, Tag: -1},
			AuraTag:          ShatteringThrowAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     ShatteringThrowDuration,
			AuraCD:           ShatteringThrowCD,
			Type:             CooldownTypeDPS,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) { stAura.Activate(sim) },
		},
		numShatteringThrows)
}

var InnervateAuraTag = "Innervate"

const InnervateDuration = time.Second * 10
const InnervateCD = time.Minute * 3

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
	var innervateAura *Aura

	character := agent.GetCharacter()
	character.Env.RegisterPostFinalizeEffect(func() {
		innervateThreshold = InnervateManaThreshold(character)
		innervateAura = InnervateAura(character, -1)
	})

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 29166, Tag: -1},
			AuraTag:          InnervateAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     InnervateDuration,
			AuraCD:           InnervateCD,
			Type:             CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only cast innervate when very low on mana, to make sure all other mana CDs are prioritized.
				return character.CurrentMana() <= innervateThreshold
			},
			AddAura: func(sim *Simulation, character *Character) {
				innervateAura.Activate(sim)
			},
		},
		numInnervates)
}

func InnervateAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 29166, Tag: actionTag}
	manaMetrics := character.NewManaMetrics(actionID)
	return character.GetOrRegisterAura(Aura{
		Label:    "Innervate-" + actionID.String(),
		Tag:      InnervateAuraTag,
		ActionID: actionID,
		Duration: InnervateDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			const manaPerTick = 3496 * 2.25 / 10 // WotLK druid's base mana
			StartPeriodicAction(sim, PeriodicActionOptions{
				Period:   InnervateDuration / 10,
				NumTicks: 10,
				OnAction: func(sim *Simulation) {
					character.AddMana(sim, manaPerTick, manaMetrics)
				},
			})
		},
	})
}

var ManaTideTotemActionID = ActionID{SpellID: 16190}
var ManaTideTotemAuraTag = "ManaTideTotem"

const ManaTideTotemDuration = time.Second * 12
const ManaTideTotemCD = time.Minute * 5

func registerManaTideTotemCD(agent Agent, numManaTideTotems int32) {
	if numManaTideTotems == 0 {
		return
	}

	initialDelay := time.Duration(0)
	var mttAura *Aura

	character := agent.GetCharacter()
	character.Env.RegisterPostFinalizeEffect(func() {
		// Use first MTT at 60s, or halfway through the fight, whichever comes first.
		initialDelay = MinDuration(character.Env.BaseDuration/2, time.Second*60)
		mttAura = ManaTideTotemAura(character, -1)
	})

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ManaTideTotemActionID.WithTag(-1),
			AuraTag:          ManaTideTotemAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     ManaTideTotemDuration,
			AuraCD:           ManaTideTotemCD,
			Type:             CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// A normal resto shaman would wait to use MTT.
				return sim.CurrentTime >= initialDelay
			},
			AddAura: func(sim *Simulation, character *Character) {
				mttAura.Activate(sim)
			},
		},
		numManaTideTotems)
}

func ManaTideTotemAura(character *Character, actionTag int32) *Aura {
	actionID := ManaTideTotemActionID.WithTag(actionTag)

	metrics := make([]*ResourceMetrics, len(character.Party.Players))
	for i, player := range character.Party.Players {
		char := player.GetCharacter()
		if char.HasManaBar() {
			metrics[i] = char.NewManaMetrics(actionID)
		}
	}

	return character.GetOrRegisterAura(Aura{
		Label:    "ManaTideTotem-" + actionID.String(),
		Tag:      ManaTideTotemAuraTag,
		ActionID: actionID,
		Duration: ManaTideTotemDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			StartPeriodicAction(sim, PeriodicActionOptions{
				Period:   ManaTideTotemDuration / 4,
				NumTicks: 4,
				OnAction: func(sim *Simulation) {
					for i, player := range character.Party.Players {
						if metrics[i] != nil {
							char := player.GetCharacter()
							char.AddMana(sim, 0.06*char.MaxMana(), metrics[i])
						}
					}
				},
			})
		},
	})
}

const ReplenishmentAuraDuration = time.Second * 15

// Creates the actual replenishment aura for a unit.
func replenishmentAura(unit *Unit, _ ActionID) *Aura {
	if unit.ReplenishmentAura != nil {
		return unit.ReplenishmentAura
	}

	replenishmentDep := unit.NewDynamicStatDependency(stats.Mana, stats.MP5, 0.01)

	unit.ReplenishmentAura = unit.RegisterAura(Aura{
		Label:    "Replenishment",
		ActionID: ActionID{SpellID: 57669},
		Duration: ReplenishmentAuraDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.EnableDynamicStatDep(sim, replenishmentDep)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.DisableDynamicStatDep(sim, replenishmentDep)
		},
	})

	return unit.ReplenishmentAura
}

type ReplenishmentSource int

// Returns a new aura whose activation will give the Replenishment buff to 10 party/raid members.
func (raid *Raid) NewReplenishmentSource(actionID ActionID) ReplenishmentSource {
	newReplSource := ReplenishmentSource(len(raid.curReplenishmentUnits))
	raid.curReplenishmentUnits = append(raid.curReplenishmentUnits, []*Unit{})

	if raid.replenishmentUnits != nil {
		return newReplSource
	}

	// Get the list of all eligible units (party/raid members + their pets, but no guardians).
	var manaUsers []*Unit
	for _, party := range raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			if character.HasManaBar() {
				manaUsers = append(manaUsers, &character.Unit)
			}
		}
		for _, petAgent := range party.Pets {
			pet := petAgent.GetPet()
			if pet.HasManaBar() && !pet.IsGuardian() {
				manaUsers = append(manaUsers, &pet.Unit)
			}
		}
	}
	raid.replenishmentUnits = manaUsers

	// Initialize replenishment aura for all applicable units.
	for _, unit := range raid.replenishmentUnits {
		replenishmentAura(unit, actionID)
	}

	return newReplSource
}

func (raid *Raid) resetReplenishment(_ *Simulation) {
	raid.leftoverReplenishmentUnits = raid.replenishmentUnits
	for i := 0; i < len(raid.curReplenishmentUnits); i++ {
		raid.curReplenishmentUnits[i] = nil
	}
}

func (raid *Raid) ProcReplenishment(sim *Simulation, src ReplenishmentSource) {
	// If the raid is fully covered by one or more replenishment sources, we can
	// skip the mana sorting.
	if len(raid.curReplenishmentUnits)*10 >= len(raid.replenishmentUnits) {
		if len(raid.curReplenishmentUnits[src]) == 0 {
			if len(raid.leftoverReplenishmentUnits) > 10 {
				raid.curReplenishmentUnits[src] = raid.leftoverReplenishmentUnits[:10]
				raid.leftoverReplenishmentUnits = raid.leftoverReplenishmentUnits[10:]
			} else {
				raid.curReplenishmentUnits[src] = raid.leftoverReplenishmentUnits
				raid.leftoverReplenishmentUnits = nil
			}
		}
		for _, unit := range raid.curReplenishmentUnits[src] {
			unit.ReplenishmentAura.Activate(sim)
		}
		return
	}

	eligible := append(raid.curReplenishmentUnits[src], raid.leftoverReplenishmentUnits...)
	slices.SortFunc(eligible, func(v1, v2 *Unit) bool {
		return v1.CurrentManaPercent() < v2.CurrentManaPercent()
	})
	raid.curReplenishmentUnits[src] = eligible[:10]
	raid.leftoverReplenishmentUnits = eligible[10:]
	for _, unit := range raid.curReplenishmentUnits[src] {
		unit.ReplenishmentAura.Activate(sim)
	}
	for _, unit := range raid.leftoverReplenishmentUnits {
		unit.ReplenishmentAura.Deactivate(sim)
	}
}

func TotemOfWrathAura(character *Character) *Aura {
	aura := character.GetOrRegisterAura(Aura{
		Label:      "Totem of Wrath",
		ActionID:   ActionID{SpellID: 57722},
		Duration:   NeverExpires,
		BuildPhase: CharacterBuildPhaseBuffs,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
	spellPowerBonusEffect(aura, 280)
	return aura
}

func FlametongueTotemAura(character *Character) *Aura {
	aura := character.GetOrRegisterAura(Aura{
		Label:      "Flametongue Totem",
		ActionID:   ActionID{SpellID: 58656},
		Duration:   NeverExpires,
		BuildPhase: CharacterBuildPhaseBuffs,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
	spellPowerBonusEffect(aura, 144)
	return aura
}

func DemonicPactAura(character *Character) *Aura {
	aura := character.GetOrRegisterAura(Aura{
		Label:      "Demonic Pact",
		ActionID:   ActionID{SpellID: 47240},
		Duration:   time.Second * 45,
		BuildPhase: CharacterBuildPhaseBuffs,
	})
	spellPowerBonusEffect(aura, 0)
	return aura
}

func spellPowerBonusEffect(aura *Aura, spellPowerBonus float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("SpellPowerBonus", false, ExclusiveEffect{
		Priority: spellPowerBonus,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.SpellPower: ee.Priority,
			})
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.SpellPower: -ee.Priority,
			})
		},
	})
}

func FocusMagicAura(caster *Character, target *Character) (*Aura, *Aura) {
	actionID := ActionID{SpellID: 54648}

	var casterAura *Aura
	var onHitCallback OnSpellHit
	casterIndex := -1
	if caster != nil {
		casterIndex = int(caster.Index)
		casterAura = caster.GetOrRegisterAura(Aura{
			Label:    "Focus Magic",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
				aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.SpellCrit: 3 * CritRatingPerCritChance,
				})
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.SpellCrit: -3 * CritRatingPerCritChance,
				})
			},
		})

		onHitCallback = func(_ *Aura, sim *Simulation, _ *Spell, result *SpellResult) {
			if result.DidCrit() {
				casterAura.Activate(sim)
			}
		}
	}

	var aura *Aura
	if target != nil {
		aura = target.GetOrRegisterAura(Aura{
			Label:      "Focus Magic" + strconv.Itoa(casterIndex),
			ActionID:   actionID.WithTag(int32(casterIndex)),
			Duration:   NeverExpires,
			BuildPhase: CharacterBuildPhaseBuffs,
			OnReset: func(aura *Aura, sim *Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: onHitCallback,
		})
		aura.NewExclusiveEffect("FocusMagic", true, ExclusiveEffect{
			OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
				ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.SpellCrit: 3 * CritRatingPerCritChance,
				})
			},
			OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
				ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.SpellCrit: -3 * CritRatingPerCritChance,
				})
			},
		})
	}

	return casterAura, aura
}
