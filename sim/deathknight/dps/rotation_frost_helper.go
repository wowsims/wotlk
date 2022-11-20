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
func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_EndOfFight_1GCD(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	frAt := dk.NormalFrostRuneReadyAt(sim)
	uhAt := dk.NormalUnholyRuneReadyAt(sim)
	obAt := core.MaxDuration(frAt, uhAt)
	//ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
	//bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	if obAt < sim.CurrentTime {
		s.Clear().NewAction(dk.RotationActionCallback_FrostSubUnh_EndOfFight_Obli)
	} else if dk.CurrentRunicPower() <= fsCost {
		s.Clear().NewAction(dk.RotationActionCallback_FS)
	} else if dk.Rime() {
		s.Clear().NewAction(dk.RotationActionCallback_HB)
	} else if dk.CurrentBloodRunes() >= 1 {
		s.Clear().NewAction(dk.RotationActionCallback_BS)
	} else if dk.CurrentFrostRunes() >= 1 {
		s.Clear().NewAction(dk.RotationActionCallback_IT)
	} else if dk.CurrentUnholyRunes() >= 1 {
		s.Clear().NewAction(dk.RotationActionCallback_PS)
	}
	return sim.CurrentTime
}

// end of fight logic in the works
func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_EndOfFight(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if sim.Log != nil && sim.GetRemainingDuration() < 1500*time.Millisecond*4 {
		sim.Log("endoffight function entered!")
	} else {
		s.Advance()
	}
	frAt := dk.NormalFrostRuneReadyAt(sim)
	uhAt := dk.NormalUnholyRuneReadyAt(sim)
	obAt := core.MaxDuration(frAt, uhAt)
	extraFS := 0.0
	if sim.GetRemainingDuration() < 1500*time.Millisecond*3 && obAt > sim.Duration {
		extraFS = 1
	} else {
		extraFS = 0
	}

	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	if sim.GetRemainingDuration() < 1500*time.Millisecond*1 {
		s.Clear().NewAction(dk.RotationActionCallback_FrostSubUnholy_EndOfFight_1GCD)
	} else if sim.GetRemainingDuration() < 1500*time.Millisecond*2 {
		if dk.CurrentRunicPower() >= fsCost*(1+extraFS) {
			s.Clear().
				NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Dump).
				NewAction(dk.RotationActionCallback_FrostSubUnh_EndOfFight_Obli)
		} else if dk.CurrentBloodRunes() > 0 {
			s.Clear().
				NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Dump).
				NewAction(dk.RotationActionCallback_BS).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli)
		}

		sim.GetRemainingDuration()
	} else if sim.GetRemainingDuration() < 1500*time.Millisecond*3 {
		if dk.CurrentRunicPower() >= fsCost*(2+extraFS) {
			s.Clear().
				NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Dump).
				NewAction(dk.RotationActionCallback_FrostSubUnh_EndOfFight_Obli)
		} else if dk.CurrentRunicPower() >= fsCost*(1+extraFS) && dk.NormalSpentBloodRuneReadyAt(sim) <= sim.CurrentTime+1500*time.Millisecond {
			s.Clear().
				NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Dump).
				NewAction(dk.RotationActionCallback_BS).
				NewAction(dk.RotationActionCallback_FrostSubUnh_EndOfFight_Obli)
		}
	} else if sim.GetRemainingDuration() < 1500*time.Millisecond*4 {
		if dk.CurrentRunicPower() >= fsCost*(1+extraFS) && dk.CurrentBloodRunes() > 0 {
			s.Clear().
				NewAction(dk.RotationActionCallback_Pesti).
				NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Dump).
				NewAction(dk.RotationActionCallback_FrostSubUnh_EndOfFight_Obli)
		}

	}
	/*{
		if (sim.GetRemainingDuration().Minutes() < 1500*time.Millisecond * 1) {
			Prio oblit > FS > HB(rime) > BS > IT > PS
		} else if (sim.GetRemainingDuration().Minutes() < 1500*time.Millisecond * 2) {
			Prio oblit > FS > HB(rime) > BS > IT > PS
		} else if (sim.GetRemainingDuration().Minutes() < 1500*time.Millisecond * 3) {
			Prio oblit > FS > HB(rime) > Pesti > BS > IT > PS
		} else if (sim.GetRemainingDuration().Minutes() < 1500*time.Millisecond * 4) {
			Prio oblit > FS > HB(rime) > Pesti > BS
		} else if (sim.GetRemainingDuration().Minutes() < 1500*time.Millisecond * 5) {
			Prio oblit > FS > HB(rime) > Pesti > BS
		}
	}*/

	return sim.CurrentTime
}

/*
else if ffExpiresAt+4500*time.Millisecond <= sim.GetMaxDuration() {
	pesti
}
*/
