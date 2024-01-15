package shaman

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (shaman *Shaman) BloodlustActionID() core.ActionID {
	return core.ActionID{
		SpellID: 2825,
		Tag:     shaman.Index,
	}
}

func (shaman *Shaman) registerBloodlustCD() {
	actionID := shaman.BloodlustActionID()

	blAuras := []*core.Aura{}
	for _, party := range shaman.Env.Raid.Parties {
		for _, partyMember := range party.Players {
			blAuras = append(blAuras, core.BloodlustAura(partyMember.GetCharacter(), actionID.Tag))
		}
	}

	shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.26,
			Multiplier: 1 - 0.02*float64(shaman.Talents.MentalQuickness),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: core.BloodlustCD,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Only cast if there is a player missing Sated.
			for _, playerUnit := range shaman.Env.Raid.AllPlayerUnits {
				if !playerUnit.HasActiveAura(core.SatedAuraLabel) {
					return true
				}
			}
			return false
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, blAura := range blAuras {
				target := blAura.Unit
				// Only activate bloodlust on units without sated.
				if !target.HasActiveAura(core.SatedAuraLabel) {
					blAura.Activate(sim)
				}
			}
		},
	})
}
