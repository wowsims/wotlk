package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

// Returns true if the regular melee swing should be used, false otherwise.
func (hunter *Hunter) TryRaptorStrike(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if hunter.curQueuedAutoSpell != nil && hunter.curQueuedAutoSpell.CanCast(sim, hunter.CurrentTarget) {
		return hunter.curQueuedAutoSpell
	}
	return mhSwingSpell
}

func (hunter *Hunter) getRaptorStrikeConfig(rank int) core.SpellConfig {
	spellId := [9]int32{0, 2973, 14260, 14261, 14262, 14263, 14264, 14265, 14266}[rank]
	baseDamage := [9]float64{0, 5, 11, 21, 34, 50, 80, 110, 140}[rank]
	manaCost := [9]float64{0, 15, 25, 35, 45, 55, 70, 80, 100}[rank]
	level := [9]int{0, 1, 8, 16, 24, 32, 40, 48, 56}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolPhysical,
		ProcMask:      core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget <= 5
		},

		BonusCritRating:  float64(hunter.Talents.SavageStrikes) * 10 * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   hunter.critMultiplier(false, hunter.CurrentTarget),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := baseDamage +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if hunter.curQueueAura != nil {
				hunter.curQueueAura.Deactivate(sim)
			}
		},
	}
}

func (hunter *Hunter) makeQueueSpellsAndAura(srcSpell *core.Spell) *core.Spell {
	queueAura := hunter.RegisterAura(core.Aura{
		Label:    "RaptorStrikeQueue" + srcSpell.ActionID.String(),
		ActionID: srcSpell.ActionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.curQueueAura != nil {
				hunter.curQueueAura.Deactivate(sim)
			}
			hunter.PseudoStats.DisableDWMissPenalty = true
			hunter.curQueueAura = aura
			hunter.curQueuedAutoSpell = srcSpell
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.PseudoStats.DisableDWMissPenalty = false
			hunter.curQueueAura = nil
			hunter.curQueuedAutoSpell = nil
		},
	})

	queueSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID:    srcSpell.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.curQueueAura != queueAura &&
				hunter.CurrentMana() >= srcSpell.DefaultCast.Cost &&
				sim.CurrentTime >= hunter.Hardcast.Expires &&
				hunter.DistanceFromTarget <= 5 &&
				srcSpell.IsReady(sim)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			queueAura.Activate(sim)
		},
	})

	return queueSpell
}

func (hunter *Hunter) registerRaptorStrikeSpell() {
	maxRank := 8

	for i := 1; i <= maxRank; i++ {
		config := hunter.getRaptorStrikeConfig(i)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.RaptorStrike = hunter.GetOrRegisterSpell(config)
			hunter.makeQueueSpellsAndAura(hunter.RaptorStrike)
		}
	}
}
