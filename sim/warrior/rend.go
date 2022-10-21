package warrior

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) RegisterRendSpell(rageThreshold float64, healthThreshold float64) {
	actionID := core.ActionID{SpellID: 47465}

	cost := 10.0
	refundAmount := cost * 0.8
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

		// 135% damage multiplier is applied at the beginning of the fight and removed when target is at 75% health
		DamageMultiplier: 1 + 0.1*float64(warrior.Talents.ImprovedRend),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier: warrior.OutcomeFuncMeleeSpecialHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					warrior.RendDots.Apply(sim)
					warrior.procBloodFrenzy(sim, spellEffect, time.Second*core.TernaryDuration(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRending), 21, 15))
					warrior.rendValidUntil = sim.CurrentTime + time.Second*core.TernaryDuration(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRending), 21, 15)
				} else {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
	target := warrior.CurrentTarget
	//snapshotCalculator := func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
	//	tickDamage := (380 + 0.2*5*warrior.AutoAttacks.MH.CalculateAverageWeaponDamage(spell.MeleeAttackPower())) / 5
	//	if sim.GetRemainingDurationPercent() > 0.75 {
	//		return tickDamage * 1.35
	//	}
	//	return tickDamage
	//}
	warrior.RendDots = core.NewDot(core.Dot{
		Spell: warrior.Rend,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Rends-" + strconv.Itoa(int(warrior.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: core.TernaryInt(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRending), 7, 5),
		TickLength:    time.Second * 3,
		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = (380 + 0.2*5*warrior.AutoAttacks.MH.CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower())) / 5
			if sim.GetRemainingDurationPercent() > 0.75 {
				dot.SnapshotBaseDamage *= 1.35
			}
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			if sim.Log != nil {
				sim.Log("Snapshot rend multiplier: %0.04f", dot.SnapshotAttackerMultiplier)
			}
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
		},
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
