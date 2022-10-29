package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerUnholyFrenzyCD() {
	if !dk.Talents.Hysteria {
		return
	}

	actionID := core.ActionID{SpellID: 49016, Tag: dk.Index}

	unholyFrenzyTargetAgent := dk.Party.Raid.GetPlayerFromRaidTarget(dk.Inputs.UnholyFrenzyTarget)
	if unholyFrenzyTargetAgent == nil {
		return
	}
	unholyFrenzyTarget := unholyFrenzyTargetAgent.GetCharacter()
	unholyFrenzyAura := core.UnholyFrenzyAura(unholyFrenzyTarget, actionID.Tag)

	unholyFrenzySpell := dk.Character.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			unholyFrenzyAura.Activate(sim)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell:    unholyFrenzySpell,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeDPS,
	})
}
