package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var TalentTreeSizes = [3]int{31, 27, 27}

type WarriorInputs struct {
	StanceSnapshot bool
}

const (
	SpellFlagBloodsurge  = core.SpellFlagAgentReserved1
	SpellFlagWhirlwindOH = core.SpellFlagAgentReserved2
	ArmsTree             = 0
	FuryTree             = 1
	ProtTree             = 2
)

type Warrior struct {
	core.Character

	Talents *proto.WarriorTalents

	WarriorInputs

	// Current state
	Stance               Stance
	RendValidUntil       time.Duration
	BloodsurgeValidUntil time.Duration
	revengeProcAura      *core.Aura
	Ymirjar4pcProcAura   *core.Aura

	// Reaction time values
	reactionTime       time.Duration
	lastBloodsurgeProc time.Duration
	lastOverpowerProc  time.Duration
	LastAMTick         time.Duration

	BattleShout     *core.Spell
	CommandingShout *core.Spell
	BattleStance    *core.Spell
	DefensiveStance *core.Spell
	BerserkerStance *core.Spell

	BerserkerRage        *core.Spell
	Bloodthirst          *core.Spell
	DemoralizingShout    *core.Spell
	Devastate            *core.Spell
	Execute              *core.Spell
	MortalStrike         *core.Spell
	Overpower            *core.Spell
	Rend                 *core.Spell
	Revenge              *core.Spell
	ShieldBlock          *core.Spell
	ShieldSlam           *core.Spell
	Slam                 *core.Spell
	SunderArmor          *core.Spell
	SunderArmorDevastate *core.Spell
	ThunderClap          *core.Spell
	Whirlwind            *core.Spell
	WhirlwindOH          *core.Spell
	DeepWounds           *core.Spell
	Shockwave            *core.Spell
	ConcussionBlow       *core.Spell
	Bladestorm           *core.Spell
	BladestormOH         *core.Spell

	HeroicStrike       *core.Spell
	Cleave             *core.Spell
	curQueueAura       *core.Aura
	curQueuedAutoSpell *core.Spell

	OverpowerAura *core.Aura

	BattleStanceAura    *core.Aura
	DefensiveStanceAura *core.Aura
	BerserkerStanceAura *core.Aura

	BloodsurgeAura  *core.Aura
	SuddenDeathAura *core.Aura
	ShieldBlockAura *core.Aura

	DemoralizingShoutAuras core.AuraArray
	SunderArmorAuras       core.AuraArray
	ThunderClapAuras       core.AuraArray
}

func (warrior *Warrior) GetCharacter() *core.Character {
	return &warrior.Character
}

func (warrior *Warrior) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if warrior.Talents.Rampage {
		raidBuffs.Rampage = true
	}
}

func (warrior *Warrior) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (warrior *Warrior) Initialize() {
	warrior.AutoAttacks.MHConfig().CritMultiplier = warrior.autoCritMultiplier(mh)
	warrior.AutoAttacks.OHConfig().CritMultiplier = warrior.autoCritMultiplier(oh)

	primaryTimer := warrior.NewTimer()
	overpowerRevengeTimer := warrior.NewTimer()

	warrior.reactionTime = time.Millisecond * 500

	warrior.registerShouts()
	warrior.registerStances()
	warrior.registerBerserkerRageSpell()
	warrior.registerBloodthirstSpell(primaryTimer)
	warrior.registerCleaveSpell()
	warrior.registerDemoralizingShoutSpell()
	warrior.registerDevastateSpell()
	warrior.registerExecuteSpell()
	warrior.registerHeroicStrikeSpell()
	warrior.registerMortalStrikeSpell(primaryTimer)
	warrior.registerOverpowerSpell(overpowerRevengeTimer)
	warrior.registerRevengeSpell(overpowerRevengeTimer)
	warrior.registerShieldSlamSpell()
	warrior.registerSlamSpell()
	warrior.registerThunderClapSpell()
	warrior.registerWhirlwindSpell()
	warrior.registerShockwaveSpell()
	warrior.registerConcussionBlowSpell()
	warrior.RegisterHeroicThrow()
	warrior.RegisterRendSpell()

	warrior.SunderArmor = warrior.newSunderArmorSpell(false)
	warrior.SunderArmorDevastate = warrior.newSunderArmorSpell(true)

	warrior.registerBloodrageCD()
}

func (warrior *Warrior) Reset(_ *core.Simulation) {
	warrior.RendValidUntil = 0
	warrior.curQueueAura = nil
	warrior.curQueuedAutoSpell = nil
}

func NewWarrior(character *core.Character, talents string, inputs WarriorInputs) *Warrior {
	warrior := &Warrior{
		Character:     *character,
		Talents:       &proto.WarriorTalents{},
		WarriorInputs: inputs,
	}
	core.FillTalentsProto(warrior.Talents.ProtoReflect(), talents, TalentTreeSizes)

	warrior.PseudoStats.CanParry = true

	warrior.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[character.Class]*core.CritRatingPerCritChance)
	warrior.AddStatDependency(stats.Agility, stats.Dodge, core.DodgeRatingPerDodgeChance/84.746)
	warrior.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	warrior.AddStatDependency(stats.Strength, stats.BlockValue, .5) // 50% block from str
	warrior.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	// Base dodge unaffected by Diminishing Returns
	warrior.PseudoStats.BaseDodge += 0.03664
	warrior.PseudoStats.BaseParry += 0.05

	return warrior
}

type hand int8

const (
	none hand = 0
	mh   hand = 1
	oh   hand = 2
)

func (warrior *Warrior) autoCritMultiplier(hand hand) float64 {
	return warrior.MeleeCritMultiplier(primary(warrior, hand), 0)
}

func primary(warrior *Warrior, hand hand) float64 {
	if warrior.Talents.PoleaxeSpecialization > 0 {
		if (hand == mh && isPoleaxe(warrior.MainHand())) || (hand == oh && isPoleaxe(warrior.OffHand())) {
			return 1 + 0.01*float64(warrior.Talents.PoleaxeSpecialization)
		}
	}
	return 1
}

func isPoleaxe(weapon *core.Item) bool {
	return weapon.WeaponType == proto.WeaponType_WeaponTypeAxe || weapon.WeaponType == proto.WeaponType_WeaponTypePolearm
}

func (warrior *Warrior) critMultiplier(hand hand) float64 {
	return warrior.MeleeCritMultiplier(primary(warrior, hand), 0.1*float64(warrior.Talents.Impale))
}

func (warrior *Warrior) HasMajorGlyph(glyph proto.WarriorMajorGlyph) bool {
	return warrior.HasGlyph(int32(glyph))
}

func (warrior *Warrior) HasMinorGlyph(glyph proto.WarriorMinorGlyph) bool {
	return warrior.HasGlyph(int32(glyph))
}

func (warrior *Warrior) intensifyRageCooldown(baseCd time.Duration) time.Duration {
	baseCd /= 100
	return []time.Duration{baseCd * 100, baseCd * 89, baseCd * 78, baseCd * 67}[warrior.Talents.IntensifyRage]
}

// Agent is a generic way to access underlying warrior on any of the agents.
type WarriorAgent interface {
	GetWarrior() *Warrior
}
