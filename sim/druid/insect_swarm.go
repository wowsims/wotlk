package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerInsectSwarmSpell() {
	actionID := core.ActionID{SpellID: 27013}
	baseCost := 0.08 * druid.BaseMana

	target := druid.CurrentTarget
	missAura := core.InsectSwarmAura(target)

	druid.InsectSwarm = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			OutcomeApplier:   druid.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.InsectSwarmDot.Apply(sim)
					missAura.Activate(sim)
				}
			},
		}),
	})

	druid.InsectSwarmDot = core.NewDot(core.Dot{
		Spell: druid.InsectSwarm,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Insect Swarm",
			ActionID: actionID,
		}),
		NumberOfTicks: 6 + core.TernaryInt(druid.Talents.NaturesSplendor, 1, 0),
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(124, 0.2),
			OutcomeApplier:   druid.OutcomeFuncTick(),
		}),
	})
}
