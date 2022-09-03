package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerExplosiveTrapSpell(timer *core.Timer) {
	actionID := core.ActionID{SpellID: 49067}
	baseCost := 0.19 * hunter.BaseMana

	applyAOEDamage := core.ApplyEffectFuncAOEDamageCapped(hunter.Env, core.SpellEffect{
		ProcMask: core.ProcMaskSpellDamage,

		BonusSpellHitRating: float64(hunter.Talents.SurvivalTactics) * 2 * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1 *
			(1 + 0.02*float64(hunter.Talents.TNT)),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				rap := hitEffect.RangedAttackPower(spell.Unit) + hitEffect.RangedAttackPowerOnTarget()
				return core.DamageRoll(sim, 523, 671) + rap*0.1
			},
		},
		OutcomeApplier: hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(false, false, hunter.CurrentTarget)),
	})

	hunter.ExplosiveTrap = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.2*float64(hunter.Talents.Resourcefulness)),

				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*30 - time.Second*2*time.Duration(hunter.Talents.Resourcefulness),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			applyAOEDamage(sim, target, spell)
			hunter.ExplosiveTrapDot.Apply(sim)
		},

		InitialDamageMultiplier: 1 +
			.10*float64(hunter.Talents.TrapMastery) +
			.02*float64(hunter.Talents.TNT),
	})

	periodicOutcomeFunc := hunter.OutcomeFuncRangedHit()
	if hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfExplosiveTrap) {
		periodicOutcomeFunc = hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(false, false, hunter.CurrentTarget))
	}

	hunter.ExplosiveTrapDot = core.NewDot(core.Dot{
		Spell: hunter.ExplosiveTrap,
		Aura: hunter.RegisterAura(core.Aura{
			Label:    "Explosive Trap",
			ActionID: actionID,
		}),
		NumberOfTicks: 10,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncAOESnapshot(hunter.Env, core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,

			BonusSpellHitRating: float64(hunter.Talents.SurvivalTactics) * 2 * core.SpellHitRatingPerHitChance,
			DamageMultiplier:    1,
			ThreatMultiplier:    1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					rap := hitEffect.RangedAttackPower(spell.Unit) + hitEffect.RangedAttackPowerOnTarget()
					return 90 + rap*0.1
				},
			},
			OutcomeApplier: periodicOutcomeFunc,
			IsPeriodic:     true,
		}),
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
