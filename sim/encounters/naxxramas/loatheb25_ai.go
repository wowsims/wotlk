package naxxramas

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addLoatheb25(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        16011,
			Name:      "Loatheb 25",
			Level:     83,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      26_286_324,
				stats.Armor:       10643,
				stats.AttackPower: 574,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.2,
			MinBaseDamage:    6727,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
		},
		AI: NewLoatheb25AI(),
	})
	core.AddPresetEncounter("Loatheb 25", []string{
		bossPrefix + "/Loatheb 25",
	})
}

type Loatheb25AI struct {
	Target *core.Target
}

func NewLoatheb25AI() core.AIFactory {
	return func() core.TargetAI {
		return &Loatheb25AI{}
	}
}

func (ai *Loatheb25AI) Initialize(target *core.Target) {
	ai.Target = target
}

func (ai *Loatheb25AI) DoAction(sim *core.Simulation) {
	ai.Target.DoNothing()
}
