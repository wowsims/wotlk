package tank

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *TankDeathknight) setupBloodTankERWOpener() {
	dk.RotationSequence.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_DS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_DS).
		NewAction(dk.RotationActionCallback_TankBlood_PrioRotation)
}

func (dk *TankDeathknight) setupBloodTankERWThreatOpener() {
	dk.RotationSequence.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_TankBlood_PrioRotation)
}

func (dk *TankDeathknight) RotationActionCallback_TankBlood_PrioRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if !dk.GCD.IsReady(sim) {
		return dk.NextGCDAt()
	}

	t := sim.CurrentTime
	ff := dk.FrostFeverDisease[target.Index].ExpiresAt() - t
	bp := dk.BloodPlagueDisease[target.Index].ExpiresAt() - t
	fd := dk.CurrentFrostRunes() + dk.CurrentDeathRunes()
	ud := dk.CurrentUnholyRunes() + dk.CurrentDeathRunes()
	b, f, u := dk.NormalCurrentRunes()

	if ff <= 0 {
		dk.IcyTouch.Cast(sim, target)
		return -1
	}

	if bp <= 0 {
		dk.PlagueStrike.Cast(sim, target)
		return -1
	}

	if ff <= 2*time.Second || bp <= 2*time.Second {
		dk.Pestilence.Cast(sim, target)
		return -1
	}

	if f > 0 && u > 0 && (dk.CurrentHealthPercent() > 0.5 && dk.CurrentHealth()+dk.AverageDSHeal() <= 1.05*dk.MaxHealth()) {
		dk.DeathStrike.Cast(sim, target)
		return -1
	} else if f > 1 && u > 1 && !dk.IsMainTank() {
		dk.DeathStrike.Cast(sim, target)
		return -1
	} else if f < 2 && u < 2 && !dk.IsMainTank() {
		dk.IcyTouch.Cast(sim, target)
		return -1
	} else if fd > 0 && ud > 0 && (dk.CurrentHealthPercent() > 0.5 && dk.CurrentHealth()+dk.AverageDSHeal() <= 1.05*dk.MaxHealth()) {
		dk.DeathStrike.Cast(sim, target)
		return -1
	}
	//else if fd > 0 && ud > 0 {
	//	dk.IcyTouch.Cast(sim, target)
	//	return -1
	//}

	if dk.BloodTap.CanCast(sim) {
		dk.BloodTap.Cast(sim, target)
		dk.IcyTouch.Cast(sim, target)
		dk.CancelBloodTap(sim)
		return -1
	}

	if b == 1 && dk.NormalSpentBloodRuneReadyAt(sim)-t < ff-2*time.Second && dk.NormalSpentBloodRuneReadyAt(sim)-t < bp-2*time.Second {
		dk.BloodStrike.Cast(sim, target)
		return -1
	}

	return -1
}
