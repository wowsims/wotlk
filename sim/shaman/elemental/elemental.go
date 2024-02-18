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
		Shield: eleShamOptions.Options.Shield,
	}

	totems := &proto.ShamanTotems{}
	if eleShamOptions.Options.Totems != nil {
		totems = eleShamOptions.Options.Totems
	}

	inRange := eleShamOptions.Options.ThunderstormRange == proto.ElementalShaman_Options_TSInRange
	ele := &ElementalShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, totems, selfBuffs, inRange),
	}

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
}

func (eleShaman *ElementalShaman) GetShaman() *shaman.Shaman {
	return eleShaman.Shaman
}

func (eleShaman *ElementalShaman) Reset(sim *core.Simulation) {
	eleShaman.Shaman.Reset(sim)
}
