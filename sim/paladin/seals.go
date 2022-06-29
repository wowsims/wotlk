package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

const TwistWindow = time.Millisecond * 400
const SealDuration = time.Second * 30

// Handles the cast, gcd, deducts the mana cost
func (paladin *Paladin) setupSealOfBlood() {
	effect := core.SpellEffect{
		ProcMask: core.ProcMaskEmpty,
		DamageMultiplier: 1 *
			paladin.WeaponSpecializationMultiplier() *
			core.TernaryFloat64(ItemSetJusticarArmor.CharacterHasSetBonus(&paladin.Character, 2), 1.1, 1),
		ThreatMultiplier: 1,

		// should deal 35% weapon deamage
		BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 0, 0.35, false),
		OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(paladin.MeleeCritMultiplier()),
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				// Add mana from Spiritual Attunement
				// 10% of damage is self-inflicted, 10% of self-inflicted damage is returned as mana
				paladin.AddMana(sim, spellEffect.Damage*0.1*0.1, paladin.SpiritualAttunementMetrics, false)
			}
		},
	}

	procActionID := core.ActionID{SpellID: 31893}
	sobProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:     procActionID,
		SpellSchool:  core.SpellSchoolHoly,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

	// Define the aura
	paladin.SealOfBloodAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Blood",
		Tag:      "Seal",
		ActionID: procActionID,
		Duration: SealDuration,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}
			sobProc.Cast(sim, spellEffect.Target)
		},
	})

	baseCost := 210.0
	cost := baseCost - paladin.sealCostReduction()
	paladin.SealOfBlood = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31892},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagSeal,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost - baseCost*(0.03*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.ApplyAOEThreat(64)
			paladin.UpdateSeal(sim, paladin.SealOfBloodAura)
		},
	})
}

func (paladin *Paladin) SetupSealOfCommand() {
	procActionID := core.ActionID{SpellID: 20424}
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 0, 0.7, false)

	socProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    procActionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: paladin.WeaponSpecializationMultiplier(),
			ThreatMultiplier: 1,
			OutcomeApplier:   paladin.OutcomeFuncMeleeSpecialHitAndCrit(paladin.MeleeCritMultiplier()),
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return weaponBaseDamage(sim, hitEffect, spell) + 0.29*hitEffect.SpellPower(spell.Unit, spell)
				},
				TargetSpellCoefficient: 0.29,
			},
		}),
	})

	ppmm := paladin.AutoAttacks.NewPPMManager(7.0, core.ProcMaskMelee)
	icd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 1,
	}

	paladin.SealOfCommandAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Command",
		Tag:      "Seal",
		ActionID: procActionID,
		Duration: SealDuration,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			if !ppmm.Proc(sim, spellEffect.ProcMask, "seal of command") {
				return
			}

			icd.Use(sim)
			socProc.Cast(sim, spellEffect.Target)
		},
	})

	baseCost := 65.0
	cost := baseCost - paladin.sealCostReduction()
	paladin.SealOfCommand = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20375},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagSeal,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost - baseCost*(0.03*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.ApplyAOEThreat(70)
			paladin.UpdateSeal(sim, paladin.SealOfCommandAura)
		},
	})
}

// TODO: Make a universal setup seals function

// Seal of the crusader has a bunch of effects that we realistically don't care about (bonus AP, faster swing speed)
// For now, we'll just use it as a setup to casting Judgement of the Crusader
func (paladin *Paladin) setupSealOfTheCrusader() {
	actionID := core.ActionID{SpellID: 27158}
	apBonus := 495.0
	if paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 23203 {
		apBonus += 48
	} else if paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 27949 {
		apBonus += 68
	}

	paladin.SealOfTheCrusaderAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of the Crusader",
		Tag:      "Seal",
		ActionID: actionID,
		Duration: SealDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.AttackPower, apBonus)
			paladin.MultiplyMeleeSpeed(sim, 1.4)
			paladin.AutoAttacks.MHEffect.DamageMultiplier *= 0.6
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.AttackPower, -apBonus)
			paladin.MultiplyMeleeSpeed(sim, 1/1.4)
			paladin.AutoAttacks.MHEffect.DamageMultiplier /= 0.6
		},
	})
	baseCost := 210.0
	cost := baseCost - paladin.sealCostReduction()
	paladin.SealOfTheCrusader = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagSeal,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost - baseCost*(0.03*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.ApplyAOEThreat(61)
			paladin.UpdateSeal(sim, paladin.SealOfTheCrusaderAura)
		},
	})
}

// Didn't encode all the functionality of seal of wisdom
// For now, we'll just use it as a setup to casting Judgement of Wisdom

