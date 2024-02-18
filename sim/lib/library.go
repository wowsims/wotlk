package main

// #include <stdlib.h>
import "C"
import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"log"
	"unsafe"

	"github.com/wowsims/wotlk/sim"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
	goproto "google.golang.org/protobuf/proto"
)

var _default_rsr = proto.RaidSimRequest{
	Raid:       &proto.Raid{},
	Encounter:  &proto.Encounter{},
	SimOptions: &proto.SimOptions{},
}
var _active_sim = core.NewSim(&_default_rsr)
var _active_seed int64 = 1
var _aura_labels = []string{}
var _target_aura_labels = []string{}

//export runSim
func runSim(json *C.char) *C.char {
	input := &proto.RaidSimRequest{}
	jsonString := C.GoString(json)
	err := protojson.Unmarshal([]byte(jsonString), input)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}
	sim.RegisterAll()
	result := core.RunSim(input, nil)
	out, err := protojson.Marshal(result)
	if err != nil {
		panic(err)
	}
	return C.CString(string(out))
}

//export computeStats
func computeStats(json *C.char) *C.char {
	input := &proto.ComputeStatsRequest{}
	jsonString := C.GoString(json)
	err := protojson.Unmarshal([]byte(jsonString), input)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}
	sim.RegisterAll()
	result := core.ComputeStats(input)
	out, err := protojson.Marshal(result)
	if err != nil {
		panic(err)
	}
	return C.CString(string(out))
}

//export encodeSettings
func encodeSettings(json *C.char) *C.char {
	input := &proto.RaidSimRequest{}
	jsonString := C.GoString(json)
	err := protojson.Unmarshal([]byte(jsonString), input)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}
	settings := &proto.IndividualSimSettings{
		Settings: &proto.SimSettings{
			Iterations: input.SimOptions.Iterations,
		},
		RaidBuffs:  input.Raid.Buffs,
		Debuffs:    input.Raid.Debuffs,
		Tanks:      input.Raid.Tanks,
		PartyBuffs: input.Raid.Parties[0].Buffs,
		Player:     input.Raid.Parties[0].Players[0],
		Encounter:  input.Encounter,
	}
	var buffer bytes.Buffer
	data, err := goproto.Marshal(settings)
	if err != nil {
		panic(err)
	}
	writer := zlib.NewWriter(&buffer)
	writer.Write(data)
	writer.Close()
	out := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return C.CString(string(out))
}

//export getDatabase
func getDatabase(itemIds *int32, numItems int32, enchantIds *int32, numEnchants int32, gemIds *int32, numGems int32) *C.char {
	ids := unsafe.Slice(itemIds, numItems)
	eids := unsafe.Slice(enchantIds, numEnchants)
	gids := unsafe.Slice(gemIds, numGems)
	simDB := &proto.SimDatabase{
		Items:    make([]*proto.SimItem, numItems),
		Enchants: make([]*proto.SimEnchant, numEnchants),
		Gems:     make([]*proto.SimGem, numGems),
	}
	for i, itemId := range ids {
		item := core.ItemsByID[itemId]
		simDB.Items[i] = &proto.SimItem{
			Id:               item.ID,
			Name:             item.Name,
			Type:             item.Type,
			ArmorType:        item.ArmorType,
			WeaponType:       item.WeaponType,
			HandType:         item.HandType,
			RangedWeaponType: item.RangedWeaponType,
			Stats:            item.Stats[:],
			GemSockets:       item.GemSockets,
			SocketBonus:      item.SocketBonus[:],
			WeaponDamageMin:  item.WeaponDamageMin,
			WeaponDamageMax:  item.WeaponDamageMax,
			WeaponSpeed:      item.SwingSpeed,
			SetName:          item.SetName,
		}
	}
	for i, enchantId := range eids {
		enchant := core.EnchantsByEffectID[enchantId]
		simDB.Enchants[i] = &proto.SimEnchant{
			EffectId: enchant.EffectID,
			Stats:    enchant.Stats[:],
		}
	}
	for i, gemId := range gids {
		gem := core.GemsByID[gemId]
		simDB.Gems[i] = &proto.SimGem{
			Id:    gem.ID,
			Name:  gem.Name,
			Color: gem.Color,
			Stats: gem.Stats[:],
		}
	}
	out, err := protojson.Marshal(simDB)
	if err != nil {
		panic(err)
	}
	return C.CString(string(out))
}

//export new
func new(json *C.char) {
	input := &proto.RaidSimRequest{}
	jsonString := C.GoString(json)
	err := protojson.Unmarshal([]byte(jsonString), input)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}
	sim.RegisterAll()
	_active_sim = core.NewSim(input)
	_active_sim.Reseed(_active_seed)
	_active_seed += 1
	_active_sim.Reset()
	_active_sim.PrePull()
}

