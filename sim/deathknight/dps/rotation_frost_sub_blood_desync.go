package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

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

func (dk *DpsDeathknight) RotationActionCallback_FS_Special(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	fitsTwo := dk.CurrentRunicPower()/float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower()) > 2.0
	km := dk.KM()

	if km || fitsTwo {
		dk.FrostStrike.Cast(sim, target)

		if !fitsTwo {
			s.Advance()
		}
		return -1
	} else {
		// TODO: Use the grace period of those runes.
		hwExtra := core.TernaryDuration(dk.HornOfWinter.CanCast(sim), dk.SpellGCD(), 0)
		ob := dk.oblitRunesAt(sim)
		fishingPeriod := ob - 1500*time.Millisecond - hwExtra
		nextSwingAt := dk.NextMHSwingAt(sim) + 1

		if nextSwingAt <= fishingPeriod {
			return nextSwingAt
		} else {
			dk.FrostStrike.Cast(sim, target)
			s.Advance()
			return -1
		}
	}
}

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
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FS_Special).
		NewAction(dk.RotationActionCallback_HW).
		NewAction(dk.RotationActionCallback_FrostSubBlood_DesyncRotation)
}

func (dk *DpsDeathknight) canCastInDesyncWindow(sim *core.Simulation, spell *deathknight.RuneSpell) bool {
	if !dk.RuneIsDeath(1) {
		return true
	}

	gcd := dk.GetGcdDuration(spell)
	u := dk.UnholyRuneReadyAt(sim)
	f := dk.FrostRuneReadyAt(sim)
	d := dk.RuneReadyAt(sim, 1)

	if f <= sim.CurrentTime || u <= sim.CurrentTime {
		return true
	}

	if !(d <= f && f < u) {
		return true
	}

	if f+gcd >= u {
		return false
	}

	return true
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_DesyncRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	t := sim.CurrentTime
	ff := dk.FrostFeverDisease[target.Index].ExpiresAt() - t
	bp := dk.BloodPlagueDisease[target.Index].ExpiresAt() - t
	ob := dk.oblitRunesAt(sim)
	spGcd := dk.SpellGCD()
	abGcd := 1500 * time.Millisecond
	b, f, u := dk.NormalCurrentRunes()
	d := dk.CurrentDeathRunes()
	km := dk.KM()
	rime := dk.Rime()

	if b+f+u+int32(d) <= 1 && sim.IsExecutePhase35() && dk.EmpowerRuneWeapon.CanCast(sim) && dk.Rotation.UseEmpowerRuneWeapon &&
		sim.GetRemainingDuration() <= time.Duration(2.0*float64(dk.UnbreakableArmorAura.Duration)) {
		dk.EmpowerRuneWeapon.Cast(sim, target)
		dk.castAllMajorCooldowns(sim)
		return t
	}

	if dk.RotationActionCallback_LastSecondsCast(sim, target) {
		return -1
	}

	if dk.RuneIsDeath(0) && dk.RuneIsDeath(1) && dk.LeftBloodRuneReady() && dk.RightBloodRuneReady() {
		dk.Pestilence.Cast(sim, target)
		return -1
	}

	if ff <= 0 {
		dk.IcyTouch.Cast(sim, target)
		return -1
	}

	if bp <= 0 {
		dk.PlagueStrike.Cast(sim, target)
		return -1
	}

	if ff <= 2*time.Second || bp < 2*time.Second {
		dk.Pestilence.Cast(sim, target)
		return -1
	}

	if dk.UnbreakableArmor.CanCast(sim) && dk.BloodTap.CanCast(sim) {

		if b > 0 {
			dk.UnbreakableArmor.Cast(sim, target)
			dk.castAllMajorCooldowns(sim)
			dk.BloodTap.Cast(sim, target)
		} else if d == 2 {
			dk.UnbreakableArmor.Cast(sim, target)
			dk.castAllMajorCooldowns(sim)
			dk.BloodTap.Cast(sim, target)
		} else if b == 0 && d == 0 {
			dk.BloodTap.Cast(sim, target)
			dk.UnbreakableArmor.Cast(sim, target)
			dk.castAllMajorCooldowns(sim)
		}

	}

	if km && dk.FrostStrike.CanCast(sim) && dk.shDiseaseCheck(sim, target, dk.FrostStrike, false, 1, 0) && dk.canCastInDesyncWindow(sim, dk.FrostStrike) {
		dk.FrostStrike.Cast(sim, target)
		return -1
	}

	if ((f > 0 && u > 0) || (f == 0 && u > 0 && d > 0) || (f > 0 && u == 0 && d > 0)) && dk.shDiseaseCheck(sim, target, dk.Obliterate, true, 1, 0) && dk.canCastInDesyncWindow(sim, dk.FrostStrike) {
		dk.Obliterate.Cast(sim, target)
		return -1
	}

	if t+abGcd <= ob && dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 100.0 && dk.canCastInDesyncWindow(sim, dk.FrostStrike) {
		dk.FrostStrike.Cast(sim, target)
		return -1
	}

	if t+spGcd <= ob && rime && dk.HowlingBlast.CanCast(sim) && dk.CurrentRunicPower() <= dk.MaxRunicPower()-5.0 && dk.canCastInDesyncWindow(sim, dk.FrostStrike) {
		dk.HowlingBlast.Cast(sim, target)
		return -1
	}

	if t+abGcd <= ob && dk.FrostStrike.CanCast(sim) && dk.canCastInDesyncWindow(sim, dk.FrostStrike) {
		dk.FrostStrike.Cast(sim, target)
		return -1
	}

	if t+spGcd <= ob && dk.HornOfWinter.CanCast(sim) && dk.CurrentRunicPower()+10.0 <= dk.MaxRunicPower() && dk.canCastInDesyncWindow(sim, dk.FrostStrike) {
		dk.HornOfWinter.Cast(sim, target)
		return -1
	}

	if dk.LeftBloodRuneReady() && !dk.RuneIsDeath(0) {
		dk.Pestilence.Cast(sim, target)
		return -1
	}

	return -1
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Sequence_Pesti_Desync(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	waitUntil := time.Duration(-1)

	ff := dk.FrostFeverDisease[target.Index].ExpiresAt() - sim.CurrentTime
	bp := dk.BloodPlagueDisease[target.Index].ExpiresAt() - sim.CurrentTime

	if dk.RotationActionCallback_LastSecondsCast(sim, target) {
		return -1
	}

	if ff <= 0 || bp <= 0 {
		return dk.RotationActionCallback_FrostSubBlood_DesyncRotation(sim, target, s)
	} else {
		casted = dk.Pestilence.Cast(sim, target)
		advance := dk.LastOutcome.Matches(core.OutcomeLanded)
		if !casted || (casted && !dk.LastOutcome.Matches(core.OutcomeLanded)) {

			if dk.SpellGCD() > ff || dk.SpellGCD() > bp {
				return dk.RotationActionCallback_FrostSubBlood_DesyncRotation(sim, target, s)
			} else {
				s.ConditionalAdvance(casted && advance)
				return core.TernaryDuration(casted, -1, waitUntil)
			}
		} else {
			s.ConditionalAdvance(casted && advance)
			return core.TernaryDuration(casted, -1, waitUntil)
		}
	}
}

func (dk *DpsDeathknight) RotationActionCallback_LastSecondsCast_ERW(sim *core.Simulation, target *core.Unit) bool {
	casted := false

	t := sim.CurrentTime
	ff := dk.FrostFeverDisease[target.Index].ExpiresAt() - t
	bp := dk.BloodPlagueDisease[target.Index].ExpiresAt() - t

	km := dk.KM()
	if core.MinDuration(ff, bp) > sim.GetRemainingDuration() {
		if dk.Obliterate.CanCast(sim) && ff > 0 && bp > 0 {
			casted = dk.Obliterate.Cast(sim, target)
		} else if dk.FrostStrike.CanCast(sim) && km {
			casted = dk.FrostStrike.Cast(sim, target)
		} else if dk.FrostStrike.CanCast(sim) {
			casted = dk.FrostStrike.Cast(sim, target)
		} else if dk.Obliterate.CanCast(sim) {
			casted = dk.Obliterate.Cast(sim, target)
		} else if dk.HowlingBlast.CanCast(sim) {
			casted = dk.HowlingBlast.Cast(sim, target)
		} else if dk.HornOfWinter.CanCast(sim) {
			casted = dk.HornOfWinter.Cast(sim, target)
		}
	}

	return casted
}
