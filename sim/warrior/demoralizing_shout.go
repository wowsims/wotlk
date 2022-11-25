package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerDemoralizingShoutSpell() {
	cost := 10.0 - float64(warrior.Talents.FocusedRage)

	dsAuras := make([]*core.Aura, warrior.Env.GetNumTargets())
	for _, target := range warrior.Env.Encounter.Targets {
		dsAuras[target.Index] = core.DemoralizingShoutAura(&target.Unit, warrior.Talents.BoomingVoice, warrior.Talents.ImprovedDemoralizingShout)
	}
	warrior.DemoralizingShoutAura = dsAuras[warrior.CurrentTarget.Index]

	warrior.DemoralizingShout = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 25203},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,

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
		FlatThreatBonus:  63.2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.Targets {
				result := spell.CalcAndDealOutcome(sim, &aoeTarget.Unit, spell.OutcomeMagicHit)
				if result.Landed() {
					dsAuras[aoeTarget.Index].Activate(sim)
				}
			}
		},
	})
}

func (warrior *Warrior) CanDemoralizingShout(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.DemoralizingShout.DefaultCast.Cost
}

func (warrior *Warrior) ShouldDemoralizingShout(sim *core.Simulation, filler bool, maintainOnly bool) bool {
	if !warrior.CanDemoralizingShout(sim) {
		return false
	}

	if filler {
		return true
	}

	return maintainOnly &&
		warrior.CurrentTarget.ShouldRefreshAuraWithTagAtPriority(sim, core.APReductionAuraTag, warrior.DemoralizingShoutAura.Priority, time.Second*2)
}
