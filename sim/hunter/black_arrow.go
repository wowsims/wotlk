package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerBlackArrowSpell(timer *core.Timer) {
	if !hunter.Talents.BlackArrow {
		return
	}

	actionID := core.ActionID{SpellID: 63672}
	baseCost := 0.06 * hunter.BaseMana

	hunter.BlackArrow = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagIgnoreResists,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.03*float64(hunter.Talents.Efficiency)) *
					(1 - 0.2*float64(hunter.Talents.Resourcefulness)),
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*30 - time.Second*2*time.Duration(hunter.Talents.Resourcefulness),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskRangedSpecial,
			ThreatMultiplier: 1,
			OutcomeApplier:   hunter.OutcomeFuncRangedHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					hunter.BlackArrowDot.Apply(sim)
				}
			},
		}),

		InitialDamageMultiplier: 1 +
			.10*float64(hunter.Talents.TrapMastery) +
			.02*float64(hunter.Talents.TNT),
	})

	target := hunter.CurrentTarget
	hunter.BlackArrowDot = core.NewDot(core.Dot{
		Spell: hunter.BlackArrow,
		Aura: target.RegisterAura(core.Aura{
			Label:    "BlackArrow-" + strconv.Itoa(int(hunter.Index)),
			ActionID: actionID,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				hunter.AttackTables[aura.Unit.UnitIndex].DamageDealtMultiplier *= 1.06
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				hunter.AttackTables[aura.Unit.UnitIndex].DamageDealtMultiplier /= 1.06
			},
		}),
		NumberOfTicks: 5,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 *
				(1.0 / 1.06), // Black Arrow is not affected by its own 1.06 aura.
			ThreatMultiplier: 1,
			IsPeriodic:       true,

			BaseDamage: core.BuildBaseDamageConfig(func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
				attackPower := spellEffect.RangedAttackPower(spell.Unit) + spellEffect.RangedAttackPowerOnTarget()
				return 553 + attackPower*0.02
			}, 0),
			OutcomeApplier: hunter.OutcomeFuncTick(),
		}),
	})
}
