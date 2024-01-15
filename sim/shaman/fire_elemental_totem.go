package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

const fireTotemDuration time.Duration = time.Second * 120

func (shaman *Shaman) registerFireElementalTotem() {
	if !shaman.Totems.UseFireElemental {
		return
	}

	actionID := core.ActionID{SpellID: 2894}

	fireElementalAura := shaman.RegisterAura(core.Aura{
		Label:    "Fire Elemental Totem",
		ActionID: actionID,
		Duration: fireTotemDuration,
	})

	shaman.FireElementalTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.23,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * time.Duration(core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFireElementalTotem), 5, 10)),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			// TODO: ToW needs a unique buff/debuff aura for each raidmember/target.
			//  Otherwise we will be possibly disabling another ele shaman's ToW debuff/buff.
			if shaman.Totems.Fire != proto.FireTotem_NoFireTotem {
				shaman.TotemExpirations[FireTotem] = sim.CurrentTime + fireTotemDuration
			}

			shaman.MagmaTotem.AOEDot().Cancel(sim)
			searingTotemDot := shaman.SearingTotem.Dot(shaman.CurrentTarget)
			if searingTotemDot != nil {
				searingTotemDot.Cancel(sim)
			}

			shaman.FireElemental.EnableWithTimeout(sim, shaman.FireElemental, fireTotemDuration)

			// Add a dummy aura to show in metrics
			fireElementalAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: shaman.FireElementalTotem,
		Type:  core.CooldownTypeDPS,
	})
}
