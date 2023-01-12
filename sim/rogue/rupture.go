package rogue

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

const RuptureEnergyCost = 25.0
const RuptureSpellID = 48672

func (rogue *Rogue) makeRupture(comboPoints int32) *core.Spell {
	numTicks := comboPoints + 3 + core.TernaryInt32(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfRupture), 2, 0)

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: RuptureSpellID, Tag: comboPoints},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | rogue.finisherFlags(),

		EnergyCost: core.EnergyCostOptions{
			Cost:          RuptureEnergyCost,
			Refund:        0.4 * float64(rogue.Talents.QuickRecovery),
			RefundMetrics: rogue.QuickRecoveryMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 +
			0.15*float64(rogue.Talents.BloodSpatter) +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.HasSetBonus(ItemSetBonescythe, 2), 0.1, 0) +
			core.TernaryFloat64(rogue.HasSetBonus(ItemSetTerrorblade, 4), 0.2, 0) +
			0.1*float64(rogue.Talents.SerratedBlades),
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				rogue.ruptureDot.Spell = spell
				rogue.ruptureDot.NumberOfTicks = numTicks
				rogue.ruptureDot.RecomputeAuraDuration()
				rogue.ruptureDot.Apply(sim)
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}

func (rogue *Rogue) RuptureDuration(comboPoints int32) time.Duration {
	return time.Second*6 +
		time.Second*2*time.Duration(comboPoints) +
		core.TernaryDuration(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfRupture), time.Second*4, 0)
}

func (rogue *Rogue) registerRupture() {
	rogue.Rupture = [6]*core.Spell{
		rogue.makeRupture(0), // Just for metrics
		rogue.makeRupture(1),
		rogue.makeRupture(2),
		rogue.makeRupture(3),
		rogue.makeRupture(4),
		rogue.makeRupture(5),
	}

	rogue.ruptureDot = core.NewDot(core.Dot{
		Spell: rogue.Rupture[0],
		Aura: rogue.CurrentTarget.RegisterAura(core.Aura{
			Label:    "Rupture-" + strconv.Itoa(int(rogue.Index)),
			Tag:      RogueBleedTag,
			ActionID: rogue.Rupture[0].ActionID,
		}),
		NumberOfTicks: 0, // Set dynamically
		TickLength:    time.Second * 2,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			comboPoints := rogue.ComboPoints()
			dot.SnapshotBaseDamage = 127 +
				18*float64(comboPoints) +
				[]float64{0, 0.06 / 4, 0.12 / 5, 0.18 / 6, 0.24 / 7, 0.30 / 8}[comboPoints]*dot.Spell.MeleeAttackPower()

			attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
			dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(target, attackTable)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
		},
	})
}
