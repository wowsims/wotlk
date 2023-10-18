package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type Shadowfiend struct {
	core.Pet

	Priest          *Priest
	ManaMetric      *core.ResourceMetrics
	Shadowcrawl     *core.Spell
	ShadowcrawlAura *core.Aura
}

var baseStats = stats.Stats{
	stats.Strength:    314,
	stats.Agility:     90,
	stats.Stamina:     348,
	stats.Intellect:   201,
	stats.AttackPower: -20,
	// with 3% crit debuff, shadowfiend crits around 9-12% (TODO: verify and narrow down)
	stats.MeleeCrit: 8 * core.CritRatingPerCritChance,
}

func (priest *Priest) NewShadowfiend() *Shadowfiend {
	shadowfiend := &Shadowfiend{
		Pet:    core.NewPet("Shadowfiend", &priest.Character, baseStats, priest.shadowfiendStatInheritance(), false, false),
		Priest: priest,
	}

	shadowfiend.ManaMetric = priest.NewManaMetrics(core.ActionID{SpellID: 34433})
	actionID := core.ActionID{SpellID: 63619}

	shadowfiend.ShadowcrawlAura = shadowfiend.GetOrRegisterAura(core.Aura{
		Label:    "Shadowcrawl",
		ActionID: actionID,
		Duration: time.Second * 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shadowfiend.PseudoStats.DamageDealtMultiplier *= 1.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shadowfiend.PseudoStats.DamageDealtMultiplier /= 1.15
		},
	})

	shadowfiend.Shadowcrawl = shadowfiend.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolMagic,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoLogs,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shadowfiend.ShadowcrawlAura.Activate(sim)
		},
	})

	shadowfiend.PseudoStats.DamageTakenMultiplier *= 0.1

	shadowfiend.EnableAutoAttacks(shadowfiend, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:        110,
			BaseDamageMax:        145,
			SwingSpeed:           1.5,
			NormalizedSwingSpeed: 1.5,
			CritMultiplier:       2,
			SpellSchool:          core.SpellSchoolShadow,
		},
		AutoSwingMelee: true,
	})

	shadowfiend.AddStatDependency(stats.Strength, stats.AttackPower, 2.0)

	core.ApplyPetConsumeEffects(&shadowfiend.Character, priest.Consumes)

	priest.AddPet(shadowfiend)

	return shadowfiend
}

func (priest *Priest) shadowfiendStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		inheritableSP := ownerStats[stats.SpellPower] - 0.04*float64(priest.Talents.TwistedFaith)*ownerStats[stats.Spirit]
		// Shadow fiend gets a "Spell Bonus" that adds bonus damage to each attack
		// for simplicity, we will just convert this added damage as if it were AP
		// Spell Bonus SP coefficient: 30%
		// Spell Bonus Damage coefficient: 106%
		// Damage to DPS coefficient: 1/1.5 (1.5 speed weapon)
		// DPS to AP coefficient: 14
		spellBonusAPEquivalent := inheritableSP * 0.3 * 1.06 * 14 / 1.5

		return stats.Stats{ //still need to nail down shadow fiend crit scaling, but removing owner crit scaling after further investigation
			stats.AttackPower: inheritableSP*0.57 + spellBonusAPEquivalent,
			// never misses
			stats.MeleeHit:  8 * core.MeleeHitRatingPerHitChance,
			stats.Expertise: 14 * core.ExpertisePerQuarterPercentReduction * 4,
		}
	}
}

func (shadowfiend *Shadowfiend) OnAutoAttack(sim *core.Simulation, _ *core.Spell) {
	priest := shadowfiend.Priest
	restoreMana := priest.MaxMana() * 0.05

	priest.AddMana(sim, restoreMana, shadowfiend.ManaMetric)
}

func (shadowfiend *Shadowfiend) Initialize() {
}

func (shadowfiend *Shadowfiend) OnGCDReady(sim *core.Simulation) {
	if shadowfiend.Shadowcrawl.IsReady(sim) {
		shadowfiend.Shadowcrawl.Cast(sim, nil)
	} else {
		shadowfiend.DoNothing()
	}
}

func (shadowfiend *Shadowfiend) Reset(sim *core.Simulation) {
	shadowfiend.ShadowcrawlAura.Deactivate(sim)
	shadowfiend.Disable(sim)
}

func (shadowfiend *Shadowfiend) OnPetDisable(sim *core.Simulation) {
	shadowfiend.ShadowcrawlAura.Deactivate(sim)
}

func (shadowfiend *Shadowfiend) GetPet() *core.Pet {
	return &shadowfiend.Pet
}
