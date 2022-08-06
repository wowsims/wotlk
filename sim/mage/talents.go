package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) ApplyTalents() {
	mage.applyArcaneConcentration()
	mage.applyIgnite()
	mage.applyMasterOfElements()
	mage.applyWintersChill()
	mage.applyMoltenFury()
	mage.applyMissileBarrage()
	mage.applyHotStreak()
	mage.registerArcanePowerCD()
	mage.registerPresenceOfMindCD()
	mage.registerCombustionCD()
	mage.registerIcyVeinsCD()
	mage.registerColdSnapCD()
	mage.registerSummonWaterElementalCD()
	// TODO: Enduring Winter

	mage.PseudoStats.SpiritRegenRateCasting += float64(mage.Talents.ArcaneMeditation) / 6

	if mage.Talents.StudentOfTheMind > 0 {
		mage.Character.AddStatDependency(stats.Spirit, stats.Spirit, 1.0+[]float64{0, .04, .07, .10}[mage.Talents.StudentOfTheMind])
	}

	if mage.Talents.FocusMagic {
		// TODO: Pretty sure this should be 2 separate effects? One from talent, another from generic raid buff from another mage.
		totalCritPercent := 3 + float64(mage.Options.FocusMagicPercentUptime*3)/100.0
		mage.AddStat(stats.SpellCrit, totalCritPercent*core.CritRatingPerCritChance)
	}

	if mage.Talents.ArcaneMind > 0 {
		mage.Character.AddStatDependency(stats.Intellect, stats.Intellect, 1.0+(0.03*float64(mage.Talents.ArcaneMind)))
	}

	if mage.Talents.MindMastery > 0 {
		mage.Character.AddStatDependency(stats.Intellect, stats.SpellPower, 1+0.03*float64(mage.Talents.MindMastery))
	}

	mage.AddStat(stats.SpellCrit, float64(mage.Talents.ArcaneInstability)*1*core.CritRatingPerCritChance)
	mage.spellDamageMultiplier += .01 * float64(mage.Talents.ArcaneInstability)

	mage.PseudoStats.CastSpeedMultiplier *= 1 + .02*float64(mage.Talents.NetherwindPresence)

	mage.spellDamageMultiplier += .01 * float64(mage.Talents.PlayingWithFire)
	mage.PseudoStats.FireDamageDealtMultiplier *= 1 + .02*float64(mage.Talents.FirePower)

	mage.AddStat(stats.SpellCrit, float64(mage.Talents.Pyromaniac)*core.CritRatingPerCritChance)
	mage.PseudoStats.SpiritRegenRateCasting += float64(mage.Talents.Pyromaniac) / 6

	if mage.Talents.SpellPower > 0 {
		mage.bonusCritDamage += .25 * float64(mage.Talents.SpellPower)
	}

	if mage.Talents.Burnout > 0 {
		mage.bonusCritDamage += .1 * float64(mage.Talents.Burnout)
	}

	mage.AddStat(stats.SpellHit, float64(mage.Talents.Precision)*core.SpellHitRatingPerHitChance)
	mage.PseudoStats.CostMultiplier *= 1 - .01*float64(mage.Talents.Precision)

	mage.PseudoStats.FrostDamageDealtMultiplier *= 1 + .02*float64(mage.Talents.PiercingIce)
	mage.PseudoStats.FrostDamageDealtMultiplier *= 1 + .01*float64(mage.Talents.ArcticWinds)
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
	})

	heatingUp := false
	mage.RegisterAura(core.Aura{
		Label:    "HeatingUp",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spell.Flags.Matches(HotStreakSpells) {
				return
			}

			if mage.HotStreakAura.IsActive() {
				return
			}

			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				heatingUp = false
				return
			} else {
				if heatingUp {
					if procChance == 1 || sim.RandomFloat("Hot Streak") < procChance {
						mage.HotStreakAura.Activate(sim)
						heatingUp = false
					}
				} else {
					heatingUp = true
				}
			}
		},
	})

}

func (mage *Mage) applyArcaneConcentration() {
	if mage.Talents.ArcaneConcentration == 0 {
		return
	}

	procChance := 0.02 * float64(mage.Talents.ArcaneConcentration)
	bonusCrit := float64(mage.Talents.ArcanePotency) * 10 * core.CritRatingPerCritChance

	// Used to make sure we don't try to roll twice for the same cast on aoe spells.
	var curCastIdx int
	var lastCheckedCastIdx int

	mage.ClearcastingAura = mage.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 12536},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
			mage.PseudoStats.NoCost = true
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, -bonusCrit)
			mage.PseudoStats.NoCost = false
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spell.Flags.Matches(SpellFlagMage) {
				return
			}
			if curCastIdx == lastCheckedCastIdx {
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
			curCastIdx = 0
			lastCheckedCastIdx = 0
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if mage.bonusAMCCCrit != 0 {
				mage.AddStatDynamic(sim, stats.SpellCrit, -mage.bonusAMCCCrit)
				mage.bonusAMCCCrit = 0
			}
			if !spell.Flags.Matches(SpellFlagMage) {
				return
			}
			curCastIdx++
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spell.Flags.Matches(SpellFlagMage) {
				return
			}

			if curCastIdx == lastCheckedCastIdx {
				// Means we already rolled for this cast.
				return
			}
			lastCheckedCastIdx = curCastIdx

			if !spellEffect.Landed() {
				return
			}

			if sim.RandomFloat("Arcane Concentration") > procChance {
				return
			}

			mage.ClearcastingAura.Activate(sim)
			mage.ClearcastingAura.Prioritize()
		},
	})
}

