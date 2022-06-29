package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warrior *Warrior) registerRampageSpell() {
	if !warrior.Talents.Rampage {
		return
	}
	actionID := core.ActionID{SpellID: 30033}

	var bonusPerStack stats.Stats
	warrior.RampageAura = warrior.RegisterAura(core.Aura{
		Label:     "Rampage",
		ActionID:  actionID,
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			bonusPerStack = warrior.ApplyStatDependencies(stats.Stats{stats.AttackPower: 50})
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			warrior.AddStatsDynamic(sim, bonusPerStack.Multiply(float64(newStacks-oldStacks)))
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Rampage Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				warrior.rampageValidUntil = sim.CurrentTime + time.Second*5
			}
		},
	})

	warrior.Rampage = warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Rage,
		BaseCost:     20,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: 20,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.rampageValidUntil = 0
			warrior.RampageAura.Activate(sim)
			warrior.RampageAura.AddStack(sim)
		},
	})
}

func (warrior *Warrior) ShouldRampage(sim *core.Simulation) bool {
	return warrior.Rampage != nil &&
		sim.CurrentTime < warrior.rampageValidUntil &&
		warrior.CurrentRage() >= 20 &&
		(warrior.RampageAura.GetStacks() < 5 || warrior.RampageAura.RemainingDuration(sim) <= warrior.RampageCDThreshold)
}
