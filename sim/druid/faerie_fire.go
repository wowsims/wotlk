package druid

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (druid *Druid) registerFaerieFireSpell() {
	actionID := core.ActionID{SpellID: 26993}
	resourceType := stats.Mana
	baseCost := 145.0
	gcd := core.GCDDefault
	ignoreHaste := false
	cd := core.Cooldown{}

	if druid.InForm(Cat | Bear) {
		if !druid.Talents.FaerieFire {
			return
		}
		actionID = core.ActionID{SpellID: 27011}
		resourceType = 0
		baseCost = 0
		gcd = time.Second
		ignoreHaste = true
		cd = core.Cooldown{
			Timer:    druid.NewTimer(),
			Duration: time.Second * 6,
		}
	}

	druid.FaerieFireAura = core.FaerieFireAura(druid.CurrentTarget, druid.Talents.ImprovedFaerieFire)

	druid.FaerieFire = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,

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

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			ThreatMultiplier: 1,
			FlatThreatBonus:  66 * 2,
			OutcomeApplier:   druid.OutcomeFuncMagicHit(),
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

	return druid.CurrentTarget.ShouldRefreshAuraWithTagAtPriority(sim, core.FaerieFireAuraTag, druid.FaerieFireAura.Priority, time.Second*3)
}
