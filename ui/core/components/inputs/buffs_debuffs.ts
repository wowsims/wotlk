import { Stat, TristateEffect } from "../../proto/common";
import { ActionId } from "../../proto_utils/action_id";

import {
  makeBooleanDebuffInput,
  makeBooleanIndividualBuffInput,
  makeBooleanPartyBuffInput,
  makeBooleanRaidBuffInput,
  makeMultistateIndividualBuffInput,
  makeMultistateMultiplierIndividualBuffInput,
  makeMultistatePartyBuffInput,
  makeMultistateRaidBuffInput,
  makeQuadstateDebuffInput,
  makeTristateDebuffInput,
  makeTristateIndividualBuffInput,
  makeTristateRaidBuffInput,
  withLabel,
} from "../icon_inputs";
import { IconPicker } from "../icon_picker";
import { MultiIconPicker } from "../multi_icon_picker";

import { IconPickerStatOption, PickerStatOptions } from "./stat_options";

import * as InputHelpers from '../input_helpers';

///////////////////////////////////////////////////////////////////////////
//                                 RAID BUFFS
///////////////////////////////////////////////////////////////////////////

export const AllStatsBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(48470), impId: ActionId.fromSpellId(17051), fieldName: 'giftOfTheWild'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromItemId(49634), fieldName: 'drumsOfTheWild'}),
], 'Stats');

export const AllStatsPercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(25898), fieldName: 'blessingOfKings'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromItemId(49633), fieldName: 'drumsOfForgottenKings'}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(25899), fieldName: 'blessingOfSanctuary'}),
], 'Stats %');

export const ArmorBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(48942), impId: ActionId.fromSpellId(20140), fieldName: 'devotionAura'}),
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(58753), impId: ActionId.fromSpellId(16293), fieldName: 'stoneskinTotem'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromItemId(43468), fieldName: 'scrollOfProtection'}),
], 'Armor');

export const AttackPowerBuff = InputHelpers.makeMultiIconInput([
	makeTristateIndividualBuffInput({actionId: ActionId.fromSpellId(48934), impId: ActionId.fromSpellId(20045), fieldName: 'blessingOfMight'}),
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(47436), impId: ActionId.fromSpellId(12861), fieldName: 'battleShout'}),
], 'AP');

export const AttackPowerPercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(53138), fieldName: 'abominationsMight'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(30809), fieldName: 'unleashedRage'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(19506), fieldName: 'trueshotAura'}),
], 'Atk Pwr %');

export const Bloodlust = withLabel(
  makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(2825), fieldName: 'bloodlust'}),
  'Lust',
);

export const DamagePercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(31869), fieldName: 'sanctifiedRetribution'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(31583), fieldName: 'arcaneEmpowerment'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(34460), fieldName: 'ferociousInspiration'}),
], 'Dmg %');

export const DamageReductionPercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(57472), fieldName: 'renewedHope'}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(25899), fieldName: 'blessingOfSanctuary'}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(50720), fieldName: 'vigilance'}),
], 'Mit %');

export const DefensiveCooldownBuff = InputHelpers.makeMultiIconInput([
	makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(6940), numStates: 11, fieldName: 'handOfSacrifices'}),
	makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(53530), numStates: 11, fieldName: 'divineGuardians'}),
	makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(33206), numStates: 11, fieldName: 'painSuppressions'}),
	makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(47788), numStates: 11, fieldName: 'guardianSpirits'}),
], 'Defensive CDs');

export const HastePercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(53648), fieldName: 'swiftRetribution'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(48396), fieldName: 'moonkinAura', value: TristateEffect.TristateEffectImproved}),
], 'Haste %');

export const HealthBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(47440), impId: ActionId.fromSpellId(12861), fieldName: 'commandingShout'}),
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(47982), impId: ActionId.fromSpellId(18696), fieldName: 'bloodPact'}),
], 'Health');

export const IntellectBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(43002), fieldName: 'arcaneBrilliance'}),
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(57567), impId: ActionId.fromSpellId(54038), fieldName: 'felIntelligence'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromItemId(37092), fieldName: 'scrollOfIntellect'}),
], 'Int');

export const MeleeCritBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(17007), impId: ActionId.fromSpellId(34300), fieldName: 'leaderOfThePack'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(29801), fieldName: 'rampage'}),
], 'Melee Crit');

export const MeleeHasteBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(55610), fieldName: 'icyTalons'}),
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(65990), impId: ActionId.fromSpellId(29193), fieldName: 'windfuryTotem'}),
], 'Melee Haste');

