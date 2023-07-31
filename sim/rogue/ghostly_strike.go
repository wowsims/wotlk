package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) registerGhostlyStrikeSpell() {
	if !rogue.Talents.GhostlyStrike {
		return
	}

	hasGlyph := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfGhostlyStrike)

	actionID := core.ActionID{SpellID: 14278}

	rogue.GhostlyStrike = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded | core.SpellFlagAPL,
		EnergyCost: core.EnergyCostOptions{
			Cost:   40.0,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second*20 + core.TernaryDuration(hasGlyph, time.Second*10, 0),
			},
			IgnoreHaste: true,
		},

		BonusCritRating: []float64{0, 2, 4, 6}[rogue.Talents.TurnTheTables] * core.CritRatingPerCritChance,

		DamageMultiplier: core.TernaryFloat64(rogue.HasDagger(core.MainHand), 1.8, 1.25) * core.TernaryFloat64(hasGlyph, 1.4, 1) * (1 + 0.02*float64(rogue.Talents.FindWeakness)),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
