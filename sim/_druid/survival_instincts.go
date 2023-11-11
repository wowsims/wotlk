package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerSurvivalInstinctsCD() {
	if !druid.Talents.SurvivalInstincts {
		return
	}

	actionID := core.ActionID{SpellID: 61336}
	healthMetrics := druid.NewHealthMetrics(actionID)

	cdTimer := druid.NewTimer()
	cd := time.Minute * 3
	healthFac := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfSurvivalInstincts), 0.45, 0.3)

	var bonusHealth float64
	druid.SurvivalInstinctsAura = druid.RegisterAura(core.Aura{
		Label:    "Survival Instincts",
		ActionID: actionID,
		Duration: time.Second*20 + core.TernaryDuration(druid.HasSetBonus(ItemSetNightsongBattlegear, 4), 8*time.Second, 0),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = druid.MaxHealth() * healthFac
			druid.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})
			druid.GainHealth(sim, bonusHealth, healthMetrics)
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.AddStatsDynamic(sim, stats.Stats{stats.Health: -bonusHealth})
		},
	})

	druid.SurvivalInstincts = druid.RegisterSpell(Cat|Bear, core.SpellConfig{
		ActionID: actionID,
		Flags:    SpellFlagOmenTrigger,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			druid.SurvivalInstinctsAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.SurvivalInstincts.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}
