package shadow

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/priest"
)

func RegisterShadowPriest() {
	core.RegisterAgentFactory(
		proto.Player_ShadowPriest{},
		proto.Spec_SpecShadowPriest,
		func(character core.Character, options *proto.Player) core.Agent {
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

func NewShadowPriest(character core.Character, options *proto.Player) *ShadowPriest {
	shadowOptions := options.GetShadowPriest()

	selfBuffs := priest.SelfBuffs{
		UseShadowfiend: shadowOptions.Options.UseShadowfiend,
		UseInnerFire:   shadowOptions.Options.Armor == proto.ShadowPriest_Options_InnerFire,
	}

	basePriest := priest.New(character, selfBuffs, options.TalentsString)
	basePriest.Latency = shadowOptions.Rotation.Latency
	spriest := &ShadowPriest{
		Priest:   basePriest,
		rotation: shadowOptions.Rotation,
		options:  shadowOptions.Options,
	}

	spriest.SelfBuffs.PowerInfusionTarget = &proto.RaidTarget{TargetIndex: -1}
	if spriest.Talents.PowerInfusion && shadowOptions.Options.PowerInfusionTarget != nil {
		spriest.SelfBuffs.PowerInfusionTarget = shadowOptions.Options.PowerInfusionTarget
	}

	spriest.EnableResumeAfterManaWait(spriest.tryUseGCD)
	spriest.CanRolloverSWP = spriest.Talents.MindFlay && spriest.Talents.PainAndSuffering > 0

	return spriest
}

type ShadowPriest struct {
	PrevTicks float64

	*priest.Priest
	rotation *proto.ShadowPriest_Rotation
	options  *proto.ShadowPriest_Options

	VTCastTime time.Duration
	AllCDs     []time.Duration
	BLUsedAt   time.Duration

	CanRolloverSWP bool
}

func (spriest *ShadowPriest) GetPriest() *priest.Priest {
	return spriest.Priest
}

func (spriest *ShadowPriest) Initialize() {
	spriest.Priest.Initialize()

	if spriest.rotation.PrecastType > 0 {
		spriest.RegisterPrepullAction(-1500*time.Millisecond, func(sim *core.Simulation) {
			if spriest.rotation.PrecastType == 1 {
				spriest.VampiricTouch.Cast(sim, spriest.CurrentTarget)
			}
			if spriest.rotation.PrecastType == 2 {
				spriest.MindBlast.Cast(sim, spriest.CurrentTarget)
			}
		})
	}
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
