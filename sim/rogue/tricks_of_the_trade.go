package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) registerTricksOfTheTradeSpell() {
	actionID := core.ActionID{SpellID: 57934, Tag: rogue.Index}
	hasShadowblades := rogue.HasSetBonus(ItemSetShadowblades, 2)
	energyMetrics := rogue.NewEnergyMetrics(actionID)
	energyCost := 15 - 5*float64(rogue.Talents.FilthyTricks)
	hasGlyph := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfTricksOfTheTrade)

	if rogue.Options.TricksOfTheTradeTarget != nil {
		target := rogue.Env.Raid.GetPlayerFromRaidTarget(rogue.Options.TricksOfTheTradeTarget).GetCharacter()
		rogue.TricksOfTheTradeAura = core.TricksOfTheTradeAura(target, actionID.Tag, hasGlyph)
	} else {
		target := rogue.GetCharacter()
		rogue.TricksOfTheTradeAura = core.TricksOfTheTradeAura(target, actionID.Tag, hasGlyph)
		rogue.TricksOfTheTradeAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {}
		rogue.TricksOfTheTradeAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {}
	}

	rogue.TricksOfTheTrade = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		BaseCost:     energyCost,
		ResourceType: stats.Energy,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: energyCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * time.Duration(30-5*rogue.Talents.FilthyTricks),
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if hasShadowblades {
					cast.Cost = 0
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			rogue.TricksOfTheTradeAura.Activate(sim)
			if hasShadowblades {
				rogue.AddEnergy(sim, energyCost, energyMetrics)
			}
		},
	})

	if rogue.Rotation.TricksOfTheTradeFrequency != proto.Rogue_Rotation_Never {
		// TODO: Support Rogue_Rotation_Once
		rogue.AddMajorCooldown(core.MajorCooldown{
			Spell:    rogue.TricksOfTheTrade,
			Priority: core.CooldownPriorityDrums,
			Type:     core.CooldownTypeDPS,
		})
	}
}
