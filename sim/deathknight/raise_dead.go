package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerRaiseDeadCD() {
	// If talented as permanent pet skip this spell
	if dk.Talents.MasterOfGhouls {
		return
	}

	raiseDeadAura := dk.RegisterAura(core.Aura{
		Label:    "Raise Dead",
		ActionID: core.ActionID{SpellID: 46584},
		Duration: time.Minute * 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.Ghoul.Pet.Enable(sim, dk.Ghoul)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.Ghoul.Pet.Disable(sim)
		},
	})

	dk.RaiseDead = dk.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 46584},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute*3 - time.Second*45*time.Duration(dk.Talents.NightOfTheDead),
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			raiseDeadAura.Activate(sim)
		},
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell: dk.RaiseDead,
			Type:  core.CooldownTypeSurvival,
			CanActivate: func(sim *core.Simulation, character *core.Character) bool {
				return dk.RaiseDead.CanCast(sim, nil) && dk.CurrentHealthPercent() < 0.5 && sim.GetRemainingDuration() > 5*time.Second
			},
		})
	}
}
