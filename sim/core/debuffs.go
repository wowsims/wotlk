package core

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func applyDebuffEffects(target *Unit, targetIdx int, debuffs *proto.Debuffs, raid *proto.Raid) {
	if debuffs.Misery && targetIdx == 0 {
		MakePermanent(MiseryAura(target, 3))
	}

	if debuffs.JudgementOfWisdom && targetIdx == 0 {
		MakePermanent(JudgementOfWisdomAura(target))
	}
	if debuffs.JudgementOfLight && targetIdx == 0 {
		MakePermanent(JudgementOfLightAura(target))
	}

	if debuffs.CurseOfElements {
		MakePermanent(CurseOfElementsAura(target))
	}
	if debuffs.EbonPlaguebringer {
		MakePermanent(EbonPlaguebringerOrCryptFeverAura(nil, target, 2, 3, 3))
	}
	if debuffs.EarthAndMoon && targetIdx == 0 {
		MakePermanent(EarthAndMoonAura(target, 3))
	}

	if debuffs.ShadowMastery && targetIdx == 0 {
		MakePermanent(ShadowMasteryAura(target))
	}

	if debuffs.ImprovedScorch && targetIdx == 0 {
		MakePermanent(ImprovedScorchAura(target))
	}

	if debuffs.WintersChill && targetIdx == 0 {
		MakePermanent(WintersChillAura(target, 5))
	}

	if debuffs.BloodFrenzy && targetIdx < 4 {
		MakePermanent(BloodFrenzyAura(target, 2))
	}
	if debuffs.SavageCombat {
		MakePermanent(SavageCombatAura(target, 2))
	}

	if debuffs.GiftOfArthas {
		MakePermanent(GiftOfArthasAura(target))
	}

	if debuffs.SporeCloud {
		MakePermanent(SporeCloudAura(target))
	}

	if debuffs.Mangle && targetIdx == 0 {
		MakePermanent(MangleAura(target))
	} else if debuffs.Trauma && targetIdx == 0 {
		MakePermanent(TraumaAura(target, 2))
	} else if debuffs.Stampede && targetIdx == 0 {
		stampedeAura := StampedeAura(target)
		target.RegisterResetEffect(func(sim *Simulation) {
			StartPeriodicAction(sim, PeriodicActionOptions{
				Period: time.Second * 60,
				OnAction: func(sim *Simulation) {
					stampedeAura.Activate(sim)
				},
			})
		})
	}

	if debuffs.ExposeArmor && targetIdx == 0 {
		aura := ExposeArmorAura(target, false)
		ScheduledMajorArmorAura(aura, PeriodicActionOptions{
			Period:   time.Second * 3,
			NumTicks: 1,
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
			},
		}, raid)
	}

	if debuffs.SunderArmor && targetIdx == 0 {
		aura := SunderArmorAura(target)
		ScheduledMajorArmorAura(aura, PeriodicActionOptions{
			Period:          time.Millisecond * 1500,
			NumTicks:        5,
			TickImmediately: true,
			Priority:        ActionPriorityDOT, // High prio so it comes before actual warrior sunders.
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			},
		}, raid)
	}

	if debuffs.AcidSpit && targetIdx == 0 {
		aura := AcidSpitAura(target)
		ScheduledMajorArmorAura(aura, PeriodicActionOptions{
			Period:          time.Second * 10,
			NumTicks:        2,
			TickImmediately: true,
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			},
		}, raid)
	}

	if debuffs.CurseOfWeakness != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(CurseOfWeaknessAura(target, GetTristateValueInt32(debuffs.CurseOfWeakness, 1, 2)))
	}
	if debuffs.Sting && targetIdx == 0 {
		MakePermanent(StingAura(target))
	}

	if debuffs.FaerieFire != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(FaerieFireAura(target, GetTristateValueInt32(debuffs.FaerieFire, 0, 3)))
	}

	if debuffs.DemoralizingRoar != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingRoarAura(target, GetTristateValueInt32(debuffs.DemoralizingRoar, 0, 5)))
	}
	if debuffs.DemoralizingShout != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingShoutAura(target, 0, GetTristateValueInt32(debuffs.DemoralizingShout, 0, 5)))
	}
	if debuffs.Vindication && targetIdx == 0 {
		MakePermanent(VindicationAura(target))
	}
	if debuffs.DemoralizingScreech {
		MakePermanent(DemoralizingScreechAura(target))
	}

	// Atk spd reduction
	if debuffs.ThunderClap != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(ThunderClapAura(target, GetTristateValueInt32(debuffs.ThunderClap, 0, 3)))
	}
	if debuffs.FrostFever != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(FrostFeverAura(target, GetTristateValueInt32(debuffs.FrostFever, 0, 3)))
	}
	if debuffs.InfectedWounds && targetIdx == 0 {
		MakePermanent(InfectedWoundsAura(target, 3))
	}
	if debuffs.JudgementsOfTheJust && targetIdx == 0 {
		MakePermanent(JudgementsOfTheJustAura(target, 2))
	}

	// Miss
	if debuffs.InsectSwarm && targetIdx == 0 {
		MakePermanent(InsectSwarmAura(target))
	}
	if debuffs.ScorpidSting && targetIdx == 0 {
		MakePermanent(ScorpidStingAura(target))
	}

	if debuffs.TotemOfWrath {
		MakePermanent(TotemOfWrathDebuff(target))
	}

	if debuffs.MasterPoisoner {
		MakePermanent(MasterPoisonerDebuff(target, 3))
	}

	if debuffs.HeartOfTheCrusader && targetIdx == 0 {
		MakePermanent(HeartOfTheCrusaderDebuff(target, 3))
	}

	if debuffs.HuntersMark > 0 && targetIdx == 0 {
		points := int32(0)
		glyphed := false
		if debuffs.HuntersMark > 1 {
			points = 3
			if debuffs.HuntersMark > 2 {
				glyphed = true
			}
		}
		MakePermanent(HuntersMarkAura(target, points, glyphed))
	}
}

