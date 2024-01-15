package encounters

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func addLevel25(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        213334, // Aku'mai
			Name:      "Level 25",
			Level:     27,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      127_393, // Aku'mai
				stats.Armor:       1104,    // Level 27 presumed
				stats.AttackPower: 574,     // TODO: Find out attack power
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2,      // TODO:
			MinBaseDamage:    400,    // TODO:
			DamageSpread:     0.3333, // TODO:
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
	})
	core.AddPresetEncounter("Level 25", []string{
		bossPrefix + "/Level 25",
	})
}

func addLevel40(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        213334, // TODO:
			Name:      "Level 40",
			Level:     42,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      127_393, // TODO:
				stats.Armor:       2053,    // TODO:
				stats.AttackPower: 574,     // TODO:
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2,      // TODO:
			MinBaseDamage:    1000,   // TODO:
			DamageSpread:     0.3333, // TODO:
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
	})
	core.AddPresetEncounter("Level 40", []string{
		bossPrefix + "/Level 40",
	})
}

func addLevel50(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        213334, // TODO:
			Name:      "Level 50",
			Level:     52,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      127_393, // TODO:
				stats.Armor:       2053,    // TODO:
				stats.AttackPower: 574,     // TODO:
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2,      // TODO:
			MinBaseDamage:    2000,   // TODO:
			DamageSpread:     0.3333, // TODO:
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
	})
	core.AddPresetEncounter("Level 50", []string{
		bossPrefix + "/Level 50",
	})
}

func addLevel60(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        213334, // TODO:
			Name:      "Level 60",
			Level:     63,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      127_393, // TODO:
				stats.Armor:       3731,    // TODO:
				stats.AttackPower: 805,     // TODO:
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2,      // TODO:
			MinBaseDamage:    3000,   // TODO:
			DamageSpread:     0.3333, // TODO:
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
	})
	core.AddPresetEncounter("Level 60", []string{
		bossPrefix + "/Level 60",
	})
}
