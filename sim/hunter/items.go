package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var ItemSetBeastLord = core.NewItemSet(core.ItemSet{
	Name: "Beast Lord Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
			// Handled in kill_command.go
		},
	},
})

var ItemSetDemonStalker = core.NewItemSet(core.ItemSet{
	Name: "Demon Stalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
			// Handled in multi_shot.go
		},
	},
})

var ItemSetRiftStalker = core.NewItemSet(core.ItemSet{
	Name: "Rift Stalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
			// Handled in steady_shot.go
		},
	},
})

var ItemSetGronnstalker = core.NewItemSet(core.ItemSet{
	Name: "Gronnstalker's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Handled in rotation.go
		},
		4: func(agent core.Agent) {
			// Handled in steady_shot.go
		},
	},
})

func (hunter *Hunter) talonOfAlarDamageMod(baseDamageConfig core.BaseDamageConfig) core.BaseDamageConfig {
	if hunter.HasTrinketEquipped(30448) {
		return core.WrapBaseDamageConfig(baseDamageConfig, func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
			return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				normalDamage := oldCalculator(sim, hitEffect, spell)
				if hunter.TalonOfAlarAura != nil && hunter.TalonOfAlarAura.IsActive() {
					return normalDamage + 40
				} else {
					return normalDamage
				}
			}
		})
	} else {
		return baseDamageConfig
	}
}

func init() {
	core.NewItemEffect(30448, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()

		procAura := hunter.GetOrRegisterAura(core.Aura{
			Label:    "Talon of Alar Proc",
			ActionID: core.ActionID{ItemID: 30448},
			// Add 1 in case we use arcane shot exactly off CD.
			Duration: time.Second*6 + 1,
		})

		hunter.TalonOfAlarAura = hunter.GetOrRegisterAura(core.Aura{
			Label:    "Talon of Alar",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell == hunter.ArcaneShot {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(30892, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()

		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.03
		hunter.pet.AddStat(stats.MeleeCrit, core.MeleeCritRatingPerCritChance*2)
	})

	core.NewItemEffect(32336, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()
		const manaGain = 8.0
		manaMetrics := hunter.NewManaMetrics(core.ActionID{SpellID: 46939})

		hunter.RegisterAura(core.Aura{
			Label:    "Black Bow of the Betrayer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskRanged) {
					return
				}
				hunter.AddMana(sim, manaGain, manaMetrics, false)
			},
		})
	})

	core.NewItemEffect(32487, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()

		procAura := hunter.NewTemporaryStatsAura("Ashtongue Talisman Proc", core.ActionID{ItemID: 32487}, stats.Stats{stats.AttackPower: 275, stats.RangedAttackPower: 275}, time.Second*8)
		const procChance = 0.15

		hunter.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell != hunter.SteadyShot {
					return
				}
				if sim.RandomFloat("Ashtongue Talisman of Swiftness") > procChance {
					return
				}
				procAura.Activate(sim)
			},
		})
	})

}
