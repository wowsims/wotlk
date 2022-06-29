package encounters

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
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

func (ai *DefaultAI) Initialize(target *core.Target) {
	ai.Target = target

	for i, _ := range ai.Abilities {
		ability := &ai.Abilities[i]
		if ability.MakeSpell != nil {
			ability.Spell = ability.MakeSpell(target)
		}
	}
}

func (ai *DefaultAI) DoAction(sim *core.Simulation) {
	for _, ability := range ai.Abilities {
		if sim.CurrentTime < ability.InitialCD {
			continue
		}

		if !ability.Spell.IsReady(sim) {
			continue
		}

		if ability.ChanceToUse == 1 || sim.RandomFloat("TargetAbility") < ability.ChanceToUse {
			ability.Spell.Cast(sim, ai.Target.CurrentTarget)
			return
		}
	}
}
