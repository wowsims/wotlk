package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func RegisterDpsDeathknight() {
	core.RegisterAgentFactory(
		proto.Player_DeathKnight{},
		proto.Spec_SpecDeathKnight,
		func(character core.Character, options proto.Player) core.Agent {
			return NewDpsDeathknight(character, options)
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

type DpsDeathKnight struct {
	*deathknight.DeathKnight

	Rotation proto.DeathKnight_Rotation
}

func NewDpsDeathknight(character core.Character, player proto.Player) *DpsDeathKnight {
	dk := player.GetDeathKnight()

	dpsDk := &DpsDeathKnight{
		DeathKnight: deathknight.NewDeathKnight(character, player),
		Rotation:    *dk.Rotation,
	}

	dpsDk.DeathKnight.RefreshHornOfWinter = dk.Rotation.RefreshHornOfWinter
	dpsDk.DeathKnight.UnholyPresenceOpener = dk.Rotation.UnholyPresenceOpener
	dpsDk.DeathKnight.ArmyOfTheDeadType = dk.Rotation.ArmyOfTheDead

	dpsDk.SetupRotationEvent = dpsDk.SetupRotations
	dpsDk.DoRotationEvent = dpsDk.DoRotations

	return dpsDk
}

func (deathKnight *DpsDeathKnight) SetupRotations() deathknight.RotationID {
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
			if deathKnight.Rotation.ArmyOfTheDead == proto.DeathKnight_Rotation_AsMajorCd {
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

func (dk *DpsDeathKnight) DoRotations(sim *core.Simulation, target *core.Unit) {
	rotationId := dk.GetRotationId()

	if rotationId == deathknight.RotationID_FrostSubBlood_Full || rotationId == deathknight.RotationID_FrostSubUnholy_Full {
		dk.doFrostRotation(sim, target)
	} else {
		dk.doUnholyRotation(sim, target)
	}
}

func (dk *DpsDeathKnight) GetDeathKnight() *deathknight.DeathKnight {
	return dk.DeathKnight
}

func (dk *DpsDeathKnight) Initialize() {
	dk.DeathKnight.Initialize()
}

func (dk *DpsDeathKnight) Reset(sim *core.Simulation) {
	dk.DeathKnight.Reset(sim)
}
