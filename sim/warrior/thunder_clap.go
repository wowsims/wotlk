package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warrior *Warrior) registerThunderClapSpell() {
	warrior.ThunderClapAuras = warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.ThunderClapAura(target, warrior.Talents.ImprovedThunderClap)
	})

	warrior.ThunderClap = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47502},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 20 -
				float64(warrior.Talents.FocusedRage) -
				[]float64{0, 1, 2, 4}[warrior.Talents.ImprovedThunderClap] -
				core.TernaryFloat64(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfResonatingPower), 5, 0),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance | DefensiveStance)
		},

		// Cruelty doesn't apply to Thunder Clap
		BonusCritRating:  (float64(warrior.Talents.Incite)*5 - float64(warrior.Talents.Cruelty)*1) * core.CritRatingPerCritChance,
		DamageMultiplier: []float64{1.0, 1.1, 1.2, 1.3}[warrior.Talents.ImprovedThunderClap],
		CritMultiplier:   warrior.critMultiplier(none),
		ThreatMultiplier: 1.85,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 300 + 0.12*spell.MeleeAttackPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeRangedHitAndCrit)
				if result.Landed() {
					warrior.ThunderClapAuras.Get(aoeTarget).Activate(sim)
				}
			}
		},

		RelatedAuras: []core.AuraArray{warrior.ThunderClapAuras},
	})
}

func (warrior *Warrior) CanThunderClapIgnoreStance(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.ThunderClap.DefaultCast.Cost && warrior.ThunderClap.IsReady(sim)
}

func (warrior *Warrior) ShouldThunderClap(sim *core.Simulation, target *core.Unit, filler bool, maintainOnly bool, ignoreStance bool) bool {
	if ignoreStance && !warrior.CanThunderClapIgnoreStance(sim) {
		return false
	} else if !ignoreStance && !warrior.ThunderClap.CanCast(sim, target) {
		return false
	}

	if filler {
		return true
	}

	return maintainOnly &&
		warrior.ThunderClapAuras.Get(target).ShouldRefreshExclusiveEffects(sim, time.Second*2)
}
