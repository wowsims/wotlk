package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) ApplyTalents() {
	mage.applyArcaneConcentration()
	mage.applyFocusMagic()
	mage.applyIgnite()
	mage.applyEmpoweredFire()
	mage.applyMasterOfElements()
	mage.applyWintersChill()
	mage.applyMoltenFury()
	mage.applyMissileBarrage()
	mage.applyHotStreak()
	mage.applyFingersOfFrost()
	mage.applyBrainFreeze()
	mage.registerArcanePowerCD()
	mage.registerPresenceOfMindCD()
	mage.registerCombustionCD()
	mage.registerIcyVeinsCD()
	mage.registerColdSnapCD()
	mage.registerSummonWaterElementalCD()

	mage.PseudoStats.SpiritRegenRateCasting += float64(mage.Talents.ArcaneMeditation) / 6

	if mage.Talents.StudentOfTheMind > 0 {
		mage.MultiplyStat(stats.Spirit, 1.0+[]float64{0, .04, .07, .10}[mage.Talents.StudentOfTheMind])
	}

	if mage.Talents.ArcaneMind > 0 {
		mage.MultiplyStat(stats.Intellect, 1.0+0.03*float64(mage.Talents.ArcaneMind))
	}

	if mage.Talents.MindMastery > 0 {
		mage.AddStatDependency(stats.Intellect, stats.SpellPower, 0.03*float64(mage.Talents.MindMastery))
	}

	mage.AddStat(stats.SpellCrit, float64(mage.Talents.ArcaneInstability)*1*core.CritRatingPerCritChance)
	mage.PseudoStats.DamageDealtMultiplier *= 1 + .01*float64(mage.Talents.ArcaneInstability)
	mage.PseudoStats.DamageDealtMultiplier *= 1 + .01*float64(mage.Talents.PlayingWithFire)
	mage.PseudoStats.CastSpeedMultiplier *= 1 + .02*float64(mage.Talents.NetherwindPresence)

	mage.AddStat(stats.SpellCrit, float64(mage.Talents.Pyromaniac)*core.CritRatingPerCritChance)
	mage.PseudoStats.SpiritRegenRateCasting += float64(mage.Talents.Pyromaniac) / 6

	mage.AddStat(stats.SpellHit, float64(mage.Talents.Precision)*core.SpellHitRatingPerHitChance)
	mage.PseudoStats.CostMultiplier *= 1 - .01*float64(mage.Talents.Precision)

	mage.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1 + .02*float64(mage.Talents.PiercingIce)
	mage.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1 + .01*float64(mage.Talents.ArcticWinds)
	mage.PseudoStats.CostMultiplier *= 1 - .04*float64(mage.Talents.FrostChanneling)

	magicAbsorptionBonus := 2 * float64(mage.Talents.MagicAbsorption)
	mage.AddStat(stats.ArcaneResistance, magicAbsorptionBonus)
	mage.AddStat(stats.FireResistance, magicAbsorptionBonus)
	mage.AddStat(stats.FrostResistance, magicAbsorptionBonus)
	mage.AddStat(stats.NatureResistance, magicAbsorptionBonus)
	mage.AddStat(stats.ShadowResistance, magicAbsorptionBonus)
}

func (mage *Mage) applyHotStreak() {
	if mage.Talents.HotStreak == 0 {
		return
	}

	procChance := float64(mage.Talents.HotStreak) / 3

	mage.HotStreakAura = mage.RegisterAura(core.Aura{
		Label:    "HotStreak",
		ActionID: core.ActionID{SpellID: 44448},
		Duration: time.Second * 10,
		// This is handled in Pyroblast.ModifyCast instead.
		//OnGain: func(aura *core.Aura, sim *core.Simulation) {
		//	if mage.Pyroblast != nil {
		//		mage.Pyroblast.CastTimeMultiplier -= 1
		//	}
		//},
		//OnExpire: func(aura *core.Aura, sim *core.Simulation) {
		//	if mage.Pyroblast != nil {
		//		mage.Pyroblast.CastTimeMultiplier += 1
		//	}
		//},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Hot Streak Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
			mage.heatingUp = false
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(HotStreakSpells) {
				return
			}

			if !result.DidCrit() {
				mage.heatingUp = false
				return
			}

			if mage.heatingUp {
				if procChance == 1 || sim.Proc(procChance, "Hot Streak") {
					mage.HotStreakAura.Activate(sim)
					mage.heatingUp = false
				}
			} else {
				mage.heatingUp = true
			}
		},
	})

}

