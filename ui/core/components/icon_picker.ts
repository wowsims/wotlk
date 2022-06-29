import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { isRightClick } from '/tbc/core/utils.js';

import { Component } from './component.js';
import { Input, InputConfig } from './input.js';

// Data for creating an icon-based input component.
// 
// E.g. one of these for arcane brilliance, another for kings, etc.
// ModObject is the object being modified (Sim, Player, or Target).
// ValueType is either number or boolean.
export interface IconPickerConfig<ModObject, ValueType> extends InputConfig<ModObject, ValueType> {
	id: ActionId;

	// The number of possible 'states' this icon can have. Most inputs will use 2
	// for a bi-state icon (on or off). 0 indicates an unlimited number of states.
	states: number;

	// Only used if states >= 3.
	improvedId?: ActionId;

	// Only used if states >= 4.
	improvedId2?: ActionId;
};

// Icon-based UI for picking buffs / consumes / etc
// ModObject is the object being modified (Sim, Player, or Target).
export class IconPicker<ModObject, ValueType> extends Input<ModObject, ValueType> {
	private readonly config: IconPickerConfig<ModObject, ValueType>;

	private readonly rootAnchor: HTMLAnchorElement;
	private readonly improvedAnchor: HTMLAnchorElement;
	private readonly improvedAnchor2: HTMLAnchorElement;
	private readonly counterElem: HTMLElement;

	private currentValue: number;

	constructor(parent: HTMLElement, modObj: ModObject, config: IconPickerConfig<ModObject, ValueType>) {
		config.rootElem = document.createElement('a');
		super(parent, 'icon-input', modObj, config);
		this.config = config;
		this.currentValue = 0;

		this.rootAnchor = this.rootElem as HTMLAnchorElement;
		this.rootAnchor.target = '_blank';

		const useImprovedIcons = Boolean(this.config.improvedId);
		if (useImprovedIcons) {
			this.rootAnchor.classList.add('use-improved-icons');
		}
		if (this.config.improvedId2) {
			this.rootAnchor.classList.add('use-improved-icons2');
		}
		if (!useImprovedIcons && this.config.states > 2) {
			this.rootAnchor.classList.add('use-counter');
		}

		this.rootAnchor.innerHTML = `
    <div class="icon-input-level-container">
      <a class="icon-input-improved icon-input-improved1"></a>
      <a class="icon-input-improved icon-input-improved2"></a>
      <span class="icon-input-counter"></span>
    </div>
    `;

		this.improvedAnchor = this.rootAnchor.getElementsByClassName('icon-input-improved1')[0] as HTMLAnchorElement;
		this.improvedAnchor2 = this.rootAnchor.getElementsByClassName('icon-input-improved2')[0] as HTMLAnchorElement;
		this.counterElem = this.rootAnchor.getElementsByClassName('icon-input-counter')[0] as HTMLElement;

		this.config.id.fillAndSet(this.rootAnchor, true, true);

		if (this.config.states >= 3 && this.config.improvedId) {
			this.config.improvedId.fillAndSet(this.improvedAnchor, true, true);
		}
		if (this.config.states >= 4 && this.config.improvedId2) {
			this.config.improvedId2.fillAndSet(this.improvedAnchor2, true, true);
		}

		this.init();

		this.rootAnchor.addEventListener('click', event => {
			event.preventDefault();
		});
		this.rootAnchor.addEventListener('contextmenu', event => {
			event.preventDefault();
		});
		this.rootAnchor.addEventListener('mousedown', event => {
			const rightClick = isRightClick(event);

			if (rightClick) {
				if (this.currentValue > 0) {
					this.currentValue--;
					this.inputChanged(TypedEvent.nextEventID());
				}
			} else {
				if (this.config.states == 0 || (this.currentValue + 1) < this.config.states) {
					this.currentValue++;
					this.inputChanged(TypedEvent.nextEventID());
				}
			}
		});

		this.rootAnchor.addEventListener('touchstart', event => {
			if (this.config.states == 0 || (this.currentValue + 1) < this.config.states) {
				this.currentValue++;
				this.inputChanged(TypedEvent.nextEventID());
			} else if (this.currentValue > 0) { // roll over
				this.currentValue = 0;
				this.inputChanged(TypedEvent.nextEventID());
			}
			event.preventDefault();
		});
		this.rootAnchor.addEventListener('touchend', event => {
			event.preventDefault();
		});
	}

	getInputElem(): HTMLElement {
		return this.rootAnchor;
	}

	getInputValue(): ValueType {
		if (this.config.states == 2) {
			return Boolean(this.currentValue) as unknown as ValueType;
		} else {
			return this.currentValue as unknown as ValueType;
		}
	}

	setInputValue(newValue: ValueType) {
		this.currentValue = Number(newValue);

		if (this.currentValue > 0) {
			this.rootAnchor.classList.add('active');
			this.counterElem.classList.add('active');
		} else {
			this.rootAnchor.classList.remove('active');
			this.counterElem.classList.remove('active');
		}

		if (this.config.states >= 3 && this.config.improvedId) {
			if (this.currentValue > 1) {
				this.improvedAnchor.classList.add('active');
			} else {
				this.improvedAnchor.classList.remove('active');
			}
		}
		if (this.config.states >= 4 && this.config.improvedId2) {
			if (this.currentValue > 2) {
				this.improvedAnchor2.classList.add('active');
			} else {
				this.improvedAnchor2.classList.remove('active');
			}
		}

		if (!this.config.improvedId && (this.config.states > 3 || this.config.states == 0)) {
			this.counterElem.textContent = String(this.currentValue);
		}
	}
}
