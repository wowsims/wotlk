package warrior

import (
	"math"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerExecuteSpell() {
	cost := 15.0 - float64(warrior.Talents.FocusedRage)
	if warrior.Talents.ImprovedExecute == 1 {
		cost -= 2
	} else if warrior.Talents.ImprovedExecute == 2 {
		cost -= 5
	}
	if warrior.HasSetBonus(ItemSetOnslaughtBattlegear, 2) {
		cost -= 3
	}
	refundAmount := cost * 0.8

	var extraRage float64
	extraRageBonus := core.TernaryFloat64(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfExecution), 10, 0)

	warrior.Execute = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47471},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.Cost = math.Min(spell.Unit.CurrentRage(), 30)
				extraRage = cast.Cost - spell.BaseCost
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: 1,
			ThreatMultiplier: 1.25,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return 1456 + hitEffect.MeleeAttackPower(spell.Unit)*0.2 + 38*(extraRage+extraRageBonus)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: warrior.OutcomeFuncMeleeSpecialHitAndCrit(warrior.critMultiplier(true)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) SpamExecute(spam bool) bool {
	return warrior.CurrentRage() >= warrior.Execute.BaseCost && spam && warrior.Talents.MortalStrike
}

func (warrior *Warrior) CanExecute() bool {
	return warrior.CurrentRage() >= warrior.Execute.BaseCost
}
