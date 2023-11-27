package mage

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
)

// TODO: Classic verify numbers such as aoe caps and base damage
// https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044#news-post-336044
// https://www.wowhead.com/classic/spell=400614/living-bomb
// https://www.wowhead.com/classic/spell=400613/living-bomb
func (mage *Mage) registerLivingBombSpell() {
	if !mage.HasRuneById(MageRuneHandsLivingBomb) {
		return
	}

	level := float64(mage.GetCharacter().Level)
	baseCalc := (13.828124 + 0.018012*level + 0.044141*level*level)
	ticks := int32(4)
	baseDotDamage := baseCalc * 3.4
	baseExplosionDamage := baseCalc * 1.71

	livingBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 400614},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage,

		BonusCritRating:          0,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1 - 0.1*float64(mage.Talents.BurningSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := baseExplosionDamage + 0.4*spell.SpellPower()
			// baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeExpectedMagicCrit)
			}
		},
	})

	mage.LivingBomb = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 400613},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.22,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusCritRating:  0,
		BonusHitRating:   float64(mage.Talents.ElementalPrecision) * 2 * core.SpellHitRatingPerHitChance,
		ThreatMultiplier: 1 - 0.15*float64(mage.Talents.BurningSoul),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "LivingBomb",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					livingBombExplosionSpell.Cast(sim, aura.Unit)
				},
			},

			NumberOfTicks: ticks,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = (baseDotDamage / float64(ticks)) + 0.2*dot.Spell.SpellPower()
				dot.SnapshotAttackerMultiplier = 1 // dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
