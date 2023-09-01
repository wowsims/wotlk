package deathknight

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/common/wotlk"
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
	IsDps  bool
	NewDrw bool

	UnholyFrenzyTarget *proto.UnitReference

	StartingRunicPower  float64
	PrecastGhoulFrenzy  bool
	PrecastHornOfWinter bool
	PetUptime           float64
	DrwPestiApply       bool
	BloodOpener         proto.Deathknight_Rotation_BloodOpener

	// Rotation Vars
	RefreshHornOfWinter bool
	ArmyOfTheDeadType   proto.Deathknight_Rotation_ArmyOfTheDead
	StartingPresence    proto.Deathknight_Rotation_Presence
	UseAMS              bool
	AvgAMSSuccessRate   float64
	AvgAMSHit           float64
	FuStrike            Rotation_FuStrike
	DiseaseDowntime     float64
	VirulenceRefresh    float64
}

type DeathknightCoeffs struct {
	runeTapHealing float64

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
	RaiseDead *core.Spell

	Gargoyle                 *GargoylePet
	SummonGargoyle           *core.Spell
	SummonGargoyleAura       *core.Aura
	GargoyleSummonDelay      time.Duration
	OnGargoyleStartFirstCast func()

	RuneWeapon        *RuneWeaponPet
	DancingRuneWeapon *core.Spell
	drwDmgSnapshot    float64
	drwPhysSnapshot   float64

	ArmyOfTheDead *core.Spell
	ArmyGhoul     []*GhoulPet

	Bloodworm []*BloodwormPet

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
	DeathStrikeHeals []float64

	Obliterate      *core.Spell
	ObliterateMhHit *core.Spell
	ObliterateOhHit *core.Spell

	BloodStrike      *core.Spell
	BloodStrikeMhHit *core.Spell
	BloodStrikeOhHit *core.Spell

	FrostStrike      *core.Spell
	FrostStrikeMhHit *core.Spell
	FrostStrikeOhHit *core.Spell

	HeartStrike       *core.Spell
	HeartStrikeOffHit *core.Spell

	RuneStrikeQueued bool
	RuneStrikeQueue  *core.Spell
	RuneStrike       *core.Spell
	RuneStrikeOh     *core.Spell
	RuneStrikeAura   *core.Aura

	GhoulFrenzy *core.Spell
	// Dummy aura for timeline metrics
	GhoulFrenzyAura *core.Aura

	LastScourgeStrikeDamage float64
	ScourgeStrike           *core.Spell

	DeathCoil *core.Spell

	DeathAndDecay *core.Spell

	HowlingBlast *core.Spell

	HasDraeneiHitAura         bool
	OtherRelevantStrAgiActive bool
	HornOfWinter              *core.Spell
	HornOfWinterAura          *core.Aura

	// "CDs"
	RuneTap     *core.Spell
	MarkOfBlood *core.Spell

	BloodTap     *core.Spell
	BloodTapAura *core.Aura

	AntiMagicShell     *core.Spell
	AntiMagicShellAura *core.Aura

	EmpowerRuneWeapon *core.Spell

	UnbreakableArmor     *core.Spell
	UnbreakableArmorAura *core.Aura

	VampiricBlood     *core.Spell
	VampiricBloodAura *core.Aura

	BoneShield     *core.Spell
	BoneShieldAura *core.Aura

	UnholyFrenzy     *core.Spell
	UnholyFrenzyAura *core.Aura

	IceboundFortitude     *core.Spell
	IceboundFortitudeAura *core.Aura

	DeathPact *core.Spell

	// Used only to proc stuff as its free GCD
	MindFreezeSpell *core.Spell

	// Diseases
	FrostFeverSpell     *core.Spell
	BloodPlagueSpell    *core.Spell
	FrostFeverExtended  []int
	BloodPlagueExtended []int

	UnholyBlightSpell *core.Spell

	// Talent Auras
	KillingMachineAura  *core.Aura
	IcyTalonsAura       *core.Aura
	DesolationAura      *core.Aura
	BloodCakedBladeAura *core.Aura
	ButcheryAura        *core.Aura
	ButcheryPA          *core.PendingAction
	FreezingFogAura     *core.Aura
	BladeBarrierAura    *core.Aura
	ScentOfBloodAura    *core.Aura
	WillOfTheNecropolis *core.Aura

	// Talent Spells
	LastDiseaseDamage float64
	LastTickTime      time.Duration
	WanderingPlague   *core.Spell
	NecrosisCoeff     float64
	Necrosis          *core.Spell
	Deathchill        *core.Spell
	DeathchillAura    *core.Aura

	// Presences
	BloodPresence      *core.Spell
	BloodPresenceAura  *core.Aura
	FrostPresence      *core.Spell
	FrostPresenceAura  *core.Aura
	UnholyPresence     *core.Spell
	UnholyPresenceAura *core.Aura

	// Debuffs
	FrostFeverDebuffAura       core.AuraArray
	EbonPlagueOrCryptFeverAura core.AuraArray

	RoRTSBonus func(*core.Unit) float64 // is either RoR or TS bonus function based on talents

	MakeTSRoRAssumptions bool
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

func (dk *Deathknight) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	dk.HasDraeneiHitAura = partyBuffs.HeroicPresence
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

	dk.RegisterAura(core.Aura{
		Label:    "Last Cast Assigner",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.DefaultCast.GCD > 0 {
				dk.LastCast = spell
			}
		},
	})

	if !dk.IsUsingAPL {
		if dk.Inputs.PrecastHornOfWinter {
			dk.RegisterPrepullAction(-1500*time.Millisecond, func(sim *core.Simulation) {
				dk.HornOfWinter.Cast(sim, nil)
			})
		}

		if dk.Inputs.ArmyOfTheDeadType == proto.Deathknight_Rotation_PreCast {
			dk.RegisterPrepullAction(-10*time.Second, func(sim *core.Simulation) {
				dk.ArmyOfTheDead.Cast(sim, nil)
			})
		}
	}

	// allows us to use these auras in the APL pre-pull actions
	wotlk.CreateBlackMagicProcAura(&dk.Character)
	CreateVirulenceProcAura(&dk.Character)
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
		runeTapHealing: 0,

		glacierRotBonusCoeff:      1,
		mercilessCombatBonusCoeff: 1,
		impurityBonusCoeff:        1,
		threatOfThassarianChance:  0,

		wanderingPlagueMultiplier:     1,
		scourgeStrikeShadowMultiplier: 1,
	}
}

