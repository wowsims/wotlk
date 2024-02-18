package tbc

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false
	// Keep these in order by item ID.

	// TODO: Crusader, Mongoose, and Executioner could also be modelled as AddWeaponEffect instead
	core.AddWeaponEffect(1897, func(agent core.Agent, slot proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.MH()
		if slot == proto.ItemSlot_ItemSlotOffHand {
			w = agent.GetCharacter().AutoAttacks.OH()
		}
		w.BaseDamageMin += 5
		w.BaseDamageMax += 5
	})

	core.NewEnchantEffect(2523, func(agent core.Agent) {
		agent.GetCharacter().AddBonusRangedHitRating(30)
	})
	core.NewEnchantEffect(2724, func(agent core.Agent) {
		agent.GetCharacter().AddBonusRangedCritRating(28)
	})

	core.NewEnchantEffect(2671, func(agent core.Agent) {
		// Sunfire
		agent.GetCharacter().OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolArcane | core.SpellSchoolFire) {
				spell.BonusSpellPower += 50
			}
		})
	})
	core.NewEnchantEffect(2672, func(agent core.Agent) {
		// Soulfrost
		agent.GetCharacter().OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolFrost | core.SpellSchoolShadow) {
				spell.BonusSpellPower += 54
			}
		})
	})

	core.AddWeaponEffect(963, func(agent core.Agent, slot proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.MH()
		if slot == proto.ItemSlot_ItemSlotOffHand {
			w = agent.GetCharacter().AutoAttacks.OH()
		}
		w.BaseDamageMin += 7
		w.BaseDamageMax += 7
	})

	// ApplyCrusaderEffect will be applied twice if there is two weapons with this enchant.
	//   However, it will automatically overwrite one of them, so it should be ok.
	//   A single application of the aura will handle both mh and oh procs.
	core.NewEnchantEffect(1900, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(1900)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		// -4 str per level over 60
		const strBonus = 100.0 - 4.0*float64(core.CharacterLevel-60)
		mhAura := character.NewTemporaryStatsAura("Crusader Enchant MH", core.ActionID{SpellID: 20007, Tag: 1}, stats.Stats{stats.Strength: strBonus}, time.Second*15)
		ohAura := character.NewTemporaryStatsAura("Crusader Enchant OH", core.ActionID{SpellID: 20007, Tag: 2}, stats.Stats{stats.Strength: strBonus}, time.Second*15)

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Crusader Enchant",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Crusader") {
					if spell.IsMH() {
						mhAura.Activate(sim)
					} else {
						ohAura.Activate(sim)
					}
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(1900, 1.0, &ppmm, aura)
	})

	core.NewEnchantEffect(2929, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.BonusDamage += 2
	})

	// ApplyMongooseEffect will be applied twice if there is two weapons with this enchant.
	//   However, it will automatically overwrite one of them, so it should be ok.
	//   A single application of the aura will handle both mh and oh procs.
	core.NewEnchantEffect(2673, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(2673)
		ppmm := character.AutoAttacks.NewPPMManager(0.73, procMask)

		mhAura := character.NewTemporaryStatsAura("Lightning Speed MH", core.ActionID{SpellID: 28093, Tag: 1}, stats.Stats{stats.MeleeHaste: 30.0, stats.Agility: 120}, time.Second*15)
		ohAura := character.NewTemporaryStatsAura("Lightning Speed OH", core.ActionID{SpellID: 28093, Tag: 2}, stats.Stats{stats.MeleeHaste: 30.0, stats.Agility: 120}, time.Second*15)

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Mongoose Enchant",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "mongoose") {
					if spell.IsMH() {
						mhAura.Activate(sim)
					} else {
						ohAura.Activate(sim)
					}
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(2673, 0.73, &ppmm, aura)
	})

	core.AddWeaponEffect(2723, func(agent core.Agent, _ proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 12
		w.BaseDamageMax += 12
	})

	core.NewEnchantEffect(2621, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 0.98
	})
	core.NewEnchantEffect(2613, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 1.02
	})

	core.NewEnchantEffect(3225, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(3225)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		procAura := character.NewTemporaryStatsAura("Executioner Proc", core.ActionID{SpellID: 42976}, stats.Stats{stats.ArmorPenetration: 120}, time.Second*15)

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Executioner",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Executioner") {
					procAura.Activate(sim)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3225, 1.0, &ppmm, aura)
	})

	// https://web.archive.org/web/20100702102132/http://elitistjerks.com/f15/t27347-deathfrost_its_mechanics/p2/#post789470
	applyDeathfrostForWeapon := func(character *core.Character, procSpell *core.Spell, isMH bool) {
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 25,
		}

		label := "Deathfrost-"
		if isMH {
			label += "MH"
		} else {
			label += "OH"
		}
		ppmm := character.AutoAttacks.NewPPMManager(2.15, core.ProcMaskMelee)

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    label,
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Damage == 0 {
					return
				}

				if spell.ProcMask.Matches(core.ProcMaskMelee) {
					if ppmm.Proc(sim, spell.ProcMask, "Deathfrost") {
						procSpell.Cast(sim, result.Target)
					}
				} else if spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
					if icd.IsReady(sim) && sim.RandomFloat("Deathfrost") < 0.5 {
						icd.Use(sim)
						procSpell.Cast(sim, result.Target)
					}
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3273, 2.15, &ppmm, aura)
	}
	core.NewEnchantEffect(3273, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(3273)
		if procMask == core.ProcMaskUnknown {
			return
		}

		actionID := core.ActionID{SpellID: 46579}
		if spell := character.GetSpell(actionID); spell != nil {
			// This function gets called twice when dual wielding this enchant, but we
			// handle both in one call.
			return
		}

		debuffs := make([]*core.Aura, len(character.Env.Encounter.TargetUnits))
		for i, target := range character.Env.Encounter.TargetUnits {
			aura := target.GetOrRegisterAura(core.Aura{
				Label:    "Deathfrost",
				ActionID: actionID,
				Duration: time.Second * 8,
			})
			core.AtkSpeedReductionEffect(aura, 1.15)
			debuffs[i] = aura
		}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFrost,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcDamage(sim, target, 150, spell.OutcomeMagicCrit)
				if result.Landed() {
					debuffs[target.Index].Activate(sim)
				}
				spell.DealDamage(sim, result)
			},
		})

		if procMask.Matches(core.ProcMaskMeleeMH) {
			applyDeathfrostForWeapon(character, procSpell, true)
		}
		if procMask.Matches(core.ProcMaskMeleeOH) {
			applyDeathfrostForWeapon(character, procSpell, false)
		}
	})

	core.AddEffectsToTest = true
}
