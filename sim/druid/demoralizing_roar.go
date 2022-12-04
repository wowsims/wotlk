package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerDemoralizingRoarSpell() {
	cost := 10.0

	drAuras := make([]*core.Aura, druid.Env.GetNumTargets())
	for _, target := range druid.Env.Encounter.Targets {
		drAuras[target.Index] = core.DemoralizingRoarAura(&target.Unit, druid.Talents.FeralAggression)
	}
	druid.DemoralizingRoarAura = drAuras[druid.CurrentTarget.Index]

	druid.DemoralizingRoar = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48560},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskEmpty,
		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  62 * 2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.Targets {
				result := spell.CalcDamage(sim, &aoeTarget.Unit, 0, spell.OutcomeMagicHit)
				spell.DealDamage(sim, result)
				if result.Landed() {
					drAuras[aoeTarget.Index].Activate(sim)
				}
			}
		},
	})
}

func (druid *Druid) CanDemoralizingRoar(_ *core.Simulation) bool {
	return druid.CurrentRage() >= druid.DemoralizingRoar.DefaultCast.Cost
}

func (druid *Druid) ShouldDemoralizingRoar(sim *core.Simulation, filler bool, maintainOnly bool) bool {
	if !druid.CanDemoralizingRoar(sim) {
		return false
	}

	if filler {
		return true
	}

	return maintainOnly &&
		druid.CurrentTarget.ShouldRefreshAuraWithTagAtPriority(sim, core.APReductionAuraTag, druid.DemoralizingRoarAura.Priority, time.Second*2)
}
