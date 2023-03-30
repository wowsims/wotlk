package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (shaman *Shaman) ApplyTalents() {
	// We are going to treat this like a snapshot if you have the glyph.
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfTotemOfWrath) {
		shaman.AddStat(stats.SpellPower, 280*0.3)
	}

	shaman.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(shaman.Talents.ThunderingStrikes))
	shaman.AddStat(stats.SpellCrit, core.CritRatingPerCritChance*1*float64(shaman.Talents.ThunderingStrikes))
	shaman.AddStat(stats.Dodge, core.DodgeRatingPerDodgeChance*1*float64(shaman.Talents.Anticipation))
	shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= []float64{1, 1.04, 1.07, 1.1}[shaman.Talents.WeaponMastery]

	shaman.AddStat(stats.Expertise, 3*core.ExpertisePerQuarterPercentReduction*float64(shaman.Talents.UnleashedRage))

	if shaman.Talents.DualWieldSpecialization > 0 && shaman.HasOHWeapon() {
		shaman.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*2*float64(shaman.Talents.DualWieldSpecialization))
	}

	shaman.AddStat(stats.SpellCrit, float64(shaman.Talents.BlessingOfTheEternals)*2*core.CritRatingPerCritChance)
	if shaman.Talents.Toughness > 0 {
		shaman.MultiplyStat(stats.Stamina, 1.0+0.02*float64(shaman.Talents.Toughness))
	}
	if shaman.Talents.UnrelentingStorm > 0 {
		shaman.AddStatDependency(stats.Intellect, stats.MP5, 0.04*float64(shaman.Talents.UnrelentingStorm))
	}
	if shaman.Talents.AncestralKnowledge > 0 {
		shaman.MultiplyStat(stats.Intellect, 1.0+0.02*float64(shaman.Talents.AncestralKnowledge))
	}
	if shaman.Talents.MentalQuickness > 0 {
		shaman.AddStatDependency(stats.AttackPower, stats.SpellPower, 0.1*float64(shaman.Talents.MentalQuickness))
	}
	if shaman.Talents.MentalDexterity > 0 {
		shaman.AddStatDependency(stats.Intellect, stats.AttackPower, 0.3333*float64(shaman.Talents.MentalDexterity))
	}
	if shaman.Talents.NaturesBlessing > 0 {
		shaman.AddStatDependency(stats.Intellect, stats.SpellPower, 0.1*float64(shaman.Talents.NaturesBlessing))
	}

	if shaman.Talents.SpiritWeapons {
		shaman.PseudoStats.CanParry = true
		shaman.AutoAttacks.MHConfig.ThreatMultiplier *= 0.7
		shaman.AutoAttacks.OHConfig.ThreatMultiplier *= 0.7
	}

	shaman.applyElementalFocus()
	shaman.applyElementalDevastation()
	shaman.applyFlurry()
	shaman.applyMaelstromWeapon()
	shaman.registerElementalMasteryCD()
	shaman.registerNaturesSwiftnessCD()
	shaman.registerShamanisticRageCD()
	shaman.registerManaTideTotemCD()
}

func (shaman *Shaman) spellThreatMultiplier() float64 {
	return []float64{1, 0.9, 0.8, 0.7}[shaman.Talents.ElementalPrecision]
}

