package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerThunderClapSpell() {
	cost := 20.0 -
		float64(warrior.Talents.FocusedRage) -
		[]float64{0, 1, 2, 4}[warrior.Talents.ImprovedThunderClap] -
		core.TernaryFloat64(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfResonatingPower), 5, 0)

	warrior.ThunderClapAuras = make([]*core.Aura, warrior.Env.GetNumTargets())
	for _, target := range warrior.Env.Encounter.Targets {
		warrior.ThunderClapAuras[target.Index] = core.ThunderClapAura(&target.Unit, warrior.Talents.ImprovedThunderClap)
	}

	warrior.ThunderClap = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47502},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagIncludeTargetBonusDamage,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		BonusCritRating:  float64(warrior.Talents.Incite) * 5 * core.CritRatingPerCritChance,
		DamageMultiplier: []float64{1.0, 1.1, 1.2, 1.3}[warrior.Talents.ImprovedThunderClap],
		CritMultiplier:   warrior.critMultiplier(none),
		ThreatMultiplier: 1.85,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 300 + 0.12*spell.MeleeAttackPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()

			for _, aoeTarget := range sim.Encounter.Targets {
				result := spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, spell.OutcomeRangedHitAndCrit)
				if result.Landed() {
					warrior.ThunderClapAuras[aoeTarget.Index].Activate(sim)
				}
			}
		},
	})
}

func (warrior *Warrior) CanThunderClap(sim *core.Simulation) bool {
	return warrior.StanceMatches(BattleStance|DefensiveStance) && warrior.CanThunderClapIgnoreStance(sim)
}
func (warrior *Warrior) CanThunderClapIgnoreStance(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.ThunderClap.DefaultCast.Cost && warrior.ThunderClap.IsReady(sim)
}

func (warrior *Warrior) ShouldThunderClap(sim *core.Simulation, filler bool, maintainOnly bool, ignoreStance bool) bool {
	if ignoreStance && !warrior.CanThunderClapIgnoreStance(sim) {
		return false
	} else if !ignoreStance && !warrior.CanThunderClap(sim) {
		return false
	}

	if filler {
		return true
	}

	return maintainOnly &&
		warrior.ThunderClapAuras[warrior.CurrentTarget.Index].ShouldRefreshExclusiveEffects(sim, time.Second*2)
}
