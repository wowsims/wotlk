package warrior

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerRendSpell() {
	actionID := core.ActionID{SpellID: 47465}

	cost := 10.0
	refundAmount := cost * 0.8
	isAbove75 := true
	warrior.Rend = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
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
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			ThreatMultiplier: 1,
			OutcomeApplier:   warrior.OutcomeFuncMeleeSpecialHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					if sim.GetRemainingDurationPercent() <= 0.75 && isAbove75 {
						isAbove75 = false
						warrior.RendDots.Spell.DamageMultiplier /= 1.35
					}
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
	tickDamage := (380 + 0.2*5*warrior.AutoAttacks.MH.AverageDamage()*warrior.PseudoStats.PhysicalDamageDealtMultiplier) / 5
	warrior.RendDots = core.NewDot(core.Dot{
		Spell: warrior.Rend,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Rends-" + strconv.Itoa(int(warrior.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: core.TernaryInt(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRending), 7, 5),
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncApplyEffectsToUnit(target, core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,
			// 135% damage multiplier is applied at the begining of the fight and removed when target is at 75% health
			DamageMultiplier: (1 + 0.1*float64(warrior.Talents.ImprovedRend)) * 1.35,
			ThreatMultiplier: 1,
			IsPeriodic:       true,

			BaseDamage:     core.BaseDamageConfigFlat(tickDamage),
			OutcomeApplier: warrior.OutcomeFuncTick(),
		})),
	})
}

func (warrior *Warrior) ShouldRend(sim *core.Simulation) bool {
	return warrior.Rend.IsReady(sim) && sim.CurrentTime >= (warrior.rendValidUntil-warrior.RendCdThreshold) && warrior.CurrentRage() >= warrior.Rend.DefaultCast.Cost
}
