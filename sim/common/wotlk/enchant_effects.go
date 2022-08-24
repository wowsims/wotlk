package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	// Keep these in order by item ID.

	core.NewItemEffect(37339, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.ID == 37339
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.ID == 37339
		if !mh && !oh {
			return
		}
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(4.0, procMask)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 44622},
			SpellSchool: core.SpellSchoolPhysical,
			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				BaseDamage:     core.BaseDamageConfigFlat(237),
				OutcomeApplier: character.OutcomeFuncMagicHitAndCrit(character.DefaultSpellCritMultiplier()),
			}),
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Giant Slayer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if spellEffect.Target.MobType != proto.MobType_MobTypeGiant {
					return
				}

				if ppmm.Proc(sim, spellEffect.ProcMask, "Giant Slayer") {
					procSpell.Cast(sim, spellEffect.Target)
				}
			},
		})
	})

	core.NewItemEffect(37344, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.ID == 37344
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.ID == 37344
		if !mh && !oh {
			return
		}
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(7.0, procMask)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 44525},
			SpellSchool: core.SpellSchoolFire,
			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				BaseDamage:     core.BaseDamageConfigRoll(185, 215),
				OutcomeApplier: character.OutcomeFuncMagicHitAndCrit(character.DefaultSpellCritMultiplier()),
			}),
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Icebreaker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if ppmm.Proc(sim, spellEffect.ProcMask, "Icebreaker") {
					procSpell.Cast(sim, spellEffect.Target)
				}
			},
		})
	})

	core.NewItemEffect(41146, func(agent core.Agent) {
		character := agent.GetCharacter()
		// TODO: This should be ranged-only haste.
		character.AddStats(stats.Stats{stats.MeleeHaste: 40, stats.SpellHaste: 40})
	})

	core.NewItemEffect(41167, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.BonusRangedCritRating += 40
	})

	core.NewItemEffect(42500, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 42500}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			Flags:       core.SpellFlagBinary,

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				BaseDamage:     core.BaseDamageConfigRoll(45, 67),
				OutcomeApplier: character.OutcomeFuncMeleeSpecialHitAndCrit(character.DefaultMeleeCritMultiplier()),
			}),
		})

		character.RegisterAura(core.Aura{
			Label:    "Titanium Shield Spike",
			ActionID: actionID,
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() && spell.SpellSchool == core.SpellSchoolPhysical {
					procSpell.Cast(sim, spell.Unit)
				}
			},
		})
	})

	core.NewItemEffect(44473, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			character.PseudoStats.MobTypeAttackPower += 140
		}
	})

	core.NewItemEffect(44485, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 1.02
	})
	core.NewItemEffect(44488, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 0.98
	})

	core.NewItemEffect(44492, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.ID == 44492
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.ID == 44492
		if !mh && !oh {
			return
		}
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		// TODO: Also reduces armor by 5%
		procAuraMH := character.NewTemporaryStatsAura("Berserking MH Proc", core.ActionID{SpellID: 59620, Tag: 1}, stats.Stats{stats.AttackPower: 400, stats.RangedAttackPower: 400}, time.Second*15)
		procAuraOH := character.NewTemporaryStatsAura("Berserking OH Proc", core.ActionID{SpellID: 59620, Tag: 2}, stats.Stats{stats.AttackPower: 400, stats.RangedAttackPower: 400}, time.Second*15)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Berserking (Enchant)",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if ppmm.Proc(sim, spellEffect.ProcMask, "Berserking") {
					if spellEffect.IsMH() {
						procAuraMH.Activate(sim)
					} else {
						procAuraOH.Activate(sim)
					}
				}
			},
		})
	})

	// TODO: These are stand-in values without any real reference.
	core.NewItemEffect(44494, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.ID == 44494
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.ID == 44494
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if ppmm.Proc(sim, spellEffect.ProcMask, "Lifeward") {
					character.GainHealth(sim, 300, healthMetrics)
				}
			},
		})
	})

	core.NewItemEffect(44495, func(agent core.Agent) {
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskSpellDamage) {
					return
				}

				if icd.IsReady(sim) && sim.RandomFloat("Black Magic") < 0.35 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	core.AddWeaponEffect(44739, func(agent core.Agent, _ proto.ItemSlot) {
		w := &agent.GetCharacter().AutoAttacks.Ranged
		w.BaseDamageMin += 15
		w.BaseDamageMax += 15
	})

	core.NewItemEffect(54998, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 54757}

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFire,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 45,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 10,
				},
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				BaseDamage:     core.BaseDamageConfigRoll(1654, 2020),
				OutcomeApplier: character.OutcomeFuncMagicCrit(character.DefaultSpellCritMultiplier()),
			}),
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
			Type:     core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(54999, func(agent core.Agent) {
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
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 12,
				},
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
		baseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 0, 1.0, 1.0, true)
		if !isMH {
			baseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1.0, 1.0, true)
		}

		return character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 50401},
			SpellSchool: core.SpellSchoolFrost,

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask: core.ProcMaskEmpty,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				BaseDamage: core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
						return baseDamage(sim, spellEffect, spell) * 0.02
					},
				},
				OutcomeApplier: character.OutcomeFuncAlwaysHit(),
			}),
		})
	}

	core.NewItemEffect(53343, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.ID == 53343
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.ID == 53343
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

				vulnAura.Activate(sim)
				isMH := spellEffect.ProcMask.Matches(core.ProcMaskMeleeMH)
				if isMH {
					mhRazoriceSpell.Cast(sim, target)
					vulnAura.AddStack(sim)
				}

				isOH := spellEffect.ProcMask.Matches(core.ProcMaskMeleeOH)
				if isOH {
					ohRazoriceSpell.Cast(sim, target)
					vulnAura.AddStack(sim)
				}
			},
		})
	})

	// TODO: Verify all of this
	newRuneOfTheFallenCrusaderAura := func(character *core.Character, auraLabel string, actionID core.ActionID) *core.Aura {
		return character.NewTemporaryStatsAuraWrapped(auraLabel, actionID, stats.Stats{}, time.Second*15, func(aura *core.Aura) {
			oldOnGain := aura.OnGain
			oldOnExpire := aura.OnExpire

			statDep := character.NewDynamicMultiplyStat(stats.Strength, 1.15)

			aura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
				oldOnGain(aura, sim)
				aura.Unit.EnableDynamicStatDep(sim, statDep)
			}

			aura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
				oldOnExpire(aura, sim)
				aura.Unit.DisableDynamicStatDep(sim, statDep)
			}
		})
	}

	// ApplyRuneOfTheFallenCrusader will be applied twice if there is two weapons with this enchant.
	//   However it will automatically overwrite one of them so it should be ok.
	//   A single application of the aura will handle both mh and oh procs.
	core.NewItemEffect(53344, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.ID == 53344
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.ID == 53344
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

				if ppmm.Proc(sim, spellEffect.ProcMask, "rune of the fallen crusader") {
					rfcAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(70164, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.ID == 70164
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.ID == 70164
		if !mh && !oh {
			return
		}

		character.AddStat(stats.Defense, 13)
		character.MultiplyStat(stats.Stamina, 1.01)
	})

	core.NewItemEffect(62158, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.ID == 62158
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.ID == 62158
		if !mh {
			return
		}

		if oh {
			return
		}

		character.AddStat(stats.Defense, 25)
		character.MultiplyStat(stats.Stamina, 1.02)
	})

	core.NewItemEffect(55642, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Lightweave Embroidery Proc", core.ActionID{SpellID: 55637}, stats.Stats{stats.SpellPower: 295, stats.HealingPower: 295}, time.Second*15)
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

	core.NewItemEffect(55768, func(agent core.Agent) {
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

	core.NewItemEffect(55777, func(agent core.Agent) {
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
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
