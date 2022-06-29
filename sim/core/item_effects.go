package core

import (
	"github.com/wowsims/tbc/sim/core/proto"
	"log"
)

// Function for applying permanent effects to an Agent.
//
// Passing Character instead of Agent would work for almost all cases,
// but there are occasionally class-specific item effects.
type ApplyEffect func(Agent)

// Function for applying permenent effects to an agent's weapon
type ApplyWeaponEffect func(Agent, proto.ItemSlot)

var itemEffects = map[int32]ApplyEffect{}
var weaponEffects = map[int32]ApplyWeaponEffect{}

func HasItemEffect(id int32) bool {
	_, ok := itemEffects[id]
	return ok
}

func HasWeaponEffect(id int32) bool {
	_, ok := weaponEffects[id]
	return ok
}

// Registers an ApplyEffect function which will be called before the Sim
// starts, for any Agent that is wearing the item.
func NewItemEffect(id int32, itemEffect ApplyEffect) {
	if HasItemEffect(id) {
		log.Fatalf("Cannot add multiple effects for one item: %d, %#v", id, itemEffect)
	}
	itemEffects[id] = itemEffect
}

func AddWeaponEffect(id int32, weaponEffect ApplyWeaponEffect) {
	if HasWeaponEffect(id) {
		log.Fatalf("Cannot add multiple effects for one item: %d, %#v", id, weaponEffect)
	}
	weaponEffects[id] = weaponEffect
}
