package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (warlock *Warlock) registerDrainSoulSpell() {
	soulSiphonMultiplier := 0.03 * float64(warlock.Talents.SoulSiphon)

	warlock.DrainSoul = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47855},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagChanneled,

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
		// For performance optimization, the execute modifier is basekit since we never use it before execute
		DamageMultiplier: (4.0 + 0.04*float64(warlock.Talents.DeathsEmbrace)) / (1 + 0.04*float64(warlock.Talents.DeathsEmbrace)),
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
				dot.SnapshotBaseDamage = baseDmg * (1.0 + float64(core.MinInt(3, numActive))*soulSiphonMultiplier)

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
			}
		},
	})
}
