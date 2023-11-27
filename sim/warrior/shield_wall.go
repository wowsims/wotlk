package warrior

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
	"github.com/wowsims/classic/sod/sim/core/proto"
)

func (warrior *Warrior) RegisterShieldWallCD() {
	if warrior.OffHand().WeaponType != proto.WeaponType_WeaponTypeShield {
		return
	}
	duration := time.Duration(10+[]float64{3, 5}[warrior.Talents.ImprovedShieldWall]) * time.Second
	//This is the inverse of the tooltip since it is a damage TAKEN coefficient
	damageTaken := 0.25

	actionID := core.ActionID{SpellID: 871}
	swAura := warrior.RegisterAura(core.Aura{
		Label:    "Shield Wall",
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier *= damageTaken
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier /= damageTaken
		},
	})

	cooldownDur := time.Minute * 30

	swSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownDur,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(DefensiveStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			swAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: swSpell,
		Type:  core.CooldownTypeSurvival,
	})
}
