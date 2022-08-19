package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) getSavageRoarMultiplier() float64 {
	return 1.3
}

func (druid *Druid) registerSavageRoarSpell() {
	actionID := core.ActionID{SpellID: 52610}
	baseCost := 25.0

	srm := druid.getSavageRoarMultiplier()

	druid.SavageRoarAura = druid.RegisterAura(core.Aura{
		Label:    "Savage Roar Aura",
		ActionID: actionID,
		Duration: 9,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.PhysicalDamageDealtMultiplier *= srm
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.PhysicalDamageDealtMultiplier /= srm
		},
	})

	srSpell := druid.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			druid.SpendComboPoints(sim, spell.ComboPointMetrics())
			druid.TigersFuryAura.Duration = time.Duration(float64(time.Second) * (9.0 + (5.0)*float64(druid.ComboPoints())))
			druid.TigersFuryAura.Activate(sim)
		},
	})

	druid.SavageRoar = srSpell
}
