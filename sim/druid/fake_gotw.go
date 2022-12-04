package druid

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// This is 'fake' because it doesnt actually account for any actual buff updating
// this is only used as a 'clearcast fisher' spell
func (druid *Druid) registerFakeGotw() {

	baseCost := druid.BaseMana * core.TernaryFloat64(druid.HasMinorGlyph(proto.DruidMinorGlyph_GlyphOfTheWild), 0.32, 0.64)

	druid.GiftOfTheWild = druid.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 48470},

		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		// Not actually 'Healing' but close enough
		ProcMask: core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
		},
	})
}
