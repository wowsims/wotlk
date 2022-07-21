package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) registerUnbreakableArmorSpell() {
	if !deathKnight.Talents.UnbreakableArmor {
		return
	}

	actionID := core.ActionID{SpellID: 51271}
	cdTimer := deathKnight.NewTimer()
	cd := time.Minute * 1

	deathKnight.UnbreakableArmorAura = deathKnight.RegisterAura(core.Aura{
		Label:    "Unbreakable Armor",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.UnbreakableArmorAura.Unit.AddStatDependencyDynamic(sim, stats.Strength, stats.Strength, 1.2)
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.UnbreakableArmorAura.Unit.AddStatDependencyDynamic(sim, stats.Strength, stats.Strength, 1/1.2)
		},
	})

	deathKnight.UnbreakableArmor = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			deathKnight.UnbreakableArmorAura.Activate(sim)
			deathKnight.UnbreakableArmorAura.Prioritize()

			amountOfRunicPower := 10.0
			deathKnight.AddRunicPower(sim, amountOfRunicPower, deathKnight.UnbreakableArmor.RunicPowerMetrics())
		},
	})
}

func (deathKnight *DeathKnight) CanUnbreakableArmor(sim *core.Simulation) bool {
	return deathKnight.UnbreakableArmor.IsReady(sim) && deathKnight.UnbreakableArmor.CD.IsReady(sim)
}

func (deathKnight *DeathKnight) CastUnbreakableArmor(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanUnbreakableArmor(sim) {
		deathKnight.UnbreakableArmor.Cast(sim, target)
		return true
	}
	return false
}
