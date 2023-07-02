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
		spell, target := hunter.singleTargetChooseSpell(sim)

		success := spell.Cast(sim, target)
		if !success {
			hunter.WaitForMana(sim, spell.CurCast.Cost)
		}
	}
}

func (hunter *Hunter) aoeChooseSpell(sim *core.Simulation) *core.Spell {
	if hunter.Rotation.TrapWeave && hunter.ExplosiveTrap.IsReady(sim) && !hunter.ExplosiveTrap.AOEDot().IsActive() {
		return hunter.TrapWeaveSpell
	} else {
		return hunter.Volley
	}
}

func (hunter *Hunter) singleTargetChooseSpell(sim *core.Simulation) (*core.Spell, *core.Unit) {
	for _, spell := range hunter.rotationPriority {
		if spell == nil {
			continue
		}

		if spell == hunter.SerpentSting && hunter.Rotation.MultiDotSerpentSting {
			for i := int32(0); i < hunter.Env.GetNumTargets(); i++ {
				if hunter.rotationConditions[spell].CanUse(sim, hunter.Env.GetTargetUnit(i)) {
					return spell, hunter.Env.GetTargetUnit(i)
				}
			}
		} else if hunter.rotationConditions[spell].CanUse(sim, hunter.CurrentTarget) {
			return spell, hunter.CurrentTarget
		}
	}
	panic("No spell found to cast!")
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
				return hunter.ExplosiveTrap.IsReady(sim) && !hunter.ExplosiveTrap.AOEDot().IsActive()
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

type RotationCondition struct {
	CanUse func(sim *core.Simulation, target *core.Unit) bool
}

func (hunter *Hunter) initRotation() {
	hunter.rotationConditions = map[*core.Spell]RotationCondition{
		hunter.KillShot: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return sim.IsExecutePhase20() && hunter.KillShot.IsReady(sim)
			},
		},
		hunter.ExplosiveShotR4: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.ExplosiveShotR4.IsReady(sim) && !hunter.ExplosiveShotR4.CurDot().IsActive()
			},
		},
		hunter.ExplosiveShotR3: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.Rotation.AllowExplosiveShotDownrank && hunter.ExplosiveShotR3.IsReady(sim) && !hunter.ExplosiveShotR3.CurDot().IsActive()
			},
		},
		hunter.ScorpidSting: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.Rotation.Sting == proto.Hunter_Rotation_ScorpidSting && !hunter.ScorpidStingAuras.Get(hunter.CurrentTarget).IsActive()
			},
		},
		hunter.SerpentSting: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.Rotation.Sting == proto.Hunter_Rotation_SerpentSting && !hunter.SerpentSting.Dot(target).IsActive()
			},
		},
		hunter.ChimeraShot: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.ChimeraShot.IsReady(sim)
			},
		},
		hunter.BlackArrow: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return !hunter.Rotation.TrapWeave && hunter.BlackArrow.IsReady(sim)
			},
		},
		hunter.TrapWeaveSpell: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.Rotation.TrapWeave && hunter.ExplosiveTrap.IsReady(sim) && !hunter.ExplosiveTrap.AOEDot().IsActive()
			},
		},
		hunter.AimedShot: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.AimedShot.IsReady(sim)
			},
		},
		hunter.MultiShot: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.MultiShot.IsReady(sim)
			},
		},
		hunter.ArcaneShot: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.ArcaneShot.IsReady(sim) && (!hunter.ExplosiveShotR4.CurDot().IsActive() && !hunter.ExplosiveShotR3.CurDot().IsActive())
			},
		},
		hunter.SteadyShot: RotationCondition{
			func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.SteadyShot.IsReady(sim)
			},
		},
	}

	if hunter.PrimaryTalentTree == 0 {
		// BM
		hunter.rotationPriority = []*core.Spell{
			hunter.KillShot,
			hunter.TrapWeaveSpell,
			hunter.SerpentSting,
			hunter.ScorpidSting,
			hunter.AimedShot,
			hunter.MultiShot,
			hunter.SteadyShot,
		}
	} else if hunter.PrimaryTalentTree == 1 {
		// MM
		hunter.rotationPriority = []*core.Spell{
			hunter.KillShot,
			hunter.SerpentSting,
			hunter.ScorpidSting,
			hunter.TrapWeaveSpell,
			hunter.ChimeraShot,
			hunter.AimedShot,
			hunter.MultiShot,
			hunter.SteadyShot,
		}
	} else {
		// SV
		hunter.rotationPriority = []*core.Spell{
			hunter.KillShot,
			hunter.ExplosiveShotR4,
			hunter.ExplosiveShotR3,
			hunter.TrapWeaveSpell,
			hunter.SerpentSting,
			hunter.ScorpidSting,
			hunter.BlackArrow,
			hunter.MultiShot,
			hunter.SteadyShot,
		}
	}
}
