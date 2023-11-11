package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warrior *Warrior) registerRevengeSpell(cdTimer *core.Timer) {
	actionID := core.ActionID{SpellID: 57823}

	warrior.revengeProcAura = warrior.RegisterAura(core.Aura{
		Label:    "Revenge",
		Duration: 5 * time.Second,
		ActionID: actionID,
	})

	var glyphOfRevengeProcAura *core.Aura
	if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRevenge) {
		glyphOfRevengeProcAura = warrior.RegisterAura(core.Aura{
			Label:    "Glyph of Revenge",
			Duration: core.NeverExpires,
			ActionID: core.ActionID{SpellID: 58398},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warrior.HeroicStrike.CostMultiplier -= 1
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warrior.HeroicStrike.CostMultiplier += 1
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell == warrior.HeroicStrike {
					aura.Deactivate(sim)
				}
			},
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

	cooldownDur := time.Second * 5
	gcdDur := core.GCDDefault

	if warrior.Talents.UnrelentingAssault == 1 {
		cooldownDur -= time.Second * 2
	} else if warrior.Talents.UnrelentingAssault == 2 {
		cooldownDur -= time.Second * 4
		gcdDur -= time.Millisecond * 500
	}

	extraHit := warrior.Talents.ImprovedRevenge > 0 && warrior.Env.GetNumTargets() > 1

	warrior.Revenge = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   5 - float64(warrior.Talents.FocusedRage),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: gcdDur,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cooldownDur,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(DefensiveStance) && warrior.revengeProcAura.IsActive()
		},

		DamageMultiplier: 1.0 + 0.1*float64(warrior.Talents.UnrelentingAssault) + 0.3*float64(warrior.Talents.ImprovedRevenge),
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,
		FlatThreatBonus:  121,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(1636, 1998) + 0.31*spell.MeleeAttackPower()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			if extraHit {
				if sim.RandomFloat("Revenge Target Roll") <= 0.5*float64(warrior.Talents.ImprovedRevenge) {
					otherTarget := sim.Environment.NextTargetUnit(target)
					baseDamage := sim.Roll(1636, 1998) + 0.31*spell.MeleeAttackPower()
					spell.CalcAndDealDamage(sim, otherTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				}
			}

			warrior.revengeProcAura.Deactivate(sim)

			if glyphOfRevengeProcAura != nil {
				glyphOfRevengeProcAura.Activate(sim)
			}
		},
	})
}
