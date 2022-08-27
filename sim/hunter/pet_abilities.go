package hunter

import (
	"math"
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type PetAbilityType int

const (
	Unknown PetAbilityType = iota
	AcidSpit
	Bite
	Claw
	DemoralizingScreech
	FireBreath
	FuriousHowl
	FroststormBreath
	Gore
	LavaBreath
	LightningBreath
	MonstrousBite
	NetherShock
	Pin
	PoisonSpit
	Rake
	Ravage
	SavageRend
	ScorpidPoison
	Smack
	Snatch
	SonicBlast
	SpiritStrike
	SporeCloud
	Stampede
	Sting
	Swipe
	TendonRip
	VenomWebSpray
)

// These IDs are needed for certain talents.
const BiteSpellID = 52474
const ClawSpellID = 52472
const SmackSpellID = 52476

type PetAbility struct {
	Type PetAbilityType

	// Focus cost
	Cost float64

	*core.Spell
}

func (ability *PetAbility) IsEmpty() bool {
	return ability.Spell == nil
}

// Returns whether the ability was successfully cast.
func (ability *PetAbility) TryCast(sim *core.Simulation, target *core.Unit, hp *HunterPet) bool {
	if ability.IsEmpty() {
		return false
	}
	if hp.currentFocus < ability.Cost {
		return false
	}
	if !ability.IsReady(sim) {
		return false
	}

	if !hp.PseudoStats.NoCost {
		hp.SpendFocus(sim, ability.Cost, ability.ActionID)
	}
	ability.Cast(sim, target)
	return true
}

func (hp *HunterPet) NewPetAbility(abilityType PetAbilityType, isPrimary bool) PetAbility {
	switch abilityType {
	case AcidSpit:
		return hp.newAcidSpit()
	case Bite:
		return hp.newBite()
	case Claw:
		return hp.newClaw()
	case DemoralizingScreech:
		return hp.newDemoralizingScreech()
	case FireBreath:
		return hp.newFireBreath()
	case FroststormBreath:
		return hp.newFroststormBreath()
	case FuriousHowl:
		return hp.newFuriousHowl()
	case Gore:
		return hp.newGore()
	case LavaBreath:
		return hp.newLavaBreath()
	case LightningBreath:
		return hp.newLightningBreath()
	case MonstrousBite:
		return hp.newMonstrousBite()
	case NetherShock:
		return hp.newNetherShock()
	case Pin:
		return hp.newPin()
	case PoisonSpit:
		return hp.newPoisonSpit()
	case Rake:
		return hp.newRake()
	case Ravage:
		return hp.newRavage()
	case SavageRend:
		return hp.newSavageRend()
	case ScorpidPoison:
		return hp.newScorpidPoison()
	case Smack:
		return hp.newSmack()
	case Snatch:
		return hp.newSnatch()
	case SonicBlast:
		return hp.newSonicBlast()
	case SpiritStrike:
		return hp.newSpiritStrike()
	case SporeCloud:
		return hp.newSporeCloud()
	case Stampede:
		return hp.newStampede()
	case Sting:
		return hp.newSting()
	case Swipe:
		return hp.newSwipe()
	case TendonRip:
		return hp.newTendonRip()
	case VenomWebSpray:
		return hp.newVenomWebSpray()
	case Unknown:
		return PetAbility{}
	default:
		panic("Invalid pet ability type")
	}
	return PetAbility{}
}

func (hp *HunterPet) newFocusDump(pat PetAbilityType, spellID int32) PetAbility {
	return PetAbility{
		Type: pat,
		Cost: 25,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: spellID},
			SpellSchool: core.SpellSchoolPhysical,
			Flags:       core.SpellFlagMeleeMetrics,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				IgnoreHaste: true,
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskMeleeMHSpecial,
				DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
				ThreatMultiplier: 1,
				BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigMelee(118, 168, 0.07)),
				OutcomeApplier:   hp.OutcomeFuncMeleeSpecialHitAndCrit(2),
			}),
		}),
	}
}

