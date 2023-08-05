package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (priest *Priest) registerPrayerOfMendingSpell() {
	actionID := core.ActionID{SpellID: 48113}

	pomAuras := make([]*core.Aura, len(priest.Env.AllUnits))
	for _, unit := range priest.Env.AllUnits {
		if !priest.IsOpponent(unit) {
			pomAuras[unit.UnitIndex] = priest.makePrayerOfMendingAura(unit)
		}
	}

	maxJumps := 5 + core.TernaryInt(priest.HasSetBonus(ItemSetRegaliaOfFaith, 2), 1, 0)

	var curTarget *core.Unit
	var remainingJumps int
	priest.ProcPrayerOfMending = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseHealing := 1043 + 0.8057*spell.HealingPower(target)
		priest.PrayerOfMending.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

		pomAuras[target.UnitIndex].Deactivate(sim)
		curTarget = nil

		// Bounce to new ally.
		if remainingJumps == 0 {
			return
		}

		// Find ally with lowest % HP and is not the current mending target.
		var newTarget *core.Unit
		for _, raidUnit := range priest.Env.Raid.AllUnits {
			if raidUnit == target {
				continue
			}

			if newTarget == nil || (raidUnit.HasHealthBar() && newTarget.HasHealthBar() && raidUnit.CurrentHealthPercent() < newTarget.CurrentHealthPercent()) {
				newTarget = raidUnit
			}
		}

		if newTarget != nil {
			pomAuras[newTarget.UnitIndex].Activate(sim)
			curTarget = newTarget
			remainingJumps--
		}
	}

	priest.PrayerOfMending = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.15,
			Multiplier: 1 *
				(1 - .1*float64(priest.Talents.HealingPrayers)) *
				(1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Duration(float64(time.Second*10) * (1 - .06*float64(priest.Talents.DivineProvidence))),
			},
		},

		BonusCritRating: float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.SpiritualHealing)) *
			(1 + .01*float64(priest.Talents.BlessedResilience)) *
			(1 + .02*float64(priest.Talents.FocusedPower)) *
			(1 + .02*float64(priest.Talents.DivineProvidence)) *
			(1 + .01*float64(priest.Talents.TwinDisciplines)) *
			core.TernaryFloat64(priest.HasSetBonus(ItemSetZabrasRaiment, 2), 1.2, 1),
		CritMultiplier:   priest.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if curTarget != nil {
				pomAuras[curTarget.UnitIndex].Deactivate(sim)
			}

			pomAuras[target.UnitIndex].Activate(sim)
			curTarget = target
			remainingJumps = maxJumps
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
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !autoProc && result.Damage > 0 {
				priest.ProcPrayerOfMending(sim, aura.Unit, priest.PrayerOfMending)
			}
		},
	})
}
