import { Component } from '/tbc/core/components/component.js';
import { IconEnumPicker, IconEnumValueConfig } from '/tbc/core/components/icon_enum_picker.js';
import { Input, InputConfig } from '/tbc/core/components/input.js';
import { NumberListPicker } from '/tbc/core/components/number_list_picker.js';
import { Player } from '/tbc/core/player.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { ActionID as ActionIdProto } from '/tbc/core/proto/common.js';
import { Cooldowns } from '/tbc/core/proto/common.js';
import { Cooldown } from '/tbc/core/proto/common.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { Class } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { getEnumValues } from '/tbc/core/utils.js';
import { wait } from '/tbc/core/utils.js';

declare var tippy: any;

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

			const deleteButton = document.createElement('span');
			deleteButton.classList.add('delete-cooldown', 'fa', 'fa-times');
			deleteButton.addEventListener('click', event => {
				const newCooldowns = this.player.getCooldowns();
				newCooldowns.cooldowns.splice(i, 1);
				this.player.setCooldowns(TypedEvent.nextEventID(), newCooldowns);
			});
			row.appendChild(deleteButton);

			const actionPicker = this.makeActionPicker(row, i);

			const label = document.createElement('span');
			label.classList.add('cooldown-picker-label');
			if (cooldown && cooldown.id) {
				ActionId.fromProto(cooldown.id).fill(this.player.getRaidIndex()).then(filledId => label.textContent = filledId.name);
			}
			row.appendChild(label);

			const timingsPicker = this.makeTimingsPicker(row, i);

			this.cooldownPickers.push(row);
		}
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
