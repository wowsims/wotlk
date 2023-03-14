package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO:
// Sanctified Wrath (Damage penetration, questions over affected stats)

func (paladin *Paladin) ApplyTalents() {
	paladin.AddStat(stats.MeleeCrit, float64(paladin.Talents.Conviction)*core.CritRatingPerCritChance)
	paladin.AddStat(stats.SpellCrit, float64(paladin.Talents.Conviction)*core.CritRatingPerCritChance)
	paladin.AddStat(stats.MeleeCrit, float64(paladin.Talents.SanctityOfBattle)*core.CritRatingPerCritChance)
	paladin.AddStat(stats.SpellCrit, float64(paladin.Talents.SanctityOfBattle)*core.CritRatingPerCritChance)

	paladin.PseudoStats.BaseParry += 0.01 * float64(paladin.Talents.Deflection)
	paladin.PseudoStats.BaseDodge += 0.01 * float64(paladin.Talents.Anticipation)

	paladin.AddStat(stats.Armor, paladin.Equip.Stats()[stats.Armor]*0.02*float64(paladin.Talents.Toughness))

	if paladin.Talents.DivineStrength > 0 {
		paladin.MultiplyStat(stats.Strength, 1.0+0.03*float64(paladin.Talents.DivineStrength))
	}
	if paladin.Talents.DivineIntellect > 0 {
		paladin.MultiplyStat(stats.Intellect, 1.0+0.02*float64(paladin.Talents.DivineIntellect))
	}

	if paladin.Talents.SheathOfLight > 0 {
		// doesn't implement HOT
		percentage := 0.10 * float64(paladin.Talents.SheathOfLight)
		paladin.AddStatDependency(stats.AttackPower, stats.SpellPower, percentage)
	}

	if paladin.Talents.TouchedByTheLight > 0 {
		percentage := 0.20 * float64(paladin.Talents.TouchedByTheLight)
		paladin.AddStatDependency(stats.Strength, stats.SpellPower, percentage)
	}

	if paladin.Talents.SacredDuty > 0 {
		paladin.MultiplyStat(stats.Stamina, 1.0+0.02*float64(paladin.Talents.SacredDuty))
	}

	if paladin.Talents.CombatExpertise > 0 {
		paladin.AddStat(stats.Expertise, core.ExpertisePerQuarterPercentReduction*2*float64(paladin.Talents.CombatExpertise))
		paladin.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*2*float64(paladin.Talents.CombatExpertise))
		paladin.AddStat(stats.SpellCrit, core.CritRatingPerCritChance*2*float64(paladin.Talents.CombatExpertise))
		paladin.MultiplyStat(stats.Stamina, 1.0+0.02*float64(paladin.Talents.CombatExpertise))
	}

	if paladin.Talents.ShieldOfTheTemplar > 0 {
		paladin.PseudoStats.DamageTakenMultiplier *= 1 - 0.01*float64(paladin.Talents.ShieldOfTheTemplar)
	}

	paladin.applyRedoubt()
	paladin.applyReckoning()
	paladin.applyArdentDefender()
	paladin.applyCrusade()
	paladin.applyWeaponSpecialization()
	paladin.applyVengeance()
	paladin.applyHeartOfTheCrusader()
	paladin.applyArtOfWar()
	paladin.applyJudgmentsOfTheWise()
	paladin.applyRighteousVengeance()
	paladin.applyMinorGlyphOfSenseUndead()
	paladin.applyGuardedByTheLight()
}

func (paladin *Paladin) getTalentSealsOfThePureBonus() float64 {
	return 0.03 * float64(paladin.Talents.SealsOfThePure)
}

func (paladin *Paladin) getTalentTwoHandedWeaponSpecializationBonus() float64 {
	return 0.02 * float64(paladin.Talents.TwoHandedWeaponSpecialization)
}

func (paladin *Paladin) getTalentSanctityOfBattleBonus() float64 {
	return 0.05 * float64(paladin.Talents.SanctityOfBattle)
}

func (paladin *Paladin) getTalentTheArtOfWarBonus() float64 {
	return 0.05 * float64(paladin.Talents.TheArtOfWar)
}

func (paladin *Paladin) getMajorGlyphSealOfRighteousnessBonus() float64 {
	return core.TernaryFloat64(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfSealOfRighteousness), .1, 0)
}

func (paladin *Paladin) getMajorGlyphOfExorcismBonus() float64 {
	return core.TernaryFloat64(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfExorcism), 0.20, 0)
}

func (paladin *Paladin) getMajorGlyphOfJudgementBonus() float64 {
	return core.TernaryFloat64(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfJudgement), 0.10, 0)
}

