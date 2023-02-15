package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var TalentTreeSizes = [3]int{31, 27, 27}

type WarriorInputs struct {
	ShoutType            proto.WarriorShout
	PrecastShout         bool
	PrecastShoutSapphire bool
	PrecastShoutT2       bool
	RendCdThreshold      time.Duration
	Munch                bool
}

const (
	SpellFlagBloodsurge = core.SpellFlagAgentReserved1
	ArmsTree            = 0
	FuryTree            = 1
	ProtTree            = 2
)

type Warrior struct {
	core.Character

	Talents *proto.WarriorTalents

	WarriorInputs

	// Current state
	Stance                Stance
	overpowerValidUntil   time.Duration
	rendValidUntil        time.Duration
	shoutExpiresAt        time.Duration
	revengeProcAura       *core.Aura
	lastTasteForBloodProc time.Duration
	Ymirjar4pcProcAura    *core.Aura

	munchedDeepWoundsProcs []*core.PendingAction

	// Reaction time values
	reactionTime       time.Duration
	lastBloodsurgeProc time.Duration

	// Cached values
	shoutDuration time.Duration

	Shout           *core.Spell
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

	HeroicStrikeOrCleave     *core.Spell
	HSOrCleaveQueueAura      *core.Aura
	HSRageThreshold          float64
	RendRageThresholdBelow   float64
	RendHealthThresholdAbove float64

	BattleStanceAura    *core.Aura
	DefensiveStanceAura *core.Aura
	BerserkerStanceAura *core.Aura

	BloodsurgeAura  *core.Aura
	SuddenDeathAura *core.Aura
	ShieldBlockAura *core.Aura

	DemoralizingShoutAuras core.AuraArray
	BloodFrenzyAuras       []*core.Aura
	TraumaAuras            []*core.Aura
	SunderArmorAuras       core.AuraArray
	ThunderClapAuras       core.AuraArray
}

func (warrior *Warrior) GetCharacter() *core.Character {
	return &warrior.Character
}

func (warrior *Warrior) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if warrior.ShoutType == proto.WarriorShout_WarriorShoutBattle {
		raidBuffs.BattleShout = core.MaxTristate(raidBuffs.BattleShout, proto.TristateEffect_TristateEffectRegular)
		if warrior.Talents.CommandingPresence == 5 {
			raidBuffs.BattleShout = proto.TristateEffect_TristateEffectImproved
		}
	} else if warrior.ShoutType == proto.WarriorShout_WarriorShoutCommanding {
		raidBuffs.CommandingShout = core.MaxTristate(raidBuffs.CommandingShout, proto.TristateEffect_TristateEffectRegular)
		if warrior.Talents.CommandingPresence == 5 {
			raidBuffs.CommandingShout = proto.TristateEffect_TristateEffectImproved
		}
	}

	if warrior.Talents.Rampage {
		raidBuffs.Rampage = true
	}
}

func (warrior *Warrior) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (warrior *Warrior) Initialize() {
	warrior.AutoAttacks.MHConfig.CritMultiplier = warrior.autoCritMultiplier(mh)
	warrior.AutoAttacks.OHConfig.CritMultiplier = warrior.autoCritMultiplier(oh)

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
	warrior.registerShockwaveSpell()
	warrior.registerConcussionBlowSpell()

	warrior.SunderArmor = warrior.newSunderArmorSpell(false)
	warrior.SunderArmorDevastate = warrior.newSunderArmorSpell(true)

	warrior.shoutDuration = time.Duration(float64(time.Minute*2) * (1 + 0.1*float64(warrior.Talents.BoomingVoice)))

	warrior.registerBloodrageCD()

	warrior.munchedDeepWoundsProcs = make([]*core.PendingAction, warrior.Env.GetNumTargets())
}

func (warrior *Warrior) Reset(_ *core.Simulation) {
	warrior.overpowerValidUntil = 0
	warrior.rendValidUntil = 0

	warrior.shoutExpiresAt = 0
	if warrior.Shout != nil && warrior.PrecastShout {
		warrior.shoutExpiresAt = warrior.shoutDuration - time.Second*10
	}
	for i := range warrior.munchedDeepWoundsProcs {
		warrior.munchedDeepWoundsProcs[i] = nil
	}
}

func NewWarrior(character core.Character, talents string, inputs WarriorInputs) *Warrior {
	warrior := &Warrior{
		Character:     character,
		Talents:       &proto.WarriorTalents{},
		WarriorInputs: inputs,
	}
	core.FillTalentsProto(warrior.Talents.ProtoReflect(), talents, TalentTreeSizes)

	warrior.PseudoStats.CanParry = true

	warrior.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/62.5)
	warrior.AddStatDependency(stats.Agility, stats.Dodge, core.DodgeRatingPerDodgeChance/84.746)
	warrior.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	warrior.AddStatDependency(stats.Strength, stats.BlockValue, .5) // 50% block from str

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
		if (hand == mh && isPoleaxe(warrior.GetMHWeapon())) || (hand == oh && isPoleaxe(warrior.GetOHWeapon())) {
			return 1 + 0.01*float64(warrior.Talents.PoleaxeSpecialization)
		}
	}
	return 1
}

func isPoleaxe(weapon *core.Item) bool {
	return weapon != nil && (weapon.WeaponType == proto.WeaponType_WeaponTypeAxe || weapon.WeaponType == proto.WeaponType_WeaponTypePolearm)
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

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9611,
		stats.Strength:    175,
		stats.Agility:     110,
		stats.Stamina:     159,
		stats.Intellect:   37,
		stats.Spirit:      61,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9651,
		stats.Strength:    179,
		stats.Agility:     109,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      58,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9581,
		stats.Strength:    169,
		stats.Agility:     116,
		stats.Stamina:     159,
		stats.Intellect:   42,
		stats.Spirit:      59,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9621,
		stats.Strength:    174,
		stats.Agility:     113,
		stats.Stamina:     159,
		stats.Intellect:   36,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9611,
		stats.Strength:    179,
		stats.Agility:     118,
		stats.Stamina:     159,
		stats.Intellect:   36,
		stats.Spirit:      59,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9641,
		stats.Strength:    177,
		stats.Agility:     110,
		stats.Stamina:     160,
		stats.Intellect:   33,
		stats.Spirit:      62,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      10047,
		stats.Strength:    179,
		stats.Agility:     108,
		stats.Stamina:     160,
		stats.Intellect:   31,
		stats.Spirit:      61,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9631,
		stats.Strength:    175,
		stats.Agility:     115,
		stats.Stamina:     159,
		stats.Intellect:   32,
		stats.Spirit:      60,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9541,
		stats.Strength:    173,
		stats.Agility:     111,
		stats.Stamina:     159,
		stats.Intellect:   34,
		stats.Spirit:      64,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
	}
}

// Agent is a generic way to access underlying warrior on any of the agents.
type WarriorAgent interface {
	GetWarrior() *Warrior
}
