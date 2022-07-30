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
	stats.Strength:  314,
	stats.Agility:   90,
	stats.Stamina:   348,
	stats.Intellect: 201,
}

func (priest *Priest) NewShadowfiend() *Shadowfiend {
	shadowfiend := &Shadowfiend{
		Pet: core.NewPet(
			"Shadowfiend",
			&priest.Character,
			baseStats,
			priest.shadowfiendStatInheritance(),
			false,
			true,
		),
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
		Flags:       core.SpellFlagNoLogs,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 6,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:       core.ProcMaskEmpty,
			OutcomeApplier: priest.OutcomeFuncAlwaysHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				shadowfiend.ShadowcrawlAura.Activate(sim)
			},
		}),
	})

	shadowfiend.PseudoStats.DamageTakenMultiplier *= 0.1

	shadowfiend.EnableAutoAttacks(shadowfiend, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:        176,
			BaseDamageMax:        210,
			SwingSpeed:           1.5,
			NormalizedSwingSpeed: 1.5,
			SwingDuration:        time.Millisecond * 1500,
			CritMultiplier:       2,
			SpellSchool:          core.SpellSchoolShadow,
		},
		AutoSwingMelee: true,
	})

	shadowfiend.AddStatDependency(stats.Strength, stats.AttackPower, 1.0)

	core.ApplyPetConsumeEffects(&shadowfiend.Character, priest.Consumes)

	priest.AddPet(shadowfiend)

	return shadowfiend
}

func (priest *Priest) shadowfiendStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		hitPercentage := ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance

		return stats.Stats{
			stats.AttackPower: ownerStats[stats.SpellPower] * 5.377,
			stats.MeleeHit:    hitPercentage * core.MeleeHitRatingPerHitChance,
			stats.SpellHit:    ownerStats[stats.SpellHit],
			stats.MeleeCrit:   ownerStats[stats.SpellCrit],
			stats.SpellCrit:   ownerStats[stats.SpellCrit],
			stats.MeleeHaste:  ownerStats[stats.SpellHaste],
			stats.SpellHaste:  ownerStats[stats.SpellHaste],
		}
	}
}

func (shadowfiend *Shadowfiend) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	priest := shadowfiend.Priest
	restoreMana := priest.MaxMana() * 0.05

	priest.AddMana(sim, restoreMana, shadowfiend.ManaMetric, false)
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
	shadowfiend.AutoAttacks.CancelAutoSwing(sim)
}

func (shadowfiend *Shadowfiend) GetPet() *core.Pet {
	return &shadowfiend.Pet
}
