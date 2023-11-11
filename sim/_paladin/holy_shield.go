package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerHolyShieldSpell() {
	actionID := core.ActionID{SpellID: 48952}
	numCharges := int32(8)

	procSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID.WithTag(1),
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Beta testing shows wowhead coeffs are probably correct
			baseDamage := 274 +
				0.0732*spell.MeleeAttackPower() +
				0.117*spell.SpellPower()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
		},
	})

	blockBonus := 30*core.BlockRatingPerBlockChance + core.TernaryFloat64(paladin.Ranged().ID == 29388, 42, 0)

	paladin.HolyShieldAura = paladin.RegisterAura(core.Aura{
		Label:     "Holy Shield",
		ActionID:  actionID,
		Duration:  time.Second * 10,
		MaxStacks: numCharges,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.Block, blockBonus)
			aura.SetStacks(sim, numCharges)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.Block, -blockBonus)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeBlock) {
				procSpell.Cast(sim, spell.Unit)
				aura.RemoveStack(sim)
			}
		},
	})

	paladin.HolyShield = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if paladin.HolyShieldAura.IsActive() {
				paladin.HolyShieldAura.SetStacks(sim, numCharges)
			}
			paladin.HolyShieldAura.Activate(sim)
		},
	})
}
