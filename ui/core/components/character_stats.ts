import { PlayerStats } from '/tbc/core/proto/api.js';
import { Stat } from '/tbc/core/proto/common.js';
import { TristateEffect } from '/tbc/core/proto/common.js'
import { statNames, statOrder } from '/tbc/core/proto_utils/names.js';
import { Stats } from '/tbc/core/proto_utils/stats.js';
import { Player } from '/tbc/core/player.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

import * as Mechanics from '/tbc/core/constants/mechanics.js';

import { Component } from './component.js';

declare var tippy: any;

const spellPowerTypeStats = [
	Stat.StatArcaneSpellPower,
	Stat.StatFireSpellPower,
	Stat.StatFrostSpellPower,
	Stat.StatHolySpellPower,
	Stat.StatNatureSpellPower,
	Stat.StatShadowSpellPower,
];

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
					<span>${statNames[stat].toUpperCase()}<span>
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
			this.valueElems[idx].textContent = CharacterStats.statDisplayString(finalStats, stat);

			tippy(this.tooltipElems[idx], {
				'content': `
					<div class="character-stats-tooltip-row">
						<span>Base:</span>
						<span>${CharacterStats.statDisplayString(baseStats, stat)}</span>
					</div>
					<div class="character-stats-tooltip-row">
						<span>Gear:</span>
						<span>${CharacterStats.statDisplayString(gearDelta, stat)}</span>
					</div>
					<div class="character-stats-tooltip-row">
						<span>Talents:</span>
						<span>${CharacterStats.statDisplayString(talentsDelta, stat)}</span>
					</div>
					<div class="character-stats-tooltip-row">
						<span>Buffs:</span>
						<span>${CharacterStats.statDisplayString(buffsDelta, stat)}</span>
					</div>
					<div class="character-stats-tooltip-row">
						<span>Consumes:</span>
						<span>${CharacterStats.statDisplayString(consumesDelta, stat)}</span>
					</div>
					${debuffStats.getStat(stat) == 0 ? '' : `
					<div class="character-stats-tooltip-row">
						<span>Debuffs:</span>
						<span>${CharacterStats.statDisplayString(debuffStats, stat)}</span>
					</div>
					`}
					<div class="character-stats-tooltip-row">
						<span>Total:</span>
						<span>${CharacterStats.statDisplayString(finalStats, stat)}</span>
					</div>
				`,
				'allowHTML': true,
			});
		});
	}

	static statDisplayString(stats: Stats, stat: Stat): string {
		let rawValue = stats.getStat(stat);
		if (spellPowerTypeStats.includes(stat)) {
			rawValue = rawValue + stats.getStat(Stat.StatSpellPower);
		}
		let displayStr = String(Math.round(rawValue));

		if (stat == Stat.StatMeleeHit) {
			displayStr += ` (${(rawValue / Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE).toFixed(2)}%)`;
		} else if (stat == Stat.StatSpellHit) {
			displayStr += ` (${(rawValue / Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE).toFixed(2)}%)`;
		} else if (stat == Stat.StatMeleeCrit || stat == Stat.StatSpellCrit) {
			displayStr += ` (${(rawValue / Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE).toFixed(2)}%)`;
		} else if (stat == Stat.StatMeleeHaste || stat == Stat.StatSpellHaste) {
			displayStr += ` (${(rawValue / Mechanics.HASTE_RATING_PER_HASTE_PERCENT).toFixed(2)}%)`;
		} else if (stat == Stat.StatExpertise) {
			displayStr += ` (${(Math.floor(rawValue / Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION)).toFixed(0)})`;
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
		if (debuffs.faerieFire == TristateEffect.TristateEffectImproved) {
			debuffStats = debuffStats.addStat(Stat.StatMeleeHit, 3 * Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE);
		}
		if (debuffs.improvedSealOfTheCrusader) {
			debuffStats = debuffStats.addStat(Stat.StatMeleeCrit, 3 * Mechanics.MELEE_CRIT_RATING_PER_CRIT_CHANCE);
			debuffStats = debuffStats.addStat(Stat.StatSpellCrit, 3 * Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE);
		}

		return debuffStats;
	}
}
