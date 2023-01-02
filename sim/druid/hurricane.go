package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerHurricaneSpell() {
	actionID := core.ActionID{SpellID: 48467}
	baseCost := 0.81 * druid.BaseMana

	hurricaneDot := core.NewDot(core.Dot{
		Aura: druid.RegisterAura(core.Aura{
			Label:    "Hurricane",
			ActionID: actionID,
		}),
		NumberOfTicks:       10,
		TickLength:          time.Second * 1,
		AffectedByCastSpeed: true,

		OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
			target := druid.CurrentTarget
			dot.SnapshotBaseDamage = 451 + 0.129*dot.Spell.SpellPower()
			dot.SnapshotBaseDamage *= sim.Encounter.AOECapMultiplier()
			dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
			attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			for _, aoeTarget := range sim.Encounter.Targets {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, &aoeTarget.Unit, dot.OutcomeMagicHitAndSnapshotCrit)
			}
		},
	})

	druid.Hurricane = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagChanneled,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:        baseCost,
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 10,
			},
		},

		DamageMultiplier: 1 +
			0.15*float64(druid.Talents.GaleWinds) +
			0.01*float64(druid.Talents.Genesis),
		ThreatMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if druid.ClearcastingAura != nil {
				druid.ClearcastingAura.Deactivate(sim)
			}

			hurricaneDot.Apply(sim)
		},
	})
	hurricaneDot.Spell = druid.Hurricane
}