func (hp *HunterPet) newBite() PetAbility {
	return hp.newFocusDump(Bite, BiteSpellID)
}
func (hp *HunterPet) newClaw() PetAbility {
	return hp.newFocusDump(Claw, ClawSpellID)
}
func (hp *HunterPet) newSmack() PetAbility {
	return hp.newFocusDump(Smack, SmackSpellID)
}

type PetSpecialAbilityConfig struct {
	Type    PetAbilityType
	Cost    float64
	SpellID int32
	School  core.SpellSchool
	GCD     time.Duration
	CD      time.Duration
	MinDmg  float64
	MaxDmg  float64
	APRatio float64

	OnSpellHitDealt func(*core.Simulation, *core.Spell, *core.SpellEffect)
}

func (hp *HunterPet) newSpecialAbility(config PetSpecialAbilityConfig) PetAbility {
	var flags core.SpellFlag
	var applyEffects core.ApplySpellEffects
	if config.School == core.SpellSchoolPhysical {
		flags = core.SpellFlagMeleeMetrics
		applyEffects = core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,
			BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigMelee(config.MinDmg, config.MaxDmg, config.APRatio)),
			OutcomeApplier:   hp.OutcomeFuncMeleeSpecialHitAndCrit(2),
			OnSpellHitDealt:  config.OnSpellHitDealt,
		})
	} else {
		applyEffects = core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,
			BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigMelee(config.MinDmg, config.MaxDmg, config.APRatio)),
			OutcomeApplier:   hp.OutcomeFuncMagicHitAndCrit(2),
			OnSpellHitDealt:  config.OnSpellHitDealt,
		})
	}

	return PetAbility{
		Type: config.Type,
		Cost: config.Cost,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: config.SpellID},
			SpellSchool: config.School,
			Flags:       flags,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: config.GCD,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(config.CD),
				},
			},
			ApplyEffects: applyEffects,
		}),
	}
}

func (hp *HunterPet) newAcidSpit() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    AcidSpit,
		Cost:    20,
		GCD:     core.GCDDefault,
		CD:      time.Second * 10,
		SpellID: 55754,
		School:  core.SpellSchoolNature,
		MinDmg:  124,
		MaxDmg:  176,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newDemoralizingScreech() PetAbility {
	var debuffs []*core.Aura
	for _, target := range hp.Env.Encounter.Targets {
		debuffs = append(debuffs, core.ScreechAura(&target.Unit))
	}

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    DemoralizingScreech,
		Cost:    20,
		GCD:     core.GCDDefault,
		CD:      time.Second * 10,
		SpellID: 55487,
		School:  core.SpellSchoolPhysical,
		MinDmg:  85,
		MaxDmg:  129,
		APRatio: 0.07,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				for _, debuff := range debuffs {
					debuff.Activate(sim)
				}
			}
		},
	})
}

func (hp *HunterPet) newFireBreath() PetAbility {
	actionID := core.ActionID{SpellID: 55485}
	var dot *core.Dot

	pa := hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    FireBreath,
		Cost:    20,
		GCD:     core.GCDDefault,
		CD:      time.Second * 10,
		SpellID: 55485,
		School:  core.SpellSchoolFire,
		MinDmg:  43,
		MaxDmg:  57,
		APRatio: 0.049,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				dot.Apply(sim)
			}
		},
	})

	target := hp.CurrentTarget
	dot = core.NewDot(core.Dot{
		Spell: pa.Spell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Fire Breath-" + strconv.Itoa(int(hp.hunterOwner.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 2,
		TickLength:    time.Second * 1,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,
			BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigRoll(44/2, 56/2)),
			OutcomeApplier:   hp.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})

	return pa
}

