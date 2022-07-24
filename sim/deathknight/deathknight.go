package deathknight

import (
	"math"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type DeathknightInputs struct {
	// Option Vars
	StartingRunicPower  float64
	PrecastGhoulFrenzy  bool
	PrecastHornOfWinter bool
	PetUptime           float64

	// Rotation Vars
	RefreshHornOfWinter  bool
	UnholyPresenceOpener bool
	ArmyOfTheDeadType    proto.Deathknight_Rotation_ArmyOfTheDead
}

type Deathknight struct {
	core.Character
	Talents proto.DeathknightTalents

	Inputs DeathknightInputs

	LastCastOutcome core.HitOutcome
	RotationHelper

	Ghoul     *GhoulPet
	RaiseDead *core.Spell

	Gargoyle       *GargoylePet
	SummonGargoyle *core.Spell

	ArmyOfTheDead *core.Spell
	ArmyGhoul     []*GhoulPet

	Presence Presence

	IcyTouch   *core.Spell
	BloodBoil  *core.Spell
	Pestilence *core.Spell

	PlagueStrike      *core.Spell
	PlagueStrikeMhHit *core.Spell
	PlagueStrikeOhHit *core.Spell

	DeathStrike      *core.Spell
	DeathStrikeMhHit *core.Spell
	DeathStrikeOhHit *core.Spell

	Obliterate      *core.Spell
	ObliterateMhHit *core.Spell
	ObliterateOhHit *core.Spell

	BloodStrike      *core.Spell
	BloodStrikeMhHit *core.Spell
	BloodStrikeOhHit *core.Spell

	FrostStrike      *core.Spell
	FrostStrikeMhHit *core.Spell
	FrostStrikeOhHit *core.Spell

	GhoulFrenzy *core.Spell
	// Dummy aura for timeline metrics
	GhoulFrenzyAura *core.Aura

	LastScourgeStrikeDamage float64
	ScourgeStrike           *core.Spell

	LastDeathCoilDamage float64
	DeathCoil           *core.Spell

	DeathAndDecay    *core.Spell
	DeathAndDecayDot *core.Dot

	HowlingBlast *core.Spell

	OtherRelevantStrAgiActive bool
	HornOfWinter              *core.Spell
	HornOfWinterAura          *core.Aura

	// "CDs"
	RuneTap *core.Spell

	BloodTap     *core.Spell
	BloodTapAura *core.Aura

	EmpowerRuneWeapon *core.Spell

	UnbreakableArmor     *core.Spell
	UnbreakableArmorAura *core.Aura

	BoneShield     *core.Spell
	BoneShieldAura *core.Aura

	// Diseases
	FrostFeverSpell    *core.Spell
	BloodPlagueSpell   *core.Spell
	FrostFeverDisease  []*core.Dot
	BloodPlagueDisease []*core.Dot

	UnholyBlightSpell *core.Spell
	UnholyBlightDot   []*core.Dot

	// Talent Auras
	KillingMachineAura  *core.Aura
	IcyTalonsAura       *core.Aura
	DesolationAura      *core.Aura
	NecrosisAura        *core.Aura
	BloodCakedBladeAura *core.Aura
	ButcheryAura        *core.Aura
	RimeAura            *core.Aura
	BladeBarrierAura    *core.Aura

	// Talent Spells
	LastDiseaseDamage float64
	WanderingPlague   *core.Spell

	// Presences
	BloodPresence      *core.Spell
	BloodPresenceAura  *core.Aura
	FrostPresence      *core.Spell
	FrostPresenceAura  *core.Aura
	UnholyPresence     *core.Spell
	UnholyPresenceAura *core.Aura

	// Debuffs
	FrostFeverDebuffAura []*core.Aura
	CryptFeverAura       []*core.Aura
	EbonPlagueAura       []*core.Aura

	// Dynamic trackers
	additiveDamageModifier float64
}

func (deathKnight *Deathknight) ModifyAdditiveDamageModifier(sim *core.Simulation, value float64) {
	deathKnight.PseudoStats.DamageDealtMultiplier /= deathKnight.additiveDamageModifier
	deathKnight.additiveDamageModifier += value
	deathKnight.PseudoStats.DamageDealtMultiplier *= deathKnight.additiveDamageModifier
}

func (deathKnight *Deathknight) GetCharacter() *core.Character {
	return &deathKnight.Character
}

func (deathKnight *Deathknight) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (deathKnight *Deathknight) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if deathKnight.Talents.AbominationsMight > 0 {
		raidBuffs.AbominationsMight = true
	}

	if deathKnight.Talents.ImprovedIcyTalons {
		raidBuffs.IcyTalons = true
	}

	raidBuffs.HornOfWinter = !deathKnight.Inputs.RefreshHornOfWinter

	if raidBuffs.StrengthOfEarthTotem == proto.TristateEffect_TristateEffectImproved ||
		raidBuffs.StrengthOfEarthTotem == proto.TristateEffect_TristateEffectRegular {
		deathKnight.OtherRelevantStrAgiActive = true
	} else {
		deathKnight.OtherRelevantStrAgiActive = false
	}
}

