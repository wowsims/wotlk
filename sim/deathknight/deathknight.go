package deathknight

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type Rotation_FuStrike int32

const (
	FuStrike_DeathStrike   Rotation_FuStrike = 0
	FuStrike_ScourgeStrike Rotation_FuStrike = 1
	FuStrike_Obliterate    Rotation_FuStrike = 2
)

var TalentTreeSizes = [3]int{28, 29, 31}

type DeathknightInputs struct {
	// Option Vars
	IsDps bool

	UnholyFrenzyTarget *proto.RaidTarget

	StartingRunicPower  float64
	PrecastGhoulFrenzy  bool
	PrecastHornOfWinter bool
	PetUptime           float64

	// Rotation Vars
	RefreshHornOfWinter bool
	ArmyOfTheDeadType   proto.Deathknight_Rotation_ArmyOfTheDead
	StartingPresence    proto.Deathknight_Rotation_Presence
	UseAMS              bool
	AvgAMSSuccessRate   float64
	AvgAMSHit           float64
	FuStrike            Rotation_FuStrike
}

type DeathknightCoeffs struct {
	glacierRotBonusCoeff      float64
	mercilessCombatBonusCoeff float64
	impurityBonusCoeff        float64
	threatOfThassarianChance  float64

	wanderingPlagueMultiplier     float64
	scourgeStrikeShadowMultiplier float64
}

type Deathknight struct {
	core.Character
	Talents *proto.DeathknightTalents

	bonusCoeffs DeathknightCoeffs

	onRuneSpendT10          core.OnRune
	onRuneSpendBladeBarrier core.OnRune

	Inputs DeathknightInputs

	RotationHelper

	Ghoul     *GhoulPet
	RaiseDead *RuneSpell

	Gargoyle                 *GargoylePet
	SummonGargoyle           *RuneSpell
	GargoyleSummonDelay      time.Duration
	OnGargoyleStartFirstCast func()

	RuneWeapon        *RuneWeaponPet
	DancingRuneWeapon *RuneSpell

	ArmyOfTheDead *RuneSpell
	ArmyGhoul     []*GhoulPet

	Bloodworm []*BloodwormPet

	Presence Presence

	IcyTouch   *RuneSpell
	BloodBoil  *RuneSpell
	Pestilence *RuneSpell

	PlagueStrike      *RuneSpell
	PlagueStrikeMhHit *RuneSpell
	PlagueStrikeOhHit *RuneSpell

	DeathStrike      *RuneSpell
	DeathStrikeMhHit *RuneSpell
	DeathStrikeOhHit *RuneSpell
	DeathStrikeHeals []float64

	Obliterate      *RuneSpell
	ObliterateMhHit *RuneSpell
	ObliterateOhHit *RuneSpell

	BloodStrike      *RuneSpell
	BloodStrikeMhHit *RuneSpell
	BloodStrikeOhHit *RuneSpell

	FrostStrike      *RuneSpell
	FrostStrikeMhHit *RuneSpell
	FrostStrikeOhHit *RuneSpell

	HeartStrike       *RuneSpell
	HeartStrikeOffHit *RuneSpell

	RuneStrike     *RuneSpell
	RuneStrikeAura *core.Aura

	GhoulFrenzy *RuneSpell
	// Dummy aura for timeline metrics
	GhoulFrenzyAura *core.Aura

	LastScourgeStrikeDamage float64
	ScourgeStrike           *RuneSpell

	DeathCoil *RuneSpell

	DeathAndDecay    *RuneSpell
	DeathAndDecayDot *core.Dot

	HowlingBlast *RuneSpell

	OtherRelevantStrAgiActive bool
	HornOfWinter              *RuneSpell
	HornOfWinterAura          *core.Aura

	// "CDs"
	RuneTap     *RuneSpell
	MarkOfBlood *RuneSpell

	BloodTap     *RuneSpell
	BloodTapAura *core.Aura

	AntiMagicShell     *RuneSpell
	AntiMagicShellAura *core.Aura

	EmpowerRuneWeapon *RuneSpell

	UnbreakableArmor     *RuneSpell
	UnbreakableArmorAura *core.Aura

	VampiricBlood     *RuneSpell
	VampiricBloodAura *core.Aura

	BoneShield     *RuneSpell
	BoneShieldAura *core.Aura

	IceboundFortitude     *RuneSpell
	IceboundFortitudeAura *core.Aura

	DeathPact *RuneSpell

	// Used only to proc stuff as its free GCD
	MindFreezeSpell *core.Spell

	// Diseases
	FrostFeverSpell     *RuneSpell
	BloodPlagueSpell    *RuneSpell
	FrostFeverDisease   []*core.Dot
	BloodPlagueDisease  []*core.Dot
	FrostFeverExtended  []int
	BloodPlagueExtended []int

	UnholyBlightSpell *core.Spell
	UnholyBlightDots  []*core.Dot

	// Talent Auras
	KillingMachineAura  *core.Aura
	IcyTalonsAura       *core.Aura
	DesolationAura      *core.Aura
	NecrosisAura        *core.Aura
	BloodCakedBladeAura *core.Aura
	ButcheryAura        *core.Aura
	ButcheryPA          *core.PendingAction
	RimeAura            *core.Aura
	BladeBarrierAura    *core.Aura
	SuddenDoomAura      *core.Aura
	ScentOfBloodAura    *core.Aura
	WillOfTheNecropolis *core.Aura

	// Talent Spells
	LastDiseaseDamage float64
	LastTickTime      time.Duration
	WanderingPlague   *core.Spell

	Deathchill     *RuneSpell
	DeathchillAura *core.Aura

	// Presences
	BloodPresence      *RuneSpell
	BloodPresenceAura  *core.Aura
	FrostPresence      *RuneSpell
	FrostPresenceAura  *core.Aura
	UnholyPresence     *RuneSpell
	UnholyPresenceAura *core.Aura

	// Debuffs
	FrostFeverDebuffAura       []*core.Aura
	EbonPlagueOrCryptFeverAura []*core.Aura

	RoRTSBonus func(*core.Unit) float64 // is either RoR or TS bonus function based on talents
}

