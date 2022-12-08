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

func (dk *DpsDeathknight) oblitRunesAt(sim *core.Simulation) time.Duration {
	_, f, u := dk.NormalCurrentRunes()
	d := dk.CurrentDeathRunes()

	if f == 0 && u == 0 && d == 0 {
		timings := [3]time.Duration{dk.NormalSpentFrostRuneReadyAt(sim), dk.NormalSpentUnholyRuneReadyAt(sim), dk.SpentDeathRuneReadyAt()}
		if timings[0] > timings[2] {
			timings[0], timings[2] = timings[2], timings[0]
		}
		if timings[0] > timings[1] {
			timings[0], timings[1] = timings[1], timings[0]
		}
		if timings[1] > timings[2] {
			timings[1], timings[2] = timings[2], timings[1]
		}
		return timings[1]
	} else if f == 0 && u == 0 && d > 0 {
		// Next Frost/Unholy
		return core.MinDuration(dk.NormalSpentFrostRuneReadyAt(sim), dk.NormalSpentUnholyRuneReadyAt(sim))
	} else if f == 0 && u > 0 && d == 0 {
		// Next death rune or next f rune
		return core.MinDuration(dk.NormalSpentFrostRuneReadyAt(sim), dk.SpentDeathRuneReadyAt())
	} else if f > 0 && u == 0 && d == 0 {
		// Next death rune or next f rune
		return core.MinDuration(dk.NormalSpentUnholyRuneReadyAt(sim), dk.SpentDeathRuneReadyAt())
	}

	return sim.CurrentTime
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
		} else if dk.KM() && dk.FrostStrike.CanCast(sim) {
			dk.FrostStrike.Cast(sim, target)
		} else if dk.Rime() {
			dk.HowlingBlast.Cast(sim, target)
		} else if dk.FrostStrike.CanCast(sim) {
			dk.FrostStrike.Cast(sim, target)
		}

		s.ConditionalAdvance(casted && advance)
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
	waitFor := 100 * time.Millisecond

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

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Desync_FS_Dump(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if !dk.AllRunesSpent() {
		s.Advance()
		return sim.CurrentTime
	}

	if dk.KM() && dk.FrostStrike.CanCast(sim) {
		dk.FrostStrike.Cast(sim, target)
	} else if dk.Rime() {
		dk.HowlingBlast.Cast(sim, target)
	} else if dk.FrostStrike.CanCast(sim) {
		dk.FrostStrike.Cast(sim, target)
	}

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
