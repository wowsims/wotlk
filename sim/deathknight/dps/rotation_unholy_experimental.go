package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) dndExperimentalOpener() {
	dk.RotationSequence.Clear().
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.getBloodRuneAction(true)).
		NewAction(dk.RotationActionCallback_DND).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UP).
		NewAction(dk.RotationActionCallback_Garg).
		NewAction(dk.RotationAction_CancelBT)

	if dk.Rotation.UseEmpowerRuneWeapon {
		dk.RotationSequence.
			NewAction(dk.RotationActionCallback_ERW)

		if dk.Inputs.ArmyOfTheDeadType == proto.Deathknight_Rotation_AsMajorCd {
			dk.RotationSequence.
				NewAction(dk.RotationActionCallback_AOTD).
				NewAction(dk.RotationActionCallback_BP)
		} else {
			dk.RotationSequence.
				NewAction(dk.RotationActionCallback_BP).
				NewAction(dk.getFirstDiseaseAction()).
				NewAction(dk.getSecondDiseaseAction()).
				NewAction(dk.RotationActionCallback_BS)
		}

		dk.RotationSequence.
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_GF).
			NewAction(dk.RotationAction_DC_Custom).
			NewAction(dk.RotationAction_DC_Custom)
	} else {
		dk.RotationSequence.
			NewAction(dk.getFirstDiseaseAction()).
			NewAction(dk.getSecondDiseaseAction()).
			NewAction(dk.RotationActionCallback_BP)
	}

	dk.dndExperimentalStartSequence()
}

func (dk *DpsDeathknight) dndExperimentalStartSequence() {
	dk.RotationSequence.Clear().
		NewAction(dk.RotationActionUH_FF_ClipCheck).
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.RotationActionUH_BP_ClipCheck).
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

func (dk *DpsDeathknight) RotationAction_Gargoyle_Custom1(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	return dk.RotationAction_Gargoyle_Custom(dk.SpellGCD()+50*time.Millisecond, sim, target, s)
}

func (dk *DpsDeathknight) RotationAction_Gargoyle_Custom2(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	return dk.RotationAction_Gargoyle_Custom(dk.SpellGCD()*2+50*time.Millisecond, sim, target, s)
}

func (dk *DpsDeathknight) RotationAction_Gargoyle_Custom(castTime time.Duration, sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.uhGargoyleCanCast(sim, castTime) {
		if !dk.PresenceMatches(deathknight.UnholyPresence) {
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
			dk.UnholyPresence.Cast(sim, dk.CurrentTarget)
		}
		if dk.SummonGargoyle.Cast(sim, target) {
			return sim.CurrentTime
		}
	}

	// Go back to Blood Presence after gargoyle cast
	if dk.PresenceMatches(deathknight.UnholyPresence) && !dk.SummonGargoyle.IsReady(sim) {
		if dk.BloodTapAura.IsActive() {
			dk.BloodTapAura.Deactivate(sim)
		}
		if dk.BloodPresence.Cast(sim, target) {
			return sim.CurrentTime
		}
	}

	s.Advance()
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationAction_Dnd_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := dk.DeathAndDecay.Cast(sim, target)
	if !casted {
		if !dk.DeathAndDecay.CD.IsReady(sim) {
			if dk.SummonGargoyle.IsReady(sim) {
				if dk.uhGargoyleCanCast(sim, dk.SpellGCD()+50*time.Millisecond) {
					if !dk.PresenceMatches(deathknight.UnholyPresence) {
						dk.BloodTap.Cast(sim, dk.CurrentTarget)
						dk.UnholyPresence.Cast(sim, dk.CurrentTarget)
					}
					if dk.SummonGargoyle.Cast(sim, target) {
						return sim.CurrentTime
					}
				} else {
					return core.MinDuration(dk.DeathAndDecay.ReadyAt(), sim.CurrentTime+100*time.Millisecond)
				}
			} else {
				return dk.DeathAndDecay.ReadyAt()
			}
		}
	} else {
		s.Advance()
	}
	return -1
}

func (dk *DpsDeathknight) RotationAction_DC_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	waitUntil := time.Duration(-1)
	if !dk.uhDeathCoilCheck(sim) {
		waitUntil = sim.CurrentTime
	} else {
		casted := dk.DeathCoil.Cast(sim, target)
		if !casted {
			waitUntil = sim.CurrentTime
		}
	}
	s.Advance()
	return waitUntil
}

func (dk *DpsDeathknight) RotationAction_BloodRune_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.PresenceMatches(deathknight.UnholyPresence) {
		return dk.RotationActionCallback_BP(sim, target, s)
	} else if dk.desolationAuraCheck(sim) {
		return dk.RotationActionCallback_BS(sim, target, s)
	} else {
		return dk.RotationActionCallback_BB(sim, target, s)
	}
}

func (dk *DpsDeathknight) RotationAction_UnholyDndRotationGhoulFrenzyCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dk.RotationSequence.Clear()

	if dk.Talents.GhoulFrenzy && (!dk.GhoulFrenzyAura.IsActive() || dk.GhoulFrenzyAura.RemainingDuration(sim) < time.Second*10) {
		if dk.sr.ffFirst {
			dk.RotationSequence.NewAction(dk.RotationActionCallback_IT).
				NewAction(dk.RotationActionCallback_GF)
		} else {
			dk.RotationSequence.NewAction(dk.RotationActionCallback_GF).
				NewAction(dk.RotationActionCallback_IT)
		}
	} else {
		if dk.Talents.ScourgeStrike {
			dk.RotationSequence.NewAction(dk.RotationActionCallback_SS)
		} else {
			dk.RotationSequence.
				NewAction(dk.RotationActionUH_FF_ClipCheck).
				NewAction(dk.getFirstDiseaseAction()).
				NewAction(dk.RotationActionUH_BP_ClipCheck).
				NewAction(dk.getSecondDiseaseAction())
		}
	}

	dk.RotationSequence.
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_BloodRune_Custom)

	dk.RotationSequence.
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_Gargoyle_Custom1).
		NewAction(dk.RotationAction_DC_Custom).
		NewAction(dk.RotationAction_Gargoyle_Custom2).
		NewAction(dk.RotationAction_UnholyDndRotationEnd)

	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationAction_UnholyDndRotationEnd(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dk.dndExperimentalStartSequence()
	return sim.CurrentTime
}
