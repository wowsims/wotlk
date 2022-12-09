package balance

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
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

	shouldRebirth := sim.GetRemainingDuration().Seconds() < moonkin.RebirthTiming
	lunarUptime := moonkin.LunarEclipseProcAura.ExpiresAt() - sim.CurrentTime

	if moonkin.HasActiveAura("Elune's Wrath") && moonkin.GetAura("Elune's Wrath").RemainingDuration(sim).Seconds() < moonkin.SpellGCD().Seconds() {
		if (rotation.UseSmartCooldowns && lunarUptime > 14*time.Second) || sim.GetRemainingDuration() < 15*time.Second {
			moonkin.castMajorCooldown(moonkin.hyperSpeedMCD, sim, target)
			moonkin.castMajorCooldown(moonkin.potionSpeedMCD, sim, target)
			moonkin.useTrinkets(stats.SpellHaste, sim, target)
		}
		return moonkin.Starfire
	}
	if rotation.UseBattleRes && shouldRebirth && moonkin.Rebirth.IsReady(sim) {
		return moonkin.Rebirth
	} else if moonkin.Talents.ForceOfNature && moonkin.ForceOfNature.IsReady(sim) {
		moonkin.useTrinkets(stats.SpellPower, sim, target)
		return moonkin.ForceOfNature
	} else if moonkin.Starfall.IsReady(sim) {
		moonkin.useTrinkets(stats.SpellPower, sim, target)
		return moonkin.Starfall
	} else if moonkin.Typhoon.IsReady(sim) && rotation.UseTyphoon {
		return moonkin.Typhoon
	} else if rotation.UseHurricane {
		return moonkin.Hurricane
	}

	moonfireUptime := moonkin.MoonfireDot.RemainingDuration(sim)
	insectSwarmUptime := moonkin.InsectSwarmDot.RemainingDuration(sim)
	// Player "brain" latency
	playerLatency := time.Duration(rotation.PlayerLatency)
	lunarICD := moonkin.LunarICD.Timer.TimeToReady(sim)
	solarICD := moonkin.SolarICD.Timer.TimeToReady(sim)
	fishingForLunar := lunarICD <= solarICD
	//fishingForSolar := solarICD < lunarICD
	useMf := moonkin.Rotation.MfUsage != proto.BalanceDruid_Rotation_NoMf
	useIs := moonkin.Rotation.IsUsage != proto.BalanceDruid_Rotation_NoIs
	maximizeMfUptime := moonkin.Rotation.MfUsage == proto.BalanceDruid_Rotation_MaximizeMf
	maximizeIsUptime := moonkin.Rotation.IsUsage == proto.BalanceDruid_Rotation_MaximizeIs
	shouldRefreshMf := moonfireUptime <= 0 && useMf
	shouldRefreshIs := insectSwarmUptime <= 0 && useIs

	if moonkin.Talents.Eclipse > 0 {

		solarUptime := moonkin.SolarEclipseProcAura.ExpiresAt() - sim.CurrentTime
		lunarIsActive := moonkin.LunarEclipseProcAura.IsActive()
		solarIsActive := moonkin.SolarEclipseProcAura.IsActive()

		// "Dispelling" eclipse effects before casting if needed
		if float64(lunarUptime-moonkin.Starfire.CurCast.CastTime) <= 0 && useMf {
			lunarIsActive = false
		}
		if float64(solarUptime-moonkin.Wrath.CurCast.CastTime) <= 0 && useIs {
			solarIsActive = false
		}

		if lunarIsActive {
			lunarIsActive = lunarUptime < (moonkin.LunarEclipseProcAura.Duration - playerLatency)
			//fishingForSolar = false
		}
		if solarIsActive {
			solarIsActive = solarUptime < (moonkin.SolarEclipseProcAura.Duration - playerLatency)
			fishingForLunar = false
		}

		// Eclipse
		if solarIsActive || lunarIsActive {
			if maximizeIsUptime && shouldRefreshIs {
				return moonkin.InsectSwarm
			}
			if maximizeMfUptime && shouldRefreshMf {
				return moonkin.Moonfire
			}
			if lunarIsActive {
				if (rotation.UseSmartCooldowns && lunarUptime > 14*time.Second) || sim.GetRemainingDuration() < 15*time.Second {
					moonkin.castMajorCooldown(moonkin.hyperSpeedMCD, sim, target)
					moonkin.castMajorCooldown(moonkin.potionSpeedMCD, sim, target)
					moonkin.useTrinkets(stats.SpellHaste, sim, target)
				}
				return moonkin.Starfire
			} else if solarIsActive {
				if rotation.UseWrath {
					if (rotation.UseSmartCooldowns && solarUptime > 14*time.Second) || sim.GetRemainingDuration() < 15*time.Second {
						moonkin.castMajorCooldown(moonkin.potionWildMagicMCD, sim, target)
						moonkin.useTrinkets(stats.SpellCrit, sim, target)
					}
					return moonkin.Wrath
				}
			}
		}
		if moonkin.Rotation.MfUsage == proto.BalanceDruid_Rotation_BeforeLunar && lunarICD < 2*time.Second && shouldRefreshMf {
			return moonkin.Moonfire
		}
		if moonkin.Rotation.IsUsage == proto.BalanceDruid_Rotation_BeforeSolar && solarICD < 2*time.Second && shouldRefreshIs {
			return moonkin.InsectSwarm
		}
	} else {
		// If Eclipse isn't talented we're not fishing
		fishingForLunar = true
		//fishingForSolar = true
	}

	// Non-Eclipse
	if maximizeMfUptime && shouldRefreshMf {
		return moonkin.Moonfire
	} else if maximizeIsUptime && shouldRefreshIs {
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

func (moonkin *BalanceDruid) useTrinkets(stat stats.Stat, sim *core.Simulation, target *core.Unit) {
	if moonkin.onUseTrinket1.Stat == stat {
		moonkin.castMajorCooldown(moonkin.onUseTrinket1.Cooldown, sim, target)
	}
	if moonkin.onUseTrinket2.Stat == stat {
		moonkin.castMajorCooldown(moonkin.onUseTrinket2.Cooldown, sim, target)
	}
}
