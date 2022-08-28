package druid

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"time"
)

type Druid struct {
	core.Character
	SelfBuffs
	Talents proto.DruidTalents

	StartingForm DruidForm

	RebirthUsed       bool
	MaulRageThreshold float64
	RebirthTiming     float64

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
	Swipe            *core.Spell
	TigersFury       *core.Spell
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

	PrimalPrecisionRecoveryMetrics *core.ResourceMetrics

	LunarICD core.Cooldown
	SolarICD core.Cooldown
	Treant1  *TreantPet
	Treant2  *TreantPet
	Treant3  *TreantPet

	form           DruidForm
	disabledMCDs   []*core.MajorCooldown
	SetBonuses     DruidTierSets
	TalentsBonuses TalentsBonuses
}

type TalentsBonuses struct {
	moonfuryMultiplier      float64
	iffBonusCrit            float64
	vengeanceModifier       float64
	genesisMultiplier       float64
	moonglowMultiplier      float64
	naturesMajestyBonusCrit float64
	naturesSplendorTick     int
	starlightWrathModifier  time.Duration
}

type DruidTierSets struct {
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
}

type SelfBuffs struct {
	InnervateTarget proto.RaidTarget
}

// Registering non-unique Talent effects
func (druid *Druid) RegisterTalentsBonuses() {
	druid.TalentsBonuses = TalentsBonuses{
		moonfuryMultiplier:      []float64{0.0, 0.03, 0.06, 0.1}[druid.Talents.Moonfury],
		genesisMultiplier:       1 + 0.01*float64(druid.Talents.Genesis),
		moonglowMultiplier:      1 - 0.03*float64(druid.Talents.Moonglow),
		iffBonusCrit:            float64(druid.Talents.ImprovedFaerieFire) * 1 * core.CritRatingPerCritChance,
		naturesMajestyBonusCrit: 2 * float64(druid.Talents.NaturesMajesty) * core.CritRatingPerCritChance,
		vengeanceModifier:       0.2 * float64(druid.Talents.Vengeance),
		naturesSplendorTick:     core.TernaryInt(druid.Talents.NaturesSplendor, 1, 0),
		starlightWrathModifier:  time.Millisecond * 100 * time.Duration(druid.Talents.StarlightWrath),
	}
}

func (druid *Druid) ResetTalentsBonuses() {
	druid.TalentsBonuses = TalentsBonuses{
		moonfuryMultiplier:      0,
		genesisMultiplier:       0,
		moonglowMultiplier:      0,
		iffBonusCrit:            0,
		naturesMajestyBonusCrit: 0,
		vengeanceModifier:       0,
		naturesSplendorTick:     0,
		starlightWrathModifier:  0,
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
		// Idol of the Raven Goddess
		if druid.Equip[items.ItemSlotRanged].ID == 32387 {
			druid.AddStat(stats.SpellCrit, 40)
		}
	}
	if druid.InForm(Cat|Bear) && druid.Talents.LeaderOfThePack {
		raidBuffs.LeaderOfThePack = core.MaxTristate(raidBuffs.LeaderOfThePack, proto.TristateEffect_TristateEffectRegular)
		if druid.Talents.ImprovedLeaderOfThePack > 0 {
			raidBuffs.LeaderOfThePack = proto.TristateEffect_TristateEffectImproved
		}
	}

}

func (druid *Druid) PrimalGoreOutcomeFuncTick() core.OutcomeApplier {
	if druid.Talents.PrimalGore {
		return druid.OutcomeFuncTickHitAndCrit(druid.MeleeCritMultiplier())
	} else {
		return druid.OutcomeFuncTick()
	}
}

func (druid *Druid) MeleeCritMultiplier() float64 {
	// Assumes that Predatory Instincts is a primary rather than secondary modifier for now, but this needs to confirmed!
	primaryModifier := 1.0
	if druid.InForm(Cat | Bear) {
		primaryModifier = 1 + ((0.1 / 3) * float64(druid.Talents.PredatoryInstincts))
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
	druid.SetBonuses = DruidTierSets{
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
	}
}

func (druid *Druid) RegisterBalanceSpells() {
	druid.registerHurricaneSpell()
	druid.registerInsectSwarmSpell()
	druid.registerMoonfireSpell()
	druid.registerStarfireSpell()
	druid.registerWrathSpell()
	druid.registerStarfallSpell()
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
	druid.registerSwipeSpell()
	druid.registerTigersFurySpell()
}

func (druid *Druid) Reset(sim *core.Simulation) {
	druid.form = druid.StartingForm
	druid.disabledMCDs = []*core.MajorCooldown{}
	druid.RebirthUsed = false
}

func New(char core.Character, form DruidForm, selfBuffs SelfBuffs, talents proto.DruidTalents) *Druid {
	druid := &Druid{
		Character:    char,
		SelfBuffs:    selfBuffs,
		Talents:      talents,
		StartingForm: form,
		form:         form,
	}
	druid.EnableManaBar()

	druid.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	druid.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/25)
	druid.AddStatDependency(stats.Agility, stats.Dodge, core.DodgeRatingPerDodgeChance/14.7059)

	// Druids get extra melee haste
	druid.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	if druid.Talents.ForceOfNature {
		druid.Treant1 = druid.NewTreant()
		druid.Treant2 = druid.NewTreant()
		druid.Treant3 = druid.NewTreant()
	}

	return druid
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassDruid}] = stats.Stats{
		stats.Health:      7237,
		stats.Strength:    85,
		stats.Agility:     86,
		stats.Stamina:     98,
		stats.Intellect:   143,
		stats.Spirit:      159,
		stats.Mana:        3496,
		stats.SpellCrit:   1.85 * core.CritRatingPerCritChance, // Class-specific constant
		stats.AttackPower: -20,                                 // accounts for the fact that the first 20 points in Str only provide 1 AP rather than 2
		stats.MeleeCrit:   0.96 * core.CritRatingPerCritChance, // 3.56% chance to crit shown on naked character screen
		stats.Dodge:       -1.87 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassDruid}] = stats.Stats{
		stats.Health:      7237,
		stats.Strength:    94,
		stats.Agility:     78,
		stats.Stamina:     99,
		stats.Intellect:   139,
		stats.Spirit:      161,
		stats.Mana:        3496,
		stats.SpellCrit:   1.85 * core.CritRatingPerCritChance, // Class-specific constant
		stats.AttackPower: -20,                                 // accounts for the fact that the first 20 points in Str only provide 1 AP rather than 2
		stats.MeleeCrit:   0.96 * core.CritRatingPerCritChance, // 3.96% chance to crit shown on naked character screen
		stats.Dodge:       -1.87 * core.DodgeRatingPerDodgeChance,
	}
}

// Agent is a generic way to access underlying druid on any of the agents (for example balance druid.)
type DruidAgent interface {
	GetDruid() *Druid
}
