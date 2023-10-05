package elemental

import (
	"github.com/wowsims/wotlk/sim/common/wotlk"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/shaman"
)

func RegisterElementalShaman() {
	core.RegisterAgentFactory(
		proto.Player_ElementalShaman{},
		proto.Spec_SpecElementalShaman,
		func(character *core.Character, options *proto.Player) core.Agent {
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

func NewElementalShaman(character *core.Character, options *proto.Player) *ElementalShaman {
	eleShamOptions := options.GetElementalShaman()

	selfBuffs := shaman.SelfBuffs{
		Bloodlust: eleShamOptions.Options.Bloodlust,
		Shield:    eleShamOptions.Options.Shield,
	}

	totems := &proto.ShamanTotems{}
	if eleShamOptions.Options.Totems != nil {
		totems = eleShamOptions.Options.Totems
		totems.UseFireMcd = true // Control fire totems as MCD.
	}

	var rotation Rotation

	switch eleShamOptions.Rotation.Type {
	case proto.ElementalShaman_Rotation_Adaptive:
		rotation = NewAdaptiveRotation(eleShamOptions.Rotation)
	case proto.ElementalShaman_Rotation_Manual:
		rotation = NewManualRotation(eleShamOptions.Rotation)
	default:
		rotation = NewAdaptiveRotation(eleShamOptions.Rotation)
	}

	ele := &ElementalShaman{
		Shaman:   shaman.NewShaman(character, options.TalentsString, totems, selfBuffs, eleShamOptions.Rotation.InThunderstormRange),
		rotation: rotation,
		has4pT6:  character.HasSetBonus(shaman.ItemSetSkyshatterRegalia, 4),
	}
	ele.EnableResumeAfterManaWait(ele.tryUseGCD)

	if mh := ele.GetMHWeapon(); mh != nil {
		ele.ApplyFlametongueImbueToItem(mh, false)
	}

	if oh := ele.GetOHWeapon(); oh != nil {
		ele.ApplyFlametongueImbueToItem(oh, false)
	}

	if ele.Talents.FeralSpirit {
		// Enable Auto Attacks but don't enable auto swinging
		ele.EnableAutoAttacks(ele, core.AutoAttackOptions{
			MainHand: ele.WeaponFromMainHand(ele.DefaultMeleeCritMultiplier()),
			OffHand:  ele.WeaponFromOffHand(ele.DefaultMeleeCritMultiplier()),
		})
		ele.SpiritWolves = &shaman.SpiritWolves{
			SpiritWolf1: ele.NewSpiritWolf(1),
			SpiritWolf2: ele.NewSpiritWolf(2),
		}
	}

	wotlk.ConstructValkyrPets(&ele.Character)
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