export const MP5Buff = InputHelpers.makeMultiIconInput([
	makeTristateIndividualBuffInput({actionId: ActionId.fromSpellId(48938), impId: ActionId.fromSpellId(20245), fieldName: 'blessingOfWisdom'}),
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(58774), impId: ActionId.fromSpellId(16206), fieldName: 'manaSpringTotem'}),
], 'MP5');

export const ReplenishmentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(48160), fieldName: 'vampiricTouch'}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(31878), fieldName: 'judgementsOfTheWise'}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(53292), fieldName: 'huntingParty'}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(54118), fieldName: 'improvedSoulLeech'}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(44561), fieldName: 'enduringWinter'}),
], 'Replen');

export const ResistanceBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(48170), fieldName: 'shadowProtection'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(58749), fieldName: 'natureResistanceTotem'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(49071), fieldName: 'aspectOfTheWild'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(48945), fieldName: 'frostResistanceAura'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(58745), fieldName: 'frostResistanceTotem'}),
], 'Resistances');

export const RevitalizeBuff = InputHelpers.makeMultiIconInput([
  makeMultistateMultiplierIndividualBuffInput(ActionId.fromSpellId(26982), 101, 10, 'revitalizeRejuvination'),
  makeMultistateMultiplierIndividualBuffInput(ActionId.fromSpellId(53251), 101, 10, 'revitalizeWildGrowth'),
], 'Revit', ActionId.fromSpellId(48545))

export const SpellCritBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(24907), impId: ActionId.fromSpellId(48396), fieldName: 'moonkinAura'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(51470), fieldName: 'elementalOath'}),
], 'Spell Crit');

export const SpellHasteBuff = withLabel(
  makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(3738), fieldName: 'wrathOfAirTotem'}),
  'Spell Haste',
);

export const SpellPowerBuff = InputHelpers.makeMultiIconInput([
	makeMultistateRaidBuffInput({actionId: ActionId.fromSpellId(47240), numStates: 2000, fieldName: 'demonicPactSp', multiplier: 20}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(57722), fieldName: 'totemOfWrath'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(58656), fieldName: 'flametongueTotem'}),
], 'Spell Power');

export const SpiritBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(48073), fieldName: 'divineSpirit'}),
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(57567), impId: ActionId.fromSpellId(54038), fieldName: 'felIntelligence'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromItemId(37098), fieldName: 'scrollOfSpirit'}),
], 'Spirit');

export const StaminaBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(48161), impId: ActionId.fromSpellId(14767), fieldName: 'powerWordFortitude'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromItemId(37094), fieldName: 'scrollOfStamina'}),
], 'Stamina');

export const StrengthAndAgilityBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(58643), impId: ActionId.fromSpellId(52456), fieldName: 'strengthOfEarthTotem'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(57623), fieldName: 'hornOfWinter'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromItemId(43464), fieldName: 'scrollOfAgility'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromItemId(43466), fieldName: 'scrollOfStrength'}),
], 'Str/Agi');

// Misc Buffs
export const StrengthOfWrynn = makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(73828), fieldName: 'strengthOfWrynn'});
export const RetributionAura = makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(54043), fieldName: 'retributionAura'});
export const BraidedEterniumChain = makeBooleanPartyBuffInput({actionId: ActionId.fromSpellId(31025), fieldName: 'braidedEterniumChain'});
export const ChainOfTheTwilightOwl = makeBooleanPartyBuffInput({actionId: ActionId.fromSpellId(31035), fieldName: 'chainOfTheTwilightOwl'});
export const HeroicPresence = makeBooleanPartyBuffInput({actionId: ActionId.fromSpellId(6562), fieldName: 'heroicPresence'});
export const EyeOfTheNight = makeBooleanPartyBuffInput({actionId: ActionId.fromSpellId(31033), fieldName: 'eyeOfTheNight'});
export const Thorns = makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(53307), impId: ActionId.fromSpellId(16840), fieldName: 'thorns'});
export const ManaTideTotem = makeMultistatePartyBuffInput(ActionId.fromSpellId(16190), 5, 'manaTideTotems');
export const Innervate = makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(29166), numStates: 11, fieldName: 'innervates'});
export const PowerInfusion = makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(10060), numStates: 11, fieldName: 'powerInfusions'});
export const FocusMagic = makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(54648), fieldName: 'focusMagic'});
export const TricksOfTheTrade = makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(57933), numStates: 20, fieldName: 'tricksOfTheTrades'});
export const UnholyFrenzy = makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(49016), numStates: 11, fieldName: 'unholyFrenzy'});

///////////////////////////////////////////////////////////////////////////
//                                 DEBUFFS
///////////////////////////////////////////////////////////////////////////

export const MajorArmorDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(47467), fieldName: 'sunderArmor'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(8647), fieldName: 'exposeArmor'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(55754), fieldName: 'acidSpit'}),
], 'Major ArP');

