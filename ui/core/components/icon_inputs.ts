import { AgilityElixir, Consumes, Debuffs, Explosive, Flask, Food, IndividualBuffs, RaidBuffs, SpellPowerBuff, StrengthBuff, WeaponBuff } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';

import { Player } from '../player.js';
import { Spec } from '../proto/common.js';
import { Raid } from '../raid.js';
import { EventID, TypedEvent } from '../typed_event.js';

import { IconEnumPicker, IconEnumValueConfig } from './icon_enum_picker.js';
import { IconPicker } from './icon_picker.js';

import * as InputHelpers from './input_helpers.js';

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

export const AllStatsBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput(ActionId.fromSpellId(21850), ActionId.fromSpellId(17055), 'giftOfTheWild'),
], 'Stats');

export const AllStatsPercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput(ActionId.fromSpellId(20217), 'blessingOfKings'),
], 'Stats %');

export const ArmorBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput(ActionId.fromSpellId(10293), ActionId.fromSpellId(20142), 'devotionAura'),
	makeBooleanRaidBuffInput(ActionId.fromItemId(43468), 'scrollOfProtection'),
], 'Armor');

export const StaminaBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput(ActionId.fromSpellId(10938), ActionId.fromSpellId(14767), 'powerWordFortitude'),
	makeBooleanRaidBuffInput(ActionId.fromItemId(10307), 'scrollOfStamina'),
], 'Stamina');

// TODO: Breakout Strength / Agi
export const StrengthAndAgilityBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput(ActionId.fromSpellId(25361), ActionId.fromSpellId(52456), 'strengthOfEarthTotem'),
	makeBooleanRaidBuffInput(ActionId.fromItemId(10309), 'scrollOfAgility'),
	makeBooleanRaidBuffInput(ActionId.fromItemId(10310), 'scrollOfStrength'),
], 'Str/Agi');

export const IntellectBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput(ActionId.fromSpellId(23028), 'arcaneBrilliance'),
	makeBooleanRaidBuffInput(ActionId.fromItemId(10308), 'scrollOfIntellect'),
], 'Int');

export const SpiritBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput(ActionId.fromSpellId(27841), 'divineSpirit'),
	makeBooleanRaidBuffInput(ActionId.fromItemId(10306), 'scrollOfSpirit'),
], 'Spirit');

export const AttackPowerBuff = InputHelpers.makeMultiIconInput([
	makeTristateIndividualBuffInput(ActionId.fromSpellId(48934), ActionId.fromSpellId(20045), 'blessingOfMight'),
	makeTristateRaidBuffInput(ActionId.fromSpellId(47436), ActionId.fromSpellId(12861), 'battleShout'),
	makeBooleanRaidBuffInput(ActionId.fromSpellId(19506), 'trueshotAura'),
], 'AP');

export const AttackPowerPercentBuff = InputHelpers.makeMultiIconInput([
], 'Atk Pwr %');

export const DamageReductionPercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput(ActionId.fromSpellId(25899), 'blessingOfSanctuary'),
], 'Mit %');

export const ResistanceBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput(ActionId.fromSpellId(48170), 'shadowProtection'),
	makeBooleanRaidBuffInput(ActionId.fromSpellId(58749), 'natureResistanceTotem'),
	makeBooleanRaidBuffInput(ActionId.fromSpellId(49071), 'aspectOfTheWild'),
	makeBooleanRaidBuffInput(ActionId.fromSpellId(48945), 'frostResistanceAura'),
	makeBooleanRaidBuffInput(ActionId.fromSpellId(58745), 'frostResistanceTotem'),
], 'Resistances');

// export const HealthBuff = InputHelpers.makeMultiIconInput([
// 	makeTristateRaidBuffInput(ActionId.fromSpellId(47982), ActionId.fromSpellId(18696), 'bloodPact'),
// ], 'Health');

export const MP5Buff = InputHelpers.makeMultiIconInput([
	makeTristateIndividualBuffInput(ActionId.fromSpellId(25290), ActionId.fromSpellId(20245), 'blessingOfWisdom'),
	makeTristateRaidBuffInput(ActionId.fromSpellId(10494), ActionId.fromSpellId(16208), 'manaSpringTotem'),
], 'MP5');

export const MeleeCritBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput(ActionId.fromSpellId(17007), ActionId.fromSpellId(34300), 'leaderOfThePack'),
], 'Melee Crit');

export const SpellCritBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput(ActionId.fromSpellId(24907), 'moonkinAura'),
], 'Spell Crit');

export const SpellIncreaseBuff = InputHelpers.makeMultiIconInput([
	// makeMultistateRaidBuffInput(ActionId.fromSpellId(47240), 2000, 'demonicPactSp', 20),
], 'Spell Power');

export const DefensiveCooldownBuff = InputHelpers.makeMultiIconInput([
], 'Defensive CDs');

