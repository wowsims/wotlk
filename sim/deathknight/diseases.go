package deathknight

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) countActiveDiseases() int {
	count := 0
	if deathKnight.FrostFeverDisease.IsActive() {
		count++
	}
	if deathKnight.BloodPlagueDisease.IsActive() {
		count++
	}
	if deathKnight.EbonPlagueAura.IsActive() {
		count++
	}
	return count
}

func (deathKnight *DeathKnight) registerDiseaseDots() {
	deathKnight.registerFrostFever()
	deathKnight.registerBloodPlague()
}

func (deathKnight *DeathKnight) registerFrostFever() {
	actionID := core.ActionID{SpellID: 55095}
	target := deathKnight.CurrentTarget

	deathKnight.FrostFeverDisease = core.NewDot(core.Dot{
		Aura: target.RegisterAura(core.Aura{
			Label:    "FrostFever-" + strconv.Itoa(int(deathKnight.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 5 + int(deathKnight.Talents.Epidemic),
		TickLength:    time.Second * 3,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if deathKnight.Talents.WanderingPlague == 0 {
					return
				}

				critRating := spell.Unit.GetStats()[stats.MeleeCrit] + spellEffect.BonusCritRating + spellEffect.Target.PseudoStats.BonusCritRatingTaken
				critRating += spell.Unit.PseudoStats.BonusMeleeCritRating
				critChance := critRating / (core.CritRatingPerCritChance * 100)
				if sim.RandomFloat("Wandering Plague Roll") < critChance {
					deathKnight.LastDiseaseDamage = spellEffect.Damage * float64(deathKnight.Talents.WanderingPlague) * 0.33
					deathKnight.WanderingPlague.Cast(sim, spellEffect.Target)
				}
			},
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return ((127.0 + 80.0*0.32) + hitEffect.MeleeAttackPower(spell.Unit)*0.055) * (1.0 +
						core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 0.02*float64(deathKnight.Talents.RageOfRivendare), 0.0) +
						core.TernaryFloat64(deathKnight.DiseasesAreActive(), 0.05*float64(deathKnight.Talents.TundraStalker), 0.0))
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})

	deathKnight.FrostFeverSpell = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFrost,
		ApplyEffects: core.ApplyEffectFuncDot(deathKnight.FrostFeverDisease),
	})

	deathKnight.FrostFeverDisease.Spell = deathKnight.FrostFeverSpell
}

func (deathKnight *DeathKnight) registerBloodPlague() {
	actionID := core.ActionID{SpellID: 55078}
	target := deathKnight.CurrentTarget

	deathKnight.BloodPlagueDisease = core.NewDot(core.Dot{
		Aura: target.RegisterAura(core.Aura{
			Label:    "BloodPlague-" + strconv.Itoa(int(deathKnight.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 5 + int(deathKnight.Talents.Epidemic),
		TickLength:    time.Second * 3,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if deathKnight.Talents.WanderingPlague == 0 {
					return
				}

				critRating := spell.Unit.GetStats()[stats.MeleeCrit] + spellEffect.BonusCritRating + spellEffect.Target.PseudoStats.BonusCritRatingTaken
				critRating += spell.Unit.PseudoStats.BonusMeleeCritRating
				critChance := critRating / (core.CritRatingPerCritChance * 100)
				if sim.RandomFloat("Wandering Plague Roll") < critChance {
					deathKnight.LastDiseaseDamage = spellEffect.Damage * float64(deathKnight.Talents.WanderingPlague) * 0.33
					deathKnight.WanderingPlague.Cast(sim, spellEffect.Target)
				}
			},
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return ((127.0 + 80.0*0.32) + hitEffect.MeleeAttackPower(spell.Unit)*0.055) * (1.0 +
						core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 0.02*float64(deathKnight.Talents.RageOfRivendare), 0.0) +
						core.TernaryFloat64(deathKnight.DiseasesAreActive(), 0.05*float64(deathKnight.Talents.TundraStalker), 0.0))
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})

	deathKnight.BloodPlagueSpell = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ApplyEffects: core.ApplyEffectFuncDot(deathKnight.BloodPlagueDisease),
	})

	deathKnight.BloodPlagueDisease.Spell = deathKnight.BloodPlagueSpell
}
