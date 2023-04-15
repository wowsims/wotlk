package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/warrior"
)

const DebuffRefreshWindow = time.Second * 2

func (war *DpsWarrior) OnGCDReady(sim *core.Simulation) {
	rendRemainingDur := war.RendValidUntil - sim.CurrentTime

	// Pause rotation on every rend tick to check if TFB procs
	if rendRemainingDur != war.Rend.CurDot().Duration && rendRemainingDur%3 == 0 && war.Talents.TasteForBlood > 0 {
		core.StartDelayedAction(sim, core.DelayedActionOptions{
			DoAt: sim.CurrentTime + time.Microsecond*1,
			OnAction: func(_ *core.Simulation) {
				war.doRotation(sim)
			},
		})
	} else {
		war.doRotation(sim)
	}

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
		if war.ThunderClap.CanCast(sim, war.CurrentTarget) {
			if war.ThunderClap.Cast(sim, war.CurrentTarget) {
				if war.ThunderClapAuras.Get(war.CurrentTarget).RemainingDuration(sim) > DebuffRefreshWindow {
					war.thunderClapNext = false

					// Switching back to berserker immediately is unrealistic because the player needs
					// to visually confirm the TC landed. Instead we add a delay to model that.
					war.canSwapStanceAt = sim.CurrentTime + time.Millisecond*300
				}
				return
			}
		}
	} else {
		if war.Rotation.StanceOption == proto.Warrior_Rotation_BerserkerStance {
			war.trySwapToBerserker(sim)
		} else if war.Rotation.StanceOption == proto.Warrior_Rotation_BattleStance {
			war.trySwapToBattle(sim)
		}
	}

	if war.shouldSunder(sim) {
		if war.Devastate != nil {
			war.Devastate.Cast(sim, war.CurrentTarget)
		} else {
			war.SunderArmor.Cast(sim, war.CurrentTarget)
		}
		war.tryQueueHsCleave(sim)
		return
	}

	IsExecutePhase20 := sim.IsExecutePhase20()
	if war.Rotation.CustomRotationOption {
		war.CustomRotation.Cast(sim)
	} else {
		if IsExecutePhase20 {
			war.executeRotation(sim)
		} else {
			war.normalRotation(sim)
		}
	}

	if war.GCD.IsReady(sim) && !war.thunderClapNext {
		// We didn't cast anything, so wait for the next CD.
		nextCD := war.Whirlwind.CD.ReadyAt()
		if war.Bloodthirst != nil && war.Bloodthirst.CD.ReadyAt() < nextCD {
			nextCD = war.Bloodthirst.CD.ReadyAt()
		} else if war.MortalStrike != nil && war.MortalStrike.CD.ReadyAt() < nextCD {
			nextCD = war.MortalStrike.CD.ReadyAt()
		}

		if war.Rotation.SunderArmor == proto.Warrior_Rotation_SunderArmorMaintain {
			nextSunderAt := war.SunderArmorAuras.Get(war.CurrentTarget).ExpiresAt() - SunderWindow
			// TODO looks fishy, nextCD is unused
			nextCD = core.MinDuration(nextCD, nextSunderAt)
		}
	}
}

func (war *DpsWarrior) normalRotation(sim *core.Simulation) {
	if war.GCD.IsReady(sim) {
		if war.Warrior.PrimaryTalentTree == warrior.FuryTree {
			war.furyNormalRotation(sim)
		} else if war.Warrior.PrimaryTalentTree == warrior.ArmsTree {
			war.armsNormalRotation(sim)
		}
	}
	war.tryQueueHsCleave(sim)
}

func (war *DpsWarrior) executeRotation(sim *core.Simulation) {
	if war.GCD.IsReady(sim) {
		if war.Warrior.PrimaryTalentTree == warrior.FuryTree {
			war.furyExecuteRotation(sim)
		} else if war.Warrior.PrimaryTalentTree == warrior.ArmsTree {
			war.armsExecuteRotation(sim)
		}
	}

	war.tryQueueHsCleave(sim)
}

