import { Component } from '../core/components/component';
import { IconEnumPicker } from '../core/components/icon_enum_picker';

import { memeSpecs } from '../core/launched_sims';
import { EventID, TypedEvent } from '../core/typed_event';

import { Class, Spec } from '../core/proto/common';
import { Blessings } from '../core/proto/paladin';
import { BlessingsAssignments } from '../core/proto/ui';
import { ActionId } from '../core/proto_utils/action_id';
import {
	makeDefaultBlessings,
	classColors,
	naturalSpecOrder,
	specNames,
	titleIcons,
} from '../core/proto_utils/utils';

import { RaidSimUI } from './raid_sim_ui';
import { implementedSpecs } from './presets';
import { Tooltip } from 'bootstrap';

const MAX_PALADINS = 4;

export class BlessingsPicker extends Component {
	readonly simUI: RaidSimUI;
	readonly changeEmitter: TypedEvent<void> = new TypedEvent<void>();

	private readonly pickers: Array<Array<IconEnumPicker<this, Blessings>>> = [];

	private assignments: BlessingsAssignments;

	constructor(parentElem: HTMLElement, raidSimUI: RaidSimUI) {
		super(parentElem, 'blessings-picker-root');
		this.simUI = raidSimUI;
		this.assignments = BlessingsAssignments.clone(makeDefaultBlessings(4));

		const specs = naturalSpecOrder
			.filter(spec => implementedSpecs.includes(spec))
			.filter(spec => !memeSpecs.includes(spec));
		const paladinIndexes = [...Array(MAX_PALADINS).keys()];

		specs.map(spec => {
			const row = document.createElement('div');
			row.classList.add('blessings-picker-row');
			this.rootElem.appendChild(row);

			row.append(this.buildSpecIcon(spec));

			const container = document.createElement('div');
			container.classList.add('blessings-picker-container');
			row.appendChild(container);

			paladinIndexes.forEach(paladinIdx => {
				if (!this.pickers[paladinIdx])
					this.pickers.push([]);

				const blessingPicker = new IconEnumPicker(container, this, {
					extraCssClasses: ['blessing-picker'],
					numColumns: 1,
					values: [
						{ color: classColors[Class.ClassPaladin], value: Blessings.BlessingUnknown },
						{ actionId: ActionId.fromSpellId(25898), value: Blessings.BlessingOfKings },
						{ actionId: ActionId.fromSpellId(48934), value: Blessings.BlessingOfMight },
						{ actionId: ActionId.fromSpellId(48938), value: Blessings.BlessingOfWisdom },
						{ actionId: ActionId.fromSpellId(25899), value: Blessings.BlessingOfSanctuary },
					],
					equals: (a: Blessings, b: Blessings) => a == b,
					zeroValue: Blessings.BlessingUnknown,
					enableWhen: (_picker: BlessingsPicker) => {
						const numPaladins = Math.min(this.simUI.getClassCount(Class.ClassPaladin), MAX_PALADINS);
						return paladinIdx < numPaladins;
					},
					changedEvent: (picker: BlessingsPicker) => picker.changeEmitter,
					getValue: (picker: BlessingsPicker) => picker.assignments.paladins[paladinIdx]?.blessings[spec] || Blessings.BlessingUnknown,
					setValue: (eventID: EventID, picker: BlessingsPicker, newValue: number) => {
						const currentValue = picker.assignments.paladins[paladinIdx].blessings[spec];
						if (currentValue != newValue) {
							picker.assignments.paladins[paladinIdx].blessings[spec] = newValue;
							this.changeEmitter.emit(eventID);
						}
					},
				});

				this.pickers[paladinIdx].push(blessingPicker);
			});

			return row;
		});

		this.updatePickers()
		this.simUI.compChangeEmitter.on(_eventID => this.updatePickers());
	}

	private updatePickers() {
		for (let i = 0; i < MAX_PALADINS; i++) {
			this.pickers[i].forEach(picker => picker.update());
		}
	}

	private buildSpecIcon(spec: Spec): HTMLElement {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="blessings-picker-spec">
				<img
					src="${titleIcons[spec]}"
					class="blessings-spec-icon"
				/>
			</div>
		`;

		const icon = fragment.querySelector('.blessings-spec-icon') as HTMLElement;
		Tooltip.getOrCreateInstance(icon, { title: specNames[spec]});

		return fragment.children[0] as HTMLElement;
	}

	getAssignments(): BlessingsAssignments {
		// Defensive copy.
		return BlessingsAssignments.clone(this.assignments);
	}

	setAssignments(eventID: EventID, newAssignments: BlessingsAssignments) {
		this.assignments = BlessingsAssignments.clone(newAssignments);
		this.changeEmitter.emit(eventID);
	}
}
