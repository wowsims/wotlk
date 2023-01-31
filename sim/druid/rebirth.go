package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

// Right now, add the additional GCD + mana cost for shifting back to Moonkin form as a hack
// Consider adding moonkin shapeshift spell / form tracking to balance rotation instead
// Then we can properly incur Rebirth cost through additional Moonkin form spell cast
func (druid *Druid) registerRebirthSpell() {
	baseCost := 0.68
	if druid.InForm(Moonkin) {
		baseCost += 0.13
	}

	druid.Rebirth = druid.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 48477},
		Flags:    SpellFlagOmenTrigger,

		ManaCost: core.ManaCostOptions{
			BaseCost:   baseCost,
			Multiplier: 1 - 0.1*float64(druid.Talents.NaturalShapeshifter),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second*3 + time.Millisecond*500,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.RebirthUsed = true
		},
	})
}
