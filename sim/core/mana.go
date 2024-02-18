package core

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const ThreatPerManaGained = 0.5

type manaBar struct {
	unit     *Unit
	BaseMana float64

	currentMana           float64
	manaCastingMetrics    *ResourceMetrics
	manaNotCastingMetrics *ResourceMetrics
	JowManaMetrics        *ResourceMetrics
	VtManaMetrics         *ResourceMetrics
	JowiseManaMetrics     *ResourceMetrics

	ReplenishmentAura *Aura

	// For keeping track of OOM status.
	waitingForMana          float64
	waitingForManaStartTime time.Duration
}

// EnableManaBar will setup caster stat dependencies (int->mana and int->spellcrit)
// as well as enable the mana gain action to regenerate mana.
// It will then enable mana gain metrics for reporting.
func (character *Character) EnableManaBar() {
	character.EnableManaBarWithModifier(1.0)
	character.Unit.SetCurrentPowerBar(ManaBar)
}

func (character *Character) EnableManaBarWithModifier(modifier float64) {
	// Assumes all units have >= 20 intellect.
	// See https://wowwiki-archive.fandom.com/wiki/Base_mana.
	// Subtract out the non-linear part of the formula separately, so that weird
	// mana values are not included when using the stat dependency manager.
	character.AddStat(stats.Mana, 20-15*20*modifier)
	character.AddStatDependency(stats.Intellect, stats.Mana, 15*modifier)

	// This conversion is now universal for
	character.AddStatDependency(stats.Intellect, stats.SpellCrit, CritRatingPerCritChance/166.66667)

	// Not a real spell, just holds metrics from mana gain threat.
	character.RegisterSpell(SpellConfig{
		ActionID: ActionID{OtherID: proto.OtherAction_OtherActionManaGain},
	})

	character.manaCastingMetrics = character.NewManaMetrics(ActionID{OtherID: proto.OtherAction_OtherActionManaRegen, Tag: 1})
	character.manaNotCastingMetrics = character.NewManaMetrics(ActionID{OtherID: proto.OtherAction_OtherActionManaRegen, Tag: 2})

	character.BaseMana = character.GetBaseStats()[stats.Mana]
	character.Unit.manaBar.unit = &character.Unit
}

func (unit *Unit) HasManaBar() bool {
	return unit.manaBar.unit != nil
}

// Gets the Maxiumum mana including bonus and temporary affects that would increase your mana pool.
func (unit *Unit) MaxMana() float64 {
	return unit.stats[stats.Mana]
}
func (unit *Unit) CurrentMana() float64 {
	return unit.currentMana
}
func (unit *Unit) CurrentManaPercent() float64 {
	return unit.CurrentMana() / unit.MaxMana()
}

func (unit *Unit) AddMana(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative mana!")
	}

	oldMana := unit.CurrentMana()
	newMana := min(oldMana+amount, unit.MaxMana())
	metrics.AddEvent(amount, newMana-oldMana)

	if sim.Log != nil {
		unit.Log(sim, "Gained %0.3f mana from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, oldMana, newMana)
	}

	unit.currentMana = newMana
	unit.Metrics.ManaGained += newMana - oldMana
}

func (unit *Unit) SpendMana(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative mana!")
	}

	newMana := unit.CurrentMana() - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		unit.Log(sim, "Spent %0.3f mana from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, unit.CurrentMana(), newMana)
	}

	unit.currentMana = newMana
	unit.Metrics.ManaSpent += amount
}

func (mb *manaBar) doneIteration(sim *Simulation) {
	if mb.unit == nil {
		return
	}

	if mb.waitingForMana != 0 {
		mb.unit.Metrics.AddOOMTime(sim, sim.CurrentTime-mb.waitingForManaStartTime)
	}

	manaGainSpell := mb.unit.GetSpell(ActionID{OtherID: proto.OtherAction_OtherActionManaGain})

	for _, resourceMetrics := range mb.unit.Metrics.resources {
		if resourceMetrics.Type != proto.ResourceType_ResourceTypeMana {
			continue
		}
		if resourceMetrics.ActionID.SameActionIgnoreTag(ActionID{OtherID: proto.OtherAction_OtherActionManaRegen}) {
			continue
		}
		if resourceMetrics.ActionID.SameActionIgnoreTag(ActionID{SpellID: 34917}) {
			// Vampiric Touch mana threat goes to the priest, so it's handled in the priest code.
			continue
		}
		if resourceMetrics.ActualGain <= 0 {
			continue
		}

		manaGainSpell.SpellMetrics[0].Casts += resourceMetrics.EventsForCurrentIteration()
		manaGainSpell.ApplyAOEThreatIgnoreMultipliers(resourceMetrics.ActualGainForCurrentIteration() * ThreatPerManaGained)
	}
}

