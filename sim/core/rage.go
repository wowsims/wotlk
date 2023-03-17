package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
)

const MaxRage = 100.0
const RageFactor = 453.3
const ThreatPerRageGained = 5

// OnRageGain is called any time rage is increased.
type OnRageGain func(sim *Simulation)

type rageBar struct {
	unit *Unit

	startingRage float64
	currentRage  float64

	onRageGain OnRageGain

	RageRefundMetrics *ResourceMetrics
}

type RageBarOptions struct {
	StartingRage   float64
	RageMultiplier float64
	MHSwingSpeed   float64
	OHSwingSpeed   float64
}

func (unit *Unit) EnableRageBar(options RageBarOptions, onRageGain OnRageGain) {
	rageFromDamageTakenMetrics := unit.NewRageMetrics(ActionID{OtherID: proto.OtherAction_OtherActionDamageTaken})

	unit.SetCurrentPowerBar(RageBar)
	unit.RegisterAura(Aura{
		Label:    "RageBar",
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if unit.GetCurrentPowerBar() != RageBar {
				return
			}
			if result.Outcome.Matches(OutcomeMiss) {
				return
			}
			if !spell.ProcMask.Matches(ProcMaskMelee) {
				return
			}
			if !spell.ProcMask.Matches(ProcMaskWhiteHit) {
				return
			}

			// Need separate check to exclude auto replacers (e.g. Heroic Strike and Cleave).
			if spell.ProcMask.Matches(ProcMaskMeleeMHSpecial) {
				return
			}

			var hitFactor float64
			var speed float64
			if spell.IsMH() {
				hitFactor = 3.5
				speed = options.MHSwingSpeed
			} else {
				hitFactor = 1.75
				speed = options.OHSwingSpeed
			}

			if result.Outcome.Matches(OutcomeCrit) {
				hitFactor *= 2
			}

			damage := result.Damage
			if result.Outcome.Matches(OutcomeDodge | OutcomeParry) {
				// Rage is still generated for dodges/parries, based on the damage it WOULD have done.
				damage = result.PreOutcomeDamage
			}

			// generatedRage is capped for very low damage swings
			generatedRage := MinFloat((damage*7.5/RageFactor+hitFactor*speed)/2, damage*15/RageFactor)

			generatedRage *= options.RageMultiplier

			var metrics *ResourceMetrics
			if spell.Cost != nil {
				metrics = spell.Cost.(*RageCost).ResourceMetrics
			} else {
				if spell.ResourceMetrics == nil {
					spell.ResourceMetrics = spell.Unit.NewRageMetrics(spell.ActionID)
				}
				metrics = spell.ResourceMetrics
			}
			unit.AddRage(sim, generatedRage, metrics)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if unit.GetCurrentPowerBar() != RageBar {
				return
			}
			generatedRage := result.Damage * 2.5 / RageFactor
			unit.AddRage(sim, generatedRage, rageFromDamageTakenMetrics)
		},
	})

	// Not a real spell, just holds metrics from rage gain threat.
	unit.RegisterSpell(SpellConfig{
		ActionID: ActionID{OtherID: proto.OtherAction_OtherActionRageGain},
	})

	unit.rageBar = rageBar{
		unit:         unit,
		startingRage: MaxFloat(0, MinFloat(options.StartingRage, MaxRage)),
		onRageGain:   onRageGain,

		RageRefundMetrics: unit.NewRageMetrics(ActionID{OtherID: proto.OtherAction_OtherActionRefund}),
	}
}

func (unit *Unit) HasRageBar() bool {
	return unit.rageBar.unit != nil
}

func (rb *rageBar) CurrentRage() float64 {
	return rb.currentRage
}

func (rb *rageBar) AddRage(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative rage!")
	}

	newRage := MinFloat(rb.currentRage+amount, MaxRage)
	metrics.AddEvent(amount, newRage-rb.currentRage)

	if sim.Log != nil {
		rb.unit.Log(sim, "Gained %0.3f rage from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, rb.currentRage, newRage)
	}

	rb.currentRage = newRage
	rb.onRageGain(sim)
}

