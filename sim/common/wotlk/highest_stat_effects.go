package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type HighestStatAura struct {
	statOptions []stats.Stat
	auras       []*core.Aura
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

	return hsa.auras[bestIdx]
}

func NewHighestStatAura(statOptions []stats.Stat, auraFactory func(stat stats.Stat) *core.Aura) HighestStatAura {
	var auras []*core.Aura
	for _, stat := range statOptions {
		auras = append(auras, auraFactory(stat))
	}

	return HighestStatAura{
		statOptions: statOptions,
		auras:       auras,
	}
}

func init() {
	core.AddEffectsToTest = false

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

			MakeProcTriggerAura(&character.Unit, ProcTrigger{
				Name:       "DMC Greatness",
				Callback:   OnSpellHitDealt | OnPeriodicDamageDealt,
				Harmful:    true,
				ProcChance: 0.35,
				ICD:        time.Second * 45,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellEffect) {
					hsa.Get(character).Activate(sim)
				},
			})
		})
	}
	newDMCGreatnessEffect(42987)
	newDMCGreatnessEffect(44253)
	newDMCGreatnessEffect(44254)
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

			MakeProcTriggerAura(&character.Unit, ProcTrigger{
				Name:       name,
				Callback:   OnSpellHitDealt | OnPeriodicDamageDealt,
				Harmful:    true,
				ProcChance: 0.35,
				ICD:        time.Second * 45,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellEffect) {
					hsa.Get(character).Activate(sim)
				},
			})
		})
	}
	newDeathsChoiceEffect(47115, "Deaths Verdict", 450)
	newDeathsChoiceEffect(47131, "Deaths Verdict H", 510)
	newDeathsChoiceEffect(47303, "Deaths Choice", 450)
	newDeathsChoiceEffect(47464, "Deaths Choice H", 510)

	core.AddEffectsToTest = true
}
