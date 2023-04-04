package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (shaman *Shaman) registerAncestralHealingSpell() {
	shaman.AncestralAwakening = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 52752},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskSpellHealing,
		Flags:            core.SpellFlagHelpful,
		DamageMultiplier: 1 * (1 + .02*float64(shaman.Talents.Purification)),
		CritMultiplier:   1,
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, target, shaman.ancestralHealingAmount, spell.OutcomeHealing)
		},
	})
}

func (shaman *Shaman) registerLesserHealingWaveSpell() {
	spellCoeff := 0.807
	bonusCoeff := 0.02 * float64(shaman.Talents.TidalWaves)
	impShieldChance := 0.2 * float64(shaman.Talents.ImprovedWaterShield)
	impShieldManaGain := 428.0 * (1 + 0.05*float64(shaman.Talents.ImprovedShields))

	hasGlyph := shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLesserHealingWave)

	bonusHeal := 0 +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 42598, 338, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 42597, 267, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 42596, 236, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 42595, 204, 0)

	shaman.LesserHealingWave = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49276},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.15,
			Multiplier: 1 *
				(1 - .01*float64(shaman.Talents.TidalFocus)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusCritRating: float64(shaman.Talents.TidalMastery) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(shaman.Talents.Purification)),
		CritMultiplier:   shaman.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healPower := spell.HealingPower(target)
			baseHealing := sim.Roll(1624, 1852) + spellCoeff*healPower + bonusCoeff*healPower + bonusHeal
			if hasGlyph {
				if shaman.EarthShield.Hot(target).IsActive() {
					baseHealing *= 1.2
				}
			}
			result := spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			if result.Outcome.Matches(core.OutcomeCrit) {
				if impShieldChance > 0 {
					if sim.RandomFloat("imp water shield") > impShieldChance {
						shaman.AddMana(sim, impShieldManaGain, shaman.waterShieldManaMetrics)
					}
				}
				if shaman.Talents.AncestralAwakening > 0 {
					shaman.ancestralHealingAmount = result.Damage * 0.3

					// TODO: this should actually target the lowest health target in the raid.
					//  does it matter in a sim? We currently only simulate tanks taking damage (multiple tanks could be handled here though.)
					shaman.AncestralAwakening.Cast(sim, target)
				}
			}

			if shaman.tidalWaveProc.IsActive() {
				shaman.tidalWaveProc.RemoveStack(sim)
			}
		},
	})
}

func (shaman *Shaman) registerRiptideSpell() {
	spellCoeff := 0.402
	hotCoeff := 0.188
	impShieldChance := []float64{0, 0.33, 0.66, 1.0}[shaman.Talents.ImprovedWaterShield]
	impShieldManaGain := 428.0 * (1 + 0.05*float64(shaman.Talents.ImprovedShields))

	shaman.Riptide = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 61301},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.18,
			Multiplier: 1 *
				(1 - .01*float64(shaman.Talents.TidalFocus)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		BonusCritRating: float64(shaman.Talents.TidalMastery) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(shaman.Talents.Purification)),
		CritMultiplier:   shaman.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label: "Riptide",
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 334 + hotCoeff*dot.Spell.HealingPower(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healPower := spell.HealingPower(target)
			baseHealing := sim.Roll(1604, 1736) + spellCoeff*healPower
			result := spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
			spell.Hot(target).Apply(sim)

			if result.Outcome.Matches(core.OutcomeCrit) {
				if impShieldChance > 0 {
					if impShieldChance > 0.9999 || sim.RandomFloat("imp water shield") > impShieldChance {
						shaman.AddMana(sim, impShieldManaGain, shaman.waterShieldManaMetrics)
					}
				}
				if shaman.Talents.AncestralAwakening > 0 {
					shaman.ancestralHealingAmount = result.Damage * 0.3
					// TODO: this should actually target the lowest health target in the raid.
					//  does it matter in a sim? We currently only simulate tanks taking damage (multiple tanks could be handled here though.)
					shaman.AncestralAwakening.Cast(sim, target)
				}
			}
			if shaman.Talents.TidalWaves > 0 {
				shaman.tidalWaveProc.Activate(sim)
				shaman.tidalWaveProc.SetStacks(sim, 2)
			}
		},
	})
}

