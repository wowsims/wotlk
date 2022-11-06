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
	refundAmount := 0.4 * float64(rogue.Talents.QuickRecovery)
	numTicks := 6
	baseCost := rogue.costModifier(50 - 10*float64(rogue.Talents.DirtyDeeds))
	totalDamageMod := 1.0
	if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfGarrote) {
		numTicks = 5
		// 20% **total** damage increase with one fewer tick
		totalDamageMod = 1.2 * 6.0 / 5.0
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
			0.15*float64(rogue.Talents.BloodSpatter) +
			0.02*float64(rogue.Talents.FindWeakness),
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				rogue.garroteDot.Spell = spell
				rogue.garroteDot.NumberOfTicks = numTicks
				rogue.garroteDot.RecomputeAuraDuration()
				rogue.garroteDot.Apply(sim)
				rogue.ApplyFinisher(sim, spell)
			} else {
				if refundAmount > 0 {
					rogue.AddEnergy(sim, spell.CurCast.Cost*refundAmount, rogue.QuickRecoveryMetrics)
				}
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
			dot.SnapshotBaseDamage = 119 + dot.Spell.MeleeAttackPower()*0.07*totalDamageMod
			attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
			dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(target, attackTable)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
		},
	})
}
