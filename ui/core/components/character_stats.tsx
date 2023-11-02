import { Stat, Class, PseudoStat, Spec } from '..//proto/common.js';
import { TristateEffect } from '..//proto/common.js'
import { getClassStatName, statOrder } from '..//proto_utils/names.js';
import { Stats } from '..//proto_utils/stats.js';
import { Player } from '..//player.js';
import { EventID, TypedEvent } from '..//typed_event.js';

import * as Mechanics from '../constants/mechanics.js';

import { NumberPicker } from './number_picker';
import { Component } from './component.js';

import { Popover, Tooltip } from 'bootstrap';

import { element, fragment } from 'tsx-vanilla';

export type StatMods = { talents: Stats };

export class CharacterStats extends Component {
	readonly stats: Array<Stat>;
	readonly valueElems: Array<HTMLTableCellElement>;
	readonly meleeCritCapValueElem: HTMLTableCellElement | undefined;

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

			const row = (
			<tr
				className='character-stats-table-row'
			>
				<td className="character-stats-table-label">{statName}</td>
				<td className="character-stats-table-value">
					{this.bonusStatsLink(stat)}
				</td>
			</tr>);

			table.appendChild(row);

			const valueElem = row.getElementsByClassName('character-stats-table-value')[0] as HTMLTableCellElement;
			this.valueElems.push(valueElem);
		});

		if(this.shouldShowMeleeCritCap(player)) {
			const row = (
			<tr
				className='character-stats-table-row'
			>
				<td className="character-stats-table-label">Melee Crit Cap</td>
				<td className="character-stats-table-value"></td>
			</tr>);

			table.appendChild(row);
			this.meleeCritCapValueElem = row.getElementsByClassName('character-stats-table-value')[0] as HTMLTableCellElement;
		} else {
			this.meleeCritCapValueElem = undefined;
		}

		this.updateStats(player);
		TypedEvent.onAny([player.currentStatsEmitter, player.sim.changeEmitter, player.talentsChangeEmitter]).on(() => {
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

		const baseDelta = baseStats;
		const gearDelta = gearStats.subtract(baseStats).subtract(bonusStats);
		const talentsDelta = talentsStats.subtract(gearStats).add(statMods.talents);
		const buffsDelta = buffsStats.subtract(talentsStats);
		const consumesDelta = consumesStats.subtract(buffsStats);

		const finalStats = Stats.fromProto(playerStats.finalStats).add(statMods.talents).add(debuffStats);

		this.stats.forEach((stat, idx) => {
			let valueElem = (
				<a
					href="javascript:void(0)"
					className="stat-value-link"
					attributes={{role:"button"}}>
					{`${this.statDisplayString(finalStats, finalStats, stat)} `}
				</a>
			)

			this.valueElems[idx].querySelector('.stat-value-link')?.remove();
			this.valueElems[idx].prepend(valueElem);

			let bonusStatValue = bonusStats.getStat(stat);

			if (bonusStatValue == 0) {
				valueElem.classList.add('text-white');
			} else if (bonusStatValue > 0) {
				valueElem.classList.add('text-success');
			} else if (bonusStatValue < 0) {
				valueElem.classList.add('text-danger');
			}

			let tooltipContent =
			<div>
				<div className="character-stats-tooltip-row">
					<span>Base:</span>
					<span>{this.statDisplayString(baseStats, baseDelta, stat)}</span>
				</div>
				<div className="character-stats-tooltip-row">
					<span>Gear:</span>
					<span>{this.statDisplayString(gearStats, gearDelta, stat)}</span>
				</div>
				<div className="character-stats-tooltip-row">
					<span>Talents:</span>
					<span>{this.statDisplayString(talentsStats, talentsDelta, stat)}</span>
				</div>
				<div className="character-stats-tooltip-row">
					<span>Buffs:</span>
					<span>{this.statDisplayString(buffsStats, buffsDelta, stat)}</span>
				</div>
				<div className="character-stats-tooltip-row">
					<span>Consumes:</span>
					<span>{this.statDisplayString(consumesStats, consumesDelta, stat)}</span>
				</div>
				{debuffStats.getStat(stat) != 0 &&
				<div className="character-stats-tooltip-row">
					<span>Debuffs:</span>
					<span>{this.statDisplayString(debuffStats, debuffStats, stat)}</span>
				</div>
				}
				{bonusStatValue != 0 &&
				<div className="character-stats-tooltip-row">
					<span>Bonus:</span>
					<span>{this.statDisplayString(bonusStats, bonusStats, stat)}</span>
				</div>
				}
				<div className="character-stats-tooltip-row">
					<span>Total:</span>
					<span>{this.statDisplayString(finalStats, finalStats, stat)}</span>
				</div>
			</div>;
			Tooltip.getOrCreateInstance(valueElem, {
				title: tooltipContent,
				html: true,
			});
		});

		if(this.meleeCritCapValueElem) {
			const meleeCritCapInfo = player.getMeleeCritCapInfo();

			const valueElem = (
				<a
					href="javascript:void(0)"
					className="stat-value-link"
					attributes={{role:"button"}}>
					{`${this.meleeCritCapDisplayString(player, finalStats)} `}
				</a>
			)

			const capDelta = meleeCritCapInfo.playerCritCapDelta;
			if (capDelta == 0) {
				valueElem.classList.add('text-white');
			} else if (capDelta > 0) {
				valueElem.classList.add('text-danger');
			} else if (capDelta < 0) {
				valueElem.classList.add('text-success');
			}

			this.meleeCritCapValueElem.querySelector('.stat-value-link')?.remove();
			this.meleeCritCapValueElem.prepend(valueElem);

			const tooltipContent = (
			<div>
				<div className="character-stats-tooltip-row">
					<span>Glancing:</span>
					<span>{`${meleeCritCapInfo.glancing.toFixed(2)}%`}</span>
				</div>
				<div className="character-stats-tooltip-row">
					<span>Suppression:</span>
					<span>{`${meleeCritCapInfo.suppression.toFixed(2)}%`}</span>
				</div>
				<div className="character-stats-tooltip-row">
					<span>To Hit Cap:</span>
					<span>{`${meleeCritCapInfo.remainingMeleeHitCap.toFixed(2)}%`}</span>
				</div>
				<div className="character-stats-tooltip-row">
					<span>To Exp Cap:</span>
					<span>{`${meleeCritCapInfo.remainingExpertiseCap.toFixed(2)}%`}</span>
				</div>
				<div className="character-stats-tooltip-row">
					<span>Debuffs:</span>
					<span>{`${meleeCritCapInfo.debuffCrit.toFixed(2)}%`}</span>
				</div>
				{meleeCritCapInfo.specSpecificOffset != 0 &&
				<div className="character-stats-tooltip-row">
					<span>Spec Offsets:</span>
					<span>{`${meleeCritCapInfo.specSpecificOffset.toFixed(2)}%`}</span>
				</div>
				}
				<div className="character-stats-tooltip-row">
					<span>Final Crit Cap:</span>
					<span>{`${meleeCritCapInfo.baseCritCap.toFixed(2)}%`}</span>
				</div>
				<hr/>
				<div className="character-stats-tooltip-row">
					<span>Can Raise By:</span>
					<span>{`${(meleeCritCapInfo.remainingExpertiseCap + meleeCritCapInfo.remainingMeleeHitCap).toFixed(2)}%`}</span>
				</div>
			</div>
			);

			Tooltip.getOrCreateInstance(valueElem, {
				title: tooltipContent,
				html: true,
			});
		}
	}

	private statDisplayString(stats: Stats, deltaStats: Stats, stat: Stat): string {
		let rawValue = deltaStats.getStat(stat);

		if (stat == Stat.StatBlockValue) {
			rawValue *= stats.getPseudoStat(PseudoStat.PseudoStatBlockValueMultiplier) || 1;
		}

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
			// As of 06/20, Blizzard has changed Expertise to no longer truncate at quarter percent intervals. Note that
			// in-game character sheet tooltips will still display the truncated values, but it has been tested to behave
			// continuously in reality since the patch.
			displayStr += ` (${(rawValue / Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION / 4).toFixed(2)}%)`;
		} else if (stat == Stat.StatDefense) {
			displayStr += ` (${(Mechanics.CHARACTER_LEVEL * 5 + Math.floor(rawValue / Mechanics.DEFENSE_RATING_PER_DEFENSE)).toFixed(0)})`;
		} else if (stat == Stat.StatBlock) {
			// TODO: Figure out how to display these differently for the components than the final value
			//displayStr += ` (${(rawValue / Mechanics.BLOCK_RATING_PER_BLOCK_CHANCE).toFixed(2)}%)`;
			displayStr += ` (${((rawValue / Mechanics.BLOCK_RATING_PER_BLOCK_CHANCE) + (Mechanics.MISS_DODGE_PARRY_BLOCK_CRIT_CHANCE_PER_DEFENSE * Math.floor(stats.getStat(Stat.StatDefense) / Mechanics.DEFENSE_RATING_PER_DEFENSE)) + 5.00).toFixed(2)}%)`;
		} else if (stat == Stat.StatDodge) {
			//displayStr += ` (${(rawValue / Mechanics.DODGE_RATING_PER_DODGE_CHANCE).toFixed(2)}%)`;
			displayStr += ` (${(stats.getPseudoStat(PseudoStat.PseudoStatDodge) * 100).toFixed(2)}%)`;
		} else if (stat == Stat.StatParry) {
			//displayStr += ` (${(rawValue / Mechanics.PARRY_RATING_PER_PARRY_CHANCE).toFixed(2)}%)`;
			displayStr += ` (${(stats.getPseudoStat(PseudoStat.PseudoStatParry) * 100).toFixed(2)}%)`;
		} else if (stat == Stat.StatResilience) {
			displayStr += ` (${(rawValue / Mechanics.RESILIENCE_RATING_PER_CRIT_REDUCTION_CHANCE).toFixed(2)}%)`;
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

		let link = (
			<a
				href="javascript:void(0)"
				className='add-bonus-stats text-white ms-2'
				dataset={{bsToggle: 'popover'}}
				attributes={{role:'button'}}
			>
				<i className="fas fa-plus-minus"></i>
			</a>
		);

		let popover: Popover | null = null;

		let picker = new NumberPicker(null, this.player, {
			label: `Bonus ${statName}`,
			extraCssClasses: ['mb-0'],
			changedEvent: (player: Player<any>) => player.bonusStatsChangeEmitter,
			getValue: (player: Player<any>) => player.getBonusStats().getStat(stat),
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const bonusStats = player.getBonusStats().withStat(stat, newValue);
				player.setBonusStats(eventID, bonusStats);
				popover?.hide();
			},
		});

		popover = new Popover(link, {
			customClass: 'bonus-stats-popover',
			placement: 'right',
			fallbackPlacement: ['left'],
			sanitize: false,
			html:true,
			content: picker.rootElem,
		});

		return link as HTMLElement;
	}

	private shouldShowMeleeCritCap(player: Player<any>): boolean {
		return [
			Spec.SpecDeathknight,
			Spec.SpecEnhancementShaman,
			Spec.SpecRetributionPaladin,
			Spec.SpecRogue,
			Spec.SpecWarrior
		].includes(player.spec);
	}

	private meleeCritCapDisplayString(player: Player<any>, finalStats: Stats): string {
		const playerCritCapDelta = player.getMeleeCritCap();

		if(playerCritCapDelta === 0.0) {
			return 'Exact';
		}

		const prefix = playerCritCapDelta > 0 ? 'Over by ' : 'Under by ';
		return `${prefix} ${Math.abs(playerCritCapDelta).toFixed(2)}%`;
	}
}
