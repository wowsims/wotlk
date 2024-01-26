import { Faction, SaygesFortune, Stat } from "../../proto/common";
import { ActionId } from "../../proto_utils/action_id";

import { IconEnumPicker, IconEnumPickerDirection } from "../icon_enum_picker";
import {
  makeBooleanDebuffInput,
  makeBooleanIndividualBuffInput,
  makeBooleanRaidBuffInput,
  makeEnumIndividualBuffInput,
  makeMultistateIndividualBuffInput,
  makeMultistateRaidBuffInput,
  makeTristateDebuffInput,
  makeTristateIndividualBuffInput,
  makeTristateRaidBuffInput,
  withLabel
} from "../icon_inputs";
import { IconPicker } from "../icon_picker";
import { MultiIconPicker } from "../multi_icon_picker";

import { PickerStatOptions } from "./stat_options";

import * as InputHelpers from '../input_helpers';

///////////////////////////////////////////////////////////////////////////
//                                 RAID BUFFS
///////////////////////////////////////////////////////////////////////////

// TODO: Classic buff icon by level
export const AllStatsBuff = withLabel(
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(9885), impId: ActionId.fromSpellId(17055), fieldName: 'giftOfTheWild'}),
	'Mark of the Wild',
);

// Separate Strength buffs allow us to use a boolean pickers for Horde specifically
export const AllStatsPercentBuffAlliance = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(20217), fieldName: 'blessingOfKings', faction: Faction.Alliance}),
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(409580), fieldName: 'aspectOfTheLion', faction: Faction.Alliance}),
], 'Stats %');

export const AllStatsPercentBuffHorde = withLabel(
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(409580), fieldName: 'aspectOfTheLion', faction: Faction.Horde}),
	'Stats %',
);

// TODO: Classic armor buff ranks
export const ArmorBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(10293), impId: ActionId.fromSpellId(20142), fieldName: 'devotionAura'}),
	makeBooleanRaidBuffInput({id: ActionId.fromItemId(1478), fieldName: 'scrollOfProtection'}),
], 'Armor');

export const StaminaBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(10938), impId: ActionId.fromSpellId(14767), fieldName: 'powerWordFortitude'}),
	// makeTristateRaidBuffInput({id: ActionId.fromSpellId(10937), impId: ActionId.fromSpellId(14767), fieldName: 'powerWordFortitude', minLevel: 48, maxLevel: 59}),
	// makeTristateRaidBuffInput({id: ActionId.fromSpellId(2791), impId: ActionId.fromSpellId(14767), fieldName: 'powerWordFortitude', minLevel: 36, maxLevel: 47}),
	// makeTristateRaidBuffInput({id: ActionId.fromSpellId(1245), impId: ActionId.fromSpellId(14767), fieldName: 'powerWordFortitude', minLevel: 24, maxLevel: 35}),
	// makeTristateRaidBuffInput({id: ActionId.fromSpellId(1244), impId: ActionId.fromSpellId(14767), fieldName: 'powerWordFortitude', minLevel: 12, maxLevel: 23}),
	// makeTristateRaidBuffInput({id: ActionId.fromSpellId(1243), impId: ActionId.fromSpellId(14767), fieldName: 'powerWordFortitude', minLevel: 1, maxLevel: 11}),
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(11767), impId: ActionId.fromSpellId(18696), fieldName: 'bloodPact'}),
	makeBooleanRaidBuffInput({id: ActionId.fromItemId(10307), fieldName: 'scrollOfStamina'}),
], 'Stamina');

// Separate Strength buffs allow us to use boolean pickers for each
export const StrengthBuffAlliance = withLabel(
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(425600), fieldName: 'hornOfLordaeron', faction: Faction.Alliance}),
	'Strength',
)

export const StrengthBuffHorde = withLabel(
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(25361), impId: ActionId.fromSpellId(16295), fieldName: 'strengthOfEarthTotem', faction: Faction.Horde}),
	'Strength',
);

