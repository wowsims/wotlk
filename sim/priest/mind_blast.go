package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerMindBlastSpell() {
	baseCost := priest.BaseMana() * 0.17

	Mult_mod := (1 + float64(priest.Talents.Darkness)*0.02) *
		core.TernaryFloat64(priest.Talents.Shadowform, 1.15, 1) *
		core.TernaryFloat64(ItemSetAbsolution.CharacterHasSetBonus(&priest.Character, 4), 1.1, 1)
	if priest.ShadowWordPainDot.IsActive() {
		Mult_mod = (1 + float64(priest.Talents.Darkness)*0.02 + float64(priest.Talents.TwistedFaith)*0.02) *
			core.TernaryFloat64(priest.Talents.Shadowform, 1.15, 1) *
			core.TernaryFloat64(ItemSetAbsolution.CharacterHasSetBonus(&priest.Character, 4), 1.1, 1)
	}

	base := core.BaseDamageConfigMagic(997, 1053, 0.429)
	//if priest.MiseryAura.IsActive() {
	if priest.MiseryAura != nil {
		base = core.BaseDamageConfigMagic(997, 1053, 0.429*float64(priest.Talents.Misery)*0.05)
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

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskSpellDamage,
			BonusSpellHitRating: 0 + float64(priest.Talents.ShadowFocus)*1*core.SpellHitRatingPerHitChance,

			BonusSpellCritRating: float64(priest.Talents.MindMelt) * 2 * core.CritRatingPerCritChance,

			DamageMultiplier: Mult_mod,

			ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),
			BaseDamage:       base,
			OutcomeApplier:   priest.OutcomeFuncMagicHitAndCrit(priest.SpellCritMultiplier(1, float64(priest.Talents.ShadowPower)/5)),
		}),
	})
}

// Need to add a check to see if VT is active, and if so, then apply replenishment
