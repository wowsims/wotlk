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

// Handles things that function for *both* cat/bear
func (druid *Druid) applyFeralShift(sim *core.Simulation, enter_form bool) {
	s := druid.GetFormShiftStats(enter_form)
	druid.AddStatsDynamic(sim, s)

	pos := core.TernaryFloat64(enter_form, 1.0, -1.0)
	druid.PseudoStats.BaseDodge += pos * 0.02 * float64(druid.Talents.FeralSwiftness) // Unaffected by Diminishing Returns
}

func (druid *Druid) GetFormShiftStats(enter_form bool) stats.Stats {
	s := stats.Stats{}

	pos := core.TernaryFloat64(enter_form, 1.0, -1.0)
	fap := 0.0
	weapAp := 0.0
	if weapon := druid.GetMHWeapon(); weapon != nil {
		dps := (weapon.WeaponDamageMax + weapon.WeaponDamageMin) / 2.0 / weapon.SwingSpeed
		weapAp = weapon.Stats[stats.AttackPower] + weapon.Enchant.Stats[stats.AttackPower]
		fap = math.Floor((dps - 54.8) * 14)
	}

	s[stats.AttackPower] = pos * fap

	if druid.Talents.PredatoryStrikes > 0 {
		s[stats.AttackPower] += pos * float64(druid.Talents.PredatoryStrikes) * 0.5 * float64(core.CharacterLevel)

		if fap > 0 {
			s[stats.AttackPower] += pos * (fap + weapAp) * ((0.2 / 3) * float64(druid.Talents.PredatoryStrikes))
		}
	}
	s[stats.MeleeCrit] += pos * float64(druid.Talents.SharpenedClaws) * 2 * core.CritRatingPerCritChance
	return s
}

type FormStatDep struct {
	Src    stats.Stat
	Dst    stats.Stat
	Amount float64
}

type FormRawStat struct {
	S      stats.Stat
	Amount float64
}

type FormBonuses struct {
	S    stats.Stats
	Deps []*FormStatDep
	Mul  []*FormRawStat
}

func (druid *Druid) GetCatFormBonuses(enable bool) FormBonuses {
	f := FormBonuses{}
	pos := core.TernaryFloat64(enable, 1.0, -1.0)

	f.S[stats.AttackPower] += pos * float64(druid.Level) * 2
	f.S[stats.MeleeCrit] += pos * 2 * float64(druid.Talents.MasterShapeshifter) * core.CritRatingPerCritChance

	// Ap dep
	f.Deps = append(f.Deps, &FormStatDep{
		Src:    stats.Agility,
		Dst:    stats.AttackPower,
		Amount: pos,
	})

	hotw := 1.0 + 0.02*float64(druid.Talents.HeartOfTheWild)
	f.Mul = append(f.Mul, &FormRawStat{
		S:      stats.AttackPower,
		Amount: core.Ternary(pos > 0, hotw, 1/hotw),
	})

	return f
}

func (druid *Druid) GetBearFormBonuses(enable bool) FormBonuses {
	f := FormBonuses{}
	pos := core.TernaryFloat64(enable, 1.0, -1.0)

	f.S[stats.Armor] = pos * druid.Equip.Stats()[stats.Armor] * 3.7
	f.S[stats.AttackPower] = pos * 3 * float64(core.CharacterLevel)

	// Stam dep
	f.Mul = append(f.Mul, &FormRawStat{
		S:      stats.Stamina,
		Amount: core.Ternary(pos > 0, 1.25, 1/1.25),
	})

	// Hotw
	hotw := 1.0 + 0.02*float64(druid.Talents.HeartOfTheWild)
	f.Mul = append(f.Mul, &FormRawStat{
		S:      stats.Stamina,
		Amount: core.Ternary(pos > 0, hotw, 1/hotw),
	})

	// Potp
	potp := 1.0 + 0.02*float64(druid.Talents.ProtectorOfThePack)
	f.Mul = append(f.Mul, &FormRawStat{
		S:      stats.AttackPower,
		Amount: core.Ternary(pos > 0, potp, 1/potp),
	})

	return f
}

