package tank

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"slices"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

const threatOpenerCastsBeforeBloodTap = 2
const normalOpenerCastsBeforeBloodTap = 3

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

func shouldUseBloodTapInOpener(dk *TankDeathknight) bool {
	bloodTapDefensiveCd := dk.GetMajorCooldown(dk.BloodTap.ActionID)
	if bloodTapDefensiveCd != nil {
		timings := bloodTapDefensiveCd.GetTimings()
		slices.Sort(timings)
		if len(timings) == 0 ||
			(len(timings) > 0 && timings[0] <
				dk.BloodTap.CD.Duration+getPlannedOpenerBloodTapUsageTime(dk)) {
			return false
		} else {
			return true
		}
	}
	return true
}

func getPlannedOpenerBloodTapUsageTime(dk *TankDeathknight) time.Duration {
	if dk.Rotation.Opener == proto.TankDeathknight_Rotation_Threat {
		return threatOpenerCastsBeforeBloodTap * core.GCDDefault
	}
	return normalOpenerCastsBeforeBloodTap * core.GCDDefault
}

func (dk *TankDeathknight) setupTankRegularERWOpener() {
	if shouldUseBloodTapInOpener(dk) {
		dk.setupTankRegularERWOpenerWithBloodTap()
	} else {
		dk.setupTankRegularERWOpenerWithoutBloodTap()
	}
}

func (dk *TankDeathknight) setupTankRegularERWOpenerWithBloodTap() {
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

func (dk *TankDeathknight) setupTankRegularERWOpenerWithoutBloodTap() {
	dk.RotationSequence.
		NewAction(dk.TankRA_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.TankRA_FuSpell).
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
	if shouldUseBloodTapInOpener(dk) {
		dk.setupTankThreatERWOpenerWithBloodTap()
	} else {
		dk.setupTankThreatERWOpenerWithoutBloodTap()
	}
}

func (dk *TankDeathknight) setupTankThreatERWOpenerWithBloodTap() {
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

// If maintaining GCD parity is desired, the threat opener without BT should cast DS with
// the GCD previously allocated to the Blood Tapped IT before using ERW
func (dk *TankDeathknight) setupTankThreatERWOpenerWithoutBloodTap() {
	dk.setupTankRegularERWOpenerWithoutBloodTap()
}