func (paladin *Paladin) applyMinorGlyphOfSenseUndead() {
	if !paladin.HasMinorGlyph(proto.PaladinMinorGlyph_GlyphOfSenseUndead) {
		return
	}

	var applied bool

	paladin.RegisterResetEffect(
		func(s *core.Simulation) {
			if !applied {
				for i := int32(0); i < paladin.Env.GetNumTargets(); i++ {
					unit := paladin.Env.GetTargetUnit(i)
					if unit.MobType == proto.MobType_MobTypeUndead {
						paladin.AttackTables[unit.UnitIndex].DamageDealtMultiplier *= 1.01
					}
				}
				applied = true
			}
		},
	)
}

func (paladin *Paladin) applyRedoubt() {
	if paladin.Talents.Redoubt == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 20132}

	paladin.PseudoStats.BlockValueMultiplier += 0.10 * float64(paladin.Talents.Redoubt)

	bonusBlockRating := 10 * core.BlockRatingPerBlockChance * float64(paladin.Talents.Redoubt)

	procAura := paladin.RegisterAura(core.Aura{
		Label:     "Redoubt Proc",
		ActionID:  actionID,
		Duration:  time.Second * 10,
		MaxStacks: 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.Block, bonusBlockRating)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.Block, -bonusBlockRating)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeBlock) {
				aura.RemoveStack(sim)
			}
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    "Redoubt",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				if sim.RandomFloat("Redoubt") < 0.1 {
					procAura.Activate(sim)
					procAura.SetStacks(sim, 5)
				}
			}
		},
	})
}

func (paladin *Paladin) applyReckoning() {
	if paladin.Talents.Reckoning == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 20182}
	procChance := 0.02 * float64(paladin.Talents.Reckoning)

	var reckoningSpell *core.Spell

	procAura := paladin.RegisterAura(core.Aura{
		Label:     "Reckoning Proc",
		ActionID:  actionID,
		Duration:  time.Second * 8,
		MaxStacks: 4,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			reckoningSpell = paladin.GetOrRegisterSpell(core.SpellConfig{
				ActionID:         actionID,
				SpellSchool:      core.SpellSchoolPhysical,
				ProcMask:         core.ProcMaskMeleeMH,
				Flags:            core.SpellFlagMeleeMetrics,
				CritMultiplier:   paladin.MeleeCritMultiplier(),
				ThreatMultiplier: 1,
				DamageMultiplier: 1,
				ApplyEffects:     paladin.AutoAttacks.MHConfig.ApplyEffects,
			})
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == paladin.AutoAttacks.MHAuto {
				reckoningSpell.Cast(sim, result.Target)
			}
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    "Reckoning",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && sim.RandomFloat("Reckoning") < procChance {
				procAura.Activate(sim)
				procAura.SetStacks(sim, 4)
			}
		},
	})
}

func (paladin *Paladin) applyArdentDefender() {
	if paladin.Talents.ArdentDefender == 0 {
		return
	}

	var ardentDamageReduction float64
	switch paladin.Talents.ArdentDefender {
	case 3:
		ardentDamageReduction = 0.07
	case 2:
		ardentDamageReduction = 0.13
	case 1:
		ardentDamageReduction = 0.20
	}

	// TBD? Buff to mark time spent fully below 35% and attribute absorbs
	// rangeAura := paladin.RegisterAura(core.Aura{
	// Label:    "Ardent Defender (Active)",
	// ActionID: core.ActionID{SpellID: 31852},
	// Duration: core.NeverExpires,
	// })

	// paladin.RegisterAura(core.Aura{
	// Label:    "Ardent Defender Talent",
	// Duration: core.NeverExpires,
	// OnReset: func(aura *core.Aura, sim *core.Simulation) {
	// aura.Activate(sim)
	// },
	// OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// if aura.Unit.CurrentHealthPercent() < 0.35 {
	// procAura.Activate(sim)
	// }
	// },
	// })

	// Debuff to show that AD has procced
	procAura := paladin.RegisterAura(core.Aura{
		Label:    "Ardent Defender",
		ActionID: core.ActionID{SpellID: 66233},
		Duration: time.Second * 120.0,
	})

	// Spell to heal you when AD has procced; fire this before fatal damage so that a Death is not detected
	procHeal := paladin.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 66233},
		SpellSchool:      core.SpellSchoolHoly,
		ProcMask:         core.ProcMaskSpellHealing,
		CritMultiplier:   1, // Assuming this can't really crit?
		ThreatMultiplier: 0.25,
		DamageMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// 540 defense (+140) yields the full heal amount
			ardentHealAmount := core.MaxFloat(1.0, float64(paladin.GetStat(stats.Defense))/core.DefenseRatingPerDefense/140.0) * 0.10 * float64(paladin.Talents.ArdentDefender)
			spell.CalcAndDealHealing(sim, &paladin.Unit, ardentHealAmount*paladin.MaxHealth(), spell.OutcomeHealingCrit)
		},
	})

	// >= 0.35, no effect
	// < 0.35, pro-rated DR
	// =< 0, proc death save
	paladin.AddDynamicDamageTakenModifier(func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
		incomingDamage := result.Damage
		if (paladin.CurrentHealth()-incomingDamage)/paladin.MaxHealth() <= 0.35 {
			//rangeAura.Activate(sim)
			result.Damage -= (paladin.MaxHealth()*0.35 - (paladin.CurrentHealth() - incomingDamage)) * ardentDamageReduction
			if sim.Log != nil {
				paladin.Log(sim, "Ardent Defender reduced damage by %d", int32(incomingDamage-result.Damage))
			}

			// Now check death save, based on the reduced damage
			if (result.Damage >= paladin.CurrentHealth()) && !procAura.IsActive() {
				result.Damage = paladin.CurrentHealth()
				procHeal.Cast(sim, &paladin.Unit)
				procAura.Activate(sim)
			}
		}
		// TODO: Metrics, attribute reduced damage as absorption
	})
}

