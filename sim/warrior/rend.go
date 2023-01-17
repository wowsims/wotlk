package warrior

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// TODO (maybe) https://github.com/magey/wotlk-warrior/issues/23 - Rend is not benefitting from Two-Handed Weapon Specialization
func (warrior *Warrior) RegisterRendSpell(rageThreshold float64, healthThreshold float64) {
	actionID := core.ActionID{SpellID: 47465}

	dotDuration := time.Second * 15
	dotTicks := int32(5)
	if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRending) {
		dotDuration += time.Second * 6
		dotTicks += 2
	}

	warrior.Rend = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagNoOnCastComplete,

		RageCost: core.RageCostOptions{
			Cost:   10 - float64(warrior.Talents.FocusedRage),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 + 0.1*float64(warrior.Talents.ImprovedRend),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				warrior.RendDots.Apply(sim)
				warrior.procBloodFrenzy(sim, result, dotDuration)
				warrior.rendValidUntil = sim.CurrentTime + dotDuration
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})

	warrior.RendDots = core.NewDot(core.Dot{
		Spell: warrior.Rend,
		Aura: warrior.CurrentTarget.RegisterAura(core.Aura{
			Label:    "Rends-" + strconv.Itoa(int(warrior.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: dotTicks,
		TickLength:    time.Second * 3,
		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = (380 + warrior.AutoAttacks.MH.CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower())) / 5
			// 135% damage multiplier is applied at the beginning of the fight and removed when target is at 75% health
			if sim.GetRemainingDurationPercent() > 0.75 {
				dot.SnapshotBaseDamage *= 1.35
			}
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
		},
	})

	warrior.RendHealthThresholdAbove = healthThreshold / 100
	warrior.RendRageThresholdBelow = core.MaxFloat(warrior.Rend.DefaultCast.Cost, rageThreshold)
}

func (warrior *Warrior) ShouldRend(sim *core.Simulation) bool {
	if warrior.PrimaryTalentTree == FuryTree {
		return warrior.Rend.IsReady(sim) && sim.CurrentTime >= (warrior.rendValidUntil-warrior.RendCdThreshold) && !warrior.Whirlwind.IsReady(sim) &&
			warrior.CurrentRage() <= warrior.RendRageThresholdBelow && warrior.RendHealthThresholdAbove < sim.GetRemainingDurationPercent() &&
			warrior.CurrentRage() >= warrior.Rend.DefaultCast.Cost
	}
	return warrior.Rend.IsReady(sim) && sim.CurrentTime >= (warrior.rendValidUntil-warrior.RendCdThreshold) && warrior.CurrentRage() >= warrior.Rend.DefaultCast.Cost
}