func (mage *Mage) applyArcaneConcentration() {
	if mage.Talents.ArcaneConcentration == 0 {
		return
	}

	procChance := 0.02 * float64(mage.Talents.ArcaneConcentration)
	bonusCrit := float64(mage.Talents.ArcanePotency) * 15 * core.CritRatingPerCritChance

	// The result that caused the proc. Used to check we don't deactivate from the same proc.
	var proccedAt time.Duration
	var proccedSpell *core.Spell

	if mage.Talents.ArcanePotency > 0 {
		mage.ArcanePotencyAura = mage.RegisterAura(core.Aura{
			Label:    "Arcane Potency",
			ActionID: core.ActionID{SpellID: 31572},
			Duration: time.Second * 15,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.SpellCrit, -bonusCrit)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !spell.Flags.Matches(SpellFlagMage) {
					return
				}
				if proccedAt == sim.CurrentTime && proccedSpell == spell {
					// Means this is another hit from the same cast that procced CC.
					return
				}
				aura.Deactivate(sim)
			},
		})
	}

	mage.ClearcastingAura = mage.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 12536},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.CostMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.CostMultiplier += 1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Flags.Matches(SpellFlagMage) {
				return
			}
			if spell == mage.ArcaneMissiles && mage.MissileBarrageAura.IsActive() {
				return
			}
			if proccedAt == sim.CurrentTime && proccedSpell == spell {
				// Means this is another hit from the same cast that procced CC.
				return
			}
			aura.Deactivate(sim)
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Arcane Concentration",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagMage) || spell == mage.ArcaneMissiles {
				return
			}

			if !result.Landed() {
				return
			}

			if sim.RandomFloat("Arcane Concentration") > procChance {
				return
			}

			proccedAt = sim.CurrentTime
			proccedSpell = spell
			mage.ClearcastingAura.Activate(sim)
			if mage.ArcanePotencyAura != nil {
				mage.ArcanePotencyAura.Activate(sim)
			}
		},
	})
}

func (mage *Mage) applyMissileBarrage() {
	if mage.Talents.MissileBarrage == 0 {
		return
	}

	procChance := float64(mage.Talents.MissileBarrage) * .04
	mage.MissileBarrageAura = mage.RegisterAura(core.Aura{
		Label:    "Missile Barrage Proc",
		ActionID: core.ActionID{SpellID: 44401},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.ArcaneMissiles.CostMultiplier -= 100
			mage.ArcaneMissiles.CastTimeMultiplier /= 2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.ArcaneMissiles.CostMultiplier += 100
			mage.ArcaneMissiles.CastTimeMultiplier *= 2
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Missile Barrage Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Flags.Matches(BarrageSpells) {
				return
			}

			roll := sim.RandomFloat("Missile Barrage")
			updChance := core.TernaryFloat64(spell.ActionID == mage.ArcaneBlast.ActionID, 2*procChance, procChance)

			if roll < updChance {
				mage.MissileBarrageAura.Activate(sim)
			}
		},
	})
}

