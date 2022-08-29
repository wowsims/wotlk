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
	const drainChance = 0.2

	core.NewItemEffect(49623, func(agent core.Agent) {
		player := agent.GetCharacter()

		tempStrProc := player.NewTemporaryStatsAura("Chaos Bane", core.ActionID{SpellID: 73422}, stats.Stats{stats.Strength: 270}, time.Second*10)
		choasBaneSpell := player.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 71904},
			SpellSchool: core.SpellSchoolShadow,
			ApplyEffects: core.ApplyEffectFuncAOEDamageCapped(player.Env, core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty, // not sure if this can proc things.
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage:       core.BaseDamageConfigMagic(1900, 2100, 0),
				OutcomeApplier:   player.OutcomeFuncMagicHit(), // can miss, can't crit
			})})

		stackingAura := player.GetOrRegisterAura(core.Aura{
			Label:     "Soul Fragment",
			Duration:  time.Minute,
			ActionID:  core.ActionID{SpellID: 71905},
			MaxStacks: 10,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				player.AddStatDynamic(sim, stats.Strength, float64(newStacks-oldStacks)*30)
			},
		})

		core.MakePermanent(player.GetOrRegisterAura(core.Aura{
			Label: "Shadowmourne",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if stackingAura.GetStacks() == 10 {
					stackingAura.Deactivate(sim)
					tempStrProc.Activate(sim)
					choasBaneSpell.Cast(sim, player.CurrentTarget)
				}

				if sim.RandomFloat("shadowmourne") > drainChance {
					return
				}

				stackingAura.Activate(sim)
				stackingAura.AddStack(sim)
			},
		}))
	})
}
