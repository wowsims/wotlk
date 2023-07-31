package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// TODO (maybe) https://github.com/magey/wotlk-warrior/issues/23 - Rend is not benefitting from Two-Handed Weapon Specialization
func (warrior *Warrior) RegisterRendSpell(rageThreshold float64, healthThreshold float64) {
	dotDuration := time.Second * 15
	dotTicks := int32(5)
	if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRending) {
		dotDuration += time.Second * 6
		dotTicks += 2
	}

	warrior.Rend = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47465},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagNoOnCastComplete,

		RageCost: core.RageCostOptions{
			Cost:   10 - float64(warrior.Talents.FocusedRage),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 + 0.1*float64(warrior.Talents.ImprovedRend),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rend",
				Tag:   "Rend",
			},
			NumberOfTicks: dotTicks,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = (380 + warrior.AutoAttacks.MH.CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower())) / 5
				// 135% damage multiplier is applied at the beginning of the fight and removed when target is at 75% health
				if sim.GetRemainingDurationPercent() > 0.75 {
					dot.SnapshotBaseDamage *= 1.35
				}
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				warrior.procBloodFrenzy(sim, result, dotDuration)
				warrior.RendValidUntil = sim.CurrentTime + dotDuration
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealOutcome(sim, result)

			// Queue Overpower to be cast at 3s after rend if talented for 3/3 TfB
			if warrior.Talents.TasteForBlood == 3 {
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + time.Second*3,
					OnAction: func(_ *core.Simulation) {
						// Force to use OP on first 3s due to AM ticks that happens before TfB procs and break that
						if warrior.Overpower.CanCast(sim, target) && (warrior.ShouldOverpower(sim) || sim.CurrentTime >= time.Second*3 && sim.CurrentTime <= time.Second*4) {
							warrior.CastFullTfbOverpower(sim, target)
						}
					},
				})
			}
		},
	})

	warrior.RendHealthThresholdAbove = healthThreshold / 100
	warrior.RendRageThresholdBelow = core.MaxFloat(warrior.Rend.DefaultCast.Cost, rageThreshold)
}

func (warrior *Warrior) ShouldRend(sim *core.Simulation) bool {
	return warrior.Rend.IsReady(sim) && sim.CurrentTime >= (warrior.RendValidUntil-warrior.RendCdThreshold) &&
		warrior.CurrentRage() >= warrior.Rend.DefaultCast.Cost && warrior.RendHealthThresholdAbove < sim.GetRemainingDurationPercent()
}
