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

	target := moonkin.CurrentTarget

	// Eclipse stuff
	lunarICD := moonkin.LunarICD.Timer.TimeToReady(sim)
	solarICD := moonkin.SolarICD.Timer.TimeToReady(sim)
	lunarIsActive := lunarICD > time.Millisecond*15000
	solarIsActive := solarICD > time.Millisecond*15000
	lunarUptime := core.TernaryDuration(lunarIsActive, lunarICD-time.Millisecond*15000, 0)
	solarUptime := core.TernaryDuration(solarIsActive, solarICD-time.Millisecond*15000, 0)

	moonfireUptime := moonkin.MoonfireDot.RemainingDuration(sim)
	insectSwarmUptime := moonkin.InsectSwarmDot.RemainingDuration(sim)

	shouldRebirth := sim.GetRemainingDuration().Seconds() < moonkin.RebirthTiming

	// "Dispelling" eclipse effects before casting if needed
	if float64(lunarUptime-moonkin.Starfire.CurCast.CastTime) <= 0 && moonkin.useIS {
		lunarIsActive = false
	}
	if float64(solarUptime-moonkin.Wrath.CurCast.CastTime) <= 0 && moonkin.useMF {
		solarIsActive = false
	}

	var spell *core.Spell
	// TODO Treants
	if moonkin.useBattleRes && shouldRebirth && moonkin.Rebirth.IsReady(sim) {
		spell = moonkin.Rebirth
	} else if moonkin.Starfall.IsReady(sim) {
		spell = moonkin.Starfall
	} else if (solarIsActive && (insectSwarmUptime > 0 || float64(moonkin.isInsideEclipseThreshold) >= solarUptime.Seconds())) || (!lunarIsActive && moonfireUptime > 13) {
		spell = moonkin.Wrath
	} else if (lunarIsActive && (moonfireUptime > 0 || float64(moonkin.mfInsideEclipseThreshold) >= lunarUptime.Seconds())) || (!solarIsActive && insectSwarmUptime > 13) {
		spell = moonkin.Starfire
	} else if (lunarIsActive || lunarICD < core.GCDDefault) && moonkin.useMF {
		spell = moonkin.Moonfire
	} else if moonkin.useIS {
		spell = moonkin.InsectSwarm
	} else {
		spell = moonkin.Wrath // Always fallback to Wrath to trigger Lunar, because yes
	}

	if success := spell.Cast(sim, target); !success {
		moonkin.WaitForMana(sim, spell.CurCast.Cost)
	}
}
