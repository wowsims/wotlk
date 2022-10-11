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
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskRangedSpecial,
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

		DamageMultiplierAdditive: 1 +
			.10*float64(hunter.Talents.TrapMastery) +
			.02*float64(hunter.Talents.TNT),
		DamageMultiplier: 1 *
			(1.0 / 1.06), // Black Arrow is not affected by its own 1.06 aura.
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)
			if result.Landed() {
				hunter.BlackArrowDot.Apply(sim)
			}
			spell.DealOutcome(sim, &result)
		},
	})

	target := hunter.CurrentTarget
	hunter.BlackArrowDot = core.NewDot(core.Dot{
		Spell: hunter.BlackArrow,
		Aura: target.RegisterAura(core.Aura{
			Label:    "BlackArrow-" + strconv.Itoa(int(hunter.Index)),
			ActionID: actionID,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= 1.06
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= 1.06
			},
		}),
		NumberOfTicks: 5,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			IsPeriodic: true,
			// scales slightly better (11.5%) than the tooltip implies (10%), but isn't affected by Hunter's Mark
			BaseDamage: core.BuildBaseDamageConfig(func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
				return 553 + 0.023*(spell.Unit.GetStat(stats.RangedAttackPower)+spell.Unit.PseudoStats.MobTypeAttackPower)
			}),
			OutcomeApplier: hunter.OutcomeFuncTick(),
		}),
	})
}
