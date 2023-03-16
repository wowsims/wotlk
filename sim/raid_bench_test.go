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
					Rotation: &proto.BalanceDruid_Rotation{
						Type: proto.BalanceDruid_Rotation_Default,
					},
					Options: &proto.BalanceDruid_Options{
						InnervateTarget: &proto.RaidTarget{
							TargetIndex: 0,
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
			Name:      "Shadow Priest 1",
			Race:      proto.Race_RaceUndead,
			Class:     proto.Class_ClassPriest,
			Equipment: ShadowEquipment,
			Spec: &proto.Player_ShadowPriest{
				ShadowPriest: &proto.ShadowPriest{
					Rotation: &proto.ShadowPriest_Rotation{
						RotationType: proto.ShadowPriest_Rotation_Ideal,
					},
					Options: &proto.ShadowPriest_Options{
						UseShadowfiend: true,
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
			Name:      "Elemental Shaman 1",
			Race:      proto.Race_RaceTroll,
			Class:     proto.Class_ClassShaman,
			Equipment: ElementalEquipment,
			Spec: &proto.Player_ElementalShaman{
				ElementalShaman: &proto.ElementalShaman{
					Rotation: &proto.ElementalShaman_Rotation{
						Totems: &proto.ShamanTotems{
							Earth: proto.EarthTotem_TremorTotem,
							Air:   proto.AirTotem_WrathOfAirTotem,
							Fire:  proto.FireTotem_TotemOfWrath,
							Water: proto.WaterTotem_ManaSpringTotem,
						},
						Type: proto.ElementalShaman_Rotation_Adaptive,
					},
					Options: &proto.ElementalShaman_Options{
						Shield:    proto.ShamanShield_WaterShield,
						Bloodlust: true,
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
					Rotation: &proto.Mage_Rotation{},
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
					Rotation: &proto.BalanceDruid_Rotation{
						Type: proto.BalanceDruid_Rotation_Default,
					},
					Options: &proto.BalanceDruid_Options{
						InnervateTarget: &proto.RaidTarget{
							TargetIndex: 6,
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
					Rotation: &proto.ShadowPriest_Rotation{
						RotationType: proto.ShadowPriest_Rotation_Ideal,
					},
					Options: &proto.ShadowPriest_Options{
						UseShadowfiend: true,
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
			Name:      "Arcane Mage 3",
			Race:      proto.Race_RaceTroll,
			Class:     proto.Class_ClassMage,
			Equipment: ArcaneEquipment,
			Spec: &proto.Player_Mage{
				Mage: &proto.Mage{
					Options: &proto.Mage_Options{
						Armor: proto.Mage_Options_MageArmor,
					},
					Rotation: &proto.Mage_Rotation{},
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
							Name:  "Enhancement Shaman 1",
							Race:  proto.Race_RaceTroll,
							Class: proto.Class_ClassShaman,
							Spec: &proto.Player_EnhancementShaman{
								EnhancementShaman: &proto.EnhancementShaman{
									Rotation: &proto.EnhancementShaman_Rotation{
										Totems: &proto.ShamanTotems{
											Earth: proto.EarthTotem_TremorTotem,
											Air:   proto.AirTotem_WrathOfAirTotem,
											Fire:  proto.FireTotem_TotemOfWrath,
											Water: proto.WaterTotem_ManaSpringTotem,
										},
									},
									Options: &proto.EnhancementShaman_Options{
										Shield:    proto.ShamanShield_LightningShield,
										Bloodlust: true,
										SyncType:  proto.ShamanSyncType_SyncMainhandOffhandSwings,
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
			Id:      31064,
			Enchant: 29191,
			Gems: []int32{
				25893,
				32215,
			},
		},
		{
			Id: 30666,
		},
		{
			Id:      31070,
			Enchant: 28886,
			Gems: []int32{
				32196,
				32196,
			},
		},
		{
			Id:      32590,
			Enchant: 33150,
		},
		{
			Id:      31065,
			Enchant: 24003,
			Gems: []int32{
				32196,
				32196,
				32196,
			},
		},
		{
			Id:      32586,
			Enchant: 22534,
		},
		{
			Id:      31061,
			Enchant: 28272,
			Gems: []int32{
				32196,
			},
		},
		{
			Id: 32256,
		},
		{
			Id:      30916,
			Enchant: 24274,
			Gems: []int32{
				32196,
				32196,
				32196,
			},
		},
		{
			Id:      32239,
			Enchant: 35297,
			Gems: []int32{
				32196,
				32196,
			},
		},
		{
			Id:      32527,
			Enchant: 22536,
		},
		{
			Id:      32527,
			Enchant: 22536,
		},
		{
			Id: 32483,
		},
		{
			Id: 29370,
		},
		{
			Id:      32374,
			Enchant: 22561,
		},
		{
			Id: 29982,
		},
	},
}

// Arcane Equipment
var ArcaneEquipment = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{
			Id:      30206,
			Enchant: 29191,
			Gems: []int32{
				34220,
				32204,
			},
		},
		{
			Id: 30015,
		},
		{
			Id:      30210,
			Enchant: 28886,
			Gems: []int32{
				32204,
				32215,
			},
		},
		{
			Id:      32331,
			Enchant: 33150,
		},
		{
			Id:      30196,
			Enchant: 24003,
			Gems: []int32{
				32204,
				32204,
				32215,
			},
		},
		{
			Id:      30870,
			Enchant: 22534,
			Gems: []int32{
				32204,
			},
		},
		{
			Id:      30205,
			Enchant: 28272,
		},
		{
			Id: 30888,
			Gems: []int32{
				32204,
				32204,
			},
		},
		{
			Id:      31058,
			Enchant: 24274,
			Gems: []int32{
				32204,
			},
		},
		{
			Id:      32239,
			Enchant: 35297,
			Gems: []int32{
				32204,
				32204,
			},
		},
		{
			Id:      32527,
			Enchant: 22536,
		},
		{
			Id:      29305,
			Enchant: 22536,
		},
		{
			Id: 32483,
		},
		{
			Id: 30720,
		},
		{
			Id:      32374,
			Enchant: 22560,
		},
		{},
		{
			Id: 28783,
		},
	},
}

// Moonkin Equipment
var MoonkinEquipment = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{
			Id:      31040,
			Enchant: 29191,
			Gems: []int32{
				32218,
				34220,
			},
		},
		{
			Id: 30015,
		},
		{
			Id:      31049,
			Enchant: 28886,
			Gems: []int32{
				32215,
				32218,
			},
		},
		{
			Id:      32331,
			Enchant: 33150,
		},
		{
			Id:      31043,
			Enchant: 24003,
			Gems: []int32{
				32196,
				32196,
				32196,
			},
		},
		{
			Id:      32586,
			Enchant: 22534,
		},
		{
			Id:      31035,
			Enchant: 28272,
			Gems: []int32{
				32218,
			},
		},
		{
			Id: 30914,
		},
		{
			Id:      30916,
			Enchant: 24274,
			Gems: []int32{
				32196,
				32196,
				32196,
			},
		},
		{
			Id:      32352,
			Enchant: 35297,
			Gems: []int32{
				32218,
				32215,
			},
		},
		{
			Id:      32527,
			Enchant: 22536,
		},
		{
			Id:      29305,
			Enchant: 22536,
		},
		{
			Id: 32486,
		},
		{
			Id: 32483,
		},
		{
			Id:      32374,
			Enchant: 22560,
		},
		{
			Id: 32387,
		},
	},
}

// Elemental Equipment
var ElementalEquipment = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{
			Id:      31014,
			Enchant: 29191,
			Gems: []int32{
				34220,
				32215,
			},
		},
		{
			Id: 30015,
		},
		{
			Id:      31023,
			Enchant: 28886,
			Gems: []int32{
				32215,
				32218,
			},
		},
		{
			Id:      32331,
			Enchant: 33150,
		},
		{
			Id:      31017,
			Enchant: 24003,
			Gems: []int32{
				32196,
				32196,
				32196,
			},
		},
		{
			Id:      32586,
			Enchant: 22534,
		},
		{
			Id:      31008,
			Enchant: 28272,
			Gems: []int32{
				32218,
			},
		},
		{
			Id: 32276,
		},
		{
			Id:      30916,
			Enchant: 24274,
			Gems: []int32{
				32196,
				32196,
				32196,
			},
		},
		{
			Id:      32352,
			Enchant: 35297,
			Gems: []int32{
				32196,
				32196,
			},
		},
		{
			Id:      32527,
			Enchant: 22536,
		},
		{
			Id:      29305,
			Enchant: 22536,
		},
		{
			Id: 32483,
		},
		{
			Id: 28785,
		},
		{
			Id:      32374,
			Enchant: 22555,
		},
		{
			Id: 32330,
		},
	},
}
