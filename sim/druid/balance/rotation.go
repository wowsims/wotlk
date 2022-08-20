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

	// "Applying" or "Dispelling" eclipse effects before casting
	solarShouldStayActive := float64(solarUptime-spell.DefaultCast.CastTime) > 0
	lunarShouldStayActive := float64(lunarUptime-spell.DefaultCast.CastTime) > 0
	moonkin.Wrath.DamageMultiplier = core.TernaryFloat64(solarShouldStayActive, moonkin.OriginalWrathDamageMultiplier+0.4, moonkin.OriginalWrathDamageMultiplier)
	moonkin.Starfire.BonusCritRating = core.TernaryFloat64(lunarShouldStayActive, moonkin.OriginalStarfireBonusCritRating+(40*core.CritRatingPerCritChance), moonkin.OriginalStarfireBonusCritRating)

	//spell.BonusCritRating += core.TernaryFloat64(target.HasActiveAura("Faerie Fire"), float64(moonkin.Talents.ImprovedFaerieFire)*1*core.CritRatingPerCritChance, 0)

	if success := spell.Cast(sim, target); !success {
		moonkin.WaitForMana(sim, spell.CurCast.Cost)
	}
}
