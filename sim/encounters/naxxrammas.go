package encounters

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func registerNaxxramas10() {
	const bossPrefix = "Naxxrammas 10"

	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        16028,
			Name:      "Patchwerk",
			Level:     83,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      5_691_835,
				stats.Armor:       10643,
				stats.AttackPower: 640,
				stats.BlockValue:  108,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.6,
			MinBaseDamage:    14135,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        true,
			DualWieldPenalty: false,
		},
		AI: NewDefaultAI([]TargetAbility{
			Patchwerk10HatefulStrike,
		}),
	})
	core.AddPresetEncounter("Patchwerk", []string{
		bossPrefix + "/Patchwerk",
	})
}

func registerNaxxramas25() {
	const bossPrefix = "Naxxrammas 25"

	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        16028,
			Name:      "Patchwerk",
			Level:     83,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      16_950_147,
				stats.Armor:       10643,
				stats.AttackPower: 640,
				stats.BlockValue:  108,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.6,
			MinBaseDamage:    14135,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        true,
			DualWieldPenalty: false,
		},
		AI: NewDefaultAI([]TargetAbility{
			Patchwerk25HatefulStrike,
		}),
	})
	core.AddPresetEncounter("Patchwerk", []string{
		bossPrefix + "/Patchwerk",
	})
}