export const AgilityBuff = withLabel(
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(25359), impId: ActionId.fromSpellId(16295), fieldName: 'graceOfAirTotem', minLevel: 42, faction: Faction.Horde}),
	'Agility',
);

export const IntellectBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(23028), fieldName: 'arcaneBrilliance'}),
	makeBooleanRaidBuffInput({id: ActionId.fromItemId(10308), fieldName: 'scrollOfIntellect'}),
], 'Intellect');

export const SpiritBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(27841), fieldName: 'divineSpirit', minLevel: 30}),
	makeBooleanRaidBuffInput({id: ActionId.fromItemId(10306), fieldName: 'scrollOfSpirit'}),
], 'Spirit');

export const BlessingOfMightBuff = withLabel(
	makeTristateIndividualBuffInput({id: ActionId.fromSpellId(25291), impId: ActionId.fromSpellId(20048), fieldName: 'blessingOfMight', faction: Faction.Alliance}),
	'Blessing of Might',
);

export const BattleShoutBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(25289), impId: ActionId.fromSpellId(12861), fieldName: 'battleShout'}),
], 'Battle Shout');

export const TrueshotAuraBuff = withLabel(
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(20906), fieldName: 'trueshotAura', minLevel: 40}),
	'Trueshot Aura',
);

// export const AttackPowerPercentBuff = InputHelpers.makeMultiIconInput([
// ], 'Attack Power %', 1, 40);

export const DamageReductionPercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(25899), fieldName: 'blessingOfSanctuary'}),
], 'Mit %');

export const ResistanceBuff = InputHelpers.makeMultiIconInput([
	// Shadow
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(10958), fieldName: 'shadowProtection'}),
	// Nature
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(10601), fieldName: 'natureResistanceTotem', faction: Faction.Horde}),
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(20190), fieldName: 'aspectOfTheWild'}),
	// Frost
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(19898), fieldName: 'frostResistanceAura', faction: Faction.Alliance}),
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(10479), fieldName: 'frostResistanceTotem', faction: Faction.Horde}),
], 'Resistances');

export const BlessingOfWisdom = withLabel(
	makeTristateIndividualBuffInput({id: ActionId.fromSpellId(25290), impId: ActionId.fromSpellId(20245), fieldName: 'blessingOfWisdom', faction: Faction.Alliance}),
	'Blessing of Wisdom',
);

export const ManaSpringTotem = withLabel(
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(10497), impId: ActionId.fromSpellId(16208), fieldName: 'manaSpringTotem', minLevel: 40, faction: Faction.Horde}),
	'Mana Spring Totem',
);

export const MeleeCritBuff = withLabel(
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(17007), fieldName: 'leaderOfThePack', minLevel: 40}),
	'Leader of the Pack',
);

export const SpellCritBuff = withLabel(
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(24907), fieldName: 'moonkinAura', minLevel: 40}),
	'Moonkin Aura',
);

export const SpellIncreaseBuff = withLabel(
	makeMultistateRaidBuffInput({id: ActionId.fromSpellId(425464), numStates: 200, fieldName: 'demonicPact', multiplier: 10}),
	'Demonic Pact',
);

export const DefensiveCooldownBuff = InputHelpers.makeMultiIconInput([
], 'Defensive CDs');

// Misc Buffs
export const RetributionAura = withLabel(
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(10301), fieldName: 'retributionAura'}),
	'Retribution Aura',
);
export const Thorns = withLabel(
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(9910), impId: ActionId.fromSpellId(16840), fieldName: 'thorns'}),
	'Thorns',
);
export const Innervate = withLabel(
	makeMultistateIndividualBuffInput({id: ActionId.fromSpellId(29166), numStates: 11, fieldName: 'innervates', minLevel: 40}),
	'Innervate',
);
export const PowerInfusion = withLabel(
	makeMultistateIndividualBuffInput({id: ActionId.fromSpellId(10060), numStates: 11, fieldName: 'powerInfusions', minLevel: 40}),
	'Power Infusion',
);

