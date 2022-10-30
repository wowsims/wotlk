package healing

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (hpriest *HealingPriest) OnGCDReady(sim *core.Simulation) {
	hpriest.tryUseGCD(sim)
}

func (hpriest *HealingPriest) tryUseGCD(sim *core.Simulation) {
	spell := hpriest.chooseSpell(sim)

	if success := spell.Cast(sim, hpriest.CurrentTarget); !success {
		hpriest.WaitForMana(sim, spell.CurCast.Cost)
	}
}

func (hpriest *HealingPriest) chooseSpell(sim *core.Simulation) *core.Spell {
	if !hpriest.RenewHots[hpriest.CurrentTarget.UnitIndex].IsActive() {
		return hpriest.Renew
	} else if hpriest.CanCastPWS(sim, hpriest.CurrentTarget) {
		return hpriest.PowerWordShield
	} else {
		for !hpriest.spellCycle[hpriest.nextCycleIndex].IsReady(sim) {
			hpriest.nextCycleIndex = (hpriest.nextCycleIndex + 1) % len(hpriest.spellCycle)
		}
		spell := hpriest.spellCycle[hpriest.nextCycleIndex]
		hpriest.nextCycleIndex = (hpriest.nextCycleIndex + 1) % len(hpriest.spellCycle)
		return spell
	}
}
