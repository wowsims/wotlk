package sim

import (
	"testing"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// 1 moonkin, 1 ele shaman, 1 spriest, 2x arcane
var castersWithElemental = &proto.Party{
	Players: []*proto.Player{
		{
			Name:      "Balance Druid 1",
			Race:      proto.Race_RaceTauren,
			Class:     proto.Class_ClassDruid,
			Equipment: MoonkinEquipment,
			Spec: &proto.Player_BalanceDruid{
				BalanceDruid: &proto.BalanceDruid{
					Options: &proto.BalanceDruid_Options{
						InnervateTarget: &proto.UnitReference{},
					},
				},
			},
			Consumes: &proto.Consumes{
				DefaultPotion: proto.Potions_SuperManaPotion,
			},
			Buffs: &proto.IndividualBuffs{
				BlessingOfKings:  true,
				BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
			},
		},
		{
			Name:      "Shadow Priest 1",
			Race:      proto.Race_RaceUndead,
			Class:     proto.Class_ClassPriest,
			Equipment: ShadowEquipment,
			Spec: &proto.Player_ShadowPriest{
				ShadowPriest: &proto.ShadowPriest{
					Options: &proto.ShadowPriest_Options{},
				},
			},
			Consumes: &proto.Consumes{
				DefaultPotion: proto.Potions_SuperManaPotion,
			},
			Buffs: &proto.IndividualBuffs{
				BlessingOfKings:  true,
				BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
			},
		},
		{
			Name:      "Elemental Shaman 1",
			Race:      proto.Race_RaceTroll,
			Class:     proto.Class_ClassShaman,
			Equipment: ElementalEquipment,
			Spec: &proto.Player_ElementalShaman{
				ElementalShaman: &proto.ElementalShaman{
					Options: &proto.ElementalShaman_Options{
						Shield: proto.ShamanShield_WaterShield,
						Totems: &proto.ShamanTotems{
							Earth: proto.EarthTotem_TremorTotem,
							Air:   proto.AirTotem_WrathOfAirTotem,
							Fire:  proto.FireTotem_TotemOfWrath,
							Water: proto.WaterTotem_ManaSpringTotem,
						},
					},
				},
			},
			Consumes: &proto.Consumes{
				DefaultPotion: proto.Potions_SuperManaPotion,
			},
			Buffs: &proto.IndividualBuffs{
				BlessingOfKings:  true,
				BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
			},
		},
		{
			Name:      "Arcane Mage 1",
			Race:      proto.Race_RaceTroll,
			Class:     proto.Class_ClassMage,
			Equipment: ArcaneEquipment,
			Spec: &proto.Player_Mage{
				Mage: &proto.Mage{
					Options: &proto.Mage_Options{
						Armor: proto.Mage_Options_MageArmor,
					},
				},
			},
			Consumes: &proto.Consumes{
				DefaultPotion: proto.Potions_SuperManaPotion,
			},
			Buffs: &proto.IndividualBuffs{
				BlessingOfKings:  true,
				BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
			},
		},
	},
	Buffs: &proto.PartyBuffs{},
}

var castersWithResto = &proto.Party{
	Players: []*proto.Player{
		// 1 moonkin, 1 spriest, 2x arcane, 1 resto shaman
		{
			Name:      "Balance Druid 2",
			Race:      proto.Race_RaceTauren,
			Class:     proto.Class_ClassDruid,
			Equipment: MoonkinEquipment,
			Spec: &proto.Player_BalanceDruid{
				BalanceDruid: &proto.BalanceDruid{
					Options: &proto.BalanceDruid_Options{
						InnervateTarget: &proto.UnitReference{
							Type:  proto.UnitReference_Player,
							Index: 6,
						},
					},
				},
			},
			Consumes: &proto.Consumes{
				DefaultPotion: proto.Potions_SuperManaPotion,
			},
			Buffs: &proto.IndividualBuffs{
				BlessingOfKings:  true,
				BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
			},
		},
		{
			Name:      "Shadow Priest 2",
			Race:      proto.Race_RaceUndead,
			Class:     proto.Class_ClassPriest,
			Equipment: ShadowEquipment,
			Spec: &proto.Player_ShadowPriest{
				ShadowPriest: &proto.ShadowPriest{
					Options: &proto.ShadowPriest_Options{},
				},
			},
			Consumes: &proto.Consumes{
				DefaultPotion: proto.Potions_SuperManaPotion,
			},
			Buffs: &proto.IndividualBuffs{
				BlessingOfKings:  true,
				BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
			},
		},
		{
			Name:      "Arcane Mage 3",
			Race:      proto.Race_RaceTroll,
			Class:     proto.Class_ClassMage,
			Equipment: ArcaneEquipment,
			Spec: &proto.Player_Mage{
				Mage: &proto.Mage{
					Options: &proto.Mage_Options{
						Armor: proto.Mage_Options_MageArmor,
					},
				},
			},
			Consumes: &proto.Consumes{
				DefaultPotion: proto.Potions_SuperManaPotion,
			},
			Buffs: &proto.IndividualBuffs{
				BlessingOfKings:  true,
				BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
			},
		},
	},
	Buffs: &proto.PartyBuffs{
		ManaTideTotems: 1,
	},
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: &proto.Raid{
			Parties: []*proto.Party{
				castersWithElemental,
				castersWithResto,
				{
					Players: []*proto.Player{
						{
							Name:      "Enhancement Shaman 1",
							Race:      proto.Race_RaceTroll,
							Class:     proto.Class_ClassShaman,
							Equipment: EnhancementEquipment,
							Spec: &proto.Player_EnhancementShaman{
								EnhancementShaman: &proto.EnhancementShaman{
									Options: &proto.EnhancementShaman_Options{
										Shield:   proto.ShamanShield_LightningShield,
										SyncType: proto.ShamanSyncType_SyncMainhandOffhandSwings,
										Totems: &proto.ShamanTotems{
											Earth: proto.EarthTotem_TremorTotem,
											Air:   proto.AirTotem_WrathOfAirTotem,
											Fire:  proto.FireTotem_TotemOfWrath,
											Water: proto.WaterTotem_ManaSpringTotem,
										},
									},
								},
							},
							Consumes: &proto.Consumes{},
							Buffs: &proto.IndividualBuffs{
								BlessingOfKings:  true,
								BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
							},
						},
					},
				},
			},
			Buffs: &proto.RaidBuffs{
				GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
				ArcaneBrilliance: true,
				Bloodlust:        true,
				WrathOfAirTotem:  true,
				ManaSpringTotem:  proto.TristateEffect_TristateEffectImproved,
			},
			Debuffs: &proto.Debuffs{
				JudgementOfWisdom: true,
				CurseOfElements:   true,
			},
		},
		Encounter: &proto.Encounter{
			Duration:             180,
			ExecuteProportion_20: 0.1,
			Targets: []*proto.Target{
				{
					Stats:   stats.Stats{stats.Armor: 7684}.ToFloatArray(),
					MobType: proto.MobType_MobTypeDemon,
				},
			},
		},
		SimOptions: core.AverageDefaultSimTestOptions,
	}

	core.RaidBenchmark(b, rsr)
}

