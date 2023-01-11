package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (hunter *Hunter) registerExplosiveTrapSpell(timer *core.Timer) {
	actionID := core.ActionID{SpellID: 49067}
	hasGlyph := hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfExplosiveTrap)

	hunter.ExplosiveTrap = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		Cost: core.NewManaCost(core.ManaCostOptions{
			BaseCost:   0.19,
			Multiplier: 1 - 0.2*float64(hunter.Talents.Resourcefulness),
		}),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*30 - time.Second*2*time.Duration(hunter.Talents.Resourcefulness),
			},
		},

		DamageMultiplierAdditive: 1 +
			.02*float64(hunter.Talents.TNT),
		CritMultiplier:   hunter.critMultiplier(false, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.Targets {
				baseDamage := sim.Roll(523, 671) + 0.1*spell.RangedAttackPower(&aoeTarget.Unit)
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, spell.OutcomeRangedHitAndCrit)
			}
			hunter.ExplosiveTrapDot.Apply(sim)
		},
	})

	hunter.ExplosiveTrapDot = core.NewDot(core.Dot{
		Spell: hunter.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplierAdditive: 1 +
				.10*float64(hunter.Talents.TrapMastery) +
				.02*float64(hunter.Talents.TNT),
			CritMultiplier:   hunter.critMultiplier(false, false),
			ThreatMultiplier: 1,
		}),
		Aura: hunter.RegisterAura(core.Aura{
			Label:    "Explosive Trap",
			ActionID: actionID,
		}),
		NumberOfTicks: 10,
		TickLength:    time.Second * 2,

		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			baseDamage := 90 + 0.1*dot.Spell.RangedAttackPower(target)
			for _, aoeTarget := range sim.Encounter.Targets {
				if hasGlyph {
					dot.Spell.CalcAndDealPeriodicDamage(sim, &aoeTarget.Unit, baseDamage, dot.Spell.OutcomeRangedHitAndCrit)
				} else {
					dot.Spell.CalcAndDealPeriodicDamage(sim, &aoeTarget.Unit, baseDamage, dot.Spell.OutcomeRangedHit)
				}
			}
		},
	})

	timeToTrapWeave := time.Millisecond * time.Duration(hunter.Rotation.TimeToTrapWeaveMs)
	halfWeaveTime := timeToTrapWeave / 2
	hunter.TrapWeaveSpell = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID.WithTag(1),
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagNoMetrics | core.SpellFlagNoLogs,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Assume we started running after the most recent ranged auto, so that time
			// can be subtracted from the run in.
			reachLocationAt := hunter.mayMoveAt + halfWeaveTime
			layTrapAt := core.MaxDuration(reachLocationAt, sim.CurrentTime)
			doneAt := layTrapAt + halfWeaveTime

			hunter.AutoAttacks.DelayRangedUntil(sim, doneAt+time.Millisecond*500)

			if layTrapAt == sim.CurrentTime {
				success := hunter.ExplosiveTrap.Cast(sim, hunter.CurrentTarget)
				if doneAt > hunter.GCD.ReadyAt() {
					hunter.GCD.Set(doneAt)
				} else if !success {
					hunter.WaitForMana(sim, hunter.ExplosiveTrap.CurCast.Cost)
				}
			} else {
				// Make sure the GCD doesn't get used while we're waiting.
				hunter.WaitUntil(sim, doneAt)

				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: layTrapAt,
					OnAction: func(sim *core.Simulation) {
						hunter.GCD.Reset()
						success := hunter.ExplosiveTrap.Cast(sim, hunter.CurrentTarget)
						if doneAt > hunter.GCD.ReadyAt() {
							hunter.GCD.Set(doneAt)
						} else if !success {
							hunter.WaitForMana(sim, hunter.ExplosiveTrap.CurCast.Cost)
						}
					},
				})
			}
		},
	})
}
