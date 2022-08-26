package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Right now, add the additional GCD + mana cost for shifting back to Moonkin form as a hack
// Consider adding moonkin shapeshift spell / form tracking to balance rotation instead
// Then we can properly incur Rebirth cost through additional Moonkin form spell cast
func (druid *Druid) registerRebirthSpell() {
	baseCost := 1611 + (521.4 * (1 - (float64(druid.Talents.NaturalShapeshifter) * 0.1)))

	druid.Rebirth = druid.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 48477},

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
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
