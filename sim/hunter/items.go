package hunter

import (
	"github.com/wowsims/sod/sim/core"
)

func init() {
	core.NewItemEffect(209823, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()
		if hunter.pet != nil {
			hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.01
		}
	})

}
