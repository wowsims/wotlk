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
		// Crypt fever is embedded in EP but due to aura stacking
		// and tags we need it as a separate aura with its unique tag
		// for the disease damage
		MakePermanent(CryptFeverAura(target, -1))
		MakePermanent(EbonPlaguebringerAura(target, -1))
	}
	if debuffs.EarthAndMoon {
		MakePermanent(EarthAndMoonAura(target))
	}

	if debuffs.ShadowMastery {
		MakePermanent(ShadowMasteryAura(target))
	}

	if debuffs.ImprovedScorch {
		MakePermanent(ImprovedScorchAura(target))
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

	if debuffs.SporeCloud {
		MakePermanent(SporeCloudAura(target))
	}

	if debuffs.Mangle {
		MakePermanent(MangleAura(target))
	} else if debuffs.Trauma {
		MakePermanent(TraumaAura(target, 2))
	} else if debuffs.Stampede {
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

	if debuffs.ExposeArmor {
		exposeArmorAura := ExposeArmorAura(target, false)
		ScheduledAura(exposeArmorAura, false, ExposeArmorPeriodicActonOptions(exposeArmorAura))
	}

	if debuffs.SunderArmor {
		sunderArmorAura := SunderArmorAura(target, 1)
		ScheduledAura(sunderArmorAura, true, SunderArmorPeriodicActionOptions(sunderArmorAura))
	}

	if debuffs.AcidSpit {
		acidSpitAura := AcidSpitAura(target, 1)
		ScheduledAura(acidSpitAura, true, AcidSpitPeriodicActionOptions(acidSpitAura))
	}

	if debuffs.CurseOfWeakness != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(CurseOfWeaknessAura(target, 2))
	}
	if debuffs.Sting {
		MakePermanent(StingAura(target))
	}

	if debuffs.FaerieFire != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(FaerieFireAura(target, debuffs.FaerieFire == proto.TristateEffect_TristateEffectImproved))
	}

	if debuffs.DemoralizingRoar != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingRoarAura(target, GetTristateValueInt32(debuffs.DemoralizingRoar, 0, 5)))
	}
	if debuffs.DemoralizingShout != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingShoutAura(target, 0, GetTristateValueInt32(debuffs.DemoralizingShout, 0, 5)))
	}
	if debuffs.Vindication {
		MakePermanent(VindicationAura(target))
	}

	// Atk spd reduction
	if debuffs.ThunderClap != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(ThunderClapAura(target, GetTristateValueInt32(debuffs.ThunderClap, 0, 3)))
	}
	if debuffs.FrostFever != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(FrostFeverAura(target, GetTristateValueInt32(debuffs.FrostFever, 0, 3)))
	}
	if debuffs.InfectedWounds {
		MakePermanent(InfectedWoundsAura(target, 3))
	}
	if debuffs.JudgementsOfTheJust {
		MakePermanent(JudgementsOfTheJustAura(target, 2))
	}

	// Miss
	if debuffs.InsectSwarm {
		MakePermanent(InsectSwarmAura(target))
	}
	if debuffs.ScorpidSting {
		MakePermanent(ScorpidStingAura(target))
	}

	if debuffs.Screech {
		MakePermanent(ScreechAura(target))
	}

	if debuffs.TotemOfWrath {
		MakePermanent(TotemOfWrathDebuff(target))
	}

	if debuffs.MasterPoisoner {
		MakePermanent(MasterPoisonerDebuff(target, 3))
	}

	if debuffs.HeartOfTheCrusader {
		MakePermanent(HeartoftheCrusaderDebuff(target, 3))
	}

	if debuffs.HuntersMark > 0 {
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

func AcidSpitPeriodicActionOptions(aura *Aura) PeriodicActionOptions {
	return PeriodicActionOptions{
		Period:   time.Second * 10,
		NumTicks: 1,
		OnAction: func(sim *Simulation) {
			if aura.IsActive() {
				aura.AddStack(sim)
			}
		},
	}
}

func ExposeArmorPeriodicActonOptions(aura *Aura) PeriodicActionOptions {
	return PeriodicActionOptions{
		Period:   time.Second * 3,
		NumTicks: 1,
		OnAction: func(sim *Simulation) {
			aura.Activate(sim)
		},
	}
}

func SunderArmorPeriodicActionOptions(aura *Aura) PeriodicActionOptions {
	return PeriodicActionOptions{
		Period:   time.Millisecond * 1500,
		NumTicks: 4,
		Priority: ActionPriorityDOT, // High prio so it comes before actual warrior sunders.
		OnAction: func(sim *Simulation) {
			if aura.IsActive() {
				aura.AddStack(sim)
			}
		},
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
			aura.Unit.PseudoStats.BonusSpellHitRatingTaken += 3 * SpellHitRatingPerHitChance
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusSpellHitRatingTaken -= 3 * SpellHitRatingPerHitChance
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
	actionID := ActionID{SpellID: 20271}

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
			if !target.HasActiveAuraWithTagExcludingAura(spelldmgtag, aura) {
				aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.FireDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.FrostDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.ShadowDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.NatureDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.HolyDamageTakenMultiplier *= multiplier
			}
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: -165, stats.FireResistance: -165, stats.FrostResistance: -165, stats.ShadowResistance: -165, stats.NatureResistance: -165})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if !target.HasActiveAuraWithTagExcludingAura(spelldmgtag, aura) {
				aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.FireDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.FrostDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.ShadowDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.NatureDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.HolyDamageTakenMultiplier /= multiplier
			}
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: 165, stats.FireResistance: 165, stats.FrostResistance: 165, stats.ShadowResistance: 165, stats.NatureResistance: 165})
		},
	})
}