func (shaman *Shaman) applyElementalFocus() {
	if !shaman.Talents.ElementalFocus {
		return
	}

	oathBonus := 1 + (0.05 * float64(shaman.Talents.ElementalOath))
	var affectedSpells []*core.Spell

	clearcastingAura := shaman.RegisterAura(core.Aura{
		Label:     "Clearcasting",
		ActionID:  core.ActionID{SpellID: 16246},
		Duration:  time.Second * 15,
		MaxStacks: 2,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice([]*core.Spell{
				shaman.LightningBolt,
				shaman.ChainLightning,
				shaman.LavaBurst,
				shaman.FireNova,
				shaman.EarthShock,
				shaman.FlameShock,
				shaman.FrostShock,
			}, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CostMultiplier -= 0.4
			}
			if oathBonus > 1 {
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= oathBonus
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= oathBonus
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= oathBonus
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CostMultiplier += 0.4
			}
			if oathBonus > 1 {
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] /= oathBonus
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= oathBonus
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] /= oathBonus
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Flags.Matches(SpellFlagShock | SpellFlagFocusable) {
				return
			}
			if spell.ActionID.Tag == 6 { // Filter LO casts
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagShock | SpellFlagFocusable) {
				return
			}
			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			clearcastingAura.Activate(sim)
			clearcastingAura.SetStacks(sim, 2)
		},
	})
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			procAura.Activate(sim)
		},
	})
}

var eleMasterActionID = core.ActionID{SpellID: 16166}

func (shaman *Shaman) registerElementalMasteryCD() {
	if !shaman.Talents.ElementalMastery {
		return
	}

	cdTimer := shaman.NewTimer()
	cd := time.Minute * 3

	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfElementalMastery) {
		cd -= time.Second * 30
	}

	// TODO: Share CD with Natures Swiftness

	buffAura := shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery Haste",
		ActionID: core.ActionID{SpellID: 64701},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1.15)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1 / 1.15)
		},
	})

	emAura := shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery",
		ActionID: eleMasterActionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.ChainLightning.CastTimeMultiplier -= 1
			shaman.LavaBurst.CastTimeMultiplier -= 1
			shaman.LightningBolt.CastTimeMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.ChainLightning.CastTimeMultiplier += 1
			shaman.LavaBurst.CastTimeMultiplier += 1
			shaman.LightningBolt.CastTimeMultiplier += 1
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != shaman.LightningBolt && spell != shaman.ChainLightning && spell != shaman.LavaBurst {
				return
			}
			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			cdTimer.Set(sim.CurrentTime + cd)
			shaman.UpdateMajorCooldowns()
		},
	})

	eleMastSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: eleMasterActionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buffAura.Activate(sim)
			emAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: eleMastSpell,
		Type:  core.CooldownTypeDPS,
	})

	if shaman.HasSetBonus(ItemSetFrostWitchRegalia, 2) {
		shaman.RegisterAura(core.Aura{
			Label:    "Shaman T10 Elemental 2P Bonus",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if (spell == shaman.LightningBolt || spell == shaman.ChainLightning) && !eleMastSpell.CD.IsReady(sim) {
					*eleMastSpell.CD.Timer = core.Timer(time.Duration(*eleMastSpell.CD.Timer) - time.Second*2)
					shaman.UpdateMajorCooldowns() // this could get expensive because it will be called all the time.
				}
			},
		})
	}
}

