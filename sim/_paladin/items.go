package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
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
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !spell.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
					if sim.RandomFloat("lightbringer 2pc") > 0.2 {
						return
					}
					paladin.AddMana(sim, 50, manaMetrics)
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

// Tier 9 ret (Alliance/Horde)
var ItemSetTuralyonsBattlegear = core.NewItemSet(core.ItemSet{
	Name:            "Turalyon's Battlegear",
	AlternativeName: "Liadrin's Battlegear",
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
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
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

// PvP ret
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

// Tier 7 prot
var ItemSetRedemptionPlate = core.NewItemSet(core.ItemSet{
	Name: "Redemption Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in hammer_of_the_righteous.go
		},
		4: func(agent core.Agent) {
			// TODO: increase duration of divine shield by 3sec
			// Implemented in divine_protection.go
		},
	},
})

func (paladin *Paladin) getItemSetRedemptionPlateBonus2() float64 {
	return core.TernaryFloat64(paladin.HasSetBonus(ItemSetRedemptionPlate, 2), .1, 0)
}

// Tier 8 prot
var ItemSetAegisPlate = core.NewItemSet(core.ItemSet{
	Name: "Aegis Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in sov.go
		},
		4: func(agent core.Agent) {
			// Implemented in shield_of_righteousness.go
		},
	},
})

func (paladin *Paladin) getItemSetAegisPlateBonus2() float64 {
	return core.TernaryFloat64(paladin.HasSetBonus(ItemSetAegisPlate, 2), .1, 0)
}

// Tier 9 prot (Alliance/Horde)
var ItemSetTuralyonsPlate = core.NewItemSet(core.ItemSet{
	Name:            "Turalyon's Plate",
	AlternativeName: "Liadrin's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in hammer_of_the_righteous.go
			// TODO: Implement Hand of Reckoning bonus, if it ever becomes relevant
		},
		4: func(agent core.Agent) {
			// Implemented in divine_protection.go
		},
	},
})

func (paladin *Paladin) getItemSetT9PlateBonus2() float64 {
	return core.TernaryFloat64(paladin.HasSetBonus(ItemSetTuralyonsPlate, 2), .05, 0)
}

// Tier 10 prot
var ItemSetLightswornPlate = core.NewItemSet(core.ItemSet{
	Name: "Lightsworn Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in hammer_of_the_righteous.go
		},
		4: func(agent core.Agent) {
			// TODO: When you activate Divine Plea, you gain 12% dodge for 10 sec
		},
	},
})

func (paladin *Paladin) getItemSetLightswornPlateBonus2() float64 {
	return core.TernaryFloat64(paladin.HasSetBonus(ItemSetLightswornPlate, 2), .2, 0)
}

func (paladin *Paladin) getItemSetGladiatorsVindicationBonusGloves() float64 {
	switch paladin.Hands().ID {
	case 40798, 40802, 40805, 40808, 40812, 51475: // S5a Hateful, S5b Hateful, S5c Deadly, S6 Furious, S7 Relentless, S8 Wrathful
		return 0.05
	default:
		return 0
	}
}

func init() {
	// Librams implemented in seals.go and judgement.go

	core.NewItemEffect(37574, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Libram of Furious Blows Proc", core.ActionID{SpellID: 48835}, stats.Stats{stats.MeleeCrit: 61, stats.SpellCrit: 61}, time.Second*5)

		paladin.RegisterAura(core.Aura{
			Label:    "Libram of Furious Blows",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if paladin.CurrentSeal == paladin.SealOfCommandAura && spell.Flags.Matches(SpellFlagSecondaryJudgement) {
					if sim.RandomFloat("Libram of Reciprocation") > 0.15 {
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.SpellID == paladin.CrusaderStrike.SpellID {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(50455, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()

		procAura := core.MakeStackingAura(paladin.GetCharacter(), core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Formidable",
				ActionID:  core.ActionID{SpellID: 71187},
				Duration:  time.Second * 15,
				MaxStacks: 5,
			},
			BonusPerStack: stats.Stats{stats.Strength: 44},
		})

		paladin.RegisterAura(core.Aura{
			Label:    "Libram Of Three Truths",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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
		procAura.Icd = &icd

		paladin.RegisterAura(core.Aura{
			Label:    "Libram Of Valiance",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == paladin.SovDotSpell {
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
		procAura := paladin.NewTemporaryStatsAura("Tome of the Lightbringer Proc", core.ActionID{SpellID: 41042}, stats.Stats{stats.BlockValue: 186}, time.Second*10)

		paladin.RegisterAura(core.Aura{
			Label:    "Tome of the Lightbringer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(40707, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Libram of Obstruction Proc", core.ActionID{SpellID: 60794}, stats.Stats{stats.BlockValue: 352}, time.Second*10)

		paladin.RegisterAura(core.Aura{
			Label:    "Libram of Obstruction",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(45145, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Libram of the Sacred Shield Proc", core.ActionID{SpellID: 65182}, stats.Stats{stats.BlockValue: 450}, time.Second*20)

		paladin.RegisterAura(core.Aura{
			Label:    "Libram of the Sacred Shield",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.SpellID == paladin.HolyShield.SpellID {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(32489, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()

		// The spell effect is https://www.wowhead.com/wotlk/spell=40472/enduring-judgement, most likely
		dotSpell := paladin.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{ItemID: 32489},
			SpellSchool:      core.SpellSchoolHoly,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "AshtongueTalismanOfZeal",
				},
				NumberOfTicks: 4,
				TickLength:    time.Second * 2,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = 480 / 4
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},
		})

		paladin.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman of Zeal",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Flags.Matches(SpellFlagPrimaryJudgement) && sim.RandomFloat("AshtongueTalismanOfZeal") < 0.5 {
					dotSpell.Dot(result.Target).Apply(sim)
				}
			},
		})
	})

}
