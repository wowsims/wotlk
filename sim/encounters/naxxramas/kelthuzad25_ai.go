package naxxramas

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addKelThuzad25(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        15990,
			Name:      "Kel'Thuzad 25",
			Level:     83,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      19_034_924,
				stats.Armor:       10643,
				stats.AttackPower: 574,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.3,
			MinBaseDamage:    28767,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
		},
		AI: NewKelThuzad25AI(),
	})
	core.AddPresetEncounter("Kel'Thuzad 25", []string{
		bossPrefix + "/Kel'Thuzad 25",
	})
}

type KelThuzad25AI struct {
	Target *core.Target
}

func NewKelThuzad25AI() core.AIFactory {
	return func() core.TargetAI {
		return &KelThuzad25AI{}
	}
}

func (ai *KelThuzad25AI) Initialize(target *core.Target) {
	ai.Target = target
}

func (ai *KelThuzad25AI) Reset(*core.Simulation) {
}

func (ai *KelThuzad25AI) DoAction(sim *core.Simulation) {
	ai.Target.DoNothing()
}
