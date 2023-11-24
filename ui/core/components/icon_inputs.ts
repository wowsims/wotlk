import { AgilityElixir, Consumes, Debuffs, Explosive, Flask, Food, IndividualBuffs, RaidBuffs, StrengthBuff, WeaponBuff } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';

import { Player } from '../player.js';
import { Spec } from '../proto/common.js';
import { Raid } from '../raid.js';
import { EventID, TypedEvent } from '../typed_event.js';

import { IconEnumPicker, IconEnumValueConfig } from './icon_enum_picker.js';
import { IconPicker } from './icon_picker.js';

import * as InputHelpers from './input_helpers.js';
import { MAX_CHARACTER_LEVEL } from '../constants/mechanics.js';

// Component Functions

export type IconInputConfig<ModObject, T> = (
	InputHelpers.TypedIconPickerConfig<ModObject, T> |
	InputHelpers.TypedIconEnumPickerConfig<ModObject, T>
);

export const buildIconInput = (parent: HTMLElement, player: Player<Spec>, inputConfig: IconInputConfig<Player<Spec>, any>) => {
	if (inputConfig.type == 'icon') {
		return new IconPicker<Player<Spec>, any>(parent, player, inputConfig);
	} else if (inputConfig.type == 'iconEnum') {
		return new IconEnumPicker<Player<Spec>, any>(parent, player, inputConfig);
	}
};

// Raid Buffs

// TODO: Classic buff icon by level
export const AllStatsBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(9885), impId: ActionId.fromSpellId(17055), fieldName: 'giftOfTheWild'}),
], 'Stats');

export const AllStatsPercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(20217), fieldName: 'blessingOfKings'}),
], 'Stats %');

// TODO: Classic armor buff ranks
export const ArmorBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(10293), impId: ActionId.fromSpellId(20142), fieldName: 'devotionAura'}),
	makeBooleanRaidBuffInput({id: ActionId.fromItemId(43468), fieldName: 'scrollOfProtection'}),
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

// TODO: Breakout Strength / Agi
export const StrengthAndAgilityBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(25361), impId: ActionId.fromSpellId(52456), fieldName: 'strengthOfEarthTotem'}),
	makeBooleanRaidBuffInput({id: ActionId.fromItemId(10309), fieldName: 'scrollOfAgility'}),
	makeBooleanRaidBuffInput({id: ActionId.fromItemId(10310), fieldName: 'scrollOfStrength'}),
], 'Str/Agi');

export const IntellectBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(23028), fieldName: 'arcaneBrilliance'}),
	makeBooleanRaidBuffInput({id: ActionId.fromItemId(10308), fieldName: 'scrollOfIntellect'}),
], 'Int');

export const SpiritBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(27841), fieldName: 'divineSpirit'}),
	makeBooleanRaidBuffInput({id: ActionId.fromItemId(10306), fieldName: 'scrollOfSpirit'}),
], 'Spirit');

export const AttackPowerBuff = InputHelpers.makeMultiIconInput([
	makeTristateIndividualBuffInput(ActionId.fromSpellId(48934), ActionId.fromSpellId(20045), 'blessingOfMight'),
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(47436), impId: ActionId.fromSpellId(12861), fieldName: 'battleShout'}),
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(19506), fieldName: 'trueshotAura'}),
], 'AP');

export const AttackPowerPercentBuff = InputHelpers.makeMultiIconInput([
], 'Atk Pwr %');

export const DamageReductionPercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput({id: ActionId.fromSpellId(25899), fieldName: 'blessingOfSanctuary'}),
], 'Mit %');

export const ResistanceBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(48170), fieldName: 'shadowProtection'}),
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(58749), fieldName: 'natureResistanceTotem'}),
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(49071), fieldName: 'aspectOfTheWild'}),
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(48945), fieldName: 'frostResistanceAura'}),
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(58745), fieldName: 'frostResistanceTotem'}),
], 'Resistances');

