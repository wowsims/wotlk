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

func (deathKnight *DeathKnight) OnManaTick(sim *core.Simulation) {
	if deathKnight.FinishedWaitingForManaAndGCDReady(sim) {
		deathKnight.tryUseGCD(sim)
	}
}

func (deathKnight *DeathKnight) tryUseGCD(sim *core.Simulation) {
	//var spell *core.Spell
	//var target = deathKnight.CurrentTarget

}
