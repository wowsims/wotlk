package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/warrior"
)

const DebuffRefreshWindow = time.Second * 2

func (war *DpsWarrior) OnGCDReady(sim *core.Simulation) {
	war.doRotation(sim)

	if war.GCD.IsReady(sim) && !war.IsWaiting() {
		// This means we did nothing
		war.DoNothing()
	}
}

func (war *DpsWarrior) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	war.tryQueueHsCleave(sim)
}

func (war *DpsWarrior) doRotation(sim *core.Simulation) {
	if war.thunderClapNext {
		if war.CanThunderClap(sim) {
			if war.ThunderClap.Cast(sim, war.CurrentTarget) {
				if war.ThunderClapAura.RemainingDuration(sim) > DebuffRefreshWindow {
					war.thunderClapNext = false

					// Switching back to berserker immediately is unrealistic because the player needs
					// to visually confirm the TC landed. Instead we add a delay to model that.
					war.canSwapStanceAt = sim.CurrentTime + time.Millisecond*300
				}
				return
			}
		}
	} else {
		if war.Talents.Bloodthirst {
			war.trySwapToBerserker(sim)
		} else if war.Talents.MortalStrike {
			war.trySwapToBattle(sim)
		}
	}

	if war.shouldSunder(sim) {
		if war.Talents.Devastate {
			war.Devastate.Cast(sim, war.CurrentTarget)
		} else {
			war.SunderArmor.Cast(sim, war.CurrentTarget)
		}
		war.tryQueueHsCleave(sim)
		return
	}

	IsExecutePhase20 := sim.IsExecutePhase20()

	if IsExecutePhase20 {
		war.executeRotation(sim)
	} else {
		war.normalRotation(sim)
	}

	if war.GCD.IsReady(sim) && !war.thunderClapNext {
		// We didn't cast anything, so wait for the next CD.
		// Note that BT/MS share a CD timer so we don't need to check MS.
		nextCD := core.MinDuration(war.Bloodthirst.CD.ReadyAt(), war.Whirlwind.CD.ReadyAt())

		if war.Rotation.SunderArmor == proto.Warrior_Rotation_SunderArmorMaintain {
			nextSunderAt := war.SunderArmorAura.ExpiresAt() - SunderWindow
			nextCD = core.MinDuration(nextCD, nextSunderAt)
		}
	}
}

func (war *DpsWarrior) normalRotation(sim *core.Simulation) {
	if war.GCD.IsReady(sim) {
		if war.ShouldInstantSlam(sim) {
			war.CastSlam(sim, war.CurrentTarget)
		} else if war.ShouldOverpower(sim) {
			if !war.StanceMatches(warrior.BattleStance) {
				if !war.BattleStance.IsReady(sim) {
					return
				}
				war.BattleStance.Cast(sim, nil)
			}
			war.Overpower.Cast(sim, war.CurrentTarget)
		} else if war.tryMaintainDebuffs(sim) {
			war.DoNothing()
		} else if war.Rotation.PrioritizeWw && war.CanWhirlwind(sim) {
			war.Whirlwind.Cast(sim, war.CurrentTarget)
		} else if war.CanBloodthirst(sim) {
			war.Bloodthirst.Cast(sim, war.CurrentTarget)
		} else if war.Rotation.UseRend && war.ShouldRend(sim) {
			if war.Talents.Bloodthirst && war.CurrentRage() >= war.Rotation.RendRageThresholdBelow {
				return
			}
			if !war.StanceMatches(warrior.BattleStance) {
				if !war.BattleStance.IsReady(sim) {
					return
				}
				war.BattleStance.Cast(sim, nil)
			}
			war.Rend.Cast(sim, war.CurrentTarget)
		} else if war.SuddenDeathAura.IsActive() && war.CanExecute() {
			war.Execute.Cast(sim, war.CurrentTarget)
		} else if war.Rotation.UseMs && war.CanMortalStrike(sim) && war.CurrentRage() >= war.Rotation.MsRageThreshold {
			war.MortalStrike.Cast(sim, war.CurrentTarget)
		} else if war.ShouldSlam(sim) && war.CurrentRage() >= war.Rotation.SlamRageThreshold {
			war.CastSlam(sim, war.CurrentTarget)
		} else if war.CanShieldSlam(sim) {
			war.ShieldSlam.Cast(sim, war.CurrentTarget)
		} else if !war.Rotation.PrioritizeWw && war.CanWhirlwind(sim) {
			war.Whirlwind.Cast(sim, war.CurrentTarget)
		} else if war.ShouldBerserkerRage(sim) {
			war.BerserkerRage.Cast(sim, nil)
		}
	}
	war.tryQueueHsCleave(sim)
}