func (hp *HunterPet) newFroststormBreath() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    FroststormBreath,
		Cost:    20,
		GCD:     0,
		CD:      time.Second * 10,
		SpellID: 55492,
		School:  core.SpellSchoolFrost,
		MinDmg:  128,
		MaxDmg:  172,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newFuriousHowl() PetAbility {
	actionID := core.ActionID{SpellID: 64495}

	petAura := hp.NewTemporaryStatsAura("FuriousHowl", actionID, stats.Stats{stats.AttackPower: 320, stats.RangedAttackPower: 320}, time.Second*20)
	ownerAura := hp.hunterOwner.NewTemporaryStatsAura("FuriousHowl", actionID, stats.Stats{stats.AttackPower: 320, stats.RangedAttackPower: 320}, time.Second*20)
	const cost = 20.0

	howlSpell := hp.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 40),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hp.SpendFocus(sim, cost, actionID)
			petAura.Activate(sim)
			ownerAura.Activate(sim)
		},
	})

	hp.hunterOwner.AddMajorCooldown(core.MajorCooldown{
		Spell: howlSpell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return hp.IsEnabled() && hp.CurrentFocus() >= cost
		},
	})

	return PetAbility{}
}

func (hp *HunterPet) newGore() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Gore,
		Cost:    20,
		GCD:     core.GCDDefault,
		CD:      time.Second * 10,
		SpellID: 35295,
		School:  core.SpellSchoolPhysical,
		MinDmg:  122,
		MaxDmg:  164,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newLavaBreath() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    LavaBreath,
		Cost:    20,
		GCD:     core.GCDDefault,
		CD:      time.Second * 10,
		SpellID: 58611,
		School:  core.SpellSchoolFire,
		MinDmg:  128,
		MaxDmg:  172,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newLightningBreath() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    LightningBreath,
		Cost:    20,
		GCD:     core.GCDDefault,
		CD:      time.Second * 10,
		SpellID: 25012,
		School:  core.SpellSchoolNature,
		MinDmg:  80,
		MaxDmg:  120,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newMonstrousBite() PetAbility {
	procAura := hp.RegisterAura(core.Aura{
		Label:     "Monstrous Bite",
		ActionID:  core.ActionID{SpellID: 55499},
		Duration:  time.Second * 12,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= math.Pow(1.03, float64(oldStacks))
			aura.Unit.PseudoStats.DamageDealtMultiplier *= math.Pow(1.03, float64(newStacks))
		},
	})

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    MonstrousBite,
		Cost:    20,
		GCD:     core.GCDDefault,
		CD:      time.Second * 10,
		SpellID: 55499,
		School:  core.SpellSchoolPhysical,
		MinDmg:  91,
		MaxDmg:  123,
		APRatio: 0.07,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			}
		},
	})
}

func (hp *HunterPet) newNetherShock() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    NetherShock,
		Cost:    20,
		GCD:     core.GCDDefault,
		CD:      time.Second * 10,
		SpellID: 53589,
		School:  core.SpellSchoolShadow,
		MinDmg:  64,
		MaxDmg:  86,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newPin() PetAbility {
	actionID := core.ActionID{SpellID: 53548}
	var dot *core.Dot

	pa := PetAbility{
		Type: Pin,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD:         core.GCDDefault,
					ChannelTime: time.Second * 4,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(time.Second * 40),
				},
			},
			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				ThreatMultiplier: 1,
				OutcomeApplier:   hp.OutcomeFuncMeleeSpecialHit(),
				OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.Landed() {
						dot.Apply(sim)
					}
				},
			}),
		}),
	}

	target := hp.CurrentTarget
	dot = core.NewDot(core.Dot{
		Spell: pa.Spell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Pin-" + strconv.Itoa(int(hp.hunterOwner.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 4,
		TickLength:    time.Second * 1,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,
			BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigMelee(112/4, 144/4, 0.07)),
			OutcomeApplier:   hp.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})

	return pa
}

func (hp *HunterPet) newPoisonSpit() PetAbility {
	actionID := core.ActionID{SpellID: 55557}
	var dot *core.Dot

	pa := PetAbility{
		Type: PoisonSpit,
		Cost: 20,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolNature,
			Flags:       core.SpellFlagIgnoreResists,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(time.Second * 10),
				},
			},
			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				ThreatMultiplier: 1,
				OutcomeApplier:   hp.OutcomeFuncMeleeSpecialHit(),
				OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.Landed() {
						dot.Apply(sim)
					}
				},
			}),
		}),
	}

	target := hp.CurrentTarget
	dot = core.NewDot(core.Dot{
		Spell: pa.Spell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "PoisonSpit-" + strconv.Itoa(int(hp.hunterOwner.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 4,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,
			BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigMelee(104/4, 136/4, 0.049/4)),
			OutcomeApplier:   hp.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})

	return pa
}

