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
	FirstDisease         proto.Deathknight_Rotation_FirstDisease
}

type DeathknightCoeffs struct {
	glacierRotBonusCoeff      float64
	mercilessCombatBonusCoeff float64
	tundraStalkerBonusCoeff   float64
	rageOfRivendareBonusCoeff float64
	impurityBonusCoeff        float64

	bloodOfTheNorthChance    float64
	threatOfThassarianChance float64
	reapingChance            float64

	additiveDamageModifier float64
}

type Deathknight struct {
	core.Character
	Talents proto.DeathknightTalents

	bonusCoeffs DeathknightCoeffs

	onRuneSpendT10          core.OnRuneSpend
	onRuneSpendBladeBarrier core.OnRuneSpend

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

	DeathCoil *core.Spell

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

	IceboundFortitude     *core.Spell
	IceboundFortitudeAura *core.Aura

	// Diseases
	FrostFeverSpell    *core.Spell
	BloodPlagueSpell   *core.Spell
	FrostFeverDisease  []*core.Dot
	BloodPlagueDisease []*core.Dot

	UnholyBlightSpell      *core.Spell
	UnholyBlightDot        []*core.Dot
	UnholyBlightTickDamage []float64

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
}

func (dk *Deathknight) ModifyAdditiveDamageModifier(sim *core.Simulation, value float64) {
	dk.PseudoStats.DamageDealtMultiplier /= dk.bonusCoeffs.additiveDamageModifier
	dk.bonusCoeffs.additiveDamageModifier += value
	dk.PseudoStats.DamageDealtMultiplier *= dk.bonusCoeffs.additiveDamageModifier
}

func (dk *Deathknight) GetCharacter() *core.Character {
	return &dk.Character
}

func (dk *Deathknight) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (dk *Deathknight) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if dk.Talents.AbominationsMight > 0 {
		raidBuffs.AbominationsMight = true
	}

	if dk.Talents.ImprovedIcyTalons {
		raidBuffs.IcyTalons = true
	}

	raidBuffs.HornOfWinter = !dk.Inputs.RefreshHornOfWinter

	if raidBuffs.StrengthOfEarthTotem == proto.TristateEffect_TristateEffectImproved ||
		raidBuffs.StrengthOfEarthTotem == proto.TristateEffect_TristateEffectRegular {
		dk.OtherRelevantStrAgiActive = true
	} else {
		dk.OtherRelevantStrAgiActive = false
	}
}

func (dk *Deathknight) ApplyTalents() {
	dk.ResetBonusCoeffs()

	dk.ApplyBloodTalents()
	dk.ApplyFrostTalents()
	dk.ApplyUnholyTalents()
}

func (dk *Deathknight) Initialize() {
	dk.registerPresences()
	dk.registerIcyTouchSpell()
	dk.registerPlagueStrikeSpell()
	dk.registerObliterateSpell()
	dk.registerBloodStrikeSpell()
	dk.registerBloodTapSpell()
	dk.registerHowlingBlastSpell()
	dk.registerScourgeStrikeSpell()
	dk.registerDeathCoilSpell()
	dk.registerFrostStrikeSpell()
	dk.registerDeathAndDecaySpell()
	dk.registerDiseaseDots()
	dk.registerGhoulFrenzySpell()
	dk.registerBoneShieldSpell()
	dk.registerUnbreakableArmorSpell()
	dk.registerBloodBoilSpell()
	dk.registerHornOfWinterSpell()
	dk.registerPestilenceSpell()
	dk.registerEmpowerRuneWeaponSpell()
	dk.registerRuneTapSpell()
	dk.registerIceboundFortitudeSpell()

	dk.registerRaiseDeadCD()
	dk.registerSummonGargoyleCD()
	dk.registerArmyOfTheDeadCD()
}

func (dk *Deathknight) ResetBonusCoeffs() {
	dk.bonusCoeffs = DeathknightCoeffs{
		glacierRotBonusCoeff:      1.0,
		mercilessCombatBonusCoeff: 1.0,
		tundraStalkerBonusCoeff:   1.0,
		impurityBonusCoeff:        1.0,
		rageOfRivendareBonusCoeff: 1.0,

		bloodOfTheNorthChance:    0.0,
		threatOfThassarianChance: 0.0,
		reapingChance:            0.0,

		additiveDamageModifier: dk.bonusCoeffs.additiveDamageModifier,
	}
}

