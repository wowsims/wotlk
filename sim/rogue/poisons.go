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
	rogue.applyWoundPoison()
}

var DeadlyPoisonActionID = core.ActionID{SpellID: 57973}

func (rogue *Rogue) registerPoisonAuras() {
	numTargets := rogue.Env.GetNumTargets()
	for i := int32(0); i < numTargets; i++ {
		target := rogue.Env.GetTargetUnit(i)
		if rogue.Talents.SavageCombat > 0 {
			rogue.savageCombatDebuffAuras = append(rogue.savageCombatDebuffAuras, core.SavageCombatAura(target, rogue.Talents.SavageCombat))
		}
		if rogue.Talents.MasterPoisoner > 0 {
			masterPoisonerAura := core.MasterPoisonerDebuff(target, float64(rogue.Talents.MasterPoisoner))
			masterPoisonerAura.Duration = core.NeverExpires
			rogue.masterPoisonerDebuffAuras = append(rogue.masterPoisonerDebuffAuras, masterPoisonerAura)
		}
	}
}

func (rogue *Rogue) registerDeadlyPoisonSpell() {
	rogue.DeadlyPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    DeadlyPoisonActionID,
		SpellSchool: core.SpellSchoolNature,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,
			OutcomeApplier:   rogue.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				target := spellEffect.Target
				if spellEffect.Landed() {
					dot := rogue.deadlyPoisonDots[target.Index]
					if dot.IsActive() {
						if dot.GetStacks() == 5 {
							if rogue.lastDeadlyPoisonProcMask.Matches(core.ProcMaskMeleeMH) {
								switch rogue.Options.OhImbue {
								case proto.Rogue_Options_DeadlyPoison:
									dot.Refresh(sim)
								case proto.Rogue_Options_InstantPoison:
									rogue.InstantPoison[1].Cast(sim, target)
								case proto.Rogue_Options_WoundPoison:
									rogue.WoundPoison[1].Cast(sim, target)
								}
							}
							if rogue.lastDeadlyPoisonProcMask.Matches(core.ProcMaskMeleeOH) {
								switch rogue.Options.MhImbue {
								case proto.Rogue_Options_DeadlyPoison:
									dot.Refresh(sim)
								case proto.Rogue_Options_InstantPoison:
									rogue.InstantPoison[1].Cast(sim, target)
								case proto.Rogue_Options_WoundPoison:
									rogue.WoundPoison[1].Cast(sim, target)
								}
							}
						}
						dot.Refresh(sim)
						dot.AddStack(sim)
					} else {
						dot.Apply(sim)
						dot.SetStacks(sim, 1)
					}
				}
				rogue.lastDeadlyPoisonProcMask = core.ProcMaskEmpty
			},
		}),
	})
	deadlyPoisonDebuffAura := core.Aura{
		Label:     "DeadlyPoison-" + strconv.Itoa(int(rogue.Index)),
		ActionID:  DeadlyPoisonActionID,
		MaxStacks: 5,
		Duration:  time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Talents.SavageCombat > 0 {
				rogue.savageCombatDebuffAuras[aura.Unit.Index].Activate(sim)
			}
			if rogue.Talents.MasterPoisoner > 0 {
				rogue.masterPoisonerDebuffAuras[aura.Unit.Index].Activate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Talents.SavageCombat > 0 {
				rogue.savageCombatDebuffAuras[aura.Unit.Index].Deactivate(sim)
			}
			if rogue.Talents.MasterPoisoner > 0 {
				rogue.masterPoisonerDebuffAuras[aura.Unit.Index].Deactivate(sim)
			}
		},
	}
	deadlyPoisonTickBaseDamage := core.BaseDamageConfig{
		Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
			return 74 + hitEffect.MeleeAttackPower(spell.Unit)*0.03
		},
		TargetSpellCoefficient: 1,
	}
	numTargets := rogue.Env.GetNumTargets()
	for i := int32(0); i < numTargets; i++ {
		target := rogue.Env.GetTargetUnit(i)
		dotAura := target.RegisterAura(deadlyPoisonDebuffAura)
		dot := core.NewDot(core.Dot{
			Spell:         rogue.DeadlyPoison,
			Aura:          dotAura,
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,
			TickEffects: core.TickFuncApplyEffectsToUnit(target, core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask: core.ProcMaskPeriodicDamage,
				DamageMultiplier: 1 +
					[]float64{0.0, 0.07, 0.14, 0.20}[rogue.Talents.VilePoisons],
				ThreatMultiplier: 1,
				IsPeriodic:       false, // hack to get attacker modifiers applied
				BaseDamage:       core.MultiplyByStacks(deadlyPoisonTickBaseDamage, dotAura),
				OutcomeApplier:   rogue.OutcomeFuncTickMagicHit(),
			})),
		})
		if rogue.HasSetBonus(ItemSetTerrorblade, 2) {
			metrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 64914})
			dot.OnPeriodicDamageDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				rogue.AddEnergy(sim, 1, metrics)
			}
		}
		// Would like to do this for the snapshotting but it also shots the aura
		//dot.TickEffects = core.TickFuncSnapshot(target, deadlyPoisonTickEffect)
		rogue.deadlyPoisonDots = append(rogue.deadlyPoisonDots, dot)
	}
}

