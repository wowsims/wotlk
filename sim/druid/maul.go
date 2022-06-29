package druid

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (druid *Druid) registerMaulSpell(rageThreshold float64) {
	cost := 15.0 - float64(druid.Talents.Ferocity)
	refundAmount := cost * 0.8

	baseDamage := 176.0
	if druid.Equip[items.ItemSlotRanged].ID == 23198 { // Idol of Brutality
		baseDamage += 50
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

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			FlatThreatBonus:  344,

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, baseDamage, 1, true),
			OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(druid.MeleeCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					druid.AddRage(sim, refundAmount, druid.RageRefundMetrics)
				}
			},
		}),
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
