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

	var targetUnit *core.Unit
	hasGlyph := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfTricksOfTheTrade)
	if rogue.Options.TricksOfTheTradeTarget != nil {
		targetUnit = rogue.GetUnit(rogue.Options.TricksOfTheTradeTarget)
	}

	tricksOfTheTradeThreatTransferAura := rogue.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 59628},
		Label:    "TricksOfTheTradeThreatTransfer",
		Duration: core.TernaryDuration(hasGlyph, time.Second*10, time.Second*6),
	})

	tricksOfTheTradeApplicationAura := rogue.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 57934},
		Label:    "TricksOfTheTradeApplication",
		Duration: 30 * time.Second,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				tricksOfTheTradeThreatTransferAura.Activate(sim)
				if targetUnit != nil {
					core.TricksOfTheTradeAura(targetUnit, rogue.Index, hasGlyph).Activate(sim)
				}
				aura.Deactivate(sim)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.TricksOfTheTrade.CD.Set(core.NeverExpires)
			rogue.UpdateMajorCooldowns()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.TricksOfTheTrade.CD.Set(sim.CurrentTime + time.Second*time.Duration(30-5*rogue.Talents.FilthyTricks))
			rogue.UpdateMajorCooldowns()
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
}
