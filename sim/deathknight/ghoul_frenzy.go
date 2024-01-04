package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerGhoulFrenzySpell() {
	if !dk.Talents.GhoulFrenzy {
		return
	}

	dk.GhoulFrenzy = dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 63560},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagAPL,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
			RunicPowerGain: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return dk.Ghoul.IsEnabled()
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Hot: core.DotConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label: "Ghoul Frenzy Hot",
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 6,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 0.06 * dk.Ghoul.MaxHealth()
				dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotHealing(sim, &dk.Ghoul.Unit, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.GhoulFrenzyAura.Activate(sim)
			dk.Ghoul.GhoulFrenzyAura.Activate(sim)
			spell.SpellMetrics[target.UnitIndex].Hits++
			spell.SelfHot().Apply(sim)
		},
	})

	dk.GhoulFrenzyAura = dk.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 63560},
		Label:    "Ghoul Frenzy",
		Duration: time.Second * 30.0,
	})

	dk.Ghoul.GhoulFrenzyAura = dk.Ghoul.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 63560},
		Label:    "Ghoul Frenzy",
		Duration: time.Second * 30.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.Ghoul.MultiplyMeleeSpeed(sim, 1.25)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.Ghoul.MultiplyMeleeSpeed(sim, 1/1.25)
		},
	})
}
