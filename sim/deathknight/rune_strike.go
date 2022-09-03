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
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: dk.darkrunedPlateRuneStrikeDamageBonus(),
			ThreatMultiplier: 1.75,
			BonusCritRating:  (dk.annihilationCritBonus() + runeStrikeGlyphCritBonus) * core.CritRatingPerCritChance,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					bonusDmg := 0.15 * (hitEffect.MeleeAttackPower(spell.Unit) + hitEffect.MeleeAttackPowerOnTarget())
					weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, bonusDmg, 1.5, 1.0, true)

					return weaponBaseDamage(sim, hitEffect, spell) *
						dk.RoRTSBonus(hitEffect.Target)
				},
				TargetSpellCoefficient: 1,
			},

			OutcomeApplier: dk.OutcomeFuncMeleeSpecialNoBlockDodgeParry(dk.critMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				rs.DoCost(sim)
				dk.RuneStrikeAura.Deactivate(sim)
			},
		}),
	}, func(sim *core.Simulation) bool {
		runeCost := core.RuneCost(dk.RuneStrike.BaseCost)
		return dk.CastCostPossible(sim, float64(runeCost.RunicPower()), 0, 0, 0) && dk.RuneStrike.IsReady(sim) && dk.RuneStrikeAura.IsActive() && dk.CurrentRunicPower() >= float64(runeCost.RunicPower())
	}, nil)

	dk.RuneStrikeAura = dk.RegisterAura(core.Aura{
		Label:    "Rune Strike",
		ActionID: actionID,
		Duration: 6 * time.Second,
	})

	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label:    "Rune Strike Trigger",
		Duration: core.NeverExpires,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeDodge | core.OutcomeParry) {
				dk.RuneStrikeAura.Activate(sim)
			}
		},
	}))
}