///////////////////////////////////////////////////////////////////////////
//                                 WORLD BUFFS
///////////////////////////////////////////////////////////////////////////

export const RallyingCryOfTheDragonslayer = withLabel(
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(22888), fieldName: 'rallyingCryOfTheDragonslayer'}),
	'Rallying Cry of the Dragonslayer',
);

export const SpiritOfZandalar = withLabel(
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(24425), fieldName: 'spiritOfZandalar'}),
	'Spirit of Zandalar',
);
export const SongflowerSerenade = withLabel(
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(15366), fieldName: 'songflowerSerenade'}),
	'Songflower Serenade',
);
export const WarchiefsBlessing = withLabel(
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(16609), fieldName: 'warchiefsBlessing'}),
	`Warchief's Blessing`,
);

export const SaygesDarkFortune = makeEnumIndividualBuffInput({
	numColumns: 6,
	direction: IconEnumPickerDirection.Horizontal,
	values: [
		{ iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_orb_02.jpg', value: SaygesFortune.SaygesUnknown, text: `Sayge's Dark Fortune` },
		{ id: ActionId.fromSpellId(23768), value: SaygesFortune.SaygesDamage, text: `Sayge's Damage` },
		{ id: ActionId.fromSpellId(23736), value: SaygesFortune.SaygesAgility, text: `Sayge's Agility` },
		{ id: ActionId.fromSpellId(23766), value: SaygesFortune.SaygesIntellect, text: `Sayge's Intellect` },
		{ id: ActionId.fromSpellId(23738), value: SaygesFortune.SaygesSpirit, text: `Sayge's Spirit` },
		{ id: ActionId.fromSpellId(23737), value: SaygesFortune.SaygesStamina, text: `Sayge's Stamina` },
	],
	fieldName: 'saygesFortune'
})

// Dire Maul Buffs
export const FengusFerocity = withLabel(
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(22817), fieldName: 'fengusFerocity'}),
	`Fengus' Ferocity`,
);
export const MoldarsMoxie = withLabel(
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(22818), fieldName: 'moldarsMoxie'}),
	`Moldar's Moxie`,
);
export const SlipKiksSavvy = withLabel(
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(22820), fieldName: 'slipkiksSavvy'}),
	`Slip'kik's Savvy`,
);

// SoD World Buffs
export const BoonOfBlackfathom = withLabel(
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(430947), fieldName: 'boonOfBlackfathom'}),
	'Boon of Blackfathom',
);

///////////////////////////////////////////////////////////////////////////
//                                 DEBUFFS
///////////////////////////////////////////////////////////////////////////

export const MajorArmorDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({id: ActionId.fromSpellId(11597), fieldName: 'sunderArmor'}),
	makeTristateDebuffInput(ActionId.fromSpellId(11198), ActionId.fromSpellId(14169), 'exposeArmor'),
], 'Armor Penetration');

export const CurseOfRecklessness = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(11717), fieldName: 'curseOfRecklessness'}),
	'Curse of Recklessness'
);

export const FaerieFire = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(9907), fieldName: 'faerieFire'})
	, 'Faerie Fire'
);

// TODO: Classic
export const MinorArmorDebuff = InputHelpers.makeMultiIconInput([
	//makeTristateDebuffInput(ActionId.fromSpellId(770), ActionId.fromSpellId(33602), 'faerieFire'),
	//makeBooleanDebuffInput({id: ActionId.fromSpellId(50511), fieldName: 'curseOfWeakness'}),
], 'Minor ArP');

export const AttackPowerDebuff = InputHelpers.makeMultiIconInput([
	makeTristateDebuffInput(ActionId.fromSpellId(11556), ActionId.fromSpellId(12879), 'demoralizingShout'),
	makeTristateDebuffInput(ActionId.fromSpellId(9898), ActionId.fromSpellId(16862), 'demoralizingRoar'),
], 'Attack Power');

