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

	Ghoul     *GhoulPet
	RaiseDead *core.Spell

	Gargoyle       *GargoylePet
	SummonGargoyle *core.Spell

	Presence Presence

	IcyTouch  *core.Spell
	BloodBoil *core.Spell

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

	HowlingBlastCostless bool
	HowlingBlast         *core.Spell

	EnhancingTotemsActive bool
	HornOfWinter          *core.Spell
	HornOfWinterAura      *core.Aura
	//ArmyOfTheDead    *core.Spell

	// "CDs"
	BloodTap     *core.Spell
	BloodTapAura *core.Aura

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
	UnholyBlightDot   *core.Dot

	// Talent Auras
	KillingMachineAura  *core.Aura
	IcyTalonsAura       *core.Aura
	DesolationAura      *core.Aura
	NecrosisAura        *core.Aura
	BloodCakedBladeAura *core.Aura

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
	IcyTouchAura   *core.Aura
	EbonPlagueAura *core.Aura

	// Dynamic trackers
	RageOfRivendareActive bool
	TundraStalkerActive   bool
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

	raidBuffs.HornOfWinter = !deathKnight.Options.RefreshHornOfWinter

	if raidBuffs.StrengthOfEarthTotem == proto.TristateEffect_TristateEffectImproved {
		deathKnight.EnhancingTotemsActive = true
	} else {
		deathKnight.EnhancingTotemsActive = false
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

	deathKnight.registerRaiseDeadCD()
	deathKnight.registerSummonGargoyleCD()
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

	if deathKnight.Options.PrecastHornOfWinter && deathKnight.Options.RefreshHornOfWinter {
		if deathKnight.HornOfWinterAura.IsActive() {
			deathKnight.HornOfWinterAura.Deactivate(sim)
			deathKnight.HornOfWinterAura.Activate(sim)
		}
	}
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
	}

	maxRunicPower := 100.0 + 15.0*float64(deathKnight.Talents.RunicPowerMastery)
	currentRunicPower := math.Min(maxRunicPower, deathKnightOptions.Options.StartingRunicPower+core.TernaryFloat64(deathKnightOptions.Options.PrecastHornOfWinter, 10.0, 0.0))

	deathKnight.EnableRunicPowerBar(
		currentRunicPower,
		maxRunicPower,
		func(sim *core.Simulation) {
			if deathKnight.GCD.IsReady(sim) {
				deathKnight.tryUseGCD(sim)
			}
		},
		func(sim *core.Simulation) {
			if deathKnight.GCD.IsReady(sim) {
				deathKnight.tryUseGCD(sim)
			}
		},
		func(sim *core.Simulation) {
			if deathKnight.GCD.IsReady(sim) {
				deathKnight.tryUseGCD(sim)
			}
		},
		func(sim *core.Simulation) {
			if deathKnight.GCD.IsReady(sim) {
				deathKnight.tryUseGCD(sim)
			}
		},
		func(sim *core.Simulation) {
			if deathKnight.GCD.IsReady(sim) {
				deathKnight.tryUseGCD(sim)
			}
		},
	)

	deathKnight.EnableAutoAttacks(deathKnight, core.AutoAttackOptions{
		MainHand:       deathKnight.WeaponFromMainHand(deathKnight.DefaultMeleeCritMultiplier()),
		OffHand:        deathKnight.WeaponFromOffHand(deathKnight.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	deathKnight.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.MeleeCrit,
		Modifier: func(agility float64, meleecrit float64) float64 {
			return meleecrit + (agility/62.5)*core.CritRatingPerCritChance
		},
	})
	deathKnight.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.Dodge,
		Modifier: func(agility float64, dodge float64) float64 {
			return dodge + (agility/84.74576271)*core.DodgeRatingPerDodgeChance
		},
	})
	deathKnight.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.AttackPower,
		Modifier: func(strength float64, attackPower float64) float64 {
			return attackPower + strength*2
		},
	})

	deathKnight.Ghoul = deathKnight.NewGhoulPet(deathKnight.Talents.MasterOfGhouls)
	if deathKnight.Talents.SummonGargoyle {
		deathKnight.Gargoyle = deathKnight.NewGargoyle()
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
func (deathKnight *DeathKnight) critMultiplier() float64 {
	applyGuile := deathKnight.Talents.GuileOfGorefiend > 0
	return deathKnight.MeleeCritMultiplier(1.0, deathKnight.secondaryCritModifier(applyGuile))
}
func (deathKnight *DeathKnight) spellCritMultiplier() float64 {
	applyGuile := deathKnight.Talents.GuileOfGorefiend > 0
	return deathKnight.MeleeCritMultiplier(1.0, deathKnight.secondaryCritModifier(applyGuile))
}
func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
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
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
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
