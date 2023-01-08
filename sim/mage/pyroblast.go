package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerPyroblastSpell() {
	if !mage.Talents.Pyroblast {
		return
	}

	actionID := core.ActionID{SpellID: 42891}
	baseCost := .22 * mage.BaseMana
	spellCoeff := 1.15 + 0.05*float64(mage.Talents.EmpoweredFire)
	tickCoeff := 0.05 + 0.05*float64(mage.Talents.EmpoweredFire)

	hasT8_4pc := mage.HasSetBonus(ItemSetKirinTorGarb, 4)
	t10ProcAura := mage.BloodmagesRegalia2pcAura()

	mage.Pyroblast = mage.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage,
		MissileSpeed: 24,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,

				GCD:      core.GCDDefault,
				CastTime: time.Second * 5,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if mage.HotStreakAura.IsActive() {
					cast.CastTime = 0
					// cast.AfterCastDelay could be used for CQS to avoid ignite munching. Going to wait to implement for now though
					if !hasT8_4pc || sim.RandomFloat("MageT84PC") > T84PcProcChance {
						mage.HotStreakAura.Deactivate(sim)
					}
					if t10ProcAura != nil {
						t10ProcAura.Activate(sim)
					}
				}
			},
		},

		BonusCritRating: 0 +
			2*float64(mage.Talents.CriticalMass)*core.CritRatingPerCritChance +
			2*float64(mage.Talents.WorldInFlames)*core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .04*float64(mage.Talents.TormentTheWeak)),
		DamageMultiplierAdditive: 1 +
			.02*float64(mage.Talents.SpellImpact) +
			.02*float64(mage.Talents.FirePower),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(1210, 1531) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					mage.PyroblastDot.Apply(sim)
				}
				spell.DealDamage(sim, result)
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

			DamageMultiplier:         mage.Pyroblast.DamageMultiplier,
			DamageMultiplierAdditive: mage.Pyroblast.DamageMultiplierAdditive,
			ThreatMultiplier:         mage.Pyroblast.ThreatMultiplier,
		}),
		Aura: target.RegisterAura(core.Aura{
			Label:    "Pyroblast-" + strconv.Itoa(int(mage.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 4,
		TickLength:    time.Second * 3,
		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 113.0 + tickCoeff*dot.Spell.SpellPower()
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
		},
	})
}
