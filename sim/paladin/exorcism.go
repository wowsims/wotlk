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

	baseModifiers := Multiplicative{
		Additive{
			paladin.getTalentSanctityOfBattleBonus(),
			paladin.getMajorGlyphOfExorcismBonus(),
			paladin.getItemSetAegisBattlegearBonus2(),
		},
	}
	baseMultiplier := baseModifiers.Get()

	scaling := hybridScaling{
		AP: 0.15,
		SP: 0.15,
	}

	paladin.Exorcism = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48801},
		SpellSchool: core.SpellSchoolHoly,

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

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,

			DamageMultiplier: baseMultiplier,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					// TODO: discuss exporting or adding to core for damageRollOptimized hybrid scaling.
					deltaDamage := 1146.0 - 1028.0
					damage := 1028.0 + deltaDamage*sim.RandomFloat("Damage Roll")
					damage += hitEffect.SpellPower(spell.Unit, spell) * scaling.SP
					damage += hitEffect.MeleeAttackPower(spell.Unit) * scaling.AP
					return damage
				},
			},

			OutcomeApplier: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, attackTable *core.AttackTable) {
				if spellEffect.MagicHitCheck(sim, spell, attackTable) {
					if spellEffect.Target.MobType == proto.MobType_MobTypeDemon || spellEffect.Target.MobType == proto.MobType_MobTypeUndead || spellEffect.MagicCritCheck(sim, spell, attackTable) {
						spellEffect.Outcome = core.OutcomeCrit
						spell.SpellMetrics[spellEffect.Target.TableIndex].Crits++
						spellEffect.Damage *= paladin.SpellCritMultiplier()
					} else {
						spellEffect.Outcome = core.OutcomeHit
						spell.SpellMetrics[spellEffect.Target.TableIndex].Hits++
					}
				} else {
					spellEffect.Outcome = core.OutcomeMiss
					spell.SpellMetrics[spellEffect.Target.TableIndex].Misses++
					spellEffect.Damage = 0
				}
			},
		}),
	})
}
