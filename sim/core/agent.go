package core

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/wowsims/tbc/sim/core/proto"
)

// Agent can be thought of as the 'Player', i.e. the thing controlling the Character.
// This is the interface implemented by each class/spec.
type Agent interface {
	// The Character controlled by this Agent.
	GetCharacter() *Character

	// Called once after all Players/Pets/Targets have finished the construction phase.
	// Use this to register spells and any initialization steps that require
	// other raid members or auras.
	Initialize()

	// Updates the input Buffs to include raid-wide buffs provided by this Agent.
	AddRaidBuffs(raidBuffs *proto.RaidBuffs)
	// Updates the input Buffs to include party-wide buffs provided by this Agent.
	AddPartyBuffs(partyBuffs *proto.PartyBuffs)

	// Applies any bonuses that should be considered part of gear. This makes sure
	// the bonus is applied at the appropriate time to be accounted in player stats.
	ApplyGearBonuses()

	// All talent stats / auras should be added within this callback. This makes sure
	// talents are applied at the right time so we can calculate groups of stats.
	ApplyTalents()

	// Returns this Agent to its initial state. Called before each Sim iteration
	// and once after the final iteration.
	Reset(sim *Simulation)

	// Called whenever the GCD becomes ready for this Agent.
	OnGCDReady(sim *Simulation)

	// Called after each mana tick, if this Agent uses mana.
	OnManaTick(sim *Simulation)

	// Called after each auto attack performed by this Agent.
	// This is different from Aura.OnSpellHit in that it is invoked fully after
	// everything related to the attack is complete, and it is only invoked for
	// auto attacks (white hits or white-hit-replacers).
	OnAutoAttack(sim *Simulation, spell *Spell)
}

type ActionID struct {
	// Only one of these should be set.
	SpellID int32
	ItemID  int32
	OtherID proto.OtherAction

	Tag int32
}

func (actionID ActionID) IsEmptyAction() bool {
	return actionID.SpellID == 0 && actionID.ItemID == 0 && actionID.OtherID == 0
}

func (actionID ActionID) IsSpellAction(spellID int32) bool {
	return actionID.SpellID == spellID
}

func (actionID ActionID) IsItemAction(itemID int32) bool {
	return actionID.ItemID == itemID
}

func (actionID ActionID) IsOtherAction(otherID proto.OtherAction) bool {
	return actionID.OtherID == otherID
}

func (actionID ActionID) SameActionIgnoreTag(other ActionID) bool {
	return actionID.SpellID == other.SpellID && actionID.ItemID == other.ItemID && actionID.OtherID == other.OtherID
}

func (actionID ActionID) SameAction(other ActionID) bool {
	return actionID.SameActionIgnoreTag(other) && actionID.Tag == other.Tag
}

func (actionID ActionID) String() string {
	var sb strings.Builder
	sb.WriteString("{")

	if actionID.SpellID != 0 {
		sb.WriteString("SpellID: ")
		sb.WriteString(strconv.Itoa(int(actionID.SpellID)))
	} else if actionID.ItemID != 0 {
		sb.WriteString("ItemID: ")
		sb.WriteString(strconv.Itoa(int(actionID.ItemID)))
	} else if actionID.OtherID != 0 {
		sb.WriteString("OtherID: ")
		sb.WriteString(strconv.Itoa(int(actionID.OtherID)))
	}
	if actionID.Tag != 0 {
		sb.WriteString(", Tag: ")
		sb.WriteString(strconv.Itoa(int(actionID.Tag)))
	}
	sb.WriteString("}")

	return sb.String()
}

// Returns a new ActionID with the corresponding Tag value.
func (actionID ActionID) WithTag(tag int32) ActionID {
	newID := actionID
	newID.Tag = tag
	return newID
}

func (actionID ActionID) ToProto() *proto.ActionID {
	protoID := &proto.ActionID{
		Tag: actionID.Tag,
	}

	if actionID.SpellID != 0 {
		protoID.RawId = &proto.ActionID_SpellId{SpellId: actionID.SpellID}
	} else if actionID.ItemID != 0 {
		protoID.RawId = &proto.ActionID_ItemId{ItemId: actionID.ItemID}
	} else if actionID.OtherID != 0 {
		protoID.RawId = &proto.ActionID_OtherId{OtherId: actionID.OtherID}
	}

	return protoID
}

func ProtoToActionID(protoID proto.ActionID) ActionID {
	return ActionID{
		ItemID:  protoID.GetItemId(),
		SpellID: protoID.GetSpellId(),
		OtherID: protoID.GetOtherId(),
		Tag:     protoID.Tag,
	}
}

type AgentFactory func(Character, proto.Player) Agent
type SpecSetter func(*proto.Player, interface{})

var agentFactories = make(map[string]AgentFactory)
var specSetters = make(map[string]SpecSetter)
var configSpecs = make(map[string]proto.Spec)

func PlayerProtoToSpec(player proto.Player) proto.Spec {
	typeName := reflect.TypeOf(player.GetSpec()).Elem().Name()
	return configSpecs[typeName]
}

func RegisterAgentFactory(emptyOptions interface{}, spec proto.Spec, factory AgentFactory, specSetter SpecSetter) {
	typeName := reflect.TypeOf(emptyOptions).Name()
	if _, ok := agentFactories[typeName]; ok {
		panic("Aleady registered agent factory: " + typeName)
	}
	//fmt.Printf("Registering type: %s", typeName)

	agentFactories[typeName] = factory
	specSetters[typeName] = specSetter
	configSpecs[typeName] = spec
}

// Constructs a new Agent.
func NewAgent(party *Party, partyIndex int, player proto.Player) Agent {
	typeName := reflect.TypeOf(player.GetSpec()).Elem().Name()

	factory, ok := agentFactories[typeName]
	if !ok {
		panic("No agent factory for type: " + typeName)
	}

	character := NewCharacter(party, partyIndex, player)
	return factory(character, player)
}

// Applies the spec options to the given player. This is only necessary because
// the generated proto code does not export oneof interface types.
// Player is returned so this function can be used in-line with player creation.
func WithSpec(player *proto.Player, spec interface{}) *proto.Player {
	typeName := reflect.TypeOf(spec).Elem().Name()

	specSetter, ok := specSetters[typeName]
	if !ok {
		panic("No spec setter for type: " + typeName)
	}

	specSetter(player, spec)
	return player
}
