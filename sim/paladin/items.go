package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Tier 7 ret
var ItemSetRedemptionBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Redemption Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in divine_storm.go
		},
		4: func(agent core.Agent) {
			// Implemented in judgement.go
		},
	},
})

// Tier 8 ret
var ItemSetAegisBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Aegis Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in exorcism.go & hammer_of_wrath.go
		},
		4: func(agent core.Agent) {
			// Implemented in divine_storm.go & crusader_strike.go
		},
	},
})

// Tier 9 ret (Alliance)
var ItemSetTuralyonsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Turalyon's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in talents.go (Righteous Vengeance)
		},
		4: func(agent core.Agent) {
			// Implemented in soc.go, sor.go, sov.go
		},
	},
})

// Tier 9 ret (Horde)
var ItemSetLiadrinsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Liadrin's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in talents.go (Righteous Vengeance)
		},
		4: func(agent core.Agent) {
			// Implemented in soc.go, sor.go, sov.go
		},
	},
})

// Tier 10 ret
var ItemSetLightswornBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Lightsworn Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			procSpell := paladin.RegisterSpell(core.SpellConfig{
				ActionID: core.ActionID{SpellID: 70765},
				ApplyEffects: func(_ *core.Simulation, _ *core.Unit, _ *core.Spell) {
					paladin.DivineStorm.CD.Reset()
				},
			})

			paladin.RegisterAura(core.Aura{
				Label:    "Lightsworn Battlegear 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
						return
					}
					if sim.RandomFloat("lightsworn 2pc") > 0.4 {
						return
					}
					procSpell.Cast(sim, &paladin.Unit)
				},
			})
		},
		4: func(agent core.Agent) {
			// Implemented in seals.go
		},
	},
})

func init() {
	// Librams implemented in seals.go and judgement.go

	// TODO: once we have judgement of command.. https://wotlk.wowhead.com/item=33503/libram-of-divine-judgement

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
				// if spell == paladin.JudgementOfBlood || spell == paladin.JudgementOfRighteousness {
				procAura.Activate(sim)
				// }
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
				if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
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
				if !spell.Flags.Matches(SpellFlagSecondaryJudgement|SpellFlagPrimaryJudgement) && spell.SpellSchool != core.SpellSchoolPhysical {
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
				if spell.Flags.Matches(SpellFlagPrimaryJudgement) && sim.RandomFloat("AshtongueTalismanOfZeal") < 0.5 {
					judgementDot.Apply(sim)
				}
			},
		})
	})

}
