package tank

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/druid"
)

func RegisterFeralTankDruid() {
	core.RegisterAgentFactory(
		proto.Player_FeralTankDruid{},
		proto.Spec_SpecFeralTankDruid,
		func(character core.Character, options *proto.Player) core.Agent {
			return NewFeralTankDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FeralTankDruid)
			if !ok {
				panic("Invalid spec value for Feral Tank Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFeralTankDruid(character core.Character, options *proto.Player) *FeralTankDruid {
	tankOptions := options.GetFeralTankDruid()
	selfBuffs := druid.SelfBuffs{}

	bear := &FeralTankDruid{
		Druid:    druid.New(character, druid.Bear, selfBuffs, options.TalentsString),
		Rotation: tankOptions.Rotation,
		Options:  tankOptions.Options,
	}

	bear.SelfBuffs.InnervateTarget = &proto.UnitReference{}
	if tankOptions.Options.InnervateTarget != nil {
		bear.SelfBuffs.InnervateTarget = tankOptions.Options.InnervateTarget
	}

	rbo := core.RageBarOptions{
		StartingRage:   bear.Options.StartingRage,
		RageMultiplier: 1,
		MHSwingSpeed:   2.5,
	}

	bear.EnableRageBar(rbo, func(sim *core.Simulation) {
		if bear.GCD.IsReady(sim) {
			bear.TryUseCooldowns(sim)
			if bear.GCD.IsReady(sim) {
				bear.doRotation(sim)
			}
		}
	})

	bear.EnableAutoAttacks(bear, core.AutoAttackOptions{
		// Base paw weapon.
		MainHand:       bear.GetBearWeapon(),
		AutoSwingMelee: true,
		ReplaceMHSwing: func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			return bear.TryMaul(sim, mhSwingSpell)
		},
	})
	bear.ReplaceBearMHFunc = bear.TryMaul

	healingModel := options.HealingModel
	if healingModel != nil {
		if healingModel.InspirationUptime > 0.0 {
			core.ApplyInspiration(bear.GetCharacter(), healingModel.InspirationUptime)
		}
	}

	return bear
}

type FeralTankDruid struct {
	*druid.Druid

	Rotation *proto.FeralTankDruid_Rotation
	Options  *proto.FeralTankDruid_Options
}

func (bear *FeralTankDruid) GetDruid() *druid.Druid {
	return bear.Druid
}

func (bear *FeralTankDruid) Initialize() {
	bear.Druid.Initialize()
	bear.RegisterFeralTankSpells(float64(bear.Rotation.MaulRageThreshold))
}

func (bear *FeralTankDruid) Reset(sim *core.Simulation) {
	bear.Druid.Reset(sim)
	bear.Druid.ClearForm(sim)
	bear.BearFormAura.Activate(sim)
	bear.Druid.PseudoStats.Stunned = false
}
