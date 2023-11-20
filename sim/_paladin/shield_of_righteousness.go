package paladin

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/stats"
)

func (paladin *Paladin) registerShieldOfRighteousnessSpell() {
	var aegisPlateProcAura *core.Aura
	if paladin.HasSetBonus(ItemSetAegisPlate, 4) {
		aegisPlateProcAura = paladin.NewTemporaryStatsAura("Aegis", core.ActionID{SpellID: 64883}, stats.Stats{stats.BlockValue: 225}, time.Second*6)
	}

	paladin.ShieldOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 61411},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.06,
			Multiplier: 1 - 0.02*float64(paladin.Talents.Benediction),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if aegisPlateProcAura != nil {
				aegisPlateProcAura.Activate(sim)
			}

			var baseDamage float64
			// TODO: Derive or find accurate source for DR curve
			bv := paladin.BlockValue()
			if bv <= 2400.0 {
				baseDamage = 520.0 + bv
			} else {
				bv = 2400.0 + (bv-2400.0)/2
				baseDamage = 520.0 + core.TernaryFloat64(bv > 2760.0, 2760.0, bv)
			}

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}
