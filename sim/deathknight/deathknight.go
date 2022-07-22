package deathknight

import (
	"math"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type DeathKnight struct {
	core.Character
	Talents  proto.DeathKnightTalents
	Options  proto.DeathKnight_Options
	Rotation proto.DeathKnight_Rotation

	LastCastOutcome core.HitOutcome
	DKRotation      DKRotation

	Ghoul     *GhoulPet
	RaiseDead *core.Spell

	Gargoyle       *GargoylePet
	SummonGargoyle *core.Spell

	ArmyOfTheDead *core.Spell
	ArmyGhoul     []*GhoulPet

	Presence Presence

	IcyTouch   *core.Spell
	BloodBoil  *core.Spell
	Pestilence *core.Spell

	PlagueStrike      *core.Spell
	PlagueStrikeMhHit *core.Spell
	PlagueStrikeOhHit *core.Spell

	Obliterate      *core.Spell
	ObliterateMhHit *core.Spell
	ObliterateOhHit *core.Spell

	BloodStrike      *core.Spell
	BloodStrikeMhHit *core.Spell
	BloodStrikeOhHit *core.Spell

	FrostStrike      *core.Spell
	FrostStrikeMhHit *core.Spell
	FrostStrikeOhHit *core.Spell

	GhoulFrenzy *core.Spell
	// Dummy aura for timeline metrics
	GhoulFrenzyAura *core.Aura

	LastScourgeStrikeDamage float64
	ScourgeStrike           *core.Spell

	LastDeathCoilDamage float64
	DeathCoil           *core.Spell

	DeathAndDecay    *core.Spell
	DeathAndDecayDot *core.Dot

	HowlingBlast *core.Spell

	OtherRelevantStrAgiActive bool
	HornOfWinter              *core.Spell
	HornOfWinterAura          *core.Aura
	//ArmyOfTheDead    *core.Spell

	// "CDs"
	BloodTap     *core.Spell
	BloodTapAura *core.Aura

	EmpowerRuneWeapon *core.Spell

	UnbreakableArmor     *core.Spell
	UnbreakableArmorAura *core.Aura

	BoneShield     *core.Spell
	BoneShieldAura *core.Aura

	// Diseases
	FrostFeverSpell    *core.Spell
	BloodPlagueSpell   *core.Spell
	FrostFeverDisease  []*core.Dot
	BloodPlagueDisease []*core.Dot

	UnholyBlightSpell *core.Spell
	UnholyBlightDot   []*core.Dot

	// Talent Auras
	KillingMachineAura  *core.Aura
	IcyTalonsAura       *core.Aura
	DesolationAura      *core.Aura
	NecrosisAura        *core.Aura
	BloodCakedBladeAura *core.Aura
	ButcheryAura        *core.Aura
	RimeAura            *core.Aura

	// Talent Spells
	LastDiseaseDamage float64
	WanderingPlague   *core.Spell

	// Presences
	BloodPressence     *core.Spell
	BloodPresenceAura  *core.Aura
	FrostPressence     *core.Spell
	FrostPresenceAura  *core.Aura
	UnholyPressence    *core.Spell
	UnholyPresenceAura *core.Aura

	// Debuffs
	IcyTouchAura   []*core.Aura
	CryptFeverAura []*core.Aura
	EbonPlagueAura []*core.Aura

	// Dynamic trackers
	additiveDamageModifier float64
}

func (deathKnight *DeathKnight) ModifyAdditiveDamageModifier(sim *core.Simulation, value float64) {
	deathKnight.PseudoStats.DamageDealtMultiplier /= deathKnight.additiveDamageModifier
	deathKnight.additiveDamageModifier += value
	deathKnight.PseudoStats.DamageDealtMultiplier *= deathKnight.additiveDamageModifier
}

func (deathKnight *DeathKnight) GetCharacter() *core.Character {
	return &deathKnight.Character
}

func (deathKnight *DeathKnight) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (deathKnight *DeathKnight) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if deathKnight.Talents.AbominationsMight > 0 {
		raidBuffs.AbominationsMight = true
	}

	if deathKnight.Talents.ImprovedIcyTalons {
		raidBuffs.IcyTalons = true
	}

	raidBuffs.HornOfWinter = !deathKnight.Rotation.RefreshHornOfWinter

	if raidBuffs.StrengthOfEarthTotem == proto.TristateEffect_TristateEffectImproved ||
		raidBuffs.StrengthOfEarthTotem == proto.TristateEffect_TristateEffectRegular {
		deathKnight.OtherRelevantStrAgiActive = true
	} else {
		deathKnight.OtherRelevantStrAgiActive = false
	}
}