// TODO: Classic
export const BleedDebuff = InputHelpers.makeMultiIconInput([
	// makeBooleanDebuffInput(ActionId.fromSpellId(48564), 'mangle'),
], 'Bleed');

export const MeleeAttackSpeedDebuff = withLabel(
	makeTristateDebuffInput(ActionId.fromSpellId(6343), ActionId.fromSpellId(12666), 'thunderClap'),
	'Thunder Clap',
);

export const MeleeHitDebuff = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(24977), fieldName: 'insectSwarm'}),
	'Insect Swarm',
);

// TODO: Classic
export const SpellISBDebuff = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(17803), fieldName: 'improvedShadowBolt'}),
	'Improved Shadow Bolt',
);

export const SpellScorchDebuff = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(12873), fieldName: 'improvedScorch', minLevel: 40}),
	'Improved Scorch',
);

export const SpellWintersChillDebuff = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(28595), fieldName: 'wintersChill', minLevel: 40}),
	'Winters Chill',
);

// TODO: Classic
// export const SpellDamageDebuff = InputHelpers.makeMultiIconInput([
// 	makeBooleanDebuffInput(ActionId.fromSpellId(47865), 'curseOfElements'),
// ], 'Spell Dmg');

// TODO: Classic
export const HuntersMark = withLabel(
	makeTristateDebuffInput(ActionId.fromSpellId(14325), ActionId.fromSpellId(19425), 'huntersMark'),
	`Hunter's Mark`,
);
export const JudgementOfWisdom = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(20355), fieldName: 'judgementOfWisdom', minLevel: 38}),
	'Judgement of Wisdom',
);
export const JudgementOfLight = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(20346), fieldName: 'judgementOfLight', minLevel: 30}),
	'Judgement of Light',
);

///////////////////////////////////////////////////////////////////////////
//                                 CONFIGS
///////////////////////////////////////////////////////////////////////////

export const RAID_BUFFS_CONFIG = [
	// Standard buffs
	{
		config: AllStatsBuff,
		picker: IconPicker,
		stats: []
	},
	{
		config: AllStatsPercentBuffAlliance,
		picker: MultiIconPicker,
		stats: []
	},
	{
		config: AllStatsPercentBuffHorde,
		picker: IconPicker,
		stats: []
	},
	{
		config: ArmorBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor]
	},
	{
		config: StaminaBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStamina]
	},
	{
		config: StrengthBuffAlliance,
		picker: IconPicker,
		stats: [Stat.StatStrength]
	},
	{
		config: StrengthBuffHorde,
		picker: IconPicker,
		stats: [Stat.StatStrength]
	},
	{
		config: AgilityBuff,
		picker: IconPicker,
		stats: [Stat.StatAgility]
	},
	{
		config: IntellectBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatIntellect]
	},
	{
		config: SpiritBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpirit]
	},
	{
		config: BattleShoutBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower]
	},
	{
		config: BlessingOfMightBuff,
		picker: IconPicker,
		stats: [Stat.StatAttackPower]
	},
	{
		config: TrueshotAuraBuff,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower]
		},
	// {
	// 	config: AttackPowerPercentBuff,
	// 	picker: MultiIconPicker,
	// 	stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower]
	// },
	{
		config: MeleeCritBuff,
		picker: IconPicker,
		stats: [Stat.StatMeleeCrit]
	},
	{
		config: SpellIncreaseBuff,
		picker: IconPicker,
		stats: [Stat.StatSpellPower]
	},
	{
		config: SpellCritBuff,
		picker: IconPicker,
		stats: [Stat.StatSpellCrit]
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
		config: BlessingOfWisdom,
		picker: IconPicker,
		stats: [Stat.StatMP5]
	},
	{
		config: ManaSpringTotem,
		picker: IconPicker,
		stats: [Stat.StatMP5]
	},

	// Misc Buffs
	{
		config: Thorns,
		picker: IconPicker,
		stats: [Stat.StatArmor]
	},
	{
		config: RetributionAura,
		picker: IconPicker,
		stats: [Stat.StatArmor]
	},
	{
		config: Innervate,
		picker: IconPicker,
		stats: [Stat.StatMP5]
	},
	{
		config: PowerInfusion,
		picker: IconPicker,
		stats: [Stat.StatMP5, Stat.StatSpellPower]
	},
] as PickerStatOptions[]

