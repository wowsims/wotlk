package mage

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
	"github.com/wowsims/classic/sod/sim/core/proto"
)

// TODO: Classic verify numbers / snapshot / travel time
// https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044#news-post-336044
// https://www.wowhead.com/classic/spell=401558/living-flame
func (mage *Mage) registerLivingFlameSpell() {
	if !mage.HasRune(proto.MageRune_RuneLegsLivingFlame) {
		return
	}

	level := float64(mage.GetCharacter().Level)
	baseCalc := (13.828124 + 0.018012*level + 0.044141*level*level)
	baseDamage := baseCalc * 1

	mage.LivingFlame = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 401556},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.11,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 120,
			},
		},

		BonusCritRating:  0,
		BonusHitRating:   float64(mage.Talents.ElementalPrecision) * 2 * core.SpellHitRatingPerHitChance,
		ThreatMultiplier: 1 - 0.15*float64(mage.Talents.BurningSoul),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "LivingFlame",
			},

			NumberOfTicks: 20,
			TickLength:    time.Second * 1,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = baseDamage + 0.143*dot.Spell.SpellPower()
				dot.SnapshotAttackerMultiplier = 1
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					// TODO: Classic verify hit check, assuming dot damage with no hit check for now
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
			spell.Dot(target).Apply(sim)
		},
	})
}
