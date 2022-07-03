import { Player } from '/wotlk/core/player.js';
import { Class } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import {
	SpecTalents,
	specToClass,
	specTypeFunctions,
} from '/wotlk/core/proto_utils/utils.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import { druidTalentsConfig, druidGlyphsConfig } from './druid.js';
import { hunterTalentsConfig, hunterGlyphsConfig } from './hunter.js';
import { mageTalentsConfig, mageGlyphsConfig } from './mage.js';
import { paladinTalentsConfig, paladinGlyphsConfig } from './paladin.js';
import { priestTalentsConfig, priestGlyphsConfig } from './priest.js';
import { rogueTalentsConfig, rogueGlyphsConfig } from './rogue.js';
import { shamanTalentsConfig, shamanGlyphsConfig } from './shaman.js';
import { warlockTalentsConfig, warlockGlyphsConfig } from './warlock.js';
import { warriorTalentsConfig, warriorGlyphsConfig } from './warrior.js';
import { deathKnightTalentsConfig, deathKnightGlyphsConfig } from './deathknight.js';
import { TalentsConfig, TalentsPicker } from './talents_picker.js';
import { GlyphsConfig, GlyphsPicker } from './glyphs_picker.js';

import * as Mechanics from '/wotlk/core/constants/mechanics.js';

export function newTalentsPicker(parent: HTMLElement, player: Player<any>): TalentsPicker<Player<any>, any> {
	return new TalentsPicker(parent, player, classTalentsConfig[player.getClass()], {
		changedEvent: (player: Player<any>) => player.talentsChangeEmitter,
		getValue: (player: Player<any>) => player.getTalentsString(),
		setValue: (eventID: EventID, player: Player<any>, newValue: string) => {
			player.setTalentsString(eventID, newValue);
		},
		numRows: 9,
		pointsPerRow: 5,
		maxPoints: Mechanics.CHARACTER_LEVEL - 9,
	});
}

export function newGlyphsPicker(parent: HTMLElement, player: Player<any>): GlyphsPicker {
	return new GlyphsPicker(parent, player, classGlyphsConfig[player.getClass()]);
}

const classTalentsConfig: Record<Class, TalentsConfig<any>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassDeathKnight]: deathKnightTalentsConfig,
	[Class.ClassDruid]: druidTalentsConfig,
	[Class.ClassShaman]: shamanTalentsConfig,
	[Class.ClassHunter]: hunterTalentsConfig,
	[Class.ClassMage]: mageTalentsConfig,
	[Class.ClassRogue]: rogueTalentsConfig,
	[Class.ClassPaladin]: paladinTalentsConfig,
	[Class.ClassPriest]: priestTalentsConfig,
	[Class.ClassWarlock]: warlockTalentsConfig,
	[Class.ClassWarrior]: warriorTalentsConfig,
};

const classGlyphsConfig: Record<Class, GlyphsConfig> = {
	[Class.ClassUnknown]: { majorGlyphs: [], minorGlyphs: [] },
	[Class.ClassDeathKnight]: deathKnightGlyphsConfig,
	[Class.ClassDruid]: druidGlyphsConfig,
	[Class.ClassShaman]: shamanGlyphsConfig,
	[Class.ClassHunter]: hunterGlyphsConfig,
	[Class.ClassMage]: mageGlyphsConfig,
	[Class.ClassRogue]: rogueGlyphsConfig,
	[Class.ClassPaladin]: paladinGlyphsConfig,
	[Class.ClassPriest]: priestGlyphsConfig,
	[Class.ClassWarlock]: warlockGlyphsConfig,
	[Class.ClassWarrior]: warriorGlyphsConfig,
};

export function talentSpellIdsToTalentString(playerClass: Class, talentIds: Array<number>): string {
	const talentsConfig = classTalentsConfig[playerClass];

	const talentsStr = talentsConfig.map(treeConfig => {
		const treeStr = treeConfig.talents.map(talentConfig => {
			const spellIdIndex = talentConfig.spellIds.findIndex(spellId => talentIds.includes(spellId));
			if (spellIdIndex == -1) {
				return '0';
			} else {
				return String(spellIdIndex + 1);
			}
		}).join('').replace(/0+$/g, '');

		return treeStr;
	}).join('-').replace(/-+$/g, '');

	return talentsStr
}

export function talentStringToProto<SpecType extends Spec>(spec: Spec, talentString: string): SpecTalents<SpecType> {
	const talentsConfig = classTalentsConfig[specToClass[spec]] as TalentsConfig<SpecTalents<SpecType>>;

	const specFunctions = specTypeFunctions[spec];
	const proto = specFunctions.talentsCreate() as SpecTalents<SpecType>;

	talentString.split('-').forEach((treeString, treeIdx) => {
		const treeConfig = talentsConfig[treeIdx];
		[...treeString].forEach((talentString, i) => {
			const talentConfig = treeConfig.talents[i];
			const points = parseInt(talentString);
			if (talentConfig.fieldName) {
				if (talentConfig.maxPoints == 1) {
					(proto[talentConfig.fieldName] as unknown as boolean) = points == 1;
				} else {
					(proto[talentConfig.fieldName] as unknown as number) = points;
				}
			}
		});
	});

	return proto;
}
