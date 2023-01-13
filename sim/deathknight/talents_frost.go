package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) ApplyFrostTalents() {
	// Improved Icy Touch
	// Implemented outside

	// Toughness
	if dk.Talents.Toughness > 0 {
		dk.AddStat(stats.Armor, dk.Equip.Stats()[stats.Armor]*0.02*float64(dk.Talents.Toughness))
	}

	// Icy Reach
	// Pointless to Implement

	// Black Ice
	dk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1.0 + 0.02*float64(dk.Talents.BlackIce)
	dk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.0 + 0.02*float64(dk.Talents.BlackIce)
	dk.modifyShadowDamageModifier(0.02 * float64(dk.Talents.BlackIce))

	// Nerves Of Cold Steel
	if dk.nervesOfColdSteelActive() {
		dk.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*float64(dk.Talents.NervesOfColdSteel))
		dk.AutoAttacks.OHConfig.DamageMultiplier *= dk.nervesOfColdSteelBonus()
	}

	// Icy Talons
	dk.applyIcyTalons()

	// Lichborne
	// Pointless to Implement

	// Annihilation

	// Killing Machine
	dk.applyKillingMachine()

	// Chill of the Grave
	// Implemented outside

	// Endless Winter
	if dk.Talents.EndlessWinter > 0 {
		strengthCoeff := 0.02 * float64(dk.Talents.EndlessWinter)
		dk.MultiplyStat(stats.Strength, 1.0+strengthCoeff)
	}

	// Frigid Dreadplate
	if dk.Talents.FrigidDreadplate > 0 {
		dk.PseudoStats.BonusMeleeHitRatingTaken -= core.MeleeHitRatingPerHitChance * float64(dk.Talents.FrigidDreadplate)
	}

	// Glacier rot
	dk.applyGlacierRot()

	// Deathchill
	dk.applyDeathchill()

	// Merciless Combat
	dk.applyMercilessCombat()

	dk.applyThreatOfThassarian()

	// Rime
	dk.applyRime()

	// Tundra Stalker
	dk.AddStat(stats.Expertise, 1.0*float64(dk.Talents.TundraStalker)*core.ExpertisePerQuarterPercentReduction)
	if dk.Talents.TundraStalker > 0 {
		dk.applyTundaStalker()
	}
}

func (dk *Deathknight) nervesOfColdSteelActive() bool {
	return dk.HasMHWeapon() && dk.HasOHWeapon() && dk.Equip[proto.ItemSlot_ItemSlotMainHand].HandType == proto.HandType_HandTypeMainHand || dk.Equip[proto.ItemSlot_ItemSlotMainHand].HandType == proto.HandType_HandTypeOneHand
}

func (dk *Deathknight) nervesOfColdSteelBonus() float64 {
	return []float64{1.0, 1.08, 1.16, 1.25}[dk.Talents.NervesOfColdSteel]
}

func (dk *Deathknight) applyGlacierRot() {
	dk.bonusCoeffs.glacierRotBonusCoeff = []float64{1.0, 1.07, 1.13, 1.20}[dk.Talents.GlacierRot]
}

func (dk *Deathknight) glacielRotBonus(target *core.Unit) float64 {
	return core.TernaryFloat64(dk.DiseasesAreActive(target), dk.bonusCoeffs.glacierRotBonusCoeff, 1.0)
}

func (dk *Deathknight) applyMercilessCombat() {
	dk.bonusCoeffs.mercilessCombatBonusCoeff = 1.0 + 0.06*float64(dk.Talents.MercilessCombat)
}

func (dk *Deathknight) mercilessCombatBonus(sim *core.Simulation) float64 {
	return core.TernaryFloat64(sim.IsExecutePhase35() && dk.Talents.MercilessCombat > 0, dk.bonusCoeffs.mercilessCombatBonusCoeff, 1.0)
}

func (dk *Deathknight) applyTundaStalker() {
	bonus := 1.0 + 0.03*float64(dk.Talents.TundraStalker)
	dk.RoRTSBonus = func(target *core.Unit) float64 {
		return core.TernaryFloat64(dk.FrostFeverDisease[target.Index].IsActive(), bonus, 1.0)
	}
}

func (dk *Deathknight) applyRime() {
	if dk.Talents.Rime == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 59057}

	dk.RimeAura = dk.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Rime",
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if dk.HowlingBlast == nil {
				return
			}
			dk.HowlingBlast.CD.Reset()
		},
	})
}

func (dk *Deathknight) rimeCritBonus() float64 {
	return 5.0 * float64(dk.Talents.Rime)
}

func (dk *Deathknight) rimeHbChanceProc() float64 {
	return 0.05 * float64(dk.Talents.Rime)
}

func (dk *Deathknight) annihilationCritBonus() float64 {
	return 1.0 * float64(dk.Talents.Annihilation)
}

func (dk *Deathknight) runeSpellComp(spell *core.Spell, hitSpell *RuneSpell) bool {
	if hitSpell == nil {
		return false
	}

	return hitSpell.Spell == spell
}

