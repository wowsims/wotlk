package druid

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (druid *Druid) registerDemoralizingRoarSpell() {
	cost := 10.0

	baseEffect := core.SpellEffect{
		ProcMask:         core.ProcMaskEmpty,
		ThreatMultiplier: 1,
		FlatThreatBonus:  62 * 2,
		OutcomeApplier:   druid.OutcomeFuncMagicHit(),
	}

	numHits := druid.Env.GetNumTargets()
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = druid.Env.GetTargetUnit(i)

		demoRoarAura := core.DemoralizingRoarAura(effects[i].Target, druid.Talents.FeralAggression)
		if i == 0 {
			druid.DemoralizingRoarAura = demoRoarAura
		}

		effects[i].OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				demoRoarAura.Activate(sim)
			}
		}
	}

	druid.DemoralizingRoar = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26998},
		SpellSchool: core.SpellSchoolPhysical,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			ModifyCast:  druid.ApplyClearcasting,
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}

func (druid *Druid) CanDemoralizingRoar(sim *core.Simulation) bool {
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
