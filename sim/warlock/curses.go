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
		ActionID:    core.ActionID{SpellID: 27228},
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:  core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 0, 1)*500*time.Millisecond,
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
	if warlock.Rotation.Curse != proto.Warlock_Rotation_Weakness {
		return
	}
	baseCost := 0.1 * warlock.BaseMana
	warlock.CurseOfWeaknessAura = core.CurseOfWeaknessAura(warlock.CurrentTarget)
	warlock.CurseOfWeaknessAura.Duration = time.Minute * 2

	warlock.CurseOfWeakness = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 27226},
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:  core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 0, 1)*500*time.Millisecond,
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

// https://wotlk.wowhead.com/spell=11719/curse-of-tongues
func (warlock *Warlock) registerCurseOfTonguesSpell() {
	if warlock.Rotation.Curse != proto.Warlock_Rotation_Tongues {
		return
	}
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
				GCD:  core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 0, 1)*500*time.Millisecond,
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

// https://wotlk.wowhead.com/spell=27218/curse-of-agony
func (warlock *Warlock) registerCurseOfAgonySpell() {
	if warlock.Rotation.Curse != proto.Warlock_Rotation_Agony && warlock.Rotation.Curse != proto.Warlock_Rotation_Doom {
		return
	}
	actionID := core.ActionID{SpellID: 27218}
	baseCost := 0.1 * warlock.BaseMana
	target := warlock.CurrentTarget
	baseDmg := 1740 / 12.0
	baseDmg *= (1 + 0.05*float64(warlock.Talents.ImprovedCurseOfAgony))

	effect := core.SpellEffect{
		DamageMultiplier: 1 *
			(1 + 0.02*float64(warlock.Talents.ShadowMastery)) *
			(1 + 0.01*float64(warlock.Talents.Contagion)),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		BaseDamage:       core.BaseDamageConfigMagicNoRoll(baseDmg, 0.1),
		OutcomeApplier:   warlock.OutcomeFuncTick(),
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
				GCD:  core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 0, 1)*500*time.Millisecond,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ThreatMultiplier: 1,
			FlatThreatBonus:  0, // TODO
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
		NumberOfTicks: 12,
		TickLength:    time.Second * 2,
		TickEffects:   core.TickFuncSnapshot(target, effect),
	})
}

func (warlock *Warlock) registerCurseOfDoomSpell() {
	if warlock.Rotation.Curse != proto.Warlock_Rotation_Doom {
		return
	}
	actionID := core.ActionID{SpellID: 30910}
	baseCost := 0.15 * warlock.BaseMana

	target := warlock.CurrentTarget
	effect := core.SpellEffect{
		DamageMultiplier: 1 *
			(1 + 0.02*float64(warlock.Talents.ShadowMastery)) *
			(1 + 0.01*float64(warlock.Talents.Contagion)),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		BaseDamage:       core.BaseDamageConfigMagicNoRoll(7300, 2),
		OutcomeApplier:   warlock.OutcomeFuncTick(),
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
				GCD:  core.GCDDefault - core.TernaryDuration(warlock.Talents.AmplifyCurse, 0, 1)*500*time.Millisecond,
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
