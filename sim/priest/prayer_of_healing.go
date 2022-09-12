package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerPrayerOfHealingSpell() {
	baseCost := .48 * priest.BaseMana

	baseEffect := core.SpellEffect{
		IsHealing: true,
		ProcMask:  core.ProcMaskSpellHealing,

		BonusCritRating: 0 +
			1*float64(priest.Talents.HolySpecialization)*core.CritRatingPerCritChance +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetSanctificationRegalia, 2), 10*core.CritRatingPerCritChance, 0),
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.DivineProvidence)),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		BaseDamage:     core.BaseDamageConfigHealing(2109, 2228, 0.526),
		OutcomeApplier: priest.OutcomeFuncHealingCrit(priest.DefaultHealingCritMultiplier()),
	}

	// Separate ApplyEffects functions for each party.
	var applyPartyEffects []core.ApplySpellEffects

	priest.PrayerOfHealing = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48072},
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 -
					.1*float64(priest.Talents.HealingPrayers) -
					core.TernaryFloat64(priest.HasSetBonus(ItemSetVestmentsOfAbsolution, 2), 0.1, 0)),
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			targetAgent := target.Env.Raid.GetPlayerFromUnitIndex(target.UnitIndex)
			party := targetAgent.GetCharacter().Party
			applyPartyEffects[party.Index](sim, target, spell)
		},
	})

	for _, party := range priest.Env.Raid.Parties {
		var effects []core.SpellEffect
		var hots []*core.Dot
		for _, target := range party.PlayersAndPets {
			effect := baseEffect
			effect.Target = &target.GetCharacter().Unit
			effects = append(effects, effect)

			if priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfPrayerOfHealing) {
				hots = append(hots, priest.makePrayerOfHealingGlyphHot(effect.Target, baseEffect))
			}
		}

		if len(effects) == 0 {
			applyPartyEffects = append(applyPartyEffects, nil)
			continue
		}

		applyDirectEffects := core.ApplyEffectFuncDamageMultiple(effects)
		if priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfPrayerOfHealing) {
			applyPartyEffects = append(applyPartyEffects, func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				applyDirectEffects(sim, target, spell)
				for _, hot := range hots {
					hot.Activate(sim)
				}
			})
		} else {
			applyPartyEffects = append(applyPartyEffects, applyDirectEffects)
		}
	}
}

func (priest *Priest) makePrayerOfHealingGlyphHot(target *core.Unit, pohEffect core.SpellEffect) *core.Dot {
	spell := priest.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{ItemID: 42409},
		SpellSchool: core.SpellSchoolHoly,
	})

	return core.NewDot(core.Dot{
		Spell: priest.PrayerOfHealing,
		Aura: target.RegisterAura(core.Aura{
			Label:    "PoH Glyph" + strconv.Itoa(int(priest.Index)),
			ActionID: spell.ActionID,
		}),
		NumberOfTicks: 2,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:   core.ProcMaskPeriodicHealing,
			IsPeriodic: true,
			IsHealing:  true,

			DamageMultiplier: pohEffect.DamageMultiplier * 0.2 / 2,
			ThreatMultiplier: pohEffect.ThreatMultiplier,

			BaseDamage:     pohEffect.BaseDamage,
			OutcomeApplier: priest.OutcomeFuncTick(),
		}),
	})
}
