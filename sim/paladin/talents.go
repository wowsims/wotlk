package paladin

import (
	"strconv"
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

	// if paladin.Talents.ShieldSpecialization > 0 {
	// 	bonus := 1 + 0.1*float64(paladin.Talents.ShieldSpecialization)
	// 	paladin.AddStatDependency(stats.StatDependency{
	// 		SourceStat:   stats.BlockValue,
	// 		ModifiedStat: stats.BlockValue,
	// 		Modifier: func(bv float64, _ float64) float64 {
	// 			return bv * bonus
	// 		},
	// 	})
	// }

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

	paladin.MultiplyStat(stats.BlockValue, 1.0+0.10*float64(paladin.Talents.Redoubt))

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

	actionID := core.ActionID{SpellID: 31854}
	damageReduction := 1.0 - 0.06*float64(paladin.Talents.ArdentDefender)

	procAura := paladin.RegisterAura(core.Aura{
		Label:    "Ardent Defender",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= damageReduction
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= damageReduction
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    "Ardent Defender Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if aura.Unit.CurrentHealthPercent() < 0.35 {
				procAura.Activate(sim)
			}
		},
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
					switch unit.MobType {
					case proto.MobType_MobTypeHumanoid, proto.MobType_MobTypeDemon, proto.MobType_MobTypeUndead, proto.MobType_MobTypeElemental:
						paladin.AttackTables[unit.UnitIndex].DamageDealtMultiplier *= 1 + (0.02 * float64(paladin.Talents.Crusade))
					default:
						paladin.AttackTables[unit.UnitIndex].DamageDealtMultiplier *= 1 + (0.01 * float64(paladin.Talents.Crusade))
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
			debuffAura := core.HeartoftheCrusaderDebuff(result.Target, float64(paladin.Talents.HeartOfTheCrusader))
			debuffAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) applyArtOfWar() {
	if paladin.Talents.TheArtOfWar == 0 {
		return
	}

	paladin.ArtOfWarInstantCast = paladin.RegisterAura(core.Aura{
		Label:    "Art Of War",
		ActionID: core.ActionID{SpellID: 53488},
		Duration: time.Second * 15,
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

	procSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 31878},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, _ *core.Spell) {
			if paladin.JowiseManaMetrics == nil {
				paladin.JowiseManaMetrics = paladin.NewManaMetrics(core.ActionID{SpellID: 31878})
			}
			paladin.AddMana(sim, paladin.BaseMana*0.25, paladin.JowiseManaMetrics, false)
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

			if paladin.Talents.JudgementsOfTheWise == 3 {
				procSpell.Cast(sim, nil)
			} else {
				if sim.RandomFloat("judgements of the wise") > (0.33)*float64(paladin.Talents.JudgementsOfTheWise) {
					return
				}
				procSpell.Cast(sim, nil)
			}
		},
	})
}

func (paladin *Paladin) makeRighteousVengeanceDot(target *core.Unit) *core.Dot {
	canCrit := paladin.HasTuralyonsOrLiadrinsBattlegear2Pc

	return core.NewDot(core.Dot{
		Spell: paladin.RighteousVengeanceSpell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Righteous Vengeance (DoT) - " + strconv.Itoa(int(paladin.Index)) + " - " + strconv.Itoa(int(target.Index)),
			ActionID: paladin.RighteousVengeanceSpell.ActionID,
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				paladin.RighteousVengeanceDamage[target.Index] = 0.0
				paladin.RighteousVengeancePools[target.Index] = 0.0
			},
		}),
		NumberOfTicks: 4,
		TickLength:    time.Second * 2,

		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			baseDmg := paladin.RighteousVengeanceDamage[target.Index]
			paladin.RighteousVengeancePools[target.Index] -= baseDmg

			if canCrit {
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDmg, dot.Spell.OutcomeMeleeSpecialCritOnly)
			} else {
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDmg, dot.Spell.OutcomeAlwaysHit)
			}
		},
	})
}

func (paladin *Paladin) registerRighteousVengeanceSpell() {
	dotActionID := core.ActionID{SpellID: 61840} // Righteous Vengeance

	paladin.RighteousVengeanceSpell = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    dotActionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreModifiers,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,
	})
}

func (paladin *Paladin) applyRighteousVengeance() {
	// Righteous Vengeance is a MAGIC debuff that pools 10/20/30% crit damage from Crusader Strike, Divine Storm, and Judgements.
	// It drains the pool every 2 seconds at a rate of 1/4 of the pool size.
	// And then deals that 1/4 as PHYSICAL damage.
	// TODO: Can crit with certain set bonuses.
	if paladin.Talents.RighteousVengeance == 0 {
		return
	}

	paladin.RegisterAura(core.Aura{
		Label:    "Righteous Vengeance",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() || !result.Landed() {
				return
			}

			if spell.SpellID == paladin.CrusaderStrike.SpellID || spell.SpellID == paladin.DivineStorm.SpellID || spell.Flags.Matches(SpellFlagSecondaryJudgement) {
				paladin.RighteousVengeancePools[result.Target.Index] += result.Damage * (0.10 * float64(paladin.Talents.RighteousVengeance))
				paladin.RighteousVengeanceDamage[result.Target.Index] = paladin.RighteousVengeancePools[result.Target.Index] / 4

				if !paladin.RighteousVengeanceDots[result.Target.Index].IsActive() {
					paladin.RighteousVengeanceDots[result.Target.Index].Apply(sim)
				}

				paladin.RighteousVengeanceDots[result.Target.Index].Refresh(sim)
			}
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
