package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (hunter *Hunter) getAspectOfTheHawkSpellConfig(rank int) core.SpellConfig {
	var impHawkAura *core.Aura
	improvedHawkProcChance := 0.01 * float64(hunter.Talents.ImprovedAspectOfTheHawk)

	spellId := [8]int32{0, 13165, 14318, 14319, 14320, 14321, 14322, 25296}[rank]
	rap := [8]float64{0, 20, 35, 50, 70, 90, 110, 120}[rank]
	//manaCost := [8]float64{0, 20, 35, 50, 70, 90, 110, 120}[rank]
	level := [8]int{0, 10, 18, 28, 38, 48, 58, 60}[rank]

	if hunter.Talents.ImprovedAspectOfTheHawk > 0 {
		improvedHawkBonus := 1.3
		impHawkAura = hunter.GetOrRegisterAura(core.Aura{
			Label:    "Quick Shots",
			ActionID: core.ActionID{SpellID: 6150},
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyRangedSpeed(sim, improvedHawkBonus)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyRangedSpeed(sim, 1/improvedHawkBonus)
			},
		})
	}

	actionID := core.ActionID{SpellID: spellId}
	hunter.AspectOfTheHawkAura = hunter.NewTemporaryStatsAuraWrapped(
		"Aspect of the Hawk",
		actionID,
		stats.Stats{
			stats.RangedAttackPower: rap,
		},
		core.NeverExpires,
		func(aura *core.Aura) {
			aura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell != hunter.AutoAttacks.RangedAuto() {
					return
				}

				if impHawkAura != nil && sim.RandomFloat("Imp Aspect of the Hawk") < improvedHawkProcChance {
					impHawkAura.Activate(sim)
				}
			}
		})

	return core.SpellConfig{
		ActionID:      actionID,
		Flags:         core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,

		// ManaCost: core.ManaCostOptions{
		// 	FlatCost: manaCost,
		// },

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if hunter.AspectOfTheHawkAura.IsActive() {
				hunter.AspectOfTheHawkAura.Deactivate(sim)
			} else {
				hunter.AspectOfTheHawkAura.Activate(sim)
			}
		},
	}
}

func (hunter *Hunter) registerAspectOfTheHawkSpell() {
	maxRank := 7

	for i := 1; i <= maxRank; i++ {
		config := hunter.getAspectOfTheHawkSpellConfig(i)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.AspectOfTheHawk = hunter.GetOrRegisterSpell(config)
		}
	}
}
