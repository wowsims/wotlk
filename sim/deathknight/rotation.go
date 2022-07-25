package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	if !dk.onOpener {
		if dk.GCD.IsReady(sim) {
			dk.tryUseGCD(sim)
		}
	}
}

func (dk *Deathknight) OnGCDReady(sim *core.Simulation) {
	dk.tryUseGCD(sim)
}

func (dk *Deathknight) tryUseGCD(sim *core.Simulation) {
	if dk.GCD.IsReady(sim) {
		dk.DoRotation(sim)
	}
}

func (o *Sequence) IsOngoing() bool {
	return o.idx < o.numActions
}

func (o *Sequence) DoAction(sim *core.Simulation, target *core.Unit, dk *Deathknight) bool {
	casted := false
	advance := true
	action := o.actions[o.idx]

	minClickLatency := time.Millisecond * 0

	switch action {
	case RotationAction_IT:
		casted = dk.CastIcyTouch(sim, target)
		// Add this line if you care about recasting a spell in the opener in
		// case it missed
		advance = dk.LastCastOutcome != core.OutcomeMiss
	case RotationAction_PS:
		casted = dk.CastPlagueStrike(sim, target)
		advance = dk.LastCastOutcome.Matches(core.OutcomeHit | core.OutcomeCrit)
	case RotationAction_UA:
		casted = dk.CastUnbreakableArmor(sim, target)
		// Add this line if your spell does not incur a GCD or you will hang!
		dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
	case RotationAction_BT:
		casted = dk.CastBloodTap(sim, target)
		dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
	case RotationAction_Obli:
		casted = dk.CastObliterate(sim, target)
	case RotationAction_FS:
		casted = dk.CastFrostStrike(sim, target)
	case RotationAction_Pesti:
		casted = dk.CastPestilence(sim, target)
		if dk.LastCastOutcome == core.OutcomeMiss {
			advance = false
		}
	case RotationAction_ERW:
		casted = dk.CastEmpowerRuneWeapon(sim, target)
		dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
	case RotationAction_HB_Ghoul_RimeCheck:
		// You can do custom actions, this is deciding whether to HB or raise dead
		if dk.RimeAura.IsActive() {
			casted = dk.CastHowlingBlast(sim, target)
		} else {
			casted = dk.CastRaiseDead(sim, target)
		}
	case RotationAction_BS:
		casted = dk.CastBloodStrike(sim, target)
		advance = dk.LastCastOutcome != core.OutcomeMiss
	case RotationAction_SS:
		casted = dk.CastScourgeStrike(sim, target)
		advance = dk.LastCastOutcome.Matches(core.OutcomeHit | core.OutcomeCrit)
	case RotationAction_DND:
		casted = dk.CastDeathAndDecay(sim, target)
	case RotationAction_GF:
		casted = dk.CastGhoulFrenzy(sim, target)
	case RotationAction_DC:
		casted = dk.CastDeathCoil(sim, target)
	case RotationAction_Garg:
		casted = dk.CastSummonGargoyle(sim, target)
	case RotationAction_AOTD:
		casted = dk.CastArmyOfTheDead(sim, target)
	case RotationAction_BP:
		casted = dk.CastBloodPresence(sim, target)
		if !casted {
			dk.WaitUntil(sim, dk.BloodPresence.CD.ReadyAt())
		} else {
			dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
		}
	case RotationAction_FP:
		casted = dk.CastFrostPresence(sim, target)
		if !casted {
			dk.WaitUntil(sim, dk.FrostPresence.CD.ReadyAt())
		} else {
			dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
		}
	case RotationAction_UP:
		casted = dk.CastUnholyPresence(sim, target)
		if !casted {
			dk.WaitUntil(sim, dk.UnholyPresence.CD.ReadyAt())
		} else {
			dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
		}
	}

	// Advances the opener
	if casted && advance {
		o.idx += 1
	}

	return casted
}

func (o *Sequence) DoNext(sim *core.Simulation, dk *Deathknight) bool {
	target := dk.CurrentTarget
	casted := &dk.CastSuccessful
	*casted = false

	if o.IsOngoing() {
		*casted = dk.opener.DoAction(sim, target, dk)
	} else if dk.sequence != nil {
		if dk.sequence.IsOngoing() {
			*casted = dk.sequence.DoAction(sim, target, dk)
			if !dk.sequence.IsOngoing() {
				dk.sequence = nil
			}
		}
	} else {
		dk.onOpener = false

		if dk.DoRotationEvent == nil {
			panic("Missing rotation event. Please assign one during spec creation")
		}
		dk.DoRotationEvent(sim, target)
	}

	return *casted
}

func (dk *Deathknight) DoRotation(sim *core.Simulation) {
	opener := dk.opener
	if !opener.DoNext(sim, dk) {
		if dk.GCD.IsReady(sim) && !dk.IsWaiting() {
			waitUntil := dk.AutoAttacks.MainhandSwingAt
			if dk.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
				waitUntil = core.MinDuration(waitUntil, dk.AutoAttacks.OffhandSwingAt)
			}
			waitUntil = core.MinDuration(waitUntil, dk.AnyRuneReadyAt(sim))
			dk.WaitUntil(sim, waitUntil)
		} else { // No resources
			waitUntil := dk.AnySpentRuneReadyAt(sim)
			dk.WaitUntil(sim, waitUntil)
		}
	}
}

func (dk *Deathknight) ResetRotation(sim *core.Simulation) {
	dk.opener.idx = 0
	dk.onOpener = true
}
