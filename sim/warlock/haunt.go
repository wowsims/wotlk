package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) HauntDebuffAura(target *core.Unit) *core.Aura {
	shadowDotMultiplier := 1.2
	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfHaunt) {
		shadowDotMultiplier += 0.03
	}

	return target.GetOrRegisterAura(core.Aura{
		Label:     "Haunt-" + warlock.Label,
		ActionID:  core.ActionID{SpellID: 59164},
		Duration:  time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AttackTables[aura.Unit.TableIndex].PeriodicShadowDamageDealtMultiplier *= shadowDotMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AttackTables[aura.Unit.TableIndex].PeriodicShadowDamageDealtMultiplier /= shadowDotMultiplier
		},
	})
}

func (warlock *Warlock) registerHauntSpell() {
	actionID := core.ActionID{SpellID: 59164}
	spellSchool := core.SpellSchoolShadow
	baseAdditiveMultiplier := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false)
	HauntDebuffAura := warlock.HauntDebuffAura(warlock.CurrentTarget)
	
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
			HauntDebuffAura.Activate(sim)
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
