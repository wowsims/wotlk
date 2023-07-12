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

	unholyFrenzyTarget := &dk.Character
	if !dk.IsUsingAPL {
		unholyFrenzyTarget = nil
	}

	unholyFrenzyTargetAgent := dk.Party.Raid.GetPlayerFromRaidTarget(dk.Inputs.UnholyFrenzyTarget)
	if unholyFrenzyTargetAgent != nil {
		unholyFrenzyTarget = unholyFrenzyTargetAgent.GetCharacter()
	}

	dk.UnholyFrenzyAura = core.UnholyFrenzyAura(unholyFrenzyTarget, actionID.Tag)

	dk.UnholyFrenzy = dk.Character.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			dk.UnholyFrenzyAura.Activate(sim)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell:    dk.UnholyFrenzy,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeDPS,
	})
}
