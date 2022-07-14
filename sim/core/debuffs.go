package core

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func applyDebuffEffects(target *Unit, debuffs proto.Debuffs) {
	if debuffs.Misery {
		MakePermanent(MiseryAura(target))
	}

	if debuffs.JudgementOfWisdom {
		MakePermanent(JudgementOfWisdomAura(target))
	}
	if debuffs.JudgementOfLight {
		MakePermanent(JudgementOfLightAura(target))
	}

	if debuffs.CurseOfElements {
		MakePermanent(CurseOfElementsAura(target))
	}
	if debuffs.EbonPlaguebringer {
		MakePermanent(EbonPlaguebringerAura(target))
	}
	if debuffs.EarthAndMoon {
		MakePermanent(EarthAndMoonAura(target))
	}

	if debuffs.ImprovedShadowBolt {
		MakePermanent(ImprovedShadowBoltAura(target))
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
	if debuffs.SavageCombat {
		MakePermanent(SavageCombatAura(target, 2))
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

	if debuffs.CurseOfWeakness != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(CurseOfWeaknessAura(target, 2))
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
		MakePermanent(FaerieFireAura(target, debuffs.FaerieFire == proto.TristateEffect_TristateEffectImproved))
	}

	if debuffs.HuntersMark != proto.TristateEffect_TristateEffectMissing {
		if debuffs.HuntersMark == proto.TristateEffect_TristateEffectImproved {
			MakePermanent(HuntersMarkAura(target, 3, true))
		} else {
			MakePermanent(HuntersMarkAura(target, 0, false))
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

func MiseryAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Misery",
		Tag:      MinorSpellHitDebuffAuraTag,
		Priority: 3,
		ActionID: ActionID{SpellID: 33198},
		Duration: time.Second * 24,
		OnGain: func(aura *Aura, sim *Simulation) {
			target.PseudoStats.BonusSpellHitRating += 3 * SpellHitRatingPerHitChance
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			target.PseudoStats.BonusSpellHitRating -= 3 * SpellHitRatingPerHitChance
		},
	})
}

var JudgementOfWisdomAuraLabel = "Judgement of Wisdom"

func JudgementOfWisdomAura(target *Unit) *Aura {
	actionID := ActionID{SpellID: 53408}

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
			if unit.HasManaBar() && sim.RandomFloat("jow") > 0.5 {
				if unit.JowManaMetrics == nil {
					unit.JowManaMetrics = unit.NewManaMetrics(actionID)
				}
				// JoW returns 2% of base mana 50% of the time.
				unit.AddMana(sim, unit.BaseMana*0.02, unit.JowManaMetrics, false)
			}

			if spell.ActionID.SpellID == 35395 { // Crusader strike
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

const spelldmgtag = `13%dmg`

func CurseOfElementsAura(target *Unit) *Aura {
	multiplier := 1.13

	return target.GetOrRegisterAura(Aura{
		Label:    "Curse of Elements",
		Tag:      spelldmgtag,
		ActionID: ActionID{SpellID: 47865},
		OnGain: func(aura *Aura, sim *Simulation) {
			if !target.HasActiveAuraWithTag(spelldmgtag) {
				aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.FireDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.FrostDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.ShadowDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.NatureDamageTakenMultiplier *= multiplier
			}
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: -165, stats.FireResistance: -165, stats.FrostResistance: -165, stats.ShadowResistance: -165, stats.NatureResistance: -165})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if !target.HasActiveAuraWithTag(spelldmgtag) {
				aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.FireDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.FrostDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.ShadowDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.NatureDamageTakenMultiplier /= multiplier
			}
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: 165, stats.FireResistance: 165, stats.FrostResistance: 165, stats.ShadowResistance: 165, stats.NatureResistance: 165})
		},
	})
}

func EarthAndMoonAura(target *Unit) *Aura {
	return earthMoonEbonPlaguebringerAura(target, "Earth And Moon", 48511)
}

func EbonPlaguebringerAura(target *Unit) *Aura {
	return earthMoonEbonPlaguebringerAura(target, "Ebon Plaguebringer", 51161)
}

