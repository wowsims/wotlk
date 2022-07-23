package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerHammerOfWrathSpell() {
	// From the perspective of max rank.
	baseCost := paladin.BaseMana * 0.12

	baseModifiers := Modifiers{
		{
			core.TernaryFloat64(paladin.HasSetBonus(ItemSetAegisBattlegear, 2), .1, 0),
		},
	}
	baseMultiplier := baseModifiers.Get()

	scaling := hybridScaling{
		AP: 0.15,
		SP: 0.15,
	}

	paladin.HammerOfWrath = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48806},
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
				Duration: time.Second * 6,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,

			DamageMultiplier: baseMultiplier,
			ThreatMultiplier: 1,
			BonusCritRating:  (25 * float64(paladin.Talents.SanctifiedWrath)) * core.CritRatingPerCritChance,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					// TODO: discuss exporting or adding to core for damageRollOptimized hybrid scaling.
					deltaDamage := 1257.0 - 1139.0
					damage := 1139.0 + deltaDamage*sim.RandomFloat("Damage Roll")
					damage += hitEffect.SpellPower(spell.Unit, spell) * scaling.SP
					damage += hitEffect.MeleeAttackPower(spell.Unit) * scaling.AP
					return damage
				},
			},
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialNoBlockDodgeParry(paladin.MeleeCritMultiplier()),
		}),
	})
}
