package mage

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/stats"
)

func (mage *Mage) ApplyTalents() {
	mage.applyIgnite()
	mage.applyArcaneConcentration()
	mage.applyMasterOfElements()
	mage.applyWintersChill()

	mage.registerArcanePowerCD()
	mage.registerPresenceOfMindCD()
	mage.registerCombustionCD()
	mage.registerColdSnapCD()

	mage.PseudoStats.SpiritRegenRateCasting += 0.05 * float64(mage.Talents.ArcaneMeditation)

	if mage.Talents.ElementalPrecision > 0 {
		bonusHit := 2 * float64(mage.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance
		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolFire) || spell.SpellSchool.Matches(core.SpellSchoolFrost) {
				spell.BonusHitRating += bonusHit
			}
		})
	}

	if mage.Talents.ArcaneFocus > 0 {
		bonusHit := 2 * float64(mage.Talents.ArcaneFocus) * core.SpellHitRatingPerHitChance
		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolArcane) {
				spell.BonusHitRating += bonusHit
			}
		})
	}

	if mage.Talents.PiercingIce > 0 {
		mage.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1 + 0.02*float64(mage.Talents.PiercingIce)
	}

	if mage.Talents.FirePower > 0 {
		mage.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1 + 0.02*float64(mage.Talents.FirePower)
	}

	if mage.Talents.CriticalMass > 0 {
		bonusCrit := 2 * float64(mage.Talents.CriticalMass) * core.SpellCritRatingPerCritChance
		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolFire) {
				spell.BonusCritRating += bonusCrit
			}
		})
	}

	if mage.Talents.ArcaneMind > 0 {
		mage.MultiplyStat(stats.Mana, 1.0+0.02*float64(mage.Talents.ArcaneMind))
	}

	mage.AddStat(stats.SpellCrit, float64(mage.Talents.ArcaneInstability)*1*core.SpellCritRatingPerCritChance)
	mage.PseudoStats.DamageDealtMultiplier *= 1 + .01*float64(mage.Talents.ArcaneInstability)
	mage.PseudoStats.CostMultiplier *= 1 - .05*float64(mage.Talents.FrostChanneling)

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

	var proccedAt time.Duration
	var proccedSpell *core.Spell

	mage.ClearcastingAura = mage.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 12577},
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
			if spell.DefaultCast.Cost == 0 {
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

			procChance := 0.02 * float64(mage.Talents.ArcaneConcentration)

			// TODO: Classic verify arcane missle proc chance
			// Arcane Missile ticks can proc CC, just at a low rate of about 1.5% with 5/5 Arcane Concentration
			if spell == mage.ArcaneMissilesTickSpell {
				procChance *= 0.15
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

// TODO: Classic allow more dynamic choice in PoM with APL
func (mage *Mage) registerPresenceOfMindCD() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	cooldown := 180.0

	actionID := core.ActionID{SpellID: 12043}

	var spellToUse *core.Spell
	mage.Env.RegisterPostFinalizeEffect(func() {
		if mage.ArcaneBlast != nil {
			spellToUse = mage.ArcaneBlast
		} else if mage.Pyroblast != nil {
			spellToUse = mage.Pyroblast
		} else if mage.PrimaryTalentTree == 1 {
			spellToUse = mage.Fireball
		} else {
			spellToUse = mage.Frostbolt
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
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += 0.3
				spell.CostMultiplier += 0.3
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive -= 0.3
				spell.CostMultiplier -= 0.3
			}
		},
	})
	core.RegisterPercentDamageModifierEffect(mage.ArcanePowerAura, 1.3)

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 180,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.ArcanePowerAura.Activate(sim)
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				return
			}
			if spell.CurCast.Cost == 0 {
				return
			}
			if result.DidCrit() {
				mage.AddMana(sim, spell.DefaultCast.Cost*refundCoeff, manaMetrics)
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

	fireCombCritMult := mage.SpellCritMultiplier(1, 0.1) / mage.SpellCritMultiplier(1, 0)

	var fireSpells []*core.Spell
	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellSchool.Matches(core.SpellSchoolFire) {
			fireSpells = append(fireSpells, spell)
		}
	})

	numCrits := 0
	const critPerStack = 10 * core.SpellCritRatingPerCritChance

	mage.CombustionAura = mage.RegisterAura(core.Aura{
		Label:     "Combustion",
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			numCrits = 0
			for _, spell := range fireSpells {
				spell.CritMultiplier *= fireCombCritMult
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			cd.Use(sim)
			mage.UpdateMajorCooldowns()
			for _, spell := range fireSpells {
				spell.CritMultiplier /= fireCombCritMult
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

func (mage *Mage) registerColdSnapCD() {
	if !mage.Talents.ColdSnap {
		return
	}

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 12472},
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Duration(time.Minute * 10),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if mage.IcyVeins != nil {
				mage.IcyVeins.CD.Reset()
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if mage.IcyVeins != nil && mage.IcyVeins.IsReady(sim) {
				return false
			}

			return true
		},
	})
}

// func (mage *Mage) hasChillEffect(spell *core.Spell) bool {
// 	return spell == mage.Frostbolt || spell == mage.FrostfireBolt || (spell == mage.Blizzard && mage.Talents.ImprovedBlizzard > 0)
// }

func (mage *Mage) applyWintersChill() {
	if mage.Talents.WintersChill == 0 {
		return
	}

	procChance := float64(mage.Talents.WintersChill) * 0.2

	wcAuras := mage.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		return core.WintersChillAura(target, 0)
	})
	mage.Env.RegisterPreFinalizeEffect(func() {
		for _, spell := range mage.GetSpellsMatchingSchool(core.SpellSchoolFrost) {
			spell.RelatedAuras = append(spell.RelatedAuras, wcAuras)
		}
	})

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
				aura := wcAuras.Get(result.Target)
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			}
		},
	})
}
