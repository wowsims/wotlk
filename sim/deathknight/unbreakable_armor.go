package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerUnbreakableArmorSpell() {
	if !dk.Talents.UnbreakableArmor {
		return
	}

	actionID := core.ActionID{SpellID: 51271}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 1

	dk.UnbreakableArmorAura = dk.RegisterAura(core.Aura{
		Label:    "Unbreakable Armor",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.UnbreakableArmorAura.Unit.AddStatDependencyDynamic(sim, stats.Strength, stats.Strength, 1.2)
			dk.UnbreakableArmorAura.Unit.AddStatDependencyDynamic(sim, stats.Armor, stats.Armor, 1.25)
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.UnbreakableArmorAura.Unit.AddStatDependencyDynamic(sim, stats.Strength, stats.Strength, 1.0/1.2)
			dk.UnbreakableArmorAura.Unit.AddStatDependencyDynamic(sim, stats.Armor, stats.Armor, 1.0/1.25)
		},
	})

	dk.UnbreakableArmor = dk.RegisterSpell(core.SpellConfig{
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
			dkSpellCost := dk.DetermineOptimalCost(sim, 0, 1, 0)
			dk.Spend(sim, spell, dkSpellCost)
			amountOfRunicPower := 10.0
			dk.AddRunicPower(sim, amountOfRunicPower, dk.UnbreakableArmor.RunicPowerMetrics())

			dk.UnbreakableArmorAura.Activate(sim)
			dk.UnbreakableArmorAura.Prioritize()
		},
	})
}

func (dk *Deathknight) CanUnbreakableArmor(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0, 0, 1, 0) && dk.UnbreakableArmor.IsReady(sim)
}

func (dk *Deathknight) CastUnbreakableArmor(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanUnbreakableArmor(sim) {
		dk.UnbreakableArmor.Cast(sim, target)
		return true
	}
	return false
}
