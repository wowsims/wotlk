import { Stat, Class } from '..//proto/common.js';
import { TristateEffect } from '..//proto/common.js'
import { getClassStatName, statOrder } from '..//proto_utils/names.js';
import { Stats } from '..//proto_utils/stats.js';
import { Player } from '..//player.js';
import { EventID, TypedEvent } from '..//typed_event.js';

import * as Mechanics from '../constants/mechanics.js';

import { NumberPicker } from './number_picker';
import { Component } from './component.js';

import { Popover, Tooltip } from 'bootstrap';

declare var tippy: any;

export type StatMods = { talents: Stats };

export class CharacterStats extends Component {
	readonly stats: Array<Stat>;
	readonly valueElems: Array<HTMLTableCellElement>;

	private readonly player: Player<any>;
	private readonly modifyDisplayStats?: (player: Player<any>) => StatMods;

	constructor(parent: HTMLElement, player: Player<any>, stats: Array<Stat>, modifyDisplayStats?: (player: Player<any>) => StatMods) {
		super(parent, 'character-stats-root');
		this.stats = statOrder.filter(stat => stats.includes(stat));
		this.player = player;
		this.modifyDisplayStats = modifyDisplayStats;

		const label = document.createElement('label');
		label.classList.add('character-stats-label');
		label.textContent = 'Stats';
		this.rootElem.appendChild(label);

		const table = document.createElement('table');
		table.classList.add('character-stats-table');
		this.rootElem.appendChild(table);

		this.valueElems = [];
		this.stats.forEach(stat => {
			let statName = getClassStatName(stat, player.getClass());

			const row = document.createElement('tr');
			row.classList.add('character-stats-table-row');
			row.innerHTML = `
				<td class="character-stats-table-label">${statName}</td>
				<td class="character-stats-table-value"></td
			`;
			table.appendChild(row);

			const valueElem = row.getElementsByClassName('character-stats-table-value')[0] as HTMLTableCellElement;
			valueElem.appendChild(this.bonusStatsLink(stat));
			this.valueElems.push(valueElem);
		});

		this.updateStats(player);
		TypedEvent.onAny([player.currentStatsEmitter, player.sim.changeEmitter]).on(() => {
			this.updateStats(player);
		});
	}

	private updateStats(player: Player<any>) {
		const playerStats = player.getCurrentStats();

		const statMods = this.modifyDisplayStats ? this.modifyDisplayStats(this.player) : {
			talents: new Stats(),
		};

		const baseStats = Stats.fromProto(playerStats.baseStats);
		const gearStats = Stats.fromProto(playerStats.gearStats);
		const talentsStats = Stats.fromProto(playerStats.talentsStats);
		const buffsStats = Stats.fromProto(playerStats.buffsStats);
		const consumesStats = Stats.fromProto(playerStats.consumesStats);
		const debuffStats = this.getDebuffStats();
		const bonusStats = player.getBonusStats();

		const baseDelta = baseStats.subtract(bonusStats);
		const gearDelta = gearStats.subtract(baseStats);
		const talentsDelta = talentsStats.subtract(gearStats).add(statMods.talents);
		const buffsDelta = buffsStats.subtract(talentsStats);
		const consumesDelta = consumesStats.subtract(buffsStats);

		const finalStats = Stats.fromProto(playerStats.finalStats).add(statMods.talents).add(debuffStats);

		this.stats.forEach((stat, idx) => {
			let fragment = document.createElement('fragment');
			fragment.innerHTML = `
				<a href="javascript:void(0)" class="stat-value-link" role="button" data-bs-toggle="tooltip" data-bs-html="true">${this.statDisplayString(finalStats, stat)}</a>
			`
			let valueElem = fragment.children[0] as HTMLElement;
			this.valueElems[idx].querySelector('.stat-value-link')?.remove()
			this.valueElems[idx].prepend(valueElem);

			let bonusStatValue = player.getBonusStats().getStat(stat);
			
			if (bonusStatValue == 0) {
				valueElem.classList.remove('text-success', 'text-danger');
				valueElem.classList.add('text-white');
			} else if (bonusStatValue > 0) {
				valueElem.classList.remove('text-white', 'text-danger');
				valueElem.classList.add('text-success');
			} else if (bonusStatValue < 0) {
				valueElem.classList.remove('text-white', 'text-success');
				valueElem.classList.add('text-danger');
			}

			valueElem.setAttribute('data-bs-title', `
				<div class="character-stats-tooltip-row">
					<span>Base:</span>
					<span>${this.statDisplayString(baseDelta, stat)}</span>
				</div>
				<div class="character-stats-tooltip-row">
					<span>Gear:</span>
					<span>${this.statDisplayString(gearDelta, stat)}</span>
				</div>
				<div class="character-stats-tooltip-row">
					<span>Talents:</span>
					<span>${this.statDisplayString(talentsDelta, stat)}</span>
				</div>
				<div class="character-stats-tooltip-row">
					<span>Buffs:</span>
					<span>${this.statDisplayString(buffsDelta, stat)}</span>
				</div>
				<div class="character-stats-tooltip-row">
					<span>Consumes:</span>
					<span>${this.statDisplayString(consumesDelta, stat)}</span>
				</div>
				${debuffStats.getStat(stat) == 0 ? '' : `
				<div class="character-stats-tooltip-row">
					<span>Debuffs:</span>
					<span>${this.statDisplayString(debuffStats, stat)}</span>
				</div>
				`}
				${bonusStatValue == 0 ? '' : `
				<div class="character-stats-tooltip-row">
					<span>Bonus:</span>
					<span>${this.statDisplayString(this.player.getBonusStats(), stat)}</span>
				</div>
				`}
				<div class="character-stats-tooltip-row">
					<span>Total:</span>
					<span>${this.statDisplayString(finalStats, stat)}</span>
				</div>
			`);

			Tooltip.getOrCreateInstance(valueElem);
		});
	}

