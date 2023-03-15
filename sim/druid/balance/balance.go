package balance

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/druid"
)

func RegisterBalanceDruid() {
	core.RegisterAgentFactory(
		proto.Player_BalanceDruid{},
		proto.Spec_SpecBalanceDruid,
		func(character core.Character, options *proto.Player) core.Agent {
			return NewBalanceDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_BalanceDruid)
			if !ok {
				panic("Invalid spec value for Balance Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewBalanceDruid(character core.Character, options *proto.Player) *BalanceDruid {
	balanceOptions := options.GetBalanceDruid()
	selfBuffs := druid.SelfBuffs{}

	moonkin := &BalanceDruid{
		Druid:    druid.New(character, druid.Moonkin, selfBuffs, options.TalentsString),
		Rotation: balanceOptions.Rotation,
	}

	moonkin.SelfBuffs.InnervateTarget = &proto.RaidTarget{TargetIndex: -1}
	if balanceOptions.Options.InnervateTarget != nil {
		moonkin.SelfBuffs.InnervateTarget = balanceOptions.Options.InnervateTarget
	}

	moonkin.EnableResumeAfterManaWait(moonkin.tryUseGCD)
	return moonkin
}

type BalanceOnUseTrinket struct {
	Cooldown *core.MajorCooldown
	Stat     stats.Stat
}

type BalanceDruid struct {
	*druid.Druid

	Rotation           *proto.BalanceDruid_Rotation
	CooldownsAvailable []*core.MajorCooldown
	LastCast           *core.Spell

	// CDS
	hyperSpeedMCD      *core.MajorCooldown
	potionSpeedMCD     *core.MajorCooldown
	potionWildMagicMCD *core.MajorCooldown
	powerInfusion      *core.MajorCooldown
	onUseTrinket1      BalanceOnUseTrinket
	onUseTrinket2      BalanceOnUseTrinket
	potionUsed         bool
}

func (moonkin *BalanceDruid) GetDruid() *druid.Druid {
	return moonkin.Druid
}

func (moonkin *BalanceDruid) Initialize() {
	moonkin.Druid.Initialize()
	moonkin.RegisterBalanceSpells()
}

func (moonkin *BalanceDruid) Reset(sim *core.Simulation) {
	moonkin.Druid.Reset(sim)
	moonkin.RebirthTiming = moonkin.Env.BaseDuration.Seconds() * sim.RandomFloat("Rebirth Timing")

	if moonkin.Rotation.Type == proto.BalanceDruid_Rotation_Default {
		moonkin.Rotation.MfUsage = proto.BalanceDruid_Rotation_BeforeLunar
		moonkin.Rotation.IsUsage = proto.BalanceDruid_Rotation_MaximizeIs
		moonkin.Rotation.WrathUsage = proto.BalanceDruid_Rotation_RegularWrath
		moonkin.Rotation.UseBattleRes = false
		moonkin.Rotation.UseStarfire = true
		moonkin.Rotation.UseTyphoon = false
		moonkin.Rotation.UseHurricane = false
		moonkin.Rotation.UseSmartCooldowns = true
		moonkin.Rotation.MaintainFaerieFire = true
		moonkin.Rotation.PlayerLatency = 200
	}

	if moonkin.Rotation.UseSmartCooldowns {
		moonkin.potionUsed = false
		consumes := moonkin.Consumes

		if consumes.DefaultPotion == proto.Potions_PotionOfSpeed {
			moonkin.potionSpeedMCD = moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: 40211})
		}
		if consumes.DefaultPotion == proto.Potions_PotionOfWildMagic {
			moonkin.potionWildMagicMCD = moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: 40212})
		}
		if moonkin.HasProfession(proto.Profession_Engineering) {
			moonkin.hyperSpeedMCD = moonkin.getBalanceMajorCooldown(core.ActionID{SpellID: 54758})
		}
		moonkin.powerInfusion = moonkin.getBalanceMajorCooldown(core.ActionID{SpellID: 10060})
		moonkin.onUseTrinket1 = BalanceOnUseTrinket{
			Cooldown: moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: moonkin.Equip[core.ItemSlotTrinket1].ID}),
			Stat:     getOnUseTrinketStat(moonkin.Equip[core.ItemSlotTrinket1].ID),
		}
		moonkin.onUseTrinket2 = BalanceOnUseTrinket{
			Cooldown: moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: moonkin.Equip[core.ItemSlotTrinket2].ID}),
			Stat:     getOnUseTrinketStat(moonkin.Equip[core.ItemSlotTrinket2].ID),
		}
	}
}

// Takes out a Cooldown from the generic MajorCooldownManager and adds it to a custom Slice of Cooldowns
func (moonkin *BalanceDruid) getBalanceMajorCooldown(actionID core.ActionID) *core.MajorCooldown {
	if majorCd := moonkin.Character.GetMajorCooldownIgnoreTag(actionID); majorCd != nil {
		majorCd.Disable()
		return majorCd
	}
	return nil
}

func getOnUseTrinketStat(itemId int32) stats.Stat {
	if itemId == 45466 || itemId == 48722 || itemId == 47726 || itemId == 47946 || itemId == 36972 {
		return stats.SpellHaste
	}
	if itemId == 50259 {
		return stats.SpellCrit
	}
	return stats.SpellPower
}
