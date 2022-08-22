package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var DeathCoilActionID = core.ActionID{SpellID: 49895}

func (dk *Deathknight) registerDeathCoilSpell() {
	baseDamage := 443.0 + dk.sigilOfTheWildBuckBonus() + dk.sigilOfTheVengefulHeartDeathCoil()
	baseCost := float64(core.NewRuneCost(40, 0, 0, 0, 0))
	dk.DeathCoil = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:    DeathCoilActionID,
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if dk.SuddenDoomAura.IsActive() {
					cast.GCD = 0
					cast.Cost = 0
				} else {
					cast.GCD = dk.getModifiedGCD()
				}
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
					return (baseDamage + dk.getImpurityBonus(hitEffect, spell.Unit)*0.15) * dk.RoRTSBonus(hitEffect.Target)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.OutcomeFuncMagicHitAndCrit(dk.spellCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				dk.LastOutcome = spellEffect.Outcome
				if spellEffect.Landed() && dk.Talents.UnholyBlight {
					dk.procUnholyBlight(sim, spellEffect.Target, spellEffect.Damage)
				}
			},
		}),
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 40.0, 0, 0, 0) && dk.DeathCoil.IsReady(sim)
	}, nil)
}

func (dk *Deathknight) registerDrwDeathCoilSpell() {
	baseDamage := 443.0 + dk.sigilOfTheWildBuckBonus() + dk.sigilOfTheVengefulHeartDeathCoil()

	dk.RuneWeapon.DeathCoil = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    DeathCoilActionID,
		SpellSchool: core.SpellSchoolShadow,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: dk.darkrunedBattlegearCritBonus() * core.CritRatingPerCritChance,
			DamageMultiplier: (1.0 + float64(dk.Talents.Morbidity)*0.05) *
				core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDarkDeath), 1.15, 1.0),
			ThreatMultiplier: 1.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (baseDamage + dk.RuneWeapon.getImpurityBonus(hitEffect, spell.Unit)*0.15)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.RuneWeapon.OutcomeFuncMagicHitAndCrit(dk.RuneWeapon.MeleeCritMultiplier(1.0, 0.0)),
		}),
	})
}
