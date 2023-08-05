package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var TalentTreeSizes = [3]int{28, 27, 27}

type Priest struct {
	core.Character
	SelfBuffs
	Talents *proto.PriestTalents

	SurgeOfLight bool

	Latency float64

	ShadowfiendAura *core.Aura
	ShadowfiendPet  *Shadowfiend

	// cached cast stuff
	// TODO: aoe multi-target situations will need multiple spells ticking for each target.
	InnerFocusAura     *core.Aura
	ShadowWeavingAura  *core.Aura
	ShadowyInsightAura *core.Aura
	ImprovedSpiritTap  *core.Aura
	DispersionAura     *core.Aura

	SurgeOfLightProcAura *core.Aura

	BindingHeal     *core.Spell
	CircleOfHealing *core.Spell
	DevouringPlague *core.Spell
	FlashHeal       *core.Spell
	GreaterHeal     *core.Spell
	HolyFire        *core.Spell
	InnerFocus      *core.Spell
	ShadowWordPain  *core.Spell
	MindBlast       *core.Spell
	MindFlay        []*core.Spell
	MindSear        []*core.Spell
	Penance         *core.Spell
	PenanceHeal     *core.Spell
	PowerWordShield *core.Spell
	PrayerOfHealing *core.Spell
	PrayerOfMending *core.Spell
	Renew           *core.Spell
	EmpoweredRenew  *core.Spell
	ShadowWordDeath *core.Spell
	Shadowfiend     *core.Spell
	Smite           *core.Spell
	VampiricTouch   *core.Spell
	Dispersion      *core.Spell

	PWSShields    []*core.Shield
	WeakenedSouls core.AuraArray

	ProcPrayerOfMending core.ApplySpellResults

	DpInitMultiplier float64

	// set bonus cache
	// The mana cost of your Mind Blast is reduced by 10%.
	T7TwoSetBonus bool
	// Your Shadow Word: Death has an additional 10% chance to critically strike.
	T7FourSetBonus bool
	// Increases the damage done by your Devouring Plague by 15%.
	T8TwoSetBonus bool
	// Your Mind Blast also grants you 240 haste for 4 sec.
	T8FourSetBonus bool
	// Increases the duration of your Vampiric Touch spell by 6 sec.
	T9TwoSetBonus bool
	// Increases the critical strike chance of your Mind Flay spell by 5%.
	T9FourSetBonus bool
	// The critical strike chance of your Shadow Word: Pain, Devouring Plague, and Vampiric Touch spells is increased by 5%
	T10TwoSetBonus bool
	// Reduces the channel duration by 0.51 sec and period by 0.17 sec on your Mind Flay spell
	T10FourSetBonus bool
}

type SelfBuffs struct {
	UseShadowfiend bool
	UseInnerFire   bool

	PowerInfusionTarget *proto.UnitReference
}

func (priest *Priest) GetCharacter() *core.Character {
	return &priest.Character
}

func (priest *Priest) HasMajorGlyph(glyph proto.PriestMajorGlyph) bool {
	return priest.HasGlyph(int32(glyph))
}
func (priest *Priest) HasMinorGlyph(glyph proto.PriestMinorGlyph) bool {
	return priest.HasGlyph(int32(glyph))
}

func (priest *Priest) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ShadowProtection = true
	raidBuffs.DivineSpirit = true

	raidBuffs.PowerWordFortitude = core.MaxTristate(raidBuffs.PowerWordFortitude, core.MakeTristateValue(
		true,
		priest.Talents.ImprovedPowerWordFortitude == 2))
}

