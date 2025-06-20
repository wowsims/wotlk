package database

import (
	_ "embed"

	"github.com/wowsims/wotlk/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

//go:embed db.bin
var dbBytes []byte

func Load() *proto.UIDatabase {
	db := &proto.UIDatabase{}
	if err := googleProto.Unmarshal(dbBytes, db); err != nil {
		panic(err)
	}
	return db
}
