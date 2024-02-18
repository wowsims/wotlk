package protection

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/paladin"
)

func RegisterProtectionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_ProtectionPaladin{},
		proto.Spec_SpecProtectionPaladin,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewProtectionPaladin(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ProtectionPaladin) // I don't really understand this line
			if !ok {
				panic("Invalid spec value for Protection Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewProtectionPaladin(character *core.Character, options *proto.Player) *ProtectionPaladin {
	protOptions := options.GetProtectionPaladin()

	prot := &ProtectionPaladin{
		Paladin: paladin.NewPaladin(character, options.TalentsString),
		Options: protOptions.Options,
		Seal:    protOptions.Options.Seal,
	}

	prot.PaladinAura = protOptions.Options.Aura

	prot.HasGlyphAS = prot.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfAvengerSShield)

	prot.EnableAutoAttacks(prot, core.AutoAttackOptions{
		MainHand:       prot.WeaponFromMainHand(0), // Set crit multiplier later when we have targets.
		AutoSwingMelee: true,
	})

	healingModel := options.HealingModel
	if healingModel != nil {
		if healingModel.InspirationUptime > 0.0 {
			core.ApplyInspiration(prot.GetCharacter(), healingModel.InspirationUptime)
		}
	}

	return prot
}

type ProtectionPaladin struct {
	*paladin.Paladin

	Options *proto.ProtectionPaladin_Options

	Judgement proto.PaladinJudgement

	Seal proto.PaladinSeal

	HasGlyphAS bool
}

func (prot *ProtectionPaladin) GetPaladin() *paladin.Paladin {
	return prot.Paladin
}

func (prot *ProtectionPaladin) Initialize() {
	prot.Paladin.Initialize()
	prot.ActivateRighteousFury()

	if prot.Options.UseAvengingWrath {
		prot.RegisterAvengingWrathCD()
	}
}

func (prot *ProtectionPaladin) Reset(sim *core.Simulation) {
	prot.Paladin.Reset(sim)

	switch prot.Seal {
	case proto.PaladinSeal_Vengeance:
		prot.CurrentSeal = prot.SealOfVengeanceAura
		prot.SealOfVengeanceAura.Activate(sim)
	case proto.PaladinSeal_Command:
		prot.CurrentSeal = prot.SealOfCommandAura
		prot.SealOfCommandAura.Activate(sim)
	case proto.PaladinSeal_Righteousness:
		prot.CurrentSeal = prot.SealOfRighteousnessAura
		prot.SealOfRighteousnessAura.Activate(sim)
	}

	prot.RighteousFuryAura.Activate(sim)
	prot.Paladin.PseudoStats.Stunned = false
}
