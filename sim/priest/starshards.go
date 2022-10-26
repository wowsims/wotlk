package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (priest *Priest) registerStarshardsSpell() {
	actionID := core.ActionID{SpellID: 25446}
	priest.Starshards = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				priest.StarshardsDot.Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})

	target := priest.CurrentTarget
	priest.StarshardsDot = core.NewDot(core.Dot{
		Spell: priest.Starshards,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Starshards-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks: 5,
		TickLength:    time.Second * 3,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 785/5 + 0.167*dot.Spell.SpellPower()
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
		},
	})
}
