import { Spec } from '/tbc/core/proto/common.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

import { Component } from '/tbc/core/components/component.js';

declare var tippy: any;

export type SavedDataManagerConfig<ModObject, T> = {
	label: string;
	storageKey: string;
	changeEmitters: Array<TypedEvent<any>>,
	equals: (a: T, b: T) => boolean;
	getData: (modObject: ModObject) => T;
	setData: (eventID: EventID, modObject: ModObject, data: T) => void;
	toJson: (a: T) => any;
	fromJson: (obj: any) => T;
};

export type SavedDataConfig<ModObject, T> = {
	name: string;
	data: T;
	tooltip?: string;
	isPreset?: boolean;

	// If set, will automatically hide the saved data when this evaluates to false.
	enableWhen?: (obj: ModObject) => boolean;
};

type SavedData<ModObject, T> = {
	name: string;
	data: T;
	elem: HTMLElement;
	enableWhen?: (obj: ModObject) => boolean;
};

export class SavedDataManager<ModObject, T> extends Component {
	private readonly modObject: ModObject;
	private readonly config: SavedDataManagerConfig<ModObject, T>;

	private readonly userData: Array<SavedData<ModObject, T>>;
	private readonly presets: Array<SavedData<ModObject, T>>;

	private readonly savedDataDiv: HTMLElement;
	private readonly saveInput: HTMLInputElement;
	private frozen: boolean;

	constructor(parent: HTMLElement, modObject: ModObject, config: SavedDataManagerConfig<ModObject, T>) {
		super(parent, 'saved-data-manager-root');
		this.modObject = modObject;
		this.config = config;

		this.userData = [];
		this.presets = [];
		this.frozen = false;

		this.rootElem.innerHTML = `
    <div class="saved-data-container">
    </div>
    <div class="saved-data-create-container">
      <input class="saved-data-save-input" type="text" placeholder="Label">
      <button class="saved-data-save-button sim-button">SAVE CURRENT ${config.label.toUpperCase()}</button>
    </div>
    `;

		this.savedDataDiv = this.rootElem.getElementsByClassName('saved-data-container')[0] as HTMLElement;

		this.saveInput = this.rootElem.getElementsByClassName('saved-data-save-input')[0] as HTMLInputElement;
		const saveButton = this.rootElem.getElementsByClassName('saved-data-save-button')[0] as HTMLButtonElement;
		saveButton.addEventListener('click', event => {
			if (this.frozen)
				return;

			const newName = this.saveInput.value;
			if (!newName) {
				alert(`Choose a label for your saved ${config.label}!`);
				return;
			}

			if (newName in this.presets) {
				alert(`${config.label} with name ${newName} already exists.`);
				return;
			}

			this.addSavedData({
				name: newName,
				data: config.getData(this.modObject),
			});
			this.saveUserData();
		});
	}

	addSavedData(config: SavedDataConfig<ModObject, T>) {
		const newData = this.makeSavedData(config);

		const dataArr = config.isPreset ? this.presets : this.userData;

		const oldIdx = dataArr.findIndex(data => data.name == config.name);
		if (oldIdx == -1) {
			if (config.isPreset || this.presets.length == 0) {
				this.savedDataDiv.appendChild(newData.elem);
			} else {
				this.savedDataDiv.insertBefore(newData.elem, this.presets[0].elem);
			}
			dataArr.push(newData);
		} else {
			this.savedDataDiv.replaceChild(newData.elem, dataArr[oldIdx].elem);
			dataArr[oldIdx] = newData;
		}
	}

	private makeSavedData(config: SavedDataConfig<ModObject, T>): SavedData<ModObject, T> {
		const dataElem = document.createElement('div');
		dataElem.classList.add('saved-data-set-chip');
		dataElem.innerHTML = `
    <span class="saved-data-set-name">${config.name}</span>
    <span class="saved-data-set-tooltip fa fa-info-circle"></span>
    <span class="saved-data-set-delete fa fa-times"></span>
    `;

		dataElem.addEventListener('click', event => {
			this.config.setData(TypedEvent.nextEventID(), this.modObject, config.data);
			this.saveInput.value = config.name;
		});

		if (config.isPreset) {
			dataElem.classList.add('saved-data-preset');
		} else {
			const deleteButton = dataElem.getElementsByClassName('saved-data-set-delete')[0] as HTMLElement;
			deleteButton.addEventListener('click', event => {
				event.stopPropagation();
				const shouldDelete = confirm(`Delete saved ${this.config.label} '${config.name}'?`);
				if (!shouldDelete)
					return;

				const idx = this.userData.findIndex(data => data.name == config.name);
				this.userData[idx].elem.remove();
				this.userData.splice(idx, 1);
				this.saveUserData();
			});
		}

		if (config.tooltip) {
			dataElem.classList.add('saved-data-has-tooltip');
			tippy(dataElem.getElementsByClassName('saved-data-set-tooltip')[0], {
				'content': config.tooltip,
				'allowHTML': true,
			});
		}

		const checkActive = () => {
			if (this.config.equals(config.data, this.config.getData(this.modObject))) {
				dataElem.classList.add('active');
			} else {
				dataElem.classList.remove('active');
			}

			if (config.enableWhen && !config.enableWhen(this.modObject)) {
				dataElem.classList.add('disabled');
			} else {
				dataElem.classList.remove('disabled');
			}
		};

		checkActive();
		this.config.changeEmitters.forEach(emitter => emitter.on(checkActive));

		return {
			name: config.name,
			data: config.data,
			elem: dataElem,
			enableWhen: config.enableWhen,
		};
	}

	// Save data to window.localStorage.
	private saveUserData() {
		const userData: Record<string, Object> = {};
		this.userData.forEach(savedData => {
			userData[savedData.name] = this.config.toJson(savedData.data);
		});

		window.localStorage.setItem(this.config.storageKey, JSON.stringify(userData));
	}

	// Load data from window.localStorage.
	loadUserData() {
		const dataStr = window.localStorage.getItem(this.config.storageKey);
		if (!dataStr)
			return;

		let jsonData;
		try {
			jsonData = JSON.parse(dataStr);
		} catch (e) {
			console.warn('Invalid json for local storage value: ' + dataStr);
		}

		for (let name in jsonData) {
			try {
				this.addSavedData({
					name: name,
					data: this.config.fromJson(jsonData[name]),
				});
			} catch (e) {
				console.warn('Failed parsing saved data: ' + jsonData[name]);
			}
		}
	}

	// Prevent user input from creating / deleting saved data.
	freeze() {
		this.frozen = true;
		this.rootElem.classList.add('frozen');
	}
}