// Because Crusade modifies unit specific attack tables, it must be applied at start of sim.
func (paladin *Paladin) applyCrusade() {
	if paladin.Talents.Crusade == 0 {
		return
	}

	var applied bool

	paladin.RegisterResetEffect(
		func(s *core.Simulation) {
			if !applied {
				for i := int32(0); i < paladin.Env.GetNumTargets(); i++ {
					unit := paladin.Env.GetTargetUnit(i)
					crusadeMod := 1.0 + (0.01 * float64(paladin.Talents.Crusade))
					switch unit.MobType {
					case proto.MobType_MobTypeHumanoid, proto.MobType_MobTypeDemon, proto.MobType_MobTypeUndead, proto.MobType_MobTypeElemental:
						paladin.AttackTables[unit.UnitIndex].DamageDealtMultiplier *= crusadeMod * crusadeMod
					default:
						paladin.AttackTables[unit.UnitIndex].DamageDealtMultiplier *= crusadeMod
					}
				}
				applied = true
			}
		},
	)
}

// Prior to WOTLK, behavior was to double dip.
func (paladin *Paladin) MeleeCritMultiplier() float64 {
	// return paladin.Character.MeleeCritMultiplier(paladin.crusadeMultiplier(), 0)
	return paladin.DefaultMeleeCritMultiplier()
}
func (paladin *Paladin) SpellCritMultiplier() float64 {
	// return paladin.Character.SpellCritMultiplier(paladin.crusadeMultiplier(), 0)
	return paladin.DefaultSpellCritMultiplier()
}

// Affects all physical damage or spells that can be rolled as physical
// It affects white, Windfury, Crusader Strike, Seals, and Judgement of Command / Blood
func (paladin *Paladin) applyWeaponSpecialization() {
	// This impacts Crusader Strike, Melee Attacks, WF attacks
	// Seals + Judgements need to be implemented separately
	mhWeapon := paladin.GetMHWeapon()

	if mhWeapon == nil {
		return
	}

	switch mhWeapon.HandType {
	case proto.HandType_HandTypeTwoHand:
		paladin.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + 0.02*float64(paladin.Talents.TwoHandedWeaponSpecialization)
	case proto.HandType_HandTypeOneHand, proto.HandType_HandTypeMainHand:
		if paladin.Talents.OneHandedWeaponSpecialization > 0 {
			// Talent points are 4%, 7%, 10%
			paladin.PseudoStats.DamageDealtMultiplier *= 1.01 + 0.03*float64(paladin.Talents.OneHandedWeaponSpecialization)
		}
	}
}

// I don't know if the new stack of vengeance applies to the crit that triggered it or not
// Need to check this
func (paladin *Paladin) applyVengeance() {
	if paladin.Talents.Vengeance == 0 {
		return
	}

	bonusPerStack := 0.01 * float64(paladin.Talents.Vengeance)
	procAura := paladin.RegisterAura(core.Aura{
		Label:     "Vengeance Proc",
		ActionID:  core.ActionID{SpellID: 20057},
		Duration:  time.Second * 30,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] /= 1 + (bonusPerStack * float64(oldStacks))
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1 + (bonusPerStack * float64(oldStacks))

			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= 1 + (bonusPerStack * float64(newStacks))
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + (bonusPerStack * float64(newStacks))
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    "Vengeance",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeCrit) {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			}
		},
	})
}

