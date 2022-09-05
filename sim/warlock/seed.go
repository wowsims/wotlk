package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerSeedSpell() {
	numTargets := int(warlock.Env.GetNumTargets())

	warlock.Seeds = make([]*core.Spell, numTargets)
	warlock.SeedDots = make([]*core.Dot, numTargets)

	// For this simulation we always assume the seed target didn't die to trigger the seed because we don't simulate health.
	// This effectively lowers the seed AOE cap using the function:
	for i := 0; i < numTargets; i++ {
		warlock.makeSeed(i, numTargets)
	}
}

func (warlock *Warlock) makeSeed(targetIdx int, numTargets int) {
	actionID := core.ActionID{SpellID: 47836}
	spellSchool := core.SpellSchoolShadow
	baseAdditiveMultiplier := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false)
	baseAdditiveMultiplierDot := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true)

	baseCost := 0.34 * warlock.BaseMana

	baseSeedExplosionEffect := core.SpellEffect{
		ProcMask: core.ProcMaskSpellDamage,

		BonusCritRating: 0 +
			warlock.masterDemonologistShadowCrit() +
			float64(warlock.Talents.ImprovedCorruption)*core.CritRatingPerCritChance,
		DamageMultiplier: baseAdditiveMultiplier,
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		BaseDamage:     core.BaseDamageConfigMagic(1633, 1897, 0.2129),
		OutcomeApplier: warlock.OutcomeFuncMagicHitAndCrit(warlock.DefaultSpellCritMultiplier()),
	}

	// Use a custom aoe effect list that does not include the seeded target.
	baseEffects := make([]core.SpellEffect, warlock.Env.GetNumTargets())
	for i := range baseEffects {
		baseEffects[i] = baseSeedExplosionEffect
		baseEffects[i].Target = warlock.Env.GetTargetUnit(int32(i))
	}

	actionID.Tag = 1

	seedExplosion := warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		Cast:         core.CastConfig{},
		ApplyEffects: core.ApplyEffectFuncMultipleDamageCapped(baseEffects, false),
	})

	effect := core.SpellEffect{
		ProcMask:        core.ProcMaskEmpty,
		OutcomeApplier:  warlock.OutcomeFuncMagicHit(),
		OnSpellHitDealt: applyDotOnLanded(&warlock.SeedDots[targetIdx]),
	}
	if warlock.Rotation.DetonateSeed {
		// Replace dot application with explosion.
		effect.OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			seedExplosion.Cast(sim, spellEffect.Target)
		}
	}

	warlock.Seeds[targetIdx] = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

	target := warlock.Env.GetTargetUnit(int32(targetIdx))

	seedDmgTracker := 0.0
	trySeedPop := func(sim *core.Simulation, dmg float64) {
		seedDmgTracker += dmg
		if seedDmgTracker > 1518 {
			warlock.SeedDots[targetIdx].Deactivate(sim)
			seedExplosion.Cast(sim, target)
			seedDmgTracker = 0
		}
	}
	warlock.SeedDots[targetIdx] = core.NewDot(core.Dot{
		Spell: warlock.Seeds[targetIdx],
		Aura: target.RegisterAura(core.Aura{
			Label:    "Seed-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}
				if spell.ActionID.SpellID == actionID.SpellID {
					return // Seed can't pop seed.
				}
				trySeedPop(sim, spellEffect.Damage)
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				trySeedPop(sim, spellEffect.Damage)
			},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				seedDmgTracker = 0
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				seedDmgTracker = 0
			},
		}),

		NumberOfTicks: 6,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:   core.ProcMaskPeriodicDamage,
			IsPeriodic: true,

			DamageMultiplier: baseAdditiveMultiplierDot,
			ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

			BaseDamage:     core.BaseDamageConfigMagicNoRoll(1518/6, 0.25),
			OutcomeApplier: warlock.OutcomeFuncTick(),
		}),
	})
}
