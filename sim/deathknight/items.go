package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/common/wotlk"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO: T7 Tank missing Icebound Fortitude
// TODO: T8 Dps missing Heart Strike
// TODO: T8 Tank missing Rune Strike and AMS
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

var ItemSetThassariansBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Thassarian's Battlegear",
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

	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label: "Unholy Might",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !icd.IsReady(sim) || (spell != dk.BloodStrike.Spell && spell != dk.HeartStrike.Spell) {
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
	Name: "Thassarian's Plate",
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

func (dk *Deathknight) thassariansPlateCooldownReduction(spell *RuneSpell) time.Duration {
	if !dk.HasSetBonus(ItemSetThassariansPlate, 4) {
		return 0
	}

	if spell == dk.UnbreakableArmor || spell == dk.BoneShield /* || spell == dk.VampiricBlood*/ {
		return 10 * time.Second
	} /* else if spell == dk.DarkCommand {
		return 2 * time.Second
	}*/
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

func (dk *Deathknight) scourgelordsBattlegearDamageBonus(spell *RuneSpell) float64 {
	if !dk.HasSetBonus(ItemSetScourgelordsBattlegear, 2) {
		return 1.0
	}

	if spell == dk.Obliterate || spell == dk.ScourgeStrike {
		return 1.1
	} else if spell == dk.HeartStrike {
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
	})

	dk.onRuneSpendT10 = func(sim *core.Simulation) {
		if dk.CurrentBloodRunes() == 0 && dk.CurrentFrostRunes() == 0 && dk.CurrentUnholyRunes() == 0 && dk.CurrentDeathRunes() == 0 {
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
		Label:    "Blood Armor",
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
		Label: "Blood Armor Proc",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == dk.BloodTap.Spell {
				bonusAura.Activate(sim)
			}
		},
	}))
}

func (dk *Deathknight) sigilOfAwarenessBonus(spell *RuneSpell) float64 {
	if dk.Equip[proto.ItemSlot_ItemSlotRanged].ID != 40207 {
		return 0
	}

	if spell == dk.Obliterate {
		return 336
	} else if spell == dk.ScourgeStrike {
		return 189
	} else if spell == dk.DeathStrike {
		return 315
	}
	return 0
}

func (dk *Deathknight) sigilOfTheFrozenConscienceBonus() float64 {
	return core.TernaryFloat64(dk.Equip[proto.ItemSlot_ItemSlotRanged].ID == 40822, 111, 0)
}

func (dk *Deathknight) sigilOfTheWildBuckBonus() float64 {
	return core.TernaryFloat64(dk.Equip[proto.ItemSlot_ItemSlotRanged].ID == 40867, 80, 0)
}

func (dk *Deathknight) sigilOfArthriticBindingBonus() float64 {
	return core.TernaryFloat64(dk.Equip[proto.ItemSlot_ItemSlotRanged].ID == 40875, 91, 0)
}

func (dk *Deathknight) sigilOfTheVengefulHeartDeathCoil() float64 {
	return core.TernaryFloat64(dk.Equip[proto.ItemSlot_ItemSlotRanged].ID == 45254, 380, 0)
}

func (dk *Deathknight) sigilOfTheVengefulHeartFrostStrike() float64 {
	return core.TernaryFloat64(dk.Equip[proto.ItemSlot_ItemSlotRanged].ID == 45254, 113, 0)
}

func init() {
	// Rune of Cinderglacier
	core.NewItemEffect(53341, func(agent core.Agent) {
		character := agent.GetCharacter()

		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.ID == 53341
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.ID == 53341
		if !mh && !oh {
			return
		}

		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		cinderBonusCoeff := 1.2

		consumeSpells := [5]core.ActionID{
			BloodBoilActionID,
			DeathCoilActionID,
			FrostStrikeActionID,
			HowlingBlastActionID,
			IcyTouchActionID,
		}

		cinderProcAura := character.GetOrRegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 53386},
			Label:     "Cinderglacier",
			Duration:  time.Second * 30,
			MaxStacks: 2,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.SetStacks(sim, aura.MaxStacks)
				aura.Unit.PseudoStats.ShadowDamageDealtMultiplier *= cinderBonusCoeff
				aura.Unit.PseudoStats.FrostDamageDealtMultiplier *= cinderBonusCoeff
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.ShadowDamageDealtMultiplier /= cinderBonusCoeff
				aura.Unit.PseudoStats.FrostDamageDealtMultiplier /= cinderBonusCoeff
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Outcome.Matches(core.OutcomeLanded) {
					return
				}

				shouldConsume := false
				for _, consumeSpell := range consumeSpells {
					if spell.ActionID.SameAction(consumeSpell) {
						shouldConsume = true
						break
					}
				}

				if shouldConsume {
					aura.RemoveStack(sim)
				}
			},
		})

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: "Rune of Cinderglacier",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}

				if mh && !oh {
					if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeMH) {
						return
					}
				} else if oh && !mh {
					if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOH) {
						return
					}
				} else if mh && oh {
					if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
				}

				if ppmm.Proc(sim, spellEffect.ProcMask, "rune of cinderglacier") {
					cinderProcAura.Activate(sim)
				}
			},
		}))
	})

	// Sigils

	core.NewItemEffect(40715, func(agent core.Agent) {
		dk := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := dk.NewTemporaryStatsAura("Sigil of Haunted Dreams Proc", core.ActionID{ItemID: 40715}, stats.Stats{stats.MeleeCrit: 173.0, stats.SpellCrit: 173.0}, time.Second*10)

		icd := core.Cooldown{
			Timer:    dk.NewTimer(),
			Duration: time.Second * 45.0,
		}

		core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
			Label: "Sigil of Haunted Dreams",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) || spell != dk.BloodStrike.Spell {
					return
				}

				if sim.RandomFloat("SigilOfHauntedDreams") < 0.15 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		}))
	})

	core.NewItemEffect(47673, func(agent core.Agent) {
		dk := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := dk.NewTemporaryStatsAura("Sigil of Virulence Proc", core.ActionID{ItemID: 47673}, stats.Stats{stats.Strength: 200.0}, time.Second*20)

		icd := core.Cooldown{
			Timer:    dk.NewTimer(),
			Duration: time.Second * 10.0,
		}

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

	core.NewItemEffect(50459, func(agent core.Agent) {
		character := agent.GetCharacter()
		dk := agent.(DeathKnightAgent).GetDeathKnight()

		procAura := wotlk.MakeStackingAura(character, wotlk.StackingProcAura{
			Aura: core.Aura{
				Label:     "Sigil of the Hanged Man Proc",
				ActionID:  core.ActionID{ItemID: 50459},
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

	CreateGladiatorsSigil(42619, "Hateful", 106, 6)
	CreateGladiatorsSigil(42620, "Deadly", 120, 10)
	CreateGladiatorsSigil(42621, "Furious", 144, 10)
	CreateGladiatorsSigil(42622, "Relentless", 172, 10)
	CreateGladiatorsSigil(51417, "Wrathful", 204, 10)
}

func CreateGladiatorsSigil(id int32, name string, ap float64, seconds time.Duration) {
	core.NewItemEffect(id, func(agent core.Agent) {
		dk := agent.(DeathKnightAgent).GetDeathKnight()
		procAura := dk.NewTemporaryStatsAura(name+" Gladiator's Sigil of Strife Proc", core.ActionID{ItemID: id}, stats.Stats{stats.AttackPower: ap}, time.Second*seconds)

		core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
			Label: name + " Gladiator's Sigil of Strife",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell != dk.PlagueStrike.Spell {
					return
				}

				procAura.Activate(sim)
			},
		}))
	})
}
