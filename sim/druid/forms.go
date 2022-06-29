package druid

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/stats"
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

func (druid *Druid) registerCatFormSpell() {
	actionID := core.ActionID{SpellID: 768}
	baseCost := druid.BaseMana() * 0.35

	furorProcChance := 0.2 * float64(druid.Talents.Furor)

	finalEnergy := 0.0
	if druid.Equip[items.ItemSlotHead].ID == 8345 { // Wolfshead Helm
		finalEnergy += 20.0
	}

	druid.CatFormAura = druid.GetOrRegisterAura(core.Aura{
		Label:    "Cat Form",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Cat
			druid.AutoAttacks.EnableAutoSwing(sim)
			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid
			druid.AutoAttacks.CancelAutoSwing(sim)
			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()
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
				Cost: baseCost * (1 - 0.1*float64(druid.Talents.NaturalShapeshifter)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			energyDelta := finalEnergy - druid.CurrentEnergy()
			if furorProcChance == 1 || (furorProcChance > 0 && sim.RandomFloat("Furor") < furorProcChance) {
				energyDelta += 40.0
			}

			if energyDelta > 0 {
				druid.AddEnergy(sim, energyDelta, energyMetrics)
			} else if energyDelta < 0 {
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
	baseCost := druid.BaseMana() * 0.35
	furorProcChance := 0.2 * float64(druid.Talents.Furor)

	previousRage := 0.0
	finalRage := 0.0
	if druid.Equip[items.ItemSlotHead].ID == 8345 { // Wolfshead Helm
		finalRage += 5.0
	}

	druid.BearFormAura = druid.GetOrRegisterAura(core.Aura{
		Label:    "Bear Form",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Bear
			druid.AutoAttacks.EnableAutoSwing(sim)
			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			previousRage = druid.CurrentRage()
			druid.form = Humanoid
			druid.AutoAttacks.CancelAutoSwing(sim)
			druid.manageCooldownsEnabled(sim)
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()
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
				Cost: baseCost * (1 - 0.1*float64(druid.Talents.NaturalShapeshifter)),
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

// A bit arbitrary
const cooldownDelayThresHold = time.Second * 10

func (druid *Druid) manageCooldownsEnabled(sim *core.Simulation) {

	// Disable cooldowns not usable in form and/or delay others
	if druid.StartingForm.Matches(Cat | Bear) {

		druid.EnableAllCooldowns(druid.disabledMCDs)
		druid.disabledMCDs = nil

		if druid.InForm(Cat | Bear) {
			// Check if any dps cooldown that requires shifting is ready soon
			// disable all cooldowns if that is the case
			nonUsableDpsMCDReadySoon := false
			for _, cd := range druid.GetMajorCooldowns() {
				if cd.TimeToReady(sim) < cooldownDelayThresHold && cd.IsEnabled() && !cd.Type.Matches(core.CooldownTypeUsableShapeShifted) && cd.Type.Matches(core.CooldownTypeDPS) {
					nonUsableDpsMCDReadySoon = true
					break
				}
			}
			for _, cd := range druid.GetMajorCooldowns() {
				if cd.IsEnabled() && (nonUsableDpsMCDReadySoon || !cd.Type.Matches(core.CooldownTypeUsableShapeShifted)) {
					druid.DisableMajorCooldown(cd.Spell.ActionID)
					druid.disabledMCDs = append(druid.disabledMCDs, cd)
				}
			}
		} else {
			// Disable cooldown that can be used in form, but incurs a gcd, so we dont get stuck out of form when we dont need to (Greater Drums)
			for _, cd := range druid.GetMajorCooldowns() {
				if cd.Type.Matches(core.CooldownTypeUsableShapeShifted) && cd.Spell.DefaultCast.GCD > 0 {
					druid.DisableMajorCooldown(cd.Spell.ActionID)
					druid.disabledMCDs = append(druid.disabledMCDs, cd)
				}
			}
		}
	}
}
