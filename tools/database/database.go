package database

import (
	"bytes"
	"fmt"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/tools"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/encoding/protojson"
	googleProto "google.golang.org/protobuf/proto"
	"log"
	"os"
)

type EnchantDBKey struct {
	EffectID int32
	ItemID   int32
	SpellID  int32
}

func EnchantToDBKey(enchant *proto.UIEnchant) EnchantDBKey {
	return EnchantDBKey{
		EffectID: enchant.EffectId,
		ItemID:   enchant.ItemId,
		SpellID:  enchant.SpellId,
	}
}

type WowDatabase struct {
	Items    map[int32]*proto.UIItem
	Enchants map[EnchantDBKey]*proto.UIEnchant
	Gems     map[int32]*proto.UIGem

	Zones map[int32]*proto.UIZone
	Npcs  map[int32]*proto.UINPC

	ItemIcons  map[int32]*proto.IconData
	SpellIcons map[int32]*proto.IconData

	Encounters []*proto.PresetEncounter
	GlyphIDs   []*proto.GlyphID
}

func NewWowDatabase() *WowDatabase {
	return &WowDatabase{
		Items:    make(map[int32]*proto.UIItem),
		Enchants: make(map[EnchantDBKey]*proto.UIEnchant),
		Gems:     make(map[int32]*proto.UIGem),
		Zones:    make(map[int32]*proto.UIZone),
		Npcs:     make(map[int32]*proto.UINPC),

		ItemIcons:  make(map[int32]*proto.IconData),
		SpellIcons: make(map[int32]*proto.IconData),
	}
}

func (db *WowDatabase) Clone() *WowDatabase {
	return &WowDatabase{
		Items:      maps.Clone(db.Items),
		Enchants:   maps.Clone(db.Enchants),
		Gems:       maps.Clone(db.Gems),
		Zones:      maps.Clone(db.Zones),
		Npcs:       maps.Clone(db.Npcs),
		ItemIcons:  maps.Clone(db.ItemIcons),
		SpellIcons: maps.Clone(db.SpellIcons),
	}
}

func (db *WowDatabase) MergeItems(arr []*proto.UIItem) {
	for _, item := range arr {
		db.MergeItem(item)
	}
}
func (db *WowDatabase) MergeItem(src *proto.UIItem) {
	if dst, ok := db.Items[src.Id]; ok {
		// googleproto.Merge concatenates lists, but we want replacement, so do them manually.
		if src.Stats != nil {
			dst.Stats = src.Stats
			src.Stats = nil
		}
		if src.SocketBonus != nil {
			dst.SocketBonus = src.SocketBonus
			src.SocketBonus = nil
		}
		googleProto.Merge(dst, src)
	} else {
		db.Items[src.Id] = src
	}
}

func (db *WowDatabase) MergeEnchants(arr []*proto.UIEnchant) {
	for _, enchant := range arr {
		db.MergeEnchant(enchant)
	}
}
func (db *WowDatabase) MergeEnchant(src *proto.UIEnchant) {
	key := EnchantToDBKey(src)
	if dst, ok := db.Enchants[key]; ok {
		// googleproto.Merge concatenates lists, but we want replacement, so do them manually.
		if src.Stats != nil {
			dst.Stats = src.Stats
			src.Stats = nil
		}
		googleProto.Merge(dst, src)
	} else {
		db.Enchants[key] = src
	}
}

func (db *WowDatabase) MergeGems(arr []*proto.UIGem) {
	for _, gem := range arr {
		db.MergeGem(gem)
	}
}
func (db *WowDatabase) MergeGem(src *proto.UIGem) {
	if dst, ok := db.Gems[src.Id]; ok {
		// googleproto.Merge concatenates lists, but we want replacement, so do them manually.
		if src.Stats != nil {
			dst.Stats = src.Stats
			src.Stats = nil
		}
		googleProto.Merge(dst, src)
	} else {
		db.Gems[src.Id] = src
	}
}

func (db *WowDatabase) MergeZones(arr []*proto.UIZone) {
	for _, zone := range arr {
		db.MergeZone(zone)
	}
}
func (db *WowDatabase) MergeZone(src *proto.UIZone) {
	if dst, ok := db.Zones[src.Id]; ok {
		googleProto.Merge(dst, src)
	} else {
		db.Zones[src.Id] = src
	}
}

func (db *WowDatabase) MergeNpcs(arr []*proto.UINPC) {
	for _, npc := range arr {
		db.MergeNpc(npc)
	}
}
func (db *WowDatabase) MergeNpc(src *proto.UINPC) {
	if dst, ok := db.Npcs[src.Id]; ok {
		googleProto.Merge(dst, src)
	} else {
		db.Npcs[src.Id] = src
	}
}

