package warlock

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) registerMetamorphosisSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsMetamorphosis) {
		return
	}

	actionID := core.ActionID{SpellID: 403789}
	warlock.MetamorphosisAura = warlock.RegisterAura(core.Aura{
		Label:    "Metamorphosis Aura",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.ApplyDynamicEquipScaling(sim, stats.Armor, 6)
			warlock.ApplyDynamicEquipScaling(sim, stats.BonusArmor, 6)
			warlock.PseudoStats.ReducedCritTakenChance += 6
			warlock.PseudoStats.ThreatMultiplier *= 1.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.RemoveDynamicEquipScaling(sim, stats.Armor, 6)
			warlock.RemoveDynamicEquipScaling(sim, stats.BonusArmor, 6)
			warlock.PseudoStats.ReducedCritTakenChance -= 6
			warlock.PseudoStats.ThreatMultiplier /= 1.5
		},
	})

	manaMetrics := warlock.NewManaMetrics(actionID)

	warlock.Metamorphosis = warlock.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ManaCost: core.ManaCostOptions{
			BaseCost: 1.0,
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if warlock.MetamorphosisAura.IsActive() {
				warlock.MetamorphosisAura.Deactivate(sim)
				warlock.AddMana(sim, warlock.BaseMana, manaMetrics)
			} else {
				warlock.MetamorphosisAura.Activate(sim)
			}
		},
	})
}
