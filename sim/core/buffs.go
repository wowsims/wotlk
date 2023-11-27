package core

import (
	"math"
	"time"

	googleProto "google.golang.org/protobuf/proto"

	"github.com/wowsims/classic/sod/sim/core/proto"
	"github.com/wowsims/classic/sod/sim/core/stats"
)

type BuffName int32

const (
	// General Buffs
	ArcaneIntellect BuffName = iota
	BattleShout
	BlessingOfMight
	BlessingOfWisdom
	BloodPact
	DivineSpirit
	GraceOfAir
	ManaSpring
	MarkOfTheWild
	PowerWordFortitude
	StrengthOfEarth
	TrueshotAura

	// Resistance
	AspectOfTheWild
	FrostResistanceTotem
	FrostResistanceAura
	NatureResistanceTotem
	ShadowProtection

	// Scrolls
	ScrollOfAgility
	ScrollOfIntellect
	ScrollOfSpirit
	ScrollOfStrength
	ScrollOfStamina
)

// Stats from buffs pre-tristate buffs
var BuffSpellByLevel = map[BuffName]map[int32]stats.Stats{
	ArcaneIntellect: {
		25: stats.Stats{
			stats.Intellect: 7, // rank 2
		},
		40: stats.Stats{
			stats.Intellect: 15,
		},
		50: stats.Stats{
			stats.Intellect: 22,
		},
		60: stats.Stats{
			stats.Intellect: 31,
		},
	},
	AspectOfTheWild: {
		25: stats.Stats{
			stats.NatureResistance: 0,
		},
		40: stats.Stats{
			stats.Stamina: 0,
		},
		50: stats.Stats{
			stats.Stamina: 45,
		},
		60: stats.Stats{
			stats.Stamina: 60,
		},
	},
	// TODO: Class Melee specific AP?
	BattleShout: {
		25: stats.Stats{
			stats.AttackPower: 60,
		},
		40: stats.Stats{
			stats.AttackPower: 94,
		},
		50: stats.Stats{
			stats.AttackPower: 139,
		},
		60: stats.Stats{
			stats.AttackPower: 193,
		},
	},
	// TODO: Class use melee AP?
	BlessingOfMight: {
		25: stats.Stats{
			stats.AttackPower: 55,
		},
		40: stats.Stats{
			stats.AttackPower: 85,
		},
		50: stats.Stats{
			stats.AttackPower: 115,
		},
		60: stats.Stats{
			stats.AttackPower: 185,
		},
	},
	BlessingOfWisdom: {
		25: stats.Stats{
			stats.MP5: 15,
		},
		40: stats.Stats{
			stats.MP5: 20,
		},
		50: stats.Stats{
			stats.MP5: 25,
		},
		60: stats.Stats{
			stats.MP5: 33,
		},
	},
	BloodPact: {
		25: stats.Stats{
			stats.Stamina: 9,
		},
		40: stats.Stats{
			stats.Stamina: 30,
		},
		50: stats.Stats{
			stats.Stamina: 42,
		},
		60: stats.Stats{
			stats.Stamina: 42,
		},
	},
	GraceOfAir: {
		25: stats.Stats{
			stats.Agility: 0,
		},
		40: stats.Stats{
			stats.Agility: 0,
		},
		50: stats.Stats{
			stats.Agility: 43,
		},
		60: stats.Stats{
			stats.Agility: 77,
		},
	},
	FrostResistanceAura: {
		25: stats.Stats{
			stats.NatureResistance: 0,
		},
		40: stats.Stats{
			stats.Stamina: 30,
		},
		50: stats.Stats{
			stats.Stamina: 45,
		},
		60: stats.Stats{
			stats.Stamina: 60,
		},
	},
	FrostResistanceTotem: {
		25: stats.Stats{
			stats.NatureResistance: 30,
		},
		40: stats.Stats{
			stats.Stamina: 45,
		},
		50: stats.Stats{
			stats.Stamina: 60,
		},
		60: stats.Stats{
			stats.Stamina: 60,
		},
	},
	ManaSpring: {
		25: stats.Stats{
			stats.MP5: 0,
		},
		40: stats.Stats{
			stats.MP5: 15,
		},
		50: stats.Stats{
			stats.MP5: 20,
		},
		60: stats.Stats{
			stats.MP5: 25,
		},
	},
	MarkOfTheWild: {
		25: stats.Stats{
			stats.Armor:            105,
			stats.Stamina:          4,
			stats.Agility:          4,
			stats.Strength:         4,
			stats.Intellect:        4,
			stats.Spirit:           4,
			stats.ArcaneResistance: 0,
			stats.ShadowResistance: 0,
			stats.NatureResistance: 0,
			stats.FireResistance:   0,
			stats.FrostResistance:  0,
		},
		40: stats.Stats{
			stats.Armor:            195,
			stats.Stamina:          8,
			stats.Agility:          8,
			stats.Strength:         8,
			stats.Intellect:        8,
			stats.Spirit:           8,
			stats.ArcaneResistance: 10,
			stats.ShadowResistance: 10,
			stats.NatureResistance: 10,
			stats.FireResistance:   10,
			stats.FrostResistance:  10,
		},
		50: stats.Stats{
			stats.Armor:            240,
			stats.Stamina:          10,
			stats.Agility:          10,
			stats.Strength:         10,
			stats.Intellect:        10,
			stats.Spirit:           10,
			stats.ArcaneResistance: 15,
			stats.ShadowResistance: 15,
			stats.NatureResistance: 15,
			stats.FireResistance:   15,
			stats.FrostResistance:  15,
		},
		60: stats.Stats{
			stats.Armor:            285,
			stats.Stamina:          12,
			stats.Agility:          12,
			stats.Strength:         12,
			stats.Intellect:        12,
			stats.Spirit:           12,
			stats.ArcaneResistance: 20,
			stats.ShadowResistance: 20,
			stats.NatureResistance: 20,
			stats.FireResistance:   20,
			stats.FrostResistance:  20,
		},
	},
	NatureResistanceTotem: {
		25: stats.Stats{
			stats.NatureResistance: 0,
		},
		40: stats.Stats{
			stats.Stamina: 30,
		},
		50: stats.Stats{
			stats.Stamina: 45,
		},
		60: stats.Stats{
			stats.Stamina: 60,
		},
	},
	PowerWordFortitude: {
		25: stats.Stats{
			stats.Stamina: 20,
		},
		40: stats.Stats{
			stats.Stamina: 32,
		},
		50: stats.Stats{
			stats.Stamina: 43,
		},
		60: stats.Stats{
			stats.Stamina: 54,
		},
	},
	ShadowProtection: {
		25: stats.Stats{
			stats.ShadowResistance: 0,
		},
		40: stats.Stats{
			stats.Stamina: 30,
		},
		50: stats.Stats{
			stats.Stamina: 45,
		},
		60: stats.Stats{
			stats.Stamina: 60,
		},
	},
	TrueshotAura: {
		25: stats.Stats{
			stats.AttackPower: 0,
		},
		40: stats.Stats{
			stats.AttackPower: 50,
		},
		50: stats.Stats{
			stats.AttackPower: 75,
		},
		60: stats.Stats{
			stats.AttackPower: 100,
		},
	},
	StrengthOfEarth: {
		25: stats.Stats{
			stats.Strength: 20,
		},
		40: stats.Stats{
			stats.Strength: 36,
		},
		50: stats.Stats{
			stats.Strength: 61,
		},
		60: stats.Stats{
			stats.Strength: 77,
		},
	},
	ScrollOfAgility: {
		25: stats.Stats{
			stats.Agility: 9,
		},
		40: stats.Stats{
			stats.Agility: 13,
		},
		50: stats.Stats{
			stats.Agility: 17,
		},
		60: stats.Stats{
			stats.Agility: 17,
		},
	},
	ScrollOfIntellect: {
		25: stats.Stats{
			stats.Intellect: 8,
		},
		40: stats.Stats{
			stats.Intellect: 12,
		},
		50: stats.Stats{
			stats.Intellect: 16,
		},
		60: stats.Stats{
			stats.Intellect: 16,
		},
	},
	ScrollOfStamina: {
		25: stats.Stats{
			stats.Stamina: 8,
		},
		40: stats.Stats{
			stats.Stamina: 12,
		},
		50: stats.Stats{
			stats.Stamina: 16,
		},
		60: stats.Stats{
			stats.Stamina: 16,
		},
	},
	ScrollOfStrength: {
		25: stats.Stats{
			stats.Strength: 9,
		},
		40: stats.Stats{
			stats.Strength: 13,
		},
		50: stats.Stats{
			stats.Strength: 13,
		},
		60: stats.Stats{
			stats.Strength: 17,
		},
	},
}

