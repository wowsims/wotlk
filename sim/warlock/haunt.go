package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerHauntSpell() {
	actionID := core.ActionID{SpellID: 59164}
	spellSchool := core.SpellSchoolShadow
	baseAdditiveMultiplier := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false)
	shadowDotMultiplier := 1.2
	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfHaunt) {
		shadowDotMultiplier += 0.03
	}

	warlock.HauntAura = warlock.RegisterAura(core.Aura{
		Label:    "Haunt Buff",
		ActionID: actionID,
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.PeriodicShadowDamageDealtMultiplier *= shadowDotMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.PeriodicShadowDamageDealtMultiplier /= shadowDotMultiplier
		},
	})

	effect := core.SpellEffect{
		ProcMask:         core.ProcMaskSpellDamage,
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		DamageMultiplier: baseAdditiveMultiplier,
		BaseDamage:       core.BaseDamageConfigMagic(645.0, 753.0, 0.4286),
		OutcomeApplier:   warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Pandemic, 1, 0))),
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			warlock.HauntAura.Activate(sim)
		},
	}

	baseCost := 0.12 * warlock.BaseMana
	warlock.Haunt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: spellSchool,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

}
