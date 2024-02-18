package smite

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/priest"
)

func RegisterSmitePriest() {
	core.RegisterAgentFactory(
		proto.Player_SmitePriest{},
		proto.Spec_SpecSmitePriest,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewSmitePriest(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_SmitePriest)
			if !ok {
				panic("Invalid spec value for Smite Priest!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewSmitePriest(character *core.Character, options *proto.Player) *SmitePriest {
	smiteOptions := options.GetSmitePriest()

	selfBuffs := priest.SelfBuffs{
		UseInnerFire:   smiteOptions.Options.UseInnerFire,
		UseShadowfiend: smiteOptions.Options.UseShadowfiend,
	}

	basePriest := priest.New(character, selfBuffs, options.TalentsString)
	spriest := &SmitePriest{
		Priest: basePriest,
	}

	spriest.SelfBuffs.PowerInfusionTarget = &proto.UnitReference{}
	if spriest.Talents.PowerInfusion && smiteOptions.Options.PowerInfusionTarget != nil {
		spriest.SelfBuffs.PowerInfusionTarget = smiteOptions.Options.PowerInfusionTarget
	}

	return spriest
}

type SmitePriest struct {
	*priest.Priest
}

func (spriest *SmitePriest) GetPriest() *priest.Priest {
	return spriest.Priest
}

func (spriest *SmitePriest) Initialize() {
	spriest.Priest.Initialize()

	spriest.RegisterHolyFireSpell()
	spriest.RegisterSmiteSpell()
	spriest.RegisterPenanceSpell()
	spriest.RegisterHymnOfHopeCD()
}

func (spriest *SmitePriest) Reset(sim *core.Simulation) {
	spriest.Priest.Reset(sim)
}
