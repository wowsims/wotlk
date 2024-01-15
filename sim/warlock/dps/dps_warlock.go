package dps

import (
	"github.com/wowsims/sod/sim/common"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/warlock"
)

func RegisterDpsWarlock() {
	core.RegisterAgentFactory(
		proto.Player_Warlock{},
		proto.Spec_SpecWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewDpsWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Warlock)
			if !ok {
				panic("Invalid spec value for Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

type DpsWarlock struct {
	*warlock.Warlock

	Options        *proto.Warlock_Options
	Rotation       *proto.Warlock_Rotation
	CustomRotation *common.CustomRotation
}

func NewDpsWarlock(character *core.Character, options *proto.Player) *DpsWarlock {
	warlockOptions := options.GetWarlock()

	warlock := &DpsWarlock{
		Warlock:  warlock.NewWarlock(character, options),
		Rotation: warlockOptions.Rotation,
		Options:  warlockOptions.Options,
	}

	return warlock
}

func (warlock *DpsWarlock) OnGCDReady(sim *core.Simulation) {
	return
}

func (warlock *DpsWarlock) GetWarlock() *warlock.Warlock {
	return warlock.Warlock
}

func (warlock *DpsWarlock) Initialize() {
	warlock.Warlock.Initialize()
}

func (warlock *DpsWarlock) Reset(sim *core.Simulation) {
	warlock.Warlock.Reset(sim)
}
