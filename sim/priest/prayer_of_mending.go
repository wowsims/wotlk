package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (priest *Priest) registerPrayerOfMendingSpell() {
	actionID := core.ActionID{SpellID: 48113}

	pomAuras := make([]*core.Aura, len(priest.Env.AllUnits))
	for _, unit := range priest.Env.AllUnits {
		if !priest.IsOpponent(unit) {
			pomAuras[unit.UnitIndex] = priest.makePrayerOfMendingAura(unit)
		}
	}

	maxJumps := 5

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
			BaseCost:   0.15,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		BonusCritRating:  float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + .02*float64(priest.Talents.SpiritualHealing),
		CritMultiplier:   priest.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1,

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
