package mage

import (
	"github.com/wowsims/classic/sod/sim/core"
	"github.com/wowsims/classic/sod/sim/core/proto"
)

const (
	SpellFlagMage       = core.SpellFlagAgentReserved1
	SpellFlagChillSpell = core.SpellFlagAgentReserved2
)

var TalentTreeSizes = [3]int{16, 16, 17}

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

	Talents  *proto.MageTalents
	Options  *proto.Mage_Options
	Rotation *proto.Mage_Rotation

	ArcaneBlast             *core.Spell
	ArcaneExplosion         *core.Spell
	ArcaneMissiles          *core.Spell
	ArcaneMissilesTickSpell *core.Spell
	BlastWave               *core.Spell
	Blizzard                *core.Spell
	Ignite                  *core.Spell
	LivingBomb              *core.Spell
	LivingFlame             *core.Spell
	Fireball                *core.Spell
	FireBlast               *core.Spell
	Flamestrike             *core.Spell
	Frostbolt               *core.Spell
	IceLance                *core.Spell
	Pyroblast               *core.Spell
	Scorch                  *core.Spell

	IcyVeins *core.Spell

	ArcaneBlastAura    *core.Aura
	ArcanePotencyAura  *core.Aura
	ArcanePowerAura    *core.Aura
	ClearcastingAura   *core.Aura
	ScorchAuras        core.AuraArray
	CombustionAura     *core.Aura
	FingersOfFrostAura *core.Aura
	EnlightenmentAura  *core.Aura

	CritDebuffCategories core.ExclusiveCategoryArray
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
	mage.registerArcaneMissilesSpell()
	mage.registerFireballSpell()
	mage.registerFireBlastSpell()
	mage.registerFrostboltSpell()
	// mage.registerManaGemsCD()
	// mage.registerPyroblastSpell()
	mage.registerScorchSpell()

	// TODO: Classic mage aoe spells
	// mage.registerArcaneExplosionSpell()
	// mage.registerBlizzardSpell()
	// mage.registerFlamestrikeSpells()
	// mage.registerBlastWaveSpell()
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

	mage.EnableManaBar()

	if mage.Options.Armor == proto.Mage_Options_MageArmor {
		mage.PseudoStats.SpiritRegenRateCasting += .3
	}

	return mage
}

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}