//export trySpell
func trySpell(act int) bool {
	player := _active_sim.Raid.Parties[0].Players[0]
	spells := player.GetCharacter().Spellbook
	if act >= len(spells) || act < 0 {
		return false
	}
	spell := spells[act]
	target := player.GetCharacter().CurrentTarget
	casted := false

	// FIXME : This is a hack to allow Heroic strike to work
	if spell.ActionID.SpellID == 47450 {
		aura := player.GetCharacter().GetAura("HS Queue Aura")
		if aura.IsActive() {
			return false
		}
		aura.Activate(_active_sim)
		return true
	}
	// End of Heroic strike hack

	if spell.CanCast(_active_sim, target) {
		casted = spell.Cast(_active_sim, target)
		if casted && spell.CurCast.GCD > 0 {
			_active_sim.NeedsInput = false
		}
	}
	return casted
}

//export doNothing
func doNothing() bool {
	return true
}

//export getRemainingDuration
func getRemainingDuration() float64 {
	return _active_sim.GetRemainingDuration().Seconds()
}

//export getEnergy
func getEnergy() float64 {
	player := _active_sim.Raid.Parties[0].Players[0]
	if !player.GetCharacter().HasEnergyBar() {
		return 0.0
	}
	return player.GetCharacter().CurrentEnergy()
}

//export getComboPoints
func getComboPoints() int {
	player := _active_sim.Raid.Parties[0].Players[0]
	if !player.GetCharacter().HasEnergyBar() {
		return 0
	}
	return int(player.GetCharacter().ComboPoints())
}

//export getUnitCount
func getUnitCount() int {
	return len(_active_sim.AllUnits)
}

//export getSpellCount
func getSpellCount() int {
	return len(_active_sim.Raid.Parties[0].Players[0].GetCharacter().Spellbook)
}

//export getSpells
func getSpells(storage *int32, n int32) {
	player := _active_sim.Raid.Parties[0].Players[0]
	spellbook := player.GetCharacter().Spellbook
	spells := unsafe.Slice(storage, n)
	for i, spell := range spellbook[:n] {
		if spell.Tag != -1 {
			spells[i] = spell.ActionID.SpellID
		} else {
			// These spells are not castable by the player
			spells[i] = -1
		}
	}
}

//export getCooldowns
func getCooldowns(storage *float64, spellbookIndices *int32, n int32) {
	player := _active_sim.Raid.Parties[0].Players[0]
	spellbook := player.GetCharacter().Spellbook
	spells := unsafe.Slice(spellbookIndices, n)
	cds := unsafe.Slice(storage, n)
	for i := int32(0); i < n; i++ {
		spellbookIndex := spells[i]
		spell := spellbook[spellbookIndex]
		cds[i] = spell.TimeToReady(_active_sim).Seconds()
	}
}

//export registerAuras
func registerAuras(strings **C.char) {
	_aura_labels = []string{}
	labels := unsafe.Slice(strings, 1<<30)
	for i := 0; labels[i] != nil; i++ {
		_aura_labels = append(_aura_labels, C.GoString(labels[i]))
	}
}

//export registerTargetAuras
func registerTargetAuras(strings **C.char) {
	_target_aura_labels = []string{}
	labels := unsafe.Slice(strings, 1<<30)
	for i := 0; labels[i] != nil; i++ {
		_target_aura_labels = append(_target_aura_labels, C.GoString(labels[i]))
	}
}

//export getAuras
func getAuras(storage *float64, n int32) {
	player := _active_sim.Raid.Parties[0].Players[0]
	auras := unsafe.Slice(storage, n)
	for i, label := range _aura_labels {
		aura := player.GetCharacter().GetAura(label)
		if aura != nil {
			auras[i] = aura.RemainingDuration(_active_sim).Seconds()
		} else {
			auras[i] = 0.0
		}
	}
}

//export getTargetAuras
func getTargetAuras(storage *float64, n int32) {
	player := _active_sim.Raid.Parties[0].Players[0]
	target := player.GetCharacter().CurrentTarget
	auras := unsafe.Slice(storage, n)
	for i, label := range _target_aura_labels {
		aura := target.GetAura(label)
		if aura != nil {
			auras[i] = aura.RemainingDuration(_active_sim).Seconds()
		} else {
			auras[i] = 0.0
		}
	}
}

//export getDamageDone
func getDamageDone() float64 {
	player := _active_sim.Raid.Parties[0].Players[0]
	spellbook := player.GetCharacter().Spellbook
	totalDamage := 0.0
	for _, spell := range spellbook {
		for _, metrics := range spell.SpellMetrics {
			totalDamage += metrics.TotalDamage
		}
	}
	return totalDamage
}

//export getSpellMetrics
func getSpellMetrics() *C.char {
	all_metrics := make(map[int32][]core.SpellMetrics)
	player := _active_sim.Raid.Parties[0].Players[0]
	spellbook := player.GetCharacter().Spellbook
	for _, spell := range spellbook {
		spell_id := spell.ActionID.SpellID
		for _, metrics := range spell.SpellMetrics {
			if metrics.Casts > 0 {
				all_metrics[spell_id] = append(all_metrics[spell_id], metrics)
			}
		}
	}
	out, err := json.Marshal(all_metrics)
	if err != nil {
		panic(err)
	}
	return C.CString(string(out))
}

//export step
func step() bool {
	return _active_sim.Step()
}

//export needsInput
func needsInput() bool {
	return _active_sim.NeedsInput
}

//export cleanup
func cleanup() {
	_active_sim.Cleanup()
}

//export FreeCString
func FreeCString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func main() {}
