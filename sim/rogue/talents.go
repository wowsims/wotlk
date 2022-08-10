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
	rogue.applySealFate()
	rogue.applyWeaponSpecializations()
	rogue.applyCombatPotency()
	rogue.applyFocusedAttacks()

	rogue.AddStat(stats.Dodge, core.DodgeRatingPerDodgeChance*2*float64(rogue.Talents.LightningReflexes))
	rogue.AddStat(stats.MeleeHaste, core.HasteRatingPerHastePercent*[]float64{0, 3, 6, 10}[rogue.Talents.LightningReflexes])
	rogue.AddStat(stats.Parry, core.ParryRatingPerParryChance*2*float64(rogue.Talents.Deflection))
	rogue.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(rogue.Talents.Malice))
	rogue.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*1*float64(rogue.Talents.Precision))
	rogue.AddStat(stats.Expertise, core.ExpertisePerQuarterPercentReduction*5*float64(rogue.Talents.WeaponExpertise))
	rogue.AddStat(stats.ArmorPenetration, core.ArmorPenPerPercentArmor*3*float64(rogue.Talents.SerratedBlades))

	if rogue.Talents.DualWieldSpecialization > 0 {
		rogue.AutoAttacks.OHEffect.BaseDamage.Calculator = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1+0.1*float64(rogue.Talents.DualWieldSpecialization), 1.0, true)
	}

	if rogue.Talents.Deadliness > 0 {
		rogue.AddStatDependency(stats.AttackPower, stats.AttackPower, 1.0+0.02*float64(rogue.Talents.Deadliness))
	}

	if rogue.Talents.SavageCombat > 0 {
		rogue.AddStatDependency(stats.AttackPower, stats.AttackPower, 1.0+0.02*float64(rogue.Talents.SavageCombat))
	}

	if rogue.Talents.SinisterCalling > 0 {
		rogue.AddStatDependency(stats.Agility, stats.Agility, 1.0+0.03*float64(rogue.Talents.SinisterCalling))
	}

	rogue.PseudoStats.AgentReserved1DamageDealtMultiplier *= (1 + float64(rogue.Talents.FindWeakness)*0.02)

	rogue.registerOverkillCD()
	rogue.registerHungerForBlood()
	rogue.registerColdBloodCD()
	rogue.registerBladeFlurryCD()
	rogue.registerAdrenalineRushCD()
	rogue.registerKillingSpreeCD()
}

func (rogue *Rogue) makeFinishingMoveEffectApplier() func(sim *core.Simulation, numPoints int32) {
	ruthlessnessChance := 0.2 * float64(rogue.Talents.Ruthlessness)
	ruthlessnessMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 14161})

	relentlessStrikes := rogue.Talents.RelentlessStrikes
	relentlessStrikesMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 14179})
	netherblade4pc := rogue.HasSetBonus(ItemSetNetherblade, 4)
	netherblade4pcMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 37168})

	return func(sim *core.Simulation, numPoints int32) {
		if ruthlessnessChance > 0 && sim.RandomFloat("Ruthlessness") < ruthlessnessChance {
			rogue.AddComboPoints(sim, 1, ruthlessnessMetrics)
		}
		if netherblade4pc && sim.RandomFloat("Netherblade 4pc") < 0.15 {
			rogue.AddComboPoints(sim, 1, netherblade4pcMetrics)
		}
		if relentlessStrikes > 0 {
			if sim.RandomFloat("RelentlessStrikes") < 0.04*float64(numPoints) {
				rogue.AddEnergy(sim, 25, relentlessStrikesMetrics)
			}
		}
	}
}

