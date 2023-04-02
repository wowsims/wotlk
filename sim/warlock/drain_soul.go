package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (warlock *Warlock) registerDrainSoulSpell() {
	soulSiphonMultiplier := 0.03 * float64(warlock.Talents.SoulSiphon)

	calcSoulSiphonMult := func(target *core.Unit) float64 {
		auras := []*core.Aura{
			warlock.HauntDebuffAuras.Get(target),
			warlock.UnstableAffliction.Dot(target).Aura,
			warlock.Corruption.Dot(target).Aura,
			warlock.Seed.Dot(target).Aura,
			warlock.CurseOfAgony.Dot(target).Aura,
			warlock.CurseOfDoom.Dot(target).Aura,
			warlock.CurseOfElementsAuras.Get(target),
			warlock.CurseOfWeaknessAuras.Get(target),
			warlock.CurseOfTonguesAuras.Get(target),
			warlock.ShadowEmbraceDebuffAura(target),
			// missing: death coil
		}
		numActive := 0
		for _, aura := range auras {
			if aura.IsActive() {
				numActive++
			}
		}
		return 1.0 + float64(core.MinInt(3, numActive))*soulSiphonMultiplier
	}

	warlock.DrainSoul = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47855},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagChanneled | core.SpellFlagHauntSE,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.14,
			Multiplier: 1 - 0.02*float64(warlock.Talents.Suppression),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
				// ChannelTime: channelTime,
			},
		},

		DamageMultiplierAdditive: 1 +
			warlock.GrandSpellstoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery),
		DamageMultiplier: 1,
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Drain Soul",
			},
			NumberOfTicks:       5,
			TickLength:          3 * time.Second,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDmg := 142 + 0.429*dot.Spell.SpellPower()
				dot.SnapshotBaseDamage = baseDmg * calcSoulSiphonMult(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				dot := spell.Dot(target)
				dot.Apply(sim)
				dot.UpdateExpires(dot.ExpiresAt())

				warlock.everlastingAfflictionRefresh(sim, target)
			}
		},
		ExpectedDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDmg := (142 + 0.429*spell.SpellPower()) * calcSoulSiphonMult(target)
				return spell.CalcPeriodicDamage(sim, target, baseDmg, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int) {
			if isExecute == 25 {
				mult := (4.0 + 0.04*float64(warlock.Talents.DeathsEmbrace)) / (1 + 0.04*float64(warlock.Talents.DeathsEmbrace))
				warlock.DrainSoul.DamageMultiplier = mult
			}
		})
	})
}
