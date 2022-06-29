package rogue

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var ItemSetAssassination = core.NewItemSet(core.ItemSet{
	Name: "Assassination Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
			// Your Eviscerate and Envenom abilities cost 10 less energy.
			// Handled in eviscerate.go.
		},
	},
})

var ItemSetNetherblade = core.NewItemSet(core.ItemSet{
	Name: "Netherblade",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the duration of your Slice and Dice ability by 3 sec.
			// Handled in slice_and_dice.go.
		},
		4: func(agent core.Agent) {
			// Your finishing moves have a 15% chance to grant you an extra combo point.
			// Handled in talents.go.
		},
	},
})

var ItemSetDeathmantle = core.NewItemSet(core.ItemSet{
	Name: "Deathmantle",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Eviscerate and Envenom abilities cause 40 extra damage per combo point.
			// Handled in eviscerate.go.
		},
		4: func(agent core.Agent) {
			// Your attacks have a chance to make your next finishing move cost no energy.
			rogue := agent.(RogueAgent).GetRogue()

			rogue.DeathmantleProcAura = rogue.RegisterAura(core.Aura{
				Label:    "Deathmantle 4pc Proc",
				ActionID: core.ActionID{SpellID: 37171},
				Duration: time.Second * 15,
			})

			ppmm := rogue.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMelee)

			rogue.RegisterAura(core.Aura{
				Label:    "Deathmantle 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if !spellEffect.Landed() {
						return
					}

					// https://tbc.wowhead.com/spell=37170/free-finisher-chance, proc mask = 20.
					if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}

					if !ppmm.Proc(sim, spellEffect.ProcMask, "Deathmantle 4pc") {
						return
					}

					rogue.DeathmantleProcAura.Activate(sim)
				},
			})
		},
	},
})

func (rogue *Rogue) deathmantleActive() bool {
	return rogue.DeathmantleProcAura != nil && rogue.DeathmantleProcAura.IsActive()
}

func (rogue *Rogue) applyDeathmantle(sim *core.Simulation, _ *core.Spell, cast *core.Cast) {
	if rogue.deathmantleActive() {
		cast.Cost = 0
		rogue.DeathmantleProcAura.Deactivate(sim)
	}
}

var ItemSetSlayers = core.NewItemSet(core.ItemSet{
	Name: "Slayer's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the haste from your Slice and Dice ability by 5%.
			// Handled in slice_and_dice.go.
		},
		4: func(agent core.Agent) {
			// Increases the damage dealt by your Backstab, Sinister Strike, Mutilate, and Hemorrhage abilities by 6%.
			// Handled in the corresponding ability files.
		},
	},
})

func init() {
	core.NewItemEffect(30450, func(agent core.Agent) {
		rogue := agent.(RogueAgent).GetRogue()
		procAura := rogue.NewTemporaryStatsAura("Warp Spring Coil Proc", core.ActionID{ItemID: 30450}, stats.Stats{stats.ArmorPenetration: 1000}, time.Second*15)
		const procChance = 0.25

		icd := core.Cooldown{
			Timer:    rogue.NewTimer(),
			Duration: time.Second * 30,
		}

		rogue.RegisterAura(core.Aura{
			Label:    "Warp Spring Coil",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}

				// https://tbc.wowhead.com/spell=37173/armor-penetration, proc mask = 16.
				if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if sim.RandomFloat("WarpSpringCoil") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(32492, func(agent core.Agent) {
		rogue := agent.(RogueAgent).GetRogue()
		procAura := rogue.NewTemporaryStatsAura("Ashtongue Talisman Proc", core.ActionID{ItemID: 32492}, stats.Stats{stats.MeleeCrit: 145}, time.Second*10)

		var numPoints int32

		rogue.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				numPoints = 0
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !spell.Flags.Matches(SpellFlagFinisher) {
					return
				}

				// Need to store the points because they get spent before OnSpellHit is called.
				numPoints = rogue.ComboPoints()

				if spell.SameActionIgnoreTag(rogue.SliceAndDice[1].ActionID) {
					// SND won't call OnSpellHit so we have to add the effect now.
					if numPoints == 5 || sim.RandomFloat("AshtongueTalismanOfLethality") < 0.2*float64(numPoints) {
						procAura.Activate(sim)
					}
				}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spell.Flags.Matches(SpellFlagFinisher) {
					return
				}

				if numPoints == 5 || sim.RandomFloat("AshtongueTalismanOfLethality") < 0.2*float64(numPoints) {
					procAura.Activate(sim)
				}
			},
		})
	})

}
