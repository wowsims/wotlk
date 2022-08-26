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
	if hunter.Rotation.TrapWeave && hunter.ExplosiveTrap.IsReady(sim) {
		return hunter.TrapWeaveSpell
	} else {
		return hunter.Volley
	}
}

func (hunter *Hunter) singleTargetChooseSpell(sim *core.Simulation) *core.Spell {
	if sim.IsExecutePhase20() && hunter.KillShot.IsReady(sim) {
		return hunter.KillShot
	} else if hunter.ExplosiveShot.IsReady(sim) && !hunter.ExplosiveShotDot.IsActive() {
		return hunter.ExplosiveShot
	} else if hunter.Rotation.Sting == proto.Hunter_Rotation_ScorpidSting && !hunter.ScorpidStingAura.IsActive() {
		return hunter.ScorpidSting
	} else if hunter.Rotation.Sting == proto.Hunter_Rotation_SerpentSting && !hunter.SerpentStingDot.IsActive() {
		return hunter.SerpentSting
	} else if hunter.ChimeraShot.IsReady(sim) {
		return hunter.ChimeraShot
	} else if !hunter.Rotation.TrapWeave && hunter.BlackArrow.IsReady(sim) {
		return hunter.BlackArrow
	} else if hunter.Rotation.TrapWeave && hunter.ExplosiveTrap.IsReady(sim) {
		return hunter.TrapWeaveSpell
	} else if hunter.AimedShot.IsReady(sim) {
		return hunter.AimedShot
	} else if hunter.MultiShot.IsReady(sim) {
		return hunter.MultiShot
	} else if hunter.ArcaneShot.IsReady(sim) && (hunter.ExplosiveShotDot == nil || !hunter.ExplosiveShotDot.IsActive()) {
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
		int32(proto.Hunter_Rotation_ArcaneShot): common.CustomSpell{
			Spell: hunter.ArcaneShot,
			Condition: func(sim *core.Simulation) bool {
				return hunter.ArcaneShot.IsReady(sim) && (hunter.ExplosiveShotDot == nil || !hunter.ExplosiveShotDot.IsActive())
			},
		},
		int32(proto.Hunter_Rotation_AimedShot): common.CustomSpell{
			Spell: hunter.AimedShot,
			Condition: func(sim *core.Simulation) bool {
				return hunter.AimedShot.IsReady(sim)
			},
		},
		int32(proto.Hunter_Rotation_BlackArrow): common.CustomSpell{
			Spell: hunter.BlackArrow,
			Condition: func(sim *core.Simulation) bool {
				return hunter.BlackArrow.IsReady(sim)
			},
		},
		int32(proto.Hunter_Rotation_ChimeraShot): common.CustomSpell{
			Spell: hunter.ChimeraShot,
			Condition: func(sim *core.Simulation) bool {
				return hunter.ChimeraShot.IsReady(sim)
			},
		},
		int32(proto.Hunter_Rotation_ExplosiveShot): common.CustomSpell{
			Spell: hunter.ExplosiveShot,
			Condition: func(sim *core.Simulation) bool {
				return hunter.ExplosiveShot.IsReady(sim) && !hunter.ExplosiveShotDot.IsActive()
			},
		},
		int32(proto.Hunter_Rotation_ExplosiveTrap): common.CustomSpell{
			Spell: hunter.TrapWeaveSpell,
			Condition: func(sim *core.Simulation) bool {
				return hunter.ExplosiveTrap.IsReady(sim)
			},
		},
		int32(proto.Hunter_Rotation_KillShot): common.CustomSpell{
			Spell: hunter.KillShot,
			Condition: func(sim *core.Simulation) bool {
				return sim.IsExecutePhase20() && hunter.KillShot.IsReady(sim)
			},
		},
		int32(proto.Hunter_Rotation_MultiShot): common.CustomSpell{
			Spell: hunter.MultiShot,
			Condition: func(sim *core.Simulation) bool {
				return hunter.MultiShot.IsReady(sim)
			},
		},
		int32(proto.Hunter_Rotation_ScorpidStingSpell): common.CustomSpell{
			Spell: hunter.ScorpidSting,
			Condition: func(sim *core.Simulation) bool {
				return hunter.Rotation.Sting == proto.Hunter_Rotation_ScorpidSting && !hunter.ScorpidStingAura.IsActive()
			},
		},
		int32(proto.Hunter_Rotation_SerpentStingSpell): common.CustomSpell{
			Spell: hunter.SerpentSting,
			Condition: func(sim *core.Simulation) bool {
				return hunter.Rotation.Sting == proto.Hunter_Rotation_SerpentSting && !hunter.SerpentStingDot.IsActive()
			},
		},
		int32(proto.Hunter_Rotation_SteadyShot): common.CustomSpell{
			Spell: hunter.SteadyShot,
			Condition: func(sim *core.Simulation) bool {
				return true
			},
		},
		int32(proto.Hunter_Rotation_Volley): common.CustomSpell{
			Spell: hunter.Volley,
			Condition: func(sim *core.Simulation) bool {
				return true
			},
		},
	})
}
