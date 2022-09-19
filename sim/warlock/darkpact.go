package warlock

import (
	"math"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerDarkPactSpell() {
	actionID := core.ActionID{SpellID: 59092}
	baseRestore := 1200.0
	manaMetrics := warlock.NewManaMetrics(actionID)
	petManaMetrics := warlock.Pet.NewManaMetrics(actionID)

	warlock.DarkPact = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		FlatThreatBonus:  80,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:       core.ProcMaskEmpty,
			OutcomeApplier: warlock.OutcomeFuncAlwaysHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// Glyph activates and applies SP before coef calculations are done
				if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap) {
					warlock.GlyphOfLifeTapAura.Activate(sim)
				}

				maxDrain := baseRestore + (warlock.GetStat(stats.SpellPower)+warlock.GetStat(stats.ShadowSpellPower))*0.96
				actualDrain := math.Min(maxDrain, warlock.Pet.CurrentMana())

				warlock.Pet.SpendMana(sim, actualDrain, petManaMetrics)
				warlock.AddMana(sim, actualDrain, manaMetrics, true)
			},
		}),
	})
}
