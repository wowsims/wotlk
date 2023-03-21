package enhancement

import (
	"time"

	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/shaman"
)

func RegisterEnhancementShaman() {
	core.RegisterAgentFactory(
		proto.Player_EnhancementShaman{},
		proto.Spec_SpecEnhancementShaman,
		func(character core.Character, options *proto.Player) core.Agent {
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

func NewEnhancementShaman(character core.Character, options *proto.Player) *EnhancementShaman {
	enhOptions := options.GetEnhancementShaman()

	selfBuffs := shaman.SelfBuffs{
		Bloodlust: enhOptions.Options.Bloodlust,
		Shield:    enhOptions.Options.Shield,
		ImbueMH:   enhOptions.Options.ImbueMh,
		ImbueOH:   enhOptions.Options.ImbueOh,
	}

	totems := &proto.ShamanTotems{}
	if enhOptions.Rotation.Totems != nil {
		totems = enhOptions.Rotation.Totems
	}

	enh := &EnhancementShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, totems, selfBuffs, true),
	}

	enh.EnableResumeAfterManaWait(enh.OnGCDReady)
	enh.rotation = NewPriorityRotation(enh, enhOptions.Rotation)

	syncType := int32(enhOptions.Options.SyncType)
	if syncType == int32(proto.ShamanSyncType_Auto) {
		syncType = enh.AutoSyncWeapons()
	}

	// Enable Auto Attacks for this spec
	enh.EnableAutoAttacks(enh, core.AutoAttackOptions{
		MainHand:       enh.WeaponFromMainHand(enh.DefaultMeleeCritMultiplier()),
		OffHand:        enh.WeaponFromOffHand(enh.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		SyncType:       syncType,
	})

	if enh.Totems.UseFireElemental && enhOptions.Rotation.EnableItemSwap {
		enh.EnableItemSwap(enhOptions.Rotation.ItemSwap, enh.DefaultMeleeCritMultiplier(), enh.DefaultMeleeCritMultiplier(), 0)
	}

	if enhOptions.Rotation.LightningboltWeave {
		enh.maelstromweaponMinStack = enhOptions.Rotation.MaelstromweaponMinStack
	} else {
		enh.maelstromweaponMinStack = 5
	}

	if !enh.HasMHWeapon() {
		enh.SelfBuffs.ImbueMH = proto.ShamanImbue_NoImbue
	}

	if !enh.HasOHWeapon() {
		enh.SelfBuffs.ImbueOH = proto.ShamanImbue_NoImbue
	}

	enh.RegisterFlametongueImbue(
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FlametongueWeapon,
		enh.SelfBuffs.ImbueOH == proto.ShamanImbue_FlametongueWeapon)
	enh.RegisterFlametongueDownrankImbue(
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FlametongueWeaponDownrank,
		enh.SelfBuffs.ImbueOH == proto.ShamanImbue_FlametongueWeaponDownrank)
	enh.RegisterWindfuryImbue(
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_WindfuryWeapon,
		enh.SelfBuffs.ImbueOH == proto.ShamanImbue_WindfuryWeapon)

	enh.SpiritWolves = &shaman.SpiritWolves{
		SpiritWolf1: enh.NewSpiritWolf(1),
		SpiritWolf2: enh.NewSpiritWolf(2),
	}

	enh.ShamanisticRageManaThreshold = enhOptions.Rotation.ShamanisticRageManaThreshold

	return enh
}

type EnhancementShaman struct {
	*shaman.Shaman

	rotation                Rotation
	maelstromweaponMinStack int32

	scheduler common.GCDScheduler
}

func (enh *EnhancementShaman) GetShaman() *shaman.Shaman {
	return enh.Shaman
}

func (enh *EnhancementShaman) Initialize() {
	enh.Shaman.Initialize()

	enh.RegisterFrostbrandImbue(
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FrostbrandWeapon,
		enh.SelfBuffs.ImbueOH == proto.ShamanImbue_FrostbrandWeapon)

	if enh.ItemSwap.IsEnabled() {
		mh := enh.ItemSwap.GetItem(proto.ItemSlot_ItemSlotMainHand)
		enh.ApplyFlametongueImbueToItem(mh, true)
		oh := enh.ItemSwap.GetItem(proto.ItemSlot_ItemSlotOffHand)
		enh.ApplyFlametongueImbueToItem(oh, false)
		enh.RegisterOnItemSwap(func(s *core.Simulation) {
			enh.AutoAttacks.SyncType = enh.AutoSyncWeapons()
		})
	}
	enh.DelayDPSCooldowns(3 * time.Second)
}

func (enh *EnhancementShaman) Reset(sim *core.Simulation) {
	enh.Shaman.Reset(sim)
	enh.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand}, false)
}

func (enh *EnhancementShaman) AutoSyncWeapons() int32 {
	mh := enh.GetMHWeapon()
	oh := enh.GetOHWeapon()

	if mh == nil || oh == nil || mh.SwingSpeed != oh.SwingSpeed {
		return int32(proto.ShamanSyncType_NoSync)
	}

	return int32(proto.ShamanSyncType_SyncMainhandOffhandSwings)
}

func (enh *EnhancementShaman) CastLightningBoltWeave(sim *core.Simulation, reactionTime time.Duration) bool {
	previousAttack := sim.CurrentTime - enh.AutoAttacks.PreviousSwingAt
	reactionTime = core.TernaryDuration(previousAttack < reactionTime, reactionTime-previousAttack, 0)

	//calculate cast times for weaving
	lbCastTime := enh.ApplyCastSpeed(enh.LightningBolt.DefaultCast.CastTime-(time.Millisecond*time.Duration(500*enh.MaelstromWeaponAura.GetStacks()))) + reactionTime
	//calculate swing times for weaving
	timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime

	if lbCastTime < timeUntilSwing {
		if reactionTime > 0 {
			reactionTime += sim.CurrentTime

			enh.HardcastWaitUntil(sim, reactionTime, func(_ *core.Simulation, _ *core.Unit) {
				enh.GCD.Reset()
				enh.LightningBolt.Cast(sim, enh.CurrentTarget)
			})

			enh.WaitUntil(sim, reactionTime)
			return true
		}
		return enh.LightningBolt.Cast(sim, enh.CurrentTarget)
	}

	return false
}

func (enh *EnhancementShaman) CastLavaBurstWeave(sim *core.Simulation, reactionTime time.Duration) bool {
	previousAttack := sim.CurrentTime - enh.AutoAttacks.PreviousSwingAt
	reactionTime = core.TernaryDuration(previousAttack < reactionTime, reactionTime-previousAttack, 0)

	//calculate cast times for weaving
	lvbCastTime := enh.ApplyCastSpeed(enh.LavaBurst.DefaultCast.CastTime) + reactionTime
	//calculate swing times for weaving
	timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime

	if lvbCastTime < timeUntilSwing {
		if reactionTime > 0 {
			reactionTime += sim.CurrentTime

			enh.HardcastWaitUntil(sim, reactionTime, func(_ *core.Simulation, _ *core.Unit) {
				enh.GCD.Reset()
				enh.LavaBurst.Cast(sim, enh.CurrentTarget)
			})

			enh.WaitUntil(sim, reactionTime)
			return true
		}

		return enh.LavaBurst.Cast(sim, enh.CurrentTarget)
	}

	return false
}