// Misc Buffs
export const RetributionAura = makeBooleanRaidBuffInput(ActionId.fromSpellId(10301), 'retributionAura');
export const Thorns = makeTristateRaidBuffInput(ActionId.fromSpellId(9910), ActionId.fromSpellId(16840), 'thorns');
export const Innervate = makeMultistateIndividualBuffInput(ActionId.fromSpellId(29166), 11, 'innervates');
export const PowerInfusion = makeMultistateIndividualBuffInput(ActionId.fromSpellId(10060), 11, 'powerInfusions');

// Debuffs

export const MajorArmorDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput(ActionId.fromSpellId(47467), 'sunderArmor'),
	makeBooleanDebuffInput(ActionId.fromSpellId(8647), 'exposeArmor'),
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
	makeBooleanDebuffInput(ActionId.fromSpellId(65855), 'insectSwarm'),
], 'Miss');

// TODO: Classic
export const SpellISBDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput(ActionId.fromSpellId(17803), 'improvedShadowBolt'),
], 'ISB');

export const SpellScorchDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput(ActionId.fromSpellId(12873), 'improvedScorch'),
], 'Scorch');

export const SpellWintersChillDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput(ActionId.fromSpellId(28595), 'wintersChill'),
], 'Winters Chill');

// TODO: Classic
// export const SpellDamageDebuff = InputHelpers.makeMultiIconInput([
// 	makeBooleanDebuffInput(ActionId.fromSpellId(47865), 'curseOfElements'),
// ], 'Spell Dmg');

// TODO: Classic
// export const HuntersMark = withLabel(makeQuadstateDebuffInput(ActionId.fromSpellId(53338), ActionId.fromSpellId(19423), ActionId.fromItemId(42907), 'huntersMark'), 'Mark');
// export const JudgementOfWisdom = withLabel(makeBooleanDebuffInput(ActionId.fromSpellId(53408), 'judgementOfWisdom'), 'JoW');
// export const JudgementOfLight = makeBooleanDebuffInput(ActionId.fromSpellId(20271), 'judgementOfLight');
export const GiftOfArthas = makeBooleanDebuffInput(ActionId.fromSpellId(11374), 'giftOfArthas');
export const CrystalYield = makeBooleanDebuffInput(ActionId.fromSpellId(15235), 'crystalYield');

// Consumes
export const Sapper = makeBooleanConsumeInput(ActionId.fromItemId(10646), 'sapper');

// TODO: Classic
// export const PetScrollOfAgilityV = makeBooleanConsumeInput(ActionId.fromItemId(27498), 'petScrollOfAgility', 5);
// export const PetScrollOfStrengthV = makeBooleanConsumeInput(ActionId.fromItemId(27503), 'petScrollOfStrength', 5);

// eslint-disable-next-line unused-imports/no-unused-vars
function withLabel<ModObject, T>(config: InputHelpers.TypedIconPickerConfig<ModObject, T>, label: string): InputHelpers.TypedIconPickerConfig<ModObject, T> {
	config.label = label;
	return config;
}

function makeBooleanRaidBuffInput(id: ActionId, fieldName: keyof RaidBuffs, value?: number): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, RaidBuffs, Raid>({
		getModObject: (player: Player<any>) => player.getRaid()!,
		getValue: (raid: Raid) => raid.getBuffs(),
		setValue: (eventID: EventID, raid: Raid, newVal: RaidBuffs) => raid.setBuffs(eventID, newVal),
		changeEmitter: (raid: Raid) => raid.buffsChangeEmitter,
	}, id, fieldName, value);
}
// function makeBooleanPartyBuffInput(id: ActionId, fieldName: keyof PartyBuffs, value?: number): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
// 	return InputHelpers.makeBooleanIconInput<any, PartyBuffs, Party>({
// 		getModObject: (player: Player<any>) => player.getParty()!,
// 		getValue: (party: Party) => party.getBuffs(),
// 		setValue: (eventID: EventID, party: Party, newVal: PartyBuffs) => party.setBuffs(eventID, newVal),
// 		changeEmitter: (party: Party) => party.buffsChangeEmitter,
// 	}, id, fieldName, value);
// }
function makeBooleanIndividualBuffInput(id: ActionId, fieldName: keyof IndividualBuffs, value?: number): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, IndividualBuffs, Player<any>>({
		getModObject: (player: Player<any>) => player,
		getValue: (player: Player<any>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<any>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<any>) => player.buffsChangeEmitter,
	}, id, fieldName, value);
}
// eslint-disable-next-line unused-imports/no-unused-vars
function makeBooleanConsumeInput(id: ActionId, fieldName: keyof Consumes, value?: number): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, Consumes, Player<any>>({
		getModObject: (player: Player<any>) => player,
		getValue: (player: Player<any>) => player.getConsumes(),
		setValue: (eventID: EventID, player: Player<any>, newVal: Consumes) => player.setConsumes(eventID, newVal),
		changeEmitter: (player: Player<any>) => player.consumesChangeEmitter,
	}, id, fieldName, value);
}
function makeBooleanDebuffInput(id: ActionId, fieldName: keyof Debuffs, value?: number): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, Debuffs, Raid>({
		getModObject: (player: Player<any>) => player.getRaid()!,
		getValue: (raid: Raid) => raid.getDebuffs(),
		setValue: (eventID: EventID, raid: Raid, newVal: Debuffs) => raid.setDebuffs(eventID, newVal),
		changeEmitter: (raid: Raid) => raid.debuffsChangeEmitter,
	}, id, fieldName, value);
}

