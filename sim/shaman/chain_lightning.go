package shaman

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (shaman *Shaman) newChainLightningSpell(isLightningOverload bool) *core.Spell {
	spellConfig := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: 25442},
		760.0,
		time.Millisecond*2000,
		isLightningOverload)

	if !isLightningOverload {
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 6,
		}
	}

	spellConfig.Cast.ModifyCast = func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
		shaman.applyElectricSpellCastInitModifiers(spell, cast)
	}

	effect := shaman.newElectricSpellEffect(734, 838, 0.651, isLightningOverload)

	makeOnSpellHit := func(hitIndex int32) func(*core.Simulation, *core.Spell, *core.SpellEffect) {
		if !isLightningOverload && shaman.Talents.LightningOverload > 0 {
			lightningOverloadChance := float64(shaman.Talents.LightningOverload) * 0.04 / 3
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

	hasTidefury := ItemSetTidefury.CharacterHasSetBonus(&shaman.Character, 2)
	numHits := core.MinInt32(3, shaman.Env.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)

	effect.Target = shaman.Env.GetTargetUnit(0)
	effect.OnSpellHitDealt = makeOnSpellHit(0)
	effects = append(effects, effect)

	for i := int32(1); i < numHits; i++ {
		bounceEffect := effects[i-1] // Makes a copy of the previous bounce
		bounceEffect.Target = shaman.Env.GetTargetUnit(i)
		if hasTidefury {
			bounceEffect.DamageMultiplier *= 0.83
		} else {
			bounceEffect.DamageMultiplier *= 0.7
		}
		bounceEffect.OnSpellHitDealt = makeOnSpellHit(i)

		effects = append(effects, bounceEffect)
	}

	spellConfig.ApplyEffects = core.ApplyEffectFuncDamageMultiple(effects)
	return shaman.RegisterSpell(spellConfig)
}