func (shaman *Shaman) registerHealingWaveSpell() {
	// TODO: finish this
	// ActionID:    core.ActionID{SpellID: 49273},

	// -79 mana totem: 39728

	spellCoeff := 0.807
	bonusCoeff := 0.02 * float64(shaman.Talents.TidalWaves)
	impShieldChance := 0.2 * float64(shaman.Talents.ImprovedWaterShield)
	impShieldManaGain := 428.0 * (1 + 0.05*float64(shaman.Talents.ImprovedShields))

	hasGlyph := shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLesserHealingWave)

	bonusHeal := 0 +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 42598, 338, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 42597, 267, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 42596, 236, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 42595, 204, 0)

	shaman.HealingWave = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49273},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.15,
			Multiplier: 1 *
				(1 - .01*float64(shaman.Talents.TidalFocus)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusCritRating: float64(shaman.Talents.TidalMastery) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(shaman.Talents.Purification)),
		CritMultiplier:   shaman.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healPower := spell.HealingPower(target)
			baseHealing := sim.Roll(1624, 1852) + spellCoeff*healPower + bonusCoeff*healPower + bonusHeal
			if hasGlyph {
				if shaman.EarthShield.Hot(target).IsActive() {
					baseHealing *= 1.2
				}
			}
			result := spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			if result.Outcome.Matches(core.OutcomeCrit) {
				if impShieldChance > 0 {
					if sim.RandomFloat("imp water shield") > impShieldChance {
						shaman.AddMana(sim, impShieldManaGain, shaman.waterShieldManaMetrics)
					}
				}
				if shaman.Talents.AncestralAwakening > 0 {
					shaman.ancestralHealingAmount = result.Damage * 0.3

					// TODO: this should actually target the lowest health target in the raid.
					//  does it matter in a sim? We currently only simulate tanks taking damage (multiple tanks could be handled here though.)
					shaman.AncestralAwakening.Cast(sim, target)
				}
			}

			if shaman.tidalWaveProc.IsActive() {
				shaman.tidalWaveProc.RemoveStack(sim)
			}
		},
	})
}

func (shaman *Shaman) registerEarthShieldSpell() {
	actionID := core.ActionID{SpellID: 49284}
	spCoeff := 0.286

	bonusHeal := 0.0
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfEarthShield) {
		bonusHeal = 0.2
	}

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond * 3500,
	}

	shaman.EarthShield = shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskEmpty,
		BonusCritRating:  float64(shaman.Talents.TidalMastery) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + 0.05*float64(shaman.Talents.ImprovedShields) + 0.05*float64(shaman.Talents.ImprovedEarthShield) + bonusHeal,
		ThreatMultiplier: 1,
		Hot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Earth Shield",
				ActionID: core.ActionID{SpellID: 379},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}
					if !icd.IsReady(sim) {
						return
					}
					icd.Use(sim)
					shaman.EarthShield.Hot(result.Target).ManualTick(sim)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				},
			},
			NumberOfTicks: 6 + shaman.Talents.ImprovedEarthShield,
			TickLength:    time.Minute*10 + 1, // tick length longer than expire time.
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 377 + dot.Spell.HealingPower(target)*spCoeff
				dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Hot(target).Apply(sim)
		},
	})
}

func (shaman *Shaman) registerChainHealSpell() {
	spellCoeff := 1.342884
	impShieldChance := 0.1 * float64(shaman.Talents.ImprovedWaterShield)
	impShieldManaGain := 428.0 * (1 + 0.05*float64(shaman.Talents.ImprovedShields))

	hasGlyph := shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainHeal)

	numHits := core.MinInt32(core.TernaryInt32(hasGlyph, 4, 3), int32(len(shaman.Env.Raid.AllUnits)))

	bonusHeal := 0 +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 28523, 87, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 38368, 102, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 45114, 257, 0)

	manaDiscount := 0 +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == 40709, 78, 0)

	shaman.ChainHeal = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 55459},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful,

		ManaCost: core.ManaCostOptions{
			FlatCost: 0.19*shaman.BaseMana - manaDiscount,
			Multiplier: 1 *
				(1 - .01*float64(shaman.Talents.TidalFocus)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},
		BonusCritRating:  float64(shaman.Talents.TidalMastery) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + .02*float64(shaman.Talents.Purification) + 0.1*float64(shaman.Talents.ImprovedChainHeal),
		CritMultiplier:   shaman.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			bounceCoeff := 1.0
			dmgReductionPerBounce := 0.6
			curTarget := target
			// TODO: This bounces to most hurt friendly...
			targets := sim.Environment.Raid.GetFirstNPlayersOrPets(numHits)
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				healPower := spell.HealingPower(target)
				baseHealing := sim.Roll(1055, 1205) + spellCoeff*healPower + bonusHeal
				baseHealing *= bounceCoeff

				riptide := shaman.Riptide.Hot(curTarget)
				if riptide.IsActive() {
					riptide.Deactivate(sim)
					baseHealing *= 1.25
				}

				result := spell.CalcAndDealHealing(sim, curTarget, baseHealing, spell.OutcomeHealingCrit)
				if result.Outcome.Matches(core.OutcomeCrit) {
					if impShieldChance > 0 {
						if sim.RandomFloat("imp water shield") > impShieldChance {
							shaman.AddMana(sim, impShieldManaGain, shaman.waterShieldManaMetrics)
						}
					}
				}
				if shaman.Talents.TidalWaves > 0 {
					shaman.tidalWaveProc.Activate(sim)
					shaman.tidalWaveProc.SetStacks(sim, 2)
				}

				bounceCoeff *= dmgReductionPerBounce
				curTarget = targets[hitIndex]
			}
		},
	})
}
