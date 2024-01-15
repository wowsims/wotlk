import {
	Consumes,
	Debuffs,
	Faction,
	IndividualBuffs,
	RaidBuffs,
} from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';

import { Player } from '../player';
import { Spec } from '../proto/common';
import { Raid } from '../raid';
import { EventID, TypedEvent } from '../typed_event';

import { IconEnumPicker, IconEnumPickerDirection, IconEnumValueConfig } from './icon_enum_picker';
import { IconPicker } from './icon_picker';

import * as InputHelpers from './input_helpers';
import { MAX_CHARACTER_LEVEL } from '../constants/mechanics';

// Component Functions

export type IconInputConfig<ModObject, T> = (
	InputHelpers.TypedIconPickerConfig<ModObject, T> |
	InputHelpers.TypedIconEnumPickerConfig<ModObject, T>
);

export const buildIconInput = (parent: HTMLElement, p: Player<Spec>, inputConfig: IconInputConfig<Player<Spec>, any>) => {
	if (inputConfig.type == 'icon') {
		return new IconPicker<Player<Spec>, any>(parent, p, inputConfig);
	} else if (inputConfig.type == 'iconEnum') {
		return new IconEnumPicker<Player<Spec>, any>(parent, p, inputConfig);
	}
};

// TODO: These should be moved to consumes
export const GiftOfArthas = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(11374), fieldName: 'giftOfArthas'}),
	'Gift of Arthas',
);
export const CrystalYield = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(15235), fieldName: 'crystalYield'}),
	'Crystal Yield',
);

// Consumes
export const Sapper = makeBooleanConsumeInput({id: ActionId.fromItemId(10646), fieldName: 'sapper', minLevel: 40});

// TODO: Classic
// export const PetScrollOfAgilityV = makeBooleanConsumeInput(ActionId.fromItemId(27498), 'petScrollOfAgility', 5);
// export const PetScrollOfStrengthV = makeBooleanConsumeInput(ActionId.fromItemId(27503), 'petScrollOfStrength', 5);

// eslint-disable-next-line unused-imports/no-unused-vars
export function withLabel<ModObject, T>(config: InputHelpers.TypedIconPickerConfig<ModObject, T>, label: string): InputHelpers.TypedIconPickerConfig<ModObject, T> {
	config.label = label;
	return config;
}

interface BooleanInputConfig<T> {
	id: ActionId
	fieldName: keyof T
	value?: number
	minLevel?: number
	maxLevel?: number
	faction?: Faction
}

export function makeBooleanRaidBuffInput<SpecType extends Spec>(config: BooleanInputConfig<RaidBuffs>): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, RaidBuffs, Player<SpecType>>({
		getModObject: (p: Player<any>) => p,
		showWhen: (p) =>
			(config.minLevel || 0) <= p.getLevel() &&
			p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL) &&
			(!config.faction || config.faction == p.getFaction()),
		getValue: (p) => p.getRaid()!.getBuffs(),
		setValue: (eventID: EventID, p: Player<SpecType>, newVal: RaidBuffs) => p.getRaid()!.setBuffs(eventID, newVal),
		changeEmitter: (p) => TypedEvent.onAny([p.getRaid()!.buffsChangeEmitter, p.levelChangeEmitter, p.raceChangeEmitter]),
	}, config.id, config.fieldName, config.value);
}
// function makeBooleanPartyBuffInput(id: ActionId, fieldName: keyof PartyBuffs, value?: number): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
// 	return InputHelpers.makeBooleanIconInput<any, PartyBuffs, Party>({
// 		getModObject: (p: Player<any>) => p.getParty()!,
// 		getValue: (party: Party) => party.getBuffs(),
// 		setValue: (eventID: EventID, party: Party, newVal: PartyBuffs) => party.setBuffs(eventID, newVal),
// 		changeEmitter: (party: Party) => party.buffsChangeEmitter,
// 	}, id, fieldName, value);
// }

export function makeBooleanIndividualBuffInput(config: BooleanInputConfig<IndividualBuffs>): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, IndividualBuffs, Player<any>>({
		getModObject: (p: Player<any>) => p,
		showWhen: (p) =>
			(config.minLevel || 0) <= p.getLevel() &&
			p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL) &&
			(!config.faction || config.faction == p.getFaction()),
		getValue: (p: Player<any>) => p.getBuffs(),
		setValue: (eventID: EventID, p: Player<any>, newVal: IndividualBuffs) => p.setBuffs(eventID, newVal),
		changeEmitter: (p: Player<any>) => TypedEvent.onAny([p.buffsChangeEmitter, p.levelChangeEmitter, p.raceChangeEmitter]),
	}, config.id, config.fieldName, config.value);
}

