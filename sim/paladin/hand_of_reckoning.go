package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (paladin *Paladin) registerHandOfReckoningSpell() {
	paladin.HandOfReckoning = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 62124},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagMeleeMetrics,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.03,
			Multiplier: 1 - 0.02*float64(paladin.Talents.Benediction),
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * time.Duration(core.TernaryInt(paladin.HasSetBonus(ItemSetTuralyonsPlate, 2), 6, 8)),
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		CritMultiplier:           paladin.SpellCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1 +
				.5*spell.MeleeAttackPower()

			bonusHit := core.TernaryFloat64(
				paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfRighteousDefense),
				8*core.SpellHitRatingPerHitChance,
				0)

			spell.BonusHitRating += bonusHit
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.BonusHitRating -= bonusHit
		},
	})
}
