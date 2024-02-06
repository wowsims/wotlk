package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

// TODO: Add level based scaling to armor reduction
// TODO: spellIDs
func (rogue *Rogue) registerExposeArmorSpell() {
	rogue.ExposeArmorAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.ExposeArmorAura(target)
	})
	rogue.exposeArmorDurations = [6]time.Duration{
		0,
		time.Second * 6,
		time.Second * 12,
		time.Second * 18,
		time.Second * 24,
		time.Second * 30,
	}

	rogue.ExposeArmor = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 8647},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | rogue.finisherFlags() | core.SpellFlagAPL,
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:   25.0,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				debuffAura := rogue.ExposeArmorAuras.Get(target)
				debuffAura.Duration = rogue.exposeArmorDurations[rogue.ComboPoints()]
				debuffAura.Activate(sim)
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},

		RelatedAuras: []core.AuraArray{rogue.ExposeArmorAuras},
	})
}
