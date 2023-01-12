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
		panic("No unholy frenzy target")
		return
	}
	unholyFrenzyTarget := unholyFrenzyTargetAgent.GetCharacter()
	unholyFrenzyAura := core.UnholyFrenzyAura(unholyFrenzyTarget, actionID.Tag)

	dk.UnholyFrenzy = dk.Character.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
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
		Spell:    dk.UnholyFrenzy,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeDPS,
	})
}