func EarthAndMoonAura(target *Unit) *Aura {
	multiplier := 1.13

	return target.GetOrRegisterAura(Aura{
		Label:    "Earth And Moon",
		Tag:      spelldmgtag,
		ActionID: ActionID{SpellID: 48511},
		OnGain: func(aura *Aura, sim *Simulation) {
			if !target.HasActiveAuraWithTagExcludingAura(spelldmgtag, aura) {
				aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.FireDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.FrostDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.ShadowDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.NatureDamageTakenMultiplier *= multiplier
				aura.Unit.PseudoStats.HolyDamageTakenMultiplier *= multiplier
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if !target.HasActiveAuraWithTagExcludingAura(spelldmgtag, aura) {
				aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.FireDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.FrostDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.ShadowDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.NatureDamageTakenMultiplier /= multiplier
				aura.Unit.PseudoStats.HolyDamageTakenMultiplier /= multiplier
			}
		},
	})
}

const diseasedmgtag = "diseasedmg"
const CryptFeverAuraLabel = "Crypt Fever-"

func CryptFeverAura(target *Unit, dkIndex int) *Aura {
	diseaseMultiplier := 1.3

	return target.GetOrRegisterAura(Aura{
		Label:    CryptFeverAuraLabel + strconv.Itoa(dkIndex), // Support multiple DKs having their CF up
		Tag:      diseasedmgtag,
		ActionID: ActionID{SpellID: 49632},
		OnGain: func(aura *Aura, sim *Simulation) {
			if !target.HasActiveAuraWithTagExcludingAura(diseasedmgtag, aura) {
				aura.Unit.PseudoStats.DiseaseDamageTakenMultiplier *= diseaseMultiplier
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if !target.HasActiveAuraWithTagExcludingAura(diseasedmgtag, aura) {
				aura.Unit.PseudoStats.DiseaseDamageTakenMultiplier /= diseaseMultiplier
			}
		},
	})
}

const EbonPlaguebringerAuraLabel = "Ebon Plaguebringer-"

func EbonPlaguebringerAura(target *Unit, dkIndex int) *Aura {
	magicMultiplier := 1.13

	return target.GetOrRegisterAura(Aura{
		Label:    EbonPlaguebringerAuraLabel + strconv.Itoa(dkIndex), // Support multiple DKs having their EP up
		Tag:      spelldmgtag,
		ActionID: ActionID{SpellID: 51161},
		OnGain: func(aura *Aura, sim *Simulation) {
			if !target.HasActiveAuraWithTagExcludingAura(spelldmgtag, aura) {
				aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier *= magicMultiplier
				aura.Unit.PseudoStats.FireDamageTakenMultiplier *= magicMultiplier
				aura.Unit.PseudoStats.FrostDamageTakenMultiplier *= magicMultiplier
				aura.Unit.PseudoStats.ShadowDamageTakenMultiplier *= magicMultiplier
				aura.Unit.PseudoStats.NatureDamageTakenMultiplier *= magicMultiplier
				aura.Unit.PseudoStats.HolyDamageTakenMultiplier *= magicMultiplier
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if !target.HasActiveAuraWithTagExcludingAura(spelldmgtag, aura) {
				aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier /= magicMultiplier
				aura.Unit.PseudoStats.FireDamageTakenMultiplier /= magicMultiplier
				aura.Unit.PseudoStats.FrostDamageTakenMultiplier /= magicMultiplier
				aura.Unit.PseudoStats.ShadowDamageTakenMultiplier /= magicMultiplier
				aura.Unit.PseudoStats.NatureDamageTakenMultiplier /= magicMultiplier
				aura.Unit.PseudoStats.HolyDamageTakenMultiplier /= magicMultiplier
			}
		},
	})
}

var BloodFrenzyActionID = ActionID{SpellID: 29859}
var phyDmgDebuff = `4%phydmg`

func BloodFrenzyAura(target *Unit, points int32) *Aura {
	return bloodFrenzySavageCombatAura(target, "Blood Frenzy", BloodFrenzyActionID, points)
}
func SavageCombatAura(target *Unit, points int32) *Aura {
	return bloodFrenzySavageCombatAura(target, "Savage Combat", ActionID{SpellID: 58413}, points)
}

func bloodFrenzySavageCombatAura(target *Unit, label string, id ActionID, points int32) *Aura {
	multiplier := 1 + 0.02*float64(points)
	return target.GetOrRegisterAura(Aura{
		Label:    label + "-" + strconv.Itoa(int(points)),
		Tag:      phyDmgDebuff,
		ActionID: id,
		Priority: multiplier,
		// No fixed duration, lasts as long as the bleed that activates it.
		Duration: NeverExpires,
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

const BleedDamageAuraTag = "BleedDamage"

func MangleAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Mangle",
		Tag:      BleedDamageAuraTag,
		ActionID: ActionID{SpellID: 33876},
		Duration: time.Minute,
		Priority: 1.3,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier *= 1.3
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier /= 1.3
		},
	})
}

func TraumaAura(target *Unit, points int) *Aura {
	multiplier := 1 + 0.15*float64(points)
	return target.GetOrRegisterAura(Aura{
		Label:    "Trauma",
		Tag:      BleedDamageAuraTag,
		ActionID: ActionID{SpellID: 46855},
		Duration: time.Second * 60,
		Priority: multiplier,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier *= multiplier
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier /= multiplier
		},
	})
}

func StampedeAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Stampede",
		Tag:      BleedDamageAuraTag,
		ActionID: ActionID{SpellID: 57393},
		Duration: time.Second * 12,
		Priority: 1.25,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier *= 1.25
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier /= 1.25
		},
	})
}