func (mage *Mage) registerPresenceOfMindCD() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	cooldown := 120.0
	if mage.Talents.ArcaneFlows > 0 {
		cooldown *= 1 - (.15 * float64(mage.Talents.ArcaneFlows))
	}

	actionID := core.ActionID{SpellID: 12043}

	var spellToUse *core.Spell
	mage.Env.RegisterPostFinalizeEffect(func() {
		if mage.Pyroblast != nil {
			spellToUse = mage.Pyroblast
		} else if mage.Rotation.Type == proto.Mage_Rotation_Fire {
			spellToUse = mage.Fireball
		} else if mage.Rotation.Type == proto.Mage_Rotation_Frost {
			spellToUse = mage.Frostbolt
		} else {
			spellToUse = mage.ArcaneBlast
		}
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Duration(cooldown) * time.Second,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			if !mage.GCD.IsReady(sim) {
				return false
			}
			if mage.ArcanePowerAura.IsActive() {
				return false
			}

			manaCost := spellToUse.DefaultCast.Cost * mage.PseudoStats.CostMultiplier
			if spellToUse == mage.ArcaneBlast {
				manaCost *= float64(mage.ArcaneBlastAura.GetStacks()) * 1.75
			}

			return mage.CurrentMana() >= manaCost
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if mage.ArcanePotencyAura != nil {
				mage.ArcanePotencyAura.Activate(sim)
			}

			normalCastTime := spellToUse.DefaultCast.CastTime
			spellToUse.DefaultCast.CastTime = 0
			spellToUse.Cast(sim, mage.CurrentTarget)
			spellToUse.DefaultCast.CastTime = normalCastTime
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerArcanePowerCD() {
	if !mage.Talents.ArcanePower {
		return
	}
	actionID := core.ActionID{SpellID: 12042}

	var affectedSpells []*core.Spell
	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagMage) {
			affectedSpells = append(affectedSpells, spell)
		}
	})

	mage.ArcanePowerAura = mage.RegisterAura(core.Aura{
		Label:    "Arcane Power",
		ActionID: actionID,
		Duration: core.TernaryDuration(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcanePower), time.Second*18, time.Second*15),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += 0.2
				spell.CostMultiplier += 0.2
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive -= 0.2
				spell.CostMultiplier -= 0.2
			}
		},
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * time.Duration(120*(1-(.15*float64(mage.Talents.ArcaneFlows)))),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.ArcanePowerAura.Activate(sim)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			if mage.ArcanePotencyAura.IsActive() {
				return false
			}
			return true
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) applyMasterOfElements() {
	if mage.Talents.MasterOfElements == 0 && mage.Talents.Burnout == 0 {
		return
	}

	refundCoeff := 0.1*float64(mage.Talents.MasterOfElements) - .01*float64(mage.Talents.Burnout)
	manaMetrics := mage.NewManaMetrics(core.ActionID{SpellID: 29076})

	mage.RegisterAura(core.Aura{
		Label:    "Master of Elements",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				return
			}
			if result.DidCrit() {
				if refundCoeff < 0 {
					mage.SpendMana(sim, -1*spell.DefaultCast.Cost*refundCoeff, manaMetrics)
				} else {
					mage.AddMana(sim, spell.DefaultCast.Cost*refundCoeff, manaMetrics)
				}
			}
		},
	})
}

