package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO see Mind Flay: Mind Sear (53023) now "periodically triggers" Mind Sear (53022).
//
//	Since Mind Flay no longer is a binary spell, Mind Sear likely isn't, either.
func (priest *Priest) MindSearActionID(numTicks int32) core.ActionID {
	return core.ActionID{SpellID: 53023, Tag: numTicks}
}

func (priest *Priest) newMindSearSpell(numTicks int32) *core.Spell {
	baseCost := priest.BaseMana * 0.28
	channelTime := time.Second * time.Duration(numTicks)

	return priest.RegisterSpell(core.SpellConfig{
		ActionID:     priest.MindSearActionID(numTicks),
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        core.SpellFlagChanneled,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:        baseCost * (1 - 0.05*float64(priest.Talents.FocusedMind)),
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

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.Targets {
				result := spell.CalcAndDealOutcome(sim, &aoeTarget.Unit, spell.OutcomeMagicHit)
				if result.Landed() {
					priest.MindSearDot[numTicks].Apply(sim)
				}
			}
		},
	})
}

func (priest *Priest) newMindSearDot(numTicks int32) *core.Dot {
	target := priest.CurrentTarget

	miseryCoeff := 0.2861 * (1 + 0.05*float64(priest.Talents.Misery))
	hasGlyphOfShadow := priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow))

	return core.NewDot(core.Dot{
		Spell: priest.MindSear[numTicks],
		Aura: target.RegisterAura(core.Aura{
			Label:    "MindSear-" + strconv.Itoa(int(numTicks)) + "-" + strconv.Itoa(int(priest.Index)),
			ActionID: priest.MindSearActionID(numTicks),
		}),

		NumberOfTicks:       numTicks,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = sim.Roll(212, 228) + miseryCoeff*dot.Spell.SpellPower()
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			for _, aoeTarget := range sim.Encounter.Targets {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, &aoeTarget.Unit, dot.Spell.OutcomeMagicHit)
				if result.Landed() {
					priest.AddShadowWeavingStack(sim)
				}
				if result.DidCrit() && hasGlyphOfShadow {
					priest.ShadowyInsightAura.Activate(sim)
				}
			}
		},
	})
}
