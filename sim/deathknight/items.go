package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO: T9 Dps missing Heart Strike
// TODO: T9 Tank missing Heart Strike and Vampiric Blood and Dark Command
// TODO: T10 Dps missing Heart Strike

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

func (dk *Deathknight) scourgeborneBattlegearCritBonus() float64 {
	return core.TernaryFloat64(dk.HasSetBonus(ItemSetScourgeborneBattlegear, 2), 5.0, 0.0)
}

func (dk *Deathknight) scourgeborneBattlegearRunicPowerBonus() float64 {
	return core.TernaryFloat64(dk.HasSetBonus(ItemSetScourgeborneBattlegear, 4), 5.0, 0.0)
}

var ItemSetScourgebornePlate = core.NewItemSet(core.ItemSet{
	Name: "Scourgeborne Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the critical strike chance of your plague
			// strike by 10%
		},
		4: func(agent core.Agent) {
			// Increases the duration of your Icebound Fortitude by 3 seconds
		},
	},
})

func (dk *Deathknight) scourgebornePlateCritBonus() float64 {
	return core.TernaryFloat64(dk.HasSetBonus(ItemSetScourgebornePlate, 2), 10.0, 0.0)
}

func (dk *Deathknight) scourgebornePlateIFDurationBonus() time.Duration {
	return core.TernaryDuration(dk.HasSetBonus(ItemSetScourgebornePlate, 4), 3*time.Second, 0*time.Second)
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

func (dk *Deathknight) darkrunedBattlegearCritBonus() float64 {
	return core.TernaryFloat64(dk.HasSetBonus(ItemSetDarkrunedBattlegear, 2), 8.0, 0.0)
}

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

func (dk *Deathknight) darkrunedPlateRuneStrikeDamageBonus() float64 {
	return core.TernaryFloat64(dk.HasSetBonus(ItemSetDarkrunedPlate, 2), 1.1, 1.0)
}

func (dk *Deathknight) darkrunedPlateAMSBonus() float64 {
	return core.TernaryFloat64(dk.HasSetBonus(ItemSetDarkrunedPlate, 4), 0.9, 1.0)
}

var ItemSetThassariansBattlegear = core.NewItemSet(core.ItemSet{
	Name:            "Thassarian's Battlegear",
	AlternativeName: "Koltira's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Blood Strike and Heart Strike abilities have a
			// chance to grant you 180 additional strength for 15 sec.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			dk.registerThassariansBattlegearProc()
		},
		4: func(agent core.Agent) {
			// Your Blood Plague ability now has a chance for its
			// damage to be critical strikes.
		},
	},
})

func (dk *Deathknight) registerThassariansBattlegearProc() {
	procAura := dk.NewTemporaryStatsAura("Unholy Might Proc", core.ActionID{SpellID: 67117}, stats.Stats{stats.Strength: 180.0}, time.Second*15)

	icd := core.Cooldown{
		Timer:    dk.NewTimer(),
		Duration: time.Second * 45.0,
	}
	procAura.Icd = &icd

	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label: "Unholy Might",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !icd.IsReady(sim) {
				return
			}

			if spell != dk.BloodStrike && spell != dk.HeartStrike {
				return
			}

			if sim.RandomFloat("UnholyMight") < 0.5 {
				icd.Use(sim)
				procAura.Activate(sim)
			}
		},
	}))
}

var ItemSetThassariansPlate = core.NewItemSet(core.ItemSet{
	Name:            "Thassarian's Plate",
	AlternativeName: "Koltira's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Decreases the cooldown on your Dark Command ability by 2 sec and
			// increases the damage done by your Blood Strike and Heart Strike abilities by 5%.
		},
		4: func(agent core.Agent) {
			// Decreases the cooldown on your Unbreakable Armor, Vampiric Blood
			// and Bone Shield abilities by 10 sec.
		},
	},
})

func (dk *Deathknight) thassariansPlateDamageBonus() float64 {
	return core.TernaryFloat64(dk.HasSetBonus(ItemSetThassariansPlate, 2), 1.05, 1.0)
}

