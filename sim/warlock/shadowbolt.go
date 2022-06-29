package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warlock *Warlock) registerShadowboltSpell() {
	has4pMal := ItemSetMaleficRaiment.CharacterHasSetBonus(&warlock.Character, 4)

	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellCritRating: float64(warlock.Talents.Devastation) * 1 * core.SpellCritRatingPerCritChance,
		DamageMultiplier:     1 * core.TernaryFloat64(has4pMal, 1.06, 1.0) * (1 + 0.02*float64(warlock.Talents.ShadowMastery)),
		ThreatMultiplier:     1 - 0.05*float64(warlock.Talents.DestructiveReach),
		BaseDamage:           core.BaseDamageConfigMagic(544.0, 607.0, 0.857+0.04*float64(warlock.Talents.ShadowAndFlame)),
		OutcomeApplier:       warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Ruin, 1, 0))),
	}
	// Don't add ISB debuff aura if the target is initialized with the 'estimated ISB uptime' debuff.
	if warlock.Talents.ImprovedShadowBolt > 0 {
		existingAura := warlock.Env.Encounter.Targets[0].GetAurasWithTag("ImprovedShadowBolt")

		if len(existingAura) == 0 || existingAura[0].Duration != core.NeverExpires {
			warlock.ImpShadowboltAura = core.ImprovedShadowBoltAura(warlock.CurrentTarget, warlock.Talents.ImprovedShadowBolt, 0)
			effect.OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spellEffect.Outcome.Matches(core.OutcomeCrit) {
					return
				}
				if !warlock.ImpShadowboltAura.IsActive() {
					warlock.ImpShadowboltAura.Activate(sim)
				}
				warlock.ImpShadowboltAura.SetStacks(sim, 4)
			}
		}
	}

	var modCast func(*core.Simulation, *core.Spell, *core.Cast)

	if warlock.Talents.Nightfall > 0 {
		modCast = func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
			warlock.applyNightfall(cast)
		}
	}

	baseCost := 420.0
	warlock.Shadowbolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27209},
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.01*float64(warlock.Talents.Cataclysm)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*3000 - (time.Millisecond * 100 * time.Duration(warlock.Talents.Bane)),
			},
			ModifyCast: modCast,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

}
