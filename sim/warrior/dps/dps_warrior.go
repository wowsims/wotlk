package dps

import (
	"time"

	"github.com/wowsims/classic/sod/sim/common"
	"github.com/wowsims/classic/sod/sim/core"
	"github.com/wowsims/classic/sod/sim/core/proto"
	"github.com/wowsims/classic/sod/sim/warrior"
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

	Options        *proto.Warrior_Options
	Rotation       *proto.Warrior_Rotation
	CustomRotation *common.CustomRotation

	// Prevent swapping stances until this time, to account for human reaction time.
	canSwapStanceAt time.Duration
	// Last time sunder was applied. Used for maintaining sunder even if sunder is enabled as debuff in individual sim
	lastSunderAt time.Duration

	maintainSunder  bool
	thunderClapNext bool

	castSlamAt time.Duration
}

func NewDpsWarrior(character *core.Character, options *proto.Player) *DpsWarrior {
	warOptions := options.GetWarrior()

	war := &DpsWarrior{
		Warrior: warrior.NewWarrior(character, options.TalentsString, warrior.WarriorInputs{
			ShoutType: warOptions.Options.Shout,
		}),
		Rotation: warOptions.Rotation,
		Options:  warOptions.Options,
	}

	rbo := core.RageBarOptions{
		StartingRage:   warOptions.Options.StartingRage,
		RageMultiplier: core.TernaryFloat64(war.HasRune(proto.WarriorRune_RuneEndlessRage), 1.25, 1),
	}
	if mh := war.GetMHWeapon(); mh != nil {
		rbo.MHSwingSpeed = mh.SwingSpeed
	}
	if oh := war.GetOHWeapon(); oh != nil {
		rbo.OHSwingSpeed = oh.SwingSpeed
	}

	war.EnableRageBar(rbo, func(sim *core.Simulation) {})
	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultMeleeCritMultiplier()),
		OffHand:        war.WeaponFromOffHand(war.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	return war
}

func (cat *DpsWarrior) OnGCDReady(sim *core.Simulation) {
	return
}

func (war *DpsWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *DpsWarrior) Initialize() {
	war.Warrior.Initialize()

	war.RegisterHSOrCleave(war.Rotation.UseCleave)
	war.RegisterRendSpell()

	if war.Options.UseRecklessness {
		war.RegisterRecklessnessCD()
	}

	// This makes the behavior of these options more intuitive in the individual sim.
	if war.Env.Raid.Size() == 1 {
		if war.Rotation.SunderArmor == proto.Warrior_Rotation_SunderArmorHelpStack {
			war.SunderArmorAuras.Get(war.CurrentTarget).Duration = core.NeverExpires
		} else if war.Rotation.SunderArmor == proto.Warrior_Rotation_SunderArmorMaintain {
			war.SunderArmorAuras.Get(war.CurrentTarget).Duration = time.Second * 30
		}
	}

	if war.Rotation.StanceOption == proto.Warrior_Rotation_DefaultStance {
		if war.Warrior.PrimaryTalentTree == warrior.FuryTree {
			war.Rotation.StanceOption = proto.Warrior_Rotation_BerserkerStance
		} else {
			war.Rotation.StanceOption = proto.Warrior_Rotation_BattleStance
		}
	}

	if war.Rotation.StanceOption == proto.Warrior_Rotation_BerserkerStance {
		war.BerserkerStanceAura.BuildPhase = core.CharacterBuildPhaseTalents
	} else if war.Rotation.StanceOption == proto.Warrior_Rotation_BattleStance {
		war.BattleStanceAura.BuildPhase = core.CharacterBuildPhaseTalents
	}

	war.DelayDPSCooldownsForArmorDebuffs(time.Second * 10)
}

func (war *DpsWarrior) Reset(sim *core.Simulation) {
	if war.Rotation.StanceOption == proto.Warrior_Rotation_BerserkerStance {
		war.Warrior.Reset(sim)
		war.BerserkerStanceAura.Activate(sim)
		war.Stance = warrior.BerserkerStance
	} else if war.Rotation.StanceOption == proto.Warrior_Rotation_BattleStance {
		war.Warrior.Reset(sim)
		war.BattleStanceAura.Activate(sim)
		war.Stance = warrior.BattleStance
	}
	war.canSwapStanceAt = 0
	war.maintainSunder = war.Rotation.SunderArmor != proto.Warrior_Rotation_SunderArmorNone
	war.castSlamAt = 0
	war.thunderClapNext = false
}
