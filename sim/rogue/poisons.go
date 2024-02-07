package rogue

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

/**
Instant Poison: 20% proc chance
25: 22 +/- 3 damage, 8679 ID, 40 charges
40: 50 +/- 6 damage, 8688 ID, 70 charges
50: 76 +/- 9 damage, 11338 ID, 85 charges
60: 130 =/- 18 damage, 11340 ID, 115 charges

Deadly Poison: 30% proc chance, 5 stacks
25: 36 damage, 2823 ID, 60 charges (Deadly Brew only)
40: 52 damage, 2824 ID, 75 charges
50: 80 damage, 11355 ID, 90 charges
60: 108 damage, 11356 ID, 105 charges (Rank 4, Rank 5 is by book)

Wound Poison: 30% proc chance, 5 stacks
25: x damage, x ID (none, first rank is level 32)
40: -75 healing, 11325 ID, 75 charges (Rank 2)
50: -105 healing, 13226 ID, 90 charges (Rank 3)
60: -135 healing, 13227 ID, 105 charges (Rank 4)
*/

// TODO: Add charges to poisons (not deadly brew)
func (rogue *Rogue) getPoisonDamageMultiplier() float64 {
	return []float64{1, 1.04, 1.08, 1.12, 1.16, 1.2}[rogue.Talents.VilePoisons]
}

func (rogue *Rogue) applyPoisons() {
	rogue.applyDeadlyPoison()
	rogue.applyInstantPoison()
	rogue.applyWoundPoison()
}

func (rogue *Rogue) registerDeadlyPoisonSpell() {
	baseDamage := map[int32]float64{
		25: 36,
		40: 52,
		50: 80,
		60: 108,
	}[rogue.Level]
	spellID := map[int32]int32{
		25: 2823,
		40: 2824,
		50: 11355,
		60: 11356,
	}[rogue.Level]

	rogue.DeadlyPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskWeaponProc,

		DamageMultiplier: rogue.getPoisonDamageMultiplier(),
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

			OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				if stacks := dot.GetStacks(); stacks > 0 {
					dot.SnapshotBaseDamage = (baseDamage + core.TernaryFloat64(rogue.HasRune(proto.RogueRune_RuneDeadlyBrew), 0.09*dot.Spell.MeleeAttackPower(), 0)) * float64(stacks)
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)

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

			/** Old Wrath code for proccing other poisons from deadly at full stacks
			May be useful to reference when doing Deadly Brew implementation
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
			}*/
			dot.Refresh(sim)
			dot.TakeSnapshot(sim, false)
		},
	})
}

func (rogue *Rogue) procDeadlyPoison(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	rogue.lastDeadlyPoisonProcMask = spell.ProcMask
	// Cast Deadly Poison (including roll to hit)
	rogue.DeadlyPoison.Cast(sim, result.Target)
}

// Get the mask for poisons procs to determine which hand is being used
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

// Applies Deadly Poison to a weapon and enables Procs to be rolled on hits
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

// Apply Wound Poison to a weapon and enable procs to be rolled on hits
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

// TODO: remove Deadly Proc and replace with Deadly Brew procs
const (
	NormalProc PoisonProcSource = iota
	ShivProc
	DeadlyBrewProc
)

// Make a source based variant of Instant Poison
func (rogue *Rogue) makeInstantPoison(procSource PoisonProcSource) *core.Spell {
	baseDamageByLevel := map[int32]float64{
		25: 19,
		40: 44,
		50: 67,
		60: 112,
	}[rogue.Level]

	damageVariance := map[int32]float64{
		25: 6,
		40: 12,
		50: 18,
		60: 36,
	}[rogue.Level]

	spellID := map[int32]int32{
		25: 8679,
		40: 8688,
		50: 11338,
		60: 11340,
	}[rogue.Level]
	isShivProc := procSource == ShivProc

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskWeaponProc,

		DamageMultiplier: rogue.getPoisonDamageMultiplier(),
		CritMultiplier:   rogue.SpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageByLevel, baseDamageByLevel+damageVariance) + core.TernaryFloat64(rogue.HasRune(proto.RogueRune_RuneDeadlyBrew), 0.09*spell.MeleeAttackPower(), 0)
			if isShivProc {
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			} else {
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

// Make a source based variant of Wound Poison
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
			// No damage dealt by classic wound poison, fix
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

// TODO: Update spell id with phase/level
var WoundPoisonActionID = core.ActionID{SpellID: 57975}

func (rogue *Rogue) registerWoundPoisonSpell() {
	woundPoisonDebuffAura := core.Aura{
		Label:     "WoundPoison-" + strconv.Itoa(int(rogue.Index)),
		MaxStacks: 5,
		ActionID:  WoundPoisonActionID,
		Duration:  time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// all healing effects used on target reduced by x, stacks 5 times
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			// Undo reduced healing effects used on target
		},
	}

	rogue.woundPoisonDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		return target.RegisterAura(woundPoisonDebuffAura)
	})
	rogue.WoundPoison = [2]*core.Spell{
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

// Get the proc chance of Deadly Poison (between 0 and 1)
func (rogue *Rogue) GetDeadlyPoisonProcChance() float64 {
	return 0.3 + 0.04*float64(rogue.Talents.ImprovedPoisons)
}

// Get the proc chance of Instant Poison (between 0 and 1)
func (rogue *Rogue) GetInstantPoisonProcChance() float64 {
	return 0.2 + 0.04*float64(rogue.Talents.ImprovedPoisons) + rogue.instantPoisonProcChanceBonus
}

// get the proc chance of Wound Poison (between 0 and 1)
func (rogue *Rogue) GetWoundPoisonProcChance() float64 {
	return 0.3 + 0.04*float64(rogue.Talents.ImprovedPoisons)
}

// Apply Instant Poisons to a weapon and enable procs
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
