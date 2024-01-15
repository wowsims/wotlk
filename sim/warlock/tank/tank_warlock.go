package tank

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/warlock"
)

func RegisterTankWarlock() {
	core.RegisterAgentFactory(
		proto.Player_TankWarlock{},
		proto.Spec_SpecTankWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewTankWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_TankWarlock)
			if !ok {
				panic("Invalid spec value for Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

type TankWarlock struct {
	*warlock.Warlock
}

func NewTankWarlock(character *core.Character, options *proto.Player) *TankWarlock {
	warlock := &TankWarlock{
		Warlock: warlock.NewWarlock(character, options, options.GetTankWarlock().Options),
	}

	warlock.PseudoStats.CanParry = false
	warlock.PseudoStats.CanBlock = false

	warlock.EnableAutoAttacks(warlock, core.AutoAttackOptions{
		MainHand:       warlock.WeaponFromMainHand(warlock.DefaultMeleeCritMultiplier()),
		OffHand:        warlock.WeaponFromOffHand(warlock.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	return warlock
}

func (warlock *TankWarlock) OnGCDReady(sim *core.Simulation) {
	return
}

func (warlock *TankWarlock) GetWarlock() *warlock.Warlock {
	return warlock.Warlock
}

func (warlock *TankWarlock) Initialize() {
	warlock.Warlock.Initialize()
}

func (warlock *TankWarlock) Reset(sim *core.Simulation) {
	warlock.Warlock.Reset(sim)
}