export const MP5Buff = InputHelpers.makeMultiIconInput([
	makeTristateIndividualBuffInput(ActionId.fromSpellId(25290), ActionId.fromSpellId(20245), 'blessingOfWisdom'),
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(10497), impId: ActionId.fromSpellId(16208), fieldName: 'manaSpringTotem', minLevel: 40}),
], 'MP5');

export const MeleeCritBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({id: ActionId.fromSpellId(17007), impId: ActionId.fromSpellId(34300), fieldName: 'leaderOfThePack'}),
], 'Melee Crit');

export const SpellCritBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({id: ActionId.fromSpellId(24907), fieldName: 'moonkinAura'}),
], 'Spell Crit');

// TODO: Classic rune
export const SpellIncreaseBuff = InputHelpers.makeMultiIconInput([
	// makeMultistateRaidBuffInput(ActionId.fromSpellId(47240), 2000, 'demonicPactSp', 20),
], 'Spell Power');

export const DefensiveCooldownBuff = InputHelpers.makeMultiIconInput([
], 'Defensive CDs');

// Misc Buffs
export const RetributionAura = makeBooleanRaidBuffInput({id: ActionId.fromSpellId(10301), fieldName: 'retributionAura'});
export const Thorns = makeTristateRaidBuffInput({id: ActionId.fromSpellId(9910), impId: ActionId.fromSpellId(16840), fieldName: 'thorns'});
export const Innervate = makeMultistateIndividualBuffInput(ActionId.fromSpellId(29166), 11, 'innervates');
export const PowerInfusion = makeMultistateIndividualBuffInput(ActionId.fromSpellId(10060), 11, 'powerInfusions');

// Debuffs
export const MajorArmorDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({id: ActionId.fromSpellId(7386), fieldName: 'sunderArmor'}),
	makeBooleanDebuffInput({id: ActionId.fromSpellId(8647), fieldName: 'exposeArmor'}),
], 'Major ArP');

// TODO: Classic
export const MinorArmorDebuff = InputHelpers.makeMultiIconInput([
	// makeTristateDebuffInput(ActionId.fromSpellId(770), ActionId.fromSpellId(33602), 'faerieFire'),
	// makeTristateDebuffInput(ActionId.fromSpellId(50511), ActionId.fromSpellId(18180), 'curseOfWeakness'),
], 'Minor ArP');

export const AttackPowerDebuff = InputHelpers.makeMultiIconInput([
	makeTristateDebuffInput(ActionId.fromSpellId(47437), ActionId.fromSpellId(12879), 'demoralizingShout'),
	makeTristateDebuffInput(ActionId.fromSpellId(48560), ActionId.fromSpellId(16862), 'demoralizingRoar'),
], 'Atk Pwr');

// TODO: Classic
export const BleedDebuff = InputHelpers.makeMultiIconInput([
	// makeBooleanDebuffInput(ActionId.fromSpellId(48564), 'mangle'),
], 'Bleed');

export const MeleeAttackSpeedDebuff = InputHelpers.makeMultiIconInput([
	makeTristateDebuffInput(ActionId.fromSpellId(47502), ActionId.fromSpellId(12666), 'thunderClap'),
], 'Atk Speed');

export const MeleeHitDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({id: ActionId.fromSpellId(65855), fieldName: 'insectSwarm'}),
], 'Miss');

// TODO: Classic
export const SpellISBDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({id: ActionId.fromSpellId(17803), fieldName: 'improvedShadowBolt'}),
], 'ISB');

export const SpellScorchDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({id: ActionId.fromSpellId(12873), fieldName: 'improvedScorch'}),
], 'Scorch');

export const SpellWintersChillDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({id: ActionId.fromSpellId(28595), fieldName: 'wintersChill'}),
], 'Winters Chill');

// TODO: Classic
// export const SpellDamageDebuff = InputHelpers.makeMultiIconInput([
// 	makeBooleanDebuffInput(ActionId.fromSpellId(47865), 'curseOfElements'),
// ], 'Spell Dmg');

