package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) ApplyTalents() {
	// warlock.PseudoStats.DamageDealtMultiplier = 10
	// warlock.PseudoStats.BonusSpellCritRating = 100 * core.CritRatingPerCritChance

	// Demonic Embrace
	if warlock.Talents.DemonicEmbrace > 0 {
		warlock.MultiplyStat(stats.Stamina, 1.01+(float64(warlock.Talents.DemonicEmbrace)*0.03))
	}

	// Molten Skin
	warlock.PseudoStats.DamageTakenMultiplier *= 1. - 0.02*float64(warlock.Talents.MoltenSkin)

	// Malediction
	maledictionMultiplier := 1. + 0.01*float64(warlock.Talents.Malediction)
	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= maledictionMultiplier
	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= maledictionMultiplier
	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] *= maledictionMultiplier
	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= maledictionMultiplier
	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= maledictionMultiplier

	warlock.setupDemonicPact()

	// Suppression (Add 1% hit per point)
	warlock.AddStat(stats.SpellHit, float64(warlock.Talents.Suppression)*core.SpellHitRatingPerHitChance)

	// Backlash (Add 1% crit per point)
	warlock.AddStat(stats.SpellCrit, float64(warlock.Talents.Backlash)*core.CritRatingPerCritChance)

	if warlock.Talents.DeathsEmbrace > 0 {
		warlock.applyDeathsEmbrace()
	}

	// Fel Vitality
	if warlock.Talents.FelVitality > 0 {
		bonus := 1.0 + 0.01*float64(warlock.Talents.FelVitality)
		warlock.MultiplyStat(stats.Mana, bonus)
		warlock.MultiplyStat(stats.Health, bonus)
	}

	// Demonic Tactics, applies even without pet out
	if warlock.Talents.DemonicTactics > 0 {
		warlock.AddStats(stats.Stats{
			stats.MeleeCrit: float64(warlock.Talents.DemonicTactics) * 2 * core.CritRatingPerCritChance,
			stats.SpellCrit: float64(warlock.Talents.DemonicTactics) * 2 * core.CritRatingPerCritChance,
		})
	}

	if warlock.Talents.Nightfall > 0 || warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfCorruption) {
		warlock.setupNightfall()
	}
	if warlock.Talents.EverlastingAffliction > 0 {
		warlock.setupEverlastingAffliction()
	}

	if warlock.Talents.ShadowEmbrace > 0 {
		warlock.setupShadowEmbrace()
	}

	if warlock.Talents.Eradication > 0 {
		warlock.setupEradication()
	}

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

	if warlock.Talents.ImprovedSoulLeech > 0 {
		warlock.setupImprovedSoulLeech()
	}

	if warlock.Talents.EmpoweredImp > 0 && warlock.Options.Summon == proto.Warlock_Options_Imp {
		warlock.Pet.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.1*float64(warlock.Talents.EmpoweredImp)
		warlock.setupEmpoweredImp()
	}

	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap) {
		warlock.registerGlyphOfLifeTapAura()
	}
}

func (warlock *Warlock) applyDeathsEmbrace() {
	multiplier := 1.0 + 0.04*float64(warlock.Talents.DeathsEmbrace)

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int) {
			if isExecute == 35 {
				warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= multiplier
			}
		})
	})
}

func (warlock *Warlock) applyWeaponImbue() {
	if warlock.Options.WeaponImbue == proto.Warlock_Options_GrandFirestone {
		warlock.AddStat(stats.SpellCrit, 49*(1+1.5*float64(warlock.Talents.MasterConjuror)))
	}
	if warlock.Options.WeaponImbue == proto.Warlock_Options_GrandSpellstone {
		warlock.AddStat(stats.SpellHaste, 60*(1+1.5*float64(warlock.Talents.MasterConjuror)))
	}
}

