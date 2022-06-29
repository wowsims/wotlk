package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warrior *Warrior) registerDemoralizingShoutSpell() {
	cost := 10.0
	cost -= float64(warrior.Talents.FocusedRage)
	if ItemSetBoldArmor.CharacterHasSetBonus(&warrior.Character, 2) {
		cost -= 2
	}

	baseEffect := core.SpellEffect{
		ProcMask:         core.ProcMaskEmpty,
		ThreatMultiplier: 1,
		FlatThreatBonus:  56,
		OutcomeApplier:   warrior.OutcomeFuncMagicHit(),
	}

	numHits := warrior.Env.GetNumTargets()
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = warrior.Env.GetTargetUnit(i)

		demoShoutAura := core.DemoralizingShoutAura(effects[i].Target, warrior.Talents.BoomingVoice, warrior.Talents.ImprovedDemoralizingShout)
		if i == 0 {
			warrior.DemoralizingShoutAura = demoShoutAura
		}

		effects[i].OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				demoShoutAura.Activate(sim)
			}
		}
	}

	warrior.DemoralizingShout = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 25203},
		SpellSchool: core.SpellSchoolPhysical,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}

func (warrior *Warrior) CanDemoralizingShout(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.DemoralizingShout.DefaultCast.Cost
}

func (warrior *Warrior) ShouldDemoralizingShout(sim *core.Simulation, filler bool, maintainOnly bool) bool {
	if !warrior.CanDemoralizingShout(sim) {
		return false
	}

	if filler {
		return true
	}

	return maintainOnly &&
		warrior.CurrentTarget.ShouldRefreshAuraWithTagAtPriority(sim, core.APReductionAuraTag, warrior.DemoralizingShoutAura.Priority, time.Second*2)
}
