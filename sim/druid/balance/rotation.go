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

	var spell *core.Spell

	spell = moonkin.rotation(sim)

	if success := spell.Cast(sim, moonkin.CurrentTarget); !success {
		moonkin.WaitForMana(sim, spell.CurCast.Cost)
	}
}

func (moonkin *BalanceDruid) rotation(sim *core.Simulation) *core.Spell {

	rotation := &moonkin.Rotation
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
		canUseCooldownsInLunar := lunarUptime.Seconds() >= float64(rotation.McdInsideLunarThreshold)-0.5 && rotation.UseSmartCooldowns
		canUseCooldownsInSolar := solarUptime.Seconds() >= float64(rotation.McdInsideSolarThreshold)-0.5 && rotation.UseSmartCooldowns

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
					moonkin.castMajorCooldown(moonkin.hyperSpeedMCD, sim, target)
					moonkin.castMajorCooldown(moonkin.potionSpeedMCD, sim, target)
					moonkin.castMajorCooldown(moonkin.onUseTrinket1, sim, target)
					moonkin.castMajorCooldown(moonkin.onUseTrinket2, sim, target)
				}
				if moonfireUptime > 0 || float64(rotation.MfInsideEclipseThreshold) >= lunarUptime.Seconds() {
					return moonkin.Starfire
				} else if rotation.UseMf {
					return moonkin.Moonfire
				}
			} else {
				if canUseCooldownsInSolar {
					moonkin.castMajorCooldown(moonkin.potionWildMagicMCD, sim, target)
					moonkin.castMajorCooldown(moonkin.onUseTrinket1, sim, target)
					moonkin.castMajorCooldown(moonkin.onUseTrinket2, sim, target)
				}
				if insectSwarmUptime > 0 || float64(rotation.IsInsideEclipseThreshold) >= solarUptime.Seconds() {
					return moonkin.Wrath
				} else if rotation.UseIs {
					return moonkin.InsectSwarm
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
			return moonkin.Rebirth
		} else if moonkin.Starfall.IsReady(sim) {
			return moonkin.Starfall
		} else if moonkin.Talents.ForceOfNature && moonkin.ForceOfNature.IsReady(sim) {
			return moonkin.ForceOfNature
		} else if rotation.UseMf && moonfireUptime <= 0 && (fishingForLunar || rotation.KeepMfUp) {
			return moonkin.Moonfire
		} else if rotation.UseIs && insectSwarmUptime <= 0 && (fishingForSolar || rotation.KeepIsUp) {
			return moonkin.InsectSwarm
		} else if fishingForLunar {
			return moonkin.Wrath
		} else {
			return moonkin.Starfire
		}
	}
	return moonkin.Starfire
}

func (moonkin *BalanceDruid) castMajorCooldown(mcd *core.MajorCooldown, sim *core.Simulation, target *core.Unit) {
	if mcd != nil && mcd.Spell.IsReady(sim) {
		isOffensivePotion := mcd.Spell.SameAction(core.ActionID{ItemID: 40211}) || mcd.Spell.SameAction(core.ActionID{ItemID: 40212})
		willUseOffensivePotion := isOffensivePotion && !moonkin.potionUsed

		// Use Potion if we can
		if isOffensivePotion && moonkin.potionUsed {
			return
		}
		mcd.Spell.Cast(sim, target)
		if willUseOffensivePotion {
			moonkin.potionUsed = true
		}
	}
}