export const MinorArmorDebuff = InputHelpers.makeMultiIconInput([
	makeTristateDebuffInput({actionId: ActionId.fromSpellId(770), impId: ActionId.fromSpellId(33602), fieldName: 'faerieFire'}),
	makeTristateDebuffInput({actionId: ActionId.fromSpellId(50511), impId: ActionId.fromSpellId(18180), fieldName: 'curseOfWeakness'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(56631), fieldName: 'sting'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(53598), fieldName: 'sporeCloud'}),
], 'Minor ArP');

export const AttackPowerDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(26016), fieldName: 'vindication'}),
	makeTristateDebuffInput({actionId: ActionId.fromSpellId(47437), impId: ActionId.fromSpellId(12879), fieldName: 'demoralizingShout'}),
	makeTristateDebuffInput({actionId: ActionId.fromSpellId(48560), impId: ActionId.fromSpellId(16862), fieldName: 'demoralizingRoar'}),
	makeTristateDebuffInput({actionId: ActionId.fromSpellId(50511), impId: ActionId.fromSpellId(18180), fieldName: 'curseOfWeakness'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(55487), fieldName: 'demoralizingScreech'}),
], 'Atk Pwr');

export const BleedDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(48564), fieldName: 'mangle'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(46855), fieldName: 'trauma'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(57393), fieldName: 'stampede'}),
], 'Bleed');

export const CritDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(30706), fieldName: 'totemOfWrath'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(20337), fieldName: 'heartOfTheCrusader'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(58410), fieldName: 'masterPoisoner'}),
], 'Crit');

export const MeleeAttackSpeedDebuff = InputHelpers.makeMultiIconInput([
	makeTristateDebuffInput({actionId: ActionId.fromSpellId(47502), impId: ActionId.fromSpellId(12666), fieldName: 'thunderClap'}),
	makeTristateDebuffInput({actionId: ActionId.fromSpellId(55095), impId: ActionId.fromSpellId(51456), fieldName: 'frostFever'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(53696), fieldName: 'judgementsOfTheJust'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(48485), fieldName: 'infectedWounds'}),
], 'Atk Speed');

export const MeleeHitDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(65855), fieldName: 'insectSwarm'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(3043), fieldName: 'scorpidSting'}),
], 'Miss');

export const PhysicalDamageDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(29859), fieldName: 'bloodFrenzy'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(58413), fieldName: 'savageCombat'}),
], 'Phys Vuln');

export const SpellCritDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(17803), fieldName: 'shadowMastery'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(12873), fieldName: 'improvedScorch'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(28593), fieldName: 'wintersChill'}),
], 'Spell Crit');

export const SpellHitDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(33198), fieldName: 'misery'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(33602), fieldName: 'faerieFire', value: TristateEffect.TristateEffectImproved}),
], 'Spell Hit');

export const SpellDamageDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(51161), fieldName: 'ebonPlaguebringer'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(48511), fieldName: 'earthAndMoon'}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(47865), fieldName: 'curseOfElements'}),
], 'Spell Dmg');

export const HuntersMark = withLabel(
  makeQuadstateDebuffInput({actionId: ActionId.fromSpellId(53338), impId: ActionId.fromSpellId(19423), impId2: ActionId.fromItemId(42907), fieldName: 'huntersMark'}),
  'Mark',
);
export const JudgementOfWisdom = withLabel(makeBooleanDebuffInput({actionId: ActionId.fromSpellId(53408), fieldName: 'judgementOfWisdom'}), 'JoW');
export const JudgementOfLight = makeBooleanDebuffInput({actionId: ActionId.fromSpellId(20271), fieldName: 'judgementOfLight'});
export const ShatteringThrow = makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(64382), numStates: 20, fieldName: 'shatteringThrows'});
export const GiftOfArthas = makeBooleanDebuffInput({actionId: ActionId.fromSpellId(11374), fieldName: 'giftOfArthas'});
export const CrystalYield = makeBooleanDebuffInput({actionId: ActionId.fromSpellId(15235), fieldName: 'crystalYield'});

///////////////////////////////////////////////////////////////////////////
//                                 CONFIGS
///////////////////////////////////////////////////////////////////////////

