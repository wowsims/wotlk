package druid

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (druid *Druid) registerMaulSpell(rageThreshold float64) {
	flatBaseDamage := 578.0
	if druid.Equip[core.ItemSlotRanged].ID == 23198 { // Idol of Brutality
		flatBaseDamage += 50
	} else if druid.Equip[core.ItemSlotRanged].ID == 38365 { // Idol of Perspicacious Attacks
		flatBaseDamage += 120
	}

	numHits := core.TernaryInt32(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMaul) && druid.Env.GetNumTargets() > 1, 2, 1)
	bleedCategory := druid.CurrentTarget.GetExclusiveEffectCategory(core.BleedEffectCategory)

	druid.Maul = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26996},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

		RageCost: core.RageCostOptions{
			Cost:   15 - float64(druid.Talents.Ferocity),
			Refund: 0.8,
		},

		DamageMultiplier: 1 + 0.1*float64(druid.Talents.SavageFury),
		CritMultiplier:   druid.MeleeCritMultiplier(Bear),
		ThreatMultiplier: 1,
		FlatThreatBonus:  424,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Need to specially deactivate CC here in case maul is cast simultaneously with another spell.
			if druid.ClearcastingAura != nil {
				druid.ClearcastingAura.Deactivate(sim)
			}

			modifier := 1.0
			if bleedCategory.AnyActive() {
				modifier += .3
			}
			if druid.AssumeBleedActive || druid.Rip.Dot(target).IsActive() || druid.Rake.Dot(target).IsActive() || druid.Lacerate.Dot(target).IsActive() {
				modifier *= 1.0 + (0.04 * float64(druid.Talents.RendAndTear))
			}

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := flatBaseDamage +
					spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
				baseDamage *= modifier

				result := spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

				if !result.Landed() {
					spell.IssueRefund(sim)
				}

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})

	druid.MaulQueueAura = druid.RegisterAura(core.Aura{
		Label:    "Maul Queue Aura",
		ActionID: druid.Maul.ActionID,
		Duration: core.NeverExpires,
	})

	druid.MaulRageThreshold = core.MaxFloat(druid.Maul.DefaultCast.Cost, rageThreshold)
}

func (druid *Druid) QueueMaul(sim *core.Simulation) {
	if druid.CurrentRage() < druid.Maul.DefaultCast.Cost {
		panic("Not enough rage for HS")
	}
	if druid.MaulQueueAura.IsActive() {
		return
	}
	druid.MaulQueueAura.Activate(sim)
}

// Returns true if the regular melee swing should be used, false otherwise.
func (druid *Druid) MaulReplaceMH(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if !druid.MaulQueueAura.IsActive() {
		return nil
	}

	if druid.CurrentRage() < druid.Maul.DefaultCast.Cost {
		druid.MaulQueueAura.Deactivate(sim)
		return nil
	} else if druid.CurrentRage() < druid.MaulRageThreshold {
		if mhSwingSpell == druid.AutoAttacks.MHAuto {
			druid.MaulQueueAura.Deactivate(sim)
			return nil
		}
	}

	druid.MaulQueueAura.Deactivate(sim)
	return druid.Maul
}

func (druid *Druid) ShouldQueueMaul(sim *core.Simulation) bool {
	return druid.CurrentRage() >= druid.MaulRageThreshold
}
