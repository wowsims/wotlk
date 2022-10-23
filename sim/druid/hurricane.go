package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerHurricaneSpell() {
	actionID := core.ActionID{SpellID: 48467}
	baseCost := 1905.0

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
			dot.SnapshotBaseDamage = 206 + 0.107*dot.Spell.SpellPower()
			//dot.SnapshotBaseDamage *= sim.Encounter.AOECapMultiplier()
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			for _, aoeTarget := range sim.Encounter.Targets {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, &aoeTarget.Unit, dot.OutcomeTick)
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
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 60,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDot(hurricaneDot),
	})
	hurricaneDot.Spell = druid.Hurricane
}

func (druid *Druid) ShouldCastHurricane(sim *core.Simulation, rotation proto.BalanceDruid_Rotation) bool {
	return len(druid.Env.Encounter.Targets) > 1 && druid.Hurricane.IsReady(sim)
}
