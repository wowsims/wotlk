package core

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func applyDebuffEffects(target *Unit, debuffs proto.Debuffs) {
	if debuffs.Misery {
		MakePermanent(MiseryAura(target, 5))
	}

	if debuffs.ShadowWeaving {
		MakePermanent(ShadowWeavingAura(target, 5))
	}

	if debuffs.JudgementOfWisdom {
		MakePermanent(JudgementOfWisdomAura(target))
	}
	if debuffs.JudgementOfLight {
		MakePermanent(JudgementOfLightAura(target))
	}

	if debuffs.ImprovedSealOfTheCrusader {
		MakePermanent(JudgementOfTheCrusaderAura(target, 3, 0, 1.0))
	}

	if debuffs.CurseOfElements != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(CurseOfElementsAura(target))
	}

	if debuffs.ImprovedShadowBolt {
		MakePermanent(ImprovedShadowBoltAura(target))
	}

	if debuffs.IsbUptime > 0.0 {
		uptime := MinFloat(1.0, debuffs.IsbUptime)
		isbAura := MakePermanent(ImprovedShadowBoltAura(target))
		if uptime != 1.0 {
			isbAura.OnDoneIteration = func(aura *Aura, _ *Simulation) {
				aura.metrics.Uptime = time.Duration(float64(aura.metrics.Uptime) * uptime)
			}
		}
	}

	if debuffs.ImprovedScorch {
		MakePermanent(ImprovedScorchAura(target, 5))
	}

	if debuffs.WintersChill {
		MakePermanent(WintersChillAura(target, 5))
	}

	if debuffs.BloodFrenzy {
		MakePermanent(BloodFrenzyAura(target, 2))
	}

	if debuffs.GiftOfArthas {
		MakePermanent(GiftOfArthasAura(target))
	}

	if debuffs.Mangle {
		MakePermanent(MangleAura(target))
	}

	if debuffs.ExposeArmor != proto.TristateEffect_TristateEffectMissing {
		exposeArmorAura := ExposeArmorAura(target, false) // TODO: check glyph
		ScheduledAura(exposeArmorAura, false, PeriodicActionOptions{
			Period:   time.Duration(10.0 * float64(time.Second)),
			NumTicks: 1,
			OnAction: func(sim *Simulation) {
				exposeArmorAura.Activate(sim)
			},
		})
	}

	if debuffs.CurseOfWeakness {
		MakePermanent(CurseOfWeaknessAura(target))
	}

	if debuffs.SunderArmor {
		sunderArmorAura := SunderArmorAura(target, 1)
		ScheduledAura(sunderArmorAura, true, PeriodicActionOptions{
			Period:   time.Duration(1.5 * float64(time.Second)),
			NumTicks: 4,
			Priority: ActionPriorityDOT, // High prio so it comes before actual warrior sunders.
			OnAction: func(sim *Simulation) {
				if sunderArmorAura.IsActive() {
					sunderArmorAura.AddStack(sim)
				}
			},
		})
	}

	if debuffs.FaerieFire != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(FaerieFireAura(target, GetTristateValueInt32(debuffs.FaerieFire, 0, 3)))
	}

	if debuffs.ExposeWeaknessUptime > 0 && debuffs.ExposeWeaknessHunterAgility > 0 {
		uptime := MinFloat(1.0, debuffs.ExposeWeaknessUptime)
		ewAura := MakePermanent(ExposeWeaknessAura(target, debuffs.ExposeWeaknessHunterAgility, uptime))
		if uptime != 1.0 {
			ewAura.OnDoneIteration = func(aura *Aura, _ *Simulation) {
				aura.metrics.Uptime = time.Duration(float64(aura.metrics.Uptime) * uptime)
			}
		}
	}

	if debuffs.HuntersMark != proto.TristateEffect_TristateEffectMissing {
		if debuffs.HuntersMark == proto.TristateEffect_TristateEffectImproved {
			MakePermanent(HuntersMarkAura(target, 5, true))
		} else {
			MakePermanent(HuntersMarkAura(target, 0, true))
		}
	}

	if debuffs.DemoralizingRoar != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingRoarAura(target, GetTristateValueInt32(debuffs.DemoralizingRoar, 0, 5)))
	}
	if debuffs.DemoralizingShout != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingShoutAura(target, 0, GetTristateValueInt32(debuffs.DemoralizingShout, 0, 5)))
	}
	if debuffs.ThunderClap != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(ThunderClapAura(target, GetTristateValueInt32(debuffs.ThunderClap, 0, 3)))
	}
	if debuffs.InsectSwarm {
		MakePermanent(InsectSwarmAura(target))
	}
	if debuffs.ScorpidSting {
		MakePermanent(ScorpidStingAura(target))
	}

	if debuffs.Screech {
		MakePermanent(ScreechAura(target))
	}
}