func (hp *HunterPet) newRake() PetAbility {
	actionID := core.ActionID{SpellID: 59886}
	var dot *core.Dot

	pa := hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Rake,
		Cost:    20,
		GCD:     core.GCDDefault,
		CD:      time.Second * 10,
		SpellID: 59886,
		School:  core.SpellSchoolPhysical,
		MinDmg:  47,
		MaxDmg:  67,
		APRatio: 0.0175,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				dot.Apply(sim)
			}
		},
	})

	target := hp.CurrentTarget
	dot = core.NewDot(core.Dot{
		Spell: pa.Spell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Rake-" + strconv.Itoa(int(hp.hunterOwner.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 3,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,
			BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigMelee(19, 25, 0.0175)),
			OutcomeApplier:   hp.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})

	return pa
}

func (hp *HunterPet) newRavage() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Ravage,
		Cost:    0,
		CD:      time.Second * 40,
		SpellID: 53562,
		School:  core.SpellSchoolPhysical,
		MinDmg:  106,
		MaxDmg:  150,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newSavageRend() PetAbility {
	actionID := core.ActionID{SpellID: 53582}
	const cost = 20.0
	var dot *core.Dot

	procAura := hp.RegisterAura(core.Aura{
		Label:    "Savage Rend",
		ActionID: actionID,
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
		},
	})

	srSpell := hp.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagApplyArmorReduction,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 60),
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,

			BaseDamage:     hp.specialDamageMod(core.BaseDamageConfigMelee(59, 83, 0.07)),
			OutcomeApplier: hp.OutcomeFuncMeleeSpecialHitAndCrit(2),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				hp.SpendFocus(sim, cost, actionID)
				if spellEffect.Landed() {
					dot.Apply(sim)
					if spellEffect.Outcome.Matches(core.OutcomeCrit) {
						procAura.Activate(sim)
					}
				}
			},
		}),
	})

	target := hp.CurrentTarget
	dot = core.NewDot(core.Dot{
		Spell: srSpell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "SavageRend-" + strconv.Itoa(int(hp.hunterOwner.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 3,
		TickLength:    time.Second * 5,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,
			BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigMelee(21, 27, 0.07)),
			OutcomeApplier:   hp.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})

	hp.hunterOwner.AddMajorCooldown(core.MajorCooldown{
		Spell: srSpell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return hp.IsEnabled() && hp.CurrentFocus() >= cost
		},
	})

	return PetAbility{}
}

func (hp *HunterPet) newScorpidPoison() PetAbility {
	actionID := core.ActionID{SpellID: 55728}
	var dot *core.Dot

	pa := PetAbility{
		Type: ScorpidPoison,
		Cost: 20,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolNature,
			Flags:       core.SpellFlagIgnoreResists,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(time.Second * 10),
				},
			},
			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				ThreatMultiplier: 1,
				OutcomeApplier:   hp.OutcomeFuncMeleeSpecialHit(),
				OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.Landed() {
						dot.Apply(sim)
					}
				},
			}),
		}),
	}

	target := hp.CurrentTarget
	dot = core.NewDot(core.Dot{
		Spell: pa.Spell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "ScorpidPoison-" + strconv.Itoa(int(hp.hunterOwner.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 5,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,
			BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigMelee(100/5, 130/5, 0.07/5)),
			OutcomeApplier:   hp.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})

	return pa
}

