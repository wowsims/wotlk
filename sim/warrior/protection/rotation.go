package protection

import (
	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (war *ProtectionWarrior) OnGCDReady(sim *core.Simulation) {
	war.doRotation(sim)
}

func (war *ProtectionWarrior) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	war.tryQueueHsCleave(sim)
}

func (war *ProtectionWarrior) doRotation(sim *core.Simulation) {
	if war.CustomRotation != nil {
		war.CustomRotation.Cast(sim)
	} else {
		if war.GCD.IsReady(sim) {
			if war.CanShieldSlam(sim) {
				war.ShieldSlam.Cast(sim, war.CurrentTarget)
			} else if war.CanRevenge(sim) {
				war.Revenge.Cast(sim, war.CurrentTarget)
			} else if war.ShouldShout(sim) {
				war.Shout.Cast(sim, nil)
			} else if war.shouldThunderClap(sim) {
				war.ThunderClap.Cast(sim, war.CurrentTarget)
			} else if war.shouldDemoShout(sim) {
				war.DemoralizingShout.Cast(sim, war.CurrentTarget)
			} else if war.CanMortalStrike(sim) {
				war.MortalStrike.Cast(sim, war.CurrentTarget)
			} else if war.CanDevastate(sim) {
				war.Devastate.Cast(sim, war.CurrentTarget)
			} else if war.CanSunderArmor(sim) {
				war.SunderArmor.Cast(sim, war.CurrentTarget)
			}
		}
	}

	// if we did nothing else, mark we intentionally did nothing here.
	if war.GCD.IsReady(sim) {
		war.DoNothing()
	}

}

func (war *ProtectionWarrior) tryQueueHsCleave(sim *core.Simulation) {
	if war.ShouldQueueHSOrCleave(sim) {
		war.QueueHSOrCleave(sim)
	}
}

func (war *ProtectionWarrior) shouldDemoShout(sim *core.Simulation) bool {
	return war.ShouldDemoralizingShout(sim,
		war.Rotation.DemoShoutChoice == proto.ProtectionWarrior_Rotation_DemoShoutChoiceFiller,
		war.Rotation.DemoShoutChoice == proto.ProtectionWarrior_Rotation_DemoShoutChoiceMaintain)
}

func (war *ProtectionWarrior) shouldThunderClap(sim *core.Simulation) bool {
	return war.ShouldThunderClap(sim,
		war.Rotation.ThunderClapChoice == proto.ProtectionWarrior_Rotation_ThunderClapChoiceOnCD,
		war.Rotation.ThunderClapChoice == proto.ProtectionWarrior_Rotation_ThunderClapChoiceMaintain,
		false)
}

func (war *ProtectionWarrior) makeCustomRotation() *common.CustomRotation {
	return common.NewCustomRotation(war.Rotation.CustomRotation, war.GetCharacter(), map[int32]common.CustomSpell{
		int32(proto.ProtectionWarrior_Rotation_Revenge): common.CustomSpell{
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := war.Revenge.CurCast.Cost
				return war.Revenge.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				if !war.Rotation.PrioSslamOnShieldBlock {
					return war.CanRevenge(sim)
				}

				if war.ShieldBlockAura.IsActive() {
					return !war.CanShieldSlam(sim) && war.CanRevenge(sim)
				} else {
					return war.CanRevenge(sim)
				}
			},
		},
		int32(proto.ProtectionWarrior_Rotation_ShieldSlam): common.CustomSpell{
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := war.ShieldSlam.CurCast.Cost
				return war.ShieldSlam.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				if !war.Rotation.PrioSslamOnShieldBlock {
					return war.CanShieldSlam(sim)
				}

				if war.ShieldBlockAura.IsActive() {
					return war.CanShieldSlam(sim)
				} else {
					return !war.CanRevenge(sim) && war.CanShieldSlam(sim)
				}
			},
		},
		int32(proto.ProtectionWarrior_Rotation_Devastate): common.CustomSpell{
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := war.Devastate.CurCast.Cost
				return war.Devastate.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return war.CanDevastate(sim)
			},
		},
		int32(proto.ProtectionWarrior_Rotation_SunderArmor): common.CustomSpell{
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := war.SunderArmor.CurCast.Cost
				return war.SunderArmor.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return war.CanSunderArmor(sim)
			},
		},
		int32(proto.ProtectionWarrior_Rotation_DemoralizingShout): common.CustomSpell{
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := war.DemoralizingShout.CurCast.Cost
				return war.DemoralizingShout.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return war.shouldDemoShout(sim)
			},
		},
		int32(proto.ProtectionWarrior_Rotation_ThunderClap): common.CustomSpell{
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := war.ThunderClap.CurCast.Cost
				return war.ThunderClap.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return war.shouldThunderClap(sim)
			},
		},
		int32(proto.ProtectionWarrior_Rotation_Shout): common.CustomSpell{
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := war.Shout.CurCast.Cost
				return war.Shout.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return war.ShouldShout(sim)
			},
		},
		int32(proto.ProtectionWarrior_Rotation_MortalStrike): common.CustomSpell{
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := war.MortalStrike.CurCast.Cost
				return war.MortalStrike.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return war.CanMortalStrike(sim)
			},
		},
		int32(proto.ProtectionWarrior_Rotation_ConcussionBlow): common.CustomSpell{
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := war.ConcussionBlow.CurCast.Cost
				return war.ConcussionBlow.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return war.CanConcussionBlow(sim)
			},
		},
		int32(proto.ProtectionWarrior_Rotation_Shockwave): common.CustomSpell{
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := war.Shockwave.CurCast.Cost
				return war.Shockwave.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return war.CanShockwave(sim)
			},
		},
	})
}
