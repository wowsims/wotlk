package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
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

	minClickLatency := time.Millisecond * 150

	switch action {
	case RotationAction_IT:
		casted = deathKnight.CastIcyTouch(sim, target)
		// Add this line if you care about recasting a spell in the opener in
		// case it missed
		advance = deathKnight.LastCastOutcome != core.OutcomeMiss
	case RotationAction_PS:
		casted = deathKnight.CastPlagueStrike(sim, target)
		advance = deathKnight.LastCastOutcome.Matches(core.OutcomeHit | core.OutcomeCrit)
	case RotationAction_UA:
		casted = deathKnight.CastUnbreakableArmor(sim, target)
		// Add this line if your spell does not incur a GCD or you will hang!
		deathKnight.WaitUntil(sim, sim.CurrentTime+minClickLatency)
	case RotationAction_BT:
		casted = deathKnight.CastBloodTap(sim, target)
		deathKnight.WaitUntil(sim, sim.CurrentTime+minClickLatency)
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
		deathKnight.WaitUntil(sim, sim.CurrentTime+minClickLatency)
	case RotationAction_HB_Ghoul_RimeCheck:
		// You can do custom actions, this is deciding whether to HB or raise dead
		if deathKnight.RimeAura.IsActive() {
			casted = deathKnight.CastHowlingBlast(sim, target)
		} else {
			casted = deathKnight.CastRaiseDead(sim, target)
		}
	case RotationAction_BS:
		casted = deathKnight.CastBloodStrike(sim, target)
		advance = deathKnight.LastCastOutcome != core.OutcomeMiss
	case RotationAction_SS:
		casted = deathKnight.CastScourgeStrike(sim, target)
		advance = deathKnight.LastCastOutcome.Matches(core.OutcomeHit | core.OutcomeCrit)
	case RotationAction_DND:
		casted = deathKnight.CastDeathAndDecay(sim, target)
	case RotationAction_GF:
		casted = deathKnight.CastGhoulFrenzy(sim, target)
	case RotationAction_DC:
		casted = deathKnight.CastDeathCoil(sim, target)
	case RotationAction_Garg:
		casted = deathKnight.CastSummonGargoyle(sim, target)
	case RotationAction_AOTD:
		casted = deathKnight.CastArmyOfTheDead(sim, target)
	case RotationAction_BP:
		casted = deathKnight.CastBloodPresence(sim, target)
		if !casted {
			deathKnight.WaitUntil(sim, deathKnight.BloodPresence.CD.ReadyAt())
		} else {
			deathKnight.WaitUntil(sim, sim.CurrentTime+minClickLatency)
		}
	case RotationAction_FP:
		casted = deathKnight.CastFrostPresence(sim, target)
		if !casted {
			deathKnight.WaitUntil(sim, deathKnight.FrostPresence.CD.ReadyAt())
		} else {
			deathKnight.WaitUntil(sim, sim.CurrentTime+minClickLatency)
		}
	case RotationAction_UP:
		casted = deathKnight.CastUnholyPresence(sim, target)
		if !casted {
			deathKnight.WaitUntil(sim, deathKnight.UnholyPresence.CD.ReadyAt())
		} else {
			deathKnight.WaitUntil(sim, sim.CurrentTime+minClickLatency)
		}
	}

	// Advances the opener
	if casted && advance {
		o.idx += 1
	}

	return casted
}

func (o *Sequence) DoNext(sim *core.Simulation, deathKnight *DeathKnight) bool {
	target := deathKnight.CurrentTarget
	casted := &deathKnight.CastSuccessful
	*casted = false

	if o.IsOngoing() {
		*casted = deathKnight.opener.DoAction(sim, target, deathKnight)
	} else if deathKnight.sequence != nil {
		if deathKnight.sequence.IsOngoing() {
			*casted = deathKnight.sequence.DoAction(sim, target, deathKnight)
			if !deathKnight.sequence.IsOngoing() {
				deathKnight.sequence = nil
			}
		}
	} else {
		deathKnight.onOpener = false

		if deathKnight.DoRotationEvent == nil {
			panic("Missing rotation event. Please assign one during spec creation")
		}
		deathKnight.DoRotationEvent(sim, target)
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
	deathKnight.onOpener = true
}
