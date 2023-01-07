package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerBlizzardSpell() {
	actionID := core.ActionID{SpellID: 42939}
	baseCost := .74 * mage.BaseMana

	results := make([]*core.SpellResult, len(mage.Env.Encounter.Targets))
	blizzardDot := core.NewDot(core.Dot{
		Aura: mage.RegisterAura(core.Aura{
			Label:    "Blizzard",
			ActionID: actionID,
		}),
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
			for i, aoeTarget := range sim.Encounter.Targets {
				results[i] = dot.CalcSnapshotDamage(sim, &aoeTarget.Unit, dot.OutcomeTick)
			}
			for i := range sim.Encounter.Targets {
				dot.Spell.DealPeriodicDamage(sim, results[i])
			}
		},
	})

	mage.Blizzard = mage.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | core.SpellFlagChanneled,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,

				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 8,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1 - (0.1/3)*float64(mage.Talents.FrostChanneling),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			blizzardDot.Apply(sim)
		},
	})
	blizzardDot.Spell = mage.Blizzard
}
