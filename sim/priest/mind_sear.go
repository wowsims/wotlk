package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// TODO see Mind Flay: Mind Sear (53023) now "periodically triggers" Mind Sear (53022).
// Since Mind Flay no longer is a binary spell, Mind Sear likely isn't, either.

func (priest *Priest) newMindSearSpell(numTicks int32) *core.Spell {
	channelTime := time.Second * time.Duration(numTicks)
	miseryCoeff := 0.2861 * (1 + 0.05*float64(priest.Talents.Misery))
	hasGlyphOfShadow := priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow))

	return priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53023, Tag: numTicks},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagChanneled,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.28,
			Multiplier: 1 - 0.05*float64(priest.Talents.FocusedMind),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
		},

		BonusHitRating:  float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating: float64(priest.Talents.MindMelt) * 2 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 +
			0.02*float64(priest.Talents.Darkness) +
			0.01*float64(priest.Talents.TwinDisciplines),
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),
		CritMultiplier:   priest.SpellCritMultiplier(1, 1),
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "MindSear-" + strconv.Itoa(int(numTicks)),
			},
			NumberOfTicks:       numTicks,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = sim.Roll(212, 228) + miseryCoeff*dot.Spell.SpellPower()
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeMagicHitAndSnapshotCrit)
				dot.Spell.DealDamage(sim, result)

				if result.Landed() {
					priest.AddShadowWeavingStack(sim)
				}
				if result.DidCrit() && hasGlyphOfShadow {
					priest.ShadowyInsightAura.Activate(sim)
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				if aoeTarget != sim.Encounter.TargetUnits[0] {

					result := spell.CalcOutcome(sim, aoeTarget, spell.OutcomeMagicHit)
					if result.Landed() {
						spell.SpellMetrics[aoeTarget.UnitIndex].Hits--
						spell.Dot(aoeTarget).Apply(sim)
					}
					spell.DealOutcome(sim, result)
				}
			}
		},
		ExpectedDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := sim.Roll(212, 228) + miseryCoeff*spell.SpellPower()
			baseDamage *= float64(numTicks)
			return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
		},
	})
}
