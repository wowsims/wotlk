package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) ApplyFrostTalents() {
	// Improved Icy Touch
	// Implemented outside

	// Toughness
	if deathKnight.Talents.Toughness > 0 {
		armorCoeff := 0.02 * float64(deathKnight.Talents.Toughness)
		deathKnight.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Armor,
			ModifiedStat: stats.Armor,
			Modifier: func(armor float64, _ float64) float64 {
				return armor * (1.0 + armorCoeff)
			},
		})
	}

	// Icy Reach
	// Pointless to Implement

	// Black Ice
	deathKnight.PseudoStats.FrostDamageDealtMultiplier += 0.02 * float64(deathKnight.Talents.BlackIce)
	deathKnight.PseudoStats.ShadowDamageDealtMultiplier += 0.02 * float64(deathKnight.Talents.BlackIce)

	// Nerves Of Cold Steel
	deathKnight.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*float64(deathKnight.Talents.NervesOfColdSteel))
	if deathKnight.Talents.NervesOfColdSteel == 1 {
		deathKnight.AutoAttacks.OHEffect.BaseDamage.Calculator = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1.08, true)
	} else if deathKnight.Talents.NervesOfColdSteel == 2 {
		deathKnight.AutoAttacks.OHEffect.BaseDamage.Calculator = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1.16, true)
	} else {
		deathKnight.AutoAttacks.OHEffect.BaseDamage.Calculator = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1.25, true)
	}

	// Icy Talons
	deathKnight.applyIcyTalons()

	// Lichborne
	// Pointless to Implement

	// Annihilation

	// TODO: Implement

	// Killing Machine
	deathKnight.applyKillingMachine()

	// Chill of the Grave
	// Implemented outside

	// Endless Winter
	if deathKnight.Talents.EndlessWinter > 0 {
		strengthCoeff := 0.02 * float64(deathKnight.Talents.EndlessWinter)
		deathKnight.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Strength,
			ModifiedStat: stats.Strength,
			Modifier: func(strength float64, _ float64) float64 {
				return strength * (1.0 + strengthCoeff)
			},
		})
	}

	// Frigid Dreadplate
	// TODO: Implement

	// Glacier rot
	// Implemented outside

	// Deathchill
	// TODO: Implement

	// Improved Icy Talons
	if deathKnight.Talents.ImprovedIcyTalons {
		deathKnight.PseudoStats.MeleeSpeedMultiplier *= 1.05
	}

	// Merciless Combat
	// Implemented Outside

	// Blood of the North

	// Tundra Stalker
	deathKnight.AddStat(stats.Expertise, 1.0*float64(deathKnight.Talents.TundraStalker)*core.ExpertisePerQuarterPercentReduction)
}

func (deathKnight *DeathKnight) glacielRotBonus() float64 {
	glacierRotCoeff := 1.0
	if deathKnight.Talents.GlacierRot == 1 {
		glacierRotCoeff = 1.07
	} else if deathKnight.Talents.GlacierRot == 2 {
		glacierRotCoeff = 1.13
	} else if deathKnight.Talents.GlacierRot == 3 {
		glacierRotCoeff = 1.20
	}

	return core.TernaryFloat64(deathKnight.DiseasesAreActive() && deathKnight.Talents.GlacierRot > 0, glacierRotCoeff, 1.0)
}

func (deathKnight *DeathKnight) mercilessCombatBonus(sim *core.Simulation) float64 {
	return core.TernaryFloat64(sim.IsExecutePhase35() && deathKnight.Talents.MercilessCombat > 0, 1.0+0.06*float64(deathKnight.Talents.MercilessCombat), 1.0)
}

func (deathKnight *DeathKnight) tundraStalkerBonus() float64 {
	return core.TernaryFloat64(deathKnight.FrostFeverDisease.IsActive(), 1.0+0.03*float64(deathKnight.Talents.TundraStalker), 1.0)
}

func (deathKnight *DeathKnight) applyKillingMachine() {
	if deathKnight.Talents.KillingMachine == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 51130}
	weaponMH := deathKnight.GetMHWeapon()
	procChance := (weaponMH.SwingSpeed * 5.0 / 60.0) * float64(deathKnight.Talents.KillingMachine)

	deathKnight.KillingMachineAura = deathKnight.RegisterAura(core.Aura{
		Label:    "Killing Machine Proc",
		ActionID: actionID,
		Duration: time.Second * 30.0,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			//TODO: add the other spells
			if spell == deathKnight.IcyTouch {
				aura.Deactivate(sim)
			}
		},
	})

	deathKnight.RegisterAura(core.Aura{
		Label:    "Killing Machine",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if spell != deathKnight.IcyTouch {
				return
			}

			if sim.RandomFloat("Killing Machine") < procChance {
				deathKnight.KillingMachineAura.Activate(sim)
			}
		},
	})
}

