package elemental

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/shaman"
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
		Bloodlust: eleShamOptions.Options.Bloodlust,
		Shield:    eleShamOptions.Options.Shield,
	}

	totems := proto.ShamanTotems{}
	if eleShamOptions.Rotation.Totems != nil {
		totems = *eleShamOptions.Rotation.Totems
	}

	var rotation Rotation

	switch eleShamOptions.Rotation.Type {
	case proto.ElementalShaman_Rotation_Adaptive:
		rotation = NewAdaptiveRotation(eleShamOptions.Talents)
	}

	ele := &ElementalShaman{
		Shaman:   shaman.NewShaman(character, *eleShamOptions.Talents, totems, selfBuffs),
		rotation: rotation,
		has4pT6:  character.HasSetBonus(shaman.ItemSetSkyshatterRegalia, 4),
	}
	ele.EnableResumeAfterManaWait(ele.tryUseGCD)

	ele.ApplyFlametongueImbue(
		ele.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueShamanFlametongue,
		ele.Consumes.OffHandImbue == proto.WeaponImbue_WeaponImbueShamanFlametongue)

	return ele
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
