package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func RegisterDpsDeathknight() {
	core.RegisterAgentFactory(
		proto.Player_Deathknight{},
		proto.Spec_SpecDeathknight,
		func(character core.Character, options proto.Player) core.Agent {
			return NewDpsDeathknight(character, options)
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

type DpsDeathknight struct {
	*deathknight.Deathknight

	UnholyRotation

	Rotation proto.Deathknight_Rotation
}

func NewDpsDeathknight(character core.Character, player proto.Player) *DpsDeathknight {
	dk := player.GetDeathknight()

	dpsDk := &DpsDeathknight{
		Deathknight: deathknight.NewDeathknight(character, player, deathknight.DeathknightInputs{
			StartingRunicPower:  dk.Options.StartingRunicPower,
			PrecastGhoulFrenzy:  dk.Options.PrecastGhoulFrenzy,
			PrecastHornOfWinter: dk.Options.PrecastHornOfWinter,
			PetUptime:           dk.Options.PetUptime,

			RefreshHornOfWinter: dk.Rotation.RefreshHornOfWinter,
			ArmyOfTheDeadType:   dk.Rotation.ArmyOfTheDead,
			FirstDisease:        dk.Rotation.FirstDisease,
		}),
		Rotation: *dk.Rotation,
	}

	return dpsDk
}

func (dk *DpsDeathknight) SetupRotations() {
	dk.ffFirst = dk.Inputs.FirstDisease == proto.Deathknight_Rotation_FrostFever

	dk.Opener.Clear()
	dk.Main.Clear()

	if dk.Talents.DarkConviction > 0 && dk.Talents.HowlingBlast {
		dk.setupFrostSubBloodOpener()
	} else if dk.Talents.BloodCakedBlade > 0 && dk.Talents.HowlingBlast {
		dk.setupFrostSubUnholyOpener()
	} else if dk.Talents.HowlingBlast {
		dk.setupFrostSubBloodOpener()
	} else if dk.Talents.SummonGargoyle {
		if dk.Rotation.UseDeathAndDecay {
			dk.setupUnholyDndOpener()
		} else {

			if dk.Rotation.ArmyOfTheDead == proto.Deathknight_Rotation_AsMajorCd {
				dk.setupUnholySsArmyOpener()
			} else {
				dk.setupUnholySsOpener()
			}
		}
	} else {
		// TODO: Add some default rotation that works without special talents
		dk.setupFrostSubBloodOpener()
	}
}

func (dk *DpsDeathknight) GetDeathknight() *deathknight.Deathknight {
	return dk.Deathknight
}

func (dk *DpsDeathknight) Initialize() {
	dk.Deathknight.Initialize()
}

func (dk *DpsDeathknight) Reset(sim *core.Simulation) {
	dk.Deathknight.Reset(sim)
	dk.SetupRotations()
	dk.ResetUnholyRotation(sim)
}
