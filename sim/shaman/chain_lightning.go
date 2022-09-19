package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (shaman *Shaman) registerChainLightningSpell() {
	numHits := core.MinInt32(core.TernaryInt32(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning), 4, 3), shaman.Env.GetNumTargets())
	shaman.ChainLightning = shaman.newChainLightningSpell(false)
	shaman.ChainLightningLOs = []*core.Spell{}
	for i := int32(0); i < numHits; i++ {
		shaman.ChainLightningLOs = append(shaman.ChainLightningLOs, shaman.newChainLightningSpell(true))
	}
}

func (shaman *Shaman) newChainLightningSpell(isLightningOverload bool) *core.Spell {
	spellConfig := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: 49271},
		baseMana*0.26,
		time.Millisecond*2000,
		isLightningOverload)

	if !isLightningOverload {
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second*6 - []time.Duration{0, 750 * time.Millisecond, 1500 * time.Millisecond, 2500 * time.Millisecond}[shaman.Talents.StormEarthAndFire],
		}
	}

	spellConfig.Cast.ModifyCast = func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
		shaman.applyElectricSpellCastInitModifiers(spell, cast)
		if shaman.NaturesSwiftnessAura.IsActive() {
			cast.CastTime = 0
		} else {
			shaman.modifyCastMaelstrom(spell, cast)
		}
	}

	effect := shaman.newElectricSpellEffect(973, 1111, 0.5714, isLightningOverload)

	makeOnSpellHit := func(hitIndex int32) func(*core.Simulation, *core.Spell, *core.SpellEffect) {
		if !isLightningOverload && shaman.Talents.LightningOverload > 0 {
			lightningOverloadChance := float64(shaman.Talents.LightningOverload) * 0.11 / 3
			return func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}
				if sim.RandomFloat("CL Lightning Overload") > lightningOverloadChance {
					return
				}
				shaman.ChainLightningLOs[hitIndex].Cast(sim, spellEffect.Target)
			}
		} else {
			return nil
		}
	}

	hasTidefury := shaman.HasSetBonus(ItemSetTidefury, 2)
	numHits := core.MinInt32(core.TernaryInt32(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning), 4, 3), shaman.Env.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)

	effect.Target = shaman.Env.GetTargetUnit(0)
	effect.OnSpellHitDealt = makeOnSpellHit(0)
	effects = append(effects, effect)

	bounceMult := 1.0
	for i := int32(1); i < numHits; i++ {
		bounceEffect := effects[i-1] // Makes a copy of the previous bounce
		bounceEffect.Target = shaman.Env.GetTargetUnit(i)
		if hasTidefury {
			bounceMult *= 0.83
		} else {
			bounceMult *= 0.7
		}
		curBounceMult := bounceMult
		bounceEffect.BaseDamage = core.WrapBaseDamageConfig(bounceEffect.BaseDamage, func(oldCalc core.BaseDamageCalculator) core.BaseDamageCalculator {
			return func(sim *core.Simulation, effect *core.SpellEffect, spell *core.Spell) float64 {
				return oldCalc(sim, effect, spell) * curBounceMult
			}
		})
		bounceEffect.OnSpellHitDealt = makeOnSpellHit(i)

		effects = append(effects, bounceEffect)
	}

	spellConfig.ApplyEffects = core.ApplyEffectFuncDamageMultiple(effects)
	return shaman.RegisterSpell(spellConfig)
}
