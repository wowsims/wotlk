import { Player } from '../player.js';
import { ActionId } from '../proto_utils/action_id.js';
import { SimUI } from '../sim_ui.js';
import { TypedEvent } from '../typed_event.js';
import { isRightClick } from '../utils.js';

import { Component } from './component.js';
import { IconPicker, IconPickerConfig } from './icon_picker.js';

export interface MultiIconPickerItemConfig<ModObject> extends IconPickerConfig<ModObject, any> {
}

export interface MultiIconPickerConfig<ModObject> {
	inputs: Array<MultiIconPickerItemConfig<ModObject>>,
	label?: string,
	categoryId?: ActionId,
	showWhen?: (obj: Player<any>) => boolean,
}

// Icon-based UI for a dropdown with multiple icon pickers.
// ModObject is the object being modified (Sim, Player, or Target).
export class MultiIconPicker<ModObject> extends Component {
	private readonly config: MultiIconPickerConfig<ModObject>;

	private currentValue: ActionId | null;

	private readonly buttonElem: HTMLAnchorElement;
	private readonly dropdownMenu: HTMLElement;

	private readonly pickers: Array<IconPicker<ModObject, any>>;

	constructor(parent: HTMLElement, modObj: ModObject, config: MultiIconPickerConfig<ModObject>, simUI: SimUI) {
		super(parent, 'multi-icon-picker-root');
		this.rootElem.classList.add('icon-picker');
		this.config = config;
		this.currentValue = null;

		this.rootElem.innerHTML = `
			<div class="dropend">
				<a
					class="icon-picker-button"
					role="button"
					data-bs-toggle="dropdown"
					aria-expanded="false"
					data-disable-wowhead-touch-tooltip='true'
					data-whtticon='false'
				></a>
				<ul class="dropdown-menu"></ul>
			</div>
			<label class="multi-icon-picker-label form-label"></label>
    `;

		const labelElem = this.rootElem.querySelector('.multi-icon-picker-label') as HTMLElement;
		if (config.label) {
			labelElem.textContent = config.label;
		} else {
			labelElem.remove();
		}

		this.buttonElem = this.rootElem.querySelector('.icon-picker-button') as HTMLAnchorElement;
		this.dropdownMenu = this.rootElem.querySelector('.dropdown-menu') as HTMLElement;

		this.buttonElem.addEventListener('hide.bs.dropdown', event => {
			if (event.hasOwnProperty('clickEvent'))
				event.preventDefault();
		})

		this.buttonElem.addEventListener('contextmenu', event => {
			event.preventDefault();
		});

		this.buttonElem.addEventListener('mousedown', event => {
			const rightClick = isRightClick(event);

			if (rightClick) {
				this.clearPicker();
			}
		});

		this.buildBlankOption();

		this.pickers = config.inputs.map((pickerConfig, _i) => {
			const optionContainer = document.createElement('li');
			optionContainer.classList.add('icon-picker-option', 'dropdown-option');
			this.dropdownMenu.appendChild(optionContainer);

			return new IconPicker(optionContainer, modObj, pickerConfig);
		});
		simUI.sim.waitForInit().then(() => this.updateButtonImage());
		simUI.changeEmitter.on(() => this.updateButtonImage());
		simUI.changeEmitter.on(() => {
			const show = !this.config.showWhen || this.config.showWhen(simUI.sim.raid.getPlayer(0)!);
			if (show) {
				this.rootElem.classList.remove('hide');
			} else {
				this.rootElem.classList.add('hide');
			}
		})
	}

	private buildBlankOption() {
		const listItem = document.createElement('li');
		this.dropdownMenu.appendChild(listItem);

		const option = document.createElement('a');
		option.classList.add('icon-dropdown-option', 'dropdown-option');
		listItem.appendChild(option);

		option.addEventListener('click', () => this.clearPicker());
	}

	private clearPicker() {
		TypedEvent.freezeAllAndDo(() => {
			this.pickers.forEach((picker) => {
				picker.setInputValue(null)
				picker.inputChanged(TypedEvent.nextEventID());
			});
			this.updateButtonImage();
		})
	}

	private updateButtonImage() {
		this.currentValue = this.getMaxValue();

		if (this.currentValue) {
			this.buttonElem.classList.add('active');
			if (this.config.categoryId != null) {
				this.config.categoryId.fillAndSet(this.buttonElem, false, true);
			} else {
				this.currentValue.fillAndSet(this.buttonElem, false, true);
			}
		} else {
			this.buttonElem.classList.remove('active');
			if (this.config.categoryId != null) {
				this.config.categoryId.fillAndSet(this.buttonElem, false, true);
			} else {
				this.buttonElem.style.backgroundImage = '';
			}
			this.buttonElem.removeAttribute("href");
		}
	}

	private getMaxValue(): ActionId | null {
		return this.pickers.map(picker => picker.getActionId()).filter(id => id != null)[0] || null;
	}
}
