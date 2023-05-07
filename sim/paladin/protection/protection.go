package protection

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/paladin"
)

func RegisterProtectionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_ProtectionPaladin{},
		proto.Spec_SpecProtectionPaladin,
		func(character core.Character, options *proto.Player) core.Agent {
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

func NewProtectionPaladin(character core.Character, options *proto.Player) *ProtectionPaladin {
	protOptions := options.GetProtectionPaladin()

	prot := &ProtectionPaladin{
		Paladin:  paladin.NewPaladin(character, options.TalentsString),
		Rotation: protOptions.Rotation,
		Options:  protOptions.Options,
		Seal:     protOptions.Options.Seal,
	}

	var rotationInput = protOptions.Rotation.CustomRotation

	if rotationInput != nil {
		prot.RotationInput = make([]int32, len(rotationInput.Spells))
		for i, customSpellProto := range rotationInput.Spells {
			prot.RotationInput[i] = customSpellProto.Spell
		}
	}

	prot.SelectedRotation = prot.customRotation

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

	Rotation *proto.ProtectionPaladin_Rotation
	Options  *proto.ProtectionPaladin_Options

	Judgement proto.PaladinJudgement

	Seal proto.PaladinSeal

	HasGlyphAS bool

	SelectedRotation func(*core.Simulation)
	RotationInput    []int32
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

	// Pre-activate Holy Shield before combat starts.
	// Assume it gets cast 3s before entering combat.
	prot.HolyShieldAura.Activate(sim)
	prot.HolyShield.CD.Timer.Set(time.Second * 7)

	sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int) {
		if isExecute == 20 {
			prot.OnGCDReady(sim)
		}
	})

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

	prot.DivinePleaAura.Activate(sim)
	prot.DivinePlea.CD.Use(sim)
	prot.RighteousFuryAura.Activate(sim)
}
