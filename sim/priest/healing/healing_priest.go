package healing

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/priest"
)

func RegisterHealingPriest() {
	core.RegisterAgentFactory(
		proto.Player_HealingPriest{},
		proto.Spec_SpecHealingPriest,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewHealingPriest(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_HealingPriest)
			if !ok {
				panic("Invalid spec value for Healing Priest!")
			}
			player.Spec = playerSpec
		},
	)
}

type HealingPriest struct {
	*priest.Priest

	Options *proto.HealingPriest_Options
}

func NewHealingPriest(character *core.Character, options *proto.Player) *HealingPriest {
	healingOptions := options.GetHealingPriest()

	selfBuffs := priest.SelfBuffs{
		UseInnerFire:   healingOptions.Options.UseInnerFire,
		UseShadowfiend: healingOptions.Options.UseShadowfiend,
	}

	basePriest := priest.New(character, selfBuffs, options.TalentsString)
	hpriest := &HealingPriest{
		Priest:  basePriest,
		Options: healingOptions.Options,
	}

	hpriest.SelfBuffs.PowerInfusionTarget = &proto.UnitReference{}
	if hpriest.Talents.PowerInfusion && hpriest.Options.PowerInfusionTarget != nil {
		hpriest.SelfBuffs.PowerInfusionTarget = hpriest.Options.PowerInfusionTarget
	}

	return hpriest
}

func (hpriest *HealingPriest) GetPriest() *priest.Priest {
	return hpriest.Priest
}

func (hpriest *HealingPriest) GetMainTarget() *core.Unit {
	target := hpriest.Env.Raid.GetFirstTargetDummy()
	if target == nil {
		return &hpriest.Unit
	} else {
		return &target.Unit
	}
}

func (hpriest *HealingPriest) Initialize() {
	hpriest.CurrentTarget = hpriest.GetMainTarget()
	hpriest.Priest.Initialize()
	hpriest.Priest.RegisterHealingSpells()

	hpriest.ApplyRapture(hpriest.Options.RapturesPerMinute)
	hpriest.RegisterHymnOfHopeCD()
}

func (hpriest *HealingPriest) Reset(sim *core.Simulation) {
	hpriest.Priest.Reset(sim)
}
