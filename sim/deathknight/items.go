package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/common/wotlk"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var ItemSetScourgeborneBattlegear = core.NewItemSet(core.ItemSet{
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

func (deathKnight *DeathKnight) scourgeborneBattlegearCritBonus() float64 {
	return core.TernaryFloat64(deathKnight.HasSetBonus(ItemSetScourgeborneBattlegear, 2), 5.0, 0.0)
}

func (deathKnight *DeathKnight) scourgeborneBattlegearRunicPowerBonus() float64 {
	return core.TernaryFloat64(deathKnight.HasSetBonus(ItemSetScourgeborneBattlegear, 4), 5.0, 0.0)
}

var ItemSetScourgebornePlate = core.NewItemSet(core.ItemSet{
	Name: "Scourgeborne Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the critical strike chance of your plague
			// strike by 10%
		},
		4: func(agent core.Agent) {
			// TODO:
			// Increases the duration of your Icebound Fortitude by 3 seconds
		},
	},
})

func (deathKnight *DeathKnight) scourgebornePlateCritBonus() float64 {
	return core.TernaryFloat64(deathKnight.HasSetBonus(ItemSetScourgebornePlate, 2), 10.0, 0.0)
}

var ItemSetDarkrunedBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Darkruned Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the critical strike chance of Death Coil
			// and Frost Strike by 8%
		},
		4: func(agent core.Agent) {
			// Increases the damage bonus done per disease by 20%
			// on Blood Strike, Heart Strike, Obliterate and Scourge Strike
		},
	},
})

func (deathKnight *DeathKnight) darkrunedBattlegearCritBonus() float64 {
	return core.TernaryFloat64(deathKnight.HasSetBonus(ItemSetDarkrunedBattlegear, 2), 8.0, 0.0)
}

func (deathKnight *DeathKnight) darkrunedBattlegearDiseaseBonus(baseMultiplier float64) float64 {
	return core.TernaryFloat64(deathKnight.HasSetBonus(ItemSetDarkrunedBattlegear, 4), baseMultiplier*1.2, baseMultiplier)
}

// TODO:
var ItemSetDarkrunedPlate = core.NewItemSet(core.ItemSet{
	Name: "Darkruned Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage done by Rune Strike by 10%
		},
		4: func(agent core.Agent) {
			// Anti-magic shell also grants you 10% reduction
			// to physical damage taken
		},
	},
})

var ItemSetThassariansBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Thassarian's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Blood Strike and Heart Strike abilities have a
			// chance to grant you 180 additional strength for 15 sec.
			deathKnight := agent.(DeathKnightAgent).GetDeathKnight()
			deathKnight.registerThassariansBattlegearProc()
		},
		4: func(agent core.Agent) {
			// Your Blood Plague ability now has a chance for its
			// damage to be critical strikes.
		},
	},
})

func (deathKnight *DeathKnight) registerThassariansBattlegearProc() {
	procAura := deathKnight.NewTemporaryStatsAura("Unholy Might Proc", core.ActionID{SpellID: 67117}, stats.Stats{stats.Strength: 180.0}, time.Second*15)

	icd := core.Cooldown{
		Timer:    deathKnight.NewTimer(),
		Duration: time.Second * 45.0,
	}

	core.MakePermanent(deathKnight.GetOrRegisterAura(core.Aura{
		Label: "Unholy Might",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !icd.IsReady(sim) || (spell != deathKnight.BloodStrike /* && spell != deathKnight.HeartStrike*/) {
				return
			}

			if sim.RandomFloat("UnholyMight") < 0.5 {
				icd.Use(sim)
				procAura.Activate(sim)
			}
		},
	}))
}

func (deathKnight *DeathKnight) sigilOfAwarenessBonus(spell *core.Spell) float64 {
	if deathKnight.Equip[proto.ItemSlot_ItemSlotRanged].ID != 40207 {
		return 0
	}

	if spell == deathKnight.Obliterate {
		return 336
	} else if spell == deathKnight.ScourgeStrike {
		return 189
	} else if spell == deathKnight.DeathStrike {
		return 315
	}
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
