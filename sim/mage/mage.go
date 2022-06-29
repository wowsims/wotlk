package mage

import (
	"github.com/wowsims/tbc/sim/common"
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

const (
	SpellFlagMage = core.SpellFlagAgentReserved1
)

func RegisterMage() {
	core.RegisterAgentFactory(
		proto.Player_Mage{},
		proto.Spec_SpecMage,
		func(character core.Character, options proto.Player) core.Agent {
			return NewMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Mage)
			if !ok {
				panic("Invalid spec value for Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

type Mage struct {
	core.Character
	Talents proto.MageTalents

	Options        proto.Mage_Options
	RotationType   proto.Mage_Rotation_Type
	ArcaneRotation proto.Mage_Rotation_ArcaneRotation
	FireRotation   proto.Mage_Rotation_FireRotation
	FrostRotation  proto.Mage_Rotation_FrostRotation
	AoeRotation    proto.Mage_Rotation_AoeRotation
	UseAoeRotation bool

	isDoingRegenRotation bool
	tryingToDropStacks   bool
	numCastsDone         int32
	isBlastSpamming      bool
	disabledMCDs         []*core.MajorCooldown

	waterElemental *WaterElemental

	hasTristfal bool

	// Cached values for a few mechanics.
	spellDamageMultiplier float64

	// Current bonus crit from AM+CC interaction.
	bonusAMCCCrit float64

	ArcaneBlast     []*core.Spell
	ArcaneExplosion *core.Spell
	ArcaneMissiles  *core.Spell
	Blizzard        *core.Spell
	Ignite          *core.Spell
	Fireball        *core.Spell
	FireBlast       *core.Spell
	Flamestrike     *core.Spell
	Frostbolt       *core.Spell
	Pyroblast       *core.Spell
	Scorch          *core.Spell
	WintersChill    *core.Spell

	IcyVeins             *core.Spell
	SummonWaterElemental *core.Spell

	ArcaneMissilesDot *core.Dot
	IgniteDots        []*core.Dot
	FireballDot       *core.Dot
	FlamestrikeDot    *core.Dot
	PyroblastDot      *core.Dot

	ArcaneBlastAura  *core.Aura
	ClearcastingAura *core.Aura
	ScorchAura       *core.Aura

	IgniteTickDamage []float64

	manaTracker common.ManaSpendingRateTracker
}

func (mage *Mage) GetCharacter() *core.Character {
	return &mage.Character
}

func (mage *Mage) GetMage() *Mage {
	return mage
}

func (mage *Mage) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ArcaneBrilliance = true
}
func (mage *Mage) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (mage *Mage) Initialize() {
	mage.ArcaneBlast = []*core.Spell{
		mage.newArcaneBlastSpell(0),
		mage.newArcaneBlastSpell(1),
		mage.newArcaneBlastSpell(2),
		mage.newArcaneBlastSpell(3),
	}
	mage.registerArcaneExplosionSpell()
	mage.registerArcaneMissilesSpell()
	mage.registerBlizzardSpell()
	mage.registerFireballSpell()
	mage.registerFireBlastSpell()
	mage.registerFlamestrikeSpell()
	mage.registerFrostboltSpell()
	mage.registerIgniteSpell()
	mage.registerPyroblastSpell()
	mage.registerScorchSpell()
	mage.registerWintersChillSpell()

	mage.registerEvocationCD()
	mage.registerManaGemsCD()

	mage.IgniteDots = []*core.Dot{}
	mage.IgniteTickDamage = []float64{}
	for i := int32(0); i < mage.Env.GetNumTargets(); i++ {
		mage.IgniteTickDamage = append(mage.IgniteTickDamage, 0)
	}
	for i := int32(0); i < mage.Env.GetNumTargets(); i++ {
		mage.IgniteDots = append(mage.IgniteDots, mage.newIgniteDot(mage.Env.GetTargetUnit(i)))
	}
}

func (mage *Mage) Reset(_ *core.Simulation) {
	mage.isDoingRegenRotation = false
	mage.tryingToDropStacks = false
	mage.numCastsDone = 0
	mage.isBlastSpamming = false
	mage.manaTracker.Reset()
	mage.disabledMCDs = nil
	mage.bonusAMCCCrit = 0
}

func NewMage(character core.Character, options proto.Player) *Mage {
	mageOptions := options.GetMage()

	mage := &Mage{
		Character:    character,
		Talents:      *mageOptions.Talents,
		Options:      *mageOptions.Options,
		RotationType: mageOptions.Rotation.Type,

		UseAoeRotation: mageOptions.Rotation.MultiTargetRotation,

		spellDamageMultiplier: 1.0,
		manaTracker:           common.NewManaSpendingRateTracker(),
	}
	mage.EnableManaBar()

	if mage.RotationType == proto.Mage_Rotation_Arcane && mageOptions.Rotation.Arcane != nil {
		mage.ArcaneRotation = *mageOptions.Rotation.Arcane
	} else if mage.RotationType == proto.Mage_Rotation_Fire && mageOptions.Rotation.Fire != nil {
		mage.FireRotation = *mageOptions.Rotation.Fire
	} else if mage.RotationType == proto.Mage_Rotation_Frost && mageOptions.Rotation.Frost != nil {
		mage.FrostRotation = *mageOptions.Rotation.Frost
	}
	if mageOptions.Rotation.Aoe != nil {
		mage.AoeRotation = *mageOptions.Rotation.Aoe
	}

	mage.Character.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Intellect,
		ModifiedStat: stats.SpellCrit,
		Modifier: func(intellect float64, spellCrit float64) float64 {
			return spellCrit + (intellect/81)*core.SpellCritRatingPerCritChance
		},
	})

	mage.Character.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.AttackPower,
		Modifier: func(strength float64, attackPower float64) float64 {
			return attackPower + strength*2
		},
	})

	if mage.Options.Armor == proto.Mage_Options_MageArmor {
		mage.PseudoStats.SpiritRegenRateCasting += 0.3
	} else if mage.Options.Armor == proto.Mage_Options_MoltenArmor {
		mage.AddStat(stats.SpellCrit, 3*core.SpellCritRatingPerCritChance)
	}

	if mage.Talents.SummonWaterElemental {
		mage.waterElemental = mage.NewWaterElemental(mage.FrostRotation.WaterElementalDisobeyChance)
	}

	mage.hasTristfal = ItemSetTirisfalRegalia.CharacterHasSetBonus(&mage.Character, 2)
	return mage
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  30,
		stats.Agility:   41,
		stats.Stamina:   50,
		stats.Intellect: 155,
		stats.Spirit:    144,
		stats.Mana:      2241,
		stats.SpellCrit: core.SpellCritRatingPerCritChance * 0.926,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  34,
		stats.Agility:   36,
		stats.Stamina:   50,
		stats.Intellect: 152,
		stats.Spirit:    147,
		stats.Mana:      2241,
		stats.SpellCrit: core.SpellCritRatingPerCritChance * 0.933,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  28,
		stats.Agility:   42,
		stats.Stamina:   50,
		stats.Intellect: 154.3, // Gnomes start with 162 int, we assume this include racial so / 1.05
		stats.Spirit:    145,
		stats.Mana:      2241,
		stats.SpellCrit: core.SpellCritRatingPerCritChance * 0.93,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  33,
		stats.Agility:   39,
		stats.Stamina:   51,
		stats.Intellect: 151,
		stats.Spirit:    159,
		stats.Mana:      2241,
		stats.SpellCrit: core.SpellCritRatingPerCritChance * 0.926,
	}
	trollStats := stats.Stats{
		stats.Health:    3213,
		stats.Strength:  34,
		stats.Agility:   41,
		stats.Stamina:   52,
		stats.Intellect: 147,
		stats.Spirit:    146,
		stats.Mana:      2241,
		stats.SpellCrit: core.SpellCritRatingPerCritChance * 0.935,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll10, Class: proto.Class_ClassMage}] = trollStats
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll30, Class: proto.Class_ClassMage}] = trollStats
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  32,
		stats.Agility:   37,
		stats.Stamina:   52,
		stats.Intellect: 149,
		stats.Spirit:    150,
		stats.Mana:      2241,
		stats.SpellCrit: core.SpellCritRatingPerCritChance * 0.930,
	}
}

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}
