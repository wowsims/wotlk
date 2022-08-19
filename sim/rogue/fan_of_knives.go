package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const FanOfKnivesSpellID int32 = 51723

func (rogue *Rogue) makeFanOfKnivesWeaponHitEffect(isMH bool) core.SpellEffect {
	var procMask core.ProcMask
	var baseDamageConfig core.BaseDamageConfig
	if isMH {
		weaponMultiplier := core.TernaryFloat64(rogue.Equip[proto.ItemSlot_ItemSlotMainHand].WeaponType == proto.WeaponType_WeaponTypeDagger, 1.05, 0.7)
		procMask = core.ProcMaskMeleeMHSpecial
		baseDamageConfig = core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 0, 1, weaponMultiplier, false)

	} else {
		weaponMultiplier := core.TernaryFloat64(rogue.Equip[proto.ItemSlot_ItemSlotOffHand].WeaponType == proto.WeaponType_WeaponTypeDagger, 1.05, 0.7)
		weaponMultiplier += 0.1 * float64(rogue.Talents.DualWieldSpecialization)
		procMask = core.ProcMaskMeleeOHSpecial
		baseDamageConfig = core.BaseDamageConfigMeleeWeapon(core.OffHand, false, 0, 1, weaponMultiplier, false)
	}
	return core.SpellEffect{
		ProcMask: procMask,
		DamageMultiplier: 1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfFanOfKnives), 0.2, 0.0),
		ThreatMultiplier: 1,
		BaseDamage:       baseDamageConfig,
		OutcomeApplier:   rogue.OutcomeFuncMeleeSpecialHitAndCrit(rogue.MeleeCritMultiplier(isMH, false)),
	}

}

func (rogue *Rogue) registerFanOfKnives() {
	energyCost := 50.0
	mhHit := rogue.makeFanOfKnivesWeaponHitEffect(true)
	ohHit := rogue.makeFanOfKnivesWeaponHitEffect(false)
	applyEffects := func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
		core.ApplyEffectFuncAOEDamageCappedWithDeferredEffects(rogue.Env, ohHit)(sim, unit, spell)
		core.ApplyEffectFuncAOEDamageCappedWithDeferredEffects(rogue.Env, mhHit)(sim, unit, spell)
	}
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
		ApplyEffects: applyEffects,
	})
}
