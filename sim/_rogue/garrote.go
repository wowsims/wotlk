package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

// TODO: Add level based scaling
func (rogue *Rogue) registerGarrote() {
	numTicks := int32(6)
	rogue.Garrote = rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48676},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.costModifier(50 - 10*float64(rogue.Talents.DirtyDeeds)),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !rogue.PseudoStats.InFrontOfTarget && rogue.IsStealthed()
		},

		DamageMultiplier: 1 + 0.04*float64(rogue.Talents.Opportunity),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Garrote",
				Tag:   RogueBleedTag,
			},
			NumberOfTicks: numTicks,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 119 + dot.Spell.MeleeAttackPower()*0.07
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit)
			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				spell.Dot(target).Apply(sim)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