func (dk *Deathknight) ModifyDamageModifier(value float64) {
	if value > 0 {
		dk.PseudoStats.DamageDealtMultiplier *= 1 + value
	} else {
		dk.PseudoStats.DamageDealtMultiplier /= 1 - value
	}
	dk.modifyShadowDamageModifier(value)
}

func (dk *Deathknight) modifyShadowDamageModifier(value float64) {
	dk.bonusCoeffs.scourgeStrikeShadowMultiplier += value
	dk.bonusCoeffs.wanderingPlagueMultiplier += value / 10
}

func (dk *Deathknight) GetCharacter() *core.Character {
	return &dk.Character
}

func (dk *Deathknight) AddPartyBuffs(_ *proto.PartyBuffs) {
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
	dk.registerDeathStrikeSpell()
	dk.registerHeartStrikeSpell()
	dk.registerMarkOfBloodSpell()
	dk.registerVampiricBloodSpell()
	dk.registerAntiMagicShellSpell()
	dk.registerRuneStrikeSpell()
	dk.registerMindFreeze()

	dk.registerRaiseDeadCD()
	dk.registerSummonGargoyleCD()
	dk.registerArmyOfTheDeadCD()
	dk.registerDancingRuneWeaponCD()
	dk.registerDeathPactSpell()

	dk.registerUnholyFrenzyCD()
}

func (dk *Deathknight) registerMindFreeze() {
	// If talented to have no cost we use it in rotation
	if dk.Talents.EndlessWinter == 2 {
		// We dont care about the kick part and only want it to trigger on harmful spell procs
		dk.MindFreezeSpell = dk.Character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 47528},
			SpellSchool: core.SpellSchoolMagic,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    dk.NewTimer(),
					Duration: time.Second * 10,
				},
			},

			DamageMultiplier: 1,
			ThreatMultiplier: 0,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				// Just deal 0 damage as the "Harmful Spell" is implemented on spell damage
				spell.CalcAndDealDamage(sim, target, 0, spell.OutcomeAlwaysHit)
			},
		})
	}
}

func (dk *Deathknight) ResetBonusCoeffs() {
	dk.bonusCoeffs = DeathknightCoeffs{
		glacierRotBonusCoeff:      1,
		mercilessCombatBonusCoeff: 1,
		impurityBonusCoeff:        1,
		threatOfThassarianChance:  0,

		wanderingPlagueMultiplier:     1,
		scourgeStrikeShadowMultiplier: 1,
	}
}

func (dk *Deathknight) Prepull(sim *core.Simulation) {
	if dk.Inputs.ArmyOfTheDeadType == proto.Deathknight_Rotation_PreCast {
		dk.PrecastArmyOfTheDead(sim)
	}

	if dk.Inputs.PrecastHornOfWinter {
		dk.HornOfWinter.CD.UsePrePull(sim, 1500*time.Millisecond)
	}
}

func (dk *Deathknight) Reset(sim *core.Simulation) {
	dk.LastTickTime = -1
	dk.LastCast = nil
	dk.NextCast = nil
	dk.DeathStrikeHeals = dk.DeathStrikeHeals[:0]
}

func (dk *Deathknight) IsFuStrike(spell *core.Spell) bool {
	return spell == dk.Obliterate.Spell || spell == dk.ScourgeStrike.Spell || spell == dk.DeathStrike.Spell
}

