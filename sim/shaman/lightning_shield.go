package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (shaman *Shaman) registerLightningShieldSpell() *core.Spell {
	actionID := core.ActionID{SpellID: 49281}
	var proc = 0.02 * float64(shaman.Talents.StaticShock)
	lsGlyph := 0.0
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLightningShield) {
		lsGlyph = 1.0
	}
	t7bonus := 0.0
	if shaman.HasSetBonus(ItemSetEarthshatterBattlegear, 2) {
		t7bonus = 1.0
	}
	//maxStackCount :=

	procSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1.0 * (0.05 * float64(shaman.Talents.ImprovedShields)) * (0.2 * lsGlyph) * (0.1 * t7bonus),
			ThreatMultiplier: 1, //fix when spirit weapons is fixed
			BaseDamage:       core.BaseDamageConfigMagic(380, 380, 0.267),
			OutcomeApplier:   shaman.OutcomeFuncMagicHitAndCrit(shaman.DefaultSpellCritMultiplier()),
		}),
	})

	switch shaman.Equip[items.ItemSlotHands].ID { //s1 and s2 enh pvp gloves, probably unnessecary but its fun
	case 26000:
		fallthrough
	case 32005:
		procSpell.DamageMultiplier *= 1.08
	}

	lightningShieldAura := shaman.RegisterAura(core.Aura{
		Label:     "Lightning Shield",
		ActionID:  actionID,
		Duration:  time.Minute * 10,
		MaxStacks: 9,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) || !spellEffect.Landed() {
				return
			}
			if sim.RandomFloat("Static Shock") > proc {
				return
			}
			procSpell.Cast(sim, spellEffect.Target)
			aura.RemoveStack(sim)
		},
	})

	shaman.LightningShield = shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			lightningShieldAura.Activate(sim)
			lightningShieldAura.SetStacks(sim, 3+(2*shaman.Talents.StaticShock))
		},
	})

	return (shaman.LightningShield)
}