func ScheduledMajorArmorAura(aura *Aura, options PeriodicActionOptions, raid *proto.Raid) {
	// Individual rogue sim rotation option messes with these debuff options,
	// so it has to be handled separately.
	allRogues := RaidPlayersWithClass(raid, proto.Class_ClassRogue)
	singleExposeDelay := len(allRogues) == 1 &&
		allRogues[0].Spec.(*proto.Player_Rogue).Rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once

	if singleExposeDelay {
		target := aura.Unit
		exposeArmorAura := ExposeArmorAura(target, false)
		exposeArmorAura.ApplyOnExpire(func(_ *Aura, sim *Simulation) {
			aura.Duration = NeverExpires
			StartPeriodicAction(sim, options)
		})
	} else {
		aura.OnReset = func(aura *Aura, sim *Simulation) {
			aura.Duration = NeverExpires
			StartPeriodicAction(sim, options)
		}
	}
}

var JudgementOfWisdomAuraLabel = "Judgement of Wisdom"

func JudgementOfWisdomAura(target *Unit) *Aura {
	actionID := ActionID{SpellID: 53408}

	return target.GetOrRegisterAura(Aura{
		Label:    JudgementOfWisdomAuraLabel,
		ActionID: actionID,
		Duration: time.Second * 20,
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			unit := spell.Unit
			if !unit.HasManaBar() {
				return
			}

			if spell.ProcMask.Matches(ProcMaskEmpty) {
				return // Phantom spells (Romulo's, Lightning Capacitor, etc) don't proc JoW.
			}

			// Melee claim that wisdom can proc on misses.
			if !spell.ProcMask.Matches(ProcMaskMeleeOrRanged) && !result.Landed() {
				return
			}

			if spell.ProcMask.Matches(ProcMaskMeleeOrRanged) {
				if !unit.AutoAttacks.PPMProc(sim, 15, spell.ProcMask, "jow") {
					return
				}
			} else {
				// TODO: Figure out if spell proc rate is also different from TBC.
				if sim.RandomFloat("jow") <= 0.5 {
					return
				}
			}

			if unit.JowManaMetrics == nil {
				unit.JowManaMetrics = unit.NewManaMetrics(actionID)
			}
			// JoW returns 2% of base mana 50% of the time.
			unit.AddMana(sim, unit.BaseMana*0.02, unit.JowManaMetrics)
		},
	})
}

