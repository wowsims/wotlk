package shadow

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/priest"
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
	basePriest := priest.New(character, options.TalentsString)
	basePriest.Latency = float64(basePriest.ChannelClipDelay.Milliseconds())
	spriest := &ShadowPriest{
		Priest:  basePriest,
		options: shadowOptions.Options,
	}

	return spriest
}

type ShadowPriest struct {
	PrevTicks float64

	*priest.Priest
	options *proto.ShadowPriest_Options

	AllCDs   []time.Duration
	BLUsedAt time.Duration
}

func (spriest *ShadowPriest) GetPriest() *priest.Priest {
	return spriest.Priest
}

func (spriest *ShadowPriest) Initialize() {
	spriest.Priest.Initialize()
}

func (spriest *ShadowPriest) Reset(sim *core.Simulation) {
	spriest.Priest.Reset(sim)

	// Save info related to blood lust timing
	spriest.BLUsedAt = 0
	if bloodlustMCD := spriest.GetMajorCooldownIgnoreTag(core.BloodlustActionID); bloodlustMCD != nil {
		timings := bloodlustMCD.GetTimings()
		if len(timings) > 0 {
			spriest.BLUsedAt = timings[0]
		}
	}
}
