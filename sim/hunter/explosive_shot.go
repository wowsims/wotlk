package hunter

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (hunter *Hunter) registerExplosiveShotSpell(timer *core.Timer) {
	if !hunter.Talents.ExplosiveShot {
		return
	}

	hunter.ExplosiveShotR4, hunter.ExplosiveShotR4Dot = hunter.makeExplosiveShotSpell(timer, false)
	hunter.ExplosiveShotR3, hunter.ExplosiveShotR3Dot = hunter.makeExplosiveShotSpell(timer, true)
}

func (hunter *Hunter) makeExplosiveShotSpell(timer *core.Timer, downrank bool) (*core.Spell, *core.Dot) {
	actionID := core.ActionID{SpellID: 60053}
	minFlatDamage := 386.0
	maxFlatDamage := 464.0
	if downrank {
		actionID = core.ActionID{SpellID: 60052}
		minFlatDamage = 325.0
		maxFlatDamage = 391.0
	}

	var esDot *core.Dot
	esSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		Cost: core.NewManaCost(core.ManaCostOptions{
			BaseCost:   0.07,
			Multiplier: 1 - 0.03*float64(hunter.Talents.Efficiency),
		}),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * 6,
			},
		},

		BonusCritRating: 0 +
			2*core.CritRatingPerCritChance*float64(hunter.Talents.SurvivalInstincts) +
			core.TernaryFloat64(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfExplosiveShot), 4*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1 +
			.02*float64(hunter.Talents.TNT),
		DamageMultiplier: 1,
		CritMultiplier:   hunter.critMultiplier(true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeRangedHit)

			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				esDot.Apply(sim)
				esDot.TickOnce(sim)
			}
		},
	})

	target := hunter.CurrentTarget
	esDot = core.NewDot(core.Dot{
		Spell: esSpell,
		Aura: target.RegisterAura(core.Aura{
			Label:    fmt.Sprintf("ExplosiveShot-%d-%d", actionID.SpellID, hunter.Index),
			ActionID: actionID,
		}),
		NumberOfTicks: 2,
		TickLength:    time.Second * 1,
		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			dot.SnapshotBaseDamage = sim.Roll(minFlatDamage, maxFlatDamage) + 0.14*dot.Spell.RangedAttackPower(target)
			attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
			dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(target, attackTable)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeRangedHitAndCritSnapshot)
		},
	})

	return esSpell, esDot
}
