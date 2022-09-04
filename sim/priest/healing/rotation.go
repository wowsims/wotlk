package healing

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (hpriest *HealingPriest) OnGCDReady(sim *core.Simulation) {
	hpriest.tryUseGCD(sim)
}

func (hpriest *HealingPriest) tryUseGCD(sim *core.Simulation) {
	var spell *core.Spell
	spell = hpriest.chooseSpell(sim)

	if success := spell.Cast(sim, hpriest.CurrentTarget); !success {
		hpriest.WaitForMana(sim, spell.CurCast.Cost)
	}
}

func (hpriest *HealingPriest) chooseSpell(sim *core.Simulation) *core.Spell {
	if !hpriest.RenewHots[hpriest.CurrentTarget.UnitIndex].IsActive() {
		return hpriest.Renew
	} else {
		return hpriest.GreaterHeal
	}
}
