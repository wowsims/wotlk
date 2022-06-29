package protection

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/paladin"
)

// Do 1 less millisecond to solve for sim order of operation problems
// Buffs are removed before melee swing is processed
const twistWindow = 399 * time.Millisecond

func RegisterProtectionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_ProtectionPaladin{},
		proto.Spec_SpecProtectionPaladin,
		func(character core.Character, options proto.Player) core.Agent {
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

func NewProtectionPaladin(character core.Character, options proto.Player) *ProtectionPaladin {
	protOptions := options.GetProtectionPaladin()

	prot := &ProtectionPaladin{
		Paladin:  paladin.NewPaladin(character, *protOptions.Talents),
		Rotation: *protOptions.Rotation,
		Options:  *protOptions.Options,
	}
	prot.PaladinAura = protOptions.Options.Aura

	prot.EnableAutoAttacks(prot, core.AutoAttackOptions{
		MainHand:       prot.WeaponFromMainHand(0), // Set crit multiplier later when we have targets.
		AutoSwingMelee: true,
	})

	return prot
}

type ProtectionPaladin struct {
	*paladin.Paladin

	Rotation proto.ProtectionPaladin_Rotation
	Options  proto.ProtectionPaladin_Options
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

	switch prot.Rotation.ConsecrationRank {
	case 6:
		prot.RegisterConsecrationSpell(6)
	case 4:
		prot.RegisterConsecrationSpell(4)
	case 1:
		prot.RegisterConsecrationSpell(1)
	}
}

func (prot *ProtectionPaladin) Reset(sim *core.Simulation) {
	prot.Paladin.Reset(sim)

	// Pre-activate seal before combat starts.
	if prot.Rotation.MaintainJudgement == proto.PaladinJudgement_JudgementOfWisdom {
		prot.UpdateSeal(sim, prot.SealOfWisdomAura)
	} else if prot.Rotation.MaintainJudgement == proto.PaladinJudgement_JudgementOfLight {
		prot.UpdateSeal(sim, prot.SealOfLightAura)
	} else {
		prot.UpdateSeal(sim, prot.SealOfRighteousnessAura)
	}

	// Pre-activate Holy Shield before combat starts.
	// Assume it gets cast 3s before entering combat.
	prot.HolyShieldAura.Activate(sim)
	prot.HolyShield.CD.Timer.Set(time.Second * 7)
}
