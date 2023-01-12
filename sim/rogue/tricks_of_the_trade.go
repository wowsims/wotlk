package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) registerTricksOfTheTradeSpell() {
	actionID := core.ActionID{SpellID: 57934, Tag: rogue.Index}
	hasShadowblades := rogue.HasSetBonus(ItemSetShadowblades, 2)
	energyMetrics := rogue.NewEnergyMetrics(actionID)
	hasGlyph := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfTricksOfTheTrade)

	if rogue.Options.TricksOfTheTradeTarget != nil {
		targetAgent := rogue.Env.Raid.GetPlayerFromRaidTarget(rogue.Options.TricksOfTheTradeTarget)
		if targetAgent != nil {
			target := targetAgent.GetCharacter()
			rogue.TricksOfTheTradeAura = core.TricksOfTheTradeAura(target, actionID.Tag, hasGlyph)
		}
	}
	if rogue.TricksOfTheTradeAura == nil {
		target := rogue.GetCharacter()
		rogue.TricksOfTheTradeAura = core.TricksOfTheTradeAura(target, actionID.Tag, hasGlyph)
		rogue.TricksOfTheTradeAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {}
		rogue.TricksOfTheTradeAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {}
	}

	rogue.TricksOfTheTrade = rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		EnergyCost: core.EnergyCostOptions{
			Cost: 15 - 5*float64(rogue.Talents.FilthyTricks),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
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
				rogue.AddEnergy(sim, spell.DefaultCast.Cost, energyMetrics)
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