func (dk *Deathknight) Reset(sim *core.Simulation) {
	dk.LastTickTime = -1
	dk.LastCast = nil
	dk.NextCast = nil
	dk.DeathStrikeHeals = dk.DeathStrikeHeals[:0]
	dk.MakeTSRoRAssumptions = sim.Raid.Size() <= 1
}

func (dk *Deathknight) IsFuStrike(spell *core.Spell) bool {
	return spell == dk.Obliterate || spell == dk.ScourgeStrike || spell == dk.DeathStrike
}

func (dk *Deathknight) HasMajorGlyph(glyph proto.DeathknightMajorGlyph) bool {
	return dk.HasGlyph(int32(glyph))
}
func (dk *Deathknight) HasMinorGlyph(glyph proto.DeathknightMinorGlyph) bool {
	return dk.HasGlyph(int32(glyph))
}

func NewDeathknight(character core.Character, inputs DeathknightInputs, talents string, preNerfedGargoyle bool) *Deathknight {
	dk := &Deathknight{
		Character:  character,
		Talents:    &proto.DeathknightTalents{},
		Inputs:     inputs,
		RoRTSBonus: func(u *core.Unit) float64 { return 1.0 }, // default to no bonus for RoR/TS
	}
	core.FillTalentsProto(dk.Talents.ProtoReflect(), talents, TalentTreeSizes)

	maxRunicPower := 100.0 + 15.0*float64(dk.Talents.RunicPowerMastery)
	currentRunicPower := math.Min(maxRunicPower, dk.Inputs.StartingRunicPower+core.TernaryFloat64(dk.Inputs.PrecastHornOfWinter, 10.0, 0.0))

	dk.EnableRunicPowerBar(
		currentRunicPower,
		maxRunicPower,
		10*time.Second,
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

	dk.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[character.Class]*core.CritRatingPerCritChance)
	dk.AddStatDependency(stats.Agility, stats.Dodge, core.DodgeRatingPerDodgeChance/84.74576271)
	dk.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	dk.AddStatDependency(stats.Strength, stats.Parry, 0.25)
	dk.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	dk.PseudoStats.CanParry = true
	dk.PseudoStats.GracefulCastCDFailures = true

	// Base dodge unaffected by Diminishing Returns
	dk.PseudoStats.BaseDodge += 0.03664
	dk.PseudoStats.BaseParry += 0.05

	dk.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	if dk.Talents.SummonGargoyle {
		dk.Gargoyle = dk.NewGargoyle(!preNerfedGargoyle)
	}

	dk.Ghoul = dk.NewGhoulPet(dk.Talents.MasterOfGhouls)
	dk.OnGargoyleStartFirstCast = func() {}
	dk.GargoyleSummonDelay = time.Millisecond * 2500

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

	if dk.Talents.DancingRuneWeapon {
		dk.RuneWeapon = dk.NewRuneWeapon()
	}

	dk.RotationSequence = &Sequence{}
	// done here so enchants that modify stats are applied before stats are calculated
	dk.registerItems()

	return dk
}

func (dk *Deathknight) AllDiseasesAreActive(target *core.Unit) bool {
	return dk.FrostFeverSpell.Dot(target).IsActive() && dk.BloodPlagueSpell.Dot(target).IsActive()
}

func (dk *Deathknight) DiseasesAreActive(target *core.Unit) bool {
	return dk.FrostFeverSpell.Dot(target).IsActive() || dk.BloodPlagueSpell.Dot(target).IsActive()
}

func (dk *Deathknight) DrwDiseasesAreActive(target *core.Unit) bool {
	return dk.Talents.DancingRuneWeapon && dk.RuneWeapon.FrostFeverSpell.Dot(target).IsActive() || dk.RuneWeapon.BloodPlagueSpell.Dot(target).IsActive()
}

func (dk *Deathknight) bonusCritMultiplier(bonusTalentPoints int32) float64 {
	return dk.MeleeCritMultiplier(1, 0.15*float64(bonusTalentPoints))
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
	core.AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassDeathknight)
	core.AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassDeathknight)
	core.AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassDeathknight)
	core.AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassDeathknight)
	core.AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassDeathknight)
	core.AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassDeathknight)
	core.AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassDeathknight)
	core.AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassDeathknight)
	core.AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassDeathknight)
	core.AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassDeathknight)
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
