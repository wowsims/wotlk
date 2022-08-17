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
		Flags:       core.SpellFlagNoOnCastComplete,
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Deep Wounds Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.ProcMask.Matches(core.ProcMaskEmpty) {
				return
			}
			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				warrior.DeepWounds.Cast(sim, nil)
				warrior.procDeepWounds(sim, spellEffect.Target)
				warrior.procBloodFrenzy(sim, spellEffect, time.Second*6)
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

func (warrior *Warrior) procDeepWounds(sim *core.Simulation, target *core.Unit) {
	deepWoundsDot := warrior.DeepWoundsDots[target.Index]

	newDeepWoundsDamage := warrior.AutoAttacks.MH.AverageDamage() * 0.16 * float64(warrior.Talents.DeepWounds)
	if deepWoundsDot.IsActive() {
		newDeepWoundsDamage += warrior.DeepWoundsTickDamage[target.Index] * float64(6-deepWoundsDot.TickCount)
	}

	newTickDamage := newDeepWoundsDamage / 6
	warrior.DeepWoundsTickDamage[target.Index] = newTickDamage

	warrior.DeepWounds.SpellMetrics[target.TableIndex].Hits++

	deepWoundsDot.TickEffects = core.TickFuncApplyEffectsToUnit(target, core.ApplyEffectFuncDirectDamage(core.SpellEffect{
		ProcMask:         core.ProcMaskPeriodicDamage,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		IsPeriodic:       true,
		BaseDamage:       core.BaseDamageConfigFlat(newTickDamage),
		OutcomeApplier:   warrior.OutcomeFuncTick(),
	}))
	deepWoundsDot.Apply(sim)
}
