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
		Label:    "Haunt-" + warlock.Label,
		ActionID: core.ActionID{SpellID: 59164},
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AttackTables[aura.Unit.UnitIndex].PeriodicShadowDamageDealtMultiplier *= shadowDotMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AttackTables[aura.Unit.UnitIndex].PeriodicShadowDamageDealtMultiplier /= shadowDotMultiplier
		},
	})
}

func (warlock *Warlock) registerHauntSpell() {
	actionID := core.ActionID{SpellID: 59164}
	spellSchool := core.SpellSchoolShadow
	HauntDebuffAura := warlock.HauntDebuffAura(warlock.CurrentTarget)
	baseCost := 0.12 * warlock.BaseMana

	warlock.Haunt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 20,

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

		BonusCritRating: 0 +
			warlock.masterDemonologistShadowCrit(),
		DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false),
		CritMultiplier:           warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Pandemic, 1, 0)),
		ThreatMultiplier:         1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(645, 753) + 0.4286*spell.SpellPower()
			result := spell.CalcDamageMagicHitAndCrit(sim, target, baseDamage)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					HauntDebuffAura.Activate(sim)
				}
				spell.DealDamage(sim, &result)
			})
		},
	})
}