export const WORLD_BUFFS_CONFIG = [
	{
		config: BoonOfBlackfathom,
		picker: IconPicker,
		stats: [
			Stat.StatMeleeCrit,
			// TODO: Stat.StatRangedCrit,
			Stat.StatSpellCrit,
			Stat.StatAttackPower
		]
	},
	{
		config: FengusFerocity,
		picker: IconPicker,
		stats: [Stat.StatAttackPower]
	},
	{
		config: MoldarsMoxie,
		picker: IconPicker,
		stats: [Stat.StatStamina]
	},
	{
		config: RallyingCryOfTheDragonslayer,
		picker: IconPicker,
		stats: [
			Stat.StatMeleeCrit,
			// TODO: Stat.StatRangedCrit,
			Stat.StatSpellCrit,
			Stat.StatAttackPower,
		]
	},
	{
		config: SaygesDarkFortune,
		picker: IconEnumPicker,
		stats: [],
	},
	{
		config: SongflowerSerenade,
		picker: IconPicker,
		stats: []
	},
	{
		config: SpiritOfZandalar,
		picker: IconPicker,
		stats: []
	},
	{
		config: WarchiefsBlessing,
		picker: IconPicker,
		stats: [
			Stat.StatHealth,
			Stat.StatMeleeHaste,
			Stat.StatMP5,
		]
	},
] as PickerStatOptions[];

export const DEBUFFS_CONFIG = [
	// Standard Debuffs
	{ 
		config: MajorArmorDebuff,
		stats: [Stat.StatArmorPenetration],
		picker: MultiIconPicker,
	},
	{ 
		config: CurseOfRecklessness,
		picker: IconPicker,
		stats: [Stat.StatAttackPower]
	},
	{ 
		config: FaerieFire,
		picker: IconPicker,
		stats: [Stat.StatAttackPower]
	},
	// { 
	// 	config: MinorArmorDebuff,
	// picker: MultiIconPicker,
	// 	stats: [Stat.StatArmorPenetration]
	// },
	{ 
		config: BleedDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower]
	},
	{ 
		config: SpellISBDebuff,
		picker: IconPicker,
		stats: [Stat.StatShadowPower]
	},
	{ 
		config: SpellScorchDebuff,
		picker: IconPicker,
		stats: [Stat.StatFirePower]
	},
	{ 
		config: SpellWintersChillDebuff,
		picker: IconPicker,
		stats: [Stat.StatFrostPower]
	},
	{ 
		config: AttackPowerDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor]
	},
	{ 
		config: MeleeAttackSpeedDebuff,
		picker: IconPicker,
		stats: [Stat.StatArmor]
	},
	{ 
		config: MeleeHitDebuff,
		picker: IconPicker,
		stats: [Stat.StatDodge]
	},

	// Other Debuffs
	{
		config: JudgementOfWisdom,
		picker: IconPicker,
		stats: [Stat.StatMP5, Stat.StatIntellect],
	},
	{
		config: HuntersMark,
		picker: IconPicker,
		stats: [Stat.StatRangedAttackPower],
	},

	// Misc Debuffs
	{
		config: JudgementOfLight,
		picker: IconPicker,
		stats: [Stat.StatStamina]
	},
] as PickerStatOptions[];
