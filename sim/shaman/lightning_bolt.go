package shaman

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// newLightningBoltSpell returns a precomputed instance of lightning bolt to use for casting.
func (shaman *Shaman) newLightningBoltSpell(isLightningOverload bool) *core.Spell {
	baseCost := baseMana * 0.1
	cost := baseCost
	if shaman.Equip[items.ItemSlotRanged].ID == TotemOfThePulsingEarth {
		cost -= 27.0
	}
	if shaman.HasSetBonus(ItemSetEarthShatterGarb, 2) {
		cost -= baseCost * 0.05
	}

	spellConfig := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: 49238},
		baseCost,
		time.Millisecond*2500,
		isLightningOverload)

	if !isLightningOverload {
		spellConfig.Cast.ModifyCast = func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
			shaman.applyElectricSpellCastInitModifiers(spell, cast)
			if shaman.NaturesSwiftnessAura.IsActive() {
				cast.CastTime = 0
			} else {
				spell.ActionID.Tag = shaman.MaelstromWeaponAura.GetStacks()
				shaman.modifyCastMaelstrom(spell, cast)
			}
		}
	}

	effect := shaman.newElectricSpellEffect(719, 819, 0.7143, isLightningOverload)
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLightningBolt) {
		effect.DamageMultiplier *= 1.04
	}

	if shaman.HasSetBonus(ItemSetSkyshatterRegalia, 4) {
		effect.DamageMultiplier *= 1.05
	}

	has4pT8 := shaman.HasSetBonus(ItemSetWorldbreakerGarb, 4)
	var lbdot *core.Dot
	var applyDot func(sim *core.Simulation, dmg float64)
	if has4pT8 && !isLightningOverload {
		lbdotDmg := 0.0 // dynamically changing dmg
		spell := shaman.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 64930},
			Flags:    core.SpellFlagIgnoreModifiers,
		})
		lbdot = core.NewDot(core.Dot{
			Spell: spell,
			Aura: shaman.CurrentTarget.RegisterAura(core.Aura{
				Label:    "Electrified-" + strconv.Itoa(int(shaman.Index)),
				ActionID: core.ActionID{SpellID: 64930},
			}),
			TickLength:    time.Second * 2,
			NumberOfTicks: 2,
			TickEffects: core.TickFuncSnapshot(shaman.CurrentTarget, core.SpellEffect{
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(_ *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
						return lbdotDmg / 2 //spread dot over 2 ticks
					},
				},
				IsPeriodic:     true,
				ProcMask:       core.ProcMaskEmpty,
				OutcomeApplier: shaman.OutcomeFuncTick(),
			}),
		})
		applyDot = func(sim *core.Simulation, dmg float64) {
			lbdotDmg = dmg * 0.08 // TODO: does this pool with a currently ticking dot?
			lbdot.Apply(sim)      // will resnapshot
		}
		effect.OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.DidCrit() {
				return
			}
			applyDot(sim, spellEffect.Damage)
		}
	}

	if !isLightningOverload && shaman.Talents.LightningOverload > 0 {
		lightningOverloadChance := float64(shaman.Talents.LightningOverload) * 0.11
		effect.OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			if has4pT8 && spellEffect.DidCrit() { // need to merge in the 4pt8 effect
				applyDot(sim, spellEffect.Damage)
			}
			if sim.RandomFloat("LB Lightning Overload") > lightningOverloadChance {
				return
			}
			shaman.LightningBoltLO.Cast(sim, spellEffect.Target)
		}
	}

	spellConfig.ApplyEffects = core.ApplyEffectFuncDirectDamage(effect)
	return shaman.RegisterSpell(spellConfig)
}