func (rb *rageBar) SpendRage(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative rage!")
	}

	newRage := rb.currentRage - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		rb.unit.Log(sim, "Spent %0.3f rage from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, rb.currentRage, newRage)
	}

	rb.currentRage = newRage
}

func (rb *rageBar) reset(sim *Simulation) {
	if rb.unit == nil {
		return
	}

	rb.currentRage = rb.startingRage
}

func (rb *rageBar) doneIteration() {
	if rb.unit == nil {
		return
	}

	rageGainSpell := rb.unit.GetSpell(ActionID{OtherID: proto.OtherAction_OtherActionRageGain})

	for _, resourceMetrics := range rb.unit.Metrics.resources {
		if resourceMetrics.Type != proto.ResourceType_ResourceTypeRage {
			continue
		}
		if resourceMetrics.ActionID.SameActionIgnoreTag(ActionID{OtherID: proto.OtherAction_OtherActionDamageTaken}) {
			continue
		}
		if resourceMetrics.ActionID.SameActionIgnoreTag(ActionID{OtherID: proto.OtherAction_OtherActionRefund}) {
			continue
		}
		if resourceMetrics.ActualGain <= 0 {
			continue
		}

		// Need to exclude rage gained from white hits. Rather than have a manual list of all IDs that would
		// apply here (autos, WF attack, sword spec procs, etc), just check if the effect caused any damage.
		sourceSpell := rb.unit.GetSpell(resourceMetrics.ActionID)
		if sourceSpell != nil && sourceSpell.SpellMetrics[0].TotalDamage > 0 {
			continue
		}

		rageGainSpell.SpellMetrics[0].Casts += resourceMetrics.EventsForCurrentIteration()
		rageGainSpell.ApplyAOEThreatIgnoreMultipliers(resourceMetrics.ActualGainForCurrentIteration() * ThreatPerRageGained)
	}
}

type RageCostOptions struct {
	Cost float64

	Refund        float64
	RefundMetrics *ResourceMetrics // Optional, will default to unit.RageRefundMetrics if not supplied.
}
type RageCost struct {
	Refund          float64
	RefundMetrics   *ResourceMetrics
	ResourceMetrics *ResourceMetrics
}

func newRageCost(spell *Spell, options RageCostOptions) *RageCost {
	spell.DefaultCast.Cost = options.Cost
	if options.Refund > 0 && options.RefundMetrics == nil {
		options.RefundMetrics = spell.Unit.RageRefundMetrics
	}

	return &RageCost{
		Refund:          options.Refund * options.Cost,
		RefundMetrics:   options.RefundMetrics,
		ResourceMetrics: spell.Unit.NewRageMetrics(spell.ActionID),
	}
}

func (rc *RageCost) MeetsRequirement(spell *Spell) bool {
	spell.CurCast.Cost = spell.ApplyCostModifiers(spell.CurCast.Cost)
	return spell.Unit.CurrentRage() >= spell.CurCast.Cost
}
func (rc *RageCost) LogCostFailure(sim *Simulation, spell *Spell) {
	spell.Unit.Log(sim,
		"Failed casting %s, not enough rage. (Current Rage = %0.03f, Rage Cost = %0.03f)",
		spell.ActionID, spell.Unit.CurrentRage(), spell.CurCast.Cost)
}
func (rc *RageCost) SpendCost(sim *Simulation, spell *Spell) {
	if spell.CurCast.Cost > 0 {
		spell.Unit.SpendRage(sim, spell.CurCast.Cost, rc.ResourceMetrics)
	}
}
func (rc *RageCost) IssueRefund(sim *Simulation, spell *Spell) {
	if rc.Refund > 0 {
		spell.Unit.AddRage(sim, rc.Refund, rc.RefundMetrics)
	}
}
