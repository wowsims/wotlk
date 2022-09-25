package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerExorcismSpell() {
	// From the perspective of max rank.
	baseCost := paladin.BaseMana * 0.08

	paladin.Exorcism = paladin.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48801},
		SpellSchool:  core.SpellSchoolHoly,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 15,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if paladin.ArtOfWarInstantCast.IsActive() {
					paladin.ArtOfWarInstantCast.Deactivate(sim)
					cast.CastTime = 0
					cast.Cost = cast.Cost * (1 - 0.02*float64(paladin.Talents.Benediction))
				}
			},
		},

		DamageMultiplierAdditive: 1 +
			paladin.getTalentSanctityOfBattleBonus() +
			paladin.getMajorGlyphOfExorcismBonus() +
			paladin.getItemSetAegisBattlegearBonus2(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					// TODO: discuss exporting or adding to core for damageRollOptimized hybrid scaling.
					deltaDamage := 1146.0 - 1028.0
					return 1028.0 + deltaDamage*sim.RandomFloat("Damage Roll") +
						.15*spell.SpellPower() +
						.15*spell.MeleeAttackPower()
				},
			},

			OutcomeApplier: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, attackTable *core.AttackTable) {
				if spell.MagicHitCheck(sim, attackTable) {
					if spellEffect.Target.MobType == proto.MobType_MobTypeDemon || spellEffect.Target.MobType == proto.MobType_MobTypeUndead || spellEffect.MagicCritCheck(sim, spell, attackTable) {
						spellEffect.Outcome = core.OutcomeCrit
						spell.SpellMetrics[spellEffect.Target.UnitIndex].Crits++
						spellEffect.Damage *= paladin.SpellCritMultiplier()
					} else {
						spellEffect.Outcome = core.OutcomeHit
						spell.SpellMetrics[spellEffect.Target.UnitIndex].Hits++
					}
				} else {
					spellEffect.Outcome = core.OutcomeMiss
					spell.SpellMetrics[spellEffect.Target.UnitIndex].Misses++
					spellEffect.Damage = 0
				}
			},
		}),
	})
}
