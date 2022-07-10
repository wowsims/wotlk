package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerExplosiveShotSpell(timer *core.Timer) {
	if !hunter.Talents.ExplosiveShot {
		return
	}

	actionID := core.ActionID{SpellID: 60053}
	baseCost := 0.07 * hunter.BaseMana()

	baseEffect := core.SpellEffect{
		ProcMask:        core.ProcMaskRangedSpecial,
		BonusCritRating: 2 * core.CritRatingPerCritChance * float64(hunter.Talents.SurvivalInstincts),
		DamageMultiplier: 1 *
			(1 + 0.02*float64(hunter.Talents.TNT)) *
			hunter.sniperTrainingMultiplier(),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return (hitEffect.RangedAttackPower(spell.Unit)+hitEffect.RangedAttackPowerOnTarget())*0.14 + 492
			},
			TargetSpellCoefficient: 1,
		},
		OutcomeApplier: hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, false, hunter.CurrentTarget)),
	}

	initialEffect := baseEffect
	initialEffect.OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		// TODO: Is this just an always-applied dot that ticks instantly?
		// if spellEffect.Landed() {
		hunter.ExplosiveShotDot.Apply(sim)
		// }
	}

	hunter.ExplosiveShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.03*float64(hunter.Talents.Efficiency)),
				GCD:  core.GCDDefault + hunter.latency,
			},
			IgnoreHaste: true,
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				if hunter.LockAndLoadAura.IsActive() {
					cast.Cost = 0
				}
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * 6,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(initialEffect),
	})

	dotEffect := baseEffect
	dotEffect.IsPeriodic = true
	dotEffect.ProcMask = core.ProcMaskPeriodicDamage

	target := hunter.CurrentTarget
	hunter.ExplosiveShotDot = core.NewDot(core.Dot{
		Spell: hunter.ExplosiveShot,
		Aura: target.RegisterAura(core.Aura{
			Label:    "ExplosiveShot-" + strconv.Itoa(int(hunter.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 2,
		TickLength:    time.Second * 1,
		TickEffects:   core.TickFuncSnapshot(target, dotEffect),
	})
}
