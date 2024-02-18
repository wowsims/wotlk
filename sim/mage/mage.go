package mage

import (
	"github.com/wowsims/wotlk/sim/common/wotlk"

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
		func(character *core.Character, options *proto.Player) core.Agent {
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

	Talents *proto.MageTalents
	Options *proto.Mage_Options

	waterElemental *WaterElemental
	mirrorImage    *MirrorImage

	// Cached values for a few mechanics.
	bonusCritDamage float64

	ArcaneBarrage           *core.Spell
	ArcaneBlast             *core.Spell
	ArcaneExplosion         *core.Spell
	ArcaneMissiles          *core.Spell
	ArcaneMissilesTickSpell *core.Spell
	Blizzard                *core.Spell
	DeepFreeze              *core.Spell
	Ignite                  *core.Spell
	LivingBomb              *core.Spell
	Fireball                *core.Spell
	FireBlast               *core.Spell
	Flamestrike             *core.Spell
	FlamestrikeRank8        *core.Spell
	Frostbolt               *core.Spell
	FrostfireBolt           *core.Spell
	IceLance                *core.Spell
	Pyroblast               *core.Spell
	Scorch                  *core.Spell
	MirrorImage             *core.Spell
	BlastWave               *core.Spell
	DragonsBreath           *core.Spell

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
	mage.registerFlamestrikeSpells()
	mage.registerFrostboltSpell()
	mage.registerIceLanceSpell()
	mage.registerPyroblastSpell()
	mage.registerScorchSpell()
	mage.registerLivingBombSpell()
	mage.registerFrostfireBoltSpell()
	mage.registerEvocation()
	mage.registerManaGemsCD()
	mage.registerMirrorImageCD()
	mage.registerBlastWaveSpell()
	mage.registerDragonsBreathSpell()
}

func (mage *Mage) Reset(sim *core.Simulation) {
}

func NewMage(character *core.Character, options *proto.Player) *Mage {
	mageOptions := options.GetMage()

	mage := &Mage{
		Character: *character,
		Talents:   &proto.MageTalents{},
		Options:   mageOptions.Options,
	}
	core.FillTalentsProto(mage.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	mage.bonusCritDamage = .25*float64(mage.Talents.SpellPower) + .1*float64(mage.Talents.Burnout)
	mage.EnableManaBar()

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
		mage.waterElemental = mage.NewWaterElemental(mage.Options.WaterElementalDisobeyChance)
	}

	wotlk.ConstructValkyrPets(&mage.Character)
	return mage
}

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}