func (warlock *Warlock) registerGlyphOfLifeTapAura() {
	statDep := warlock.NewDynamicStatDependency(stats.Spirit, stats.SpellPower, 0.2)
	warlock.GlyphOfLifeTapAura = warlock.RegisterAura(core.Aura{
		Label:    "Glyph Of LifeTap Aura",
		ActionID: core.ActionID{SpellID: 63321},
		Duration: time.Second * 40,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.EnableDynamicStatDep(sim, statDep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.DisableDynamicStatDep(sim, statDep)
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeCrit) {
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeCrit) {
				warlock.EmpoweredImpAura.Activate(sim)
				warlock.EmpoweredImpAura.Refresh(sim)
			}
		},
	})
}

func (warlock *Warlock) setupDecimation() {
	decimationMod := 0.2 * float64(warlock.Talents.Decimation)
	warlock.DecimationAura = warlock.RegisterAura(core.Aura{
		Label:    "Decimation Proc Aura",
		ActionID: core.ActionID{SpellID: 63167},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.SoulFire.CastTimeMultiplier -= decimationMod
			warlock.SoulFire.DefaultCast.GCD = time.Duration(float64(warlock.SoulFire.DefaultCast.GCD) * (1 - decimationMod))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.SoulFire.CastTimeMultiplier += decimationMod
			warlock.SoulFire.DefaultCast.GCD = time.Duration(float64(warlock.SoulFire.DefaultCast.GCD) / (1 - decimationMod))
		},
	})

	decimation := warlock.RegisterAura(core.Aura{
		Label:    "Decimation Talent Hidden Aura",
		Duration: core.NeverExpires,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && (spell == warlock.ShadowBolt || spell == warlock.Incinerate || spell == warlock.SoulFire) {
				warlock.DecimationAura.Activate(sim)
			}
		},
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int) {
			if isExecute == 35 {
				decimation.Activate(sim)
			}
		})
	})
}

func (warlock *Warlock) setupPyroclasm() {
	pyroclasmDamageBonus := 1 + 0.02*float64(warlock.Talents.Pyroclasm)

	warlock.PyroclasmAura = warlock.RegisterAura(core.Aura{
		Label:    "Pyroclasm",
		ActionID: core.ActionID{SpellID: 63244},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= pyroclasmDamageBonus
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= pyroclasmDamageBonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= pyroclasmDamageBonus
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= pyroclasmDamageBonus
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Pyroclasm Talent Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warlock.Conflagrate && result.Outcome.Matches(core.OutcomeCrit) { // || spell == warlock.SearingPain
				warlock.PyroclasmAura.Activate(sim)
			}
		},
	})
}

func (warlock *Warlock) setupEradication() {
	castSpeedMultiplier := 1 + 0.06*float64(warlock.Talents.Eradication)
	if warlock.Talents.Eradication == 3 {
		castSpeedMultiplier += 0.02
	}
	warlock.EradicationAura = warlock.RegisterAura(core.Aura{
		Label:    "Eradication",
		ActionID: core.ActionID{SpellID: 64371},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(castSpeedMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(1 / castSpeedMultiplier)
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Eradication Talent Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warlock.Corruption {
				if sim.RandomFloat("Eradication") < 0.06 {
					warlock.EradicationAura.Activate(sim)
				}
			}
		},
	})
}

func (warlock *Warlock) ShadowEmbraceDebuffAura(target *core.Unit) *core.Aura {
	shadowEmbraceBonus := 0.01 * float64(warlock.Talents.ShadowEmbrace)

	return target.GetOrRegisterAura(core.Aura{
		Label:     "Shadow Embrace-" + warlock.Label,
		ActionID:  core.ActionID{SpellID: 32391},
		Duration:  time.Second * 12,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			warlock.AttackTables[aura.Unit.UnitIndex].PeriodicShadowDamageTakenMultiplier /= 1.0 + shadowEmbraceBonus*float64(oldStacks)
			warlock.AttackTables[aura.Unit.UnitIndex].PeriodicShadowDamageTakenMultiplier *= 1.0 + shadowEmbraceBonus*float64(newStacks)
		},
	})
}

func (warlock *Warlock) setupShadowEmbrace() {
	shadowEmbraceAuras := warlock.NewEnemyAuraArray(warlock.ShadowEmbraceDebuffAura)

	warlock.RegisterAura(core.Aura{
		Label:    "Shadow Embrace Talent Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warlock.ShadowBolt || spell == warlock.Haunt {
				aura := shadowEmbraceAuras.Get(result.Target)
				aura.Activate(sim)
				aura.AddStack(sim)
			}
		},
	})
}

