package shaman

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (shaman *Shaman) ApplyTalents() {
	if shaman.Talents.NaturesGuidance > 0 {
		shaman.AddStat(stats.SpellHit, float64(shaman.Talents.NaturesGuidance)*1*core.SpellHitRatingPerHitChance)
		shaman.AddStat(stats.MeleeHit, float64(shaman.Talents.NaturesGuidance)*1*core.MeleeHitRatingPerHitChance)
	}

	if shaman.Talents.ThunderingStrikes > 0 {
		shaman.AddStat(stats.MeleeCrit, core.MeleeCritRatingPerCritChance*1*float64(shaman.Talents.ThunderingStrikes))
	}

	shaman.AddStat(stats.Dodge, core.DodgeRatingPerDodgeChance*1*float64(shaman.Talents.Anticipation))
	shaman.AddStat(stats.Block, core.BlockRatingPerBlockChance*1*float64(shaman.Talents.ShieldSpecialization))
	shaman.AddStat(stats.Armor, shaman.Equip.Stats()[stats.Armor]*0.02*float64(shaman.Talents.Toughness))
	shaman.PseudoStats.PhysicalDamageDealtMultiplier *= 1 + 0.02*float64(shaman.Talents.WeaponMastery)

	if shaman.Talents.ShieldSpecialization > 0 {
		bonus := 1 + 0.05*float64(shaman.Talents.ShieldSpecialization)
		shaman.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.BlockValue,
			ModifiedStat: stats.BlockValue,
			Modifier: func(bv float64, _ float64) float64 {
				return bv * bonus
			},
		})
	}

	if shaman.Talents.DualWieldSpecialization > 0 && shaman.HasOHWeapon() {
		shaman.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*2*float64(shaman.Talents.DualWieldSpecialization))
	}

	if shaman.Talents.UnrelentingStorm > 0 {
		coeff := 0.02 * float64(shaman.Talents.UnrelentingStorm)
		shaman.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.MP5,
			Modifier: func(intellect float64, mp5 float64) float64 {
				return mp5 + intellect*coeff
			},
		})
	}

	if shaman.Talents.AncestralKnowledge > 0 {
		coeff := 0.01 * float64(shaman.Talents.AncestralKnowledge)
		shaman.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Mana,
			ModifiedStat: stats.Mana,
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
	shaman.applyShamanisticFocus()
	shaman.applyUnleashedRage()
	shaman.registerElementalMasteryCD()
	shaman.registerNaturesSwiftnessCD()
	shaman.registerShamanisticRageCD()
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
	if shaman.ClearcastingAura != nil && shaman.ClearcastingAura.IsActive() {
		// Reduces mana cost by 40%
		cast.Cost -= spell.BaseCost * 0.4
	}
}

func (shaman *Shaman) applyElementalDevastation() {
	if shaman.Talents.ElementalDevastation == 0 {
		return
	}

	critBonus := 3.0 * float64(shaman.Talents.ElementalDevastation) * core.SpellCritRatingPerCritChance
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

	shaman.ElementalMasteryAura = shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.AddStatDynamic(sim, stats.SpellCrit, 100*core.SpellCritRatingPerCritChance)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.AddStatDynamic(sim, stats.SpellCrit, -100*core.SpellCritRatingPerCritChance)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spell.Flags.Matches(SpellFlagShock | SpellFlagElectric) {
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
			shaman.ElementalMasteryAura.Activate(sim)
			shaman.ElementalMasteryAura.Prioritize()
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

func (shaman *Shaman) applyUnleashedRage() {
	if shaman.Talents.UnleashedRage == 0 {
		return
	}
	level := shaman.Talents.UnleashedRage

	bonusCoeff := 0.02 * float64(level)
	var currentAPBonuses []float64
	var urAuras = make([]*core.Aura, len(shaman.Party.PlayersAndPets))

	for i, playerOrPet := range shaman.Party.PlayersAndPets {
		char := playerOrPet.GetCharacter()
		idx := i
		urAuras[i] = char.GetOrRegisterAura(core.Aura{
			Label:    "Unleahed Rage Proc",
			ActionID: core.ActionID{SpellID: 30811},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				buffs := char.ApplyStatDependencies(stats.Stats{stats.AttackPower: currentAPBonuses[idx]})
				char.AddStatsDynamic(sim, buffs)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				buffs := char.ApplyStatDependencies(stats.Stats{stats.AttackPower: currentAPBonuses[idx]})
				unbuffs := buffs.Multiply(-1)
				char.AddStatsDynamic(sim, unbuffs)
			},
		})
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Unleashed Rage",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			currentAPBonuses = make([]float64, len(shaman.Party.PlayersAndPets))
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			// proc mask = 20 (melee auto & special)
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			for i, playerOrPet := range shaman.Party.PlayersAndPets {
				char := playerOrPet.GetCharacter()
				prevBonus := currentAPBonuses[i]
				newBonus := (char.GetStat(stats.AttackPower) - prevBonus) * bonusCoeff

				if prevBonus != newBonus {
					urAuras[i].Deactivate(sim)
					currentAPBonuses[i] = newBonus
					urAuras[i].Activate(sim)
				} else if newBonus != 0 {
					// If the bonus is the same, we can just refresh.
					urAuras[i].Refresh(sim)
				}
			}
		},
	})
}

func (shaman *Shaman) applyShamanisticFocus() {
	if !shaman.Talents.ShamanisticFocus {
		return
	}

	shaman.ShamanisticFocusAura = shaman.RegisterAura(core.Aura{
		Label:    "Shamanistic Focus Proc",
		ActionID: core.ActionID{SpellID: 43338},
		Duration: core.NeverExpires,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Flags.Matches(SpellFlagShock) {
				aura.Deactivate(sim)
			}
		},
	})

	shaman.RegisterAura(core.Aura{
		Label:    "Shamanistic Focus",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			shaman.ShamanisticFocusAura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyFlurry() {
	if shaman.Talents.Flurry == 0 {
		return
	}

	bonus := 1.05 + 0.05*float64(shaman.Talents.Flurry)
	if ItemSetCataclysmHarness.CharacterHasSetBonus(&shaman.Character, 4) {
		bonus += 0.05
	}
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
