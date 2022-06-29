package core

// Pets used by core effects/buffs/consumes.

import (
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Creates any pets that come from item / consume effects.
func (character *Character) addEffectPets() {
	if character.Consumes.FillerExplosive == proto.Explosive_ExplosiveGnomishFlameTurret {
		character.NewGnomishFlameTurret()
	}
}

func (character *Character) newGnomishFlameTurretSpell() *Spell {
	gft := character.GetPet(GnomishFlameTurretName).(*GnomishFlameTurret)

	return character.GetOrRegisterSpell(SpellConfig{
		ActionID: ActionID{ItemID: 23841},

		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			gft.EnableWithTimeout(sim, gft, time.Second*45)
		},
	})
}

const GnomishFlameTurretName = "Gnomish Flame Turret"

type GnomishFlameTurret struct {
	Pet

	FlameCannon *Spell
}

func (character *Character) NewGnomishFlameTurret() *GnomishFlameTurret {
	gft := &GnomishFlameTurret{
		Pet: NewPet(
			GnomishFlameTurretName,
			character,
			stats.Stats{
				stats.SpellCrit: 1 * SpellCritRatingPerCritChance,
			},
			func(ownerStats stats.Stats) stats.Stats {
				return stats.Stats{}
			},
			false,
		),
	}

	character.AddPet(gft)

	return gft
}

func (gft *GnomishFlameTurret) GetPet() *Pet {
	return &gft.Pet
}

func (gft *GnomishFlameTurret) Initialize() {
	gft.registerFlameCannonSpell()
}

func (gft *GnomishFlameTurret) registerFlameCannonSpell() {
	gft.FlameCannon = gft.RegisterSpell(SpellConfig{
		ActionID:    ActionID{SpellID: 30527},
		SpellSchool: SpellSchoolFire,

		Cast: CastConfig{
			DefaultCast: Cast{
				// Pretty sure this works the same way as Searing Totem, where the next shot
				// fires once the previous missile has hit the target. Just give some static
				// value for now.
				GCD: time.Millisecond * 800,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: ApplyEffectFuncDirectDamage(SpellEffect{
			ProcMask:         ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage:     BaseDamageConfigRoll(31, 36),
			OutcomeApplier: gft.OutcomeFuncMagicHitAndCrit(gft.DefaultSpellCritMultiplier()),
		}),
	})
}

func (gft *GnomishFlameTurret) Reset(sim *Simulation) {
}

func (gft *GnomishFlameTurret) OnGCDReady(sim *Simulation) {
	gft.FlameCannon.Cast(sim, gft.CurrentTarget)
}
