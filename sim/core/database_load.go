// Only include this file in the build when we specify the 'with_db' tag.
// Without the tag, the database will start out completely empty.
//go:build with_db

package core

import (
	"github.com/wowsims/wotlk/assets/database"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	db := database.Load()
	WITH_DB = true

	simDB := &proto.SimDatabase{
		Items:    make([]*proto.SimItem, len(db.Items)),
		Enchants: make([]*proto.SimEnchant, len(db.Enchants)),
		Gems:     make([]*proto.SimGem, len(db.Gems)),
	}

	for i, item := range db.Items {
		simDB.Items[i] = &proto.SimItem{
			Id:               item.Id,
			Name:             item.Name,
			Type:             item.Type,
			ArmorType:        item.ArmorType,
			WeaponType:       item.WeaponType,
			HandType:         item.HandType,
			RangedWeaponType: item.RangedWeaponType,
			Stats:            item.Stats,
			GemSockets:       item.GemSockets,
			SocketBonus:      item.SocketBonus,
			WeaponDamageMin:  item.WeaponDamageMin,
			WeaponDamageMax:  item.WeaponDamageMax,
			WeaponSpeed:      item.WeaponSpeed,
			SetName:          item.SetName,
		}
	}

	for i, enchant := range db.Enchants {
		simDB.Enchants[i] = &proto.SimEnchant{
			EffectId: enchant.EffectId,
			Stats:    enchant.Stats,
		}
	}

	for i, gem := range db.Gems {
		simDB.Gems[i] = &proto.SimGem{
			Id:    gem.Id,
			Name:  gem.Name,
			Color: gem.Color,
			Stats: gem.Stats,
		}
	}

	addToDatabase(simDB)
}
