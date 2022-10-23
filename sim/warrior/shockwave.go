package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerShockwaveSpell() {
	cost := 15.0 - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8
	cd := 20*time.Second - core.TernaryDuration(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfShockwave), 3*time.Second, 0)

	warrior.Shockwave = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 46968},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRanged, // TODO: Is this correct?
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
				Timer:    warrior.NewTimer(),
				Duration: cd,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(none),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.75 * spell.MeleeAttackPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.Targets {
				result := spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				// TODO: Should this be main target only?
				if !result.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			}
		},
	})
}

func (warrior *Warrior) CanShockwave(sim *core.Simulation) bool {
	return warrior.StanceMatches(DefensiveStance) &&
		warrior.CurrentRage() >= warrior.Shockwave.DefaultCast.Cost &&
		warrior.Shockwave.IsReady(sim)
}
