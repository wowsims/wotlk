package shadow

import (
	"github.com/wowsims/wotlk/sim/common/wotlk"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/priest"
)

func RegisterShadowPriest() {
	core.RegisterAgentFactory(
		proto.Player_ShadowPriest{},
		proto.Spec_SpecShadowPriest,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewShadowPriest(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ShadowPriest)
			if !ok {
				panic("Invalid spec value for Shadow Priest!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewShadowPriest(character *core.Character, options *proto.Player) *ShadowPriest {
	shadowOptions := options.GetShadowPriest()

	selfBuffs := priest.SelfBuffs{
		UseShadowfiend: true,
		UseInnerFire:   shadowOptions.Options.Armor == proto.ShadowPriest_Options_InnerFire,
	}

	basePriest := priest.New(character, selfBuffs, options.TalentsString)
	basePriest.Latency = float64(basePriest.ChannelClipDelay.Milliseconds())
	spriest := &ShadowPriest{
		Priest:  basePriest,
		options: shadowOptions.Options,
	}

	spriest.SelfBuffs.PowerInfusionTarget = &proto.UnitReference{}
	if spriest.Talents.PowerInfusion && shadowOptions.Options.PowerInfusionTarget != nil {
		spriest.SelfBuffs.PowerInfusionTarget = shadowOptions.Options.PowerInfusionTarget
	}

	wotlk.ConstructValkyrPets(&spriest.Character)

	return spriest
}

type ShadowPriest struct {
	*priest.Priest
	options *proto.ShadowPriest_Options
}

func (spriest *ShadowPriest) GetPriest() *priest.Priest {
	return spriest.Priest
}

func (spriest *ShadowPriest) Initialize() {
	spriest.Priest.Initialize()
}

func (spriest *ShadowPriest) Reset(sim *core.Simulation) {
	spriest.Priest.Reset(sim)
}
