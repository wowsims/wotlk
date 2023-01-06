package restoration

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/shaman"
)

func RegisterRestorationShaman() {
	core.RegisterAgentFactory(
		proto.Player_RestorationShaman{},
		proto.Spec_SpecRestorationShaman,
		func(character core.Character, options *proto.Player) core.Agent {
			return NewRestorationShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_RestorationShaman)
			if !ok {
				panic("Invalid spec value for Restoration Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewRestorationShaman(character core.Character, options *proto.Player) *RestorationShaman {
	restoShamOptions := options.GetRestorationShaman()

	selfBuffs := shaman.SelfBuffs{
		Bloodlust: restoShamOptions.Options.Bloodlust,
		Shield:    restoShamOptions.Options.Shield,
	}

	totems := &proto.ShamanTotems{}
	if restoShamOptions.Rotation.Totems != nil {
		totems = restoShamOptions.Rotation.Totems
		totems.UseFireMcd = true // Control fire totems as MCD.
	}

	resto := &RestorationShaman{
		Shaman: shaman.NewShaman(character, restoShamOptions.Talents, totems, selfBuffs, false),
	}
	resto.EnableResumeAfterManaWait(resto.tryUseGCD)

	return resto
}

type RestorationShaman struct {
	*shaman.Shaman
}

func (resto *RestorationShaman) GetShaman() *shaman.Shaman {
	return resto.Shaman
}

func (resto *RestorationShaman) Reset(sim *core.Simulation) {
	resto.Shaman.Reset(sim)
}
