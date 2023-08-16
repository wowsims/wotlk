import { Player } from '../player.js';
import { Class } from '../proto/common.js';
import { Spec } from '../proto/common.js';
import {
	SpecTalents,
	specToClass,
	specTypeFunctions,
} from '../proto_utils/utils.js';
import { EventID, TypedEvent } from '../typed_event.js';

import { druidTalentsConfig, druidGlyphsConfig } from './druid.js';
import { hunterTalentsConfig, hunterGlyphsConfig } from './hunter.js';
import { mageTalentsConfig, mageGlyphsConfig } from './mage.js';
import { paladinTalentsConfig, paladinGlyphsConfig } from './paladin.js';
import { priestTalentsConfig, priestGlyphsConfig } from './priest.js';
import { rogueTalentsConfig, rogueGlyphsConfig } from './rogue.js';
import { shamanTalentsConfig, shamanGlyphsConfig } from './shaman.js';
import { warlockTalentsConfig, warlockGlyphsConfig } from './warlock.js';
import { warriorTalentsConfig, warriorGlyphsConfig } from './warrior.js';
import { deathknightTalentsConfig, deathknightGlyphsConfig } from './deathknight.js';
import { TalentsConfig } from './talents_picker.js';
import { GlyphsConfig } from './glyphs_picker.js';

export const classTalentsConfig: Record<Class, TalentsConfig<any>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassDeathknight]: deathknightTalentsConfig,
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

export const classGlyphsConfig: Record<Class, GlyphsConfig> = {
	[Class.ClassUnknown]: { majorGlyphs: [], minorGlyphs: [] },
	[Class.ClassDeathknight]: deathknightGlyphsConfig,
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

export function playerTalentStringToProto<SpecType extends Spec>(spec: Spec, talentString: string): SpecTalents<SpecType> {
	const specFunctions = specTypeFunctions[spec];
	const proto = specFunctions.talentsCreate() as SpecTalents<SpecType>;
	const talentsConfig = classTalentsConfig[specToClass[spec]] as TalentsConfig<SpecTalents<SpecType>>;

	return talentStringToProto(proto, talentString, talentsConfig);
}

export function talentStringToProto<TalentsProto>(proto: TalentsProto, talentString: string, talentsConfig: TalentsConfig<TalentsProto>): TalentsProto {
	talentString.split('-').forEach((treeString, treeIdx) => {
		const treeConfig = talentsConfig[treeIdx];
		[...treeString].forEach((talentString, i) => {
			const talentConfig = treeConfig.talents[i];
			const points = parseInt(talentString);
			if (talentConfig.fieldName) {
				if (talentConfig.maxPoints == 1) {
					(proto[talentConfig.fieldName as keyof TalentsProto] as unknown as boolean) = points == 1;
				} else {
					(proto[talentConfig.fieldName as keyof TalentsProto] as unknown as number) = points;
				}
			}
		});
	});

	return proto;
}

// Note that this function will fail if any of the talent names are not defined. TODO: Remove that condition
// once all talents are migrated to wrath and use all fields.
export function protoToTalentString<TalentsProto>(proto: TalentsProto, talentsConfig: TalentsConfig<TalentsProto>): string {
	return talentsConfig.map(treeConfig => {
		return treeConfig.talents
			.map(talentConfig => String(Number(proto[(talentConfig.fieldName as keyof TalentsProto)!])))
			.join('').replace(/0+$/g, '');
	}).join('-').replace(/-+$/g, '');
}
