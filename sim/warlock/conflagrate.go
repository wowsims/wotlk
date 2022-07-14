package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) CanConflagrate(sim *core.Simulation) bool {
	return warlock.Talents.Conflagrate && warlock.ImmolateDot.IsActive() && warlock.Conflagrate.IsReady(sim)
}

func (warlock *Warlock) registerConflagrateSpell() {

	baseCost := 0.16 * warlock.BaseMana
	costReduction := 0.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReduction += 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}
	actionID := core.ActionID{SpellID: 17962}
	target := warlock.CurrentTarget

	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellCritRating: core.TernaryFloat64(warlock.Talents.Devastation, 1, 0) * 5 * core.CritRatingPerCritChance,
		DamageMultiplier: 0.6 * (1 + (0.1 * float64(warlock.Talents.ImprovedImmolate))) *
			(1 + 0.03*float64(warlock.Talents.Aftermath)) * (1 + 0.03*float64(warlock.Talents.Emberstorm)),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),
		BaseDamage:       core.BaseDamageConfigMagicNoRoll(785, 0.2*5),
		OutcomeApplier:   warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
		OnSpellHitDealt:  applyDotOnLanded(&warlock.ConflagrateDot),
	}

	warlock.Conflagrate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - costReduction),
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 10,
			},
			OnCastComplete: func(sim *core.Simulation, spell *core.Spell) {
				if !warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate) {
					warlock.ImmolateDot.Deactivate(sim)
					//warlock.ShadowflameDot.Deactivate(sim)
				}
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

	warlock.ConflagrateDot = core.NewDot(core.Dot{
		Spell: warlock.Conflagrate,
		Aura: target.RegisterAura(core.Aura{
			Label:    "conflagrate-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 3,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			DamageMultiplier: 0.4 * (1 + 0.03*float64(warlock.Talents.Aftermath)) * (1 + 0.03*float64(warlock.Talents.Emberstorm)),
			ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(785/3, 0.2*5/3),
			OutcomeApplier:   warlock.OutcomeFuncTick(),
			IsPeriodic:       true,
			ProcMask:         core.ProcMaskPeriodicDamage,
		}),
	})
}