func (dk *Deathknight) HasMajorGlyph(glyph proto.DeathknightMajorGlyph) bool {
	return dk.HasGlyph(int32(glyph))
}
func (dk *Deathknight) HasMinorGlyph(glyph proto.DeathknightMinorGlyph) bool {
	return dk.HasGlyph(int32(glyph))
}

func NewDeathknight(character core.Character, inputs DeathknightInputs, talents string) *Deathknight {
	dk := &Deathknight{
		Character:  character,
		Talents:    &proto.DeathknightTalents{},
		Inputs:     inputs,
		RoRTSBonus: func(u *core.Unit) float64 { return 1.0 }, // default to no bonus for RoR/TS
	}
	core.FillTalentsProto(dk.Talents.ProtoReflect(), talents, TalentTreeSizes)

	maxRunicPower := 100.0 + 15.0*float64(dk.Talents.RunicPowerMastery)
	currentRunicPower := math.Min(maxRunicPower, dk.Inputs.StartingRunicPower+core.TernaryFloat64(dk.Inputs.PrecastHornOfWinter, 10.0, 0.0))

	runeCD := 10 * time.Second
	if dk.Talents.ImprovedUnholyPresence > 0 {
		runeCD = time.Duration(float64(runeCD) * (1.0 - 0.05*float64(dk.Talents.ImprovedUnholyPresence)))
	}

	dk.EnableRunicPowerBar(
		currentRunicPower,
		maxRunicPower,
		runeCD,
		func(sim *core.Simulation) {
			if dk.onRuneSpendT10 != nil {
				dk.onRuneSpendT10(sim)
			}
			if dk.onRuneSpendBladeBarrier != nil {
				dk.onRuneSpendBladeBarrier(sim)
			}
		},
		func(sim *core.Simulation) {
		},
		func(sim *core.Simulation) {
		},
		func(sim *core.Simulation) {
		},
		func(sim *core.Simulation) {
		},
		func(sim *core.Simulation) {
		},
	)

	dk.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/62.5)
	dk.AddStatDependency(stats.Agility, stats.Dodge, core.DodgeRatingPerDodgeChance/84.74576271)
	dk.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	dk.AddStatDependency(stats.Strength, stats.Parry, 0.25)

	dk.PseudoStats.CanParry = true

	// Base dodge unaffected by Diminishing Returns
	dk.PseudoStats.BaseDodge += 0.03664
	dk.PseudoStats.BaseParry += 0.05

	dk.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	dk.Ghoul = dk.NewGhoulPet(dk.Talents.MasterOfGhouls)
	dk.OnGargoyleStartFirstCast = func() {}
	dk.GargoyleSummonDelay = time.Millisecond * 1000

	dk.ArmyGhoul = make([]*GhoulPet, 8)
	for i := 0; i < 8; i++ {
		dk.ArmyGhoul[i] = dk.NewArmyGhoulPet(i)
	}

	if dk.Talents.Bloodworms > 0 {
		dk.Bloodworm = make([]*BloodwormPet, 4)
		for i := 0; i < 4; i++ {
			dk.Bloodworm[i] = dk.NewBloodwormPet(i)
		}
	}

	dk.RuneWeapon = dk.NewRuneWeapon()

	dk.RotationSequence = &Sequence{}

	return dk
}

func (dk *Deathknight) AllDiseasesAreActive(target *core.Unit) bool {
	return dk.FrostFeverDisease[target.Index].IsActive() && dk.BloodPlagueDisease[target.Index].IsActive()
}

func (dk *Deathknight) DiseasesAreActive(target *core.Unit) bool {
	return dk.FrostFeverDisease[target.Index].IsActive() || dk.BloodPlagueDisease[target.Index].IsActive()
}

func (dk *Deathknight) bonusCritMultiplier(bonusTalentPoints int32) float64 {
	return dk.MeleeCritMultiplier(1, 0.15*float64(bonusTalentPoints))
}

func (dk *Deathknight) KM() bool {
	if dk.KillingMachineAura != nil {
		return dk.KillingMachineAura.IsActive()
	} else {
		return false
	}
}

func (dk *Deathknight) Rime() bool {
	if dk.RimeAura != nil {
		return dk.RimeAura.IsActive()
	} else {
		return false
	}
}

func (dk *Deathknight) AverageDSHeal() float64 {
	count := len(dk.DeathStrikeHeals)
	if count >= 5 {
		sum := dk.DeathStrikeHeals[count-1]
		sum += dk.DeathStrikeHeals[count-2]
		sum += dk.DeathStrikeHeals[count-3]
		sum += dk.DeathStrikeHeals[count-4]
		sum += dk.DeathStrikeHeals[count-5]
		return sum / 5.0
	} else if count > 0 {
		sum := dk.DeathStrikeHeals[count-1]
		for i := 1; i < count; i++ {
			sum += dk.DeathStrikeHeals[count-i-1]
		}
		return sum / float64(count)
	} else {
		return 0
	}
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
	}
}

