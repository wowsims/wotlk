package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
)

// newLightningBoltSpell returns a precomputed instance of lightning bolt to use for casting.
func (shaman *Shaman) newLightningBoltSpell(isLightningOverload bool) *core.Spell {
	baseCost := baseMana * 0.1
	if shaman.Equip[items.ItemSlotRanged].ID == TotemOfThePulsingEarth {
		baseCost -= 27.0
	}

	spellConfig := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: 49238},
		baseCost,
		time.Millisecond*2500,
		isLightningOverload)

	spellConfig.Cast.ModifyCast = func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
		shaman.applyElectricSpellCastInitModifiers(spell, cast)
		if shaman.NaturesSwiftnessAura != nil && shaman.NaturesSwiftnessAura.IsActive() {
			cast.CastTime = 0
		} else {
			shaman.modifyCastMaelstrom(spell, cast)
		}
	}

	effect := shaman.newElectricSpellEffect(719, 819, 0.7143, isLightningOverload)

	if ItemSetSkyshatterRegalia.CharacterHasSetBonus(&shaman.Character, 4) {
		effect.DamageMultiplier *= 1.05
	}

	if !isLightningOverload && shaman.Talents.LightningOverload > 0 {
		lightningOverloadChance := float64(shaman.Talents.LightningOverload) * 0.11
		effect.OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
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
