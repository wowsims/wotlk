package deathknight

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
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

func (deathKnight *DeathKnight) diseaseMultiplierBonus(multiplier float64) float64 {
	return 1.0 + float64(deathKnight.countActiveDiseases())*multiplier
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
			DamageMultiplier: core.TernaryFloat64(deathKnight.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfIcyTouch), 1.2, 1.0),
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				deathKnight.doWanderingPlague(sim, spell, spellEffect)
			},
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return ((127.0 + 80.0*0.32) + deathKnight.applyImpurity(hitEffect, spell.Unit)*0.055) *
						deathKnight.rageOfRivendareBonus() *
						deathKnight.tundraStalkerBonus()
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})

	deathKnight.FrostFeverSpell = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFrost,
		Flags:        core.SpellFlagDisease,
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
				deathKnight.doWanderingPlague(sim, spell, spellEffect)
			},
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return ((127.0 + 80.0*0.32) + deathKnight.applyImpurity(hitEffect, spell.Unit)*0.055) *
						deathKnight.rageOfRivendareBonus() *
						deathKnight.tundraStalkerBonus()
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})

	deathKnight.BloodPlagueSpell = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		Flags:        core.SpellFlagDisease,
		ApplyEffects: core.ApplyEffectFuncDot(deathKnight.BloodPlagueDisease),
	})

	deathKnight.BloodPlagueDisease.Spell = deathKnight.BloodPlagueSpell
}

func (deathKnight *DeathKnight) doWanderingPlague(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
	if deathKnight.Talents.WanderingPlague == 0 {
		return
	}

	critRating := spell.Unit.GetStats()[stats.MeleeCrit] + spellEffect.BonusCritRating + spellEffect.Target.PseudoStats.BonusCritRatingTaken
	critRating += spell.Unit.PseudoStats.BonusMeleeCritRating
	critChance := critRating / (core.CritRatingPerCritChance * 100)
	if sim.RandomFloat("Wandering Plague Roll") < critChance {
		deathKnight.LastDiseaseDamage = spellEffect.Damage
		deathKnight.WanderingPlague.Cast(sim, spellEffect.Target)
	}
}
