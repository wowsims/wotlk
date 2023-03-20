package rogue

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) ApplyTalents() {
	rogue.applyMurder()
	rogue.applySlaughterFromTheShadows()
	rogue.applySealFate()
	rogue.applyWeaponSpecializations()
	rogue.applyCombatPotency()
	rogue.applyFocusedAttacks()
	rogue.applyInitiative()

	rogue.AddStat(stats.Dodge, core.DodgeRatingPerDodgeChance*2*float64(rogue.Talents.LightningReflexes))
	rogue.PseudoStats.MeleeSpeedMultiplier *= []float64{1, 1.03, 1.06, 1.10}[rogue.Talents.LightningReflexes]
	rogue.AddStat(stats.Parry, core.ParryRatingPerParryChance*2*float64(rogue.Talents.Deflection))
	rogue.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(rogue.Talents.Malice))
	rogue.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*1*float64(rogue.Talents.Precision))
	rogue.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*1*float64(rogue.Talents.Precision))
	rogue.AddStat(stats.Expertise, core.ExpertisePerQuarterPercentReduction*5*float64(rogue.Talents.WeaponExpertise))
	rogue.AddStat(stats.ArmorPenetration, core.ArmorPenPerPercentArmor*3*float64(rogue.Talents.SerratedBlades))
	rogue.AutoAttacks.OHConfig.DamageMultiplier *= rogue.dwsMultiplier()

	if rogue.Talents.Deadliness > 0 {
		rogue.MultiplyStat(stats.AttackPower, 1.0+0.02*float64(rogue.Talents.Deadliness))
	}

	if rogue.Talents.SavageCombat > 0 {
		rogue.MultiplyStat(stats.AttackPower, 1.0+0.02*float64(rogue.Talents.SavageCombat))
	}

	if rogue.Talents.SinisterCalling > 0 {
		rogue.MultiplyStat(stats.Agility, 1.0+0.03*float64(rogue.Talents.SinisterCalling))
	}

	rogue.registerOverkillCD()
	rogue.registerHungerForBlood()
	rogue.registerColdBloodCD()
	rogue.registerBladeFlurryCD()
	rogue.registerAdrenalineRushCD()
	rogue.registerKillingSpreeCD()
	rogue.registerShadowstepCD()
	rogue.registerShadowDanceCD()
	rogue.registerMasterOfSubtletyCD()
	rogue.registerPreparationCD()
	rogue.registerPremeditation()
	rogue.registerGhostlyStrikeSpell()
	rogue.registerDirtyDeeds()
	rogue.registerHonorAmongThieves()
}

// dwsMultiplier returns the offhand damage multiplier
func (rogue *Rogue) dwsMultiplier() float64 {
	return 1 + 0.1*float64(rogue.Talents.DualWieldSpecialization)
}

func getRelentlessStrikesSpellID(talentPoints int32) int32 {
	if talentPoints == 1 {
		return 14179
	}
	return 58420 + talentPoints
}

func (rogue *Rogue) makeFinishingMoveEffectApplier() func(sim *core.Simulation, numPoints int32) {
	ruthlessnessMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 14161})
	relentlessStrikesMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: getRelentlessStrikesSpellID(rogue.Talents.RelentlessStrikes)})

	return func(sim *core.Simulation, numPoints int32) {
		if t := rogue.Talents.Ruthlessness; t > 0 {
			if sim.RandomFloat("Ruthlessness") < 0.2*float64(t) {
				rogue.AddComboPoints(sim, 1, ruthlessnessMetrics)
			}
		}
		if t := rogue.Talents.RelentlessStrikes; t > 0 {
			if sim.RandomFloat("RelentlessStrikes") < 0.04*float64(t*numPoints) {
				rogue.AddEnergy(sim, 25, relentlessStrikesMetrics)
			}
		}
	}
}

func (rogue *Rogue) makeCostModifier() func(baseCost float64) float64 {
	if rogue.HasSetBonus(ItemSetBonescythe, 4) {
		return func(baseCost float64) float64 {
			return math.RoundToEven(0.95 * baseCost)
		}
	}
	return func(baseCost float64) float64 {
		return baseCost
	}
}

func (rogue *Rogue) applyMurder() {
	rogue.PseudoStats.DamageDealtMultiplier *= rogue.murderMultiplier()
}

func (rogue *Rogue) murderMultiplier() float64 {
	return 1.0 + 0.02*float64(rogue.Talents.Murder)
}

func (rogue *Rogue) applySlaughterFromTheShadows() {
	rogue.PseudoStats.DamageDealtMultiplier *= rogue.slaughterFromTheShadowsMultiplier()
}

