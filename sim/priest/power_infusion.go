package priest

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (priest *Priest) registerPowerInfusionCD() {
	if !priest.Talents.PowerInfusion {
		return
	}

	actionID := core.ActionID{SpellID: 10060, Tag: int32(priest.Index)}
	baseCost := priest.BaseMana() * 0.16

	powerInfusionTargetAgent := priest.Party.Raid.GetPlayerFromRaidTarget(priest.SelfBuffs.PowerInfusionTarget)
	if powerInfusionTargetAgent == nil {
		return
	}
	powerInfusionTarget := powerInfusionTargetAgent.GetCharacter()
	powerInfusionAura := core.PowerInfusionAura(powerInfusionTarget, actionID.Tag)

	piSpell := priest.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: core.PowerInfusionCD,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			powerInfusionAura.Activate(sim)
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    piSpell,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeMana,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentMana() >= baseCost
		},
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// How can we determine the target will be able to continue casting
			// 	for the next 15s at 20% reduced mana cost? Arbitrary value until then.
			//if powerInfusionTarget.CurrentMana() < 3000 {
			//	return false
			//}
			return !powerInfusionTarget.HasActiveAuraWithTag(core.BloodlustAuraTag)
		},
	})
}
