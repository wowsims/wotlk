package core

import (
	"testing"

	"github.com/wowsims/wotlk/sim/core/stats"
)

func TestSunderArmorStacks(t *testing.T) {
	sim := Simulation{}
	baseArmor := 10643.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        83,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	target.stats = target.initialStats
	expectedArmor := baseArmor
	if target.Armor() != expectedArmor {
		t.Fatalf("Armor value for target should be %f but found %f", 10643.0, target.Armor())
	}
	stacks := int32(1)
	sunderAura := SunderArmorAura(&target)
	sunderAura.Activate(&sim)
	sunderAura.SetStacks(&sim, stacks)
	tolerance := 0.001
	for stacks <= 5 {
		expectedArmor = baseArmor * (1.0 - float64(stacks)*0.04)
		if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
			t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
		}
		stacks++
		sunderAura.AddStack(&sim)
	}
}

func TestAcidSpitStacks(t *testing.T) {
	sim := Simulation{}
	baseArmor := 10643.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        83,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	target.stats = target.initialStats
	expectedArmor := baseArmor
	if target.Armor() != expectedArmor {
		t.Fatalf("Armor value for target should be %f but found %f", 10643.0, target.Armor())
	}
	stacks := int32(1)
	acidSpitAura := AcidSpitAura(&target)
	acidSpitAura.Activate(&sim)
	acidSpitAura.SetStacks(&sim, stacks)
	tolerance := 0.001
	for stacks <= 2 {
		expectedArmor = baseArmor * (1.0 - float64(stacks)*0.1)
		if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
			t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
		}
		stacks++
		acidSpitAura.AddStack(&sim)
	}
}

func TestExposeArmor(t *testing.T) {
	sim := Simulation{}
	baseArmor := 10643.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        83,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	target.stats = target.initialStats
	expectedArmor := baseArmor
	if target.Armor() != expectedArmor {
		t.Fatalf("Armor value for target should be %f but found %f", 10643.0, target.Armor())
	}
	exposeAura := ExposeArmorAura(&target, false)
	exposeAura.Activate(&sim)
	tolerance := 0.001
	expectedArmor = baseArmor * (1.0 - 0.2)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}
}

func TestMajorArmorReductionAurasDoNotStack(t *testing.T) {
	sim := Simulation{}
	baseArmor := 10643.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        83,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	target.stats = target.initialStats
	expectedArmor := baseArmor
	if target.Armor() != expectedArmor {
		t.Fatalf("Armor value for target should be %f but found %f", 10643.0, target.Armor())
	}
	stacks := int32(1)
	acidSpitAura := AcidSpitAura(&target)
	acidSpitAura.Activate(&sim)
	acidSpitAura.SetStacks(&sim, stacks)
	tolerance := 0.001
	expectedArmor = baseArmor * (1.0 - float64(stacks)*0.1)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}
	exposeArmorAura := ExposeArmorAura(&target, false)
	exposeArmorAura.Activate(&sim)
	expectedArmor = baseArmor * (1.0 - 0.2)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}
}

func TestMajorAndMinorArmorReductionsApplyMultiplicatively(t *testing.T) {
	sim := Simulation{}
	baseArmor := 10643.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        83,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	target.stats = target.initialStats
	expectedArmor := baseArmor
	if target.Armor() != expectedArmor {
		t.Fatalf("Armor value for target should be %f but found %f", 10643.0, target.Armor())
	}
	stacks := int32(2)
	acidSpitAura := AcidSpitAura(&target)
	acidSpitAura.Activate(&sim)
	acidSpitAura.SetStacks(&sim, stacks)
	tolerance := 0.001
	expectedArmor = baseArmor * (1.0 - float64(stacks)*0.1)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}
	faerieFireAura := FaerieFireAura(&target, 0)
	impFaerieFireAura := FaerieFireAura(&target, 3)
	faerieFireAura.Activate(&sim)
	expectedArmor = baseArmor * (1.0 - 0.2) * (1.0 - 0.05)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}
	impFaerieFireAura.Activate(&sim)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}
	impFaerieFireAura.Deactivate(&sim)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}

	faerieFireAura.Deactivate(&sim)
	expectedArmor = baseArmor * (1.0 - 0.2)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}
}

func TestDamageReductionFromArmor(t *testing.T) {
	sim := Simulation{}
	baseArmor := 10643.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        83,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	attacker := Unit{
		Type:  PlayerUnit,
		Level: 80,
	}
	spell := &Spell{}
	target.stats = target.initialStats
	expectedDamageReduction := 0.41132
	attackTable := NewAttackTable(&attacker, &target)
	tolerance := 0.0001
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected no armor modifiers to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

	// Major
	acidSpitAura := AcidSpitAura(&target)
	acidSpitAura.Activate(&sim)
	acidSpitAura.SetStacks(&sim, 2)
	expectedDamageReduction = 0.3585
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

	// Major + Minor
	faerieFireAura := FaerieFireAura(&target, 3)
	faerieFireAura.Activate(&sim)
	expectedDamageReduction = 0.3468
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

	// Major + Minor + Spore
	sporeCloudAura := SporeCloudAura(&target)
	sporeCloudAura.Activate(&sim)
	expectedDamageReduction = 0.3468
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

	// Major + Minor + Spore + Throw
	shatteringThrowAura := ShatteringThrowAura(&target)
	shatteringThrowAura.Activate(&sim)
	expectedDamageReduction = 0.2981
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

	// Just Major minor again; testing Deactivate
	sporeCloudAura.Deactivate(&sim)
	shatteringThrowAura.Deactivate(&sim)
	expectedDamageReduction = 0.3468
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

	// Cap armor pen
	attacker.stats[stats.ArmorPenetration] = 1400
	expectedDamageReduction = 0.02026
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

	// Verify going past Cap doesn't help
	attacker.stats[stats.ArmorPenetration] = 1600
	expectedDamageReduction = 0.02026
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

	// Add spore back
	sporeCloudAura.Activate(&sim)
	expectedDamageReduction = 0.02026
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

	// Fully debuffs
	shatteringThrowAura.Activate(&sim)
	expectedDamageReduction = 0.0
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

}
