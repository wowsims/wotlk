package warrior

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

const ShoutExpirationThreshold = time.Second * 3

func (warrior *Warrior) makeShoutSpellHelper(actionID core.ActionID, allyAuras core.AuraArray) *core.Spell {
	return warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | core.SpellFlagHelpful,

		RageCost: core.RageCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, aura := range allyAuras {
				if aura != nil {
					aura.Activate(sim)
				}
			}
		},

		RelatedAuras: []core.AuraArray{allyAuras},
	})
}

func (warrior *Warrior) makeShoutSpell() *core.Spell {
	battleShout := warrior.makeShoutSpellHelper(core.ActionID{SpellID: 47436}, warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		return core.BattleShoutAura(unit, warrior.Talents.CommandingPresence, warrior.Talents.BoomingVoice, warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfBattle))
	}))

	commandingShout := warrior.makeShoutSpellHelper(core.ActionID{SpellID: 47440}, warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		return core.CommandingShoutAura(unit, warrior.Talents.CommandingPresence, warrior.Talents.BoomingVoice, warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfCommand))
	}))

	if warrior.ShoutType == proto.WarriorShout_WarriorShoutBattle {
		return battleShout
	} else if warrior.ShoutType == proto.WarriorShout_WarriorShoutCommanding {
		return commandingShout
	} else {
		return nil
	}
}

func (warrior *Warrior) ShouldShout(sim *core.Simulation) bool {
	return warrior.Shout != nil && warrior.CurrentRage() >= warrior.Shout.DefaultCast.Cost && warrior.Shout.ShouldRefreshExclusiveEffects(sim, &warrior.Unit, ShoutExpirationThreshold)
}
