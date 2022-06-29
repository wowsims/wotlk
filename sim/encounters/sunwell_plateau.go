package encounters

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func registerSunwellPlateau() {
	const bossPrefix = "Sunwell Plateau"

	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        24850,
			Name:      "Kalecgos",
			Level:     73,
			MobType:   proto.MobType_MobTypeDragonkin,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      2_785_000,
				stats.Armor:       6193,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    14135,
			CanCrush:         false,
			SuppressDodge:    true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})
	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        24892,
			Name:      "Sathrovarr the Corruptor",
			Level:     73,
			MobType:   proto.MobType_MobTypeDemon,
			TankIndex: 1,

			Stats: stats.Stats{
				stats.Health:      2_785_000,
				stats.Armor:       6193,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    14135,
			CanCrush:         false,
			SuppressDodge:    true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})
	core.AddPresetEncounter("Kalecgos", []string{
		bossPrefix + "/Kalecgos",
		bossPrefix + "/Sathrovarr the Corruptor",
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        24882,
			Name:      "Brutallus",
			Level:     73,
			MobType:   proto.MobType_MobTypeDemon,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      10_130_000,
				stats.Armor:       7684,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.0,
			MinBaseDamage:    17250,
			CanCrush:         false,
			SuppressDodge:    true,
			ParryHaste:       false,
			DualWield:        true,
			DualWieldPenalty: false,
		},
		AI: NewDefaultAI([]TargetAbility{
			BrutallusStomp,
		}),
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        25038,
			Name:      "Felmyst",
			Level:     73,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      6_840_000,
				stats.Armor:       6193,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    13793.85,
			CanCrush:         false,
			SuppressDodge:    true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        25165,
			Name:      "Lady Sacrolash",
			Level:     73,
			MobType:   proto.MobType_MobTypeDemon,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      2_874_000,
				stats.Armor:       6193,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    14127.7,
			CanCrush:         true,
			SuppressDodge:    true,
			ParryHaste:       true,
			DualWield:        true,
			DualWieldPenalty: false,
		},
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        25840,
			Name:      "Entropius",
			Level:     73,
			MobType:   proto.MobType_MobTypeDemon,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      2_520_000,
				stats.Armor:       7684,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.5,
			MinBaseDamage:    18338.5,
			CanCrush:         false,
			SuppressDodge:    true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        25315,
			Name:      "Kil'jaeden",
			Level:     73,
			MobType:   proto.MobType_MobTypeDemon,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      12_600_000,
				stats.Armor:       6193,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    18338.5,
			CanCrush:         false,
			SuppressDodge:    true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})
}
