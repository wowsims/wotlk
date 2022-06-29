import { Sim } from '/tbc/core/sim.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

import { Component } from './component.js';

declare var tippy: any;

/**
 * Data for creating a new input UI element.
 */
export interface InputConfig<ModObject, T> {
	label?: string,
	labelTooltip?: string,
	extraCssClasses?: Array<string>,

	defaultValue?: T,

	// Returns the event indicating the mapped value has changed.
	changedEvent: (obj: ModObject) => TypedEvent<any>,

	// Get and set the mapped value.
	getValue: (obj: ModObject) => T,
	setValue: (eventID: EventID, obj: ModObject, newValue: T) => void,

	// If set, will automatically disable the input when this evaluates to false.
	enableWhen?: (obj: ModObject) => boolean,

	// If set, will automatically hide the input when this evaluates to false.
	showWhen?: (obj: ModObject) => boolean,

	// Overrides the default root element (new div).
	rootElem?: HTMLElement,
}

// Shared logic for UI elements that are mapped to a value for some modifiable object.
export abstract class Input<ModObject, T> extends Component {
	private readonly inputConfig: InputConfig<ModObject, T>;
	readonly modObject: ModObject;

	readonly changeEmitter = new TypedEvent<void>();

	constructor(parent: HTMLElement, cssClass: string, modObject: ModObject, config: InputConfig<ModObject, T>) {
		super(parent, 'input-root', config.rootElem);
		this.inputConfig = config;
		this.modObject = modObject;
		this.rootElem.classList.add(cssClass);
		if (config.extraCssClasses) {
			this.rootElem.classList.add(...config.extraCssClasses);
		}

		if (config.label) {
			const labelDiv = document.createElement('div');
			labelDiv.classList.add('input-label-div');
			this.rootElem.appendChild(labelDiv);

			const label = document.createElement('span');
			label.classList.add('input-label');
			label.textContent = config.label;
			labelDiv.appendChild(label);

			if (config.labelTooltip) {
				const tooltip = document.createElement('span');
				tooltip.classList.add('input-tooltip', 'fa', 'fa-info-circle');
				tippy(tooltip, {
					'content': config.labelTooltip,
					'allowHTML': true,
				});
				labelDiv.appendChild(tooltip);
			}
		}

		config.changedEvent(this.modObject).on(eventID => {
			this.setInputValue(config.getValue(this.modObject));
			this.update();
		});
	}

	private update() {
		const enable = !this.inputConfig.enableWhen || this.inputConfig.enableWhen(this.modObject);
		if (enable) {
			this.rootElem.classList.remove('disabled');
			this.getInputElem().removeAttribute('disabled');
		} else {
			this.rootElem.classList.add('disabled');
			this.getInputElem().setAttribute('disabled', '');
		}

		const show = !this.inputConfig.showWhen || this.inputConfig.showWhen(this.modObject);
		if (show) {
			this.rootElem.classList.remove('hide');
		} else {
			this.rootElem.classList.add('hide');
		}
	}

	// Can't call abstract functions in constructor, so need an init() call.
	init() {
		if (this.inputConfig.defaultValue) {
			this.setInputValue(this.inputConfig.defaultValue);
		} else {
			this.setInputValue(this.inputConfig.getValue(this.modObject));
		}
		this.update();
	}

	abstract getInputElem(): HTMLElement;

	abstract getInputValue(): T;

	abstract setInputValue(newValue: T): void;

	// Child classes should call this method when the value in the input element changes.
	inputChanged(eventID: EventID) {
		this.inputConfig.setValue(eventID, this.modObject, this.getInputValue());
		this.changeEmitter.emit(eventID);
	}

	// Sets the underlying value directly.
	setValue(eventID: EventID, newValue: T) {
		this.inputConfig.setValue(eventID, this.modObject, newValue);
	}
}
