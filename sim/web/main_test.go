package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	_ "github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

var basicSpec = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Options: &proto.ElementalShaman_Options{
			Shield: proto.ShamanShield_WaterShield,
		},
	},
}

var p1Equip = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{Id: 40516, Enchant: 3820, Gems: []int32{41285, 40027}},
		{Id: 44661, Gems: []int32{39998}},
		{Id: 40286, Enchant: 3810},
		{Id: 44005, Enchant: 3722, Gems: []int32{40027}},
		{Id: 40514, Enchant: 3832, Gems: []int32{42144, 42144}},
		{Id: 40324, Enchant: 2332, Gems: []int32{42144, 0}},
		{Id: 40302, Enchant: 3246, Gems: []int32{0}},
		{Id: 40301, Gems: []int32{40014}},
		{Id: 40560, Enchant: 3721},
		{Id: 40519, Enchant: 3826},
		{Id: 37694},
		{Id: 40399},
		{Id: 40432},
		{Id: 40255},
		{Id: 40395, Enchant: 3834},
		{Id: 40401, Enchant: 1128},
		{Id: 40267},
	},
}

func init() {
	s := &server{
		progMut:         sync.RWMutex{},
		asyncProgresses: map[string]*asyncProgress{},
	}
	go func() {
		s.runServer(true, "localhost:3339", false, "", false, bufio.NewReader(bytes.NewBuffer([]byte{})))
	}()

	time.Sleep(time.Second) // hack so we have time for server to startup. Probably could repeatedly curl the endpoint until it responds.
}

// TestIndividualSim is just a smoke test to make sure the http server works as expected.
//
//	Don't modify this test unless the proto defintions change and this no longer compiles.
func TestIndividualSim(t *testing.T) {
	req := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:      proto.Race_RaceTroll,
				Class:     proto.Class_ClassShaman,
				Equipment: p1Equip,
				Spec:      basicSpec,
			},
			&proto.PartyBuffs{},
			&proto.RaidBuffs{},
			&proto.Debuffs{}),
		Encounter: &proto.Encounter{
			Duration: 120,
			Targets: []*proto.Target{
				{},
			},
		},
		SimOptions: &proto.SimOptions{
			Iterations: 5000,
			RandomSeed: 1,
			Debug:      false,
		},
	}

	msgBytes, err := googleProto.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to encode request: %s", err.Error())
	}

	r, err := http.Post("http://localhost:3339/raidSim", "application/x-protobuf", bytes.NewReader(msgBytes))
	if err != nil {
		t.Fatalf("Failed to POST request: %s", err.Error())
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Failed to read result body: %s", err.Error())
		return
	}

	rsr := &proto.RaidSimResult{}
	if err := googleProto.Unmarshal(body, rsr); err != nil {
		t.Fatalf("Failed to parse request: %s", err.Error())
		return
	}

	log.Printf("RESULT: %#v", rsr)
}
