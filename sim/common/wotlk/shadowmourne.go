package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// https://www.wowhead.com/wotlk/item=49623/shadowmourne

// Your melee attacks have a chance to drain a Soul Fragment granting you 30 Strength.
// When you have acquired 10 Soul Fragments you will unleash Chaos Bane,
//   dealing 1900 to 2100 Shadow damage split between all enemies within 15 yards and
//   granting you 270 Strength for 10 sec.

func init() {
	// https://web.archive.org/web/20120509024819/http://elitistjerks.com/f81/t37680-depth_fury_dps_discussion/p129/
	// has some testing, and arrives at ~12 ppm (~75% for 3.7 speed)
	// https://web.archive.org/web/20100508065259/http://elitistjerks.com/f81/t37462-warrior_dps_calculation_spreadsheet/p109/
	// arrives at ~80% with "2000 white swings" on a dummy.
	core.NewItemEffect(49623, func(agent core.Agent) {
		player := agent.GetCharacter()

		ppmm := player.AutoAttacks.NewPPMManager(12, core.ProcMaskMeleeOrProc)

		chaosBaneAura := player.NewTemporaryStatsAura("Chaos Bane", core.ActionID{SpellID: 73422}, stats.Stats{stats.Strength: 270}, time.Second*10)

		choasBaneSpell := player.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 71904},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty, // not sure if this can proc things.

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(1900, 2100) / float64(sim.GetNumTargets())
				for _, target := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit) // probably has a very low crit rate
				}
			},
		})

		stackingAura := player.GetOrRegisterAura(core.Aura{
			Label:     "Soul Fragment",
			Duration:  time.Minute,
			ActionID:  core.ActionID{SpellID: 71905},
			MaxStacks: 10,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				player.AddStatDynamic(sim, stats.Strength, float64(newStacks-oldStacks)*30)

				if newStacks == aura.MaxStacks {
					choasBaneSpell.Cast(sim, nil)
					chaosBaneAura.Activate(sim)
					aura.SetStacks(sim, 0)
					return
				}
			},
		})

		core.MakePermanent(player.GetOrRegisterAura(core.Aura{
			Label: "Shadowmourne",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if chaosBaneAura.IsActive() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Shadowmourne") {
					stackingAura.Activate(sim)
					stackingAura.AddStack(sim)
				}
			},
		}))
	})
}
