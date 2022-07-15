package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (shaman *Shaman) ApplyTalents() {
	if shaman.Talents.ThunderingStrikes > 0 {
		shaman.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(shaman.Talents.ThunderingStrikes))
		shaman.AddStat(stats.SpellCrit, core.CritRatingPerCritChance*1*float64(shaman.Talents.ThunderingStrikes))
	}

	shaman.AddStat(stats.Dodge, core.DodgeRatingPerDodgeChance*1*float64(shaman.Talents.Anticipation))
	shaman.PseudoStats.PhysicalDamageDealtMultiplier *= []float64{0, 1.04, 1.07, 1.1}[shaman.Talents.WeaponMastery]

	if shaman.Talents.DualWieldSpecialization > 0 && shaman.HasOHWeapon() {
		shaman.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*2*float64(shaman.Talents.DualWieldSpecialization))
	}

	if shaman.Talents.BlessingOfTheEternals > 0 {
		shaman.AddStat(stats.SpellCrit, float64(shaman.Talents.BlessingOfTheEternals)*2*core.CritRatingPerCritChance)
	}

	if shaman.Talents.Toughness > 0 {
		coeff := 1 + 0.02*float64(shaman.Talents.Toughness)
		shaman.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Stamina,
			ModifiedStat: stats.Stamina,
			Modifier: func(stm float64, _ float64) float64 {
				return stm * coeff
			},
		})
	}

	if shaman.Talents.UnrelentingStorm > 0 {
		coeff := 0.04 * float64(shaman.Talents.UnrelentingStorm)
		shaman.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.MP5,
			Modifier: func(intellect float64, mp5 float64) float64 {
				return mp5 + intellect*coeff
			},
		})
	}

	if shaman.Talents.AncestralKnowledge > 0 {
		coeff := 0.02 * float64(shaman.Talents.AncestralKnowledge)
		shaman.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.Intellect,
			Modifier: func(mana float64, _ float64) float64 {
				return mana + mana*coeff
			},
		})
	}

	if shaman.Talents.MentalQuickness > 0 {
		coeff := 0.1 * float64(shaman.Talents.MentalQuickness)
		shaman.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.AttackPower,
			ModifiedStat: stats.SpellPower,
			Modifier: func(attackPower float64, spellPower float64) float64 {
				return spellPower + attackPower*coeff
			},
		})
	}

	if shaman.Talents.MentalDexterity > 0 {
		coeff := 0.3333 * float64(shaman.Talents.MentalDexterity)
		shaman.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.AttackPower,
			Modifier: func(intellect float64, attackPower float64) float64 {
				return attackPower + intellect*coeff
			},
		})
	}

	if shaman.Talents.NaturesBlessing > 0 {
		coeff := 0.1 * float64(shaman.Talents.NaturesBlessing)
		shaman.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.SpellPower,
			Modifier: func(intellect float64, spellPower float64) float64 {
				return spellPower + intellect*coeff
			},
		})
	}

	if shaman.Talents.SpiritWeapons {
		shaman.PseudoStats.CanParry = true
		shaman.AutoAttacks.MHEffect.ThreatMultiplier *= 0.7
		shaman.AutoAttacks.OHEffect.ThreatMultiplier *= 0.7
	}

	shaman.applyElementalFocus()
	shaman.applyElementalDevastation()
	shaman.applyFlurry()
	shaman.applyMaelstromWeapon()
	shaman.registerElementalMasteryCD()
	shaman.registerNaturesSwiftnessCD()
	shaman.registerShamanisticRageCD()

	// TODO: FeralSpirit Spirit summons
}

func (shaman *Shaman) applyElementalFocus() {
	if !shaman.Talents.ElementalFocus {
		return
	}

	shaman.ClearcastingAura = shaman.RegisterAura(core.Aura{
		Label:     "Clearcasting",
		ActionID:  core.ActionID{SpellID: 16246},
		Duration:  time.Second * 15,
		MaxStacks: 2,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Flags.Matches(SpellFlagShock | SpellFlagElectric) {
				return
			}
			if spell.ActionID.Tag != 0 { // Filter LO casts
				return
			}
			aura.RemoveStack(sim)
		},
	})

	shaman.RegisterAura(core.Aura{
		Label:    "Elemental Focus",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spell.Flags.Matches(SpellFlagShock | SpellFlagElectric) {
				return
			}
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			shaman.ClearcastingAura.Activate(sim)
			shaman.ClearcastingAura.SetStacks(sim, 2)
		},
	})
}

func (shaman *Shaman) modifyCastClearcasting(spell *core.Spell, cast *core.Cast) {
	if shaman.ClearcastingAura.IsActive() {
		// Reduces mana cost by 40%
		cast.Cost -= spell.BaseCost * 0.4
	}
}

