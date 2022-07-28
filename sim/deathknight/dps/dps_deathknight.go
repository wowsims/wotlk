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
			FirstDisease:         dk.Rotation.FirstDisease,
		}),
		Rotation: *dk.Rotation,
	}

	dpsDk.SetupRotations()

	return dpsDk
}

func (dk *DpsDeathknight) SetupRotations() {
	if dk.Talents.DarkConviction > 0 && dk.Talents.HowlingBlast {
		dk.setupFrostSubBloodOpener()
	} else if dk.Talents.BloodCakedBlade > 0 && dk.Talents.HowlingBlast {
		dk.setupFrostSubUnholyOpener()
	} else if dk.Talents.HowlingBlast {
		dk.setupFrostSubBloodOpener()
	} else if dk.Talents.SummonGargoyle {
		if dk.Rotation.UseDeathAndDecay {
			if dk.Rotation.UnholyPresenceOpener {
				dk.setupUnholyDndUnholyPresenceOpener()
			} else {
				dk.setupUnholyDndBloodPresenceOpener()
			}
		} else {

			if dk.Rotation.ArmyOfTheDead == proto.Deathknight_Rotation_AsMajorCd {
				if dk.Rotation.UnholyPresenceOpener {
					dk.setupUnholySsArmyUnholyPresenceOpener()
				} else {
					dk.setupUnholySsArmyBloodPresenceOpener()
				}
			} else {
				if dk.Rotation.UnholyPresenceOpener {
					dk.setupUnholySsUnholyPresenceOpener()
				} else {
					dk.setupUnholySsBloodPresenceOpener()
				}
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
}
