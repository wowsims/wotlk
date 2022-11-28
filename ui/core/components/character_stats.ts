import { PlayerStats } from '..//proto/api.js';
import { Stat, Class } from '..//proto/common.js';
import { TristateEffect } from '..//proto/common.js'
import { getClassStatName, statOrder } from '..//proto_utils/names.js';
import { Stats } from '..//proto_utils/stats.js';
import { Player } from '..//player.js';
import { EventID, TypedEvent } from '..//typed_event.js';

import * as Mechanics from '../constants/mechanics.js';

import { Component } from './component.js';

declare var tippy: any;

export type StatMods = { talents: Stats };

export class CharacterStats extends Component {
	readonly stats: Array<Stat>;
	readonly valueElems: Array<HTMLTableCellElement>;
	readonly tooltipElems: Array<HTMLElement>;

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
		this.tooltipElems = [];
		this.stats.forEach(stat => {
			const row = document.createElement('tr');
			row.classList.add('character-stats-table-row');
			row.innerHTML = `
				<td class="character-stats-table-label">
					<span>${getClassStatName(stat, player.getClass()).toUpperCase()}<span>
					<span class="character-stats-table-tooltip fas fa-search"></span>
				</td>
				<td class="character-stats-table-value"></td>
			`;
			table.appendChild(row);

			const valueElem = row.getElementsByClassName('character-stats-table-value')[0] as HTMLTableCellElement;
			this.valueElems.push(valueElem);

			const tooltipElem = row.getElementsByClassName('character-stats-table-tooltip')[0] as HTMLElement;
			this.tooltipElems.push(tooltipElem);
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

		const baseStats = new Stats(playerStats.baseStats);
		const gearStats = new Stats(playerStats.gearStats);
		const talentsStats = new Stats(playerStats.talentsStats);
		const buffsStats = new Stats(playerStats.buffsStats);
		const consumesStats = new Stats(playerStats.consumesStats);
		const debuffStats = CharacterStats.getDebuffStats(player);
		const finalStats = new Stats(playerStats.finalStats).add(statMods.talents).add(debuffStats);

		const gearDelta = gearStats.subtract(baseStats);
		const talentsDelta = talentsStats.subtract(gearStats).add(statMods.talents);
		const buffsDelta = buffsStats.subtract(talentsStats);
		const consumesDelta = consumesStats.subtract(buffsStats);

		this.stats.forEach((stat, idx) => {
			this.valueElems[idx].textContent = CharacterStats.statDisplayString(player, finalStats, stat);

			tippy(this.tooltipElems[idx], {
				'content': `
					<div class="character-stats-tooltip-row">
						<span>Base:</span>
						<span>${CharacterStats.statDisplayString(player, baseStats, stat)}</span>
					</div>
					<div class="character-stats-tooltip-row">
						<span>Gear:</span>
						<span>${CharacterStats.statDisplayString(player, gearDelta, stat)}</span>
					</div>
					<div class="character-stats-tooltip-row">
						<span>Talents:</span>
						<span>${CharacterStats.statDisplayString(player, talentsDelta, stat)}</span>
					</div>
					<div class="character-stats-tooltip-row">
						<span>Buffs:</span>
						<span>${CharacterStats.statDisplayString(player, buffsDelta, stat)}</span>
					</div>
					<div class="character-stats-tooltip-row">
						<span>Consumes:</span>
						<span>${CharacterStats.statDisplayString(player, consumesDelta, stat)}</span>
					</div>
					${debuffStats.getStat(stat) == 0 ? '' : `
					<div class="character-stats-tooltip-row">
						<span>Debuffs:</span>
						<span>${CharacterStats.statDisplayString(player, debuffStats, stat)}</span>
					</div>
					`}
					<div class="character-stats-tooltip-row">
						<span>Total:</span>
						<span>${CharacterStats.statDisplayString(player, finalStats, stat)}</span>
					</div>
				`,
				'allowHTML': true,
			});
		});
	}

	static statDisplayString(player: Player<any>, stats: Stats, stat: Stat): string {
		const rawValue = stats.getStat(stat);
		let displayStr = String(Math.round(rawValue));

		if (stat == Stat.StatMeleeHit) {
			displayStr += ` (${(rawValue / Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE).toFixed(2)}%)`;
		} else if (stat == Stat.StatSpellHit) {
			displayStr += ` (${(rawValue / Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE).toFixed(2)}%)`;
		} else if (stat == Stat.StatMeleeCrit || stat == Stat.StatSpellCrit) {
			displayStr += ` (${(rawValue / Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE).toFixed(2)}%)`;
		} else if (stat == Stat.StatMeleeHaste) {
			if ([Class.ClassDruid, Class.ClassShaman, Class.ClassPaladin, Class.ClassDeathknight].includes(player.getClass())) {
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

	static getDebuffStats(player: Player<any>): Stats {
		let debuffStats = new Stats();

		const debuffs = player.sim.raid.getDebuffs();
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
}