	private statDisplayString(stats: Stats, stat: Stat): string {
		const rawValue = stats.getStat(stat);
		let displayStr = String(Math.round(rawValue));

		if (stat == Stat.StatMeleeHit) {
			displayStr += ` (${(rawValue / Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE).toFixed(2)}%)`;
		} else if (stat == Stat.StatSpellHit) {
			displayStr += ` (${(rawValue / Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE).toFixed(2)}%)`;
		} else if (stat == Stat.StatMeleeCrit || stat == Stat.StatSpellCrit) {
			displayStr += ` (${(rawValue / Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE).toFixed(2)}%)`;
		} else if (stat == Stat.StatMeleeHaste) {
			if ([Class.ClassDruid, Class.ClassShaman, Class.ClassPaladin, Class.ClassDeathknight].includes(this.player.getClass())) {
				displayStr += ` (${(rawValue / Mechanics.SPECIAL_MELEE_HASTE_RATING_PER_HASTE_PERCENT).toFixed(2)}%)`;
			} else {
				displayStr += ` (${(rawValue / Mechanics.HASTE_RATING_PER_HASTE_PERCENT).toFixed(2)}%)`;
			}
		} else if (stat == Stat.StatSpellHaste) {
			displayStr += ` (${(rawValue / Mechanics.HASTE_RATING_PER_HASTE_PERCENT).toFixed(2)}%)`;
		} else if (stat == Stat.StatArmorPenetration) {
			displayStr += ` (${(rawValue / Mechanics.ARMOR_PEN_PER_PERCENT_ARMOR).toFixed(2)}%)`;
		} else if (stat == Stat.StatExpertise) {
			displayStr += ` (${(Math.floor(rawValue / Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION) / 4).toFixed(2)}%)`;
		} else if (stat == Stat.StatDefense) {
			displayStr += ` (${(Mechanics.CHARACTER_LEVEL * 5 + rawValue / Mechanics.DEFENSE_RATING_PER_DEFENSE).toFixed(1)})`;
		} else if (stat == Stat.StatBlock) {
			displayStr += ` (${(rawValue / Mechanics.BLOCK_RATING_PER_BLOCK_CHANCE).toFixed(2)}%)`;
		} else if (stat == Stat.StatDodge) {
			displayStr += ` (${(rawValue / Mechanics.DODGE_RATING_PER_DODGE_CHANCE).toFixed(2)}%)`;
		} else if (stat == Stat.StatParry) {
			displayStr += ` (${(rawValue / Mechanics.PARRY_RATING_PER_PARRY_CHANCE).toFixed(2)}%)`;
		}

		return displayStr;
	}

	private getDebuffStats(): Stats {
		let debuffStats = new Stats();

		const debuffs = this.player.sim.raid.getDebuffs();
		if (debuffs.misery || debuffs.faerieFire == TristateEffect.TristateEffectImproved) {
			debuffStats = debuffStats.addStat(Stat.StatSpellHit, 3 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);
		}
		if (debuffs.totemOfWrath || debuffs.heartOfTheCrusader || debuffs.masterPoisoner) {
			debuffStats = debuffStats.addStat(Stat.StatSpellCrit, 3 * Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE);
			debuffStats = debuffStats.addStat(Stat.StatMeleeCrit, 3 * Mechanics.MELEE_CRIT_RATING_PER_CRIT_CHANCE);
		}
		if (debuffs.improvedScorch || debuffs.wintersChill || debuffs.shadowMastery) {
			debuffStats = debuffStats.addStat(Stat.StatSpellCrit, 5 * Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE);
		}

		return debuffStats;
	}

	private bonusStatsLink(stat: Stat): HTMLElement {
		let statName = getClassStatName(stat, this.player.getClass());
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<a
				href="javascript:void(0)"
				class="add-bonus-stats text-white ms-2"
				role="button"
				data-bs-toggle="popover"
				data-bs-content="
					<div class='input-root number-picker-root'>
						<label class='form-label'>Bonus Health</label>
						<input type='text' class='form-control number-picker-input' value=${this.player.getBonusStats().getStat(stat)}>
					</div>
				"
				data-bs-placement="right"
				data-bs-html="true"
			>
				<i class="fas fa-plus-minus"></i>
			</a>
		`;

		let link = fragment.children[0] as HTMLElement;
		let popover = Popover.getOrCreateInstance(link, {
			customClass: 'bonus-stats-popover',
			fallbackPlacement: ['left'],
			sanitize: false,
		});

		link.addEventListener('shown.bs.popover', (event) => {
			let popoverBody = document.querySelector('.popover.bonus-stats-popover .popover-body') as HTMLElement;
			popoverBody.innerHTML = '';
			let picker = new NumberPicker(popoverBody, this.player, {
				label: `Bonus ${statName}`,
				changedEvent: (player: Player<any>) => player.bonusStatsChangeEmitter,
				getValue: (player: Player<any>) => player.getBonusStats().getStat(stat),
				setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
					const bonusStats = player.getBonusStats().withStat(stat, newValue);
					player.setBonusStats(eventID, bonusStats);
					popover.hide();
				},
			});
		});

		return link as HTMLElement;
	}
}
