package balance

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (moonkin *BalanceDruid) OnGCDReady(sim *core.Simulation) {
	moonkin.tryUseGCD(sim)
}

func (moonkin *BalanceDruid) tryUseGCD(sim *core.Simulation) {
	spell := moonkin.rotation(sim)

	if success := spell.Cast(sim, moonkin.CurrentTarget); !success {
		moonkin.WaitForMana(sim, spell.CurCast.Cost)
	}
}

func (moonkin *BalanceDruid) rotation(sim *core.Simulation) *core.Spell {
	rotation := moonkin.Rotation
	target := moonkin.CurrentTarget

	if rotation.MaintainFaerieFire && moonkin.ShouldFaerieFire(sim) {
		if moonkin.Talents.ImprovedFaerieFire > 0 {
			if aura := target.GetActiveAuraWithTag(core.MinorSpellHitDebuffAuraTag); aura == nil {
				return moonkin.FaerieFire
			}
		}
	}

	moonfireUptime := moonkin.MoonfireDot.RemainingDuration(sim)
	insectSwarmUptime := moonkin.InsectSwarmDot.RemainingDuration(sim)
	shouldRebirth := sim.GetRemainingDuration().Seconds() < moonkin.RebirthTiming

	// Player "brain" latency
	playerLatency := time.Duration(rotation.PlayerLatency)
	lunarICD := moonkin.LunarICD.Timer.TimeToReady(sim)
	solarICD := moonkin.SolarICD.Timer.TimeToReady(sim)
	fishingForLunar := lunarICD <= solarICD
	fishingForSolar := solarICD < lunarICD
	maximizeIsUptime := rotation.MaximizeIsUptime && rotation.UseIs
	maximizeMfUptime := rotation.MaximizeMfUptime && rotation.UseMf

	if rotation.UseBattleRes && shouldRebirth && moonkin.Rebirth.IsReady(sim) {
		return moonkin.Rebirth
	} else if moonkin.Talents.ForceOfNature && moonkin.ForceOfNature.IsReady(sim) {
		moonkin.castMajorCooldown(moonkin.onUseTrinket1, sim, target)
		moonkin.castMajorCooldown(moonkin.onUseTrinket2, sim, target)
		return moonkin.ForceOfNature
	} else if moonkin.Starfall.IsReady(sim) {
		moonkin.castMajorCooldown(moonkin.onUseTrinket1, sim, target)
		moonkin.castMajorCooldown(moonkin.onUseTrinket2, sim, target)
		return moonkin.Starfall
	} else if moonkin.Typhoon.IsReady(sim) && rotation.UseTyphoon {
		return moonkin.Typhoon
	} else if rotation.UseHurricane {
		return moonkin.Hurricane
	}

	if moonkin.Talents.Eclipse > 0 {

		lunarUptime := moonkin.LunarEclipseProcAura.ExpiresAt() - sim.CurrentTime
		solarUptime := moonkin.SolarEclipseProcAura.ExpiresAt() - sim.CurrentTime
		lunarIsActive := moonkin.LunarEclipseProcAura.IsActive()
		solarIsActive := moonkin.SolarEclipseProcAura.IsActive()

		if lunarIsActive {
			lunarIsActive = lunarUptime < (moonkin.LunarEclipseProcAura.Duration - playerLatency)
			fishingForSolar = false
		}
		if solarIsActive {
			solarIsActive = solarUptime < (moonkin.SolarEclipseProcAura.Duration - playerLatency)
			fishingForLunar = false
		}

		// "Dispelling" eclipse effects before casting if needed
		if float64(lunarUptime-moonkin.Starfire.CurCast.CastTime) <= 0 && rotation.UseMf {
			lunarIsActive = false
		}
		if float64(solarUptime-moonkin.Wrath.CurCast.CastTime) <= 0 && rotation.UseIs {
			solarIsActive = false
		}

		// Eclipse
		if solarIsActive || lunarIsActive {
			if maximizeIsUptime && insectSwarmUptime <= 0 {
				return moonkin.InsectSwarm
			}
			if maximizeMfUptime && moonfireUptime <= 0 {
				return moonkin.Moonfire
			}
			if lunarIsActive {
				if (moonfireUptime > 0 || float64(rotation.MfInsideEclipseThreshold) >= lunarUptime.Seconds()) && rotation.UseStarfire {
					if (rotation.UseSmartCooldowns && lunarUptime > 14*time.Second) || sim.GetRemainingDuration() < 15*time.Second {
						moonkin.castMajorCooldown(moonkin.hyperSpeedMCD, sim, target)
						moonkin.castMajorCooldown(moonkin.potionSpeedMCD, sim, target)
					}
					return moonkin.Starfire
				} else if rotation.UseMf {
					return moonkin.Moonfire
				}
			} else if solarIsActive {
				if insectSwarmUptime > 0 || float64(rotation.IsInsideEclipseThreshold) >= solarUptime.Seconds() && rotation.UseWrath {
					if (rotation.UseSmartCooldowns && solarUptime > 14*time.Second) || sim.GetRemainingDuration() < 15*time.Second {
						moonkin.castMajorCooldown(moonkin.potionWildMagicMCD, sim, target)
					}
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
	if rotation.UseMf && moonfireUptime <= 0 && (fishingForLunar || maximizeMfUptime) {
		return moonkin.Moonfire
	} else if rotation.UseIs && insectSwarmUptime <= 0 && (fishingForSolar || maximizeIsUptime) {
		return moonkin.InsectSwarm
	} else if fishingForLunar && rotation.UseWrath {
		return moonkin.Wrath
	} else {
		return moonkin.Starfire
	}
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
		moonkin.UpdateMajorCooldowns()
	}
}
