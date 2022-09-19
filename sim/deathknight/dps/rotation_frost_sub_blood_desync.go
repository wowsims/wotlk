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
		NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti_Desync).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_HW).
		NewAction(dk.RotationActionCallback_FrostSubBlood_DesyncRotation)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_DesyncRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	t := sim.CurrentTime
	ff := dk.FrostFeverDisease[target.Index].ExpiresAt() - t
	bp := dk.BloodPlagueDisease[target.Index].ExpiresAt() - t
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

	if ff <= 0 {
		dk.IcyTouch.Cast(sim, target)
		return -1
	}

	if bp <= 0 {
		dk.PlagueStrike.Cast(sim, target)
		return -1
	}

	if ff < 2*time.Second || bp < 2*time.Second {
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

	if km && dk.FrostStrike.CanCast(sim) && dk.shDiseaseCheck(sim, target, dk.FrostStrike, false, 1, 0) {
		dk.FrostStrike.Cast(sim, target)
		return -1
	}

	if f > 0 && u > 0 || ((f == 0 && u > 0 && d > 0) || (f > 0 && u == 0 && d > 0)) && dk.shDiseaseCheck(sim, target, dk.Obliterate, true, 1, 0) {
		dk.Obliterate.Cast(sim, target)
		return -1
	}

	if dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 100.0 {
		dk.FrostStrike.Cast(sim, target)
		return -1
	}

	if rime && dk.HowlingBlast.CanCast(sim) && dk.CurrentRunicPower() <= dk.MaxRunicPower()-5.0 {
		dk.HowlingBlast.Cast(sim, target)
		return -1
	}

	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())
	if dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 2.0*(fsCost-dk.fr.oblitRPRegen) {
		dk.FrostStrike.Cast(sim, target)
		return -1
	}

	if dk.HornOfWinter.CanCast(sim) {
		dk.HornOfWinter.Cast(sim, target)
		return -1
	}

	if dk.LeftBloodRuneReady() {
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