func earthMoonEbonPlaguebringerAura(target *Unit, label string, id int32) *Aura {
	multiplier := 1.13

	return target.GetOrRegisterAura(Aura{
		Label:    label,
		Tag:      spelldmgtag,
		ActionID: ActionID{SpellID: id},
		OnGain: func(aura *Aura, sim *Simulation) {
			if !target.HasActiveAuraWithTag(spelldmgtag) {
				aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.FireDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.FrostDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.ShadowDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.NatureDamageTakenMultiplier *= multiplier
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if !target.HasActiveAuraWithTag(spelldmgtag) {
				aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.FireDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.FrostDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.ShadowDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.NatureDamageTakenMultiplier /= multiplier
			}
		},
	})
}

func ImprovedShadowBoltAura(target *Unit) *Aura {
	bonusSpellCrit := 5.0 * CritRatingPerCritChance
	config := Aura{
		Label:    "Shadow Mastery",
		Tag:      "Shadow Mastery",
		ActionID: ActionID{SpellID: 17800},
		Duration: time.Second * 30,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusCritRating += bonusSpellCrit
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusCritRating -= bonusSpellCrit
		},
	}

	return target.GetOrRegisterAura(config)
}

var BloodFrenzyActionID = ActionID{SpellID: 29859}
var phyDmgDebuff = `4%phydmg`

func BloodFrenzyAura(target *Unit, points int32) *Aura {
	return bloodFrenzySavageCombatAura(target, "Blood Frenzy", BloodFrenzyActionID, float64(points))
}
func SavageCombatAura(target *Unit, points int32) *Aura {
	return bloodFrenzySavageCombatAura(target, "Savage Combat", ActionID{SpellID: 58413}, float64(points))
}

func bloodFrenzySavageCombatAura(target *Unit, label string, id ActionID, points float64) *Aura {
	multiplier := 1 + 0.02*points
	return target.GetOrRegisterAura(Aura{
		Label:    label + "-" + strconv.Itoa(int(points)),
		Tag:      phyDmgDebuff,
		ActionID: id,
		// No fixed duration, lasts as long as the bleed that activates it.
		OnGain: func(aura *Aura, sim *Simulation) {
			if oAura := target.GetActiveAuraWithTag(phyDmgDebuff); oAura == nil {
				aura.Unit.PseudoStats.PhysicalDamageTakenMultiplier *= multiplier
			} else if oAura.Priority < points {
				// remove weaker debuff
				aura.Unit.PseudoStats.PhysicalDamageTakenMultiplier /= 1.0 + (0.02 * oAura.Priority)
				aura.Unit.PseudoStats.PhysicalDamageTakenMultiplier *= multiplier
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if oAura := target.GetActiveAuraWithTag(phyDmgDebuff); oAura == nil {
				aura.Unit.PseudoStats.PhysicalDamageTakenMultiplier /= multiplier
			} else if oAura.Priority < points {
				aura.Unit.PseudoStats.PhysicalDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.PhysicalDamageTakenMultiplier *= 1.0 + (0.02 * oAura.Priority)
			}
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
var MinorSpellHitDebuffAuraTag = "sphit3%"

func FaerieFireAura(target *Unit, imp bool) *Aura {
	const armorReduction = 0.05

	var mainAura *Aura
	var secondaryAura *Aura

	label := "Faerie Fire"
	if imp {
		label = "Improved " + label
		secondaryAura = target.GetOrRegisterAura(Aura{
			Label:    "Improved Faerie Fire Secondary",
			Tag:      MinorSpellHitDebuffAuraTag,
			Duration: time.Minute * 5,
			Priority: 3,
			// no ActionID to hide this secondary effect from stats
			OnGain: func(aura *Aura, sim *Simulation) {
				aura.Unit.PseudoStats.BonusSpellHitRating += 3 * SpellHitRatingPerHitChance
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				aura.Unit.PseudoStats.BonusSpellHitRating -= 3 * SpellHitRatingPerHitChance
				if mainAura.IsActive() {
					mainAura.Deactivate(sim)
				}
			},
		})

	}
	mainAura = target.GetOrRegisterAura(Aura{
		Label:    label,
		Tag:      MinorArmorReductionAuraTag,
		Priority: armorReduction,
		ActionID: ActionID{SpellID: 770},
		Duration: time.Minute * 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
			aura.Unit.updateArmor()
			if imp {
				secondaryAura.Activate(sim)
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
			aura.Unit.updateArmor()
			if secondaryAura.IsActive() {
				secondaryAura.Deactivate(sim)
			}
		},
	})
	return mainAura
}

var SunderArmorAuraLabel = "Sunder Armor"
var MajorArmorReductionTag = "MajorArmorReductionAura"

func SunderArmorAura(target *Unit, startingStacks int32) *Aura {
	armorReductionPerStack := 0.04

	return target.GetOrRegisterAura(Aura{
		Label:     SunderArmorAuraLabel,
		Tag:       MajorArmorReductionTag,
		ActionID:  ActionID{SpellID: 47467},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		Priority:  armorReductionPerStack * 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, startingStacks)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			oldMultiplier := 1.0 - float64(oldStacks)*armorReductionPerStack
			newMultiplier := 1.0 - float64(newStacks)*armorReductionPerStack
			aura.Unit.PseudoStats.ArmorMultiplier *= newMultiplier / oldMultiplier
			aura.Unit.updateArmor()
		},
	})
}

var AcidSpitAuraLabel = "Acid Spit"

func AcidSpitAura(target *Unit, startingStacks int32) *Aura {
	armorReductionPerStack := 0.1

	return target.GetOrRegisterAura(Aura{
		Label:     AcidSpitAuraLabel,
		Tag:       MajorArmorReductionTag,
		ActionID:  ActionID{SpellID: 55754},
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
			aura.Unit.updateArmor()
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
		ActionID: ActionID{SpellID: 48669},
		Duration: duration,
		Priority: armorReduction,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
			aura.Unit.updateArmor()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
			aura.Unit.updateArmor()
		},
	})
}

func CurseOfWeaknessAura(target *Unit, points int32) *Aura {
	APReduction := 478 * (1 + 0.1*float64(points))
	bonus := stats.Stats{stats.AttackPower: -APReduction}
	armorReduction := 0.05

	return target.GetOrRegisterAura(Aura{
		Label:    "Curse of Weakness",
		Tag:      MinorArmorReductionAuraTag,
		Priority: armorReduction,
		ActionID: ActionID{SpellID: 50511},
		Duration: time.Minute * 2,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, bonus)
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
			aura.Unit.updateArmor()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, bonus.Multiply(-1))
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
			aura.Unit.updateArmor()
		},
	})
}

