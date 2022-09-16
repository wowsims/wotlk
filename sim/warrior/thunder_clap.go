package warrior

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerThunderClapSpell() {
	cost := 20.0 - float64(warrior.Talents.FocusedRage) - []float64{0, 1, 2, 4}[warrior.Talents.ImprovedThunderClap]
	cost -= core.TernaryFloat64(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfResonatingPower), 5, 0)
	impTCDamageMult := []float64{1.0, 1.1, 1.2, 1.3}[warrior.Talents.ImprovedThunderClap]

	warrior.ThunderClapAura = core.ThunderClapAura(warrior.CurrentTarget, warrior.Talents.ImprovedThunderClap)

	baseEffect := core.SpellEffect{
		ProcMask:         core.ProcMaskRangedSpecial,
		DamageMultiplier: impTCDamageMult,
		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return warrior.attackPowerMultiplier(hitEffect, spell.Unit, 0.12) + 300
			},
			TargetSpellCoefficient: 1,
		},
		OutcomeApplier: warrior.OutcomeFuncRangedHitAndCrit(warrior.critMultiplier(none)),
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				core.ThunderClapAura(spellEffect.Target, warrior.Talents.ImprovedThunderClap).Activate(sim)
			}
		},
	}

	warrior.ThunderClap = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47502},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagBinary,

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
		ThreatMultiplier: 1.85,

		ApplyEffects: core.ApplyEffectFuncAOEDamageCapped(warrior.Env, baseEffect),
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
		warrior.CurrentTarget.ShouldRefreshAuraWithTagAtPriority(sim, core.AtkSpeedReductionAuraTag, warrior.ThunderClapAura.Priority, time.Second*2)
}
