import { Spec } from '../../proto/common.js';
import { CustomRotation, CustomSpell } from '../../proto/common.js';
import { EventID } from '../../typed_event.js';
import { Player } from '../../player.js';
import { IconEnumPicker, IconEnumValueConfig } from '../icon_enum_picker.js';
import { ListItemPickerConfig, ListPicker } from '../list_picker.js';
import { NumberPicker } from '../number_picker.js';

import { Component } from '../component.js';
import { Input } from '../input.js';
import { SimUI } from 'ui/core/sim_ui.js';

export interface CustomRotationPickerConfig<SpecType extends Spec, T> {
	getValue: (player: Player<SpecType>) => CustomRotation,
	setValue: (eventID: EventID, player: Player<SpecType>, newValue: CustomRotation) => void,

	numColumns: number,
	extraCssClasses?: string[];
	showCastsPerMinute?: boolean,
	values: Array<IconEnumValueConfig<Player<SpecType>, T>>;

	showWhen?: (player: Player<SpecType>) => boolean,
}

export class CustomRotationPicker<SpecType extends Spec, T> extends Component {
	constructor(parent: HTMLElement, simUI: SimUI, modPlayer: Player<SpecType>, config: CustomRotationPickerConfig<SpecType, T>) {
		super(parent, 'custom-rotation-picker-root');

		if (config.extraCssClasses)
			this.rootElem.classList.add(...config.extraCssClasses);

		new ListPicker<Player<SpecType>, CustomSpell>(this.rootElem, modPlayer, {
			extraCssClasses: ['custom-spells-picker'],
			title: 'Spell Priority',
			titleTooltip: 'Spells at the top of the list are prioritized first. Safely ignores untalented options.',
			itemLabel: 'Spell',
			changedEvent: (player: Player<SpecType>) => player.changeEmitter,
			getValue: (player: Player<SpecType>) => config.getValue(player).spells,
			setValue: (eventID: EventID, player: Player<SpecType>, newValue: Array<CustomSpell>) => {
				config.setValue(eventID, player, CustomRotation.create({
					spells: newValue,
				}));
			},
			newItem: () => CustomSpell.create(),
			copyItem: (oldItem: CustomSpell) => CustomSpell.clone(oldItem),
			newItemPicker: (parent: HTMLElement, listPicker: ListPicker<Player<SpecType>, CustomSpell>, index: number, itemConfig: ListItemPickerConfig<Player<SpecType>, CustomSpell>) => new CustomSpellPicker(parent, modPlayer, index, itemConfig, listPicker, config),
			inlineMenuBar: true,
			showWhen: config.showWhen,
		});
	}
}

class CustomSpellPicker<SpecType extends Spec, T> extends Input<Player<SpecType>, CustomSpell> {
	private readonly player: Player<SpecType>;
	private readonly listPicker: ListPicker<Player<SpecType>, CustomSpell>;
	private readonly spellIndex: number;

	private readonly spellPicker: Input<null, number>;
	private readonly cpmPicker: Input<null, number> | null;

	getSpell(): CustomSpell {
		return this.listPicker.config.getValue(this.player)[this.spellIndex] || CustomSpell.create();
	}

	constructor(parent: HTMLElement, player: Player<SpecType>, spellIndex: number, config: ListItemPickerConfig<Player<SpecType>, CustomSpell>, listPicker: ListPicker<Player<SpecType>, CustomSpell>, crConfig: CustomRotationPickerConfig<SpecType, T>) {
		super(parent, 'custom-spell-picker-root', player, config);
		this.player = player;
		this.listPicker = listPicker;
		this.spellIndex = spellIndex;

		this.spellPicker = new IconEnumPicker<null, number>(this.rootElem, null, {
			numColumns: crConfig.numColumns,
			values: crConfig.values.map(value => {
				if (value.showWhen) {
					const oldShowWhen = value.showWhen;
					value.showWhen = (() => oldShowWhen(player)) as unknown as ((player: Player<SpecType>) => boolean);
				}
				return value;
			}) as unknown as Array<IconEnumValueConfig<null, number>>,
			equals: (a: number, b: number) => a == b,
			zeroValue: 0,
			changedEvent: () => player.changeEmitter,
			getValue: () => this.getSpell().spell,
			setValue: (eventID: EventID, _: null, newValue: number) => {
				const spell = this.getSpell();
				spell.spell = newValue;
				this.setSpell(eventID, spell);
			},
		});

		this.cpmPicker = null;
		if (crConfig.showCastsPerMinute) {
			this.cpmPicker = new NumberPicker<null>(this.rootElem, null, {
				label: 'CPM',
				labelTooltip: 'Desired Casts-Per-Minute for this spell.',
				float: true,
				positive: true,
				changedEvent: () => player.changeEmitter,
				getValue: () => this.getSpell().castsPerMinute,
				setValue: (eventID: EventID, _: null, newValue: number) => {
					const spell = this.getSpell();
					spell.castsPerMinute = newValue;
					this.setSpell(eventID, spell);
				},
			});
		}
	}

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

	getInputValue(): CustomSpell {
		return CustomSpell.create({
			spell: this.spellPicker.getInputValue(),
			castsPerMinute: this.cpmPicker?.getInputValue(),
		});
	}

	setInputValue(newValue: CustomSpell) {
		if (!newValue) {
			return;
		}
		this.spellPicker.setInputValue(newValue.spell);
		if (this.cpmPicker) {
			this.cpmPicker.setInputValue(newValue.castsPerMinute);
		}
	}

	private setSpell(eventID: EventID, spell: CustomSpell) {
		const customSpells = this.listPicker.config.getValue(this.player);
		customSpells[this.spellIndex] = CustomSpell.clone(spell);
		this.listPicker.config.setValue(eventID, this.player, customSpells);
	}
}