// P3 gear for each class

// Shadow Priest Equipment
var ShadowEquipment = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{
			Id:      46172,
			Enchant: 3820,
			Gems:    []int32{41285, 45883},
		},
		{
			Id:   45243,
			Gems: []int32{39998},
		},
		{
			Id:      46165,
			Enchant: 3810,
			Gems:    []int32{39998},
		},
		{
			Id:      45242,
			Enchant: 3722,
			Gems:    []int32{40049},
		},
		{
			Id:      46168,
			Enchant: 1144,
			Gems:    []int32{39998, 39998},
		},
		{
			Id:      45446,
			Enchant: 2332,
			Gems:    []int32{39998},
		},
		{
			Id:      45665,
			Enchant: 3604,
			Gems:    []int32{39998, 39998},
		},
		{
			Id:      45619,
			Enchant: 3601,
			Gems:    []int32{39998, 39998, 39998},
		},
		{
			Id:      46170,
			Enchant: 3719,
			Gems:    []int32{39998, 40049},
		},
		{
			Id:      45135,
			Enchant: 3606,
			Gems:    []int32{39998, 40049},
		},
		{
			Id:   45495,
			Gems: []int32{40026},
		},
		{
			Id:   46046,
			Gems: []int32{39998},
		},
		{
			Id: 45518,
		},
		{
			Id: 45466,
		},
		{
			Id:      45620,
			Enchant: 3834,
			Gems:    []int32{40026},
		},
		{
			Id: 45617,
		},
		{
			Id:   45294,
			Gems: []int32{39998},
		},
	},
}

// Arcane Equipment
var ArcaneEquipment = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{
			Id:      45497,
			Enchant: 3820,
			Gems:    []int32{41285, 45883},
		},
		{
			Id:   45243,
			Gems: []int32{39998},
		},
		{
			Id:      46134,
			Enchant: 3810,
			Gems:    []int32{39998},
		},
		{
			Id:      45618,
			Enchant: 3722,
			Gems:    []int32{40026},
		},
		{
			Id:      46130,
			Enchant: 3832,
			Gems:    []int32{39998, 39998},
		},
		{
			Id:      45446,
			Enchant: 2332,
			Gems:    []int32{39998},
		},
		{
			Id:      45665,
			Enchant: 3604,
			Gems:    []int32{39998, 39998},
		},
		{
			Id:   45619,
			Gems: []int32{39998, 39998, 39998},
		},
		{
			Id:      45488,
			Enchant: 3719,
			Gems:    []int32{39998, 40051, 40026},
		},
		{
			Id:      45135,
			Enchant: 3606,
			Gems:    []int32{39998, 39998},
		},
		{
			Id:   46046,
			Gems: []int32{39998},
		},
		{
			Id:   45495,
			Gems: []int32{39998},
		},
		{
			Id: 45466,
		},
		{
			Id: 45518,
		},
		{
			Id:      45620,
			Enchant: 3834,
			Gems:    []int32{39998},
		},
		{
			Id: 45617,
		},
		{
			Id:      45294,
			Enchant: 0,
			Gems:    []int32{39998},
		},
	},
}

