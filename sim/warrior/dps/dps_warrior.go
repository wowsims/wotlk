package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/warrior"
)

func RegisterDpsWarrior() {
	core.RegisterAgentFactory(
		proto.Player_Warrior{},
		proto.Spec_SpecWarrior,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewDpsWarrior(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Warrior)
			if !ok {
				panic("Invalid spec value for Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

type DpsWarrior struct {
	*warrior.Warrior

	Options *proto.Warrior_Options
}

func NewDpsWarrior(character *core.Character, options *proto.Player) *DpsWarrior {
	warOptions := options.GetWarrior()

	war := &DpsWarrior{
		Warrior: warrior.NewWarrior(character, options.TalentsString, warrior.WarriorInputs{
			StanceSnapshot: warOptions.Options.StanceSnapshot,
		}),
		Options: warOptions.Options,
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

	war.EnableRageBar(rbo)
	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultMeleeCritMultiplier()),
		OffHand:        war.WeaponFromOffHand(war.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		ReplaceMHSwing: war.TryHSOrCleave,
	})

	return war
}

func (war *DpsWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *DpsWarrior) Initialize() {
	war.Warrior.Initialize()

	if war.Options.UseRecklessness {
		war.RegisterRecklessnessCD()
	}

	if war.Options.UseShatteringThrow {
		war.RegisterShatteringThrowCD()
	}

	if war.PrimaryTalentTree == warrior.FuryTree {
		war.BerserkerStanceAura.BuildPhase = core.CharacterBuildPhaseTalents
	} else if war.PrimaryTalentTree == warrior.ArmsTree {
		war.BattleStanceAura.BuildPhase = core.CharacterBuildPhaseTalents
	}
}

func (war *DpsWarrior) Reset(sim *core.Simulation) {
	if war.PrimaryTalentTree == warrior.FuryTree {
		war.Warrior.Reset(sim)
		war.BerserkerStanceAura.Activate(sim)
		war.Stance = warrior.BerserkerStance
	} else if war.PrimaryTalentTree == warrior.ArmsTree {
		war.Warrior.Reset(sim)
		war.BattleStanceAura.Activate(sim)
		war.Stance = warrior.BattleStance
	}
}
