import { Stat } from '../proto/common.js';
import { statNames, statOrder } from '../proto_utils/names.js';
import { Stats } from '../proto_utils/stats.js';
import { Player } from '../player.js';
import { EventID, TypedEvent } from '../typed_event.js';

import { Component } from './component.js';
import { NumberPicker } from './number_picker.js';

declare var tippy: any;

export class BonusStatsPicker extends Component {
	readonly stats: Array<Stat>;
	readonly statPickers: Array<NumberPicker<Player<any>>>;

	constructor(parent: HTMLElement, player: Player<any>, stats: Array<Stat>) {
		super(parent, 'bonus-stats-root');
		this.stats = stats;

		const label = document.createElement('span');
		label.classList.add('bonus-stats-label');
		label.textContent = 'Bonus Stats';
		tippy(label, {
			'content': 'Extra stats to add on top of gear, buffs, etc.',
			'allowHTML': true,
		});
		this.rootElem.appendChild(label);

		this.statPickers = statOrder.filter(stat => this.stats.includes(stat)).map(stat => new NumberPicker(this.rootElem, player, {
			label: statNames[stat],
			changedEvent: (player: Player<any>) => player.bonusStatsChangeEmitter,
			getValue: (player: Player<any>) => player.getBonusStats().getStat(stat),
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const bonusStats = player.getBonusStats().withStat(stat, newValue);
				player.setBonusStats(eventID, bonusStats);
			},
		}));

		player.bonusStatsChangeEmitter.on(() => {
			this.statPickers.forEach(statPicker => {
				if (statPicker.getInputValue() > 0) {
					statPicker.rootElem.classList.remove('negative');
					statPicker.rootElem.classList.add('positive');
				} else if (statPicker.getInputValue() < 0) {
					statPicker.rootElem.classList.remove('positive');
					statPicker.rootElem.classList.add('negative');
				} else {
					statPicker.rootElem.classList.remove('negative');
					statPicker.rootElem.classList.remove('positive');
				}
			});
		});
	}
}
