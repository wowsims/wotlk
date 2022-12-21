package rogue

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const GarroteSpellID = 48676

func (rogue *Rogue) registerGarrote() {
	refundAmount := 0.8
	baseCost := rogue.costModifier(50 - 10*float64(rogue.Talents.DirtyDeeds))

	numTicks := int32(6)
	var glyphMultiplier float64
	if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfGarrote) {
		numTicks = 5
		glyphMultiplier = 0.44 // cp. https://www.wowhead.com/wotlk/spell=56812/glyph-of-garrote
	}
	rogue.Garrote = rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: GarroteSpellID},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | SpellFlagBuilder,
		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 +
			glyphMultiplier +
			0.15*float64(rogue.Talents.BloodSpatter) +
			0.10*float64(rogue.Talents.Opportunity) +
			0.02*float64(rogue.Talents.FindWeakness),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// #FIXME Can be dodged by boss
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit)
			if result.Landed() {
				comboPoints = 1 + (1 * 1/3 * float64(rogue.Talents.Initiative))
				rogue.AddComboPoints(sim, comboPoints, spell.ComboPointMetrics())
				rogue.garroteDot.Apply(sim)
			} else {
				rogue.AddEnergy(sim, spell.CurCast.Cost*refundAmount, rogue.EnergyRefundMetrics)
			}
			spell.DealOutcome(sim, result)
		},
	})

	rogue.garroteDot = core.NewDot(core.Dot{
		Spell: rogue.Garrote,
		Aura: rogue.CurrentTarget.RegisterAura(core.Aura{
			Label:    "Garrote-" + strconv.Itoa(int(rogue.Index)),
			Tag:      RogueBleedTag,
			ActionID: rogue.Garrote.ActionID,
		}),
		NumberOfTicks: numTicks,
		TickLength:    time.Second * 3,
		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 119 + dot.Spell.MeleeAttackPower()*0.07
			attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
		},
	})
}
