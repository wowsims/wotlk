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
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamageMagicHitAndCrit(sim, target, 237)
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Giant Slayer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if spellEffect.Target.MobType != proto.MobType_MobTypeGiant {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Giant Slayer") {
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
		ppmm := character.AutoAttacks.NewPPMManager(4.0, procMask)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 44525},
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamageMagicHitAndCrit(sim, target, sim.Roll(185, 215))
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Icebreaker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Icebreaker") {
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
		agent.GetCharacter().AddBonusRangedCritRating(40)
	})

	core.NewItemEffect(42500, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 42500}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				BaseDamage:     core.BaseDamageConfigRoll(45, 67),
				OutcomeApplier: character.OutcomeFuncMeleeSpecialHitAndCrit(),
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
				if spellEffect.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
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

	// Apply for Wisdom to Cloak (itemid) but NOT for Enchant Gloves Precision (spellid)
	core.NewItemEffect(44488, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.Equip[proto.ItemSlot_ItemSlotBack].Enchant.ID == 44488 {
			if character.Equip[proto.ItemSlot_ItemSlotHands].Enchant.ID == 44488 {
				// If someone has both of these enchants for some reason, this will get called twice.
				character.PseudoStats.ThreatMultiplier *= 0.98995
			} else {
				character.PseudoStats.ThreatMultiplier *= 0.98
			}
		}
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
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
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
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Lifeward") {
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
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
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
				spell.CalcAndDealDamageMagicCrit(sim, target, sim.Roll(1654, 2020))
			},
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
				spell.CalcAndDealDamageAlwaysHit(sim, target, dmg)
			},
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
				if isMH {
					mhRazoriceSpell.Cast(sim, target)
					vulnAura.AddStack(sim)
				}

				isOH := spell.ProcMask.Matches(core.ProcMaskMeleeOH)
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

	core.NewItemEffect(70164, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.ID == 70164
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.ID == 70164
		if !mh && !oh {
			return
		}

		character.AddStat(stats.Defense, 13*core.DefenseRatingPerDefense)
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

		character.AddStat(stats.Defense, 25*core.DefenseRatingPerDefense)
		character.MultiplyStat(stats.Stamina, 1.02)
	})

	core.NewItemEffect(55642, func(agent core.Agent) {
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
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
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
