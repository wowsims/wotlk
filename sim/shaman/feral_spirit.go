package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (shaman *Shaman) registerFeralSpirit() {
	if !shaman.Talents.FeralSpirit {
		return
	}

	manaCost := 0.12 * shaman.BaseMana

	spiritWolvesActiveAura := shaman.RegisterAura(core.Aura{
		Label:    "Feral Spirit",
		ActionID: core.ActionID{SpellID: 51533},
		Duration: time.Second * 45,
	})

	shaman.FeralSpirit = shaman.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 51533},

		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: manaCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shaman.SpiritWolves.EnableWithTimeout(sim)
			shaman.SpiritWolves.CancelGCDTimer(sim)

			// Add a dummy aura to show in metrics
			spiritWolvesActiveAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell:    shaman.FeralSpirit,
		Priority: core.CooldownPriorityDrums + 1, // Always prefer to use wolves before bloodlust/drums so wolves gain haste buff
		Type:     core.CooldownTypeDPS,
		ActivationFactory: func(sim *core.Simulation) core.CooldownActivation {
			return func(sim *core.Simulation, unit *core.Character) {
				if shaman.FeralSpirit.Cast(sim, unit.CurrentTarget) {
					// https://github.com/JamminL/wotlk-classic-bugs/issues/280
					// instant casts (e.g. shocks) usually don't reset a shaman's swing timer
					shaman.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
				}
			}
		},
	})
}
