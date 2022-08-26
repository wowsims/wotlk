package feral

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/druid"
)

func RegisterFeralDruid() {
	core.RegisterAgentFactory(
		proto.Player_FeralDruid{},
		proto.Spec_SpecFeralDruid,
		func(character core.Character, options proto.Player) core.Agent {
			return NewFeralDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FeralDruid)
			if !ok {
				panic("Invalid spec value for Feral Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFeralDruid(character core.Character, options proto.Player) *FeralDruid {
	feralOptions := options.GetFeralDruid()

	selfBuffs := druid.SelfBuffs{}
	if feralOptions.Options.InnervateTarget != nil {
		selfBuffs.InnervateTarget = *feralOptions.Options.InnervateTarget
	} else {
		selfBuffs.InnervateTarget.TargetIndex = -1
	}

	cat := &FeralDruid{
		Druid:   druid.New(character, druid.Cat, selfBuffs, *feralOptions.Talents),
		latency: time.Duration(feralOptions.Options.LatencyMs) * time.Millisecond,
	}

	cat.setupRotation(feralOptions.Rotation)

	// Passive Cat Form threat reduction
	cat.PseudoStats.ThreatMultiplier *= 0.71

	// Prevents Windfury from applying.
	cat.HasMHWeaponImbue = true

	cat.EnableEnergyBar(100.0, func(sim *core.Simulation) {
		if cat.GCD.IsReady(sim) {
			cat.doRotation(sim)
		}
	})

	cat.EnableAutoAttacks(cat, core.AutoAttackOptions{
		// Base paw weapon.
		MainHand: core.Weapon{
			BaseDamageMin:        72,
			BaseDamageMax:        95,
			SwingSpeed:           1.0,
			NormalizedSwingSpeed: 1.0,
			SwingDuration:        time.Second,
			CritMultiplier:       cat.MeleeCritMultiplier(),
		},
		AutoSwingMelee: true,
	})

	return cat
}

type FeralDruid struct {
	*druid.Druid

	Rotation FeralDruidRotation

	readyToShift   bool
	waitingForTick bool
	latency        time.Duration
}

func (cat *FeralDruid) GetDruid() *druid.Druid {
	return cat.Druid
}

func (cat *FeralDruid) Initialize() {
	cat.Druid.Initialize()
	cat.RegisterFeralSpells(15)
	cat.DelayDPSCooldownsForArmorDebuffs()
}

func (cat *FeralDruid) Reset(sim *core.Simulation) {
	cat.Druid.Reset(sim)
	cat.Druid.ClearForm(sim)
	cat.CatFormAura.Activate(sim)
	cat.readyToShift = false
	cat.waitingForTick = false
}