func (deathKnight *Deathknight) ApplyTalents() {
	deathKnight.ApplyBloodTalents()
	deathKnight.ApplyFrostTalents()
	deathKnight.ApplyUnholyTalents()
}

func (deathKnight *Deathknight) Initialize() {
	deathKnight.registerPresences()
	deathKnight.registerIcyTouchSpell()
	deathKnight.registerPlagueStrikeSpell()
	deathKnight.registerObliterateSpell()
	deathKnight.registerBloodStrikeSpell()
	deathKnight.registerBloodTapSpell()
	deathKnight.registerHowlingBlastSpell()
	deathKnight.registerScourgeStrikeSpell()
	deathKnight.registerDeathCoilSpell()
	deathKnight.registerFrostStrikeSpell()
	deathKnight.registerDeathAndDecaySpell()
	deathKnight.registerDiseaseDots()
	deathKnight.registerGhoulFrenzySpell()
	deathKnight.registerBoneShieldSpell()
	deathKnight.registerUnbreakableArmorSpell()
	deathKnight.registerBloodBoilSpell()
	deathKnight.registerHornOfWinterSpell()
	deathKnight.registerPestilenceSpell()
	deathKnight.registerEmpowerRuneWeaponSpell()
	deathKnight.registerRuneTapSpell()

	deathKnight.registerRaiseDeadCD()
	deathKnight.registerSummonGargoyleCD()
	deathKnight.registerArmyOfTheDeadCD()

	deathKnight.SetupRotation()
}

func (deathKnight *Deathknight) Reset(sim *core.Simulation) {
	deathKnight.Presence = UnsetPresence
	if deathKnight.Inputs.UnholyPresenceOpener {
		deathKnight.ChangePresence(sim, UnholyPresence)
	} else {
		deathKnight.ChangePresence(sim, BloodPresence)
	}

	if deathKnight.Inputs.ArmyOfTheDeadType == proto.Deathknight_Rotation_PreCast {
		deathKnight.PrecastArmyOfTheDead(sim)
	}

	deathKnight.ResetRotation(sim)
}

func (deathKnight *Deathknight) IsFuStrike(spell *core.Spell) bool {
	return spell == deathKnight.Obliterate || spell == deathKnight.ScourgeStrike // || spell == deathKnight.DeathStrike
}

func (deathKnight *Deathknight) HasMajorGlyph(glyph proto.DeathknightMajorGlyph) bool {
	return deathKnight.HasGlyph(int32(glyph))
}
func (deathKnight *Deathknight) HasMinorGlyph(glyph proto.DeathknightMinorGlyph) bool {
	return deathKnight.HasGlyph(int32(glyph))
}

