package deathknight

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/common/wotlk"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const (
	PetSpellHitScale   = 17.0 / 8.0 * core.SpellHitRatingPerHitChance / core.MeleeHitRatingPerHitChance    // 1.7
	PetExpertiseScale  = 3.25 * core.ExpertisePerQuarterPercentReduction / core.MeleeHitRatingPerHitChance // 0.8125
	PetSpellHasteScale = 1.3
)

var TalentTreeSizes = [3]int{28, 29, 31}

type DeathknightInputs struct {
	// Option Vars
	IsDps bool

	UnholyFrenzyTarget *proto.UnitReference

	StartingRunicPower float64
	PetUptime          float64
	DrwPestiApply      bool

	// Rotation Vars
	UseAMS            bool
	AvgAMSSuccessRate float64
	AvgAMSHit         float64
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

	onRuneSpendT10          core.OnRuneChange
	onRuneSpendBladeBarrier core.OnRuneChange

	Inputs DeathknightInputs

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

	HasDraeneiHitAura bool
	HornOfWinter      *core.Spell

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

	raidBuffs.HornOfWinter = true
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

	// allows us to use these auras in the APL pre-pull actions
	wotlk.CreateBlackMagicProcAura(&dk.Character)
	CreateVirulenceProcAura(&dk.Character)

	// for some reason re-using the same label as DMC:G proc causes tests to fail
	dk.NewTemporaryStatsAura("DMC Greatness Pre-Pull Strength Proc", core.ActionID{SpellID: 60229}, stats.Stats{stats.Strength: 300}, time.Second*15)
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

func NewDeathknight(character *core.Character, inputs DeathknightInputs, talents string) *Deathknight {
	dk := &Deathknight{
		Character:  *character,
		Talents:    &proto.DeathknightTalents{},
		Inputs:     inputs,
		RoRTSBonus: func(u *core.Unit) float64 { return 1.0 }, // default to no bonus for RoR/TS
	}
	core.FillTalentsProto(dk.Talents.ProtoReflect(), talents, TalentTreeSizes)

	maxRunicPower := 100.0 + 15.0*float64(dk.Talents.RunicPowerMastery)
	currentRunicPower := math.Min(maxRunicPower, dk.Inputs.StartingRunicPower)

	dk.EnableRunicPowerBar(
		currentRunicPower,
		maxRunicPower,
		10*time.Second,
		func(sim *core.Simulation, changeType core.RuneChangeType) {
			if dk.onRuneSpendT10 != nil {
				dk.onRuneSpendT10(sim, changeType)
			}
			if dk.onRuneSpendBladeBarrier != nil {
				dk.onRuneSpendBladeBarrier(sim, changeType)
			}
		},
		nil,
	)

	dk.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[character.Class]*core.CritRatingPerCritChance)
	dk.AddStatDependency(stats.Agility, stats.Dodge, core.DodgeRatingPerDodgeChance/84.74576271)
	dk.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	dk.AddStatDependency(stats.Strength, stats.Parry, 0.25)
	dk.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	dk.PseudoStats.CanParry = true

	// Base dodge unaffected by Diminishing Returns
	dk.PseudoStats.BaseDodge += 0.03664
	dk.PseudoStats.BaseParry += 0.05

	dk.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	if dk.Talents.SummonGargoyle {
		dk.Gargoyle = dk.NewGargoyle()
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

// Agent is a generic way to access underlying warrior on any of the agents.

func (dk *Deathknight) GetDeathKnight() *Deathknight {
	return dk
}

type DeathKnightAgent interface {
	GetDeathKnight() *Deathknight
}
