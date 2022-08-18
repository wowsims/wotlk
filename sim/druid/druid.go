package druid

import (
	"github.com/wowsims/wotlk/sim/core"
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

	DemoralizingRoar *core.Spell
	FaerieFire       *core.Spell
	FerociousBite    *core.Spell
	Hurricane        *core.Spell
	InsectSwarm      *core.Spell
	Lacerate         *core.Spell
	Mangle           *core.Spell
	Maul             *core.Spell
	Moonfire         *core.Spell
	Rebirth          *core.Spell
	Rake             *core.Spell
	Rip              *core.Spell
	Shred            *core.Spell
	Starfire         *core.Spell
	Swipe            *core.Spell
	Wrath            *core.Spell

	CatForm  *core.Spell
	BearForm *core.Spell

	InsectSwarmDot *core.Dot
	LacerateDot    *core.Dot
	MoonfireDot    *core.Dot
	RakeDot        *core.Dot
	RipDot         *core.Dot

	ClearcastingAura     *core.Aura
	DemoralizingRoarAura *core.Aura
	FaerieFireAura       *core.Aura
	MangleAura           *core.Aura
	MaulQueueAura        *core.Aura
	NaturesGraceProcAura *core.Aura
	NaturesSwiftnessAura *core.Aura
	CatFormAura          *core.Aura
	BearFormAura         *core.Aura

	LunarICD core.Cooldown
	SolarICD core.Cooldown

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
	if druid.Talents.ImprovedMarkOfTheWild == 5 { // probably could work on actually calculating the fraction effect later if we care.
		raidBuffs.GiftOfTheWild = proto.TristateEffect_TristateEffectImproved
	}

	raidBuffs.Thorns = core.MaxTristate(raidBuffs.Thorns, proto.TristateEffect_TristateEffectRegular)
	if druid.Talents.Brambles == 3 {
		raidBuffs.Thorns = proto.TristateEffect_TristateEffectImproved
	}

	if druid.InForm(Moonkin) && druid.Talents.MoonkinForm {
		raidBuffs.MoonkinAura = core.MaxTristate(raidBuffs.MoonkinAura, proto.TristateEffect_TristateEffectRegular)
		// if druid.Talents.ImprovedMoonkinForm > 0 {
		// 	raidBuffs.LeaderOfThePack = proto.TristateEffect_TristateEffectImproved
		// }
	}
	if druid.InForm(Cat|Bear) && druid.Talents.LeaderOfThePack {
		raidBuffs.LeaderOfThePack = core.MaxTristate(raidBuffs.LeaderOfThePack, proto.TristateEffect_TristateEffectRegular)
		if druid.Talents.ImprovedLeaderOfThePack > 0 {
			raidBuffs.LeaderOfThePack = proto.TristateEffect_TristateEffectImproved
		}
	}

}

const ravenGoddessItemID = 32387

func (druid *Druid) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (druid *Druid) MeleeCritMultiplier() float64 {
	// Assumes that Predatory Instincts is a primary rather than secondary modifier for now, but this needs to confirmed!
	primaryModifier := 1.0
	if druid.InForm(Cat | Bear) {
		primaryModifier = 1 + 0.02*float64(druid.Talents.PredatoryInstincts)
	}
	return druid.Character.MeleeCritMultiplier(primaryModifier, 0)
}

func (druid *Druid) Initialize() {
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
}

func (druid *Druid) RegisterBearSpells(maulRageThreshold float64) {
	druid.registerBearFormSpell()
	druid.registerMangleBearSpell()
	druid.registerMaulSpell(maulRageThreshold)
	druid.registerLacerateSpell()
	druid.registerSwipeSpell()
	druid.registerDemoralizingRoarSpell()
}

func (druid *Druid) RegisterCatSpells() {
	druid.registerCatFormSpell()
	druid.registerFerociousBiteSpell()
	druid.registerMangleCatSpell()
	druid.registerRipSpell()
	druid.registerShredSpell()
	druid.registerRakeSpell()
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
