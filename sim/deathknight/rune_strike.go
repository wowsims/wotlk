package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerRuneStrikeSpell() {
	actionID := core.ActionID{SpellID: 56815}

	runeStrikeGlyphCritBonus := core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfRuneStrike), 10.0, 0.0)

	baseCost := float64(core.NewRuneCost(20, 0, 0, 0, 0))
	rs := &RuneSpell{}
	dk.RuneStrike = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			IgnoreHaste: true,
		},

		BonusCritRating: (dk.annihilationCritBonus() + runeStrikeGlyphCritBonus) * core.CritRatingPerCritChance,
		DamageMultiplier: 1.5 *
			dk.darkrunedPlateRuneStrikeDamageBonus(),
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.75,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				0.15*spell.MeleeAttackPower() +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()
			baseDamage *= dk.RoRTSBonus(target)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)

			rs.DoCost(sim)
			dk.RuneStrikeAura.Deactivate(sim)
		},
	}, func(sim *core.Simulation) bool {
		runeCost := core.RuneCost(dk.RuneStrike.BaseCost)
		return dk.CastCostPossible(sim, float64(runeCost.RunicPower()), 0, 0, 0) && dk.RuneStrike.IsReady(sim) && dk.RuneStrikeAura.IsActive() && dk.CurrentRunicPower() >= float64(runeCost.RunicPower())
	})

	dk.RuneStrikeAura = dk.RegisterAura(core.Aura{
		Label:    "Rune Strike",
		ActionID: actionID,
		Duration: 6 * time.Second,
	})

	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label:    "Rune Strike Trigger",
		Duration: core.NeverExpires,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeDodge | core.OutcomeParry) {
				dk.RuneStrikeAura.Activate(sim)
			}
		},
	}))
}
