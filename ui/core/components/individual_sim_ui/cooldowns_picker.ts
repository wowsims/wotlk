import { Component } from '../component.js';
import { IconEnumPicker, IconEnumValueConfig } from '../icon_enum_picker.js';
import { Input, InputConfig } from '../input.js';
import { NumberListPicker } from '../number_list_picker.js';
import { Player } from '../../player.js';
import { EventID, TypedEvent } from '../../typed_event.js';
import { ActionID as ActionIdProto, ItemSlot } from '../../proto/common.js';
import { Cooldowns } from '../../proto/common.js';
import { Cooldown } from '../../proto/common.js';
import { ActionId } from '../../proto_utils/action_id.js';
import { Class } from '../../proto/common.js';
import { Spec } from '../../proto/common.js';
import { getEnumValues } from '../../utils.js';
import { wait } from '../../utils.js';
import { Tooltip } from 'bootstrap';
import { NumberPicker } from '../number_picker.js';
import { Sim } from 'ui/core/sim.js';

export class CooldownsPicker extends Component {
	readonly player: Player<any>;

	private cooldownPickers: Array<HTMLElement>;

	constructor(parentElem: HTMLElement, player: Player<any>) {
		super(parentElem, 'cooldowns-picker-root');
		this.player = player;
		this.cooldownPickers = [];

		TypedEvent.onAny([this.player.currentStatsEmitter]).on(eventID => {
			this.update();
		});
		this.update();
	}

	private update() {
		this.rootElem.innerHTML = '';
		const cooldowns = this.player.getCooldowns().cooldowns;

		this.cooldownPickers = [];
		for (let i = 0; i < cooldowns.length + 1; i++) {
			const cooldown = cooldowns[i];

			const row = document.createElement('div');
			row.classList.add('cooldown-picker');
			if (i == cooldowns.length) {
				row.classList.add('add-cooldown-picker');
			}
			this.rootElem.appendChild(row);

			const actionPicker = this.makeActionPicker(row, i);

			const label = document.createElement('label');
			label.classList.add('cooldown-picker-label', 'form-label');
			if (cooldown && cooldown.id) {
				ActionId.fromProto(cooldown.id).fill(this.player.getRaidIndex()).then(filledId => label.textContent = filledId.name);
			}
			row.appendChild(label);

			const timingsPicker = this.makeTimingsPicker(row, i);

			let deleteButtonFragment = document.createElement('fragment');
			deleteButtonFragment.innerHTML = `
				<a
					href="javascript:void(0)"
					class="delete-cooldown link-danger"
					role="button"
					data-bs-toggle="tooltip"
					data-bs-title="Delete Cooldown"
				>
					<i class="fa fa-times fa-xl"></i>
				</a>
			`
			const deleteButton = deleteButtonFragment.children[0] as HTMLElement;
			const deleteButtonTooltip = Tooltip.getOrCreateInstance(deleteButton);
			deleteButton.addEventListener('click', event => {
				const newCooldowns = this.player.getCooldowns();
				newCooldowns.cooldowns.splice(i, 1);
				this.player.setCooldowns(TypedEvent.nextEventID(), newCooldowns);
				deleteButtonTooltip.hide();
			});
			row.appendChild(deleteButton);

			this.cooldownPickers.push(row);
		}

		this.addTrinketDesyncPicker(ItemSlot.ItemSlotTrinket1);
		this.addTrinketDesyncPicker(ItemSlot.ItemSlotTrinket2);
	}

