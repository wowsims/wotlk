package druid

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// This is 'fake' because it doesnt actually account for any actual buff updating
// this is only used as a 'clearcast fisher' spell
func (druid *Druid) registerFakeGotw() {
	baseCost := core.TernaryFloat64(druid.HasMinorGlyph(proto.DruidMinorGlyph_GlyphOfTheWild), 0.32, 0.64)

	druid.GiftOfTheWild = druid.RegisterSpell(Humanoid|Moonkin|Tree, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 48470},
		Flags:    SpellFlagOmenTrigger | core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   baseCost,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})
}