var JudgementOfLightAuraLabel = "Judgement of Light"

func JudgementOfLightAura(target *Unit) *Aura {
	actionID := ActionID{SpellID: 20271}

	return target.GetOrRegisterAura(Aura{
		Label:    JudgementOfLightAuraLabel,
		ActionID: actionID,
		Duration: time.Second * 20,
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if !spell.ProcMask.Matches(ProcMaskMelee) || !result.Landed() {
				return
			}
		},
	})
}

func CurseOfElementsAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Curse of Elements",
		ActionID: ActionID{SpellID: 47865},
		Duration: time.Minute * 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: -165, stats.FireResistance: -165, stats.FrostResistance: -165, stats.ShadowResistance: -165, stats.NatureResistance: -165})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: 165, stats.FireResistance: 165, stats.FrostResistance: 165, stats.ShadowResistance: 165, stats.NatureResistance: 165})
		},
	})
	spellDamageEffect(aura, 1.13)
	return aura
}

func EarthAndMoonAura(target *Unit, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Earth And Moon" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 48511},
		Duration: time.Second * 12,
	})
	spellDamageEffect(aura, []float64{1, 1.04, 1.09, 1.13}[points])
	return aura
}

func EbonPlaguebringerOrCryptFeverAura(caster *Character, target *Unit, epidemicPoints, cryptFeverPoints, ebonPlaguebringerPoints int32) *Aura {
	casterIndex := -1
	// On application, Crypt Fever and Ebon Plaguebringer trigger extra 'ghost' procs.
	var ghostSpell *Spell
	if caster != nil {
		casterIndex = int(caster.Index)
		ghostSpell = caster.RegisterSpell(SpellConfig{
			ActionID:    ActionID{SpellID: 52789},
			SpellSchool: SpellSchoolMagic,
			ProcMask:    ProcMaskSpellDamage,
			Flags:       SpellFlagNoLogs | SpellFlagNoMetrics | SpellFlagNoOnCastComplete | SpellFlagIgnoreModifiers,

			DamageMultiplier: 1,
			ThreatMultiplier: 0,

			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				// Just deal 0 damage as the "Harmful Spell" is implemented on spell damage
				spell.CalcAndDealDamage(sim, target, 0, spell.OutcomeAlwaysHit)
			},
		})
	}

	aura := target.GetOrRegisterAura(Aura{
		Label: "EbonPlaguebringer" + strconv.Itoa(casterIndex), // Support multiple DKs having their EP up
		// ActionID: ActionID{SpellID: 49632}, // Crypt Fever spellID if we ever care
		ActionID: ActionID{SpellID: 51161},
		Duration: time.Second * (15 + 3*time.Duration(epidemicPoints)),
		OnGain: func(aura *Aura, sim *Simulation) {
			if ghostSpell != nil {
				if cryptFeverPoints > 0 {
					ghostSpell.Cast(sim, aura.Unit)
				}
				if ebonPlaguebringerPoints > 0 {
					ghostSpell.Cast(sim, aura.Unit)
				}
			}
		},
	})

	diseaseMultiplier := 1.0 + 0.1*float64(cryptFeverPoints)
	aura.NewExclusiveEffect("diseaseDmg", false, ExclusiveEffect{
		Priority: diseaseMultiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.DiseaseDamageTakenMultiplier *= diseaseMultiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.DiseaseDamageTakenMultiplier /= diseaseMultiplier
		},
	})

	if ebonPlaguebringerPoints > 0 {
		spellDamageEffect(aura, []float64{1, 1.04, 1.09, 1.13}[ebonPlaguebringerPoints])
	}
	return aura
}