func (dk *Deathknight) thassariansPlateCooldownReduction(spell *core.Spell) time.Duration {
	if !dk.HasSetBonus(ItemSetThassariansPlate, 4) {
		return 0
	}

	if spell == dk.UnbreakableArmor || spell == dk.BoneShield || spell == dk.VampiricBlood {
		return 10 * time.Second
	}
	return 0
}

var ItemSetScourgelordsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Scourgelord's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Obliterate and Scourge Strike abilities deal 10% increased damage
			// and your Heart Strike ability deals 7% increased damage.
		},
		4: func(agent core.Agent) {
			// Whenever all your runes are on cooldown, you gain 3% increased
			// damage done with weapons, spells, and abilities for the next 15 sec.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			dk.registerScourgelordsBattlegearProc()
		},
	},
})

type ScourgelordBonusSpell int8

const (
	ScourgelordBonusSpellOB = iota + 1
	ScourgelordBonusSpellSS
	ScourgelordBonusSpellHS
)

func (dk *Deathknight) scourgelordsBattlegearDamageBonus(spell ScourgelordBonusSpell) float64 {
	if !dk.HasSetBonus(ItemSetScourgelordsBattlegear, 2) {
		return 1.0
	}

	if spell == ScourgelordBonusSpellOB || spell == ScourgelordBonusSpellSS {
		return 1.1
	} else if spell == ScourgelordBonusSpellHS {
		return 1.07
	}
	return 1.0
}

func (dk *Deathknight) registerScourgelordsBattlegearProc() {
	bonusCoeff := 1.03

	damageAura := dk.RegisterAura(core.Aura{
		Label:    "Advantage",
		ActionID: core.ActionID{SpellID: 70657},
		Duration: 15 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= bonusCoeff
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= bonusCoeff
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.DefaultCast.GCD > 0 && dk.AllRunesSpent() {
				aura.Refresh(sim)
			}
		},
	})

	dk.onRuneSpendT10 = func(sim *core.Simulation, changeType core.RuneChangeType) {
		if changeType.Matches(core.SpendRune) && dk.AllRunesSpent() {
			damageAura.Activate(sim)
		}
	}
}

var ItemSetScourgelordsPlate = core.NewItemSet(core.ItemSet{
	Name: "Scourgelord's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage done by your Death and Decay ability by 20%.
		},
		4: func(agent core.Agent) {
			// When you activate Blood Tap, you gain 12% damage reduction from all attacks for 10 sec.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			dk.registerScourgelordsPlateProc()
		},
	},
})

func (dk *Deathknight) scourgelordsPlateDamageBonus() float64 {
	return core.TernaryFloat64(dk.HasSetBonus(ItemSetScourgelordsPlate, 2), 1.2, 1.0)
}

func (dk *Deathknight) registerScourgelordsPlateProc() {
	damageTakenMult := 0.88

	bonusAura := dk.RegisterAura(core.Aura{
		Label:    "Blood Armor Proc",
		ActionID: core.ActionID{SpellID: 70654},
		Duration: time.Second * 10.0,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTakenMult
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTakenMult
		},
	})

	core.MakePermanent(dk.RegisterAura(core.Aura{
		Label: "Blood Armor",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == dk.BloodTap {
				bonusAura.Activate(sim)
			}
		},
	}))
}

func (dk *Deathknight) sigilOfTheDarkRiderBonus() float64 {
	return core.TernaryFloat64(dk.Ranged().ID == 39208, 90, 0)
}

func (dk *Deathknight) sigilOfAwarenessBonus() float64 {
	return core.TernaryFloat64(dk.Ranged().ID == 40207, 420, 0)
}

func (dk *Deathknight) sigilOfTheFrozenConscienceBonus() float64 {
	return core.TernaryFloat64(dk.Ranged().ID == 40822, 111, 0)
}

func (dk *Deathknight) sigilOfTheWildBuckBonus() float64 {
	return core.TernaryFloat64(dk.Ranged().ID == 40867, 80, 0)
}

