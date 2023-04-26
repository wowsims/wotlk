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

const HasteRatingPerHastePercent = 32.79
const CritRatingPerCritChance = 45.91

const SpellHitRatingPerHitChance = 26.232
const MeleeHitRatingPerHitChance = 32.79

const DefenseRatingPerDefense = 4.92
const DodgeRatingPerDodgeChance = 45.25
const ParryRatingPerParryChance = 45.25
const BlockRatingPerBlockChance = 16.39
const MissDodgeParryBlockCritChancePerDefense = 0.04

const DefenseRatingToChanceReduction = (1.0 / DefenseRatingPerDefense) * MissDodgeParryBlockCritChancePerDefense / 100

const ResilienceRatingPerCritReductionChance = 94.27
const ResilienceRatingPerCritDamageReductionPercent = 94.27 / 2.2

// TODO: More log scraping to verify this value for WOTLK.
// Assuming 574 AP debuffs go to exactly zero and achieve -14.2%
const EnemyAutoAttackAPCoefficient = 0.0002883296

const AverageMagicPartialResistMultiplier = 0.94

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

type Hand bool

const MainHand Hand = true
const OffHand Hand = false
