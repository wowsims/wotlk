package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (priest *Priest) getMindSearMiseryCoefficient() float64 {
	return 0.2861 * (1 + 0.05*float64(priest.Talents.Misery))
}

func (priest *Priest) getMindSearBaseConfig() core.SpellConfig {
	return core.SpellConfig{
		SpellSchool:     core.SpellSchoolShadow,
		ProcMask:        core.ProcMaskProc,
		BonusHitRating:  float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating: float64(priest.Talents.MindMelt) * 2 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 +
			0.02*float64(priest.Talents.Darkness) +
			0.01*float64(priest.Talents.TwinDisciplines),
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),
		CritMultiplier:   priest.DefaultSpellCritMultiplier(),
	}
}

func (priest *Priest) getMindSearTickSpell(numTicks int32) *core.Spell {
	hasGlyphOfShadow := priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow))
	miseryCoeff := priest.getMindSearMiseryCoefficient()

	config := priest.getMindSearBaseConfig()
	config.ActionID = core.ActionID{SpellID: 53022}.WithTag(numTicks)
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		damage := sim.Roll(212, 228) + miseryCoeff*spell.SpellPower()
		result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

		if result.Landed() {
			priest.AddShadowWeavingStack(sim)
		}
		if result.DidCrit() && hasGlyphOfShadow {
			priest.ShadowyInsightAura.Activate(sim)
		}
	}
	return priest.GetOrRegisterSpell(config)
}

func (priest *Priest) newMindSearSpell(numTicksIdx int32) *core.Spell {
	numTicks := numTicksIdx
	flags := core.SpellFlagChanneled | core.SpellFlagNoMetrics
	if numTicksIdx == 0 {
		numTicks = 5
		flags |= core.SpellFlagAPL
	}

	miseryCoeff := priest.getMindSearMiseryCoefficient()
	mindSearTickSpell := priest.getMindSearTickSpell(numTicksIdx)

	config := priest.getMindSearBaseConfig()
	config.ActionID = core.ActionID{SpellID: 53023}.WithTag(numTicksIdx)
	config.Flags = flags
	config.ManaCost = core.ManaCostOptions{
		BaseCost:   0.28,
		Multiplier: 1 - 0.05*float64(priest.Talents.FocusedMind),
	}
	config.Cast = core.CastConfig{
		DefaultCast: core.Cast{
			GCD: core.GCDDefault,
		},
	}
	config.Dot = core.DotConfig{
		Aura: core.Aura{
			Label: "MindSear-" + strconv.Itoa(int(numTicksIdx)),
		},
		NumberOfTicks:       numTicks,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				if aoeTarget != target {
					mindSearTickSpell.Cast(sim, aoeTarget)
					mindSearTickSpell.SpellMetrics[target.UnitIndex].Casts -= 1
				}
			}
		},
	}
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
		if result.Landed() {
			spell.Dot(target).Apply(sim)
			mindSearTickSpell.SpellMetrics[target.UnitIndex].Casts += 1
		}
	}
	config.ExpectedTickDamage = func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
		baseDamage := sim.Roll(212, 228) + miseryCoeff*spell.SpellPower()
		return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
	}
	return priest.GetOrRegisterSpell(config)
}
