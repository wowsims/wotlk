package druid

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerMaulSpell(rageThreshold float64) {
	cost := 15.0 - float64(druid.Talents.Ferocity)
	refundAmount := cost * 0.8

	flatBaseDamage := 578.0
	if druid.Equip[core.ItemSlotRanged].ID == 23198 { // Idol of Brutality
		flatBaseDamage += 50
	} else if druid.Equip[core.ItemSlotRanged].ID == 38365 { // Idol of Perspicacious Attacks
		flatBaseDamage += 120
	}

	numHits := core.TernaryInt32(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMaul) && druid.Env.GetNumTargets() > 1, 2, 1)

	druid.Maul = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 26996},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,
		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
			},
		},

		DamageMultiplier: 1 + 0.1*float64(druid.Talents.SavageFury),
		CritMultiplier:   druid.MeleeCritMultiplier(),
		ThreatMultiplier: 1,
		FlatThreatBonus:  344,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			modifier := 1.0
			if druid.CurrentTarget.HasActiveAuraWithTag(core.BleedDamageAuraTag) {
				modifier += .3
			}
			if druid.AssumeBleedActive || druid.RipDot.IsActive() || druid.RakeDot.IsActive() || druid.LacerateDot.IsActive() {
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
					druid.AddRage(sim, refundAmount, druid.RageRefundMetrics)
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

func (druid *Druid) DequeueMaul(sim *core.Simulation) {
	druid.MaulQueueAura.Deactivate(sim)
}

// Returns true if the regular melee swing should be used, false otherwise.
func (druid *Druid) TryMaul(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if !druid.MaulQueueAura.IsActive() {
		return nil
	}

	if druid.CurrentRage() < druid.Maul.DefaultCast.Cost {
		druid.DequeueMaul(sim)
		return nil
	} else if druid.CurrentRage() < druid.MaulRageThreshold {
		if mhSwingSpell == druid.AutoAttacks.MHAuto {
			druid.DequeueMaul(sim)
			return nil
		}
	}

	druid.DequeueMaul(sim)
	return druid.Maul
}

func (druid *Druid) ShouldQueueMaul(sim *core.Simulation) bool {
	return druid.CurrentRage() >= druid.MaulRageThreshold
}