func spellDamageEffect(aura *Aura, multiplier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("spellDamage", false, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= multiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] /= multiplier
		},
	})
}

func BloodFrenzyAura(target *Unit, points int32) *Aura {
	return bloodFrenzySavageCombatAura(target, "Blood Frenzy", ActionID{SpellID: 29859}, points)
}
func SavageCombatAura(target *Unit, points int32) *Aura {
	return bloodFrenzySavageCombatAura(target, "Savage Combat", ActionID{SpellID: 58413}, points)
}

func bloodFrenzySavageCombatAura(target *Unit, label string, id ActionID, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    label + "-" + strconv.Itoa(int(points)),
		ActionID: id,
		// No fixed duration, lasts as long as the bleed that activates it.
		Duration: NeverExpires,
	})

	multiplier := 1 + 0.02*float64(points)
	aura.NewExclusiveEffect("PhysicalDmg", true, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= multiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= multiplier
		},
	})
	return aura
}

func GiftOfArthasAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Gift of Arthas",
		ActionID: ActionID{SpellID: 11374},
		Duration: time.Minute * 3,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusPhysicalDamageTaken += 8
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusPhysicalDamageTaken -= 8
		},
	})
}

func MangleAura(target *Unit) *Aura {
	return bleedDamageAura(target, Aura{
		Label:    "Mangle",
		ActionID: ActionID{SpellID: 33876},
		Duration: time.Minute,
	}, 1.3)
}

func TraumaAura(target *Unit, points int) *Aura {
	return bleedDamageAura(target, Aura{
		Label:    "Trauma",
		ActionID: ActionID{SpellID: 46855},
		Duration: time.Second * 60,
	}, 1+0.15*float64(points))
}

func StampedeAura(target *Unit) *Aura {
	return bleedDamageAura(target, Aura{
		Label:    "Stampede",
		ActionID: ActionID{SpellID: 57393},
		Duration: time.Second * 12,
	}, 1.25)
}

const BleedEffectCategory = "BleedDamage"

func bleedDamageAura(target *Unit, config Aura, multiplier float64) *Aura {
	aura := target.GetOrRegisterAura(config)
	aura.NewExclusiveEffect(BleedEffectCategory, true, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier *= multiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier /= multiplier
		},
	})
	return aura
}

func ShadowMasteryAura(target *Unit) *Aura {
	return majorSpellCritDebuffAura(target, "Shadow Mastery", ActionID{SpellID: 17800}, 5)
}

func ImprovedScorchAura(target *Unit) *Aura {
	return majorSpellCritDebuffAura(target, "Improved Scorch", ActionID{SpellID: 12873}, 5)
}

const SpellCritEffectCategory = "spellcritdebuff"

func WintersChillAura(target *Unit, startingStacks int32) *Aura {
	var effect *ExclusiveEffect
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Winter's Chill",
		ActionID:  ActionID{SpellID: 28593},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, startingStacks)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks, newStacks int32) {
			effect.SetPriority(sim, float64(newStacks))
		},
	})

	effect = aura.NewExclusiveEffect(SpellCritEffectCategory, true, ExclusiveEffect{
		Priority: 0,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusSpellCritRatingTaken += ee.Priority * CritRatingPerCritChance
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusSpellCritRatingTaken -= ee.Priority * CritRatingPerCritChance
		},
	})
	return aura
}

