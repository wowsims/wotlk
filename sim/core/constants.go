package core

import (
	"time"
)

const CharacterLevel = 70

const GCDMin = time.Second * 1
const GCDDefault = time.Millisecond * 1500

const HasteRatingPerHastePercent = 15.77

const MeleeCritRatingPerCritChance = 22.08
const MeleeHitRatingPerHitChance = 15.77
const MeleeAttackRatingPerDamage = 14.0

const ExpertisePerQuarterPercentReduction = 3.94
const ArmorPenPerPercentArmor = 5.92

const SpellCritRatingPerCritChance = 22.08
const SpellHitRatingPerHitChance = 12.62

const DefenseRatingPerDefense = 2.3654
const MissDodgeParryBlockCritChancePerDefense = 0.04
const BlockRatingPerBlockChance = 7.8846
const DodgeRatingPerDodgeChance = 18.9231
const ParryRatingPerParryChance = 23.6538
const ResilienceRatingPerCritReductionChance = 39.4231
const ResilienceRatingPerCritDamageReductionPercent = 39.4231 / 2
const DefenseRatingToChanceReduction = (1.0 / DefenseRatingPerDefense) * MissDodgeParryBlockCritChancePerDefense / 100

const LevelBasedNPCSpellResistancePerLevel = 28.0 / 3

const EnemyAutoAttackAPCoefficient = 0.000649375
const CrushChance = 0.15

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
