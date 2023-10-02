package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (shaman *Shaman) registerChainLightningSpell() {
	numHits := min(core.TernaryInt32(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning), 4, 3), shaman.Env.GetNumTargets())
	shaman.ChainLightning = shaman.newChainLightningSpell(false)
	shaman.ChainLightningLOs = []*core.Spell{}
	for i := int32(0); i < numHits; i++ {
		shaman.ChainLightningLOs = append(shaman.ChainLightningLOs, shaman.newChainLightningSpell(true))
	}
}

func (shaman *Shaman) newChainLightningSpell(isLightningOverload bool) *core.Spell {
	spellConfig := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: 49271},
		0.26,
		time.Millisecond*2000,
		isLightningOverload)

	if !isLightningOverload {
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second*6 - []time.Duration{0, 750 * time.Millisecond, 1500 * time.Millisecond, 2500 * time.Millisecond}[shaman.Talents.StormEarthAndFire],
		}
	}

	numHits := min(core.TernaryInt32(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning), 4, 3), shaman.Env.GetNumTargets())
	dmgReductionPerBounce := core.TernaryFloat64(shaman.HasSetBonus(ItemSetTidefury, 2), 0.83, 0.7)
	dmgBonus := shaman.electricSpellBonusDamage(0.5714)
	spellCoeff := 0.5714 + 0.04*float64(shaman.Talents.Shamanism)

	canLO := !isLightningOverload && shaman.Talents.LightningOverload > 0
	lightningOverloadChance := float64(shaman.Talents.LightningOverload) * 0.11 / 3

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		bounceCoeff := 1.0
		curTarget := target
		for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
			baseDamage := dmgBonus + sim.Roll(973, 1111) + spellCoeff*spell.SpellPower()
			baseDamage *= bounceCoeff
			result := spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

			if canLO && result.Landed() && sim.RandomFloat("CL Lightning Overload") <= lightningOverloadChance {
				shaman.ChainLightningLOs[hitIndex].Cast(sim, curTarget)
			}

			spell.DealDamage(sim, result)

			bounceCoeff *= dmgReductionPerBounce
			curTarget = sim.Environment.NextTargetUnit(curTarget)
		}
	}

	return shaman.RegisterSpell(spellConfig)
}
