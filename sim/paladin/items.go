package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Tier 6 ret
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
			// Implemented in hammer_of_wrath.go
		},
	},
})

func (paladin *Paladin) getItemSetLightbringerBattlegearBonus4() float64 {
	return core.TernaryFloat64(paladin.HasSetBonus(ItemSetLightbringerBattlegear, 4), .1, 0)
}

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

func (paladin *Paladin) getItemSetRedemptionBattlegearBonus2() float64 {
	return core.TernaryFloat64(paladin.HasSetBonus(ItemSetRedemptionBattlegear, 2), .1, 0)
}

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

func (paladin *Paladin) getItemSetAegisBattlegearBonus2() float64 {
	return core.TernaryFloat64(paladin.HasSetBonus(ItemSetAegisBattlegear, 2), .1, 0)
}

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
			// Implemented in soc.go, sor.go, sov.go
		},
	},
})

func (paladin *Paladin) getItemSetLightswornBattlegearBonus4() float64 {
	return core.TernaryFloat64(paladin.HasSetBonus(ItemSetLightswornBattlegear, 4), .1, 0)
}

// Tier 10 ret
var ItemSetGladiatorsVindication = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Vindication",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.AddStat(stats.AttackPower, 50)
			paladin.AddStat(stats.Resilience, 100)
		},
		4: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.AddStat(stats.AttackPower, 150)
			// Rest implemented in judgement.go
		},
	},
})

func (paladin *Paladin) getItemSetGladiatorsVindicationBonusGloves() float64 {
	hasGloves := (paladin.Equip[proto.ItemSlot_ItemSlotHands].ID == 40798) || // S5a Hateful
		(paladin.Equip[proto.ItemSlot_ItemSlotHands].ID == 40802) || // S5b Hateful
		(paladin.Equip[proto.ItemSlot_ItemSlotHands].ID == 40805) || // S5c Deadly
		(paladin.Equip[proto.ItemSlot_ItemSlotHands].ID == 40808) || // S6 Furious
		(paladin.Equip[proto.ItemSlot_ItemSlotHands].ID == 40812) || // S7 Relentless
		(paladin.Equip[proto.ItemSlot_ItemSlotHands].ID == 51475) // S8 Wrathful
	return core.TernaryFloat64(hasGloves, .05, 0)
}

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
				if spell.Flags.Matches(SpellFlagSecondaryJudgement) {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(37574, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Libram of Furious Blows Proc", core.ActionID{SpellID: 48835}, stats.Stats{stats.MeleeCrit: 61, stats.SpellCrit: 61}, time.Second*5)

		paladin.RegisterAura(core.Aura{
			Label:    "Libram of Furious Blows",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell.Flags.Matches(SpellFlagSecondaryJudgement) {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(40706, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Libram of Reciprocation Proc", core.ActionID{SpellID: 60819}, stats.Stats{stats.MeleeCrit: 173, stats.SpellCrit: 173}, time.Second*10)

		paladin.RegisterAura(core.Aura{
			Label:    "Libram of Reciprocation",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if paladin.CurrentSeal == paladin.SealOfCommandAura && spell.Flags.Matches(SpellFlagSecondaryJudgement) {
					if sim.RandomFloat("Libram of Reciprocation") > 0.15 {
						return
					}
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(33503, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Libram of Divine Judgement Proc", core.ActionID{SpellID: 43745}, stats.Stats{stats.AttackPower: 200}, time.Second*10)

		paladin.RegisterAura(core.Aura{
			Label:    "Libram of Divine Judgement",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if paladin.CurrentSeal == paladin.SealOfCommandAura && spell.Flags.Matches(SpellFlagSecondaryJudgement) {
					if sim.RandomFloat("Libram of Divine Judgement") > 0.40 {
						return
					}
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(42611, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Savage Gladiator's Libram of Fortitude Proc", core.ActionID{SpellID: 60577}, stats.Stats{stats.AttackPower: 94}, time.Second*6)

		paladin.RegisterAura(core.Aura{
			Label:    "Savage Gladiator's Libram of Fortitude",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell.SpellID == paladin.CrusaderStrike.SpellID {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(42851, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Hateful Gladiator's Libram of Fortitude Proc", core.ActionID{SpellID: 60632}, stats.Stats{stats.AttackPower: 106}, time.Second*6)

		paladin.RegisterAura(core.Aura{
			Label:    "Savage Gladiator's Libram of Fortitude",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell.SpellID == paladin.CrusaderStrike.SpellID {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(42852, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Deadly Gladiator's Libram of Fortitude Proc", core.ActionID{SpellID: 60633}, stats.Stats{stats.AttackPower: 120}, time.Second*10)

		paladin.RegisterAura(core.Aura{
			Label:    "Deadly Gladiator's Libram of Fortitude",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell.SpellID == paladin.CrusaderStrike.SpellID {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(42853, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Furious Gladiator's Libram of Fortitude Proc", core.ActionID{SpellID: 60634}, stats.Stats{stats.AttackPower: 144}, time.Second*10)

		paladin.RegisterAura(core.Aura{
			Label:    "Furious Gladiator's Libram of Fortitude",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell.SpellID == paladin.CrusaderStrike.SpellID {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(42854, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Relentless Gladiator's Libram of Fortitude Proc", core.ActionID{SpellID: 60635}, stats.Stats{stats.AttackPower: 172}, time.Second*10)

		paladin.RegisterAura(core.Aura{
			Label:    "Relentless Gladiator's Libram of Fortitude",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell.SpellID == paladin.CrusaderStrike.SpellID {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(51478, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Wrathful Gladiator's Libram of Fortitude Proc", core.ActionID{SpellID: 60636}, stats.Stats{stats.AttackPower: 204}, time.Second*10)

		paladin.RegisterAura(core.Aura{
			Label:    "Wrathful Gladiator's Libram of Fortitude",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell.SpellID == paladin.CrusaderStrike.SpellID {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(50455, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Libram Of Three Truths Proc", core.ActionID{SpellID: 71186}, stats.Stats{stats.Strength: 44}, time.Second*15)
		procAura.MaxStacks = 5

		paladin.RegisterAura(core.Aura{
			Label:    "Libram Of Three Truths",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell.SpellID == paladin.CrusaderStrike.SpellID {
					procAura.Activate(sim)
					procAura.AddStack(sim)
				}
			},
		})
	})

	core.NewItemEffect(47661, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Libram Of Valiance Proc", core.ActionID{SpellID: 67365}, stats.Stats{stats.Strength: 200}, time.Second*15)

		icd := core.Cooldown{
			Timer:    paladin.NewTimer(),
			Duration: time.Second * 8,
		}

		paladin.RegisterAura(core.Aura{
			Label:    "Libram Of Valiance",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				isVengeanceDot := false
				for _, vengeanceDot := range paladin.SealOfVengeanceDots {
					if spell == vengeanceDot.Spell {
						isVengeanceDot = true
					}
				}
				if isVengeanceDot {
					if !icd.IsReady(sim) || sim.RandomFloat("Libram of Valiance") > 0.70 {
						return
					}
					icd.Use(sim)

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
