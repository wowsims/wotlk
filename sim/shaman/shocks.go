package shaman

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (shaman *Shaman) ShockCD() time.Duration {
	return time.Second*6 - time.Millisecond*200*time.Duration(shaman.Talents.Reverberation)
}

// Shared logic for all shocks.
func (shaman *Shaman) newShockSpellConfig(spellID int32, spellSchool core.SpellSchool, baseCost float64, shockTimer *core.Timer) core.SpellConfig {
	actionID := core.ActionID{SpellID: spellID}
	cost := baseCost

	return core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagShock,
		ResourceType: stats.Mana,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost -
					baseCost*(core.TernaryFloat64(shaman.Talents.ShamanisticFocus, 0.45, 0)+
						float64(shaman.Talents.Convection)*0.02+
						float64(shaman.Talents.MentalQuickness)*0.02+
						core.TernaryFloat64(shaman.HasSetBonus(ItemSetSkyshatterHarness, 2), 0.1, 0)),
				GCD: core.GCDDefault,
			},
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				shaman.modifyCastClearcasting(spell, cast)
			},
			CD: core.Cooldown{
				Timer:    shockTimer,
				Duration: shaman.ShockCD(),
			},
		},

		BonusHitRating: float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1 +
			0.01*float64(shaman.Talents.Concussion) +
			core.TernaryFloat64(shaman.HasSetBonus(ItemSetThrallsBattlegear, 4), 0.25, 0),
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: shaman.spellThreatMultiplier(),
	}
}

func (shaman *Shaman) registerEarthShockSpell(shockTimer *core.Timer) {
	config := shaman.newShockSpellConfig(49231, core.SpellSchoolNature, 0.18*shaman.BaseMana, shockTimer)
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := sim.Roll(854, 900) + 0.386*spell.SpellPower()
		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	}

	shaman.EarthShock = shaman.RegisterSpell(config)
}

const FlameshockID = 49233

func (shaman *Shaman) registerFlameShockSpell(shockTimer *core.Timer) {
	actionID := core.ActionID{SpellID: FlameshockID}
	config := shaman.newShockSpellConfig(FlameshockID, core.SpellSchoolFire, 0.17*shaman.BaseMana, shockTimer)

	config.Cast.CD.Duration -= time.Duration(shaman.Talents.BoomingEchoes) * time.Second
	config.CritMultiplier = shaman.ElementalCritMultiplier(core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFlameShock), 0.6, 0))

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := 500 + 0.214*spell.SpellPower()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		if result.Landed() {
			shaman.FlameShockDot.Apply(sim)
		}
		spell.DealDamage(sim, result)
	}

	target := shaman.CurrentTarget
	shaman.FlameShockDot = core.NewDot(core.Dot{
		Spell: shaman.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: config.DamageMultiplier +
				float64(shaman.Talents.StormEarthAndFire)*0.2 +
				core.TernaryFloat64(shaman.HasSetBonus(ItemSetWorldbreakerGarb, 2), 0.2, 0),
			CritMultiplier:   config.CritMultiplier,
			ThreatMultiplier: config.ThreatMultiplier,
		}),
		Aura: target.RegisterAura(core.Aura{
			Label:    "FlameShock-" + strconv.Itoa(int(shaman.Index)),
			ActionID: actionID,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.LavaBurst.BonusCritRating += 100 * core.CritRatingPerCritChance
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.LavaBurst.BonusCritRating -= 100 * core.CritRatingPerCritChance
			},
		}),
		NumberOfTicks:       6 + core.TernaryInt32(shaman.HasSetBonus(ItemSetThrallsRegalia, 2), 3, 0),
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: true,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 834/6 + 0.1*dot.Spell.SpellPower()
			dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
		},
	})

	// Apply this talent after creating DoT spell, so it doesn't get copied into periodic DamageMultiplier.
	config.DamageMultiplier += 0.1 * float64(shaman.Talents.BoomingEchoes)
	shaman.FlameShock = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerFrostShockSpell(shockTimer *core.Timer) {
	config := shaman.newShockSpellConfig(49236, core.SpellSchoolFrost, 0.18*shaman.BaseMana, shockTimer)
	config.Cast.CD.Duration -= time.Duration(shaman.Talents.BoomingEchoes) * time.Second
	config.DamageMultiplier += 0.1 * float64(shaman.Talents.BoomingEchoes)
	config.ThreatMultiplier *= 2
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := sim.Roll(812, 858) + 0.386*spell.SpellPower()
		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	}

	shaman.FrostShock = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerShocks() {
	shockTimer := shaman.NewTimer()
	shaman.registerEarthShockSpell(shockTimer)
	shaman.registerFlameShockSpell(shockTimer)
	shaman.registerFrostShockSpell(shockTimer)
}
