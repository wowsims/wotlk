package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warrior *Warrior) registerOverpowerSpell(cdTimer *core.Timer) {
	outcomeMask := core.OutcomeDodge
	if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfOverpower) {
		outcomeMask |= core.OutcomeParry
	}
	warrior.RegisterAura(core.Aura{
		Label:    "Overpower Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(outcomeMask) {
				warrior.OverpowerAura.Activate(sim)
				warrior.lastOverpowerProc = sim.CurrentTime
			}
		},
	})

	warrior.OverpowerAura = warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 68051},
		Label:    "Overpower Aura",
		Duration: time.Second * 5,
	})

	cooldownDur := time.Second * 5
	gcdDur := core.GCDDefault

	if warrior.Talents.UnrelentingAssault == 1 {
		cooldownDur -= time.Second * 2
	} else if warrior.Talents.UnrelentingAssault == 2 {
		cooldownDur -= time.Second * 4
		gcdDur -= time.Millisecond * 500
	}
	warrior.Overpower = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 7384},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

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

		BonusCritRating:  25 * core.CritRatingPerCritChance * float64(warrior.Talents.ImprovedOverpower),
		DamageMultiplier: 1 + 0.1*float64(warrior.Talents.UnrelentingAssault),
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 0.75,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warrior.OverpowerAura.Deactivate(sim)

			baseDamage := 0 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (warrior *Warrior) ShouldOverpower(sim *core.Simulation) bool {
	return warrior.OverpowerAura.IsActive() && warrior.Overpower.IsReady(sim) &&
		warrior.CurrentRage() >= warrior.Overpower.DefaultCast.Cost &&
		sim.CurrentTime > (warrior.lastOverpowerProc+warrior.reactionTime)
}

// Queue Overpower to be cast at every 6s if talented for 3/3 TfB
func (warrior *Warrior) CastFullTfbOverpower(sim *core.Simulation, target *core.Unit) bool {
	if warrior.Talents.TasteForBlood < 3 {
		return false
	}

	core.StartDelayedAction(sim, core.DelayedActionOptions{
		DoAt: sim.CurrentTime + time.Second*6,
		OnAction: func(_ *core.Simulation) {
			if warrior.Overpower.CanCast(sim, target) && warrior.ShouldOverpower(sim) {
				warrior.CastFullTfbOverpower(sim, target)
			}
		},
	})

	return warrior.Overpower.Cast(sim, target)
}