// Applies buffs that affect individual players.
// TODO: Classic Maximum buff based on character level
func applyBuffEffects(agent Agent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individualBuffs *proto.IndividualBuffs) {
	character := agent.GetCharacter()
	level := character.Level
	bonusResist := float64(0)

	if raidBuffs.ArcaneBrilliance {
		character.AddStats(BuffSpellByLevel[ArcaneIntellect][level])
	} else if raidBuffs.ScrollOfIntellect {
		character.AddStats(BuffSpellByLevel[ScrollOfIntellect][level])
	}

	if raidBuffs.GiftOfTheWild > 0 {
		updateStats := BuffSpellByLevel[MarkOfTheWild][level]
		if raidBuffs.GiftOfTheWild == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.35)
		}
		character.AddStats(updateStats)
		bonusResist = updateStats[NatureResistanceTotem]
	}

	if raidBuffs.NatureResistanceTotem {
		updateStats := BuffSpellByLevel[NatureResistanceTotem][level]
		updateStats[stats.NatureResistance] = updateStats[stats.NatureResistance] - bonusResist
		character.AddStats(updateStats)
	} else if raidBuffs.AspectOfTheWild {
		updateStats := BuffSpellByLevel[AspectOfTheWild][level]
		updateStats[stats.NatureResistance] = updateStats[stats.NatureResistance] - bonusResist
		character.AddStats(updateStats)
	}

	if raidBuffs.FrostResistanceAura || raidBuffs.FrostResistanceTotem {
		character.AddStat(stats.FrostResistance, 60-bonusResist)
	}

	// TODO: Classic Thorns
	if raidBuffs.Thorns == proto.TristateEffect_TristateEffectImproved {
		ThornsAura(character, 3)
	} else if raidBuffs.Thorns == proto.TristateEffect_TristateEffectRegular {
		ThornsAura(character, 0)
	}

	if raidBuffs.MoonkinAura {
		character.AddStat(stats.SpellCrit, 3*SpellCritRatingPerCritChance)
	}

	if raidBuffs.LeaderOfThePack {
		character.AddStats(stats.Stats{
			stats.MeleeCrit: 3 * CritRatingPerCritChance,
		})
	}

	if raidBuffs.TrueshotAura {
		character.AddStats(BuffSpellByLevel[TrueshotAura][level])
	}

	if raidBuffs.PowerWordFortitude > 0 {
		updateStats := BuffSpellByLevel[PowerWordFortitude][level]
		if raidBuffs.PowerWordFortitude == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.3)
		}
		character.AddStats(updateStats)
	} else if raidBuffs.BloodPact > 0 {
		updateStats := BuffSpellByLevel[BloodPact][level]
		if raidBuffs.BloodPact == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.3)
		}
		character.AddStats(updateStats)
	} else if raidBuffs.ScrollOfStamina {
		character.AddStats(BuffSpellByLevel[ScrollOfStamina][level])
	}

	if raidBuffs.ShadowProtection {
		updateStats := BuffSpellByLevel[ShadowProtection][level]
		updateStats[stats.ShadowResistance] = updateStats[stats.ShadowResistance] - bonusResist
		character.AddStats(updateStats)
	}

	if raidBuffs.DivineSpirit {
		character.AddStats(BuffSpellByLevel[DivineSpirit][level])
	} else if raidBuffs.ScrollOfSpirit {
		character.AddStats(BuffSpellByLevel[ScrollOfSpirit][level])
	}

	kingsAgiIntSpiAmount := 1.0
	kingsStrStamAmount := 1.0
	if individualBuffs.BlessingOfKings {
		kingsAgiIntSpiAmount = 1.1
		kingsStrStamAmount = 1.1
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

	// TODO: Classic
	if individualBuffs.BlessingOfSanctuary {
		character.PseudoStats.DamageTakenMultiplier *= 0.97
		BlessingOfSanctuaryAura(character)
	}

	// TODO: Classic
	if raidBuffs.DevotionAura != proto.TristateEffect_TristateEffectMissing {
		character.AddStats(stats.Stats{
			stats.Armor: GetTristateValueFloat(raidBuffs.DevotionAura, 735, 735*1.25),
		})
	}

	// TODO: Classic
	if raidBuffs.ScrollOfProtection && raidBuffs.DevotionAura == proto.TristateEffect_TristateEffectMissing {
		character.AddStats(stats.Stats{
			stats.Armor: 240,
		})
	}

	// TODO: Classic version
	// if raidBuffs.RetributionAura {
	// 	RetributionAura(character)
	// }

	if raidBuffs.BattleShout > 0 {
		MakePermanent(BattleShoutAura(&character.Unit, GetTristateValueInt32(raidBuffs.BattleShout, 0, 5), 0, level))
	}
	if individualBuffs.BlessingOfMight > 0 {
		MakePermanent(BlessingOfMightAura(&character.Unit, GetTristateValueInt32(individualBuffs.BlessingOfMight, 0, 5), level))
	}

	// TODO: Rune Demo Pact
	// if raidBuffs.DemonicPactOld > 0 || raidBuffs.DemonicPact > 0 || raidBuffs.DemonicPactSp > 0 {
	// 	// Use DemonicPactSp if set.
	// 	power := raidBuffs.DemonicPactSp
	// 	if power == 0 {
	// 		power = raidBuffs.DemonicPact // fallback to old setting.
	// 	}
	// 	if power == 0 {
	// 		power = raidBuffs.DemonicPactOld
	// 	}

	// 	dpAura := DemonicPactAura(character)
	// 	dpAura.ExclusiveEffects[0].Priority = float64(power)
	// 	MakePermanent(dpAura)
	// }

	if raidBuffs.StrengthOfEarthTotem > 0 {
		updateStats := BuffSpellByLevel[StrengthOfEarth][level]
		if raidBuffs.StrengthOfEarthTotem == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.15)
		}
		character.AddStats(updateStats)
	} else if raidBuffs.ScrollOfStrength {
		character.AddStats(BuffSpellByLevel[ScrollOfStrength][level])
	}

	if raidBuffs.GraceOfAirTotem > 0 {
		updateStats := BuffSpellByLevel[GraceOfAir][level]
		if raidBuffs.GraceOfAirTotem == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.15)
		}
		character.AddStats(updateStats)
	} else if raidBuffs.ScrollOfAgility {
		character.AddStats(BuffSpellByLevel[ScrollOfAgility][level])
	}

	if individualBuffs.BlessingOfWisdom > 0 {
		updateStats := BuffSpellByLevel[BlessingOfWisdom][level]
		if individualBuffs.BlessingOfWisdom == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.2)
		}
		character.AddStats(updateStats)
	} else if raidBuffs.ManaSpringTotem > 0 {
		updateStats := BuffSpellByLevel[ManaSpring][level]
		if raidBuffs.ManaSpringTotem == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.25)
		}
		character.AddStats(updateStats)
	}

	// TODO: Classic windfury Totem update
	// if raidBuffs.WindfuryTotem > 0 {
	// 	character.PseudoStats.MeleeSpeedMultiplier *= GetTristateValueFloat(raidBuffs.WindfuryTotem, 1.16, 1.20)
	// }

	// TODO: Classic provide in APL?
	registerPowerInfusionCD(agent, individualBuffs.PowerInfusions)
	registerManaTideTotemCD(agent, partyBuffs.ManaTideTotems)
	registerInnervateCD(agent, individualBuffs.Innervates)

	character.AddStats(stats.Stats{
		stats.SpellCrit: 2 * SpellCritRatingPerCritChance * float64(partyBuffs.AtieshMage),
	})
	character.AddStats(stats.Stats{
		stats.SpellPower: 33 * float64(partyBuffs.AtieshWarlock),
	})
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
	// the owner during combat or don't make sense for a pet.
	individualBuffs.Innervates = 0
	individualBuffs.PowerInfusions = 0

	if !petAgent.GetPet().enabledOnStart {
		raidBuffs.ArcaneBrilliance = false
		raidBuffs.DivineSpirit = false
		raidBuffs.GiftOfTheWild = 0
		raidBuffs.PowerWordFortitude = 0
		raidBuffs.Thorns = 0
		raidBuffs.ShadowProtection = false
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

	applyBuffEffects(petAgent, raidBuffs, partyBuffs, individualBuffs)
}

