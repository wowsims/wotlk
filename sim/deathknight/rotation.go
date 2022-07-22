package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
	//"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	if !deathKnight.onOpener {
		if deathKnight.GCD.IsReady(sim) {
			deathKnight.tryUseGCD(sim)
		}
	}
}

func (deathKnight *DeathKnight) OnGCDReady(sim *core.Simulation) {
	deathKnight.tryUseGCD(sim)
}

func (deathKnight *DeathKnight) shouldWaitForDnD(sim *core.Simulation, blood bool, frost bool, unholy bool) bool {
	return deathKnight.Rotation.UseDeathAndDecay && !(deathKnight.Talents.Morbidity == 0 || !(deathKnight.DeathAndDecay.CD.IsReady(sim) || deathKnight.DeathAndDecay.CD.TimeToReady(sim) < 4*time.Second) || ((!blood || deathKnight.CurrentBloodRunes() > 1) && (!frost || deathKnight.CurrentFrostRunes() > 1) && (!unholy || deathKnight.CurrentUnholyRunes() > 1)))
}

var recastedFF = false
var recastedBP = false

func (deathKnight *DeathKnight) shouldSpreadDisease(sim *core.Simulation) bool {
	return recastedFF && recastedBP && deathKnight.Env.GetNumTargets() > 1
}

func (deathKnight *DeathKnight) spreadDiseases(sim *core.Simulation, target *core.Unit) {
	deathKnight.Pestilence.Cast(sim, target)
	recastedFF = false
	recastedBP = false
}

func (deathKnight *DeathKnight) tryUseGCD(sim *core.Simulation) {
	if deathKnight.GCD.IsReady(sim) {
		deathKnight.DoRotation(sim)
	}
}

func (o *Sequence) IsOngoing() bool {
	return o.idx < o.numActions
}

func (o *Sequence) DoAction(sim *core.Simulation, target *core.Unit, deathKnight *DeathKnight) bool {
	casted := false
	advance := true
	action := o.actions[o.idx]

	switch action {
	case RotationAction_IT:
		casted = deathKnight.CastIcyTouch(sim, target)
		// Add this line if you care about recasting a spell in the opener in
		// case it missed
		advance = deathKnight.LastCastOutcome != core.OutcomeMiss
	case RotationAction_PS:
		casted = deathKnight.CastPlagueStrike(sim, target)
		advance = deathKnight.LastCastOutcome != core.OutcomeMiss
	case RotationAction_UA:
		casted = deathKnight.CastUnbreakableArmor(sim, target)
		// Add this line if your spell does not incur a GCD or you will hang!
		deathKnight.WaitUntil(sim, sim.CurrentTime)
	case RotationAction_BT:
		casted = deathKnight.CastBloodTap(sim, target)
		deathKnight.WaitUntil(sim, sim.CurrentTime)
	case RotationAction_Obli:
		casted = deathKnight.CastObliterate(sim, target)
	case RotationAction_FS:
		casted = deathKnight.CastFrostStrike(sim, target)
	case RotationAction_Pesti:
		casted = deathKnight.CastPestilence(sim, target)
		if deathKnight.LastCastOutcome == core.OutcomeMiss {
			advance = false
		}
	case RotationAction_ERW:
		casted = deathKnight.CastEmpowerRuneWeapon(sim, target)
		deathKnight.WaitUntil(sim, sim.CurrentTime)
	case RotationAction_HB_Ghoul_RimeCheck:
		// You can do custom actions, this is deciding whether to HB or raise dead
		if deathKnight.RimeAura.IsActive() {
			casted = deathKnight.CastHowlingBlast(sim, target)
		} else {
			casted = deathKnight.CastRaiseDead(sim, target)
		}
	case RotationAction_BS:
		casted = deathKnight.CastBloodStrike(sim, target)
	}

	// Advances the opener
	if casted && advance {
		o.idx += 1
	}

	return casted
}

func (o *Sequence) DoNext(sim *core.Simulation, deathKnight *DeathKnight) bool {
	target := deathKnight.CurrentTarget
	casted := &deathKnight.castSuccessful
	*casted = false

	if deathKnight.sequence != nil {
		*casted = deathKnight.sequence.DoAction(sim, target, deathKnight)
		if !deathKnight.sequence.IsOngoing() {
			deathKnight.sequence = nil
		}
	} else if o.IsOngoing() {
		*casted = deathKnight.opener.DoAction(sim, target, deathKnight)
	} else {
		deathKnight.onOpener = false

		if deathKnight.opener.id == RotationID_FrostSubBlood_Full || deathKnight.opener.id == RotationID_FrostSubUnholy_Full {
			deathKnight.doFrostRotation(sim, target)
		} else if deathKnight.opener.id == RotationID_Unholy_Full {
			deathKnight.doUnholyRotation(sim, target)
		}
		// Other prio lists for other specs here just else if {...
	}

	return *casted
}

func (deathKnight *DeathKnight) DoRotation(sim *core.Simulation) {
	opener := deathKnight.opener
	if !opener.DoNext(sim, deathKnight) {
		if deathKnight.GCD.IsReady(sim) && !deathKnight.IsWaiting() {
			waitUntil := deathKnight.AutoAttacks.MainhandSwingAt
			if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
				waitUntil = core.MinDuration(waitUntil, deathKnight.AutoAttacks.OffhandSwingAt)
			}
			waitUntil = core.MinDuration(waitUntil, deathKnight.AnyRuneReadyAt(sim))
			deathKnight.WaitUntil(sim, waitUntil)
		} else { // No resources
			waitUntil := deathKnight.AnySpentRuneReadyAt(sim)
			deathKnight.WaitUntil(sim, waitUntil)
		}
	}
}

func (deathKnight *DeathKnight) ResetRotation(sim *core.Simulation) {
	deathKnight.opener.idx = 0
}
