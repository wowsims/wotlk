package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerExecuteSpell() {
	const maxRage = 30

	cost := 15 - float64(warrior.Talents.FocusedRage) - []float64{0, 2, 5}[warrior.Talents.ImprovedExecute]
	if warrior.HasSetBonus(ItemSetOnslaughtBattlegear, 2) {
		cost -= 3
	}

	refundAmount := 0.8 * cost

	gcd := core.GCDDefault
	if warrior.HasSetBonus(ItemSetYmirjarLordsBattlegear, 4) {
		gcd = time.Second
	}

	var extraRageBonus float64
	if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfExecution) {
		extraRageBonus = 10
	}

	warrior.Execute = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47471},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  gcd,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1.25,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			extraRage := spell.Unit.CurrentRage()
			if extraRage > maxRage-spell.CurCast.Cost {
				extraRage = maxRage - spell.CurCast.Cost
			}
			warrior.SpendRage(sim, extraRage, spell.ResourceMetrics)
			spell.ResourceMetrics.Events--

			baseDamage := 1456 + 0.2*spell.MeleeAttackPower() + 38*(extraRage+extraRageBonus)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
			}
		},
	})
}

func (warrior *Warrior) SpamExecute(spam bool) bool {
	return warrior.CurrentRage() >= warrior.Execute.BaseCost && spam && warrior.Talents.MortalStrike
}

func (warrior *Warrior) CanExecute() bool {
	return warrior.CurrentRage() >= warrior.Execute.BaseCost
}

func (warrior *Warrior) CanSuddenDeathExecute() bool {
	return warrior.CurrentRage() >= warrior.Execute.BaseCost && warrior.isSuddenDeathActive()
}

func (warrior *Warrior) CastExecute(sim *core.Simulation, target *core.Unit) bool {
	return warrior.Execute.Cast(sim, target)
}
