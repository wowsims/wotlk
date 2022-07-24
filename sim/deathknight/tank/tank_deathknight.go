package tank

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func RegisterTankDeathknight() {
	core.RegisterAgentFactory(
		proto.Player_DeathKnightTank{},
		proto.Spec_SpecDeathKnightTank,
		func(character core.Character, options proto.Player) core.Agent {
			return NewTankDeathknight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DeathKnight)
			if !ok {
				panic("Invalid spec value for Deathknight!")
			}
			player.Spec = playerSpec
		},
	)
}

type TankDeathknight struct {
	*deathknight.DeathKnight

	Options  proto.DeathKnightTank_Options
	Rotation proto.DeathKnightTank_Rotation
}

func NewTankDeathknight(character core.Character, options proto.Player) *TankDeathknight {
	dkOptions := options.GetDeathKnightTank()

	dk := &TankDeathknight{
		DeathKnight: deathknight.NewDeathKnight(character, options),
		Rotation:    *dkOptions.Rotation,
		Options:     *dkOptions.Options,
	}

	return dk
}

func (dk *TankDeathknight) GetDeathKnight() *deathknight.DeathKnight {
	return dk.DeathKnight
}

func (dk *TankDeathknight) Initialize() {
	dk.DeathKnight.Initialize()
}

func (dk *TankDeathknight) Reset(sim *core.Simulation) {
	dk.DeathKnight.Reset(sim)
}