function makeTristateRaidBuffInput(id: ActionId, impId: ActionId, fieldName: keyof RaidBuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeTristateIconInput<any, RaidBuffs, Raid>({
		getModObject: (player: Player<any>) => player.getRaid()!,
		getValue: (raid: Raid) => raid.getBuffs(),
		setValue: (eventID: EventID, raid: Raid, newVal: RaidBuffs) => raid.setBuffs(eventID, newVal),
		changeEmitter: (raid: Raid) => raid.buffsChangeEmitter,
	}, id, impId, fieldName);
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
	onSet: (eventID: EventID, player: Player<any>, newValue: Flask) => {
		if (newValue) {
			const newConsumes = player.getConsumes();
			player.setConsumes(eventID, newConsumes);
		}
	}
});

export const makeWeaponBuffsInput = makeConsumeInputFactory({
	consumesFieldName: 'weaponBuff',
	allOptions: [
		{ actionId: ActionId.fromItemId(20749), value: WeaponBuff.BrillianWizardOil },
		{ actionId: ActionId.fromItemId(20748), value: WeaponBuff.BrilliantManaOil },
		{ actionId: ActionId.fromItemId(12404), value: WeaponBuff.DenseSharpeningStone },
		{ actionId: ActionId.fromItemId(18262), value: WeaponBuff.ElementalSharpeningStone },
	] as Array<IconEnumValueConfig<Player<any>, WeaponBuff>>,
	onSet: (eventID: EventID, player: Player<any>, newValue: WeaponBuff) => {
		if (newValue) {
			const newConsumes = player.getConsumes();
			player.setConsumes(eventID, newConsumes);
		}
	}
});

export const makeFoodInput = makeConsumeInputFactory({
	consumesFieldName: 'food',
	allOptions: [
		{ actionId: ActionId.fromItemId(13931), value: Food.FoodNightfinSoup },
		{ actionId: ActionId.fromItemId(13928), value: Food.FoodGrilledSquid },
		{ actionId: ActionId.fromItemId(20452), value: Food.FoodSmokedDesertDumpling },
		{ actionId: ActionId.fromItemId(18254), value: Food.FoodRunnTumTuberSurprise },
		{ actionId: ActionId.fromItemId(21023), value: Food.FoodDirgesKickChimaerokChops },
		{ actionId: ActionId.fromItemId(13813), value: Food.FoodBlessedSunfruitJuice },
		{ actionId: ActionId.fromItemId(13810), value: Food.FoodBlessSunfruit },
	] as Array<IconEnumValueConfig<Player<any>, Food>>
});

export const AgilityBuffInput = makeConsumeInput('agilityElixir', [
	{ actionId: ActionId.fromItemId(13452), value: AgilityElixir.ElixirOfTheMongoose },
	{ actionId: ActionId.fromItemId(9187), value: AgilityElixir.ElixirOfGreaterAgility},
] as Array<IconEnumValueConfig<Player<any>, AgilityElixir>>);

export const StrengthBuffInput = makeConsumeInput('strengthBuff', [
	{ actionId: ActionId.fromItemId(12451), value: StrengthBuff.JujuPower },
	{ actionId: ActionId.fromItemId(9206), value: StrengthBuff.ElixirOfGiants },
] as Array<IconEnumValueConfig<Player<any>, StrengthBuff>>);

export const SpellDamageBuff = makeBooleanConsumeInput(ActionId.fromItemId(13454), 'spellPowerBuff');
export const ShadowDamageBuff = makeBooleanConsumeInput(ActionId.fromItemId(9264), 'shadowPowerBuff');
export const FireDamageBuff = makeBooleanConsumeInput(ActionId.fromItemId(21546), 'firePowerBuff');
export const FrostDamageBuff = makeBooleanConsumeInput(ActionId.fromItemId(17708), 'frostPowerBuff');

export const FillerExplosiveInput = makeConsumeInput('fillerExplosive', [
	{ actionId: ActionId.fromItemId(18641), value: Explosive.ExplosiveDenseDynamite },
	{ actionId: ActionId.fromItemId(15993), value: Explosive.ExplosiveThoriumGrenade },
] as Array<IconEnumValueConfig<Player<any>, Explosive>>);

export interface ConsumeInputFactoryArgs<T extends number> {
	consumesFieldName: keyof Consumes,
	allOptions: Array<IconEnumValueConfig<Player<any>, T>>,
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
			changedEvent: (player: Player<any>) => player.consumesChangeEmitter,
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