func (paladin *Paladin) applyHeartOfTheCrusader() {
	if paladin.Talents.HeartOfTheCrusader == 0 {
		return
	}

	paladin.RegisterAura(core.Aura{
		Label:    "Heart of the Crusader",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagSecondaryJudgement) {
				return
			}
			debuffAura := core.HeartOfTheCrusaderDebuff(result.Target, float64(paladin.Talents.HeartOfTheCrusader))
			debuffAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) applyArtOfWar() {
	if paladin.Talents.TheArtOfWar == 0 {
		return
	}

	castTimeReduction := 0.5 * float64(paladin.Talents.TheArtOfWar)
	paladin.ArtOfWarInstantCast = paladin.RegisterAura(core.Aura{
		Label:    "Art Of War",
		ActionID: core.ActionID{SpellID: 53488},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.Exorcism.CastTimeMultiplier -= castTimeReduction
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.Exorcism.CastTimeMultiplier += castTimeReduction
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == paladin.Exorcism {
				aura.Deactivate(sim)
			}
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    "The Art of War",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.IsMelee() && !spell.Flags.Matches(SpellFlagSecondaryJudgement) {
				return
			}

			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			paladin.ArtOfWarInstantCast.Activate(sim)
		},
	})
}

func (paladin *Paladin) applyJudgmentsOfTheWise() {
	if paladin.Talents.JudgementsOfTheWise == 0 {
		return
	}

	procChance := float64(paladin.Talents.JudgementsOfTheWise) / 3
	paladin.JowiseManaMetrics = paladin.NewManaMetrics(core.ActionID{SpellID: 31878})
	replSrc := paladin.Env.Raid.NewReplenishmentSource(core.ActionID{SpellID: 31878})

	procSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 31878},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, _ *core.Spell) {
			paladin.AddMana(sim, paladin.BaseMana*0.25, paladin.JowiseManaMetrics)
			paladin.Env.Raid.ProcReplenishment(sim, replSrc)
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    "Judgements of the Wise",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagSecondaryJudgement) || !result.Landed() {
				return
			}

			if procChance == 1 || sim.RandomFloat("judgements of the wise") < procChance {
				procSpell.Cast(sim, nil)
			}
		},
	})
}

func (paladin *Paladin) applyRighteousVengeance() {
	// Righteous Vengeance is a MAGIC debuff that pools 10/20/30% crit damage from Crusader Strike, Divine Storm, and Judgements.
	// It drains the pool every 2 seconds at a rate of 1/4 of the pool size.
	// And then deals that 1/4 as PHYSICAL damage.
	if paladin.Talents.RighteousVengeance == 0 {
		return
	}

	dotActionID := core.ActionID{SpellID: 61840} // Righteous Vengeance

	rvDot := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    dotActionID.WithTag(2),
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreModifiers,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Righteous Vengeance",
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 2,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if paladin.HasTuralyonsOrLiadrinsBattlegear2Pc {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeMeleeSpecialCritOnly)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
				}
			},
		},
	})

	rvSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    dotActionID.WithTag(1),
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rvDot.Dot(target).ApplyOrReset(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    "Righteous Vengeance",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}
			if spell != paladin.CrusaderStrike && spell != paladin.DivineStorm && !spell.Flags.Matches(SpellFlagSecondaryJudgement) {
				return
			}

			dot := rvDot.Dot(result.Target)

			newDamage := result.Damage * (0.10 * float64(paladin.Talents.RighteousVengeance))
			outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)

			dot.SnapshotAttackerMultiplier = 1
			dot.SnapshotBaseDamage = (outstandingDamage + newDamage) / float64(dot.NumberOfTicks)
			rvSpell.Cast(sim, result.Target)
		},
	})
}

//nolint:unused
func (paladin *Paladin) applyFanaticism() {
	// TODO: Possibly implement as aura.
	if paladin.Talents.Fanaticism == 0 {
		return
	}

	paladin.PseudoStats.ThreatMultiplier *= 1 - 0.10*float64(paladin.Talents.Fanaticism)
}

func (paladin *Paladin) applyGuardedByTheLight() {
	if paladin.Talents.GuardedByTheLight == 0 {
		return
	}

	paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 1 - 0.03*float64(paladin.Talents.GuardedByTheLight)
	paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 1 - 0.03*float64(paladin.Talents.GuardedByTheLight)
	paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 1 - 0.03*float64(paladin.Talents.GuardedByTheLight)
	paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= 1 - 0.03*float64(paladin.Talents.GuardedByTheLight)
	paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 1 - 0.03*float64(paladin.Talents.GuardedByTheLight)
	paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 1 - 0.03*float64(paladin.Talents.GuardedByTheLight)

	paladin.RegisterAura(core.Aura{
		Label:    "Guarded By The Light",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if paladin.DivinePleaAura.IsActive() {
				paladin.DivinePleaAura.Refresh(sim)
			}
		},
	})
}
