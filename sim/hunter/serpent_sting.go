package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (hunter *Hunter) registerSerpentStingSpell() {
	actionID := core.ActionID{SpellID: 27016}
	baseCost := 275.0

	hunter.SerpentSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(hunter.Talents.Efficiency)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskRangedSpecial,
			ThreatMultiplier: 1,
			OutcomeApplier:   hunter.OutcomeFuncRangedHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					hunter.SerpentStingDot.Apply(sim)
				}
			},
		}),
	})

	target := hunter.CurrentTarget
	hunter.SerpentStingDot = core.NewDot(core.Dot{
		Spell: hunter.SerpentSting,
		Aura: target.RegisterAura(core.Aura{
			Label:    "SerpentSting-" + strconv.Itoa(int(hunter.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 5,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 + 0.06*float64(hunter.Talents.ImprovedStings),
			ThreatMultiplier: 1,
			IsPeriodic:       true,

			BaseDamage: core.BuildBaseDamageConfig(func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
				attackPower := spellEffect.RangedAttackPower(spell.Unit) + spellEffect.RangedAttackPowerOnTarget()
				return 132 + attackPower*0.02
			}, 0),
			OutcomeApplier: hunter.OutcomeFuncTick(),
		}),
	})
}
