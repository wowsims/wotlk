package balance

import (
	"github.com/wowsims/wotlk/sim/core"
	"time"
)

func (moonkin *BalanceDruid) OnGCDReady(sim *core.Simulation) {
	moonkin.tryUseGCD(sim)
}

func (moonkin *BalanceDruid) tryUseGCD(sim *core.Simulation) {
	// TODO add rotation choice here
	moonkin.rotation(sim)
}

func (moonkin *BalanceDruid) rotation(sim *core.Simulation) {
	// Activate shared druid behaviors
	// Use Rebirth at the beginning of the fight if flagged in rotation settings
	// Potentially allow options for "Time of cast" in future or default cast like 1 min into fight
	// Currently just casts at the beginning of encounter (with all CDs popped)
	if moonkin.useBattleRes && moonkin.TryRebirth(sim) {
		return
	}

	target := moonkin.CurrentTarget

	lunarICD := moonkin.LunarICD.Timer.TimeToReady(sim)
	solarICD := moonkin.SolarICD.Timer.TimeToReady(sim)

	lunarIsActive := lunarICD > time.Millisecond*15000
	solarIsActive := solarICD > time.Millisecond*15000

	lunarUptime := core.TernaryDuration(lunarIsActive, lunarICD-time.Millisecond*15000, 0)
	solarUptime := core.TernaryDuration(solarIsActive, solarICD-time.Millisecond*15000, 0)

	//TODO These temp stats will be based on Cerdiwyn's new spell values
	moonkin.Wrath.DamageMultiplier = core.TernaryFloat64(solarIsActive, 1+0.02*float64(moonkin.Talents.Moonfury)+0.4, 1+0.02*float64(moonkin.Talents.Moonfury))
	moonkin.Starfire.BonusCritRating = core.TernaryFloat64(lunarIsActive, 40*core.CritRatingPerCritChance, 0)

	moonfireUptime := moonkin.MoonfireDot.RemainingDuration(sim)
	insectSwarmUptime := moonkin.InsectSwarmDot.RemainingDuration(sim)

	var spell *core.Spell
	//TODO Treants
	//TODO Starfall
	if (solarIsActive && insectSwarmUptime > time.Second*3) || (solarIsActive && solarUptime < time.Second*13) || (lunarICD < 2 && moonfireUptime > 0) {
		spell = moonkin.Wrath
	} else if (lunarIsActive && moonfireUptime > time.Second*3) || (lunarIsActive && lunarUptime < time.Second*13) || (solarICD < 2 && insectSwarmUptime > 0) {
		spell = moonkin.Starfire
	} else if lunarIsActive || lunarICD < time.Second*2 {
		spell = moonkin.Moonfire
	} else {
		spell = moonkin.InsectSwarm
	}

	if success := spell.Cast(sim, target); !success {
		moonkin.WaitForMana(sim, spell.CurCast.Cost)
	}
}
