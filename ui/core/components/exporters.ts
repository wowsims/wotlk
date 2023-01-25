import { Popup } from './popup';
import { IndividualSimUI } from '../individual_sim_ui';
import { SimUI } from '../sim_ui';
import {
	PseudoStat,
	Spec,
	Stat
} from '../proto/common';
import { IndividualSimSettings } from '../proto/ui';
import { classNames, raceNames } from '../proto_utils/names';
import { UnitStat } from '../proto_utils/stats';
import { specNames } from '../proto_utils/utils';
import { downloadString } from '../utils';
import { IndividualWowheadGearPlannerImporter } from './importers';

import * as Mechanics from '../constants/mechanics';

export abstract class Exporter extends Popup {
	private readonly textElem: HTMLElement;

	constructor(parent: HTMLElement, simUI: SimUI, title: string, allowDownload: boolean) {
		super(parent);

		this.rootElem.classList.add('exporter');
		this.rootElem.innerHTML = `
			<span class="exporter-title">${title}</span>
			<div class="export-content">
				<textarea class="exporter-textarea form-control" readonly></textarea>
			</div>
			<div class="actions-row">
				<button class="exporter-button btn btn-${simUI.cssScheme} clipboard-button">COPY TO CLIPBOARD</button>
				<button class="exporter-button btn btn-${simUI.cssScheme} download-button">DOWNLOAD</button>
			</div>
		`;

		this.addCloseButton();

		this.textElem = this.rootElem.getElementsByClassName('exporter-textarea')[0] as HTMLElement;

		const clipboardButton = this.rootElem.getElementsByClassName('clipboard-button')[0] as HTMLElement;
		clipboardButton.addEventListener('click', event => {
			const data = this.textElem.textContent!;
			if (navigator.clipboard == undefined) {
				alert(data);
			} else {
				navigator.clipboard.writeText(data);
			}
		});

		const downloadButton = this.rootElem.getElementsByClassName('download-button')[0] as HTMLElement;
		if (allowDownload) {
			downloadButton.addEventListener('click', event => {
				const data = this.textElem.textContent!;
				downloadString(data, 'wowsims.json');
			});
		} else {
			downloadButton.remove();
		}
	}

	protected init() {
		this.textElem.textContent = this.getData();
	}

	abstract getData(): string;
}

export class IndividualLinkExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, 'Sharable Link', false);
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		return this.simUI.toLink();
	}
}

export class IndividualJsonExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, 'JSON Export', true);
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		return JSON.stringify(IndividualSimSettings.toJson(this.simUI.toProto()), null, 2);
	}
}

export class IndividualWowheadGearPlannerExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, 'WoWHead Export', true);
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		const player = this.simUI.player;

		const classStr = classNames[player.getClass()].replace(' ', '-').toLowerCase();
		const raceStr = raceNames[player.getRace()].replace(' ', '-').toLowerCase();
		let url = `https://www.wowhead.com/wotlk/gear-planner/${classStr}/${raceStr}/`;

		// See comments on the importer for how the binary formatting is structured.
		let bytes: Array<number> = [];
		bytes.push(6);
		bytes.push(0);
		bytes.push(Mechanics.CHARACTER_LEVEL);

		let talentsStr = player.getTalentsString().replace('-', 'f') + 'f';
		if (talentsStr.length % 2 == 1) {
			talentsStr += '0';
		}
		bytes.push(talentsStr.length / 2);
		for (let i = 0; i < talentsStr.length; i += 2) {
			bytes.push(parseInt(talentsStr.substring(i, i + 2), 16));
		}

		//console.log('With talents: ' + btoa(String.fromCharCode(...bytes)).replaceAll('/', '_').replaceAll('+', '-'));

		//let glyphBytes: Array<number> = [];
		//let glyphStr = '';
		//const glyphs = player.getGlyphs();
		//const d = "0123456789abcdefghjkmnpqrstvwxyz";
		//const addGlyph = (glyphItemId: number, glyphPosition: number) => {
		//	const spellId = this.simUI.sim.db.glyphItemToSpellId(glyphItemId);
		//	if (!spellId) {
		//		return;
		//	}
		//	glyphStr += d[glyphPosition];
		//	glyphStr += d[(spellId >> 15) & 0b00011111];
		//	glyphStr += d[(spellId >> 10) & 0b00011111];
		//	glyphStr += d[(spellId >>  5) & 0b00011111];
		//	glyphStr += d[(spellId >>  0) & 0b00011111];
		//};
		//addGlyph(glyphs.major1, 0);
		//addGlyph(glyphs.major2, 1);
		//addGlyph(glyphs.major3, 2);
		//addGlyph(glyphs.minor1, 3);
		//addGlyph(glyphs.minor2, 4);
		//addGlyph(glyphs.minor3, 5);
		//if (glyphStr) {
		//	glyphBytes.push(0x30);
		//	for (let i = 0; i < glyphStr.length; i++) {
		//		glyphBytes.push(glyphStr.charCodeAt(i));
		//	}
		//}
		//bytes.push(glyphBytes.length);
		//bytes = bytes.concat(glyphBytes)
		bytes.push(0);

		const to2Bytes = (val: number): Array<number> => {
			return [
				(val & 0xff00) >> 8,
				val & 0x00ff,
			];
		};

		const gear = player.getGear();
		const isBlacksmithing = player.isBlacksmithing();
		gear.getItemSlots()
				.sort((slot1, slot2) => IndividualWowheadGearPlannerImporter.slotIDs[slot1] - IndividualWowheadGearPlannerImporter.slotIDs[slot2])
				.forEach(itemSlot => {
			const item = gear.getEquippedItem(itemSlot);
			if (!item) {
				return;
			}

			let slotId = IndividualWowheadGearPlannerImporter.slotIDs[itemSlot];
			if (item.enchant) {
				slotId = slotId | 0b10000000;
			}
			bytes.push(slotId);
			bytes.push(item.curGems(isBlacksmithing).length << 5);
			bytes = bytes.concat(to2Bytes(item.item.id));

			if (item.enchant) {
				bytes.push(0);
				bytes = bytes.concat(to2Bytes(item.enchant.spellId));
			}

			item.gems.slice(0, item.numSockets(isBlacksmithing)).forEach((gem, i) => {
				if (gem) {
					bytes.push(i << 5);
					bytes = bytes.concat(to2Bytes(gem.id));
				}
			});
		});

		const binaryString = String.fromCharCode(...bytes);
		const b64encoded = btoa(binaryString);
		const b64converted = b64encoded.replaceAll('/', '_').replaceAll('+', '-');

		return url + b64converted;
	}
}

export class Individual80UEPExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, '80Upgrades EP Export', true);
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		const player = this.simUI.player;
		const epValues = player.getEpWeights();
		const allUnitStats = UnitStat.getAll();

		const namesToWeights: Record<string, number> = {};
		allUnitStats
		.forEach(stat => {
			const statName = Individual80UEPExporter.getName(stat);
			const weight = epValues.getUnitStat(stat);
			if (weight == 0 || statName == '') {
				return;
			}

			// Need to add together stats with the same name (e.g. hit/crit/haste).
			if (namesToWeights[statName]) {
				namesToWeights[statName] += weight;
			} else {
				namesToWeights[statName] = weight;
			}
		});

		return `https://eightyupgrades.com/ep/import?name=${encodeURIComponent(`${specNames[player.spec]} WoWSims Weights`)}` +
			Object.keys(namesToWeights)
				.map(statName => `&${statName}=${namesToWeights[statName].toFixed(3)}`).join('');
	}

	static getName(stat: UnitStat): string {
		if (stat.isStat()) {
			return Individual80UEPExporter.statNames[stat.getStat()];
		} else {
			return Individual80UEPExporter.pseudoStatNames[stat.getPseudoStat()] || '';
		}
	}

	static statNames: Record<Stat, string> = {
		[Stat.StatStrength]: 'strength',
		[Stat.StatAgility]: 'agility',
		[Stat.StatStamina]: 'stamina',
		[Stat.StatIntellect]: 'intellect',
		[Stat.StatSpirit]: 'spirit',
		[Stat.StatSpellPower]: 'spellDamage',
		[Stat.StatMP5]: 'mp5',
		[Stat.StatSpellHit]: 'hitRating',
		[Stat.StatSpellCrit]: 'critRating',
		[Stat.StatSpellHaste]: 'hasteRating',
		[Stat.StatSpellPenetration]: 'spellPen',
		[Stat.StatAttackPower]: 'attackPower',
		[Stat.StatMeleeHit]: 'hitRating',
		[Stat.StatMeleeCrit]: 'critRating',
		[Stat.StatMeleeHaste]: 'hasteRating',
		[Stat.StatArmorPenetration]: 'armorPenRating',
		[Stat.StatExpertise]: 'expertiseRating',
		[Stat.StatMana]: 'mana',
		[Stat.StatEnergy]: 'energy',
		[Stat.StatRage]: 'rage',
		[Stat.StatArmor]: 'armor',
		[Stat.StatRangedAttackPower]: 'attackPower',
		[Stat.StatDefense]: 'defenseRating',
		[Stat.StatBlock]: 'blockRating',
		[Stat.StatBlockValue]: 'blockValue',
		[Stat.StatDodge]: 'dodgeRating',
		[Stat.StatParry]: 'parryRating',
		[Stat.StatResilience]: 'resilienceRating',
		[Stat.StatHealth]: 'health',
		[Stat.StatArcaneResistance]: 'arcaneResistance',
		[Stat.StatFireResistance]: 'fireResistance',
		[Stat.StatFrostResistance]: 'frostResistance',
		[Stat.StatNatureResistance]: 'natureResistance',
		[Stat.StatShadowResistance]: 'shadowResistance',
		[Stat.StatBonusArmor]: 'armorBonus',
	}
	static pseudoStatNames: Partial<Record<PseudoStat, string>> = {
		[PseudoStat.PseudoStatMainHandDps]: 'dps',
		[PseudoStat.PseudoStatRangedDps]: 'rangedDps',
	}
}

