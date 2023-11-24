package warrior

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (warrior *Warrior) registerThunderClapSpell() {
	warrior.ThunderClapAuras = warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.ThunderClapAura(target, warrior.Talents.ImprovedThunderClap)
	})

	baseDamage := map[int32]float64{
		25: 23,
		40: 55,
		50: 82,
		60: 103,
	}[warrior.Level]

	numHits := min(4, warrior.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	warrior.ThunderClap = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47502},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 20 - []float64{0, 1, 2, 4}[warrior.Talents.ImprovedThunderClap],
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
		BonusCritRating:  (0 - float64(warrior.Talents.Cruelty)*1),
		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(none),
		ThreatMultiplier: 1.85,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				if results[hitIndex].Landed() {
					warrior.ThunderClapAuras.Get(target).Activate(sim)
				}
				curTarget = sim.Environment.NextTargetUnit(curTarget)
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
