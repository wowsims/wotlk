package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) ApplyTalents() {
	// Demonic Embrace
	if warlock.Talents.DemonicEmbrace > 0 {
		warlock.AddStatDependency(stats.Stamina, stats.Stamina, 1.01+(float64(warlock.Talents.DemonicEmbrace)*0.03))
	}

	// Suppression
	warlock.AddStat(stats.SpellHit, float64(warlock.Talents.Suppression)*core.SpellHitRatingPerHitChance)

	// Shadow Mastery
	warlock.PseudoStats.ShadowDamageDealtMultiplier *= 1.0 + (0.03 * float64(warlock.Talents.ShadowMastery))

	// Backlash (Add 1% crit per level)
	warlock.AddStat(stats.SpellCrit, float64(warlock.Talents.Backlash)*core.CritRatingPerCritChance)

	// Malediction (SP bonus)
	if warlock.Talents.Malediction > 0 {
		warlock.AddStatDependency(stats.SpellPower, stats.SpellPower, 1.0+(0.01*float64(warlock.Talents.Malediction)))
	}

	// Fel Vitality
	if warlock.Talents.FelVitality > 0 {
		bonus := 0.01 * float64(warlock.Talents.FelVitality)
		// Adding a second 3% bonus int->mana dependency
		warlock.AddStatDependency(stats.Intellect, stats.Mana, 1.0+bonus)
	}

	if warlock.Options.Summon != proto.Warlock_Options_NoSummon {
		if warlock.Talents.MasterDemonologist > 0 {
			switch warlock.Options.Summon {
			case proto.Warlock_Options_Imp:
				warlock.PseudoStats.FireDamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
				warlock.PseudoStats.BonusFireCritRating *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
			case proto.Warlock_Options_Succubus:
				warlock.PseudoStats.ShadowDamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
				warlock.PseudoStats.BonusShadowCritRating *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
			case proto.Warlock_Options_Felguard:
				warlock.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
			}
		}
		// Extract stats for demonic knowledge
		if warlock.Talents.DemonicKnowledge > 0 {
			petChar := warlock.Pets[0].GetCharacter()
			bonus := (petChar.GetStat(stats.Stamina) + petChar.GetStat(stats.Intellect)) * (0.04 * float64(warlock.Talents.DemonicKnowledge))
			warlock.AddStat(stats.SpellPower, bonus)
		}
	}

	// Demonic Tactics, applies even without pet out
	if warlock.Talents.DemonicTactics > 0 {
		warlock.AddStats(stats.Stats{
			stats.MeleeCrit: float64(warlock.Talents.DemonicTactics) * 2 * core.CritRatingPerCritChance,
			stats.SpellCrit: float64(warlock.Talents.DemonicTactics) * 2 * core.CritRatingPerCritChance,
		})
	}

	if warlock.Talents.DemonicPact > 0 {
		spellDmgBonus := 0.02 * float64(warlock.Talents.DemonicPact)
		warlock.PseudoStats.ShadowDamageDealtMultiplier *= 1 + spellDmgBonus
		warlock.PseudoStats.FireDamageDealtMultiplier *= 1 + spellDmgBonus
	}

	if warlock.Talents.MoltenSkin > 0 {
		warlock.PseudoStats.DamageTakenMultiplier /= 1 + 0.02*float64(warlock.Talents.MoltenSkin)
	}

	if warlock.Talents.Nightfall > 0 || warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfCorruption) {
		warlock.setupNightfall()
	}

	if warlock.Talents.ShadowEmbrace > 0 {
		warlock.setupShadowEmbrace()
	}

	if warlock.Talents.Eradication > 0 {
		warlock.setupEradication()
	}

	if warlock.Talents.DeathsEmbrace > 0 {
		warlock.applyDeathsEmbrace()
	}

	warlock.setupDrainSoulExecutePhase()

	if warlock.Talents.MoltenCore > 0 {
		warlock.setupMoltenCore()
	}

	if warlock.Talents.Decimation > 0 {
		warlock.setupDecimation()
	}

	if warlock.Talents.Pyroclasm > 0 {
		warlock.setupPyroclasm()
	}

	if warlock.Talents.Backdraft > 0 {
		warlock.setupBackdraft()
	}

	if warlock.Talents.EmpoweredImp > 0 && warlock.Options.Summon == proto.Warlock_Options_Imp {
		warlock.Pet.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.1*float64(warlock.Talents.EmpoweredImp)
		warlock.setupEmpoweredImp()
	}

	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap) {
		warlock.registerGlyphOfLifeTapAura()
	}
}

