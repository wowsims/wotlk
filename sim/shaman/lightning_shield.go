package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (shaman *Shaman) registerLightningShieldSpell() *core.Spell {
	actionID := core.ActionID{SpellID: 49281}

	t9Bonus := false
	if shaman.HasSetBonus(ItemSetThrallsBattlegear, 2) || shaman.HasSetBonus(ItemSetNobundosBattlegear, 2) {
		t9Bonus = true
	}
	var proc = 0.02*float64(shaman.Talents.StaticShock) + core.TernaryFloat64(t9Bonus, 0.03, 0)

	dmgMultBonus := 1.0

	switch shaman.Equip[items.ItemSlotHands].ID { //s1 and s2 enh pvp gloves, probably unnessecary but its fun
	case 26000:
		fallthrough
	case 32005:
		dmgMultBonus = 1.08
	}

	procSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskEmpty,
			DamageMultiplier: 1 * (1 + 0.05*float64(shaman.Talents.ImprovedShields) +
				core.TernaryFloat64(shaman.HasSetBonus(ItemSetEarthshatterBattlegear, 2), 0.1, 0)) *
				core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLightningShield), 1.2, 1) * dmgMultBonus, //possibly additive?

			ThreatMultiplier: 1, //fix when spirit weapons is fixed
			BaseDamage:       core.BaseDamageConfigMagic(380, 380, 0.267),
			OutcomeApplier:   shaman.OutcomeFuncMagicHitAndCrit(shaman.DefaultSpellCritMultiplier()),
		}),
	})

	shaman.LightningShieldAura = shaman.RegisterAura(core.Aura{
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
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shaman.LightningShieldAura.Activate(sim)
			shaman.LightningShieldAura.SetStacks(sim, 3+(2*shaman.Talents.StaticShock))
		},
	})

	return (shaman.LightningShield)
}
