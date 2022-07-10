package enhancement

import (
	"time"

	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
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
			if !enh.FlameShockDot.IsActive() {
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
	schedule5MTotem := func(castFactory func(sim *core.Simulation) *core.Spell) {
		scheduleTotem(time.Minute*5, true, true, func(sim *core.Simulation) (bool, float64) {
			spell := castFactory(sim)
			return spell.Cast(sim, enh.CurrentTarget), spell.CurCast.Cost
		})
	}

	switch enh.Totems.Fire {
	case proto.FireTotem_MagmaTotem:
		scheduleSpellTotem(time.Second*20+1, enh.MagmaTotem)
	case proto.FireTotem_SearingTotem:
		scheduleSpellTotem(time.Minute*1+1, enh.SearingTotem)
	case proto.FireTotem_TotemOfWrath:
		schedule5MTotem(func(sim *core.Simulation) *core.Spell { return enh.TotemOfWrath })
	}

	if enh.Totems.Air != proto.AirTotem_NoAirTotem {
		var defaultCastFactory func(sim *core.Simulation) *core.Spell
		switch enh.Totems.Air {
		case proto.AirTotem_TranquilAirTotem:
			defaultCastFactory = func(sim *core.Simulation) *core.Spell { return enh.TranquilAirTotem }
		case proto.AirTotem_WindfuryTotem:
			defaultCastFactory = func(sim *core.Simulation) *core.Spell { return enh.WindfuryTotem }
		case proto.AirTotem_WrathOfAirTotem:
			defaultCastFactory = func(sim *core.Simulation) *core.Spell { return enh.WrathOfAirTotem }
		}

		schedule5MTotem(defaultCastFactory)
	}

	if enh.Totems.Earth != proto.EarthTotem_NoEarthTotem {
		switch enh.Totems.Earth {
		case proto.EarthTotem_StrengthOfEarthTotem:
			schedule5MTotem(func(sim *core.Simulation) *core.Spell { return enh.StrengthOfEarthTotem })
		case proto.EarthTotem_TremorTotem:
			schedule5MTotem(func(sim *core.Simulation) *core.Spell { return enh.TremorTotem })
		}
	}

	if enh.Totems.Water != proto.WaterTotem_NoWaterTotem {
		if enh.Totems.Water == proto.WaterTotem_ManaSpringTotem {
			schedule5MTotem(func(sim *core.Simulation) *core.Spell { return enh.ManaSpringTotem })
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
