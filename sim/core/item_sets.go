package core

import (
	"fmt"
	"golang.org/x/exp/slices"
)

type ItemSet struct {
	Name            string
	AlternativeName string

	// IDs of items that are part of this set. map[key]struct{} is roughly a set in go.
	Items map[int32]struct{}

	// Maps set piece requirement to an ApplyEffect function that will be called
	// before the Sim starts.
	//
	// The function should apply any benefits provided by the set bonus.
	Bonuses map[int32]ApplyEffect
}

func (set ItemSet) ItemIDs() []int32 {
	ids := make([]int32, 0, len(set.Items))
	for id := range set.Items {
		ids = append(ids, id)
	}
	// Sort so the order of IDs is always consistent, for tests.
	slices.Sort(ids)
	return ids
}

func (set ItemSet) ItemIsInSet(itemID int32) bool {
	_, ok := set.Items[itemID]
	return ok
}

var sets []*ItemSet

func GetAllItemSets() []*ItemSet {
	// Defensive copy to prevent modifications.
	tmp := make([]*ItemSet, len(sets))
	copy(tmp, sets)
	return tmp
}

// cache for mapping item to set for fast resetting of sim.
var itemSetLookup = map[int32]*ItemSet{}

// Registers a new ItemSet with item IDs populated.
func NewItemSet(setStruct ItemSet) *ItemSet {
	set := &ItemSet{}
	*set = setStruct

	if len(set.Items) > 0 {
		panic(set.Name + " (" + set.AlternativeName + ") supplied item IDs, set items are detected automatically!")
	}

	set.Items = make(map[int32]struct{})
	foundName := false
	foundAlternativeName := false
	for _, item := range ItemsByID {
		if item.SetName == set.Name || (len(set.AlternativeName) > 0 && item.SetName == set.AlternativeName) {
			//fmt.Printf("Adding item %s-%d to set %s\n", item.Name, item.ID, item.SetName)
			set.Items[item.ID] = struct{}{}

			if item.SetName == set.Name {
				foundName = true
			} else {
				foundAlternativeName = true
			}
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

	sets = append(sets, set)
	for itemID := range set.Items {
		itemSetLookup[itemID] = set
	}
	return set
}

func AddItemToSets(item Item) {
	if item.SetName == "" {
		return
	}

	for _, set := range sets {
		if set.Name == item.SetName || set.AlternativeName == item.SetName {
			set.Items[item.ID] = struct{}{}
			itemSetLookup[item.ID] = set
		}
	}
}

func (character *Character) HasSetBonus(itemSet *ItemSet, numItems int32) bool {
	if character.Env != nil && character.Env.IsFinalized() {
		panic("HasSetBonus is very slow and should never be called after finalization. Try caching the value during construction instead!")
	}

	if _, ok := itemSet.Bonuses[numItems]; !ok {
		panic(fmt.Sprintf("Item set %s does not have a bonus with %d pieces.", itemSet.Name, numItems))
	}

	var count int32
	for _, item := range character.Equipment {
		if itemSet.ItemIsInSet(item.ID) {
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
	activeBonuses := []ActiveSetBonus{}
	setItemCount := map[string]int32{}

	for _, i := range character.Equipment {
		set := itemSetLookup[i.ID]
		if set != nil {
			setItemCount[set.Name]++
			if setBonusFunc, ok := set.Bonuses[setItemCount[set.Name]]; ok {
				activeBonuses = append(activeBonuses, ActiveSetBonus{
					Name:        set.Name,
					NumPieces:   setItemCount[set.Name],
					BonusEffect: setBonusFunc,
				})
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
	names := []string{}

	for _, activeSetBonus := range activeSetBonuses {
		names = append(names, fmt.Sprintf("%s (%dpc)", activeSetBonus.Name, activeSetBonus.NumPieces))
	}

	return names
}
