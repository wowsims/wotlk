package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	SpellFlagBleed = core.SpellFlagAgentReserved1
)

var TalentTreeSizes = [3]int{31, 27, 27}

type WarriorInputs struct {
	ShoutType      proto.WarriorShout
	StanceSnapshot bool
}

const (
	ArmsTree = 0
	FuryTree = 1
	ProtTree = 2
)

type Warrior struct {
	core.Character

	Talents *proto.WarriorTalents

	WarriorInputs

	// Current state
	Stance               Stance
	revengeProcAura      *core.Aura
	OverpowerAura        *core.Aura
	BerserkerRageAura    *core.Aura
	BloodrageAura        *core.Aura
	ConsumedByRageAura   *core.Aura
	Above80RageCBRActive bool

	// Reaction time values
	reactionTime       time.Duration
	lastBloodsurgeProc time.Duration
	lastOverpowerProc  time.Duration
	LastAMTick         time.Duration

	Shout           *core.Spell
	BattleStance    *core.Spell
	DefensiveStance *core.Spell
	BerserkerStance *core.Spell

	Bloodrage            *core.Spell
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
	DeepWounds           *core.Spell
	ConcussionBlow       *core.Spell
	RagingBlow           *core.Spell

	HeroicStrike         *core.Spell
	QuickStrike          *core.Spell
	Cleave               *core.Spell
	hsOrCleaveQueueSpell *core.Spell
	curQueueAura         *core.Aura
	curQueuedAutoSpell   *core.Spell

	BattleStanceAura    *core.Aura
	DefensiveStanceAura *core.Aura
	BerserkerStanceAura *core.Aura

	ShieldBlockAura *core.Aura

	DemoralizingShoutAuras core.AuraArray
	SunderArmorAuras       core.AuraArray
	ThunderClapAuras       core.AuraArray
}

func (warrior *Warrior) GetCharacter() *core.Character {
	return &warrior.Character
}

func (warrior *Warrior) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
}

func (warrior *Warrior) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (warrior *Warrior) Initialize() {
	warrior.AutoAttacks.MHConfig().CritMultiplier = warrior.autoCritMultiplier(mh)
	warrior.AutoAttacks.OHConfig().CritMultiplier = warrior.autoCritMultiplier(oh)

	warrior.Shout = warrior.makeShoutSpell()

	primaryTimer := warrior.NewTimer()
	overpowerRevengeTimer := warrior.NewTimer()

	warrior.reactionTime = time.Millisecond * 500

	warrior.registerStances()
	warrior.registerBerserkerRageSpell()
	warrior.registerBloodthirstSpell(primaryTimer)
	warrior.registerDemoralizingShoutSpell()
	warrior.registerDevastateSpell()
	warrior.registerExecuteSpell()
	warrior.registerMortalStrikeSpell(primaryTimer)
	warrior.registerOverpowerSpell(overpowerRevengeTimer)
	warrior.registerRevengeSpell(overpowerRevengeTimer)
	warrior.registerShieldSlamSpell()
	warrior.registerSlamSpell()
	warrior.registerThunderClapSpell()
	warrior.registerWhirlwindSpell()
	warrior.registerConcussionBlowSpell()

	warrior.SunderArmor = warrior.newSunderArmorSpell(false)
	warrior.SunderArmorDevastate = warrior.newSunderArmorSpell(true)

	warrior.registerBloodrageCD()
}

func (warrior *Warrior) Reset(_ *core.Simulation) {
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

	warrior.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(warrior.Level)])
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
	return 1
}

func isPoleaxe(weapon *core.Item) bool {
	return weapon.WeaponType == proto.WeaponType_WeaponTypeAxe || weapon.WeaponType == proto.WeaponType_WeaponTypePolearm
}

func (warrior *Warrior) critMultiplier(hand hand) float64 {
	return warrior.MeleeCritMultiplier(primary(warrior, hand), 0.1*float64(warrior.Talents.Impale))
}

// Agent is a generic way to access underlying warrior on any of the agents.
type WarriorAgent interface {
	GetWarrior() *Warrior
}

func (warrior *Warrior) HasRune(rune proto.WarriorRune) bool {
	return warrior.HasRuneById(int32(rune))
}