func (dk *Deathknight) Reset(sim *core.Simulation) {
	dk.Presence = UnsetPresence
	if dk.Inputs.UnholyPresenceOpener {
		dk.ChangePresence(sim, UnholyPresence)
	} else {
		dk.ChangePresence(sim, BloodPresence)
	}

	if dk.Inputs.ArmyOfTheDeadType == proto.Deathknight_Rotation_PreCast {
		dk.PrecastArmyOfTheDead(sim)
	}

	dk.ResetRotation(sim)
}

func (dk *Deathknight) IsFuStrike(spell *core.Spell) bool {
	return spell == dk.Obliterate || spell == dk.ScourgeStrike // || spell == dk.DeathStrike
}

func (dk *Deathknight) HasMajorGlyph(glyph proto.DeathknightMajorGlyph) bool {
	return dk.HasGlyph(int32(glyph))
}
func (dk *Deathknight) HasMinorGlyph(glyph proto.DeathknightMinorGlyph) bool {
	return dk.HasGlyph(int32(glyph))
}

func NewDeathknight(character core.Character, options proto.Player, inputs DeathknightInputs) *Deathknight {
	deathKnightOptions := options.GetDeathknight()

	dk := &Deathknight{
		Character: character,
		Talents:   *deathKnightOptions.Talents,

		Inputs: inputs,
	}

	dk.bonusCoeffs.additiveDamageModifier = 1

	maxRunicPower := 100.0 + 15.0*float64(dk.Talents.RunicPowerMastery)
	currentRunicPower := math.Min(maxRunicPower, dk.Inputs.StartingRunicPower+core.TernaryFloat64(dk.Inputs.PrecastHornOfWinter, 10.0, 0.0))

	dk.EnableRunicPowerBar(
		currentRunicPower,
		maxRunicPower,
		func(sim *core.Simulation) {
			if dk.onRuneSpendT10 != nil {
				dk.onRuneSpendT10(sim)
			}
			if dk.onRuneSpendBladeBarrier != nil {
				dk.onRuneSpendBladeBarrier(sim)
			}
		},
		func(sim *core.Simulation) {
			// I change this here because when using the opener sequence
			// you do not want these to trigger a tryUseGCD, so after the opener
			// its fine since you're running off a prio system, and rune generation
			// can change your logic which we want.
			if !dk.onOpener {
				if dk.GCD.IsReady(sim) {
					dk.tryUseGCD(sim)
				}
			}
		},
		func(sim *core.Simulation) {
			if !dk.onOpener {
				if dk.GCD.IsReady(sim) {
					dk.tryUseGCD(sim)
				}
			}
		},
		func(sim *core.Simulation) {
			if !dk.onOpener {
				if dk.GCD.IsReady(sim) {
					dk.tryUseGCD(sim)
				}
			}
		},
		func(sim *core.Simulation) {
			if !dk.onOpener {
				if dk.GCD.IsReady(sim) {
					dk.tryUseGCD(sim)
				}
			}
		},
		func(sim *core.Simulation) {
			if !dk.onOpener {
				if dk.GCD.IsReady(sim) {
					dk.tryUseGCD(sim)
				}
			}
		},
	)

	dk.EnableAutoAttacks(dk, core.AutoAttackOptions{
		MainHand:       dk.WeaponFromMainHand(dk.DefaultMeleeCritMultiplier()),
		OffHand:        dk.WeaponFromOffHand(dk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	dk.AddStatDependency(stats.Agility, stats.MeleeCrit, 1.0+(core.CritRatingPerCritChance/62.5))
	dk.AddStatDependency(stats.Agility, stats.Dodge, 1.0+(core.DodgeRatingPerDodgeChance/84.74576271))
	dk.AddStatDependency(stats.Strength, stats.AttackPower, 1.0+2)

	dk.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	dk.Ghoul = dk.NewGhoulPet(dk.Talents.MasterOfGhouls)
	if dk.Talents.SummonGargoyle {
		dk.Gargoyle = dk.NewGargoyle()
	}

	dk.ArmyGhoul = make([]*GhoulPet, 8)
	for i := 0; i < 8; i++ {
		dk.ArmyGhoul[i] = dk.NewArmyGhoulPet(i)
	}

	return dk
}

func (dk *Deathknight) AllDiseasesAreActive(target *core.Unit) bool {
	return dk.FrostFeverDisease[target.Index].IsActive() && dk.BloodPlagueDisease[target.Index].IsActive()
}

func (dk *Deathknight) DiseasesAreActive(target *core.Unit) bool {
	return dk.FrostFeverDisease[target.Index].IsActive() || dk.BloodPlagueDisease[target.Index].IsActive()
}

func (dk *Deathknight) secondaryCritModifier(applyGuile bool, applyMoM bool) float64 {
	secondaryModifier := 0.0
	if applyGuile {
		secondaryModifier += 0.15 * float64(dk.Talents.GuileOfGorefiend)
	}
	if applyMoM {
		secondaryModifier += 0.15 * float64(dk.Talents.MightOfMograine)
	}
	return secondaryModifier
}

// TODO: DKs have x2 modifier on spell crit as a passive. Is this the best way to do it?
func (dk *Deathknight) spellCritMultiplier() float64 {
	return dk.MeleeCritMultiplier(1.0, 0)
}

func (dk *Deathknight) spellCritMultiplierGoGandMoM() float64 {
	applyGuile := dk.Talents.GuileOfGorefiend > 0
	applyMightOfMograine := dk.Talents.MightOfMograine > 0
	return dk.MeleeCritMultiplier(1.0, dk.secondaryCritModifier(applyGuile, applyMightOfMograine))
}

func (dk *Deathknight) critMultiplier() float64 {
	return dk.MeleeCritMultiplier(1.0, 0)
}

func (dk *Deathknight) critMultiplierGoGandMoM() float64 {
	applyGuile := dk.Talents.GuileOfGorefiend > 0
	applyMightOfMograine := dk.Talents.MightOfMograine > 0
	return dk.MeleeCritMultiplier(1.0, dk.secondaryCritModifier(applyGuile, applyMightOfMograine))
}

func (dk *Deathknight) RuneAmountForSpell(spell *core.Spell) core.RuneAmount {
	blood := 0
	frost := 0
	unholy := 0
	switch spell {
	case dk.DeathAndDecay:
		blood = 1
		frost = 1
		unholy = 1
	case dk.ArmyOfTheDead:
		blood = 1
		frost = 1
		unholy = 1
	case dk.Pestilence:
		blood = 1
	case dk.BloodStrike:
		blood = 1
	case dk.BloodBoil:
		blood = 1
	case dk.UnbreakableArmor:
		frost = 1
	case dk.IcyTouch:
		frost = 1
	case dk.PlagueStrike:
		unholy = 1
	case dk.GhoulFrenzy:
		unholy = 1
	case dk.BoneShield:
		unholy = 1
	case dk.ScourgeStrike:
		frost = 1
		unholy = 1
	case dk.Obliterate:
		frost = 1
		unholy = 1
	case dk.HowlingBlast:
		frost = 1
		unholy = 1
	}

	return core.RuneAmount{blood, frost, unholy, 0}
}

func (dk *Deathknight) CanCast(sim *core.Simulation, spell *core.Spell) bool {
	switch spell {
	case dk.DeathAndDecay:
		return dk.CanDeathAndDecay(sim)
	case dk.ArmyOfTheDead:
		return dk.CanArmyOfTheDead(sim)
	case dk.Pestilence:
		return dk.CanPestilence(sim)
	case dk.BloodStrike:
		return dk.CanBloodStrike(sim)
	case dk.BloodBoil:
		return dk.CanBloodBoil(sim)
	case dk.UnbreakableArmor:
		return dk.CanUnbreakableArmor(sim)
	case dk.IcyTouch:
		return dk.CanIcyTouch(sim)
	case dk.PlagueStrike:
		return dk.CanPlagueStrike(sim)
	case dk.GhoulFrenzy:
		return dk.CanGhoulFrenzy(sim)
	case dk.BoneShield:
		return dk.CanBoneShield(sim)
	case dk.ScourgeStrike:
		return dk.CanScourgeStrike(sim)
	case dk.Obliterate:
		return dk.CanObliterate(sim)
	case dk.HowlingBlast:
		return dk.CanHowlingBlast(sim)
	case dk.FrostStrike:
		return dk.CanFrostStrike(sim)
	case dk.DeathCoil:
		return dk.CanDeathCoil(sim)
	case dk.BloodTap:
		return dk.CanBloodTap(sim)
	case dk.EmpowerRuneWeapon:
		return dk.CanEmpowerRuneWeapon(sim)
	case dk.HornOfWinter:
		return dk.CanHornOfWinter(sim)
	case dk.RaiseDead:
		return dk.CanRaiseDead(sim)
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

func (dk *Deathknight) GetDeathKnight() *Deathknight {
	return dk
}

type DeathKnightAgent interface {
	GetDeathKnight() *Deathknight
}
