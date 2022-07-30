package enhancement

import (
	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/shaman"
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
		Bloodlust: enhOptions.Options.Bloodlust,
		Shield:    enhOptions.Options.Shield,
		ImbueMH:   enhOptions.Options.ImbueMH,
		ImbueOH:   enhOptions.Options.ImbueOH,
	}

	totems := proto.ShamanTotems{}
	if enhOptions.Rotation.Totems != nil {
		totems = *enhOptions.Rotation.Totems
	}

	var rotation Rotation
	switch enhOptions.Rotation.Type {
	case proto.EnhancementShaman_Rotation_Basic:
		rotation = NewPriorityRotation(enhOptions.Talents, nil)
	case proto.EnhancementShaman_Rotation_Priority:
		rotation = NewPriorityRotation(enhOptions.Talents, enhOptions.Rotation.Order)
	}

	enh := &EnhancementShaman{
		Shaman:   shaman.NewShaman(character, *enhOptions.Talents, totems, selfBuffs, true),
		rotation: rotation,
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
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_WindfuryWeapon,
		enh.SelfBuffs.ImbueOH == proto.ShamanImbue_WindfuryWeapon)
	enh.ApplyFlametongueImbue(
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FlametongueWeapon,
		enh.SelfBuffs.ImbueOH == proto.ShamanImbue_FlametongueWeapon)
	enh.ApplyFrostbrandImbue(
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FrostbrandWeapon,
		enh.SelfBuffs.ImbueOH == proto.ShamanImbue_FrostbrandWeapon)

	if enh.SelfBuffs.ImbueMH == proto.ShamanImbue_WindfuryWeapon ||
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FlametongueWeapon ||
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FrostbrandWeapon {
		enh.HasMHWeaponImbue = true
	}

	enh.SpiritWolves = &shaman.SpiritWolves{
		SpiritWolf1: enh.NewSpiritWolf(1),
		SpiritWolf2: enh.NewSpiritWolf(2),
	}

	return enh
}

type EnhancementShaman struct {
	*shaman.Shaman

	rotation Rotation

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
	//enh.Env.RegisterPostFinalizeEffect(enh.SetupRotationSchedule)
}

func (enh *EnhancementShaman) Reset(sim *core.Simulation) {
	enh.Shaman.Reset(sim)
	enh.scheduler.Reset(sim, enh.GetCharacter())
}