func (war *DpsWarrior) executeRotation(sim *core.Simulation) {
	if war.GCD.IsReady(sim) {
		if war.ShouldInstantSlam(sim) {
			war.CastSlam(sim, war.CurrentTarget)
		} else if war.ShouldOverpower(sim) {
			if !war.StanceMatches(warrior.BattleStance) {
				if !war.BattleStance.IsReady(sim) {
					return
				}
				war.BattleStance.Cast(sim, nil)
			}
			war.Overpower.Cast(sim, war.CurrentTarget)
		} else if war.tryMaintainDebuffs(sim) {
			war.DoNothing()
		} else if war.SpamExecute(war.Rotation.SpamExecute) {
			war.Execute.Cast(sim, war.CurrentTarget)
		} else if war.Rotation.PrioritizeWw && war.Rotation.UseWwDuringExecute && war.CanWhirlwind(sim) {
			war.Whirlwind.Cast(sim, war.CurrentTarget)
		} else if war.Rotation.UseBtDuringExecute && war.CanBloodthirst(sim) {
			war.Bloodthirst.Cast(sim, war.CurrentTarget)
		} else if war.Rotation.UseRend && war.ShouldRend(sim) {
			if war.Talents.Bloodthirst && war.CurrentRage() >= war.Rotation.RendRageThresholdBelow {
				return
			}
			if !war.StanceMatches(warrior.BattleStance) {
				if !war.BattleStance.IsReady(sim) {
					return
				}
				war.BattleStance.Cast(sim, nil)
			}
			war.Rend.Cast(sim, war.CurrentTarget)
		} else if war.SuddenDeathAura.IsActive() && war.CanExecute() {
			war.Execute.Cast(sim, war.CurrentTarget)
		} else if war.Rotation.UseMs && war.CanMortalStrike(sim) && war.CurrentRage() >= war.Rotation.MsRageThreshold {
			war.MortalStrike.Cast(sim, war.CurrentTarget)
		} else if war.Rotation.UseSlamOverExecute && war.ShouldSlam(sim) && war.CurrentRage() >= war.Rotation.SlamRageThreshold {
			war.CastSlam(sim, war.CurrentTarget)
		} else if war.CanExecute() {
			war.Execute.Cast(sim, war.CurrentTarget)
		} else if war.ShouldBerserkerRage(sim) {
			war.BerserkerRage.Cast(sim, nil)
		}
	}

	war.tryQueueHsCleave(sim)
}

func (war *DpsWarrior) trySwapToBerserker(sim *core.Simulation) bool {
	if !war.StanceMatches(warrior.BerserkerStance) && war.BerserkerStance.IsReady(sim) && sim.CurrentTime >= war.canSwapStanceAt {
		war.BerserkerStance.Cast(sim, nil)
		return true
	}
	return false
}

func (war *DpsWarrior) trySwapToBattle(sim *core.Simulation) bool {
	if !war.StanceMatches(warrior.BattleStance) && war.BattleStance.IsReady(sim) && sim.CurrentTime >= war.canSwapStanceAt {
		war.BattleStance.Cast(sim, nil)
		return true
	}
	return false
}

const SunderWindow = time.Second * 3

func (war *DpsWarrior) shouldSunder(sim *core.Simulation) bool {
	if !war.maintainSunder {
		return false
	}

	if !war.CanSunderArmor(sim) {
		return false
	}

	stacks := war.SunderArmorAura.GetStacks()
	if war.Rotation.SunderArmor == proto.Warrior_Rotation_SunderArmorHelpStack && stacks == 5 {
		war.maintainSunder = false
	}

	return stacks < 5 || war.SunderArmorAura.RemainingDuration(sim) <= SunderWindow
}

// Returns whether any ability was cast.
func (war *DpsWarrior) tryMaintainDebuffs(sim *core.Simulation) bool {
	if war.ShouldShout(sim) {
		war.Shout.Cast(sim, nil)
		return true
	} else if war.Rotation.MaintainDemoShout && war.ShouldDemoralizingShout(sim, false, true) {
		war.DemoralizingShout.Cast(sim, war.CurrentTarget)
		return true
	} else if war.Rotation.MaintainThunderClap && war.ShouldThunderClap(sim, false, true, true) {
		war.thunderClapNext = true
		if !war.StanceMatches(warrior.BattleStance) {
			if !war.BattleStance.IsReady(sim) {
				return false
			}
			war.BattleStance.Cast(sim, nil)
		}
		// Need to check again because we might have lost rage from switching stances.
		if war.CanThunderClap(sim) {
			war.ThunderClap.Cast(sim, war.CurrentTarget)
			if war.ThunderClapAura.RemainingDuration(sim) > DebuffRefreshWindow {
				war.thunderClapNext = false
			}
		}
		return true
	}
	return false
}

func (war *DpsWarrior) tryQueueHsCleave(sim *core.Simulation) {
	if sim.IsExecutePhase20() && !war.Rotation.UseHsDuringExecute {
		return
	}

	if war.ShouldQueueHSOrCleave(sim) {
		war.QueueHSOrCleave(sim)
	}
}
