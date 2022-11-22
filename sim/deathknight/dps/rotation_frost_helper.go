package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

// end of fight oblit does not check diseases, it just presses it regardless, but will retry if fails to land.
func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnh_EndOfFight_Obli(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	advance := true
	waitTime := time.Duration(-1)
	if dk.Obliterate.CanCast(sim) {
		if dk.Deathchill != nil && dk.Deathchill.IsReady(sim) {
			dk.Deathchill.Cast(sim, target)
		}
		casted = dk.Obliterate.Cast(sim, target)
		advance = dk.LastOutcome.Matches(core.OutcomeLanded)
	}
	s.ConditionalAdvance(casted && advance)
	return core.TernaryDuration(casted, -1, waitTime)
}

func (dk *DpsDeathknight) RegularPrioPickSpell(sim *core.Simulation, target *core.Unit, untilTime time.Duration) *deathknight.RuneSpell {
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	abGcd := 1500 * time.Millisecond
	spGcd := dk.SpellGCD()

	km := dk.KM()
	rime := dk.Rime()
	if sim.CurrentTime+abGcd <= untilTime && dk.FrostStrike.CanCast(sim) && km {
		return dk.FrostStrike
	} else if sim.CurrentTime+abGcd <= untilTime && dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 100.0 {
		return dk.FrostStrike
	} else if sim.CurrentTime+spGcd <= untilTime && dk.HowlingBlast.CanCast(sim) && rime {
		return dk.HowlingBlast
	} else if sim.CurrentTime+abGcd <= untilTime && dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 2.0*(fsCost-dk.fr.oblitRPRegen) {
		return dk.FrostStrike
	} else if sim.CurrentTime+spGcd <= untilTime && dk.HornOfWinter.CanCast(sim) {
		return dk.HornOfWinter
	} else {
		return nil
	}
}

// more end of fight functions
func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_EndOfFight_1and2_GCD(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	frAt := dk.NormalFrostRuneReadyAt(sim)
	uhAt := dk.NormalUnholyRuneReadyAt(sim)
	obAt := core.MaxDuration(frAt, uhAt)
	//ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
	//bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	if obAt < sim.CurrentTime {
		s.Clear().NewAction(dk.RotationActionCallback_FrostSubUnh_EndOfFight_Obli)
	} else if dk.CurrentRunicPower() >= fsCost {
		s.Clear().NewAction(dk.RotationActionCallback_FS)
	} else if dk.Rime() {
		s.Clear().NewAction(dk.RotationActionCallback_HB)
	} else if dk.CurrentBloodRunes() >= 1 {
		s.Clear().NewAction(dk.RotationActionCallback_BS)
	} else if dk.CurrentFrostRunes() >= 1 {
		s.Clear().NewAction(dk.RotationActionCallback_IT)
	} else if dk.CurrentUnholyRunes() >= 1 {
		s.Clear().NewAction(dk.RotationActionCallback_PS)
	} else {
		s.Advance()
	}
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_EndOfFight_3GCD(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	frAt := dk.NormalFrostRuneReadyAt(sim)
	uhAt := dk.NormalUnholyRuneReadyAt(sim)
	obAt := core.MaxDuration(frAt, uhAt)
	ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
	bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	if obAt < sim.CurrentTime {
		s.Clear().NewAction(dk.RotationActionCallback_FrostSubUnh_EndOfFight_Obli)
	} else if dk.CurrentRunicPower() >= fsCost {
		s.Clear().NewAction(dk.RotationActionCallback_FS)
	} else if dk.Rime() {
		s.Clear().NewAction(dk.RotationActionCallback_HB)
	} else if ffExpiresAt < sim.GetMaxDuration() && bpExpiresAt < sim.GetMaxDuration() {
		s.Clear().NewAction(dk.RotationActionCallback_Pesti)
	} else if dk.CurrentBloodRunes() >= 1 {
		s.Clear().NewAction(dk.RotationActionCallback_BS)
	} else if dk.CurrentFrostRunes() >= 1 {
		s.Clear().NewAction(dk.RotationActionCallback_IT)
	} else if dk.CurrentUnholyRunes() >= 1 {
		s.Clear().NewAction(dk.RotationActionCallback_PS)
	} else {
		s.Advance()
	}
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_EndOfFight_4and5_GCD(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	frAt := dk.NormalFrostRuneReadyAt(sim)
	uhAt := dk.NormalUnholyRuneReadyAt(sim)
	obAt := core.MaxDuration(frAt, uhAt)
	ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
	bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	if obAt < sim.CurrentTime {
		s.Clear().NewAction(dk.RotationActionCallback_FrostSubUnh_EndOfFight_Obli)
	} else if dk.CurrentRunicPower() >= fsCost {
		s.Clear().NewAction(dk.RotationActionCallback_FS)
	} else if dk.Rime() {
		s.Clear().NewAction(dk.RotationActionCallback_HB)
	} else if ffExpiresAt < sim.GetMaxDuration() && bpExpiresAt < sim.GetMaxDuration() {
		s.Clear().NewAction(dk.RotationActionCallback_Pesti)
	} else if dk.CurrentBloodRunes() >= 1 {
		s.Clear().NewAction(dk.RotationActionCallback_BS)
	} else {
		s.Advance()
	}
	return sim.CurrentTime
}

// end of fight logic in the works
func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_EndOfFight(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if sim.Log != nil && sim.GetRemainingDuration() <= 1500*time.Millisecond*5 {
		sim.Log("endoffight function entered!")

		if sim.GetRemainingDuration() <= 1500*time.Millisecond*1 {
			s.Clear().NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_1and2_GCD)
		} else if sim.GetRemainingDuration() <= 1500*time.Millisecond*2 {
			s.Clear().
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_1and2_GCD).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_1and2_GCD)
		} else if sim.GetRemainingDuration() <= 1500*time.Millisecond*3 {
			s.Clear().
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_3GCD).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_1and2_GCD).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_1and2_GCD)
		} else if sim.GetRemainingDuration() <= 1500*time.Millisecond*4 {
			s.Clear().
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_4and5_GCD).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_3GCD).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_1and2_GCD).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_1and2_GCD)
		} else if sim.GetRemainingDuration() <= 1500*time.Millisecond*5 {
			s.Clear().
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_4and5_GCD).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_4and5_GCD).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_3GCD).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_1and2_GCD).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_1and2_GCD)
		}
	} else {
		s.Advance()
	}
	return sim.CurrentTime
}

/*
else if ffExpiresAt+4500*time.Millisecond <= sim.GetMaxDuration() {
	pesti
}
*/
