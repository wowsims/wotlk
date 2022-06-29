package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warlock *Warlock) registerSeedSpell() {
	numTargets := int(warlock.Env.GetNumTargets())

	warlock.Seeds = make([]*core.Spell, numTargets)
	warlock.SeedDots = make([]*core.Dot, numTargets)

	numHit := float64(numTargets - 1)
	// For this simulation we always assume the seed target didn't die to trigger the seed because we don't simulate health.
	// This effectively lowers the seed AOE cap using the function:
	cap := 13580.0 * numHit / (numHit + 1)

	for i := 0; i < numTargets; i++ {
		warlock.makeSeed(i, cap)
	}
}

func (warlock *Warlock) makeSeed(targetIdx int, cap float64) {
	baseCost := 882.0

	flatBonus := 0.0
	if ItemSetOblivionRaiment.CharacterHasSetBonus(&warlock.Character, 4) {
		flatBonus += 180
	}
	baseSeedExplosionEffect := core.SpellEffect{
		ProcMask:         core.ProcMaskSpellDamage,
		DamageMultiplier: 1 * (1 + 0.02*float64(warlock.Talents.ShadowMastery)) * (1 + 0.01*float64(warlock.Talents.Contagion)),
		ThreatMultiplier: 1 - 0.05*float64(warlock.Talents.ImprovedDrainSoul),
		BaseDamage:       core.BaseDamageConfigMagic(1110+flatBonus, 1290+flatBonus, 0.143),
		OutcomeApplier:   warlock.OutcomeFuncMagicHitAndCrit(1.5),
	}

	// Use a custom aoe effect list that does not include the seeded target.
	baseEffects := make([]core.SpellEffect, warlock.Env.GetNumTargets()-1)
	skipped := false
	for i := range baseEffects {
		baseEffects[i] = baseSeedExplosionEffect
		expTarget := i
		if i == targetIdx {
			skipped = true
		}
		if skipped {
			expTarget++
		}
		baseEffects[i].Target = warlock.Env.GetTargetUnit(int32(expTarget))
	}
	seedActionID := core.ActionID{SpellID: 27243}

	explosionId := seedActionID
	explosionId.Tag = 1

	seedExplosion := warlock.RegisterSpell(core.SpellConfig{
		ActionID:     explosionId,
		SpellSchool:  core.SpellSchoolShadow,
		Cast:         core.CastConfig{},
		ApplyEffects: core.ApplyEffectFuncMultipleDamageCapped(baseEffects, cap),
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
		ActionID:     seedActionID,
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
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
		if seedDmgTracker > 1044 {
			warlock.SeedDots[targetIdx].Deactivate(sim)
			seedExplosion.Cast(sim, target)
			seedDmgTracker = 0
		}
	}
	warlock.SeedDots[targetIdx] = core.NewDot(core.Dot{
		Spell: warlock.Seeds[targetIdx],
		Aura: target.RegisterAura(core.Aura{
			Label:    "Seed-" + strconv.Itoa(int(warlock.Index)),
			ActionID: seedActionID,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}
				if spell.ActionID.SpellID == seedActionID.SpellID {
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
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * (1 + 0.02*float64(warlock.Talents.ShadowMastery)) * (1 + 0.01*float64(warlock.Talents.Contagion)),
			ThreatMultiplier: 1 - 0.05*float64(warlock.Talents.ImprovedDrainSoul),
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(174, 0.25),
			OutcomeApplier:   warlock.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})
}