func (rogue *Rogue) makeCastModifier() func(*core.Simulation, *core.Spell, *core.Cast) {
	builderCostMultiplier := 1.0
	costMultiplier := 1.0
	costReduction := 40.0
	hasDeathmantle := rogue.HasSetBonus(ItemSetDeathmantle, 4)
	if rogue.HasSetBonus(ItemSetBonescythe, 4) {
		builderCostMultiplier -= 0.05
	}
	return func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
		if spell.Flags.Matches(SpellFlagBuilder) {
			costMultiplier *= builderCostMultiplier
		}
		if spell.Flags.Matches(SpellFlagFinisher) {
			if hasDeathmantle && rogue.DeathmantleProcAura.IsActive() {
				cast.Cost = 0
				rogue.DeathmantleProcAura.Deactivate(sim)
				return
			}
		}
		cast.Cost *= costMultiplier
		cast.Cost = math.Ceil(cast.Cost)
		if rogue.VanCleefsProcAura.IsActive() {
			cast.Cost = core.MaxFloat(0, cast.Cost-costReduction)
			rogue.VanCleefsProcAura.Deactivate(sim)
		}
	}
}

func (rogue *Rogue) applyMurder() {
	rogue.PseudoStats.DamageDealtMultiplier *= rogue.murderMultiplier()
}

func (rogue *Rogue) murderMultiplier() float64 {
	return 1.0 + 0.02*float64(rogue.Talents.Murder)
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
		ActionID:     actionID,
		ResourceType: stats.Energy,
		BaseCost:     15,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: 15,
				GCD:  time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.HungerForBloodAura.Activate(sim)
		},
	})

}

func (rogue *Rogue) preyOnTheWeakMultiplier(target *core.Unit) float64 {
	// TODO: Use the following predicate if/when health values are modeled
	//if rogue.CurrentTarget != nil &&
	//rogue.CurrentTarget.HasHealthBar() &&
	//rogue.CurrentTarget.CurrentHealthPercent() < rogue.CurrentHealthPercent()
	return (1 + 0.04*float64(rogue.Talents.PreyOnTheWeak))
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
			aura.Unit.PseudoStats.BonusCritRatingAgentReserved1 += 100 * core.CritRatingPerCritChance
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.BonusCritRatingAgentReserved1 -= 100 * core.CritRatingPerCritChance
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			aura.Deactivate(sim)
		},
	})

	coldBloodSpell := rogue.RegisterSpell(core.SpellConfig{
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
		Spell: coldBloodSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (rogue *Rogue) applySealFate() {
	if rogue.Talents.SealFate == 0 {
		return
	}

	procChance := 0.2 * float64(rogue.Talents.SealFate)
	cpMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 14195})

	rogue.RegisterAura(core.Aura{
		Label:    "Seal Fate",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spell.Flags.Matches(SpellFlagBuilder) {
				return
			}

			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if procChance == 1 || sim.RandomFloat("Seal Fate") < procChance {
				rogue.AddComboPoints(sim, 1, cpMetrics)
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
			rogue.PseudoStats.BonusMHCritRating += 1 * core.CritRatingPerCritChance * float64(rogue.Talents.CloseQuartersCombat)
		case proto.WeaponType_WeaponTypeMace:
			rogue.PseudoStats.BonusMHArmorPenRating += 3 * core.ArmorPenPerPercentArmor * float64(rogue.Talents.MaceSpecialization)
		}
	}
	if ohWeapon != nil && ohWeapon.ID != 0 {
		switch ohWeapon.WeaponType {
		case proto.WeaponType_WeaponTypeSword, proto.WeaponType_WeaponTypeAxe:
			hackAndSlashMask |= core.ProcMaskMeleeOH
		case proto.WeaponType_WeaponTypeDagger, proto.WeaponType_WeaponTypeFist:
			rogue.PseudoStats.BonusOHCritRating += 1 * core.CritRatingPerCritChance * float64(rogue.Talents.CloseQuartersCombat)
		case proto.WeaponType_WeaponTypeMace:
			rogue.PseudoStats.BonusOHArmorPenRating += 3 * core.ArmorPenPerPercentArmor * float64(rogue.Talents.MaceSpecialization)
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}

			// https://wotlk.wowhead.com/spell=35553/combat-potency, proc mask = 8838608.
			if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOH) {
				return
			}

			if sim.RandomFloat("Combat Potency") > procChance {
				return
			}

			rogue.AddEnergy(sim, energyBonus, energyMetrics)
		},
	})
}