// TODO: Classic
export const HuntersMark = withLabel(makeTristateDebuffInput(ActionId.fromSpellId(14325), ActionId.fromSpellId(19425), 'huntersMark'), 'Mark');
export const JudgementOfWisdom = withLabel(makeBooleanDebuffInput({id: ActionId.fromSpellId(20355), fieldName: 'judgementOfWisdom'}), 'JoW');
export const JudgementOfLight = makeBooleanDebuffInput({id: ActionId.fromSpellId(20346), fieldName: 'judgementOfLight'});
export const GiftOfArthas = makeBooleanDebuffInput({id: ActionId.fromSpellId(11374), fieldName: 'giftOfArthas'});
export const CrystalYield = makeBooleanDebuffInput({id: ActionId.fromSpellId(15235), fieldName: 'crystalYield'});

// Consumes
export const Sapper = makeBooleanConsumeInput({id: ActionId.fromItemId(10646), fieldName: 'sapper', minLevel: 40});

// TODO: Classic
// export const PetScrollOfAgilityV = makeBooleanConsumeInput(ActionId.fromItemId(27498), 'petScrollOfAgility', 5);
// export const PetScrollOfStrengthV = makeBooleanConsumeInput(ActionId.fromItemId(27503), 'petScrollOfStrength', 5);

// eslint-disable-next-line unused-imports/no-unused-vars
function withLabel<ModObject, T>(config: InputHelpers.TypedIconPickerConfig<ModObject, T>, label: string): InputHelpers.TypedIconPickerConfig<ModObject, T> {
	config.label = label;
	return config;
}

interface BooleanInputConfig<T> {
	id: ActionId, 
	fieldName: keyof T, 
	value?: number, 
	minLevel?: number,
	maxLevel?: number,
}

function makeBooleanRaidBuffInput<SpecType extends Spec>(config: BooleanInputConfig<RaidBuffs>): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, RaidBuffs, Player<SpecType>>({
		getModObject: (player: Player<any>) => player,
		showWhen: (p) => (config.minLevel || 0) <= p.getLevel() && p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL),
		getValue: (p) => p.getRaid()!.getBuffs(),
		setValue: (eventID: EventID, p: Player<SpecType>, newVal: RaidBuffs) => p.getRaid()!.setBuffs(eventID, newVal),
		changeEmitter: (p) => TypedEvent.onAny([p.getRaid()!.buffsChangeEmitter, p.levelChangeEmitter]),
	}, config.id, config.fieldName, config.value);
}
// function makeBooleanPartyBuffInput(id: ActionId, fieldName: keyof PartyBuffs, value?: number): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
// 	return InputHelpers.makeBooleanIconInput<any, PartyBuffs, Party>({
// 		getModObject: (player: Player<any>) => player.getParty()!,
// 		getValue: (party: Party) => party.getBuffs(),
// 		setValue: (eventID: EventID, party: Party, newVal: PartyBuffs) => party.setBuffs(eventID, newVal),
// 		changeEmitter: (party: Party) => party.buffsChangeEmitter,
// 	}, id, fieldName, value);
// }

function makeBooleanIndividualBuffInput(config: BooleanInputConfig<IndividualBuffs>): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, IndividualBuffs, Player<any>>({
		getModObject: (player: Player<any>) => player,
		showWhen: (p) => (config.minLevel || 0) <= p.getLevel() && p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL),
		getValue: (player: Player<any>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<any>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<any>) => player.buffsChangeEmitter,
	}, config.id, config.fieldName, config.value);
}

