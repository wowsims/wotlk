package dps

import (
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
}

func NewDpsWarlock(character *core.Character, options *proto.Player) *DpsWarlock {
	warlock := &DpsWarlock{
		Warlock: warlock.NewWarlock(character, options, options.GetWarlock().Options),
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