func (rogue *Rogue) slaughterFromTheShadowsMultiplier() float64 {
	return 1.0 + 0.01*float64(rogue.Talents.SlaughterFromTheShadows)
}

func (rogue *Rogue) registerHungerForBlood() {
	if !rogue.Talents.HungerForBlood {
		return
	}
	actionID := core.ActionID{SpellID: 51662}
	multiplier := 1.05
	if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfHungerForBlood) {
		multiplier += 0.03
	}
	rogue.HungerForBloodAura = rogue.RegisterAura(core.Aura{
		Label:    "Hunger for Blood",
		ActionID: actionID,
		Duration: time.Minute,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier *= multiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier *= 1 / multiplier
		},
	})

	rogue.HungerForBlood = rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		EnergyCost: core.EnergyCostOptions{
			Cost: 15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.HungerForBloodAura.Activate(sim)
		},
	})
}

func (rogue *Rogue) preyOnTheWeakMultiplier(_ *core.Unit) float64 {
	// TODO: Use the following predicate if/when health values are modeled,
	//  but note that this would have to be applied dynamically in that case.
	//if rogue.CurrentTarget != nil &&
	//rogue.CurrentTarget.HasHealthBar() &&
	//rogue.CurrentTarget.CurrentHealthPercent() < rogue.CurrentHealthPercent()
	return 1 + 0.04*float64(rogue.Talents.PreyOnTheWeak)
}

func (rogue *Rogue) registerDirtyDeeds() {
	if rogue.Talents.DirtyDeeds == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 14083}

	rogue.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int) {
			if isExecute == 35 {
				rogue.DirtyDeedsAura.Activate(sim)
			}
		})
	})

	rogue.DirtyDeedsAura = rogue.RegisterAura(core.Aura{
		Label:    "Dirty Deeds",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range rogue.Spellbook {
				if spell.Flags.Matches(SpellFlagBuilder|SpellFlagFinisher) && spell.DamageMultiplier > 0 {
					spell.DamageMultiplier *= rogue.DirtyDeedsMultiplier()
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range rogue.Spellbook {
				if spell.Flags.Matches(SpellFlagBuilder|SpellFlagFinisher) && spell.DamageMultiplier > 0 {
					spell.DamageMultiplier /= rogue.DirtyDeedsMultiplier()
				}
			}
		},
	})
}

func (rogue *Rogue) DirtyDeedsMultiplier() float64 {
	if rogue.Talents.DirtyDeeds == 0 {
		return 1
	}

	return 1 + 0.1*float64(rogue.Talents.DirtyDeeds)
}

func (rogue *Rogue) registerColdBloodCD() {
	if !rogue.Talents.ColdBlood {
		return
	}

	actionID := core.ActionID{SpellID: 14177}

	coldBloodAura := rogue.RegisterAura(core.Aura{
		Label:    "Cold Blood",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range rogue.Spellbook {
				if spell.Flags.Matches(SpellFlagColdBlooded) {
					spell.BonusCritRating += 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range rogue.Spellbook {
				if spell.Flags.Matches(SpellFlagColdBlooded) {
					spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// for Fan of Knives and Mutilate, the offhand hit comes first and is ignored, so the aura doesn't fade too early
			if spell.Flags.Matches(SpellFlagColdBlooded) && spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
				aura.Deactivate(sim)
			}
		},
	})

	rogue.ColdBlood = rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			coldBloodAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: rogue.ColdBlood,
		Type:  core.CooldownTypeDPS,
	})
}

func (rogue *Rogue) applySealFate() {
	if rogue.Talents.SealFate == 0 {
		return
	}

	procChance := 0.2 * float64(rogue.Talents.SealFate)
	cpMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 14195})

	icd := core.Cooldown{
		Timer:    rogue.NewTimer(),
		Duration: 500 * time.Millisecond,
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Seal Fate",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagBuilder) {
				return
			}

			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if icd.IsReady(sim) && sim.Proc(procChance, "Seal Fate") {
				rogue.AddComboPoints(sim, 1, cpMetrics)
				icd.Use(sim)
			}
		},
	})
}

func (rogue *Rogue) applyInitiative() {
	if rogue.Talents.Initiative == 0 {
		return
	}

	procChance := []float64{0, 0.33, 0.66, 1.0}[rogue.Talents.Initiative]
	cpMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 13980})

	rogue.RegisterAura(core.Aura{
		Label:    "Initiative",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == rogue.Garrote || spell == rogue.Ambush {
				if result.Landed() {
					if sim.Proc(procChance, "Initiative") {
						rogue.AddComboPoints(sim, 1, cpMetrics)
					}
				}
			}
		},
	})
}

