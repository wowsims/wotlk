package warrior

import (
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

func (warrior *Warrior) registerHeroicStrikeSpell() *core.Spell {
	damage := map[int32]float64{
		25: 44,
		40: 80,
		50: 111,
		60: 138,
	}[warrior.Level]

	spellID := map[int32]int32{
		25: 1608,
		40: 11565,
		50: 11566,
		60: 11567,
	}[warrior.Level]

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

		RageCost: core.RageCostOptions{
			Cost:   15 - float64(warrior.Talents.ImprovedHeroicStrike),
			Refund: 0.8,
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,
		FlatThreatBonus:  259,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := damage +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
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
	flatDamageBonus := map[int32]float64{
		25: 5,
		40: 18,
		50: 32,
		60: 50,
	}[warrior.Level]

	spellID := map[int32]int32{
		25: 845,
		40: 11608,
		50: 11609,
		60: 20569,
	}[warrior.Level]

	targets := int32(2)
	numHits := min(targets, warrior.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RageCost: core.RageCostOptions{
			Cost: 20,
		},

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

func (warrior *Warrior) RegisterHSOrCleave(useCleave bool) {

	if useCleave {
		warrior.Cleave = warrior.registerCleaveSpell()
		cleaveQueueSpell := warrior.makeQueueSpellsAndAura(warrior.Cleave)
		warrior.hsOrCleaveQueueSpell = cleaveQueueSpell
	} else if !warrior.HasRune(proto.WarriorRune_RuneQuickStrike) {
		warrior.HeroicStrike = warrior.registerHeroicStrikeSpell()
		hsQueueSpell := warrior.makeQueueSpellsAndAura(warrior.HeroicStrike)
		warrior.hsOrCleaveQueueSpell = hsQueueSpell
	} else {
		warrior.hsOrCleaveQueueSpell = nil
	}
}