func StingAura(target *Unit) *Aura {
	armorReduction := 0.05

	return target.GetOrRegisterAura(Aura{
		Label:    "Sting",
		Tag:      MinorArmorReductionAuraTag,
		Priority: armorReduction,
		ActionID: ActionID{SpellID: 56631},
		Duration: time.Second * 20,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
			aura.Unit.updateArmor()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
			aura.Unit.updateArmor()
		},
	})
}

func SporeCloudAura(target *Unit) *Aura {
	armorReduction := 0.03

	return target.GetOrRegisterAura(Aura{
		Label:    "Spore Cloud",
		ActionID: ActionID{SpellID: 53598},
		Duration: time.Second * 9,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
			aura.Unit.updateArmor()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
			aura.Unit.updateArmor()
		},
	})
}

func ShatteringThrowAura(target *Unit) *Aura {
	armorReduction := 0.2

	return target.GetOrRegisterAura(Aura{
		Label:    "Shattering Throw",
		ActionID: ActionID{SpellID: 64382},
		Duration: time.Second * 10,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
			aura.Unit.updateArmor()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
			aura.Unit.updateArmor()
		},
	})
}

func HuntersMarkAura(target *Unit, points int32, glyphed bool) *Aura {
	bonus := 500.0 * (1 + 0.1*float64(points))
	priority := float64(points)

	if glyphed {
		bonus += 500.0 * 0.2
		priority += 2
	}

	return target.GetOrRegisterAura(Aura{
		Label:    "HuntersMark-" + strconv.Itoa(int(priority)),
		Tag:      "HuntersMark",
		ActionID: ActionID{SpellID: 53338},
		Duration: NeverExpires,
		Priority: priority,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusRangedAttackPower += bonus
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusRangedAttackPower -= bonus
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

const IcyTouchAuraTag = "IcyTouch"

func IcyTouchAura(target *Unit, impIcyTouch int32) *Aura {
	speedMultiplier := 0.85
	if impIcyTouch > 0 {
		speedMultiplier -= 0.02 * float64(impIcyTouch)
	}

	inverseMult := 1 / speedMultiplier
	return target.GetOrRegisterAura(Aura{
		Label:    "IcyTouch",
		Tag:      IcyTouchAuraTag,
		ActionID: ActionID{SpellID: 49909},
		Duration: time.Second * 15,
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