func (rogue *Rogue) procDeadlyPoison(sim *core.Simulation, spellEffect *core.SpellEffect) {
	rogue.lastDeadlyPoisonProcMask = spellEffect.ProcMask
	rogue.DeadlyPoison.Cast(sim, spellEffect.Target)
}

func (rogue *Rogue) applyDeadlyPoison() {
	procMask := core.GetMeleeProcMaskForHands(
		rogue.Options.MhImbue == proto.Rogue_Options_DeadlyPoison,
		rogue.Options.OhImbue == proto.Rogue_Options_DeadlyPoison)

	if procMask == core.ProcMaskUnknown {
		return
	}
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
			if sim.RandomFloat("Deadly Poison") > rogue.GetDeadlyPoisonProcChance(procMask) {
				return
			}
			rogue.procDeadlyPoison(sim, spellEffect)
		},
	})
}

func (rogue *Rogue) GetWoundPoisonProcChance(mask core.ProcMask) float64 {
	if mask.Matches(core.ProcMaskMeleeMH) && rogue.Options.MhImbue != proto.Rogue_Options_WoundPoison {
		return 0.0
	}
	if mask.Matches(core.ProcMaskMeleeOH) && rogue.Options.OhImbue != proto.Rogue_Options_WoundPoison {
		return 0.0
	}
	return 0.5

}

func (rogue *Rogue) applyWoundPoison() {
	procMask := core.GetMeleeProcMaskForHands(
		rogue.Options.MhImbue == proto.Rogue_Options_WoundPoison,
		rogue.Options.OhImbue == proto.Rogue_Options_WoundPoison,
	)
	if procMask == core.ProcMaskUnknown {
		return
	}
	rogue.RegisterAura(core.Aura{
		Label:    "Wound Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Wound Poison") > rogue.GetWoundPoisonProcChance(procMask) {
				return
			}
			rogue.WoundPoison[0].Cast(sim, spellEffect.Target)
		},
	})
}

type PoisonProcSource int

const (
	NormalProc PoisonProcSource = iota
	DeadlyProc
	ShivProc
)

func (rogue *Rogue) makeInstantPoison(procSource PoisonProcSource) *core.Spell {
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 57968, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskEmpty,
			DamageMultiplier: 1 +
				[]float64{0.0, 0.07, 0.14, 0.20}[rogue.Talents.VilePoisons],
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

func (rogue *Rogue) makeWoundPoison(procSource PoisonProcSource) *core.Spell {
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 57975, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskEmpty,
			DamageMultiplier: 1 +
				[]float64{0.0, 0.07, 0.14, 0.20}[rogue.Talents.VilePoisons],
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return 231 + hitEffect.MeleeAttackPower(spell.Unit)*0.04
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: rogue.OutcomeFuncMagicHitAndCrit(rogue.SpellCritMultiplier()),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}
				rogue.woundPoisonDebuffAuras[spellEffect.Target.Index].Activate(sim)
			},
		}),
	})
}

