import { Class } from '../proto/common.js';
import { EquipmentSpec } from '../proto/common.js';
import { ItemSpec } from '../proto/common.js';
import { Race } from '../proto/common.js';
import { Spec } from '../proto/common.js';
import { Stat } from '../proto/common.js';
import { IndividualSimSettings } from '../proto/ui.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { Player } from '../player.js';
import { classNames, nameToClass, nameToRace } from '../proto_utils/names.js';
import { specNames } from '../proto_utils/utils.js';
import { talentSpellIdsToTalentString } from '../talents/factory.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { downloadString, getEnumValues } from '../utils.js';

import { Popup } from './popup.js';

declare var $: any;
declare var tippy: any;
declare var pako: any;

export function newIndividualExporters<SpecType extends Spec>(simUI: IndividualSimUI<SpecType>): HTMLElement {
	const exportSettings = document.createElement('div');
	exportSettings.classList.add('export-settings', 'sim-dropdown-menu');
	exportSettings.innerHTML = `
		<span id="exportMenuLink" class="dropdown-toggle fas fa-file-export" role="button" data-bs-toggle="dropdown" aria-haspopup="true" aria-expanded="false"></span>
		<div class="dropdown-menu dropdown-menu-end" aria-labelledby="exportMenuLink">
		</div>
	`;
	const linkElem = exportSettings.getElementsByClassName('dropdown-toggle')[0] as HTMLElement;
	tippy(linkElem, {
		'content': 'Export',
		'allowHTML': true,
	});

	const menuElem = exportSettings.getElementsByClassName('dropdown-menu')[0] as HTMLElement;
	const addMenuItem = (label: string, onClick: () => void, showInRaidSim: boolean) => {
		const itemElem = document.createElement('span');
		itemElem.classList.add('dropdown-item');
		if (!showInRaidSim) {
			itemElem.classList.add('within-raid-sim-hide');
		}
		itemElem.textContent = label;
		itemElem.addEventListener('click', onClick);
		menuElem.appendChild(itemElem);
	};

	addMenuItem('Link', () => new IndividualLinkExporter(menuElem, simUI), false);
	addMenuItem('Json', () => new IndividualJsonExporter(menuElem, simUI), true);
	addMenuItem('80U EP', () => new Individual80UEPExporter(menuElem, simUI), false);
	addMenuItem('Pawn EP', () => new IndividualPawnEPExporter(menuElem, simUI), false);

	return exportSettings;
}

export abstract class Exporter extends Popup {
	private readonly textElem: HTMLElement;

	constructor(parent: HTMLElement, title: string, allowDownload: boolean) {
		super(parent);

		this.rootElem.classList.add('exporter');
		this.rootElem.innerHTML = `
			<span class="exporter-title">${title}</span>
			<div class="export-content">
				<textarea class="exporter-textarea form-control" readonly></textarea>
			</div>
			<div class="actions-row">
				<button class="exporter-button sim-button clipboard-button">COPY TO CLIPBOARD</button>
				<button class="exporter-button sim-button download-button">DOWNLOAD</button>
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

class IndividualLinkExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, 'Sharable Link', false);
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		return this.simUI.toLink();
	}
}

class IndividualJsonExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, 'JSON Export', true);
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		return JSON.stringify(IndividualSimSettings.toJson(this.simUI.toProto()), null, 2);
	}
}

class Individual80UEPExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, '80Upgrades EP Export', true);
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		const player = this.simUI.player;
		const epValues = player.getEpWeights();
		const allStats = (getEnumValues(Stat) as Array<Stat>).filter(stat => ![Stat.StatEnergy, Stat.StatRage].includes(stat));

		const namesToWeights: Record<string, number> = {};
		allStats
		.filter(stat => epValues.getStat(stat) != 0)
		.forEach(stat => {
			const statName = Individual80UEPExporter.linkNames[stat];
			const weight = epValues.getStat(stat);

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

	static linkNames: Record<Stat, string> = {
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
		[Stat.StatArmorPenetration]: 'armorPen',
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
	}
}

class IndividualPawnEPExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, 'Pawn EP Export', true);
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		const player = this.simUI.player;
		const epValues = player.getEpWeights();
		const allStats = (getEnumValues(Stat) as Array<Stat>).filter(stat => ![Stat.StatEnergy, Stat.StatRage].includes(stat));

		const namesToWeights: Record<string, number> = {};
		allStats
		.filter(stat => epValues.getStat(stat) != 0)
		.forEach(stat => {
			const statName = IndividualPawnEPExporter.statNames[stat];
			const weight = epValues.getStat(stat);

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
	}
}
