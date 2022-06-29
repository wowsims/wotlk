package shaman

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (shaman *Shaman) BloodlustActionID() core.ActionID {
	return core.ActionID{
		SpellID: 2825,
		Tag:     int32(shaman.Index),
	}
}

func (shaman *Shaman) registerBloodlustCD() {
	if !shaman.SelfBuffs.Bloodlust {
		return
	}
	actionID := shaman.BloodlustActionID()

	blAuras := []*core.Aura{}
	for _, partyMember := range shaman.Party.Players {
		blAuras = append(blAuras, core.BloodlustAura(partyMember.GetCharacter(), actionID.Tag))
	}

	baseCost := 750.0
	bloodlustSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(shaman.Talents.MentalQuickness)),
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: core.BloodlustCD,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, blAura := range blAuras {
				blAura.Activate(sim)
			}
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell:    bloodlustSpell,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			if character.CurrentMana() < bloodlustSpell.DefaultCast.Cost {
				return false
			}

			// Need to check if any party member has lust, not just self, because of
			// major CD ordering issues with the shared bloodlust.
			for _, partyMember := range character.Party.Players {
				if partyMember.GetCharacter().HasActiveAuraWithTag(core.BloodlustAuraTag) {
					return false
				}
			}
			return true
		},
	})
}
