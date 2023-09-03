package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/tools"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/encoding/protojson"
	googleProto "google.golang.org/protobuf/proto"
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
	other := NewWowDatabase()

	for k, v := range db.Items {
		other.Items[k] = v
	}
	for k, v := range db.Enchants {
		other.Enchants[k] = v
	}
	for k, v := range db.Gems {
		other.Gems[k] = v
	}
	for k, v := range db.Zones {
		other.Zones[k] = v
	}
	for k, v := range db.Npcs {
		other.Npcs[k] = v
	}
	for k, v := range db.ItemIcons {
		other.ItemIcons[k] = v
	}
	for k, v := range db.SpellIcons {
		other.SpellIcons[k] = v
	}

	return other
}

func (db *WowDatabase) MergeItems(arr []*proto.UIItem) {
	for _, item := range arr {
		db.MergeItem(item)
	}
}
func (db *WowDatabase) MergeItem(src *proto.UIItem) {
	if dst, ok := db.Items[src.Id]; ok {
		// googleproto.Merge concatenates lists but we want replacement, so do them manually.
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
		// googleproto.Merge concatenates lists but we want replacement, so do them manually.
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
		// googleproto.Merge concatenates lists but we want replacement, so do them manually.
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

func (db *WowDatabase) MergeUIProto(dbProto *proto.UIDatabase) {
	db.MergeItems(dbProto.Items)
	db.MergeEnchants(dbProto.Enchants)
	db.MergeGems(dbProto.Gems)
	db.MergeZones(dbProto.Zones)
	db.MergeNpcs(dbProto.Npcs)
	for _, icon := range dbProto.ItemIcons {
		db.ItemIcons[icon.Id] = icon
	}
	for _, icon := range dbProto.SpellIcons {
		db.SpellIcons[icon.Id] = icon
	}
}

func (db *WowDatabase) ToUIProto() *proto.UIDatabase {
	uidb := &proto.UIDatabase{
		Encounters: db.Encounters,
		GlyphIds:   db.GlyphIDs,
	}

	for _, v := range db.Items {
		uidb.Items = append(uidb.Items, v)
	}
	for _, v := range db.Enchants {
		uidb.Enchants = append(uidb.Enchants, v)
	}
	for _, v := range db.Gems {
		uidb.Gems = append(uidb.Gems, v)
	}
	for _, v := range db.Zones {
		uidb.Zones = append(uidb.Zones, v)
	}
	for _, v := range db.Npcs {
		uidb.Npcs = append(uidb.Npcs, v)
	}
	for _, v := range db.ItemIcons {
		uidb.ItemIcons = append(uidb.ItemIcons, v)
	}
	for _, v := range db.SpellIcons {
		uidb.SpellIcons = append(uidb.SpellIcons, v)
	}

	slices.SortStableFunc(uidb.Items, func(v1, v2 *proto.UIItem) bool {
		return v1.Id < v2.Id
	})
	slices.SortStableFunc(uidb.Enchants, func(v1, v2 *proto.UIEnchant) bool {
		return v1.EffectId < v2.EffectId || v1.EffectId == v2.EffectId && v1.Type < v2.Type
	})
	slices.SortStableFunc(uidb.Gems, func(v1, v2 *proto.UIGem) bool {
		return v1.Id < v2.Id
	})
	slices.SortStableFunc(uidb.Zones, func(v1, v2 *proto.UIZone) bool {
		return v1.Id < v2.Id
	})
	slices.SortStableFunc(uidb.Npcs, func(v1, v2 *proto.UINPC) bool {
		return v1.Id < v2.Id
	})
	slices.SortStableFunc(uidb.ItemIcons, func(v1, v2 *proto.IconData) bool {
		return v1.Id < v2.Id
	})
	slices.SortStableFunc(uidb.SpellIcons, func(v1, v2 *proto.IconData) bool {
		return v1.Id < v2.Id
	})

	return uidb
}

func ReadDatabaseFromJson(jsonStr string) *WowDatabase {
	dbProto := &proto.UIDatabase{}
	if err := protojson.Unmarshal([]byte(jsonStr), dbProto); err != nil {
		panic(err)
	}

	db := NewWowDatabase()
	db.MergeUIProto(dbProto)
	return db
}

func (db *WowDatabase) WriteBinaryAndJson(binFilePath, jsonFilePath string) {
	db.WriteBinary(binFilePath)
	db.WriteJson(jsonFilePath)
}

func (db *WowDatabase) WriteBinary(binFilePath string) {
	uidb := db.ToUIProto()

	// Write database as a binary file.
	outbytes, err := googleProto.Marshal(uidb)
	if err != nil {
		log.Fatalf("[ERROR] Failed to marshal db: %s", err.Error())
	}
	os.WriteFile(binFilePath, outbytes, 0666)
}

func (db *WowDatabase) WriteJson(jsonFilePath string) {
	// Also write in JSON format so we can manually inspect the contents.
	// Write it out line-by-line so we can have 1 line / item, making it more human-readable.
	uidb := db.ToUIProto()
	builder := &strings.Builder{}
	builder.WriteString("{\n")

	tools.WriteProtoArrayToBuilder(uidb.Items, builder, "items")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uidb.Enchants, builder, "enchants")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uidb.Gems, builder, "gems")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uidb.Zones, builder, "zones")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uidb.Npcs, builder, "npcs")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uidb.ItemIcons, builder, "itemIcons")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uidb.SpellIcons, builder, "spellIcons")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uidb.Encounters, builder, "encounters")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uidb.GlyphIds, builder, "glyphIds")
	builder.WriteString("\n")

	builder.WriteString("}")
	os.WriteFile(jsonFilePath, []byte(builder.String()), 0666)
}

func toSlice(stats Stats) []float64 {
	return stats[:]
}
