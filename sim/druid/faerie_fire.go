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
	formMask := Humanoid | Moonkin

	if druid.InForm(Cat | Bear) {
		actionID = core.ActionID{SpellID: 16857}
		manaCostOptions = core.ManaCostOptions{}
		gcd = time.Second
		ignoreHaste = true
		flags = core.SpellFlagNone
		formMask = Cat | Bear
		cd = core.Cooldown{
			Timer:    druid.NewTimer(),
			Duration: time.Second * 6,
		}
		flatThreatBonus = 632.
	}
	flags |= core.SpellFlagAPL

	druid.FaerieFireAuras = druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.FaerieFireAura(target, druid.Talents.ImprovedFaerieFire)
	})

	druid.FaerieFire = druid.RegisterSpell(formMask, core.SpellConfig{
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
				druid.FaerieFireAuras.Get(target).Activate(sim)
			}
		},

		RelatedAuras: []core.AuraArray{druid.FaerieFireAuras},
	})
}

func (druid *Druid) ShouldFaerieFire(sim *core.Simulation, target *core.Unit) bool {
	if druid.FaerieFire == nil {
		return false
	}

	if !druid.FaerieFire.IsReady(sim) {
		return false
	}

	return druid.FaerieFireAuras.Get(target).ShouldRefreshExclusiveEffects(sim, time.Second*3)
}
