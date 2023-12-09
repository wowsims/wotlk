package core

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func applyDebuffEffects(target *Unit, targetIdx int, debuffs *proto.Debuffs, raid *proto.Raid) {

	if debuffs.JudgementOfWisdom && targetIdx == 0 {
		MakePermanent(JudgementOfWisdomAura(target))
	}

	if debuffs.ImprovedShadowBolt {
		//TODO: Apply periodically
		MakePermanent(ImprovedShadowBoltAura(target, 5))
	}

	if debuffs.CurseOfElements {
		MakePermanent(CurseOfElementsAura(target))
	}

	if debuffs.CurseOfShadow {
		MakePermanent(CurseOfShadowAura(target))
	}

	if debuffs.ImprovedScorch && targetIdx == 0 {
		MakePermanent(ImprovedScorchAura(target))
	}

	if debuffs.WintersChill && targetIdx == 0 {
		MakePermanent(WintersChillAura(target, 5))
	}

	if debuffs.GiftOfArthas {
		MakePermanent(GiftOfArthasAura(target))
	}

	if debuffs.CrystalYield {
		MakePermanent(CrystalYieldAura(target))
	}

	if debuffs.ExposeArmor && targetIdx == 0 {
		aura := ExposeArmorAura(target, false, raid.Parties[0].Players[0].Level)
		ScheduledMajorArmorAura(aura, PeriodicActionOptions{
			Period:   time.Second * 3,
			NumTicks: 1,
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
			},
		}, raid)
	}

	if debuffs.SunderArmor && targetIdx == 0 {
		// Get the first player level as the debuff level
		aura := SunderArmorAura(target, raid.Parties[0].Players[0].Level)
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

	if debuffs.CurseOfWeakness != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(CurseOfWeaknessAura(target, GetTristateValueInt32(debuffs.CurseOfWeakness, 1, 2), raid.Parties[0].Players[0].Level))
	}

	if debuffs.DemoralizingRoar != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingRoarAura(target, GetTristateValueInt32(debuffs.DemoralizingRoar, 0, 5), raid.Parties[0].Players[0].Level))
	}
	if debuffs.DemoralizingShout != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingShoutAura(target, 0, GetTristateValueInt32(debuffs.DemoralizingShout, 0, 5), raid.Parties[0].Players[0].Level))
	}

	// Atk spd reduction
	if debuffs.ThunderClap != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(ThunderClapAura(target, GetTristateValueInt32(debuffs.ThunderClap, 0, 3), raid.Parties[0].Players[0].Level))
	}

	// Miss
	if debuffs.InsectSwarm && targetIdx == 0 {
		MakePermanent(InsectSwarmAura(target))
	}
	if debuffs.ScorpidSting && targetIdx == 0 {
		MakePermanent(ScorpidStingAura(target))
	}

	// if debuffs.HuntersMark > 0 && targetIdx == 0 {
	// 	points := int32(0)
	// 	glyphed := false
	// 	if debuffs.HuntersMark > 1 {
	// 		points = 3
	// 		if debuffs.HuntersMark > 2 {
	// 			glyphed = true
	// 		}
	// 	}
	// 	MakePermanent(HuntersMarkAura(target, points, glyphed))
	// }
}

func ImprovedShadowBoltAura(unit *Unit, level int32) *Aura {
	damageMulti := 1. + 0.04*float64(level)
	return unit.RegisterAura(Aura{
		Label:     "Improved Shadow Bolt",
		ActionID:  ActionID{SpellID: 17800},
		Duration:  12 * time.Second,
		MaxStacks: 4,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= damageMulti
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= damageMulti
		},
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if spell.SpellSchool != SpellSchoolShadow {
				return
			}

			if !result.Landed() {
				return
			}

			aura.RemoveStack(sim)
		},
	})
}

