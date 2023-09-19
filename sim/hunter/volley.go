package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (hunter *Hunter) registerVolleySpell() {
	hunter.Volley = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58434},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagChanneled | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.17,
			Multiplier: core.TernaryFloat64(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfVolley), 0.8, 1),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 6,
			},
		},

		DamageMultiplier: 1 *
			(1 + 0.04*float64(hunter.Talents.Barrage)),
		CritMultiplier:   hunter.critMultiplier(true, false, false),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Volley",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					hunter.AutoAttacks.DelayRangedUntil(sim, sim.CurrentTime+time.Millisecond*500)
				},
			},
			NumberOfTicks:       6,
			TickLength:          time.Second * 1,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
				target := hunter.CurrentTarget
				dot.SnapshotBaseDamage = 353 + 0.0837*dot.Spell.RangedAttackPower(target)
				dot.SnapshotBaseDamage *= sim.Encounter.AOECapMultiplier()

				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeRangedHitAndCritSnapshot)
				}
				hunter.AutoAttacks.DelayRangedUntil(sim, sim.CurrentTime+dot.TickPeriod()+time.Millisecond*500)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
			hunter.AutoAttacks.CancelAutoSwing(sim)
		},
	})
}
