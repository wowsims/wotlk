package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) dndExperimentalOpener() {
	dk.Opener.Clear().
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.getBloodRuneAction(true)).
		NewAction(dk.RotationActionCallback_DND).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UP).
		NewAction(dk.RotationActionCallback_Garg).
		NewAction(dk.RotationAction_CancelBT)

	if dk.Rotation.UseEmpowerRuneWeapon {
		dk.Opener.
			NewAction(dk.RotationActionCallback_ERW)

		if dk.Rotation.ArmyOfTheDead == proto.Deathknight_Rotation_AsMajorCd {
			dk.Opener.
				NewAction(dk.RotationActionCallback_AOTD).
				NewAction(dk.RotationActionCallback_BP)
		} else {
			dk.Opener.
				NewAction(dk.RotationActionCallback_BP).
				NewAction(dk.getFirstDiseaseAction()).
				NewAction(dk.getSecondDiseaseAction()).
				NewAction(dk.RotationActionCallback_BS)
		}

		dk.Opener.
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_GF).
			NewAction(dk.RotationAction_DC_Custom).
			NewAction(dk.RotationAction_DC_Custom)
	} else {
		dk.Opener.
			NewAction(dk.getFirstDiseaseAction()).
			NewAction(dk.getSecondDiseaseAction()).
			NewAction(dk.RotationActionCallback_BP)
	}

	dk.dndExperimentalStartSequence()
}

func (dk *DpsDeathknight) dndExperimentalStartSequence() {
	dk.Main.Clear().
		NewAction(dk.RotationAction_FF_ClipCheck).
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.RotationAction_BP_ClipCheck).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_BloodRune_Custom).
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_Dnd_Custom).
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_Gargoyle_Custom2).
		NewAction(dk.RotationAction_UnholyDndRotationGhoulFrenzyCheck)
}

func (dk *DpsDeathknight) RotationAction_Gargoyle_Custom1(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	return dk.RotationAction_Gargoyle_Custom(dk.SpellGCD()+50*time.Millisecond, sim, target, s)
}

func (dk *DpsDeathknight) RotationAction_Gargoyle_Custom2(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	return dk.RotationAction_Gargoyle_Custom(dk.SpellGCD()*2+50*time.Millisecond, sim, target, s)
}

func (dk *DpsDeathknight) RotationAction_Gargoyle_Custom(castTime time.Duration, sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	if dk.uhGargoyleCanCast(sim, castTime) {
		if !dk.PresenceMatches(deathknight.UnholyPresence) {
			dk.CastBloodTap(sim, dk.CurrentTarget)
			dk.CastUnholyPresence(sim, dk.CurrentTarget)
		}
		if dk.CastSummonGargoyle(sim, target) {
			return true
		}
	}

	// Go back to Blood Presence after gargoyle cast
	if dk.PresenceMatches(deathknight.UnholyPresence) && !dk.CanSummonGargoyle(sim) {
		if dk.BloodTapAura.IsActive() {
			dk.BloodTapAura.Deactivate(sim)
		}
		if dk.CastBloodPresence(sim, target) {
			dk.WaitUntil(sim, sim.CurrentTime)
			return true
		}
	}

	s.Advance()
	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationAction_Dnd_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.CastDeathAndDecay(sim, target)
	if !casted {
		if !dk.DeathAndDecay.CD.IsReady(sim) {
			if dk.SummonGargoyle.IsReady(sim) {
				if dk.uhGargoyleCanCast(sim, dk.SpellGCD()+50*time.Millisecond) {
					if !dk.PresenceMatches(deathknight.UnholyPresence) {
						dk.CastBloodTap(sim, dk.CurrentTarget)
						dk.CastUnholyPresence(sim, dk.CurrentTarget)
					}
					if dk.CastSummonGargoyle(sim, target) {
						return true
					}
				} else {
					waitUntil := core.MinDuration(dk.DeathAndDecay.ReadyAt(), sim.CurrentTime+100*time.Millisecond)
					dk.WaitUntil(sim, waitUntil)
					return true
				}
			} else {
				dk.WaitUntil(sim, dk.DeathAndDecay.ReadyAt())
				return true
			}
		}
	} else {
		s.Advance()
	}
	return casted
}

func (dk *DpsDeathknight) RotationAction_DC_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	if !dk.uhDeathCoilCheck(sim) {
		dk.WaitUntil(sim, sim.CurrentTime)
	} else {
		casted := dk.CastDeathCoil(sim, target)
		if !casted {
			dk.WaitUntil(sim, sim.CurrentTime)
		}
	}
	s.Advance()
	return true
}

func (dk *DpsDeathknight) RotationAction_BloodRune_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	if dk.PresenceMatches(deathknight.UnholyPresence) {
		return dk.RotationActionCallback_BP(sim, target, s)
	} else if dk.desolationAuraCheck(sim) {
		return dk.RotationActionCallback_BS(sim, target, s)
	} else {
		return dk.RotationActionCallback_BB(sim, target, s)
	}
}

func (dk *DpsDeathknight) RotationAction_UnholyDndRotationGhoulFrenzyCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.Main.Clear()

	if dk.Talents.GhoulFrenzy && (!dk.GhoulFrenzyAura.IsActive() || dk.GhoulFrenzyAura.RemainingDuration(sim) < time.Second*10) {
		if dk.ur.ffFirst {
			dk.Main.NewAction(dk.RotationActionCallback_IT).
				NewAction(dk.RotationActionCallback_GF)
		} else {
			dk.Main.NewAction(dk.RotationActionCallback_GF).
				NewAction(dk.RotationActionCallback_IT)
		}
	} else {
		if dk.Talents.ScourgeStrike {
			dk.Main.NewAction(dk.RotationActionCallback_SS)
		} else {
			dk.Main.
				NewAction(dk.RotationAction_FF_ClipCheck).
				NewAction(dk.getFirstDiseaseAction()).
				NewAction(dk.RotationAction_BP_ClipCheck).
				NewAction(dk.getSecondDiseaseAction())
		}
	}

	dk.Main.
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_BloodRune_Custom)

	dk.Main.
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_Gargoyle_Custom2).
		NewAction(dk.RotationAction_UnholyDndRotationEnd)

	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationAction_UnholyDndRotationEnd(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.dndExperimentalStartSequence()
	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}