func ScheduledAura(aura *Aura, preActivate bool, options PeriodicActionOptions) *Aura {
	aura.Duration = NeverExpires
	aura.OnReset = func(aura *Aura, sim *Simulation) {
		if preActivate {
			aura.Activate(sim)
		}
		StartPeriodicAction(sim, options)
	}
	return aura
}

func MiseryAura(target *Unit, numPoints int32) *Aura {
	multiplier := float64(numPoints)

	return target.GetOrRegisterAura(Aura{
		Label:    "Misery-" + strconv.Itoa(int(numPoints)),
		Tag:      "Misery",
		ActionID: ActionID{SpellID: 33195},
		Duration: time.Second * 24,
		Priority: float64(numPoints),
		OnGain: func(aura *Aura, sim *Simulation) {
			target.PseudoStats.BonusSpellHitRating += multiplier * SpellHitRatingPerHitChance
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			target.PseudoStats.BonusSpellHitRating -= multiplier * SpellHitRatingPerHitChance
		},
	})
}

func ShadowWeavingAura(target *Unit, startingStacks int32) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:     "Shadow Weaving",
		ActionID:  ActionID{SpellID: 15334},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, startingStacks)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.ShadowDamageTakenMultiplier /= 1.0 + 0.02*float64(oldStacks)
			aura.Unit.PseudoStats.ShadowDamageTakenMultiplier *= 1.0 + 0.02*float64(newStacks)
		},
	})
}

var JudgementOfWisdomAuraLabel = "Judgement of Wisdom"

func JudgementOfWisdomAura(target *Unit) *Aura {
	const mana = 74 / 2 // 50% proc
	actionID := ActionID{SpellID: 27164}

	return target.GetOrRegisterAura(Aura{
		Label:    JudgementOfWisdomAuraLabel,
		ActionID: actionID,
		Duration: time.Second * 20,
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			if spellEffect.ProcMask.Matches(ProcMaskEmpty) {
				return // Phantom spells (Romulo's, Lightning Capacitor, etc) don't proc JoW.
			}

			// Melee claim that wisdom can proc on misses.
			if !spellEffect.ProcMask.Matches(ProcMaskMeleeOrRanged) && !spellEffect.Landed() {
				return
			}

			unit := spell.Unit
			if unit.HasManaBar() {
				if unit.JowManaMetrics == nil {
					unit.JowManaMetrics = unit.NewManaMetrics(actionID)
				}
				unit.AddMana(sim, mana, unit.JowManaMetrics, false)
			}

			if spell.ActionID.SpellID == 35395 {
				aura.Refresh(sim)
			}
		},
	})
}

var JudgementOfLightAuraLabel = "Judgement of Light"

func JudgementOfLightAura(target *Unit) *Aura {
	actionID := ActionID{SpellID: 27163}

	return target.GetOrRegisterAura(Aura{
		Label:    JudgementOfLightAuraLabel,
		ActionID: actionID,
		Duration: time.Second * 20,
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			if !spellEffect.ProcMask.Matches(ProcMaskMelee) || !spellEffect.Landed() {
				return
			}

			if spell.ActionID.SpellID == 35395 {
				aura.Refresh(sim)
			}
		},
	})
}

func JudgementOfTheCrusaderAura(target *Unit, level int32, flatBonus float64, percentBonus float64) *Aura {
	bonusCrit := float64(level) * CritRatingPerCritChance

	totalSP := 219*percentBonus + flatBonus
	return target.GetOrRegisterAura(Aura{
		Label:    "Judgement of the Crusader-" + strconv.Itoa(int(level)),
		Tag:      "Judgement of the Crusader",
		ActionID: ActionID{SpellID: 27159},
		Duration: time.Second * 20,
		Priority: float64(level),
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusHolyDamageTaken += totalSP
			aura.Unit.PseudoStats.BonusCritRating += bonusCrit
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusHolyDamageTaken -= totalSP
			aura.Unit.PseudoStats.BonusCritRating -= bonusCrit
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			if spell.ActionID.SpellID == 35395 {
				aura.Refresh(sim)
			}
		},
	})
}

