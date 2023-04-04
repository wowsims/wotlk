package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) registerCurseOfElementsSpell() {
	warlock.CurseOfElementsAuras = warlock.NewEnemyAuraArray(core.CurseOfElementsAura)

	warlock.CurseOfElements = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47865},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.1,
			Multiplier: 1 - 0.02*float64(warlock.Talents.Suppression),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 1, 0)*500*time.Millisecond,
			},
		},

		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		FlatThreatBonus:  156,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfElementsAuras.Get(target).Activate(sim)
			}
		},
	})
}

func (warlock *Warlock) ShouldCastCurseOfElements(sim *core.Simulation, target *core.Unit, curse proto.Warlock_Rotation_Curse) bool {
	return curse == proto.Warlock_Rotation_Elements && !warlock.CurseOfElementsAuras.Get(target).IsActive()
}

func (warlock *Warlock) registerCurseOfWeaknessSpell() {
	warlock.CurseOfWeaknessAuras = warlock.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.CurseOfWeaknessAura(target, warlock.Talents.ImprovedCurseOfWeakness)
	})

	warlock.CurseOfWeakness = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50511},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.1,
			Multiplier: 1 - 0.02*float64(warlock.Talents.Suppression),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 1, 0)*500*time.Millisecond,
			},
		},

		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		FlatThreatBonus:  142,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfWeaknessAuras.Get(target).Activate(sim)
			}
		},
	})
}

func (warlock *Warlock) registerCurseOfTonguesSpell() {
	actionID := core.ActionID{SpellID: 11719}

	// Empty aura so we can simulate cost/time to keep tongues up
	warlock.CurseOfTonguesAuras = warlock.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Curse of Tongues",
			ActionID: actionID,
			Duration: time.Second * 30,
		})
	})

	warlock.CurseOfTongues = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.04,
			Multiplier: 1 - 0.02*float64(warlock.Talents.Suppression),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 1, 0)*500*time.Millisecond,
			},
		},

		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		FlatThreatBonus:  100,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfTonguesAuras.Get(target).Activate(sim)
			}
		},
	})
}

func (warlock *Warlock) registerCurseOfAgonySpell() {
	numberOfTicks := core.TernaryInt32(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfCurseOfAgony), 14, 12)
	baseTickDmg := 145.0

	warlock.CurseOfAgony = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47864},
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagHauntSE,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.1,
			Multiplier: 1 - 0.02*float64(warlock.Talents.Suppression),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 1, 0)*500*time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1 +
			0.03*float64(warlock.Talents.ShadowMastery) +
			0.01*float64(warlock.Talents.Contagion) +
			0.05*float64(warlock.Talents.ImprovedCurseOfAgony),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		FlatThreatBonus:  0,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "CurseofAgony",
			},
			NumberOfTicks: numberOfTicks,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 0.5*baseTickDmg + 0.1*dot.Spell.SpellPower()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
				if dot.TickCount%4 == 0 { // CoA ramp up
					dot.SnapshotBaseDamage += 0.5 * baseTickDmg
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				warlock.CurseOfDoom.Dot(target).Cancel(sim)
				spell.Dot(target).Apply(sim)
			}
		},
	})
}

func (warlock *Warlock) registerCurseOfDoomSpell() {
	warlock.CurseOfDoom = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47867},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.15,
			Multiplier: 1 - 0.02*float64(warlock.Talents.Suppression),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 1, 0)*500*time.Millisecond,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute,
			},
		},

		DamageMultiplierAdditive: 1 +
			0.03*float64(warlock.Talents.ShadowMastery),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		FlatThreatBonus:  160,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "CurseofDoom",
			},
			NumberOfTicks: 1,
			TickLength:    time.Minute,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 7300 + 2*dot.Spell.SpellPower()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfAgony.Dot(target).Cancel(sim)
				spell.Dot(target).Apply(sim)
			}
		},
	})
}
