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
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry) {
				warrior.revengeProcAura.Activate(sim)
			}
		},
	})

	cost := 5.0 - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8

	extraHit := warrior.Talents.ImprovedRevenge > 0 && warrior.Env.GetNumTargets() > 1
	rollMultiplier := 1 + 0.3*float64(warrior.Talents.ImprovedRevenge)

	warrior.Revenge = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

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

		DamageMultiplier: 1.0 + 0.1*float64(warrior.Talents.UnrelentingAssault),
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,
		FlatThreatBonus:  121,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmgFromAP := 0.31 * spell.MeleeAttackPower()
			baseDamage := sim.Roll(1636, 1998)*rollMultiplier + dmgFromAP
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			if !result.Landed() {
				warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
			}

			if extraHit {
				if sim.RandomFloat("Revenge Target Roll") <= 0.5*float64(warrior.Talents.ImprovedRevenge) {
					otherTarget := sim.Environment.NextTargetUnit(target)
					baseDamage := sim.Roll(1636, 1998)*rollMultiplier + dmgFromAP
					spell.CalcAndDealDamage(sim, otherTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				}
			}

			warrior.revengeProcAura.Deactivate(sim)

			if warrior.glyphOfRevengeProcAura != nil {
				warrior.glyphOfRevengeProcAura.Activate(sim)
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
