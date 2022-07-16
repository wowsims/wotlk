package hunter

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (hunter *Hunter) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	hunter.TryUseCooldowns(sim)
	if hunter.GCD.IsReady(sim) {
		hunter.rotation(sim)
	}
}

func (hunter *Hunter) OnGCDReady(sim *core.Simulation) {
	hunter.rotation(sim)
}

func (hunter *Hunter) rotation(sim *core.Simulation) {
	hunter.trySwapAspect(sim)

	spell := hunter.chooseSpell(sim)

	if spell != nil {
		success := spell.Cast(sim, hunter.CurrentTarget)
		if !success {
			hunter.WaitForMana(sim, spell.CurCast.Cost)
		}
	}
}

func (hunter *Hunter) chooseSpell(sim *core.Simulation) *core.Spell {
	if hunter.Rotation.Sting == proto.Hunter_Rotation_ScorpidSting && !hunter.ScorpidStingAura.IsActive() {
		return hunter.ScorpidSting
	} else if hunter.Rotation.Sting == proto.Hunter_Rotation_SerpentSting && !hunter.SerpentStingDot.IsActive() {
		return hunter.SerpentSting
	} else if sim.IsExecutePhase20() && hunter.KillShot.IsReady(sim) {
		return hunter.KillShot
	} else if hunter.ChimeraShot != nil && hunter.ChimeraShot.IsReady(sim) {
		return hunter.ChimeraShot
	} else if hunter.BlackArrow != nil && hunter.BlackArrow.IsReady(sim) {
		return hunter.BlackArrow
	} else if hunter.ExplosiveShot != nil && hunter.ExplosiveShot.IsReady(sim) && !hunter.ExplosiveShotDot.IsActive() {
		return hunter.ExplosiveShot
	} else if hunter.AimedShot != nil && hunter.AimedShot.IsReady(sim) {
		return hunter.AimedShot
	} else if hunter.ArcaneShot != nil && hunter.ArcaneShot.IsReady(sim) {
		return hunter.ArcaneShot
	} else {
		return hunter.SteadyShot
	}
}

// Returns whether an aspect was swapped.
func (hunter *Hunter) trySwapAspect(sim *core.Simulation) bool {
	currentMana := hunter.CurrentManaPercent()
	if hunter.currentAspect == hunter.AspectOfTheViperAura && hunter.Rotation.ViperStartManaPercent < 1 {
		if !hunter.permaHawk &&
			hunter.CurrentMana() > hunter.manaSpentPerSecondAtFirstAspectSwap*sim.GetRemainingDuration().Seconds() {
			hunter.permaHawk = true
		}
		if hunter.permaHawk || currentMana > hunter.Rotation.ViperStopManaPercent {
			hunter.AspectOfTheDragonhawk.Cast(sim, nil)
			return true
		}
	} else if hunter.currentAspect != hunter.AspectOfTheViperAura && !hunter.permaHawk && currentMana < hunter.Rotation.ViperStartManaPercent {
		if hunter.manaSpentPerSecondAtFirstAspectSwap == 0 {
			hunter.manaSpentPerSecondAtFirstAspectSwap = (hunter.Metrics.ManaSpent - hunter.Metrics.ManaGained) / sim.CurrentTime.Seconds()
		}
		if !hunter.permaHawk &&
			hunter.CurrentMana() > hunter.manaSpentPerSecondAtFirstAspectSwap*sim.GetRemainingDuration().Seconds() {
			hunter.permaHawk = true
		} else {
			hunter.AspectOfTheViper.Cast(sim, nil)
			return true
		}
	}
	return false
}