// eslint-disable-next-line unused-imports/no-unused-vars
export function makeBooleanConsumeInput<SpecType extends Spec>(config: BooleanInputConfig<Consumes>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, Consumes, Player<any>>({
		getModObject: (p: Player<SpecType>) => p,
		showWhen: (p) => p.getLevel() >= (config.minLevel || 0),
		getValue: (p: Player<any>) => p.getConsumes(),
		setValue: (eventID: EventID, p: Player<any>, newVal: Consumes) => p.setConsumes(eventID, newVal),
		changeEmitter: (p: Player<any>) => TypedEvent.onAny([p.consumesChangeEmitter, p.levelChangeEmitter])
	}, config.id, config.fieldName, config.value);
}
export function makeBooleanDebuffInput<SpecType extends Spec>(config: BooleanInputConfig<Debuffs>): InputHelpers.TypedIconPickerConfig<Player<any>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, Debuffs, Player<SpecType>>({
		getModObject: (p: Player<SpecType>) => p,
		showWhen: (p) => (config.minLevel || 0) <= p.getLevel() && p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL),
		getValue: (p) => p.getRaid()!.getDebuffs(),
		setValue: (eventID: EventID, p: Player<SpecType>, newVal: Debuffs) => p.getRaid()!.setDebuffs(eventID, newVal),
		changeEmitter: (p) => TypedEvent.onAny([p.getRaid()!.debuffsChangeEmitter, p.levelChangeEmitter]),
	}, config.id, config.fieldName, config.value);
}

interface TristateInputConfig<T> {
	id: ActionId
	impId: ActionId
	fieldName: keyof T
	minLevel?: number
	maxLevel?: number
	faction?: Faction
}

export function makeTristateRaidBuffInput<SpecType extends Spec>(config: TristateInputConfig<RaidBuffs>): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeTristateIconInput<any, RaidBuffs, Player<SpecType>>({
		getModObject: (p: Player<SpecType>) => p,
		showWhen: (p: Player<SpecType>) =>
			(config.minLevel || 0) <= p.getLevel() &&
			p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL) &&
			(!config.faction || config.faction == p.getFaction()),
		getValue: (p: Player<SpecType>) => p.getRaid()!.getBuffs(),
		setValue: (eventID: EventID, p: Player<SpecType>, newVal: RaidBuffs) => p.getRaid()!.setBuffs(eventID, newVal),
		changeEmitter: (p: Player<SpecType>) => TypedEvent.onAny([p.getRaid()!.buffsChangeEmitter, p.levelChangeEmitter, p.raceChangeEmitter]),
	}, config.id, config.impId, config.fieldName);
}

export function makeTristateIndividualBuffInput(config: TristateInputConfig<IndividualBuffs>): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeTristateIconInput<any, IndividualBuffs, Player<any>>({
		getModObject: (p: Player<any>) => p,
		showWhen: (p: Player<any>) =>
			(config.minLevel || 0) <= p.getLevel() &&
			p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL) &&
			(!config.faction || config.faction == p.getFaction()),
		getValue: (p: Player<any>) => p.getBuffs(),
		setValue: (eventID: EventID, p: Player<any>, newVal: IndividualBuffs) => p.setBuffs(eventID, newVal),
		changeEmitter: (p: Player<any>) => TypedEvent.onAny([p.buffsChangeEmitter, p.levelChangeEmitter, p.raceChangeEmitter])
	}, config.id, config.impId, config.fieldName);
}

export function makeTristateDebuffInput(id: ActionId, impId: ActionId, fieldName: keyof Debuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeTristateIconInput<any, Debuffs, Raid>({
		getModObject: (p: Player<any>) => p.getRaid()!,
		getValue: (raid: Raid) => raid.getDebuffs(),
		setValue: (eventID: EventID, raid: Raid, newVal: Debuffs) => raid.setDebuffs(eventID, newVal),
		changeEmitter: (raid: Raid) => raid.debuffsChangeEmitter,
	}, id, impId, fieldName);
}
// function makeQuadstateDebuffInput(id: ActionId, impId: ActionId, impId2: ActionId, fieldName: keyof Debuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
// 	return InputHelpers.makeQuadstateIconInput<any, Debuffs, Raid>({
// 		getModObject: (p: Player<any>) => p.getRaid()!,
// 		getValue: (raid: Raid) => raid.getDebuffs(),
// 		setValue: (eventID: EventID, raid: Raid, newVal: Debuffs) => raid.setDebuffs(eventID, newVal),
// 		changeEmitter: (raid: Raid) => raid.debuffsChangeEmitter,
// 	}, id, impId, impId2, fieldName);
// }

