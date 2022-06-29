package enhancement

import (
	"github.com/wowsims/tbc/sim/common"
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/shaman"
)

func RegisterEnhancementShaman() {
	core.RegisterAgentFactory(
		proto.Player_EnhancementShaman{},
		proto.Spec_SpecEnhancementShaman,
		func(character core.Character, options proto.Player) core.Agent {
			return NewEnhancementShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_EnhancementShaman)
			if !ok {
				panic("Invalid spec value for Enhancement Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewEnhancementShaman(character core.Character, options proto.Player) *EnhancementShaman {
	enhOptions := options.GetEnhancementShaman()

	selfBuffs := shaman.SelfBuffs{
		Bloodlust:        enhOptions.Options.Bloodlust,
		WaterShield:      enhOptions.Options.WaterShield,
		SnapshotSOET42Pc: enhOptions.Options.SnapshotT4_2Pc,
	}

	totems := proto.ShamanTotems{}
	if enhOptions.Rotation.Totems != nil {
		totems = *enhOptions.Rotation.Totems
	}
	enh := &EnhancementShaman{
		Shaman:   shaman.NewShaman(character, *enhOptions.Talents, totems, selfBuffs),
		Rotation: *enhOptions.Rotation,
	}
	// Enable Auto Attacks for this spec
	enh.EnableAutoAttacks(enh, core.AutoAttackOptions{
		MainHand:       enh.WeaponFromMainHand(enh.DefaultMeleeCritMultiplier()),
		OffHand:        enh.WeaponFromOffHand(enh.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		DelayOHSwings:  enhOptions.Options.DelayOffhandSwings,
	})

	if !enh.HasMHWeapon() {
		enh.Consumes.MainHandImbue = proto.WeaponImbue_WeaponImbueUnknown
	}
	if !enh.HasOHWeapon() {
		enh.Consumes.OffHandImbue = proto.WeaponImbue_WeaponImbueUnknown
	}
	enh.ApplyWindfuryImbue(
		enh.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueShamanWindfury,
		enh.Consumes.OffHandImbue == proto.WeaponImbue_WeaponImbueShamanWindfury)
	enh.ApplyFlametongueImbue(
		enh.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueShamanFlametongue,
		enh.Consumes.OffHandImbue == proto.WeaponImbue_WeaponImbueShamanFlametongue)
	enh.ApplyFrostbrandImbue(
		enh.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueShamanFrostbrand,
		enh.Consumes.OffHandImbue == proto.WeaponImbue_WeaponImbueShamanFrostbrand)
	enh.ApplyRockbiterImbue(
		enh.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueShamanRockbiter,
		enh.Consumes.OffHandImbue == proto.WeaponImbue_WeaponImbueShamanRockbiter)

	if enh.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueShamanWindfury ||
		enh.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueShamanFlametongue ||
		enh.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueShamanFrostbrand ||
		enh.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueShamanRockbiter {
		enh.HasMHWeaponImbue = true
	}

	return enh
}

type EnhancementShaman struct {
	*shaman.Shaman

	Rotation proto.EnhancementShaman_Rotation

	scheduler common.GCDScheduler
}

func (enh *EnhancementShaman) GetShaman() *shaman.Shaman {
	return enh.Shaman
}

func (enh *EnhancementShaman) Initialize() {
	enh.Shaman.Initialize()
	enh.DelayDPSCooldownsForArmorDebuffs()

	// This needs to be called after DPS cooldowns are delayed, which also happens
	// after finalization.
	enh.Env.RegisterPostFinalizeEffect(enh.SetupRotationSchedule)
}

func (enh *EnhancementShaman) Reset(sim *core.Simulation) {
	enh.Shaman.Reset(sim)
	enh.scheduler.Reset(sim, enh.GetCharacter())
}