// eslint-disable-next-line unused-imports/no-unused-vars
function makeBooleanConsumeInput<SpecType extends Spec>(config: BooleanInputConfig<Consumes>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, Consumes, Player<any>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (p) => p.getLevel() >= (config.minLevel || 0),
		getValue: (player: Player<any>) => player.getConsumes(),
		setValue: (eventID: EventID, player: Player<any>, newVal: Consumes) => player.setConsumes(eventID, newVal),
		changeEmitter: (player: Player<any>) => TypedEvent.onAny([player.consumesChangeEmitter, player.levelChangeEmitter])
	}, config.id, config.fieldName, config.value);
}
function makeBooleanDebuffInput<SpecType extends Spec>(config: BooleanInputConfig<Debuffs>): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, Debuffs, Player<SpecType>>({
		getModObject: (player) => player,
		showWhen: (p) => (config.minLevel || 0) <= p.getLevel() && p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL),
		getValue: (p) => p.getRaid()!.getDebuffs(),
		setValue: (eventID: EventID, p: Player<SpecType>, newVal: Debuffs) => p.getRaid()!.setDebuffs(eventID, newVal),
		changeEmitter: (p) => TypedEvent.onAny([p.getRaid()!.debuffsChangeEmitter, p.levelChangeEmitter]),
	}, config.id, config.fieldName, config.value);
}

interface TristateInputConfig<T> {
	id: ActionId, 
	impId: ActionId, 
	fieldName: keyof T,
	minLevel?: number
	maxLevel?: number
}

function makeTristateRaidBuffInput<SpecType extends Spec>(config: TristateInputConfig<RaidBuffs>): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeTristateIconInput<any, RaidBuffs, Player<SpecType>>({
		getModObject: (player) => player,
		showWhen: (p) => (config.minLevel || 0) <= p.getLevel() && p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL),
		getValue: (p) => p.getRaid()!.getBuffs(),
		setValue: (eventID: EventID, p: Player<SpecType>, newVal: RaidBuffs) => p.getRaid()!.setBuffs(eventID, newVal),
		changeEmitter: (p: Player<SpecType>) => TypedEvent.onAny([p.getRaid()!.buffsChangeEmitter, p.levelChangeEmitter]),
	}, config.id, config.impId, config.fieldName);
}
function makeTristateIndividualBuffInput(id: ActionId, impId: ActionId, fieldName: keyof IndividualBuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeTristateIconInput<any, IndividualBuffs, Player<any>>({
		getModObject: (player: Player<any>) => player,
		getValue: (player: Player<any>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<any>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<any>) => player.buffsChangeEmitter,
	}, id, impId, fieldName);
}
function makeTristateDebuffInput(id: ActionId, impId: ActionId, fieldName: keyof Debuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeTristateIconInput<any, Debuffs, Raid>({
		getModObject: (player: Player<any>) => player.getRaid()!,
		getValue: (raid: Raid) => raid.getDebuffs(),
		setValue: (eventID: EventID, raid: Raid, newVal: Debuffs) => raid.setDebuffs(eventID, newVal),
		changeEmitter: (raid: Raid) => raid.debuffsChangeEmitter,
	}, id, impId, fieldName);
}
// function makeQuadstateDebuffInput(id: ActionId, impId: ActionId, impId2: ActionId, fieldName: keyof Debuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
// 	return InputHelpers.makeQuadstateIconInput<any, Debuffs, Raid>({
// 		getModObject: (player: Player<any>) => player.getRaid()!,
// 		getValue: (raid: Raid) => raid.getDebuffs(),
// 		setValue: (eventID: EventID, raid: Raid, newVal: Debuffs) => raid.setDebuffs(eventID, newVal),
// 		changeEmitter: (raid: Raid) => raid.debuffsChangeEmitter,
// 	}, id, impId, impId2, fieldName);
// }
// function makeMultistateRaidBuffInput(id: ActionId, numStates: number, fieldName: keyof RaidBuffs, multiplier?: number): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
// 	return InputHelpers.makeMultistateIconInput<any, RaidBuffs, Raid>({
// 		getModObject: (player: Player<any>) => player.getRaid()!,
// 		getValue: (raid: Raid) => raid.getBuffs(),
// 		setValue: (eventID: EventID, raid: Raid, newVal: RaidBuffs) => raid.setBuffs(eventID, newVal),
// 		changeEmitter: (raid: Raid) => raid.buffsChangeEmitter,
// 	}, id, numStates, fieldName, multiplier);
// }
// function makeMultistatePartyBuffInput(id: ActionId, numStates: number, fieldName: keyof PartyBuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
// 	return InputHelpers.makeMultistateIconInput<any, PartyBuffs, Party>({
// 		getModObject: (player: Player<any>) => player.getParty()!,
// 		getValue: (party: Party) => party.getBuffs(),
// 		setValue: (eventID: EventID, party: Party, newVal: PartyBuffs) => party.setBuffs(eventID, newVal),
// 		changeEmitter: (party: Party) => party.buffsChangeEmitter,
// 	}, id, numStates, fieldName);
// }
function makeMultistateIndividualBuffInput(id: ActionId, numStates: number, fieldName: keyof IndividualBuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeMultistateIconInput<any, IndividualBuffs, Player<any>>({
		getModObject: (player: Player<any>) => player,
		getValue: (player: Player<any>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<any>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<any>) => player.buffsChangeEmitter,
	}, id, numStates, fieldName);
}
// function makeMultistateMultiplierIndividualBuffInput(id: ActionId, numStates: number, multiplier: number, fieldName: keyof IndividualBuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
// 	return InputHelpers.makeMultistateIconInput<any, IndividualBuffs, Player<any>>({
// 		getModObject: (player: Player<any>) => player,
// 		getValue: (player: Player<any>) => player.getBuffs(),
// 		setValue: (eventID: EventID, player: Player<any>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
// 		changeEmitter: (player: Player<any>) => player.buffsChangeEmitter,
// 	}, id, numStates, fieldName, multiplier);
// }


