package warrior

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO (maybe) https://github.com/magey/wotlk-warrior/issues/23 - Rend is not benefitting from Two-Handed Weapon Specialization
func (warrior *Warrior) RegisterRendSpell(rageThreshold float64, healthThreshold float64) {
	actionID := core.ActionID{SpellID: 47465}

	cost := 10.0
	refundAmount := cost * 0.8

	dotDuration := time.Second * 15
	dotTicks := 5
	if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRending) {
		dotDuration += time.Second * 6
		dotTicks += 2
	}

	warrior.Rend = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagNoOnCastComplete,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 + 0.1*float64(warrior.Talents.ImprovedRend),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier: warrior.OutcomeFuncMeleeSpecialHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					warrior.RendDots.Apply(sim)
					warrior.procBloodFrenzy(sim, spellEffect, dotDuration)
					warrior.rendValidUntil = sim.CurrentTime + dotDuration
				} else {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
	target := warrior.CurrentTarget
	warrior.RendDots = core.NewDot(core.Dot{
		Spell: warrior.Rend,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Rends-" + strconv.Itoa(int(warrior.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: dotTicks,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, _ *core.SpellEffect, spell *core.Spell) float64 {
					tickDamage := (380 + warrior.AutoAttacks.MH.CalculateAverageWeaponDamage(spell.MeleeAttackPower())) / 5
					// 135% damage multiplier is applied at the beginning of the fight and removed when target is at 75% health
					if sim.GetRemainingDurationPercent() > 0.75 {
						return tickDamage * 1.35
					}
					return tickDamage
				},
			},
			OutcomeApplier: warrior.OutcomeFuncTick(),
			IsPeriodic:     true,
		}),
	})

	warrior.RendRageThresholdBelow = core.MaxFloat(warrior.Rend.DefaultCast.Cost, rageThreshold)
	warrior.RendHealthThresholdAbove = healthThreshold / 100
}

func (warrior *Warrior) ShouldRend(sim *core.Simulation) bool {
	if warrior.Talents.Bloodthirst {
		return warrior.Rend.IsReady(sim) && sim.CurrentTime >= (warrior.rendValidUntil-warrior.RendCdThreshold) && !warrior.Whirlwind.IsReady(sim) &&
			warrior.CurrentRage() <= warrior.RendRageThresholdBelow && warrior.RendHealthThresholdAbove < sim.GetRemainingDurationPercent()
	}
	return warrior.Rend.IsReady(sim) && sim.CurrentTime >= (warrior.rendValidUntil-warrior.RendCdThreshold) && warrior.CurrentRage() >= warrior.Rend.DefaultCast.Cost
}