func CurseOfElementsAura(target *Unit) *Aura {
	multiplier := 1.13

	return target.GetOrRegisterAura(Aura{
		Label:    "Curse of Elements",
		Tag:      "Curse of Elements",
		ActionID: ActionID{SpellID: 27228},

		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier *= multiplier
			aura.Unit.PseudoStats.FireDamageTakenMultiplier *= multiplier
			aura.Unit.PseudoStats.FrostDamageTakenMultiplier *= multiplier
			aura.Unit.PseudoStats.ShadowDamageTakenMultiplier *= multiplier
			aura.Unit.PseudoStats.NatureDamageTakenMultiplier *= multiplier
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: -165, stats.FireResistance: -165, stats.FrostResistance: -165, stats.ShadowResistance: -165, stats.NatureResistance: -165})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier /= multiplier
			aura.Unit.PseudoStats.FireDamageTakenMultiplier /= multiplier
			aura.Unit.PseudoStats.FrostDamageTakenMultiplier /= multiplier
			aura.Unit.PseudoStats.ShadowDamageTakenMultiplier /= multiplier
			aura.Unit.PseudoStats.NatureDamageTakenMultiplier /= multiplier
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: 165, stats.FireResistance: 165, stats.FrostResistance: 165, stats.ShadowResistance: 165, stats.NatureResistance: 165})
		},
	})
}

func ImprovedShadowBoltAura(target *Unit) *Aura {
	bonusSpellCrit := 5.0 * CritRatingPerCritChance
	config := Aura{
		Label:    "ImprovedShadowBolt",
		Tag:      "ImprovedShadowBolt",
		ActionID: ActionID{SpellID: 17803},
		Duration: time.Second * 30,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusCritRating += bonusSpellCrit
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusCritRating -= bonusSpellCrit
		},
	}

	//	if uptime == 0 {
	//		config.OnSpellHitTaken = func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
	//			if spell.SpellSchool != SpellSchoolShadow {
	//				return
	//			}
	//			if !spellEffect.Landed() || spellEffect.Damage == 0 || !spellEffect.ProcMask.Matches(ProcMaskSpellDamage) {
	//				return
	//			}
	//			aura.RemoveStack(sim)
	//		}
	//	}

	return target.GetOrRegisterAura(config)
}

var BloodFrenzyActionID = ActionID{SpellID: 29859}

func BloodFrenzyAura(target *Unit, points int32) *Aura {
	multiplier := 1 + 0.02*float64(points)
	return target.GetOrRegisterAura(Aura{
		Label:    "Blood Frenzy-" + strconv.Itoa(int(points)),
		Tag:      "Blood Frenzy",
		ActionID: BloodFrenzyActionID,
		// No fixed duration, lasts as long as the bleed that activates it.
		Priority: float64(points),
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.PhysicalDamageTakenMultiplier *= multiplier
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.PhysicalDamageTakenMultiplier /= multiplier
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
	return target.GetOrRegisterAura(Aura{
		Label:    "Mangle",
		ActionID: ActionID{SpellID: 33876},
		Duration: time.Second * 12,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier *= 1.3
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier /= 1.3
		},
	})
}

var ImprovedScorchAuraLabel = "Improved Scorch"

func ImprovedScorchAura(target *Unit, startingStacks int32) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:     ImprovedScorchAuraLabel,
		ActionID:  ActionID{SpellID: 12873},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, startingStacks)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.FireDamageTakenMultiplier /= 1.0 + 0.03*float64(oldStacks)
			aura.Unit.PseudoStats.FireDamageTakenMultiplier *= 1.0 + 0.03*float64(newStacks)
		},
	})
}

var WintersChillAuraLabel = "Winter's Chill"

func WintersChillAura(target *Unit, startingStacks int32) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:     WintersChillAuraLabel,
		ActionID:  ActionID{SpellID: 28595},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, startingStacks)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.BonusFrostCritRating += 2 * CritRatingPerCritChance * float64(newStacks-oldStacks)
		},
	})
}

var MinorArmorReductionAuraTag = "MinorArmorReductionAura"

func FaerieFireAura(target *Unit, level int32) *Aura {
	const armorReduction = 0.05

	return target.GetOrRegisterAura(Aura{
		Label:    "Faerie Fire-" + strconv.Itoa(int(level)),
		Tag:      MinorArmorReductionAuraTag,
		ActionID: ActionID{SpellID: 26993},
		Duration: time.Second * 40,
		Priority: float64(level),
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
			aura.Unit.PseudoStats.BonusSpellHitRating += float64(level) * SpellHitRatingPerHitChance
			aura.Unit.PseudoStats.BonusCritRating += float64(level) * CritRatingPerCritChance
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
			aura.Unit.PseudoStats.BonusSpellHitRating -= float64(level) * SpellHitRatingPerHitChance
			aura.Unit.PseudoStats.BonusCritRating += float64(level) * CritRatingPerCritChance
		},
	})
}