//////////////////////////////////////////////////////////////////////
// Custom buffs that don't fit into any of the helper functions above.
//////////////////////////////////////////////////////////////////////

// function makePotionInputFactory(consumesFieldName: keyof Consumes): (options: Array<Potions>, tooltip?: string) => InputHelpers.TypedIconEnumPickerConfig<Player<any>, Potions> {
// 	return makeConsumeInputFactory({
// 		consumesFieldName: consumesFieldName,
// 		allOptions: [
// 			{ actionId: ActionId.fromItemId(33447), value: Potions.RunicHealingPotion },
// 			{ actionId: ActionId.fromItemId(41166), value: Potions.RunicHealingInjector },
// 			{ actionId: ActionId.fromItemId(33448), value: Potions.RunicManaPotion },
// 			{ actionId: ActionId.fromItemId(42545), value: Potions.RunicManaInjector },
// 			{ actionId: ActionId.fromItemId(40093), value: Potions.IndestructiblePotion },
// 			{ actionId: ActionId.fromItemId(40211), value: Potions.PotionOfSpeed },
// 			{ actionId: ActionId.fromItemId(40212), value: Potions.PotionOfWildMagic },

// 			{ actionId: ActionId.fromItemId(22839), value: Potions.DestructionPotion },
// 			{ actionId: ActionId.fromItemId(22838), value: Potions.HastePotion },
// 			{ actionId: ActionId.fromItemId(13442), value: Potions.MightyRagePotion },
// 			{ actionId: ActionId.fromItemId(22832), value: Potions.SuperManaPotion },
// 			{ actionId: ActionId.fromItemId(31677), value: Potions.FelManaPotion },
// 			{ actionId: ActionId.fromItemId(22828), value: Potions.InsaneStrengthPotion },
// 			{ actionId: ActionId.fromItemId(22849), value: Potions.IronshieldPotion },
// 			{ actionId: ActionId.fromItemId(22837), value: Potions.HeroicPotion },
// 		] as Array<IconEnumValueConfig<Player<any>, Potions>>,
// 	});
// }