func majorSpellCritDebuffAura(target *Unit, label string, actionID ActionID, percent float64) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: actionID,
		Duration: time.Second * 30,
	})

	bonusSpellCrit := percent * CritRatingPerCritChance
	aura.NewExclusiveEffect(SpellCritEffectCategory, true, ExclusiveEffect{
		Priority: percent,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusSpellCritRatingTaken += bonusSpellCrit
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusSpellCritRatingTaken -= bonusSpellCrit
		},
	})
	return aura
}

func MiseryAura(target *Unit, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Misery",
		ActionID: ActionID{SpellID: 33198},
		Duration: time.Second * 24,
	})
	spellHitBonusEffect(aura, float64(points)*SpellHitRatingPerHitChance)
	return aura
}

var FaerieFireAuraTag = "Faerie Fire"

func FaerieFireAura(target *Unit, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Faerie Fire" + strconv.Itoa(int(points)),
		Tag:      FaerieFireAuraTag,
		ActionID: ActionID{SpellID: 770},
		Duration: time.Minute * 5,
	})

	minorArmorReductionEffect(aura, 0.05)

	if points > 0 {
		spellHitBonusEffect(aura, float64(points)*SpellHitRatingPerHitChance)
	}

	return aura
}

func spellHitBonusEffect(aura *Aura, spellHitBonus float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("BonusSpellHit", false, ExclusiveEffect{
		Priority: spellHitBonus,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusSpellHitRatingTaken += spellHitBonus
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusSpellHitRatingTaken -= spellHitBonus
		},
	})
}

var majorArmorReductionEffectCategory = "MajorArmorReduction"

func SunderArmorAura(target *Unit) *Aura {
	var effect *ExclusiveEffect
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Sunder Armor",
		ActionID:  ActionID{SpellID: 47467},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			effect.SetPriority(sim, 0.04*float64(newStacks))
		},
	})

	effect = aura.NewExclusiveEffect(majorArmorReductionEffectCategory, true, ExclusiveEffect{
		Priority: 0,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.ArmorMultiplier *= 1 - ee.Priority
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.ArmorMultiplier /= 1 - ee.Priority
		},
	})

	return aura
}

func AcidSpitAura(target *Unit) *Aura {
	var effect *ExclusiveEffect
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Acid Spit",
		ActionID:  ActionID{SpellID: 55754},
		Duration:  time.Second * 10,
		MaxStacks: 2,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			effect.SetPriority(sim, 0.1*float64(newStacks))
		},
	})

	effect = aura.NewExclusiveEffect(majorArmorReductionEffectCategory, true, ExclusiveEffect{
		Priority: 0,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.ArmorMultiplier *= 1 - ee.Priority
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.ArmorMultiplier /= 1 - ee.Priority
		},
	})

	return aura
}

func ExposeArmorAura(target *Unit, hasGlyph bool) *Aura {
	const armorReduction = 0.2
	aura := target.GetOrRegisterAura(Aura{
		Label:    "ExposeArmor",
		ActionID: ActionID{SpellID: 8647},
		Duration: time.Second * TernaryDuration(hasGlyph, 42, 30),
	})

	aura.NewExclusiveEffect(majorArmorReductionEffectCategory, true, ExclusiveEffect{
		Priority: armorReduction,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.ArmorMultiplier *= 1 - armorReduction
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.ArmorMultiplier /= 1 - armorReduction
		},
	})

	return aura
}

func CurseOfWeaknessAura(target *Unit, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Curse of Weakness" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 50511},
		Duration: time.Minute * 2,
	})
	minorArmorReductionEffect(aura, 0.05)
	apReductionEffect(aura, 478*(1+0.1*float64(points)))
	return aura
}

func StingAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Sting",
		ActionID: ActionID{SpellID: 56631},
		Duration: time.Second * 20,
	})
	minorArmorReductionEffect(aura, 0.05)
	return aura
}

func SporeCloudAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Spore Cloud",
		ActionID: ActionID{SpellID: 53598},
		Duration: time.Second * 9,
	})
	minorArmorReductionEffect(aura, 0.03)
	return aura
}