func (deathKnight *DeathKnight) ApplyTalents() {
	deathKnight.ApplyBloodTalents()
	deathKnight.ApplyFrostTalents()
	deathKnight.ApplyUnholyTalents()
}

func (deathKnight *DeathKnight) Initialize() {
	deathKnight.registerPresences()
	deathKnight.registerIcyTouchSpell()
	deathKnight.registerPlagueStrikeSpell()
	deathKnight.registerObliterateSpell()
	deathKnight.registerBloodStrikeSpell()
	deathKnight.registerBloodTapSpell()
	deathKnight.registerHowlingBlastSpell()
	deathKnight.registerScourgeStrikeSpell()
	deathKnight.registerDeathCoilSpell()
	deathKnight.registerFrostStrikeSpell()
	deathKnight.registerDeathAndDecaySpell()
	deathKnight.registerDiseaseDots()
	deathKnight.registerGhoulFrenzySpell()
	deathKnight.registerBoneShieldSpell()
	deathKnight.registerUnbreakableArmorSpell()
	deathKnight.registerBloodBoilSpell()
	deathKnight.registerHornOfWinterSpell()
	deathKnight.registerPestilenceSpell()
	deathKnight.registerEmpowerRuneWeaponSpell()

	deathKnight.registerRaiseDeadCD()
	deathKnight.registerSummonGargoyleCD()
	deathKnight.registerArmyOfTheDeadCD()

	deathKnight.SetupRotation()
}

func (deathKnight *DeathKnight) Reset(sim *core.Simulation) {
	deathKnight.ResetRunicPowerBar(sim)

	if deathKnight.Rotation.UnholyPresenceOpener {
		deathKnight.UnholyPresenceAura.Activate(sim)
		deathKnight.Presence = UnholyPresence
	} else {
		deathKnight.BloodPresenceAura.Activate(sim)
		deathKnight.Presence = BloodPresence
	}

	if deathKnight.Rotation.ArmyOfTheDead == proto.DeathKnight_Rotation_PreCast {
		deathKnight.PrecastArmyOfTheDead(sim)
	}

	deathKnight.resetDKRotation(sim)
}

func (deathKnight *DeathKnight) IsFuStrike(spell *core.Spell) bool {
	return spell == deathKnight.Obliterate || spell == deathKnight.ScourgeStrike // || spell == deathKnight.DeathStrike
}

func (deathKnight *DeathKnight) HasMajorGlyph(glyph proto.DeathKnightMajorGlyph) bool {
	return deathKnight.HasGlyph(int32(glyph))
}
func (deathKnight *DeathKnight) HasMinorGlyph(glyph proto.DeathKnightMinorGlyph) bool {
	return deathKnight.HasGlyph(int32(glyph))
}