func NewDeathknight(character core.Character, options proto.Player, inputs DeathknightInputs) *Deathknight {
	deathKnightOptions := options.GetDeathknight()

	deathKnight := &Deathknight{
		Character: character,
		Talents:   *deathKnightOptions.Talents,

		Inputs: inputs,

		additiveDamageModifier: 1,
	}

	maxRunicPower := 100.0 + 15.0*float64(deathKnight.Talents.RunicPowerMastery)
	currentRunicPower := math.Min(maxRunicPower, deathKnight.Inputs.StartingRunicPower+core.TernaryFloat64(deathKnight.Inputs.PrecastHornOfWinter, 10.0, 0.0))

	deathKnight.EnableRunicPowerBar(
		deathKnight.Talents.BladeBarrier > 0,
		currentRunicPower,
		maxRunicPower,
		func(sim *core.Simulation) {
			// I change this here because when using the opener sequence
			// you do not want these to trigger a tryUseGCD, so after the opener
			// its fine since you're running off a prio system, and rune generation
			// can change your logic which we want.
			if !deathKnight.onOpener {
				if deathKnight.GCD.IsReady(sim) {
					deathKnight.tryUseGCD(sim)
				}
			}
		},
		func(sim *core.Simulation) {
			if !deathKnight.onOpener {
				if deathKnight.GCD.IsReady(sim) {
					deathKnight.tryUseGCD(sim)
				}
			}
		},
		func(sim *core.Simulation) {
			if !deathKnight.onOpener {
				if deathKnight.GCD.IsReady(sim) {
					deathKnight.tryUseGCD(sim)
				}
			}
		},
		func(sim *core.Simulation) {
			if !deathKnight.onOpener {
				if deathKnight.GCD.IsReady(sim) {
					deathKnight.tryUseGCD(sim)
				}
			}
		},
		func(sim *core.Simulation) {
			if !deathKnight.onOpener {
				if deathKnight.GCD.IsReady(sim) {
					deathKnight.tryUseGCD(sim)
				}
			}
		},
	)

	deathKnight.EnableAutoAttacks(deathKnight, core.AutoAttackOptions{
		MainHand:       deathKnight.WeaponFromMainHand(deathKnight.DefaultMeleeCritMultiplier()),
		OffHand:        deathKnight.WeaponFromOffHand(deathKnight.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	deathKnight.AddStatDependency(stats.Agility, stats.MeleeCrit, 1.0+(core.CritRatingPerCritChance/62.5))
	deathKnight.AddStatDependency(stats.Agility, stats.Dodge, 1.0+(core.DodgeRatingPerDodgeChance/84.74576271))
	deathKnight.AddStatDependency(stats.Strength, stats.AttackPower, 1.0+2)

	deathKnight.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	deathKnight.Ghoul = deathKnight.NewGhoulPet(deathKnight.Talents.MasterOfGhouls)
	if deathKnight.Talents.SummonGargoyle {
		deathKnight.Gargoyle = deathKnight.NewGargoyle()
	}

	deathKnight.ArmyGhoul = make([]*GhoulPet, 8)
	for i := 0; i < 8; i++ {
		deathKnight.ArmyGhoul[i] = deathKnight.NewArmyGhoulPet(i)
	}

	return deathKnight
}

func (deathKnight *Deathknight) AllDiseasesAreActive(target *core.Unit) bool {
	return deathKnight.FrostFeverDisease[target.Index].IsActive() && deathKnight.BloodPlagueDisease[target.Index].IsActive()
}

func (deathKnight *Deathknight) DiseasesAreActive(target *core.Unit) bool {
	return deathKnight.FrostFeverDisease[target.Index].IsActive() || deathKnight.BloodPlagueDisease[target.Index].IsActive()
}

func (deathKnight *Deathknight) secondaryCritModifier(applyGuile bool, applyMoM bool) float64 {
	secondaryModifier := 0.0
	if applyGuile {
		secondaryModifier += 0.15 * float64(deathKnight.Talents.GuileOfGorefiend)
	}
	if applyMoM {
		secondaryModifier += 0.15 * float64(deathKnight.Talents.MightOfMograine)
	}
	return secondaryModifier
}

// TODO: DKs have x2 modifier on spell crit as a passive. Is this the best way to do it?
func (deathKnight *Deathknight) spellCritMultiplier() float64 {
	return deathKnight.MeleeCritMultiplier(1.0, 0)
}

func (deathKnight *Deathknight) spellCritMultiplierGoGandMoM() float64 {
	applyGuile := deathKnight.Talents.GuileOfGorefiend > 0
	applyMightOfMograine := deathKnight.Talents.MightOfMograine > 0
	return deathKnight.MeleeCritMultiplier(1.0, deathKnight.secondaryCritModifier(applyGuile, applyMightOfMograine))
}

func (deathKnight *Deathknight) critMultiplier() float64 {
	return deathKnight.MeleeCritMultiplier(1.0, 0)
}

func (deathKnight *Deathknight) critMultiplierGoGandMoM() float64 {
	applyGuile := deathKnight.Talents.GuileOfGorefiend > 0
	applyMightOfMograine := deathKnight.Talents.MightOfMograine > 0
	return deathKnight.MeleeCritMultiplier(1.0, deathKnight.secondaryCritModifier(applyGuile, applyMightOfMograine))
}

func (deathKnight *Deathknight) RuneAmountForSpell(spell *core.Spell) core.RuneAmount {
	blood := 0
	frost := 0
	unholy := 0
	switch spell {
	case deathKnight.DeathAndDecay:
		blood = 1
		frost = 1
		unholy = 1
	case deathKnight.ArmyOfTheDead:
		blood = 1
		frost = 1
		unholy = 1
	case deathKnight.Pestilence:
		blood = 1
	case deathKnight.BloodStrike:
		blood = 1
	case deathKnight.BloodBoil:
		blood = 1
	case deathKnight.UnbreakableArmor:
		frost = 1
	case deathKnight.IcyTouch:
		frost = 1
	case deathKnight.PlagueStrike:
		unholy = 1
	case deathKnight.GhoulFrenzy:
		unholy = 1
	case deathKnight.BoneShield:
		unholy = 1
	case deathKnight.ScourgeStrike:
		frost = 1
		unholy = 1
	case deathKnight.Obliterate:
		frost = 1
		unholy = 1
	case deathKnight.HowlingBlast:
		frost = 1
		unholy = 1
	}

	return core.RuneAmount{blood, frost, unholy, 0}
}

func (deathKnight *Deathknight) CanCast(sim *core.Simulation, spell *core.Spell) bool {
	switch spell {
	case deathKnight.DeathAndDecay:
		return deathKnight.CanDeathAndDecay(sim)
	case deathKnight.ArmyOfTheDead:
		return deathKnight.CanArmyOfTheDead(sim)
	case deathKnight.Pestilence:
		return deathKnight.CanPestilence(sim)
	case deathKnight.BloodStrike:
		return deathKnight.CanBloodStrike(sim)
	case deathKnight.BloodBoil:
		return deathKnight.CanBloodBoil(sim)
	case deathKnight.UnbreakableArmor:
		return deathKnight.CanUnbreakableArmor(sim)
	case deathKnight.IcyTouch:
		return deathKnight.CanIcyTouch(sim)
	case deathKnight.PlagueStrike:
		return deathKnight.CanPlagueStrike(sim)
	case deathKnight.GhoulFrenzy:
		return deathKnight.CanGhoulFrenzy(sim)
	case deathKnight.BoneShield:
		return deathKnight.CanBoneShield(sim)
	case deathKnight.ScourgeStrike:
		return deathKnight.CanScourgeStrike(sim)
	case deathKnight.Obliterate:
		return deathKnight.CanObliterate(sim)
	case deathKnight.HowlingBlast:
		return deathKnight.CanHowlingBlast(sim)
	case deathKnight.FrostStrike:
		return deathKnight.CanFrostStrike(sim)
	case deathKnight.DeathCoil:
		return deathKnight.CanDeathCoil(sim)
	case deathKnight.BloodTap:
		return deathKnight.CanBloodTap(sim)
	case deathKnight.EmpowerRuneWeapon:
		return deathKnight.CanEmpowerRuneWeapon(sim)
	case deathKnight.HornOfWinter:
		return deathKnight.CanHornOfWinter(sim)
	case deathKnight.RaiseDead:
		return deathKnight.CanRaiseDead(sim)
	default:
		panic("Not in cost list.")
	}

	return false
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassDeathknight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    176,
		stats.Agility:     109,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      61,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassDeathknight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     108,
		stats.Stamina:     161,
		stats.Intellect:   34,
		stats.Spirit:      58,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassDeathknight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    170,
		stats.Agility:     114,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassDeathknight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    175,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassDeathknight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    171,
		stats.Agility:     116,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      59,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassDeathknight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    178,
		stats.Agility:     109,
		stats.Stamina:     161,
		stats.Intellect:   32,
		stats.Spirit:      61,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassDeathknight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     108,
		stats.Stamina:     161,
		stats.Intellect:   31,
		stats.Spirit:      61,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassDeathknight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    176,
		stats.Agility:     114,
		stats.Stamina:     160,
		stats.Intellect:   31,
		stats.Spirit:      60,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassDeathknight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    174,
		stats.Agility:     110,
		stats.Stamina:     160,
		stats.Intellect:   33,
		stats.Spirit:      64,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassDeathknight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    172,
		stats.Agility:     114,
		stats.Stamina:     160,
		stats.Intellect:   38,
		stats.Spirit:      57,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
}

// Agent is a generic way to access underlying warrior on any of the agents.

func (deathKnight *Deathknight) GetDeathKnight() *Deathknight {
	return deathKnight
}

type DeathKnightAgent interface {
	GetDeathKnight() *Deathknight
}