func (mage *Mage) registerCombustionCD() {
	if !mage.Talents.Combustion {
		return
	}
	actionID := core.ActionID{SpellID: 11129}
	cd := core.Cooldown{
		Timer:    mage.NewTimer(),
		Duration: time.Minute * 2,
	}

	fireCombCritMult := mage.SpellCritMultiplier(1, mage.bonusCritDamage+.5) / mage.SpellCritMultiplier(1, mage.bonusCritDamage)

	frostfireCombCritMult := mage.SpellCritMultiplier(1, mage.bonusCritDamage+float64(mage.Talents.IceShards)/3+.5) /
		mage.SpellCritMultiplier(1, mage.bonusCritDamage+float64(mage.Talents.IceShards)/3)

	var fireSpells []*core.Spell
	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellSchool.Matches(core.SpellSchoolFire) {
			fireSpells = append(fireSpells, spell)
		}
	})

	numCrits := 0
	const critPerStack = 10 * core.CritRatingPerCritChance

	mage.CombustionAura = mage.RegisterAura(core.Aura{
		Label:     "Combustion",
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			numCrits = 0
			for _, spell := range fireSpells {
				spell.CritMultiplier *= core.TernaryFloat64(spell != mage.FrostfireBolt, fireCombCritMult, frostfireCombCritMult)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			cd.Use(sim)
			mage.UpdateMajorCooldowns()
			for _, spell := range fireSpells {
				spell.CritMultiplier /= core.TernaryFloat64(spell != mage.FrostfireBolt, fireCombCritMult, frostfireCombCritMult)
			}
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			bonusCrit := critPerStack * float64(newStacks-oldStacks)
			for _, spell := range fireSpells {
				spell.BonusCritRating += bonusCrit
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.SpellSchool.Matches(core.SpellSchoolFire) || !spell.Flags.Matches(SpellFlagMage) {
				return
			}
			if spell == mage.Ignite || spell == mage.LivingBomb { //LB dot action should be ignored
				return
			}
			if !result.Landed() {
				return
			}
			if numCrits >= 3 {
				return
			}

			// TODO: This wont work properly with flamestrike
			aura.AddStack(sim)

			if result.DidCrit() {
				numCrits++
				if numCrits == 3 {
					aura.Deactivate(sim)
				}
			}
		},
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: cd,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mage.CombustionAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.CombustionAura.Activate(sim)
			mage.CombustionAura.AddStack(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerIcyVeinsCD() {
	if !mage.Talents.IcyVeins {
		return
	}

	actionID := core.ActionID{SpellID: 12472}

	icyVeinsAura := mage.RegisterAura(core.Aura{
		Label:    "Icy Veins",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(1.2)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(1 / 1.2)
		},
	})

	mage.IcyVeins = mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.03,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * time.Duration(180*[]float64{1, .93, .86, .80}[mage.Talents.IceFloes]),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Need to check for icy veins already active in case Cold Snap is used right after.
			return !icyVeinsAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			icyVeinsAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.IcyVeins,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerColdSnapCD() {
	if !mage.Talents.ColdSnap {
		return
	}

	cooldown := time.Duration(float64(time.Minute*8) * (1.0 - float64(mage.Talents.ColdAsIce)*0.1))
	actionID := core.ActionID{SpellID: 11958}

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: cooldown,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Don't use if there are no cooldowns to reset.
			return (mage.IcyVeins != nil && !mage.IcyVeins.IsReady(sim)) ||
				(mage.SummonWaterElemental != nil && !mage.SummonWaterElemental.IsReady(sim))
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if mage.IcyVeins != nil {
				mage.IcyVeins.CD.Reset()
			}
			if mage.SummonWaterElemental != nil {
				mage.SummonWaterElemental.CD.Reset()
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Ideally wait for both water ele and icy veins so we can reset both.
			if mage.IcyVeins != nil && mage.IcyVeins.IsReady(sim) {
				return false
			}
			if mage.SummonWaterElemental != nil && mage.SummonWaterElemental.IsReady(sim) {
				return false
			}

			return true
		},
	})
}

func (mage *Mage) applyMoltenFury() {
	if mage.Talents.MoltenFury == 0 {
		return
	}

	multiplier := 1.0 + 0.06*float64(mage.Talents.MoltenFury)

	mage.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int) {
			mage.GetMajorCooldown(core.ActionID{SpellID: EvocationId}).Disable()
			if isExecute == 35 {
				mage.PseudoStats.DamageDealtMultiplier *= multiplier
				// For some reason Molten Fury doesn't apply to living bomb DoT, so cancel it out.
				if mage.LivingBomb != nil {
					mage.LivingBomb.DamageMultiplier /= multiplier
				}
			}
		})
	})
}

