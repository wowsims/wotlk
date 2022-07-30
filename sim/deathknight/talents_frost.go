package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) ApplyFrostTalents() {
	// Improved Icy Touch
	// Implemented outside

	// Toughness
	if dk.Talents.Toughness > 0 {
		armorCoeff := 0.02 * float64(dk.Talents.Toughness)
		dk.AddStatDependency(stats.Armor, stats.Armor, 1.0+armorCoeff)
	}

	// Icy Reach
	// Pointless to Implement

	// Black Ice
	dk.PseudoStats.FrostDamageDealtMultiplier *= 1.0 + 0.02*float64(dk.Talents.BlackIce)
	dk.PseudoStats.ShadowDamageDealtMultiplier *= 1.0 + 0.02*float64(dk.Talents.BlackIce)

	// Nerves Of Cold Steel
	dk.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*float64(dk.Talents.NervesOfColdSteel))
	dk.AutoAttacks.OHEffect.BaseDamage.Calculator = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, dk.nervesOfColdSteelBonus(), true)

	// Icy Talons
	dk.applyIcyTalons()

	// Lichborne
	// Pointless to Implement

	// Annihilation

	// Killing Machine
	dk.applyKillingMachine()

	// Chill of the Grave
	// Implemented outside

	// Endless Winter
	if dk.Talents.EndlessWinter > 0 {
		strengthCoeff := 0.02 * float64(dk.Talents.EndlessWinter)
		dk.AddStatDependency(stats.Strength, stats.Strength, 1.0+strengthCoeff)
	}

	// Frigid Dreadplate
	if dk.Talents.FrigidDreadplate > 0 {
		dk.PseudoStats.BonusMeleeHitRatingTaken -= core.MeleeHitRatingPerHitChance * float64(dk.Talents.FrigidDreadplate)
	}

	// Glacier rot
	dk.applyGlacierRot()

	// Deathchill
	// TODO: Implement

	// Improved Icy Talons
	if dk.Talents.ImprovedIcyTalons {
		dk.PseudoStats.MeleeSpeedMultiplier *= 1.05
	}

	// Merciless Combat
	dk.applyMercilessCombat()

	// Blood of the North
	dk.applyBloodOfTheNorth()

	dk.applyThreatOfThassarian()

	// Rime
	dk.applyRime()

	// Tundra Stalker
	dk.AddStat(stats.Expertise, 1.0*float64(dk.Talents.TundraStalker)*core.ExpertisePerQuarterPercentReduction)
	dk.applyTundaStalker()
}

func (dk *Deathknight) nervesOfColdSteelBonus() float64 {
	return []float64{1.0, 1.08, 1.16, 1.25}[dk.Talents.NervesOfColdSteel]
}

func (dk *Deathknight) applyGlacierRot() {
	dk.bonusCoeffs.glacierRotBonusCoeff = []float64{1.0, 1.07, 1.13, 1.20}[dk.Talents.GlacierRot]
}

func (dk *Deathknight) glacielRotBonus(target *core.Unit) float64 {
	return core.TernaryFloat64(dk.DiseasesAreActive(target), dk.bonusCoeffs.glacierRotBonusCoeff, 1.0)
}

func (dk *Deathknight) applyMercilessCombat() {
	dk.bonusCoeffs.mercilessCombatBonusCoeff = 1.0 + 0.06*float64(dk.Talents.MercilessCombat)
}

func (dk *Deathknight) mercilessCombatBonus(sim *core.Simulation) float64 {
	return core.TernaryFloat64(sim.IsExecutePhase35() && dk.Talents.MercilessCombat > 0, dk.bonusCoeffs.mercilessCombatBonusCoeff, 1.0)
}

func (dk *Deathknight) applyTundaStalker() {
	dk.bonusCoeffs.tundraStalkerBonusCoeff = 1.0 + 0.03*float64(dk.Talents.TundraStalker)
}

func (dk *Deathknight) tundraStalkerBonus(target *core.Unit) float64 {
	return core.TernaryFloat64(dk.FrostFeverDisease[target.Index].IsActive(), dk.bonusCoeffs.tundraStalkerBonusCoeff, 1.0)
}

func (dk *Deathknight) applyRime() {
	if dk.Talents.Rime == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 59057}

	dk.RimeAura = dk.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Rime",
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.HowlingBlast.CD.Reset()
		},
	})
}

func (dk *Deathknight) rimeCritBonus() float64 {
	return 0.05 * float64(dk.Talents.Rime)
}

func (dk *Deathknight) rimeHbChanceProc() float64 {
	return 5.0 * float64(dk.Talents.Rime)
}

func (dk *Deathknight) annihilationCritBonus() float64 {
	return 1.0 * float64(dk.Talents.Annihilation)
}

func (dk *Deathknight) applyKillingMachine() {
	if dk.Talents.KillingMachine == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 51130}
	//weaponMH := dk.GetMHWeapon()
	//procChance := (weaponMH.SwingSpeed * 5.0 / 60.0) * float64(dk.Talents.KillingMachine)

	ppmm := dk.AutoAttacks.NewPPMManager(float64(dk.Talents.KillingMachine), core.ProcMaskMeleeMHAuto|core.ProcMaskMeleeMHSpecial)

	dk.KillingMachineAura = dk.RegisterAura(core.Aura{
		Label:    "Killing Machine Proc",
		ActionID: actionID,
		Duration: time.Second * 30.0,
	})

	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label: "Killing Machine",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}

			if !ppmm.Proc(sim, spellEffect.ProcMask, "killing machine") {
				return
			}

			if !dk.KillingMachineAura.IsActive() {
				dk.KillingMachineAura.Activate(sim)
			} else {
				dk.KillingMachineAura.Refresh(sim)
			}
		},
	}))
}

