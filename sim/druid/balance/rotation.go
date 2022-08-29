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
	var spell *core.Spell

	moonfireUptime := moonkin.MoonfireDot.RemainingDuration(sim)
	insectSwarmUptime := moonkin.InsectSwarmDot.RemainingDuration(sim)
	shouldRebirth := sim.GetRemainingDuration().Seconds() < moonkin.RebirthTiming

	lunarICD := moonkin.LunarICD.Timer.TimeToReady(sim)
	solarICD := moonkin.SolarICD.Timer.TimeToReady(sim)
	fishingForLunar := lunarICD <= solarICD
	fishingForSolar := solarICD < lunarICD

	if moonkin.Talents.Eclipse > 0 {
		// Eclipse stuff
		lunarIsActive := lunarICD > time.Millisecond*15000
		solarIsActive := solarICD > time.Millisecond*15000
		lunarUptime := core.TernaryDuration(lunarIsActive, lunarICD-time.Millisecond*15000, 0)
		solarUptime := core.TernaryDuration(solarIsActive, solarICD-time.Millisecond*15000, 0)

		// "Dispelling" eclipse effects before casting if needed
		if float64(lunarUptime-moonkin.Starfire.CurCast.CastTime) <= 0 && moonkin.useIS {
			moonkin.GetAura("Lunar Eclipse proc").Deactivate(sim)
			lunarIsActive = false
		}
		if float64(solarUptime-moonkin.Wrath.CurCast.CastTime) <= 0 && moonkin.useMF {
			moonkin.GetAura("Solar Eclipse proc").Deactivate(sim)
			solarIsActive = false
		}

		// Eclipse
		if solarIsActive || lunarIsActive {
			if lunarIsActive {
				if lunarUptime.Seconds() >= float64(moonkin.mcdInsideLunarThreshold) && moonkin.mcdInsideLunarThreshold >= 0 {
					moonkin.castAllMajorCooldowns(sim)
				}
				if moonfireUptime > 0 || float64(moonkin.mfInsideEclipseThreshold) >= lunarUptime.Seconds() {
					spell = moonkin.Starfire
				} else if moonkin.useMF {
					spell = moonkin.Moonfire
				}
			} else {
				if solarUptime.Seconds() >= float64(moonkin.mcdInsideSolarThreshold) && moonkin.mcdInsideSolarThreshold >= 0 {
					moonkin.castAllMajorCooldowns(sim)
				}
				if insectSwarmUptime > 0 || float64(moonkin.isInsideEclipseThreshold) >= solarUptime.Seconds() {
					spell = moonkin.Wrath
				} else if moonkin.useIS {
					spell = moonkin.InsectSwarm
				}
			}
		}
	} else {
		fishingForLunar, fishingForSolar = true, true // If Eclipse isn't talented we're not fishing
	}

	// Non-Eclipse
	if spell == nil {
		// We're not gonna rez someone during eclipse, are we ?
		if moonkin.useBattleRes && shouldRebirth && moonkin.Rebirth.IsReady(sim) {
			spell = moonkin.Rebirth
		} else if moonkin.Starfall.IsReady(sim) {
			spell = moonkin.Starfall
		} else if moonkin.Talents.ForceOfNature && moonkin.ForceOfNature.IsReady(sim) {
			spell = moonkin.ForceOfNature
		} else if moonkin.useMF && moonfireUptime <= 0 && fishingForLunar {
			spell = moonkin.Moonfire
		} else if moonkin.useIS && insectSwarmUptime <= 0 && fishingForSolar {
			spell = moonkin.InsectSwarm
		} else if fishingForLunar {
			spell = moonkin.Wrath
		} else {
			spell = moonkin.Starfire
		}
	}

	if success := spell.Cast(sim, target); !success {
		moonkin.WaitForMana(sim, spell.CurCast.Cost)
	}
}

func (moonkin *BalanceDruid) castMajorCooldown(mcd *core.MajorCooldown, sim *core.Simulation, target *core.Unit) {
	if mcd != nil {
		if mcd.Spell.IsReady(sim) && moonkin.GCD.IsReady(sim) {
			mcd.Spell.Cast(sim, target)
			if mcd.Spell.SameAction(core.ActionID{ItemID: 40211}) || mcd.Spell.SameAction(core.ActionID{ItemID: 40212}) {
				moonkin.potionUsed = true
			}
		}
	}
}

func (moonkin *BalanceDruid) castAllMajorCooldowns(sim *core.Simulation) {
	target := moonkin.CurrentTarget
	moonkin.castMajorCooldown(moonkin.hyperSpeedMCD, sim, target)
	if !moonkin.potionUsed {
		moonkin.castMajorCooldown(moonkin.potionMCD, sim, target)
	}
	moonkin.castMajorCooldown(moonkin.onUseTrinket1, sim, target)
	moonkin.castMajorCooldown(moonkin.onUseTrinket2, sim, target)
}