func (warlock *Warlock) setupNightfall() {
	nightfallProcChance := 0.02*float64(warlock.Talents.Nightfall) +
		0.04*core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfCorruption), 1, 0)

	warlock.NightfallProcAura = warlock.RegisterAura(core.Aura{
		Label:    "Nightfall Shadow Trance",
		ActionID: core.ActionID{SpellID: 17941},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.ShadowBolt.CastTimeMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.ShadowBolt.CastTimeMultiplier += 1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Check if the shadowbolt was instant cast and not a normal one
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
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warlock.Corruption { // TODO: also works on drain life...
				if sim.RandomFloat("Nightfall") < nightfallProcChance {
					warlock.NightfallProcAura.Activate(sim)
				}
			}
		},
	})
}

func (warlock *Warlock) setupMoltenCore() {
	castReduction := 0.1 * float64(warlock.Talents.MoltenCore)
	moltenCoreDamageBonus := 1 + 0.06*float64(warlock.Talents.MoltenCore)
	moltenCoreCritBonus := 5 * float64(warlock.Talents.MoltenCore) * core.CritRatingPerCritChance

	warlock.MoltenCoreAura = warlock.RegisterAura(core.Aura{
		Label:     "Molten Core Proc Aura",
		ActionID:  core.ActionID{SpellID: 71165},
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.Incinerate.DamageMultiplier *= moltenCoreDamageBonus
			warlock.Incinerate.CastTimeMultiplier -= castReduction
			warlock.Incinerate.DefaultCast.GCD = time.Duration(float64(warlock.Incinerate.DefaultCast.GCD) * (1 - castReduction))
			warlock.SoulFire.DamageMultiplier *= moltenCoreDamageBonus
			warlock.SoulFire.BonusCritRating += moltenCoreCritBonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.Incinerate.DamageMultiplier /= moltenCoreDamageBonus
			warlock.Incinerate.CastTimeMultiplier += castReduction
			warlock.Incinerate.DefaultCast.GCD = time.Duration(float64(warlock.Incinerate.DefaultCast.GCD) / (1 - castReduction))
			warlock.SoulFire.DamageMultiplier /= moltenCoreDamageBonus
			warlock.SoulFire.BonusCritRating -= moltenCoreCritBonus
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == warlock.Incinerate || spell == warlock.SoulFire {
				aura.RemoveStack(sim)
			}
		},
	})

	warlock.RegisterAura(core.Aura{
		Label: "Molten Core Hidden Aura",
		// ActionID: core.ActionID{SpellID: 47247},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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
	castTimeModifier := 0.1 * float64(warlock.Talents.Backdraft)
	var affectedSpells []*core.Spell

	warlock.BackdraftAura = warlock.RegisterAura(core.Aura{
		Label:     "Backdraft Proc Aura",
		ActionID:  core.ActionID{SpellID: 54277},
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice([]*core.Spell{
				warlock.Incinerate,
				warlock.SoulFire,
				warlock.ShadowBolt,
				warlock.ChaosBolt,
				warlock.Immolate,
			}, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, destroSpell := range affectedSpells {
				destroSpell.CastTimeMultiplier -= castTimeModifier
				destroSpell.DefaultCast.GCD = time.Duration(float64(destroSpell.DefaultCast.GCD) * (1 - castTimeModifier))
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, destroSpell := range affectedSpells {
				destroSpell.CastTimeMultiplier += castTimeModifier
				destroSpell.DefaultCast.GCD = time.Duration(float64(destroSpell.DefaultCast.GCD) / (1 - castTimeModifier))
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			for _, destroSpell := range affectedSpells {
				if spell == destroSpell {
					aura.RemoveStack(sim)
					return
				}
			}
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Backdraft Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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

func (warlock *Warlock) everlastingAfflictionRefresh(sim *core.Simulation, target *core.Unit) {
	procChance := 0.2 * float64(warlock.Talents.EverlastingAffliction)

	if warlock.Corruption.Dot(target).IsActive() && sim.Proc(procChance, "EverlastingAffliction") {
		warlock.Corruption.Dot(target).Rollover(sim)
	}
}

func (warlock *Warlock) setupEverlastingAffliction() {
	warlock.RegisterAura(core.Aura{
		Label:    "Everlasting Affliction Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// TODO: also works on drain life...
			// drain soul is handled inside the spell itself
			if (spell != warlock.ShadowBolt && spell != warlock.Haunt) || !result.Landed() {
				return
			}

			warlock.everlastingAfflictionRefresh(sim, result.Target)
		},
	})
}

func (warlock *Warlock) setupImprovedSoulLeech() {
	soulLeechProcChance := 0.1 * float64(warlock.Talents.SoulLeech)
	improvedSoulLeechProcChance := float64(warlock.Talents.ImprovedSoulLeech) / 2.
	actionID := core.ActionID{SpellID: 54118}
	improvedSoulLeechManaMetric := warlock.NewManaMetrics(actionID)
	improvedSoulLeechPetManaMetric := warlock.Pets[0].GetCharacter().NewManaMetrics(actionID)

	replSrc := warlock.Env.Raid.NewReplenishmentSource(core.ActionID{SpellID: 54118})

	warlock.RegisterAura(core.Aura{
		Label:    "Improved Soul Leech Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warlock.Conflagrate || spell == warlock.ShadowBolt || spell == warlock.ChaosBolt || spell == warlock.SoulFire || spell == warlock.Incinerate {
				if sim.RandomFloat("SoulLeech") < soulLeechProcChance {
					warlock.AddMana(sim, warlock.MaxMana()*float64(warlock.Talents.ImprovedSoulLeech)/100, improvedSoulLeechManaMetric)
					warlock.Pets[0].GetCharacter().AddMana(sim, warlock.Pets[0].GetCharacter().MaxMana()*float64(warlock.Talents.ImprovedSoulLeech)/100, improvedSoulLeechPetManaMetric)
					if warlock.Talents.ImprovedSoulLeech == 2 || sim.RandomFloat("ImprovedSoulLeech") < improvedSoulLeechProcChance {
						warlock.Env.Raid.ProcReplenishment(sim, replSrc)
					}
				}
			}
		},
	})
}

func (warlock *Warlock) setupDemonicPact() {
	if warlock.Talents.DemonicPact == 0 {
		return
	}

	demonicPactMultiplier := 0.02 * float64(warlock.Talents.DemonicPact)
	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1. + demonicPactMultiplier
	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1. + demonicPactMultiplier
	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] *= 1. + demonicPactMultiplier
	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= 1. + demonicPactMultiplier
	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= 1. + demonicPactMultiplier

	if warlock.Options.Summon == proto.Warlock_Options_NoSummon {
		return
	}

	icd := core.Cooldown{
		Timer:    warlock.NewTimer(),
		Duration: time.Second * 5,
	}

	var demonicPactAuras [25]*core.Aura
	for _, party := range warlock.Party.Raid.Parties {
		for _, player := range party.Players {
			demonicPactAuras[player.GetCharacter().Index] = core.DemonicPactAura(player.GetCharacter())
		}
	}
	warlock.DemonicPactAura = demonicPactAuras[warlock.Index]

	warlock.Pets[0].GetCharacter().RegisterAura(core.Aura{
		Label:    "Demonic Pact Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() || !icd.IsReady(sim) {
				return
			}

			icd.Use(sim)
			newSPBonus := warlock.GetStat(stats.SpellPower) * demonicPactMultiplier

			shouldRefresh := !warlock.DemonicPactAura.IsActive() ||
				warlock.DemonicPactAura.RemainingDuration(sim) < time.Second*10 ||
				newSPBonus > warlock.DemonicPactAura.ExclusiveEffects[0].Priority

			if shouldRefresh {
				for _, dpAura := range demonicPactAuras {
					if dpAura != nil {
						dpAura.ExclusiveEffects[0].SetPriority(sim, newSPBonus)
						dpAura.Activate(sim)
					}
				}
			}
		},
	})
}