func (mage *Mage) applyMissileBarrage() {
	if mage.Talents.MissileBarrage == 0 {
		return
	}

	// countBarrageChances := 0
	missileBarrageActionId := core.ActionID{SpellID: 44401}

	procChance := float64(mage.Talents.MissileBarrage) * .04
	mage.MissileBarrageAura = mage.RegisterAura(core.Aura{
		Label:    "Missile Barrage Proc",
		ActionID: missileBarrageActionId,
		Duration: time.Second * 15,
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

			// countBarrageChances++
			// mage.Log(sim, "Total missile barrage opportunities %d", countBarrageChances)
			roll := sim.RandomFloat("Missile Barrage")

			updChance := core.TernaryFloat64(spell.ActionID == mage.ArcaneBlast.ActionID, 2*procChance, procChance)

			if roll < updChance {
				mage.MissileBarrageAura.Activate(sim)
				mage.MissileBarrageAura.Prioritize()
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

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Duration(cooldown) * time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			var spell *core.Spell
			if mage.Talents.Pyroblast {
				spell = mage.Pyroblast
			} else if mage.Rotation.Type == proto.Mage_Rotation_Fire {
				spell = mage.Fireball
			} else if mage.Rotation.Type == proto.Mage_Rotation_Frost {
				spell = mage.Frostbolt
			} else {
				spell = mage.ArcaneBlast
			}

			normalCastTime := spell.DefaultCast.CastTime
			spell.DefaultCast.CastTime = 0
			spell.Cast(sim, mage.CurrentTarget)
			spell.DefaultCast.CastTime = normalCastTime
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			var manaCost float64
			if mage.Talents.Pyroblast {
				manaCost = mage.Pyroblast.DefaultCast.Cost
			} else if mage.Rotation.Type == proto.Mage_Rotation_Fire {
				manaCost = mage.Fireball.DefaultCast.Cost
			} else if mage.Rotation.Type == proto.Mage_Rotation_Frost {
				manaCost = mage.Frostbolt.DefaultCast.Cost
			} else {
				manaCost = mage.ArcaneBlast.DefaultCast.Cost * float64(mage.ArcaneBlastAura.GetStacks()) * 1.75
			}
			manaCost *= character.PseudoStats.CostMultiplier

			if character.CurrentMana() < manaCost {
				return false
			}

			return true
		},
	})
}

func (mage *Mage) registerArcanePowerCD() {
	if !mage.Talents.ArcanePower {
		return
	}
	actionID := core.ActionID{SpellID: 12042}

	apAura := mage.RegisterAura(core.Aura{
		Label:    "Arcane Power",
		ActionID: actionID,
		Duration: core.TernaryDuration(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcanePower), time.Second*18, time.Second*15),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.PseudoStats.DamageDealtMultiplier *= 1.2
			mage.PseudoStats.CostMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.PseudoStats.DamageDealtMultiplier /= 1.2
			mage.PseudoStats.CostMultiplier /= 1.2
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
			apAura.Activate(sim)
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				return
			}
			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				if refundCoeff < 0 {
					mage.SpendMana(sim, spell.BaseCost*refundCoeff, manaMetrics)
				} else {
					mage.AddMana(sim, spell.BaseCost*refundCoeff, manaMetrics, false)
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
		Duration: time.Minute * 3,
	}

	numCrits := 0
	const critPerStack = 10 * core.CritRatingPerCritChance

	aura := mage.RegisterAura(core.Aura{
		Label:     "Combustion",
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			numCrits = 0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			cd.Use(sim)
			// mage.UpdateMajorCooldowns()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.BonusFireCritRating += critPerStack * float64(newStacks-oldStacks)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell.SpellSchool != core.SpellSchoolFire {
				return
			}
			if spell.SameAction(IgniteActionID) {
				return
			}
			if !spellEffect.Landed() {
				return
			}
			if numCrits >= 3 {
				return
			}

			// TODO: This wont work properly with flamestrike
			aura.AddStack(sim)

			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
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
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			aura.Activate(sim)
			aura.Prioritize()
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return !aura.IsActive()
		},
	})
}

func (mage *Mage) registerIcyVeinsCD() {
	if !mage.Talents.IcyVeins {
		return
	}

	actionID := core.ActionID{SpellID: 12472}
	manaCost := mage.BaseMana * 0.03

	cooldown := 180.0
	if mage.Talents.IceFloes > 0 {
		cooldown *= 1 - []float64{0, .7, .14, .20}[mage.Talents.IceFloes]
	}

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

		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: manaCost,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Duration(cooldown) * time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			icyVeinsAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.IcyVeins,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Need to check for icy veins already active in case Cold Snap is used right after.
			if icyVeinsAura.IsActive() {
				return false
			}

			if character.CurrentMana() < manaCost {
				return false
			}

			return true
		},
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
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Don't use if there are no cooldowns to reset.
			return (mage.IcyVeins != nil && !mage.IcyVeins.IsReady(sim)) ||
				(mage.SummonWaterElemental != nil && !mage.SummonWaterElemental.IsReady(sim))
		},
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
			if isExecute == 20 {
				mage.PseudoStats.DamageDealtMultiplier *= multiplier
			}
		})
	})
}