interface MultiStateInputConfig<T> {
	id: ActionId
	numStates: number
	fieldName: keyof T
	multiplier?: number
	minLevel?: number
	maxLevel?: number
	faction?: Faction
}

export function makeMultistateRaidBuffInput<SpecType extends Spec>(config: MultiStateInputConfig<RaidBuffs>): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeMultistateIconInput<any, RaidBuffs, Player<SpecType>>({
		getModObject: (p: Player<SpecType>) => p,
		showWhen: (p: Player<SpecType>) =>
			(config.minLevel || 0) <= p.getLevel() &&
			p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL) &&
			(!config.faction || config.faction == p.getFaction()),
		getValue: (p: Player<SpecType>) => p.getRaid()!.getBuffs(),
		setValue: (eventID: EventID, p: Player<SpecType>, newVal: RaidBuffs) => p.getRaid()!.setBuffs(eventID, newVal),
		changeEmitter: (p: Player<SpecType>) => TypedEvent.onAny([p.getRaid()!.buffsChangeEmitter, p.levelChangeEmitter, p.raceChangeEmitter]),
	}, config.id, config.numStates, config.fieldName, config.multiplier);
}
// function makeMultistatePartyBuffInput(id: ActionId, numStates: number, fieldName: keyof PartyBuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
// 	return InputHelpers.makeMultistateIconInput<any, PartyBuffs, Party>({
// 		getModObject: (p: Player<any>) => p.getParty()!,
// 		getValue: (party: Party) => party.getBuffs(),
// 		setValue: (eventID: EventID, party: Party, newVal: PartyBuffs) => party.setBuffs(eventID, newVal),
// 		changeEmitter: (party: Party) => party.buffsChangeEmitter,
// 	}, id, numStates, fieldName);
// }
export function makeMultistateIndividualBuffInput<SpecType extends Spec>(config: MultiStateInputConfig<IndividualBuffs>): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeMultistateIconInput<any, IndividualBuffs, Player<SpecType>>({
		getModObject: (p: Player<SpecType>) => p,
		showWhen: (p: Player<SpecType>) =>
			(config.minLevel || 0) <= p.getLevel() &&
			p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL) &&
			(!config.faction || config.faction == p.getFaction()),
		getValue: (p: Player<SpecType>) => p.getBuffs(),
		setValue: (eventID: EventID, p: Player<SpecType>, newVal: IndividualBuffs) => p.setBuffs(eventID, newVal),
		changeEmitter: (p: Player<SpecType>) => TypedEvent.onAny([p.buffsChangeEmitter, p.levelChangeEmitter, p.raceChangeEmitter]),
	}, config.id, config.numStates, config.fieldName, config.multiplier);
}
// function makeMultistateMultiplierIndividualBuffInput(id: ActionId, numStates: number, multiplier: number, fieldName: keyof IndividualBuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
// 	return InputHelpers.makeMultistateIconInput<any, IndividualBuffs, Player<any>>({
// 		getModObject: (p: Player<any>) => p,
// 		getValue: (p: Player<any>) => p.getBuffs(),
// 		setValue: (eventID: EventID, p: Player<any>, newVal: IndividualBuffs) => p.setBuffs(eventID, newVal),
// 		changeEmitter: (p: Player<any>) => p.buffsChangeEmitter,
// 	}, id, numStates, fieldName, multiplier);
// }

interface EnumInputConfig<ModObject, Message, T> {
	fieldName: keyof Message
	values: Array<IconEnumValueConfig<ModObject, T>>
	direction?: IconEnumPickerDirection
	numColumns?: number
	minLevel?: number
	maxLevel?: number
	faction?: Faction
}

export function makeEnumIndividualBuffInput<SpecType extends Spec>(config: EnumInputConfig<Player<any>, IndividualBuffs, number>): InputHelpers.TypedIconEnumPickerConfig<Player<any>, number> {
	return InputHelpers.makeEnumIconInput<any, IndividualBuffs, Player<SpecType>, number>({
		getModObject: (p: Player<SpecType>) => p,
		showWhen: (p: Player<SpecType>) =>
			(config.minLevel || 0) <= p.getLevel() &&
			p.getLevel() <= (config.maxLevel || MAX_CHARACTER_LEVEL) &&
			(!config.faction || config.faction == p.getFaction()),
		getValue: (p: Player<SpecType>) => p.getBuffs(),
		setValue: (eventID: EventID, p: Player<SpecType>, newVal: IndividualBuffs) => p.setBuffs(eventID, newVal),
		changeEmitter: (p: Player<SpecType>) => TypedEvent.onAny([p.buffsChangeEmitter, p.levelChangeEmitter, p.raceChangeEmitter]),
	}, config.fieldName, config.values, config.numColumns || 1, config.direction || IconEnumPickerDirection.Vertical)
};
