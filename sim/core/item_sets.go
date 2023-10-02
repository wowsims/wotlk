package core

import (
	"fmt"
	"slices"
)

type ItemSet struct {
	Name            string
	AlternativeName string

	// Maps set piece requirement to an ApplyEffect function that will be called
	// before the Sim starts.
	//
	// The function should apply any benefits provided by the set bonus.
	Bonuses map[int32]ApplyEffect
}

func (set ItemSet) Items() []Item {
	var items []Item
	for _, item := range ItemsByID {
		if item.SetName == "" {
			continue
		}
		if item.SetName == set.Name || item.SetName == set.AlternativeName {
			items = append(items, item)
		}
	}
	// Sort so the order of IDs is always consistent, for tests.
	slices.SortFunc(items, func(a, b Item) int {
		return int(a.ID - b.ID)
	})
	return items
}

var sets []*ItemSet

// Registers a new ItemSet with item IDs populated.
func NewItemSet(set ItemSet) *ItemSet {
	foundName := false
	foundAlternativeName := set.AlternativeName == ""
	for _, item := range ItemsByID {
		if item.SetName == "" {
			continue
		}
		foundName = foundName || item.SetName == set.Name
		foundAlternativeName = foundAlternativeName || item.SetName == set.AlternativeName
		if foundName && foundAlternativeName {
			break
		}
	}

	if WITH_DB {
		if !foundName {
			panic("No items found for set " + set.Name)
		}
		if len(set.AlternativeName) > 0 && !foundAlternativeName {
			panic("No items found for set alternative " + set.AlternativeName)
		}
	}

	sets = append(sets, &set)
	return &set
}

func (character *Character) HasSetBonus(set *ItemSet, numItems int32) bool {
	if character.Env != nil && character.Env.IsFinalized() {
		panic("HasSetBonus is very slow and should never be called after finalization. Try caching the value during construction instead!")
	}

	if _, ok := set.Bonuses[numItems]; !ok {
		panic(fmt.Sprintf("Item set %s does not have a bonus with %d pieces.", set.Name, numItems))
	}

	var count int32
	for _, item := range character.Equipment {
		if item.SetName == "" {
			continue
		}
		if item.SetName == set.Name || item.SetName == set.AlternativeName {
			count++
			if count >= numItems {
				return true
			}
		}
	}

	return false
}

type ActiveSetBonus struct {
	// Name of the set.
	Name string

	// Number of pieces required for this bonus.
	NumPieces int32

	// Function for applying the effects of this set bonus.
	BonusEffect ApplyEffect
}

// Returns a list describing all active set bonuses.
func (character *Character) GetActiveSetBonuses() []ActiveSetBonus {
	var activeBonuses []ActiveSetBonus

	setItemCount := make(map[*ItemSet]int32)
	for _, item := range character.Equipment {
		if item.SetName == "" {
			continue
		}
		for _, set := range sets {
			if set.Name == item.SetName || set.AlternativeName == item.SetName {
				setItemCount[set]++
				if bonusEffect, ok := set.Bonuses[setItemCount[set]]; ok {
					activeBonuses = append(activeBonuses, ActiveSetBonus{
						Name:        set.Name,
						NumPieces:   setItemCount[set],
						BonusEffect: bonusEffect,
					})
				}
				break
			}
		}
	}

	return activeBonuses
}

// Apply effects from item set bonuses.
func (character *Character) applyItemSetBonusEffects(agent Agent) {
	activeSetBonuses := character.GetActiveSetBonuses()

	for _, activeSetBonus := range activeSetBonuses {
		activeSetBonus.BonusEffect(agent)
	}
}

// Returns the names of all active set bonuses.
func (character *Character) GetActiveSetBonusNames() []string {
	activeSetBonuses := character.GetActiveSetBonuses()

	names := make([]string, len(activeSetBonuses))
	for i, activeSetBonus := range activeSetBonuses {
		names[i] = fmt.Sprintf("%s (%dpc)", activeSetBonus.Name, activeSetBonus.NumPieces)
	}
	return names
}
