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

	ShadowfiendAura *core.Aura
	ShadowfiendPet  *Shadowfiend

	// cached cast stuff
	// TODO: aoe multi-target situations will need multiple spells ticking for each target.
	InnerFocusAura     *core.Aura
	MiseryAura         *core.Aura
	ShadowWeavingAura  *core.Aura
	ShadowyInsightAura *core.Aura
	ImprovedSpiritTap  *core.Aura

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
	StarshardsDot      *core.Dot
	VampiricTouchDot   *core.Dot

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
}

func (priest *Priest) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (priest *Priest) Initialize() {

	if priest.Talents.Misery > 0 {
		priest.MiseryAura = core.MiseryAura(priest.CurrentTarget)
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

	if priest.Talents.ImprovedSpiritTap > 0 {
		increase := 1 + 0.05*float64(priest.Talents.ImprovedSpiritTap)
		statDep := priest.NewDynamicMultiplyStat(stats.Spirit, increase)
		priest.ImprovedSpiritTap = priest.GetOrRegisterAura(core.Aura{
			Label:    "Improved Spirit Tap",
			ActionID: core.ActionID{SpellID: 59000},
			Duration: time.Second * 8,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				priest.EnableDynamicStatDep(sim, statDep)
				priest.PseudoStats.SpiritRegenRateCasting += 0.33
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				priest.DisableDynamicStatDep(sim, statDep)
				priest.PseudoStats.SpiritRegenRateCasting -= 0.33
			},
		})
	}

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
		priest.ShadowWeavingAura.AddStack(sim)
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
	priest.ShadowfiendPet = priest.NewShadowfiend()

	if selfBuffs.UseInnerFire {
		multi := 1 + float64(priest.Talents.ImprovedInnerFire)*0.15
		sp := 120.0 * multi
		armor := 2440 * multi
		priest.AddStat(stats.SpellPower, sp)
		priest.AddStat(stats.Armor, armor)
	}

	return priest
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Strength:  39,
		stats.Agility:   45,
		stats.Stamina:   58,
		stats.Intellect: 145,
		stats.Spirit:    166,
		stats.Mana:      3863,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Strength:  41,
		stats.Agility:   41,
		stats.Stamina:   61,
		stats.Intellect: 144,
		stats.Spirit:    150,
		stats.Mana:      3863,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Strength:  36,
		stats.Agility:   50,
		stats.Stamina:   57,
		stats.Intellect: 145,
		stats.Spirit:    151,
		stats.Mana:      3863,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Strength:  40,
		stats.Agility:   42,
		stats.Stamina:   57,
		stats.Intellect: 146,
		stats.Spirit:    153,
		stats.Mana:      3863,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Strength:  38,
		stats.Agility:   43,
		stats.Stamina:   59,
		stats.Intellect: 143,
		stats.Spirit:    156,
		stats.Mana:      3863,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Strength:  40,
		stats.Agility:   47,
		stats.Stamina:   59,
		stats.Intellect: 141,
		stats.Spirit:    152,
		stats.Mana:      3863,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    6960,
		stats.Strength:  36,
		stats.Agility:   47,
		stats.Stamina:   57,
		stats.Intellect: 149,
		stats.Spirit:    150,
		stats.Mana:      3863,
		stats.SpellCrit: core.CritRatingPerCritChance * 1.24,
	}
}

// Agent is a generic way to access underlying priest on any of the agents.
type PriestAgent interface {
	GetPriest() *Priest
}
