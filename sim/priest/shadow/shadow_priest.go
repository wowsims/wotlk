package shadow

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/priest"
)

func RegisterShadowPriest() {
	core.RegisterAgentFactory(
		proto.Player_ShadowPriest{},
		proto.Spec_SpecShadowPriest,
		func(character core.Character, options proto.Player) core.Agent {
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

func NewShadowPriest(character core.Character, options proto.Player) *ShadowPriest {
	shadowOptions := options.GetShadowPriest()

	selfBuffs := priest.SelfBuffs{
		UseShadowfiend: shadowOptions.Options.UseShadowfiend,
		UseInnerFire:   shadowOptions.Options.Armor == proto.ShadowPriest_Options_InnerFire,
	}

	basePriest := priest.New(character, selfBuffs, *shadowOptions.Talents)
	basePriest.Latency = shadowOptions.Rotation.Latency
	spriest := &ShadowPriest{
		Priest:   basePriest,
		rotation: *shadowOptions.Rotation,
	}

	spriest.ApplyShadowOnHitEffects()

	spriest.EnableResumeAfterManaWait(spriest.tryUseGCD)

	return spriest
}

type ShadowPriest struct {
	DPstatH  float64
	DPstatpH float64
	DPstatSp float64

	VTstatH  float64
	VTstatpH float64
	VTstatSp float64

	*priest.Priest
	rotation proto.ShadowPriest_Rotation
}

func (spriest *ShadowPriest) GetPriest() *priest.Priest {
	return spriest.Priest
}

func (spriest *ShadowPriest) Reset(sim *core.Simulation) {
	spriest.Priest.Reset(sim)
}
