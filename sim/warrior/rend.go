package warrior

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerRendSpell(rageThreshold float64) {
	actionID := core.ActionID{SpellID: 47465}

	cost := 10.0
	refundAmount := cost * 0.8
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
					warrior.RendDots.Apply(sim)
					warrior.procBloodFrenzy(sim, spellEffect, time.Second*15)
					warrior.rendValidUntil = sim.CurrentTime + time.Second*15
				} else {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
	target := warrior.CurrentTarget
	tickDamage := 380 + 0.2*5*warrior.AutoAttacks.MH.AverageDamage()/15
	warrior.RendDots = core.NewDot(core.Dot{
		Spell: warrior.Rend,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Rends-" + strconv.Itoa(int(warrior.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: core.TernaryInt(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRending), 7, 5),
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 + 0.01*float64(warrior.Talents.ImprovedRend),
			ThreatMultiplier: 1,
			IsPeriodic:       true,

			BaseDamage:     core.BaseDamageConfigFlat(tickDamage),
			OutcomeApplier: warrior.OutcomeFuncTick(),
		}),
	})
	warrior.RendRageThreshold = core.MaxFloat(warrior.Rend.DefaultCast.Cost, rageThreshold)
}

func (warrior *Warrior) ShouldRend(sim *core.Simulation) bool {
	if !warrior.Rend.IsReady(sim) {
		return false
	}

	if warrior.Talents.MortalStrike {
		return sim.CurrentTime >= (warrior.rendValidUntil-warrior.RendCdThreshold) && warrior.CurrentRage() >= warrior.Rend.DefaultCast.Cost
	} else if warrior.Talents.Bloodthirst {
		return warrior.CurrentRage() >= warrior.RendRageThreshold
	}
	return false
}