func (war *DpsWarrior) furyNormalRotation(sim *core.Simulation) {
	if war.tryMaintainDebuffs(sim) {
		war.DoNothing()
	} else if war.Rotation.MainGcd == proto.Warrior_Rotation_Slam && war.ShouldInstantSlam(sim) {
		war.CastSlam(sim, war.CurrentTarget)
	} else if war.Rotation.MainGcd == proto.Warrior_Rotation_Bloodthirst && war.Bloodthirst.CanCast(sim, war.CurrentTarget) {
		war.Bloodthirst.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.MainGcd == proto.Warrior_Rotation_Whirlwind && war.CanWhirlwind(sim) {
		war.Whirlwind.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.MainGcd != proto.Warrior_Rotation_Slam && war.ShouldInstantSlam(sim) {
		war.CastSlam(sim, war.CurrentTarget)
	} else if war.Rotation.MainGcd != proto.Warrior_Rotation_Bloodthirst && war.Bloodthirst.CanCast(sim, war.CurrentTarget) {
		war.Bloodthirst.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.MainGcd != proto.Warrior_Rotation_Whirlwind && war.CanWhirlwind(sim) {
		war.Whirlwind.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.UseRend && war.ShouldRend(sim) {
		if !war.StanceMatches(warrior.BattleStance) {
			if !war.BattleStance.IsReady(sim) {
				return
			}
			war.BattleStance.Cast(sim, nil)
		}
		war.Rend.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.UseOverpower && war.ShouldOverpower(sim) {
		if !war.StanceMatches(warrior.BattleStance) {
			if !war.BattleStance.IsReady(sim) {
				return
			}
			war.BattleStance.Cast(sim, nil)
		}
		war.Overpower.Cast(sim, war.CurrentTarget)
	}
}

func (war *DpsWarrior) armsNormalRotation(sim *core.Simulation) {
	if war.tryMaintainDebuffs(sim) {
		war.DoNothing()
	} else if war.Execute.CanCast(sim, war.CurrentTarget) {
		war.CastExecute(sim, war.CurrentTarget)
	} else if war.Rotation.UseRend && war.ShouldRend(sim) {
		if !war.StanceMatches(warrior.BattleStance) {
			if !war.BattleStance.IsReady(sim) {
				return
			}
			war.BattleStance.Cast(sim, nil)
		}
		war.Rend.Cast(sim, war.CurrentTarget)
	} else if war.ShouldOverpower(sim) {
		if !war.StanceMatches(warrior.BattleStance) {
			if !war.BattleStance.IsReady(sim) {
				return
			}
			war.BattleStance.Cast(sim, nil)
		}
		war.Overpower.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.UseMs && war.MortalStrike.CanCast(sim, war.CurrentTarget) && war.CurrentRage() >= war.Rotation.MsRageThreshold {
		war.MortalStrike.Cast(sim, war.CurrentTarget)
	} else if war.Slam.CanCast(sim, war.CurrentTarget) && war.CurrentRage() >= war.Rotation.SlamRageThreshold {
		war.CastSlam(sim, war.CurrentTarget)
	}
}

func (war *DpsWarrior) furyExecuteRotation(sim *core.Simulation) {
	if war.tryMaintainDebuffs(sim) {
		war.DoNothing()
	} else if war.SpamExecute(war.Rotation.SpamExecute) {
		war.CastExecute(sim, war.CurrentTarget)
	} else if war.Rotation.MainGcd == proto.Warrior_Rotation_Slam &&
		war.Rotation.UseSlamOverExecute && war.ShouldInstantSlam(sim) {
		war.CastSlam(sim, war.CurrentTarget)
	} else if war.Rotation.MainGcd == proto.Warrior_Rotation_Bloodthirst &&
		war.Rotation.UseBtDuringExecute && war.Bloodthirst.CanCast(sim, war.CurrentTarget) {
		war.Bloodthirst.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.MainGcd == proto.Warrior_Rotation_Whirlwind &&
		war.Rotation.UseWwDuringExecute && war.CanWhirlwind(sim) {
		war.Whirlwind.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.MainGcd != proto.Warrior_Rotation_Slam &&
		war.Rotation.UseSlamOverExecute && war.ShouldInstantSlam(sim) {
		war.CastSlam(sim, war.CurrentTarget)
	} else if war.Rotation.MainGcd != proto.Warrior_Rotation_Bloodthirst &&
		war.Rotation.UseBtDuringExecute && war.Bloodthirst.CanCast(sim, war.CurrentTarget) {
		war.Bloodthirst.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.MainGcd != proto.Warrior_Rotation_Whirlwind &&
		war.Rotation.UseWwDuringExecute && war.CanWhirlwind(sim) {
		war.Whirlwind.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.UseRend && war.ShouldRend(sim) {
		if !war.StanceMatches(warrior.BattleStance) {
			if !war.BattleStance.IsReady(sim) {
				return
			}
			war.BattleStance.Cast(sim, nil)
		}
		war.Rend.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.UseOverpower && war.Rotation.ExecutePhaseOverpower &&
		war.ShouldOverpower(sim) {
		if !war.StanceMatches(warrior.BattleStance) {
			if !war.BattleStance.IsReady(sim) {
				return
			}
			war.BattleStance.Cast(sim, nil)
		}
		war.Overpower.Cast(sim, war.CurrentTarget)
	} else if war.Execute.CanCast(sim, war.CurrentTarget) {
		war.CastExecute(sim, war.CurrentTarget)
	}
}

func (war *DpsWarrior) armsExecuteRotation(sim *core.Simulation) {
	if war.tryMaintainDebuffs(sim) {
		war.DoNothing()
	} else if war.IsSuddenDeathActive() && war.Execute.CanCast(sim, war.CurrentTarget) {
		war.CastExecute(sim, war.CurrentTarget)
	} else if war.ShouldOverpower(sim) {
		if !war.StanceMatches(warrior.BattleStance) {
			if !war.BattleStance.IsReady(sim) {
				return
			}
			war.BattleStance.Cast(sim, nil)
		}
		war.Overpower.Cast(sim, war.CurrentTarget)
	} else if war.SpamExecute(war.Rotation.SpamExecute) {
		war.CastExecute(sim, war.CurrentTarget)
	} else if war.Rotation.UseRend && war.ShouldRend(sim) {
		if !war.StanceMatches(warrior.BattleStance) {
			if !war.BattleStance.IsReady(sim) {
				return
			}
			war.BattleStance.Cast(sim, nil)
		}
		war.Rend.Cast(sim, war.CurrentTarget)
	} else if war.Rotation.UseMs && war.MortalStrike.CanCast(sim, war.CurrentTarget) && war.CurrentRage() >= war.Rotation.MsRageThreshold {
		war.MortalStrike.Cast(sim, war.CurrentTarget)
	} else if war.Execute.CanCast(sim, war.CurrentTarget) {
		war.CastExecute(sim, war.CurrentTarget)
	}
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

	if !war.SunderArmor.CanCast(sim, war.CurrentTarget) {
		return false
	}

	saAura := war.SunderArmorAuras.Get(war.CurrentTarget)
	stacks := saAura.GetStacks()
	if war.Rotation.SunderArmor == proto.Warrior_Rotation_SunderArmorHelpStack && stacks == 5 {
		war.maintainSunder = false
	}

	return stacks < 5 || saAura.RemainingDuration(sim) <= SunderWindow
}

// Returns whether any ability was cast.
func (war *DpsWarrior) tryMaintainDebuffs(sim *core.Simulation) bool {
	if war.ShouldShout(sim) {
		war.Shout.Cast(sim, nil)
		return true
	} else if war.Rotation.MaintainDemoShout && war.ShouldDemoralizingShout(sim, war.CurrentTarget, false, true) {
		war.DemoralizingShout.Cast(sim, war.CurrentTarget)
		return true
	} else if war.Rotation.MaintainThunderClap && war.ShouldThunderClap(sim, war.CurrentTarget, false, true, true) {
		war.thunderClapNext = true
		if !war.StanceMatches(warrior.BattleStance) {
			if !war.BattleStance.IsReady(sim) {
				return false
			}
			war.BattleStance.Cast(sim, nil)
		}
		// Need to check again because we might have lost rage from switching stances.
		if war.ThunderClap.CanCast(sim, war.CurrentTarget) {
			war.ThunderClap.Cast(sim, war.CurrentTarget)
			if war.ThunderClapAuras.Get(war.CurrentTarget).RemainingDuration(sim) > DebuffRefreshWindow {
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

func (war *DpsWarrior) makeCustomRotation() *common.CustomRotation {
	return common.NewCustomRotation(war.Rotation.CustomRotation, war.GetCharacter(), map[int32]common.CustomSpell{
		int32(proto.Warrior_Rotation_BloodthirstCustom): {
			Spell: war.Bloodthirst,
			Condition: func(sim *core.Simulation) bool {
				if sim.IsExecutePhase20() && !war.Rotation.UseBtDuringExecute {
					return false
				}
				return war.Bloodthirst.CanCast(sim, war.CurrentTarget)
			},
		},
		int32(proto.Warrior_Rotation_MortalStrike): {
			Spell: war.MortalStrike,
			Condition: func(sim *core.Simulation) bool {
				return war.MortalStrike.CanCast(sim, war.CurrentTarget) && war.CurrentRage() >= war.Rotation.MsRageThreshold
			},
		},
		int32(proto.Warrior_Rotation_WhirlwindCustom): {
			Spell: war.Whirlwind,
			Condition: func(sim *core.Simulation) bool {
				if sim.IsExecutePhase20() && !war.Rotation.UseWwDuringExecute {
					return false
				}

				if !war.StanceMatches(warrior.BerserkerStance) {
					if !war.BerserkerStance.IsReady(sim) {
						return false
					}
					war.BerserkerStance.Cast(sim, nil)
				}
				return war.Whirlwind.CanCast(sim, war.CurrentTarget)
			},
		},
		int32(proto.Warrior_Rotation_SlamCustom): {
			Spell: war.Slam,
			Condition: func(sim *core.Simulation) bool {
				if sim.IsExecutePhase20() && !war.Rotation.UseSlamOverExecute {
					return false
				}

				if (war.ShouldSlam(sim) && war.CurrentRage() >= war.Rotation.SlamRageThreshold || war.ShouldInstantSlam(sim)) &&
					war.Slam.CanCast(sim, war.CurrentTarget) {
					war.AutoAttacks.DelayMeleeBy(sim, war.Slam.CurCast.CastTime)
					return true
				}
				return false
			},
		},

		int32(proto.Warrior_Rotation_SlamExpiring): {
			Spell: war.Slam,
			Condition: func(sim *core.Simulation) bool {
				if !war.ShouldInstantSlam(sim) {
					return false
				}

				if (war.BloodsurgeValidUntil - sim.CurrentTime) > war.BloodsurgeDurationThreshold {
					return false
				}

				if sim.IsExecutePhase20() && !war.Rotation.UseSlamOverExecute {
					return false
				}

				if war.CurrentRage() >= war.Rotation.SlamRageThreshold && war.Slam.CanCast(sim, war.CurrentTarget) {
					war.AutoAttacks.DelayMeleeBy(sim, war.Slam.CurCast.CastTime)
					return true
				}
				return false
			},
		},

		int32(proto.Warrior_Rotation_Rend): {
			Spell: war.Rend,
			Condition: func(sim *core.Simulation) bool {
				if !war.ShouldRend(sim) {
					return false
				}

				if !war.StanceMatches(warrior.BattleStance) {
					if !war.BattleStance.IsReady(sim) {
						return false
					}
					war.BattleStance.Cast(sim, nil)
				}
				return war.Rend.CanCast(sim, war.CurrentTarget)
			},
		},
		int32(proto.Warrior_Rotation_Overpower): {
			Spell: war.Overpower,
			Condition: func(sim *core.Simulation) bool {
				if !war.ShouldOverpower(sim) {
					return false
				}
				if sim.IsExecutePhase20() && !war.Rotation.ExecutePhaseOverpower {
					return false
				}

				if !war.StanceMatches(warrior.BattleStance) {
					if !war.BattleStance.IsReady(sim) {
						return false
					}
					war.BattleStance.Cast(sim, nil)
				}
				return war.Overpower.CanCast(sim, war.CurrentTarget)
			},
		},
		int32(proto.Warrior_Rotation_Execute): {
			Spell: war.Execute,
		},
		int32(proto.Warrior_Rotation_ThunderClap): {
			Spell: war.ThunderClap,
			Condition: func(sim *core.Simulation) bool {
				if !war.StanceMatches(warrior.BattleStance) {
					if !war.BattleStance.IsReady(sim) {
						return false
					}
					war.BattleStance.Cast(sim, nil)
				}
				return war.ThunderClap.CanCast(sim, war.CurrentTarget)
			},
		},
	})
}
