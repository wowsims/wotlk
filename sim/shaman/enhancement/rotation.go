package enhancement

import (
	"fmt"
	"time"

	"github.com/wowsims/tbc/sim/common"
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/shaman"
)

func (enh *EnhancementShaman) SetupRotationSchedule() {
	// Fill the GCD schedule based on our settings.
	maxDuration := enh.Env.GetMaxDuration()

	var curTime time.Duration

	if enh.Talents.Stormstrike {
		ssAction := common.ScheduledAbility{
			Duration: core.GCDDefault,
			TryCast: func(sim *core.Simulation) bool {
				ss := enh.Stormstrike
				success := ss.Cast(sim, enh.CurrentTarget)
				if !success {
					enh.WaitForMana(sim, ss.CurCast.Cost)
				}
				return success
			},
		}
		curTime = core.DurationFromSeconds(enh.Rotation.FirstStormstrikeDelay)
		for curTime <= maxDuration {
			ability := ssAction
			ability.DesiredCastAt = curTime
			castAt := enh.scheduler.Schedule(ability)
			curTime = castAt + time.Second*10
		}
	}

	shockCD := enh.ShockCD()
	shockAction := common.ScheduledAbility{
		Duration: core.GCDDefault,
		TryCast: func(sim *core.Simulation) bool {
			var shock *core.Spell
			if enh.Rotation.WeaveFlameShock && !enh.FlameShockDot.IsActive() {
				shock = enh.FlameShock
			} else if enh.Rotation.PrimaryShock == proto.EnhancementShaman_Rotation_Earth {
				shock = enh.EarthShock
			} else if enh.Rotation.PrimaryShock == proto.EnhancementShaman_Rotation_Frost {
				shock = enh.FrostShock
			}

			success := shock.Cast(sim, enh.CurrentTarget)
			if !success {
				enh.WaitForMana(sim, shock.CurCast.Cost)
			}
			return success
		},
	}
	if enh.Rotation.PrimaryShock != proto.EnhancementShaman_Rotation_None {
		curTime = 0
		for curTime <= maxDuration {
			ability := shockAction
			ability.DesiredCastAt = curTime
			ability.MinCastAt = curTime
			ability.MaxCastAt = curTime + time.Second*10
			castAt := enh.scheduler.Schedule(ability)
			curTime = castAt + shockCD
		}
	} else if enh.Rotation.WeaveFlameShock {
		// Flame shock but no regular shock, so only use it once every 12s.
		curTime = 0
		for curTime <= maxDuration {
			ability := shockAction
			ability.DesiredCastAt = curTime
			ability.MinCastAt = curTime
			ability.MaxCastAt = curTime + time.Second*10
			castAt := enh.scheduler.Schedule(ability)
			curTime = castAt + time.Second*12
		}
	}

	// We need to directly manage all GCD-bound CDs ourself.
	if enh.Consumes.Drums == proto.Drums_DrumsOfBattle {
		enh.scheduler.ScheduleMCD(enh.GetCharacter(), core.DrumsOfBattleActionID)
	} else if enh.Consumes.Drums == proto.Drums_DrumsOfRestoration {
		enh.scheduler.ScheduleMCD(enh.GetCharacter(), core.DrumsOfRestorationActionID)
	} else if enh.Consumes.Drums == proto.Drums_DrumsOfWar {
		enh.scheduler.ScheduleMCD(enh.GetCharacter(), core.DrumsOfWarActionID)
	}
	if enh.SelfBuffs.Bloodlust {
		enh.scheduler.ScheduleMCD(enh.GetCharacter(), enh.BloodlustActionID())
	}

	scheduleTotem := func(duration time.Duration, prioritizeEarlier bool, precast bool, tryCast func(sim *core.Simulation) (bool, float64)) {
		totemAction := common.ScheduledAbility{
			Duration: time.Second * 1,
			TryCast: func(sim *core.Simulation) bool {
				success, manaCost := tryCast(sim)
				if !success {
					enh.WaitForMana(sim, manaCost)
				}
				return success
			},
			PrioritizeEarlierForConflicts: prioritizeEarlier,
		}

		curTime := time.Duration(0)
		if precast {
			curTime = duration
		}
		for curTime <= maxDuration {
			ability := totemAction
			ability.DesiredCastAt = curTime
			if prioritizeEarlier {
				ability.MinCastAt = curTime - time.Second*30
				ability.MaxCastAt = curTime + time.Second*15
			} else {
				ability.MinCastAt = curTime - time.Second*5
				ability.MaxCastAt = curTime + time.Second*30
			}
			castAt := enh.scheduler.Schedule(ability)
			if castAt == common.Unresolved {
				panic("No timeslot found for totem")
			}
			curTime = castAt + duration
		}
	}
	scheduleSpellTotem := func(duration time.Duration, spell *core.Spell) {
		scheduleTotem(duration, false, false, func(sim *core.Simulation) (bool, float64) {
			success := spell.Cast(sim, enh.CurrentTarget)
			return success, spell.CurCast.Cost
		})
	}
	schedule2MTotem := func(castFactory func(sim *core.Simulation) *core.Spell) {
		scheduleTotem(time.Minute*2, true, true, func(sim *core.Simulation) (bool, float64) {
			spell := castFactory(sim)
			return spell.Cast(sim, enh.CurrentTarget), spell.CurCast.Cost
		})
	}

	if enh.Totems.TwistFireNova {
		var defaultCastFactory func(sim *core.Simulation)
		switch enh.Totems.Fire {
		case proto.FireTotem_MagmaTotem:
			defaultCastFactory = func(sim *core.Simulation) {
				if enh.SearingTotemDot.IsActive() || enh.MagmaTotemDot.IsActive() || enh.FireNovaTotemDot.IsActive() {
					return
				}

				cast := enh.MagmaTotem
				success := cast.Cast(sim, nil)
				if !success {
					enh.WaitForMana(sim, cast.CurCast.Cost)
				}
			}
		case proto.FireTotem_SearingTotem:
			defaultCastFactory = func(sim *core.Simulation) {
				if enh.SearingTotemDot.IsActive() || enh.MagmaTotemDot.IsActive() || enh.FireNovaTotemDot.IsActive() {
					return
				}

				cast := enh.SearingTotem
				success := cast.Cast(sim, enh.CurrentTarget)
				if !success {
					enh.WaitForMana(sim, cast.CurCast.Cost)
				}
			}
		case proto.FireTotem_TotemOfWrath:
			defaultCastFactory = func(sim *core.Simulation) {
				if enh.NextTotemDrops[shaman.FireTotem] > sim.CurrentTime+time.Second*5 {
					// Skip dropping if we've gone OOM reverted to dropping default only, and have plenty of time left.
					return
				}

				cast := enh.TotemOfWrath
				success := cast.Cast(sim, nil)
				if !success {
					enh.WaitForMana(sim, cast.CurCast.Cost)
				}
			}
		}

		fntAction := common.ScheduledAbility{
			Duration: time.Second * 1,
			TryCast: func(sim *core.Simulation) bool {
				if enh.Metrics.WentOOM && enh.CurrentManaPercent() < 0.2 {
					return false
				}

				cast := enh.FireNovaTotem
				success := cast.Cast(sim, nil)
				if !success {
					enh.WaitForMana(sim, cast.CurCast.Cost)
				}
				return success
			},
		}
		defaultAction := common.ScheduledAbility{
			Duration: time.Second * 1,
			TryCast: func(sim *core.Simulation) bool {
				defaultCastFactory(sim)
				return true
			},
		}

		curTime := time.Duration(0)
		nextNovaCD := time.Duration(0)
		defaultNext := false
		for curTime <= maxDuration {
			ability := fntAction
			if defaultNext {
				ability = defaultAction
			}
			ability.DesiredCastAt = curTime
			ability.MinCastAt = curTime
			ability.MaxCastAt = curTime + time.Second*15

			castAt := enh.scheduler.Schedule(ability)

			if defaultNext {
				curTime = nextNovaCD
				defaultNext = false
			} else {
				nextNovaCD = castAt + time.Second*15 + 1
				if defaultCastFactory == nil {
					curTime = nextNovaCD
				} else {
					curTime = castAt + enh.FireNovaTickLength() + 1
					defaultNext = true
				}
			}
		}
	} else {
		switch enh.Totems.Fire {
		case proto.FireTotem_MagmaTotem:
			scheduleSpellTotem(time.Second*20+1, enh.MagmaTotem)
		case proto.FireTotem_SearingTotem:
			scheduleSpellTotem(time.Minute*1+1, enh.SearingTotem)
		case proto.FireTotem_TotemOfWrath:
			schedule2MTotem(func(sim *core.Simulation) *core.Spell { return enh.TotemOfWrath })
		}
	}

	if enh.Totems.Air != proto.AirTotem_NoAirTotem {
		var defaultCastFactory func(sim *core.Simulation) *core.Spell
		switch enh.Totems.Air {
		case proto.AirTotem_GraceOfAirTotem:
			defaultCastFactory = func(sim *core.Simulation) *core.Spell { return enh.GraceOfAirTotem }
		case proto.AirTotem_TranquilAirTotem:
			defaultCastFactory = func(sim *core.Simulation) *core.Spell { return enh.TranquilAirTotem }
		case proto.AirTotem_WindfuryTotem:
			defaultCastFactory = func(sim *core.Simulation) *core.Spell { return enh.WindfuryTotem }
		case proto.AirTotem_WrathOfAirTotem:
			defaultCastFactory = func(sim *core.Simulation) *core.Spell { return enh.WrathOfAirTotem }
		}

		if enh.Totems.TwistWindfury {
			wfAction := common.ScheduledAbility{
				Duration: time.Second * 1,
				TryCast: func(sim *core.Simulation) bool {
					if enh.Metrics.WentOOM && enh.CurrentManaPercent() < 0.2 {
						return false
					}

					cast := enh.WindfuryTotem
					success := cast.Cast(sim, nil)
					if !success {
						enh.WaitForMana(sim, cast.CurCast.Cost)
					}
					return success
				},
				PrioritizeEarlierForConflicts: true,
			}
			defaultAction := common.ScheduledAbility{
				Duration: time.Second * 1,
				TryCast: func(sim *core.Simulation) bool {
					if enh.NextTotemDrops[shaman.AirTotem] > sim.CurrentTime+time.Second*10 {
						// Skip dropping if we've gone OOM reverted to dropping default only, and have plenty of time left.
						return true
					}

					cast := defaultCastFactory(sim)
					success := cast.Cast(sim, enh.CurrentTarget)
					if !success {
						enh.WaitForMana(sim, cast.CurCast.Cost)
					}
					return success
				},
			}

			curTime := time.Second * 10
			for curTime <= maxDuration {
				ability := wfAction
				ability.DesiredCastAt = curTime
				ability.MinCastAt = curTime - time.Second*8
				ability.MaxCastAt = curTime + time.Second*20
				defaultAbility := defaultAction
				castAt := enh.scheduler.ScheduleGroup([]common.ScheduledAbility{ability, defaultAbility})
				if castAt == common.Unresolved {
					panic(fmt.Sprintf("No timeslot found for air totem, desired: %s", curTime))
				}
				curTime = castAt + time.Second*10
			}
		} else {
			schedule2MTotem(defaultCastFactory)
		}
	}

	if enh.Totems.Earth != proto.EarthTotem_NoEarthTotem {
		switch enh.Totems.Earth {
		case proto.EarthTotem_StrengthOfEarthTotem:
			schedule2MTotem(func(sim *core.Simulation) *core.Spell { return enh.StrengthOfEarthTotem })
		case proto.EarthTotem_TremorTotem:
			schedule2MTotem(func(sim *core.Simulation) *core.Spell { return enh.TremorTotem })
		}
	}

	if enh.Totems.Water != proto.WaterTotem_NoWaterTotem {
		if enh.Totems.Water == proto.WaterTotem_ManaSpringTotem {
			schedule2MTotem(func(sim *core.Simulation) *core.Spell { return enh.ManaSpringTotem })
		}
	}
}

func (enh *EnhancementShaman) OnGCDReady(sim *core.Simulation) {
	enh.scheduler.DoNextAbility(sim, &enh.Character)
}

func (enh *EnhancementShaman) OnManaTick(sim *core.Simulation) {
	if enh.IsWaitingForMana() && !enh.DoneWaitingForMana(sim) {
		// Do nothing, just need to check so metrics get updated.
	}
}