func minorArmorReductionEffect(aura *Aura, reduction float64) *ExclusiveEffect {
	multiplier := 1 - reduction
	return aura.NewExclusiveEffect("MinorArmorReduction", false, ExclusiveEffect{
		Priority: reduction,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.ArmorMultiplier *= multiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.ArmorMultiplier /= multiplier
		},
	})
}

var ShatteringThrowAuraTag = "ShatteringThrow"

var ShatteringThrowDuration = time.Second * 10

func ShatteringThrowAura(target *Unit) *Aura {
	armorReduction := 0.2

	return target.GetOrRegisterAura(Aura{
		Label:    "Shattering Throw",
		Tag:      ShatteringThrowAuraTag,
		ActionID: ActionID{SpellID: 64382},
		Duration: ShatteringThrowDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
		},
	})
}

const HuntersMarkAuraTag = "HuntersMark"

func HuntersMarkAura(target *Unit, points int32, glyphed bool) *Aura {
	bonus := 500.0 * (1 + 0.1*float64(points) + TernaryFloat64(glyphed, 0.2, 0))

	aura := target.GetOrRegisterAura(Aura{
		Label:    "HuntersMark-" + strconv.Itoa(int(bonus)),
		Tag:      HuntersMarkAuraTag,
		ActionID: ActionID{SpellID: 53338},
		Duration: NeverExpires,
	})

	aura.NewExclusiveEffect("HuntersMark", true, ExclusiveEffect{
		Priority: bonus,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusRangedAttackPowerTaken += bonus
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusRangedAttackPowerTaken -= bonus
		},
	})

	return aura
}

func DemoralizingRoarAura(target *Unit, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingRoar-" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 48560},
		Duration: time.Second * 30,
	})
	apReductionEffect(aura, 411*(1+0.08*float64(points)))
	return aura
}

func DemoralizingShoutAura(target *Unit, boomingVoicePts int32, impDemoShoutPts int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingShout-" + strconv.Itoa(int(impDemoShoutPts)),
		ActionID: ActionID{SpellID: 47437},
		Duration: time.Duration(float64(time.Second*30) * (1 + 0.1*float64(boomingVoicePts))),
	})
	apReductionEffect(aura, 411*(1+0.08*float64(impDemoShoutPts)))
	return aura
}

func VindicationAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Vindication",
		ActionID: ActionID{SpellID: 26016},
		Duration: time.Second * 10,
	})
	apReductionEffect(aura, 574)
	return aura
}

func DemoralizingScreechAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingScreech",
		ActionID: ActionID{SpellID: 55487},
		Duration: time.Second * 4,
	})
	apReductionEffect(aura, 576)
	return aura
}

func apReductionEffect(aura *Aura, apReduction float64) *ExclusiveEffect {
	statReduction := stats.Stats{stats.AttackPower: -apReduction}
	return aura.NewExclusiveEffect("APReduction", false, ExclusiveEffect{
		Priority: apReduction,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, statReduction)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, statReduction.Multiply(-1))
		},
	})
}

func ThunderClapAura(target *Unit, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "ThunderClap-" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 47502},
		Duration: time.Second * 30,
	})
	AtkSpeedReductionEffect(aura, []float64{1.1, 1.14, 1.17, 1.2}[points])
	return aura
}

func InfectedWoundsAura(target *Unit, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "InfectedWounds-" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 48485},
		Duration: time.Second * 12,
	})
	AtkSpeedReductionEffect(aura, []float64{1.0, 1.06, 1.14, 1.20}[points])
	return aura
}

// Note: Paladin code might apply this as part of their judgement auras instead
// of using another separate aura.
func JudgementsOfTheJustAura(target *Unit, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "JudgementsOfTheJust-" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 53696},
		Duration: time.Second * 30,
	})
	AtkSpeedReductionEffect(aura, 1.0+0.1*float64(points))
	return aura
}

