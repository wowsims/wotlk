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

	lunarICD := moonkin.LunarICD.Timer.TimeToReady(sim)
	solarICD := moonkin.SolarICD.Timer.TimeToReady(sim)

	lunarIsActive := lunarICD > time.Millisecond*15000
	solarIsActive := solarICD > time.Millisecond*15000

	//lunarUptime := core.TernaryDuration(lunarIsActive, lunarICD-time.Millisecond*15000, 0)
	//solarUptime := core.TernaryDuration(solarIsActive, solarICD-time.Millisecond*15000, 0)

	// TODO Store those original values elsewhere and use them in the spells aswell
	originalWrathDamageMultiplier := (1 + 0.02*float64(moonkin.Talents.Moonfury)) * (1 + 0.01*float64(moonkin.Talents.ImprovedInsectSwarm))
	originalStarfireBonusCritRating := float64(2 * float64(moonkin.Talents.NaturesMajesty) * 45.91)
	moonkin.Wrath.DamageMultiplier = core.TernaryFloat64(solarIsActive, originalWrathDamageMultiplier+0.4, originalWrathDamageMultiplier)
	moonkin.Starfire.BonusCritRating = core.TernaryFloat64(lunarIsActive, originalStarfireBonusCritRating+40*core.CritRatingPerCritChance, originalStarfireBonusCritRating)

	moonfireUptime := moonkin.MoonfireDot.RemainingDuration(sim)
	insectSwarmUptime := moonkin.InsectSwarmDot.RemainingDuration(sim)

	shouldRebirth := sim.GetRemainingDuration().Seconds() < moonkin.RebirthTiming

	var spell *core.Spell
	// TODO Treants
	if moonkin.useBattleRes && shouldRebirth && moonkin.Rebirth.IsReady(sim) {
		spell = moonkin.Rebirth
	} else if moonkin.Starfall.IsReady(sim) {
		spell = moonkin.Starfall
	} else if (solarIsActive && (insectSwarmUptime > 0 || !moonkin.canIsInsideEclipse)) || (!lunarIsActive && moonfireUptime > 13) {
		spell = moonkin.Wrath
	} else if (lunarIsActive && (moonfireUptime > 0 || !moonkin.canMfInsideEclipse)) || (!solarIsActive && insectSwarmUptime > 13) {
		spell = moonkin.Starfire
	} else if (lunarIsActive || lunarICD < core.GCDDefault) && moonkin.useMF {
		spell = moonkin.Moonfire
	} else if moonkin.useIS {
		spell = moonkin.InsectSwarm
	} else {
		spell = moonkin.Starfire // Always fallback to Starfire for beautiful Classic memories
	}

	if success := spell.Cast(sim, target); !success {
		moonkin.WaitForMana(sim, spell.CurCast.Cost)
	}
}
