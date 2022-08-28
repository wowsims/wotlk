package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerRevengeSpell(cdTimer *core.Timer) {
	warrior.RegisterAura(core.Aura{
		Label:    "Revenge Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry) {
				warrior.RevengeValidUntil = sim.CurrentTime + time.Second*5
			}
		},
	})

	cost := 5.0 - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8

	// TODO: This janky stuff is working but making an array of the enemy units does not??
	baseEffect := core.SpellEffect{}
	targets := core.TernaryInt32(warrior.Talents.ImprovedRevenge > 0, 2, 1)
	numHits := core.MinInt32(targets, warrior.Env.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = warrior.Env.GetTargetUnit(i)
	}

	applyEffect := core.ApplyEffectFuncDirectDamage(core.SpellEffect{
		ProcMask: core.ProcMaskMeleeMHSpecial,

		DamageMultiplier: 1.0 + 0.3*float64(warrior.Talents.ImprovedRevenge),
		ThreatMultiplier: 1,
		FlatThreatBonus:  200,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				roll := (1998.0-1636.0)*sim.RandomFloat("Revenge Roll") + 1636.0
				return roll + hitEffect.MeleeAttackPower(spell.Unit)*0.31
			},
			TargetSpellCoefficient: 1,
		},
		OutcomeApplier: warrior.OutcomeFuncMeleeSpecialHitAndCrit(warrior.critMultiplier(true)),

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			warrior.RevengeValidUntil = 0
			if !spellEffect.Landed() {
				warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
			}
		},
	})

	warrior.Revenge = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 30357},
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
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: time.Second * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			applyEffect(sim, target, spell)

			if target == warrior.CurrentTarget && numHits > 1 {
				if sim.RandomFloat("Revenge Target Roll") <= 0.5*float64(warrior.Talents.ImprovedRevenge) {
					applyEffect(sim, effects[1].Target, spell)
				}
			}
		},
	})
}

func (warrior *Warrior) CanRevenge(sim *core.Simulation) bool {
	return sim.CurrentTime < warrior.RevengeValidUntil &&
		warrior.StanceMatches(DefensiveStance) &&
		warrior.CurrentRage() >= warrior.Revenge.DefaultCast.Cost &&
		warrior.Revenge.IsReady(sim)
}
