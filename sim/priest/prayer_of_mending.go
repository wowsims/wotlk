package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerPrayerOfMendingSpell() {
	actionID := core.ActionID{SpellID: 48113}
	baseCost := 0.15 * priest.BaseMana

	pomAuras := make([]*core.Aura, len(priest.Env.AllUnits))
	for _, unit := range priest.Env.AllUnits {
		if !priest.IsOpponent(unit) {
			pomAuras[unit.UnitIndex] = priest.makePrayerOfMendingAura(unit)
		}
	}

	var curTarget *core.Unit
	priest.ProcPrayerOfMending = core.ApplyEffectFuncDirectDamage(core.SpellEffect{
		IsHealing: true,
		ProcMask:  core.ProcMaskSpellHealing,

		BonusCritRating: float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.DivineProvidence)) *
			(1 + .01*float64(priest.Talents.TwinDisciplines)),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		BaseDamage:     core.BaseDamageConfigHealingNoRoll(1043, 0.8057),
		OutcomeApplier: priest.OutcomeFuncHealingCrit(priest.DefaultHealingCritMultiplier()),

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			pomAuras[spellEffect.Target.UnitIndex].Deactivate(sim)
			curTarget = nil

			// Bounce to new ally.

			// Find ally with lowest % HP and is not the current mending target.
			var newTarget *core.Unit
			for _, raidUnit := range priest.Env.Raid.AllUnits {
				if raidUnit == spellEffect.Target {
					continue
				}

				if newTarget == nil || (raidUnit.HasHealthBar() && newTarget.HasHealthBar() && raidUnit.CurrentHealthPercent() < newTarget.CurrentHealthPercent()) {
					newTarget = raidUnit
				}
			}

			if newTarget != nil {
				pomAuras[newTarget.UnitIndex].Activate(sim)
				curTarget = newTarget
			}
		},
	})

	priest.PrayerOfMending = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - .1*float64(priest.Talents.HealingPrayers)) *
					(1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Duration(float64(time.Second*10) * (1 - .06*float64(priest.Talents.DivineProvidence))),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if curTarget != nil {
				pomAuras[curTarget.UnitIndex].Deactivate(sim)
			}

			pomAuras[target.UnitIndex].Activate(sim)
			curTarget = target
		},
	})
}

func (priest *Priest) makePrayerOfMendingAura(target *core.Unit) *core.Aura {
	autoProc := true

	return target.RegisterAura(core.Aura{
		Label:    "PrayerOfMending" + strconv.Itoa(int(priest.Index)),
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if autoProc {
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + time.Second*5,
					OnAction: func(sim *core.Simulation) {
						priest.ProcPrayerOfMending(sim, aura.Unit, priest.PrayerOfMending)
					},
				})
			}
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !autoProc && spellEffect.Damage > 0 {
				priest.ProcPrayerOfMending(sim, aura.Unit, priest.PrayerOfMending)
			}
		},
	})
}