	private addTrinketDesyncPicker(slot: ItemSlot) {
		const index = slot - ItemSlot.ItemSlotTrinket1 + 1;
		const picker = new NumberPicker(this.rootElem, this.player.sim, {
			label: `Desync Trinket ${index} (seconds)`,
			labelTooltip: ' Put the trinket on a cooldown before pull by re-equipping it. Must be between 0 and 30 seconds.',
			extraCssClasses: [
				'within-raid-sim-hide',
			],
			inline: true,
			changedEvent: (_: Sim) => this.player.cooldownsChangeEmitter,
			getValue: (_: Sim) => {
				const cooldowns = this.player.getCooldowns();
				return (slot == ItemSlot.ItemSlotTrinket1) ? cooldowns.desyncProcTrinket1Seconds : cooldowns.desyncProcTrinket2Seconds;
			},
			setValue: (eventID: EventID, _: Sim, newValue: number) => {
				if (newValue >= 0 && newValue <= 30) {
					const newCooldowns = this.player.getCooldowns();
					if (slot == ItemSlot.ItemSlotTrinket1) {
						newCooldowns.desyncProcTrinket1Seconds = newValue;
					} else {
						newCooldowns.desyncProcTrinket2Seconds = newValue
					}
					this.player.setCooldowns(eventID, newCooldowns);
				}
			},
			enableWhen: (sim: Sim) => {
				// TODO(Riotdog-GehennasEU): Only show if the slot is non-empty and the
				// trinket has a proc effect?
				return true;
			},
		});

		const pickerInput = picker.rootElem.querySelector('.number-picker-input') as HTMLInputElement;
		pickerInput.type = 'number';
		pickerInput.min = "0";
		pickerInput.max = "30";

		const validator = () => {
			if (!pickerInput.checkValidity()) {
				pickerInput.reportValidity();
			}
		};
		pickerInput.addEventListener('change', validator);
		pickerInput.addEventListener('focusout', validator);
	}

	private makeActionPicker(parentElem: HTMLElement, cooldownIndex: number): IconEnumPicker<Player<any>, ActionIdProto> {
		const availableCooldowns = this.player.getCurrentStats().cooldowns;

		const actionPicker = new IconEnumPicker<Player<any>, ActionIdProto>(parentElem, this.player, {
			extraCssClasses: [
				'cooldown-action-picker',
			],
			numColumns: 3,
			values: ([
				{ color: '#grey', value: ActionIdProto.create() },
			] as Array<IconEnumValueConfig<Player<any>, ActionIdProto>>).concat(availableCooldowns.map(cooldownAction => {
				return { actionId: ActionId.fromProto(cooldownAction), value: cooldownAction };
			})),
			equals: (a: ActionIdProto, b: ActionIdProto) => ActionIdProto.equals(a, b),
			zeroValue: ActionIdProto.create(),
			backupIconUrl: (value: ActionIdProto) => ActionId.fromProto(value),
			changedEvent: (player: Player<any>) => player.cooldownsChangeEmitter,
			getValue: (player: Player<any>) => player.getCooldowns().cooldowns[cooldownIndex]?.id || ActionIdProto.create(),
			setValue: (eventID: EventID, player: Player<any>, newValue: ActionIdProto) => {
				const newCooldowns = player.getCooldowns();

				while (newCooldowns.cooldowns.length < cooldownIndex) {
					newCooldowns.cooldowns.push(Cooldown.create());
				}
				newCooldowns.cooldowns[cooldownIndex] = Cooldown.create({
					id: newValue,
					timings: [],
				});

				player.setCooldowns(eventID, newCooldowns);
			},
		});
		return actionPicker;
	}

	private makeTimingsPicker(parentElem: HTMLElement, cooldownIndex: number): NumberListPicker<Player<any>> {
		const actionPicker = new NumberListPicker(parentElem, this.player, {
			extraCssClasses: [
				'cooldown-timings-picker',
			],
			placeholder: '20, 40, ...',
			changedEvent: (player: Player<any>) => player.cooldownsChangeEmitter,
			getValue: (player: Player<any>) => {
				return player.getCooldowns().cooldowns[cooldownIndex]?.timings || [];
			},
			setValue: (eventID: EventID, player: Player<any>, newValue: Array<number>) => {
				const newCooldowns = player.getCooldowns();
				newCooldowns.cooldowns[cooldownIndex].timings = newValue;
				player.setCooldowns(eventID, newCooldowns);
			},
			enableWhen: (player: Player<any>) => {
				const curCooldown = player.getCooldowns().cooldowns[cooldownIndex];
				return curCooldown && !ActionIdProto.equals(curCooldown.id, ActionIdProto.create());
			},
		});
		return actionPicker;
	}
}
