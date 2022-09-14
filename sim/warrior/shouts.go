package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const ShoutExpirationThreshold = time.Second * 3

func (warrior *Warrior) makeShoutSpellHelper(actionID core.ActionID, extraDuration time.Duration) *core.Spell {
	cost := 10.0

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			// Actual shout effects are handled in core/buffs.go
			warrior.shoutExpiresAt = sim.CurrentTime + warrior.shoutDuration + extraDuration
		},
	})
}

func (warrior *Warrior) makeShoutSpell() *core.Spell {
	if warrior.ShoutType == proto.WarriorShout_WarriorShoutBattle {
		extraDur := core.TernaryDuration(warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfBattle), 2*time.Minute, 0)
		return warrior.makeShoutSpellHelper(core.ActionID{SpellID: 47436}, extraDur)
	} else if warrior.ShoutType == proto.WarriorShout_WarriorShoutCommanding {
		extraDur := core.TernaryDuration(warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfCommand), 2*time.Minute, 0)
		return warrior.makeShoutSpellHelper(core.ActionID{SpellID: 47440}, extraDur)
	} else {
		return nil
	}
}

func (warrior *Warrior) ShouldShout(sim *core.Simulation) bool {
	return warrior.Shout != nil && warrior.CurrentRage() >= warrior.Shout.DefaultCast.Cost && sim.CurrentTime+ShoutExpirationThreshold > warrior.shoutExpiresAt
}
