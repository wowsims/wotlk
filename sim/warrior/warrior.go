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
	RendCdThreshold      time.Duration
}

type Warrior struct {
	core.Character

	Talents proto.WarriorTalents

	WarriorInputs

	// Current state
	Stance               Stance
	overpowerValidUntil  time.Duration
	rendValidUntil       time.Duration
	RevengeValidUntil    time.Duration
	shoutExpiresAt       time.Duration
	disableHsCleaveUntil time.Duration

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
	DeepWounds           *core.Spell

	RendDots               *core.Dot
	DeepWoundsDots         []*core.Dot
	DeepWoundsTickDamage   []float64
	DeepwoundsDamageBuffer []float64

	HeroicStrikeOrCleave *core.Spell
	HSOrCleaveQueueAura  *core.Aura
	HSRageThreshold      float64

	BattleStanceAura    *core.Aura
	DefensiveStanceAura *core.Aura
	BerserkerStanceAura *core.Aura

	BloodsurgeAura  *core.Aura
	SuddenDeathAura *core.Aura

	DemoralizingShoutAura *core.Aura
	BloodFrenzyAuras      []*core.Aura
	TraumaAuras           []*core.Aura
	ExposeArmorAura       *core.Aura // Warriors don't cast this but they need to check it.
	AcidSpitAura          *core.Aura // Warriors don't cast this but they need to check it.
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

	if warrior.Talents.Rampage {
		raidBuffs.Rampage = true
	}
}

func (warrior *Warrior) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (warrior *Warrior) Initialize() {
	warrior.AutoAttacks.MHEffect.OutcomeApplier = warrior.OutcomeFuncMeleeWhite(warrior.critMultiplier(false))
	warrior.AutoAttacks.OHEffect.OutcomeApplier = warrior.OutcomeFuncMeleeWhite(warrior.critMultiplier(false))

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
	warrior.registerShieldSlamSpell(primaryTimer)
	warrior.registerSlamSpell()
	warrior.registerThunderClapSpell()
	warrior.registerWhirlwindSpell()
	warrior.registerRendSpell()

	warrior.SunderArmor = warrior.newSunderArmorSpell(false)
	warrior.SunderArmorDevastate = warrior.newSunderArmorSpell(true)

	warrior.shoutDuration = time.Duration(float64(time.Minute*2) * (1 + 0.1*float64(warrior.Talents.BoomingVoice)))

	warrior.registerBloodrageCD()

	warrior.DeepwoundsDamageBuffer = []float64{}
	for i := int32(0); i < warrior.Env.GetNumTargets(); i++ {
		warrior.DeepwoundsDamageBuffer = append(warrior.DeepwoundsDamageBuffer, 0)
	}
	warrior.DeepWoundsTickDamage = []float64{}
	for i := int32(0); i < warrior.Env.GetNumTargets(); i++ {
		warrior.DeepWoundsTickDamage = append(warrior.DeepWoundsTickDamage, 0)
	}
	warrior.DeepWoundsDots = []*core.Dot{}
	for i := int32(0); i < warrior.Env.GetNumTargets(); i++ {
		warrior.DeepWoundsDots = append(warrior.DeepWoundsDots, warrior.newDeepWoundsDot(warrior.Env.GetTargetUnit(i)))
	}
}

func (warrior *Warrior) Reset(sim *core.Simulation) {
	warrior.overpowerValidUntil = 0
	warrior.rendValidUntil = 0
	warrior.RevengeValidUntil = 0
	warrior.disableHsCleaveUntil = 0
	warrior.lastBloodsurgeProc = 0

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

	warrior.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/62.5)
	warrior.AddStatDependency(stats.Agility, stats.Dodge, core.DodgeRatingPerDodgeChance/85.1)
	warrior.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	warrior.AddStatDependency(stats.Strength, stats.BlockValue, .05) // 5% block from str

	return warrior
}

func (warrior *Warrior) secondaryCritModifier(applyImpale bool) float64 {
	secondaryModifier := 0.0
	if applyImpale {
		secondaryModifier += 0.1 * float64(warrior.Talents.Impale)
	}
	if warrior.Talents.PoleaxeSpecialization > 0 {
		secondaryModifier += 0.01 * float64(warrior.Talents.PoleaxeSpecialization)
	}
	return secondaryModifier
}
func (warrior *Warrior) critMultiplier(applyImpale bool) float64 {
	return warrior.MeleeCritMultiplier(1.0, warrior.secondaryCritModifier(applyImpale))
}
func (warrior *Warrior) spellCritMultiplier(applyImpale bool) float64 {
	return warrior.SpellCritMultiplier(1.0, warrior.secondaryCritModifier(applyImpale))
}

func (warrior *Warrior) HasMajorGlyph(glyph proto.WarriorMajorGlyph) bool {
	return warrior.HasGlyph(int32(glyph))
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9611,
		stats.Strength:    185,
		stats.Agility:     110,
		stats.Stamina:     167,
		stats.Intellect:   37,
		stats.Spirit:      61,
		stats.AttackPower: 590,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9651,
		stats.Strength:    186,
		stats.Agility:     109,
		stats.Stamina:     171,
		stats.Intellect:   35,
		stats.Spirit:      58,
		stats.AttackPower: 592,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9581,
		stats.Strength:    175,
		stats.Agility:     116,
		stats.Stamina:     164,
		stats.Intellect:   42,
		stats.Spirit:      59,
		stats.AttackPower: 570,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9621,
		stats.Strength:    184,
		stats.Agility:     113,
		stats.Stamina:     168,
		stats.Intellect:   36,
		stats.Spirit:      63,
		stats.AttackPower: 588,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9611,
		stats.Strength:    181,
		stats.Agility:     118,
		stats.Stamina:     167,
		stats.Intellect:   36,
		stats.Spirit:      59,
		stats.AttackPower: 582,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9641,
		stats.Strength:    187,
		stats.Agility:     110,
		stats.Stamina:     170,
		stats.Intellect:   33,
		stats.Spirit:      62,
		stats.AttackPower: 594,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      10047,
		stats.Strength:    179,
		stats.Agility:     108,
		stats.Stamina:     170,
		stats.Intellect:   31,
		stats.Spirit:      61,
		stats.AttackPower: 578,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9631,
		stats.Strength:    185,
		stats.Agility:     115,
		stats.Stamina:     169,
		stats.Intellect:   32,
		stats.Spirit:      60,
		stats.AttackPower: 590,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      9541,
		stats.Strength:    173,
		stats.Agility:     111,
		stats.Stamina:     160,
		stats.Intellect:   34,
		stats.Spirit:      64,
		stats.AttackPower: 566,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
}

// Agent is a generic way to access underlying warrior on any of the agents.
type WarriorAgent interface {
	GetWarrior() *Warrior
}
