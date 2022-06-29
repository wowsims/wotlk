package druid

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Right now, add the additional GCD + mana cost for shifting back to Moonkin form as a hack
// Consider adding moonkin shapeshift spell / form tracking to balance rotation instead
// Then we can properly incur Rebirth cost through additional Moonkin form spell cast
func (druid *Druid) registerRebirthSpell() {
	baseCost := 1611 + (521.4 * (1 - (float64(druid.Talents.NaturalShapeshifter) * 0.1)))

	druid.Rebirth = druid.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 26994},

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Second*3 + time.Millisecond*500,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.RebirthUsed = true
		},
	})
}

func (druid *Druid) TryRebirth(sim *core.Simulation) bool {
	if druid.RebirthUsed {
		return false
	}

	if success := druid.Rebirth.Cast(sim, nil); !success {
		druid.WaitForMana(sim, druid.Rebirth.CurCast.Cost)
	}
	return true
}
