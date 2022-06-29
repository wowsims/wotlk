package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warrior *Warrior) registerThunderClapSpell() {
	cost := 20.0 - float64(warrior.Talents.FocusedRage)
	impTCDamageMult := 1.0
	if warrior.Talents.ImprovedThunderClap == 1 {
		cost -= 1
		impTCDamageMult = 1.4
	} else if warrior.Talents.ImprovedThunderClap == 2 {
		cost -= 2
		impTCDamageMult = 1.7
	} else if warrior.Talents.ImprovedThunderClap == 3 {
		cost -= 4
		impTCDamageMult = 2
	}

	baseEffect := core.SpellEffect{
		ProcMask:         core.ProcMaskSpellDamage,
		DamageMultiplier: impTCDamageMult,
		ThreatMultiplier: 1.75,
		BaseDamage:       core.BaseDamageConfigFlat(123),
		OutcomeApplier:   warrior.OutcomeFuncMagicHitAndCrit(warrior.spellCritMultiplier(true)),
	}

	numHits := core.MinInt32(4, warrior.Env.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = warrior.Env.GetTargetUnit(i)

		tcAura := core.ThunderClapAura(effects[i].Target, warrior.Talents.ImprovedThunderClap)
		if i == 0 {
			warrior.ThunderClapAura = tcAura
		}

		effects[i].OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				tcAura.Activate(sim)
			}
		}
	}

	warrior.ThunderClap = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 25264},
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
				Duration: time.Second * 4,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
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
		warrior.CurrentTarget.ShouldRefreshAuraWithTagAtPriority(sim, core.ThunderClapAuraTag, warrior.ThunderClapAura.Priority, time.Second*2)
}
