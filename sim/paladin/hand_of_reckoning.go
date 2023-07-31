package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (paladin *Paladin) registerHandOfReckoningSpell() {
	if !paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfReckoning) {
		return
	}

	paladin.HandOfReckoning = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 67485}, // 62124 is the "taunt" part
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

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
			baseDamage := 1 + .5*spell.MeleeAttackPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit) // cannot miss
		},
	})
}