func ScheduledMajorArmorAura(aura *Aura, options PeriodicActionOptions, raid *proto.Raid) {
	// Individual rogue sim rotation option messes with these debuff options,
	// so it has to be handled separately.
	allRogues := RaidPlayersWithClass(raid, proto.Class_ClassRogue)
	singleExposeDelay := len(allRogues) == 1 &&
		allRogues[0].Spec.(*proto.Player_Rogue).Rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once

	if singleExposeDelay {
		target := aura.Unit
		exposeArmorAura := ExposeArmorAura(target, false, raid.Parties[0].Players[0].Level)
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

// TODO: Classic verify logic
func JudgementOfWisdomAura(target *Unit) *Aura {
	actionID := ActionID{SpellID: 20357}

	return target.GetOrRegisterAura(Aura{
		Label:    JudgementOfWisdomAuraLabel,
		ActionID: actionID,
		Duration: time.Second * 10,
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			unit := spell.Unit
			if !unit.HasManaBar() {
				return
			}

			if spell.ProcMask.Matches(ProcMaskEmpty | ProcMaskProc | ProcMaskWeaponProc) {
				return // Phantom spells (Romulo's, Lightning Capacitor, etc.) don't proc JoW.
			}

			if spell.ProcMask.Matches(ProcMaskWhiteHit | ProcMaskRanged) {
				// Apparently ranged/melee can still proc on miss
				if !unit.AutoAttacks.PPMProc(sim, 15, ProcMaskWhiteHit|ProcMaskRanged, "jow", spell) {
					return
				}
			} else { // spell casting
				if !result.Landed() {
					return
				}

				ct := spell.CurCast.CastTime.Seconds()
				if ct == 0 {
					// Current theory is that insta-cast is treated as min GCD from retail.
					// Perhaps this is a bug introduced in classic when converting JoW to wotlk.
					ct = 0.75
				}
				procChance := ct * 0.25 // ct / 60.0 * 15.0PPM (algebra) = ct*0.25
				if sim.RandomFloat("jow") > procChance {
					return
				}
			}

			if unit.JowManaMetrics == nil {
				unit.JowManaMetrics = unit.NewManaMetrics(actionID)
			}
			// JoW returns flat mana
			unit.AddMana(sim, 59, unit.JowManaMetrics)
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
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.FireResistance: -75, stats.FrostResistance: -75})
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 1.1
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 1.1
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.FireResistance: 75, stats.FrostResistance: 75})
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= 1.1
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= 1.1
		},
	})
	return aura
}

func CurseOfShadowAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Curse of Shadow",
		ActionID: ActionID{SpellID: 17937},
		Duration: time.Minute * 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: -75, stats.ShadowResistance: -75})
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 1.1
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 1.1
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: 75, stats.ShadowResistance: 75})
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= 1.1
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= 1.1
		},
	})
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
		ActionID: ActionID{SpellID: 409828},
		Duration: time.Minute,
	}, 1.3)
}

// Bleed Damage Multiplier category
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

const SpellFirePowerEffectCategory = "spellFirePowerdebuff"

func ImprovedScorchAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Improved Scorch",
		ActionID: ActionID{SpellID: 12873},
		Duration: time.Second * 30,
	})

	aura.NewExclusiveEffect(SpellFirePowerEffectCategory, true, ExclusiveEffect{
		Priority: 0.15,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 0.15
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= 0.15
		},
	})

	return aura
}

const SpellCritEffectCategory = "spellcritdebuff"

func WintersChillAura(target *Unit, startingStacks int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Winter's Chill",
		ActionID:  ActionID{SpellID: 28593},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, startingStacks)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks, newStacks int32) {
			aura.Unit.PseudoStats.SchoolCritTakenMultiplier[stats.SchoolIndexFrost] /= 1 + 0.2*float64(oldStacks)
			aura.Unit.PseudoStats.SchoolCritTakenMultiplier[stats.SchoolIndexFrost] *= 1 + 0.2*float64(newStacks)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolCritTakenMultiplier[stats.SchoolIndexFrost] /= 1 + 0.2*float64(aura.stacks)
		},
	})

	// effect = aura.NewExclusiveEffect(SpellCritEffectCategory, true, ExclusiveEffect{
	// 	Priority: 0,
	// 	OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
	// 		ee.Aura.Unit.PseudoStats.BonusSpellCritRatingTaken += ee.Priority * CritRatingPerCritChance
	// 	},
	// 	OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
	// 		ee.Aura.Unit.PseudoStats.BonusSpellCritRatingTaken -= ee.Priority * CritRatingPerCritChance
	// 	},
	// })
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

