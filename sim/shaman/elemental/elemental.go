package elemental

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/shaman"
)

func RegisterElementalShaman() {
	core.RegisterAgentFactory(
		proto.Player_ElementalShaman{},
		proto.Spec_SpecElementalShaman,
		func(character core.Character, options proto.Player) core.Agent {
			return NewElementalShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ElementalShaman)
			if !ok {
				panic("Invalid spec value for Elemental Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewElementalShaman(character core.Character, options proto.Player) *ElementalShaman {
	eleShamOptions := options.GetElementalShaman()

	selfBuffs := shaman.SelfBuffs{
		Bloodlust:        eleShamOptions.Options.Bloodlust,
		WaterShield:      eleShamOptions.Options.WaterShield,
		SnapshotWOAT42Pc: eleShamOptions.Options.SnapshotT4_2Pc,
	}

	totems := proto.ShamanTotems{}
	if eleShamOptions.Rotation.Totems != nil {
		totems = *eleShamOptions.Rotation.Totems
	}

	var rotation Rotation

	switch eleShamOptions.Rotation.Type {
	case proto.ElementalShaman_Rotation_Adaptive:
		rotation = NewAdaptiveRotation(eleShamOptions.Talents)
	case proto.ElementalShaman_Rotation_CLOnClearcast:
		if eleShamOptions.Talents.ElementalFocus {
			rotation = NewCLOnClearcastRotation()
		} else {
			rotation = NewCLOnCDRotation()
		}
	case proto.ElementalShaman_Rotation_CLOnCD:
		rotation = NewCLOnCDRotation()
	case proto.ElementalShaman_Rotation_FixedLBCL:
		rotation = NewFixedRotation(eleShamOptions.Rotation.LbsPerCl)
	case proto.ElementalShaman_Rotation_LBOnly:
		rotation = NewLBOnlyRotation()
	}

	return &ElementalShaman{
		Shaman:   shaman.NewShaman(character, *eleShamOptions.Talents, totems, selfBuffs),
		rotation: rotation,
		has4pT6:  shaman.ItemSetSkyshatterRegalia.CharacterHasSetBonus(&character, 4),
	}
}

type ElementalShaman struct {
	*shaman.Shaman

	rotation Rotation

	has4pT6 bool
}

func (eleShaman *ElementalShaman) GetShaman() *shaman.Shaman {
	return eleShaman.Shaman
}

func (eleShaman *ElementalShaman) Reset(sim *core.Simulation) {
	eleShaman.Shaman.Reset(sim)
	eleShaman.rotation.Reset(eleShaman, sim)
}