func (mage *Mage) hasChillEffect(spell *core.Spell) bool {
	return spell == mage.Frostbolt || spell == mage.FrostfireBolt || (spell == mage.Blizzard && mage.Talents.ImprovedBlizzard > 0)
}

func (mage *Mage) applyFingersOfFrost() {
	if mage.Talents.FingersOfFrost == 0 {
		return
	}

	bonusCrit := []float64{0, 17, 34, 50}[mage.Talents.Shatter] * core.CritRatingPerCritChance
	iceLanceMultiplier := core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfIceLance), 4, 3)

	var proccedAt time.Duration

	mage.FingersOfFrostAura = mage.RegisterAura(core.Aura{
		Label:     "Fingers of Frost Proc",
		ActionID:  core.ActionID{SpellID: 44545},
		Duration:  time.Second * 15,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
			mage.IceLance.DamageMultiplier *= iceLanceMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, -bonusCrit)
			mage.IceLance.DamageMultiplier /= iceLanceMultiplier
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if proccedAt != sim.CurrentTime {
				aura.RemoveStack(sim)
			}
		},
	})

	procChance := []float64{0, .07, .15}[mage.Talents.FingersOfFrost]
	mage.RegisterAura(core.Aura{
		Label:    "Fingers of Frost Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if mage.hasChillEffect(spell) && sim.RandomFloat("Fingers of Frost") < procChance {
				mage.FingersOfFrostAura.Activate(sim)
				mage.FingersOfFrostAura.SetStacks(sim, 2)
				proccedAt = sim.CurrentTime
			}
		},
	})
}

func (mage *Mage) applyBrainFreeze() {
	if mage.Talents.BrainFreeze == 0 {
		return
	}

	hasT8_4pc := mage.HasSetBonus(ItemSetKirinTorGarb, 4)
	t10ProcAura := mage.BloodmagesRegalia2pcAura()

	mage.BrainFreezeAura = mage.RegisterAura(core.Aura{
		Label:    "Brain Freeze Proc",
		ActionID: core.ActionID{SpellID: 44549},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.Fireball.CostMultiplier -= 100
			mage.Fireball.CastTimeMultiplier -= 1
			mage.FrostfireBolt.CostMultiplier -= 100
			mage.FrostfireBolt.CastTimeMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.Fireball.CostMultiplier += 100
			mage.Fireball.CastTimeMultiplier += 1
			mage.FrostfireBolt.CostMultiplier += 100
			mage.FrostfireBolt.CastTimeMultiplier += 1
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == mage.FrostfireBolt || spell == mage.Fireball {
				if !hasT8_4pc || sim.RandomFloat("MageT84PC") > T84PcProcChance {
					aura.Deactivate(sim)
				}
				if t10ProcAura != nil {
					t10ProcAura.Activate(sim)
				}
			}
		},
	})

	procChance := .05 * float64(mage.Talents.BrainFreeze)
	mage.RegisterAura(core.Aura{
		Label:    "Brain Freeze Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if mage.hasChillEffect(spell) && sim.RandomFloat("Brain Freeze") < procChance {
				mage.BrainFreezeAura.Activate(sim)
			}
		},
	})
}

func (mage *Mage) applyWintersChill() {
	if mage.Talents.WintersChill == 0 {
		return
	}

	procChance := []float64{0, 0.33, 0.66, 1}[mage.Talents.WintersChill]

	wcAuras := make([]*core.Aura, len(mage.Env.Encounter.TargetUnits))
	for i, target := range mage.Env.Encounter.TargetUnits {
		wcAuras[i] = core.WintersChillAura(target, 0)
	}

	mage.RegisterAura(core.Aura{
		Label:    "Winters Chill Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.SpellSchool.Matches(core.SpellSchoolFrost) {
				return
			}

			if sim.Proc(procChance, "Winters Chill") {
				aura := wcAuras[result.Target.Index]
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			}
		},
	})
}
