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

func (hsa HighestStatAura) Init(character *core.Character) {
	for i, stat := range hsa.statOptions {
		hsa.auras[i] = hsa.factory(stat)
	}
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

			auraIDs := map[stats.Stat]core.ActionID{
				stats.Strength:  core.ActionID{SpellID: 60229},
				stats.Agility:   core.ActionID{SpellID: 60233},
				stats.Intellect: core.ActionID{SpellID: 60234},
				stats.Spirit:    core.ActionID{SpellID: 60235},
			}

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
					actionId := auraIDs[stat]
					return character.NewTemporaryStatsAura("DMC Greatness "+stat.StatName()+" Proc", actionId, bonus, time.Second*15)
				})

			hsa.Init(character)
			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "DMC Greatness",
				Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
				ProcMask:   core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
				Harmful:    true,
				ProcChance: 0.35,
				ICD:        time.Second * 45,
				ActionID:   core.ActionID{ItemID: itemID},
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					hsa.Get(character).Activate(sim)
				},
			})
			hsa.Get(character).Icd = triggerAura.Icd
		})
	}
	core.AddEffectsToTest = false
	newDMCGreatnessEffect(42987)
	newDMCGreatnessEffect(44253)
	newDMCGreatnessEffect(44254)
	core.AddEffectsToTest = true
	newDMCGreatnessEffect(44255)

	newDeathsChoiceEffect := func(itemID int32, auraIDs map[stats.Stat]core.ActionID, name string, amount float64) {
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
					actionId := auraIDs[stat]
					return character.NewTemporaryStatsAura(name+" "+stat.StatName()+" Proc", actionId, bonus, time.Second*15)
				})

			hsa.Init(character)
			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       name,
				Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
				ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
				Harmful:    true,
				ProcChance: 0.35,
				ActionID:   core.ActionID{ItemID: itemID},
				ICD:        time.Second * 45,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					hsa.Get(character).Activate(sim)
				},
			})
			hsa.Get(character).Icd = triggerAura.Icd
		})
	}
	core.AddEffectsToTest = false

	normalAuraIDs := map[stats.Stat]core.ActionID{
		stats.Strength: core.ActionID{SpellID: 67708},
		stats.Agility:  core.ActionID{SpellID: 67703},
	}

	heroicAuraIDs := map[stats.Stat]core.ActionID{
		stats.Strength: core.ActionID{SpellID: 67773},
		stats.Agility:  core.ActionID{SpellID: 67772},
	}

	newDeathsChoiceEffect(47115, normalAuraIDs, "Deaths Verdict", 450)
	newDeathsChoiceEffect(47131, heroicAuraIDs, "Deaths Verdict H", 510)

	newDeathsChoiceEffect(47303, normalAuraIDs, "Deaths Choice", 450)
	core.AddEffectsToTest = true
	newDeathsChoiceEffect(47464, heroicAuraIDs, "Deaths Choice H", 510)
}
