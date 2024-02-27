package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/common/wotlk"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var ItemSetGladiatorsPursuit = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.AddStats(stats.Stats{
				stats.AttackPower:       50,
				stats.RangedAttackPower: 50,
				stats.Resilience:        50,
			})
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.AddStats(stats.Stats{
				stats.AttackPower:       150,
				stats.RangedAttackPower: 150,
			})
		},
	},
})

var ItemSetCryptstalkerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Cryptstalker Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			if hunter.pet != nil {
				hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.05
			}
		},
		4: func(agent core.Agent) {
		},
	},
})

var ItemSetScourgestalkerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Scourgestalker Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()

			procAura := hunter.NewTemporaryStatsAura("Scourgestalker 4pc Proc", core.ActionID{SpellID: 64860}, stats.Stats{stats.AttackPower: 600, stats.RangedAttackPower: 600}, time.Second*15)
			const procChance = 0.1

			icd := core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 45,
			}
			procAura.Icd = &icd

			hunter.RegisterAura(core.Aura{
				Label:    "Windrunner 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || spell != hunter.SteadyShot {
						return
					}
					if !icd.IsReady(sim) {
						return
					}
					if sim.RandomFloat("Scourgestalker 4pc") > procChance {
						return
					}

					icd.Use(sim)
					procAura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetWindrunnersPursuit = core.NewItemSet(core.ItemSet{
	Name:            "Windrunner's Pursuit",
	AlternativeName: "Windrunner's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			if hunter.pet == nil {
				return
			}

			procAura := hunter.pet.NewTemporaryStatsAura("Windrunner 4pc Proc", core.ActionID{SpellID: 68130}, stats.Stats{stats.AttackPower: 600}, time.Second*15)
			const procChance = 0.35

			icd := core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 45,
			}
			procAura.Icd = &icd

			hunter.RegisterAura(core.Aura{
				Label:    "Windrunner 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskRanged) {
						return
					}
					if !icd.IsReady(sim) {
						return
					}
					if sim.RandomFloat("Windrunner 4pc") > procChance {
						return
					}

					icd.Use(sim)
					procAura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetAhnKaharBloodHuntersBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Ahn'Kahar Blood Hunter's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			const procChance = 0.05
			actionID := core.ActionID{SpellID: 70727}

			procAura := hunter.RegisterAura(core.Aura{
				Label:    "AhnKahar 2pc Proc",
				ActionID: actionID,
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.15
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.15
				},
			})

			hunter.RegisterAura(core.Aura{
				Label:    "AhnKahar 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == hunter.AutoAttacks.RangedAuto() && sim.RandomFloat("AhnKahar 2pc") < procChance {
						procAura.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			const procChance = 0.05
			actionID := core.ActionID{SpellID: 70730}

			var curBonus stats.Stats
			procAura := hunter.RegisterAura(core.Aura{
				Label:    "AhnKahar 4pc Proc",
				ActionID: actionID,
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					curBonus = stats.Stats{
						stats.AttackPower:       aura.Unit.GetStat(stats.AttackPower) * 0.1,
						stats.RangedAttackPower: aura.Unit.GetStat(stats.RangedAttackPower) * 0.1,
					}

					aura.Unit.AddStatsDynamic(sim, curBonus)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatsDynamic(sim, curBonus.Invert())
				},
			})

			hunter.RegisterAura(core.Aura{
				Label:    "AhnKahar 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == hunter.SerpentSting && sim.RandomFloat("AhnKahar 4pc") < procChance {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

func init() {

	wotlk.NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Zod's Repeating Longbow"
		itemID := int32(50034)
		procChance := 0.04
		if isHeroic {
			name += " H"
			itemID = 50638
			procChance = 0.05
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			if agent.GetCharacter().Class != proto.Class_ClassHunter {
				return
			}

			hunter := agent.(HunterAgent).GetHunter()

			var rangedSpell *core.Spell
			initSpell := func() {
				rangedSpell = hunter.RegisterSpell(core.SpellConfig{
					ActionID:    core.ActionID{ItemID: itemID},
					SpellSchool: core.SpellSchoolPhysical,
					ProcMask:    core.ProcMaskRangedAuto,
					Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

					DamageMultiplier: 0.5,
					CritMultiplier:   hunter.AutoAttacks.RangedConfig().CritMultiplier,
					ThreatMultiplier: 1,

					ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						baseDamage := hunter.RangedWeaponDamage(sim, spell.RangedAttackPower(target)) +
							hunter.AmmoDamageBonus +
							spell.BonusWeaponDamage()

						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
					},
				})
			}

			triggerAura := core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
				Name:       name + " Trigger",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskRanged,
				Outcome:    core.OutcomeLanded,
				ProcChance: procChance,
				ActionID:   core.ActionID{ItemID: itemID},
				Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					rangedSpell.Cast(sim, result.Target)
				},
			})
			triggerAura.OnInit = func(aura *core.Aura, sim *core.Simulation) {
				initSpell()
			}
		})
	})

}
