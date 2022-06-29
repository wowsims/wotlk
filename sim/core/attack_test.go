package core

import (
	//"math/rand"
	//"testing"

	//"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

// These tests are currently broken, seems like target PseudoStats are somehow not being set.
//func TestAutoSwing(t *testing.T) {
//	a := &FakeAgent{}
//	c := &Character{
//		Metrics: NewCharacterMetrics(),
//		Equip: items.Equipment{
//			proto.ItemSlot_ItemSlotMainHand: items.ByID[32262],
//			proto.ItemSlot_ItemSlotOffHand:  items.ByID[32262],
//		},
//	}
//	sim := &Simulation{
//		rand:    rand.New(rand.NewSource(1)),
//		Options: proto.SimOptions{},
//		encounter: Encounter{
//			Targets: []*Target{NewTarget(proto.Target{}, 0)},
//		},
//		isTest:            true,
//		testRands:         make(map[string]*rand.Rand),
//		emptyAuras:        make([]Aura, numAuraIDs),
//		pendingActionPool: newPAPool(),
//	}
//
//	c.EnableAutoAttacks(a, AutoAttackOptions{
//		MainHand: c.WeaponFromMainHand(c.DefaultMeleeCritMultiplier()),
//		OffHand:  c.WeaponFromOffHand(c.DefaultMeleeCritMultiplier()),
//	})
//	c.AutoAttacks.TrySwingMH(sim, sim.GetPrimaryTarget())
//	c.AutoAttacks.TrySwingOH(sim, sim.GetPrimaryTarget())
//
//	metricTests := []struct {
//		key   ActionKey
//		value float64
//	}{
//		{key: NewActionKey(ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 1}), value: 323.355012},
//		{key: NewActionKey(ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 2}), value: 126.601614},
//	}
//	tolerance := 0.0001
//	for _, mt := range metricTests {
//		if c.Metrics.actions[mt.key].Damage < mt.value-tolerance || c.Metrics.actions[mt.key].Damage > mt.value+tolerance {
//			t.Fatalf("Failed... Expected: %f, Actual: %f", mt.value, c.Metrics.actions[mt.key].Damage)
//		}
//	}
//}
//
//func TestRangedAutoSwing(t *testing.T) {
//	a := &FakeAgent{}
//	c := &Character{
//		Metrics: NewCharacterMetrics(),
//		Equip: items.Equipment{
//			proto.ItemSlot_ItemSlotRanged:   items.ByID[28772], // sunfury bow phoenix
//			proto.ItemSlot_ItemSlotMainHand: items.ByID[28435], // mooncleaver
//		},
//	}
//	sim := &Simulation{
//		rand:    rand.New(rand.NewSource(1)),
//		Options: proto.SimOptions{},
//		encounter: Encounter{
//			Targets: []*Target{{}},
//		},
//		isTest:            true,
//		testRands:         make(map[string]*rand.Rand),
//		emptyAuras:        make([]Aura, numAuraIDs),
//		pendingActionPool: newPAPool(),
//	}
//
//	c.EnableAutoAttacks(a, AutoAttackOptions{
//		MainHand: c.WeaponFromMainHand(c.DefaultMeleeCritMultiplier()),
//		Ranged:   c.WeaponFromRanged(0),
//	})
//	c.AutoAttacks.TrySwingMH(sim, sim.GetPrimaryTarget())
//	c.AutoAttacks.Ranged.CritMultiplier = 2.0 // technically hunters actually calculate this.
//	c.AutoAttacks.RangedAuto.CritMultiplier = 2.0
//	// Ranged autos require a windup, so we just skip that here.
//	ama := c.AutoAttacks.RangedAuto
//	ama.Effect.Target = sim.GetPrimaryTarget()
//	ama.Cast(sim)
//
//	metricTests := []struct {
//		name  string
//		key   ActionKey
//		value float64
//	}{
//		{name: "main hand attack", key: NewActionKey(ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 1}), value: 483.630072},
//		{name: "ranged attack", key: NewActionKey(ActionID{OtherID: proto.OtherAction_OtherActionShoot}), value: 218.079693},
//	}
//	tolerance := 0.0001
//
//	for _, mt := range metricTests {
//		if c.Metrics.actions[mt.key].Damage < mt.value-tolerance || c.Metrics.actions[mt.key].Damage > mt.value+tolerance {
//			t.Fatalf("Failed (%s) Expected: %f, Actual: %f", mt.name, mt.value, c.Metrics.actions[mt.key].Damage)
//		}
//	}
//}

// The Character controlled by this Agent.
func (fa *FakeAgent) GetCharacter() *Character {
	panic("not implemented") // TODO: Implement
}

// Updates the input Buffs to include raid-wide buffs provided by this Agent.
func (fa *FakeAgent) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	panic("not implemented") // TODO: Implement
}

// Updates the input Buffs to include party-wide buffs provided by this Agent.
func (fa *FakeAgent) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	panic("not implemented") // TODO: Implement
}

func (fa *FakeAgent) ApplyTalents() {
	panic("not implemented") // TODO: Implement
}

// Called once before the first iteration, after all Agents and Targets are finalized.
// Use this to do any precomputations that require access to Sim or Target fields.
func (fa *FakeAgent) Init(sim *Simulation) {
	panic("not implemented") // TODO: Implement
}

// Returns this Agent to its initial state. Called before each Sim iteration
// and once after the final iteration.
func (fa *FakeAgent) Reset(sim *Simulation) {
	panic("not implemented") // TODO: Implement
}

// Called whenever the GCD becomes ready for this Agent.
func (fa *FakeAgent) OnGCDReady(sim *Simulation) {
	panic("not implemented") // TODO: Implement
}

// Called after each mana tick, if this Agent uses mana.
func (fa *FakeAgent) OnManaTick(sim *Simulation) {
	panic("not implemented") // TODO: Implement
}

type FakeAgent struct {
	Character
}
