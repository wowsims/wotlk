package tank

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *TankDeathknight) setupBloodTankERWOpener() {
	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_IT)

	dk.Main.
		NewAction(dk.RotationActionCallback_TankBlood_PrioRotation)
}

func (dk *TankDeathknight) RotationActionCallback_TankBlood_PrioRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	attackGcd := 1500 * time.Millisecond
	spellGcd := dk.SpellGCD()
	ff := dk.FrostFeverDisease[target.Index].IsActive()
	bp := dk.BloodPlagueDisease[target.Index].IsActive()
	fbAt := core.MinDuration(dk.FrostFeverDisease[target.Index].ExpiresAt(), dk.BloodPlagueDisease[target.Index].ExpiresAt())

	if !ff {
		return dk.CastIcyTouch(sim, target)
	} else if !bp {
		return dk.CastPlagueStrike(sim, target)
	} else if dk.CanDeathStrike(sim) && sim.CurrentTime+attackGcd < fbAt {
		return dk.CastDeathStrike(sim, target)
	} else {
		if sim.CurrentTime < fbAt-2*spellGcd {
			dk.WaitUntil(sim, fbAt-2*spellGcd)
			return false
		} else {
			return dk.CastPestilence(sim, target)
		}
	}
}
