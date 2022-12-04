package druid

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// We register two spells to apply two different dot effects and get two entries in Damage/Detailed results
func (druid *Druid) registerStarfallSpell() {
	if !druid.Talents.Starfall {
		return
	}

	baseCost := druid.BaseMana * 0.35
	target := druid.CurrentTarget
	numberOfTicks := core.TernaryInt32(druid.Env.GetNumTargets() > 1, 20, 10)
	tickLength := core.TernaryDuration(druid.Env.GetNumTargets() > 1, time.Millisecond*500, time.Millisecond*1000)

	druid.Starfall = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 53201},
		SpellSchool:  core.SpellSchoolArcane,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagNaturesGrace,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.03*float64(druid.Talents.Moonglow)),
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * (90 - core.TernaryDuration(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfStarfall), 30, 0)),
			},
		},

		BonusCritRating:  2 * float64(druid.Talents.NaturesMajesty) * core.CritRatingPerCritChance,
		DamageMultiplier: 1 * (1 + core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFocus), 0.1, 0)),
		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				druid.StarfallDot.Apply(sim)
				druid.StarfallDotSplash.Apply(sim)
			}
		},
	})

	druid.StarfallSplash = druid.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 53190},
		SpellSchool:      core.SpellSchoolArcane,
		ProcMask:         core.ProcMaskSpellDamage,
		BonusCritRating:  2 * float64(druid.Talents.NaturesMajesty) * core.CritRatingPerCritChance,
		DamageMultiplier: 1 * (1 + core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFocus), 0.1, 0)),
		CritMultiplier:   druid.BalanceCritMultiplier(),
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

		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			baseDamage := sim.Roll(563, 653) + 0.3*dot.Spell.SpellPower()
			dot.Spell.CalcAndDealDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
		},
	})

	druid.StarfallDotSplash = core.NewDot(core.Dot{
		Spell: druid.StarfallSplash,
		Aura: target.RegisterAura(core.Aura{
			Label:    "StarfallSplash-" + strconv.Itoa(int(druid.Index)),
			ActionID: core.ActionID{SpellID: 53190},
		}),
		NumberOfTicks: numberOfTicks,
		TickLength:    tickLength,
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			baseDamage := 101 + 0.13*dot.Spell.SpellPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.Targets {
				dot.Spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