// TODO: Classic? 
// export const makeConjuredInput = makeConsumeInputFactory({
// 	consumesFieldName: 'defaultConjured',
// 	allOptions: [
// 		{ actionId: ActionId.fromItemId(12662), value: Conjured.ConjuredDarkRune },
// 		{ actionId: ActionId.fromItemId(22105), value: Conjured.ConjuredHealthstone },
// 		{ actionId: ActionId.fromItemId(7676), value: Conjured.ConjuredRogueThistleTea },
// 	] as Array<IconEnumValueConfig<Player<any>, Conjured>>
// });

export const makeFlasksInput = makeConsumeInputFactory({
	consumesFieldName: 'flask',
	allOptions: [
		{ actionId: ActionId.fromItemId(13510), value: Flask.FlaskOfTheTitans },
		{ actionId: ActionId.fromItemId(13511), value: Flask.FlaskOfDistilledWisdom },
		{ actionId: ActionId.fromItemId(13512), value: Flask.FlaskOfSupremePower },
		{ actionId: ActionId.fromItemId(13513), value: Flask.FlaskOfChromaticResistance },
	] as Array<IconEnumValueConfig<Player<any>, Flask>>,
});

export const makeWeaponBuffsInput = makeConsumeInputFactory({
	consumesFieldName: 'weaponBuff',
	allOptions: [
		{ actionId: ActionId.fromItemId(20749), value: WeaponBuff.BrillianWizardOil, showWhen: (p) =>  p.getLevel() >= 45},
		{ actionId: ActionId.fromItemId(20748), value: WeaponBuff.BrilliantManaOil, showWhen: (p) =>  p.getLevel() >= 45 },
		{ actionId: ActionId.fromItemId(12404), value: WeaponBuff.DenseSharpeningStone, showWhen: (p) =>  p.getLevel() >= 35 },
		{ actionId: ActionId.fromItemId(18262), value: WeaponBuff.ElementalSharpeningStone, showWhen: (p) =>  p.getLevel() >= 50 },
	] as Array<IconEnumValueConfig<Player<any>, WeaponBuff>>,
});

export const makeFoodInput = makeConsumeInputFactory({
	consumesFieldName: 'food',
	allOptions: [
		{ actionId: ActionId.fromItemId(15856), value: Food.FoodHotWolfRibs, showWhen: (p) => p.getLevel() >= 25 },
		{ actionId: ActionId.fromItemId(22480), value: Food.FoodTenderWolfSteak, showWhen: (p) => p.getLevel() >= 40 },
		{ actionId: ActionId.fromItemId(13931), value: Food.FoodNightfinSoup, showWhen: (p) => p.getLevel() >= 35 },
		{ actionId: ActionId.fromItemId(13931), value: Food.FoodNightfinSoup, showWhen: (p) => p.getLevel() >= 35 },
		{ actionId: ActionId.fromItemId(13928), value: Food.FoodGrilledSquid, showWhen: (p) => p.getLevel() >= 35 },
		{ actionId: ActionId.fromItemId(20452), value: Food.FoodSmokedDesertDumpling, showWhen: (p) => p.getLevel() >= 45 },
		{ actionId: ActionId.fromItemId(18254), value: Food.FoodRunnTumTuberSurprise, showWhen: (p) => p.getLevel() >= 45 },
		{ actionId: ActionId.fromItemId(13813), value: Food.FoodBlessedSunfruitJuice, showWhen: (p) => p.getLevel() >= 45 },
		{ actionId: ActionId.fromItemId(13810), value: Food.FoodBlessSunfruit, showWhen: (p) => p.getLevel() >= 45 },
		{ actionId: ActionId.fromItemId(21023), value: Food.FoodDirgesKickChimaerokChops, showWhen: (p) => p.getLevel() >= 55 },
	] as Array<IconEnumValueConfig<Player<any>, Food>>
});