func (rogue *Rogue) applyWeaponSpecializations() {
	mhWeapon := rogue.GetMHWeapon()
	ohWeapon := rogue.GetOHWeapon()
	// https://wotlk.wowhead.com/spell=13964/sword-specialization, proc mask = 20.
	hackAndSlashMask := core.ProcMaskUnknown
	if mhWeapon != nil && mhWeapon.ID != 0 {
		switch mhWeapon.WeaponType {
		case proto.WeaponType_WeaponTypeSword, proto.WeaponType_WeaponTypeAxe:
			hackAndSlashMask |= core.ProcMaskMeleeMH
		case proto.WeaponType_WeaponTypeDagger, proto.WeaponType_WeaponTypeFist:
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
					spell.BonusCritRating += 1 * core.CritRatingPerCritChance * float64(rogue.Talents.CloseQuartersCombat)
				}
			})
		case proto.WeaponType_WeaponTypeMace:
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
					spell.BonusArmorPenRating += 3 * core.ArmorPenPerPercentArmor * float64(rogue.Talents.MaceSpecialization)
				}
			})
		}
	}
	if ohWeapon != nil && ohWeapon.ID != 0 {
		switch ohWeapon.WeaponType {
		case proto.WeaponType_WeaponTypeSword, proto.WeaponType_WeaponTypeAxe:
			hackAndSlashMask |= core.ProcMaskMeleeOH
		case proto.WeaponType_WeaponTypeDagger, proto.WeaponType_WeaponTypeFist:
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating += 1 * core.CritRatingPerCritChance * float64(rogue.Talents.CloseQuartersCombat)
				}
			})
		case proto.WeaponType_WeaponTypeMace:
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusArmorPenRating += 3 * core.ArmorPenPerPercentArmor * float64(rogue.Talents.MaceSpecialization)
				}
			})
		}
	}

	rogue.registerHackAndSlash(hackAndSlashMask)
}

func (rogue *Rogue) applyCombatPotency() {
	if rogue.Talents.CombatPotency == 0 {
		return
	}

	const procChance = 0.2
	energyBonus := 3.0 * float64(rogue.Talents.CombatPotency)
	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 35553})

	rogue.RegisterAura(core.Aura{
		Label:    "Combat Potency",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// from 3.0.3 patch notes: "Combat Potency: Now only works with auto attacks"
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOHAuto) {
				return
			}

			if sim.RandomFloat("Combat Potency") < procChance {
				rogue.AddEnergy(sim, energyBonus, energyMetrics)
			}
		},
	})
}

func (rogue *Rogue) applyFocusedAttacks() {
	if rogue.Talents.FocusedAttacks == 0 {
		return
	}

	procChance := []float64{0, 0.33, 0.66, 1}[rogue.Talents.FocusedAttacks]
	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 51637})

	rogue.RegisterAura(core.Aura{
		Label:    "Focused Attacks",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) || !result.DidCrit() {
				return
			}
			// Fan of Knives OH hits do not trigger focused attacks
			if spell.ProcMask.Matches(core.ProcMaskMeleeOH) && spell.IsSpellAction(FanOfKnivesSpellID) {
				return
			}
			if sim.Proc(procChance, "Focused Attacks") {
				rogue.AddEnergy(sim, 2, energyMetrics)
			}
		},
	})
}

var BladeFlurryActionID = core.ActionID{SpellID: 13877}

