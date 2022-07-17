package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) registerDeathCoilSpell() {
	deathKnight.DeathCoil = deathKnight.RegisterSpell(core.SpellConfig{
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
				cast.GCD = deathKnight.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 0.0,
			DamageMultiplier: (1.0 + float64(deathKnight.Talents.Morbidity)*0.05) *
				core.TernaryFloat64(deathKnight.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfDarkDeath), 1.15, 1.0),
			ThreatMultiplier: 1.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (443.0 + deathKnight.applyImpurity(hitEffect, spell.Unit)*0.15) *
						deathKnight.rageOfRivendareBonus() *
						deathKnight.tundraStalkerBonus()
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncMagicHitAndCrit(deathKnight.spellCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() && deathKnight.Talents.UnholyBlight {
					deathKnight.LastDeathCoilDamage = spellEffect.Damage
					deathKnight.UnholyBlightSpell.Cast(sim, spellEffect.Target)
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanDeathCoil(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 40.0, 0, 0, 0) && deathKnight.DeathCoil.IsReady(sim)
}
