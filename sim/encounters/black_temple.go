package encounters

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func registerBlackTemple() {
	const bossPrefix = "Black Temple"

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22887,
			Name:      "High Warlord Naj'entus",
			Level:     73,
			MobType:   proto.MobType_MobTypeHumanoid,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      3_790_000,
				stats.Armor:       7684,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    9232.5,
			CanCrush:         true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22898,
			Name:      "Supremus",
			Level:     73,
			MobType:   proto.MobType_MobTypeDemon,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      4_552_800,
				stats.Armor:       7684,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.5,
			MinBaseDamage:    12365.25,
			CanCrush:         false,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22841,
			Name:      "Shade of Akama",
			Level:     73,
			MobType:   proto.MobType_MobTypeHumanoid,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      1_001_616,
				stats.Armor:       7684,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    19784.5,
			CanCrush:         true,
			ParryHaste:       true,
			DualWield:        true,
			DualWieldPenalty: false,
		},
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22871,
			Name:      "Teron Gorefiend",
			Level:     73,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      5_007_750,
				stats.Armor:       6193,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    16301,
			CanCrush:         true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22948,
			Name:      "Gurtogg Bloodboil",
			Level:     73,
			MobType:   proto.MobType_MobTypeHumanoid,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      5_691_000,
				stats.Armor:       7684,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    9892.5,
			CanCrush:         true,
			ParryHaste:       true,
			DualWield:        true,
			DualWieldPenalty: false,
		},
	})

	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        23418,
			Name:      "Essence of Suffering",
			Level:     73,
			MobType:   proto.MobType_MobTypeUnknown,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:           3_034_700,
				stats.Armor:            0,
				stats.AttackPower:      320,
				stats.BlockValue:       54,
				stats.MeleeCrit:        -1,    // Disable crit
				stats.ArmorPenetration: 99999, // Ignore player armor
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.0,
			MinBaseDamage:    915,
			CanCrush:         true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})
	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        23419,
			Name:      "Essence of Desire",
			Level:     73,
			MobType:   proto.MobType_MobTypeUnknown,
			TankIndex: 1,

			Stats: stats.Stats{
				stats.Health:      3_034_700,
				stats.Armor:       7684,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    9807,
			CanCrush:         true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})
	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        23420,
			Name:      "Essence of Anger",
			Level:     73,
			MobType:   proto.MobType_MobTypeUnknown,
			TankIndex: 2,

			Stats: stats.Stats{
				stats.Health:      2_276_400,
				stats.Armor:       7684,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    11875,
			CanCrush:         true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})
	core.AddPresetEncounter("Reliquary of Souls", []string{
		bossPrefix + "/Essence of Suffering",
		bossPrefix + "/Essence of Desire",
		bossPrefix + "/Essence of Anger",
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22947,
			Name:      "Mother Shahraz",
			Level:     73,
			MobType:   proto.MobType_MobTypeDemon,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      4_552_500,
				stats.Armor:       6193,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    18338.5,
			CanCrush:         false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})

	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22949,
			Name:      "Gathios the Shatterer",
			Level:     73,
			MobType:   proto.MobType_MobTypeHumanoid,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      1_746_500,
				stats.Armor:       6193,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    15281.85,
			CanCrush:         true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})
	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22952,
			Name:      "Veras Darkshadow",
			Level:     73,
			MobType:   proto.MobType_MobTypeHumanoid,
			TankIndex: 1,

			Stats: stats.Stats{
				stats.Health:           1_746_500,
				stats.Armor:            137,
				stats.ArcaneResistance: 75,
				stats.FireResistance:   75,
				stats.FrostResistance:  75,
				stats.NatureResistance: 75,
				stats.ShadowResistance: 75,
				stats.AttackPower:      320,
				stats.BlockValue:       54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    8792.95,
			CanCrush:         true,
			ParryHaste:       true,
			DualWield:        true,
			DualWieldPenalty: true,
		},
	})
	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22951,
			Name:      "Lady Malande",
			Level:     73,
			MobType:   proto.MobType_MobTypeHumanoid,
			TankIndex: 2,

			Stats: stats.Stats{
				stats.Health:      1_746_500,
				stats.Armor:       6193,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    5093.7,
			CanCrush:         true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})
	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22950,
			Name:      "High Nethermancer Zerevor",
			Level:     73,
			MobType:   proto.MobType_MobTypeHumanoid,
			TankIndex: 3,

			Stats: stats.Stats{
				stats.Health:      1_746_500,
				stats.Armor:       6193,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       0.0,
			MinBaseDamage:    0,
			CanCrush:         true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})
	core.AddPresetEncounter("Illidari Council", []string{
		bossPrefix + "/Gathios the Shatterer",
		bossPrefix + "/Veras Darkshadow",
		bossPrefix + "/Lady Malande",
		bossPrefix + "/High Nethermancer Zerevor",
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22917,
			Name:      "Illidan Stormrage",
			Level:     73,
			MobType:   proto.MobType_MobTypeDemon,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      6_070_400,
				stats.Armor:       7684,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.5,
			MinBaseDamage:    16486.65,
			CanCrush:         false,
			ParryHaste:       true,
			DualWield:        true,
			DualWieldPenalty: false,
		},
	})

	AddSingleTargetBossEncounter(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        22997,
			Name:      "Flame of Azzinoth",
			Level:     73,
			MobType:   proto.MobType_MobTypeDemon,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      1_138_200,
				stats.Armor:       7684,
				stats.AttackPower: 320,
				stats.BlockValue:  54,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolFire,
			SwingSpeed:       1.5,
			MinBaseDamage:    6594.9,
			CanCrush:         false,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
	})
}