var SunderArmorAuraLabel = "Sunder Armor"
var MajorArmorReductionTag = "MajorArmorReductionAura"

func SunderArmorAura(target *Unit, startingStacks int32) *Aura {
	armorReductionPerStack := 0.04

	return target.GetOrRegisterAura(Aura{
		Label:     SunderArmorAuraLabel,
		Tag:       MajorArmorReductionTag,
		ActionID:  ActionID{SpellID: 25225}, // TODO: Fix spell id
		Duration:  time.Second * 30,
		MaxStacks: 5,
		Priority:  armorReductionPerStack * 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, startingStacks)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			oldMultiplier := (1.0 - float64(oldStacks)*armorReductionPerStack)
			newMultiplier := (1.0 - float64(newStacks)*armorReductionPerStack)
			aura.Unit.PseudoStats.ArmorMultiplier *= (newMultiplier / oldMultiplier)
		},
	})
}

var AcidSpitAuraLabel = "Acid Spit"

func AcidSpitAura(target *Unit, startingStacks int32) *Aura {
	armorReductionPerStack := 0.1

	return target.GetOrRegisterAura(Aura{
		Label:     AcidSpitAuraLabel,
		Tag:       MajorArmorReductionTag,
		ActionID:  ActionID{SpellID: 55745}, // TODO: Spell id
		Duration:  time.Second * 10,
		MaxStacks: 2,
		Priority:  armorReductionPerStack * 2,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, startingStacks)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			oldMultiplier := (1.0 - float64(oldStacks)*armorReductionPerStack)
			newMultiplier := (1.0 - float64(newStacks)*armorReductionPerStack)
			aura.Unit.PseudoStats.ArmorMultiplier *= (newMultiplier / oldMultiplier)
		},
	})
}

func ExposeArmorAura(target *Unit, hasGlyph bool) *Aura {
	armorReduction := 0.2
	duration := time.Second * 30
	if hasGlyph {
		duration += 12
	}
	return target.GetOrRegisterAura(Aura{
		Label:    "ExposeArmor",
		Tag:      MajorArmorReductionTag,
		ActionID: ActionID{SpellID: 26866}, // TODO: Spell id
		Duration: duration,
		Priority: armorReduction,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
		},
	})
}

func CurseOfWeaknessAura(target *Unit) *Aura {
	bonus := stats.Stats{stats.AttackPower: -478}
	armorReduction := 0.05

	return target.GetOrRegisterAura(Aura{
		Label:    "Curse of Weakness",
		Tag:      MinorArmorReductionAuraTag,
		ActionID: ActionID{SpellID: 27226}, // TODO: Fix spell id
		Duration: time.Minute * 2,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, bonus)
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, bonus.Multiply(-1))
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
		},
	})
}

func StingAura(target *Unit) *Aura {
	armorReduction := 0.05

	return target.GetOrRegisterAura(Aura{
		Label:    "Sting",
		Tag:      MinorArmorReductionAuraTag,
		ActionID: ActionID{SpellID: 27226}, // TODO: Fix spell id
		Duration: time.Second * 20,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
		},
	})
}

func SporeCloudAura(target *Unit) *Aura {
	armorReduction := 0.03

	return target.GetOrRegisterAura(Aura{
		Label:    "Spore Cloud",
		ActionID: ActionID{SpellID: 27226}, // TODO: Fix spell id
		Duration: time.Second * 9,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
		},
	})
}

func ShatteringThrowAura(target *Unit) *Aura {
	armorReduction := 0.2

	return target.GetOrRegisterAura(Aura{
		Label:    "Shattering Throw",
		ActionID: ActionID{SpellID: 27226}, // TODO: Fix spell id
		Duration: time.Second * 10,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
		},
	})
}

// Multiplier is for accomodating uptime %. For a real hunter, always pass 1.0
func ExposeWeaknessAura(target *Unit, hunterAgility float64, multiplier float64) *Aura {
	apBonus := hunterAgility * 0.25 * multiplier

	return target.GetOrRegisterAura(Aura{
		Label:    "ExposeWeakness-" + strconv.Itoa(int(hunterAgility)),
		Tag:      "ExposeWeakness",
		ActionID: ActionID{SpellID: 34503},
		Duration: time.Second * 7,
		Priority: apBonus,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusMeleeAttackPower += aura.Priority
			aura.Unit.PseudoStats.BonusRangedAttackPower += aura.Priority
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusMeleeAttackPower -= aura.Priority
			aura.Unit.PseudoStats.BonusRangedAttackPower -= aura.Priority
		},
	})
}

