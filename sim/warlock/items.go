package warlock

func init() {
	// core.NewItemEffect(32493, func(agent core.Agent) {
	// 	warlock := agent.(WarlockAgent).GetWarlock()
	// 	procAura := warlock.NewTemporaryStatsAura("Ashtongue Talisman Proc", core.ActionID{SpellID: 40478}, stats.Stats{stats.SpellPower: 220}, time.Second*5)

	// 	warlock.RegisterAura(core.Aura{
	// 		Label:    "Ashtongue Talisman",
	// 		Duration: core.NeverExpires,
	// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
	// 			aura.Activate(sim)
	// 		},
	// 		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 			if spell == warlock.Corruption && sim.Proc(0.2, "Ashtongue Talisman of Insight") {
	// 				procAura.Activate(sim)
	// 			}
	// 		},
	// 	})
	// })
}