const MajorSpellCritDebuffAuraTag = "majorspellcritdebuff"

func ShadowMasteryAura(target *Unit) *Aura {
	return majorSpellCritDebuffAura(target, "Shadow Mastery", ActionID{SpellID: 17800}, 5)
}

var ImprovedScorchAuraLabel = "Improved Scorch"

func ImprovedScorchAura(target *Unit) *Aura {
	return majorSpellCritDebuffAura(target, ImprovedScorchAuraLabel, ActionID{SpellID: 12873}, 5)
}

var WintersChillAuraLabel = "Winter's Chill"

func WintersChillAura(target *Unit, startingStacks int32) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:     WintersChillAuraLabel,
		Tag:       MajorSpellCritDebuffAuraTag,
		Priority:  1,
		ActionID:  ActionID{SpellID: 28595},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, startingStacks)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks, newStacks int32) {
			aura.Unit.PseudoStats.BonusSpellCritRatingTaken += float64(newStacks-oldStacks) * CritRatingPerCritChance
		},
	})
}

func majorSpellCritDebuffAura(target *Unit, label string, actionID ActionID, percent float64) *Aura {
	bonusSpellCrit := percent * CritRatingPerCritChance
	return target.GetOrRegisterAura(Aura{
		Label:    label,
		Tag:      MajorSpellCritDebuffAuraTag,
		Priority: percent,
		ActionID: actionID,
		Duration: time.Second * 30,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusSpellCritRatingTaken += bonusSpellCrit
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusSpellCritRatingTaken -= bonusSpellCrit
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
				aura.Unit.PseudoStats.BonusSpellHitRatingTaken += 3 * SpellHitRatingPerHitChance
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				aura.Unit.PseudoStats.BonusSpellHitRatingTaken -= 3 * SpellHitRatingPerHitChance
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
			aura.Unit.PseudoStats.ArmorMultiplier *= 1.0 - armorReduction
			aura.Unit.updateArmor()
			if imp {
				secondaryAura.Activate(sim)
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier /= 1.0 - armorReduction
			aura.Unit.updateArmor()
			if imp && secondaryAura.IsActive() {
				secondaryAura.Deactivate(sim)
			}
		},
	})
	return mainAura
}

