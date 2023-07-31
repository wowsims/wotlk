package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const (
	SpellFlagMage   = core.SpellFlagAgentReserved1
	BarrageSpells   = core.SpellFlagAgentReserved2
	HotStreakSpells = core.SpellFlagAgentReserved3
)

var TalentTreeSizes = [3]int{30, 28, 28}

func RegisterMage() {
	core.RegisterAgentFactory(
		proto.Player_Mage{},
		proto.Spec_SpecMage,
		func(character core.Character, options *proto.Player) core.Agent {
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

	Talents  *proto.MageTalents
	Options  *proto.Mage_Options
	Rotation *proto.Mage_Rotation

	ReactionTime     time.Duration
	PyroblastDelayMs time.Duration

	arcaneBlastStreak int32
	arcanePowerMCD    *core.MajorCooldown
	delayedPyroAt     time.Duration

	waterElemental *WaterElemental
	mirrorImage    *MirrorImage

	// Cached values for a few mechanics.
	bonusCritDamage float64

	ArcaneBarrage   *core.Spell
	ArcaneBlast     *core.Spell
	ArcaneExplosion *core.Spell
	ArcaneMissiles  *core.Spell
	Blizzard        *core.Spell
	DeepFreeze      *core.Spell
	Ignite          *core.Spell
	LivingBomb      *core.Spell
	Fireball        *core.Spell
	FireBlast       *core.Spell
	Flamestrike     *core.Spell
	Frostbolt       *core.Spell
	FrostfireBolt   *core.Spell
	IceLance        *core.Spell
	Pyroblast       *core.Spell
	Scorch          *core.Spell
	MirrorImage     *core.Spell

	IcyVeins             *core.Spell
	SummonWaterElemental *core.Spell

	ArcaneBlastAura    *core.Aura
	ArcanePotencyAura  *core.Aura
	ArcanePowerAura    *core.Aura
	MissileBarrageAura *core.Aura
	ClearcastingAura   *core.Aura
	ScorchAuras        core.AuraArray
	hotStreakCritAura  *core.Aura
	HotStreakAura      *core.Aura
	CombustionAura     *core.Aura
	FingersOfFrostAura *core.Aura
	BrainFreezeAura    *core.Aura

	CritDebuffCategories core.ExclusiveCategoryArray
}

func (mage *Mage) GetCharacter() *core.Character {
	return &mage.Character
}

func (mage *Mage) GetMage() *Mage {
	return mage
}

func (mage *Mage) HasMajorGlyph(glyph proto.MageMajorGlyph) bool {
	return mage.HasGlyph(int32(glyph))
}
func (mage *Mage) HasMinorGlyph(glyph proto.MageMinorGlyph) bool {
	return mage.HasGlyph(int32(glyph))
}

func (mage *Mage) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ArcaneBrilliance = true

	if mage.Talents.ArcaneEmpowerment == 3 {
		raidBuffs.ArcaneEmpowerment = true
	}
}
func (mage *Mage) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (mage *Mage) Initialize() {
	mage.registerArcaneBarrageSpell()
	mage.registerArcaneBlastSpell()
	mage.registerArcaneExplosionSpell()
	mage.registerArcaneMissilesSpell()
	mage.registerBlizzardSpell()
	mage.registerDeepFreezeSpell()
	mage.registerFireballSpell()
	mage.registerFireBlastSpell()
	mage.registerFlamestrikeSpell()
	mage.registerFrostboltSpell()
	mage.registerIceLanceSpell()
	mage.registerPyroblastSpell()
	mage.registerScorchSpell()
	mage.registerLivingBombSpell()
	mage.registerFrostfireBoltSpell()

	mage.registerEvocationCD()
	mage.registerManaGemsCD()
	mage.registerMirrorImageCD()

	if mirrorImageMCD := mage.GetMajorCooldownIgnoreTag(mage.MirrorImage.ActionID); mirrorImageMCD != nil {
		if len(mirrorImageMCD.GetTimings()) == 0 {
			mage.RegisterPrepullAction(-1500*time.Millisecond, func(sim *core.Simulation) {
				mage.MirrorImage.Cast(sim, nil)
			})
		}
	}
}

func (mage *Mage) Reset(sim *core.Simulation) {
	mage.arcaneBlastStreak = 0
	mage.arcanePowerMCD = mage.GetMajorCooldown(core.ActionID{SpellID: 12042})
	mage.delayedPyroAt = 0
}

func NewMage(character core.Character, options *proto.Player) *Mage {
	mageOptions := options.GetMage()

	mage := &Mage{
		Character: character,
		Talents:   &proto.MageTalents{},
		Options:   mageOptions.Options,
		Rotation:  mageOptions.Rotation,

		ReactionTime:     time.Millisecond * time.Duration(mageOptions.Options.ReactionTimeMs),
		PyroblastDelayMs: time.Millisecond * time.Duration(mageOptions.Rotation.PyroblastDelayMs),
	}
	core.FillTalentsProto(mage.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	mage.bonusCritDamage = .25*float64(mage.Talents.SpellPower) + .1*float64(mage.Talents.Burnout)
	mage.EnableManaBar()
	mage.EnableResumeAfterManaWait(mage.tryUseGCD)

	if !mage.Talents.ArcaneBarrage {
		mage.Rotation.UseArcaneBarrage = false
	}
	if mage.Talents.ImprovedScorch == 0 {
		mage.Rotation.MaintainImprovedScorch = false
	}
	if !mage.Options.IgniteMunching || mage.Rotation.PrimaryFireSpell != proto.Mage_Rotation_Fireball {
		mage.PyroblastDelayMs = 0
	}

	if mage.Options.Armor == proto.Mage_Options_MageArmor {
		mage.PseudoStats.SpiritRegenRateCasting += .5
		if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfMageArmor) {
			mage.PseudoStats.SpiritRegenRateCasting += .2
		}
		if mage.HasSetBonus(ItemSetKhadgarsRegalia, 2) {
			mage.PseudoStats.SpiritRegenRateCasting += .1
		}
	} else if mage.Options.Armor == proto.Mage_Options_MoltenArmor {
		//Need to switch to spirit crit calc
		multi := 0.35
		if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfMoltenArmor) {
			multi += .2
		}
		if mage.HasSetBonus(ItemSetKhadgarsRegalia, 2) {
			multi += .15
		}
		mage.Character.AddStatDependency(stats.Spirit, stats.SpellCrit, multi)
	}

	mage.mirrorImage = mage.NewMirrorImage()

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
		stats.Mana:      3268,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.93,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  33,
		stats.Agility:   39,
		stats.Stamina:   51,
		stats.Intellect: 181,
		stats.Spirit:    179,
		stats.Mana:      3268,
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
