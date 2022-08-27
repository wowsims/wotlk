package smite

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/priest"
)

func RegisterSmitePriest() {
	core.RegisterAgentFactory(
		proto.Player_SmitePriest{},
		proto.Spec_SpecSmitePriest,
		func(character core.Character, options proto.Player) core.Agent {
			return NewSmitePriest(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_SmitePriest)
			if !ok {
				panic("Invalid spec value for Smite Priest!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewSmitePriest(character core.Character, options proto.Player) *SmitePriest {
	smiteOptions := options.GetSmitePriest()

	selfBuffs := priest.SelfBuffs{
		UseInnerFire:   smiteOptions.Options.UseInnerFire,
		UseShadowfiend: smiteOptions.Options.UseShadowfiend,
	}

	if smiteOptions.Options.PowerInfusionTarget != nil {
		selfBuffs.PowerInfusionTarget = *smiteOptions.Options.PowerInfusionTarget
	} else {
		selfBuffs.PowerInfusionTarget.TargetIndex = -1
	}

	basePriest := priest.New(character, selfBuffs, *smiteOptions.Talents)

	spriest := &SmitePriest{
		Priest:   basePriest,
		rotation: *smiteOptions.Rotation,

		allowedHFDelay: time.Millisecond * time.Duration(smiteOptions.Rotation.AllowedHolyFireDelayMs),
	}

	spriest.EnableResumeAfterManaWait(spriest.tryUseGCD)

	return spriest
}

type SmitePriest struct {
	*priest.Priest

	rotation proto.SmitePriest_Rotation

	allowedHFDelay time.Duration
}

func (spriest *SmitePriest) GetPriest() *priest.Priest {
	return spriest.Priest
}

func (spriest *SmitePriest) Initialize() {
	spriest.Priest.Initialize()

	spriest.RegisterHolyFireSpell(spriest.rotation.MemeDream)
	spriest.RegisterSmiteSpell(spriest.rotation.MemeDream)
	spriest.RegisterPenanceSpell()
	spriest.RegisterHymnOfHopeCD()
}

func (spriest *SmitePriest) Reset(sim *core.Simulation) {
	spriest.Priest.Reset(sim)
}
