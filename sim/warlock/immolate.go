package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) registerImmolateSpell() {
	fireAndBrimstoneBonus := 0.02 * float64(warlock.Talents.FireAndBrimstone)
	bonusPeriodicDamageMultiplier := 0 +
		0.03*float64(warlock.Talents.Aftermath) +
		core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImmolate), 0.1, 0) +
		warlock.GrandSpellstoneBonus() -
		warlock.GrandFirestoneBonus()

	warlock.Immolate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47811},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.17,
			Multiplier: 1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm],
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (2000 - 100*time.Duration(warlock.Talents.Bane)),
			},
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.Emberstorm) +
			0.1*float64(warlock.Talents.ImprovedImmolate) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDeathbringerGarb, 2), 0.1, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetGuldansRegalia, 4), 0.1, 0),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Immolate",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.ChaosBolt.DamageMultiplierAdditive += fireAndBrimstoneBonus
					warlock.Incinerate.DamageMultiplierAdditive += fireAndBrimstoneBonus
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.ChaosBolt.DamageMultiplierAdditive -= fireAndBrimstoneBonus
					warlock.Incinerate.DamageMultiplierAdditive -= fireAndBrimstoneBonus
				},
			},
			NumberOfTicks: 5 + warlock.Talents.MoltenCore,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 785/5 + 0.2*dot.Spell.SpellPower()
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)

				dot.Spell.DamageMultiplierAdditive += bonusPeriodicDamageMultiplier
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				dot.Spell.DamageMultiplierAdditive -= bonusPeriodicDamageMultiplier
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 460 + 0.2*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealDamage(sim, result)
		},
	})
}
