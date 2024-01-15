package enhancement

import (
	"time"

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
		Shield:  enhOptions.Options.Shield,
		ImbueMH: enhOptions.Options.ImbueMh,
		ImbueOH: enhOptions.Options.ImbueOh,
	}

	totems := &proto.ShamanTotems{}
	if enhOptions.Options.Totems != nil {
		totems = enhOptions.Options.Totems
	}

	enh := &EnhancementShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, totems, selfBuffs, true),
	}

	// Enable Auto Attacks for this spec
	enh.EnableAutoAttacks(enh, core.AutoAttackOptions{
		MainHand:       enh.WeaponFromMainHand(enh.DefaultMeleeCritMultiplier()),
		OffHand:        enh.WeaponFromOffHand(enh.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	enh.ApplySyncType(enhOptions.Options.SyncType)
	enh.ApplyFlametongueImbue(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeapon), false)
	enh.ApplyFlametongueImbue(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeaponDownrank), true)

	if !enh.HasMHWeapon() {
		enh.SelfBuffs.ImbueMH = proto.ShamanImbue_NoImbue
	}

	if !enh.HasOHWeapon() {
		enh.SelfBuffs.ImbueOH = proto.ShamanImbue_NoImbue
	}

	enh.SpiritWolves = &shaman.SpiritWolves{
		SpiritWolf1: enh.NewSpiritWolf(1),
		SpiritWolf2: enh.NewSpiritWolf(2),
	}

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
}

func (enh *EnhancementShaman) GetShaman() *shaman.Shaman {
	return enh.Shaman
}

func (enh *EnhancementShaman) Initialize() {
	enh.Shaman.Initialize()
	// In the Initialize due to frost brand adding the aura to the enemy
	enh.RegisterFrostbrandImbue(enh.getImbueProcMask(proto.ShamanImbue_FrostbrandWeapon))
	enh.RegisterFlametongueImbue(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeapon), false)
	enh.RegisterFlametongueImbue(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeaponDownrank), true)
	enh.RegisterWindfuryImbue(enh.getImbueProcMask(proto.ShamanImbue_WindfuryWeapon))

	if enh.ItemSwap.IsEnabled() {
		mh := enh.ItemSwap.GetItem(proto.ItemSlot_ItemSlotMainHand)
		enh.ApplyFlametongueImbueToItem(mh, true)
		oh := enh.ItemSwap.GetItem(proto.ItemSlot_ItemSlotOffHand)
		enh.ApplyFlametongueImbueToItem(oh, false)
		enh.RegisterOnItemSwap(func(_ *core.Simulation) {
			enh.ApplySyncType(proto.ShamanSyncType_Auto)
		})
	}
}

func (enh *EnhancementShaman) Reset(sim *core.Simulation) {
	enh.Shaman.Reset(sim)
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
		enh.AutoAttacks.SetReplaceMHSwing(func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if aa := &enh.AutoAttacks; aa.OffhandSwingAt()-sim.CurrentTime > FlurryICD {
				if nextMHSwingAt := sim.CurrentTime + aa.MainhandSwingSpeed(); nextMHSwingAt > aa.OffhandSwingAt() {
					aa.SetOffhandSwingAt(nextMHSwingAt)
				}
			}
			return mhSwingSpell
		})
	case proto.ShamanSyncType_DelayOffhandSwings:
		enh.AutoAttacks.SetReplaceMHSwing(func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if aa := &enh.AutoAttacks; aa.OffhandSwingAt()-sim.CurrentTime > FlurryICD {
				if nextMHSwingAt := sim.CurrentTime + aa.MainhandSwingSpeed() + 100*time.Millisecond; nextMHSwingAt > aa.OffhandSwingAt() {
					aa.SetOffhandSwingAt(nextMHSwingAt)
				}
			}
			return mhSwingSpell
		})
	default:
		enh.AutoAttacks.SetReplaceMHSwing(nil)
	}
}
