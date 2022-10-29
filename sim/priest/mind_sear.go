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
func (priest *Priest) MindSearActionID(numTicks int) core.ActionID {
	return core.ActionID{SpellID: 53023, Tag: int32(numTicks)}
}

func (priest *Priest) newMindSearSpell(numTicks int) *core.Spell {
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

		BonusHitRating:   float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  float64(priest.Talents.MindMelt) * 2 * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
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

func (priest *Priest) newMindSearDot(numTicks int) *core.Dot {
	target := priest.CurrentTarget

	miseryCoeff := 0.2861 * (1 + 0.05*float64(priest.Talents.Misery))
	hasGlyphOfShadow := priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow))

	normMod := (1 + float64(priest.Talents.Darkness)*0.02 + float64(priest.Talents.TwinDisciplines)*0.01) // initialize modifier

	return core.NewDot(core.Dot{
		Spell: priest.MindSear[numTicks],
		Aura: target.RegisterAura(core.Aura{
			Label:    "MindSear-" + strconv.Itoa(numTicks) + "-" + strconv.Itoa(int(priest.Index)),
			ActionID: priest.MindSearActionID(numTicks),
		}),

		NumberOfTicks:       numTicks,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dmg := sim.Roll(212, 228) + miseryCoeff*dot.Spell.SpellPower()

			dot.SnapshotBaseDamage = dmg * normMod
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
