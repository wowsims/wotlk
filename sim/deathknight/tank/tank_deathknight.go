package tank

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func RegisterTankDeathknight() {
	core.RegisterAgentFactory(
		proto.Player_TankDeathknight{},
		proto.Spec_SpecTankDeathknight,
		func(character core.Character, options proto.Player) core.Agent {
			return NewTankDeathknight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Deathknight)
			if !ok {
				panic("Invalid spec value for Deathknight!")
			}
			player.Spec = playerSpec
		},
	)
}

type TankDeathknight struct {
	*deathknight.Deathknight

	Options  proto.TankDeathknight_Options
	Rotation proto.TankDeathknight_Rotation
}

func NewTankDeathknight(character core.Character, options proto.Player) *TankDeathknight {
	dkOptions := options.GetTankDeathknight()

	tankDk := &TankDeathknight{
		Deathknight: deathknight.NewDeathknight(character, *dkOptions.Talents, deathknight.DeathknightInputs{
			StartingRunicPower: dkOptions.Options.StartingRunicPower,
		}),
		Rotation: *dkOptions.Rotation,
		Options:  *dkOptions.Options,
	}

	return tankDk
}

func (dk *TankDeathknight) GetDeathknight() *deathknight.Deathknight {
	return dk.Deathknight
}

func (dk *TankDeathknight) Initialize() {
	dk.Deathknight.Initialize()
}

func (dk *TankDeathknight) Reset(sim *core.Simulation) {
	dk.Deathknight.Reset(sim)
}
