package balance

import (
	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
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
	}

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
	// These are only used when primary spell is set to 'Adaptive'. When the mana
	// tracker tells us we have extra mana to spare, use surplusRotation instead of
	// primaryRotation.
	useSurplusRotation bool
	surplusRotation    proto.BalanceDruid_Rotation
	manaTracker        common.ManaSpendingRateTracker
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
}
