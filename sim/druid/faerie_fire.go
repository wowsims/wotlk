package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (druid *Druid) registerFaerieFireSpell() {
	actionID := core.ActionID{SpellID: 770}
	manaCostOptions := core.ManaCostOptions{
		BaseCost: 0.08,
	}
	gcd := core.GCDDefault
	ignoreHaste := false
	cd := core.Cooldown{}
	flatThreatBonus := 66. * 2.
	flags := SpellFlagOmenTrigger

	if druid.InForm(Cat | Bear) {
		actionID = core.ActionID{SpellID: 16857}
		manaCostOptions = core.ManaCostOptions{}
		gcd = time.Second
		ignoreHaste = true
		flags = core.SpellFlagNone
		cd = core.Cooldown{
			Timer:    druid.NewTimer(),
			Duration: time.Second * 6,
		}
		flatThreatBonus = 632.
	}

	druid.FaerieFireAura = core.FaerieFireAura(druid.CurrentTarget, druid.Talents.ImprovedFaerieFire)

	druid.FaerieFire = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       flags,

		ManaCost: manaCostOptions,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: gcd,
			},
			IgnoreHaste: ignoreHaste,
			CD:          cd,
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  flatThreatBonus,
		DamageMultiplier: 1,
		CritMultiplier:   druid.BalanceCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.0
			outcome := spell.OutcomeMagicHit
			if druid.InForm(Bear) {
				baseDamage = 1 + 0.15*spell.MeleeAttackPower()
				outcome = spell.OutcomeMagicHitAndCrit
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, outcome)
			if result.Landed() {
				druid.FaerieFireAura.Activate(sim)
			}
		},
	})
}

func (druid *Druid) ShouldFaerieFire(sim *core.Simulation) bool {
	if druid.FaerieFire == nil {
		return false
	}

	if !druid.FaerieFire.IsReady(sim) {
		return false
	}

	return druid.FaerieFireAura.ShouldRefreshExclusiveEffects(sim, time.Second*3)
}
