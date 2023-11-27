package priest

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
)

func (priest *Priest) registerPowerWordShieldSpell() {
	coeff := 0.8057

	wsDuration := time.Second * 15

	cd := core.Cooldown{}

	var glyphHeal *core.Spell

	priest.PowerWordShield = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48066},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.23,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: cd,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !priest.WeakenedSouls.Get(target).IsActive()
		},

		DamageMultiplier: 1 *
			(1 + .05*float64(priest.Talents.ImprovedPowerWordShield)) *
			(1 + .02*float64(priest.Talents.SpiritualHealing)),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		Shield: core.ShieldConfig{
			Aura: core.Aura{
				Label:    "Power Word Shield",
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shieldAmount := 2230.0 + coeff*spell.HealingPower(target)
			shield := spell.Shield(target)
			shield.Apply(sim, shieldAmount)

			weakenedSoul := priest.WeakenedSouls.Get(target)
			weakenedSoul.Duration = wsDuration
			weakenedSoul.Activate(sim)

			if glyphHeal != nil {
				glyphHeal.Cast(sim, target)
			}
		},
	})

	priest.WeakenedSouls = priest.NewAllyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Weakened Soul",
			ActionID: core.ActionID{SpellID: 6788},
			Duration: time.Second * 15,
		})
	})
}
