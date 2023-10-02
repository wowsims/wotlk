package feral

import (
	"slices"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type PoolingAction struct {
	refreshTime time.Duration
	cost        float64
}

type PoolingActions struct {
	actions []PoolingAction
}

func (pa *PoolingActions) create(prealloc uint) {
	pa.actions = make([]PoolingAction, 0, prealloc)
}

func (pa *PoolingActions) addAction(t time.Duration, cost float64) {
	pa.actions = append(pa.actions, PoolingAction{t, cost})
}

func (pa *PoolingActions) sort() {
	slices.SortStableFunc(pa.actions, func(p1, p2 PoolingAction) int {
		return int(p1.refreshTime - p2.refreshTime)
	})
}

func (pa *PoolingActions) calcFloatingEnergy(cat *FeralDruid, sim *core.Simulation) float64 {
	floatingEnergy := 0.0
	previousTime := sim.CurrentTime
	tfPending := false

	for _, s := range pa.actions {
		delta_t := float64((s.refreshTime - previousTime) / core.EnergyTickDuration)
		if !tfPending {
			tfPending = cat.tfExpectedBefore(sim, s.refreshTime)
			if tfPending {
				s.cost -= 60
			}
		}

		if delta_t < s.cost {
			floatingEnergy += s.cost - delta_t
			previousTime = s.refreshTime
		} else {
			previousTime += time.Duration(s.cost * float64(core.EnergyTickDuration))
		}
	}

	return floatingEnergy
}

func (pa *PoolingActions) nextRefreshTime() (bool, time.Duration) {
	if len(pa.actions) > 0 {
		return true, pa.actions[0].refreshTime
	}
	return false, 0
}
