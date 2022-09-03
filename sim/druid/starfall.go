package druid

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// We register two spells to apply two different dot effects and get two entries in Damage/Detailed results
func (druid *Druid) registerStarfallSpell() {
	if !druid.Talents.Starfall {
		return
	}

	baseCost := druid.BaseMana * 0.35
	target := druid.CurrentTarget
	numberOfTicks := core.TernaryInt(druid.Env.GetNumTargets() > 1, 20, 10)
	tickLength := core.TernaryDuration(druid.Env.GetNumTargets() > 1, time.Millisecond*500, time.Millisecond*1000)

	// Improved Faerie Fire and Nature's Majesty
	iffCritBonus := core.TernaryFloat64(druid.CurrentTarget.HasAura("Improved Faerie Fire"), druid.TalentsBonuses.iffBonusCrit, 0)
	naturesMajestyCritBonus := druid.TalentsBonuses.naturesMajestyBonusCrit

	druid.Starfall = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53201},
		SpellSchool: core.SpellSchoolArcane,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * druid.TalentsBonuses.moonglowMultiplier,
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * (90 - core.TernaryDuration(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfStarfall), 30, 0)),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:       core.ProcMaskSpellDamage,
			OutcomeApplier: druid.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.StarfallDot.Apply(sim)
					druid.StarfallDotSplash.Apply(sim)
				}
			},
		}),
	})

	druid.StarfallSplash = druid.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 53190},
	})

	druid.StarfallDot = core.NewDot(core.Dot{
		Spell: druid.Starfall,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Starfall-" + strconv.Itoa(int(druid.Index)),
			ActionID: core.ActionID{SpellID: 53201},
		}),
		NumberOfTicks: numberOfTicks,
		TickLength:    tickLength,
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * (1 + core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFocus), 0.1, 0)),
			ThreatMultiplier: 1,
			IsPeriodic:       false,
			BaseDamage:       core.BaseDamageConfigMagic(563, 653, 0.3),
			OutcomeApplier:   druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, druid.TalentsBonuses.vengeanceModifier)),
			BonusCritRating:  iffCritBonus + naturesMajestyCritBonus,
		})),
	})

	druid.StarfallDotSplash = core.NewDot(core.Dot{
		Spell: druid.StarfallSplash,
		Aura: target.RegisterAura(core.Aura{
			Label:    "StarfallSplash-" + strconv.Itoa(int(druid.Index)),
			ActionID: core.ActionID{SpellID: 53190},
		}),
		NumberOfTicks: numberOfTicks,
		TickLength:    tickLength,
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncAOEDamageCapped(druid.Env, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * (1 + core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFocus), 0.1, 0)),
			ThreatMultiplier: 1,
			IsPeriodic:       false,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(101, 0.13),
			OutcomeApplier:   druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, druid.TalentsBonuses.vengeanceModifier)),
			BonusCritRating:  iffCritBonus + naturesMajestyCritBonus,
		})),
	})
}