// TODO: Classic
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
	uptime = min(1, uptime)

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
	actionID := ActionID{SpellID: 20914}
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
	CooldownPriority int32
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

const SatedAuraLabel = "Sated"
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

		ApplyEffects: func(sim *Simulation, target *Unit, _ *Spell) {
			if !target.HasActiveAura(SatedAuraLabel) {
				bloodlustAura.Activate(sim)
			}
		},
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell:    spell,
		Priority: CooldownPriorityBloodlust,
		Type:     CooldownTypeDPS,
		ShouldActivate: func(sim *Simulation, character *Character) bool {
			// Haste portion doesn't stack with Power Infusion, so prefer to wait.
			return !character.HasActiveAuraWithTag(PowerInfusionAuraTag) && !character.HasActiveAura(SatedAuraLabel)
		},
	})
}

func BloodlustAura(character *Character, actionTag int32) *Aura {
	actionID := BloodlustActionID.WithTag(actionTag)

	sated := character.GetOrRegisterAura(Aura{
		Label:    SatedAuraLabel,
		ActionID: ActionID{SpellID: 57724},
		Duration: time.Minute * 10,
	})

	aura := character.GetOrRegisterAura(Aura{
		Label:    "Bloodlust-" + actionID.String(),
		Tag:      BloodlustAuraTag,
		ActionID: actionID,
		Duration: BloodlustDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.MultiplyAttackSpeed(sim, 1.3)
			for _, pet := range character.Pets {
				if pet.IsEnabled() && !pet.IsGuardian() {
					BloodlustAura(&pet.Character, actionTag).Activate(sim)
				}
			}

			if character.HasActiveAura(SatedAuraLabel) {
				aura.Deactivate(sim) // immediately remove it person already has sated.
				return
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.MultiplyAttackSpeed(sim, 1.0/1.3)
			sated.Activate(sim)
		},
	})
	multiplyCastSpeedEffect(aura, 1.3)
	return aura
}