var SunderArmorAuraLabel = "Sunder Armor"
var MajorArmorReductionTag = "MajorArmorReductionAura"
var SunderArmorActionID = ActionID{SpellID: 47467}

func SunderArmorAura(target *Unit, startingStacks int32) *Aura {
	armorReductionPerStack := 0.04

	return target.GetOrRegisterAura(Aura{
		Label:     SunderArmorAuraLabel,
		Tag:       MajorArmorReductionTag,
		ActionID:  SunderArmorActionID,
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
var AcidSpitActionID = ActionID{SpellID: 55754}

func AcidSpitAura(target *Unit, startingStacks int32) *Aura {
	armorReductionPerStack := 0.1

	return target.GetOrRegisterAura(Aura{
		Label:     AcidSpitAuraLabel,
		Tag:       MajorArmorReductionTag,
		ActionID:  AcidSpitActionID,
		Duration:  time.Second * 10,
		MaxStacks: 2,
		Priority:  armorReductionPerStack * 2,
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
			aura.Unit.PseudoStats.ArmorMultiplier *= 1.0 - armorReduction
			aura.Unit.updateArmor()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, bonus.Multiply(-1))
			aura.Unit.PseudoStats.ArmorMultiplier /= 1.0 - armorReduction
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
			aura.Unit.PseudoStats.ArmorMultiplier *= 1.0 - armorReduction
			aura.Unit.updateArmor()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier /= 1.0 - armorReduction
			aura.Unit.updateArmor()
		},
	})
}

func SporeCloudAura(target *Unit) *Aura {
	armorReduction := 0.03

	return target.GetOrRegisterAura(Aura{
		Label:    "Spore Cloud",
		Tag:      MinorArmorReductionAuraTag,
		ActionID: ActionID{SpellID: 53598},
		Priority: armorReduction,
		Duration: time.Second * 9,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= 1.0 - armorReduction
			aura.Unit.updateArmor()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier /= 1.0 - armorReduction
			aura.Unit.updateArmor()
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
			aura.Unit.updateArmor()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 / (1.0 - armorReduction))
			aura.Unit.updateArmor()
		},
	})
}

var HuntersMarkAuraTag = "HuntersMark"

