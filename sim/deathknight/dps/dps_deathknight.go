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

			RefreshHornOfWinter:  dk.Rotation.RefreshHornOfWinter,
			UnholyPresenceOpener: dk.Rotation.UnholyPresenceOpener,
			ArmyOfTheDeadType:    dk.Rotation.ArmyOfTheDead,
		}),
		Rotation: *dk.Rotation,
	}

	dpsDk.SetupRotationEvent = dpsDk.SetupRotations
	dpsDk.DoRotationEvent = dpsDk.DoRotations

	return dpsDk
}

func (deathKnight *DpsDeathknight) SetupRotations() deathknight.RotationID {
	deathKnight.setupFrostRotations()
	deathKnight.setupUnholyRotations()

	// IMPORTANT
	rotationId := deathknight.RotationID_Unknown
	// Also you need to update this to however you define spec
	if deathKnight.Talents.DarkConviction > 0 && deathKnight.Talents.HowlingBlast {
		rotationId = deathknight.RotationID_FrostSubBlood_Full
	} else if deathKnight.Talents.BloodCakedBlade > 0 && deathKnight.Talents.HowlingBlast {
		rotationId = deathknight.RotationID_FrostSubUnholy_Full
	} else if deathKnight.Talents.HowlingBlast {
		rotationId = deathknight.RotationID_FrostSubBlood_Full
	} else if deathKnight.Talents.SummonGargoyle {
		if deathKnight.Rotation.UseDeathAndDecay {
			rotationId = deathknight.RotationID_UnholyDnd_Full
		} else {
			if deathKnight.Rotation.ArmyOfTheDead == proto.Deathknight_Rotation_AsMajorCd {
				if deathKnight.Rotation.UnholyPresenceOpener {
					rotationId = deathknight.RotationID_UnholySsArmyUnholyPresence_Full
				} else {
					rotationId = deathknight.RotationID_UnholySsArmyBloodPresence_Full
				}
			} else {
				if deathKnight.Rotation.UnholyPresenceOpener {
					rotationId = deathknight.RotationID_UnholySsUnholyPresence_Full
				} else {
					rotationId = deathknight.RotationID_UnholySsBloodPresence_Full
				}
			}
		}
	} else {
		rotationId = deathknight.RotationID_Default
	}

	return rotationId
}

func (dk *DpsDeathknight) DoRotations(sim *core.Simulation, target *core.Unit) {
	rotationId := dk.GetRotationId()

	if rotationId == deathknight.RotationID_FrostSubBlood_Full || rotationId == deathknight.RotationID_FrostSubUnholy_Full {
		dk.doFrostRotation(sim, target)
	} else {
		dk.doUnholyRotation(sim, target)
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
}
