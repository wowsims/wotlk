package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (druid *Druid) registerDemoralizingRoarSpell() {
	druid.DemoralizingRoarAuras = make([]*core.Aura, druid.Env.GetNumTargets())
	for _, target := range druid.Env.Encounter.Targets {
		druid.DemoralizingRoarAuras[target.Index] = core.DemoralizingRoarAura(&target.Unit, druid.Talents.FeralAggression)
	}

	druid.DemoralizingRoar = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48560},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,

		RageCost: core.RageCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  62 * 2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.Targets {
				result := spell.CalcAndDealOutcome(sim, &aoeTarget.Unit, spell.OutcomeMagicHit)
				if result.Landed() {
					druid.DemoralizingRoarAuras[aoeTarget.Index].Activate(sim)
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

	refreshWindow := time.Second * 2

	if (druid.MangleBear != nil) && (!druid.MangleBear.IsReady(sim)) {
		refreshWindow = druid.MangleBear.ReadyAt() - sim.CurrentTime + core.GCDDefault
	}

	return maintainOnly &&
		druid.DemoralizingRoarAuras[druid.CurrentTarget.Index].ShouldRefreshExclusiveEffects(sim, refreshWindow)
}
