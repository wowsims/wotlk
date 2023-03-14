package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (mage *Mage) registerBlizzardSpell() {
	results := make([]*core.SpellResult, len(mage.Env.Encounter.TargetUnits))

	mage.Blizzard = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 42939},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagChanneled,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.74,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 8,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1 - (0.1/3)*float64(mage.Talents.FrostChanneling),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Blizzard",
			},
			NumberOfTicks:       8,
			TickLength:          time.Second * 1,
			AffectedByCastSpeed: true,
			OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
				target := mage.CurrentTarget
				dot.SnapshotBaseDamage = 426 + (4.0/3.5/8)*dot.Spell.SpellPower()
				dot.SnapshotBaseDamage *= sim.Encounter.AOECapMultiplier()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for i, aoeTarget := range sim.Encounter.TargetUnits {
					results[i] = dot.CalcSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
				}
				for i := range sim.Encounter.TargetUnits {
					dot.Spell.DealPeriodicDamage(sim, results[i])
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
