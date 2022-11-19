package main

import (
	"log"
	"os"
	"reflect"
	"strings"

	protojson "google.golang.org/protobuf/encoding/protojson"
	googleProto "google.golang.org/protobuf/proto"
)

func writeDatabaseFile(db *WowDatabase) {
	uiDB := db.toUIDatabase()

	// Write database as a binary file.
	outbytes, err := googleProto.Marshal(uiDB)
	if err != nil {
		log.Fatalf("[ERROR] Failed to marshal db: %s", err.Error())
	}
	os.WriteFile("./assets/database/db.bin", outbytes, 0666)

	// Also write in JSON format so we can manually inspect the contents.
	// Write it out line-by-line so we can have 1 line / item, making it more human-readable.
	builder := &strings.Builder{}
	builder.WriteString("{\n")

	writeProtoArrayToBuilder(uiDB.Items, builder, "items")
	builder.WriteString(",\n")
	writeProtoArrayToBuilder(uiDB.Enchants, builder, "enchants")
	builder.WriteString(",\n")
	writeProtoArrayToBuilder(uiDB.Gems, builder, "gems")
	builder.WriteString(",\n")
	writeProtoArrayToBuilder(uiDB.ItemIcons, builder, "itemIcons")
	builder.WriteString(",\n")
	writeProtoArrayToBuilder(uiDB.SpellIcons, builder, "spellIcons")
	builder.WriteString("\n")

	builder.WriteString("}")
	os.WriteFile("./assets/database/db.json", []byte(builder.String()), 0666)
}

func writeProtoArrayToBuilder(arrInterface interface{}, builder *strings.Builder, name string) {
	arr := InterfaceSlice(arrInterface)
	builder.WriteString("\"")
	builder.WriteString(name)
	builder.WriteString("\":[\n")

	for i, elem := range arr {
		json, err := protojson.MarshalOptions{}.Marshal(elem.(googleProto.Message))
		if err != nil {
			log.Printf("[ERROR] Failed to marshal: %s", err.Error())
		}
		builder.WriteString(string(json))
		if i != len(arr)-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}
	builder.WriteString("]")
}

// Needed because Go won't let us cast from []FooProto --> []googleProto.Message
// https://stackoverflow.com/questions/12753805/type-converting-slices-of-interfaces
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
