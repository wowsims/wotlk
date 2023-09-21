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
	if !shaman.SelfBuffs.Bloodlust && !shaman.IsUsingAPL {
		return
	}
	actionID := shaman.BloodlustActionID()

	blAuras := []*core.Aura{}
	for _, party := range shaman.Env.Raid.Parties {
		for _, partyMember := range party.Players {
			blAuras = append(blAuras, core.BloodlustAura(partyMember.GetCharacter(), actionID.Tag))
		}
	}

	bloodlustSpell := shaman.RegisterSpell(core.SpellConfig{
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
			// Need to check if any raid member has lust, not just self, because of
			// major CD ordering issues with the shared bloodlust.
			for _, party := range shaman.Env.Raid.Parties {
				for _, partyMember := range party.Players {
					if partyMember.GetCharacter().HasActiveAuraWithTag(core.BloodlustAuraTag) {
						return false
					}
				}
			}
			return true
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, blAura := range blAuras {
				blAura.Activate(sim)
			}
		},
	})

	if !shaman.IsUsingAPL {
		shaman.AddMajorCooldown(core.MajorCooldown{
			Spell:    bloodlustSpell,
			Priority: core.CooldownPriorityBloodlust,
			Type:     core.CooldownTypeDPS,
		})
	}
}
