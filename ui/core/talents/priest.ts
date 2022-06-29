import { Spec } from '/tbc/core/proto/common.js';
import { PriestTalents as PriestTalents } from '/tbc/core/proto/priest.js';
import { PriestSpecs } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';

import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

// Talents are the same for all Priest specs, so its ok to just use ShadowPriest here
export class PriestTalentsPicker extends TalentsPicker<Spec.SpecShadowPriest> {
	constructor(parent: HTMLElement, player: Player<Spec.SpecShadowPriest>) {
		super(parent, player, priestTalentsConfig);
	}
}

export const priestTalentsConfig: TalentsConfig<Spec.SpecShadowPriest> = newTalentsConfig([
	{
		name: 'Discipline',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/201.jpg',
		talents: [
			{
				//fieldName: 'unbreakableWill',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [14522, 14788],
				maxPoints: 5,
			},
			{
				fieldName: 'wandSpecialization',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [14524],
				maxPoints: 5,
			},
			{
				fieldName: 'silentResolve',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [14523, 14784],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedPowerWordFortitude',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [14749, 14767],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedPowerWordShield',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [14748, 14768],
				maxPoints: 3,
			},
			{
				//fieldName: 'martyrdom',
				location: {
					rowIdx: 1,
					colIdx: 3,
				},
				spellIds: [14531, 14774],
				maxPoints: 2,
			},
			{
				//fieldName: 'absolution',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [33167, 33171],
				maxPoints: 3,
			},
			{
				fieldName: 'innerFocus',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [14751],
				maxPoints: 1,
			},
			{
				fieldName: 'meditation',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [14521, 14776],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedInnerFire',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [14747, 14770],
				maxPoints: 3,
			},
			{
				fieldName: 'mentalAgility',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [14520, 14780],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedManaBurn',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [14750, 14772],
				maxPoints: 2,
			},
			{
				fieldName: 'mentalStrength',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [18551],
				maxPoints: 5,
			},
			{
				fieldName: 'divineSpirit',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [14752],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedDivineSpirit',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [33174, 33182],
				maxPoints: 2,
			},
			{
				fieldName: 'focusedPower',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [33186, 33190],
				maxPoints: 2,
			},
			{
				fieldName: 'forceOfWill',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [18544, 18547],
				maxPoints: 5,
			},
			{
				//fieldName: 'focusedWill',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [45234, 45243],
				maxPoints: 3,
			},
			{
				fieldName: 'powerInfusion',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [10060],
				maxPoints: 1,
			},
			{
				//fieldName: 'reflectiveShield',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [33201],
				maxPoints: 5,
			},
			{
				fieldName: 'enlightenment',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [34908],
				maxPoints: 5,
			},
			{
				//fieldName: 'painSuppresion',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [33206],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Holy',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/202.jpg',
		talents: [
			{
				//fieldName: 'healingFocus',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [14913, 15012],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedRenew',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [14908, 15020, 17191],
				maxPoints: 3,
			},
			{
				fieldName: 'holySpecialization',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [14889, 15008],
				maxPoints: 5,
			},

			{
				//fieldName: 'spellWarding',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [27900],
				maxPoints: 5,
			},
			{
				fieldName: 'divineFury',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [18530, 18531, 18533],
				maxPoints: 5,
			},
			{
				fieldName: 'holyNova',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [15237],
				maxPoints: 1,
			},
			{
				//fieldName: 'blessedRecovery',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [27811, 27815],
				maxPoints: 3,
			},
			{
				//fieldName: 'inspiration',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [14892, 15362],
				maxPoints: 3,
			},
			{
				//fieldName: 'holyReach',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [27789],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedHealing',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [14912, 15013],
				maxPoints: 3,
			},
			{
				fieldName: 'searingLight',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [14909, 15017],
				maxPoints: 2,
			},
			{
				//fieldName: 'healingPrayers',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [14911, 15018],
				maxPoints: 2,
			},
			{
				fieldName: 'spiritOfRedemption',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [20711],
				maxPoints: 1,
			},
			{
				fieldName: 'spiritualGuidance',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [14901, 15028],
				maxPoints: 5,
			},
			{
				fieldName: 'surgeOfLight',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [33150, 33154],
				maxPoints: 2,
			},
			{
				//fieldName: 'spiritualHealing',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [14898, 15349, 15354],
				maxPoints: 5,
			},
			{
				//fieldName: 'holyConcentration',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [34753, 34859],
				maxPoints: 3,
			},
			{
				//fieldName: 'lightwell',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [724],
				maxPoints: 1,
			},
			{
				//fieldName: 'blessedResilience',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [33142, 33145],
				maxPoints: 3,
			},
			{
				//fieldName: 'empoweredHealing',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [33158],
				maxPoints: 5,
			},
			{
				//fieldName: 'circleOfHealing',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [34861],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Shadow',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/203.jpg',
		talents: [
			{
				//fieldName: 'spiritTap',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [15270, 15335],
				maxPoints: 5,
			},
			{
				//fieldName: 'blackout',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [15268, 15323],
				maxPoints: 5,
			},
			{
				fieldName: 'shadowAffinity',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [15318, 15272, 15320],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedShadowWordPain',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [15275, 15317],
				maxPoints: 2,
			},
			{
				fieldName: 'shadowFocus',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [15260, 15327],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedPsychicScream',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [15392, 15448],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedMindBlast',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [15273, 15312, 15313, 15314, 15316],
				maxPoints: 5,
			},
			{
				fieldName: 'mindFlay',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [15407],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedFade',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [15274, 15311],
				maxPoints: 2,
			},
			{
				//fieldName: 'shadowReach',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [17322],
				maxPoints: 2,
			},
			{
				fieldName: 'shadowWeaving',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [15257, 15331],
				maxPoints: 5,
			},
			{
				//fieldName: 'silence',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [15487],
				maxPoints: 1,
			},
			{
				fieldName: 'vampiricEmbrace',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [15286],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedVampiricEmbrace',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [27839],
				maxPoints: 2,
			},
			{
				fieldName: 'focusedMind',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [33213],
				maxPoints: 3,
			},
			{
				//fieldName: 'shadowResilience',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [14910, 33371],
				maxPoints: 2,
			},
			{
				fieldName: 'darkness',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [15259, 15307],
				maxPoints: 5,
			},
			{
				fieldName: 'shadowform',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [15473],
				maxPoints: 1,
			},
			{
				fieldName: 'shadowPower',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [33221],
				maxPoints: 5,
			},
			{
				fieldName: 'misery',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [33191],
				maxPoints: 5,
			},
			{
				fieldName: 'vampiricTouch',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [34914],
				maxPoints: 1,
			},
		],
	},
]);