func (dk *Deathknight) sigilOfArthriticBindingBonus() float64 {
	return core.TernaryFloat64(dk.Ranged().ID == 40875, 203, 0)
}

func (dk *Deathknight) sigilOfTheVengefulHeartDeathCoil() float64 {
	return core.TernaryFloat64(dk.Ranged().ID == 45254, 403, 0)
}

func (dk *Deathknight) sigilOfTheVengefulHeartFrostStrike() float64 {
	return core.TernaryFloat64(dk.Ranged().ID == 45254, 218, 0) // (1 / 0.55) * 120
}

func addEnchantEffect(id int32, effect func(core.Agent)) {
	if core.HasEnchantEffect(id) {
		return
	}
	core.NewEnchantEffect(id, effect)
}

func addItemEffect(id int32, effect func(core.Agent)) {
	if core.HasItemEffect(id) {
		return
	}
	core.NewItemEffect(id, effect)
}

func CreateVirulenceProcAura(character *core.Character) *core.Aura {
	return character.NewTemporaryStatsAura("Sigil of Virulence Proc", core.ActionID{SpellID: 67383}, stats.Stats{stats.Strength: 200.0}, time.Second*20)
}

func (dk *Deathknight) registerItems() {
	// Rune of Razorice
	newRazoriceHitSpell := func(character *core.Character, isMH bool) *core.Spell {
		dmg := 0.0

		if weapon := character.GetMHWeapon(); isMH && weapon != nil {
			dmg = 0.5 * (weapon.WeaponDamageMin + weapon.WeaponDamageMax) * 0.02
		} else if weapon := character.GetOHWeapon(); !isMH && weapon != nil {
			dmg = 0.5 * (weapon.WeaponDamageMin + weapon.WeaponDamageMax) * 0.02
		} else {
			return nil
		}

		return character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 50401},
			SpellSchool: core.SpellSchoolFrost,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeAlwaysHit)
			},
		})
	}

	addEnchantEffect(3370, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 50401}
		if spell := character.GetSpell(actionID); spell != nil {
			// This function gets called twice when dual wielding this enchant, but we
			// handle both in one call.
			return
		}

		procMask := character.GetProcMaskForEnchant(3370)

		vulnAuras := character.NewEnemyAuraArray(core.RuneOfRazoriceVulnerabilityAura)
		mhRazoriceSpell := newRazoriceHitSpell(character, true)
		ohRazoriceSpell := newRazoriceHitSpell(character, false)
		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Razor Frost",
			ActionID: core.ActionID{SpellID: 50401},
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}

				vulnAura := vulnAuras.Get(result.Target)
				vulnAura.Activate(sim)
				if spell.IsMH() {
					mhRazoriceSpell.Cast(sim, result.Target)
					vulnAura.AddStack(sim)
				} else {
					ohRazoriceSpell.Cast(sim, result.Target)
					vulnAura.AddStack(sim)
				}
			},
		})

		character.RegisterOnItemSwap(func(sim *core.Simulation) {
			if character.GetProcMaskForEnchant(3370) == core.ProcMaskUnknown {
				aura.Deactivate(sim)
			} else {
				aura.Activate(sim)
			}
		})
	})

	// Rune of the Fallen Crusader
	newRuneOfTheFallenCrusaderAura := func(character *core.Character, auraLabel string, actionID core.ActionID) *core.Aura {
		return character.NewTemporaryStatsAuraWrapped(auraLabel, actionID, stats.Stats{}, time.Second*15, func(aura *core.Aura) {
			statDep := character.NewDynamicMultiplyStat(stats.Strength, 1.15)

			aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.EnableDynamicStatDep(sim, statDep)
			})

			aura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.DisableDynamicStatDep(sim, statDep)
			})
		})
	}

	// ApplyRuneOfTheFallenCrusader will be applied twice if there is two weapons with this enchant.
	//   However, it will automatically overwrite one of them, so it should be ok.
	//   A single application of the aura will handle both mh and oh procs.
	addEnchantEffect(3368, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(3368)
		ppmm := character.AutoAttacks.NewPPMManager(2.0, procMask)

		rfcAura := newRuneOfTheFallenCrusaderAura(character, "Rune Of The Fallen Crusader Proc", core.ActionID{SpellID: 53365})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Rune Of The Fallen Crusader",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "rune of the fallen crusader") {
					rfcAura.Activate(sim)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3368, 2.0, &ppmm, aura)
	})

	// Rune of the Nerubian Carapace
	addEnchantEffect(3883, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.AddStat(stats.Defense, 13*core.DefenseRatingPerDefense)
		character.MultiplyStat(stats.Stamina, 1.01)
	})

	// Rune of the Stoneskin Gargoyle
	addEnchantEffect(3847, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.AddStat(stats.Defense, 25*core.DefenseRatingPerDefense)
		character.MultiplyStat(stats.Stamina, 1.02)
	})

	// Rune of the Swordbreaking
	addEnchantEffect(3594, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.AddStat(stats.Parry, 2*core.ParryRatingPerParryChance)
	})

	// Rune of Swordshattering
	addEnchantEffect(3365, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.AddStat(stats.Parry, 4*core.ParryRatingPerParryChance)
	})

	// Rune of the Spellbreaking
	addEnchantEffect(3595, func(agent core.Agent) {
		// TODO:
		// Add 2% magic deflection
	})

	// Rune of Spellshattering
	addEnchantEffect(3367, func(agent core.Agent) {
		// TODO:
		// Add 4% magic deflection
	})

	cinderBonusCoeff := 1.2

	consumeSpells := [5]core.ActionID{
		BloodBoilActionID,
		DeathCoilActionID,
		FrostStrikeMHActionID,
		HowlingBlastActionID,
		IcyTouchActionID,
	}

	targetsHit := 0

	dk.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 53386},
		Label:     "Cinderglacier",
		Duration:  time.Second * 30,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.SetStacks(sim, aura.MaxStacks)
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= cinderBonusCoeff
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= cinderBonusCoeff
			dk.modifyShadowDamageModifier(0.2)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= cinderBonusCoeff
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] /= cinderBonusCoeff
			dk.modifyShadowDamageModifier(-0.2)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ActionID == HowlingBlastActionID || spell.ActionID == BloodBoilActionID {
				if result.Target.Index == 0 {
					targetsHit = 0
				}
				if result.Landed() {
					targetsHit++
				}
				if result.Target.Index == sim.GetNumTargets()-1 {
					// Last target, consume a stack for every target hit
					for i := 0; i < targetsHit; i++ {
						if aura.IsActive() {
							aura.RemoveStack(sim)
						}
					}
				}
				return
			}

			if !result.Landed() {
				return
			}

			shouldConsume := false
			for _, consumeSpell := range consumeSpells {
				if spell.ActionID == consumeSpell {
					shouldConsume = true
					break
				}
			}

			if shouldConsume {
				aura.RemoveStack(sim)
			}
		},
	})

	// Rune of Cinderglacier
	addEnchantEffect(3369, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(3369)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)
		// have to fetch it dynamically, otherwise aura reference becomes stale? not quite sure why
		proc := character.GetAura("Cinderglacier")

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: "Rune of Cinderglacier",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "rune of cinderglacier") {
					proc.Activate(sim)
				}
			},
		}))
	})

	// Sigils

	addItemEffect(40714, func(agent core.Agent) {
		dk := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := dk.NewTemporaryStatsAura("Sigil of the Unfaltering Knight Proc", core.ActionID{SpellID: 62146}, stats.Stats{stats.Defense: 53.0 / core.DefenseRatingPerDefense}, time.Second*30)

		core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
			Label: "Sigil of the Unfaltering Knight",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell != dk.IcyTouch {
					return
				}

				procAura.Activate(sim)
			},
		}))
	})

	addItemEffect(40715, func(agent core.Agent) {
		dk := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := dk.NewTemporaryStatsAura("Sigil of Haunted Dreams Proc", core.ActionID{SpellID: 60828}, stats.Stats{stats.MeleeCrit: 173.0, stats.SpellCrit: 173.0}, time.Second*10)

		icd := core.Cooldown{
			Timer:    dk.NewTimer(),
			Duration: time.Second * 45.0,
		}
		procAura.Icd = &icd

		core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
			Label: "Sigil of Haunted Dreams",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) || (spell != dk.BloodStrike && spell != dk.HeartStrike) {
					return
				}

				if sim.RandomFloat("SigilOfHauntedDreams") < 0.15 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		}))
	})

	addItemEffect(45144, func(agent core.Agent) {
		dk := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := dk.NewTemporaryStatsAura("Sigil of Deflection Proc", core.ActionID{SpellID: 64963}, stats.Stats{stats.Dodge: 144.0}, time.Second*5)

		core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
			Label: "Sigil of Deflection",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell != dk.RuneStrike {
					return
				}

				procAura.Activate(sim)
			},
		}))
	})

	addItemEffect(47672, func(agent core.Agent) {
		dk := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := dk.NewTemporaryStatsAura("Sigil of Insolence Proc", core.ActionID{SpellID: 67380}, stats.Stats{stats.Dodge: 200.0}, time.Second*20)

		icd := core.Cooldown{
			Timer:    dk.NewTimer(),
			Duration: time.Second * 10.0,
		}
		procAura.Icd = &icd

		core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
			Label: "Sigil of Insolence",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) || spell != dk.RuneStrike {
					return
				}

				if sim.RandomFloat("SigilOfInsolence") < 0.80 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		}))
	})

	addItemEffect(47673, func(agent core.Agent) {
		dk := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := CreateVirulenceProcAura(dk.GetCharacter())

		icd := core.Cooldown{
			Timer:    dk.NewTimer(),
			Duration: time.Second * 10.0,
		}
		procAura.Icd = &icd

		core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
			Label: "Sigil of Virulence",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) || !dk.IsFuStrike(spell) {
					return
				}

				if sim.RandomFloat("SigilOfVirulence") < 0.80 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		}))
	})

	addItemEffect(50459, func(agent core.Agent) {
		character := agent.GetCharacter()
		dk := agent.(DeathKnightAgent).GetDeathKnight()

		procAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Sigil of the Hanged Man Proc",
				ActionID:  core.ActionID{SpellID: 71227},
				Duration:  time.Second * 15,
				MaxStacks: 3,
			},
			BonusPerStack: stats.Stats{stats.Strength: 73},
		})

		core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
			Label: "Sigil of the Hanged Man",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !dk.IsFuStrike(spell) {
					return
				}

				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		}))
	})

	addItemEffect(50462, func(agent core.Agent) {
		character := agent.GetCharacter()
		dk := agent.(DeathKnightAgent).GetDeathKnight()

		procAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Sigil of the Bone Gryphon Proc",
				ActionID:  core.ActionID{SpellID: 71229},
				Duration:  time.Second * 15,
				MaxStacks: 5,
			},
			BonusPerStack: stats.Stats{stats.Dodge: 44},
		})

		core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
			Label: "Sigil of the Bone Gryphon",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell != dk.RuneStrike {
					return
				}

				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		}))
	})

	CreateGladiatorsSigil(42618, "Savage", 94, 6)
	CreateGladiatorsSigil(42619, "Hateful", 106, 6)
	CreateGladiatorsSigil(42620, "Deadly", 120, 10)
	CreateGladiatorsSigil(42621, "Furious", 144, 10)
	CreateGladiatorsSigil(42622, "Relentless", 172, 10)
	CreateGladiatorsSigil(51417, "Wrathful", 204, 10)
}

func CreateGladiatorsSigil(id int32, name string, ap float64, seconds time.Duration) {
	addItemEffect(id, func(agent core.Agent) {
		dk := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := dk.NewTemporaryStatsAura(name+" Gladiator's Sigil of Strife Proc", core.ActionID{ItemID: id}, stats.Stats{stats.AttackPower: ap}, time.Second*seconds)

		core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
			Label: name + " Gladiator's Sigil of Strife",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell != dk.PlagueStrike {
					return
				}

				procAura.Activate(sim)
			},
		}))
	})
}
