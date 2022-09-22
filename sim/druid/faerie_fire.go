package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerFaerieFireSpell() {
	actionID := core.ActionID{SpellID: 770}
	resourceType := stats.Mana
	baseCost := 0.08 * druid.BaseMana
	gcd := core.GCDDefault
	ignoreHaste := false
	cd := core.Cooldown{}
	baseDamage := core.BaseDamageConfigMelee(0, 0, 0)

	if druid.InForm(Cat | Bear) {
		actionID = core.ActionID{SpellID: 16857}
		resourceType = 0
		baseCost = 0
		gcd = time.Second
		ignoreHaste = true
		cd = core.Cooldown{
			Timer:    druid.NewTimer(),
			Duration: time.Second * 6,
		}
	}

	if druid.InForm(Bear) {
		baseDamage = core.BaseDamageConfigMelee(1, 1, 0.15)
	}

	druid.FaerieFireAura = core.FaerieFireAura(druid.CurrentTarget, druid.Talents.ImprovedFaerieFire > 0)

	if druid.Talents.ImprovedFaerieFire > 0 {
		bonusCrit := core.CritRatingPerCritChance * float64(druid.Talents.ImprovedFaerieFire)
		oldOnGain := druid.FaerieFireAura.OnGain
		oldOnExpire := druid.FaerieFireAura.OnExpire
		druid.FaerieFireAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
			oldOnGain(aura, sim)
			druid.Moonfire.BonusCritRating += bonusCrit
			druid.Starfall.BonusCritRating += bonusCrit
			druid.StarfallSplash.BonusCritRating += bonusCrit
			druid.Starfire.BonusCritRating += bonusCrit
			druid.Wrath.BonusCritRating += bonusCrit
		}
		druid.FaerieFireAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
			oldOnExpire(aura, sim)
			druid.Moonfire.BonusCritRating -= bonusCrit
			druid.Starfall.BonusCritRating -= bonusCrit
			druid.StarfallSplash.BonusCritRating -= bonusCrit
			druid.Starfire.BonusCritRating -= bonusCrit
			druid.Wrath.BonusCritRating -= bonusCrit
		}
	}

	druid.FaerieFire = druid.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: resourceType,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  gcd,
			},
			IgnoreHaste: ignoreHaste,
			CD:          cd,
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  66 * 2,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage:     baseDamage,
			OutcomeApplier: druid.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.FaerieFireAura.Activate(sim)
				}
			},
		}),
	})
}

func (druid *Druid) ShouldFaerieFire(sim *core.Simulation) bool {
	if druid.FaerieFire == nil {
		return false
	}

	if !druid.FaerieFire.IsReady(sim) {
		return false
	}

	return druid.CurrentTarget.ShouldRefreshAuraWithTagAtPriority(sim, core.MinorArmorReductionAuraTag, druid.FaerieFireAura.Priority, time.Second*3)
}
