package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (shaman *Shaman) registerFeralSpirit() {
	if !shaman.Talents.FeralSpirit {
		return
	}

	spiritWolvesActiveAura := shaman.RegisterAura(core.Aura{
		Label:    "Feral Spirit",
		ActionID: core.ActionID{SpellID: 51533},
		Duration: time.Second * 45,
	})

	shaman.FeralSpirit = shaman.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 51533},

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.12,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
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

			// https://github.com/JamminL/wotlk-classic-bugs/issues/280
			// instant casts (e.g. shocks) usually don't reset a shaman's swing timer
			shaman.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell:    shaman.FeralSpirit,
		Priority: core.CooldownPriorityDrums + 1000, // Always prefer to use wolves before bloodlust/drums so wolves gain haste buff
		Type:     core.CooldownTypeDPS,
	})
}
