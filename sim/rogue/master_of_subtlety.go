package rogue

import (
	"time"
	"math"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wotsims/wotlk/sim/core/stats"
)

// Each rank is a different ID. 31223 is 3/3
switch rogue.Talents.MasterOfSubtlety {
case 2:
	id = 31222
case 3:
	id = 31223
default:
	id = 31221
}
var MasterOfSubtletyID = core.ActionID{SpellID: id}

func (rogue *Rogue) registerMasterOfSubtletyCD() {
	if !rogue.Talents.MasterOfSubtlety {
		return
	}
	rogue.MasterOfSubtletyAura = rogue.RegisterAura(core.Aura{
		Label:		"Master of Subtlety",
		ActionID:	MasterOfSubtletyID,
		Duration:	time.Second * 6,
		OnGain: func(aura *core.aura, sim *core.Simulation) {
			switch rogue.Talents.MasterOfSubtlety {
			case 2:
				percent = 0.07
			case 3:
				percent = 0.1
			default:
				percent = 0.04
			}
			rogue.MulitplyStat(stats.AttackPower, 1.0+percent)
		},
		OnExpire: func(aura *core.aura, sim *core.Simulation) {
			switch rogue.Talents.MasterOfSubtlety {
			case 2:
				percent = 0.07
			case 3:
				percent = 0.1
			default:
				percent = 0.04
			}
			rogue.MultiplyStat(stats.AttackPower, 1.0 /(1.0 + percent))
		},
	})
	masterOfSubtletySpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID: MasterOfSubtletyID,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:		rogue.NewTimer(),
				Duration:	time.Minute * (3 - (time.Second * 30 * rogue.Talents.Elusiveness)),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.MasterOfSubtletyAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:	masterOfSubtletySpell,
		Type:	core.CooldownTypeDPS,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return rogue.CurrentEnergy() > 90
		},
	})
}