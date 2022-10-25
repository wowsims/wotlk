package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerLivingBombSpell() {

	actionID := core.ActionID{SpellID: 55360}
	actionIDDot := core.ActionID{SpellID: 55359} // I want the dot to be separately trackable for metrics
	actionIDSpell := core.ActionID{SpellID: 44457}
	baseCost := .22 * mage.BaseMana
	bonusCrit := float64(mage.Talents.WorldInFlames+mage.Talents.CriticalMass) * 2 * core.CritRatingPerCritChance

	livingBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | HotStreakSpells,

		BonusCritRating:  bonusCrit,
		DamageMultiplier: mage.spellDamageMultiplier,
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 690 + (1.5/3.5)*spell.SpellPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.Targets {
				spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	target := mage.CurrentTarget
	hasGlyphOfLivingBomb := mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfLivingBomb)

	mage.LivingBomb = mage.RegisterSpell(core.SpellConfig{
		ActionID:     actionIDSpell,
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        SpellFlagMage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,

				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			// TODO: Uncomment this
			// if result.Landed() {
			mage.LivingBombDots[mage.CurrentTarget.Index].Apply(sim)
			// }
			spell.DealOutcome(sim, result)
		},
	})

	livingBombDotSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:         actionIDDot,
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            SpellFlagMage,
		Cast:             core.CastConfig{},
		BonusCritRating:  bonusCrit,
		DamageMultiplier: mage.spellDamageMultiplier,
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),
	})

	mage.LivingBombDots[target.Index] = core.NewDot(core.Dot{
		Spell: livingBombDotSpell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "LivingBomb-" + strconv.Itoa(int(mage.Index)),
			ActionID: actionID,
			Tag:      "LivingBomb",
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				mage.LivingBombNotActive.Dequeue()
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				livingBombExplosionSpell.Cast(sim, target)
				mage.LivingBombNotActive.Enqueue(target)
			},
		}),

		NumberOfTicks:       4,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: false,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 345 + 0.2*dot.Spell.SpellPower()
			dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			if hasGlyphOfLivingBomb {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeMagicHitAndSnapshotCrit)
			} else {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			}
		},
	})
}
