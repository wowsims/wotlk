package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerPyroblastSpell() {
	actionID := core.ActionID{SpellID: 42891}
	baseCost := .22 * mage.BaseMana
	spellCoeff := 1.15 + 0.05*float64(mage.Talents.EmpoweredFire)

	mage.Pyroblast = mage.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage,
		MissileSpeed: 22,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,

				GCD:      core.GCDDefault,
				CastTime: time.Second * 6,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if mage.HotStreakAura.IsActive() {
					if mage.MageTier.t10_2 {
						bloodmageHasteAura.Activate(sim)
					}
					cast.CastTime = 0
					// cast.AfterCastDelay could be used for CQS to avoid ignite munching. Going to wait to implement for now though
					if !mage.MageTier.t8_4 || sim.RandomFloat("MageT84PC") > .1 {
						mage.HotStreakAura.Deactivate(sim)
					}
				}
			},
		},

		BonusCritRating: 0 +
			float64(mage.Talents.CriticalMass+mage.Talents.WorldInFlames)*2*core.CritRatingPerCritChance,
		DamageMultiplier: mage.spellDamageMultiplier * (1 + .04*float64(mage.Talents.TormentTheWeak)),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(1210, 1531) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					mage.PyroblastDot.Apply(sim)
				}
				spell.DealDamage(sim, &result)
			})
		},
	})

	target := mage.CurrentTarget
	mage.PyroblastDot = core.NewDot(core.Dot{
		Spell: mage.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       SpellFlagMage,

			DamageMultiplier: mage.Pyroblast.DamageMultiplier,
			ThreatMultiplier: mage.Pyroblast.ThreatMultiplier,
		}),
		Aura: target.RegisterAura(core.Aura{
			Label:    "Pyroblast-" + strconv.Itoa(int(mage.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 4,
		TickLength:    time.Second * 3,
		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 113.0 + 0.02*dot.Spell.SpellPower()
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
		},
	})
}
