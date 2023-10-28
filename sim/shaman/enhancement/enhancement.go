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
		func(character *core.Character, options *proto.Player) core.Agent {
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

func NewEnhancementShaman(character *core.Character, options *proto.Player) *EnhancementShaman {
	enhOptions := options.GetEnhancementShaman()

	selfBuffs := shaman.SelfBuffs{
		Bloodlust: enhOptions.Options.Bloodlust,
		Shield:    enhOptions.Options.Shield,
		ImbueMH:   enhOptions.Options.ImbueMh,
		ImbueOH:   enhOptions.Options.ImbueOh,
	}

	// Override with new rotation option bloodlust.
	if enhOptions.Rotation.Bloodlust != proto.EnhancementShaman_Rotation_UnsetBloodlust {
		selfBuffs.Bloodlust = enhOptions.Rotation.Bloodlust == proto.EnhancementShaman_Rotation_UseBloodlust
	}

	totems := &proto.ShamanTotems{}
	if enhOptions.Options.Totems != nil {
		totems = enhOptions.Options.Totems
	}

	enh := &EnhancementShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, totems, selfBuffs, true),
	}

	enh.EnableResumeAfterManaWait(enh.OnGCDReady)
	enh.rotation = NewPriorityRotation(enh, enhOptions.Rotation)

	// Enable Auto Attacks for this spec
	enh.EnableAutoAttacks(enh, core.AutoAttackOptions{
		MainHand:       enh.WeaponFromMainHand(enh.DefaultMeleeCritMultiplier()),
		OffHand:        enh.WeaponFromOffHand(enh.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	enh.ApplySyncType(enhOptions.Options.SyncType)

	if enh.Totems.UseFireElemental && enhOptions.Rotation.EnableItemSwap {
		enh.EnableItemSwap(enhOptions.Rotation.ItemSwap, enh.DefaultMeleeCritMultiplier(), enh.DefaultMeleeCritMultiplier(), 0)
	}

	if enhOptions.Rotation.LightningboltWeave {
		enh.maelstromWeaponMinStack = enhOptions.Rotation.MaelstromweaponMinStack
	} else {
		enh.maelstromWeaponMinStack = 5
	}

	if !enh.HasMHWeapon() {
		enh.SelfBuffs.ImbueMH = proto.ShamanImbue_NoImbue
	}

	if !enh.HasOHWeapon() {
		enh.SelfBuffs.ImbueOH = proto.ShamanImbue_NoImbue
	}

	enh.RegisterFlametongueImbue(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeapon), false)
	enh.RegisterFlametongueImbue(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeaponDownrank), true)
	enh.RegisterWindfuryImbue(enh.getImbueProcMask(proto.ShamanImbue_WindfuryWeapon))

	enh.SpiritWolves = &shaman.SpiritWolves{
		SpiritWolf1: enh.NewSpiritWolf(1),
		SpiritWolf2: enh.NewSpiritWolf(2),
	}

	enh.ShamanisticRageManaThreshold = enhOptions.Rotation.ShamanisticRageManaThreshold

	return enh
}

func (enh *EnhancementShaman) getImbueProcMask(imbue proto.ShamanImbue) core.ProcMask {
	var mask core.ProcMask
	if enh.SelfBuffs.ImbueMH == imbue {
		mask |= core.ProcMaskMeleeMH
	}
	if enh.SelfBuffs.ImbueOH == imbue {
		mask |= core.ProcMaskMeleeOH
	}
	return mask
}

type EnhancementShaman struct {
	*shaman.Shaman

	rotation                Rotation
	maelstromWeaponMinStack int32

	// for weaving Lava Burst or Lightning Bolt
	previousSwingAt time.Duration

	scheduler common.GCDScheduler
}

func (enh *EnhancementShaman) GetShaman() *shaman.Shaman {
	return enh.Shaman
}

func (enh *EnhancementShaman) Initialize() {
	enh.Shaman.Initialize()
	// In the Initialize due to frost brand adding the aura to the enemy
	enh.RegisterFrostbrandImbue(enh.getImbueProcMask(proto.ShamanImbue_FrostbrandWeapon))

	if enh.ItemSwap.IsEnabled() {
		mh := enh.ItemSwap.GetItem(proto.ItemSlot_ItemSlotMainHand)
		enh.ApplyFlametongueImbueToItem(mh, true)
		oh := enh.ItemSwap.GetItem(proto.ItemSlot_ItemSlotOffHand)
		enh.ApplyFlametongueImbueToItem(oh, false)
		enh.RegisterOnItemSwap(func(_ *core.Simulation) {
			enh.ApplySyncType(proto.ShamanSyncType_Auto)
		})
	}
	enh.DelayDPSCooldowns(3 * time.Second)
}

func (enh *EnhancementShaman) Reset(sim *core.Simulation) {
	enh.previousSwingAt = 0
	enh.Shaman.Reset(sim)
	enh.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand}, false)
}

func (enh *EnhancementShaman) AutoSyncWeapons() proto.ShamanSyncType {
	if mh, oh := enh.MainHand(), enh.OffHand(); mh.SwingSpeed != oh.SwingSpeed {
		return proto.ShamanSyncType_NoSync
	}
	return proto.ShamanSyncType_SyncMainhandOffhandSwings
}

func (enh *EnhancementShaman) ApplySyncType(syncType proto.ShamanSyncType) {
	const FlurryICD = time.Millisecond * 500

	if syncType == proto.ShamanSyncType_Auto {
		syncType = enh.AutoSyncWeapons()
	}

	switch syncType {
	case proto.ShamanSyncType_SyncMainhandOffhandSwings:
		enh.AutoAttacks.ReplaceMHSwing = func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if aa := &enh.AutoAttacks; aa.OffhandSwingAt()-sim.CurrentTime > FlurryICD {
				if nextMHSwingAt := sim.CurrentTime + aa.MainhandSwingSpeed(); nextMHSwingAt > aa.OffhandSwingAt() {
					aa.SetOffhandSwingAt(nextMHSwingAt)
				}
			}
			return mhSwingSpell
		}
	case proto.ShamanSyncType_DelayOffhandSwings:
		enh.AutoAttacks.ReplaceMHSwing = func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if aa := &enh.AutoAttacks; aa.OffhandSwingAt()-sim.CurrentTime > FlurryICD {
				if nextMHSwingAt := sim.CurrentTime + aa.MainhandSwingSpeed() + 100*time.Millisecond; nextMHSwingAt > aa.OffhandSwingAt() {
					aa.SetOffhandSwingAt(nextMHSwingAt)
				}
			}
			return mhSwingSpell
		}
	default:
		enh.AutoAttacks.ReplaceMHSwing = nil
	}
}

func (enh *EnhancementShaman) CastLightningBoltWeave(sim *core.Simulation, reactionTime time.Duration) bool {
	previousAttack := sim.CurrentTime - enh.previousSwingAt
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
	previousAttack := sim.CurrentTime - enh.previousSwingAt
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
