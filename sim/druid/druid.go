package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const (
	SpellFlagNaturesGrace = core.SpellFlagAgentReserved1
)

type Druid struct {
	core.Character
	SelfBuffs
	Talents *proto.DruidTalents

	StartingForm DruidForm

	RebirthUsed       bool
	MaulRageThreshold float64
	RebirthTiming     float64
	BleedsActive      int
	AssumeBleedActive bool

	Berserk          *core.Spell
	DemoralizingRoar *core.Spell
	Enrage           *core.Spell
	FaerieFire       *core.Spell
	FerociousBite    *core.Spell
	ForceOfNature    *core.Spell
	Hurricane        *core.Spell
	InsectSwarm      *core.Spell
	Lacerate         *core.Spell
	MangleBear       *core.Spell
	MangleCat        *core.Spell
	Maul             *core.Spell
	Moonfire         *core.Spell
	Rebirth          *core.Spell
	Rake             *core.Spell
	Rip              *core.Spell
	SavageRoar       *core.Spell
	Shred            *core.Spell
	Starfire         *core.Spell
	Starfall         *core.Spell
	StarfallSplash   *core.Spell
	SwipeBear        *core.Spell
	SwipeCat         *core.Spell
	TigersFury       *core.Spell
	Typhoon          *core.Spell
	Wrath            *core.Spell

	CatForm  *core.Spell
	BearForm *core.Spell

	InsectSwarmDot    *core.Dot
	LacerateDot       *core.Dot
	LasherweaveDot    *core.Dot
	MoonfireDot       *core.Dot
	RakeDot           *core.Dot
	RipDot            *core.Dot
	StarfallDot       *core.Dot
	StarfallDotSplash *core.Dot

	BearFormAura         *core.Aura
	BerserkAura          *core.Aura
	CatFormAura          *core.Aura
	ClearcastingAura     *core.Aura
	SwiftStarfireAura    *core.Aura
	DemoralizingRoarAura *core.Aura
	EnrageAura           *core.Aura
	FaerieFireAura       *core.Aura
	MangleAura           *core.Aura
	MaulQueueAura        *core.Aura
	NaturesGraceProcAura *core.Aura
	NaturesSwiftnessAura *core.Aura
	TigersFuryAura       *core.Aura
	SavageRoarAura       *core.Aura
	SolarEclipseProcAura *core.Aura
	LunarEclipseProcAura *core.Aura

	PrimalPrecisionRecoveryMetrics *core.ResourceMetrics

	LunarICD core.Cooldown
	SolarICD core.Cooldown
	Treant1  *TreantPet
	Treant2  *TreantPet
	Treant3  *TreantPet

	form          DruidForm
	disabledMCDs  []*core.MajorCooldown
	setBonuses    druidTierSets
	talentBonuses talentBonuses
}

type talentBonuses struct {
	galeWinds       float64
	genesis         float64
	moonfury        float64
	moonglow        float64
	naturesMajesty  float64
	vengeance       float64
	naturesSplendor int32
	starlightWrath  time.Duration
}

type druidTierSets struct {
	balance_t6_2  bool
	balance_t7_2  bool
	balance_t7_4  bool
	balance_t8_2  bool
	balance_t8_4  bool
	balance_t9_2  bool
	balance_t9_4  bool
	balance_t10_2 bool
	balance_t10_4 bool
	balance_pvp_2 bool
	balance_pvp_4 bool

	feral_t8_2 bool
	feral_t8_4 bool
}

type SelfBuffs struct {
	InnervateTarget *proto.RaidTarget
}

// Registering non-unique Talent effects
func (druid *Druid) RegisterTalentsBonuses() {
	druid.talentBonuses = talentBonuses{
		galeWinds:       0.15 * float64(druid.Talents.GaleWinds),
		genesis:         0.01 * float64(druid.Talents.Genesis),                   // additive damage bonus
		moonfury:        []float64{0.0, 0.03, 0.06, 0.1}[druid.Talents.Moonfury], // additive damage bonus
		moonglow:        1 - 0.03*float64(druid.Talents.Moonglow),                // cost reduction
		naturesMajesty:  2 * float64(druid.Talents.NaturesMajesty) * core.CritRatingPerCritChance,
		vengeance:       0.2 * float64(druid.Talents.Vengeance),
		naturesSplendor: core.TernaryInt32(druid.Talents.NaturesSplendor, 1, 0),
		starlightWrath:  time.Millisecond * 100 * time.Duration(druid.Talents.StarlightWrath),
	}
}

