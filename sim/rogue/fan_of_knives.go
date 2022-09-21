package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const FanOfKnivesSpellID int32 = 51723

func (rogue *Rogue) makeFanOfKnivesWeaponHitSpell(isMH bool) (*core.Spell, core.SpellEffect) {
	var procMask core.ProcMask
	var baseDamageConfig core.BaseDamageConfig
	var weaponMultiplier float64
	if isMH {
		weaponMultiplier = core.TernaryFloat64(rogue.Equip[proto.ItemSlot_ItemSlotMainHand].WeaponType == proto.WeaponType_WeaponTypeDagger, 1.05, 0.7)
		procMask = core.ProcMaskMeleeMHSpecial
		baseDamageConfig = core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 0, false)
	} else {
		weaponMultiplier = core.TernaryFloat64(rogue.Equip[proto.ItemSlot_ItemSlotOffHand].WeaponType == proto.WeaponType_WeaponTypeDagger, 1.05, 0.7)
		weaponMultiplier += 0.1 * float64(rogue.Talents.DualWieldSpecialization)
		procMask = core.ProcMaskMeleeOHSpecial
		baseDamageConfig = core.BaseDamageConfigMeleeWeapon(core.OffHand, false, 0, false)
	}

	baseDamageConfig = core.WrapBaseDamageConfig(baseDamageConfig, func(oldCalc core.BaseDamageCalculator) core.BaseDamageCalculator {
		return func(sim *core.Simulation, effect *core.SpellEffect, spell *core.Spell) float64 {
			return oldCalc(sim, effect, spell) * weaponMultiplier
		}
	})

	effect := core.SpellEffect{
		BaseDamage:     baseDamageConfig,
		OutcomeApplier: rogue.OutcomeFuncMeleeSpecialHitAndCrit(rogue.MeleeCritMultiplier(isMH, false)),
	}

	spell := rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: FanOfKnivesSpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfFanOfKnives), 0.2, 0.0),
		ThreatMultiplier: 1,
	})

	return spell, effect
}

func (rogue *Rogue) registerFanOfKnives() {
	energyCost := 50.0
	mhSpell, mhHit := rogue.makeFanOfKnivesWeaponHitSpell(true)
	ohSpell, ohHit := rogue.makeFanOfKnivesWeaponHitSpell(false)
	rogue.FanOfKnives = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: FanOfKnivesSpellID},
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ResourceType: stats.Energy,
		BaseCost:     energyCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: energyCost,
				GCD:  time.Second,
			},
			ModifyCast:  rogue.CastModifier,
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfFanOfKnives), 0.2, 0.0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			core.ApplyEffectFuncAOEDamageCappedWithDeferredEffects(rogue.Env, ohHit)(sim, unit, ohSpell)
			core.ApplyEffectFuncAOEDamageCappedWithDeferredEffects(rogue.Env, mhHit)(sim, unit, mhSpell)
		},
	})
}
