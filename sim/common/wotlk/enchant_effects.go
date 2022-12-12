package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	// Keep these in order by item ID.

	core.NewEnchantEffect(3251, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.EffectID == 3251
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.EffectID == 3251
		if !mh && !oh {
			return
		}
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(4.0, procMask)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 44622},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 237, spell.OutcomeMagicHitAndCrit)
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Giant Slayer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if result.Target.MobType != proto.MobType_MobTypeGiant {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Giant Slayer") {
					procSpell.Cast(sim, result.Target)
				}
			},
		})
	})

	core.NewEnchantEffect(3239, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.EffectID == 3239
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.EffectID == 3239
		if !mh && !oh {
			return
		}
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(4.0, procMask)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 44525},
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(185, 215), spell.OutcomeMagicHitAndCrit)
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Icebreaker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Icebreaker") {
					procSpell.Cast(sim, result.Target)
				}
			},
		})
	})

	core.NewEnchantEffect(3607, func(agent core.Agent) {
		character := agent.GetCharacter()
		// TODO: This should be ranged-only haste. For now just make it hunter-only.
		if character.Class == proto.Class_ClassHunter {
			character.AddStats(stats.Stats{stats.MeleeHaste: 40, stats.SpellHaste: 40})
		}
	})

	core.NewEnchantEffect(3608, func(agent core.Agent) {
		agent.GetCharacter().AddBonusRangedCritRating(40)
	})

	core.NewEnchantEffect(3748, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 42500}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(45, 67)
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "Titanium Shield Spike",
			ActionID: actionID,
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
					procSpell.Cast(sim, spell.Unit)
				}
			},
		})
	})

	core.NewEnchantEffect(3247, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			character.PseudoStats.MobTypeAttackPower += 140
		}
	})

	core.NewEnchantEffect(3253, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 1.02
	})

	core.NewEnchantEffect(3296, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 0.98
	})

	core.NewEnchantEffect(3789, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.EffectID == 3789
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.EffectID == 3789
		if !mh && !oh {
			return
		}
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		// Modify only gear armor, including from agility
		fivePercentOfArmor := (character.Equip.Stats()[stats.Armor] + 2.0*character.Equip.Stats()[stats.Agility]) * 0.05
		procAuraMH := character.NewTemporaryStatsAura("Berserking MH Proc", core.ActionID{SpellID: 59620, Tag: 1}, stats.Stats{stats.AttackPower: 400, stats.RangedAttackPower: 400, stats.Armor: -fivePercentOfArmor}, time.Second*15)
		procAuraOH := character.NewTemporaryStatsAura("Berserking OH Proc", core.ActionID{SpellID: 59620, Tag: 2}, stats.Stats{stats.AttackPower: 400, stats.RangedAttackPower: 400, stats.Armor: -fivePercentOfArmor}, time.Second*15)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Berserking (Enchant)",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Berserking") {
					if spell.IsMH() {
						procAuraMH.Activate(sim)
					} else {
						procAuraOH.Activate(sim)
					}
				}
			},
		})
	})

	// TODO: These are stand-in values without any real reference.
	core.NewEnchantEffect(3241, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.EffectID == 3241
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.EffectID == 3241
		if !mh && !oh {
			return
		}
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(3.0, procMask)

		healthMetrics := character.NewHealthMetrics(core.ActionID{ItemID: 44494})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Lifeward",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Lifeward") {
					character.GainHealth(sim, 300, healthMetrics)
				}
			},
		})
	})

	core.NewEnchantEffect(3790, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Black Magic Proc", core.ActionID{SpellID: 59626}, stats.Stats{stats.MeleeHaste: 250, stats.SpellHaste: 250}, time.Second*10)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 35,
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Black Magic",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskSpellDamage) && spell.ActionID.SpellID != 47465 && spell.ActionID.SpellID != 12867 {
					return
				}

				if icd.IsReady(sim) && sim.RandomFloat("Black Magic") < 0.35 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	core.AddWeaponEffect(3843, func(agent core.Agent, _ proto.ItemSlot) {
		w := &agent.GetCharacter().AutoAttacks.Ranged
		w.BaseDamageMin += 15
		w.BaseDamageMax += 15
	})

	core.NewEnchantEffect(3603, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 54757}

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 45,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 10,
				},
			},

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(1654, 2020), spell.OutcomeMagicCrit)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
			Type:     core.CooldownTypeDPS,
		})
	})

	core.NewEnchantEffect(3604, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 54758}

		procAura := character.NewTemporaryStatsAura("Hyperspeed Acceleration", actionID, stats.Stats{stats.MeleeHaste: 340, stats.SpellHaste: 340}, time.Second*12)

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolMagic,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
				// Shared CD with Offensive trinkets has been removed.
				// https://twitter.com/AggrendWoW/status/1579664462843633664
				// Change possibly temporary, but developers have confirmed it was intended.

				// SharedCD: core.Cooldown{
				// 	Timer:    character.GetOffensiveTrinketCD(),
				// 	Duration: time.Second * 12,
				// },
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				procAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
			Type:     core.CooldownTypeDPS,
		})
	})

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

	core.NewEnchantEffect(3370, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.EffectID == 3370
		oh := character.HasOHWeapon() && character.Equip[proto.ItemSlot_ItemSlotOffHand].HandType != proto.HandType_HandTypeTwoHand && character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.EffectID == 3370
		if !mh && !oh {
			return
		}

		actionID := core.ActionID{SpellID: 50401}
		if spell := character.GetSpell(actionID); spell != nil {
			// This function gets called twice when dual wielding this enchant, but we
			// handle both in one call.
			return
		}

		target := character.CurrentTarget

		vulnAura := core.RuneOfRazoriceVulnerabilityAura(target)
		mhRazoriceSpell := newRazoriceHitSpell(character, true)
		ohRazoriceSpell := newRazoriceHitSpell(character, false)
		character.GetOrRegisterAura(core.Aura{
			Label:    "Razor Frost",
			ActionID: core.ActionID{SpellID: 50401},
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if mh && !oh {
					if !spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
						return
					}
				} else if oh && !mh {
					if !spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
						return
					}
				} else if mh && oh {
					if !spell.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
				}

				vulnAura.Activate(sim)
				isMH := spell.ProcMask.Matches(core.ProcMaskMeleeMH)
				isOH := spell.ProcMask.Matches(core.ProcMaskMeleeOH)
				if isMH {
					mhRazoriceSpell.Cast(sim, target)
					vulnAura.AddStack(sim)
				} else if isOH {
					ohRazoriceSpell.Cast(sim, target)
					vulnAura.AddStack(sim)
				}
			},
		})
	})

	// TODO: Verify all of this
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
	//   However it will automatically overwrite one of them so it should be ok.
	//   A single application of the aura will handle both mh and oh procs.
	core.NewEnchantEffect(3368, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.EffectID == 3368
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.EffectID == 3368
		if !mh && !oh {
			return
		}

		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(2.0, procMask)

		rfcAura := newRuneOfTheFallenCrusaderAura(character, "Rune Of The Fallen Crusader Proc", core.ActionID{SpellID: 53344})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Rune Of The Fallen Crusader",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if mh && !oh {
					if !spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
						return
					}
				} else if oh && !mh {
					if !spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
						return
					}
				} else if mh && oh {
					if !spell.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
				}

				if ppmm.Proc(sim, spell.ProcMask, "rune of the fallen crusader") {
					rfcAura.Activate(sim)
				}
			},
		})
	})

	core.NewEnchantEffect(3883, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.EffectID == 3883
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.EffectID == 3883
		if !mh && !oh {
			return
		}

		character.AddStat(stats.Defense, 13*core.DefenseRatingPerDefense)
		character.MultiplyStat(stats.Stamina, 1.01)
	})

	core.NewEnchantEffect(3847, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.EffectID == 3847
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.EffectID == 3847
		if !mh {
			return
		}

		if oh {
			return
		}

		character.AddStat(stats.Defense, 25*core.DefenseRatingPerDefense)
		character.MultiplyStat(stats.Stamina, 1.02)
	})

	core.NewEnchantEffect(3722, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Lightweave Embroidery Proc", core.ActionID{SpellID: 55637}, stats.Stats{stats.SpellPower: 295}, time.Second*15)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 60,
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Lightweave Embroidery",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if icd.IsReady(sim) && sim.RandomFloat("Lightweave") < 0.35 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewEnchantEffect(3728, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.HasManaBar() {
			return
		}

		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 55767})
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Darkglow Embroidery",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if icd.IsReady(sim) && sim.RandomFloat("Darkglow") < 0.35 {
					icd.Use(sim)
					character.AddMana(sim, 400, manaMetrics, false)
				}
			},
		})
	})

	core.NewEnchantEffect(3730, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Swordguard Embroidery Proc", core.ActionID{SpellID: 55775}, stats.Stats{stats.AttackPower: 400, stats.RangedAttackPower: 400}, time.Second*15)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 55,
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Swordguard Embroidery",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}

				if icd.IsReady(sim) && sim.RandomFloat("Swordguard") < 0.2 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})
}
