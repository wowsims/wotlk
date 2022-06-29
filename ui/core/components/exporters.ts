import { Class } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Race } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { Stat } from '/tbc/core/proto/common.js';
import { IndividualSimSettings } from '/tbc/core/proto/ui.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { Player } from '/tbc/core/player.js';
import { classNames, nameToClass, nameToRace } from '/tbc/core/proto_utils/names.js';
import { talentSpellIdsToTalentString } from '/tbc/core/talents/factory.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { downloadString, getEnumValues } from '/tbc/core/utils.js';

import { Popup } from './popup.js';

declare var $: any;
declare var tippy: any;
declare var pako: any;

export function newIndividualExporters<SpecType extends Spec>(simUI: IndividualSimUI<SpecType>): HTMLElement {
	const exportSettings = document.createElement('div');
	exportSettings.classList.add('export-settings', 'sim-dropdown-menu');
	exportSettings.innerHTML = `
		<span id="exportMenuLink" class="dropdown-toggle fas fa-file-export" role="button" data-toggle="dropdown" aria-haspopup="true" arai-expanded="false"></span>
		<div class="dropdown-menu dropdown-menu-right" aria-labelledby="exportMenuLink">
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
	addMenuItem('70U EP', () => new Individual70UEPExporter(menuElem, simUI), false);
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
				<textarea class="exporter-textarea" readonly></textarea>
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
		const proto = this.simUI.toProto();
		// When sharing links, people generally don't intend to share settings/ep weights.
		proto.settings = undefined;
		proto.epWeights = [];

		const protoBytes = IndividualSimSettings.toBinary(proto);
		const deflated = pako.deflate(protoBytes, { to: 'string' });
		const encoded = btoa(String.fromCharCode(...deflated));

		const linkUrl = new URL(window.location.href);
		linkUrl.hash = encoded;
		return linkUrl.toString();
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

class Individual70UEPExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, '70Upgrades EP Export', true);
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		const epValues = this.simUI.player.getEpWeights();
		const allStats = (getEnumValues(Stat) as Array<Stat>).filter(stat => ![Stat.StatEnergy, Stat.StatRage].includes(stat));
		return `https://seventyupgrades.com/ep/import?name=${encodeURIComponent('WoWSims Weights')}` +
			allStats
				.filter(stat => epValues.getStat(stat) != 0)
				.map(stat => `&${Individual70UEPExporter.linkNames[stat]}=${epValues.getStat(stat).toFixed(3)}`).join('');
	}

	static linkNames: Record<Stat, string> = {
		[Stat.StatStrength]: 'strength',
		[Stat.StatAgility]: 'agility',
		[Stat.StatStamina]: 'stamina',
		[Stat.StatIntellect]: 'intellect',
		[Stat.StatSpirit]: 'spirit',
		[Stat.StatSpellPower]: 'spellDamage',
		[Stat.StatHealingPower]: 'healing',
		[Stat.StatArcaneSpellPower]: 'arcaneDamage',
		[Stat.StatFireSpellPower]: 'fireDamage',
		[Stat.StatFrostSpellPower]: 'frostDamage',
		[Stat.StatHolySpellPower]: 'holyDamage',
		[Stat.StatNatureSpellPower]: 'natureDamage',
		[Stat.StatShadowSpellPower]: 'shadowDamage',
		[Stat.StatMP5]: 'mp5',
		[Stat.StatSpellHit]: 'spellHitRating',
		[Stat.StatSpellCrit]: 'spellCritRating',
		[Stat.StatSpellHaste]: 'spellHasteRating',
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
		[Stat.StatRangedAttackPower]: 'rangedAttackPower',
		[Stat.StatFeralAttackPower]: 'feralAttackPower',
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
		const epValues = this.simUI.player.getEpWeights();
		const allStats = (getEnumValues(Stat) as Array<Stat>).filter(stat => ![Stat.StatEnergy, Stat.StatRage].includes(stat));
		return `( Pawn: v1: "WoWSims Weights": Class=${classNames[this.simUI.player.getClass()]},` +
			allStats
				.filter(stat => epValues.getStat(stat) != 0)
				.map(stat => `${IndividualPawnEPExporter.statNames[stat]}=${epValues.getStat(stat).toFixed(3)}`).join(',') +
			' )';
	}

	static statNames: Record<Stat, string> = {
		[Stat.StatStrength]: 'Strength',
		[Stat.StatAgility]: 'Agility',
		[Stat.StatStamina]: 'Stamina',
		[Stat.StatIntellect]: 'Intellect',
		[Stat.StatSpirit]: 'Spirit',
		[Stat.StatSpellPower]: 'SpellDamage',
		[Stat.StatHealingPower]: 'Healing',
		[Stat.StatArcaneSpellPower]: 'ArcaneSpellDamage',
		[Stat.StatFireSpellPower]: 'FireSpellDamage',
		[Stat.StatFrostSpellPower]: 'FrostSpellDamage',
		[Stat.StatHolySpellPower]: 'HolySpellDamage',
		[Stat.StatNatureSpellPower]: 'NatureSpellDamage',
		[Stat.StatShadowSpellPower]: 'ShadowSpellDamage',
		[Stat.StatMP5]: 'Mp5',
		[Stat.StatSpellHit]: 'SpellHitRating',
		[Stat.StatSpellCrit]: 'SpellCritRating',
		[Stat.StatSpellHaste]: 'SpellHasteRating',
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
		[Stat.StatRangedAttackPower]: 'Rap',
		[Stat.StatFeralAttackPower]: 'FeralAp',
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