var WoundPoisonActionID = core.ActionID{SpellID: 57975}

func (rogue *Rogue) registerWoundPoisonSpell() {
	woundPoisonDebuffAura := core.Aura{
		Label:    "WoundPoison-" + strconv.Itoa(int(rogue.Index)),
		ActionID: WoundPoisonActionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Talents.SavageCombat > 0 {
				rogue.savageCombatDebuffAuras[aura.Unit.Index].Activate(sim)
			}
			if rogue.Talents.MasterPoisoner > 0 {
				rogue.masterPoisonerDebuffAuras[aura.Unit.Index].Activate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Talents.SavageCombat > 0 {
				rogue.savageCombatDebuffAuras[aura.Unit.Index].Deactivate(sim)
			}
			if rogue.Talents.MasterPoisoner > 0 {
				rogue.masterPoisonerDebuffAuras[aura.Unit.Index].Deactivate(sim)
			}
		},
	}
	numTargets := rogue.Env.GetNumTargets()
	for i := int32(0); i < numTargets; i++ {
		target := rogue.Env.GetTargetUnit(i)
		rogue.woundPoisonDebuffAuras = append(rogue.woundPoisonDebuffAuras, target.RegisterAura(woundPoisonDebuffAura))
	}
	rogue.WoundPoison = [3]*core.Spell{
		rogue.makeWoundPoison(NormalProc),
		rogue.makeWoundPoison(DeadlyProc),
		rogue.makeWoundPoison(ShivProc),
	}
}

func (rogue *Rogue) registerInstantPoisonSpell() {
	rogue.InstantPoison = [3]*core.Spell{
		rogue.makeInstantPoison(NormalProc),
		rogue.makeInstantPoison(DeadlyProc),
		rogue.makeInstantPoison(ShivProc),
	}
}

func (rogue *Rogue) GetDeadlyPoisonProcChance(mask core.ProcMask) float64 {
	if mask.Matches(core.ProcMaskMeleeMH) && rogue.Options.MhImbue != proto.Rogue_Options_DeadlyPoison {
		return 0.0
	}
	if mask.Matches(core.ProcMaskMeleeOH) && rogue.Options.OhImbue != proto.Rogue_Options_DeadlyPoison {
		return 0.0
	}
	return 0.3 + 0.04*float64(rogue.Talents.ImprovedPoisons) + rogue.deadlyPoisonProcChanceBonus
}

func (rogue *Rogue) UpdateInstantPoisonPPM(bonusChance float64) {
	procMask := core.GetMeleeProcMaskForHands(
		rogue.Options.MhImbue == proto.Rogue_Options_InstantPoison,
		rogue.Options.OhImbue == proto.Rogue_Options_InstantPoison)

	ppm := 8.57 * (1 + float64(rogue.Talents.ImprovedPoisons)*0.1 + bonusChance)
	rogue.instantPoisonPPMM = rogue.AutoAttacks.NewPPMManager(ppm, procMask)
}

func (rogue *Rogue) applyInstantPoison() {
	procMask := core.GetMeleeProcMaskForHands(
		rogue.Options.MhImbue == proto.Rogue_Options_InstantPoison,
		rogue.Options.OhImbue == proto.Rogue_Options_InstantPoison)

	if procMask == core.ProcMaskUnknown {
		return
	}
	rogue.UpdateInstantPoisonPPM(0)
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
			if rogue.instantPoisonPPMM.Proc(sim, spellEffect.ProcMask, "Instant Poison") {
				rogue.InstantPoison[0].Cast(sim, spellEffect.Target)
			}
		},
	})
}
