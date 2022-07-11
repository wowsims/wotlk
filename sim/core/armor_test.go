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
	sunderAura := SunderArmorAura(&target, stacks)
	sunderAura.Activate(&sim)
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
	acidSpitAura := AcidSpitAura(&target, stacks)
	acidSpitAura.Activate(&sim)
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
	acidSpitAura := AcidSpitAura(&target, stacks)
	acidSpitAura.Activate(&sim)
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
	acidSpitAura := AcidSpitAura(&target, stacks)
	acidSpitAura.Activate(&sim)
	tolerance := 0.001
	expectedArmor = baseArmor * (1.0 - float64(stacks)*0.1)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}
	faerieFireAura := FaerieFireAura(&target, 5)
	faerieFireAura.Activate(&sim)
	expectedArmor = baseArmor * (1.0 - 0.2) * (1.0 - 0.05)
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
	target.stats = target.initialStats
	expectedDamageReduction := 0.41132
	attackTable := NewAttackTable(&attacker, &target)
	tolerance := 0.0001
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.ArmorDamageModifier, tolerance) {
		t.Fatalf("Expected no armor modifiers to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.ArmorDamageModifier)
	}

	// Major
	acidSpitAura := AcidSpitAura(&target, 2)
	acidSpitAura.Activate(&sim)
	expectedDamageReduction = 0.3585
	attackTable.UpdateArmorDamageReduction()
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.ArmorDamageModifier, tolerance) {
		t.Fatalf("Expected major armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.ArmorDamageModifier)
	}

	// Major + Minor
	faerieFireAura := FaerieFireAura(&target, 3)
	faerieFireAura.Activate(&sim)
	attackTable.UpdateArmorDamageReduction()
	expectedDamageReduction = 0.3468
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.ArmorDamageModifier, tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.ArmorDamageModifier)
	}

	// Major + Minor + Spore
	sporeCloudAura := SporeCloudAura(&target)
	sporeCloudAura.Activate(&sim)
	attackTable.UpdateArmorDamageReduction()
	expectedDamageReduction = 0.34
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.ArmorDamageModifier, tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.ArmorDamageModifier)
	}

	// Major + Minor + Spore + Throw
	shatteringThrowAura := ShatteringThrowAura(&target)
	shatteringThrowAura.Activate(&sim)
	attackTable.UpdateArmorDamageReduction()
	expectedDamageReduction = 0.2918
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.ArmorDamageModifier, tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.ArmorDamageModifier)
	}

	// Just Major minor again; testing Deactivate
	sporeCloudAura.Deactivate(&sim)
	shatteringThrowAura.Deactivate(&sim)
	attackTable.UpdateArmorDamageReduction()
	expectedDamageReduction = 0.3468
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.ArmorDamageModifier, tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.ArmorDamageModifier)
	}

	// Cap armor pen
	attacker.stats[stats.ArmorPenetration] = 1400
	attackTable.UpdateArmorDamageReduction()
	expectedDamageReduction = 0.0203
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.ArmorDamageModifier, tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.ArmorDamageModifier)
	}

	// Verify going past Cap doesn't help
	attacker.stats[stats.ArmorPenetration] = 1600
	attackTable.UpdateArmorDamageReduction()
	expectedDamageReduction = 0.0203
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.ArmorDamageModifier, tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.ArmorDamageModifier)
	}

	// Add spore back
	sporeCloudAura.Activate(&sim)
	attackTable.UpdateArmorDamageReduction()
	expectedDamageReduction = 0.0100
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.ArmorDamageModifier, tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.ArmorDamageModifier)
	}

	// Fully debuffs
	shatteringThrowAura.Activate(&sim)
	attackTable.UpdateArmorDamageReduction()
	expectedDamageReduction = 0.0
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.ArmorDamageModifier, tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.ArmorDamageModifier)
	}

}
