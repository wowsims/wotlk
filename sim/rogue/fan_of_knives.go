package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var FanOfKnivesActionID = core.ActionID{SpellID: 51723}

func (rogue *Rogue) makeFanOfKnifesWeaponHitEffect(isMH bool) core.SpellEffect {
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
		ProcMask:         procMask,
		DamageMultiplier: 1 + core.TernaryFloat64(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfFanOfKnives), 0.2, 0.0),
		ThreatMultiplier: 1,
		BaseDamage:       baseDamageConfig,
		OutcomeApplier:   rogue.OutcomeFuncMeleeSpecialHitAndCrit(rogue.MeleeCritMultiplier(isMH, false)),
	}

}

func (rogue *Rogue) registerFanOfKnives() {
	mhWeaponHitSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID:     FanOfKnivesActionID.WithTag(1),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics | SpellFlagRogueAbility,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(rogue.makeFanOfKnifesWeaponHitEffect(true)),
	})
	ohWeaponHitSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID:     FanOfKnivesActionID.WithTag(2),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics | SpellFlagRogueAbility,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(rogue.makeFanOfKnifesWeaponHitEffect(false)),
	})
	energyCost := 50.0

	rogue.FanOfKnifes = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     FanOfKnivesActionID,
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagNoMetrics | SpellFlagBuilder,
		ResourceType: stats.Energy,
		BaseCost:     energyCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: energyCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: core.ApplyEffectFuncAOEDamage(rogue.Env, core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,
			OutcomeApplier:   rogue.OutcomeFuncAlwaysHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				mhWeaponHitSpell.Cast(sim, spellEffect.Target)
				ohWeaponHitSpell.Cast(sim, spellEffect.Target)
			},
		}),
	})
}
