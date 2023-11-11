package protection

import (
	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/warrior"
)

func (war *ProtectionWarrior) OnGCDReady(sim *core.Simulation) {
	war.doRotation(sim)
}

func (war *ProtectionWarrior) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	war.tryQueueHsCleave(sim)
}

func (war *ProtectionWarrior) doRotation(sim *core.Simulation) {
	war.trySwapToDefensive(sim)
	war.CustomRotation.Cast(sim)

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
	return war.ShouldDemoralizingShout(sim, war.CurrentTarget,
		war.Rotation.DemoShoutChoice == proto.ProtectionWarrior_Rotation_DemoShoutChoiceFiller,
		war.Rotation.DemoShoutChoice == proto.ProtectionWarrior_Rotation_DemoShoutChoiceMaintain)
}

func (war *ProtectionWarrior) shouldThunderClap(sim *core.Simulation) bool {
	return war.ShouldThunderClap(sim, war.CurrentTarget,
		war.Rotation.ThunderClapChoice == proto.ProtectionWarrior_Rotation_ThunderClapChoiceOnCD,
		war.Rotation.ThunderClapChoice == proto.ProtectionWarrior_Rotation_ThunderClapChoiceMaintain,
		false)
}

func (war *ProtectionWarrior) trySwapToDefensive(sim *core.Simulation) bool {
	if !war.StanceMatches(warrior.DefensiveStance) && war.DefensiveStance.IsReady(sim) {
		war.DefensiveStance.Cast(sim, nil)
		return true
	}
	return false
}

func (war *ProtectionWarrior) makeCustomRotation() *common.CustomRotation {
	return common.NewCustomRotation(war.Rotation.CustomRotation, war.GetCharacter(), map[int32]common.CustomSpell{
		int32(proto.ProtectionWarrior_Rotation_Revenge): {
			Spell: war.Revenge,
			Condition: func(sim *core.Simulation) bool {
				if !war.Rotation.PrioSslamOnShieldBlock {
					return war.Revenge.CanCast(sim, war.CurrentTarget)
				}

				if war.ShieldBlockAura.IsActive() {
					return !war.ShieldSlam.CanCast(sim, war.CurrentTarget) && war.Revenge.CanCast(sim, war.CurrentTarget)
				} else {
					return war.Revenge.CanCast(sim, war.CurrentTarget)
				}
			},
		},
		int32(proto.ProtectionWarrior_Rotation_ShieldSlam): {
			Spell: war.ShieldSlam,
			Condition: func(sim *core.Simulation) bool {
				if !war.Rotation.PrioSslamOnShieldBlock {
					return war.ShieldSlam.CanCast(sim, war.CurrentTarget)
				}

				if war.ShieldBlockAura.IsActive() {
					return war.ShieldSlam.CanCast(sim, war.CurrentTarget)
				} else {
					return !war.Revenge.CanCast(sim, war.CurrentTarget) && war.ShieldSlam.CanCast(sim, war.CurrentTarget)
				}
			},
		},
		int32(proto.ProtectionWarrior_Rotation_Devastate): {
			Spell: war.Devastate,
		},
		int32(proto.ProtectionWarrior_Rotation_SunderArmor): {
			Spell: war.SunderArmor,
		},
		int32(proto.ProtectionWarrior_Rotation_DemoralizingShout): {
			Spell:     war.DemoralizingShout,
			Condition: war.shouldDemoShout,
		},
		int32(proto.ProtectionWarrior_Rotation_ThunderClap): {
			Spell:     war.ThunderClap,
			Condition: war.shouldThunderClap,
		},
		int32(proto.ProtectionWarrior_Rotation_Shout): {
			Spell:     war.Shout,
			Condition: war.ShouldShout,
		},
		int32(proto.ProtectionWarrior_Rotation_MortalStrike): {
			Spell: war.MortalStrike,
		},
		int32(proto.ProtectionWarrior_Rotation_ConcussionBlow): {
			Spell: war.ConcussionBlow,
		},
		int32(proto.ProtectionWarrior_Rotation_Shockwave): {
			Spell: war.Shockwave,
		},
	})
}
