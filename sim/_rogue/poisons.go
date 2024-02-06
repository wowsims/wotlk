package rogue

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) applyPoisons() {
	rogue.applyDeadlyPoison()
	rogue.applyInstantPoison()
	rogue.applyWoundPoison()
}

// TODO: Add Deadly Brew AP scaling
// TODO: Add level based damage for each poison
func (rogue *Rogue) registerDeadlyPoisonSpell() {

	rogue.DeadlyPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 57970},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskWeaponProc,

		DamageMultiplier: []float64{1, 1.04, 1.08, 1.12, 1.16, 1.2}[rogue.Talents.VilePoisons],
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "DeadlyPoison",
				MaxStacks: 5,
				Duration:  time.Second * 12,
				OnGain:    func(aura *core.Aura, sim *core.Simulation) {},
				OnExpire:  func(aura *core.Aura, sim *core.Simulation) {},
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			// TODO: Update Deadly Poison damage
			OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				if stacks := dot.GetStacks(); stacks > 0 {
					dot.SnapshotBaseDamage = (74 + 0.027*dot.Spell.MeleeAttackPower()) * float64(stacks)
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if !result.Landed() {
				return
			}

			dot := spell.Dot(target)
			if !dot.IsActive() {
				dot.Apply(sim)
				dot.SetStacks(sim, 1)
				dot.TakeSnapshot(sim, false)
				return
			}

			if dot.GetStacks() < 5 {
				dot.Refresh(sim)
				dot.AddStack(sim)
				dot.TakeSnapshot(sim, false)
				return
			}

			if rogue.lastDeadlyPoisonProcMask.Matches(core.ProcMaskMeleeMH) {
				switch rogue.Options.OhImbue {
				case proto.Rogue_Options_InstantPoison:
					rogue.InstantPoison[DeadlyProc].Cast(sim, target)
				case proto.Rogue_Options_WoundPoison:
					rogue.WoundPoison[DeadlyProc].Cast(sim, target)
				}
			}
			if rogue.lastDeadlyPoisonProcMask.Matches(core.ProcMaskMeleeOH) {
				switch rogue.Options.MhImbue {
				case proto.Rogue_Options_InstantPoison:
					rogue.InstantPoison[DeadlyProc].Cast(sim, target)
				case proto.Rogue_Options_WoundPoison:
					rogue.WoundPoison[DeadlyProc].Cast(sim, target)
				}
			}
			dot.Refresh(sim)
			dot.TakeSnapshot(sim, false)
		},
	})
}

func (rogue *Rogue) procDeadlyPoison(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	rogue.lastDeadlyPoisonProcMask = spell.ProcMask
	rogue.DeadlyPoison.Cast(sim, result.Target)
}

func (rogue *Rogue) getPoisonProcMask(imbue proto.Rogue_Options_PoisonImbue) core.ProcMask {
	var mask core.ProcMask
	if rogue.Options.MhImbue == imbue {
		mask |= core.ProcMaskMeleeMH
	}
	if rogue.Options.OhImbue == imbue {
		mask |= core.ProcMaskMeleeOH
	}
	return mask
}

func (rogue *Rogue) applyDeadlyPoison() {
	procMask := rogue.getPoisonProcMask(proto.Rogue_Options_DeadlyPoison)
	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Deadly Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Deadly Poison") < rogue.GetDeadlyPoisonProcChance() {
				rogue.procDeadlyPoison(sim, spell, result)
			}
		},
	})
}

func (rogue *Rogue) applyWoundPoison() {
	procMask := rogue.getPoisonProcMask(proto.Rogue_Options_WoundPoison)
	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Wound Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if sim.RandomFloat("Wound Poison") < rogue.GetWoundPoisonProcChance() {
				rogue.WoundPoison[NormalProc].Cast(sim, result.Target)
			}
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
	isShivProc := procSource == ShivProc

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 57965, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskWeaponProc,

		DamageMultiplier: []float64{1, 1.07, 1.14, 1.20}[rogue.Talents.VilePoisons],
		CritMultiplier:   rogue.SpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(300, 400) + 0.09*spell.MeleeAttackPower()
			if isShivProc {
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			} else {
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

func (rogue *Rogue) makeWoundPoison(procSource PoisonProcSource) *core.Spell {
	isShivProc := procSource == ShivProc

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 57975, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskWeaponProc,

		DamageMultiplier: []float64{1, 1.04, 1.08, 1.12, 1.16, 1.20}[rogue.Talents.VilePoisons],
		CritMultiplier:   rogue.SpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 231 + 0.036*spell.MeleeAttackPower()

			var result *core.SpellResult
			if isShivProc {
				result = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			} else {
				result = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			if result.Landed() {
				rogue.woundPoisonDebuffAuras.Get(target).Activate(sim)
			}
		},
	})
}

var WoundPoisonActionID = core.ActionID{SpellID: 57975}

func (rogue *Rogue) registerWoundPoisonSpell() {
	woundPoisonDebuffAura := core.Aura{
		Label:    "WoundPoison-" + strconv.Itoa(int(rogue.Index)),
		ActionID: WoundPoisonActionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
		},
	}

	rogue.woundPoisonDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.RegisterAura(woundPoisonDebuffAura)
	})
	rogue.WoundPoison = [3]*core.Spell{
		rogue.makeWoundPoison(NormalProc),
		rogue.makeWoundPoison(ShivProc),
	}
}

func (rogue *Rogue) registerInstantPoisonSpell() {
	rogue.InstantPoison = [3]*core.Spell{
		rogue.makeInstantPoison(NormalProc),
		rogue.makeInstantPoison(ShivProc),
	}
}

func (rogue *Rogue) GetDeadlyPoisonProcChance() float64 {
	return 0.3 + 0.04*float64(rogue.Talents.ImprovedPoisons)
}

func (rogue *Rogue) GetInstantPoisonProcChance() float64 {
	return 0.2 + 0.04*float64(rogue.Talents.ImprovedPoisons) + rogue.instantPoisonProcChanceBonus
}

func (rogue *Rogue) GetWoundPoisonProcChance() float64 {
	return 0.3 + 0.04*float64(rogue.Talents.ImprovedPoisons)
}

func (rogue *Rogue) applyInstantPoison() {
	procMask := rogue.getPoisonProcMask(proto.Rogue_Options_InstantPoison)
	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.instantPoisonProcChanceBonus = 0.0

	rogue.RegisterAura(core.Aura{
		Label:    "Instant Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if sim.RandomFloat("Instant Poison") < rogue.GetInstantPoisonProcChance() {
				rogue.InstantPoison[NormalProc].Cast(sim, result.Target)
			}
		},
	})
}
