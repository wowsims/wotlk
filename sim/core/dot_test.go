package core

import (
	"testing"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

func NewFakeElementalShaman(char Character, options proto.Player) Agent {
	fa := &FakeAgent{
		Character: char,
	}

	fa.Init = func() {
		fa.Spell = fa.RegisterSpell(SpellConfig{
			ActionID:    ActionID{SpellID: 42},
			SpellSchool: SpellSchoolShadow,
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
				BaseDamage:           BaseDamageConfigMagicNoRoll(1000/6, 1),
				BonusSpellCritRating: 3 * CritRatingPerCritChance,
				DamageMultiplier:     1.5,
				OutcomeApplier:       fa.OutcomeFuncMagicCrit(2),
				IsPeriodic:           true,
			}),
		})
	}

	return fa
}

func TestDotSnapshot(t *testing.T) {
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

	fa := sim.Raid.Parties[0].Players[0].(*FakeAgent)

	fa.Dot.Apply(sim)
	sim.advance(time.Second * 3)
	sim.pendingActions[0].OnAction(sim)

	expectedDmg := 373.5
	if !WithinToleranceFloat64(expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage, 0.01) {
		t.Fatalf("Incorrect damage applied: Expected: %0.3f, Actual: %0.3f", expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage)
	}

	fa.Dot.Rollover(sim)
	sim.advance(time.Second * 3)
	sim.pendingActions[0].OnAction(sim)

	expectedDmg = 373.5 + 249.0
	if !WithinToleranceFloat64(expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage, 0.01) {
		t.Fatalf("Incorrect damage applied: Expected: %0.3f, Actual: %0.3f", expectedDmg, fa.Dot.Spell.SpellMetrics[0].TotalDamage)
	}

}

func TestDotRollover(t *testing.T) {

}
