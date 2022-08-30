package warlock

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerLifeTapSpell() {
	actionID := core.ActionID{SpellID: 57946}
	baseRestore := 2000.0 * (1.0 + 0.1*float64(warlock.Talents.ImprovedLifeTap))
	manaMetrics := warlock.NewManaMetrics(actionID)

	petRestore := core.TernaryFloat64(warlock.Talents.ManaFeed, 1, 0)
	var petManaMetrics []*core.ResourceMetrics
	if warlock.Talents.ManaFeed {
		for _, pet := range warlock.Pets {
			petManaMetrics = append(petManaMetrics, pet.GetPet().NewManaMetrics(actionID))
		}
	}

	warlock.LifeTap = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,
			OutcomeApplier:   warlock.OutcomeFuncAlwaysHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// Life tap adds 0.5*sp to mana restore
				restore := baseRestore + (warlock.GetStat(stats.SpellPower)+warlock.GetStat(stats.ShadowSpellPower))*0.5
				warlock.AddMana(sim, restore, manaMetrics, true)

				if warlock.Talents.ManaFeed {
					for i, pet := range warlock.Pets {
						pet.GetPet().AddMana(sim, restore*petRestore, petManaMetrics[i], true)
					}
				}
				if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap) {
					if warlock.GlyphOfLifeTapAura.IsActive() {
						warlock.GlyphOfLifeTapAura.Refresh(sim)
					} else {
						warlock.GlyphOfLifeTapAura.Activate(sim)
					}
				}
			},
		}),
	})
}
