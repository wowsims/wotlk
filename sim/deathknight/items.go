package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/common/wotlk"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
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

func (deathKnight *DeathKnight) sigilOfAwarenessBonus(spell *core.Spell) float64 {
	if deathKnight.Equip[proto.ItemSlot_ItemSlotRanged].ID != 40207 {
		return 0
	}

	if spell == deathKnight.Obliterate {
		return 336
	} else if spell == deathKnight.ScourgeStrike {
		return 189
	} // else if spell == deathKnight.DeathStrike {
	// 	return 315
	// }
	return 0
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

	CreateGladiatorsSigil(42619, "Hateful", 106, 6)
	CreateGladiatorsSigil(42620, "Deadly", 120, 10)
	CreateGladiatorsSigil(42621, "Furious", 144, 10)
	CreateGladiatorsSigil(42622, "Relentless", 172, 10)
	CreateGladiatorsSigil(51417, "Wrathful", 204, 10)
}

func CreateGladiatorsSigil(id int32, name string, ap float64, seconds time.Duration) {
	core.NewItemEffect(id, func(agent core.Agent) {
		deathKnight := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := deathKnight.NewTemporaryStatsAura(name+" Gladiator's Sigil of Strife Proc", core.ActionID{ItemID: id}, stats.Stats{stats.AttackPower: ap}, time.Second*seconds)

		core.MakePermanent(deathKnight.GetOrRegisterAura(core.Aura{
			Label: name + " Gladiator's Sigil of Strife",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell != deathKnight.PlagueStrike {
					return
				}

				procAura.Activate(sim)
			},
		}))
	})
}
