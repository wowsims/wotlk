package balance

import (
	"github.com/wowsims/wotlk/sim/common/wotlk"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/druid"
)

func RegisterBalanceDruid() {
	core.RegisterAgentFactory(
		proto.Player_BalanceDruid{},
		proto.Spec_SpecBalanceDruid,
		func(character *core.Character, options *proto.Player) core.Agent {
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

func NewBalanceDruid(character *core.Character, options *proto.Player) *BalanceDruid {
	balanceOptions := options.GetBalanceDruid()
	selfBuffs := druid.SelfBuffs{}

	moonkin := &BalanceDruid{
		Druid:    druid.New(character, druid.Moonkin, selfBuffs, options.TalentsString),
		Options:  balanceOptions.Options,
		Rotation: balanceOptions.Rotation,
	}

	moonkin.SelfBuffs.InnervateTarget = &proto.UnitReference{}
	if balanceOptions.Options.InnervateTarget != nil {
		moonkin.SelfBuffs.InnervateTarget = balanceOptions.Options.InnervateTarget
	}

	moonkin.EnableResumeAfterManaWait(moonkin.tryUseGCD)
	wotlk.ConstructValkyrPets(&moonkin.Character)
	return moonkin
}

type BalanceOnUseTrinket struct {
	Cooldown *core.MajorCooldown
	Stat     stats.Stat
}

type BalanceDruid struct {
	*druid.Druid

	Options            *proto.BalanceDruid_Options
	Rotation           *proto.BalanceDruid_Rotation
	CooldownsAvailable []*core.MajorCooldown
	LastCast           *druid.DruidSpell

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
