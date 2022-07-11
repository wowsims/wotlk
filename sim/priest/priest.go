package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type Priest struct {
	core.Character
	SelfBuffs
	Talents proto.PriestTalents

	SurgeOfLight bool

	Latency float64

	// cached cast stuff
	// TODO: aoe multi-target situations will need multiple spells ticking for each target.

	InnerFocusAura       *core.Aura
	MiseryAura           *core.Aura
	ShadowWeavingAura    *core.Aura
	SurgeOfLightProcAura *core.Aura

	DevouringPlague *core.Spell
	HolyFire        *core.Spell
	InnerFocus      *core.Spell
	ShadowWordPain  *core.Spell
	MindBlast       *core.Spell
	MindFlay        []*core.Spell
	ShadowWordDeath *core.Spell
	Shadowfiend     *core.Spell
	Smite           *core.Spell
	Starshards      *core.Spell
	VampiricTouch   *core.Spell

	ShadowWordPainDot  *core.Dot
	DevouringPlagueDot *core.Dot
	HolyFireDot        *core.Dot
	MindFlayDot        []*core.Dot
	ShadowfiendDot     *core.Dot
	StarshardsDot      *core.Dot
	VampiricTouchDot   *core.Dot
}

type SelfBuffs struct {
	UseShadowfiend bool

	PowerInfusionTarget proto.RaidTarget
}

func (priest *Priest) GetCharacter() *core.Character {
	return &priest.Character
}

func (priest *Priest) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ShadowProtection = true

	raidBuffs.PowerWordFortitude = core.MaxTristate(raidBuffs.PowerWordFortitude, core.MakeTristateValue(
		true,
		priest.Talents.ImprovedPowerWordFortitude == 2))

	//raidBuffs.DivineSpirit = core.MaxTristate(raidBuffs.DivineSpirit, core.MakeTristateValue(
	//priest.Talents.DivineSpirit,
	//priest.Talents.ImprovedDivineSpirit == 2))
}

func (priest *Priest) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (priest *Priest) Initialize() {

	if priest.Talents.Misery > 0 {
		priest.MiseryAura = core.MiseryAura(priest.CurrentTarget, priest.Talents.Misery)
	}

	if priest.Talents.ShadowWeaving > 0 {
		priest.ShadowWeavingAura = priest.GetOrRegisterAura(core.Aura{
			Label:     "Shadow Weaving",
			ActionID:  core.ActionID{SpellID: 15258},
			Duration:  time.Second * 15,
			MaxStacks: 5,
			// TODO: This affects all spells not just direct damage. Dot damage should omit multipliers since it's snapshot at cast time.
			// OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			// 	aura.Unit.PseudoStats.ShadowDamageDealtMultiplier /= 1.0 + 0.02*float64(oldStacks)
			// 	aura.Unit.PseudoStats.ShadowDamageDealtMultiplier *= 1.0 + 0.02*float64(newStacks)
			// },
		})
	}

	priest.registerDevouringPlagueSpell()
	priest.registerHolyFireSpell()
	priest.registerShadowWordPainSpell()
	priest.registerMindBlastSpell()
	priest.registerShadowWordDeathSpell()
	priest.registerShadowfiendSpell()
	priest.registerSmiteSpell()
	priest.registerStarshardsSpell()
	priest.registerVampiricTouchSpell()

	priest.registerPowerInfusionCD()

	priest.MindFlay = []*core.Spell{
		nil, // So we can use # of ticks as the index
		priest.newMindFlaySpell(1),
		priest.newMindFlaySpell(2),
		priest.newMindFlaySpell(3),
	}
	priest.MindFlayDot = []*core.Dot{
		nil, // So we can use # of ticks as the index
		priest.newMindFlayDot(1),
		priest.newMindFlayDot(2),
		priest.newMindFlayDot(3),
	}
}

func (priest *Priest) AddShadowWeavingStack(sim *core.Simulation) {
	if priest.Talents.ShadowWeaving == 0 {
		return
	}

	if priest.ShadowWeavingAura.IsActive() {
		priest.ShadowWeavingAura.AddStack(sim)
		priest.ShadowWeavingAura.Refresh(sim)
	} else {
		priest.ShadowWeavingAura.Activate(sim)
	}
}

func (priest *Priest) OnSpellHitAddShadowWeaving() func(*core.Simulation, *core.Spell, *core.SpellEffect) {
	return func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		if spellEffect.Landed() {
			priest.AddShadowWeavingStack(sim)
		}
	}
}

func (priest *Priest) Reset(_ *core.Simulation) {
}

func New(char core.Character, selfBuffs SelfBuffs, talents proto.PriestTalents) *Priest {
	priest := &Priest{
		Character: char,
		SelfBuffs: selfBuffs,
		Talents:   talents,
	}
	priest.EnableManaBar()

	return priest
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    3211,
		stats.Strength:  39,
		stats.Agility:   45,
		stats.Stamina:   58,
		stats.Intellect: 145,
		stats.Spirit:    166,
		stats.Mana:      2620,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    3211,
		stats.Strength:  41,
		stats.Agility:   41,
		stats.Stamina:   61,
		stats.Intellect: 144,
		stats.Spirit:    150,
		stats.Mana:      2620,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    3211,
		stats.Strength:  36,
		stats.Agility:   50,
		stats.Stamina:   57,
		stats.Intellect: 145,
		stats.Spirit:    151,
		stats.Mana:      2620,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    3211,
		stats.Strength:  40,
		stats.Agility:   42,
		stats.Stamina:   57,
		stats.Intellect: 146,
		stats.Spirit:    153,
		stats.Mana:      2620,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    3211,
		stats.Strength:  38,
		stats.Agility:   43,
		stats.Stamina:   59,
		stats.Intellect: 143,
		stats.Spirit:    156,
		stats.Mana:      2620,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    3211,
		stats.Strength:  40,
		stats.Agility:   47,
		stats.Stamina:   59,
		stats.Intellect: 141,
		stats.Spirit:    152,
		stats.Mana:      2620,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    3211,
		stats.Strength:  36,
		stats.Agility:   47,
		stats.Stamina:   57,
		stats.Intellect: 149,
		stats.Spirit:    150,
		stats.Mana:      2620,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
}

// Agent is a generic way to access underlying priest on any of the agents.
type PriestAgent interface {
	GetPriest() *Priest
}
