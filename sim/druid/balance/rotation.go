package balance

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (moonkin *BalanceDruid) OnGCDReady(sim *core.Simulation) {
	moonkin.tryUseGCD(sim)
}

func (moonkin *BalanceDruid) tryUseGCD(sim *core.Simulation) {
	spell, target := moonkin.rotation(sim)
	if success := spell.Cast(sim, target); !success {
		moonkin.WaitForMana(sim, spell.CurCast.Cost)
	}
	moonkin.LastCast = spell
}

func (moonkin *BalanceDruid) rotation(sim *core.Simulation) (*core.Spell, *core.Unit) {
	moonkin.CurrentTarget = sim.Environment.GetTargetUnit(0)
	rotation := moonkin.Rotation
	target := moonkin.CurrentTarget

	if rotation.MaintainFaerieFire && moonkin.ShouldFaerieFire(sim, moonkin.CurrentTarget) {
		return moonkin.FaerieFire, target
	}

	var lunarUptime time.Duration
	if moonkin.LunarEclipseProcAura != nil {
		lunarUptime = moonkin.LunarEclipseProcAura.RemainingDuration(sim)
	}

	if moonkin.MoonkinT84PCAura.IsActive() && moonkin.MoonkinT84PCAura.RemainingDuration(sim) < moonkin.SpellGCD() {
		if (rotation.UseSmartCooldowns && lunarUptime > 14*time.Second) || sim.GetRemainingDuration() < 15*time.Second {
			moonkin.castMajorCooldown(moonkin.hyperSpeedMCD, sim, target)
			moonkin.castMajorCooldown(moonkin.potionSpeedMCD, sim, target)
			moonkin.useTrinkets(stats.SpellHaste, sim, target)
		}
		return moonkin.Starfire, target
	} else if rotation.UseBattleRes && sim.GetRemainingDuration().Seconds() < moonkin.RebirthTiming && moonkin.Rebirth.IsReady(sim) {
		return moonkin.Rebirth, target
	} else if moonkin.Talents.ForceOfNature && moonkin.ForceOfNature.IsReady(sim) {
		moonkin.useTrinkets(stats.SpellPower, sim, target)
		return moonkin.ForceOfNature, target
	} else if moonkin.Starfall.IsReady(sim) {
		moonkin.useTrinkets(stats.SpellPower, sim, target)
		return moonkin.Starfall, target
	} else if moonkin.Typhoon.IsReady(sim) && rotation.UseTyphoon {
		return moonkin.Typhoon, target
	} else if rotation.UseHurricane {
		return moonkin.Hurricane, target
	}

	lunarIsActive := moonkin.LunarEclipseProcAura.IsActive()
	shouldHoldIs := false
	if lunarIsActive && moonkin.MoonkinT84PCAura == nil {
		shouldHoldIs = lunarUptime.Seconds() < (moonkin.InsectSwarm.DamageMultiplier-1)/0.042
	}

	// Max IS uptime
	if rotation.IsUsage == proto.BalanceDruid_Rotation_MaximizeIs && !shouldHoldIs {
		if moonkin.InsectSwarm.CurDot().RemainingDuration(sim) <= 0 {
			return moonkin.InsectSwarm, target
		}
	} else if rotation.IsUsage == proto.BalanceDruid_Rotation_MultidotIs {
		targetIndex := 0
		for targetIndex < len(sim.Encounter.Targets) {
			if moonkin.InsectSwarm.CurDot().RemainingDuration(sim) <= 0 {
				return moonkin.InsectSwarm, moonkin.CurrentTarget
			}
			moonkin.CurrentTarget = sim.Environment.NextTargetUnit(moonkin.CurrentTarget)
			targetIndex += 1
		}
	}

	// Max MF uptime
	shouldRefreshMf := moonkin.Moonfire.CurDot().RemainingDuration(sim) <= 0
	if rotation.MfUsage == proto.BalanceDruid_Rotation_MaximizeMf && shouldRefreshMf {
		return moonkin.Moonfire, target
	}

	// Player "brain" latency
	playerLatency := time.Duration(core.MaxInt32(rotation.PlayerLatency, 0)) * time.Millisecond
	lunarICD := moonkin.LunarICD.Timer.TimeToReady(sim)
	solarICD := moonkin.SolarICD.Timer.TimeToReady(sim)

	if moonkin.Talents.Eclipse > 0 {
		solarUptime := moonkin.SolarEclipseProcAura.ExpiresAt() - sim.CurrentTime
		solarIsActive := moonkin.SolarEclipseProcAura.IsActive()

		//"Dispelling" eclipse effects before casting if needed
		if float64(lunarUptime-moonkin.Starfire.CurCast.CastTime) <= 0 {
			moonkin.LunarEclipseProcAura.Deactivate(sim)
			lunarIsActive = false
		}
		if float64(solarUptime-moonkin.Wrath.CurCast.CastTime) <= 0 {
			moonkin.SolarEclipseProcAura.Deactivate(sim)
			solarIsActive = false
		}
		// Player latency adjustments
		if lunarIsActive {
			lunarIsActive = lunarUptime < (moonkin.LunarEclipseProcAura.Duration - playerLatency)
		}
		if solarIsActive {
			solarIsActive = solarUptime < (moonkin.SolarEclipseProcAura.Duration - playerLatency)
			solarICD = 0
		}

		// Eclipse
		if solarIsActive || lunarIsActive {
			if lunarIsActive {
				if (rotation.UseSmartCooldowns && lunarUptime > 10*time.Second) || sim.GetRemainingDuration() < 15*time.Second {
					moonkin.castMajorCooldown(moonkin.hyperSpeedMCD, sim, target)
					moonkin.castMajorCooldown(moonkin.potionSpeedMCD, sim, target)
					moonkin.useTrinkets(stats.SpellHaste, sim, target)
					if !moonkin.HasActiveAuraWithTag(core.BloodlustAuraTag) {
						moonkin.castMajorCooldown(moonkin.powerInfusion, sim, target)
					}
				}
				return moonkin.Starfire, target
			} else if solarIsActive && rotation.WrathUsage == proto.BalanceDruid_Rotation_RegularWrath {
				if moonkin.MoonkinT84PCAura.IsActive() {
					if moonkin.MoonkinT84PCAura.RemainingDuration(sim) < solarUptime {
						return moonkin.Starfire, target
					}
				}
				if (rotation.UseSmartCooldowns && solarUptime > 10*time.Second) || sim.GetRemainingDuration() < 15*time.Second {
					moonkin.castMajorCooldown(moonkin.potionWildMagicMCD, sim, target)
					moonkin.useTrinkets(stats.SpellCrit, sim, target)
				}
				return moonkin.Wrath, target
			}
		}
		if rotation.MfUsage == proto.BalanceDruid_Rotation_BeforeLunar && lunarICD < 2*time.Second && shouldRefreshMf {
			return moonkin.Moonfire, target
		}
		shouldRefreshIs := moonkin.InsectSwarm.CurDot().RemainingDuration(sim) <= 0
		if rotation.IsUsage == proto.BalanceDruid_Rotation_BeforeSolar && solarICD < 2*time.Second && shouldRefreshIs {
			return moonkin.InsectSwarm, target
		}
	}

	fishingForLunar := lunarICD <= solarICD
	if rotation.EclipsePrio == proto.BalanceDruid_Rotation_Solar {
		fishingForLunar = lunarICD < solarICD
	}
	// Non-Eclipse
	eclipseShuffle := rotation.EclipseShuffling && lunarICD == 0 && solarICD == 0
	if eclipseShuffle && moonkin.LastCast == moonkin.Wrath && rotation.UseStarfire {
		return moonkin.Starfire, target
	}
	if (fishingForLunar || eclipseShuffle) && rotation.WrathUsage != proto.BalanceDruid_Rotation_NoWrath {
		return moonkin.Wrath, target
	} else {
		return moonkin.Starfire, target
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

func (moonkin *BalanceDruid) useTrinkets(stat stats.Stat, sim *core.Simulation, target *core.Unit) {
	if moonkin.onUseTrinket1.Stat == stat {
		moonkin.castMajorCooldown(moonkin.onUseTrinket1.Cooldown, sim, target)
	}
	if moonkin.onUseTrinket2.Stat == stat {
		moonkin.castMajorCooldown(moonkin.onUseTrinket2.Cooldown, sim, target)
	}
}
