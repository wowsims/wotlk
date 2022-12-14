package druid

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type DruidForm uint8

const (
	Humanoid DruidForm = 1 << iota
	Bear
	Cat
	Moonkin
)

// Converts from 0.009327 to 0.0085
const AnimalSpiritRegenSuppression = 0.911337

func (form DruidForm) Matches(other DruidForm) bool {
	return (form & other) != 0
}

func (druid *Druid) GetForm() DruidForm {
	return druid.form
}

func (druid *Druid) InForm(form DruidForm) bool {
	return druid.form.Matches(form)
}

func (druid *Druid) ClearForm(sim *core.Simulation) {
	if druid.InForm(Cat) {
		druid.CatFormAura.Deactivate(sim)
	} else if druid.InForm(Bear) {
		druid.BearFormAura.Deactivate(sim)
	} else if druid.InForm(Moonkin) {
		panic("cant clear moonkin form")
	}
	druid.form = Humanoid
	druid.SetCurrentPowerBar(core.ManaBar)
}

// Bonus stats for both cat and bear.
func (druid *Druid) GetFormShiftStats() stats.Stats {
	s := stats.Stats{
		stats.AttackPower: float64(druid.Talents.PredatoryStrikes) * 0.5 * float64(core.CharacterLevel),
		stats.MeleeCrit:   float64(druid.Talents.SharpenedClaws) * 2 * core.CritRatingPerCritChance,
	}

	if weapon := druid.GetMHWeapon(); weapon != nil {
		dps := (weapon.WeaponDamageMax + weapon.WeaponDamageMin) / 2.0 / weapon.SwingSpeed
		weapAp := weapon.Stats[stats.AttackPower] + weapon.Enchant.Stats[stats.AttackPower]
		fap := math.Floor((dps - 54.8) * 14)

		s[stats.AttackPower] += fap
		s[stats.AttackPower] += (fap + weapAp) * ((0.2 / 3) * float64(druid.Talents.PredatoryStrikes))
	}

	return s
}

func (druid *Druid) registerCatFormSpell() {
	actionID := core.ActionID{SpellID: 768}
	baseCost := druid.BaseMana * 0.35

	srm := druid.getSavageRoarMultiplier()

	statBonus := druid.GetFormShiftStats().Add(stats.Stats{
		stats.AttackPower: float64(druid.Level) * 2,
		stats.MeleeCrit:   2 * float64(druid.Talents.MasterShapeshifter) * core.CritRatingPerCritChance,
	})

	agiApDep := druid.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 1)

	var hotwDep *stats.StatDependency
	if druid.Talents.HeartOfTheWild > 0 {
		hotwDep = druid.NewDynamicMultiplyStat(stats.AttackPower, 1.0+0.02*float64(druid.Talents.HeartOfTheWild))
	}

	catCritMult := druid.MeleeCritMultiplier(Cat)
	regCritMult := druid.MeleeCritMultiplier(Humanoid)

	druid.CatFormAura = druid.RegisterAura(core.Aura{
		Label:      "Cat Form",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(druid.StartingForm.Matches(Cat), core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !druid.Env.MeasuringStats && druid.form != Humanoid {
				druid.ClearForm(sim)
			}
			druid.form = Cat
			druid.SetCurrentPowerBar(core.EnergyBar)

			druid.AutoAttacks.MH = core.Weapon{
				BaseDamageMin:              43,
				BaseDamageMax:              66,
				SwingSpeed:                 1.0,
				NormalizedSwingSpeed:       1.0,
				SwingDuration:              time.Second,
				CritMultiplier:             catCritMult,
				MeleeAttackRatingPerDamage: core.MeleeAttackRatingPerDamage,
			}

			druid.PseudoStats.ThreatMultiplier *= 0.71
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.PseudoStats.BaseDodge += 0.02 * float64(druid.Talents.FeralSwiftness)
			druid.AddStatsDynamic(sim, statBonus)
			druid.EnableDynamicStatDep(sim, agiApDep)
			if hotwDep != nil {
				druid.EnableDynamicStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.ReplaceMHSwing = nil
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.manageCooldownsEnabled()
				druid.UpdateManaRegenRates()

				// These buffs stay up, but corresponding changes don't
				if druid.SavageRoarAura.IsActive() {
					druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= srm
				}
				if druid.BerserkAura.IsActive() {
					druid.PseudoStats.CostMultiplier /= 2.0
				}

				if druid.PredatoryInstinctsAura != nil {
					druid.PredatoryInstinctsAura.Activate(sim)
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid
			druid.AutoAttacks.MH = druid.WeaponFromMainHand(regCritMult)

			druid.PseudoStats.ThreatMultiplier /= 0.71
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.PseudoStats.BaseDodge -= 0.02 * float64(druid.Talents.FeralSwiftness)
			druid.AddStatsDynamic(sim, statBonus.Multiply(-1))
			druid.DisableDynamicStatDep(sim, agiApDep)
			if hotwDep != nil {
				druid.DisableDynamicStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.ReplaceMHSwing = nil
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.manageCooldownsEnabled()
				druid.UpdateManaRegenRates()

				druid.TigersFuryAura.Deactivate(sim)

				// These buffs stay up, but corresponding changes don't
				if druid.SavageRoarAura.IsActive() {
					druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= srm
				}
				if druid.BerserkAura.IsActive() {
					druid.PseudoStats.CostMultiplier *= 2.0
				}

				if druid.PredatoryInstinctsAura != nil {
					druid.PredatoryInstinctsAura.Deactivate(sim)
				}
			}
		},
	})

	energyMetrics := druid.NewEnergyMetrics(actionID)

	druid.CatForm = druid.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.2*float64(druid.Talents.KingOfTheJungle)) * (1 - 0.1*float64(druid.Talents.NaturalShapeshifter)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			maxShiftEnergy := float64(20 * druid.Talents.Furor)

			energyDelta := maxShiftEnergy - druid.CurrentEnergy()

			if energyDelta < 0 {
				druid.SpendEnergy(sim, -energyDelta, energyMetrics)
			}
			druid.CatFormAura.Activate(sim)
		},
	})
}

