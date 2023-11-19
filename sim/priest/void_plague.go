package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (priest *Priest) getVoidPlagueConfig() core.SpellConfig {
	// TODO: Classic SOD live check
	spellCoeff := 0.2
	baseDamage := float64(39 * 6) // https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: 425204},
		SpellSchool:   core.SpellSchoolShadow,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL,
		Rank:          1,
		RequiredLevel: 1,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.13,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		BonusHitRating:   float64(priest.Talents.ShadowFocus) * 2 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   1,
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "VoidPlague-" + strconv.Itoa(1),
			},

			NumberOfTicks: 8,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = baseDamage/6 + (spellCoeff * dot.Spell.SpellPower())
				dot.SnapshotAttackerMultiplier = 1
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				priest.AddShadowWeavingStack(sim)
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDamage := baseDamage/6 + (spellCoeff * spell.SpellPower())
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	}
}

func (priest *Priest) registerVoidPlagueSpell() {
	if !priest.HasRuneById(PriestRuneChestVoidPlague) {
		return
	}
	priest.VoidPlague = priest.GetOrRegisterSpell(priest.getVoidPlagueConfig())
}