func (hp *HunterPet) newSnatch() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Snatch,
		Cost:    20,
		CD:      time.Second * 60,
		SpellID: 53543,
		School:  core.SpellSchoolPhysical,
		MinDmg:  89,
		MaxDmg:  125,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newSonicBlast() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    SonicBlast,
		Cost:    80,
		CD:      time.Second * 60,
		SpellID: 53568,
		School:  core.SpellSchoolNature,
		MinDmg:  62,
		MaxDmg:  88,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newSpiritStrike() PetAbility {
	actionID := core.ActionID{SpellID: 61198}
	var dot *core.Dot

	pa := hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    SpiritStrike,
		Cost:    20,
		GCD:     0,
		CD:      time.Second * 10,
		SpellID: 61198,
		School:  core.SpellSchoolArcane,
		MinDmg:  49,
		MaxDmg:  65,
		APRatio: 0.04,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				dot.Apply(sim)
			}
		},
	})

	target := hp.CurrentTarget
	dot = core.NewDot(core.Dot{
		Spell: pa.Spell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "SpiritStrike-" + strconv.Itoa(int(hp.hunterOwner.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 1,
		TickLength:    time.Second * 6,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,
			BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigMelee(49, 65, 0.04)),
			OutcomeApplier:   hp.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})

	return pa
}

func (hp *HunterPet) newSporeCloud() PetAbility {
	actionID := core.ActionID{SpellID: 53598}
	var dot *core.Dot

	var debuffs []*core.Aura
	for _, target := range hp.Env.Encounter.Targets {
		debuffs = append(debuffs, core.SporeCloudAura(&target.Unit))
	}

	pa := PetAbility{
		Type: SporeCloud,
		Cost: 20,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolNature,
			Flags:       core.SpellFlagIgnoreResists,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(time.Second * 10),
				},
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				dot.Apply(sim)
				for _, debuff := range debuffs {
					debuff.Activate(sim)
				}
			},
		}),
	}

	dot = core.NewDot(core.Dot{
		Spell: pa.Spell,
		Aura: hp.RegisterAura(core.Aura{
			Label:    "SporeCloud-" + strconv.Itoa(int(hp.hunterOwner.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 3,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncAOESnapshot(hp.Env, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,
			BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigMelee(22, 28, 0.049/3)),
			OutcomeApplier:   hp.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})

	return pa
}

func (hp *HunterPet) newStampede() PetAbility {
	debuff := core.StampedeAura(hp.CurrentTarget)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Stampede,
		Cost:    0,
		CD:      time.Second * 60,
		SpellID: 57393,
		School:  core.SpellSchoolPhysical,
		MinDmg:  182,
		MaxDmg:  264,
		APRatio: 0.07,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				debuff.Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) newSting() PetAbility {
	debuff := core.StingAura(hp.CurrentTarget)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Sting,
		Cost:    20,
		GCD:     core.GCDDefault,
		CD:      time.Second * 6,
		SpellID: 56631,
		School:  core.SpellSchoolNature,
		MinDmg:  64,
		MaxDmg:  86,
		APRatio: 0.049,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				debuff.Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) newSwipe() PetAbility {
	// TODO: This is frontal cone, but might be more realistic as single-target
	// since pets are hard to control.
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Swipe,
		Cost:    20,
		GCD:     core.GCDDefault,
		CD:      time.Second * 5,
		SpellID: 53533,
		School:  core.SpellSchoolPhysical,
		MinDmg:  90,
		MaxDmg:  126,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newTendonRip() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    TendonRip,
		Cost:    20,
		CD:      time.Second * 20,
		SpellID: 53575,
		School:  core.SpellSchoolPhysical,
		MinDmg:  49,
		MaxDmg:  69,
		APRatio: 0,
	})
}

func (hp *HunterPet) newVenomWebSpray() PetAbility {
	actionID := core.ActionID{SpellID: 55509}
	var dot *core.Dot

	pa := PetAbility{
		Type: VenomWebSpray,
		Cost: 0,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolNature,
			Flags:       core.SpellFlagIgnoreResists,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(time.Second * 40),
				},
			},
			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				ThreatMultiplier: 1,
				OutcomeApplier:   hp.OutcomeFuncMeleeSpecialHit(),
				OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.Landed() {
						dot.Apply(sim)
					}
				},
			}),
		}),
	}

	target := hp.CurrentTarget
	dot = core.NewDot(core.Dot{
		Spell: pa.Spell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "VenomWebSpray-" + strconv.Itoa(int(hp.hunterOwner.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 4,
		TickLength:    time.Second * 1,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,
			BaseDamage:       hp.specialDamageMod(core.BaseDamageConfigMelee(46, 46, 0.07)),
			OutcomeApplier:   hp.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})

	return pa
}
