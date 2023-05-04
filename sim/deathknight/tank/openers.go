package tank

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *TankDeathknight) TankRA_BloodSpell(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := dk.BloodSpell.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *TankDeathknight) TankRA_FuSpell(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := dk.FuSpell.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *TankDeathknight) TankRA_IT(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := dk.IcyTouch.Cast(sim, target)
	if !casted && dk.Talents.UnbreakableArmor && !dk.UnbreakableArmor.IsReady(sim) {
		s.Advance()
		return -1
	}
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *TankDeathknight) setupTankRegularERWOpener() {
	dk.RotationSequence.
		NewAction(dk.TankRA_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.TankRA_FuSpell).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.TankRA_IT).
		NewAction(dk.TankRA_BloodSpell).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.TankRA_IT).
		NewAction(dk.TankRA_IT).
		NewAction(dk.TankRA_IT).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.TankRA_FuSpell)
}

func (dk *TankDeathknight) setupTankThreatERWOpener() {
	dk.RotationSequence.
		NewAction(dk.TankRA_IT).
		NewAction(dk.TankRA_IT).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.TankRA_IT).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.TankRA_IT).
		NewAction(dk.TankRA_IT).
		NewAction(dk.TankRA_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.TankRA_BloodSpell)
}
