package shaman

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (shaman *Shaman) ShockCD() time.Duration {
	return time.Second*6 - time.Millisecond*200*time.Duration(shaman.Talents.Reverberation)
}

// Shared logic for all shocks.
func (shaman *Shaman) newShockSpellConfig(spellID int32, spellSchool core.SpellSchool, baseCost float64, shockTimer *core.Timer) (core.SpellConfig, core.SpellEffect) {
	actionID := core.ActionID{SpellID: spellID}

	cost := baseCost
	if shaman.Talents.ShamanisticFocus {
		cost -= baseCost * 0.45
	}
	return core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: spellSchool,
			Flags:       SpellFlagShock,

			ResourceType: stats.Mana,
			BaseCost:     cost,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					Cost: baseCost -
						baseCost*float64(shaman.Talents.Convection)*0.02 -
						baseCost*float64(shaman.Talents.MentalQuickness)*0.02 -
						core.TernaryFloat64(shaman.HasSetBonus(ItemSetSkyshatterHarness, 2), baseCost*0.1, 0),
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
		}, core.SpellEffect{
			ProcMask:            core.ProcMaskSpellDamage,
			BonusSpellHitRating: float64(shaman.Talents.ElementalPrecision) * 2 * core.SpellHitRatingPerHitChance,
			BonusSpellPower: 0 +
				core.TernaryFloat64(shaman.Equip[items.ItemSlotRanged].ID == TotemOfRage, 30, 0) +
				core.TernaryFloat64(shaman.Equip[items.ItemSlotRanged].ID == TotemOfImpact, 46, 0),
			DamageMultiplier: 1 * (1 + 0.01*float64(shaman.Talents.Concussion)),
			ThreatMultiplier: 1 - (0.1/3)*float64(shaman.Talents.ElementalPrecision),
		}
}

func (shaman *Shaman) registerEarthShockSpell(shockTimer *core.Timer) {
	config, effect := shaman.newShockSpellConfig(25454, core.SpellSchoolNature, 535.0, shockTimer)
	config.Flags |= core.SpellFlagBinary

	effect.BaseDamage = core.BaseDamageConfigMagic(661, 696, 0.386)
	effect.OutcomeApplier = shaman.OutcomeFuncMagicHitAndCritBinary(shaman.ElementalCritMultiplier())
	config.ApplyEffects = core.ApplyEffectFuncDirectDamage(effect)

	shaman.EarthShock = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerFlameShockSpell(shockTimer *core.Timer) {
	const flameshockID = 49233
	config, effect := shaman.newShockSpellConfig(flameshockID, core.SpellSchoolFire, baseMana*0.17, shockTimer)

	config.Cast.CD.Duration -= time.Duration(shaman.Talents.BoomingEchoes) * time.Second

	effect.BaseDamage = core.BaseDamageConfigMagic(500, 500, 0.214)
	critMult := shaman.ElementalCritMultiplier()
	effect.OutcomeApplier = shaman.OutcomeFuncMagicHitAndCrit(critMult)
	if effect.OnSpellHitDealt == nil {
		effect.OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				shaman.FlameShockDot.Apply(sim)
			}
		}
	} else {
		oldSpellHit := effect.OnSpellHitDealt
		effect.OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			oldSpellHit(sim, spell, spellEffect)
			if spellEffect.Landed() {
				shaman.FlameShockDot.Apply(sim)
			}
		}
	}

	config.ApplyEffects = core.ApplyEffectFuncDirectDamage(effect)
	shaman.FlameShock = shaman.RegisterSpell(config)

	target := shaman.CurrentTarget
	shaman.FlameShockDot = core.NewDot(core.Dot{
		Spell: shaman.FlameShock,
		Aura: target.RegisterAura(core.Aura{
			Label:    "FlameShock-" + strconv.Itoa(int(shaman.Index)),
			ActionID: core.ActionID{SpellID: flameshockID},
		}),
		NumberOfTicks:       4,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: true,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			DamageMultiplier: 1 * (1 + 0.01*float64(shaman.Talents.Concussion)) * (1.0 + float64(shaman.Talents.StormEarthAndFire)*0.2), // 20% bonus dmg per SE&F
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(834/6, 0.1),
			OutcomeApplier:   shaman.OutcomeFuncMagicCrit(critMult),
			IsPeriodic:       true,
			ProcMask:         core.ProcMaskPeriodicDamage,
		}),
	})
}

func (shaman *Shaman) registerFrostShockSpell(shockTimer *core.Timer) {
	config, effect := shaman.newShockSpellConfig(25464, core.SpellSchoolFrost, 525.0, shockTimer)
	config.Flags |= core.SpellFlagBinary
	config.Cast.CD.Duration -= time.Duration(shaman.Talents.BoomingEchoes) * time.Second

	effect.ThreatMultiplier *= 2
	effect.BaseDamage = core.BaseDamageConfigMagic(647, 683, 0.386)
	effect.OutcomeApplier = shaman.OutcomeFuncMagicHitAndCritBinary(shaman.ElementalCritMultiplier())
	config.ApplyEffects = core.ApplyEffectFuncDirectDamage(effect)

	shaman.FrostShock = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerShocks() {
	shockTimer := shaman.NewTimer()
	shaman.registerEarthShockSpell(shockTimer)
	shaman.registerFlameShockSpell(shockTimer)
	shaman.registerFrostShockSpell(shockTimer)
}
