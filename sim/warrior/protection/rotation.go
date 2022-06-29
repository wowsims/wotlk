package protection

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

func (war *ProtectionWarrior) OnGCDReady(sim *core.Simulation) {
	war.doRotation(sim)
}

func (war *ProtectionWarrior) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	war.tryQueueHsCleave(sim)
}

func (war *ProtectionWarrior) doRotation(sim *core.Simulation) {
	if war.GCD.IsReady(sim) {
		if war.CanShieldSlam(sim) {
			war.ShieldSlam.Cast(sim, war.CurrentTarget)
		} else if war.CanBloodthirst(sim) {
			war.Bloodthirst.Cast(sim, war.CurrentTarget)
		} else if war.CanMortalStrike(sim) {
			war.MortalStrike.Cast(sim, war.CurrentTarget)
		} else if war.CanRevenge(sim) {
			war.Revenge.Cast(sim, war.CurrentTarget)
		} else if war.ShouldShout(sim) {
			war.Shout.Cast(sim, nil)
		} else if war.shouldThunderClap(sim) {
			war.ThunderClap.Cast(sim, war.CurrentTarget)
		} else if war.shouldDemoShout(sim) {
			war.DemoralizingShout.Cast(sim, war.CurrentTarget)
		} else if war.CanDevastate(sim) {
			war.Devastate.Cast(sim, war.CurrentTarget)
		} else if war.CanSunderArmor(sim) {
			war.SunderArmor.Cast(sim, war.CurrentTarget)
		}
	}

	war.tryShieldBlock(sim)
	war.tryQueueHsCleave(sim)
}

func (war *ProtectionWarrior) tryShieldBlock(sim *core.Simulation) {
	if war.Rotation.ShieldBlock == proto.ProtectionWarrior_Rotation_ShieldBlockOnCD ||
		(war.Rotation.ShieldBlock == proto.ProtectionWarrior_Rotation_ShieldBlockToProcRevenge && war.Revenge.IsReady(sim) && sim.CurrentTime >= war.RevengeValidUntil) {
		if war.CanShieldBlock(sim) {
			war.ShieldBlock.Cast(sim, nil)
		}
	}
}

func (war *ProtectionWarrior) tryQueueHsCleave(sim *core.Simulation) {
	if war.ShouldQueueHSOrCleave(sim) {
		war.QueueHSOrCleave(sim)
	}
}

func (war *ProtectionWarrior) shouldDemoShout(sim *core.Simulation) bool {
	return war.ShouldDemoralizingShout(sim,
		war.Rotation.DemoShout == proto.ProtectionWarrior_Rotation_DemoShoutFiller,
		war.Rotation.DemoShout == proto.ProtectionWarrior_Rotation_DemoShoutMaintain)
}

func (war *ProtectionWarrior) shouldThunderClap(sim *core.Simulation) bool {
	return war.ShouldThunderClap(sim,
		war.Rotation.ThunderClap == proto.ProtectionWarrior_Rotation_ThunderClapOnCD,
		war.Rotation.ThunderClap == proto.ProtectionWarrior_Rotation_ThunderClapMaintain,
		false)
}
