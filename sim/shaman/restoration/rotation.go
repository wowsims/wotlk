package restoration

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (resto *RestorationShaman) OnGCDReady(sim *core.Simulation) {
	resto.tryUseGCD(sim)
}

func (resto *RestorationShaman) tryUseGCD(sim *core.Simulation) {
	if resto.TryDropTotems(sim) {
		return
	}

	spell := resto.LesserHealingWave

	if !spell.Cast(sim, resto.CurrentTarget) {
		resto.WaitForMana(sim, spell.CurCast.Cost)
	}
	resto.WaitUntil(sim, sim.CurrentTime+time.Second*5)
}