export class IndividualPawnEPExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, 'Pawn EP Export', true);
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		const player = this.simUI.player;
		const epValues = player.getEpWeights();
		const allUnitStats = UnitStat.getAll();

		const namesToWeights: Record<string, number> = {};
		allUnitStats
		.forEach(stat => {
			const statName = IndividualPawnEPExporter.getName(stat);
			const weight = epValues.getUnitStat(stat);
			if (weight == 0 || statName == '') {
				return;
			}

			// Need to add together stats with the same name (e.g. hit/crit/haste).
			if (namesToWeights[statName]) {
				namesToWeights[statName] += weight;
			} else {
				namesToWeights[statName] = weight;
			}
		});

		return `( Pawn: v1: "${specNames[player.spec]} WoWSims Weights": Class=${classNames[player.getClass()]},` +
			Object.keys(namesToWeights)
				.map(statName => `${statName}=${namesToWeights[statName].toFixed(3)}`).join(',') +
			' )';
	}

	static getName(stat: UnitStat): string {
		if (stat.isStat()) {
			return IndividualPawnEPExporter.statNames[stat.getStat()];
		} else {
			return IndividualPawnEPExporter.pseudoStatNames[stat.getPseudoStat()] || '';
		}
	}

	static statNames: Record<Stat, string> = {
		[Stat.StatStrength]: 'Strength',
		[Stat.StatAgility]: 'Agility',
		[Stat.StatStamina]: 'Stamina',
		[Stat.StatIntellect]: 'Intellect',
		[Stat.StatSpirit]: 'Spirit',
		[Stat.StatSpellPower]: 'SpellDamage',
		[Stat.StatMP5]: 'Mp5',
		[Stat.StatSpellHit]: 'HitRating',
		[Stat.StatSpellCrit]: 'CritRating',
		[Stat.StatSpellHaste]: 'HasteRating',
		[Stat.StatSpellPenetration]: 'SpellPen',
		[Stat.StatAttackPower]: 'Ap',
		[Stat.StatMeleeHit]: 'HitRating',
		[Stat.StatMeleeCrit]: 'CritRating',
		[Stat.StatMeleeHaste]: 'HasteRating',
		[Stat.StatArmorPenetration]: 'ArmorPenetration',
		[Stat.StatExpertise]: 'ExpertiseRating',
		[Stat.StatMana]: 'Mana',
		[Stat.StatEnergy]: 'Energy',
		[Stat.StatRage]: 'Rage',
		[Stat.StatArmor]: 'Armor',
		[Stat.StatRangedAttackPower]: 'Ap',
		[Stat.StatDefense]: 'DefenseRating',
		[Stat.StatBlock]: 'BlockRating',
		[Stat.StatBlockValue]: 'BlockValue',
		[Stat.StatDodge]: 'DodgeRating',
		[Stat.StatParry]: 'ParryRating',
		[Stat.StatResilience]: 'ResilienceRating',
		[Stat.StatHealth]: 'Health',
		[Stat.StatArcaneResistance]: 'ArcaneResistance',
		[Stat.StatFireResistance]: 'FireResistance',
		[Stat.StatFrostResistance]: 'FrostResistance',
		[Stat.StatNatureResistance]: 'NatureResistance',
		[Stat.StatShadowResistance]: 'ShadowResistance',
		[Stat.StatBonusArmor]: 'Armor2',
	}
	static pseudoStatNames: Partial<Record<PseudoStat, string>> = {
		[PseudoStat.PseudoStatMainHandDps]: 'MeleeDps',
		[PseudoStat.PseudoStatRangedDps]: 'RangedDps',
	}
}
