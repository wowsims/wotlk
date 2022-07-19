package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/common/wotlk"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Used for filtering sigil bonuses based on skill
type DeathKnightAbility uint8

const (
	Ability_Obliterate DeathKnightAbility = iota
	Ability_ScourgeStrike
	Ability_DeathStrike
)

var ItemSetScourgeborne = core.NewItemSet(core.ItemSet{
	Name: "Scourgeborne Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the critical strike chance of your obliterate
			// scourge strike and death strike abilities by 5%
		},
		4: func(agent core.Agent) {
			// Your obliterate, scourge strike and death strike
			// generate 5 additional runic power
		},
	},
})

func (deathKnight *DeathKnight) scourgeborneCritBonus() float64 {
	return core.TernaryFloat64(deathKnight.HasSetBonus(ItemSetScourgeborne, 2), 5.0, 0.0)
}

func (deathKnight *DeathKnight) scourgeborneRunicPowerBonus() float64 {
	return core.TernaryFloat64(deathKnight.HasSetBonus(ItemSetScourgeborne, 4), 5.0, 0.0)
}

func (deathKnight *DeathKnight) sigilOfAwarenessBonus(skill DeathKnightAbility) float64 {
	if deathKnight.Equip[proto.ItemSlot_ItemSlotRanged].ID != 40207 {
		return 0
	}

	switch skill {
	case Ability_Obliterate:
		return 336
	case Ability_ScourgeStrike:
		return 189
	case Ability_DeathStrike:
		return 315
	default:
		return 0
	}
}

func (deathKnight *DeathKnight) sigilOfTheFrozenConscienceBonus() float64 {
	return core.TernaryFloat64(deathKnight.Equip[proto.ItemSlot_ItemSlotRanged].ID == 40822, 111, 0)
}

func (deathKnight *DeathKnight) sigilOfTheWildBuckBonus() float64 {
	return core.TernaryFloat64(deathKnight.Equip[proto.ItemSlot_ItemSlotRanged].ID == 40867, 80, 0)
}

func (deathKnight *DeathKnight) sigilOfArthriticBindingBonus() float64 {
	return core.TernaryFloat64(deathKnight.Equip[proto.ItemSlot_ItemSlotRanged].ID == 40875, 91, 0)
}

func init() {
	core.NewItemEffect(40715, func(agent core.Agent) {
		deathKnight := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := deathKnight.NewTemporaryStatsAura("Sigil of Haunted Dreams Proc", core.ActionID{ItemID: 40715}, stats.Stats{stats.MeleeCrit: 173.0, stats.SpellCrit: 173.0}, time.Second*10)

		icd := core.Cooldown{
			Timer:    deathKnight.NewTimer(),
			Duration: time.Second * 45.0,
		}

		core.MakePermanent(deathKnight.GetOrRegisterAura(core.Aura{
			Label: "Sigil of Haunted Dreams",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) || spell != deathKnight.BloodStrike {
					return
				}

				if sim.RandomFloat("SigilOfHauntedDreams") < 0.15 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		}))
	})

	core.NewItemEffect(47673, func(agent core.Agent) {
		deathKnight := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := deathKnight.NewTemporaryStatsAura("Sigil of Virulence Proc", core.ActionID{ItemID: 47673}, stats.Stats{stats.Strength: 200.0}, time.Second*20)

		icd := core.Cooldown{
			Timer:    deathKnight.NewTimer(),
			Duration: time.Second * 10.0,
		}

		core.MakePermanent(deathKnight.GetOrRegisterAura(core.Aura{
			Label: "Sigil of Virulence",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) || !deathKnight.IsFuStrike(spell) {
					return
				}

				if sim.RandomFloat("SigilOfVirulence") < 0.80 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		}))
	})

	core.NewItemEffect(50459, func(agent core.Agent) {
		character := agent.GetCharacter()
		deathKnight := agent.(DeathKnightAgent).GetDeathKnight()

		procAura := wotlk.MakeStackingAura(character, wotlk.StackingProcAura{
			Aura: core.Aura{
				Label:     "Sigil of the Hanged Man Proc",
				ActionID:  core.ActionID{ItemID: 50459},
				Duration:  time.Second * 15,
				MaxStacks: 3,
			},
			BonusPerStack: stats.Stats{stats.Strength: 73},
		})

		core.MakePermanent(deathKnight.GetOrRegisterAura(core.Aura{
			Label: "Sigil of the Hanged Man",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !deathKnight.IsFuStrike(spell) {
					return
				}

				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		}))
	})
}
