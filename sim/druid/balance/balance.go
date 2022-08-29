package balance

import (
	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/druid"
)

func RegisterBalanceDruid() {
	core.RegisterAgentFactory(
		proto.Player_BalanceDruid{},
		proto.Spec_SpecBalanceDruid,
		func(character core.Character, options proto.Player) core.Agent {
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

func NewBalanceDruid(character core.Character, options proto.Player) *BalanceDruid {
	balanceOptions := options.GetBalanceDruid()
	selfBuffs := druid.SelfBuffs{}

	if balanceOptions.Options.InnervateTarget != nil {
		selfBuffs.InnervateTarget = *balanceOptions.Options.InnervateTarget
	} else {
		selfBuffs.InnervateTarget.TargetIndex = -1
	}

	moonkin := &BalanceDruid{
		Druid:                    druid.New(character, druid.Moonkin, selfBuffs, *balanceOptions.Talents),
		primaryRotation:          *balanceOptions.Rotation,
		useBattleRes:             balanceOptions.Options.BattleRes,
		useIS:                    balanceOptions.Options.UseIs,
		useMF:                    balanceOptions.Options.UseMf,
		mfInsideEclipseThreshold: balanceOptions.Options.MfInsideEclipseThreshold,
		isInsideEclipseThreshold: balanceOptions.Options.IsInsideEclipseThreshold,
		mcdInsideLunarThreshold:  balanceOptions.Options.McdInsideLunarThreshold - 0.5,
		mcdInsideSolarThreshold:  balanceOptions.Options.McdInsideSolarThreshold - 0.5,
	}

	moonkin.ResetTalentsBonuses()
	moonkin.RegisterTalentsBonuses()
	moonkin.EnableResumeAfterManaWait(moonkin.tryUseGCD)
	return moonkin
}

type BalanceDruid struct {
	*druid.Druid

	primaryRotation          proto.BalanceDruid_Rotation
	useBattleRes             bool
	useIS                    bool
	useMF                    bool
	mfInsideEclipseThreshold float32
	isInsideEclipseThreshold float32
	mcdInsideLunarThreshold  float32
	mcdInsideSolarThreshold  float32
	// These are only used when primary spell is set to 'Adaptive'. When the mana
	// tracker tells us we have extra mana to spare, use surplusRotation instead of
	// primaryRotation.
	useSurplusRotation bool
	surplusRotation    proto.BalanceDruid_Rotation
	manaTracker        common.ManaSpendingRateTracker
	// CDS
	hyperSpeedMCD *core.MajorCooldown
	potionMCD     *core.MajorCooldown
	onUseTrinket1 *core.MajorCooldown
	onUseTrinket2 *core.MajorCooldown
	potionUsed    bool
}

// GetDruid is to implement druid.Agent (supports nordrassil set bonus)
func (moonkin *BalanceDruid) GetDruid() *druid.Druid {
	return moonkin.Druid
}

func (moonkin *BalanceDruid) Initialize() {
	moonkin.Druid.Initialize()
	moonkin.RegisterBalanceSpells()
}

func (moonkin *BalanceDruid) Reset(sim *core.Simulation) {
	if moonkin.useSurplusRotation {
		moonkin.manaTracker.Reset()
	}
	moonkin.Druid.Reset(sim)
	moonkin.RebirthTiming = moonkin.Env.BaseDuration.Seconds() * sim.RandomFloat("Rebirth Timing")

	if moonkin.mcdInsideLunarThreshold > 0 || moonkin.mcdInsideSolarThreshold > 0 {
		moonkin.potionUsed = false
		consumes := &moonkin.Consumes
		if consumes.DefaultPotion == proto.Potions_PotionOfSpeed {
			moonkin.potionMCD = moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: 40211})
		}
		if consumes.DefaultPotion == proto.Potions_PotionOfWildMagic {
			moonkin.potionMCD = moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: 40212})
		}
		moonkin.hyperSpeedMCD = moonkin.getBalanceMajorCooldown(core.ActionID{SpellID: 54758})
		moonkin.onUseTrinket1 = moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: moonkin.Equip[items.ItemSlotTrinket1].ID})
		moonkin.onUseTrinket2 = moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: moonkin.Equip[items.ItemSlotTrinket2].ID})
	}
}

func (moonkin *BalanceDruid) getBalanceMajorCooldown(actionID core.ActionID) *core.MajorCooldown {
	if moonkin.Character.HasMajorCooldown(actionID) {
		majorCd := moonkin.Character.GetMajorCooldown(actionID)
		majorCd.ShouldActivate = func(sim *core.Simulation, character *core.Character) bool {
			return false
		}
		return majorCd
	}
	return nil
}
