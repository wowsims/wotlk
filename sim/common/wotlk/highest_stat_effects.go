package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type HighestStatAura struct {
	statOptions []stats.Stat
	auras       []*core.Aura
	factory     func(stat stats.Stat) *core.Aura
}

func (hsa HighestStatAura) Get(character *core.Character) *core.Aura {
	bestValue := 0.0
	bestIdx := 0

	for i, stat := range hsa.statOptions {
		value := character.GetStat(stat)
		if value > bestValue {
			bestValue = value
			bestIdx = i
		}
	}

	a := hsa.auras[bestIdx]
	if a == nil {
		a = hsa.factory(hsa.statOptions[bestIdx])
		hsa.auras[bestIdx] = a
	}
	return a
}

func NewHighestStatAura(statOptions []stats.Stat, auraFactory func(stat stats.Stat) *core.Aura) HighestStatAura {
	return HighestStatAura{
		statOptions: statOptions,
		factory:     auraFactory,
		auras:       make([]*core.Aura, len(statOptions)),
	}
}

func init() {
	newDMCGreatnessEffect := func(itemID int32) {
		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			hsa := NewHighestStatAura(
				[]stats.Stat{
					stats.Strength,
					stats.Agility,
					stats.Intellect,
					stats.Spirit,
				},
				func(stat stats.Stat) *core.Aura {
					bonus := stats.Stats{}
					bonus[stat] = 300
					return character.NewTemporaryStatsAura("DMC Greatness "+stat.StatName()+" Proc", core.ActionID{ItemID: itemID}, bonus, time.Second*15)
				})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "DMC Greatness",
				Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
				Harmful:    true,
				ProcChance: 0.35,
				ICD:        time.Second * 45,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					hsa.Get(character).Activate(sim)
				},
			})
		})
	}
	newDMCGreatnessEffect(42987)
	newDMCGreatnessEffect(44253)
	newDMCGreatnessEffect(44254)
	core.AddEffectsToTest = false
	newDMCGreatnessEffect(44255)

	newDeathsChoiceEffect := func(itemID int32, name string, amount float64) {
		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			hsa := NewHighestStatAura(
				[]stats.Stat{
					stats.Strength,
					stats.Agility,
				},
				func(stat stats.Stat) *core.Aura {
					bonus := stats.Stats{}
					bonus[stat] = amount
					return character.NewTemporaryStatsAura(name+" "+stat.StatName()+" Proc", core.ActionID{ItemID: itemID}, bonus, time.Second*15)
				})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       name,
				Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
				Harmful:    true,
				ProcChance: 0.35,
				ICD:        time.Second * 45,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					hsa.Get(character).Activate(sim)
				},
			})
		})
	}
	newDeathsChoiceEffect(47115, "Deaths Verdict", 450)
	core.AddEffectsToTest = false
	newDeathsChoiceEffect(47131, "Deaths Verdict H", 510)
	newDeathsChoiceEffect(47303, "Deaths Choice", 450)
	newDeathsChoiceEffect(47464, "Deaths Choice H", 510)
}
