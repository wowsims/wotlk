package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerRevengeSpell(cdTimer *core.Timer) {
	actionID := core.ActionID{SpellID: 30357}
	warrior.revengeProcAura = warrior.RegisterAura(core.Aura{
		Label:    "Revenge",
		Duration: 5 * time.Second,
		ActionID: actionID,
	})

	hasGlyph := warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRevenge)
	if hasGlyph {
		warrior.glyphOfRevengeProcAura = warrior.RegisterAura(core.Aura{
			Label:    "Glyph of Revenge",
			Duration: core.NeverExpires,
			ActionID: core.ActionID{SpellID: 58398},
		})
	}

	warrior.RegisterAura(core.Aura{
		Label:    "Revenge Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry) {
				warrior.revengeProcAura.Activate(sim)
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

		DamageMultiplier: 1.0 + 0.1*float64(warrior.Talents.UnrelentingAssault),
		ThreatMultiplier: 1,
		FlatThreatBonus:  121,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return core.DamageRoll(sim, 1636, 1998)*(1.0+0.3*float64(warrior.Talents.ImprovedRevenge)) + warrior.attackPowerMultiplier(hitEffect, spell.Unit, 0.31)
			},
			TargetSpellCoefficient: 1,
		},
		OutcomeApplier: warrior.OutcomeFuncMeleeSpecialHitAndCrit(warrior.critMultiplier(mh)),

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
			}
		},
	})

	warrior.Revenge = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
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
				Duration: time.Second*5 - 2*time.Second*time.Duration(warrior.Talents.UnrelentingAssault),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			applyEffect(sim, target, spell)
			warrior.revengeProcAura.Deactivate(sim)

			if warrior.glyphOfRevengeProcAura != nil {
				warrior.glyphOfRevengeProcAura.Activate(sim)
			}

			if target == warrior.CurrentTarget && numHits > 1 {
				if sim.RandomFloat("Revenge Target Roll") <= 0.5*float64(warrior.Talents.ImprovedRevenge) {
					applyEffect(sim, effects[1].Target, spell)
				}
			}
		},
	})
}

func (warrior *Warrior) CanRevenge(sim *core.Simulation) bool {
	return warrior.revengeProcAura.IsActive() &&
		warrior.StanceMatches(DefensiveStance) &&
		warrior.CurrentRage() >= warrior.Revenge.DefaultCast.Cost &&
		warrior.Revenge.IsReady(sim)
}