func (druid *Druid) ResetTalentsBonuses() {
	druid.talentBonuses = talentBonuses{
		moonfury:        0,
		genesis:         0,
		moonglow:        0,
		naturesMajesty:  0,
		vengeance:       0,
		naturesSplendor: 0,
		starlightWrath:  0,
	}
}

func (druid *Druid) GetCharacter() *core.Character {
	return &druid.Character
}

func (druid *Druid) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.GiftOfTheWild = core.MaxTristate(raidBuffs.GiftOfTheWild, proto.TristateEffect_TristateEffectRegular)
	if druid.Talents.ImprovedMarkOfTheWild == 2 { // probably could work on actually calculating the fraction effect later if we care.
		raidBuffs.GiftOfTheWild = proto.TristateEffect_TristateEffectImproved
	}

	raidBuffs.Thorns = core.MaxTristate(raidBuffs.Thorns, proto.TristateEffect_TristateEffectRegular)
	if druid.Talents.Brambles == 3 {
		raidBuffs.Thorns = proto.TristateEffect_TristateEffectImproved
	}

	if druid.InForm(Moonkin) && druid.Talents.MoonkinForm {
		raidBuffs.MoonkinAura = core.MaxTristate(raidBuffs.MoonkinAura, proto.TristateEffect_TristateEffectRegular)
		if druid.Talents.ImprovedMoonkinForm > 0 {
			// For now, we assume Improved Moonkin Form is maxed-out
			raidBuffs.MoonkinAura = proto.TristateEffect_TristateEffectImproved
		}
	}
	if druid.InForm(Cat|Bear) && druid.Talents.LeaderOfThePack {
		raidBuffs.LeaderOfThePack = core.MaxTristate(raidBuffs.LeaderOfThePack, proto.TristateEffect_TristateEffectRegular)
		if druid.Talents.ImprovedLeaderOfThePack > 0 {
			raidBuffs.LeaderOfThePack = proto.TristateEffect_TristateEffectImproved
		}
	}
}

func (druid *Druid) MeleeCritMultiplier() float64 {
	// Assumes that Predatory Instincts is a primary rather than secondary modifier for now, but this needs to confirmed!
	primaryModifier := 1.0
	if druid.InForm(Cat | Bear) {
		primaryModifier = []float64{1, 1.03, 1.07, 1.10}[druid.Talents.PredatoryInstincts]
	}
	return druid.Character.MeleeCritMultiplier(primaryModifier, 0)
}

func (druid *Druid) HasMajorGlyph(glyph proto.DruidMajorGlyph) bool {
	return druid.HasGlyph(int32(glyph))
}
func (druid *Druid) HasMinorGlyph(glyph proto.DruidMinorGlyph) bool {
	return druid.HasGlyph(int32(glyph))
}

func (druid *Druid) Initialize() {
	if druid.Talents.PrimalPrecision > 0 {
		druid.PrimalPrecisionRecoveryMetrics = druid.NewEnergyMetrics(core.ActionID{SpellID: 48410})
	}
	druid.registerFaerieFireSpell()
	druid.registerRebirthSpell()
	druid.registerInnervateCD()

	// Bonus sets
	druid.setBonuses = druidTierSets{
		druid.HasSetBonus(ItemSetThunderheartRegalia, 2),
		druid.HasSetBonus(ItemSetDreamwalkerGarb, 2),
		druid.HasSetBonus(ItemSetDreamwalkerGarb, 4),
		druid.HasSetBonus(ItemSetNightsongGarb, 2),
		druid.HasSetBonus(ItemSetNightsongGarb, 4),
		druid.HasSetBonus(ItemSetMalfurionsRegalia, 2) || druid.HasSetBonus(ItemSetRunetotemsRegalia, 2),
		druid.HasSetBonus(ItemSetMalfurionsRegalia, 4) || druid.HasSetBonus(ItemSetRunetotemsRegalia, 4),
		druid.HasSetBonus(ItemSetLasherweaveRegalia, 2),
		druid.HasSetBonus(ItemSetLasherweaveRegalia, 4),
		druid.HasSetBonus(ItemSetGladiatorsWildhide, 2),
		druid.HasSetBonus(ItemSetGladiatorsWildhide, 4),
		druid.HasSetBonus(ItemSetNightsongBattlegear, 2),
		druid.HasSetBonus(ItemSetNightsongBattlegear, 4),
	}
}

