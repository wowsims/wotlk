package core

import (
	"fmt"
	"log"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type TargetAI interface {
	Initialize(*Target)

	DoAction(*Simulation)
}

func (target *Target) initialize(config *proto.Target) {
	if config == nil || target.CurrentTarget == nil {
		return
	}

	if config.SwingSpeed > 0 {
		aaOptions := AutoAttackOptions{
			MainHand: Weapon{
				BaseDamageMin:  config.MinBaseDamage,
				SwingSpeed:     config.SwingSpeed,
				SwingDuration:  time.Duration(float64(time.Second) * config.SwingSpeed),
				CritMultiplier: 2,
				SpellSchool:    SpellSchoolFromProto(config.SpellSchool),
			},
			AutoSwingMelee: true,
		}
		if config.DualWield {
			aaOptions.OffHand = aaOptions.MainHand
			if !config.DualWieldPenalty {
				target.PseudoStats.DisableDWMissPenalty = true
			}
		}
		target.EnableAutoAttacks(target, aaOptions)
	}

	if target.AI != nil {
		target.AI.Initialize(target)

		target.gcdAction = &PendingAction{
			Priority: ActionPriorityGCD,
			OnAction: func(sim *Simulation) {
				if target.GCD.IsReady(sim) {
					target.OnGCDReady(sim)

					if !target.doNothing && target.GCD.IsReady(sim) && (!target.IsWaiting() && !target.IsWaitingForMana()) {
						msg := fmt.Sprintf("Target `%s` did not perform any actions. Either this is a bug or agent should use 'WaitUntil' or 'WaitForMana' to explicitly wait.\n\tIf character has no action to perform use 'DoNothing'.", target.Label)
						panic(msg)
					}
					target.doNothing = false
				}
			},
		}
	}
}

// Empty Agent interface functions.
func (target *Target) AddRaidBuffs(raidBuffs *proto.RaidBuffs)    {}
func (target *Target) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {}
func (target *Target) ApplyGearBonuses()                          {}
func (target *Target) ApplyTalents()                              {}
func (target *Target) GetCharacter() *Character                   { return nil }
func (target *Target) Initialize()                                {}

func (target *Target) DoNothing() {
	target.doNothing = true
}

func (target *Target) OnAutoAttack(sim *Simulation, spell *Spell) {
	if target.GCD.IsReady(sim) {
		if target.AI != nil {
			target.AI.DoAction(sim)
		}
	}
}
func (target *Target) OnGCDReady(sim *Simulation) {
	if target.AI != nil {
		target.AI.DoAction(sim)
	}
}

type AIFactory func() TargetAI

type PresetTarget struct {
	// String in folder-structure format identifying a category for this unit, e.g. "Black Temple/Bosses".
	PathPrefix string

	Config proto.Target

	AI AIFactory
}

func (pt PresetTarget) Path() string {
	return pt.PathPrefix + "/" + pt.Config.Name
}
func (pt PresetTarget) ToProto() *proto.PresetTarget {
	target := &proto.Target{}
	*target = pt.Config

	return &proto.PresetTarget{
		Path:   pt.Path(),
		Target: target,
	}
}

var presetTargets = []PresetTarget{}
var presetEncounters = []*proto.PresetEncounter{}

func AddPresetTarget(newPreset PresetTarget) {
	for _, preset := range presetTargets {
		if preset.Path() == newPreset.Path() {
			log.Fatalf("Preset Target with path %s already added!", newPreset.Path())
		}
	}
	presetTargets = append(presetTargets, newPreset)
}

func GetPresetTargetWithPath(path string) *PresetTarget {
	for i, _ := range presetTargets {
		preset := &presetTargets[i]
		if preset.Path() == path {
			return preset
		}
	}
	return nil
}

func GetPresetTargetWithID(id int32) *PresetTarget {
	for i, _ := range presetTargets {
		preset := &presetTargets[i]
		if preset.Config.Id == id {
			return preset
		}
	}
	return nil
}

func AddPresetEncounter(name string, targetPaths []string) {
	if len(targetPaths) == 0 {
		log.Fatalf("Encounter must have targets!")
	}

	var path string
	var targetProtos []*proto.PresetTarget

	for i, targetPath := range targetPaths {
		presetTarget := GetPresetTargetWithPath(targetPath)
		if presetTarget == nil {
			log.Fatalf("No preset target with path: %s", targetPath)
		}
		targetProtos = append(targetProtos, presetTarget.ToProto())

		if i == 0 {
			path = presetTarget.PathPrefix + "/" + name
		}
	}

	for _, preset := range presetEncounters {
		if preset.Path == path {
			log.Fatalf("Preset Encounter with path %s already added!", path)
		}
	}

	presetEncounters = append(presetEncounters, &proto.PresetEncounter{
		Path:    path,
		Targets: targetProtos,
	})
}
