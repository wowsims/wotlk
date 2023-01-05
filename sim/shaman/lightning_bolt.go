package shaman

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (shaman *Shaman) registerLightningBoltSpell() {
	shaman.LightningBolt = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(false))
	shaman.LightningBoltLO = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(true))
}

func (shaman *Shaman) RegisterMaelstromLightningBoltSpells(minStacks int32) []*core.Spell {
	var spells []*core.Spell

	spellConfig := shaman.newLightningBoltSpellConfig(false)

	for i := minStacks; i <= 5; i++ {
		spellConfig.ActionID.Tag = i
		spell := shaman.RegisterSpell(spellConfig)
		spells = append(spells, spell)
	}

	return spells
}

func (shaman *Shaman) newLightningBoltSpellConfig(isLightningOverload bool) core.SpellConfig {
	baseCost := 0.1 * shaman.BaseMana
	if shaman.HasSetBonus(ItemSetEarthShatterGarb, 2) {
		baseCost -= baseCost * 0.05
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

	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLightningBolt) {
		spellConfig.DamageMultiplier += 0.04
	}

	if shaman.HasSetBonus(ItemSetSkyshatterRegalia, 4) {
		spellConfig.DamageMultiplier += 0.05
	}

	if !isLightningOverload && shaman.HasSetBonus(ItemSetWorldbreakerGarb, 4) && shaman.LightningBoltDot == nil {
		lbDotSpell := shaman.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 64930},
			SpellSchool:      core.SpellSchoolNature,
			ProcMask:         core.ProcMaskEmpty,
			Flags:            core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				shaman.LightningBoltDot.ApplyOrReset(sim)
				spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			},
		})

		shaman.LightningBoltDot = core.NewDot(core.Dot{
			Spell: lbDotSpell,
			Aura: shaman.CurrentTarget.RegisterAura(core.Aura{
				Label:    "Electrified-" + strconv.Itoa(int(shaman.Index)),
				ActionID: core.ActionID{SpellID: 64930},
			}),
			TickLength:    time.Second * 2,
			NumberOfTicks: 2,

			SnapshotAttackerMultiplier: 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		})
	}

	dmgBonus := shaman.electricSpellBonusDamage(0.7143)
	spellCoeff := 0.7143 + 0.04*float64(shaman.Talents.Shamanism)

	canLO := !isLightningOverload && shaman.Talents.LightningOverload > 0
	lightningOverloadChance := float64(shaman.Talents.LightningOverload) * 0.11
	lbDot := shaman.LightningBoltDot
	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := dmgBonus + sim.Roll(719, 819) + spellCoeff*spell.SpellPower()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		if !isLightningOverload && lbDot != nil && result.DidCrit() {

			var outstandingDamage float64
			if lbDot.IsActive() {
				outstandingDamage = lbDot.SnapshotBaseDamage * float64(lbDot.NumberOfTicks-lbDot.TickCount)
			}

			newDamage := result.Damage * 0.08
			lbDot.SnapshotBaseDamage = (outstandingDamage + newDamage) / float64(lbDot.NumberOfTicks)
			lbDot.Spell.Cast(sim, target)
		}

		if canLO && result.Landed() && sim.RandomFloat("LB Lightning Overload") < lightningOverloadChance {
			shaman.LightningBoltLO.Cast(sim, target)
		}

		spell.DealDamage(sim, result)
	}

	return spellConfig
}