// Moonkin Equipment
var MoonkinEquipment = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{
			Id:      45497,
			Enchant: 3820,
			Gems:    []int32{41285, 42144},
		},
		{
			Id:   45133,
			Gems: []int32{40048},
		},
		{
			Id:      46196,
			Enchant: 3810,
			Gems:    []int32{39998},
		},
		{
			Id:      45242,
			Enchant: 3859,
			Gems:    []int32{40048},
		},
		{
			Id:      45519,
			Enchant: 3832,
			Gems:    []int32{40051, 42144, 40026},
		},
		{
			Id:      45446,
			Enchant: 2332,
			Gems:    []int32{42144},
		},
		{
			Id:      45665,
			Enchant: 3604,
			Gems:    []int32{39998, 39998},
		},
		{
			Id:   45619,
			Gems: []int32{39998, 39998, 39998},
		},
		{
			Id:      46192,
			Enchant: 3719,
			Gems:    []int32{39998, 39998},
		},
		{
			Id:      45537,
			Enchant: 3606,
			Gems:    []int32{39998, 40026},
		},
		{
			Id:   46046,
			Gems: []int32{39998},
		},
		{
			Id:   45495,
			Gems: []int32{39998},
		},
		{
			Id: 45466,
		},
		{
			Id: 45518,
		},
		{
			Id:      45620,
			Enchant: 3834,
			Gems:    []int32{39998},
		},
		{
			Id: 45617,
		},
		{
			Id: 40321,
		},
	},
}

var EnhancementEquipment = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{
			Id:      45610,
			Enchant: 3817,
			Gems:    []int32{41398, 42702},
		},
		{
			Id:   45517,
			Gems: []int32{39999},
		},
		{
			Id:      46203,
			Enchant: 3808,
			Gems:    []int32{39999},
		},
		{
			Id:      45461,
			Enchant: 3831,
			Gems:    []int32{40014},
		},
		{
			Id:      46205,
			Enchant: 3832,
			Gems:    []int32{40058, 40053},
		},
		{
			Id:      45460,
			Enchant: 3845,
			Gems:    []int32{39999},
		},
		{
			Id:      46200,
			Enchant: 3604,
			Gems:    []int32{40014},
		},
		{
			Id:      45553,
			Enchant: 0,
			Gems:    []int32{36766, 36766, 36766},
		},
		{
			Id:      46208,
			Enchant: 3823,
			Gems:    []int32{39999, 39999},
		},
		{
			Id:      45989,
			Enchant: 3606,
			Gems:    []int32{40053, 39999},
		},
		{
			Id:   45456,
			Gems: []int32{39999},
		},
		{
			Id:   46046,
			Gems: []int32{40053},
		},
		{
			Id: 45609,
		},
		{
			Id: 46038,
		},
		{
			Id:      45612,
			Enchant: 3789,
			Gems:    []int32{39999},
		},
		{
			Id:      46097,
			Enchant: 3789,
			Gems:    []int32{40003},
		},
		{
			Id: 40322,
		},
	},
}

// Elemental Equipment
var ElementalEquipment = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{
			Id:      46209,
			Enchant: 3820,
			Gems:    []int32{41285, 40048},
		},
		{
			Id:   45933,
			Gems: []int32{39998},
		},
		{
			Id:      46211,
			Enchant: 3810,
			Gems:    []int32{39998},
		},
		{
			Id:      45242,
			Enchant: 3722,
			Gems:    []int32{39998},
		},
		{
			Id:      46206,
			Enchant: 3832,
			Gems:    []int32{39998, 39998},
		},
		{
			Id:      45460,
			Enchant: 2332,
			Gems:    []int32{39998},
		},
		{
			Id:      45665,
			Enchant: 3604,
			Gems:    []int32{39998, 39998},
		},
		{
			Id:      45616,
			Enchant: 3599,
			Gems:    []int32{39998, 39998, 39998},
		},
		{
			Id:      46210,
			Enchant: 3721,
			Gems:    []int32{39998, 40027},
		},
		{
			Id:      45537,
			Enchant: 3606,
			Gems:    []int32{39998, 40027},
		},
		{
			Id:   46046,
			Gems: []int32{39998},
		},
		{
			Id:   45495,
			Gems: []int32{39998},
		},
		{
			Id: 45518,
		},
		{
			Id: 40255,
		},
		{
			Id:      45612,
			Enchant: 3834,
			Gems:    []int32{39998},
		},
		{
			Id:      45470,
			Enchant: 1128,
			Gems:    []int32{39998},
		},
		{
			Id: 40267,
		},
	},
}
