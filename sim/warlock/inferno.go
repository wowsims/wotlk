package warlock

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerInfernoSpell() {
	if !warlock.Rotation.UseInfernal {
		return
	}

	summonInfernalAura := warlock.RegisterAura(core.Aura{
		Label:    "Summon Infernal",
		ActionID: core.ActionID{SpellID: 1122},
		Duration: time.Second * 60,
	})

	warlock.Inferno = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1122},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Millisecond * 1500,
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * time.Duration(600),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   warlock.SpellCritMultiplier(1, 0),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: add fire spell damage
			baseDmg := (200 + 1*spell.SpellPower()) * sim.Encounter.AOECapMultiplier()

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDmg, spell.OutcomeMagicHitAndCrit)
			}

			if warlock.Pet != nil {
				warlock.Pet.Disable(sim)
			}
			warlock.Infernal.EnableWithTimeout(sim, warlock.Infernal, time.Second*60)

			// fake aura to show duration
			summonInfernalAura.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: warlock.Inferno,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return sim.GetRemainingDuration() <= 61*time.Second
		},
	})
}

type InfernalPet struct {
	core.Pet
	owner          *Warlock
	immolationAura *core.Spell
}

func (warlock *Warlock) NewInfernal() *InfernalPet {
	statInheritance :=
		func(ownerStats stats.Stats) stats.Stats {
			ownerHitChance := math.Floor(ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance)

			// TODO: account for fire spell damage
			return stats.Stats{
				stats.Stamina:          ownerStats[stats.Stamina] * 0.75,
				stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
				stats.Armor:            ownerStats[stats.Armor] * 0.35,
				stats.AttackPower:      ownerStats[stats.SpellPower] * 0.57,
				stats.SpellPower:       ownerStats[stats.SpellPower] * 0.15,
				stats.SpellPenetration: ownerStats[stats.SpellPenetration],
				stats.MeleeHit:         ownerHitChance * core.MeleeHitRatingPerHitChance,
				stats.SpellHit:         ownerHitChance * core.SpellHitRatingPerHitChance,
				stats.Expertise: (ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance) *
					PetExpertiseScale * core.ExpertisePerQuarterPercentReduction,
			}
		}

	infernal := &InfernalPet{
		Pet: core.NewPet("Infernal", &warlock.Character, stats.Stats{
			stats.Strength:  331,
			stats.Agility:   113,
			stats.Stamina:   361,
			stats.Intellect: 65,
			stats.Spirit:    109,
			stats.Mana:      0,
			stats.MeleeCrit: 3.192 * core.CritRatingPerCritChance,
		}, statInheritance, nil, false, false),
		owner: warlock,
	}

	infernal.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	infernal.AddStat(stats.AttackPower, -20)

	// infernal is classified as a warrior class, so we assume it gets the
	// same agi crit coefficient
	infernal.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance*1/62.5)

	// command doesn't apply to infernal
	if warlock.Race == proto.Race_RaceOrc {
		infernal.PseudoStats.DamageDealtMultiplier /= 1.05
	}

	infernal.EnableAutoAttacks(infernal, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  330,
			BaseDamageMax:  494.9,
			SwingSpeed:     2,
			SwingDuration:  time.Second * 2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})
	infernal.AutoAttacks.MHConfig.DamageMultiplier *= 3.2

	core.ApplyPetConsumeEffects(&infernal.Character, warlock.Consumes)

	warlock.AddPet(infernal)

	return infernal
}

func (infernal *InfernalPet) GetPet() *core.Pet {
	return &infernal.Pet
}

func (infernal *InfernalPet) Initialize() {
	felarmor_coef := core.TernaryFloat64(infernal.owner.Options.Armor == proto.Warlock_Options_FelArmor,
		0.3*(1+float64(infernal.owner.Talents.DemonicAegis)*0.1), 0)

	infernal.immolationAura = infernal.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20153},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label:    "Immolation",
				ActionID: core.ActionID{SpellID: 19483},
			},
			NumberOfTicks:       31,
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: false,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// TODO: use highest SP amount of all schools
				// base formula is 25 + (lvl-50)*0.5 * Warlock_SP*0.2
				// note this scales with the warlocks SP, NOT with the pets

				// we remove all the spirit based sp since immolation aura doesn't benefit from it, see
				// JamminL/wotlk-classic-bugs#329
				coef := core.TernaryFloat64(infernal.owner.GlyphOfLifeTapAura.IsActive(), 0.2, 0) + felarmor_coef

				warlockSP := infernal.owner.Unit.GetStat(stats.SpellPower) - infernal.owner.Unit.GetStat(stats.Spirit)*coef
				baseDmg := (40 + warlockSP*0.2) * sim.Encounter.AOECapMultiplier()

				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, baseDmg, dot.Spell.OutcomeMagicHit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}

func (infernal *InfernalPet) Reset(sim *core.Simulation) {
}

func (infernal *InfernalPet) OnGCDReady(sim *core.Simulation) {
	infernal.immolationAura.Cast(sim, nil)
	infernal.DoNothing()
}
