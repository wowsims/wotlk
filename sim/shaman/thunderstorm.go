package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var thunderstormActionID = core.ActionID{SpellID: 51490}

// newThunderstormSpell returns a precomputed instance of lightning bolt to use for casting.
func (shaman *Shaman) newThunderstormSpell() *core.Spell {
	manaRestore := 0.08
	if shaman.HasMinorGlyph(proto.ShamanMinorGlyph_GlyphOfThunderstorm) {
		manaRestore = 0.1
	}

	manaMetrics := shaman.NewManaMetrics(thunderstormActionID)

	spellConfig := core.SpellConfig{
		ActionID:    thunderstormActionID,
		SpellSchool: core.SpellSchoolNature,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 45,
			},
		},
		ApplyEffects: func(sim *core.Simulation, u *core.Unit, s2 *core.Spell) {
			shaman.AddMana(sim, shaman.MaxMana()*manaRestore, manaMetrics, true)
		},
	}

	// TODO: Add option to specify if in range of thunderstorm.

	// effect := core.SpellEffect{
	// 	ProcMask:             core.ProcMaskSpellDamage,
	// 	BonusSpellHitRating:  float64(shaman.Talents.ElementalPrecision) * 2 * core.SpellHitRatingPerHitChance,
	// 	BonusSpellCritRating: 0,
	// 	BonusSpellPower: 0 +
	// 		core.TernaryFloat64(shaman.Equip[items.ItemSlotRanged].ID == TotemOfHex, 165, 0),
	// 	DamageMultiplier: 1 * (1 + 0.01*float64(shaman.Talents.Concussion)),
	// 	ThreatMultiplier: 1 - (0.1/3)*float64(shaman.Talents.ElementalPrecision),
	// 	BaseDamage:       core.BaseDamageConfigMagic(1192, 1518, 0.5714),
	// }

	return shaman.RegisterSpell(spellConfig)
}