func (druid *Druid) calcArmorBonus() float64 {
	// Armor calculation: Dire Bear Form, Thick Hide, and Survival of the Fittest
	// scale multiplicatively with each other. But part of the Thick Hide
	// contribution was already calculated in ApplyTalents(), so we need to subtract
	// that part out from the overall scaling factor given to ScaleBaseArmor().
	sotfMulti := 1.0 + 0.33/3.0*float64(druid.Talents.SurvivalOfTheFittest)
	thickHideMulti := 1.0

	if druid.Talents.ThickHide > 0 {
		thickHideMulti += 0.04 + 0.03*float64(druid.Talents.ThickHide-1)
	}

	totalBearMulti := 4.7 * sotfMulti * thickHideMulti
	return druid.ScaleBaseArmor(totalBearMulti - thickHideMulti)
}

func (druid *Druid) registerBearFormSpell() {
	actionID := core.ActionID{SpellID: 9634}
	baseCost := druid.BaseMana * 0.35

	statBonus := druid.GetFormShiftStats().Add(stats.Stats{
		stats.Armor:       druid.calcArmorBonus(),
		stats.AttackPower: 3 * float64(core.CharacterLevel),
	})

	stamDep := druid.NewDynamicMultiplyStat(stats.Stamina, 1.25)

	var potpDep *stats.StatDependency
	if druid.Talents.ProtectorOfThePack > 0 {
		potpDep = druid.NewDynamicMultiplyStat(stats.AttackPower, 1.0+0.02*float64(druid.Talents.ProtectorOfThePack))
	}

	var hotwDep *stats.StatDependency
	if druid.Talents.HeartOfTheWild > 0 {
		hotwDep = druid.NewDynamicMultiplyStat(stats.Stamina, 1.0+0.02*float64(druid.Talents.HeartOfTheWild))
	}

	potpdtm := 1 - 0.04*float64(druid.Talents.ProtectorOfThePack)

	bearCritMult := druid.MeleeCritMultiplier(Bear)
	regCritMult := druid.MeleeCritMultiplier(Humanoid)

	druid.BearFormAura = druid.RegisterAura(core.Aura{
		Label:      "Bear Form",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(druid.StartingForm.Matches(Bear), core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !druid.Env.MeasuringStats && druid.form != Humanoid {
				druid.ClearForm(sim)
			}
			druid.form = Bear
			druid.SetCurrentPowerBar(core.RageBar)

			druid.AutoAttacks.MH = core.Weapon{
				BaseDamageMin:              109,
				BaseDamageMax:              165,
				SwingSpeed:                 2.5,
				NormalizedSwingSpeed:       2.5,
				SwingDuration:              time.Millisecond * 2500,
				CritMultiplier:             bearCritMult,
				MeleeAttackRatingPerDamage: core.MeleeAttackRatingPerDamage,
			}

			druid.PseudoStats.ThreatMultiplier *= 29. / 14.
			druid.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.02*float64(druid.Talents.MasterShapeshifter)
			druid.PseudoStats.DamageTakenMultiplier *= potpdtm
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.PseudoStats.BaseDodge += 0.02 * float64(druid.Talents.FeralSwiftness+druid.Talents.NaturalReaction)
			druid.AddStatsDynamic(sim, statBonus)
			druid.EnableDynamicStatDep(sim, stamDep)
			if potpDep != nil {
				druid.EnableDynamicStatDep(sim, potpDep)
			}
			if hotwDep != nil {
				druid.EnableDynamicStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.ReplaceMHSwing = func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
					return druid.TryMaul(sim, mhSwingSpell)
				}
				druid.AutoAttacks.EnableAutoSwing(sim)

				druid.manageCooldownsEnabled()
				druid.UpdateManaRegenRates()
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid
			druid.AutoAttacks.MH = druid.WeaponFromMainHand(regCritMult)

			druid.PseudoStats.ThreatMultiplier /= 29. / 14.
			druid.PseudoStats.DamageDealtMultiplier /= 1.0 + 0.02*float64(druid.Talents.MasterShapeshifter)
			druid.PseudoStats.DamageTakenMultiplier /= potpdtm
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.PseudoStats.BaseDodge -= 0.02 * float64(druid.Talents.FeralSwiftness+druid.Talents.NaturalReaction)
			druid.AddStatsDynamic(sim, statBonus.Multiply(-1))
			druid.DisableDynamicStatDep(sim, stamDep)
			if potpDep != nil {
				druid.DisableDynamicStatDep(sim, potpDep)
			}
			if hotwDep != nil {
				druid.DisableDynamicStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.ReplaceMHSwing = nil
				druid.AutoAttacks.EnableAutoSwing(sim)

				druid.manageCooldownsEnabled()
				druid.UpdateManaRegenRates()
				druid.EnrageAura.Deactivate(sim)
				druid.MaulQueueAura.Deactivate(sim)
			}
		},
	})

	rageMetrics := druid.NewRageMetrics(actionID)

	furorProcChance := []float64{0, 0.2, 0.4, 0.6, 0.8, 1}[druid.Talents.Furor]

	druid.BearForm = druid.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.2*float64(druid.Talents.KingOfTheJungle)) * (1 - 0.1*float64(druid.Talents.NaturalShapeshifter)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rageDelta := 0 - druid.CurrentRage()
			if sim.Proc(furorProcChance, "Furor") {
				rageDelta += 10
			}
			if rageDelta > 0 {
				druid.AddRage(sim, rageDelta, rageMetrics)
			} else if rageDelta < 0 {
				druid.SpendRage(sim, -rageDelta, rageMetrics)
			}
			druid.BearFormAura.Activate(sim)
		},
	})
}

