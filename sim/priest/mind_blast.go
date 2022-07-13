package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerMindBlastSpell() {
	baseCost := priest.BaseMana * 0.17

	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellHitRating:  0 + float64(priest.Talents.ShadowFocus)*1*core.SpellHitRatingPerHitChance,
		BonusSpellCritRating: float64(priest.Talents.MindMelt) * 2 * core.CritRatingPerCritChance,
		DamageMultiplier:     1,
		ThreatMultiplier:     1 - 0.08*float64(priest.Talents.ShadowAffinity),
		OutcomeApplier:       priest.OutcomeFuncMagicHitAndCrit(priest.SpellCritMultiplier(1, float64(priest.Talents.ShadowPower)/5)),
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				priest.AddShadowWeavingStack(sim)
			}
			if spellEffect.DidCrit() && priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow)) {
				priest.ShadowyInsightAura.Activate(sim)
			}
		},
	}

	normalCalc := core.BaseDamageFuncMagic(997, 1053, 0.429)
	miseryCalc := core.BaseDamageFuncMagic(997, 1053, (1+float64(priest.Talents.Misery)*0.05)*0.429)

	normMod := (1 + float64(priest.Talents.Darkness)*0.02) * // initialize modifier
		core.TernaryFloat64(priest.HasSetBonus(ItemSetAbsolution, 4), 1.1, 1)

	swpMod := (1 + float64(priest.Talents.Darkness)*0.02 + float64(priest.Talents.TwistedFaith)*0.02) * // update modifier if SWP active
		core.TernaryFloat64(priest.HasSetBonus(ItemSetAbsolution, 4), 1.1, 1)

	effect.BaseDamage = core.BaseDamageConfig{
		Calculator: func(sim *core.Simulation, effect *core.SpellEffect, spell *core.Spell) float64 {
			var dmg float64
			shadowWeavingMod := 1 + float64(priest.ShadowWeavingAura.GetStacks())*0.02

			if priest.MiseryAura.IsActive() { // priest.MiseryAura != nil
				dmg = miseryCalc(sim, effect, spell)
			} else {
				dmg = normalCalc(sim, effect, spell)
			}
			if priest.ShadowWordPainDot.IsActive() {
				dmg *= swpMod // multiply the damage
			} else {
				dmg *= normMod // multiply the damage
			}
			return dmg * shadowWeavingMod
		},
		TargetSpellCoefficient: 0.0,
	}

	priest.MindBlast = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48127},
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.05*float64(priest.Talents.FocusedMind)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second*8 - time.Millisecond*500*time.Duration(priest.Talents.ImprovedMindBlast),
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}
