import { Component } from '../core/components/component';
import { UnitReferencePicker } from '../core/components/raid_target_picker';

import { Raid } from '../core/raid';
import { EventID } from '../core/typed_event';

import { UnitReference } from '../core/proto/common';
import { emptyUnitReference } from '../core/proto_utils/utils';

import { RaidSimUI } from './raid_sim_ui';

const MAX_TANKS = 4;

export class TanksPicker extends Component {
	readonly raidSimUI: RaidSimUI;

	constructor(parentElem: HTMLElement, raidSimUI: RaidSimUI) {
		super(parentElem, 'tanks-picker-root');
		this.raidSimUI = raidSimUI;

		const raid = this.raidSimUI.sim.raid;

		for (let i = 0; i < MAX_TANKS; i++) {
			const row = document.createElement('div');
			row.classList.add('tank-picker-row', 'input-inline');
			this.rootElem.appendChild(row);

			const labelElem = document.createElement('label');
			labelElem.textContent = i == 0 ? 'Main Tank' : `Tank ${i + 1}`;
			labelElem.classList.add('tank-picker-label', 'form-label');
			row.appendChild(labelElem);

			new UnitReferencePicker<Raid>(row, raid, raid, {
				extraCssClasses: ['tank-picker'],
				noTargetLabel: 'Unassigned',
				compChangeEmitter: raid.compChangeEmitter,

				changedEvent: (raid: Raid) => raid.tanksChangeEmitter,
				getValue: (raid: Raid) => raid.getTanks()[i] || emptyUnitReference(),
				setValue: (eventID: EventID, raid: Raid, newValue: UnitReference) => {
					const tanks = raid.getTanks();
					for (let j = 0; j < i; j++) {
						if (!tanks[j]) {
							tanks.push(emptyUnitReference());
						}
					}
					tanks[i] = newValue;
					raid.setTanks(eventID, tanks);
				},
			});
		}
	}
}
