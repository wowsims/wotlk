package mage

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (mage *Mage) ApplyTalents() {
	mage.applyArcaneConcentration()
	mage.applyIgnite()
	mage.applyMasterOfElements()
	mage.applyWintersChill()
	mage.applyMoltenFury()
	mage.registerArcanePowerCD()
	mage.registerPresenceOfMindCD()
	mage.registerCombustionCD()
	mage.registerIcyVeinsCD()
	mage.registerColdSnapCD()
	mage.registerSummonWaterElementalCD()

	if mage.Talents.ArcaneMeditation > 0 {
		mage.PseudoStats.SpiritRegenRateCasting += float64(mage.Talents.ArcaneMeditation) * 0.1
	}

	if mage.Talents.ArcaneMind > 0 {
		mage.Character.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.Intellect,
			Modifier: func(intellect float64, _ float64) float64 {
				return intellect * (1.0 + 0.03*float64(mage.Talents.ArcaneMind))
			},
		})
	}

	if mage.Talents.MindMastery > 0 {
		mage.Character.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.SpellPower,
			Modifier: func(intellect float64, spellPower float64) float64 {
				return spellPower + intellect*0.05*float64(mage.Talents.MindMastery)
			},
		})
	}

	if mage.Talents.ArcaneInstability > 0 {
		mage.AddStat(stats.SpellCrit, float64(mage.Talents.ArcaneInstability)*1*core.SpellCritRatingPerCritChance)
		mage.spellDamageMultiplier += float64(mage.Talents.ArcaneInstability) * 0.01
	}

	if mage.Talents.PlayingWithFire > 0 {
		mage.spellDamageMultiplier += float64(mage.Talents.PlayingWithFire) * 0.01
	}

	magicAbsorptionBonus := 2 * float64(mage.Talents.MagicAbsorption)
	mage.AddStat(stats.ArcaneResistance, magicAbsorptionBonus)
	mage.AddStat(stats.FireResistance, magicAbsorptionBonus)
	mage.AddStat(stats.FrostResistance, magicAbsorptionBonus)
	mage.AddStat(stats.NatureResistance, magicAbsorptionBonus)
	mage.AddStat(stats.ShadowResistance, magicAbsorptionBonus)
}

func (mage *Mage) applyArcaneConcentration() {
	if mage.Talents.ArcaneConcentration == 0 {
		return
	}

	procChance := 0.02 * float64(mage.Talents.ArcaneConcentration)
	bonusCrit := float64(mage.Talents.ArcanePotency) * 10 * core.SpellCritRatingPerCritChance

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

func (mage *Mage) registerPresenceOfMindCD() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	cooldown := time.Minute * 3
	if ItemSetAldorRegalia.CharacterHasSetBonus(&mage.Character, 4) {
		cooldown -= time.Second * 24
	}

	actionID := core.ActionID{SpellID: 12043}

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
			var spell *core.Spell
			if mage.Talents.Pyroblast {
				spell = mage.Pyroblast
			} else if mage.RotationType == proto.Mage_Rotation_Fire {
				spell = mage.Fireball
			} else if mage.RotationType == proto.Mage_Rotation_Frost {
				spell = mage.Frostbolt
			} else {
				numStacks := mage.ArcaneBlastAura.GetStacks()
				spell = mage.ArcaneBlast[numStacks]
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
			} else if mage.RotationType == proto.Mage_Rotation_Fire {
				manaCost = mage.Fireball.DefaultCast.Cost
			} else if mage.RotationType == proto.Mage_Rotation_Frost {
				manaCost = mage.Frostbolt.DefaultCast.Cost
			} else {
				numStacks := mage.ArcaneBlastAura.GetStacks()
				manaCost = mage.ArcaneBlast[numStacks].DefaultCast.Cost
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
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.PseudoStats.DamageDealtMultiplier *= 1.3
			mage.PseudoStats.CostMultiplier *= 1.3
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.PseudoStats.DamageDealtMultiplier /= 1.3
			mage.PseudoStats.CostMultiplier /= 1.3
		},
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 3,
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
	if mage.Talents.MasterOfElements == 0 {
		return
	}

	refundCoeff := 0.1 * float64(mage.Talents.MasterOfElements)
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
				mage.AddMana(sim, spell.BaseCost*refundCoeff, manaMetrics, false)
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
	const critPerStack = 10 * core.SpellCritRatingPerCritChance

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
			mage.UpdateMajorCooldowns()
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
	manaCost := mage.BaseMana() * 0.03

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
				Duration: time.Minute * 3,
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

	cooldown := time.Duration(float64(time.Minute*8) * (1.0 - float64(mage.Talents.IceFloes)*0.1))
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

	multiplier := 1.0 + 0.1*float64(mage.Talents.MoltenFury)

	mage.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation) {
			mage.PseudoStats.DamageDealtMultiplier *= multiplier
		})
	})
}