// Agent is a generic way to access underlying warrior on any of the agents.

func (dk *Deathknight) GetDeathKnight() *Deathknight {
	return dk
}

type DeathKnightAgent interface {
	GetDeathKnight() *Deathknight
}

func PointsInTalents(talents *proto.DeathknightTalents) (int, int, int) {
	blood := 0
	blood += int(talents.Butchery)
	blood += int(talents.Subversion)
	blood += int(talents.BladeBarrier)
	blood += int(talents.BladedArmor)
	blood += int(talents.ScentOfBlood)
	blood += int(talents.TwoHandedWeaponSpecialization)
	blood += int(talents.DarkConviction)
	blood += int(talents.DeathRuneMastery)
	blood += int(talents.ImprovedRuneTap)
	blood += int(talents.SpellDeflection)
	blood += int(talents.Vendetta)
	blood += int(talents.BloodyStrikes)
	blood += int(talents.VeteranOfTheThirdWar)
	blood += int(talents.BloodyVengeance)
	blood += int(talents.AbominationsMight)
	blood += int(talents.Bloodworms)
	blood += int(talents.ImprovedBloodPresence)
	blood += int(talents.ImprovedDeathStrike)
	blood += int(talents.SuddenDoom)
	blood += int(talents.WillOfTheNecropolis)
	blood += int(talents.MightOfMograine)
	blood += int(talents.BloodGorged)
	if talents.RuneTap {
		blood++
	}
	if talents.Hysteria {
		blood++
	}
	if talents.MarkOfBlood {
		blood++
	}
	if talents.VampiricBlood {
		blood++
	}
	if talents.HeartStrike {
		blood++
	}
	if talents.DancingRuneWeapon {
		blood++
	}

	frost := 0

	frost += int(talents.ImprovedIcyTouch)
	frost += int(talents.RunicPowerMastery)
	frost += int(talents.Toughness)
	frost += int(talents.IcyReach)
	frost += int(talents.BlackIce)
	frost += int(talents.NervesOfColdSteel)
	frost += int(talents.IcyTalons)
	frost += int(talents.Annihilation)
	frost += int(talents.KillingMachine)
	frost += int(talents.ChillOfTheGrave)
	frost += int(talents.EndlessWinter)
	frost += int(talents.FrigidDreadplate)
	frost += int(talents.GlacierRot)
	frost += int(talents.MercilessCombat)
	frost += int(talents.Rime)
	frost += int(talents.Chilblains)
	frost += int(talents.ImprovedFrostPresence)
	frost += int(talents.ThreatOfThassarian)
	frost += int(talents.BloodOfTheNorth)
	frost += int(talents.Acclimation)
	frost += int(talents.GuileOfGorefiend)
	frost += int(talents.TundraStalker)
	if talents.HowlingBlast {
		frost++
	}
	if talents.Lichborne {
		frost++
	}
	if talents.Deathchill {
		frost++
	}
	if talents.ImprovedIcyTalons {
		frost++
	}
	if talents.HungeringCold {
		frost++
	}
	if talents.UnbreakableArmor {
		frost++
	}
	if talents.FrostStrike {
		frost++
	}

	unholy := 0

	unholy += int(talents.ViciousStrikes)
	unholy += int(talents.Virulence)
	unholy += int(talents.Anticipation)
	unholy += int(talents.Epidemic)
	unholy += int(talents.Morbidity)
	unholy += int(talents.UnholyCommand)
	unholy += int(talents.RavenousDead)
	unholy += int(talents.Outbreak)
	unholy += int(talents.Necrosis)
	unholy += int(talents.OnAPaleHorse)
	unholy += int(talents.BloodCakedBlade)
	unholy += int(talents.NightOfTheDead)
	unholy += int(talents.Impurity)
	unholy += int(talents.Dirge)
	unholy += int(talents.Desecration)
	unholy += int(talents.MagicSuppression)
	unholy += int(talents.Reaping)
	unholy += int(talents.Desolation)
	unholy += int(talents.ImprovedUnholyPresence)
	unholy += int(talents.CryptFever)
	unholy += int(talents.WanderingPlague)
	unholy += int(talents.EbonPlaguebringer)
	unholy += int(talents.RageOfRivendare)
	if talents.CorpseExplosion {
		unholy++
	}
	if talents.UnholyBlight {
		unholy++
	}
	if talents.MasterOfGhouls {
		unholy++
	}
	if talents.AntiMagicZone {
		unholy++
	}
	if talents.GhoulFrenzy {
		unholy++
	}
	if talents.BoneShield {
		unholy++
	}
	if talents.ScourgeStrike {
		unholy++
	}
	if talents.SummonGargoyle {
		unholy++
	}

	return blood, frost, unholy
}