func (deathKnight *DeathKnight) applyIcyTalons() {
	if deathKnight.Talents.IcyTalons == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 50887}

	deathKnight.IcyTalonsAura = deathKnight.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Icy Talons",
		Duration: time.Second * 20.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier *= 1.0 + 0.04*float64(deathKnight.Talents.IcyTalons)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier /= 1.0 + 0.04*float64(deathKnight.Talents.IcyTalons)
		},
	})
}

func (deathKnight *DeathKnight) outcomeEitherWeaponHitOrCrit(mhOutcome core.HitOutcome, ohOutcome core.HitOutcome) bool {
	return mhOutcome == core.OutcomeHit || mhOutcome == core.OutcomeCrit || ohOutcome == core.OutcomeHit || ohOutcome == core.OutcomeCrit
}

func (deathKnight *DeathKnight) bloodOfTheNorthCoeff() float64 {
	bloodOfTheNorthCoeff := 1.0
	if deathKnight.Talents.BloodOfTheNorth == 1 {
		bloodOfTheNorthCoeff = 1.03
	} else if deathKnight.Talents.BloodOfTheNorth == 2 {
		bloodOfTheNorthCoeff = 1.06
	} else if deathKnight.Talents.BloodOfTheNorth == 3 {
		bloodOfTheNorthCoeff = 1.1
	}
	return bloodOfTheNorthCoeff
}

func (deathKnight *DeathKnight) bloodOfTheNorthChance() float64 {
	botnChance := 0.0
	if deathKnight.Talents.BloodOfTheNorth == 1 {
		botnChance = 0.3
	} else if deathKnight.Talents.BloodOfTheNorth == 2 {
		botnChance = 0.6
	} else if deathKnight.Talents.BloodOfTheNorth == 3 {
		botnChance = 1.0
	}
	return botnChance
}

func (deathKnight *DeathKnight) bloodOfTheNorthWillProc(sim *core.Simulation, botnChance float64) bool {
	ohWillCast := sim.RandomFloat("Blood of The North") <= botnChance
	return ohWillCast
}

func (deathKnight *DeathKnight) bloodOfTheNorthProc(sim *core.Simulation, spell *core.Spell, runeCost core.DKRuneCost) {
	if deathKnight.Talents.BloodOfTheNorth > 0 {
		if runeCost.Blood > 0 {
			botnChance := deathKnight.bloodOfTheNorthChance()

			if deathKnight.bloodOfTheNorthWillProc(sim, botnChance) {
				slot := deathKnight.SpendBloodRune(sim, spell.BloodRuneMetrics())
				deathKnight.SetRuneAtSlotToState(0, slot, core.RuneState_DeathSpent, core.RuneKind_Death)
				deathKnight.SetAsGeneratedByReapingOrBoTN(slot)
			} else {
				deathKnight.Spend(sim, spell, runeCost)
			}
		} else {
			deathKnight.Spend(sim, spell, runeCost)
		}
	}
}

func (deathKnight *DeathKnight) threatOfThassarianChance() float64 {
	threatOfThassarianChance := 0.0
	if deathKnight.Talents.ThreatOfThassarian == 1 {
		threatOfThassarianChance = 0.30
	} else if deathKnight.Talents.ThreatOfThassarian == 2 {
		threatOfThassarianChance = 0.60
	} else if deathKnight.Talents.ThreatOfThassarian == 3 {
		threatOfThassarianChance = 1.0
	}
	return threatOfThassarianChance
}

func (deathKnight *DeathKnight) threatOfThassarianWillProc(sim *core.Simulation, totChance float64) bool {
	ohWillCast := sim.RandomFloat("Threat of Thassarian") <= totChance
	return ohWillCast
}

func (deathKnight *DeathKnight) threatOfThassarianAdjustMetrics(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, mhOutcome core.HitOutcome) {
	spell.SpellMetrics[spellEffect.Target.TableIndex].Casts -= 1
	if mhOutcome == core.OutcomeHit {
		spell.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 1
	} else if mhOutcome == core.OutcomeCrit {
		spell.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 1
	} else {
		spell.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 2
	}
}

func (deathKnight *DeathKnight) threatOfThassarianProcMasks(isMH bool, effect *core.SpellEffect, guileOfGorefiend bool) {
	if isMH {
		effect.ProcMask = core.ProcMaskMeleeMHSpecial
		effect.OutcomeApplier = deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.critMultiplier(guileOfGorefiend))
	} else {
		effect.ProcMask = core.ProcMaskMeleeOHSpecial
		effect.OutcomeApplier = deathKnight.OutcomeFuncMeleeSpecialCritOnly(deathKnight.critMultiplier(guileOfGorefiend))
	}
}

func (deathKnight *DeathKnight) threatOfThassarianProc(sim *core.Simulation, spellEffect *core.SpellEffect, mhSpell *core.Spell, ohSpell *core.Spell) {
	totChance := deathKnight.threatOfThassarianChance()

	mhSpell.Cast(sim, spellEffect.Target)
	totProcced := deathKnight.threatOfThassarianWillProc(sim, totChance)
	if totProcced {
		ohSpell.Cast(sim, spellEffect.Target)
	}
}
