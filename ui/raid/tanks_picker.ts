import { Component } from '/tbc/core/components/component.js';
import { Input, InputConfig } from '/tbc/core/components/input.js';
import { RaidTargetPicker } from '/tbc/core/components/raid_target_picker.js';
import { Player } from '/tbc/core/player.js';
import { Raid } from '/tbc/core/raid.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { Class } from '/tbc/core/proto/common.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { emptyRaidTarget } from '/tbc/core/proto_utils/utils.js';

import { RaidSimUI } from './raid_sim_ui.js';

declare var tippy: any;

const MAX_TANKS = 4;

export class TanksPicker extends Component {
	readonly raidSimUI: RaidSimUI;

	constructor(parentElem: HTMLElement, raidSimUI: RaidSimUI) {
		super(parentElem, 'tanks-picker-root');
		this.raidSimUI = raidSimUI;

		this.rootElem.innerHTML = `
			<fieldset class="tanks-picker-container settings-section">
				<legend>TANKS</legend>
			</fieldset>
		`;

		const tanksContainer = this.rootElem.getElementsByClassName('tanks-picker-container')[0] as HTMLElement;
		const raid = this.raidSimUI.sim.raid;

		for (let i = 0; i < MAX_TANKS; i++) {
			const row = document.createElement('div');
			row.classList.add('tank-picker-row');
			tanksContainer.appendChild(row);

			const sourceElem = document.createElement('span');
			sourceElem.textContent = i == 0 ? 'MAIN TANK' : `TANK ${i + 1}`;
			sourceElem.classList.add('tank-picker-label');
			row.appendChild(sourceElem);

			const arrow = document.createElement('span');
			arrow.classList.add('fa', 'fa-arrow-right');
			row.appendChild(arrow);

			const tankIndex = i;
			const raidTargetPicker = new RaidTargetPicker<Raid>(row, raid, raid, {
				extraCssClasses: [
					'tank-picker',
				],
				noTargetLabel: 'Unassigned',
				compChangeEmitter: raid.compChangeEmitter,

				changedEvent: (raid: Raid) => raid.tanksChangeEmitter,
				getValue: (raid: Raid) => raid.getTanks()[tankIndex] || emptyRaidTarget(),
				setValue: (eventID: EventID, raid: Raid, newValue: RaidTarget) => {
					const tanks = raid.getTanks();
					for (let i = 0; i < tankIndex; i++) {
						if (!tanks[i]) {
							tanks.push(emptyRaidTarget());
						}
					}
					tanks[tankIndex] = newValue;
					raid.setTanks(eventID, tanks);
				},
			});
		}
	}
}
