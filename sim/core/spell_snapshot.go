package core

import (
	"time"
)

type procTracker struct {
	aura        *Aura
	didActivate bool
	isActive    bool
	expiresAt   time.Duration
}

type SnapshotManager struct {
	procTrackers   []*procTracker
	majorCooldowns []*MajorCooldown
	character      *Character
}

func NewSnapshotManager(character *Character) *SnapshotManager {
	return &SnapshotManager{
		procTrackers:   make([]*procTracker, 0),
		majorCooldowns: make([]*MajorCooldown, 0),
		character:      character,
	}
}

func (manager *SnapshotManager) AddProc(id int32, label string, isActive bool) bool {
	character := manager.character

	if !character.HasAura(label) {
		return false
	}

	manager.procTrackers = append(manager.procTrackers, &procTracker{
		didActivate: false,
		isActive:    isActive,
		expiresAt:   -1,
		aura:        character.GetAura(label),
	})
	return true
}

func (manager *SnapshotManager) CanSnapShot(sim *Simulation, castTime time.Duration) bool {
	success := true

	for _, procTracker := range manager.procTrackers {
		if !procTracker.didActivate && procTracker.aura.IsActive() {
			procTracker.didActivate = true
			procTracker.expiresAt = procTracker.aura.ExpiresAt()
		}

		// A proc is about to drop
		if procTracker.didActivate && procTracker.expiresAt <= sim.CurrentTime+castTime {
			if sim.Log != nil {
				sim.Log("Proc dropping " + procTracker.aura.Label)
			}
			return true
		}

		if !procTracker.didActivate && !procTracker.isActive {
			success = false
		}
	}

	return success
}

func (manager *SnapshotManager) ActivateMajorCooldowns(sim *Simulation) {
	for _, majorCd := range manager.majorCooldowns {
		if majorCd.IsReady(sim) {
			majorCd.TryActivate(sim, manager.character)
		}
	}
}

func (manager *SnapshotManager) ResetProcTrackers() {
	for _, procTracker := range manager.procTrackers {
		procTracker.didActivate = false
		procTracker.expiresAt = -1
	}
}

func (manager *SnapshotManager) ClearMajorCooldowns() {
	manager.majorCooldowns = make([]*MajorCooldown, 0)
}

func (manager *SnapshotManager) AddMajorCooldown(majorCd *MajorCooldown) {
	manager.majorCooldowns = append(manager.majorCooldowns, majorCd)
}