func NewDeathKnight(character core.Character, options proto.Player) *DeathKnight {
	deathKnightOptions := options.GetDeathKnight()

	deathKnight := &DeathKnight{
		Character: character,
		Talents:   *deathKnightOptions.Talents,
		Options:   *deathKnightOptions.Options,
		Rotation:  *deathKnightOptions.Rotation,

		additiveDamageModifier: 1,
	}

	maxRunicPower := 100.0 + 15.0*float64(deathKnight.Talents.RunicPowerMastery)
	currentRunicPower := math.Min(maxRunicPower, deathKnightOptions.Options.StartingRunicPower+core.TernaryFloat64(deathKnightOptions.Options.PrecastHornOfWinter, 10.0, 0.0))

	deathKnight.EnableRunicPowerBar(
		currentRunicPower,
		maxRunicPower,
		func(sim *core.Simulation) {
			if !deathKnight.Talents.HowlingBlast {
				if deathKnight.GCD.IsReady(sim) {
					deathKnight.tryUseGCD(sim)
				}
			} else {
				if !deathKnight.DKRotation.onOpener {
					if deathKnight.GCD.IsReady(sim) {
						deathKnight.tryUseGCD(sim)
					}
				}
			}
		},
		func(sim *core.Simulation) {
			if !deathKnight.Talents.HowlingBlast {
				if deathKnight.GCD.IsReady(sim) {
					deathKnight.tryUseGCD(sim)
				}
			} else {
				if !deathKnight.DKRotation.onOpener {
					if deathKnight.GCD.IsReady(sim) {
						deathKnight.tryUseGCD(sim)
					}
				}
			}
		},
		func(sim *core.Simulation) {
			if !deathKnight.Talents.HowlingBlast {
				if deathKnight.GCD.IsReady(sim) {
					deathKnight.tryUseGCD(sim)
				}
			} else {
				if !deathKnight.DKRotation.onOpener {
					if deathKnight.GCD.IsReady(sim) {
						deathKnight.tryUseGCD(sim)
					}
				}
			}
		},
		func(sim *core.Simulation) {
			if !deathKnight.Talents.HowlingBlast {
				if deathKnight.GCD.IsReady(sim) {
					deathKnight.tryUseGCD(sim)
				}
			} else {
				if !deathKnight.DKRotation.onOpener {
					if deathKnight.GCD.IsReady(sim) {
						deathKnight.tryUseGCD(sim)
					}
				}
			}
		},
		func(sim *core.Simulation) {
			if !deathKnight.Talents.HowlingBlast {
				if deathKnight.GCD.IsReady(sim) {
					deathKnight.tryUseGCD(sim)
				}
			} else {
				if !deathKnight.DKRotation.onOpener {
					if deathKnight.GCD.IsReady(sim) {
						deathKnight.tryUseGCD(sim)
					}
				}
			}
		},
	)

	deathKnight.EnableAutoAttacks(deathKnight, core.AutoAttackOptions{
		MainHand:       deathKnight.WeaponFromMainHand(deathKnight.DefaultMeleeCritMultiplier()),
		OffHand:        deathKnight.WeaponFromOffHand(deathKnight.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	deathKnight.AddStatDependency(stats.Agility, stats.MeleeCrit, 1.0+(core.CritRatingPerCritChance/62.5))
	deathKnight.AddStatDependency(stats.Agility, stats.Dodge, 1.0+(core.DodgeRatingPerDodgeChance/84.74576271))
	deathKnight.AddStatDependency(stats.Strength, stats.AttackPower, 1.0+2)

	deathKnight.Ghoul = deathKnight.NewGhoulPet(deathKnight.Talents.MasterOfGhouls)
	if deathKnight.Talents.SummonGargoyle {
		deathKnight.Gargoyle = deathKnight.NewGargoyle()
	}

	deathKnight.ArmyGhoul = make([]*GhoulPet, 8)
	for i := 0; i < 8; i++ {
		deathKnight.ArmyGhoul[i] = deathKnight.NewArmyGhoulPet(i)
	}

	return deathKnight
}

func RegisterDeathKnight() {
	core.RegisterAgentFactory(
		proto.Player_DeathKnight{},
		proto.Spec_SpecDeathKnight,
		func(character core.Character, options proto.Player) core.Agent {
			return NewDeathKnight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DeathKnight)
			if !ok {
				panic("Invalid spec value for DeathKnight!")
			}
			player.Spec = playerSpec
		},
	)
}

func (deathKnight *DeathKnight) AllDiseasesAreActive(target *core.Unit) bool {
	return deathKnight.FrostFeverDisease[target.Index].IsActive() && deathKnight.BloodPlagueDisease[target.Index].IsActive()
}

func (deathKnight *DeathKnight) DiseasesAreActive(target *core.Unit) bool {
	return deathKnight.FrostFeverDisease[target.Index].IsActive() || deathKnight.BloodPlagueDisease[target.Index].IsActive()
}

func (deathKnight *DeathKnight) secondaryCritModifier(applyGuile bool) float64 {
	secondaryModifier := 0.0
	if applyGuile {
		secondaryModifier += 0.15 * float64(deathKnight.Talents.GuileOfGorefiend)
	}
	return secondaryModifier
}

// TODO: DKs have x2 modifier on spell crit as a passive. Is this the best way to do it?
func (deathKnight *DeathKnight) spellCritMultiplier() float64 {
	return deathKnight.MeleeCritMultiplier(1.0, 0)
}
func (deathKnight *DeathKnight) spellCritMultiplierGuile() float64 {
	applyGuile := deathKnight.Talents.GuileOfGorefiend > 0
	return deathKnight.MeleeCritMultiplier(1.0, deathKnight.secondaryCritModifier(applyGuile))
}
func (deathKnight *DeathKnight) critMultiplier() float64 {
	return deathKnight.MeleeCritMultiplier(1.0, 0)
}
func (deathKnight *DeathKnight) critMultiplierGuile() float64 {
	applyGuile := deathKnight.Talents.GuileOfGorefiend > 0
	return deathKnight.MeleeCritMultiplier(1.0, deathKnight.secondaryCritModifier(applyGuile))
}

func DetermineOptimalCostForSpell(rp *core.CalcRunicPowerBar, sim *core.Simulation, deathKnight *DeathKnight, spell *core.Spell) core.DKRuneCost {
	blood := 0
	frost := 0
	unholy := 0
	switch spell {
	case deathKnight.DeathAndDecay:
		blood = 1
		frost = 1
		unholy = 1
	case deathKnight.ArmyOfTheDead:
		blood = 1
		frost = 1
		unholy = 1
	case deathKnight.Pestilence:
		blood = 1
	case deathKnight.BloodStrike:
		blood = 1
	case deathKnight.BloodBoil:
		blood = 1
	case deathKnight.UnbreakableArmor:
		frost = 1
	case deathKnight.IcyTouch:
		frost = 1
	case deathKnight.PlagueStrike:
		unholy = 1
	case deathKnight.GhoulFrenzy:
		unholy = 1
	case deathKnight.BoneShield:
		unholy = 1
	case deathKnight.ScourgeStrike:
		frost = 1
		unholy = 1
	case deathKnight.Obliterate:
		frost = 1
		unholy = 1
	case deathKnight.HowlingBlast:
		frost = 1
		unholy = 1
	}

	return rp.DetermineOptimalCost(sim, blood, frost, unholy)
}

func (deathKnight *DeathKnight) CanCast(sim *core.Simulation, spell *core.Spell) bool {
	switch spell {
	case deathKnight.DeathAndDecay:
		return deathKnight.CanDeathAndDecay(sim)
	case deathKnight.ArmyOfTheDead:
		return deathKnight.CanArmyOfTheDead(sim)
	case deathKnight.Pestilence:
		return deathKnight.CanPestilence(sim)
	case deathKnight.BloodStrike:
		return deathKnight.CanBloodStrike(sim)
	case deathKnight.BloodBoil:
		return deathKnight.CanBloodBoil(sim)
	case deathKnight.UnbreakableArmor:
		return deathKnight.CanUnbreakableArmor(sim)
	case deathKnight.IcyTouch:
		return deathKnight.CanIcyTouch(sim)
	case deathKnight.PlagueStrike:
		return deathKnight.CanPlagueStrike(sim)
	case deathKnight.GhoulFrenzy:
		return deathKnight.CanGhoulFrenzy(sim)
	case deathKnight.BoneShield:
		return deathKnight.CanBoneShield(sim)
	case deathKnight.ScourgeStrike:
		return deathKnight.CanScourgeStrike(sim)
	case deathKnight.Obliterate:
		return deathKnight.CanObliterate(sim)
	case deathKnight.HowlingBlast:
		return deathKnight.CanHowlingBlast(sim)
	case deathKnight.FrostStrike:
		return deathKnight.CanFrostStrike(sim)
	case deathKnight.DeathCoil:
		return deathKnight.CanDeathCoil(sim)
	case deathKnight.BloodTap:
		return deathKnight.CanBloodTap(sim)
	case deathKnight.EmpowerRuneWeapon:
		return deathKnight.CanEmpowerRuneWeapon(sim)
	case deathKnight.HornOfWinter:
		return deathKnight.CanHornOfWinter(sim)
	case deathKnight.RaiseDead:
		return deathKnight.CanRaiseDead(sim)
	default:
		panic("Not in cost list.")
	}

	return false
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    176,
		stats.Agility:     109,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      61,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     108,
		stats.Stamina:     161,
		stats.Intellect:   34,
		stats.Spirit:      58,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    170,
		stats.Agility:     114,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    175,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    171,
		stats.Agility:     116,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      59,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    178,
		stats.Agility:     109,
		stats.Stamina:     161,
		stats.Intellect:   32,
		stats.Spirit:      61,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     108,
		stats.Stamina:     161,
		stats.Intellect:   31,
		stats.Spirit:      61,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    176,
		stats.Agility:     114,
		stats.Stamina:     160,
		stats.Intellect:   31,
		stats.Spirit:      60,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    174,
		stats.Agility:     110,
		stats.Stamina:     160,
		stats.Intellect:   33,
		stats.Spirit:      64,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    172,
		stats.Agility:     114,
		stats.Stamina:     160,
		stats.Intellect:   38,
		stats.Spirit:      57,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
}

// Agent is a generic way to access underlying warrior on any of the agents.

func (deathKnight *DeathKnight) GetDeathKnight() *DeathKnight {
	return deathKnight
}

type DeathKnightAgent interface {
	GetDeathKnight() *DeathKnight
}
