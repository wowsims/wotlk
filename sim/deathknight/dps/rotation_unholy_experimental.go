package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) dndStartSequence() {
	dk.Main.Clear().NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.getBloodRuneAction(true)).
		NewAction(dk.RotationAction_Dnd_Custom).
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
			dk.Main.NewAction(dk.getFirstDiseaseAction()).
				NewAction(dk.getSecondDiseaseAction())
		}
	}

	if dk.desolationAuraCheck(sim) {
		dk.Main.NewAction(dk.RotationActionCallback_BS)
	} else {
		dk.Main.NewAction(dk.RotationActionCallback_BB)
	}
	dk.Main.NewAction(dk.RotationAction_UnholyDndRotationEnd)

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
