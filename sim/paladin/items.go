package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var ItemSetJusticarBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Justicar Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// sim/debuffs.go handles this (and paladin/judgement.go)
		},
		4: func(agent core.Agent) {
			// TODO: if we ever implemented judgement of command, add bonus from 4p
		},
	},
})

var ItemSetJusticarArmor = core.NewItemSet(core.ItemSet{
	Name: "Justicar Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage dealt by your Seal of Righteousness, Seal of
			// Vengeance, and Seal of Blood by 10%.
			// Implemented in seals.go.
		},
		4: func(agent core.Agent) {
			// Increases the damage dealt by Holy Shield by 15.
			// Implemented in holy_shield.go.
		},
	},
})

var ItemSetCrystalforgeBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Crystalforge Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// judgement.go
		},
		4: func(agent core.Agent) {
			// TODO: if we implement healing, this heals party.
		},
	},
})

var ItemSetCrystalforgeArmor = core.NewItemSet(core.ItemSet{
	Name: "Crystalforge Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage from your Retribution Aura by 15.
			// TODO
		},
		4: func(agent core.Agent) {
			// Each time you use your Holy Shield ability, you gain 100 Block Value
			// against a single attack in the next 6 seconds.
			paladin := agent.(PaladinAgent).GetPaladin()

			procAura := paladin.RegisterAura(core.Aura{
				Label:    "Crystalforge 2pc Proc",
				ActionID: core.ActionID{SpellID: 37191},
				Duration: time.Second * 6,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					paladin.AddStatDynamic(sim, stats.BlockValue, 100)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					paladin.AddStatDynamic(sim, stats.BlockValue, -100)
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.Outcome.Matches(core.OutcomeBlock) {
						aura.Deactivate(sim)
					}
				},
			})

			paladin.RegisterAura(core.Aura{
				Label:    "Crystalforge 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == paladin.HolyShield {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetLightbringerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Lightbringer Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			manaMetrics := paladin.NewManaMetrics(core.ActionID{SpellID: 38428})

			paladin.RegisterAura(core.Aura{
				Label:    "Lightbringer Battlegear 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
					if sim.RandomFloat("lightbringer 2pc") > 0.2 {
						return
					}
					paladin.AddMana(sim, 50, manaMetrics, true)
				},
			})
		},
		4: func(agent core.Agent) {
			// TODO: if we implemented hammer of wrath.. this ups dmg
		},
	},
})

var ItemSetLightbringerArmor = core.NewItemSet(core.ItemSet{
	Name: "Lightbringer Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the mana gained from your Spiritual Attunement ability by 10%.
		},
		4: func(agent core.Agent) {
			// Increases the damage dealt by Consecration by 10%.
		},
	},
})

func init() {
	// Librams implemented in seals.go and judgement.go

	// TODO: once we have judgement of command.. https://tbc.wowhead.com/item=33503/libram-of-divine-judgement

	core.NewItemEffect(27484, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Libram of Avengement Proc", core.ActionID{SpellID: 34260}, stats.Stats{stats.MeleeCrit: 53, stats.SpellCrit: 53}, time.Second*5)

		paladin.RegisterAura(core.Aura{
			Label:    "Libram of Avengement",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell == paladin.JudgementOfBlood || spell == paladin.JudgementOfRighteousness {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(32368, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Tome of the Lightbringer Proc", core.ActionID{SpellID: 41042}, stats.Stats{stats.BlockValue: 186}, time.Second*5)

		paladin.RegisterAura(core.Aura{
			Label:    "Tome of the Lightbringer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell.Flags.Matches(SpellFlagJudgement) {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(30447, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Tome of Fiery Redemption Proc", core.ActionID{ItemID: 30447}, stats.Stats{stats.SpellPower: 290}, time.Second*15)

		icd := core.Cooldown{
			Timer:    paladin.NewTimer(),
			Duration: time.Second * 45,
		}

		paladin.RegisterAura(core.Aura{
			Label:    "Tome of Fiery Redemption",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !spell.Flags.Matches(SpellFlagSeal|SpellFlagJudgement) && spell.SpellSchool != core.SpellSchoolPhysical {
					return
				}
				if !icd.IsReady(sim) || sim.RandomFloat("TomeOfFieryRedemption") > 0.15 {
					return
				}
				icd.Use(sim)

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(32489, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		actionID := core.ActionID{ItemID: 32489}

		dotSpell := paladin.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
		})

		target := paladin.CurrentTarget
		judgementDot := core.NewDot(core.Dot{
			Spell: dotSpell,
			Aura: target.RegisterAura(core.Aura{
				Label:    "AshtongueTalismanOfZeal-" + strconv.Itoa(int(paladin.Index)),
				ActionID: actionID,
			}),
			NumberOfTicks: 4,
			TickLength:    time.Second * 2,
			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:         core.ProcMaskPeriodicDamage,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				BaseDamage:     core.BaseDamageConfigFlat(480 / 4),
				OutcomeApplier: paladin.OutcomeFuncTick(),
				IsPeriodic:     true,
			}),
		})

		paladin.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman of Zeal",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell.Flags.Matches(SpellFlagJudgement) && sim.RandomFloat("AshtongueTalismanOfZeal") < 0.5 {
					judgementDot.Apply(sim)
				}
			},
		})
	})

}