func (db *WowDatabase) AddItemIcon(id int32, tooltips map[int32]WowheadItemResponse) {
	if tooltip, ok := tooltips[id]; ok {
		if tooltip.GetName() == "" || tooltip.GetIcon() == "" {
			return
		}
		db.ItemIcons[id] = &proto.IconData{Id: id, Name: tooltip.GetName(), Icon: tooltip.GetIcon()}
	} else {
		panic(fmt.Sprintf("No item tooltip with id %d", id))
	}
}

func (db *WowDatabase) AddSpellIcon(id int32, tooltips map[int32]WowheadItemResponse) {
	if tooltip, ok := tooltips[id]; ok {
		if tooltip.GetName() == "" || tooltip.GetIcon() == "" {
			return
		}
		db.SpellIcons[id] = &proto.IconData{Id: id, Name: tooltip.GetName(), Icon: tooltip.GetIcon()}
	} else {
		panic(fmt.Sprintf("No spell tooltip with id %d", id))
	}
}

type idKeyed interface {
	GetId() int32
}

func mapToSlice[T idKeyed](m map[int32]T) []T {
	vs := make([]T, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	slices.SortFunc(vs, func(a, b T) bool {
		return a.GetId() < b.GetId()
	})
	return vs
}

func (db *WowDatabase) ToUIProto() *proto.UIDatabase {
	enchants := make([]*proto.UIEnchant, 0, len(db.Enchants))
	for _, v := range db.Enchants {
		enchants = append(enchants, v)
	}
	slices.SortFunc(enchants, func(v1, v2 *proto.UIEnchant) bool {
		return v1.EffectId < v2.EffectId || v1.EffectId == v2.EffectId && v1.Type < v2.Type
	})

	return &proto.UIDatabase{
		Items:      mapToSlice(db.Items),
		Enchants:   enchants,
		Gems:       mapToSlice(db.Gems),
		Encounters: db.Encounters,
		Zones:      mapToSlice(db.Zones),
		Npcs:       mapToSlice(db.Npcs),
		ItemIcons:  mapToSlice(db.ItemIcons),
		SpellIcons: mapToSlice(db.SpellIcons),
		GlyphIds:   db.GlyphIDs,
	}
}

func sliceToMap[T idKeyed](vs []T) map[int32]T {
	m := make(map[int32]T, len(vs))
	for _, v := range vs {
		m[v.GetId()] = v
	}
	return m
}

func ReadDatabaseFromJson(jsonStr string) *WowDatabase {
	dbProto := &proto.UIDatabase{}
	if err := protojson.Unmarshal([]byte(jsonStr), dbProto); err != nil {
		panic(err)
	}

	enchants := make(map[EnchantDBKey]*proto.UIEnchant, len(dbProto.Enchants))
	for _, v := range dbProto.Enchants {
		enchants[EnchantToDBKey(v)] = v
	}

	return &WowDatabase{
		Items:      sliceToMap(dbProto.Items),
		Enchants:   enchants,
		Gems:       sliceToMap(dbProto.Gems),
		Zones:      sliceToMap(dbProto.Zones),
		Npcs:       sliceToMap(dbProto.Npcs),
		ItemIcons:  sliceToMap(dbProto.ItemIcons),
		SpellIcons: sliceToMap(dbProto.SpellIcons),
	}
}

func (db *WowDatabase) WriteBinaryAndJson(binFilePath, jsonFilePath string) {
	db.WriteBinary(binFilePath)
	db.WriteJson(jsonFilePath)
}

func (db *WowDatabase) WriteBinary(binFilePath string) {
	uidb := db.ToUIProto()

	// Write database as a binary file.
	protoBytes, err := googleProto.Marshal(uidb)
	if err != nil {
		log.Fatalf("[ERROR] Failed to marshal db: %s", err.Error())
	}
	os.WriteFile(binFilePath, protoBytes, 0666)
}

func (db *WowDatabase) WriteJson(jsonFilePath string) {
	// Also write in JSON format, so we can manually inspect the contents.
	// Write it out line-by-line, so we can have 1 line / item, making it more human-readable.
	uidb := db.ToUIProto()

	buffer := new(bytes.Buffer)

	buffer.WriteString("{\n")

	tools.WriteProtoArrayToBuffer(uidb.Items, buffer, "items")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.Enchants, buffer, "enchants")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.Gems, buffer, "gems")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.Zones, buffer, "zones")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.Npcs, buffer, "npcs")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.ItemIcons, buffer, "itemIcons")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.SpellIcons, buffer, "spellIcons")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.Encounters, buffer, "encounters")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.GlyphIds, buffer, "glyphIds")
	buffer.WriteString("\n")

	buffer.WriteString("}")
	os.WriteFile(jsonFilePath, buffer.Bytes(), 0666)
}

func toSlice(stats Stats) []float64 {
	return stats[:]
}
