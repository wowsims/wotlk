package protection

import (
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
	// Always save rage for Shield Slam
	if war.ShieldSlam.IsReady(sim) && !war.HasEnoughRageForShieldSlam() {
		war.DoNothing()
		return
	}

	if war.GCD.IsReady(sim) {
		if war.Rotation.SpamRevengeHs {
			if war.CanRevenge(sim) {
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
		} else {
			if war.Rotation.PrioRevengeOverShieldSlam && war.CanRevenge(sim) {
				war.Revenge.Cast(sim, war.CurrentTarget)
			} else if war.CanShieldSlam(sim) {
				war.ShieldSlam.Cast(sim, war.CurrentTarget)
			} else if !war.Rotation.PrioRevengeOverShieldSlam && war.CanRevenge(sim) {
				war.Revenge.Cast(sim, war.CurrentTarget)
			} else if war.ShouldShout(sim) {
				war.Shout.Cast(sim, nil)
			} else if war.shouldThunderClap(sim) {
				war.ThunderClap.Cast(sim, war.CurrentTarget)
			} else if war.shouldDemoShout(sim) {
				war.DemoralizingShout.Cast(sim, war.CurrentTarget)
			} else if war.Rotation.UseShockwaveSt && war.CanShockwave(sim) {
				war.Shockwave.Cast(sim, war.CurrentTarget)
			} else if war.Rotation.UseConcussionBlowSt && war.CanConcussionBlow(sim) {
				war.ConcussionBlow.Cast(sim, war.CurrentTarget)
			} else if war.CanMortalStrike(sim) {
				war.MortalStrike.Cast(sim, war.CurrentTarget)
			} else if war.CanDevastate(sim) {
				war.Devastate.Cast(sim, war.CurrentTarget)
			} else if war.CanSunderArmor(sim) {
				war.SunderArmor.Cast(sim, war.CurrentTarget)
			}
		}
	}
	war.tryQueueHsCleave(sim)

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
		war.Rotation.DemoShout == proto.ProtectionWarrior_Rotation_DemoShoutFiller,
		war.Rotation.DemoShout == proto.ProtectionWarrior_Rotation_DemoShoutMaintain)
}

func (war *ProtectionWarrior) shouldThunderClap(sim *core.Simulation) bool {
	return war.ShouldThunderClap(sim,
		war.Rotation.ThunderClap == proto.ProtectionWarrior_Rotation_ThunderClapOnCD,
		war.Rotation.ThunderClap == proto.ProtectionWarrior_Rotation_ThunderClapMaintain,
		false)
}
