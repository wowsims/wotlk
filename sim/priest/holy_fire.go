package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) RegisterHolyFireSpell(memeDream bool) {
	actionID := core.ActionID{SpellID: 48135}
	baseCost := .11 * priest.BaseMana

	priest.HolyFire = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*2000 - time.Millisecond*100*time.Duration(priest.Talents.DivineFury),
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,

			BonusSpellCritRating: float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
			DamageMultiplier:     1 + 0.05*float64(priest.Talents.SearingLight),
			ThreatMultiplier:     1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

			BaseDamage:     core.BaseDamageConfigMagic(900, 1140, 0.5711),
			OutcomeApplier: priest.OutcomeFuncMagicHitAndCrit(priest.DefaultSpellCritMultiplier()),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					priest.HolyFireDot.Apply(sim)
				}
			},
		}),
	})

	hasGlyph := !memeDream && priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfSmite)

	target := priest.CurrentTarget
	priest.HolyFireDot = core.NewDot(core.Dot{
		Spell: priest.HolyFire,
		Aura: target.RegisterAura(core.Aura{
			Label:    "HolyFire-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if hasGlyph {
					priest.Smite.DamageMultiplier *= 1.2
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if hasGlyph {
					priest.Smite.DamageMultiplier /= 1.2
				}
			},
		}),
		NumberOfTicks: 7,
		TickLength:    time.Second * 1,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,

			DamageMultiplier: 1 + 0.05*float64(priest.Talents.SearingLight),
			ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

			BaseDamage:     core.BaseDamageConfigMagicNoRoll(50, 0.024),
			OutcomeApplier: priest.OutcomeFuncTick(),
			IsPeriodic:     true,
		}),
	})
}
