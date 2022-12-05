package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerHolyWrathSpell() {
	// From the perspective of max rank.
	baseCost := paladin.BaseMana * 0.20
	results := make([]*core.SpellResult, len(paladin.Env.Encounter.Targets))

	paladin.HolyWrath = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48817},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second*30 - core.TernaryDuration(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHolyWrath), time.Second*15, 0),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   paladin.SpellCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			constBaseDamage := .07*spell.SpellPower() + .07*spell.MeleeAttackPower()

			for i, aoeTarget := range sim.Encounter.Targets {
				aoeUnit := &aoeTarget.Unit
				baseDamage := constBaseDamage + sim.Roll(1050, 1234)

				if aoeUnit.MobType == proto.MobType_MobTypeDemon || aoeUnit.MobType == proto.MobType_MobTypeUndead {
					results[i] = spell.CalcDamage(sim, aoeUnit, baseDamage, spell.OutcomeMagicHitAndCrit)
				} else {
					results[i] = spell.CalcDamage(sim, aoeUnit, baseDamage, spell.OutcomeAlwaysMiss)
				}
			}

			for i := range sim.Encounter.Targets {
				spell.DealDamage(sim, results[i])
			}
		},
	})
}
