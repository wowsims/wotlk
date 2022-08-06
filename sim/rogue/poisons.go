package rogue

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) applyPoisons() {
	rogue.applyDeadlyPoison()
	rogue.applyInstantPoison()
}

func (rogue *Rogue) registerDeadlyPoisonSpell() {
	actionID := core.ActionID{SpellID: 57973}

	rogue.DeadlyPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskEmpty,
			BonusSpellHitRating: 5 * core.SpellHitRatingPerHitChance * float64(rogue.Talents.Precision),
			ThreatMultiplier:    1,
			OutcomeApplier:      rogue.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					if rogue.DeadlyPoisonDot.IsActive() {
						if rogue.DeadlyPoisonDot.GetStacks() == 5 {
							if rogue.LastDeadlyPoisonProcMask.Matches(core.ProcMaskMeleeMH) {
								switch rogue.Options.OhImbue {
								case proto.Rogue_Options_DeadlyPoison:
									rogue.DeadlyPoisonDot.Refresh(sim)
								case proto.Rogue_Options_InstantPoison:
									rogue.InstantPoison.Cast(sim, spellEffect.Target)
								}
							}
							if rogue.LastDeadlyPoisonProcMask.Matches(core.ProcMaskMeleeOH) {
								switch rogue.Options.MhImbue {
								case proto.Rogue_Options_DeadlyPoison:
									rogue.DeadlyPoisonDot.Refresh(sim)
								case proto.Rogue_Options_InstantPoison:
									rogue.InstantPoison.Cast(sim, spellEffect.Target)
								}
							}
						}
						rogue.DeadlyPoisonDot.Refresh(sim)
						rogue.DeadlyPoisonDot.AddStack(sim)
					} else {
						rogue.DeadlyPoisonDot.Apply(sim)
						rogue.DeadlyPoisonDot.SetStacks(sim, 1)
					}
				}
				rogue.LastDeadlyPoisonProcMask = core.ProcMaskEmpty
			},
		}),
	})

	target := rogue.CurrentTarget
	dotAura := target.RegisterAura(core.Aura{
		Label:     "DeadlyPoison-" + strconv.Itoa(int(rogue.Index)),
		ActionID:  actionID,
		MaxStacks: 5,
		Duration:  time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Talents.SavageCombat < 1 {
				return
			}
			savageCombatAura := core.SavageCombatAura(target, rogue.Talents.SavageCombat)
			savageCombatAura.Activate(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Talents.SavageCombat < 1 {
				return
			}
			savageCombatAura := core.SavageCombatAura(target, rogue.Talents.SavageCombat)
			savageCombatAura.Deactivate(sim)
		},
	})
	rogue.DeadlyPoisonDot = core.NewDot(core.Dot{
		Spell:         rogue.DeadlyPoison,
		Aura:          dotAura,
		NumberOfTicks: 4,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 + []float64{0.0, 0.07, 0.14, 0.20}[rogue.Talents.VilePoisons],
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.MultiplyByStacks(
				core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						return 74 + hitEffect.MeleeAttackPower(spell.Unit)*0.03
					},
					TargetSpellCoefficient: 1,
				},
				dotAura),
			OutcomeApplier: rogue.OutcomeFuncTickMagicHitAndCrit(rogue.SpellCritMultiplier()),
		})),
	})
}

func (rogue *Rogue) applyDeadlyPoison() {
	procMask := core.GetMeleeProcMaskForHands(
		rogue.Options.MhImbue == proto.Rogue_Options_DeadlyPoison,
		rogue.Options.OhImbue == proto.Rogue_Options_DeadlyPoison)

	if procMask == core.ProcMaskUnknown {
		return
	}

	procChance := 0.3 + 0.04*float64(rogue.Talents.ImprovedPoisons)

	rogue.RegisterAura(core.Aura{
		Label:    "Deadly Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Deadly Poison") > procChance {
				return
			}
			rogue.LastDeadlyPoisonProcMask = spellEffect.ProcMask
			rogue.DeadlyPoison.Cast(sim, spellEffect.Target)
		},
	})
}

func (rogue *Rogue) registerInstantPoisonSpell() {
	rogue.InstantPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 57968},
		SpellSchool: core.SpellSchoolNature,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1 + []float64{0.0, 0.07, 0.14, 0.20}[rogue.Talents.VilePoisons],
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return 300 + hitEffect.MeleeAttackPower(spell.Unit)*0.1
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: rogue.OutcomeFuncMagicHitAndCrit(rogue.SpellCritMultiplier()),
		}),
	})
}

func (rogue *Rogue) applyInstantPoison() {
	procMask := core.GetMeleeProcMaskForHands(
		rogue.Options.MhImbue == proto.Rogue_Options_InstantPoison,
		rogue.Options.OhImbue == proto.Rogue_Options_InstantPoison)

	if procMask == core.ProcMaskUnknown {
		return
	}

	var mhProcChance float64
	var ohProcChance float64
	if rogue.Options.MhImbue == proto.Rogue_Options_InstantPoison {
		mhProcChance = (rogue.GetMHWeapon().SwingSpeed * 8.57 * (1 + float64(rogue.Talents.ImprovedPoisons)*0.1)) / 60
	}
	if rogue.Options.OhImbue == proto.Rogue_Options_InstantPoison {
		ohProcChance = (rogue.GetOHWeapon().SwingSpeed * 8.57 * (1 + float64(rogue.Talents.ImprovedPoisons)*0.1)) / 60
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Instant Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(procMask) {
				return
			}
			if spellEffect.ProcMask.Matches(core.ProcMaskMeleeMH) && sim.RandomFloat("Instant Poison") > mhProcChance {
				return
			}
			if spellEffect.ProcMask.Matches(core.ProcMaskMeleeOH) && sim.RandomFloat("Instant Poison") > ohProcChance {
				return
			}
			rogue.procInstantPoison(sim, spellEffect)
		},
	})
}

func (rogue *Rogue) procInstantPoison(sim *core.Simulation, spellEffect *core.SpellEffect) {
	rogue.InstantPoison.Cast(sim, spellEffect.Target)
}
