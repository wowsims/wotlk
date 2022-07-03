import { Player } from '/wotlk/core/player.js';
import { Class } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import {
	SpecTalents,
	specToClass,
	specTypeFunctions,
} from '/wotlk/core/proto_utils/utils.js';

import { druidTalentsConfig, DruidTalentsPicker, DruidGlyphsPicker } from './druid.js';
import { hunterTalentsConfig, HunterTalentsPicker, HunterGlyphsPicker } from './hunter.js';
import { mageTalentsConfig, MageTalentsPicker, MageGlyphsPicker } from './mage.js';
import { paladinTalentsConfig, PaladinTalentsPicker, PaladinGlyphsPicker } from './paladin.js';
import { priestTalentsConfig, PriestTalentsPicker, PriestGlyphsPicker } from './priest.js';
import { rogueTalentsConfig, RogueTalentsPicker, RogueGlyphsPicker } from './rogue.js';
import { shamanTalentsConfig, ShamanTalentsPicker, ShamanGlyphsPicker } from './shaman.js';
import { warlockTalentsConfig, WarlockTalentsPicker, WarlockGlyphsPicker } from './warlock.js';
import { warriorTalentsConfig, WarriorTalentsPicker, WarriorGlyphsPicker } from './warrior.js';
import { deathKnightTalentsConfig, DeathKnightTalentsPicker, DeathKnightGlyphsPicker } from './deathknight.js';
import { TalentsConfig, TalentsPicker } from './talents_picker.js';
import { GlyphsPicker } from './glyphs_picker.js';

export function newTalentsPicker<SpecType extends Spec>(parent: HTMLElement, player: Player<SpecType>): TalentsPicker<SpecType> {
	switch (player.getClass()) {
		case Class.ClassDruid:
			return new DruidTalentsPicker(parent, player as Player<Spec.SpecBalanceDruid>) as TalentsPicker<SpecType>;
			break;
		case Class.ClassShaman:
			return new ShamanTalentsPicker(parent, player as Player<Spec.SpecElementalShaman>) as TalentsPicker<SpecType>;
			break;
		case Class.ClassHunter:
			return new HunterTalentsPicker(parent, player as Player<Spec.SpecHunter>) as TalentsPicker<SpecType>;
			break;
		case Class.ClassMage:
			return new MageTalentsPicker(parent, player as Player<Spec.SpecMage>) as TalentsPicker<SpecType>;
			break;
		case Class.ClassPaladin:
			return new PaladinTalentsPicker(parent, player as Player<Spec.SpecRetributionPaladin>) as TalentsPicker<SpecType>;
			break;
		case Class.ClassRogue:
			return new RogueTalentsPicker(parent, player as Player<Spec.SpecRogue>) as TalentsPicker<SpecType>;
			break;
		case Class.ClassPriest:
			return new PriestTalentsPicker(parent, player as Player<Spec.SpecShadowPriest>) as TalentsPicker<SpecType>;
			break;
		case Class.ClassWarlock:
			return new WarlockTalentsPicker(parent, player as Player<Spec.SpecWarlock>) as TalentsPicker<SpecType>;
			break;
		case Class.ClassWarrior:
			return new WarriorTalentsPicker(parent, player as Player<Spec.SpecWarrior>) as TalentsPicker<SpecType>;
			break;
		default:
			throw new Error('Unimplemented class talents: ' + player.getClass());
	}
}

export function newGlyphsPicker(parent: HTMLElement, player: Player<any>): GlyphsPicker | null {
	switch (player.getClass()) {
		case Class.ClassDruid:
			return new DruidGlyphsPicker(parent, player);
			break;
		case Class.ClassShaman:
			return new ShamanGlyphsPicker(parent, player);
			break;
		case Class.ClassHunter:
			return new HunterGlyphsPicker(parent, player);
			break;
		case Class.ClassMage:
			return new MageGlyphsPicker(parent, player);
			break;
		case Class.ClassPaladin:
			return new PaladinGlyphsPicker(parent, player);
			break;
		case Class.ClassRogue:
			return new RogueGlyphsPicker(parent, player);
			break;
		case Class.ClassPriest:
			return new PriestGlyphsPicker(parent, player);
			break;
		case Class.ClassWarlock:
			return new WarlockGlyphsPicker(parent, player);
			break;
		case Class.ClassWarrior:
			return new WarriorGlyphsPicker(parent, player);
			break;
	}
	return null;
	//throw new Error('Unimplemented class glyphs: ' + player.getClass());
}

const classTalentsConfig: Record<Class, TalentsConfig<any>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassDruid]: druidTalentsConfig,
	[Class.ClassShaman]: shamanTalentsConfig,
	[Class.ClassHunter]: hunterTalentsConfig,
	[Class.ClassMage]: mageTalentsConfig,
	[Class.ClassRogue]: rogueTalentsConfig,
	[Class.ClassPaladin]: paladinTalentsConfig,
	[Class.ClassPriest]: priestTalentsConfig,
	[Class.ClassWarlock]: warlockTalentsConfig,
	[Class.ClassWarrior]: warriorTalentsConfig,
	[Class.ClassDeathKnight]: deathKnightTalentsConfig,
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
	const talentsConfig = classTalentsConfig[specToClass[spec]] as TalentsConfig<SpecType>;

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
