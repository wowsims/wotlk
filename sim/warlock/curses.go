package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerCurseOfElementsSpell() {
	if warlock.Rotation.Curse != proto.Warlock_Rotation_Elements {
		return
	}
	baseCost := 0.1 * warlock.BaseMana
	warlock.CurseOfElementsAura = core.CurseOfElementsAura(warlock.CurrentTarget)
	warlock.CurseOfElementsAura.Duration = time.Minute * 5

	warlock.CurseOfElements = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47865},
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:  core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 1, 0)*500*time.Millisecond,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
			FlatThreatBonus:  0, // TODO
			OutcomeApplier:   warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt:  applyAuraOnLanded(warlock.CurseOfElementsAura),
			ProcMask:         core.ProcMaskEmpty,
		}),
	})
}

func (warlock *Warlock) ShouldCastCurseOfElements(sim *core.Simulation, target *core.Unit, curse proto.Warlock_Rotation_Curse) bool {
	return curse == proto.Warlock_Rotation_Elements && !warlock.CurseOfElementsAura.IsActive()
}

func (warlock *Warlock) registerCurseOfWeaknessSpell() {
	baseCost := 0.1 * warlock.BaseMana
	warlock.CurseOfWeaknessAura = core.CurseOfWeaknessAura(warlock.CurrentTarget, warlock.Talents.ImprovedCurseOfWeakness)
	warlock.CurseOfWeaknessAura.Duration = time.Minute * 2

	warlock.CurseOfWeakness = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 50511},
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:  core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 1, 0)*500*time.Millisecond,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
			FlatThreatBonus:  0, // TODO
			OutcomeApplier:   warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt:  applyAuraOnLanded(warlock.CurseOfWeaknessAura),
			ProcMask:         core.ProcMaskEmpty,
		}),
	})
}

func (warlock *Warlock) registerCurseOfTonguesSpell() {
	actionID := core.ActionID{SpellID: 11719}
	baseCost := 0.04 * warlock.BaseMana

	// Empty aura so we can simulate cost/time to keep tongues up
	warlock.CurseOfTonguesAura = warlock.CurrentTarget.GetOrRegisterAura(core.Aura{
		Label:    "Curse of Tongues",
		ActionID: actionID,
		Duration: time.Second * 30,
	})

	warlock.CurseOfTongues = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:  core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 1, 0)*500*time.Millisecond,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
			FlatThreatBonus:  0, // TODO
			OutcomeApplier:   warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt:  applyAuraOnLanded(warlock.CurseOfTonguesAura),
			ProcMask:         core.ProcMaskEmpty,
		}),
	})
}

func (warlock *Warlock) registerCurseOfAgonySpell() {
	actionID := core.ActionID{SpellID: 47864}
	baseCost := 0.1 * warlock.BaseMana
	target := warlock.CurrentTarget
	hasGoCoA := warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfCurseOfAgony)
	numberOfTicks := 12
	totalBaseDmg := 1740.0
	agonyEffect := totalBaseDmg * 0.056
	if hasGoCoA {
		numberOfTicks += 2
		totalBaseDmg += 2 * agonyEffect // Glyphed ticks
	}
	applier := warlock.OutcomeFuncTick()
	if warlock.Talents.Pandemic {
		applier = warlock.OutcomeFuncMagicCrit(warlock.SpellCritMultiplier(1, 1))
	}

	effect := core.SpellEffect{
		DamageMultiplier: (1 + 0.01*float64(warlock.Talents.Contagion)) * (1.0 + 0.05*float64(warlock.Talents.ImprovedCurseOfAgony)),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		BaseDamage:       core.BaseDamageConfigMagicNoRoll(totalBaseDmg/float64(numberOfTicks), 0.1), //TODO : CoA ramp up effect 
		OutcomeApplier:   applier,
		IsPeriodic:       true,
		ProcMask:         core.ProcMaskPeriodicDamage,
	}
	warlock.CurseOfAgony = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:  core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 1, 0)*500*time.Millisecond,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ThreatMultiplier: 1,
			FlatThreatBonus:  0, // TODO : curses flat threat on application
			OutcomeApplier:   warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					warlock.CurseOfAgonyDot.Apply(sim)
				}
			},
			ProcMask: core.ProcMaskEmpty,
		}),
	})
	warlock.CurseOfAgonyDot = core.NewDot(core.Dot{
		Spell: warlock.CurseOfAgony,
		Aura: target.RegisterAura(core.Aura{
			Label:    "CurseofAgony-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: numberOfTicks,
		TickLength:    time.Second * 2,
		TickEffects:   core.TickFuncSnapshot(target, effect),
	})
}

func (warlock *Warlock) registerCurseOfDoomSpell() {
	actionID := core.ActionID{SpellID: 47867}
	baseCost := 0.15 * warlock.BaseMana
		applier = 

	target := warlock.CurrentTarget
	effect := core.SpellEffect{
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		DamageMultiplier: 1,
		BaseDamage:       core.BaseDamageConfigMagicNoRoll(7300, 2),
		OutcomeApplier:   warlock.OutcomeFuncMagicCrit(warlock.DefaultSpellCritMultiplier()),
		IsPeriodic:       true,
		ProcMask:         core.ProcMaskPeriodicDamage,
	}

	warlock.CurseOfDoom = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:  core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 1, 0)*500*time.Millisecond,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ThreatMultiplier: 1,
			FlatThreatBonus:  0, // TODO
			OutcomeApplier:   warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					warlock.CurseOfDoomDot.Apply(sim)
				}
			},
			ProcMask: core.ProcMaskEmpty,
		}),
	})

	warlock.CurseOfDoomDot = core.NewDot(core.Dot{
		Spell: warlock.CurseOfDoom,
		Aura: target.RegisterAura(core.Aura{
			Label:    "CurseofDoom-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 1,
		TickLength:    time.Minute,
		TickEffects:   core.TickFuncSnapshot(target, effect),
	})
}

func applyAuraOnLanded(aura *core.Aura) func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
	return func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		if spellEffect.Landed() {
			aura.Activate(sim)
		}
	}
}
