package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// We register two spells to apply two different dot effects and get two entries in Damage/Detailed results
func (druid *Druid) registerStarfallSpell() {
	if !druid.Talents.Starfall {
		return
	}

	numberOfTicks := core.TernaryInt32(druid.Env.GetNumTargets() > 1, 20, 10)
	tickLength := time.Second

	starfallTickSpell := druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 53195},
		SpellSchool:      core.SpellSchoolArcane,
		ProcMask:         core.ProcMaskSuppressedProc,
		Flags:            SpellFlagNaturesGrace,
		BonusCritRating:  2 * float64(druid.Talents.NaturesMajesty) * core.CritRatingPerCritChance,
		DamageMultiplier: 1 * (1 + core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFocus), 0.1, 0)),
		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(563, 653) + 0.3*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})

	druid.Starfall = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53201},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | SpellFlagOmenTrigger,
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.35,
			Multiplier: 1 - 0.03*float64(druid.Talents.Moonglow),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * (90 - core.TernaryDuration(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfStarfall), 30, 0)),
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Starfall",
			},
			NumberOfTicks: numberOfTicks,
			TickLength:    tickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				starfallTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				druid.StarfallSplash.Dot(target).Apply(sim)
			}
		},
	})

	starfallSplashTickSpell := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 53190},
		SpellSchool:      core.SpellSchoolArcane,
		ProcMask:         core.ProcMaskSuppressedProc,
		BonusCritRating:  2 * float64(druid.Talents.NaturesMajesty) * core.CritRatingPerCritChance,
		DamageMultiplier: 1 * (1 + core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFocus), 0.1, 0)),
		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 101 + 0.13*spell.SpellPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	druid.StarfallSplash = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 53190},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "StarfallSplash",
			},
			NumberOfTicks: numberOfTicks,
			TickLength:    tickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				starfallSplashTickSpell.Cast(sim, target)
			},
		},
	})
}