func HuntersMarkAura(target *Unit, points int32, glyphed bool) *Aura {
	bonus := 500.0 * (1 + 0.1*float64(points))
	priority := float64(points)

	if glyphed {
		bonus += 500.0 * 0.2
		priority += 2
	}

	return target.GetOrRegisterAura(Aura{
		Label:    "HuntersMark-" + strconv.Itoa(int(priority)),
		Tag:      HuntersMarkAuraTag,
		ActionID: ActionID{SpellID: 53338},
		Duration: NeverExpires,
		Priority: priority,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusRangedAttackPowerTaken += bonus
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusRangedAttackPowerTaken -= bonus
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

func VindicationAura(target *Unit) *Aura {
	apReduction := 574.0

	return target.GetOrRegisterAura(Aura{
		Label:    "Vindication",
		Tag:      APReductionAuraTag,
		ActionID: ActionID{SpellID: 26016},
		Duration: time.Second * 10,
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

const AtkSpeedReductionAuraTag = "AtkSpdReduction"

func ThunderClapAura(target *Unit, points int32) *Aura {
	speedMultiplier := []float64{0.9, 0.86, 0.83, 0.8}[points]
	inverseMult := 1 / speedMultiplier

	return target.GetOrRegisterAura(Aura{
		Label:    "ThunderClap-" + strconv.Itoa(int(points)),
		Tag:      AtkSpeedReductionAuraTag,
		ActionID: ActionID{SpellID: 47502},
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

func InfectedWoundsAura(target *Unit, points int32) *Aura {
	speedMultiplier := []float64{1.0, 0.94, 0.86, 0.80}[points]
	inverseMult := 1 / speedMultiplier

	return target.GetOrRegisterAura(Aura{
		Label:    "InfectedWounds-" + strconv.Itoa(int(points)),
		Tag:      AtkSpeedReductionAuraTag,
		ActionID: ActionID{SpellID: 48485},
		Duration: time.Second * 12,
		Priority: inverseMult,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, speedMultiplier)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, inverseMult)
		},
	})
}

// Note: Paladin code might apply this as part of their judgement auras instead
// of using another separate aura.
func JudgementsOfTheJustAura(target *Unit, points int32) *Aura {
	speedMultiplier := 1.0 - 0.1*float64(points)
	inverseMult := 1 / speedMultiplier

	return target.GetOrRegisterAura(Aura{
		Label:    "JudgementsOfTheJust-" + strconv.Itoa(int(points)),
		Tag:      AtkSpeedReductionAuraTag,
		ActionID: ActionID{SpellID: 53696},
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

func FrostFeverAura(target *Unit, impIcyTouch int32) *Aura {
	speedMultiplier := 0.86 - 0.02*float64(impIcyTouch)

	inverseMult := 1 / speedMultiplier
	return target.GetOrRegisterAura(Aura{
		Label:    "FrostFever",
		Tag:      AtkSpeedReductionAuraTag,
		ActionID: ActionID{SpellID: 55095},
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

const MarkOfBloodTag = "MarkOfBlood"

func MarkOfBloodAura(target *Unit) *Aura {
	actionId := ActionID{SpellID: 49005}

	var healthMetrics *ResourceMetrics
	aura := target.GetOrRegisterAura(Aura{
		Label:     "MarkOfBlood",
		Tag:       MarkOfBloodTag,
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
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
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

const RuneOfRazoriceVulnerabilityTag = "RuneOfRazoriceVulnerability"

func RuneOfRazoriceVulnerabilityAura(target *Unit) *Aura {
	frostVulnPerStack := 0.02
	aura := target.GetOrRegisterAura(Aura{
		Label:     "RuneOfRazoriceVulnerability",
		Tag:       RuneOfRazoriceVulnerabilityTag,
		ActionID:  ActionID{SpellID: 50401},
		Duration:  NeverExpires,
		MaxStacks: 5,
		Priority:  1.0 / 1.1,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, 0)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			oldMultiplier := 1.0 + float64(oldStacks)*frostVulnPerStack
			newMultiplier := 1.0 + float64(newStacks)*frostVulnPerStack
			aura.Unit.PseudoStats.FrostDamageTakenMultiplier *= newMultiplier / oldMultiplier
		},
	})
	return aura
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

const MinorCritDebuffAuraTag = "minorcritdebuff"

func TotemOfWrathDebuff(target *Unit) *Aura {
	return minorCritDebuffAura(target, "Totem of Wrath Debuff", ActionID{SpellID: 30708}, 3, time.Minute*5)
}

func MasterPoisonerDebuff(target *Unit, points float64) *Aura {
	return minorCritDebuffAura(target, "Master Poisoner", ActionID{SpellID: 58410}, points, time.Second*20)
}

func HeartoftheCrusaderDebuff(target *Unit, points float64) *Aura {
	return minorCritDebuffAura(target, "Heart of the Crusader", ActionID{SpellID: 20337}, points, time.Second*20)
}

func minorCritDebuffAura(target *Unit, label string, actionID ActionID, points float64, duration time.Duration) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    label,
		Tag:      MinorCritDebuffAuraTag,
		Priority: points,
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusCritRatingTaken += points * CritRatingPerCritChance
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusCritRatingTaken -= points * CritRatingPerCritChance
		},
	})
}
