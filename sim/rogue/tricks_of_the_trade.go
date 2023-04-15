package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) registerTricksOfTheTradeSpell() {
	actionID := core.ActionID{SpellID: 57934}
	energyMetrics := rogue.NewEnergyMetrics(actionID)
	hasShadowblades := rogue.HasSetBonus(Tier10, 2)
	energyCost := 15 - 5*float64(rogue.Talents.FilthyTricks)

	hasGlyph := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfTricksOfTheTrade)
	if rogue.Options.TricksOfTheTradeTarget != nil {
		targetAgent := rogue.Env.Raid.GetPlayerFromRaidTarget(rogue.Options.TricksOfTheTradeTarget)
		if targetAgent != nil {
			target := targetAgent.GetCharacter()
			rogue.TricksOfTheTradeAura = core.TricksOfTheTradeAura(target, rogue.Index, hasGlyph)
		}
	}
	if rogue.TricksOfTheTradeAura == nil {
		target := rogue.GetCharacter()
		rogue.TricksOfTheTradeAura = core.TricksOfTheTradeAura(target, rogue.Index, hasGlyph)
		rogue.TricksOfTheTradeAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {}
		rogue.TricksOfTheTradeAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {}
	}

	rogue.TricksOfTheTrade = rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		EnergyCost: core.EnergyCostOptions{
			Cost: core.TernaryFloat64(hasShadowblades, 0, energyCost),
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
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			rogue.TricksOfTheTradeAura.Activate(sim)
			if hasShadowblades {
				rogue.AddEnergy(sim, 15, energyMetrics)
			}
		},
	})

	if rogue.Rotation.TricksOfTheTradeFrequency != proto.Rogue_Rotation_Never {
		// TODO: Support Rogue_Rotation_Once
		rogue.AddMajorCooldown(core.MajorCooldown{
			Spell:    rogue.TricksOfTheTrade,
			Priority: core.CooldownPriorityDrums,
			Type:     core.CooldownTypeDPS,
			ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
				if hasShadowblades {
					return rogue.CurrentEnergy() <= rogue.maxEnergy-15-rogue.EnergyTickMultiplier*10
				} else {
					return rogue.CurrentEnergy() >= rogue.TricksOfTheTrade.DefaultCast.Cost
				}
			},
		})
	}
}
