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
		targetUnit := rogue.GetUnit(rogue.Options.TricksOfTheTradeTarget)
		if targetUnit != nil {
			rogue.TricksOfTheTradeAura = core.TricksOfTheTradeAura(targetUnit, rogue.Index, hasGlyph)
		}
	}
	if rogue.TricksOfTheTradeAura == nil {
		target := &rogue.GetCharacter().Unit
		rogue.TricksOfTheTradeAura = core.TricksOfTheTradeAura(target, rogue.Index, hasGlyph)
		rogue.TricksOfTheTradeAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {}
		rogue.TricksOfTheTradeAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {}
	}

	rogue.TricksOfTheTrade = rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

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
		tricksSpell := rogue.TricksOfTheTrade
		rogue.AddMajorCooldown(core.MajorCooldown{
			Spell:    tricksSpell,
			Priority: core.CooldownPriorityDrums,
			Type:     core.CooldownTypeDPS,
			ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
				if hasShadowblades {
					return rogue.CurrentEnergy() <= rogue.maxEnergy-15-rogue.EnergyTickMultiplier*10
				} else if sim.CurrentTime < (tricksSpell.CD.Duration) {
					// This assumes you precast a Tricks before combat, and activated it (and the cooldown) at 0.00 on the sim.
					// This was put intentionally below the hasShadowblades check, because once you have that set a precast is no longer optimal.
					return false
				} else {
					return rogue.CurrentEnergy() >= rogue.TricksOfTheTrade.DefaultCast.Cost
				}
			},
		})
	}
}
