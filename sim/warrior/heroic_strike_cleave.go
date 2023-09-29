package warrior

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warrior *Warrior) registerHeroicStrikeSpell() *core.Spell {
	hasGlyph := warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfHeroicStrike)
	var rageMetrics *core.ResourceMetrics
	if hasGlyph {
		rageMetrics = warrior.NewRageMetrics(core.ActionID{ItemID: 43418})
	}

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47450},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete | SpellFlagBloodsurge,

		RageCost: core.RageCostOptions{
			Cost:   15 - float64(warrior.Talents.ImprovedHeroicStrike) - float64(warrior.Talents.FocusedRage),
			Refund: 0.8,
		},

		BonusCritRating:  (5*float64(warrior.Talents.Incite) + core.TernaryFloat64(warrior.HasSetBonus(ItemSetWrynnsBattlegear, 4), 5, 0)) * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,
		FlatThreatBonus:  259,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 495 +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.DidCrit() && hasGlyph {
				warrior.AddRage(sim, 10, rageMetrics)
			} else if !result.Landed() {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
		},
	})
}

func (warrior *Warrior) registerCleaveSpell() *core.Spell {
	flatDamageBonus := 222 * (1 + 0.4*float64(warrior.Talents.ImprovedCleave))

	targets := core.TernaryInt32(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfCleaving), 3, 2)
	numHits := min(targets, warrior.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47520},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RageCost: core.RageCostOptions{
			Cost: 20 - float64(warrior.Talents.FocusedRage),
		},

		BonusCritRating:  float64(warrior.Talents.Incite) * 5 * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,
		FlatThreatBonus:  225,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := flatDamageBonus +
					spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
		},
	})
}

func (warrior *Warrior) makeQueueSpellsAndAura(srcSpell *core.Spell) *core.Spell {
	queueAura := warrior.RegisterAura(core.Aura{
		Label:    "HS/Cleave Queue Aura-" + srcSpell.ActionID.String(),
		ActionID: srcSpell.ActionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
			warrior.PseudoStats.DisableDWMissPenalty = true
			warrior.curQueueAura = aura
			warrior.curQueuedAutoSpell = srcSpell
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DisableDWMissPenalty = false
			warrior.curQueueAura = nil
			warrior.curQueuedAutoSpell = nil
		},
	})

	queueSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:    srcSpell.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.curQueueAura != queueAura &&
				warrior.CurrentRage() >= srcSpell.DefaultCast.Cost &&
				sim.CurrentTime >= warrior.Hardcast.Expires
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			queueAura.Activate(sim)
		},
	})

	return queueSpell
}

func (warrior *Warrior) QueueHSOrCleave(sim *core.Simulation) {
	if warrior.hsOrCleaveQueueSpell.CanCast(sim, warrior.CurrentTarget) {
		warrior.hsOrCleaveQueueSpell.Cast(sim, warrior.CurrentTarget)
	}
}

// Returns true if the regular melee swing should be used, false otherwise.
func (warrior *Warrior) TryHSOrCleave(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if !warrior.curQueueAura.IsActive() {
		return mhSwingSpell
	}

	if warrior.CurrentRage() < warrior.HSRageThreshold {
		warrior.curQueueAura.Deactivate(sim)
		return mhSwingSpell
	}

	if !warrior.curQueuedAutoSpell.CanCast(sim, warrior.CurrentTarget) {
		warrior.curQueueAura.Deactivate(sim)
		return mhSwingSpell
	}

	return warrior.curQueuedAutoSpell
}

func (warrior *Warrior) ShouldQueueHSOrCleave(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.HSRageThreshold && sim.CurrentTime >= warrior.Hardcast.Expires
}

func (warrior *Warrior) RegisterHSOrCleave(useCleave bool, rageThreshold float64) {
	warrior.HeroicStrike = warrior.registerCleaveSpell()
	hsQueueSpell := warrior.makeQueueSpellsAndAura(warrior.HeroicStrike)
	warrior.Cleave = warrior.registerHeroicStrikeSpell()
	cleaveQueueSpell := warrior.makeQueueSpellsAndAura(warrior.Cleave)

	var autoSpell *core.Spell
	if useCleave {
		autoSpell = warrior.HeroicStrike
		warrior.hsOrCleaveQueueSpell = hsQueueSpell
	} else {
		autoSpell = warrior.Cleave
		warrior.hsOrCleaveQueueSpell = cleaveQueueSpell
	}

	warrior.HSRageThreshold = max(autoSpell.DefaultCast.Cost, rageThreshold)
	if warrior.IsUsingAPL {
		warrior.HSRageThreshold = 0
	}
}