func (druid *Druid) RegisterBalanceSpells() {
	druid.registerHurricaneSpell()
	druid.registerInsectSwarmSpell()
	druid.registerMoonfireSpell()
	druid.registerStarfireSpell()
	druid.registerWrathSpell()
	druid.registerStarfallSpell()
	druid.registerTyphoonSpell()
	druid.registerForceOfNatureCD()
	druid.registerLasherweaveDot()
}

func (druid *Druid) RegisterFeralSpells(maulRageThreshold float64) {
	druid.registerBerserkCD()
	druid.registerCatFormSpell()
	druid.registerBearFormSpell()
	druid.registerDemoralizingRoarSpell()
	druid.registerEnrageSpell()
	druid.registerFerociousBiteSpell()
	druid.registerMangleBearSpell()
	druid.registerMangleCatSpell()
	druid.registerMaulSpell(maulRageThreshold)
	druid.registerLacerateSpell()
	druid.registerRakeSpell()
	druid.registerRipSpell()
	druid.registerSavageRoarSpell()
	druid.registerShredSpell()
	druid.registerSwipeBearSpell()
	druid.registerSwipeCatSpell()
	druid.registerTigersFurySpell()
}

func (druid *Druid) Reset(_ *core.Simulation) {
	druid.BleedsActive = 0
	druid.form = druid.StartingForm
	druid.disabledMCDs = []*core.MajorCooldown{}
	druid.RebirthUsed = false
	druid.LunarICD.Timer.Reset()
	druid.SolarICD.Timer.Reset()
}

func New(char core.Character, form DruidForm, selfBuffs SelfBuffs, talents *proto.DruidTalents) *Druid {
	druid := &Druid{
		Character:    char,
		SelfBuffs:    selfBuffs,
		Talents:      talents,
		StartingForm: form,
		form:         form,
	}
	druid.EnableManaBar()

	druid.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	// Druids get 0.012 crit per agi at level 80, roughly 1 per 83.33
	druid.AddStatDependency(stats.Agility, stats.MeleeCrit, (1.0/83.33)*core.CritRatingPerCritChance)
	// Druid get 0.0209 dodge per agi (before dr), roughly 1 per 47.16
	druid.AddStatDependency(stats.Agility, stats.Dodge, (1.0/47.16)*core.DodgeRatingPerDodgeChance)

	// Druids get extra melee haste
	druid.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	// Base dodge is unaffected by Diminishing Returns
	druid.PseudoStats.BaseDodge += 0.0559

	if druid.Talents.ForceOfNature {
		druid.Treant1 = druid.NewTreant()
		druid.Treant2 = druid.NewTreant()
		druid.Treant3 = druid.NewTreant()
	}

	return druid
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassDruid}] = stats.Stats{
		stats.Health:      6892, // 8227 health shown on naked character (would include tauren bonus)
		stats.Strength:    94,
		stats.Agility:     78,
		stats.Stamina:     99,
		stats.Intellect:   139,
		stats.Spirit:      161,
		stats.Mana:        3496,                                // 5301 mana shown on naked character
		stats.SpellCrit:   1.85 * core.CritRatingPerCritChance, // Class-specific constant
		stats.MeleeCrit:   7.48 * core.CritRatingPerCritChance, // 8.41% chance to crit shown on naked character screen
		stats.AttackPower: -20,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassDruid}] = stats.Stats{
		stats.Health:      7237, // 8217 health shown on naked character
		stats.Strength:    85,
		stats.Agility:     86,
		stats.Stamina:     98,
		stats.Intellect:   143,
		stats.Spirit:      159,
		stats.Mana:        3496,                                // 5361 mana shown on naked character
		stats.SpellCrit:   1.85 * core.CritRatingPerCritChance, // Class-specific constant
		stats.MeleeCrit:   7.48 * core.CritRatingPerCritChance, // 8.51% chance to crit shown on naked character screen
		stats.AttackPower: -20,
	}
}

// Agent is a generic way to access underlying druid on any of the agents (for example balance druid.)
type DruidAgent interface {
	GetDruid() *Druid
}