func (warlock *Warlock) registerGlyphOfLifeTapAura() {
	warlock.GlyphOfLifeTapAura = warlock.RegisterAura(core.Aura{
		Label:    "Glyph Of LifeTap Aura",
		Tag:      "Glyph Of LifeTap Aura",
		ActionID: core.ActionID{SpellID: 63941},
		Duration: time.Second * 40,
		Priority: float64(warlock.GetStat(stats.Spirit)),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AddStatDynamic(sim, stats.SpellPower, 0.2*float64(warlock.GetStat(stats.Spirit)))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AddStatDynamic(sim, stats.SpellPower, -0.2*float64(warlock.GetStat(stats.Spirit)))
		},
	})
}

func (warlock *Warlock) setupEmpoweredImp() {
	warlock.EmpoweredImpAura = warlock.RegisterAura(core.Aura{
		Label:    "Empowered Imp Proc Aura",
		ActionID: core.ActionID{SpellID: 47283},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AddStatDynamic(sim, stats.SpellCrit, 100*core.CritRatingPerCritChance)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AddStatDynamic(sim, stats.SpellCrit, -100*core.CritRatingPerCritChance)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				aura.Deactivate(sim)
			}
		},
	})

	warlock.Pet.RegisterAura(core.Aura{
		Label:    "Empowered Imp Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				warlock.EmpoweredImpAura.Activate(sim)
				warlock.EmpoweredImpAura.Refresh(sim)
			}
		},
	})
}

func (warlock *Warlock) applyDeathsEmbrace() {
	multiplier := 1.0 + 0.04*float64(warlock.Talents.DeathsEmbrace)

	deathsEmbraceAura := warlock.RegisterAura(core.Aura{
		Label:    "Death's Embrace Hidden Aura",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.ShadowDamageDealtMultiplier *= multiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.ShadowDamageDealtMultiplier /= multiplier
		},
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute35 bool) {
			if isExecute35 {
				deathsEmbraceAura.Activate(sim)
			}
		})
	})
}

func (warlock *Warlock) setupDecimation() {
	warlock.DecimationAura = warlock.RegisterAura(core.Aura{
		Label:    "Decimation Proc Aura",
		ActionID: core.ActionID{SpellID: 63167},
		Duration: time.Second * 10,
	})

	decimation := warlock.RegisterAura(core.Aura{
		Label:    "Decimation Talent Hidden Aura",
		Duration: core.NeverExpires,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warlock.ShadowBolt || spell == warlock.Incinerate || spell == warlock.SoulFire {
				warlock.DecimationAura.Activate(sim)
				warlock.DecimationAura.Refresh(sim)
			}
		},
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute35 bool) {
			if isExecute35 {
				decimation.Activate(sim)
			}
		})
	})
}

func (warlock *Warlock) setupPyroclasm() {
	warlock.PyroclasmAura = warlock.RegisterAura(core.Aura{
		Label:    "Pyroclasm",
		ActionID: core.ActionID{SpellID: 63244},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ShadowDamageDealtMultiplier *= 1 + 0.02*float64(warlock.Talents.Pyroclasm)
			aura.Unit.PseudoStats.FireDamageDealtMultiplier *= 1 + 0.02*float64(warlock.Talents.Pyroclasm)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ShadowDamageDealtMultiplier /= 1 + 0.02*float64(warlock.Talents.Pyroclasm)
			aura.Unit.PseudoStats.FireDamageDealtMultiplier /= 1 + 0.02*float64(warlock.Talents.Pyroclasm)
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Pyroclasm Talent Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warlock.Conflagrate && spellEffect.Outcome.Matches(core.OutcomeCrit) { // || spell == warlock.SearingPain
				warlock.PyroclasmAura.Activate(sim)
			}
		},
	})
}

func (warlock *Warlock) setupEradication() {
	castSpeedPercentBonus := float64(warlock.Talents.Eradication) * 6
	if warlock.Talents.Eradication == 3 {
		castSpeedPercentBonus += 2
	}
	warlock.EradicationAura = warlock.RegisterAura(core.Aura{
		Label:    "Eradication",
		ActionID: core.ActionID{SpellID: 64371},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(1 + castSpeedPercentBonus/100)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(1 / (1 + castSpeedPercentBonus/100))
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Eradication Talent Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warlock.Corruption {
				if sim.RandomFloat("Eradication") < 0.06 {
					warlock.EradicationAura.Activate(sim)
				}
			}
		},
	})
}

