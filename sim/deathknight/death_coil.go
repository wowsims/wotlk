package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerDeathCoilSpell() {
	baseDamage := 443.0 + dk.sigilOfTheWildBuckBonus()

	dk.DeathCoil = dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49895},
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.RunicPower,
		BaseCost:     40.0,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: 40,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: dk.darkrunedBattlegearCritBonus() * core.CritRatingPerCritChance,
			DamageMultiplier: (1.0 + float64(dk.Talents.Morbidity)*0.05) *
				core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDarkDeath), 1.15, 1.0),
			ThreatMultiplier: 1.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (baseDamage + dk.applyImpurity(hitEffect, spell.Unit)*0.15) *
						dk.rageOfRivendareBonus(hitEffect.Target) *
						dk.tundraStalkerBonus(hitEffect.Target)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.OutcomeFuncMagicHitAndCrit(dk.spellCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				dk.LastCastOutcome = spellEffect.Outcome
				if spellEffect.Landed() && dk.Talents.UnholyBlight {
					dk.LastDeathCoilDamage = spellEffect.Damage
					dk.UnholyBlightSpell.Cast(sim, spellEffect.Target)
				}
			},
		}),
	})
}

func (dk *Deathknight) CanDeathCoil(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 40.0, 0, 0, 0) && dk.DeathCoil.IsReady(sim)
}

func (dk *Deathknight) CastDeathCoil(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanDeathCoil(sim) {
		dk.DeathCoil.Cast(sim, target)
		return true
	}
	return false
}
