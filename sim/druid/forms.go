package druid

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type DruidForm uint8

const (
	Humanoid = 1 << iota
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

func (druid *Druid) PowerShiftCat(sim *core.Simulation) bool {

	if !druid.GCD.IsReady(sim) {
		panic("Trying to powershift during gcd")
	}

	druid.CatFormAura.Deactivate(sim)
	druid.TryUseCooldowns(sim)

	return druid.CatForm.Cast(sim, nil)
}

// Handles things that function for *both* cat/bear
func (druid *Druid) applyFeralShift(sim *core.Simulation, enter_form bool) {
	pos := core.TernaryFloat64(enter_form, 1.0, -1.0)
	fap := 0.0
	weapAp := 0.0
	if weapon := druid.GetMHWeapon(); weapon != nil {
		dps := (weapon.WeaponDamageMax + weapon.WeaponDamageMin) / 2.0 / weapon.SwingSpeed
		weapAp = weapon.Stats[stats.AttackPower]
		fap = math.Floor((dps - 54.8) * 14)
	}
	druid.AddStatDynamic(sim, stats.AttackPower, pos*fap)

	if druid.Talents.PredatoryStrikes > 0 {
		druid.AddStatDynamic(sim, stats.AttackPower, pos*float64(druid.Talents.PredatoryStrikes)*0.5*float64(core.CharacterLevel))

		if fap > 0 {
			druid.AddStatDynamic(sim, stats.AttackPower, pos*(fap+weapAp)*((0.2/3)*float64(druid.Talents.PredatoryStrikes)))
		}
	}
	druid.AddStatDynamic(sim, stats.MeleeCrit, pos*float64(druid.Talents.SharpenedClaws)*2*core.CritRatingPerCritChance)
	druid.PseudoStats.BaseDodge += pos * 0.02 * float64(druid.Talents.FeralSwiftness) // Unaffected by Diminishing Returns
}

func (druid *Druid) registerCatFormSpell() {
	actionID := core.ActionID{SpellID: 768}
	baseCost := druid.BaseMana * 0.35

	srm := druid.getSavageRoarMultiplier()

	apDep := druid.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 1)
	catHotw := druid.NewDynamicMultiplyStat(stats.AttackPower, 1.0+0.02*float64(druid.Talents.HeartOfTheWild))

	druid.CatFormAura = druid.GetOrRegisterAura(core.Aura{
		Label:    "Cat Form",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if druid.form != Humanoid {
				druid.ClearForm(sim)
			}
			druid.form = Cat
			druid.SetCurrentPowerBar(core.EnergyBar)
			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()

			druid.applyFeralShift(sim, true)
			druid.AddStatDynamic(sim, stats.AttackPower, float64(druid.Level)*2)
			druid.EnableDynamicStatDep(sim, apDep)
			druid.EnableDynamicStatDep(sim, catHotw)
			druid.AddStatDynamic(sim, stats.MeleeCrit, 2*float64(druid.Talents.MasterShapeshifter)*core.CritRatingPerCritChance)

			// These buffs stay up, but corresponding changes don't
			if druid.SavageRoarAura.IsActive() {
				druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= srm
			}
			if druid.BerserkAura.IsActive() {
				druid.PseudoStats.CostMultiplier /= 2.0
			}

			druid.AutoAttacks.MH = core.Weapon{
				BaseDamageMin:              43,
				BaseDamageMax:              66,
				SwingSpeed:                 1.0,
				NormalizedSwingSpeed:       1.0,
				SwingDuration:              time.Second,
				CritMultiplier:             druid.MeleeCritMultiplier(),
				MeleeAttackRatingPerDamage: core.MeleeAttackRatingPerDamage,
			}
			druid.AutoAttacks.ReplaceMHSwing = nil
			druid.AutoAttacks.EnableAutoSwing(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid
			druid.AutoAttacks.CancelAutoSwing(sim)
			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()

			druid.AddStatDynamic(sim, stats.MeleeCrit, -2*float64(druid.Talents.MasterShapeshifter)*core.CritRatingPerCritChance)
			druid.DisableDynamicStatDep(sim, catHotw)
			druid.DisableDynamicStatDep(sim, apDep)
			druid.AddStatDynamic(sim, stats.AttackPower, -(float64(druid.Level) * 2))
			druid.applyFeralShift(sim, false)

			druid.TigersFuryAura.Deactivate(sim)

			// These buffs stay up, but corresponding changes don't
			if druid.SavageRoarAura.IsActive() {
				druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= srm
			}
			if druid.BerserkAura.IsActive() {
				druid.PseudoStats.CostMultiplier *= 2.0
			}

			druid.AutoAttacks.MH = druid.WeaponFromMainHand(0)
			druid.AutoAttacks.ReplaceMHSwing = nil
			druid.AutoAttacks.EnableAutoSwing(sim)
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

func (druid *Druid) PowerShiftBear(sim *core.Simulation) {

	if !druid.GCD.IsReady(sim) {
		panic("Trying to powershift during gcd")
	}

	druid.BearFormAura.Deactivate(sim)
	druid.TryUseCooldowns(sim)

	if druid.GCD.IsReady(sim) {
		druid.BearForm.Cast(sim, nil)
	}
}

func (druid *Druid) registerBearFormSpell() {
	actionID := core.ActionID{SpellID: 9634}
	baseCost := druid.BaseMana * 0.35

	stamdep := druid.NewDynamicMultiplyStat(stats.Stamina, 1.25)
	bearHotw := druid.NewDynamicMultiplyStat(stats.Stamina, 1.0+0.02*float64(druid.Talents.HeartOfTheWild))

	potpdtm := 0.04 * float64(druid.Talents.ProtectorOfThePack)
	potpap := druid.NewDynamicMultiplyStat(stats.AttackPower, 1.0+0.02*float64(druid.Talents.ProtectorOfThePack))

	druid.BearFormAura = druid.GetOrRegisterAura(core.Aura{
		Label:    "Bear Form",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if druid.form != Humanoid {
				druid.ClearForm(sim)
			}
			druid.form = Bear

			druid.SetCurrentPowerBar(core.RageBar)
			druid.applyFeralShift(sim, true)
			druid.AddStatDynamic(sim, stats.AttackPower, 3*float64(core.CharacterLevel))
			druid.EnableDynamicStatDep(sim, stamdep)
			druid.EnableDynamicStatDep(sim, bearHotw)
			druid.EnableDynamicStatDep(sim, potpap)

			druid.PseudoStats.ThreatMultiplier *= 1.3
			druid.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.02*float64(druid.Talents.MasterShapeshifter)
			druid.PseudoStats.DamageTakenMultiplier *= 1.0 - potpdtm

			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()

			druid.AutoAttacks.MH = core.Weapon{
				BaseDamageMin:              109,
				BaseDamageMax:              165,
				SwingSpeed:                 2.5,
				NormalizedSwingSpeed:       2.5,
				SwingDuration:              time.Millisecond * 2500,
				CritMultiplier:             druid.MeleeCritMultiplier(),
				MeleeAttackRatingPerDamage: core.MeleeAttackRatingPerDamage,
			}

			druid.AutoAttacks.ReplaceMHSwing = func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
				return druid.TryMaul(sim, mhSwingSpell)
			}
			druid.AutoAttacks.EnableAutoSwing(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid

			druid.DisableDynamicStatDep(sim, potpap)
			druid.DisableDynamicStatDep(sim, bearHotw)
			druid.DisableDynamicStatDep(sim, stamdep)
			druid.AddStatDynamic(sim, stats.AttackPower, -3*float64(core.CharacterLevel))
			druid.applyFeralShift(sim, false)

			druid.PseudoStats.ThreatMultiplier /= 1.3
			druid.PseudoStats.DamageDealtMultiplier /= 1.0 + 0.02*float64(druid.Talents.MasterShapeshifter)
			druid.PseudoStats.DamageTakenMultiplier /= 1.0 - potpdtm

			druid.AutoAttacks.CancelAutoSwing(sim)
			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()
			druid.EnrageAura.Deactivate(sim)
			druid.MaulQueueAura.Deactivate(sim)

			druid.AutoAttacks.MH = druid.WeaponFromMainHand(0)
			druid.AutoAttacks.ReplaceMHSwing = nil
			druid.AutoAttacks.EnableAutoSwing(sim)
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

func (druid *Druid) manageCooldownsEnabled(sim *core.Simulation) {
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
