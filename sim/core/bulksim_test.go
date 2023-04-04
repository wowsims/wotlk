package core

import (
	"testing"

	"github.com/wowsims/wotlk/sim/core/proto"
)

const (
	itemStarshardEdge     = 45620
	itemPillarOfFortitude = 46350
	itemIronmender        = 45271
)

var (
	starshardEdge1 = &itemWithSlot{
		Item:  &proto.ItemSpec{Id: itemStarshardEdge},
		Slot:  proto.ItemSlot_ItemSlotMainHand,
		Index: 1,
	}
	starshardEdge2 = &itemWithSlot{
		Item:  &proto.ItemSpec{Id: itemStarshardEdge},
		Slot:  proto.ItemSlot_ItemSlotMainHand,
		Index: 2,
	}
	pillarOfFortitude = &itemWithSlot{
		Item:  &proto.ItemSpec{Id: itemPillarOfFortitude},
		Slot:  proto.ItemSlot_ItemSlotMainHand,
		Index: 3,
	}
	ironmender = &itemWithSlot{
		Item:  &proto.ItemSpec{Id: itemIronmender},
		Slot:  proto.ItemSlot_ItemSlotOffHand,
		Index: 4,
	}

	tinyItemDatabase = &proto.SimDatabase{
		Items: []*proto.SimItem{
			{Id: itemStarshardEdge, Type: proto.ItemType_ItemTypeWeapon, HandType: proto.HandType_HandTypeMainHand},
			{Id: itemPillarOfFortitude, Type: proto.ItemType_ItemTypeWeapon, HandType: proto.HandType_HandTypeTwoHand},
			{Id: itemIronmender, Type: proto.ItemType_ItemTypeWeapon, HandType: proto.HandType_HandTypeOffHand},
		},
	}
)

func TestEquipmentSubstitutionIsValid(t *testing.T) {
	for _, tc := range []struct {
		comment string
		items   []*itemWithSlot
		want    bool
	}{
		{
			comment: "empty replacement is valid (1)",
			items:   nil,
			want:    true,
		},
		{
			comment: "empty replacement is valid (2)",
			items:   []*itemWithSlot{},
			want:    true,
		},
		{
			comment: "mainhand replacement is valid",
			items:   []*itemWithSlot{starshardEdge1},
			want:    true,
		},
		{
			comment: "same item cannot occurr twice in a substitution",
			items:   []*itemWithSlot{starshardEdge1, starshardEdge1},
			want:    false,
		},
		{
			comment: "cannot use two items for the same item slot",
			items:   []*itemWithSlot{starshardEdge1, starshardEdge2},
			want:    false,
		},
	} {
		sub := &equipmentSubstitution{Items: tc.items}
		if got := sub.IsValid(); got != tc.want {
			t.Fatalf("%s: equipmentSubstitution.IsValid(%v) = %v, want %v", tc.comment, sub, got, tc.want)
		}
	}
}

func TestIsValidEquipment(t *testing.T) {
	addToDatabase(tinyItemDatabase)

	for _, tc := range []struct {
		comment string
		spec    *proto.EquipmentSpec
		want    bool
	}{
		{
			comment: "simple equipment set with just one mainhand weapon is valid",
			spec:    createEquipmentFromItems(starshardEdge1),
			want:    true,
		},
		{
			comment: "cannot equip offhand and two-hander",
			spec:    createEquipmentFromItems(pillarOfFortitude, ironmender),
			want:    false,
		},
	} {
		if got := isValidEquipment(tc.spec); got != tc.want {
			t.Fatalf("%s: isValidEquipment(%v) = %v, want %v", tc.comment, tc.spec, got, tc.want)
		}
	}
}

func createEquipmentFromItems(items ...*itemWithSlot) *proto.EquipmentSpec {
	spec := &proto.EquipmentSpec{
		Items: make([]*proto.ItemSpec, 17),
	}
	for _, is := range items {
		spec.Items[is.Slot] = is.Item
	}
	return spec
}

func TestGenerateAllequipmentSubstitutions(t *testing.T) {
	// TODO(Riotdog-GehennasEU): Implement.
}

func TestCreateNewRequestWithSubstitution(t *testing.T) {
	// TODO(Riotdog-GehennasEU): Implement.
}
