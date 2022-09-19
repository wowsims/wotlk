package druid

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerMaulSpell(rageThreshold float64) {
	cost := 15.0 - float64(druid.Talents.Ferocity)
	refundAmount := cost * 0.8

	baseDamage := 578.0
	if druid.Equip[items.ItemSlotRanged].ID == 23198 { // Idol of Brutality
		baseDamage += 50
	} else if druid.Equip[items.ItemSlotRanged].ID == 38365 { // Idol of Perspicacious Attacks
		baseDamage += 120
	}

	baseEffect := core.SpellEffect{
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,

		BaseDamage: core.WrapBaseDamageConfig(
			core.BaseDamageConfigMeleeWeapon(core.MainHand, false, baseDamage, true),
			func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
				return func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
					normalDamage := oldCalculator(sim, spellEffect, spell)
					modifier := 1.0
					if druid.CurrentTarget.HasActiveAuraWithTag(core.BleedDamageAuraTag) {
						modifier += .3
					}
					if druid.RipDot.IsActive() || druid.RakeDot.IsActive() || druid.LacerateDot.IsActive() {
						modifier *= 1.0 + (0.04 * float64(druid.Talents.RendAndTear))
					}

					return normalDamage * modifier
				}
			}),
		OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(druid.MeleeCritMultiplier()),

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				druid.AddRage(sim, refundAmount, druid.RageRefundMetrics)
			}
		},
	}

	maxTargets := core.TernaryInt(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMaul), 2, 1)
	numHits := core.MinInt(maxTargets, int(druid.Env.GetNumTargets()))
	effects := make([]core.SpellEffect, 0, numHits)
	for i := 0; i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = druid.Env.GetTargetUnit(int32(i))
	}

	druid.Maul = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26996},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
			},
			ModifyCast: druid.ApplyClearcasting,
		},

		DamageMultiplier: 1 + 0.1*float64(druid.Talents.SavageFury),
		ThreatMultiplier: 1,
		FlatThreatBonus:  344,

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
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
