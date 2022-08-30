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

	sr SharedRotation
	fr FrostRotation
	ur UnholyRotation

	Rotation proto.Deathknight_Rotation
}

func NewDpsDeathknight(character core.Character, player proto.Player) *DpsDeathknight {
	dk := player.GetDeathknight()

	dpsDk := &DpsDeathknight{
		Deathknight: deathknight.NewDeathknight(character, *dk.Talents, deathknight.DeathknightInputs{
			StartingRunicPower:  dk.Options.StartingRunicPower,
			PrecastGhoulFrenzy:  dk.Options.PrecastGhoulFrenzy,
			PrecastHornOfWinter: dk.Options.PrecastHornOfWinter,
			PetUptime:           dk.Options.PetUptime,
			IsDps:               true,

			RefreshHornOfWinter: dk.Rotation.RefreshHornOfWinter,
			ArmyOfTheDeadType:   dk.Rotation.ArmyOfTheDead,
			StartingPresence:    dk.Rotation.StartingPresence,
			UseAMS:              dk.Rotation.UseAms,
			AvgAMSSuccessRate:   dk.Rotation.AvgAmsSuccessRate,
			AvgAMSHit:           dk.Rotation.AvgAmsHit,
		}),
		Rotation: *dk.Rotation,
	}

	dpsDk.EnableAutoAttacks(dpsDk, core.AutoAttackOptions{
		MainHand:       dpsDk.WeaponFromMainHand(dpsDk.DefaultMeleeCritMultiplier()),
		OffHand:        dpsDk.WeaponFromOffHand(dpsDk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		ReplaceMHSwing: func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if dpsDk.RuneStrike.CanCast(sim) {
				return dpsDk.RuneStrike.Spell
			} else {
				return nil
			}
		},
	})

	dpsDk.sr.dk = dpsDk
	dpsDk.ur.dk = dpsDk

	return dpsDk
}

func (dk *DpsDeathknight) FrostPointsInBlood() int32 {
	return dk.Talents.Butchery + dk.Talents.Subversion + dk.Talents.BladeBarrier + dk.Talents.DarkConviction
}

func (dk *DpsDeathknight) FrostPointsInUnholy() int32 {
	return dk.Talents.ViciousStrikes + dk.Talents.Virulence + dk.Talents.Epidemic + dk.Talents.RavenousDead + dk.Talents.Necrosis + dk.Talents.BloodCakedBlade
}

func (dk *DpsDeathknight) SetupRotations() {
	dk.ur.ffFirst = dk.Rotation.FirstDisease == proto.Deathknight_Rotation_FrostFever
	dk.ur.hasGod = dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDisease)

	dk.RotationSequence.Clear()

	if dk.Talents.HowlingBlast && (dk.FrostPointsInBlood() > dk.FrostPointsInUnholy()) {
		if dk.Rotation.UseEmpowerRuneWeapon {
			dk.setupFrostSubBloodERWOpener()
		} else {
			dk.setupFrostSubBloodNoERWOpener()
		}
	} else if dk.Talents.HowlingBlast && (dk.FrostPointsInBlood() <= dk.FrostPointsInUnholy()) {
		if dk.Rotation.UseEmpowerRuneWeapon {
			dk.setupFrostSubUnholyERWOpener()
		} else {
			dk.setupFrostSubUnholyERWOpener()
		}
	} else if dk.Talents.SummonGargoyle {
		dk.setupUnholyRotations()
	} else if dk.Talents.DancingRuneWeapon {
		dk.setupBloodRotations()
	} else {
		// TODO: Add some default rotation that works without special talents
		dk.RotationSequence.Clear().
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_PS).
			NewAction(dk.RotationActionCallback_BS)
	}
}

func (dk *DpsDeathknight) GetDeathknight() *deathknight.Deathknight {
	return dk.Deathknight
}

func (dk *DpsDeathknight) Initialize() {
	dk.Deathknight.Initialize()
	dk.initProcTrackers()
	dk.fr.Initialize(dk)
}

func (dk *DpsDeathknight) Reset(sim *core.Simulation) {
	dk.Deathknight.Reset(sim)

	dk.Presence = deathknight.UnsetPresence

	if dk.Inputs.StartingPresence == proto.Deathknight_Rotation_Unholy && dk.Talents.SummonGargoyle {
		dk.ChangePresence(sim, deathknight.UnholyPresence)
	} else {
		dk.ChangePresence(sim, deathknight.BloodPresence)
	}

	dk.sr.Reset(sim)
	dk.fr.Reset(sim)
	dk.ur.Reset(sim)

	dk.SetupRotations()
}
