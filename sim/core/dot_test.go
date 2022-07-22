package core

import (
	"testing"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

func TestDotSnapshot(t *testing.T) {
	sim := NewSim(proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{},
		Raid:       &proto.Raid{},
		Encounter: &proto.Encounter{
			Targets: []*proto.Target{
				{Name: "target", Level: 83, MobType: proto.MobType_MobTypeDemon},
			},
			Duration: 180,
		},
	})

	sim.Raid.Parties = []*Party{
		{
			Raid:           sim.Raid,
			Index:          0,
			Players:        []Agent{},
			Pets:           []PetAgent{},
			PlayersAndPets: []Agent{},
			dpsMetrics:     NewDistributionMetrics(),
		},
	}
	caster := NewCharacter(sim.Raid.Parties[0], 0, proto.Player{Equipment: &proto.EquipmentSpec{}})
	caster.Unit.Env = sim.Environment
	caster.Unit.updateCastSpeed()

	target := &sim.Encounter.Targets[0].Unit

	var dot *Dot

	spell := caster.RegisterSpell(SpellConfig{
		ActionID:    ActionID{SpellID: 42},
		SpellSchool: SpellSchoolShadow,
		Cast:        CastConfig{},
		ApplyEffects: ApplyEffectFuncDirectDamage(SpellEffect{
			ProcMask:       ProcMaskSpellDamage,
			OutcomeApplier: caster.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
				dot.Apply(sim)
			},
		}),
	})

	dot = NewDot(Dot{
		Spell: spell,
		Aura: target.RegisterAura(Aura{
			Label:    "fakdot",
			ActionID: ActionID{SpellID: 42},
		}),
		NumberOfTicks:       6,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: true,
		TickEffects: TickFuncSnapshot(target, SpellEffect{
			ProcMask:             ProcMaskPeriodicDamage,
			ThreatMultiplier:     1,
			BaseDamage:           BaseDamageConfigMagicNoRoll(1000/6, 1),
			BonusSpellCritRating: 3 * CritRatingPerCritChance,
			DamageMultiplier:     1.5,
			OutcomeApplier:       caster.OutcomeFuncMagicCrit(2),
			IsPeriodic:           true,
		}),
	})

	dot.Apply(sim)
	sim.advance(time.Second * 3)
	sim.pendingActions[0].OnAction(sim)
	// validate damage

	dot.Rollover(sim)
	sim.pendingActions[0].OnAction(sim)
	// validate damage
}

func TestDotRollover(t *testing.T) {

}
