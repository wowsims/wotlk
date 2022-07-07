package core

import (
	"time"
)

const CharacterLevel = 80

const GCDMin = time.Second * 1
const GCDDefault = time.Millisecond * 1500

const MeleeAttackRatingPerDamage = 14.0
const ExpertisePerQuarterPercentReduction = 32.79 / 4 // TODO: Does it still cutoff at 1/4 percents?
const ArmorPenPerPercentArmor = 13.99
const ReducibleArmorConstant = 15232.5

const HasteRatingPerHastePercent = 32.79
const CritRatingPerCritChance = 45.91

const SpellHitRatingPerHitChance = 26.23199272
const MeleeHitRatingPerHitChance = 26.23199272

const DefenseRatingPerDefense = 4.92
const DodgeRatingPerDodgeChance = 45.25
const ParryRatingPerParryChance = 45.25
const BlockRatingPerBlockChance = 16.39
const MissDodgeParryBlockCritChancePerDefense = 0.0325 // TODO: verify this.

const DefenseRatingToChanceReduction = (1.0 / DefenseRatingPerDefense) * MissDodgeParryBlockCritChancePerDefense / 100

const ResilienceRatingPerCritReductionChance = 82.0
const ResilienceRatingPerCritDamageReductionPercent = 39.4231 / 2.2

// With a level 80 attacker you get 32/3 bonus resist per level above attacker
const LevelBasedNPCSpellResistancePerLevel = 32.0 / 3

// TODO: Find these numbers for WOTLK
const EnemyAutoAttackAPCoefficient = 0.000649375

// IDs for items used in core
const (
	ItemIDAtieshMage            = 22589
	ItemIDAtieshWarlock         = 22630
	ItemIDBraidedEterniumChain  = 24114
	ItemIDChainOfTheTwilightOwl = 24121
	ItemIDEyeOfTheNight         = 24116
	ItemIDJadePendantOfBlasting = 20966
	ItemIDTheLightningCapacitor = 28785
)