func (rogue *Rogue) applyFocusedAttacks() {
	if rogue.Talents.FocusedAttacks == 0 {
		return
	}

	procChance := float64(rogue.Talents.FocusedAttacks) / 3.0
	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 51637})

	rogue.RegisterAura(core.Aura{
		Label:    "Focused Attacks",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) || !spellEffect.DidCrit() {
				return
			}
			if procChance == 1 || sim.RandomFloat("Focused Attacks") <= procChance {
				rogue.AddEnergy(sim, 2, energyMetrics)
			}
		},
	})
}

func (rogue *Rogue) registerBladeFlurryCD() {
	if !rogue.Talents.BladeFlurry {
		return
	}

	actionID := core.ActionID{SpellID: 13877}

	var curDmg float64
	bfHit := rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		ApplyEffects: core.ApplyEffectFuncDirectDamageTargetModifiersOnly(core.SpellEffect{
			// No proc mask, so it won't proc itself.
			ProcMask: core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(_ *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
					return curDmg
				},
			},
			OutcomeApplier: rogue.OutcomeFuncAlwaysHit(),
		}),
	})

	const hasteBonus = 1.2
	const inverseHasteBonus = 1 / 1.2
	const energyCost = 25.0

	dur := time.Second * 15

	rogue.BladeFlurryAura = rogue.RegisterAura(core.Aura{
		Label:    "Blade Flurry",
		ActionID: actionID,
		Duration: dur,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, hasteBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, inverseHasteBonus)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if sim.GetNumTargets() < 2 {
				return
			}
			if spellEffect.Damage == 0 || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			// Undo armor reduction to get the raw damage value.
			curDmg = spellEffect.Damage / rogue.AttackTables[spellEffect.Target.Index].ArmorDamageModifier

			bfHit.Cast(sim, rogue.Env.NextTargetUnit(spellEffect.Target))
			bfHit.SpellMetrics[spellEffect.Target.TableIndex].Casts--
		},
	})

	cooldownDur := time.Minute * 2
	bladeFlurrySpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Energy,
		BaseCost:     energyCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: energyCost,
				GCD:  time.Second,
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
		Spell:    bladeFlurrySpell,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityLow,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return rogue.CurrentEnergy() >= energyCost
		},
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if sim.GetRemainingDuration() > cooldownDur+dur {
				// We'll have enough time to cast another BF, so use it immediately to make sure we get the 2nd one.
				return true
			}

			// Since this is our last BF, wait until we have SND / procs up.
			sndTimeRemaining := rogue.SliceAndDiceAura.RemainingDuration(sim)
			if sndTimeRemaining >= time.Second {
				return true
			}

			// TODO: Wait for dst/mongoose procs

			return false
		},
	})
}

func (rogue *Rogue) registerAdrenalineRushCD() {
	if !rogue.Talents.AdrenalineRush {
		return
	}

	actionID := core.ActionID{SpellID: 13750}

	rogue.AdrenalineRushAura = rogue.RegisterAura(core.Aura{
		Label:    "Adrenaline Rush",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.ResetEnergyTick(sim)
			rogue.EnergyTickMultiplier *= 2.0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.ResetEnergyTick(sim)
			rogue.EnergyTickMultiplier *= 1 / 2.0
		},
	})

	adrenalineRushSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

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
			// Make sure we have plenty of room so the big ticks dont get wasted.
			thresh := 85.0
			if rogue.NextEnergyTickAt() < sim.CurrentTime+time.Second*1 {
				thresh = 60.0
			}
			if rogue.CurrentEnergy() > thresh {
				return false
			}
			return true
		},
	})
}

func (rogue *Rogue) registerKillingSpreeCD() {
	if !rogue.Talents.KillingSpree {
		return
	}
	rogue.registerKillingSpreeSpell()
}