// Returns the rate of mana regen per second from mp5.
func (unit *Unit) MP5ManaRegenPerSecond() float64 {
	return unit.stats[stats.MP5] / 5.0
}

// Returns the rate of mana regen per second from spirit.
func (unit *Unit) SpiritManaRegenPerSecond() float64 {
	return 0.001 + unit.stats[stats.Spirit]*math.Sqrt(unit.stats[stats.Intellect])*0.003345
}

// Returns the rate of mana regen per second, assuming this unit is
// considered to be casting.
func (unit *Unit) ManaRegenPerSecondWhileCasting() float64 {
	regenRate := unit.MP5ManaRegenPerSecond()

	spiritRegenRate := 0.0
	if unit.PseudoStats.SpiritRegenRateCasting != 0 || unit.PseudoStats.ForceFullSpiritRegen {
		spiritRegenRate = unit.SpiritManaRegenPerSecond() * unit.PseudoStats.SpiritRegenMultiplier
		if !unit.PseudoStats.ForceFullSpiritRegen {
			spiritRegenRate *= unit.PseudoStats.SpiritRegenRateCasting
		}
	}
	regenRate += spiritRegenRate

	return regenRate
}

// Returns the rate of mana regen per second, assuming this unit is
// considered to be not casting.
func (unit *Unit) ManaRegenPerSecondWhileNotCasting() float64 {
	regenRate := unit.MP5ManaRegenPerSecond()

	regenRate += unit.SpiritManaRegenPerSecond() * unit.PseudoStats.SpiritRegenMultiplier

	return regenRate
}

func (unit *Unit) UpdateManaRegenRates() {
	unit.manaTickWhileCasting = unit.ManaRegenPerSecondWhileCasting() * 2
	unit.manaTickWhileNotCasting = unit.ManaRegenPerSecondWhileNotCasting() * 2
}

// Applies 1 'tick' of mana regen, which worth 2s of regeneration based on mp5/int/spirit/etc.
func (unit *Unit) ManaTick(sim *Simulation) {
	if sim.CurrentTime < unit.PseudoStats.FiveSecondRuleRefreshTime {
		regen := unit.manaTickWhileCasting
		unit.AddMana(sim, max(0, regen), unit.manaCastingMetrics)
	} else {
		regen := unit.manaTickWhileNotCasting
		unit.AddMana(sim, max(0, regen), unit.manaNotCastingMetrics)
	}
}

// Returns the amount of time this Unit would need to wait in order to reach
// the desired amount of mana, via mana regen.
//
// Assumes that desiredMana > currentMana. Calculation assumes the Unit
// will not take any actions during this period that would reset the 5-second rule.
func (unit *Unit) TimeUntilManaRegen(desiredMana float64) time.Duration {
	// +1 at the end is to deal with floating point math rounding errors.
	manaNeeded := desiredMana - unit.CurrentMana()
	regenTime := NeverExpires

	regenWhileCasting := unit.ManaRegenPerSecondWhileCasting()
	if regenWhileCasting != 0 {
		regenTime = DurationFromSeconds(manaNeeded/regenWhileCasting) + 1
	}

	// TODO: this needs to have access to the sim to see current time vs unit.PseudoStats.FiveSecondRule.
	//  it is possible that we have been waiting.
	//  In practice this function is always used right after a previous cast so no big deal for now.
	if regenTime > time.Second*5 {
		regenTime = time.Second * 5
		manaNeeded -= regenWhileCasting * 5
		// now we move into spirit based regen.
		regenTime += DurationFromSeconds(manaNeeded / unit.ManaRegenPerSecondWhileNotCasting())
	}

	return regenTime
}