func (shaman *Shaman) registerNaturesSwiftnessCD() {
	if !shaman.Talents.NaturesSwiftness {
		return
	}
	actionID := core.ActionID{SpellID: 16188}
	cdTimer := shaman.NewTimer()
	cd := time.Minute * 3

	nsAura := shaman.RegisterAura(core.Aura{
		Label:    "Natures Swiftness",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.ChainLightning.CastTimeMultiplier -= 1
			shaman.LavaBurst.CastTimeMultiplier -= 1
			shaman.LightningBolt.CastTimeMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.ChainLightning.CastTimeMultiplier += 1
			shaman.LavaBurst.CastTimeMultiplier += 1
			shaman.LightningBolt.CastTimeMultiplier += 1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell != shaman.LightningBolt && spell != shaman.ChainLightning && spell != shaman.LavaBurst {
				return
			}

			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			cdTimer.Set(sim.CurrentTime + cd)
			shaman.UpdateMajorCooldowns()
		},
	})

	nsSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Don't use NS unless we're casting a full-length lightning bolt, which is
			// the only spell shamans have with a cast longer than GCD.
			return !shaman.HasTemporarySpellCastSpeedIncrease()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			nsAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: nsSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (shaman *Shaman) applyFlurry() {
	if shaman.Talents.Flurry == 0 {
		return
	}

	bonus := 1.0 + 0.06*float64(shaman.Talents.Flurry)

	if shaman.HasSetBonus(ItemSetEarthshatterBattlegear, 4) {
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if result.Outcome.Matches(core.OutcomeCrit) {
				procAura.Activate(sim)
				procAura.SetStacks(sim, 3)
				icd.Reset() // the "charge protection" ICD isn't up yet
				return
			}

			// Remove a stack.
			if procAura.IsActive() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) && icd.IsReady(sim) {
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

	var t10BonusAura *core.Aura
	enhT10Bonus := false
	if shaman.HasSetBonus(ItemSetFrostWitchBattlegear, 4) {
		enhT10Bonus = true

		statDep := shaman.NewDynamicMultiplyStat(stats.AttackPower, 1.2)
		t10BonusAura = shaman.RegisterAura(core.Aura{
			Label:    "Maelstrom Power",
			ActionID: core.ActionID{SpellID: 70831},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.EnableDynamicStatDep(sim, statDep)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.DisableDynamicStatDep(sim, statDep)
			},
		})
	}

	// TODO: Don't forget to make it so that AA don't reset when casting when MW is active
	// for LB / CL / LvB
	// They can't actually hit while casting, but the AA timer doesnt reset if you cast during the AA timer.

	// For sim purposes maelstrom weapon only impacts CL / LB
	shaman.MaelstromWeaponAura = shaman.RegisterAura(core.Aura{
		Label:     "MaelstromWeapon Proc",
		ActionID:  core.ActionID{SpellID: 53817},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			multDiff := 0.2 * float64(newStacks-oldStacks)
			shaman.LightningBolt.CastTimeMultiplier -= multDiff
			shaman.ChainLightning.CastTimeMultiplier -= multDiff

			if enhT10Bonus && shaman.MaelstromWeaponAura.GetStacks() == 5 {
				if sim.RandomFloat("Maelstrom Power") < 0.15 {
					t10BonusAura.Activate(sim)
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagElectric) {
				return
			}
			shaman.MaelstromWeaponAura.Deactivate(sim)
		},
	})

	ppmm := shaman.AutoAttacks.NewPPMManager(core.TernaryFloat64(shaman.HasSetBonus(ItemSetWorldbreakerBattlegear, 4), 2.4, 2.0)*
		float64(shaman.Talents.MaelstromWeapon), core.ProcMaskMelee)
	// This aura is hidden, just applies stacks of the proc aura.
	shaman.RegisterAura(core.Aura{
		Label:    "MaelstromWeapon",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) || !result.Landed() {
				return
			}
			if !ppmm.Proc(sim, spell.ProcMask, "Maelstrom Weapon") {
				return
			}
			shaman.MaelstromWeaponAura.Activate(sim)
			shaman.MaelstromWeaponAura.AddStack(sim)
		},
	})
}

func (shaman *Shaman) registerManaTideTotemCD() {
	if !shaman.Talents.ManaTideTotem {
		return
	}

	mttAura := core.ManaTideTotemAura(shaman.GetCharacter(), shaman.Index)
	mttSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: core.ManaTideTotemActionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 5,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mttAura.Activate(sim)

			// If healing stream is active, cancel it while mana tide is up.
			if shaman.HealingStreamTotem.Hot(&shaman.Unit).IsActive() {
				for _, agent := range shaman.Party.Players {
					shaman.HealingStreamTotem.Hot(&agent.GetCharacter().Unit).Cancel(sim)
				}
			}

			// TODO: Current water totem buff needs to be removed from party/raid.
			if shaman.Totems.Water != proto.WaterTotem_NoWaterTotem {
				shaman.NextTotemDrops[WaterTotem] = sim.CurrentTime + time.Second*12
			}
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: mttSpell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return sim.CurrentTime > time.Second*30
		},
	})
}