func (rogue *Rogue) registerBladeFlurryCD() {
	if !rogue.Talents.BladeFlurry {
		return
	}

	var curDmg float64
	bfHit := rogue.RegisterSpell(core.SpellConfig{
		ActionID:    BladeFlurryActionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty, // No proc mask, so it won't proc itself.
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreAttackerModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, curDmg, spell.OutcomeAlwaysHit)
		},
	})

	const hasteBonus = 1.2
	const inverseHasteBonus = 1 / 1.2

	dur := time.Second * 15

	rogue.BladeFlurryAura = rogue.RegisterAura(core.Aura{
		Label:    "Blade Flurry",
		ActionID: BladeFlurryActionID,
		Duration: dur,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, hasteBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, inverseHasteBonus)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.GetNumTargets() < 2 {
				return
			}
			if result.Damage == 0 || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}
			// Fan of Knives off-hand hits are not cloned
			if spell.IsSpellAction(FanOfKnivesSpellID) && spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
				return
			}

			// Undo armor reduction to get the raw damage value.
			curDmg = result.Damage / result.ResistanceMultiplier

			bfHit.Cast(sim, rogue.Env.NextTargetUnit(result.Target))
			bfHit.SpellMetrics[result.Target.UnitIndex].Casts--
		},
	})

	cooldownDur := time.Minute * 2
	rogue.BladeFlurry = rogue.RegisterSpell(core.SpellConfig{
		ActionID: BladeFlurryActionID,

		EnergyCost: core.EnergyCostOptions{
			Cost: core.TernaryFloat64(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfBladeFlurry), 0, 25),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: cooldownDur,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.BladeFlurryAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.BladeFlurry,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if sim.GetRemainingDuration() > cooldownDur+dur {
				// We'll have enough time to cast another BF, so use it immediately to make sure we get the 2nd one.
				return true
			}

			// Since this is our last BF, wait until we have SND / procs up.
			sndTimeRemaining := rogue.SliceAndDiceAura.RemainingDuration(sim)
			// TODO: Wait for dst/mongoose procs
			return sndTimeRemaining >= time.Second
		},
	})
}

var AdrenalineRushActionID = core.ActionID{SpellID: 13750}

func (rogue *Rogue) registerAdrenalineRushCD() {
	if !rogue.Talents.AdrenalineRush {
		return
	}

	rogue.AdrenalineRushAura = rogue.RegisterAura(core.Aura{
		Label:    "Adrenaline Rush",
		ActionID: AdrenalineRushActionID,
		Duration: core.TernaryDuration(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfAdrenalineRush), time.Second*20, time.Second*15),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.ResetEnergyTick(sim)
			rogue.ApplyEnergyTickMultiplier(1.0)
			rogue.rotationItems = rogue.planRotation(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.ResetEnergyTick(sim)
			rogue.ApplyEnergyTickMultiplier(-1.0)
			rogue.rotationItems = rogue.planRotation(sim)
		},
	})

	adrenalineRushSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID: AdrenalineRushActionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.AdrenalineRushAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    adrenalineRushSpell,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityBloodlust,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			thresh := 45.0
			return rogue.CurrentEnergy() <= thresh
		},
	})
}

func (rogue *Rogue) registerKillingSpreeCD() {
	if !rogue.Talents.KillingSpree {
		return
	}
	rogue.registerKillingSpreeSpell()
}

func (rogue *Rogue) registerHonorAmongThieves() {
	// When anyone in your group critically hits with a damage or healing spell or ability,
	// you have a [33%/66%/100%] chance to gain a combo point on your current target.
	// This effect cannot occur more than once per second.
	if rogue.Talents.HonorAmongThieves == 0 {
		return
	}

	procChance := []float64{0, 0.33, 0.66, 1}[rogue.Talents.HonorAmongThieves]
	comboMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 51701})
	honorAmongThievesID := core.ActionID{SpellID: 51701}

	icd := core.Cooldown{
		Timer:    rogue.NewTimer(),
		Duration: time.Second,
	}

	maybeProc := func(sim *core.Simulation) {
		if icd.IsReady(sim) && sim.Proc(procChance, "honor of thieves") {
			rogue.AddComboPoints(sim, 1, comboMetrics)
			icd.Use(sim)
		}
	}

	rogue.HonorAmongThieves = rogue.RegisterAura(core.Aura{
		Label:    "Honor Among Thieves",
		ActionID: honorAmongThievesID,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			// In an ideal party, you'd probably get up to 6 ability crits/s (Rate = 600).
			//  Survival Hunters, Enhancement Shamans, and Assassination Rogues are particularly good.
			if rogue.Options.HonorOfThievesCritRate <= 0 {
				return
			}

			if rogue.Options.HonorOfThievesCritRate > 2000 {
				rogue.Options.HonorOfThievesCritRate = 2000 // limited, so performance doesn't suffer
			}

			rateToDuration := float64(time.Second) * 100 / float64(rogue.Options.HonorOfThievesCritRate)

			pa := &core.PendingAction{}
			pa.OnAction = func(sim *core.Simulation) {
				maybeProc(sim)
				pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomExpFloat("next party crit")*rateToDuration)
				sim.AddPendingAction(pa)
			}
			pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomExpFloat("next party crit")*rateToDuration)
			sim.AddPendingAction(pa)
		},
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() && !spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto|core.ProcMaskMeleeOHAuto|core.ProcMaskRangedAuto) {
				maybeProc(sim)
			}
		},
		OnPeriodicDamageDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() {
				maybeProc(sim)
			}
		},
	})
}
