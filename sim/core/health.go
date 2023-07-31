package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type healthBar struct {
	unit *Unit

	currentHealth float64

	DamageTakenHealthMetrics *ResourceMetrics
}

func (unit *Unit) EnableHealthBar() {
	unit.healthBar = healthBar{
		unit:                     unit,
		DamageTakenHealthMetrics: unit.NewHealthMetrics(ActionID{OtherID: proto.OtherAction_OtherActionDamageTaken}),
	}
}

func (unit *Unit) HasHealthBar() bool {
	return unit.healthBar.unit != nil
}

func (hb *healthBar) reset(_ *Simulation) {
	if hb.unit == nil {
		return
	}
	hb.currentHealth = hb.MaxHealth()
}

func (hb *healthBar) MaxHealth() float64 {
	return hb.unit.stats[stats.Health]
}

func (hb *healthBar) CurrentHealth() float64 {
	return hb.currentHealth
}

func (hb *healthBar) CurrentHealthPercent() float64 {
	return hb.currentHealth / hb.unit.stats[stats.Health]
}

func (hb *healthBar) GainHealth(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to gain negative health!")
	}

	oldHealth := hb.currentHealth
	newHealth := MinFloat(oldHealth+amount, hb.unit.MaxHealth())
	metrics.AddEvent(amount, newHealth-oldHealth)

	if sim.Log != nil {
		hb.unit.Log(sim, "Gained %0.3f health from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, oldHealth, newHealth)
	}

	hb.currentHealth = newHealth
}

func (hb *healthBar) RemoveHealth(sim *Simulation, amount float64) {
	if amount < 0 {
		panic("Trying to remove negative health!")
	}

	oldHealth := hb.currentHealth
	newHealth := MaxFloat(oldHealth-amount, 0)
	metrics := hb.DamageTakenHealthMetrics
	metrics.AddEvent(-amount, newHealth-oldHealth)

	// TMI calculations need timestamps and Max HP information for each damage taken event
	if hb.unit.Metrics.isTanking {
		entry := tmiListItem{
			Timestamp:      sim.CurrentTime,
			WeightedDamage: amount / hb.MaxHealth(),
		}
		hb.unit.Metrics.tmiList = append(hb.unit.Metrics.tmiList, entry)
	}

	if sim.Log != nil {
		hb.unit.Log(sim, "Spent %0.3f health from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, oldHealth, newHealth)
	}

	hb.currentHealth = newHealth
}

var ChanceOfDeathAuraLabel = "Chance of Death"

func (character *Character) trackChanceOfDeath(healingModel *proto.HealingModel) {
	character.Unit.Metrics.isTanking = false
	for _, target := range character.Env.Encounter.TargetUnits {
		if target.CurrentTarget == &character.Unit {
			character.Unit.Metrics.isTanking = true
		}
	}
	if !character.Unit.Metrics.isTanking {
		return
	}

	if healingModel == nil {
		return
	}

	character.Unit.Metrics.tmiBin = healingModel.BurstWindow

	character.RegisterAura(Aura{
		Label:    ChanceOfDeathAuraLabel,
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if result.Damage > 0 {
				aura.Unit.RemoveHealth(sim, result.Damage)

				if aura.Unit.CurrentHealth() <= 0 && !aura.Unit.Metrics.Died {
					aura.Unit.Metrics.Died = true
					if sim.Log != nil {
						character.Log(sim, "Dead")
					}
				}
			}
		},
		OnPeriodicDamageTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if result.Damage > 0 {
				aura.Unit.RemoveHealth(sim, result.Damage)

				if aura.Unit.CurrentHealth() <= 0 && !aura.Unit.Metrics.Died {
					aura.Unit.Metrics.Died = true
					if sim.Log != nil {
						character.Log(sim, "Dead")
					}
				}
			}
		},
	})

	if healingModel.Hps != 0 {
		character.applyHealingModel(healingModel)
	}
}

func (character *Character) applyHealingModel(healingModel *proto.HealingModel) {
	// Store variance parameters for healing cadence. Note that low rolls on
	// cadence are special cased here so that the model is still well-behaved
	// when CadenceVariation exceeds CadenceSeconds.
	medianCadence := healingModel.CadenceSeconds
	if medianCadence == 0 {
		medianCadence = 2.0
	}
	minCadence := MaxFloat(0.0, medianCadence-healingModel.CadenceVariation)
	cadenceVariationLow := medianCadence - minCadence

	healthMetrics := character.NewHealthMetrics(ActionID{OtherID: proto.OtherAction_OtherActionHealingModel})

	character.RegisterResetEffect(func(sim *Simulation) {
		// Hack since we don't have OnHealingReceived aura handlers yet.
		//ardentDefenderAura := character.GetAura("Ardent Defender")
		willOfTheNecropolisAura := character.GetAura("Will of The Necropolis")

		// Initialize randomized cadence model
		timeToNextHeal := DurationFromSeconds(0.0)
		healPerTick := 0.0
		pa := &PendingAction{
			NextActionAt: timeToNextHeal,
		}

		pa.OnAction = func(sim *Simulation) {
			// Use modeled HPS to scale heal per tick based on random cadence
			healPerTick = healingModel.Hps * (float64(timeToNextHeal) / float64(time.Second))

			// Execute the heal
			character.GainHealth(sim, healPerTick*character.PseudoStats.HealingTakenMultiplier, healthMetrics)

			// Might use this again in the future to track "absorb" metrics but currently disabled
			//if ardentDefenderAura != nil && character.CurrentHealthPercent() >= 0.35 {
			//	ardentDefenderAura.Deactivate(sim)
			//}

			if willOfTheNecropolisAura != nil && character.CurrentHealthPercent() > 0.35 {
				willOfTheNecropolisAura.Deactivate(sim)
			}

			// Random roll for time to next heal. In the case where CadenceVariation exceeds CadenceSeconds, then
			// CadenceSeconds is treated as the median, with two separate uniform distributions to the left and right
			// of it.
			signRoll := sim.RandomFloat("Healing Cadence Variation Sign")
			magnitudeRoll := sim.RandomFloat("Healing Cadence Variation Magnitude")

			if signRoll < 0.5 {
				timeToNextHeal = DurationFromSeconds(minCadence + magnitudeRoll*cadenceVariationLow)
			} else {
				timeToNextHeal = DurationFromSeconds(medianCadence + magnitudeRoll*healingModel.CadenceVariation)
			}

			// Refresh action
			pa.NextActionAt = sim.CurrentTime + timeToNextHeal
			sim.AddPendingAction(pa)
		}

		sim.AddPendingAction(pa)
	})
}

func (character *Character) GetPresimOptions(playerConfig *proto.Player) *PresimOptions {
	healingModel := playerConfig.HealingModel
	if healingModel == nil || healingModel.Hps != 0 || healingModel.CadenceSeconds == 0 {
		// If Hps is not 0, then we don't need to run the presim.
		// Tank sims should always have nonzero Cadence set, even if disabled
		return nil
	}
	return &PresimOptions{
		SetPresimPlayerOptions: func(player *proto.Player) {
			player.HealingModel = nil
		},
		OnPresimResult: func(presimResult *proto.UnitMetrics, iterations int32, duration time.Duration) bool {
			character.applyHealingModel(&proto.HealingModel{
				Hps:            presimResult.Dtps.Avg * 1.50,
				CadenceSeconds: healingModel.CadenceSeconds,
			})
			return true
		},
	}
}