func (dk *Deathknight) killingMachineOutcomeMod(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
	return func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, attackTable *core.AttackTable) {
		if dk.KillingMachineAura.IsActive() {
			spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			outcomeApplier(sim, spell, spellEffect, attackTable)
			spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
		} else {
			outcomeApplier(sim, spell, spellEffect, attackTable)
		}
	}
}

func (dk *Deathknight) applyIcyTalons() {
	if dk.Talents.IcyTalons == 0 {
		return
	}

	icyTalonsCoeff := 1.0 + 0.04*float64(dk.Talents.IcyTalons)

	dk.IcyTalonsAura = dk.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 50887}, // This probably doesnt need to be in metrics.
		Label:    "Icy Talons",
		Duration: time.Second * 20.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier *= icyTalonsCoeff
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier /= icyTalonsCoeff
		},
	})
}

func (dk *Deathknight) outcomeEitherWeaponLanded(mhOutcome core.HitOutcome, ohOutcome core.HitOutcome) bool {
	return mhOutcome.Matches(core.OutcomeLanded) // || ohOutcome.Matches(core.OutcomeLanded) TODO: Confirm that this should be gone
}

func (dk *Deathknight) bloodOfTheNorthCoeff() float64 {
	return []float64{1.0, 1.03, 1.06, 1.10}[dk.Talents.BloodOfTheNorth]
}

func (dk *Deathknight) applyBloodOfTheNorth() {
	dk.bonusCoeffs.bloodOfTheNorthChance = []float64{0.0, 0.3, 0.6, 1.0}[dk.Talents.BloodOfTheNorth]
}

func (dk *Deathknight) bloodOfTheNorthWillProc(sim *core.Simulation) bool {
	ohWillCast := sim.RandomFloat("Blood of The North") <= dk.bonusCoeffs.bloodOfTheNorthChance
	return ohWillCast
}

func (dk *Deathknight) bloodOfTheNorthProc(sim *core.Simulation, spell *core.Spell, runeCost core.RuneAmount) bool {
	if dk.Talents.BloodOfTheNorth > 0 {
		if runeCost.Blood > 0 {
			if dk.bloodOfTheNorthWillProc(sim) {
				slot := dk.SpendBloodRune(sim, spell.BloodRuneMetrics())
				dk.SetRuneAtIdxSlotToState(0, slot, core.RuneState_DeathSpent, core.RuneKind_Death)
				dk.SetAsGeneratedByReapingOrBoTN(slot)
				return true
			}
		}
	}
	return false
}

func (dk *Deathknight) applyThreatOfThassarian() {
	dk.bonusCoeffs.threatOfThassarianChance = []float64{0.0, 0.3, 0.6, 1.0}[dk.Talents.ThreatOfThassarian]
}

func (dk *Deathknight) threatOfThassarianWillProc(sim *core.Simulation) bool {
	ohWillCast := sim.RandomFloat("Threat of Thassarian") <= dk.bonusCoeffs.threatOfThassarianChance
	return ohWillCast
}

func (dk *Deathknight) threatOfThassarianAdjustMetrics(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, mhOutcome core.HitOutcome) {
	spell.SpellMetrics[spellEffect.Target.TableIndex].Casts -= 1
	if mhOutcome == core.OutcomeHit {
		spell.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 1
	} else if mhOutcome == core.OutcomeCrit {
		spell.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 1
	} else {
		spell.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 2
	}
}

func (dk *Deathknight) threatOfThassarianProcMasks(isMH bool, effect *core.SpellEffect, isGuileOfGorefiendStrike bool, isMightOfMograineStrike bool, wrapper func(outcomeApplier core.OutcomeApplier) core.OutcomeApplier) {
	critMultiplier := dk.critMultiplier()
	if isGuileOfGorefiendStrike || isMightOfMograineStrike {
		critMultiplier = dk.critMultiplierGoGandMoM()
	}

	if isMH {
		effect.ProcMask = core.ProcMaskMeleeMHSpecial
		effect.OutcomeApplier = wrapper(dk.OutcomeFuncMeleeSpecialHitAndCrit(critMultiplier))
	} else {
		effect.ProcMask = core.ProcMaskMeleeOHSpecial
		effect.OutcomeApplier = wrapper(dk.OutcomeFuncMeleeSpecialCritOnly(critMultiplier))
	}
}

func (dk *Deathknight) threatOfThassarianProc(sim *core.Simulation, spellEffect *core.SpellEffect, mhSpell *core.Spell, ohSpell *core.Spell) {
	mhSpell.Cast(sim, spellEffect.Target)
	if dk.Talents.ThreatOfThassarian > 0 && dk.threatOfThassarianWillProc(sim) {
		ohSpell.Cast(sim, spellEffect.Target)
	}
}
