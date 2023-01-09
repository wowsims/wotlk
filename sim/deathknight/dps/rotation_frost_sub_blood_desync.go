package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) setupFrostSubBloodDesyncERWOpener() {
	dk.setupUnbreakableArmorCooldowns()

	dk.RotationSequence.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_UA_Frost).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Obli).
		NewAction(dk.RotationAction_CancelBT).
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_Frost_FS_HB).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_Sequence1)
}

// We need to keep the B2 and F1 runes in sync and immediately use them for obliterate
// otherwise if an unholy rune comes up then we can't continue the Desync rotation without
// re-casting IT + PS
func (dk *DpsDeathknight) firstOblitAt(sim *core.Simulation) time.Duration {
	return core.MaxDuration(dk.RuneReadyAt(sim, 1), dk.RuneReadyAt(sim, 2))
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Desync_Obli(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	advance := true

	ff := dk.FrostFeverDisease[target.Index].IsActive()
	bp := dk.BloodPlagueDisease[target.Index].IsActive()

	if ff && bp {
		if dk.Obliterate.CanCast(sim) {
			if dk.Deathchill != nil && dk.Deathchill.IsReady(sim) {
				dk.Deathchill.Cast(sim, target)
			}
			casted = dk.Obliterate.Cast(sim, target)
			advance = dk.LastOutcome.Matches(core.OutcomeLanded)
			s.ConditionalAdvance(casted && advance)
		} else {
			dk.desync_Filler(sim, target)
		}
	} else if !ff {
		casted = dk.IcyTouch.Cast(sim, target)
		advance = dk.LastOutcome.Matches(core.OutcomeLanded)
		s.ConditionalAdvance(casted && advance)
	} else {
		casted = dk.PlagueStrike.Cast(sim, target)
		advance = dk.LastOutcome.Matches(core.OutcomeLanded)
		s.ConditionalAdvance(casted && advance)
	}

	return -1
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Desync_UA(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	runeGrace := dk.RuneGraceAt(0, sim.CurrentTime)
	waitFor := 5 * time.Millisecond

	if dk.UnbreakableArmor.IsReady(sim) && dk.BloodTap.IsReady(sim) {
		dk.BloodTap.Cast(sim, target)
		return sim.CurrentTime + waitFor
	} else if dk.UnbreakableArmor.IsReady(sim) && runeGrace >= waitFor {
		dk.UnbreakableArmor.Cast(sim, target)
	}

	s.Advance()
	return -1
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Desync_Sequence1(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	s.Clear().
		// f1 u1
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_Obli).
		// f2 u2
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_Obli).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_FS_Dump).
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_Sequence2)
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Desync_Sequence2(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	s.Clear().
		// d2 f1
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_Obli).
		// f2 u1
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_Obli).
		// u2 d1
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_Obli).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_UA).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_FS_Dump).
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_ERW).
		NewAction(dk.RotationAction_CancelBT).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Desync_Sequence1)
	return sim.CurrentTime
}

func (dk *DpsDeathknight) desync_Filler(sim *core.Simulation, target *core.Unit) {
	spell := dk.RegularPrioPickSpell(sim, target, dk.firstOblitAt(sim))
	if spell != nil {
		spell.Cast(sim, target)
	}
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Desync_FS_Dump(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if !dk.AllRunesSpent() {
		s.Advance()
		return sim.CurrentTime
	}

	dk.desync_Filler(sim, target)

	return -1
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Desync_ERW(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	bothDeath := dk.RuneIsDeath(0) && dk.RuneIsDeath(1)

	if sim.IsExecutePhase35() && dk.UnbreakableArmorAura.IsActive() && dk.Rotation.UseEmpowerRuneWeapon && bothDeath {
		dk.castAllMajorCooldowns(sim)

		// go to normal rotation!
		s.Clear().
			NewAction(dk.RotationActionCallback_ERW).
			NewAction(dk.RotationActionCallback_Obli).
			NewAction(dk.RotationActionCallback_Obli).
			NewAction(dk.RotationActionCallback_Obli).
			NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
	} else {
		s.Advance()
	}
	return sim.CurrentTime
}
