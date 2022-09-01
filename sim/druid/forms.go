package druid

import (
	"math"

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
	druid.CatFormAura.Deactivate(sim)
	druid.BearFormAura.Deactivate(sim)
	druid.form = Humanoid
}

func (druid *Druid) PowerShiftCat(sim *core.Simulation) bool {

	if !druid.GCD.IsReady(sim) {
		panic("Trying to powershift during gcd")
	}

	druid.CatFormAura.Deactivate(sim)
	druid.TryUseCooldowns(sim)

	if druid.GCD.IsReady(sim) {
		return druid.CatForm.Cast(sim, nil)
	}

	return true
}

// Handles things that function for *both* cat/bear
func (druid *Druid) applyFeralShift(sim *core.Simulation, enter_form bool) {
	weap := druid.GetMHWeapon()
	pos := core.TernaryFloat64(enter_form, 1.0, -1.0)
	fap := 0.0
	if weap != nil {
		dps := (((weap.WeaponDamageMax - weap.WeaponDamageMin) / 2.0) + weap.WeaponDamageMin) / weap.SwingSpeed
		fap = math.Floor((dps - 54.8) * 14)
	}
	druid.AddStatDynamic(sim, stats.AttackPower, pos*fap)

	if druid.Talents.PredatoryStrikes > 0 {
		druid.AddStatDynamic(sim, stats.AttackPower, pos*float64(druid.Talents.PredatoryStrikes)*0.5*float64(core.CharacterLevel))

		if fap > 0 {
			druid.AddStatDynamic(sim, stats.AttackPower, pos*fap*((0.2/3)*float64(druid.Talents.PredatoryStrikes)))
		}
	}
	druid.AddStatDynamic(sim, stats.MeleeCrit, pos*float64(druid.Talents.SharpenedClaws)*2*core.CritRatingPerCritChance)
	druid.AddStatDynamic(sim, stats.Dodge, pos*core.DodgeRatingPerDodgeChance*2*float64(druid.Talents.FeralSwiftness))
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
				panic("must leave form first")
			}
			druid.form = Cat
			druid.AutoAttacks.EnableAutoSwing(sim)
			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()

			druid.applyFeralShift(sim, true)
			druid.AddStatDynamic(sim, stats.AttackPower, float64(druid.Level)*2)
			druid.EnableDynamicStatDep(sim, catHotw)
			druid.EnableDynamicStatDep(sim, apDep)
			druid.AddStatDynamic(sim, stats.MeleeCrit, 2*float64(druid.Talents.MasterShapeshifter)*core.CritRatingPerCritChance)

			// These buffs stay up, but corresponding changes don't
			if druid.SavageRoarAura.IsActive() {
				druid.PseudoStats.PhysicalDamageDealtMultiplier *= srm
			}
			if druid.BerserkAura.IsActive() {
				druid.PseudoStats.CostMultiplier /= 2.0
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid
			druid.AutoAttacks.CancelAutoSwing(sim)
			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()

			druid.applyFeralShift(sim, false)
			druid.AddStatDynamic(sim, stats.AttackPower, -(float64(druid.Level) * 2))
			druid.DisableDynamicStatDep(sim, catHotw)
			druid.DisableDynamicStatDep(sim, apDep)
			druid.AddStatDynamic(sim, stats.MeleeCrit, -2*float64(druid.Talents.MasterShapeshifter)*core.CritRatingPerCritChance)

			druid.TigersFuryAura.Deactivate(sim)

			// These buffs stay up, but corresponding changes don't
			if druid.SavageRoarAura.IsActive() {
				druid.PseudoStats.PhysicalDamageDealtMultiplier /= srm
			}
			if druid.BerserkAura.IsActive() {
				druid.PseudoStats.CostMultiplier *= 2.0
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
	furorProcChance := 0.2 * float64(druid.Talents.Furor)

	previousRage := 0.0
	finalRage := 0.0

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
				panic("must leave form first")
			}
			druid.form = Bear

			druid.AddStatDynamic(sim, stats.AttackPower, 3*float64(core.CharacterLevel))
			druid.EnableDynamicStatDep(sim, stamdep)
			druid.EnableDynamicStatDep(sim, bearHotw)
			druid.EnableDynamicStatDep(sim, potpap)
			druid.PseudoStats.ThreatMultiplier *= 1.3
			druid.PseudoStats.DamageDealtMultiplier += 0.02 * float64(druid.Talents.MasterShapeshifter)
			druid.PseudoStats.DamageTakenMultiplier += -1.0 * potpdtm

			druid.applyFeralShift(sim, true)
			druid.AutoAttacks.EnableAutoSwing(sim)
			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()

			rb := druid.GetAura("RageBar")
			if rb != nil {
				rb.Activate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			previousRage = druid.CurrentRage()
			druid.form = Humanoid

			druid.AddStatDynamic(sim, stats.AttackPower, -3*float64(core.CharacterLevel))
			druid.DisableDynamicStatDep(sim, bearHotw)
			druid.DisableDynamicStatDep(sim, stamdep)
			druid.DisableDynamicStatDep(sim, potpap)
			druid.PseudoStats.ThreatMultiplier /= 1.3
			druid.PseudoStats.DamageDealtMultiplier -= 0.02 * float64(druid.Talents.MasterShapeshifter)
			druid.PseudoStats.DamageTakenMultiplier += 1.0 * potpdtm

			druid.applyFeralShift(sim, false)
			druid.AutoAttacks.CancelAutoSwing(sim)
			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()
			druid.EnrageAura.Deactivate(sim)
			rb := druid.GetAura("RageBar")
			if rb != nil {
				rb.Deactivate(sim)
			}
		},
	})

	rageMetrics := druid.NewRageMetrics(actionID)

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
			rageDelta := finalRage - previousRage
			if furorProcChance == 1 || (furorProcChance > 0 && sim.RandomFloat("Furor") < furorProcChance) {
				finalRage += 10.0
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

		druid.EnableAllCooldowns(druid.disabledMCDs)
		druid.disabledMCDs = nil

		if druid.InForm(Humanoid) {
			// Disable cooldown that incurs a gcd, so we dont get stuck out of form when we dont need to (Greater Drums)
			for _, cd := range druid.GetMajorCooldowns() {
				if cd.Spell.DefaultCast.GCD > 0 {
					druid.DisableMajorCooldown(cd.Spell.ActionID)
					druid.disabledMCDs = append(druid.disabledMCDs, cd)
				}
			}
		}
	}
}
