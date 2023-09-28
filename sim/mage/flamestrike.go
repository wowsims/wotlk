package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (mage *Mage) registerFlamestrikeSpell(rank8 bool) *core.Spell {
	actionID := core.ActionID{SpellID: 42926}.WithTag(9)
	dotDamage := 780.0 / 4
	minDamage := 876.0
	maxDamage := 1071.0
	spCoeffMultiplier := 1.0
	label := "Flamestrike (Rank 9)"
	if rank8 {
		actionID = core.ActionID{SpellID: 42925}.WithTag(8)
		dotDamage = 620.0 / 4
		minDamage = 699.0
		maxDamage = 854.0
		label = "Flamestrike (Rank 8)"
		// Flamestrike (Rank 8) has a 90% SP coefficient penalty
		// https://wowpedia.fandom.com/wiki/Downranking#Wrath_of_the_Lich_King
		spCoeffMultiplier = 0.9
	}

	return mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.30,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
		},

		BonusCritRating: float64(mage.Talents.CriticalMass+mage.Talents.WorldInFlames) * 2 * core.CritRatingPerCritChance,
		DamageMultiplierAdditive: 1 +
			.02*float64(mage.Talents.SpellImpact) +
			.02*float64(mage.Talents.FirePower),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.05*float64(mage.Talents.BurningSoul),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: label,
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
				target := mage.CurrentTarget
				dot.SnapshotBaseDamage = dotDamage + 0.122*dot.Spell.SpellPower()*spCoeffMultiplier
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmgFromSP := 0.243 * spell.SpellPower() * spCoeffMultiplier
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(minDamage, maxDamage) + dmgFromSP
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
			spell.AOEDot().Apply(sim)
		},
	})
}

func (mage *Mage) registerFlamestrikeSpells() {
	mage.Flamestrike = mage.registerFlamestrikeSpell(false)
	mage.FlamestrikeRank8 = mage.registerFlamestrikeSpell(true)
}
