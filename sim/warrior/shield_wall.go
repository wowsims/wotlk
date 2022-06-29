package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

func (warrior *Warrior) RegisterShieldWallCD() {
	if warrior.Equip[proto.ItemSlot_ItemSlotOffHand].WeaponType != proto.WeaponType_WeaponTypeShield {
		return
	}

	duration := time.Second*10 + time.Second*2*time.Duration(warrior.Talents.ImprovedDisciplines)
	if warrior.Talents.ImprovedShieldWall == 1 {
		duration += time.Second * 3
	} else if warrior.Talents.ImprovedShieldWall == 2 {
		duration += time.Second * 5
	}

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

	cooldownDur := time.Minute * 30
	if warrior.Talents.ImprovedDisciplines == 1 {
		cooldownDur -= time.Minute * 4
	} else if warrior.Talents.ImprovedDisciplines == 2 {
		cooldownDur -= time.Minute * 7
	} else if warrior.Talents.ImprovedDisciplines == 3 {
		cooldownDur -= time.Minute * 10
	}
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