func (shaman *Shaman) modifyCastMaelstrom(spell *core.Spell, cast *core.Cast) {
	if shaman.MaelstromWeaponAura.GetStacks() > 0 {
		castReduction := float64(shaman.MaelstromWeaponAura.GetStacks()) * 0.2
		cast.CastTime -= time.Duration(float64(cast.CastTime) * castReduction)
	}
}

func (shaman *Shaman) applyElementalDevastation() {
	if shaman.Talents.ElementalDevastation == 0 {
		return
	}

	critBonus := 3.0 * float64(shaman.Talents.ElementalDevastation) * core.CritRatingPerCritChance
	procAura := shaman.NewTemporaryStatsAura("Elemental Devastation Proc", core.ActionID{SpellID: 30160}, stats.Stats{stats.MeleeCrit: critBonus}, time.Second*10)

	shaman.RegisterAura(core.Aura{
		Label:    "Elemental Devastation",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			procAura.Activate(sim)
		},
	})
}

func (shaman *Shaman) registerElementalMasteryCD() {
	if !shaman.Talents.ElementalMastery {
		return
	}
	actionID := core.ActionID{SpellID: 16166}
	cdTimer := shaman.NewTimer()
	cd := time.Minute * 3

	// TODO: Share CD with Natures Swiftness

	shaman.ElementalMasteryBuffAura = shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery Haste",
		ActionID: core.ActionID{SpellID: 64701},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.AddStatDynamic(sim, stats.SpellHaste, 15*core.HasteRatingPerHastePercent)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.AddStatDynamic(sim, stats.SpellHaste, -15*core.HasteRatingPerHastePercent)
		},
	})

	shaman.ElementalMasteryAura = shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spell.Flags.Matches(SpellFlagElectric) {
				// Only LB / CL / LvB use EM
				if spell.ActionID.SpellID != lavaBurstActionID.SpellID {
					return
				}
			}
			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			cdTimer.Set(sim.CurrentTime + cd)
			shaman.UpdateMajorCooldowns()
		},
	})

	spell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shaman.ElementalMasteryBuffAura.Activate(sim)
			shaman.ElementalMasteryAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (shaman *Shaman) registerNaturesSwiftnessCD() {
	if !shaman.Talents.NaturesSwiftness {
		return
	}
	actionID := core.ActionID{SpellID: 16188}
	cdTimer := shaman.NewTimer()
	cd := time.Minute * 3

	shaman.NaturesSwiftnessAura = shaman.RegisterAura(core.Aura{
		Label:    "Natures Swiftness",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell != shaman.LightningBolt {
				return
			}

			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			cdTimer.Set(sim.CurrentTime + cd)
			shaman.UpdateMajorCooldowns()
		},
	})

	spell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shaman.NaturesSwiftnessAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Don't use NS unless we're casting a full-length lightning bolt, which is
			// the only spell shamans have with a cast longer than GCD.
			return !character.HasTemporarySpellCastSpeedIncrease()
		},
	})
}

func (shaman *Shaman) applyFlurry() {
	if shaman.Talents.Flurry == 0 {
		return
	}

	bonus := 1.0 + 0.06*float64(shaman.Talents.Flurry)

	// I believe there is a set in wotlk that improves flurry.

	// if shaman.HasSetBonus(ItemSetCataclysmHarness, 4) {
	// 	bonus += 0.05
	// }

	inverseBonus := 1 / bonus

	procAura := shaman.RegisterAura(core.Aura{
		Label:     "Flurry Proc",
		ActionID:  core.ActionID{SpellID: 16280},
		Duration:  core.NeverExpires,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, bonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, inverseBonus)
		},
	})

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond * 500,
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Flurry",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				procAura.Activate(sim)
				procAura.SetStacks(sim, 3)
				icd.Reset() // the "charge protection" ICD isn't up yet
				return
			}

			// Remove a stack.
			if procAura.IsActive() && spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) && icd.IsReady(sim) {
				icd.Use(sim)
				procAura.RemoveStack(sim)
			}
		},
	})
}

func (shaman *Shaman) applyMaelstromWeapon() {
	if shaman.Talents.MaelstromWeapon == 0 {
		return
	}

	// TODO: Don't forget to make it so that AA don't reset when casting when MW is active
	// for LB / CL / LvB
	// They can't actually hit while casting, but the AA timer doesnt reset if you cast during the AA timer.

	// For sim purposes maelstrom weapon only impacts CL / LB
	procAura := shaman.RegisterAura(core.Aura{
		Label:     "MaelstromWeapon Proc",
		ActionID:  core.ActionID{SpellID: 53817},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spell.Flags.Matches(SpellFlagElectric) {
				return
			}
			shaman.MaelstromWeaponAura.Deactivate(sim)
		},
	})
	shaman.MaelstromWeaponAura = procAura

	// This aura is hidden, just applies stacks of the proc aura.
	shaman.RegisterAura(core.Aura{
		Label:    "MaelstromWeapon",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}
			if !procAura.IsActive() {
				procAura.Activate(sim)
			}
			procAura.AddStack(sim)
		},
	})
}
