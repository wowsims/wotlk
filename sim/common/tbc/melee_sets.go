package tbc

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Keep these in alphabetical order.

var ItemSetFistsOfFury = core.NewItemSet(core.ItemSet{
	Name: "The Fists of Fury",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			procSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 41989},
				SpellSchool: core.SpellSchoolFire,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultSpellCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, sim.Roll(100, 150), spell.OutcomeMagicHitAndCrit)
				},
			})

			ppmm := character.AutoAttacks.NewPPMManager(2.0, core.ProcMaskMelee)

			character.RegisterAura(core.Aura{
				Label:    "Fists of Fury",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}

					if ppmm.Proc(sim, spell.ProcMask, "The Fists of Fury") {
						procSpell.Cast(sim, result.Target)
					}
				},
			})
		},
	},
})

var ItemSetStormshroud = core.NewItemSet(core.ItemSet{
	Name: "Stormshroud Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(a core.Agent) {
			char := a.GetCharacter()
			proc := char.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 18980},
				SpellSchool: core.SpellSchoolNature,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				CritMultiplier:   char.DefaultSpellCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, sim.Roll(15, 25), spell.OutcomeMagicHitAndCrit)
				},
			})
			char.RegisterAura(core.Aura{
				Label:    "Stormshround Armor 2pc",
				ActionID: core.ActionID{SpellID: 18979},
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
					if sim.RandomFloat("Stormshroud Armor 2pc") < 0.05 {
						proc.Cast(sim, result.Target)
					}
				},
			})
		},
		3: func(a core.Agent) {
			char := a.GetCharacter()
			if !char.HasEnergyBar() {
				return
			}
			metrics := char.NewEnergyMetrics(core.ActionID{SpellID: 23863})
			proc := char.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 23864},
				SpellSchool: core.SpellSchoolNature,
				ApplyEffects: func(sim *core.Simulation, u *core.Unit, spell *core.Spell) {
					char.AddEnergy(sim, 30, metrics)
				},
			})
			char.RegisterAura(core.Aura{
				Label:    "Stormshround Armor 3pc",
				ActionID: core.ActionID{SpellID: 18979},
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
					if sim.RandomFloat("Stormshroud Armor 3pc") < 0.02 {
						proc.Cast(sim, result.Target)
					}
				},
			})

		},
		4: func(a core.Agent) {
			a.GetCharacter().AddStat(stats.AttackPower, 14)
		},
	},
})

var ItemSetTwinBladesOfAzzinoth = core.NewItemSet(core.ItemSet{
	Name: "The Twin Blades of Azzinoth",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
				character.PseudoStats.MobTypeAttackPower += 200
			}
			procAura := character.NewTemporaryStatsAura("Twin Blade of Azzinoth Proc", core.ActionID{SpellID: 41435}, stats.Stats{stats.MeleeHaste: 450}, time.Second*10)

			ppmm := character.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMelee)
			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Second * 45,
			}

			character.RegisterAura(core.Aura{
				Label:    "Twin Blades of Azzinoth",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}

					if !icd.IsReady(sim) {
						return
					}

					if ppmm.Proc(sim, spell.ProcMask, "Twin Blades of Azzinoth") {
						icd.Use(sim)
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})
