package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warrior *Warrior) RegisterShieldWallCD() {
	if warrior.Equip[proto.ItemSlot_ItemSlotOffHand].WeaponType != proto.WeaponType_WeaponTypeShield {
		return
	}

	duration := time.Second*10 + time.Second*2*time.Duration(warrior.Talents.ImprovedDisciplines)

	actionID := core.ActionID{SpellID: 871}
	swAura := warrior.RegisterAura(core.Aura{
		Label:    "Shield Wall",
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier *= 0.25
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier /= 0.25
		},
	})

	cooldownDur := time.Minute*5 - 30*time.Second*time.Duration(warrior.Talents.ImprovedDisciplines)
	swSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownDur,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			swAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: swSpell,
		Type:  core.CooldownTypeSurvival,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return warrior.StanceMatches(DefensiveStance)
		},
	})
}
