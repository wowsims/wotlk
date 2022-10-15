package warrior

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

var DeepWoundsActionID = core.ActionID{SpellID: 12867}

func (warrior *Warrior) applyDeepWounds() {
	if warrior.Talents.DeepWounds == 0 {
		return
	}

	warrior.DeepWounds = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    DeepWoundsActionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1 * (1 + 0.16*float64(warrior.Talents.DeepWounds)),
		ThreatMultiplier: 1,
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Deep Wounds Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell.ProcMask.Matches(core.ProcMaskEmpty) || spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				warrior.DeepWounds.Cast(sim, nil)
				warrior.procDeepWounds(sim, spellEffect.Target, spell.IsMH())
				warrior.procBloodFrenzy(sim, spellEffect, time.Second*6)
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warrior.DeepWounds {
				warrior.DeepwoundsDamageBuffer[spellEffect.Target.Index] -= warrior.DeepWoundsTickDamage[spellEffect.Target.Index]
			}
		},
	})
}

func (warrior *Warrior) newDeepWoundsDot(target *core.Unit) *core.Dot {
	return core.NewDot(core.Dot{
		Spell: warrior.DeepWounds,
		Aura: target.RegisterAura(core.Aura{
			Label:    "DeepWounds-" + strconv.Itoa(int(warrior.Index)),
			ActionID: DeepWoundsActionID,
		}),
		NumberOfTicks: 6,
		TickLength:    time.Second * 1,
	})
}

func (warrior *Warrior) procDeepWounds(sim *core.Simulation, target *core.Unit, isMh bool) {
	deepWoundsDot := warrior.DeepWoundsDots[target.Index]

	if isMh {
		warrior.DeepwoundsDamageBuffer[target.Index] += warrior.AutoAttacks.MH.AverageDamage() * warrior.PseudoStats.PhysicalDamageDealtMultiplier
	} else {
		warrior.DeepwoundsDamageBuffer[target.Index] += warrior.AutoAttacks.OH.AverageDamage() * warrior.PseudoStats.PhysicalDamageDealtMultiplier
	}

	newTickDamage := warrior.DeepwoundsDamageBuffer[target.Index] / 6
	warrior.DeepWoundsTickDamage[target.Index] = newTickDamage
	warrior.DeepWounds.SpellMetrics[target.UnitIndex].Hits++

	deepWoundsDot.TickEffects = core.TickFuncApplyEffectsToUnit(target, core.ApplyEffectFuncDirectDamage(core.SpellEffect{
		IsPeriodic:     true,
		BaseDamage:     core.BaseDamageConfigFlat(newTickDamage),
		OutcomeApplier: warrior.OutcomeFuncTick(),
	}))
	deepWoundsDot.Apply(sim)
}
