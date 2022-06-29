package core

import (
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

type healthBar struct {
	unit *Unit

	currentHealth  float64

	DamageTakenHealthMetrics *ResourceMetrics
}

func (unit *Unit) EnableHealthBar() {
	unit.healthBar = healthBar{
		unit: unit,

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

	if sim.Log != nil {
		hb.unit.Log(sim, "Spent %0.3f health from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, oldHealth, newHealth)
	}

	hb.currentHealth = newHealth
}

var ChanceOfDeathAuraLabel = "Chance of Death"

func (character *Character) trackChanceOfDeath(healingModel *proto.HealingModel) {
	if healingModel == nil {
		return
	}

	isTanking := false
	for _, target := range character.Env.Encounter.Targets {
		if target.CurrentTarget == &character.Unit {
			isTanking = true
		}
	}
	if !isTanking {
		return
	}

	character.RegisterAura(Aura{
		Label:    ChanceOfDeathAuraLabel,
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			if spellEffect.Damage > 0 {
				aura.Unit.RemoveHealth(sim, spellEffect.Damage)

				if aura.Unit.CurrentHealth() <= 0 && !aura.Unit.Metrics.Died {
					aura.Unit.Metrics.Died = true
					if sim.Log != nil {
						character.Log(sim, "Dead")
					}
				}
			}
		},
		OnPeriodicDamageTaken: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			if spellEffect.Damage > 0 {
				aura.Unit.RemoveHealth(sim, spellEffect.Damage)

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
		character.applyHealingModel(*healingModel)
	}
}

func (character *Character) applyHealingModel(healingModel proto.HealingModel) {
	cadence := DurationFromSeconds(healingModel.CadenceSeconds)
	if cadence == 0 {
		cadence = time.Millisecond * 2500
	}
	healPerTick := healingModel.Hps * (float64(cadence) / float64(time.Second))

	healthMetrics := character.NewHealthMetrics(ActionID{OtherID: proto.OtherAction_OtherActionHealingModel})

	character.RegisterResetEffect(func(sim *Simulation) {
		// Hack since we don't have OnHealingReceived aura handlers yet.
		ardentDefenderAura := character.GetAura("Ardent Defender")

		StartPeriodicAction(sim, PeriodicActionOptions{
			Period: cadence,
			OnAction: func(sim *Simulation) {
				character.GainHealth(sim, healPerTick, healthMetrics)

				if ardentDefenderAura != nil && character.CurrentHealthPercent() >= 0.35 {
					ardentDefenderAura.Deactivate(sim)
				}
			},
		})
	})
}

func (character *Character) GetPresimOptions(playerConfig proto.Player) *PresimOptions {
	healingModel := playerConfig.HealingModel
	if healingModel == nil || healingModel.Hps != 0 {
		// If Hps is not 0, then we don't need to run the presim.
		return nil
	}

	return &PresimOptions{
		SetPresimPlayerOptions: func(player *proto.Player) {
			player.HealingModel = nil
		},

		OnPresimResult: func(presimResult proto.UnitMetrics, iterations int32, duration time.Duration) bool {
			character.applyHealingModel(proto.HealingModel{
				Hps:            presimResult.Dtps.Avg * 1.25,
				CadenceSeconds: healingModel.CadenceSeconds,
			})
			return true
		},
	}
}
