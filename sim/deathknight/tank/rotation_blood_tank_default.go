package tank

import (
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

	return false
}