func (dk *Deathknight) applyKillingMachine() {
	if dk.Talents.KillingMachine == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 51130}

	attackSpeed := 2.0
	if dk.HasMHWeapon() {
		attackSpeed = dk.GetMHWeapon().SwingSpeed
	}

	procChance := attackSpeed * float64(dk.Talents.KillingMachine) / 60.0

	dk.KillingMachineAura = dk.RegisterAura(core.Aura{
		Label:    "Killing Machine Proc",
		ActionID: actionID,
		Duration: time.Second * 30.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.IcyTouch.Spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			dk.FrostStrikeMhHit.Spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			dk.FrostStrikeOhHit.Spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			if dk.HowlingBlast != nil {
				dk.HowlingBlast.Spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.IcyTouch.Spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			dk.FrostStrikeMhHit.Spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			dk.FrostStrikeOhHit.Spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			if dk.HowlingBlast != nil {
				dk.HowlingBlast.Spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			}
		},
	})

	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label: "Killing Machine",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// KM is consumed even if it's a miss
			if dk.KillingMachineAura.IsActive() && (dk.runeSpellComp(spell, dk.IcyTouch) ||
				dk.runeSpellComp(spell, dk.FrostStrike)) {
				dk.KillingMachineAura.Deactivate(sim)
			}

			if !result.Landed() {
				return
			}

			if !spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto) {
				return
			}

			if sim.RandomFloat("Killing Machine Proc Chance") <= procChance {
				dk.KillingMachineAura.Activate(sim)
			}
		},
	}))
}

func (dk *Deathknight) applyDeathchill() {
	if !dk.Talents.Deathchill {
		return
	}

	actionID := core.ActionID{SpellID: 49796}

	rs := &RuneSpell{}
	dk.Deathchill = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: 2 * time.Minute,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.DeathchillAura.Activate(sim)
		},
	}, nil, nil)

	dk.DeathchillAura = dk.RegisterAura(core.Aura{
		Label:    "Deathchill",
		ActionID: actionID,
		Duration: time.Second * 30.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.IcyTouch.Spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			dk.FrostStrikeMhHit.Spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			dk.FrostStrikeOhHit.Spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			dk.ObliterateMhHit.Spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			dk.ObliterateOhHit.Spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			if dk.HowlingBlast != nil {
				dk.HowlingBlast.Spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.IcyTouch.Spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			dk.FrostStrikeMhHit.Spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			dk.FrostStrikeOhHit.Spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			dk.ObliterateMhHit.Spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			dk.ObliterateOhHit.Spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			if dk.HowlingBlast != nil {
				dk.HowlingBlast.Spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if dk.runeSpellComp(spell, dk.IcyTouch) ||
				dk.runeSpellComp(spell, dk.HowlingBlast) ||
				dk.runeSpellComp(spell, dk.FrostStrike) ||
				dk.runeSpellComp(spell, dk.Obliterate) {
				dk.DeathchillAura.Deactivate(sim)
			}
		},
	})
}

func (dk *Deathknight) applyIcyTalons() {
	if dk.Talents.IcyTalons == 0 {
		return
	}

	// Improved Icy Talons
	if dk.Talents.ImprovedIcyTalons {
		dk.PseudoStats.MeleeSpeedMultiplier *= 1.05
	}

	icyTalonsCoeff := 1.0 + 0.04*float64(dk.Talents.IcyTalons)

	dk.IcyTalonsAura = dk.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 50887}, // This probably doesnt need to be in metrics.
		Label:    "Icy Talons",
		Duration: time.Second * 20.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyAttackSpeed(sim, icyTalonsCoeff)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyAttackSpeed(sim, 1/icyTalonsCoeff)
		},
	})
}

func (dk *Deathknight) bloodOfTheNorthCoeff() float64 {
	return []float64{1.0, 1.03, 1.06, 1.10}[dk.Talents.BloodOfTheNorth]
}

func (dk *Deathknight) applyThreatOfThassarian() {
	dk.bonusCoeffs.threatOfThassarianChance = []float64{0.0, 0.3, 0.6, 1.0}[dk.Talents.ThreatOfThassarian]
}

func (dk *Deathknight) threatOfThassarianWillProc(sim *core.Simulation) bool {
	switch dk.bonusCoeffs.threatOfThassarianChance {
	case 0.0:
		return false
	case 1.0:
		return true
	default:
		return sim.RandomFloat("Threat of Thassarian") < dk.bonusCoeffs.threatOfThassarianChance
	}
}

func (dk *Deathknight) threatOfThassarianProcMask(isMH bool) core.ProcMask {
	if isMH {
		return core.ProcMaskMeleeMHSpecial
	} else {
		return core.ProcMaskMeleeOHSpecial
	}
}

func (dk *Deathknight) threatOfThassarianOutcomeApplier(spell *core.Spell) core.OutcomeApplier {
	if spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
		return spell.OutcomeMeleeSpecialHitAndCrit
	} else {
		return spell.OutcomeMeleeSpecialCritOnly
	}
}

func (dk *Deathknight) threatOfThassarianProc(sim *core.Simulation, result *core.SpellResult, ohSpell *RuneSpell) {
	if dk.threatOfThassarianWillProc(sim) && dk.GetOHWeapon() != nil {
		ohSpell.Cast(sim, result.Target)
	}
}
