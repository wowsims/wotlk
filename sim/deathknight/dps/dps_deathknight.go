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

type RotationID uint8

const (
	RotationID_Default RotationID = iota
	RotationID_FrostSubBlood_Full
	RotationID_FrostSubUnholy_Full

	RotationID_UnholySsUnholyPresence_Full
	RotationID_UnholySsArmyUnholyPresence_Full
	RotationID_UnholySsBloodPresence_Full
	RotationID_UnholySsArmyBloodPresence_Full
	RotationID_UnholyDnd_Full
	RotationID_Unknown
	RotationID_Count
)

type DpsDeathknight struct {
	*deathknight.Deathknight

	Rotation        proto.Deathknight_Rotation
	CurrentRotation RotationID
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

	dpsDk.SetupRotations()
	dpsDk.DoRotationEvent = dpsDk.DoRotations

	return dpsDk
}

func (dk *DpsDeathknight) SetupRotations() {
	// IMPORTANT
	rotationId := RotationID_Unknown
	// Also you need to update this to however you define spec
	if dk.Talents.DarkConviction > 0 && dk.Talents.HowlingBlast {
		rotationId = RotationID_FrostSubBlood_Full
		dk.setupFrostRotations(RotationID_FrostSubBlood_Full)
	} else if dk.Talents.BloodCakedBlade > 0 && dk.Talents.HowlingBlast {
		rotationId = RotationID_FrostSubUnholy_Full
		dk.setupFrostRotations(RotationID_FrostSubUnholy_Full)
	} else if dk.Talents.HowlingBlast {
		rotationId = RotationID_FrostSubBlood_Full
		dk.setupFrostRotations(RotationID_FrostSubBlood_Full)
	} else if dk.Talents.SummonGargoyle {
		if dk.Rotation.UseDeathAndDecay {
			rotationId = RotationID_UnholyDnd_Full
			dk.setupUnholyRotations(RotationID_UnholyDnd_Full)
		} else {
			if dk.Rotation.ArmyOfTheDead == proto.Deathknight_Rotation_AsMajorCd {
				if dk.Rotation.UnholyPresenceOpener {
					rotationId = RotationID_UnholySsArmyUnholyPresence_Full
					dk.setupUnholyRotations(RotationID_UnholySsArmyUnholyPresence_Full)
				} else {
					rotationId = RotationID_UnholySsArmyBloodPresence_Full
					dk.setupUnholyRotations(RotationID_UnholySsArmyBloodPresence_Full)
				}
			} else {
				if dk.Rotation.UnholyPresenceOpener {
					rotationId = RotationID_UnholySsUnholyPresence_Full
					dk.setupUnholyRotations(RotationID_UnholySsUnholyPresence_Full)
				} else {
					rotationId = RotationID_UnholySsBloodPresence_Full
					dk.setupUnholyRotations(RotationID_UnholySsBloodPresence_Full)
				}
			}
		}
	} else {
		rotationId = RotationID_Default
	}

	dk.CurrentRotation = rotationId
}

func (dk *DpsDeathknight) DoRotations(sim *core.Simulation, target *core.Unit) {
	if dk.CurrentRotation == RotationID_FrostSubBlood_Full || dk.CurrentRotation == RotationID_FrostSubUnholy_Full {
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
