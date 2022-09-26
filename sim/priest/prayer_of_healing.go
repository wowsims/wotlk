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

	var glyphHots []*core.Dot

	priest.PrayerOfHealing = priest.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48072},
		SpellSchool:  core.SpellSchoolHoly,
		ProcMask:     core.ProcMaskSpellHealing,
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

		BonusCritRating: 0 +
			1*float64(priest.Talents.HolySpecialization)*core.CritRatingPerCritChance +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetSanctificationRegalia, 2), 10*core.CritRatingPerCritChance, 0),
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.DivineProvidence)),
		CritMultiplier:   priest.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			targetAgent := target.Env.Raid.GetPlayerFromUnitIndex(target.UnitIndex)
			party := targetAgent.GetCharacter().Party

			healFromSP := 0.526 * spell.HealingPower()
			for _, partyAgent := range party.PlayersAndPets {
				partyTarget := partyAgent.GetCharacter()
				baseHealing := sim.Roll(2109, 2228) + healFromSP
				spell.CalcAndDealHealingCrit(sim, &partyTarget.Unit, baseHealing)
				if glyphHots != nil {
					glyphHots[partyTarget.UnitIndex].Activate(sim)
				}
			}
		},
	})

	if priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfPrayerOfHealing) {
		glyphHots := make([]*core.Dot, len(priest.Env.AllUnits))
		for _, unit := range priest.Env.AllUnits {
			if !priest.IsOpponent(unit) {
				glyphHots[unit.UnitIndex] = priest.makePrayerOfHealingGlyphHot(unit)
			}
		}
	}
}

func (priest *Priest) makePrayerOfHealingGlyphHot(target *core.Unit) *core.Dot {
	spell := priest.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{ItemID: 42409},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,

		DamageMultiplier: priest.PrayerOfHealing.DamageMultiplier * 0.2 / 2,
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],
	})

	return core.NewDot(core.Dot{
		Spell: spell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "PoH Glyph" + strconv.Itoa(int(priest.Index)),
			ActionID: spell.ActionID,
		}),
		NumberOfTicks: 2,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			IsPeriodic: true,
			IsHealing:  true,

			BaseDamage:     core.BaseDamageConfigHealing(2109, 2228, 0.526),
			OutcomeApplier: priest.OutcomeFuncTick(),
		}),
	})
}
