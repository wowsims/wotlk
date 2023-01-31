package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warrior *Warrior) registerExecuteSpell() {
	const maxRage = 30

	gcd := core.GCDDefault
	if warrior.HasSetBonus(ItemSetYmirjarLordsBattlegear, 4) {
		gcd = time.Second
	}

	var extraRageBonus float64
	if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfExecution) {
		extraRageBonus = 10
	}

	var rageMetrics *core.ResourceMetrics
	warrior.Execute = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47471},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RageCost: core.RageCostOptions{
			Cost: 15 -
				float64(warrior.Talents.FocusedRage) -
				[]float64{0, 2, 5}[warrior.Talents.ImprovedExecute] -
				core.TernaryFloat64(warrior.HasSetBonus(ItemSetOnslaughtBattlegear, 2), 3, 0),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: gcd,
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
			warrior.SpendRage(sim, extraRage, rageMetrics)
			rageMetrics.Events--

			baseDamage := 1456 + 0.2*spell.MeleeAttackPower() + 38*(extraRage+extraRageBonus)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
	rageMetrics = warrior.Execute.Cost.(*core.RageCost).ResourceMetrics
}

func (warrior *Warrior) SpamExecute(spam bool) bool {
	return warrior.CurrentRage() >= warrior.Execute.BaseCost && spam && warrior.PrimaryTalentTree == ArmsTree
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
