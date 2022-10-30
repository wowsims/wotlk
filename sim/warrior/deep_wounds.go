package warrior

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
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
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreAttackerModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Deep Wounds Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskEmpty) || !spell.SpellSchool.Matches(core.SpellSchoolPhysical) {
				return
			}
			if result.Outcome.Matches(core.OutcomeCrit) {
				warrior.DeepWounds.Cast(sim, nil)
				warrior.procDeepWounds(sim, result.Target, spell.IsMH())
				warrior.procBloodFrenzy(sim, result, time.Second*6)
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

		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			baseDamage := warrior.DeepWoundsTickDamage[target.Index]
			dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDamage, dot.OutcomeTick)
			// TODO: This probably can go before the calc now (after we assign baseDamage) but it breaks 1 test.
			warrior.DeepWoundsDamageBuffer[target.Index] -= warrior.DeepWoundsTickDamage[target.Index]
		},
	})
}

// TODO (maybe) https://github.com/magey/wotlk-warrior/issues/26 - Deep Wounds is not benefitting from Blood Frenzy
func (warrior *Warrior) procDeepWounds(sim *core.Simulation, target *core.Unit, isMh bool) {
	dot := warrior.DeepWoundsDots[target.Index]

	dotDamageMultiplier := 0.16 * float64(warrior.Talents.DeepWounds) * warrior.PseudoStats.DamageDealtMultiplier * warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]
	if isMh {
		dotDamage := (warrior.AutoAttacks.MH.CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower()) + dot.Spell.BonusWeaponDamage()) * dotDamageMultiplier
		warrior.DeepWoundsDamageBuffer[target.Index] += dotDamage
	} else {
		dwsMultiplier := 1 + 0.05*float64(warrior.Talents.DualWieldSpecialization)
		dotDamage := ((warrior.AutoAttacks.OH.CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower()) * 0.5) + dot.Spell.BonusWeaponDamage()) * dwsMultiplier * dotDamageMultiplier
		warrior.DeepWoundsDamageBuffer[target.Index] += dotDamage
	}

	warrior.DeepWoundsTickDamage[target.Index] = warrior.DeepWoundsDamageBuffer[target.Index] / 6
	warrior.DeepWounds.SpellMetrics[target.UnitIndex].Hits++

	dot.Apply(sim)
}
