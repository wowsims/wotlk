package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (hunter *Hunter) registerVolleySpell() {
	actionID := core.ActionID{SpellID: 58434}

	volleyDot := core.NewDot(core.Dot{
		Aura: hunter.RegisterAura(core.Aura{
			Label:    "Volley",
			ActionID: actionID,
		}),
		NumberOfTicks:       6,
		TickLength:          time.Second * 1,
		AffectedByCastSpeed: true,

		OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
			target := hunter.CurrentTarget
			dot.SnapshotBaseDamage = 353 + 0.0837*dot.Spell.RangedAttackPower(target)
			dot.SnapshotBaseDamage *= sim.Encounter.AOECapMultiplier()

			attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
			dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(target, attackTable)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			for _, aoeTarget := range sim.Encounter.Targets {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, &aoeTarget.Unit, dot.OutcomeRangedHitAndCritSnapshot)
			}
		},
	})

	hunter.Volley = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagChanneled,

		Cost: core.NewManaCost(core.ManaCostOptions{
			BaseCost:   0.17,
			Multiplier: core.TernaryFloat64(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfVolley), 0.8, 1),
		}),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 6,
			},
		},

		DamageMultiplier: 1 *
			(1 + 0.04*float64(hunter.Talents.Barrage)),
		CritMultiplier:   hunter.critMultiplier(true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			channelDoneAt := sim.CurrentTime + hunter.Volley.CurCast.ChannelTime
			hunter.mayMoveAt = channelDoneAt
			hunter.AutoAttacks.DelayRangedUntil(sim, channelDoneAt+time.Millisecond*500)
			volleyDot.Apply(sim)
		},
	})
	volleyDot.Spell = hunter.Volley
}