func (warlock *Warlock) setupShadowEmbrace() {
	warlock.ShadowEmbraceAura = warlock.RegisterAura(core.Aura{
		Label:     "Shadow Embrace",
		ActionID:  core.ActionID{SpellID: 32391},
		Duration:  time.Second * 12,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.PeriodicShadowDamageDealtMultiplier /= 1.0 + 0.01*float64(warlock.Talents.ShadowEmbrace)*float64(oldStacks)
			aura.Unit.PseudoStats.PeriodicShadowDamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.ShadowEmbrace)*float64(newStacks)
			// TO DO : Healing over time reduction part
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Shadow Embrace Talent Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warlock.ShadowBolt || spell == warlock.Haunt {
				if !warlock.ShadowEmbraceAura.IsActive() {
					warlock.ShadowEmbraceAura.Activate(sim)
				}
				warlock.ShadowEmbraceAura.AddStack(sim)
				warlock.ShadowEmbraceAura.Refresh(sim)
			}
		},
	})
}

func (warlock *Warlock) setupNightfall() {
	warlock.NightfallProcAura = warlock.RegisterAura(core.Aura{
		Label:    "Nightfall Shadow Trance",
		ActionID: core.ActionID{SpellID: 17941},
		Duration: time.Second * 10,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Check for an instant cast shadowbolt to disable aura
			if spell == warlock.ShadowBolt && spell.CurCast.CastTime == 0 {
				aura.Deactivate(sim)
			}
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Nightfall Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warlock.Corruption { // TODO: also works on drain life...
				if sim.RandomFloat("Nightfall") < 0.02*float64(warlock.Talents.Nightfall)+
					0.04*core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfCorruption), 1, 0) {
					warlock.NightfallProcAura.Activate(sim)
				}
			}
		},
	})
}

func (warlock *Warlock) applyNightfall(cast *core.Cast) {
	if warlock.NightfallProcAura.IsActive() {
		cast.CastTime = 0
	}
}

func (warlock *Warlock) setupMoltenCore() {
	warlock.MoltenCoreAura = warlock.RegisterAura(core.Aura{
		Label:     "Molten Core Proc Aura",
		ActionID:  core.ActionID{SpellID: 71165},
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == warlock.Incinerate || spell == warlock.SoulFire {
				aura.RemoveStack(sim)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.Incinerate.DamageMultiplier *= 1 + 0.06*float64(warlock.Talents.MoltenCore)
			warlock.SoulFire.DamageMultiplier *= 1 + 0.06*float64(warlock.Talents.MoltenCore)
			warlock.SoulFire.BonusCritRating += 5 * float64(warlock.Talents.MoltenCore) * core.CritRatingPerCritChance
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.Incinerate.DamageMultiplier /= 1 + 0.06*float64(warlock.Talents.MoltenCore)
			warlock.SoulFire.DamageMultiplier /= 1 + 0.06*float64(warlock.Talents.MoltenCore)
			warlock.SoulFire.BonusCritRating -= 5 * float64(warlock.Talents.MoltenCore) * core.CritRatingPerCritChance
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Molten Core Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warlock.Corruption {
				if sim.RandomFloat("Molten Core") < 0.04*float64(warlock.Talents.MoltenCore) {
					warlock.MoltenCoreAura.Activate(sim)
					warlock.MoltenCoreAura.SetStacks(sim, 3)
				}
			}
		},
	})
}

func (warlock *Warlock) setupBackdraft() {
	warlock.BackdraftAura = warlock.RegisterAura(core.Aura{
		Label:     "Backdraft Proc Aura",
		ActionID:  core.ActionID{SpellID: 54277},
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == warlock.Incinerate || spell == warlock.SoulFire || spell == warlock.ShadowBolt ||
				spell == warlock.ChaosBolt || spell == warlock.Immolate {
				aura.RemoveStack(sim)
			}
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Backdraft Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warlock.Conflagrate {
				warlock.BackdraftAura.Activate(sim)
				warlock.BackdraftAura.SetStacks(sim, 3)
			}
		},
	})
}

func (warlock *Warlock) backdraftModifier() float64 {
	castTimeModifier := 1.0
	if warlock.BackdraftAura.IsActive() {
		castTimeModifier *= (1.0 - 0.1*float64(warlock.Talents.Backdraft))
	}
	return castTimeModifier
}


