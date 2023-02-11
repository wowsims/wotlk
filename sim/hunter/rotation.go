package hunter

import (
	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (hunter *Hunter) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	hunter.mayMoveAt = sim.CurrentTime
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

	if hunter.SilencingShot.IsReady(sim) {
		hunter.SilencingShot.Cast(sim, hunter.CurrentTarget)
	}

	if hunter.Rotation.Type == proto.Hunter_Rotation_Custom {
		hunter.CustomRotation.Cast(sim)
	} else if hunter.Rotation.Type == proto.Hunter_Rotation_Aoe {
		spell := hunter.aoeChooseSpell(sim)

		success := spell.Cast(sim, hunter.CurrentTarget)
		if !success {
			hunter.WaitForMana(sim, spell.CurCast.Cost)
		}
	} else {
		spell := hunter.singleTargetChooseSpell(sim)

		success := spell.Cast(sim, hunter.CurrentTarget)
		if !success {
			hunter.WaitForMana(sim, spell.CurCast.Cost)
		}
	}
}

func (hunter *Hunter) aoeChooseSpell(sim *core.Simulation) *core.Spell {
	if hunter.Rotation.TrapWeave && hunter.ExplosiveTrap.IsReady(sim) && !hunter.ExplosiveTrapDot.IsActive() {
		return hunter.TrapWeaveSpell
	} else {
		return hunter.Volley
	}
}

func (hunter *Hunter) singleTargetChooseSpell(sim *core.Simulation) *core.Spell {
	if sim.IsExecutePhase20() && hunter.KillShot.IsReady(sim) {
		return hunter.KillShot
	} else if hunter.ExplosiveShotR4.IsReady(sim) && !hunter.ExplosiveShotR4.CurDot().IsActive() {
		return hunter.ExplosiveShotR4
	} else if hunter.Rotation.AllowExplosiveShotDownrank && hunter.ExplosiveShotR3.IsReady(sim) && !hunter.ExplosiveShotR3.CurDot().IsActive() {
		return hunter.ExplosiveShotR3
	} else if hunter.Rotation.Sting == proto.Hunter_Rotation_ScorpidSting && !hunter.ScorpidStingAuras.Get(hunter.CurrentTarget).IsActive() {
		return hunter.ScorpidSting
	} else if hunter.Rotation.Sting == proto.Hunter_Rotation_SerpentSting && !hunter.SerpentSting.CurDot().IsActive() {
		return hunter.SerpentSting
	} else if hunter.ChimeraShot.IsReady(sim) {
		return hunter.ChimeraShot
	} else if !hunter.Rotation.TrapWeave && hunter.BlackArrow.IsReady(sim) {
		return hunter.BlackArrow
	} else if hunter.Rotation.TrapWeave && hunter.ExplosiveTrap.IsReady(sim) && !hunter.ExplosiveTrapDot.IsActive() {
		return hunter.TrapWeaveSpell
	} else if hunter.AimedShot.IsReady(sim) {
		return hunter.AimedShot
	} else if hunter.MultiShot.IsReady(sim) {
		return hunter.MultiShot
	} else if hunter.ArcaneShot.IsReady(sim) && (hunter.ExplosiveShotR4 == nil || (!hunter.ExplosiveShotR4.CurDot().IsActive() && !hunter.ExplosiveShotR3.CurDot().IsActive())) {
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
		}
		hunter.AspectOfTheViper.Cast(sim, nil)
		return true
	}
	return false
}

func (hunter *Hunter) makeCustomRotation() *common.CustomRotation {
	return common.NewCustomRotation(hunter.Rotation.CustomRotation, hunter.GetCharacter(), map[int32]common.CustomSpell{
		int32(proto.Hunter_Rotation_ArcaneShot): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.ArcaneShot.CurCast.Cost
				return hunter.ArcaneShot.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return hunter.ArcaneShot.IsReady(sim) && (hunter.ExplosiveShotR4 == nil || (!hunter.ExplosiveShotR4.CurDot().IsActive() && !hunter.ExplosiveShotR3.CurDot().IsActive()))
			},
		},
		int32(proto.Hunter_Rotation_AimedShot): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.AimedShot.CurCast.Cost
				return hunter.AimedShot.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return hunter.AimedShot.IsReady(sim)
			},
		},
		int32(proto.Hunter_Rotation_BlackArrow): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.BlackArrow.CurCast.Cost
				return hunter.BlackArrow.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return hunter.BlackArrow.IsReady(sim)
			},
		},
		int32(proto.Hunter_Rotation_ChimeraShot): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.ChimeraShot.CurCast.Cost
				return hunter.ChimeraShot.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return hunter.ChimeraShot.IsReady(sim)
			},
		},
		int32(proto.Hunter_Rotation_ExplosiveShot): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.ExplosiveShotR4.CurCast.Cost
				return hunter.ExplosiveShotR4.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return hunter.ExplosiveShotR4.IsReady(sim) && !hunter.ExplosiveShotR4.CurDot().IsActive()
			},
		},
		int32(proto.Hunter_Rotation_ExplosiveShotDownrank): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.ExplosiveShotR3.CurCast.Cost
				return hunter.ExplosiveShotR3.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return hunter.ExplosiveShotR3.IsReady(sim) && !hunter.ExplosiveShotR3.CurDot().IsActive()
			},
		},
		int32(proto.Hunter_Rotation_ExplosiveTrap): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.TrapWeaveSpell.CurCast.Cost
				return hunter.TrapWeaveSpell.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return hunter.ExplosiveTrap.IsReady(sim) && !hunter.ExplosiveTrapDot.IsActive()
			},
		},
		int32(proto.Hunter_Rotation_KillShot): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.KillShot.CurCast.Cost
				return hunter.KillShot.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return sim.IsExecutePhase20() && hunter.KillShot.IsReady(sim)
			},
		},
		int32(proto.Hunter_Rotation_MultiShot): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.MultiShot.CurCast.Cost
				return hunter.MultiShot.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return hunter.MultiShot.IsReady(sim)
			},
		},
		int32(proto.Hunter_Rotation_ScorpidStingSpell): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.ScorpidSting.CurCast.Cost
				return hunter.ScorpidSting.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return hunter.Rotation.Sting == proto.Hunter_Rotation_ScorpidSting && !hunter.ScorpidStingAuras.Get(hunter.CurrentTarget).IsActive()
			},
		},
		int32(proto.Hunter_Rotation_SerpentStingSpell): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.SerpentSting.CurCast.Cost
				return hunter.SerpentSting.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return hunter.Rotation.Sting == proto.Hunter_Rotation_SerpentSting && !hunter.SerpentSting.CurDot().IsActive()
			},
		},
		int32(proto.Hunter_Rotation_SteadyShot): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.SteadyShot.CurCast.Cost
				return hunter.SteadyShot.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return true
			},
		},
		int32(proto.Hunter_Rotation_Volley): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := hunter.Volley.CurCast.Cost
				return hunter.Volley.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return true
			},
		},
	})
}
