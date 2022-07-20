package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type WarriorInputs struct {
	ShoutType            proto.WarriorShout
	PrecastShout         bool
	PrecastShoutSapphire bool
	PrecastShoutT2       bool
	RampageCDThreshold   time.Duration
}

type Warrior struct {
	core.Character

	Talents proto.WarriorTalents

	WarriorInputs

	// Current state
	Stance              Stance
	overpowerValidUntil time.Duration
	rampageValidUntil   time.Duration
	RevengeValidUntil   time.Duration
	shoutExpiresAt      time.Duration

	// Cached values
	shoutDuration time.Duration
	canShieldSlam bool

	Shout           *core.Spell
	BattleStance    *core.Spell
	DefensiveStance *core.Spell
	BerserkerStance *core.Spell

	BerserkerRage        *core.Spell
	Bloodthirst          *core.Spell
	DemoralizingShout    *core.Spell
	Devastate            *core.Spell
	Execute              *core.Spell
	Hamstring            *core.Spell
	MortalStrike         *core.Spell
	Overpower            *core.Spell
	Rampage              *core.Spell
	Revenge              *core.Spell
	ShieldBlock          *core.Spell
	ShieldSlam           *core.Spell
	Slam                 *core.Spell
	SunderArmor          *core.Spell
	SunderArmorDevastate *core.Spell
	ThunderClap          *core.Spell
	Whirlwind            *core.Spell

	HeroicStrikeOrCleave *core.Spell
	HSOrCleaveQueueAura  *core.Aura
	HSRageThreshold      float64

	BattleStanceAura    *core.Aura
	DefensiveStanceAura *core.Aura
	BerserkerStanceAura *core.Aura

	DemoralizingShoutAura *core.Aura
	BloodFrenzyAuras      []*core.Aura
	ExposeArmorAura       *core.Aura // Warriors don't cast this but they need to check it.
	AcidSpitAura          *core.Aura // Warriors don't cast this but they need to check it.
	RampageAura           *core.Aura
	SunderArmorAura       *core.Aura
	ThunderClapAura       *core.Aura
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
}

func (warrior *Warrior) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (warrior *Warrior) Initialize() {
	warrior.Shout = warrior.makeShoutSpell()

	primaryTimer := warrior.NewTimer()
	overpowerRevengeTimer := warrior.NewTimer()

	warrior.registerStances()
	warrior.registerBerserkerRageSpell()
	warrior.registerBloodthirstSpell(primaryTimer)
	warrior.registerDemoralizingShoutSpell()
	warrior.registerDevastateSpell()
	warrior.registerExecuteSpell()
	warrior.registerHamstringSpell()
	warrior.registerMortalStrikeSpell(primaryTimer)
	warrior.registerOverpowerSpell(overpowerRevengeTimer)
	warrior.registerRampageSpell()
	warrior.registerRevengeSpell(overpowerRevengeTimer)
	warrior.registerShieldBlockSpell()
	warrior.registerShieldSlamSpell(primaryTimer)
	warrior.registerSlamSpell()
	warrior.registerThunderClapSpell()
	warrior.registerWhirlwindSpell()

	warrior.SunderArmor = warrior.newSunderArmorSpell(false)
	warrior.SunderArmorDevastate = warrior.newSunderArmorSpell(true)

	warrior.shoutDuration = time.Duration(float64(time.Minute*2) * (1 + 0.1*float64(warrior.Talents.BoomingVoice)))

	warrior.registerBloodrageCD()
}

func (warrior *Warrior) Reset(sim *core.Simulation) {
	warrior.overpowerValidUntil = 0
	warrior.rampageValidUntil = 0
	warrior.RevengeValidUntil = 0

	warrior.shoutExpiresAt = 0
	if warrior.Shout != nil && warrior.PrecastShout {
		warrior.shoutExpiresAt = warrior.shoutDuration - time.Second*10
	}
}

func NewWarrior(character core.Character, talents proto.WarriorTalents, inputs WarriorInputs) *Warrior {
	warrior := &Warrior{
		Character:     character,
		Talents:       talents,
		WarriorInputs: inputs,
	}

	warrior.PseudoStats.CanParry = true

	warrior.AddStatDependency(stats.Agility, stats.MeleeCrit, (core.CritRatingPerCritChance / 33.0))
	warrior.AddStatDependency(stats.Agility, stats.Dodge, (core.DodgeRatingPerDodgeChance / 30.0))
	warrior.AddStatDependency(stats.Strength, stats.AttackPower, 2.0)
	warrior.AddStatDependency(stats.Strength, stats.BlockValue, 0.05) // 5% block from str

	return warrior
}

func (warrior *Warrior) secondaryCritModifier(applyImpale bool) float64 {
	secondaryModifier := 0.0
	if applyImpale {
		secondaryModifier += 0.1 * float64(warrior.Talents.Impale)
	}
	return secondaryModifier
}
func (warrior *Warrior) critMultiplier(applyImpale bool) float64 {
	return warrior.MeleeCritMultiplier(1.0, warrior.secondaryCritModifier(applyImpale))
}
func (warrior *Warrior) spellCritMultiplier(applyImpale bool) float64 {
	return warrior.SpellCritMultiplier(1.0, warrior.secondaryCritModifier(applyImpale))
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    146,
		stats.Agility:     93,
		stats.Stamina:     132,
		stats.Intellect:   34,
		stats.Spirit:      53,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.CritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    147,
		stats.Agility:     92,
		stats.Stamina:     136,
		stats.Intellect:   32,
		stats.Spirit:      50,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.CritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    140,
		stats.Agility:     99,
		stats.Stamina:     132,
		stats.Intellect:   38,
		stats.Spirit:      51,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.CritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    145,
		stats.Agility:     96,
		stats.Stamina:     133,
		stats.Intellect:   33,
		stats.Spirit:      56,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.CritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    142,
		stats.Agility:     101,
		stats.Stamina:     132,
		stats.Intellect:   33,
		stats.Spirit:      51,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.CritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    148,
		stats.Agility:     93,
		stats.Stamina:     135,
		stats.Intellect:   30,
		stats.Spirit:      54,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.CritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    150,
		stats.Agility:     91,
		stats.Stamina:     135,
		stats.Intellect:   28,
		stats.Spirit:      53,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.CritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    146,
		stats.Agility:     98,
		stats.Stamina:     134,
		stats.Intellect:   29,
		stats.Spirit:      52,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.CritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    144,
		stats.Agility:     94,
		stats.Stamina:     134,
		stats.Intellect:   31,
		stats.Spirit:      56,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.CritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
}

// Agent is a generic way to access underlying warrior on any of the agents.
type WarriorAgent interface {
	GetWarrior() *Warrior
}
