package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) RegisterShieldBlockCD() {
	actionID := core.ActionID{SpellID: 2565}

	statDep := warrior.NewDynamicMultiplyStat(stats.BlockValue, 2)
	warrior.ShieldBlockAura = warrior.RegisterAura(core.Aura{
		Label:    "Shield Block",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, 100*core.BlockRatingPerBlockChance)
			// TODO: The innate block value from the shield item should not be multiplied
			// as shown here https://youtu.be/LYJdkimJgn8?t=476
			warrior.EnableDynamicStatDep(sim, statDep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, -100*core.BlockRatingPerBlockChance)
			warrior.DisableDynamicStatDep(sim, statDep)
		},
	})

	warrior.ShieldBlock = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second*60 - time.Second*10*time.Duration(warrior.Talents.ShieldMastery),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.ShieldBlockAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: warrior.ShieldBlock,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return warrior.PseudoStats.CanBlock && warrior.StanceMatches(DefensiveStance)
		},
	})
}