var PowerInfusionActionID = ActionID{SpellID: 10060}
var PowerInfusionAuraTag = "PowerInfusion"

const PowerInfusionDuration = time.Second * 15
const PowerInfusionCD = time.Minute * 3

func registerPowerInfusionCD(agent Agent, numPowerInfusions int32) {
	if numPowerInfusions == 0 {
		return
	}

	piAura := PowerInfusionAura(&agent.GetCharacter().Unit, -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         PowerInfusionActionID.WithTag(-1),
			AuraTag:          PowerInfusionAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     PowerInfusionDuration,
			AuraCD:           time.Duration(float64(PowerInfusionCD) * 0.8), // All disc priests take Ascension talent.
			Type:             CooldownTypeDPS,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Haste portion doesn't stack with Bloodlust, so prefer to wait.
				return !character.HasActiveAura("Bloodlust-" + BloodlustActionID.WithTag(-1).String())
			},
			AddAura: func(sim *Simulation, character *Character) { piAura.Activate(sim) },
		},
		numPowerInfusions)
}

func PowerInfusionAura(character *Unit, actionTag int32) *Aura {
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
	TotTAura := TricksOfTheTradeAura(&agent.GetCharacter().Unit, -1, true)

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
				return !agent.GetCharacter().GetExclusiveEffectCategory("PercentDamageModifier").AnyActive()
			},
			AddAura: func(sim *Simulation, character *Character) { TotTAura.Activate(sim) },
		},
		numTricksOfTheTrades)
}

