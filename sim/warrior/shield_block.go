package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) RegisterShieldBlockCD() {
	actionID := core.ActionID{SpellID: 2565}
	blockValueMult := 2.0

	shieldBlockAura := warrior.RegisterAura(core.Aura{
		Label:    "Shield Block",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, 100*core.BlockRatingPerBlockChance)
			// TODO: The innate block value from the shield item should not be multiplied
			// as shown here https://youtu.be/LYJdkimJgn8?t=476
			warrior.AddStatDependencyDynamic(sim, stats.BlockValue, stats.BlockValue, blockValueMult)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, -100*core.BlockRatingPerBlockChance)
			warrior.AddStatDependencyDynamic(sim, stats.BlockValue, stats.BlockValue, 1.0/blockValueMult)
		},
	})

	cooldownDuration := time.Second * 60
	cooldownDuration -= time.Second * 10 * time.Duration(warrior.Talents.ShieldMastery)

	warrior.ShieldBlock = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownDuration,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shieldBlockAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: warrior.ShieldBlock,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return warrior.PseudoStats.CanBlock &&
				warrior.StanceMatches(DefensiveStance) &&
				warrior.ShieldBlock.IsReady(sim)
		},
	})
}
