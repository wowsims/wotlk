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

	if deathKnight.GCD.IsReady(sim) {
		if deathKnight.CanIcyTouch(sim) {
			deathKnight.IcyTouch.Cast(sim, target)
		}
	}
}

func (deathKnight *DeathKnight) CanIcyTouch(sim *core.Simulation) bool {
	sim.Log(sim, "%d | %d | %d", deathKnight.CurrentRunicPower(), deathKnight.CurrentFrostRunes(), deathKnight.IcyTouch.IsReady(sim))

	return deathKnight.CurrentRunicPower() >= deathKnight.IcyTouch.DefaultCast.Cost && deathKnight.CurrentFrostRunes() > 0 && deathKnight.IcyTouch.IsReady(sim)
}