func FrostFeverAura(target *Unit, impIcyTouch int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "FrostFever",
		ActionID: ActionID{SpellID: 55095},
		Duration: time.Second * 15,
	})
	AtkSpeedReductionEffect(aura, 1.14+0.02*float64(impIcyTouch))
	return aura
}

func AtkSpeedReductionEffect(aura *Aura, speedMultiplier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("AtkSpdReduction", false, ExclusiveEffect{
		Priority: speedMultiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(sim, 1/speedMultiplier)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(sim, speedMultiplier)
		},
	})
}

func MarkOfBloodAura(target *Unit) *Aura {
	actionId := ActionID{SpellID: 49005}

	var healthMetrics *ResourceMetrics
	aura := target.GetOrRegisterAura(Aura{
		Label:     "MarkOfBlood",
		ActionID:  actionId,
		Duration:  20 * time.Second,
		MaxStacks: 20,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, aura.MaxStacks)

			target := aura.Unit.CurrentTarget

			if healthMetrics == nil && target != nil {
				healthMetrics = target.NewHealthMetrics(actionId)
			}
		},
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			target := aura.Unit.CurrentTarget

			// TODO: Does vampiric blood make it so this health gain is increased?
			if target != nil {
				target.GainHealth(sim, target.MaxHealth()*0.04, healthMetrics)
				aura.RemoveStack(sim)

				if aura.GetStacks() == 0 {
					aura.Deactivate(sim)
				}
			}
		},
	})
	return aura
}

func RuneOfRazoriceVulnerabilityAura(target *Unit) *Aura {
	frostVulnPerStack := 0.02
	aura := target.GetOrRegisterAura(Aura{
		Label:     "RuneOfRazoriceVulnerability",
		ActionID:  ActionID{SpellID: 50401},
		Duration:  NeverExpires,
		MaxStacks: 5,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			oldMultiplier := 1.0 + float64(oldStacks)*frostVulnPerStack
			newMultiplier := 1.0 + float64(newStacks)*frostVulnPerStack
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= newMultiplier / oldMultiplier
		},
	})
	return aura
}

func InsectSwarmAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "InsectSwarmMiss",
		ActionID: ActionID{SpellID: 27013},
		Duration: time.Second * 12,
	})
	increasedMissEffect(aura, 0.03)
	return aura
}

func ScorpidStingAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Scorpid Sting",
		ActionID: ActionID{SpellID: 3043},
		Duration: time.Second * 20,
	})
	increasedMissEffect(aura, 0.03)
	return aura
}

func increasedMissEffect(aura *Aura, increasedMissChance float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("IncreasedMiss", false, ExclusiveEffect{
		Priority: increasedMissChance,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.IncreasedMissChance += increasedMissChance
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.IncreasedMissChance -= increasedMissChance
		},
	})
}

func TotemOfWrathDebuff(target *Unit) *Aura {
	return minorCritDebuffAura(target, "Totem of Wrath Debuff", ActionID{SpellID: 30708}, time.Minute*5, 3*CritRatingPerCritChance)
}

func MasterPoisonerDebuff(target *Unit, points float64) *Aura {
	return minorCritDebuffAura(target, "Master Poisoner", ActionID{SpellID: 58410}, time.Second*20, points*CritRatingPerCritChance)
}

func HeartOfTheCrusaderDebuff(target *Unit, points float64) *Aura {
	return minorCritDebuffAura(target, "Heart of the Crusader", ActionID{SpellID: 20337}, time.Second*20, points*CritRatingPerCritChance)
}

func minorCritDebuffAura(target *Unit, label string, actionID ActionID, duration time.Duration, critBonus float64) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: actionID,
		Duration: duration,
	})
	critBonusEffect(aura, critBonus)
	return aura
}

func critBonusEffect(aura *Aura, critBonus float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("CritBonus", false, ExclusiveEffect{
		Priority: critBonus,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusCritRatingTaken += critBonus
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusCritRatingTaken -= critBonus
		},
	})
}
