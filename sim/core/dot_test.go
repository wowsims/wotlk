package core

import (
	"testing"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	RegisterAgentFactory(
		proto.Player_ElementalShaman{},
		proto.Spec_SpecElementalShaman,
		func(character Character, options proto.Player) Agent {
			return NewFakeElementalShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ElementalShaman)
			if !ok {
				panic("Invalid spec value for Elemental Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFakeElementalShaman(char Character, options proto.Player) Agent {
	fa := &FakeAgent{
		Character: char,
	}

	fa.Init = func() {
		fa.Spell = fa.RegisterSpell(SpellConfig{
			ActionID:    ActionID{SpellID: 42},
			SpellSchool: SpellSchoolShadow,
			Flags:       SpellFlagIgnoreResists,
			Cast:        CastConfig{},
			ApplyEffects: ApplyEffectFuncDirectDamage(SpellEffect{
				ProcMask:       ProcMaskSpellDamage,
				OutcomeApplier: fa.OutcomeFuncMagicHit(),
				OnSpellHitDealt: func(sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
					fa.Dot.Apply(sim)
				},
			}),
		})

		fa.Dot = NewDot(Dot{
			Spell: fa.Spell,
			Aura: fa.CurrentTarget.RegisterAura(Aura{
				Label:    "fakdot",
				ActionID: ActionID{SpellID: 42},
			}),
			NumberOfTicks:       6,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,
			TickEffects: TickFuncSnapshot(fa.CurrentTarget, SpellEffect{
				ProcMask:             ProcMaskPeriodicDamage,
				ThreatMultiplier:     1,
				BaseDamage:           BaseDamageConfigMagicNoRoll(1000.0/6, 1),
				BonusSpellCritRating: 3 * CritRatingPerCritChance,
				DamageMultiplier:     1.5,
				OutcomeApplier:       fa.OutcomeFuncAlwaysHit(),
				IsPeriodic:           true,
			}),
		})
	}

	return fa
}

func SetupFakeSim() *Simulation {
	sim := NewSim(proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{
			RandomSeed: 100,
		},
		Raid: &proto.Raid{
			Parties: []*proto.Party{
				{
					Players: []*proto.Player{
						{
							Name:      "Caster",
							Class:     proto.Class_ClassShaman,
							Consumes:  &proto.Consumes{},
							Buffs:     &proto.IndividualBuffs{},
							Spec:      &proto.Player_ElementalShaman{},
							Equipment: &proto.EquipmentSpec{},
						},
					},
					Buffs: &proto.PartyBuffs{},
				},
			},
		},
		Encounter: &proto.Encounter{
			Targets: []*proto.Target{
				{Name: "target", Level: 83, MobType: proto.MobType_MobTypeDemon},
			},
			Duration: 180,
		},
	})
	sim.Reset()

	return sim
}

func TestDotSnapshot(t *testing.T) {
	sim := SetupFakeSim()
	fa := sim.Raid.Parties[0].Players[0].(*FakeAgent)

	fa.Dot.Apply(sim)
	sim.advance(time.Second * 3)
	sim.pendingActions[0].OnAction(sim)

	// (1000/6) * 1.5 = 250
	expectedDmg := 250.0
	if !WithinToleranceFloat64(expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage, 0.01) {
		t.Fatalf("Incorrect damage applied: Expected: %0.3f, Actual: %0.3f", expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage)
	}

	fa.Dot.Rollover(sim)
	sim.advance(time.Second * 3)
	sim.pendingActions[0].OnAction(sim)

	expectedDmg = 250.0 + 250.0
	if !WithinToleranceFloat64(expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage, 0.01) {
		t.Fatalf("Incorrect damage applied: Expected: %0.3f, Actual: %0.3f", expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage)
	}
}

func TestDotSnapshotSpellPower(t *testing.T) {
	sim := SetupFakeSim()
	fa := sim.Raid.Parties[0].Players[0].(*FakeAgent)

	fa.Dot.Apply(sim)
	sim.advance(time.Second * 3)
	sim.pendingActions[0].OnAction(sim)

	// (1000/6) * 1.5 = 250
	expectedDmg := 250.0
	if !WithinToleranceFloat64(expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage, 0.01) {
		t.Fatalf("Incorrect damage applied: Expected: %0.3f, Actual: %0.3f", expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage)
	}

	fa.GetCharacter().AddStatDynamic(sim, stats.SpellPower, 1000)
	sim.advance(time.Second * 3)
	sim.pendingActions[0].OnAction(sim)

	// Spell power shouldn't get applied because dot was already snapshot.
	expectedDmg = 250.0 + 250.0
	if !WithinToleranceFloat64(expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage, 0.01) {
		t.Fatalf("Incorrect damage applied: Expected: %0.3f, Actual: %0.3f", expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage)
	}

	fa.Dot.Deactivate(sim)
	fa.Dot.Activate(sim)
	sim.advance(time.Second * 3)
	sim.pendingActions[0].OnAction(sim)

	expectedDmg = 250.0 + 250.0 + (250.0 + 1000*1.5)
	if !WithinToleranceFloat64(expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage, 0.01) {
		t.Fatalf("Incorrect damage applied: Expected: %0.3f, Actual: %0.3f", expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage)
	}
}

func TestDotSnapshotSpellMultiplier(t *testing.T) {
	sim := SetupFakeSim()
	fa := sim.Raid.Parties[0].Players[0].(*FakeAgent)
	spell := fa.GetCharacter().Spellbook[0]
	spell.DamageMultiplier = 2

	fa.Dot.Apply(sim)
	sim.advance(time.Second * 3)
	sim.pendingActions[0].OnAction(sim)

	// (1000/6) * 1.5 * 2 = 500
	expectedDmg := 500.0
	if !WithinToleranceFloat64(expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage, 0.01) {
		t.Fatalf("Incorrect damage applied: Expected: %0.3f, Actual: %0.3f", expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage)
	}

	fa.Dot.Rollover(sim)
	sim.advance(time.Second * 3)
	sim.pendingActions[0].OnAction(sim)

	expectedDmg = 500.0 + 500.0
	if !WithinToleranceFloat64(expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage, 0.01) {
		t.Fatalf("Incorrect damage applied: Expected: %0.3f, Actual: %0.3f", expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage)
	}
}

func TestDotRollover(t *testing.T) {

}
