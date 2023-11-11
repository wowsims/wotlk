package priest

// var ItemSetCrimsonAcolytesRaiment = core.NewItemSet(core.ItemSet{
// 	Name: "Crimson Acolyte's Raiment",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			priest := agent.(PriestAgent).GetPriest()

// 			var curAmount float64
// 			procSpell := priest.RegisterSpell(core.SpellConfig{
// 				ActionID:    core.ActionID{SpellID: 70770},
// 				SpellSchool: core.SpellSchoolHoly,
// 				ProcMask:    core.ProcMaskEmpty,
// 				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers | core.SpellFlagHelpful,

// 				DamageMultiplier: 1,
// 				ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

// 				Hot: core.DotConfig{
// 					Aura: core.Aura{
// 						Label: "CrimsonAcolyteRaiment2pc",
// 					},
// 					NumberOfTicks: 3,
// 					TickLength:    time.Second * 3,
// 					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
// 						dot.SnapshotBaseDamage = curAmount * 0.33
// 						dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
// 					},
// 					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
// 						dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
// 					},
// 				},
// 			})

// 			priest.RegisterAura(core.Aura{
// 				Label:    "Crimson Acolytes Raiment 2pc",
// 				Duration: core.NeverExpires,
// 				OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 					aura.Activate(sim)
// 				},
// 				OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 					if spell != priest.FlashHeal || sim.RandomFloat("Crimson Acolytes Raiment 2pc") >= 0.33 {
// 						return
// 					}

// 					curAmount = result.Damage
// 					hot := procSpell.Hot(result.Target)
// 					hot.Apply(sim)
// 				},
// 			})
// 		},
// 		4: func(agent core.Agent) {
// 			// Implemented in power_word_shield.go and circle_of_healing.go
// 		},
// 	},
// })
