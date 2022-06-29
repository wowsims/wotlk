package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warrior *Warrior) registerShieldBlockSpell() {
	actionID := core.ActionID{SpellID: 2565}

	shieldBlockAura := warrior.RegisterAura(core.Aura{
		Label:     "Shield Block",
		ActionID:  actionID,
		Duration:  time.Second*5 + core.TernaryDuration(warrior.Talents.ImprovedShieldBlock, time.Second, 0),
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, 75*core.BlockRatingPerBlockChance)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, -75*core.BlockRatingPerBlockChance)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeBlock) {
				aura.RemoveStack(sim)
			}
		},
	})

	cost := 10.0
	initialStacks := int32(1)
	if warrior.Talents.ImprovedShieldBlock {
		initialStacks++
	}

	warrior.ShieldBlock = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shieldBlockAura.Activate(sim)
			shieldBlockAura.SetStacks(sim, initialStacks)
		},
	})
}

func (warrior *Warrior) CanShieldBlock(sim *core.Simulation) bool {
	return warrior.PseudoStats.CanBlock &&
		warrior.StanceMatches(DefensiveStance) &&
		warrior.CurrentRage() >= warrior.ShieldBlock.DefaultCast.Cost &&
		warrior.ShieldBlock.IsReady(sim)
}
