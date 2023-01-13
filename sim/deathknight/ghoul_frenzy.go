package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerGhoulFrenzySpell() {
	if !dk.Talents.GhoulFrenzy {
		return
	}

	gfHeal := dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 63560},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,

		Cast: core.CastConfig{},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SpellMetrics[target.UnitIndex].Hits++
		},
	}, nil)

	gfHealHot := core.NewDot(core.Dot{
		Spell: gfHeal.Spell,
		Aura: dk.Ghoul.RegisterAura(core.Aura{
			Label:    "Ghoul Frenzy Hot",
			ActionID: gfHeal.ActionID,
		}),
		NumberOfTicks: 5,
		TickLength:    time.Second * 6,
		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 0.06 * dk.Ghoul.MaxHealth()
			dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotHealing(sim, &dk.Ghoul.Unit, dot.OutcomeTick)
		},
	})

	baseCost := float64(core.NewRuneCost(10, 0, 0, 1, 0))
	dk.GhoulFrenzy = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 63560},
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.GetModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.GhoulFrenzyAura.Activate(sim)
			gfHealHot.Apply(sim)
			dk.Ghoul.GhoulFrenzyAura.Activate(sim)
		},
	}, func(sim *core.Simulation) bool {
		return dk.Talents.GhoulFrenzy && dk.Ghoul.IsEnabled() && dk.CastCostPossible(sim, 0.0, 0, 0, 1) && dk.GhoulFrenzy.IsReady(sim)
	})

	dk.GhoulFrenzyAura = dk.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 63560},
		Label:    "Ghoul Frenzy",
		Duration: time.Second * 30.0,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if dk.Inputs.PrecastGhoulFrenzy {
				dk.GhoulFrenzyAura.Activate(sim)
				dk.GhoulFrenzyAura.UpdateExpires(sim.CurrentTime + time.Second*20)
			}
		},
	})

	dk.Ghoul.GhoulFrenzyAura = dk.Ghoul.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 63560},
		Label:    "Ghoul Frenzy",
		Duration: time.Second * 30.0,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if dk.Inputs.PrecastGhoulFrenzy {
				dk.Ghoul.GhoulFrenzyAura.Activate(sim)
				dk.Ghoul.GhoulFrenzyAura.UpdateExpires(sim.CurrentTime + time.Second*20)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.Ghoul.MultiplyMeleeSpeed(sim, 1.25)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.Ghoul.MultiplyMeleeSpeed(sim, 1/1.25)
		},
	})
}
