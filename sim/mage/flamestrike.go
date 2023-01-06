package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerFlamestrikeSpell() {
	actionID := core.ActionID{SpellID: 42926}
	baseCost := .30 * mage.BaseMana

	mage.Flamestrike = mage.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.01*float64(mage.Talents.Pyromaniac)),

				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		BonusCritRating: 0 +
			float64(mage.Talents.CriticalMass)*2*core.CritRatingPerCritChance +
			float64(mage.Talents.Pyromaniac)*1*core.CritRatingPerCritChance,
		DamageMultiplier: mage.spellDamageMultiplier,
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.05*float64(mage.Talents.BurningSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmgFromSP := 0.2357 * spell.SpellPower()
			for _, aoeTarget := range sim.Encounter.Targets {
				baseDamage := sim.Roll(876, 1071) + dmgFromSP
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
			mage.FlamestrikeDot.Apply(sim)
		},
	})

	mage.FlamestrikeDot = core.NewDot(core.Dot{
		Spell: mage.Flamestrike,
		Aura: mage.RegisterAura(core.Aura{
			Label:    "Flamestrike",
			ActionID: actionID,
		}),
		NumberOfTicks: 4,
		TickLength:    time.Second * 2,
		OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
			target := mage.CurrentTarget
			dot.SnapshotBaseDamage = 780.0/4 + 0.122*dot.Spell.SpellPower()
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			for _, aoeTarget := range sim.Encounter.Targets {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, &aoeTarget.Unit, dot.OutcomeTick)
			}
		},
	})
}
