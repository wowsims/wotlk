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
					return war.CanRevenge(sim)
				}

				if war.ShieldBlockAura.IsActive() {
					return !war.CanShieldSlam(sim) && war.CanRevenge(sim)
				} else {
					return war.CanRevenge(sim)
				}
			},
		},
		int32(proto.ProtectionWarrior_Rotation_ShieldSlam): {
			Spell: war.ShieldSlam,
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
		int32(proto.ProtectionWarrior_Rotation_Devastate): {
			Spell:     war.Devastate,
			Condition: war.CanDevastate,
		},
		int32(proto.ProtectionWarrior_Rotation_SunderArmor): {
			Spell:     war.SunderArmor,
			Condition: war.CanSunderArmor,
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
			Spell:     war.MortalStrike,
			Condition: war.CanMortalStrike,
		},
		int32(proto.ProtectionWarrior_Rotation_ConcussionBlow): {
			Spell:     war.ConcussionBlow,
			Condition: war.CanConcussionBlow,
		},
		int32(proto.ProtectionWarrior_Rotation_Shockwave): {
			Spell:     war.Shockwave,
			Condition: war.CanShockwave,
		},
	})
}
