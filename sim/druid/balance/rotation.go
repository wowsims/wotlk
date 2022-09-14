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
	rotation := &moonkin.Rotation
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
		canUseCooldownsInLunar := lunarUptime.Seconds() >= float64(rotation.McdInsideLunarThreshold)-0.5 && rotation.McdInsideLunarThreshold > 0
		canUseCooldownsInSolar := solarUptime.Seconds() >= float64(rotation.McdInsideSolarThreshold)-0.5 && rotation.McdInsideSolarThreshold > 0

		// "Dispelling" eclipse effects before casting if needed
		if float64(lunarUptime-moonkin.Starfire.CurCast.CastTime) <= 0 && rotation.UseMf {
			moonkin.GetAura("Lunar Eclipse proc").Deactivate(sim)
			lunarIsActive = false
		}
		if float64(solarUptime-moonkin.Wrath.CurCast.CastTime) <= 0 && rotation.UseIs {
			moonkin.GetAura("Solar Eclipse proc").Deactivate(sim)
			solarIsActive = false
		}

		// Eclipse
		if solarIsActive || lunarIsActive {
			if lunarIsActive {
				if canUseCooldownsInLunar {
					moonkin.castAllMajorCooldowns(sim)
				}
				if moonfireUptime > 0 || float64(rotation.MfInsideEclipseThreshold) >= lunarUptime.Seconds() {
					spell = moonkin.Starfire
				} else if rotation.UseMf {
					spell = moonkin.Moonfire
				}
			} else {
				if canUseCooldownsInSolar {
					moonkin.castAllMajorCooldowns(sim)
				}
				if insectSwarmUptime > 0 || float64(rotation.IsInsideEclipseThreshold) >= solarUptime.Seconds() {
					spell = moonkin.Wrath
				} else if rotation.UseIs {
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
		if rotation.UseBattleRes && shouldRebirth && moonkin.Rebirth.IsReady(sim) {
			spell = moonkin.Rebirth
		} else if moonkin.Starfall.IsReady(sim) {
			spell = moonkin.Starfall
		} else if moonkin.Talents.ForceOfNature && moonkin.ForceOfNature.IsReady(sim) {
			spell = moonkin.ForceOfNature
		} else if rotation.UseMf && moonfireUptime <= 0 && fishingForLunar {
			spell = moonkin.Moonfire
		} else if rotation.UseIs && insectSwarmUptime <= 0 && fishingForSolar {
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

func (moonkin *BalanceDruid) castAllMajorCooldowns(sim *core.Simulation) {
	target := moonkin.CurrentTarget
	moonkin.castMajorCooldown(moonkin.hyperSpeedMCD, sim, target)
	moonkin.castMajorCooldown(moonkin.potionMCD, sim, target)
	moonkin.castMajorCooldown(moonkin.onUseTrinket1, sim, target)
	moonkin.castMajorCooldown(moonkin.onUseTrinket2, sim, target)
}

func (moonkin *BalanceDruid) castMajorCooldown(mcd *core.MajorCooldown, sim *core.Simulation, target *core.Unit) {
	if mcd != nil {
		isOffensivePotion := mcd.Spell.SameAction(core.ActionID{ItemID: 40211}) || mcd.Spell.SameAction(core.ActionID{ItemID: 40212})
		shouldUseOffensivePotion := isOffensivePotion && !moonkin.potionUsed

		if isOffensivePotion && moonkin.potionUsed {
			return
		}

		if mcd.Spell.IsReady(sim) && moonkin.GCD.IsReady(sim) {
			mcd.Spell.Cast(sim, target)
			if shouldUseOffensivePotion {
				moonkin.potionUsed = true
			}
		}
	}
}
