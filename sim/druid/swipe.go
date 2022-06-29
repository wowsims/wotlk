package druid

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (druid *Druid) registerSwipeSpell() {
	cost := 20.0 - float64(druid.Talents.Ferocity)

	baseDamage := 84.0
	if druid.Equip[items.ItemSlotRanged].ID == 23198 { // Idol of Brutality
		baseDamage += 10
	}

	baseEffect := core.SpellEffect{
		ProcMask: core.ProcMaskMeleeMHSpecial,

		DamageMultiplier: 1 + core.TernaryFloat64(druid.InForm(Bear) && ItemSetThunderheartHarness.CharacterHasSetBonus(&druid.Character, 4), 0.15, 0),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return baseDamage + 0.07*hitEffect.MeleeAttackPower(spell.Unit)
			},
			TargetSpellCoefficient: 1,
		},
		OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(druid.MeleeCritMultiplier()),
	}

	numHits := core.MinInt32(3, druid.Env.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effect := baseEffect
		effect.Target = druid.Env.GetTargetUnit(i)
		effects = append(effects, effect)
	}

	druid.Swipe = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26997},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			ModifyCast:  druid.ApplyClearcasting,
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}

func (druid *Druid) CanSwipe() bool {
	return druid.CurrentRage() >= druid.Swipe.DefaultCast.Cost
}