func (druid *Druid) manageCooldownsEnabled() {
	// Disable cooldowns not usable in form and/or delay others
	if druid.StartingForm.Matches(Cat | Bear) {
		for _, mcd := range druid.disabledMCDs {
			mcd.Enable()
		}
		druid.disabledMCDs = nil

		if druid.InForm(Humanoid) {
			// Disable cooldown that incurs a gcd, so we dont get stuck out of form when we dont need to (Greater Drums)
			for _, mcd := range druid.GetMajorCooldowns() {
				if mcd.Spell.DefaultCast.GCD > 0 {
					mcd.Disable()
					druid.disabledMCDs = append(druid.disabledMCDs, mcd)
				}
			}
		}
	}
}

func (druid *Druid) applyMoonkinForm() {
	if !druid.InForm(Moonkin) || !druid.Talents.MoonkinForm {
		return
	}

	druid.MultiplyStat(stats.Intellect, 1+(0.02*float64(druid.Talents.Furor)))
	druid.PseudoStats.DamageDealtMultiplier *= 1 + (float64(druid.Talents.MasterShapeshifter) * 0.02)
	if druid.Talents.ImprovedMoonkinForm > 0 {
		druid.AddStatDependency(stats.Spirit, stats.SpellPower, 0.1*float64(druid.Talents.ImprovedMoonkinForm))
	}

	manaMetrics := druid.NewManaMetrics(core.ActionID{SpellID: 24858})
	druid.RegisterAura(core.Aura{
		Label:    "Moonkin Form",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() {
				if spell == druid.Moonfire || spell == druid.Starfire || spell == druid.Wrath {
					druid.AddMana(sim, 0.02*druid.MaxMana(), manaMetrics, false)
				}
			}
		},
	})
}