func (paladin *Paladin) setupSealOfWisdom() {
	actionID := core.ActionID{SpellID: 27166}
	paladin.SealOfWisdomAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Wisdom",
		Tag:      "Seal",
		ActionID: actionID,
		Duration: SealDuration,
	})

	baseCost := 270.0
	cost := baseCost - paladin.sealCostReduction()
	paladin.SealOfWisdom = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagSeal,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost - baseCost*(0.03*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.ApplyAOEThreat(68)
			paladin.UpdateSeal(sim, paladin.SealOfWisdomAura)
		},
	})
}

func (paladin *Paladin) setupSealOfLight() {
	actionID := core.ActionID{SpellID: 27160}
	paladin.SealOfLightAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Light",
		Tag:      "Seal",
		ActionID: actionID,
		Duration: SealDuration,
	})

	baseCost := 280.0
	cost := baseCost - paladin.sealCostReduction()
	paladin.SealOfLight = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagSeal,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost - baseCost*(0.03*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.ApplyAOEThreat(69)
			paladin.UpdateSeal(sim, paladin.SealOfLightAura)
		},
	})
}

func (paladin *Paladin) setupSealOfRighteousness() {
	procActionID := core.ActionID{SpellID: 27156}

	baseCoeff := 0.85
	spCoeff := 0.09
	mhWeapon := paladin.GetMHWeapon()
	if mhWeapon != nil && mhWeapon.HandType == proto.HandType_HandTypeTwoHand {
		baseCoeff = 1.2
		spCoeff = 0.108
	}
	spCoeff *= paladin.AutoAttacks.MH.SwingSpeed
	flatBonusDamage := baseCoeff * 28.72464 * paladin.AutoAttacks.MH.SwingSpeed
	if mhWeapon != nil {
		flatBonusDamage += -mhWeapon.QualityModifier * paladin.AutoAttacks.MH.SwingSpeed * 0.03
	}

	baseDamage := core.BaseDamageConfig{
		Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
			weaponDamage := spell.Unit.AutoAttacks.MH.CalculateAverageWeaponDamage(hitEffect.MeleeAttackPower(spell.Unit))
			return flatBonusDamage + 0.03*weaponDamage + spCoeff*hitEffect.SpellPower(spell.Unit, spell)
		},
		TargetSpellCoefficient: spCoeff,
	}

	sorProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    procActionID,
		SpellSchool: core.SpellSchoolHoly,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskEmpty,

			BonusSpellPower: core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 33504, 94, 0),
			DamageMultiplier: 1 *
				paladin.WeaponSpecializationMultiplier() *
				(1 + 0.03*float64(paladin.Talents.ImprovedSealOfRighteousness)) *
				core.TernaryFloat64(ItemSetJusticarArmor.CharacterHasSetBonus(&paladin.Character, 2), 1.1, 1),
			ThreatMultiplier: 1,

			BaseDamage:     baseDamage,
			OutcomeApplier: paladin.OutcomeFuncAlwaysHit(),
		}),
	})

	spellActionID := core.ActionID{SpellID: 27155}
	paladin.SealOfRighteousnessAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Righteousness",
		Tag:      "Seal",
		ActionID: spellActionID,
		Duration: SealDuration,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}
			sorProc.Cast(sim, spellEffect.Target)
		},
	})

	baseCost := 260.0 - paladin.sealCostReduction()

	paladin.SealOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    spellActionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagSeal,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.03*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.ApplyAOEThreat(66)
			paladin.UpdateSeal(sim, paladin.SealOfRighteousnessAura)
		},
	})
}

func (paladin *Paladin) UpdateSeal(sim *core.Simulation, newSeal *core.Aura) {
	if paladin.SealOfCommandAura != nil && paladin.CurrentSeal == paladin.SealOfCommandAura {
		// Technically the current expiration could be shorter than 0.4 seconds
		// TO-DO: Lookup behavior when seal of command is twisted at shorter than 0.4 seconds duration
		expiresAt := sim.CurrentTime + TwistWindow
		paladin.CurrentSeal.UpdateExpires(expiresAt)

		// This is a hack to get the sim to process and log the SoC aura expiring at the right time
		if sim.Options.Iterations == 1 {
			sim.AddPendingAction(&core.PendingAction{
				NextActionAt: expiresAt,
				OnAction:     func(_ *core.Simulation) {},
			})
		}
	} else if paladin.CurrentSeal != nil {
		paladin.CurrentSeal.Deactivate(sim)
	}

	paladin.CurrentSeal = newSeal
	newSeal.Activate(sim)
}

func (paladin *Paladin) sealCostReduction() float64 {
	switch paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID {
	case 22401: // libram of hope
		return -20
	case 186067: // communal book of righteousness
		return -5
	}
	return 0
}