func (sim *Simulation) initManaTickAction() {
	var unitsWithManaBars []*Unit

	for _, party := range sim.Raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			if character.HasManaBar() {
				unitsWithManaBars = append(unitsWithManaBars, &player.GetCharacter().Unit)
			}

			for _, petAgent := range character.PetAgents {
				if petAgent.GetPet().HasManaBar() {
					unitsWithManaBars = append(unitsWithManaBars, &petAgent.GetCharacter().Unit)
				}
			}
		}
	}

	if len(unitsWithManaBars) == 0 {
		return
	}

	interval := time.Second * 2
	pa := &PendingAction{
		NextActionAt: sim.Environment.PrepullStartTime() + interval,
		Priority:     ActionPriorityRegen,
	}
	pa.OnAction = func(sim *Simulation) {
		for _, unit := range unitsWithManaBars {
			if unit.IsEnabled() {
				unit.ManaTick(sim)
			}
		}

		pa.NextActionAt = sim.CurrentTime + interval
		sim.AddPendingAction(pa)
	}
	sim.AddPendingAction(pa)
}

func (mb *manaBar) reset() {
	if mb.unit == nil {
		return
	}

	mb.currentMana = mb.unit.MaxMana()
	mb.waitingForMana = 0
	mb.waitingForManaStartTime = 0
}

func (mb *manaBar) IsOOM() bool {
	return mb.waitingForMana != 0
}
func (mb *manaBar) StartOOMEvent(sim *Simulation, requiredMana float64) {
	mb.waitingForManaStartTime = sim.CurrentTime
	mb.waitingForMana = requiredMana
	mb.unit.Metrics.MarkOOM(sim)
}
func (mb *manaBar) EndOOMEvent(sim *Simulation) {
	eventDuration := sim.CurrentTime - mb.waitingForManaStartTime
	mb.unit.Metrics.AddOOMTime(sim, eventDuration)
	mb.waitingForManaStartTime = 0
	mb.waitingForMana = 0
}

type ManaCostOptions struct {
	BaseCost   float64
	FlatCost   float64 // Alternative to BaseCost for giving a flat value.
	Multiplier float64 // It's OK to leave this at 0, will default to 1.
}
type ManaCost struct {
	ResourceMetrics *ResourceMetrics
}

func newManaCost(spell *Spell, options ManaCostOptions) *ManaCost {
	baseCost := TernaryFloat64(options.FlatCost > 0, options.FlatCost, options.BaseCost*spell.Unit.BaseMana)
	if player := spell.Unit.Env.Raid.GetPlayerFromUnit(spell.Unit); player != nil {
		if player.GetCharacter().HasTrinketEquipped(45703) { // Spark of Hope
			baseCost = max(0, baseCost-44)
		}
	}

	spell.DefaultCast.Cost = baseCost * TernaryFloat64(options.Multiplier == 0, 1, options.Multiplier)

	return &ManaCost{
		ResourceMetrics: spell.Unit.NewManaMetrics(spell.ActionID),
	}
}

func (mc *ManaCost) MeetsRequirement(sim *Simulation, spell *Spell) bool {
	spell.CurCast.Cost = spell.ApplyCostModifiers(spell.CurCast.Cost)
	meetsRequirement := spell.Unit.CurrentMana() >= spell.CurCast.Cost

	if spell.CurCast.Cost > 0 {
		if meetsRequirement {
			if spell.Unit.IsOOM() {
				spell.Unit.EndOOMEvent(sim)
			}
		} else {
			if spell.Unit.IsOOM() {
				// Continuation of OOM event.
				spell.Unit.waitingForMana = min(spell.Unit.waitingForMana, spell.CurCast.Cost)
			} else {
				spell.Unit.StartOOMEvent(sim, spell.CurCast.Cost)
			}
		}
	}

	return meetsRequirement
}
func (mc *ManaCost) CostFailureReason(sim *Simulation, spell *Spell) string {
	return fmt.Sprintf("not enough mana (Current Mana = %0.03f, Mana Cost = %0.03f)", spell.Unit.CurrentMana(), spell.CurCast.Cost)
}
func (mc *ManaCost) SpendCost(sim *Simulation, spell *Spell) {
	if spell.CurCast.Cost > 0 {
		spell.Unit.SpendMana(sim, spell.CurCast.Cost, mc.ResourceMetrics)
		spell.Unit.PseudoStats.FiveSecondRuleRefreshTime = max(sim.CurrentTime+time.Second*5, spell.Unit.Hardcast.Expires)
	}
}
func (mc *ManaCost) IssueRefund(_ *Simulation, _ *Spell) {}
