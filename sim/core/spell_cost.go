package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/stats"
)

// Handles computing the cost of spells and checking whether the Unit
// meets them.
type SpellCost interface {
	// For initialization logic that requires a reference to the spell
	// to which this cost will apply.
	Init(*Spell)

	// Whether the Unit associated with the spell meets the resource cost
	// requirements to cast the spell.
	MeetsRequirement(*Spell) bool

	// Logs a message for when the cast fails due to lack of resources.
	LogCostFailure(*Simulation, *Spell)

	// Subtracts the resources used from a cast from the Unit.
	SpendCost(*Simulation, *Spell)
}

type ManaCostOptions struct {
	BaseCost   float64
	Multiplier float64 // It's OK to leave this at 0, will default to 1.
}
type ManaCost struct {
	BaseCost   float64
	Multiplier float64

	ResourceMetrics *ResourceMetrics
}

func NewManaCost(options ManaCostOptions) *ManaCost {
	if options.Multiplier == 0 {
		options.Multiplier = 1
	}
	return &ManaCost{
		BaseCost:   options.BaseCost,
		Multiplier: options.Multiplier,
	}
}

func (mc *ManaCost) Init(spell *Spell) {
	mc.BaseCost = mc.BaseCost * spell.Unit.BaseMana
	spell.ResourceType = stats.Mana
	spell.BaseCost = mc.BaseCost
	spell.DefaultCast.Cost = mc.BaseCost * mc.Multiplier
	mc.ResourceMetrics = spell.Unit.NewManaMetrics(spell.ActionID)
}
func (mc *ManaCost) MeetsRequirement(spell *Spell) bool {
	spell.CurCast.Cost = spell.ApplyCostModifiers(spell.CurCast.Cost)
	return spell.Unit.CurrentMana() >= spell.CurCast.Cost
}
func (mc *ManaCost) LogCostFailure(sim *Simulation, spell *Spell) {
	spell.Unit.Log(sim,
		"Failed casting %s, not enough mana. (Current Mana = %0.03f, Mana Cost = %0.03f)",
		spell.ActionID, spell.Unit.CurrentMana(), spell.CurCast.Cost)
}
func (mc *ManaCost) SpendCost(sim *Simulation, spell *Spell) {
	if spell.CurCast.Cost > 0 {
		spell.Unit.SpendMana(sim, spell.CurCast.Cost, mc.ResourceMetrics)
		spell.Unit.PseudoStats.FiveSecondRuleRefreshTime = sim.CurrentTime + time.Second*5
	}
}
