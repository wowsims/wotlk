package core

import (
	"testing"
)

func TestNothing(t *testing.T) {

}

/*
const (
	itemStarshardEdge     = 45620
	itemPillarOfFortitude = 46350
	itemIronmender        = 45271
)

var (
	starshardEdge1 = &itemWithSlot{
		Item:  &proto.ItemSpec{Id: itemStarshardEdge},
		Slot:  proto.ItemSlot_ItemSlotMainHand,
		Index: 0,
	}
	ironmender = &itemWithSlot{
		Item:  &proto.ItemSpec{Id: itemIronmender},
		Slot:  proto.ItemSlot_ItemSlotOffHand,
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
	// This is a bit awkward because code everywhere accesses the global database maps. Hopefully
	// this won't mess with any other unit tests that need existing item/gem/enchant databases?
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

func TestGenerateAllEquipmentSubstitutions(t *testing.T) {
	for _, tc := range []struct {
		comment string
		spec    *proto.BulkEquipmentSpec
		want    []*equipmentSubstitution
	}{
		{
			comment: "empty spec returns empty base equipment substitution only",
			spec:    createBulkSpecFromItems(),
			want: []*equipmentSubstitution{
				{},
			},
		},
		{
			comment: "spec with one item returns empty base equipment substitution plus one item substitution",
			spec:    createBulkSpecFromItems(starshardEdge1),
			want: []*equipmentSubstitution{
				{},
				{Items: []*itemWithSlot{starshardEdge1}},
			},
		},
		{
			comment: "spec with one item returns empty base equipment substitution plus one item substitution",
			spec:    createBulkSpecFromItems(starshardEdge1),
			want: []*equipmentSubstitution{
				{},
				{Items: []*itemWithSlot{starshardEdge1}},
			},
		},
		{
			comment: "spec with two items returns empty base equipment substitution plus all item combos",
			spec:    createBulkSpecFromItems(starshardEdge1, ironmender),
			want: []*equipmentSubstitution{
				{},
				{Items: []*itemWithSlot{starshardEdge1}},
				{Items: []*itemWithSlot{ironmender}},
				{Items: []*itemWithSlot{starshardEdge1, ironmender}},
			},
		},
		{
			comment: "spec with a duplicate item slot returns only valid substitutions",
			spec:    createBulkSpecFromItems(starshardEdge1, ironmender, starshardEdge2),
			want: []*equipmentSubstitution{
				{},
				{Items: []*itemWithSlot{starshardEdge1}},
				{Items: []*itemWithSlot{ironmender}},
				{Items: []*itemWithSlot{starshardEdge1, ironmender}},
				{Items: []*itemWithSlot{starshardEdge2}},
				{Items: []*itemWithSlot{ironmender, starshardEdge2}},
			},
		},
	} {
		var got []*equipmentSubstitution
		for sub := range generateAllEquipmentSubstitutions(context.Background(), tc.spec) {
			got = append(got, sub)
		}

		if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreUnexported(proto.ItemSpec{})); diff != "" {
			t.Fatalf("%s: generateAllEquipmentSubstitutions(%v) returned diff (-want +got):\n%s", tc.comment, tc.spec, diff)
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

func createBulkSpecFromItems(items ...*itemWithSlot) *proto.BulkEquipmentSpec {
	spec := &proto.BulkEquipmentSpec{}
	for _, is := range items {
		spec.Items = append(spec.Items, &proto.BulkEquipmentSpec_ItemSpecWithSlots{
			Item:  is.Item,
			Slots: []proto.ItemSlot{is.Slot},
		})
	}
	return spec
}
*/
