package protection

import (
	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/warrior"
)

func RegisterProtectionWarrior() {
	core.RegisterAgentFactory(
		proto.Player_ProtectionWarrior{},
		proto.Spec_SpecProtectionWarrior,
		func(character core.Character, options *proto.Player) core.Agent {
			return NewProtectionWarrior(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ProtectionWarrior)
			if !ok {
				panic("Invalid spec value for Protection Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

type ProtectionWarrior struct {
	*warrior.Warrior

	Rotation *proto.ProtectionWarrior_Rotation
	Options  *proto.ProtectionWarrior_Options

	CustomRotation *common.CustomRotation
}

func NewProtectionWarrior(character core.Character, options *proto.Player) *ProtectionWarrior {
	warOptions := options.GetProtectionWarrior()

	war := &ProtectionWarrior{
		Warrior: warrior.NewWarrior(character, options.TalentsString, warrior.WarriorInputs{
			ShoutType: warOptions.Options.Shout,
			Munch:     warOptions.Options.Munch,
		}),
		Rotation: warOptions.Rotation,
		Options:  warOptions.Options,
	}

	rbo := core.RageBarOptions{
		StartingRage:   warOptions.Options.StartingRage,
		RageMultiplier: core.TernaryFloat64(war.Talents.EndlessRage, 1.25, 1),
	}
	if mh := war.GetMHWeapon(); mh != nil {
		rbo.MHSwingSpeed = mh.SwingSpeed
	}
	if oh := war.GetOHWeapon(); oh != nil {
		rbo.OHSwingSpeed = oh.SwingSpeed
	}

	war.EnableRageBar(rbo, func(sim *core.Simulation) {
		if war.GCD.IsReady(sim) {
			war.TryUseCooldowns(sim)
			if war.GCD.IsReady(sim) {
				war.doRotation(sim)
			}
		}
	})
	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultMeleeCritMultiplier()),
		OffHand:        war.WeaponFromOffHand(war.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		ReplaceMHSwing: func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			return war.TryHSOrCleave(sim, mhSwingSpell)
		},
	})

	healingModel := options.HealingModel
	if healingModel != nil {
		if healingModel.InspirationUptime > 0.0 {
			core.ApplyInspiration(war.GetCharacter(), healingModel.InspirationUptime)
		}
	}

	return war
}

func (war *ProtectionWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *ProtectionWarrior) Initialize() {
	war.Warrior.Initialize()

	war.RegisterHSOrCleave(false, war.Rotation.HsRageThreshold)
	war.RegisterShieldWallCD()
	war.RegisterShieldBlockCD()
	war.DefensiveStanceAura.BuildPhase = core.CharacterBuildPhaseTalents

	war.CustomRotation = war.makeCustomRotation()
	if war.Options.UseShatteringThrow {
		war.RegisterShatteringThrowCD()
	}
}

func (war *ProtectionWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
	war.DefensiveStanceAura.Activate(sim)
	war.Stance = warrior.DefensiveStance
}
