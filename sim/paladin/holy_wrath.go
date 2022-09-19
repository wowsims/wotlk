package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerHolyWrathSpell() {
	// From the perspective of max rank.
	baseCost := paladin.BaseMana * 0.20

	baseEffect := core.SpellEffect{
		ProcMask: core.ProcMaskSpellDamage,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				// TODO: discuss exporting or adding to core for damageRollOptimized hybrid scaling.
				deltaDamage := 1234.0 - 1050.0
				damage := 1050.0 + deltaDamage*sim.RandomFloat("Damage Roll")
				damage += hitEffect.SpellPower(spell.Unit, spell) * 0.07
				damage += hitEffect.MeleeAttackPower(spell.Unit) * 0.07
				return damage
			},
		},

		OutcomeApplier: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, attackTable *core.AttackTable) {
			// HW misses on non-undead/demons
			if !(spellEffect.Target.MobType == proto.MobType_MobTypeDemon || spellEffect.Target.MobType == proto.MobType_MobTypeUndead) {
				spellEffect.Outcome = core.OutcomeMiss
				spell.SpellMetrics[spellEffect.Target.UnitIndex].Misses++
				spellEffect.Damage = 0
				return
			}

			if spellEffect.MagicHitCheck(sim, spell, attackTable) {
				if spellEffect.MagicCritCheck(sim, spell, attackTable) {
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
	}

	numTargets := paladin.Env.GetNumTargets()
	effects := make([]core.SpellEffect, 0, paladin.Env.GetNumTargets())

	for i := int32(0); i < numTargets; i++ {
		effect := baseEffect
		effect.Target = paladin.Env.GetTargetUnit(i)
		effects = append(effects, effect)
	}

	paladin.HolyWrath = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48817},
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second*30 - core.TernaryDuration(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHolyWrath), time.Second*15, 0),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}
