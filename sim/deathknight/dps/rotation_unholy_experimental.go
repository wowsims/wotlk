package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) dndStartOpener() {
	// Static opener with no Proc checks for gargoyle
	dk.Opener.
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.getBloodRuneAction(true)).
		NewAction(dk.RotationActionCallback_DND).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UP).
		NewAction(dk.RotationActionCallback_Garg).
		NewAction(dk.RotationAction_CancelBT).
		NewAction(dk.RotationActionCallback_ERW)

	if dk.Rotation.ArmyOfTheDead == proto.Deathknight_Rotation_AsMajorCd {
		dk.Opener.
			NewAction(dk.RotationActionCallback_AOTD).
			NewAction(dk.RotationActionCallback_BP)
	} else {
		dk.Opener.
			NewAction(dk.RotationActionCallback_BP).
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_PS).
			NewAction(dk.RotationActionCallback_BS)
	}

	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_GF).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_DC_Custom)

	// Experimental rotation with sequences
	dk.dndStartSequence()
}

func (dk *DpsDeathknight) dndStartSequence() {
	dk.Main.Clear().
		NewAction(dk.RotationAction_FF_ClipCheck).
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.RotationAction_BP_ClipCheck).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.getBloodRuneAction(true)).
		NewAction(dk.RotationAction_Dnd_Custom).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_UnholyDndRotationGhoulFrenzyCheck)
}

// Custom Dnd callback with delay
func (dk *DpsDeathknight) RotationAction_Dnd_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.CastDeathAndDecay(sim, target)
	if !casted {
		if !dk.DeathAndDecay.CD.IsReady(sim) {
			dk.WaitUntil(sim, dk.DeathAndDecay.ReadyAt())
			return true
		}
	} else {
		s.Advance()
	}
	return casted
}

func (dk *DpsDeathknight) RotationAction_DC_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.CastDeathCoil(sim, target)
	if !casted {
		dk.WaitUntil(sim, sim.CurrentTime)
	}
	s.Advance()
	return true
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

	if dk.desolationAuraCheck(sim) {
		dk.Main.NewAction(dk.RotationActionCallback_BS)
	} else {
		dk.Main.NewAction(dk.RotationActionCallback_BB)
	}
	dk.Main.
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_UnholyDndRotationEnd)

	if dk.uhGargoyleCanCast(sim) {
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

	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationAction_UnholyDndRotationEnd(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.dndStartSequence()
	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}
