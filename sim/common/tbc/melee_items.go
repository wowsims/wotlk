package tbc

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false
	core.NewSimpleStatItemEffect(28484, stats.Stats{stats.Health: 1500, stats.Strength: 150}, time.Second*15, time.Minute*30) // Bulwark of Kings
	core.NewSimpleStatItemEffect(28485, stats.Stats{stats.Health: 1500, stats.Strength: 150}, time.Second*15, time.Minute*30) // Bulwark of Ancient Kings

	// Proc effects. Keep these in order by item ID.

	core.NewItemEffect(9449, func(agent core.Agent) {
		character := agent.GetCharacter()

		// Assumes that the user will swap pummelers to have the buff for the whole fight.
		character.AddStat(stats.MeleeHaste, 500)
	})

	core.NewItemEffect(12632, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 12632},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamageMagicHitAndCrit(sim, target, 3)
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "Storm Gauntlets",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// https://wotlk.wowhead.com/spell=16615/add-lightning-dam-weap-03, proc mask = 20.
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				procSpell.Cast(sim, spellEffect.Target)
			},
		})
	})

	core.NewItemEffect(17111, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 17111},
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamageMagicHitAndCrit(sim, target, 2)
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "Blazefury Medallion",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// https://wotlk.wowhead.com/spell=7711/add-fire-dam-weap-02, proc mask = 20.
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				procSpell.Cast(sim, spellEffect.Target)
			},
		})
	})

	core.NewItemEffect(17112, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(17112)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		const procChance = 2.8 / 60.0

		procAura := character.NewTemporaryStatsAura("Empyrean Demolisher Proc", core.ActionID{ItemID: 17112}, stats.Stats{stats.MeleeHaste: 212}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Empyrean Demolisher",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if sim.RandomFloat("EmpyreanDemolisher") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(19019, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(19019)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(6.0, procMask)

		procActionID := core.ActionID{SpellID: 21992}

		singleTargetSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(1),
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 0.5,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamageMagicHitAndCrit(sim, target, 300)
			},
		})

		makeDebuffAura := func(target *core.Unit) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				Label:    "Thunderfury",
				ActionID: procActionID,
				Duration: time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, -25)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, 25)
				},
			})
		}

		numHits := core.MinInt32(5, character.Env.GetNumTargets())
		debuffAuras := make([]*core.Aura, character.Env.GetNumTargets())
		for _, target := range character.Env.Encounter.Targets {
			debuffAuras[target.Index] = makeDebuffAura(&target.Unit)
		}

		bounceSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(2),
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			ThreatMultiplier: 1,
			FlatThreatBonus:  63,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					result := spell.CalcDamage(sim, curTarget, 0, spell.OutcomeMagicHit)
					if result.Landed() {
						debuffAuras[target.Index].Activate(sim)
					}
					spell.DealDamage(sim, result)
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "Thunderfury",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if !ppmm.Proc(sim, spell.ProcMask, "Thunderfury") {
					return
				}

				singleTargetSpell.Cast(sim, spellEffect.Target)
				bounceSpell.Cast(sim, spellEffect.Target)
			},
		})
	})

	core.NewItemEffect(23541, func(agent core.Agent) {
		character := agent.GetCharacter()

		const procChance = 0.5 * 3.3 / 60.0
		procAura := character.NewTemporaryStatsAura("Khorium Champion Proc", core.ActionID{ItemID: 23541}, stats.Stats{stats.Strength: 120}, time.Second*30)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Khorium Champion",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// https://wotlk.wowhead.com/spell=16916/strength-of-the-champion, proc mask = 0. Handled in-game via script.
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}
				if sim.RandomFloat("KhoriumChampion") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(24114, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.BonusDamage += 5
	})

	core.NewItemEffect(27901, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(27901)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)

		const procChance = 1.5 * 0.8 / 60.0
		procAura := character.NewTemporaryStatsAura("Blackout Truncheon Proc", core.ActionID{ItemID: 27901}, stats.Stats{stats.MeleeHaste: 132}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Blackout Truncheon",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// https://wotlk.wowhead.com/spell=33489/blinding-speed, proc mask = 0. Handled in-game via script.
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if sim.RandomFloat("BlackoutTruncheon") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(28429, func(agent core.Agent) {
		character := agent.GetCharacter()

		const procChance = 3.6 / 60.0
		procAura := character.NewTemporaryStatsAura("Lionheart Champion Proc", core.ActionID{ItemID: 28429}, stats.Stats{stats.Strength: 100}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Lionheart Champion",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// https://wotlk.wowhead.com/spell=34513/lionheart, proc mask = 0. Handled in-game via script.
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}
				if sim.RandomFloat("LionheartChampion") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(28430, func(agent core.Agent) {
		character := agent.GetCharacter()

		const procChance = 3.6 / 60.0
		procAura := character.NewTemporaryStatsAura("Lionheart Executioner Proc", core.ActionID{ItemID: 28430}, stats.Stats{stats.Strength: 100}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Lionheart Executioner",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}
				if sim.RandomFloat("LionheartExecutioner") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(28437, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(28437)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)

		const procChance = 2.7 / 60.0
		procAura := character.NewTemporaryStatsAura("Drakefist Hammer Proc", core.ActionID{ItemID: 28437}, stats.Stats{stats.MeleeHaste: 212}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Drakefist Hammer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if sim.RandomFloat("DrakefistHammer") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(28438, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(28438)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)

		const procChance = 2.7 / 60.0
		procAura := character.NewTemporaryStatsAura("Dragonmaw Proc", core.ActionID{ItemID: 28438}, stats.Stats{stats.MeleeHaste: 212}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Dragonmaw",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if sim.RandomFloat("Dragonmaw") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(28439, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(28439)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)

		const procChance = 2.7 / 60.0
		procAura := character.NewTemporaryStatsAura("Dragonstrike Proc", core.ActionID{ItemID: 28439}, stats.Stats{stats.MeleeHaste: 212}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Dragonstrike",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if sim.RandomFloat("Dragonstrike") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(28573, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 34580}

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagIgnoreResists,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := 600.0
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			},
		})

		const procChance = 0.5 * 3.5 / 60.0
		character.GetOrRegisterAura(core.Aura{
			Label:    "Despair",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// ProcMask: 20
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}
				if sim.RandomFloat("Despair") > procChance {
					return
				}

				procSpell.Cast(sim, spellEffect.Target)
			},
		})
	})

	core.NewItemEffect(28767, func(agent core.Agent) { // The Decapitator
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 28767}

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagIgnoreResists | core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 10,
				},
			},

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(513, 567)
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
			Type:     core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(28774, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 34696},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamageMagicHitAndCrit(sim, target, sim.Roll(285, 315))
			},
		})

		const hasteBonus = 212.0
		const procChance = 3.7 / 60.0

		character.GetOrRegisterAura(core.Aura{
			Label:    "Glaive of the Pit",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}
				if sim.RandomFloat("GlaiveOfThePit") > procChance {
					return
				}

				procSpell.Cast(sim, spellEffect.Target)
			},
		})
	})

	core.NewItemEffect(29297, func(agent core.Agent) {
		character := agent.GetCharacter()

		const procChance = 0.03
		procAura := character.NewTemporaryStatsAura("Band of the Eternal Defender Proc", core.ActionID{ItemID: 29297}, stats.Stats{stats.Armor: 800}, time.Second*10)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 60,
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Band of the Eternal Defender",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || spell.SpellSchool != core.SpellSchoolPhysical {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Band of the Eternal Defender") < procChance {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(29301, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Band of the Eternal Champion Proc", core.ActionID{ItemID: 29301}, stats.Stats{stats.AttackPower: 160, stats.RangedAttackPower: 160}, time.Second*10)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMeleeOrRanged)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 60,
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Band of the Eternal Champion",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// mask 340
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if !ppmm.Proc(sim, spell.ProcMask, "Band of the Eternal Champion") {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(29348, func(agent core.Agent) {
		character := agent.GetCharacter()

		const procChance = 2.7 / 60.0
		procAura := character.NewTemporaryStatsAura("The Bladefist Proc", core.ActionID{ItemID: 29348}, stats.Stats{stats.MeleeHaste: 180}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "The Bladefist",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
					return
				}
				if sim.RandomFloat("The Bladefist") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(29962, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(29962)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)

		procAura := character.NewTemporaryStatsAura("Heartrazor Proc", core.ActionID{ItemID: 29962}, stats.Stats{stats.AttackPower: 270, stats.RangedAttackPower: 270}, time.Second*10)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Heartrazor",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}

				if !ppmm.Proc(sim, spell.ProcMask, "Heartrazor") {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(29996, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(29996)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)

		const procChance = 2.7 / 60.0
		actionID := core.ActionID{ItemID: 29996}

		var resourceMetricsRage *core.ResourceMetrics
		var resourceMetricsEnergy *core.ResourceMetrics
		if character.HasRageBar() {
			resourceMetricsRage = character.NewRageMetrics(actionID)
		}
		if character.HasEnergyBar() {
			resourceMetricsEnergy = character.NewEnergyMetrics(actionID)
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Rod of the Sun King",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}

				cpb := spell.Unit.GetCurrentPowerBar()
				if cpb == core.RageBar {
					if sim.RandomFloat("Rod of the Sun King") > procChance {
						return
					}
					spell.Unit.AddRage(sim, 5, resourceMetricsRage)
				} else if cpb == core.EnergyBar {
					if sim.RandomFloat("Rod of the Sun King") > procChance {
						return
					}
					spell.Unit.AddEnergy(sim, 10, resourceMetricsEnergy)
				}
			},
		})
	})

	core.NewItemEffect(30090, func(agent core.Agent) {
		character := agent.GetCharacter()

		const procChance = 3.7 / 60.0
		procAura := character.NewTemporaryStatsAura("World Breaker Proc", core.ActionID{ItemID: 30090}, stats.Stats{stats.MeleeCrit: 900}, time.Second*4)

		character.RegisterAura(core.Aura{
			Label:    "World Breaker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					procAura.Deactivate(sim)
					return
				}
				if sim.RandomFloat("World Breaker") > procChance {
					procAura.Deactivate(sim)
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(30311, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(30311)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)

		const bonus = 1.2
		const inverseBonus = 1 / 1.2
		const procChance = 0.5

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Warp Slicer Proc",
			ActionID: core.ActionID{ItemID: 30311},
			Duration: time.Second * 30,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyMeleeSpeed(sim, bonus)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyMeleeSpeed(sim, inverseBonus)
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Warp Slicer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if sim.RandomFloat("WarpSlicer") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(30316, func(agent core.Agent) {
		character := agent.GetCharacter()

		const bonus = 1.2
		const inverseBonus = 1 / 1.2
		const procChance = 0.5

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Devastation Proc",
			ActionID: core.ActionID{ItemID: 30316},
			Duration: time.Second * 30,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyMeleeSpeed(sim, bonus)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyMeleeSpeed(sim, inverseBonus)
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Devastation",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}
				if sim.RandomFloat("Devastation") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(31193, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(31193)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 31193},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamageMagicHitAndCrit(sim, target, sim.Roll(48, 54)+1*spell.SpellPower())
			},
		})

		const procChance = 0.02
		character.GetOrRegisterAura(core.Aura{
			Label:    "Blade of Unquenched Thirst",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if sim.RandomFloat("BladeOfUnquenchedThirst") > procChance {
					return
				}

				procSpell.Cast(sim, spellEffect.Target)
			},
		})
	})

	core.NewItemEffect(31318, func(agent core.Agent) {
		character := agent.GetCharacter()

		const procChance = 3.5 / 60.0
		procAura := character.NewTemporaryStatsAura("Singing Crystal Axe Proc", core.ActionID{ItemID: 31318}, stats.Stats{stats.MeleeHaste: 400}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Singing Crystal Axe",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}
				if sim.RandomFloat("SingingCrystalAxe") > procChance {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(31332, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(31332)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		var blinkstrikeSpell *core.Spell
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Millisecond,
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Blinkstrike",
			Duration: core.NeverExpires,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				blinkstrikeSpell = character.GetOrRegisterSpell(core.SpellConfig{
					ActionID:     core.ActionID{ItemID: 31332},
					SpellSchool:  core.SpellSchoolPhysical,
					ProcMask:     core.ProcMaskMeleeMHAuto,
					Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
					ApplyEffects: character.AutoAttacks.MHConfig.ApplyEffects,
				})
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if !ppmm.Proc(sim, spell.ProcMask, "Blinkstrike") {
					return
				}
				icd.Use(sim)

				aura.Unit.AutoAttacks.MaybeReplaceMHSwing(sim, blinkstrikeSpell).Cast(sim, spellEffect.Target)
			},
		})
	})

	core.NewItemEffect(31331, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(31331)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:     "The Night Blade Proc",
			ActionID:  core.ActionID{ItemID: 31331},
			Duration:  time.Second * 10,
			MaxStacks: 3,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				character.AddStatDynamic(sim, stats.ArmorPenetration, 62*float64(newStacks-oldStacks))
			},
		})

		const procChance = 2 * 1.8 / 60.0
		character.GetOrRegisterAura(core.Aura{
			Label:    "The Night Blade",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if sim.RandomFloat("The Night Blade") > procChance {
					return
				}

				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		})
	})

	core.NewItemEffect(32262, func(agent core.Agent) {
		character := agent.GetCharacter()
		mh, oh := character.GetWeaponHands(32262)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 40291},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamageMagicHitAndCrit(sim, target, 20)
			},
		})

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Siphon Essence",
			ActionID: core.ActionID{SpellID: 40291},
			Duration: time.Second * 6,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				procSpell.Cast(sim, spellEffect.Target)
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Syphon of the Nathrezim",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Syphon Of The Nathrezim") {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(32375, func(agent core.Agent) {
		character := agent.GetCharacter()

		const procChance = 0.02
		procAura := character.NewTemporaryStatsAura("Bulwark Of Azzinoth Proc", core.ActionID{ItemID: 32375}, stats.Stats{stats.Armor: 2000}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Bulwark Of Azzinoth",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() && spell.SpellSchool == core.SpellSchoolPhysical && sim.RandomFloat("Bulwark of Azzinoth") < procChance {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(34473, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Commendation of Kael'Thas Proc", core.ActionID{ItemID: 34473}, stats.Stats{stats.Dodge: 152}, time.Second*10)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 30,
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Commendation of Kael'Thas",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if aura.Unit.CurrentHealthPercent() >= 0.35 {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(34679, func(agent core.Agent) {
		character := agent.GetCharacter()
		const proc = 0.15

		var aldorAura *core.Aura
		var scryerSpell *core.Spell

		if character.ShattFaction == proto.ShattrathFaction_ShattrathFactionAldor {
			aldorAura = character.NewTemporaryStatsAura("Light's Strength", core.ActionID{SpellID: 45480}, stats.Stats{stats.AttackPower: 200}, time.Second*10)
		} else if character.ShattFaction == proto.ShattrathFaction_ShattrathFactionScryer {
			scryerSpell = character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 45428},
				SpellSchool: core.SpellSchoolArcane,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultMeleeCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := sim.Roll(333, 367)
					// TODO: validate this is a melee hit roll
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				},
			})
		}

		// Gives a chance when your harmful spells land to increase the damage of your spells and effects by up to 130 for 10 sec. (Proc chance: 20%, 50s cooldown)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}

		character.RegisterAura(core.Aura{
			Label:    "Shattered Sun Pendant of Acumen",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !spellEffect.Landed() {
					return
				}
				if !icd.IsReady(sim) || sim.RandomFloat("pendant of acumen") > proc { // can't activate if on CD or didn't proc
					return
				}
				icd.Use(sim)

				if character.ShattFaction == proto.ShattrathFaction_ShattrathFactionAldor {
					aldorAura.Activate(sim)
				} else if character.ShattFaction == proto.ShattrathFaction_ShattrathFactionScryer {
					scryerSpell.Cast(sim, spellEffect.Target)
				}
			},
		})
	})
	core.NewItemEffect(12590, func(agent core.Agent) {
		character := agent.GetCharacter()
		effectAura := character.NewTemporaryStatsAura("Felstriker Proc", core.ActionID{SpellID: 16551}, stats.Stats{stats.MeleeCrit: 100 * core.CritRatingPerCritChance}, time.Second*3)
		mh, oh := character.GetWeaponHands(12590)
		procMask := core.GetMeleeProcMaskForHands(mh, oh)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Felstriker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if !ppmm.Proc(sim, spell.ProcMask, "Felstriker") {
					return
				}
				effectAura.Activate(sim)
			},
		})
	})

	core.AddEffectsToTest = true
}
