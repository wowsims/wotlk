package items

import (
	"github.com/wowsims/wotlk/assets/database"
)

func loadDatabase() {
	db := database.Load()

	Items = make([]Item, len(db.Items))
	for i, item := range db.Items {
		Items[i] = ItemFromProto(item)
	}

	Enchants = make([]Enchant, len(db.Enchants))
	for i, enchant := range db.Enchants {
		Enchants[i] = EnchantFromProto(enchant)
	}

	Gems = make([]Gem, len(db.Gems))
	for i, gem := range db.Gems {
		Gems[i] = GemFromProto(gem)
	}
}
