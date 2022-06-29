package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warlock *Warlock) registerCurseOfElementsSpell() {
	if warlock.Rotation.Curse != proto.Warlock_Rotation_Elements {
		return
	}
	baseCost := 145.0
	auras := warlock.CurrentTarget.GetAurasWithTag("Curse of Elements")
	for _, aura := range auras {
		if int32(aura.Priority) >= warlock.Talents.Malediction {
			// Someone else with at least as good of curse is already doing it... lets not.
			warlock.Rotation.Curse = proto.Warlock_Rotation_NoCurse // TODO: swap to agony for dps?
			return
		}
	}
	warlock.CurseOfElementsAura = core.CurseOfElementsAura(warlock.CurrentTarget, warlock.Talents.Malediction)
	warlock.CurseOfElementsAura.Duration = time.Minute * 5

	warlock.CurseOfElements = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27228},
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ThreatMultiplier: 1,
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

func (warlock *Warlock) registerCurseOfRecklessnessSpell() {
	if warlock.Rotation.Curse != proto.Warlock_Rotation_Recklessness {
		return
	}
	baseCost := 160.0
	warlock.CurseOfRecklessnessAura = core.CurseOfRecklessnessAura(warlock.CurrentTarget)
	warlock.CurseOfRecklessnessAura.Duration = time.Minute * 2

	warlock.CurseOfRecklessness = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 27226},
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ThreatMultiplier: 1,
			FlatThreatBonus:  0, // TODO
			OutcomeApplier:   warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt:  applyAuraOnLanded(warlock.CurseOfRecklessnessAura),
			ProcMask:         core.ProcMaskEmpty,
		}),
	})
}

// https://tbc.wowhead.com/spell=11719/curse-of-tongues
func (warlock *Warlock) registerCurseOfTonguesSpell() {
	if warlock.Rotation.Curse != proto.Warlock_Rotation_Tongues {
		return
	}
	actionID := core.ActionID{SpellID: 11719}
	baseCost := 110.0

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
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ThreatMultiplier: 1,
			FlatThreatBonus:  0, // TODO
			OutcomeApplier:   warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt:  applyAuraOnLanded(warlock.CurseOfTonguesAura),
			ProcMask:         core.ProcMaskEmpty,
		}),
	})
}

// https://tbc.wowhead.com/spell=27218/curse-of-agony
func (warlock *Warlock) registerCurseOfAgonySpell() {
	if warlock.Rotation.Curse != proto.Warlock_Rotation_Agony && warlock.Rotation.Curse != proto.Warlock_Rotation_Doom {
		return
	}
	actionID := core.ActionID{SpellID: 27218}
	baseCost := 265.0
	target := warlock.CurrentTarget
	baseDmg := 1356.0 / 12.0
	baseDmg *= (1 + 0.02*float64(warlock.Talents.ImprovedCurseOfAgony))

	effect := core.SpellEffect{
		DamageMultiplier: 1 *
			(1 + 0.02*float64(warlock.Talents.ShadowMastery)) *
			(1 + 0.01*float64(warlock.Talents.Contagion)),
		ThreatMultiplier: 1 - 0.05*float64(warlock.Talents.ImprovedDrainSoul),
		BaseDamage:       core.BaseDamageConfigMagicNoRoll(baseDmg, 0.1),
		OutcomeApplier:   warlock.OutcomeFuncTick(),
		IsPeriodic:       true,
		ProcMask:         core.ProcMaskPeriodicDamage,
	}
	// Amplify Curse talent
	if warlock.Talents.AmplifyCurse {
		effect.BaseDamage = core.WrapBaseDamageConfig(effect.BaseDamage, func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
			return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				if warlock.AmplifyCurseAura.IsActive() {
					return oldCalculator(sim, hitEffect, spell) * 1.5
				} else {
					return oldCalculator(sim, hitEffect, spell)
				}
			}
		})
	}
	warlock.CurseOfAgony = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
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
				if warlock.AmplifyCurseAura != nil && warlock.AmplifyCurseAura.IsActive() {
					warlock.AmplifyCurseAura.Deactivate(sim)
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
	baseCost := 380.0

	target := warlock.CurrentTarget
	effect := core.SpellEffect{
		DamageMultiplier: 1 *
			(1 + 0.02*float64(warlock.Talents.ShadowMastery)) *
			(1 + 0.01*float64(warlock.Talents.Contagion)),
		ThreatMultiplier: 1 - 0.05*float64(warlock.Talents.ImprovedDrainSoul),
		BaseDamage:       core.BaseDamageConfigMagicNoRoll(4200, 2),
		OutcomeApplier:   warlock.OutcomeFuncTick(),
		IsPeriodic:       true,
		ProcMask:         core.ProcMaskPeriodicDamage,
	}
	// Amplify Curse talent
	if warlock.Talents.AmplifyCurse {
		effect.BaseDamage = core.WrapBaseDamageConfig(effect.BaseDamage, func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
			return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				if warlock.AmplifyCurseAura.IsActive() {
					return oldCalculator(sim, hitEffect, spell) * 1.5
				} else {
					return oldCalculator(sim, hitEffect, spell)
				}
			}
		})
	}

	warlock.CurseOfDoom = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
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
				if warlock.AmplifyCurseAura != nil && warlock.AmplifyCurseAura.IsActive() {
					warlock.AmplifyCurseAura.Deactivate(sim)
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
