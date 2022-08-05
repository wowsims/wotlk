package mage

import (
	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const (
	SpellFlagMage   = core.SpellFlagAgentReserved1
	BarrageSpells   = core.SpellFlagAgentReserved2
	HotStreakSpells = core.SpellFlagAgentReserved3
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

type MageTierSets struct {
	t7_2  bool
	t7_4  bool
	t8_2  bool
	t8_4  bool
	t9_2  bool
	t9_4  bool
	t10_2 bool
	t10_4 bool
}

type Mage struct {
	core.Character
	Talents proto.MageTalents

	Options  proto.Mage_Options
	Rotation proto.Mage_Rotation

	isMissilesBarrage bool
	numCastsDone      int32
	disabledMCDs      []*core.MajorCooldown

	waterElemental *WaterElemental

	// Cached values for a few mechanics.
	spellDamageMultiplier float64

	// Current bonus crit from AM+CC interaction.
	bonusAMCCCrit   float64
	bonusCritDamage float64

	ArcaneBlast     *core.Spell
	ArcaneExplosion *core.Spell
	ArcaneMissiles  *core.Spell
	Blizzard        *core.Spell
	Ignite          *core.Spell
	LivingBomb      *core.Spell
	Fireball        *core.Spell
	FireBlast       *core.Spell
	Flamestrike     *core.Spell
	Frostbolt       *core.Spell
	FrostfireBolt   *core.Spell
	Pyroblast       *core.Spell
	Scorch          *core.Spell
	WintersChill    *core.Spell

	IcyVeins             *core.Spell
	SummonWaterElemental *core.Spell

	ArcaneMissilesDot   *core.Dot
	IgniteDots          []*core.Dot
	LivingBombDots      []*core.Dot
	LivingBombNotActive *llq.Queue
	FireballDot         *core.Dot
	FlamestrikeDot      *core.Dot
	FrostfireDot        *core.Dot
	PyroblastDot        *core.Dot

	ArcaneBlastAura    *core.Aura
	MissileBarrageAura *core.Aura
	ClearcastingAura   *core.Aura
	ScorchAura         *core.Aura
	HotStreakAura      *core.Aura

	IgniteTickDamage []float64

	MageTier MageTierSets

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

	//if mage.Talents.ArcaneEmpowerment == 3 {
	//	raidBuffs.ArcaneEmpowerment = true
	//}
}
func (mage *Mage) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (mage *Mage) Initialize() {

	mage.LivingBombDots = make([]*core.Dot, mage.Env.GetNumTargets())
	mage.LivingBombNotActive = llq.New()

	mage.MageTier = MageTierSets{
		mage.HasSetBonus(ItemSetFrostfireGarb, 2),
		mage.HasSetBonus(ItemSetFrostfireGarb, 4),
		mage.HasSetBonus(ItemSetKirinTorGarb, 2),
		mage.HasSetBonus(ItemSetKirinTorGarb, 4),
		mage.HasSetBonus(ItemSetKhadgarsRegalia, 2) || mage.HasSetBonus(ItemSetSunstridersRegalia, 2),
		mage.HasSetBonus(ItemSetKhadgarsRegalia, 4) || mage.HasSetBonus(ItemSetSunstridersRegalia, 4),
		false,
		false,
	}

	mage.registerArcaneBlastSpell()
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
	mage.registerLivingBombSpell()
	mage.registerFrostfireBoltSpell()

	mage.registerEvocationCD()
	mage.registerManaGemsCD()

	mage.IgniteDots = []*core.Dot{}
	mage.IgniteTickDamage = []float64{}
	for i := int32(0); i < mage.Env.GetNumTargets(); i++ {
		mage.IgniteTickDamage = append(mage.IgniteTickDamage, 0)
		mage.IgniteDots = append(mage.IgniteDots, mage.newIgniteDot(mage.Env.GetTargetUnit(i)))
		mage.LivingBombNotActive.Enqueue(mage.Env.GetTargetUnit(i))
	}
}

func (mage *Mage) Reset(_ *core.Simulation) {
	mage.numCastsDone = 0
	mage.manaTracker.Reset()
	mage.disabledMCDs = nil
	mage.bonusAMCCCrit = 0
	for i := int32(0); i < mage.Env.GetNumTargets(); i++ {
		mage.LivingBombNotActive.Clear()
		mage.LivingBombNotActive.Enqueue(mage.Env.GetTargetUnit(i))
	}
}

func NewMage(character core.Character, options proto.Player) *Mage {
	mageOptions := options.GetMage()

	mage := &Mage{
		Character: character,
		Talents:   *mageOptions.Talents,
		Options:   *mageOptions.Options,
		Rotation:  *mageOptions.Rotation,

		spellDamageMultiplier: 1.0,
		manaTracker:           common.NewManaSpendingRateTracker(),
	}
	mage.EnableManaBar()
	mage.EnableResumeAfterManaWait(mage.tryUseGCD)

	if mage.Options.Armor == proto.Mage_Options_MageArmor {
		mage.PseudoStats.SpiritRegenRateCasting += 0.5
		if mage.HasSetBonus(ItemSetKhadgarsRegalia, 2) {
			mage.PseudoStats.SpiritRegenRateCasting += .1
		}
	} else if mage.Options.Armor == proto.Mage_Options_MoltenArmor {
		//Need to switch to spirit crit calc
		multi := 0.35
		if mage.HasGlyph(int32(proto.MageMajorGlyph_GlyphOfMoltenArmor.Number())) {
			multi += .2
		}
		if mage.HasSetBonus(ItemSetKhadgarsRegalia, 2) {
			multi += .15
		}
		mage.Character.AddStatDependency(stats.Spirit, stats.SpellCrit, 1+multi)
	}

	if mage.Talents.SummonWaterElemental {
		mage.waterElemental = mage.NewWaterElemental(mage.Rotation.WaterElementalDisobeyChance)
	}

	return mage
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  30,
		stats.Agility:   41,
		stats.Stamina:   50,
		stats.Intellect: 185,
		stats.Spirit:    173,
		stats.Mana:      3268,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.926,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  34,
		stats.Agility:   36,
		stats.Stamina:   50,
		stats.Intellect: 182,
		stats.Spirit:    176,
		stats.Mana:      3268,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.933,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  28,
		stats.Agility:   42,
		stats.Stamina:   50,
		stats.Intellect: 193, // Gnomes start with 162 int, we assume this include racial so / 1.05
		stats.Spirit:    174,
		stats.Mana:      2241,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.93,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  33,
		stats.Agility:   39,
		stats.Stamina:   51,
		stats.Intellect: 181,
		stats.Spirit:    179,
		stats.Mana:      2241,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.926,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  34,
		stats.Agility:   41,
		stats.Stamina:   52,
		stats.Intellect: 177,
		stats.Spirit:    175,
		stats.Mana:      3268,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.935,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  32,
		stats.Agility:   37,
		stats.Stamina:   52,
		stats.Intellect: 179,
		stats.Spirit:    179,
		stats.Mana:      3268,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.930,
	}
}

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}
