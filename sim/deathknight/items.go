package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
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

func init() {
	core.NewItemEffect(40715, func(agent core.Agent) {
		deathKnight := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := deathKnight.NewTemporaryStatsAura("Sigil of Haunted Dreams Proc", core.ActionID{ItemID: 40715}, stats.Stats{stats.MeleeCrit: 173.0, stats.SpellCrit: 173.0}, time.Second*10)

		icd := core.Cooldown{
			Timer:    deathKnight.NewTimer(),
			Duration: time.Second * 45.0,
		}

		deathKnight.RegisterAura(core.Aura{
			Label:    "Sigil of Haunted Dreams",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !icd.IsReady(sim) || spell != deathKnight.BloodStrike {
					return
				}

				if sim.RandomFloat("SigilOfHauntedDreams") < 0.15 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

}