var majorArmorReductionEffectCategory = "MajorArmorReduction"

func SunderArmorAura(target *Unit, playerLevel int32) *Aura {
	spellID := map[int32]int32{
		25: 7405,
		40: 8380,
		50: 11596,
		60: 11597,
	}[playerLevel]

	arpen := map[int32]int32{
		25: 180,
		40: 270,
		50: 360,
		60: 450,
	}[playerLevel]

	aura := target.GetOrRegisterAura(Aura{
		Label:     "Sunder Armor",
		ActionID:  ActionID{SpellID: spellID},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.stats[stats.Armor] -= float64(arpen) * float64(oldStacks)
			aura.Unit.stats[stats.Armor] += float64(arpen) * float64(newStacks)

		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.stats[stats.Armor] -= float64(arpen) * float64(aura.stacks)
		},
	})

	return aura
}

// TODO: Classic (Flat amount)
func ExposeArmorAura(target *Unit, hasGlyph bool, playerLevel int32) *Aura {
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

// TODO: Classic
func CurseOfWeaknessAura(target *Unit, points int32, playerLevel int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Curse of Weakness" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 50511},
		Duration: time.Minute * 2,
	})
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

const HuntersMarkAuraTag = "HuntersMark"

// TODO: Classic
func HuntersMarkAura(target *Unit, points int32, glyphed bool, playerLevel int32) *Aura {
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

// TODO: Classic
func DemoralizingRoarAura(target *Unit, points int32, playerLevel int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingRoar-" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 48560},
		Duration: time.Second * 30,
	})
	apReductionEffect(aura, 411*(1+0.08*float64(points)))
	return aura
}

// TODO: Classic
func DemoralizingShoutAura(target *Unit, boomingVoicePts int32, impDemoShoutPts int32, playerLevel int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingShout-" + strconv.Itoa(int(impDemoShoutPts)),
		ActionID: ActionID{SpellID: 1160},
		Duration: time.Duration(float64(time.Second*30) * (1 + 0.1*float64(boomingVoicePts))),
	})
	apReduction := map[int32]float64{
		25: 56,
		40: 76,
		50: 111,
		60: 146,
	}[playerLevel]
	apReductionEffect(aura, apReduction*(1+0.08*float64(impDemoShoutPts)))
	return aura
}

// TODO: Classic
func VindicationAura(target *Unit, points int32, playerLevel int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Vindication",
		ActionID: ActionID{SpellID: 26016},
		Duration: time.Second * 10,
	})
	apReductionEffect(aura, 287*float64(points))
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
			ee.Aura.Unit.AddStatsDynamic(sim, statReduction.Invert())
		},
	})
}

// TODO: Classic
func ThunderClapAura(target *Unit, points int32, playerLevel int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "ThunderClap-" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 47502},
		Duration: time.Second * 30,
	})
	AtkSpeedReductionEffect(aura, []float64{1.1, 1.14, 1.17, 1.2}[points])
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

func InsectSwarmAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "InsectSwarmMiss",
		ActionID: ActionID{SpellID: 24977},
		Duration: time.Second * 12,
	})
	increasedMissEffect(aura, 0.02)
	return aura
}

func ScorpidStingAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Scorpid Sting",
		ActionID: ActionID{SpellID: 3043},
		Duration: time.Second * 20,
	})
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

func CrystalYieldAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Crystal Yield",
		ActionID: ActionID{SpellID: 15235},
		Duration: 2 * time.Minute,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.stats[stats.Armor] -= 200
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.stats[stats.Armor] += 200
		},
	})
}