func TricksOfTheTradeAura(character *Unit, actionTag int32, glyphed bool) *Aura {
	actionID := ActionID{SpellID: 57933, Tag: actionTag}

	aura := character.GetOrRegisterAura(Aura{
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

	RegisterPercentDamageModifierEffect(aura, 1.15)
	return aura
}

var UnholyFrenzyAuraTag = "UnholyFrenzy"

const UnholyFrenzyDuration = time.Second * 30
const UnholyFrenzyCD = time.Minute * 3

func RegisterPercentDamageModifierEffect(aura *Aura, percentDamageModifier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("PercentDamageModifier", true, ExclusiveEffect{
		Priority: percentDamageModifier,
	})
}

var DivineGuardianAuraTag = "DivineGuardian"

const DivineGuardianDuration = time.Second * 6
const DivineGuardianCD = time.Minute * 2

var HandOfSacrificeAuraTag = "HandOfSacrifice"

const HandOfSacrificeDuration = time.Millisecond * 10500 // subtract Divine Shield GCD
const HandOfSacrificeCD = time.Minute * 5                // use Divine Shield CD here

func registerHandOfSacrificeCD(agent Agent, numSacs int32) {
	if numSacs == 0 {
		return
	}

	hosAura := HandOfSacrificeAura(agent.GetCharacter(), -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 6940, Tag: -1},
			AuraTag:          HandOfSacrificeAuraTag,
			CooldownPriority: CooldownPriorityLow,
			AuraDuration:     HandOfSacrificeDuration,
			AuraCD:           HandOfSacrificeCD,
			Type:             CooldownTypeSurvival,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) {
				hosAura.Activate(sim)
			},
		},
		numSacs)
}

func HandOfSacrificeAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 6940, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "HandOfSacrifice-" + actionID.String(),
		Tag:      HandOfSacrificeAuraTag,
		ActionID: actionID,
		Duration: HandOfSacrificeDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier *= 0.7
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier /= 0.7
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

var GuardianSpiritAuraTag = "GuardianSpirit"

const GuardianSpiritDuration = time.Second * 10
const GuardianSpiritCD = time.Minute * 3

func registerGuardianSpiritCD(agent Agent, numGuardianSpirits int32) {
	if numGuardianSpirits == 0 {
		return
	}

	character := agent.GetCharacter()
	gsAura := GuardianSpiritAura(character, -1)
	healthMetrics := character.NewHealthMetrics(ActionID{SpellID: 47788})

	character.AddDynamicDamageTakenModifier(func(sim *Simulation, _ *Spell, result *SpellResult) {
		if (result.Damage >= character.CurrentHealth()) && gsAura.IsActive() {
			result.Damage = character.CurrentHealth()
			character.GainHealth(sim, 0.5*character.MaxHealth(), healthMetrics)
			gsAura.Deactivate(sim)
		}
	})

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 47788, Tag: -1},
			AuraTag:          GuardianSpiritAuraTag,
			CooldownPriority: CooldownPriorityLow,
			AuraDuration:     GuardianSpiritDuration,
			AuraCD:           GuardianSpiritCD,
			Type:             CooldownTypeSurvival,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) {
				gsAura.Activate(sim)
			},
		},
		numGuardianSpirits)
}

func GuardianSpiritAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 47788, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "GuardianSpirit-" + actionID.String(),
		Tag:      GuardianSpiritAuraTag,
		ActionID: actionID,
		Duration: GuardianSpiritDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.HealingTakenMultiplier *= 1.4
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.HealingTakenMultiplier /= 1.4
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
		initialDelay = min(character.Env.BaseDuration/2, time.Second*60)
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

// TODO: Classic Runes
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

func BattleShoutAura(unit *Unit, impBattleShout int32, boomingVoicePts int32, level int32) *Aura {
	aura := unit.GetOrRegisterAura(Aura{
		Label:      "Battle Shout",
		ActionID:   ActionID{SpellID: 25289},
		Duration:   time.Duration(float64(time.Minute*2) * (1 + 0.1*float64(boomingVoicePts))),
		BuildPhase: CharacterBuildPhaseBuffs,
	})
	attackPowerBonusEffect(aura, math.Floor(BuffSpellByLevel[BattleShout][level][stats.AttackPower]*(1+0.05*float64(impBattleShout))))
	return aura
}

func BlessingOfMightAura(unit *Unit, impBomPts int32, level int32) *Aura {
	aura := unit.GetOrRegisterAura(Aura{
		Label:      "Blessing of Might",
		ActionID:   ActionID{SpellID: 25291},
		Duration:   NeverExpires,
		BuildPhase: CharacterBuildPhaseBuffs,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
	attackPowerBonusEffect(aura, math.Floor(BuffSpellByLevel[BlessingOfMight][level][stats.AttackPower]*(1+0.4*float64(impBomPts))))
	return aura
}

func attackPowerBonusEffect(aura *Aura, apBonus float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("AttackPowerBonus", false, ExclusiveEffect{
		Priority: apBonus,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.AttackPower:       ee.Priority,
				stats.RangedAttackPower: ee.Priority,
			})
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.AttackPower:       -ee.Priority,
				stats.RangedAttackPower: -ee.Priority,
			})
		},
	})
}

func CommandingShoutAura(unit *Unit, commandingPresencePts int32, boomingVoicePts int32, minorGlyph bool) *Aura {
	aura := unit.GetOrRegisterAura(Aura{
		Label:      "Commanding Shout",
		ActionID:   ActionID{SpellID: 47440},
		Duration:   time.Duration(float64(time.Minute*2)*(1+0.25*float64(boomingVoicePts))) + TernaryDuration(minorGlyph, 2*time.Minute, 0),
		BuildPhase: CharacterBuildPhaseBuffs,
	})
	healthBonusEffect(aura, 2255*(1+0.05*float64(commandingPresencePts)))
	return aura
}

func BloodPactAura(unit *Unit, impImpPts int32) *Aura {
	aura := unit.GetOrRegisterAura(Aura{
		Label:      "Blood Pact",
		ActionID:   ActionID{SpellID: 11767},
		Duration:   NeverExpires,
		BuildPhase: CharacterBuildPhaseBuffs,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.Stamina: 42 * (0.1 * float64(impImpPts)),
			})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.Stamina: -42 * (0.1 * float64(impImpPts)),
			})
		},
	})
	return aura
}

func healthBonusEffect(aura *Aura, healthBonus float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("HealthBonus", false, ExclusiveEffect{
		Priority: healthBonus,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.Health: ee.Priority,
			})
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.Health: -ee.Priority,
			})
		},
	})
}