func (priest *Priest) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (priest *Priest) Initialize() {
	// Shadow Insight gained from Glyph of Shadow
	// Finalized spirit off gear and not dynamic spirit (e.g. Spirit Tap does not increase this)
	priest.ShadowyInsightAura = priest.NewTemporaryStatsAura(
		"Shadowy Insight",
		core.ActionID{SpellID: 61792},
		stats.Stats{stats.SpellPower: priest.GetStat(stats.Spirit) * 0.30},
		time.Second*10,
	)

	priest.registerSetBonuses()
	priest.registerDevouringPlagueSpell()
	priest.registerShadowWordPainSpell()
	priest.registerMindBlastSpell()
	priest.registerShadowWordDeathSpell()
	priest.registerShadowfiendSpell()
	priest.registerVampiricTouchSpell()
	priest.registerDispersionSpell()

	priest.registerPowerInfusionCD()

	priest.MindFlay = []*core.Spell{
		nil, // So we can use # of ticks as the index
		priest.newMindFlaySpell(1),
		priest.newMindFlaySpell(2),
		priest.newMindFlaySpell(3),
	}
	priest.MindSear = []*core.Spell{
		nil, // So we can use # of ticks as the index
		priest.newMindSearSpell(1),
		priest.newMindSearSpell(2),
		priest.newMindSearSpell(3),
		priest.newMindSearSpell(4),
		priest.newMindSearSpell(5),
	}
}

func (priest *Priest) RegisterHealingSpells() {
	priest.registerPenanceHealSpell()
	priest.registerBindingHealSpell()
	priest.registerCircleOfHealingSpell()
	priest.registerFlashHealSpell()
	priest.registerGreaterHealSpell()
	priest.registerPowerWordShieldSpell()
	priest.registerPrayerOfHealingSpell()
	priest.registerPrayerOfMendingSpell()
	priest.registerRenewSpell()
}

func (priest *Priest) AddShadowWeavingStack(sim *core.Simulation) {
	if priest.ShadowWeavingAura != nil {
		priest.ShadowWeavingAura.Activate(sim)
		priest.ShadowWeavingAura.AddStack(sim)
	}
}

func (priest *Priest) Reset(_ *core.Simulation) {
}

func New(char core.Character, selfBuffs SelfBuffs, talents string) *Priest {
	priest := &Priest{
		Character: char,
		SelfBuffs: selfBuffs,
		Talents:   &proto.PriestTalents{},
	}
	core.FillTalentsProto(priest.Talents.ProtoReflect(), talents, TalentTreeSizes)

	priest.EnableManaBar()
	priest.ShadowfiendPet = priest.NewShadowfiend()

	if selfBuffs.UseInnerFire {
		multi := 1 + float64(priest.Talents.ImprovedInnerFire)*0.15
		sp := 120.0 * multi
		armor := 2440 * multi * core.TernaryFloat64(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfInnerFire), 1.5, 1)
		priest.AddStat(stats.SpellPower, sp)
		priest.AddStat(stats.Armor, armor)
	}

	return priest
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Mana:      3863,
		stats.Strength:  39,
		stats.Agility:   45,
		stats.Stamina:   58,
		stats.Intellect: 145,
		stats.Spirit:    166,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Mana:      3863,
		stats.Strength:  41,
		stats.Agility:   41,
		stats.Stamina:   61,
		stats.Intellect: 144,
		stats.Spirit:    150,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Mana:      3863,
		stats.Strength:  36,
		stats.Agility:   50,
		stats.Stamina:   57,
		stats.Intellect: 145,
		stats.Spirit:    151,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Mana:      3863,
		stats.Strength:  40,
		stats.Agility:   42,
		stats.Stamina:   57,
		stats.Intellect: 146,
		stats.Spirit:    153,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Mana:      3863,
		stats.Strength:  42,
		stats.Agility:   49,
		stats.Stamina:   67,
		stats.Intellect: 172,
		stats.Spirit:    186,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Mana:      3863,
		stats.Strength:  44,
		stats.Agility:   53,
		stats.Stamina:   67,
		stats.Intellect: 170,
		stats.Spirit:    182,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Mana:      3863,
		stats.Strength:  40,
		stats.Agility:   53,
		stats.Stamina:   67,
		stats.Intellect: 177,
		stats.Spirit:    179,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
}

// Agent is a generic way to access underlying priest on any of the agents.
type PriestAgent interface {
	GetPriest() *Priest
}
