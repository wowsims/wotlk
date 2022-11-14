package dps

import (
	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func RegisterDpsDeathknight() {
	core.RegisterAgentFactory(
		proto.Player_Deathknight{},
		proto.Spec_SpecDeathknight,
		func(character core.Character, options *proto.Player) core.Agent {
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

	CustomRotation *common.CustomRotation

	Rotation *proto.Deathknight_Rotation
}

func NewDpsDeathknight(character core.Character, player *proto.Player) *DpsDeathknight {
	dk := player.GetDeathknight()

	dpsDk := &DpsDeathknight{
		Deathknight: deathknight.NewDeathknight(character, dk.Talents, deathknight.DeathknightInputs{
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
		Rotation: dk.Rotation,
	}

	dpsDk.Inputs.UnholyFrenzyTarget = &proto.RaidTarget{TargetIndex: -1}
	if dk.Options.UnholyFrenzyTarget != nil {
		dpsDk.Inputs.UnholyFrenzyTarget = dk.Options.UnholyFrenzyTarget
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
	if dk.Rotation.AutoRotation {
		bl, fr, uh := deathknight.PointsInTalents(dk.Talents)

		if uh > fr && uh > bl {
			// Unholy
			dk.Rotation.BtGhoulFrenzy = true
			dk.Rotation.UseEmpowerRuneWeapon = true
			dk.Rotation.HoldErwArmy = false
			dk.Rotation.UseGargoyle = true
			dk.Rotation.ArmyOfTheDead = proto.Deathknight_Rotation_AsMajorCd
			dk.Rotation.BloodTap = proto.Deathknight_Rotation_GhoulFrenzy
			dk.Rotation.FirstDisease = proto.Deathknight_Rotation_FrostFever
			dk.Rotation.StartingPresence = proto.Deathknight_Rotation_Unholy
			dk.Rotation.BlPresence = proto.Deathknight_Rotation_Blood
			dk.Rotation.Presence = proto.Deathknight_Rotation_Blood

			mh := dk.GetMHWeapon()
			oh := dk.GetOHWeapon()

			if mh != nil && oh != nil {
				// DW
				dk.Rotation.BloodRuneFiller = proto.Deathknight_Rotation_BloodBoil
				dk.Rotation.UseDeathAndDecay = true
			} else {
				// 2h
				if dk.Env.GetNumTargets() > 1 {
					dk.Rotation.BloodRuneFiller = proto.Deathknight_Rotation_BloodBoil
					dk.Rotation.UseDeathAndDecay = true
				} else {
					dk.Rotation.BloodRuneFiller = proto.Deathknight_Rotation_BloodStrike
					dk.Rotation.UseDeathAndDecay = false
				}
			}
			// Always use DnD if you have the glyph.
			if dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDeathAndDecay) {
				dk.Rotation.UseDeathAndDecay = true
			}
		} else if fr > uh && fr > bl {
			// Frost rotations here.
		} else if bl > fr && bl > uh {
			// Blood rotations here.
		} else {
			// some weird spec where two trees are equal...
		}
	}
	dk.ur.ffFirst = dk.Rotation.FirstDisease == proto.Deathknight_Rotation_FrostFever
	dk.ur.hasGod = dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDisease)

	dk.RotationSequence.Clear()

	dk.Inputs.FuStrike = deathknight.FuStrike_Obliterate

	dk.CustomRotation = dk.makeCustomRotation()
	if dk.CustomRotation == nil || dk.Rotation.FrostRotationType == proto.Deathknight_Rotation_SingleTarget {
		dk.Rotation.FrostRotationType = proto.Deathknight_Rotation_SingleTarget
		if dk.Talents.HowlingBlast && (dk.FrostPointsInBlood() > dk.FrostPointsInUnholy()) {
			if dk.Rotation.UseEmpowerRuneWeapon {
				if dk.Rotation.DesyncRotation {
					dk.setupFrostSubBloodDesyncERWOpener()
				} else {
					dk.setupFrostSubBloodERWOpener()
				}
			} else {
				dk.setupFrostSubBloodNoERWOpener()
			}
		} else if dk.Talents.HowlingBlast && (dk.FrostPointsInBlood() <= dk.FrostPointsInUnholy()) {
			dk.Rotation.FrostRotationType = proto.Deathknight_Rotation_SingleTarget
			if dk.Rotation.UseEmpowerRuneWeapon {
				dk.setupFrostSubUnholyERWOpener()
			} else {
				dk.setupFrostSubUnholyERWOpener()
			}
		} else if dk.Talents.SummonGargoyle {
			dk.setupUnholyRotations()
		} else if dk.Talents.DancingRuneWeapon {
			dk.setupBloodRotations()
		}
	} else {
		dk.setupCustomRotations()
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

	dk.sr.Reset(sim)
	dk.fr.Reset(sim)
	dk.ur.Reset(sim)

	dk.SetupRotations()

	dk.Presence = deathknight.UnsetPresence

	b, f, u := deathknight.PointsInTalents(dk.Talents)

	if f > u && f > b {
		if dk.Rotation.Presence == proto.Deathknight_Rotation_Blood {
			dk.ChangePresence(sim, deathknight.BloodPresence)
		} else if dk.Rotation.Presence == proto.Deathknight_Rotation_Frost {
			dk.ChangePresence(sim, deathknight.FrostPresence)
		} else if dk.Rotation.Presence == proto.Deathknight_Rotation_Unholy {
			dk.ChangePresence(sim, deathknight.UnholyPresence)
		}
	}

	if u > f && u > b {
		if dk.Rotation.StartingPresence == proto.Deathknight_Rotation_Unholy {
			dk.ChangePresence(sim, deathknight.UnholyPresence)
		} else if dk.Talents.SummonGargoyle {
			dk.ChangePresence(sim, deathknight.BloodPresence)
		}
	}
}