export const AgilityBuffInput = makeConsumeInput('agilityElixir', [
	{ actionId: ActionId.fromItemId(13452), value: AgilityElixir.ElixirOfTheMongoose, showWhen: (p) => p.getLevel() >= 46 },
	{ actionId: ActionId.fromItemId(9187), value: AgilityElixir.ElixirOfGreaterAgility, showWhen: (p) => p.getLevel() >= 38},
] as Array<IconEnumValueConfig<Player<any>, AgilityElixir>>);

export const StrengthBuffInput = makeConsumeInput('strengthBuff', [
	{ actionId: ActionId.fromItemId(12451), value: StrengthBuff.JujuPower },
	{ actionId: ActionId.fromItemId(9206), value: StrengthBuff.ElixirOfGiants },
] as Array<IconEnumValueConfig<Player<any>, StrengthBuff>>);

export const SpellDamageBuff = makeBooleanConsumeInput({id: ActionId.fromItemId(13454), fieldName: 'spellPowerBuff', minLevel: 47});
export const ShadowDamageBuff = makeBooleanConsumeInput({id: ActionId.fromItemId(9264), fieldName: 'shadowPowerBuff', minLevel: 40});
export const FireDamageBuff = makeBooleanConsumeInput({id: ActionId.fromItemId(21546), fieldName: 'firePowerBuff', minLevel: 40});
export const FrostDamageBuff = makeBooleanConsumeInput({id: ActionId.fromItemId(17708), fieldName: 'frostPowerBuff', minLevel: 40});

export const FillerExplosiveInput = makeConsumeInput('fillerExplosive', [
	{ actionId: ActionId.fromItemId(18641), value: Explosive.ExplosiveDenseDynamite, showWhen: (p) => p.getLevel() >= 40 },
	{ actionId: ActionId.fromItemId(15993), value: Explosive.ExplosiveThoriumGrenade, showWhen: (p) => p.getLevel() >= 40 },
] as Array<IconEnumValueConfig<Player<any>, Explosive>>);

export interface ConsumeInputFactoryArgs<T extends number> {
	consumesFieldName: keyof Consumes,
	allOptions: Array<IconEnumValueConfig<Player<any>, T>>,
	// Additional callback if logic besides syncing consumes is required
	onSet?: (eventID: EventID, player: Player<any>, newValue: T) => void
}
function makeConsumeInputFactory<T extends number>(args: ConsumeInputFactoryArgs<T>): (options: Array<T>, tooltip?: string) => InputHelpers.TypedIconEnumPickerConfig<Player<any>, T> {
	return (options: Array<T>, tooltip?: string) => {
		return {
			type: 'iconEnum',
			tooltip: tooltip,
			numColumns: options.length > 5 ? 2 : 1,
			values: [
				{ value: 0 } as unknown as IconEnumValueConfig<Player<any>, T>,
			].concat(options.map(option => args.allOptions.find(allOption => allOption.value == option)!)),
			equals: (a: T, b: T) => a == b,
			zeroValue: 0 as T,
			changedEvent: (player: Player<any>) => TypedEvent.onAny([player.consumesChangeEmitter, player.levelChangeEmitter]),
			getValue: (player: Player<any>) => player.getConsumes()[args.consumesFieldName] as T,
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const newConsumes = player.getConsumes();
				(newConsumes[args.consumesFieldName] as number) = newValue;
				TypedEvent.freezeAllAndDo(() => {
					player.setConsumes(eventID, newConsumes);
					if (args.onSet) {
						args.onSet(eventID, player, newValue as T);
					}
				});
			},
		};
	};
}

function makeConsumeInput<T extends number>(consumesFieldName: keyof Consumes, allOptions: Array<IconEnumValueConfig<Player<any>, T>>, onSet?: (eventID: EventID, player: Player<any>, newValue: T) => void): InputHelpers.TypedIconEnumPickerConfig<Player<any>, T> {
	const factory = makeConsumeInputFactory({
		consumesFieldName: consumesFieldName,
		allOptions: allOptions,
		onSet: onSet
	});
	return factory(allOptions.map(option => option.value));
}
