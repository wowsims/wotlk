package druid

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
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
	MoonfireDot       *core.Dot
	RakeDot           *core.Dot
	RipDot            *core.Dot
	StarfallDot       *core.Dot
	StarfallDotSplash *core.Dot

	BearFormAura         *core.Aura
	BerserkAura          *core.Aura
	CatFormAura          *core.Aura
	ClearcastingAura     *core.Aura
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

	form         DruidForm
	disabledMCDs []*core.MajorCooldown
}

type SelfBuffs struct {
	Omen bool

	InnervateTarget proto.RaidTarget
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
}

func (druid *Druid) RegisterBalanceSpells() {
	druid.registerHurricaneSpell()
	druid.registerInsectSwarmSpell()
	druid.registerMoonfireSpell()
	druid.Starfire = druid.newStarfireSpell()
	druid.registerWrathSpell()
	druid.registerStarfallSpell()
	druid.registerForceOfNatureCD()
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
		stats.Health:      3434, // 4498 health shown on naked character (would include tauren bonus)
		stats.Strength:    81,
		stats.Agility:     65,
		stats.Stamina:     85,
		stats.Intellect:   115,
		stats.Spirit:      135,
		stats.Mana:        2370,
		stats.SpellCrit:   40.66,                               // 3.29% chance to crit shown on naked character screen
		stats.AttackPower: -20,                                 // accounts for the fact that the first 20 points in Str only provide 1 AP rather than 2
		stats.MeleeCrit:   0.96 * core.CritRatingPerCritChance, // 3.56% chance to crit shown on naked character screen
		stats.Dodge:       -1.87 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassDruid}] = stats.Stats{
		stats.Health:      3434, // 4254 health shown on naked character
		stats.Strength:    73,
		stats.Agility:     75,
		stats.Stamina:     82,
		stats.Intellect:   120,
		stats.Spirit:      133,
		stats.Mana:        2370,
		stats.SpellCrit:   40.60,                               // 3.35% chance to crit shown on naked character screen
		stats.AttackPower: -20,                                 // accounts for the fact that the first 20 points in Str only provide 1 AP rather than 2
		stats.MeleeCrit:   0.96 * core.CritRatingPerCritChance, // 3.96% chance to crit shown on naked character screen
		stats.Dodge:       -1.87 * core.DodgeRatingPerDodgeChance,
	}
}

// Agent is a generic way to access underlying druid on any of the agents (for example balance druid.)
type DruidAgent interface {
	GetDruid() *Druid
}
