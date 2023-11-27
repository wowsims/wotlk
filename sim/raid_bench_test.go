package sim

import (
	"testing"

	"github.com/wowsims/classic/sod/sim/core"
	"github.com/wowsims/classic/sod/sim/core/proto"
	"github.com/wowsims/classic/sod/sim/core/stats"
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
						InnervateTarget: &proto.UnitReference{},
					},
				},
			},
			Consumes: &proto.Consumes{},
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
			Consumes: &proto.Consumes{},
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
						Type: proto.ElementalShaman_Rotation_Adaptive,
					},
					Options: &proto.ElementalShaman_Options{
						Shield:    proto.ShamanShield_WaterShield,
						Bloodlust: true,
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
			Consumes: &proto.Consumes{},
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
						InnervateTarget: &proto.UnitReference{
							Type:  proto.UnitReference_Player,
							Index: 6,
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
			Consumes: &proto.Consumes{},
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
			Consumes: &proto.Consumes{},
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
									Rotation: &proto.EnhancementShaman_Rotation{},
									Options: &proto.EnhancementShaman_Options{
										Shield:    proto.ShamanShield_LightningShield,
										Bloodlust: true,
										SyncType:  proto.ShamanSyncType_SyncMainhandOffhandSwings,
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
		},
		{
			Id: 45243,
		},
		{
			Id:      46165,
			Enchant: 3810,
		},
		{
			Id:      45242,
			Enchant: 3722,
		},
		{
			Id:      46168,
			Enchant: 1144,
		},
		{
			Id:      45446,
			Enchant: 2332,
		},
		{
			Id:      45665,
			Enchant: 3604,
		},
		{
			Id:      45619,
			Enchant: 3601,
		},
		{
			Id:      46170,
			Enchant: 3719,
		},
		{
			Id:      45135,
			Enchant: 3606,
		},
		{
			Id: 45495,
		},
		{
			Id: 46046,
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
		},
		{
			Id: 45617,
		},
		{
			Id: 45294,
		},
	},
}

// Arcane Equipment
var ArcaneEquipment = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{
			Id:      45497,
			Enchant: 3820,
		},
		{
			Id: 45243,
		},
		{
			Id:      46134,
			Enchant: 3810,
		},
		{
			Id:      45618,
			Enchant: 3722,
		},
		{
			Id:      46130,
			Enchant: 3832,
		},
		{
			Id:      45446,
			Enchant: 2332,
		},
		{
			Id:      45665,
			Enchant: 3604,
		},
		{
			Id: 45619,
		},
		{
			Id:      45488,
			Enchant: 3719,
		},
		{
			Id:      45135,
			Enchant: 3606,
		},
		{
			Id: 46046,
		},
		{
			Id: 45495,
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
		},
		{
			Id: 45617,
		},
		{
			Id:      45294,
			Enchant: 0,
		},
	},
}

// Moonkin Equipment
var MoonkinEquipment = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{
			Id:      45497,
			Enchant: 3820,
		},
		{
			Id: 45133,
		},
		{
			Id:      46196,
			Enchant: 3810,
		},
		{
			Id:      45242,
			Enchant: 3859,
		},
		{
			Id:      45519,
			Enchant: 3832,
		},
		{
			Id:      45446,
			Enchant: 2332,
		},
		{
			Id:      45665,
			Enchant: 3604,
		},
		{
			Id: 45619,
		},
		{
			Id:      46192,
			Enchant: 3719,
		},
		{
			Id:      45537,
			Enchant: 3606,
		},
		{
			Id: 46046,
		},
		{
			Id: 45495,
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
		},
		{
			Id: 45517,
		},
		{
			Id:      46203,
			Enchant: 3808,
		},
		{
			Id:      45461,
			Enchant: 3831,
		},
		{
			Id:      46205,
			Enchant: 3832,
		},
		{
			Id:      45460,
			Enchant: 3845,
		},
		{
			Id:      46200,
			Enchant: 3604,
		},
		{
			Id:      45553,
			Enchant: 0,
		},
		{
			Id:      46208,
			Enchant: 3823,
		},
		{
			Id:      45989,
			Enchant: 3606,
		},
		{
			Id: 45456,
		},
		{
			Id: 46046,
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
		},
		{
			Id:      46097,
			Enchant: 3789,
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
		},
		{
			Id: 45933,
		},
		{
			Id:      46211,
			Enchant: 3810,
		},
		{
			Id:      45242,
			Enchant: 3722,
		},
		{
			Id:      46206,
			Enchant: 3832,
		},
		{
			Id:      45460,
			Enchant: 2332,
		},
		{
			Id:      45665,
			Enchant: 3604,
		},
		{
			Id:      45616,
			Enchant: 3599,
		},
		{
			Id:      46210,
			Enchant: 3721,
		},
		{
			Id:      45537,
			Enchant: 3606,
		},
		{
			Id: 46046,
		},
		{
			Id: 45495,
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
		},
		{
			Id:      45470,
			Enchant: 1128,
		},
		{
			Id: 40267,
		},
	},
}
