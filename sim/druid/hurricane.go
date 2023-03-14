package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (druid *Druid) registerHurricaneSpell() {
	druid.Hurricane = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48467},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagChanneled | SpellFlagOmenTrigger,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.81,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 10,
			},
		},

		DamageMultiplier: 1 +
			0.15*float64(druid.Talents.GaleWinds) +
			0.01*float64(druid.Talents.Genesis),
		ThreatMultiplier: 1,
		CritMultiplier:   1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Hurricane",
			},
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
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeMagicHitAndSnapshotCrit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
