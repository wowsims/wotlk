package tank

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func RegisterTankDeathknight() {
	core.RegisterAgentFactory(
		proto.Player_TankDeathknight{},
		proto.Spec_SpecTankDeathknight,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewTankDeathknight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_TankDeathknight)
			if !ok {
				panic("Invalid spec value for Tank Deathknight!")
			}
			player.Spec = playerSpec
		},
	)
}

type TankDeathknight struct {
	*deathknight.Deathknight
}

func NewTankDeathknight(character *core.Character, options *proto.Player) *TankDeathknight {
	dkOptions := options.GetTankDeathknight()

	tankDk := &TankDeathknight{
		Deathknight: deathknight.NewDeathknight(character, deathknight.DeathknightInputs{
			IsDps:              false,
			StartingRunicPower: dkOptions.Options.StartingRunicPower,
		}, options.TalentsString),
	}

	tankDk.Inputs.UnholyFrenzyTarget = dkOptions.Options.UnholyFrenzyTarget

	tankDk.EnableAutoAttacks(tankDk, core.AutoAttackOptions{
		MainHand:       tankDk.WeaponFromMainHand(tankDk.DefaultMeleeCritMultiplier()),
		OffHand:        tankDk.WeaponFromOffHand(tankDk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		ReplaceMHSwing: func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if tankDk.RuneStrikeQueued && tankDk.RuneStrike.CanCast(sim, nil) {
				return tankDk.RuneStrike
			} else {
				return mhSwingSpell
			}
		},
	})

	healingModel := options.HealingModel
	if healingModel != nil {
		if healingModel.InspirationUptime > 0.0 {
			core.ApplyInspiration(tankDk.GetCharacter(), healingModel.InspirationUptime)
		}
	}

	return tankDk
}

func (dk *TankDeathknight) GetDeathknight() *deathknight.Deathknight {
	return dk.Deathknight
}

func (dk *TankDeathknight) Initialize() {
	dk.Deathknight.Initialize()
}

func (dk *TankDeathknight) Reset(sim *core.Simulation) {
	dk.Deathknight.Reset(sim)

	dk.Presence = deathknight.UnsetPresence
	dk.Deathknight.PseudoStats.Stunned = false
}
