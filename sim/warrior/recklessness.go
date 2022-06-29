package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warrior *Warrior) RegisterRecklessnessCD() {
	actionID := core.ActionID{SpellID: 1719}
	reckAura := warrior.RegisterAura(core.Aura{
		Label:    "Recklessness",
		ActionID: actionID,
		Duration: time.Second*15 + time.Second*2*time.Duration(warrior.Talents.ImprovedDisciplines),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.MeleeCrit, 100*core.MeleeCritRatingPerCritChance)
			warrior.PseudoStats.DamageTakenMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.MeleeCrit, -100*core.MeleeCritRatingPerCritChance)
			warrior.PseudoStats.DamageTakenMultiplier /= 1.2
		},
	})

	cooldownDur := time.Minute * 30
	if warrior.Talents.ImprovedDisciplines == 1 {
		cooldownDur -= time.Minute * 4
	} else if warrior.Talents.ImprovedDisciplines == 2 {
		cooldownDur -= time.Minute * 7
	} else if warrior.Talents.ImprovedDisciplines == 3 {
		cooldownDur -= time.Minute * 10
	}
	reckSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownDur,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			reckAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: reckSpell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return warrior.StanceMatches(BerserkerStance)
		},
	})
}
