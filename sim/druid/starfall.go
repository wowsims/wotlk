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

	// Nature's Majesty
	naturesMajestyCritBonus := druid.TalentsBonuses.naturesMajestyBonusCrit

	druid.Starfall = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 53201},
		SpellSchool:  core.SpellSchoolArcane,
		ProcMask:     core.ProcMaskSpellDamage,
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

		BonusCritRating:  naturesMajestyCritBonus,
		DamageMultiplier: 1 * (1 + core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFocus), 0.1, 0)),
		CritMultiplier:   druid.SpellCritMultiplier(1, druid.TalentsBonuses.vengeanceModifier),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
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
		ActionID:         core.ActionID{SpellID: 53190},
		ProcMask:         core.ProcMaskSpellDamage,
		BonusCritRating:  naturesMajestyCritBonus,
		DamageMultiplier: 1 * (1 + core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFocus), 0.1, 0)),
		CritMultiplier:   druid.SpellCritMultiplier(1, druid.TalentsBonuses.vengeanceModifier),
		ThreatMultiplier: 1,
	})

	druid.StarfallDot = core.NewDot(core.Dot{
		Spell: druid.Starfall,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Starfall-" + strconv.Itoa(int(druid.Index)),
			ActionID: core.ActionID{SpellID: 53201},
		}),
		NumberOfTicks: numberOfTicks,
		TickLength:    tickLength,
		TickEffects: core.TickFuncApplyEffects(func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(563, 653) + 0.3*spell.SpellPower()
			spell.CalcAndDealDamageMagicHitAndCrit(sim, target, baseDamage)
		}),
	})

	druid.StarfallDotSplash = core.NewDot(core.Dot{
		Spell: druid.StarfallSplash,
		Aura: target.RegisterAura(core.Aura{
			Label:    "StarfallSplash-" + strconv.Itoa(int(druid.Index)),
			ActionID: core.ActionID{SpellID: 53190},
		}),
		NumberOfTicks: numberOfTicks,
		TickLength:    tickLength,
		TickEffects: core.TickFuncApplyEffects(func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 101 + 0.13*spell.SpellPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.Targets {
				spell.CalcAndDealDamageMagicHitAndCrit(sim, &aoeTarget.Unit, baseDamage)
			}
		}),
	})
}
