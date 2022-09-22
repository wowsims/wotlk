package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) registerHemorrhageSpell() {
	actionID := core.ActionID{SpellID: 26864}
	target := rogue.CurrentTarget
	bonusDamage := 75.0
	if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfHemorrhage) {
		bonusDamage *= 1.4
	}
	hemoAura := target.GetOrRegisterAura(core.Aura{
		Label:     "Hemorrhage",
		ActionID:  actionID,
		Duration:  time.Second * 15,
		MaxStacks: 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			target.PseudoStats.BonusPhysicalDamageTaken += bonusDamage
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			target.PseudoStats.BonusPhysicalDamageTaken -= bonusDamage
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell.SpellSchool != core.SpellSchoolPhysical {
				return
			}
			if !spellEffect.Landed() || spellEffect.Damage == 0 {
				return
			}

			aura.RemoveStack(sim)
		},
	})

	baseCost := 35.0
	refundAmount := baseCost * 0.8
	daggerMH := rogue.Equip[proto.ItemSlot_ItemSlotMainHand].WeaponType == proto.WeaponType_WeaponTypeDagger
	weaponDamageBonus := core.TernaryFloat64(daggerMH, 1.6, 1.1) + float64(rogue.Talents.SinisterCalling)*0.02
	rogue.Hemorrhage = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | SpellFlagBuilder,
		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
			ModifyCast:  rogue.CastModifier,
		},

		BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(ItemSetVanCleefs, 4), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplier: (1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.HasSetBonus(ItemSetSlayers, 4), 0.06, 0)) * weaponDamageBonus,
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfigMeleeWeapon(
				core.MainHand, true, 0, true),
			OutcomeApplier: rogue.OutcomeFuncMeleeSpecialHitAndCrit(rogue.MeleeCritMultiplier(true, true)),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
					hemoAura.Activate(sim)
					hemoAura.SetStacks(sim, 10)
				} else {
					rogue.AddEnergy(sim, refundAmount, rogue.EnergyRefundMetrics)
				}
			},
		}),
	})
}
