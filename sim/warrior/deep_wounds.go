package warrior

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (warrior *Warrior) applyDeepWounds() {
	if warrior.Talents.DeepWounds == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 12867}

	deepWoundsSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagNoOnCastComplete,
	})

	var dwDots []*core.Dot

	warrior.RegisterAura(core.Aura{
		Label:    "Deep Wounds",
		Duration: core.NeverExpires,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if len(dwDots) > 0 {
				return
			}

			tickDamage := warrior.AutoAttacks.MH.AverageDamage()
			for i := int32(0); i < sim.GetNumTargets(); i++ {
				target := sim.GetTarget(i)
				dotAura := target.GetOrRegisterAura(core.Aura{
					Label:    "DeepWounds-" + strconv.Itoa(int(warrior.Index)),
					ActionID: actionID,
					Duration: time.Second * 12,
				})
				dot := core.NewDot(core.Dot{
					Spell:         deepWoundsSpell,
					Aura:          dotAura,
					NumberOfTicks: 4,
					TickLength:    time.Second * 3,
					TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
						ProcMask:         core.ProcMaskPeriodicDamage,
						DamageMultiplier: 0.2 * float64(warrior.Talents.DeepWounds),
						ThreatMultiplier: 1,
						IsPeriodic:       true,
						BaseDamage:       core.BaseDamageConfigFlat(tickDamage),
						OutcomeApplier:   warrior.OutcomeFuncTick(),
					})),
				})
				dwDots = append(dwDots, dot)
			}
		},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				deepWoundsSpell.Cast(sim, nil)
				deepWoundsSpell.SpellMetrics[spellEffect.Target.TableIndex].Hits++
				dwDots[spellEffect.Target.Index].Apply(sim)
				warrior.procBloodFrenzy(sim, spellEffect, time.Second*12)
			}
		},
	})
}
