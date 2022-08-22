package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (druid *Druid) registerTigersFurySpell() {
	actionID := core.ActionID{SpellID: 50213}
	energyMetrics := druid.NewEnergyMetrics(actionID)
	instantEnergy := 20.0 * float64(druid.Talents.KingOfTheJungle)

	dmgBonus := 80.0

	druid.TigersFuryAura = druid.RegisterAura(core.Aura{
		Label:    "Tiger's Fury Aura",
		ActionID: actionID,
		Duration: 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.BonusDamage += dmgBonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.BonusDamage -= dmgBonus
		},
	})

	spell := druid.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.AddEnergy(sim, instantEnergy, energyMetrics)

			druid.TigersFuryAura.Activate(sim)
		},
	})

	druid.TigersFury = spell
}
