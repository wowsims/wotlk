package restoration

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/druid"
)

func RegisterRestorationDruid() {
	core.RegisterAgentFactory(
		proto.Player_RestorationDruid{},
		proto.Spec_SpecRestorationDruid,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewRestorationDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_RestorationDruid)
			if !ok {
				panic("Invalid spec value for Restoration Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewRestorationDruid(character *core.Character, options *proto.Player) *RestorationDruid {
	restoOptions := options.GetRestorationDruid()
	selfBuffs := druid.SelfBuffs{}

	resto := &RestorationDruid{
		Druid: druid.New(character, druid.Tree, selfBuffs, options.TalentsString),
	}

	resto.SelfBuffs.InnervateTarget = &proto.UnitReference{}
	if restoOptions.Options.InnervateTarget != nil {
		resto.SelfBuffs.InnervateTarget = restoOptions.Options.InnervateTarget
	}

	return resto
}

type RestorationDruid struct {
	*druid.Druid
}

func (resto *RestorationDruid) GetDruid() *druid.Druid {
	return resto.Druid
}

func (resto *RestorationDruid) Initialize() {
	resto.Druid.Initialize()
}

func (resto *RestorationDruid) Reset(sim *core.Simulation) {
	resto.Druid.Reset(sim)
}
