package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (druid *Druid) registerFrenziedRegenerationCD() {
	actionID := core.ActionID{SpellID: 22842}
	healthMetrics := druid.NewHealthMetrics(actionID)
	rageMetrics := druid.NewRageMetrics(actionID)

	cdTimer := druid.NewTimer()
	cd := time.Minute * 3
	healingMulti := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFrenziedRegeneration), 1.2, 1.0)

	druid.FrenziedRegenerationAura = druid.RegisterAura(core.Aura{
		Label:    "Frenzied Regeneration",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.HealingTakenMultiplier *= healingMulti
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.HealingTakenMultiplier /= healingMulti
		},
	})

	druid.FrenziedRegeneration = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID: actionID,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 10,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					rageDumped := min(druid.CurrentRage(), 10.0)
					healthGained := rageDumped * 0.3 / 100 * druid.MaxHealth() * druid.PseudoStats.HealingTakenMultiplier

					if druid.FrenziedRegenerationAura.IsActive() {
						druid.SpendRage(sim, rageDumped, rageMetrics)
						druid.GainHealth(sim, healthGained, healthMetrics)
					}
				},
			})

			druid.FrenziedRegenerationAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.FrenziedRegeneration.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}
