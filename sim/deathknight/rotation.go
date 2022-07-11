package deathknight

import (
	//"math"
	//"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
	//"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) OnGCDReady(sim *core.Simulation) {
	deathKnight.tryUseGCD(sim)
}

func (deathKnight *DeathKnight) tryUseGCD(sim *core.Simulation) {
	//var spell *core.Spell
	var target = deathKnight.CurrentTarget

	if sim.Log != nil {
		deathKnight.Unit.Log(sim, "Trying to use GCD")
	}

	if deathKnight.GCD.IsReady(sim) {
		if deathKnight.CanIcyTouch(sim) {
			deathKnight.IcyTouch.Cast(sim, target)
		} else if deathKnight.CanPlagueStrike(sim) {
			deathKnight.PlagueStrike.Cast(sim, target)
		} else {
			nextCD := deathKnight.IcyTouch.ReadyAt()

			if nextCD > sim.CurrentTime {
				deathKnight.WaitUntil(sim, nextCD)
			}
		}
	}
}

func (deathKnight *DeathKnight) CanIcyTouch(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 10.0, 0, 1, 0, 0) && deathKnight.IcyTouch.IsReady(sim)
}

func (deathKnight *DeathKnight) CanPlagueStrike(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 10.0, 0, 0, 1, 0) && deathKnight.PlagueStrike.IsReady(sim)
}
