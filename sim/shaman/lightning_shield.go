package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (shaman *Shaman) registerLightningShieldSpell() {
	if shaman.SelfBuffs.Shield != proto.ShamanShield_LightningShield {
		return
	}

	actionID := core.ActionID{SpellID: 49281}
	procChance := 0.02*float64(shaman.Talents.StaticShock) + core.TernaryFloat64(shaman.HasSetBonus(ItemSetThrallsBattlegear, 2), 0.03, 0)

	dmgMultBonus := 0.0
	switch shaman.Equip[items.ItemSlotHands].ID { //s1 and s2 enh pvp gloves, probably unnessecary but its fun
	case 26000:
		fallthrough
	case 32005:
		dmgMultBonus = 0.08
	}

	procSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49279},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplier: 1*(1+0.05*float64(shaman.Talents.ImprovedShields)+
			core.TernaryFloat64(shaman.HasSetBonus(ItemSetEarthshatterBattlegear, 2), 0.1, 0)) +
			core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLightningShield), 0.2, 0) + dmgMultBonus,
		ThreatMultiplier: 1, //fix when spirit weapons is fixed

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 380 + 0.267*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
		},
	})

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond * 3500,
	}

	shaman.LightningShieldAura = shaman.RegisterAura(core.Aura{
		Label:     "Lightning Shield",
		ActionID:  actionID,
		Duration:  time.Minute * 10,
		MaxStacks: 9,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.SetStacks(sim, 3+(2*shaman.Talents.StaticShock))
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) || !spellEffect.Landed() {
				return
			}
			if sim.RandomFloat("Static Shock") > procChance {
				return
			}
			procSpell.Cast(sim, spellEffect.Target)
			aura.RemoveStack(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) || !spellEffect.Landed() {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			icd.Use(sim)

			procSpell.Cast(sim, spell.Unit)
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
		},
	})
}
