package core

import (
	"testing"
	"time"
)

func TestSingleAuraExclusiveDurationNoOverwrite(t *testing.T) {
	sim := &Simulation{}

	target := Unit{
		Type:        EnemyUnit,
		Index:       0,
		Level:       83,
		auraTracker: newAuraTracker(),
	}
	mangle := MangleAura(&target)
	trauma := MakePermanent(TraumaAura(&target, 2))

	// Trauma in this case should *never* be overwritten
	// as its duration from 'MakePermanent' should make it non overwritable by 1 min duration mangles
	trauma.Activate(sim)

	sim.CurrentTime = 1 * time.Second

	mangle.Activate(sim)

	if !(trauma.IsActive() && !mangle.IsActive()) {
		t.Fatalf("lower duration exclusive aura overwrote previous!")
	}
}

func TestSingleAuraExclusiveDurationOverwrite(t *testing.T) {
	sim := &Simulation{}

	target := Unit{
		Type:        EnemyUnit,
		Index:       0,
		Level:       83,
		auraTracker: newAuraTracker(),
	}
	mangle := MangleAura(&target)
	trauma := TraumaAura(&target, 2)

	trauma.Activate(sim)

	sim.CurrentTime = 1 * time.Second

	mangle.Activate(sim)

	// In this case mangle should overwrite trauma as mangle will give a greater duration

	if !(mangle.IsActive() && !trauma.IsActive()) {
		t.Fatalf("longer duration exclusive aura failed to overwrite")
	}
}