func (druid *Druid) registerCatFormSpell() {
	actionID := core.ActionID{SpellID: 768}
	baseCost := druid.BaseMana * 0.35

	srm := druid.getSavageRoarMultiplier()

	bonuses := druid.GetCatFormBonuses(true)

	var dynamicDeps []*stats.StatDependency
	for _, d := range bonuses.Deps {
		dynamicDeps = append(dynamicDeps, druid.NewDynamicStatDependency(d.Src, d.Dst, d.Amount))
	}

	for _, stat := range bonuses.Mul {
		dynamicDeps = append(dynamicDeps, druid.NewDynamicMultiplyStat(stat.S, stat.Amount))
	}

	enabledStats := bonuses.S
	disabledStats := druid.GetCatFormBonuses(false).S

	catCritMult := druid.MeleeCritMultiplier(Cat)
	regCritMult := druid.MeleeCritMultiplier(Humanoid)

	druid.CatFormAura = druid.RegisterAura(core.Aura{
		Label:    "Cat Form",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if druid.form != Humanoid {
				druid.ClearForm(sim)
			}
			druid.form = Cat
			druid.AutoAttacks.MH = core.Weapon{
				BaseDamageMin:              43,
				BaseDamageMax:              66,
				SwingSpeed:                 1.0,
				NormalizedSwingSpeed:       1.0,
				SwingDuration:              time.Second,
				CritMultiplier:             catCritMult,
				MeleeAttackRatingPerDamage: core.MeleeAttackRatingPerDamage,
			}
			druid.AutoAttacks.ReplaceMHSwing = nil
			druid.AutoAttacks.EnableAutoSwing(sim)

			druid.SetCurrentPowerBar(core.EnergyBar)
			druid.manageCooldownsEnabled()
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()

			druid.applyFeralShift(sim, true)
			druid.AddStatsDynamic(sim, enabledStats)
			for _, d := range dynamicDeps {
				druid.EnableDynamicStatDep(sim, d)
			}

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
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid
			druid.AutoAttacks.MH = druid.WeaponFromMainHand(regCritMult)
			druid.AutoAttacks.ReplaceMHSwing = nil
			druid.AutoAttacks.EnableAutoSwing(sim)

			druid.manageCooldownsEnabled()
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()

			for _, d := range dynamicDeps {
				druid.DisableDynamicStatDep(sim, d)
			}
			druid.AddStatsDynamic(sim, disabledStats)
			druid.applyFeralShift(sim, false)

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

func (druid *Druid) registerBearFormSpell() {
	actionID := core.ActionID{SpellID: 9634}
	baseCost := druid.BaseMana * 0.35

	bonuses := druid.GetBearFormBonuses(true)

	var dynamicDeps []*stats.StatDependency
	for _, d := range bonuses.Deps {
		dynamicDeps = append(dynamicDeps, druid.NewDynamicStatDependency(d.Src, d.Dst, d.Amount))
	}

	for _, stat := range bonuses.Mul {
		dynamicDeps = append(dynamicDeps, druid.NewDynamicMultiplyStat(stat.S, stat.Amount))
	}

	enabledStats := bonuses.S
	disabledStats := druid.GetBearFormBonuses(false).S

	potpdtm := 0.04 * float64(druid.Talents.ProtectorOfThePack)

	bearCritMult := druid.MeleeCritMultiplier(Bear)
	regCritMult := druid.MeleeCritMultiplier(Humanoid)

	druid.BearFormAura = druid.RegisterAura(core.Aura{
		Label:    "Bear Form",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if druid.form != Humanoid {
				druid.ClearForm(sim)
			}
			druid.form = Bear

			druid.AutoAttacks.MH = core.Weapon{
				BaseDamageMin:              109,
				BaseDamageMax:              165,
				SwingSpeed:                 2.5,
				NormalizedSwingSpeed:       2.5,
				SwingDuration:              time.Millisecond * 2500,
				CritMultiplier:             bearCritMult,
				MeleeAttackRatingPerDamage: core.MeleeAttackRatingPerDamage,
			}

			druid.AutoAttacks.ReplaceMHSwing = func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
				return druid.TryMaul(sim, mhSwingSpell)
			}
			druid.AutoAttacks.EnableAutoSwing(sim)

			druid.SetCurrentPowerBar(core.RageBar)
			druid.applyFeralShift(sim, true)
			druid.AddStatsDynamic(sim, enabledStats)

			for _, d := range dynamicDeps {
				druid.EnableDynamicStatDep(sim, d)
			}

			druid.PseudoStats.ThreatMultiplier *= 1.3
			druid.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.02*float64(druid.Talents.MasterShapeshifter)
			druid.PseudoStats.DamageTakenMultiplier *= 1.0 - potpdtm

			druid.manageCooldownsEnabled()
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid
			druid.AutoAttacks.MH = druid.WeaponFromMainHand(regCritMult)
			druid.AutoAttacks.ReplaceMHSwing = nil
			druid.AutoAttacks.EnableAutoSwing(sim)

			for _, d := range dynamicDeps {
				druid.DisableDynamicStatDep(sim, d)
			}
			druid.AddStatsDynamic(sim, disabledStats)
			druid.applyFeralShift(sim, false)

			druid.PseudoStats.ThreatMultiplier /= 1.3
			druid.PseudoStats.DamageDealtMultiplier /= 1.0 + 0.02*float64(druid.Talents.MasterShapeshifter)
			druid.PseudoStats.DamageTakenMultiplier /= 1.0 - potpdtm

			druid.manageCooldownsEnabled()
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.UpdateManaRegenRates()
			druid.EnrageAura.Deactivate(sim)
			druid.MaulQueueAura.Deactivate(sim)
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
