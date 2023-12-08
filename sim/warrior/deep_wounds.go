package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warrior *Warrior) applyDeepWounds() {
	if warrior.Talents.DeepWounds == 0 {
		return
	}

	spellID := map[int32]int32{
		1: 12834,
		2: 12849,
		3: 12867,
	}[warrior.Talents.DeepWounds]

	warrior.DeepWounds = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers | SpellFlagBleed,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DeepWounds",
			},
			NumberOfTicks: 12,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				attackTable := warrior.AttackTables[target.UnitIndex]
				adm := warrior.AutoAttacks.MHAuto().AttackerDamageMultiplier(attackTable)
				tdm := warrior.AutoAttacks.MHAuto().TargetDamageMultiplier(attackTable, false)
				awd := (warrior.AutoAttacks.MH().CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower()) + dot.Spell.BonusWeaponDamage()) * adm * tdm
				newDamage := awd * 0.2 * float64(warrior.Talents.DeepWounds)
				dot.SnapshotBaseDamage = newDamage / 12
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.SnapshotAttackerMultiplier = target.PseudoStats.PeriodicPhysicalDamageTakenMultiplier
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
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
				warrior.DeepWounds.Cast(sim, result.Target)
			}
		},
	})
}

func (warrior *Warrior) procDeepWounds(sim *core.Simulation, target *core.Unit, isOh bool) {
	dot := warrior.DeepWounds.Dot(target)

	outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)

	attackTable := warrior.AttackTables[target.UnitIndex]
	var awd float64
	if isOh {
		adm := warrior.AutoAttacks.OHAuto().AttackerDamageMultiplier(attackTable)
		tdm := warrior.AutoAttacks.OHAuto().TargetDamageMultiplier(attackTable, false)
		awd = ((warrior.AutoAttacks.OH().CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower()) * 0.5) + dot.Spell.BonusWeaponDamage()) * adm * tdm
	} else { // MH, Ranged (e.g. Thunder Clap)
		adm := warrior.AutoAttacks.MHAuto().AttackerDamageMultiplier(attackTable)
		tdm := warrior.AutoAttacks.MHAuto().TargetDamageMultiplier(attackTable, false)
		awd = (warrior.AutoAttacks.MH().CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower()) + dot.Spell.BonusWeaponDamage()) * adm * tdm
	}
	newDamage := awd * 0.16 * float64(warrior.Talents.DeepWounds)

	dot.SnapshotBaseDamage = (outstandingDamage + newDamage) / float64(dot.NumberOfTicks)
	dot.SnapshotAttackerMultiplier = 1
	warrior.DeepWounds.Cast(sim, target)
}