export const RAID_BUFFS_CONFIG = [
	// Standard buffs
	{
		config: AllStatsBuff,
		picker: MultiIconPicker,
		stats: [],
	},
  {
		config: AllStatsPercentBuff,
		picker: MultiIconPicker,
		stats: [],
	},
  {
    config: HealthBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatHealth],
  },
	{
		config: ArmorBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: StaminaBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStamina],
	},
	{
		config: StrengthAndAgilityBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStrength, Stat.StatAgility],
	},
	{
		config: IntellectBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatIntellect],
	},
	{
		config: SpiritBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpirit],
	},
	{
		config: AttackPowerBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: AttackPowerPercentBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: MeleeCritBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatMeleeCrit],
	},
  {
    config: MeleeHasteBuff,
    picker: MultiIconPicker,
    stats: [Stat.StatMeleeHaste],
  },
	{
		config: SpellPowerBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellPower],
	},
	{
		config: SpellCritBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellCrit],
	},
  {
    config: HastePercentBuff,
    picker: MultiIconPicker,
    stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste],
  },
  {
    config: DamagePercentBuff,
    picker: MultiIconPicker,
    stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
  },
  {
    config: DamageReductionPercentBuff,
    picker: MultiIconPicker,
    stats: [Stat.StatArmor],
  },
	{
		config: ResistanceBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatNatureResistance, Stat.StatShadowResistance, Stat.StatFrostResistance]
		},
	{
		config: DefensiveCooldownBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor]
	},
	{
		config: MP5Buff,
		picker: MultiIconPicker,
		stats: [Stat.StatMP5]
	},
	{
		config: ReplenishmentBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatMP5]
	},
  {
    config: Bloodlust,
    picker: IconPicker,
    stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste],
  },
  {
    config: SpellHasteBuff,
    picker: IconPicker,
    stats: [Stat.StatSpellHaste],
  },
  {
    config: RevitalizeBuff,
    picker: MultiIconPicker,
    stats: [],
  },
] as PickerStatOptions[]

export const RAID_BUFFS_MISC_CONFIG = [
  {
    config: StrengthOfWrynn,
    picker: IconPicker,
    stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
  },
  {
    config: HeroicPresence,
    picker: IconPicker,
    stats: [Stat.StatMeleeHit, Stat.StatSpellHit],
  },
  {
    config: BraidedEterniumChain,
    picker: IconPicker,
    stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit],
  },
  {
    config: ChainOfTheTwilightOwl,
    picker: IconPicker,
    stats: [Stat.StatSpellCrit, Stat.StatMeleeCrit],
  },
  {
    config: FocusMagic,
    picker: IconPicker,
    stats: [Stat.StatSpellCrit],
  },
  {
    config: EyeOfTheNight,
    picker: IconPicker,
    stats: [Stat.StatSpellPower],
  },
  {
    config: Thorns,
    picker: IconPicker,
    stats: [Stat.StatArmor],
  },
  {
    config: RetributionAura,
    picker: IconPicker,
    stats: [Stat.StatArmor],
  },
  {
    config: ManaTideTotem,
    picker: IconPicker,
    stats: [Stat.StatMP5],
  },
  {
    config: Innervate,
    picker: IconPicker,
    stats: [Stat.StatMP5],
  },
  {
    config: PowerInfusion,
    picker: IconPicker,
    stats: [Stat.StatMP5, Stat.StatSpellPower],
  },
  {
    config: TricksOfTheTrade,
    picker: IconPicker,
    stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
  },
  {
    config: UnholyFrenzy,
    picker: IconPicker,
    stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
  },
] as IconPickerStatOption[]

export const DEBUFFS_CONFIG = [
	{
		config: MajorArmorDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmorPenetration]
	},
  {
		config: MinorArmorDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmorPenetration]
	},
  {
		config: PhysicalDamageDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower]
	},
  {
		config: BleedDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower]
	},
  {
		config: SpellDamageDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellPower]
	},
  {
		config: SpellHitDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellHit]
	},
  {
		config: SpellCritDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellCrit]
	},
  {
		config: CritDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit]
	},
  {
		config: AttackPowerDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor]
	},
  {
		config: MeleeAttackSpeedDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor]
	},
  {
		config: MeleeHitDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatDodge]
	},
  {
		config: JudgementOfWisdom,
		picker: IconPicker,
		stats: [Stat.StatMP5, Stat.StatIntellect]
	},
	{
		config: HuntersMark,
		picker: IconPicker,
		stats: [Stat.StatRangedAttackPower]
	},
] as PickerStatOptions[];

export const DEBUFFS_MISC_CONFIG = [
  {
    config: JudgementOfLight,
    picker: IconPicker,
    stats: [Stat.StatStamina],
  },
  {
    config: ShatteringThrow,
    picker: IconPicker,
    stats: [Stat.StatArmorPenetration],
  },
  {
    config: GiftOfArthas,
    picker: IconPicker,
    stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
  },
  {
    config: CrystalYield,
    picker: IconPicker,
    stats: [Stat.StatArmorPenetration],
  },
] as IconPickerStatOption[];
