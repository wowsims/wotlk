package encounters

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// Default implementation of TargetAI which takes a list of abilities as input
// in order of priority.
type DefaultAI struct {
	Target *core.Target

	Abilities []TargetAbility
}

type TargetAbility struct {
	// Puts ability on CD at the start of each iteration.
	InitialCD time.Duration

	// Probability (0-1) that this ability will be used when available.
	ChanceToUse float64

	// Factory function for creating the spell. Can use this or supply Spell
	// directly.
	MakeSpell func(*core.Target) *core.Spell

	Spell *core.Spell
}

func NewDefaultAI(abilities []TargetAbility) core.AIFactory {
	return func() core.TargetAI {
		return &DefaultAI{
			Abilities: abilities,
		}
	}
}

func (ai *DefaultAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target

	for i := range ai.Abilities {
		ability := &ai.Abilities[i]
		if ability.MakeSpell != nil {
			ability.Spell = ability.MakeSpell(target)
		}
	}
}

func (ai *DefaultAI) Reset(sim *core.Simulation) {

}

func (ai *DefaultAI) ExecuteCustomRotation(sim *core.Simulation) {
	for _, ability := range ai.Abilities {
		if sim.CurrentTime < ability.InitialCD {
			continue
		}

		if !ability.Spell.IsReady(sim) {
			continue
		}

		if sim.Proc(ability.ChanceToUse, "TargetAbility") {
			ability.Spell.Cast(sim, ai.Target.CurrentTarget)
			return
		}
	}
}
