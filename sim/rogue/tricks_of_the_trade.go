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

	tricksOfTheTradeThreatTransferAura := rogue.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 59628},
		Label:    "TricksOfTheTradeThreat",
		Duration: 6 * time.Second,
	})

	tricksOfTheTradeApplicationAura := rogue.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 57934},
		Label:    "TricksOfTheTradeApplication",
		Duration: 30 * time.Second,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				tricksOfTheTradeThreatTransferAura.Activate(sim)
				aura.Deactivate(sim)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.TricksOfTheTrade.CD.Set(core.NeverExpires)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.TricksOfTheTrade.CD.Set(sim.CurrentTime + time.Second*time.Duration(30-5*rogue.Talents.FilthyTricks))
		},
	})

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
				Duration: time.Second * time.Duration(30-5*rogue.Talents.FilthyTricks), // CD is handled by application aura
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			tricksOfTheTradeApplicationAura.Activate(sim)
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
			Priority: core.CooldownPriorityBloodlust,
			Type:     core.CooldownTypeDPS,
			ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
				if hasShadowblades {
					return rogue.CurrentEnergy() <= rogue.maxEnergy-15-rogue.EnergyTickMultiplier*10
				} else {
					return true
				}
			},
		})
	}
}