func HuntersMarkAura(target *Unit, points int32, fullyStacked bool) *Aura {
	const baseRangedBonus = 110.0
	const bonusPerStack = 11.0
	const maxStacks = 30
	meleeBonus := baseRangedBonus * 0.2 * float64(points)

	startingStacks := int32(0)
	if fullyStacked {
		startingStacks = maxStacks
	}

	priority := float64(points)
	if fullyStacked {
		// Add a half point so that permanent versions always win.
		priority += 0.5
	}

	return target.GetOrRegisterAura(Aura{
		Label:     "HuntersMark-" + strconv.Itoa(int(points)),
		Tag:       "HuntersMark",
		ActionID:  ActionID{SpellID: 14325},
		Duration:  NeverExpires,
		MaxStacks: 30,
		Priority:  priority,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusMeleeAttackPower += meleeBonus
			aura.Unit.PseudoStats.BonusRangedAttackPower += baseRangedBonus
			aura.SetStacks(sim, startingStacks)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusMeleeAttackPower -= meleeBonus
			aura.Unit.PseudoStats.BonusRangedAttackPower -= baseRangedBonus
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.BonusRangedAttackPower += bonusPerStack * float64(newStacks-oldStacks)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			if spellEffect.ProcMask.Matches(ProcMaskRanged) && spellEffect.Landed() {
				aura.AddStack(sim)
			}
		},
	})
}

const APReductionAuraTag = "APReduction"

func DemoralizingRoarAura(target *Unit, points int32) *Aura {
	apReduction := 248 * (1 + 0.08*float64(points))

	return target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingRoar-" + strconv.Itoa(int(points)),
		Tag:      APReductionAuraTag,
		ActionID: ActionID{SpellID: 26998},
		Duration: time.Second * 30,
		Priority: apReduction,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, -apReduction)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, apReduction)
		},
	})
}

func DemoralizingShoutAura(target *Unit, boomingVoicePts int32, impDemoShoutPts int32) *Aura {
	duration := time.Duration(float64(time.Second*30) * (1 + 0.1*float64(boomingVoicePts)))
	apReduction := 300 * (1 + 0.08*float64(impDemoShoutPts))

	return target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingShout-" + strconv.Itoa(int(impDemoShoutPts)),
		Tag:      APReductionAuraTag,
		ActionID: ActionID{SpellID: 25203},
		Duration: duration,
		Priority: apReduction,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, -apReduction)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, apReduction)
		},
	})
}

func ScreechAura(target *Unit) *Aura {
	const apReduction = 210.0

	return target.GetOrRegisterAura(Aura{
		Label:    "Screech",
		ActionID: ActionID{SpellID: 27051},
		Duration: time.Second * 4,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, -apReduction)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, apReduction)
		},
	})
}

const ThunderClapAuraTag = "ThunderClap"

func ThunderClapAura(target *Unit, points int32) *Aura {
	speedMultiplier := 0.9
	if points == 1 {
		speedMultiplier = 0.86
	} else if points == 2 {
		speedMultiplier = 0.83
	} else if points == 3 {
		speedMultiplier = 0.8
	}
	inverseMult := 1 / speedMultiplier

	return target.GetOrRegisterAura(Aura{
		Label:    "ThunderClap-" + strconv.Itoa(int(points)),
		Tag:      ThunderClapAuraTag,
		ActionID: ActionID{SpellID: 25264},
		Duration: time.Second * 30,
		Priority: inverseMult,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, speedMultiplier)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, inverseMult)
		},
	})
}

func InsectSwarmAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "InsectSwarmMiss",
		ActionID: ActionID{SpellID: 27013},
		Duration: time.Second * 12,
		OnGain: func(aura *Aura, sim *Simulation) {
			if !aura.Unit.HasActiveAura("ScorpidSting") {
				aura.Unit.PseudoStats.IncreasedMissChance += 0.02
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if !aura.Unit.HasActiveAura("ScorpidSting") {
				aura.Unit.PseudoStats.IncreasedMissChance -= 0.02
			}
		},
	})
}

func ScorpidStingAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Scorpid Sting",
		ActionID: ActionID{SpellID: 3043},
		Duration: time.Second * 20,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.IncreasedMissChance += 0.05
			if aura.Unit.HasActiveAura("InsectSwarmMiss") {
				aura.Unit.PseudoStats.IncreasedMissChance -= 0.02
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.IncreasedMissChance -= 0.05
			if aura.Unit.HasActiveAura("InsectSwarmMiss") {
				aura.Unit.PseudoStats.IncreasedMissChance += 0.02
			}
		},
	})
}
