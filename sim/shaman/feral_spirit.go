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
			OnCastComplete: func(sim *core.Simulation, spell *core.Spell) {
				attackSpeed := shaman.AutoAttacks.MainhandSwingSpeed()

				if shaman.AutoAttacks.IsDualWielding {
					attackSpeed = core.MinDuration(attackSpeed, shaman.AutoAttacks.OffhandSwingSpeed())
				}

				shaman.AutoAttacks.DelayMeleeUntil(sim, sim.CurrentTime+attackSpeed)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shaman.SpiritWolves.EnableWithTimeout(sim)
			shaman.SpiritWolves.CancelGCDTimer(sim)

			// Add a dummy aura to show in metrics
			spiritWolvesActiveAura.Activate(sim)
		},
	})

	//TODO: reset swing timer on cast, unless it ends up being fixed

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell:    shaman.FeralSpirit,
		Priority: core.CooldownPriorityDrums + 1, // Always prefer to use wolves before bloodlust/drums so wolves gain haste buff
		Type:     core.CooldownTypeDPS,
	})
}
