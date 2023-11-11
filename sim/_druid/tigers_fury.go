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
	cdReduction := core.TernaryDuration(druid.HasSetBonus(ItemSetDreamwalkerBattlegear, 4), time.Second*3, 0)

	druid.TigersFuryAura = druid.RegisterAura(core.Aura{
		Label:    "Tiger's Fury Aura",
		ActionID: actionID,
		Duration: 6 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.BonusDamage += dmgBonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.BonusDamage -= dmgBonus
		},
	})

	spell := druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second*30 - cdReduction,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !druid.BerserkAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.AddEnergy(sim, instantEnergy, energyMetrics)

			druid.TigersFuryAura.Activate(sim)
		},
	})

	druid.TigersFury = spell
}
