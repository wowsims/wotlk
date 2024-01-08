package naxxramas

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addThaddius25(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        15990,
			Name:      "Thaddius",
			Level:     83,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      39_520_129,
				stats.Armor:       10643,
				stats.AttackPower: 805,
				stats.BlockValue:  76,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.25,
			MinBaseDamage:    23442,
			DamageSpread:     0.3333,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
		AI: NewThaddius25AI(),
	})
	core.AddPresetEncounter("Thaddius", []string{
		bossPrefix + "/Thaddius",
	})
}

type Thaddius25AI struct {
	Target *core.Target
}

func NewThaddius25AI() core.AIFactory {
	return func() core.TargetAI {
		return &Thaddius25AI{}
	}
}

func (ai *Thaddius25AI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
}

func (ai *Thaddius25AI) Reset(*core.Simulation) {
}

func (ai *Thaddius25AI) ExecuteCustomRotation(sim *core.Simulation) {
}
